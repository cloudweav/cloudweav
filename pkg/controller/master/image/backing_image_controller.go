package image

import (
	"fmt"
	"reflect"

	lhv1beta2 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	"github.com/rancher/norman/condition"
	"k8s.io/apimachinery/pkg/api/errors"

	cloudweavv1beta1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlcloudweavv1beta1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	ctllhv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/longhorn.io/v1beta2"
	"github.com/cloudweav/cloudweav/pkg/ref"
	"github.com/cloudweav/cloudweav/pkg/util"
)

// backingImageHandler syncs upload progress from backing image to vm image status
type backingImageHandler struct {
	vmImages          ctlcloudweavv1beta1.VirtualMachineImageClient
	vmImageCache      ctlcloudweavv1beta1.VirtualMachineImageCache
	backingImages     ctllhv1.BackingImageClient
	backingImageCache ctllhv1.BackingImageCache
}

func (h *backingImageHandler) OnChanged(_ string, backingImage *lhv1beta2.BackingImage) (*lhv1beta2.BackingImage, error) {
	if backingImage == nil || backingImage.DeletionTimestamp != nil {
		return nil, nil
	}
	if backingImage.Annotations[util.AnnotationImageID] == "" || len(backingImage.Status.DiskFileStatusMap) != 1 {
		return nil, nil
	}
	namespace, name := ref.Parse(backingImage.Annotations[util.AnnotationImageID])
	vmImage, err := h.vmImageCache.Get(namespace, name)
	if errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// There are two states that we care about here:
	// - ImageInitialized
	// - ImageImported
	// If ImageInitialized isn't yet true, it means there's no backing
	// image or storage class, so we've got nothing to work with yet and
	// should return immediately.
	if !cloudweavv1beta1.ImageInitialized.IsTrue(vmImage) {
		return nil, nil
	}
	// If ImageImported is not unknown, it means the backing image has
	// been imported, and we think we know everything about it, i.e. we've
	// now been through a series of progress updates during image download,
	// and those are finally done, so let's not worry about further updates.
	// The problem with this logic is that when we add new fields (e.g.
	// VirtualSize), existing images won't pick up those newly added fields
	// if we return here immediately.  So, now there's an additional check
	// for that new field.  Another, simpler, alternative would be to just
	// drop the ImageImported.IsUnknown check entirely, and let the following
	// loop run through on every OnChanged event.
	if !cloudweavv1beta1.ImageImported.IsUnknown(vmImage) && vmImage.Status.VirtualSize == backingImage.Status.VirtualSize {
		return nil, nil
	}
	toUpdate := vmImage.DeepCopy()
	for _, status := range backingImage.Status.DiskFileStatusMap {
		if status.State == lhv1beta2.BackingImageStateFailed {
			toUpdate = handleFail(toUpdate, condition.Cond(cloudweavv1beta1.ImageImported), fmt.Errorf(status.Message))
			toUpdate.Status.Progress = status.Progress
		} else if status.State == lhv1beta2.BackingImageStateReady {
			// We can't set ImageImported to True until we know the VirtualSize,
			// which will happen only after stats.State == lhv1beta2.BackingImageStateReady
			// (it's not there yet for lhv1beta2.BackingImageStateReadyForTransfer)
			cloudweavv1beta1.ImageImported.True(toUpdate)
			cloudweavv1beta1.ImageImported.Reason(toUpdate, "Imported")
			cloudweavv1beta1.ImageImported.Message(toUpdate, status.Message)
			// Clear the ImageRetryLimitExceeded reason and message to prevent the error message
			// from lingering in the Cloudweav dashboard after multiple image import retries
			// have failed but eventually succeeded.
			cloudweavv1beta1.ImageRetryLimitExceeded.False(toUpdate)
			cloudweavv1beta1.ImageRetryLimitExceeded.Reason(toUpdate, "")
			cloudweavv1beta1.ImageRetryLimitExceeded.Message(toUpdate, "")
			toUpdate.Status.Progress = status.Progress
			toUpdate.Status.Size = backingImage.Status.Size
			toUpdate.Status.VirtualSize = backingImage.Status.VirtualSize
		} else if status.Progress != toUpdate.Status.Progress {
			cloudweavv1beta1.ImageImported.Unknown(toUpdate)
			cloudweavv1beta1.ImageImported.Reason(toUpdate, "Importing")
			cloudweavv1beta1.ImageImported.Message(toUpdate, status.Message)
			// backing image file upload progress can be 100 before it is ready
			// Set VM image progress to be 99 for better UX in this case
			if status.Progress == 100 {
				toUpdate.Status.Progress = 99
			} else {
				toUpdate.Status.Progress = status.Progress
			}
		}
	}

	if !reflect.DeepEqual(vmImage, toUpdate) {
		if _, err := h.vmImages.Update(toUpdate); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
