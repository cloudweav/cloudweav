package addon

import (
	"encoding/base64"
	"fmt"
	"time"

	helmv1 "github.com/k3s-io/helm-controller/pkg/apis/helm.cattle.io/v1"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/util"
)

const (
	// addon enqueue self interval, defaults to 5s
	enqueueInterval = 5
)

// get the current addon related helmchart
// bool: if addonOwned or not
func (h *Handler) getAddonHelmChart(aObj *cloudweavv1.Addon) (*helmv1.HelmChart, bool, error) {
	hc, err := h.helm.Get(aObj.Namespace, aObj.Name, metav1.GetOptions{})
	if err != nil {
		// chart is gone
		if apierrors.IsNotFound(err) {
			logrus.Debugf("helmChart not found to addon %v", aObj.Name)
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("error querying helmchart %v", err)
	}

	addonOwned := false
	for _, v := range hc.GetOwnerReferences() {
		if v.Kind == aObj.Kind && v.APIVersion == aObj.APIVersion && v.UID == aObj.UID && v.Name == aObj.Name {
			addonOwned = true
			break
		}
	}
	return hc, addonOwned, nil
}

// check if update is needed, when needed, also return values related string for further use
func (h *Handler) isHelmchartUpdateNeeded(aObj *cloudweavv1.Addon, hc *helmv1.HelmChart) (bool, string, error) {
	vals, err := defaultValues(aObj)
	if err != nil {
		return false, "", fmt.Errorf("error generating default values of addon %s/%s: %v", aObj.Namespace, aObj.Name, err)
	}

	return (hc.Spec.ValuesContent != vals || hc.Spec.Version != aObj.Spec.Version || hc.Spec.Chart != aObj.Spec.Chart || hc.Spec.Repo != aObj.Spec.Repo), vals, nil
}

// deploy a new chart
func (h *Handler) deployHelmChart(aObj *cloudweavv1.Addon) error {
	vals, err := defaultValues(aObj)
	if err != nil {
		return err
	}

	hc := &helmv1.HelmChart{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aObj.Name,
			Namespace: aObj.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: aObj.APIVersion,
					Kind:       aObj.Kind,
					Name:       aObj.Name,
					UID:        aObj.UID,
				},
			},
		},
		Spec: helmv1.HelmChartSpec{
			Chart:         aObj.Spec.Chart,
			Repo:          aObj.Spec.Repo,
			ValuesContent: vals,
			Version:       aObj.Spec.Version,
			BackOffLimit:  &cloudweavv1.DefaultJobBackOffLimit,
		},
	}
	_, err = h.helm.Create(hc)
	if err != nil {
		return fmt.Errorf("error creating helmchart object %v", err)
	}

	return nil
}

func defaultValues(a *cloudweavv1.Addon) (string, error) {
	if a.Spec.ValuesContent != "" {
		return a.Spec.ValuesContent, nil
	}

	valsEncoded, ok := a.Annotations[util.AddonValuesAnnotation]
	if ok {
		valByte, err := base64.StdEncoding.DecodeString(valsEncoded)
		if err != nil {
			return "", fmt.Errorf("error decoding addon defaults: %v", err)
		}

		return string(valByte), nil
	}
	// no overrides. Use packaged chart defaults
	return "", nil
}

func isJobComplete(j *batchv1.Job) bool {
	if j.Status.CompletionTime != nil {
		return true
	}

	for _, v := range j.Status.Conditions {
		if v.Type == batchv1.JobFailed && v.Reason == "BackoffLimitExceeded" {
			return true
		}
	}

	return false
}

func (h *Handler) enqueueAfter(aObj *cloudweavv1.Addon) (*cloudweavv1.Addon, error) {
	h.addon.EnqueueAfter(aObj.Namespace, aObj.Name, enqueueInterval*time.Second)
	return aObj, nil
}

func (h *Handler) getJob(hc *helmv1.HelmChart) (*batchv1.Job, error) {
	if hc.Status.JobName == "" {
		return nil, fmt.Errorf("waiting for job to be populated on helmchart %s", hc.Name)
	}
	return h.job.Cache().Get(hc.Namespace, hc.Status.JobName)
}

func (h *Handler) currentDeletionJob(hc *helmv1.HelmChart) (*batchv1.Job, bool, error) {
	j, err := h.getJob(hc)
	if err != nil {
		return nil, false, err
	}

	// job creation timestamp should be after deletion timestamp of
	// helm chart to ensure that we are checking the correct job
	if j.CreationTimestamp.After(hc.DeletionTimestamp.Time) || j.CreationTimestamp.Equal(hc.DeletionTimestamp) {
		return j, true, nil
	}

	return j, false, nil
}

func (h *Handler) currentInstallationJob(hc *helmv1.HelmChart, a *cloudweavv1.Addon) (*batchv1.Job, bool, error) {
	logrus.Debugf("querying current installation job for addon %s", a.Name)

	j, err := h.getJob(hc)
	if err != nil {
		return nil, false, err
	}

	lastUpdatedTime, err := time.Parse(time.RFC3339, cloudweavv1.AddonOperationInProgress.GetLastUpdated(a))
	if err != nil {
		return nil, false, fmt.Errorf("error parsing last updated time for AddonOperationInProgress: %v", err)
	}

	metav1LastUpdatedTime := metav1.NewTime(lastUpdatedTime)
	logrus.Debugf("last updated time on the addon: %s", lastUpdatedTime)
	// job creation timestamp should be after the last updated time stamp
	// on the inprogress condition
	if j.CreationTimestamp.After(lastUpdatedTime) || j.CreationTimestamp.Equal(&metav1LastUpdatedTime) {
		return j, true, nil
	}

	return j, false, nil
}

func markErrorCondition(aObj *cloudweavv1.Addon, msg error) {
	now := time.Now().UTC().Format(time.RFC3339)
	cloudweavv1.AddonOperationFailed.SetError(aObj, "", msg)
	cloudweavv1.AddonOperationFailed.True(aObj)
	cloudweavv1.AddonOperationFailed.LastUpdated(aObj, now)
	cloudweavv1.AddonOperationInProgress.False(aObj)
	cloudweavv1.AddonOperationCompleted.False(aObj)

}

func markInProgressCondition(aObj *cloudweavv1.Addon) {
	now := time.Now().UTC().Format(time.RFC3339)
	cloudweavv1.AddonOperationCompleted.False(aObj)
	cloudweavv1.AddonOperationInProgress.LastUpdated(aObj, now)
	cloudweavv1.AddonOperationInProgress.True(aObj)
	cloudweavv1.AddonOperationFailed.False(aObj)
	cloudweavv1.AddonOperationFailed.Reason(aObj, "")
	cloudweavv1.AddonOperationFailed.Message(aObj, "")
}

func markCompletedCondition(aObj *cloudweavv1.Addon) {
	now := time.Now().UTC().Format(time.RFC3339)
	cloudweavv1.AddonOperationCompleted.True(aObj)
	cloudweavv1.AddonOperationCompleted.LastUpdated(aObj, now)
	cloudweavv1.AddonOperationInProgress.False(aObj)
	cloudweavv1.AddonOperationFailed.False(aObj)
	cloudweavv1.AddonOperationFailed.Reason(aObj, "")
	cloudweavv1.AddonOperationFailed.Message(aObj, "")
}
