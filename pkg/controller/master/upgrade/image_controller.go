package upgrade

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
)

// vmImageHandler syncs upgrade repo image creation
type vmImageHandler struct {
	namespace     string
	upgradeClient ctlcloudweavv1.UpgradeClient
	upgradeCache  ctlcloudweavv1.UpgradeCache
}

func (h *vmImageHandler) OnChanged(_ string, image *cloudweavv1.VirtualMachineImage) (*cloudweavv1.VirtualMachineImage, error) {
	if image == nil || image.DeletionTimestamp != nil || image.Labels == nil || image.Namespace != upgradeNamespace || image.Labels[cloudweavUpgradeLabel] == "" {
		return image, nil
	}

	upgradeControllerLock.Lock()
	defer upgradeControllerLock.Unlock()

	upgrade, err := h.upgradeCache.Get(upgradeNamespace, image.Labels[cloudweavUpgradeLabel])
	if err != nil {
		if apierrors.IsNotFound(err) {
			return image, nil
		}
		return nil, err
	}

	toUpdate := upgrade.DeepCopy()

	switch {
	case cloudweavv1.ImageImported.IsTrue(image):
		setImageReadyCondition(toUpdate, corev1.ConditionTrue, "", "")
	case cloudweavv1.ImageImported.IsFalse(image):
		setImageReadyCondition(toUpdate, corev1.ConditionFalse, cloudweavv1.ImageImported.GetReason(image), cloudweavv1.ImageImported.GetMessage(image))
	case cloudweavv1.ImageRetryLimitExceeded.IsTrue(image):
		setImageReadyCondition(toUpdate, corev1.ConditionFalse, cloudweavv1.ImageRetryLimitExceeded.GetReason(image), cloudweavv1.ImageRetryLimitExceeded.GetMessage(image))
	case isUponRetryFailure(image, upgrade):
		setImageReadyCondition(toUpdate, corev1.ConditionUnknown, cloudweavv1.ImageRetryLimitExceeded.GetReason(image), cloudweavv1.ImageRetryLimitExceeded.GetMessage(image))
	default:
		return image, nil
	}

	if !reflect.DeepEqual(toUpdate, upgrade) {
		_, err := h.upgradeClient.Update(toUpdate)
		return image, err
	}

	return image, nil
}

func isUponRetryFailure(image *cloudweavv1.VirtualMachineImage, upgrade *cloudweavv1.Upgrade) bool {
	return cloudweavv1.ImageRetryLimitExceeded.IsFalse(image) &&
		(cloudweavv1.ImageRetryLimitExceeded.GetReason(image) != cloudweavv1.ImageReady.GetReason(upgrade) ||
			cloudweavv1.ImageRetryLimitExceeded.GetMessage(image) != cloudweavv1.ImageReady.GetMessage(upgrade))
}
