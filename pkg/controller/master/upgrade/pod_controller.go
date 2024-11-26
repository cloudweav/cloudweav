package upgrade

import (
	v1 "k8s.io/api/core/v1"

	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	upgradev1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/upgrade.cattle.io/v1"
)

// podHandler syncs upgrade CRD status on upgrade pod status changes
type podHandler struct {
	namespace     string
	planCache     upgradev1.PlanCache
	upgradeClient ctlcloudweavv1.UpgradeClient
	upgradeCache  ctlcloudweavv1.UpgradeCache
}

func (h *podHandler) OnChanged(_ string, pod *v1.Pod) (*v1.Pod, error) {
	if pod == nil || pod.DeletionTimestamp != nil || pod.Labels == nil || pod.Namespace != upgradeNamespace || pod.Labels[cloudweavUpgradeLabel] == "" {
		return pod, nil
	}

	upgradeControllerLock.Lock()
	defer upgradeControllerLock.Unlock()

	upgrade, err := h.upgradeCache.Get(upgradeNamespace, pod.Labels[cloudweavUpgradeLabel])
	if err != nil {
		return nil, err
	}

	component := pod.Labels[cloudweavUpgradeComponentLabel]
	switch upgrade.Labels[upgradeStateLabel] {
	case StatePreparingRepo:
		if component == upgradeComponentRepo && len(pod.Status.ContainerStatuses) > 0 {
			if pod.Status.ContainerStatuses[0].Ready {
				toUpdate := upgrade.DeepCopy()
				toUpdate.Labels[upgradeStateLabel] = StateRepoPrepared
				setRepoProvisionedCondition(toUpdate, v1.ConditionTrue, "", "")
				_, err = h.upgradeClient.Update(toUpdate)
				return pod, err
			}
		}
	}

	return pod, nil
}

func getPodWaitingStatus(pod *v1.Pod) (reason string, message string) {
	var containerStatuses []v1.ContainerStatus
	containerStatuses = append(containerStatuses, pod.Status.InitContainerStatuses...)
	containerStatuses = append(containerStatuses, pod.Status.ContainerStatuses...)

	for _, status := range containerStatuses {
		if status.State.Waiting != nil && len(status.State.Waiting.Reason) > 0 && status.State.Waiting.Reason != "PodInitializing" {
			reason = status.State.Waiting.Reason
			message = status.State.Waiting.Message
			return
		}
	}
	return
}
