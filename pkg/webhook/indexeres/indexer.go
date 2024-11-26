package indexeres

import (
	"fmt"
	"strconv"

	lhv1beta2 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	longhorntypes "github.com/longhorn/longhorn-manager/types"
	kubevirtv1 "kubevirt.io/api/core/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/controller/master/backup"
	"github.com/cloudweav/cloudweav/pkg/util"
	indexeresutil "github.com/cloudweav/cloudweav/pkg/util/indexeres"
	"github.com/cloudweav/cloudweav/pkg/webhook/clients"
)

const (
	VMBackupBySourceUIDIndex              = "cloudweavhci.io/vmbackup-by-source-uid"
	VMBackupByIsProgressing               = "cloudweavhci.io/vmbackup-by-is-progressing"
	VMBackupByStorageClassNameIndex       = "cloudweavhci.io/vmbackup-by-storage-class-name"
	VMRestoreByTargetNamespaceAndName     = "cloudweavhci.io/vmrestore-by-target-namespace-and-name"
	VMRestoreByVMBackupNamespaceAndName   = "cloudweavhci.io/vmrestore-by-vmbackup-namespace-and-name"
	VMBackupSnapshotByPVCNamespaceAndName = "cloudweavhci.io/vmbackup-snapshot-by-pvc-namespace-and-name"
	VolumeByReplicaCountIndex             = "cloudweavhci.io/volume-by-replica-count"
	ImageByExportSourcePVCIndex           = "cloudweavhci.io/image-by-export-source-pvc"
	ScheduleVMBackupBySourceVM            = "cloudweavhci.io/svmbackup-by-source-vm"
	ScheduleVMBackupByCronGranularity     = "cloudweavhci.io/svmbackup-by-cron-granularity"
	ScheduleVMBackupBySuspended           = "cloudweavhci.io/svmbackup-by-suspended"
	ImageByStorageClass                   = "cloudweavhci.io/image-by-storage-class"
	VMInstanceMigrationByVM               = "cloudweavhci.io/vmim-by-vm"
)

func RegisterIndexers(clients *clients.Clients) {
	vmBackupCache := clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache()
	vmBackupCache.AddIndexer(VMBackupBySourceUIDIndex, vmBackupBySourceUID)
	vmBackupCache.AddIndexer(VMBackupSnapshotByPVCNamespaceAndName, vmBackupSnapshotByPVCNamespaceAndName)
	vmBackupCache.AddIndexer(VMBackupByIsProgressing, vmBackupByIsProgressing)
	vmBackupCache.AddIndexer(VMBackupByStorageClassNameIndex, vmBackupByStorageClassName)

	vmRestoreCache := clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineRestore().Cache()
	vmRestoreCache.AddIndexer(VMRestoreByTargetNamespaceAndName, vmRestoreByTargetNamespaceAndName)
	vmRestoreCache.AddIndexer(VMRestoreByVMBackupNamespaceAndName, vmRestoreByVMBackupNamespaceAndName)

	podCache := clients.CoreFactory.Core().V1().Pod().Cache()
	podCache.AddIndexer(indexeresutil.PodByVMNameIndex, indexeresutil.PodByVMName)

	volumeCache := clients.LonghornFactory.Longhorn().V1beta2().Volume().Cache()
	volumeCache.AddIndexer(VolumeByReplicaCountIndex, VolumeByReplicaCount)

	vmImageInformer := clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage().Cache()
	vmImageInformer.AddIndexer(ImageByExportSourcePVCIndex, imageByExportSourcePVC)
	vmImageInformer.AddIndexer(ImageByStorageClass, imageByStorageClass)

	vmInformer := clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache()
	vmInformer.AddIndexer(indexeresutil.VMByPVCIndex, indexeresutil.VMByPVC)

	svmBackupCache := clients.CloudweavFactory.Cloudweavhci().V1beta1().ScheduleVMBackup().Cache()
	svmBackupCache.AddIndexer(ScheduleVMBackupBySourceVM, scheduleVMBackupBySourceVM)
	svmBackupCache.AddIndexer(ScheduleVMBackupByCronGranularity, scheduleVMBackupByCronGranularity)
	svmBackupCache.AddIndexer(ScheduleVMBackupBySuspended, scheduleVMBackupBySuspended)

	scInformer := clients.StorageFactory.Storage().V1().StorageClass().Cache()
	scInformer.AddIndexer(indexeresutil.StorageClassBySecretIndex, indexeresutil.StorageClassBySecret)

	vmimCache := clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstanceMigration().Cache()
	vmimCache.AddIndexer(VMInstanceMigrationByVM, vmInstanceMigrationByVM)
}

