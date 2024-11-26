package util

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlbackup "github.com/cloudweav/cloudweav/pkg/controller/master/backup"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/webhook/indexeres"
)

func HasInProgressingVMBackupBySourceUID(cache ctlcloudweavv1.VirtualMachineBackupCache, sourceUID string) (bool, error) {
	vmBackups, err := cache.GetByIndex(indexeres.VMBackupBySourceUIDIndex, sourceUID)
	if err != nil && !apierrors.IsNotFound(err) {
		return false, err
	}
	for _, vmBackup := range vmBackups {
		if ctlbackup.IsBackupProgressing(vmBackup) || ctlbackup.GetVMBackupError(vmBackup) != nil {
			return true, nil
		}
	}
	return false, nil
}

func HasInProgressingVMRestoreOnSameTarget(cache ctlcloudweavv1.VirtualMachineRestoreCache, targetNamespace, targetName string) (bool, error) {
	vmRestores, err := cache.GetByIndex(indexeres.VMRestoreByTargetNamespaceAndName, fmt.Sprintf("%s-%s", targetNamespace, targetName))
	if err != nil && !apierrors.IsNotFound(err) {
		return false, err
	}

	for _, vmRestore := range vmRestores {
		if vmRestore != nil && vmRestore.Status != nil {
			for _, condition := range vmRestore.Status.Conditions {
				if condition.Type == cloudweavv1.BackupConditionProgressing && condition.Status == v1.ConditionTrue {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
