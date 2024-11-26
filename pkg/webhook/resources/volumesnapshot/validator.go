package volumesnapshot

import (
	"fmt"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	ctlkubevirtv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/kubevirt.io/v1"
	ctllonghornv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/longhorn.io/v1beta2"
	"github.com/cloudweav/cloudweav/pkg/ref"
	"github.com/cloudweav/cloudweav/pkg/util"
	indexeresutil "github.com/cloudweav/cloudweav/pkg/util/indexeres"
	werror "github.com/cloudweav/cloudweav/pkg/webhook/error"
	"github.com/cloudweav/cloudweav/pkg/webhook/types"
	webhookutil "github.com/cloudweav/cloudweav/pkg/webhook/util"
)

const (
	fieldSourceName = "spec.source.name"
	fieldTypeName   = "spec.type"
)

func NewValidator(
	pvcCache ctlcorev1.PersistentVolumeClaimCache,
	engineCache ctllonghornv1.EngineCache,
	resourceQuotaCache ctlcloudweavv1.ResourceQuotaCache,
	vmCache ctlkubevirtv1.VirtualMachineCache,
) types.Validator {
	return &volumeSnapshotValidator{
		pvcCache:           pvcCache,
		engineCache:        engineCache,
		resourceQuotaCache: resourceQuotaCache,
		vmCache:            vmCache,
	}
}

type volumeSnapshotValidator struct {
	types.DefaultValidator

	pvcCache           ctlcorev1.PersistentVolumeClaimCache
	engineCache        ctllonghornv1.EngineCache
	resourceQuotaCache ctlcloudweavv1.ResourceQuotaCache
	vmCache            ctlkubevirtv1.VirtualMachineCache
}

func (v *volumeSnapshotValidator) Resource() types.Resource {
	return types.Resource{
		Names:      []string{"volumesnapshots"},
		Scope:      admissionregv1.NamespacedScope,
		APIGroup:   snapshotv1.SchemeGroupVersion.Group,
		APIVersion: snapshotv1.SchemeGroupVersion.Version,
		ObjectType: &snapshotv1.VolumeSnapshot{},
		OperationTypes: []admissionregv1.OperationType{
			admissionregv1.Create,
		},
	}
}

func (v *volumeSnapshotValidator) Create(_ *types.Request, newObj runtime.Object) error {
	newVolumeSnapshot := newObj.(*snapshotv1.VolumeSnapshot)

	if newVolumeSnapshot.Spec.Source.PersistentVolumeClaimName == nil {
		return nil
	}

	for _, owner := range newVolumeSnapshot.OwnerReferences {
		// resource quota is already checked in the VMBackup webhook, skip it here
		if owner.Kind == "VirtualMachineBackup" {
			continue
		}
	}

	resourceQuota, err := v.resourceQuotaCache.Get(newVolumeSnapshot.Namespace, util.DefaultResourceQuotaName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return werror.NewInternalError(fmt.Sprintf("failed to get resource quota %s/%s, err: %s", newVolumeSnapshot.Namespace, util.DefaultResourceQuotaName, err))
	}

	vms, err := v.vmCache.GetByIndex(indexeresutil.VMByPVCIndex, ref.Construct(newVolumeSnapshot.Namespace, *newVolumeSnapshot.Spec.Source.PersistentVolumeClaimName))
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return werror.NewInternalError(fmt.Sprintf("failed to get VM by PVC %s/%s, err: %s", newVolumeSnapshot.Namespace, *newVolumeSnapshot.Spec.Source.PersistentVolumeClaimName, err))
		}
	} else if len(vms) > 0 {
		vm := vms[0]
		if err = webhookutil.CheckTotalSnapshotSizeOnVM(v.pvcCache, v.engineCache, vm, resourceQuota.Spec.SnapshotLimit.VMTotalSnapshotSizeQuota[vm.Name]); err != nil {
			return err
		}
	}

	if err = webhookutil.CheckTotalSnapshotSizeOnNamespace(v.pvcCache, v.engineCache, newVolumeSnapshot.Namespace, resourceQuota.Spec.SnapshotLimit.NamespaceTotalSnapshotSizeQuota); err != nil {
		return err
	}
	return nil
}