func vmBackupBySourceUID(obj *cloudweavv1.VirtualMachineBackup) ([]string, error) {
	if obj.Status != nil && obj.Status.SourceUID != nil {
		return []string{string(*obj.Status.SourceUID)}, nil
	}
	return []string{}, nil
}

func vmBackupSnapshotByPVCNamespaceAndName(obj *cloudweavv1.VirtualMachineBackup) ([]string, error) {
	if obj.Spec.Type == cloudweavv1.Backup || obj.Status == nil {
		return []string{}, nil
	}

	result := make([]string, 0, len(obj.Status.VolumeBackups))
	for _, volumeBackup := range obj.Status.VolumeBackups {
		pvc := volumeBackup.PersistentVolumeClaim
		result = append(result, fmt.Sprintf("%s/%s", pvc.ObjectMeta.Namespace, pvc.ObjectMeta.Name))
	}
	return result, nil
}

func vmBackupByIsProgressing(obj *cloudweavv1.VirtualMachineBackup) ([]string, error) {
	isProgressingStr := strconv.FormatBool(backup.IsBackupProgressing(obj))
	return []string{string(isProgressingStr)}, nil
}

func vmBackupByStorageClassName(obj *cloudweavv1.VirtualMachineBackup) ([]string, error) {
	storageClassNames := []string{}
	if obj.Status == nil {
		return storageClassNames, nil
	}

	for _, volumeBackup := range obj.Status.VolumeBackups {
		storageClassNames = append(storageClassNames, *volumeBackup.PersistentVolumeClaim.Spec.StorageClassName)
	}
	return storageClassNames, nil
}

func vmRestoreByTargetNamespaceAndName(obj *cloudweavv1.VirtualMachineRestore) ([]string, error) {
	if obj == nil {
		return []string{}, nil
	}
	return []string{fmt.Sprintf("%s-%s", obj.Namespace, obj.Spec.Target.Name)}, nil
}

func vmRestoreByVMBackupNamespaceAndName(obj *cloudweavv1.VirtualMachineRestore) ([]string, error) {
	if obj == nil {
		return []string{}, nil
	}
	return []string{fmt.Sprintf("%s-%s", obj.Spec.VirtualMachineBackupNamespace, obj.Spec.VirtualMachineBackupName)}, nil
}

func VolumeByReplicaCount(obj *lhv1beta2.Volume) ([]string, error) {
	replicaCount := strconv.Itoa(obj.Spec.NumberOfReplicas)
	return []string{replicaCount}, nil
}

func imageByExportSourcePVC(obj *cloudweavv1.VirtualMachineImage) ([]string, error) {
	if obj.Spec.SourceType != longhorntypes.LonghornLabelExportFromVolume ||
		obj.Spec.PVCNamespace == "" || obj.Spec.PVCName == "" {
		return nil, nil
	}

	return []string{fmt.Sprintf("%s/%s", obj.Spec.PVCNamespace, obj.Spec.PVCName)}, nil
}

func scheduleVMBackupBySourceVM(obj *cloudweavv1.ScheduleVMBackup) ([]string, error) {
	return []string{fmt.Sprintf("%s/%s", obj.Namespace, obj.Spec.VMBackupSpec.Source.Name)}, nil
}

func scheduleVMBackupByCronGranularity(obj *cloudweavv1.ScheduleVMBackup) ([]string, error) {
	if obj == nil {
		return []string{}, nil
	}

	granularity, err := util.GetCronGranularity(obj)
	if err != nil {
		return []string{}, err
	}

	return []string{granularity.String()}, nil
}

func scheduleVMBackupBySuspended(obj *cloudweavv1.ScheduleVMBackup) ([]string, error) {
	suspenedStr := strconv.FormatBool(obj.Status.Suspended)
	return []string{string(suspenedStr)}, nil
}

func imageByStorageClass(obj *cloudweavv1.VirtualMachineImage) ([]string, error) {
	sc, ok := obj.Annotations[util.AnnotationStorageClassName]
	if !ok {
		return []string{}, nil
	}
	return []string{sc}, nil
}

func vmInstanceMigrationByVM(obj *kubevirtv1.VirtualMachineInstanceMigration) ([]string, error) {
	return []string{fmt.Sprintf("%s/%s", obj.Namespace, obj.Spec.VMIName)}, nil
}
