package virtualmachineimage

import (
	"encoding/json"
	"errors"
	"fmt"

	longhorntypes "github.com/longhorn/longhorn-manager/types"
	ctlstoragev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/storage/v1"
	"github.com/rancher/wrangler/v3/pkg/slice"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/util"
	"github.com/cloudweav/cloudweav/pkg/webhook/types"
)

func NewMutator(storageClassCache ctlstoragev1.StorageClassCache) types.Mutator {
	return &virtualMachineImageMutator{
		storageClassCache: storageClassCache,
	}
}

type virtualMachineImageMutator struct {
	types.DefaultMutator
	storageClassCache ctlstoragev1.StorageClassCache
}

func (m *virtualMachineImageMutator) Resource() types.Resource {
	return types.Resource{
		Names:      []string{cloudweavv1.VirtualMachineImageResourceName},
		Scope:      admissionregv1.NamespacedScope,
		APIGroup:   cloudweavv1.SchemeGroupVersion.Group,
		APIVersion: cloudweavv1.SchemeGroupVersion.Version,
		ObjectType: &cloudweavv1.VirtualMachineImage{},
		OperationTypes: []admissionregv1.OperationType{
			admissionregv1.Create,
		},
	}
}

func (m *virtualMachineImageMutator) Create(_ *types.Request, newObj runtime.Object) (types.PatchOps, error) {
	newImage := newObj.(*cloudweavv1.VirtualMachineImage)

	return m.patchImageStorageClassParams(newImage)
}

func (m *virtualMachineImageMutator) patchImageStorageClassParams(newImage *cloudweavv1.VirtualMachineImage) ([]string, error) {
	var patchOps types.PatchOps

	storageClassName := newImage.Annotations[util.AnnotationStorageClassName]
	storageClass, err := m.getStorageClass(storageClassName)
	if err != nil {
		return patchOps, err
	}

	parameters := mergeStorageClassParams(newImage, storageClass)
	valueBytes, err := json.Marshal(parameters)
	if err != nil {
		return patchOps, err
	}

	verb := "add"
	if newImage.Spec.StorageClassParameters != nil {
		verb = "replace"
	}

	patchOps = append(patchOps, fmt.Sprintf(`{"op": "%s", "path": "/spec/storageClassParameters", "value": %s}`, verb, string(valueBytes)))
	return patchOps, nil
}

func (m *virtualMachineImageMutator) getStorageClass(storageClassName string) (*storagev1.StorageClass, error) {
	if storageClassName != "" {
		storageClass, err := m.storageClassCache.Get(storageClassName)
		if err != nil {
			return nil, err
		}
		if storageClass.Provisioner != longhorntypes.LonghornDriverName {
			return nil, fmt.Errorf("the provisioner of storageClass must be %s, not %s", longhorntypes.LonghornDriverName, storageClass.Provisioner)
		}
		if storageClass.Parameters[util.LonghornOptionBackingImageName] != "" {
			return nil, errors.New("can not use a backing image storageClass as the base storageClass template")
		}
		return storageClass, nil
	}

	storageClasses, err := m.storageClassCache.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	for _, storageClass := range storageClasses {
		if storageClass.Annotations[util.AnnotationIsDefaultStorageClassName] == "true" &&
			storageClass.Provisioner == longhorntypes.LonghornDriverName {
			return storageClass, nil
		}
	}

	return nil, nil
}

func mergeStorageClassParams(image *cloudweavv1.VirtualMachineImage, storageClass *storagev1.StorageClass) map[string]string {
	params := util.GetImageDefaultStorageClassParameters()
	var mergeParams map[string]string
	if storageClass != nil {
		mergeParams = storageClass.Parameters
	} else if image.Spec.StorageClassParameters != nil {
		mergeParams = image.Spec.StorageClassParameters
	}
	var allowPatchParams = []string{
		longhorntypes.OptionNodeSelector, longhorntypes.OptionDiskSelector,
		longhorntypes.OptionNumberOfReplicas, longhorntypes.OptionStaleReplicaTimeout,
		util.LonghornDataLocality,
		util.LonghornOptionEncrypted,
		util.CSIProvisionerSecretNameKey, util.CSIProvisionerSecretNamespaceKey,
		util.CSINodeStageSecretNameKey, util.CSINodeStageSecretNamespaceKey,
		util.CSINodePublishSecretNameKey, util.CSINodePublishSecretNamespaceKey,
	}

	for k, v := range mergeParams {
		if slice.ContainsString(allowPatchParams, k) {
			params[k] = v
		}
	}
	return params
}
