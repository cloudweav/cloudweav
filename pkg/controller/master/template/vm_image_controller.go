package template

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/indexeres"
)

// vmImageHandler watch vm image and enqueue related vm template versions.
type vmImageHandler struct {
	templateVersionCache      ctlcloudweavv1.VirtualMachineTemplateVersionCache
	templateVersionController ctlcloudweavv1.VirtualMachineTemplateVersionController
}

func (h *vmImageHandler) OnChanged(_ string, vmImage *cloudweavv1.VirtualMachineImage) (*cloudweavv1.VirtualMachineImage, error) {
	if vmImage == nil || vmImage.DeletionTimestamp != nil {
		return nil, nil
	}

	vmTemplateVersions, err := h.templateVersionCache.GetByIndex(indexeres.VMTemplateVersionByImageIDIndex, fmt.Sprintf("%s/%s", vmImage.Namespace, vmImage.Name))
	if err != nil {
		if apierrors.IsNotFound(err) {
			return vmImage, nil
		}
		return vmImage, err
	}

	for _, vmTemplateVersion := range vmTemplateVersions {
		if cloudweavv1.TemplateVersionReady.IsTrue(vmTemplateVersion) {
			continue
		}
		h.templateVersionController.Enqueue(vmTemplateVersion.Namespace, vmTemplateVersion.Name)
	}
	return vmImage, nil
}
