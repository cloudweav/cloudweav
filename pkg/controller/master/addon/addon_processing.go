package addon

import (
	"fmt"
	"strings"

	helmv1 "github.com/k3s-io/helm-controller/pkg/apis/helm.cattle.io/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
)

// Addons can in be the possible states:
// - AddonEnabling: Used to indicate that Addon is being enabled and Helm chart deploy is triggered
// - AddonEnabled: Used to ensure Addon is enabled and HelmChart is running
// - AddonUpgrading: Used to indicate Addon spec has changed and an upgrade is being triggered
// - AddonUpgraded: Used to indicate that addon upgrade was successful
// - AddonDisabling: Used to indicate addon is being disabled
// - AddonDisabled: Used to indicate that addon has been disabled and associated helmChart has been cleaned up

// In addition each state will have 3 conditions associated to indicate what state processing the object is in
// - InProgress: An operation is in progress
// - Completed: Last requested operation has completed
// - Error: Last request operation has errored along with the error message

func (h *Handler) processEnabledAddons(a *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {
	aObj := a.DeepCopy()
	switch aObj.Status.Status {
	case cloudweavv1.AddonDisabled, cloudweavv1.AddonInitState:
		// addon enabled for the first time and needs deploying
		return h.processEnableAddon(aObj)
	case cloudweavv1.AddonEnabling, cloudweavv1.AddonUpdating:
		//addon enabling or updating, need to wait and watch status of job
		return h.processEnablingAddon(aObj)
	case cloudweavv1.AddonDeployed:
		// addon is enabled, wait and watch for upgrades
		return h.processAddonUpgrades(aObj)
	}
	return aObj, nil
}

// processDisabledAddon will temporarily set the addon status to AddonDisabling
// it will then trigger deletion of helm chart and accordingly setup conditions
// to ensure that job status is tracked to completed
// when delete job is complete, the addon status is updated to AddonDisabled
// by default a newly created Addon will be placed into AddonDisabled state
func (h *Handler) processDisabledAddons(a *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {
	aObj := a.DeepCopy()

	if aObj.Status.Status == cloudweavv1.AddonDisabled {
		return aObj, nil
	}

	// if status is not AddonDisabling and an operation is in progress
	// then wait for last operation to be completed
	if aObj.Status.Status != cloudweavv1.AddonDisabling && cloudweavv1.AddonOperationInProgress.IsTrue(a) {
		return h.enqueueAfter(aObj)
	}

	// if status is not AddonDisabling then set status and process object
	if aObj.Status.Status != cloudweavv1.AddonDisabling {
		aObj.Status.Status = cloudweavv1.AddonDisabling
		return h.addon.UpdateStatus(aObj)
	}

	return h.waitForAddonDisable(aObj)
}

func (h *Handler) waitForAddonDisable(a *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {

	aObj := a.DeepCopy()
	// get HelmChart and check for processing state
	hc, owned, err := h.getAddonHelmChart(aObj)
	if err != nil {
		return aObj, err
	}

	// if chart is not owned by AddonController.
	// disable Addon and place error.
	if !owned {
		// no helm chart exists.
		// no further action required and exit processing of addon
		if hc == nil {
			markCompletedCondition(aObj)
			aObj.Status.Status = cloudweavv1.AddonDisabled
			return h.addon.UpdateStatus(aObj)
		}

		// if condition already exists, ignore chart
		// as no further action is possible without manual intervention
		if cloudweavv1.AddonOperationFailed.IsTrue(aObj) {
			// keep requeing object, as when chart is eventually deleted, this condition will not be met
			// and chart will be removed
			return h.enqueueAfter(aObj)
		}

		// set error condition for addon. no further processing is now possible
		markErrorCondition(aObj, fmt.Errorf("helmChart %s exists and is now owned by Addon %s. cannot proceed further", hc.Name, aObj.Name))
		return h.addon.UpdateStatus(aObj)
	}

	// if helmChart exists trigger deletion of the same
	if hc.DeletionTimestamp.IsZero() {
		cloudweavv1.AddonOperationInProgress.True(aObj)
		err := h.helm.Delete(hc.Namespace, hc.Name, &metav1.DeleteOptions{})
		if err != nil {
			return aObj, fmt.Errorf("error deleting helm chart: %v", err)
		}
		return h.addon.UpdateStatus(aObj)
	}

	// if deletion job has not been set then requeue and retry again
	if !strings.Contains(hc.Status.JobName, helmChartJobDeletePrefix) {
		return h.enqueueAfter(aObj)
	}

	j, ok, err := h.currentDeletionJob(hc)
	if err != nil {
		return aObj, err
	}

	if ok {
		if !isJobComplete(j) {
			return h.enqueueAfter(aObj)
		}

		if j.Status.Failed > 0 {
			markErrorCondition(aObj, fmt.Errorf("addon disable job %s failed", j.Name))
		}

		if j.Status.Succeeded > 0 {
			markCompletedCondition(aObj)
			aObj.Status.Status = cloudweavv1.AddonDisabled
		}

		return h.addon.UpdateStatus(aObj)
	}

	return h.enqueueAfter(aObj)
}

// enableAddon will check for addon status and create a helmChart if needed and setup
// appropriate conditions to ensure that addon - helmchart state can reconcile
func (h *Handler) processEnableAddon(a *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {
	aObj := a.DeepCopy()

	// since re-enable of an Addon is possible from AddonDisabling state after a failure
	// we need to ensure the underlying HC has been manually cleared up
	hc, owned, err := h.getAddonHelmChart(aObj)
	if err != nil {
		return aObj, err
	}

	if owned && hc != nil && !hc.DeletionTimestamp.IsZero() {
		logrus.Warnf("helmChart %s has a deletion timestamp. likely from an older install, wait for GC before re-installing", hc.Name)
		return h.enqueueAfter(aObj)
	}

	if !owned && hc != nil {
		logrus.Warnf("cannot create helmChart %s, as another un-owned chart with same name already exists. requeuing until this is removed", hc.Name)
		return h.enqueueAfter(aObj)
	}
	// ensure condition is set first. This is to ensure the timestamp is before helm chart creation
	// which allows us to be certain that the job creation timestamp should be after this condition
	if !cloudweavv1.AddonOperationInProgress.IsTrue(aObj) {
		markInProgressCondition(aObj)
		return h.addon.UpdateStatus(aObj)
	}

	if hc == nil {
		if err := h.deployHelmChart(aObj); err != nil {
			return aObj, fmt.Errorf("error deploying helmChart: %v", err)
		}
	}

	aObj.Status.Status = cloudweavv1.AddonEnabling
	return h.addon.UpdateStatus(aObj)
}

// processEnablingAddon will check if the job for deploying the helmChart is completed
// based on status of job, the addon deployment will be marked successful or errored
// and status will be updated accordingly
func (h *Handler) processEnablingAddon(a *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {
	aObj := a.DeepCopy()
	hc, _, err := h.getAddonHelmChart(aObj)
	if err != nil {
		return aObj, err
	}

	// if addon is failed, no further processing is needed. Ignore until user
	// disables and fixes the addon
	if cloudweavv1.AddonOperationFailed.IsTrue(aObj) {
		return aObj, nil
	}

	logrus.Debugf("processing addon: %v", a)
	if hc == nil {
		// unable to find helmChart. Likely manually GC'd. Reset to initial state
		// this will allow addon to be processed again and redeployed
		markCompletedCondition(aObj)
		aObj.Status.Status = cloudweavv1.AddonDisabled
		logrus.Debugf("hc not found so disabling and enabling: %v", a)
		return h.addon.UpdateStatus(aObj)
	}

	logrus.Debugf("checking status of addon: %s", aObj.Name)
	return h.reconcileCurrentInstallationJob(hc, aObj, cloudweavv1.AddonDeployed)
}

// processAddonUpgrades will check if an update is needed for the helm chart
// and if needed will update the helmChart object spec
func (h *Handler) processAddonUpgrades(a *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {
	aObj := a.DeepCopy()

	hc, _, err := h.getAddonHelmChart(aObj)

	if err != nil {
		return aObj, err
	}

	if hc == nil {
		return aObj, h.deployHelmChart(aObj)
	}
	needed, vals, err := h.isHelmchartUpdateNeeded(aObj, hc)
	if err != nil {
		return aObj, err
	}

	if needed {
		// set InProgress condition as this is needed to ensure that the in progress time stamp is setup before
		// triggering the helm chart upgrade
		if aObj.Status.Status != cloudweavv1.AddonUpdating && !cloudweavv1.AddonOperationInProgress.IsTrue(aObj) {
			markInProgressCondition(aObj)
			return h.addon.UpdateStatus(aObj)
		}

		hc, _, err := h.getAddonHelmChart(aObj)
		if err != nil {
			return aObj, err
		}
		hc.Spec.Chart = aObj.Spec.Chart
		hc.Spec.Repo = aObj.Spec.Repo
		hc.Spec.Version = aObj.Spec.Version
		hc.Spec.ValuesContent = vals
		_, err = h.helm.Update(hc)
		if err != nil {
			return aObj, err
		}

		if aObj.Status.Status != cloudweavv1.AddonUpdating && cloudweavv1.AddonOperationInProgress.IsTrue(aObj) {
			aObj.Status.Status = cloudweavv1.AddonUpdating
			return h.addon.UpdateStatus(aObj)
		}
	}

	return aObj, nil
}

// reconcileCurrentInstallationJob attempts to query the helm chart initiated job to reconcile state for corresponding addon job
func (h *Handler) reconcileCurrentInstallationJob(hc *helmv1.HelmChart, aObj *cloudweavv1.Addon, state cloudweavv1.AddonState) (*cloudweavv1.Addon, error) {
	j, ok, err := h.currentInstallationJob(hc, aObj)
	if err != nil {
		return aObj, err
	}

	// the current job on hc is from the previous deploy/update
	// or job is not complete. then wait for this to be completed
	if !ok || !isJobComplete(j) {
		return h.enqueueAfter(aObj)
	}

	if j.Status.Failed > 0 {
		if cloudweavv1.AddonOperationFailed.IsTrue(aObj) {
			return aObj, nil // no further action possible, until the job is successful or hc is removed
		}
		markErrorCondition(aObj, fmt.Errorf("addon deployment job %s failed", j.Name))
	}

	if j.Status.Succeeded > 0 {
		markCompletedCondition(aObj)
		aObj.Status.Status = state
	}

	return h.addon.UpdateStatus(aObj)
}
