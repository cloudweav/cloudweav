package template

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/ref"
	"github.com/cloudweav/cloudweav/pkg/util"
)

const (
	TemplateLabel = "template.cloudweavhci.io/templateID"
)

// templateVersionHandler sets metadata and status to templateVersion objects,
// including labels, ownerReference and status.Version.
type templateVersionHandler struct {
	templateCache      ctlcloudweavv1.VirtualMachineTemplateCache
	templateVersions   ctlcloudweavv1.VirtualMachineTemplateVersionClient
	templateController ctlcloudweavv1.VirtualMachineTemplateController
	vmImageCache       ctlcloudweavv1.VirtualMachineImageCache
	mu                 sync.RWMutex //use mutex to avoid create duplicated version
}

func (h *templateVersionHandler) OnChanged(_ string, tv *cloudweavv1.VirtualMachineTemplateVersion) (*cloudweavv1.VirtualMachineTemplateVersion, error) {
	if tv == nil || tv.DeletionTimestamp != nil {
		return nil, nil
	}

	ns, templateName := ref.Parse(tv.Spec.TemplateID)
	template, err := h.templateCache.Get(ns, templateName)
	if err != nil {
		return nil, err
	}

	copyObj := tv.DeepCopy()

	//set labels
	if copyObj.Labels == nil {
		copyObj.Labels = make(map[string]string)
	}
	if _, ok := copyObj.Labels[TemplateLabel]; !ok {
		copyObj.Labels[TemplateLabel] = templateName
	}

	//set ownerReference
	flagTrue := true
	ownerRef := []metav1.OwnerReference{{
		Name:               template.Name,
		APIVersion:         template.APIVersion,
		UID:                template.UID,
		Kind:               template.Kind,
		BlockOwnerDeletion: &flagTrue,
		Controller:         &flagTrue,
	}}

	if len(copyObj.OwnerReferences) == 0 {
		copyObj.OwnerReferences = ownerRef
	} else if !isVersionOwnedByTemplate(copyObj, template) {
		copyObj.OwnerReferences = append(copyObj.OwnerReferences, ownerRef...)
	}

	if isReady, err := h.isVMImagesReady(tv); err != nil {
		return nil, err
	} else if isReady {
		cloudweavv1.TemplateVersionReady.True(copyObj)
	} else {
		cloudweavv1.TemplateVersionReady.False(copyObj)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	//set version
	if !cloudweavv1.VersionAssigned.IsTrue(copyObj) {
		existLatestVersion, _, err := getTemplateLatestVersion(tv.Namespace, tv.Spec.TemplateID, h.templateVersions)
		if err != nil {
			return nil, err
		}

		latestVersion := existLatestVersion + 1
		copyObj.Status.Version = latestVersion
		cloudweavv1.VersionAssigned.True(copyObj)
	}

	if !reflect.DeepEqual(copyObj, tv) {
		if _, err = h.templateVersions.Update(copyObj); err != nil {
			return copyObj, err
		}
		h.templateController.Enqueue(ns, templateName)
	}

	return copyObj, nil
}

func (h *templateVersionHandler) isVMImagesReady(tv *cloudweavv1.VirtualMachineTemplateVersion) (bool, error) {
	volumeClaimTemplatesStr, ok := tv.Spec.VM.ObjectMeta.Annotations[util.AnnotationVolumeClaimTemplates]
	if !ok || volumeClaimTemplatesStr == "" {
		return true, nil
	}

	var volumeClaimTemplates []corev1.PersistentVolumeClaim
	if err := json.Unmarshal([]byte(volumeClaimTemplatesStr), &volumeClaimTemplates); err != nil {
		return false, fmt.Errorf("can't unmarshal %s annotation, err: %w", util.AnnotationVolumeClaimTemplates, err)
	}

	for _, volumeClaimTemplate := range volumeClaimTemplates {
		imageID, ok := volumeClaimTemplate.Annotations[util.AnnotationImageID]
		if !ok || imageID == "" {
			continue
		}

		imageNs, imageName := ref.Parse(imageID)
		if image, err := h.vmImageCache.Get(imageNs, imageName); err != nil {
			return false, err
		} else if !cloudweavv1.ImageImported.IsTrue(image) {
			return false, nil
		}
	}
	return true, nil
}

func getTemplateLatestVersion(templateVersionNs, templateID string, templateVersions ctlcloudweavv1.VirtualMachineTemplateVersionClient) (int, *cloudweavv1.VirtualMachineTemplateVersion, error) {
	var latestVersion int
	list, err := templateVersions.List(templateVersionNs, metav1.ListOptions{})
	if err != nil {
		return latestVersion, nil, err
	}

	var tvs []cloudweavv1.VirtualMachineTemplateVersion
	for _, v := range list.Items {
		if v.Spec.TemplateID == templateID {
			tvs = append(tvs, v)
		}
	}

	if len(tvs) == 0 {
		return 0, nil, nil
	}

	sort.Sort(templateVersionByCreationTimestamp(tvs))
	for _, v := range tvs {
		if cloudweavv1.VersionAssigned.IsTrue(v) {
			return v.Status.Version, &v, nil
		}
	}

	return 0, nil, nil
}

// templateVersionByCreationTimestamp sorts a list of TemplateVersion by creation timestamp, using their names as a tie breaker.
type templateVersionByCreationTimestamp []cloudweavv1.VirtualMachineTemplateVersion

func (o templateVersionByCreationTimestamp) Len() int      { return len(o) }
func (o templateVersionByCreationTimestamp) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o templateVersionByCreationTimestamp) Less(i, j int) bool {
	if o[i].CreationTimestamp.Equal(&o[j].CreationTimestamp) {
		return o[i].Name < o[j].Name
	}
	return o[j].CreationTimestamp.Before(&o[i].CreationTimestamp)
}

func isVersionOwnedByTemplate(version *cloudweavv1.VirtualMachineTemplateVersion, template *cloudweavv1.VirtualMachineTemplate) bool {
	for _, v := range version.OwnerReferences {
		if v.UID == template.UID {
			return true
		}
	}
	return false
}
