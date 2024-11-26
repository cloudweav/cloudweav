package schedulevmbackup

import (
	"context"

	catalogv1 "github.com/rancher/rancher/pkg/generated/controllers/catalog.cattle.io/v1"
	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/config"
	ctlharvbatchv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/batch/v1"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	ctllonghornv2 "github.com/cloudweav/cloudweav/pkg/generated/controllers/longhorn.io/v1beta2"
	ctlsnapshotv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/snapshot.storage.k8s.io/v1"
)

const (
	scheduleVMBackupControllerName = "schedule-vm-bakcup-controller"
	cronJobControllerName          = "cron-job-controller"
	vmBackupControllerName         = "vm-backup-controller"
	longhornBackupControllerName   = "longhorn-backup-controller"

	vmBackupKindName = "VirtualMachineBackup"
)

var vmBackupKind = cloudweavv1.SchemeGroupVersion.WithKind(vmBackupKindName)

type svmbackupHandler struct {
	svmbackupController  ctlcloudweavv1.ScheduleVMBackupController
	svmbackupClient      ctlcloudweavv1.ScheduleVMBackupClient
	svmbackupCache       ctlcloudweavv1.ScheduleVMBackupCache
	cronJobsClient       ctlharvbatchv1.CronJobClient
	cronJobCache         ctlharvbatchv1.CronJobCache
	vmBackupController   ctlcloudweavv1.VirtualMachineBackupController
	vmBackupClient       ctlcloudweavv1.VirtualMachineBackupClient
	vmBackupCache        ctlcloudweavv1.VirtualMachineBackupCache
	snapshotCache        ctlsnapshotv1.VolumeSnapshotCache
	lhsnapshotClient     ctllonghornv2.SnapshotClient
	lhsnapshotCache      ctllonghornv2.SnapshotCache
	settingCache         ctlcloudweavv1.SettingCache
	secretCache          ctlcorev1.SecretCache
	namespace            string
	appCache             catalogv1.AppCache
	lhbackupCache        ctllonghornv2.BackupCache
	lhbackupClient       ctllonghornv2.BackupClient
	snapshotContentCache ctlsnapshotv1.VolumeSnapshotContentCache
}

func Register(ctx context.Context, management *config.Management, options config.Options) error {
	svmbackups := management.CloudweavFactory.Cloudweavhci().V1beta1().ScheduleVMBackup()
	cronJobs := management.CloudweavBatchFactory.Batch().V1().CronJob()
	appCache := management.CatalogFactory.Catalog().V1().App().Cache()
	vmBackups := management.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup()
	snapshots := management.SnapshotFactory.Snapshot().V1().VolumeSnapshot()
	lhsnapshots := management.LonghornFactory.Longhorn().V1beta2().Snapshot()
	settings := management.CloudweavFactory.Cloudweavhci().V1beta1().Setting()
	secrets := management.CoreFactory.Core().V1().Secret()
	lhbackups := management.LonghornFactory.Longhorn().V1beta2().Backup()
	snapshotContents := management.SnapshotFactory.Snapshot().V1().VolumeSnapshotContent()

	svmbackupHandler := &svmbackupHandler{
		svmbackupController:  svmbackups,
		svmbackupClient:      svmbackups,
		svmbackupCache:       svmbackups.Cache(),
		cronJobsClient:       cronJobs,
		cronJobCache:         cronJobs.Cache(),
		vmBackupController:   vmBackups,
		vmBackupClient:       vmBackups,
		vmBackupCache:        vmBackups.Cache(),
		snapshotCache:        snapshots.Cache(),
		lhsnapshotClient:     lhsnapshots,
		lhsnapshotCache:      lhsnapshots.Cache(),
		settingCache:         settings.Cache(),
		secretCache:          secrets.Cache(),
		namespace:            options.Namespace,
		appCache:             appCache,
		lhbackupCache:        lhbackups.Cache(),
		lhbackupClient:       lhbackups,
		snapshotContentCache: snapshotContents.Cache(),
	}

	svmbackups.OnChange(ctx, scheduleVMBackupControllerName, svmbackupHandler.OnChanged)
	svmbackups.OnRemove(ctx, scheduleVMBackupControllerName, svmbackupHandler.OnRemove)
	cronJobs.OnChange(ctx, cronJobControllerName, svmbackupHandler.OnCronjobChanged)
	vmBackups.OnChange(ctx, vmBackupControllerName, svmbackupHandler.OnVMBackupChange)
	vmBackups.OnRemove(ctx, vmBackupControllerName, svmbackupHandler.OnVMBackupRemove)
	lhbackups.OnChange(ctx, longhornBackupControllerName, svmbackupHandler.OnLHBackupChanged)
	return nil
}
