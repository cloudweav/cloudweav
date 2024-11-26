package server

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rancher/wrangler/v3/pkg/webhook"

	"github.com/cloudweav/cloudweav/pkg/webhook/clients"
	"github.com/cloudweav/cloudweav/pkg/webhook/config"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/addon"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/bundle"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/bundledeployment"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/keypair"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/managedchart"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/namespace"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/node"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/persistentvolumeclaim"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/resourcequota"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/schedulevmbackup"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/secret"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/setting"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/storageclass"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/templateversion"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/upgrade"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/version"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachine"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachinebackup"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachineimage"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/virtualmachinerestore"
	"github.com/cloudweav/cloudweav/pkg/webhook/resources/volumesnapshot"
	"github.com/cloudweav/cloudweav/pkg/webhook/types"
	"github.com/cloudweav/cloudweav/pkg/webhook/util"
)

func Validation(clients *clients.Clients, options *config.Options) (http.Handler, []types.Resource, error) {
	bearToken, err := ioutil.ReadFile(clients.RESTConfig.BearerTokenFile)
	if err != nil {
		return nil, nil, err
	}
	transport, err := util.GetHTTPTransportWithCertificates(clients.RESTConfig)
	if err != nil {
		return nil, nil, err
	}

	resources := []types.Resource{}
	validators := []types.Validator{
		node.NewValidator(
			clients.Core.Node().Cache(),
			clients.Batch.Job().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstance().Cache()),
		persistentvolumeclaim.NewValidator(
			clients.Core.PersistentVolumeClaim().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage().Cache()),
		keypair.NewValidator(clients.CloudweavFactory.Cloudweavhci().V1beta1().KeyPair().Cache()),
		virtualmachine.NewValidator(
			clients.Core.Namespace().Cache(),
			clients.Core.Pod().Cache(),
			clients.Core.PersistentVolumeClaim().Cache(),
			clients.CloudweavCoreFactory.Core().V1().ResourceQuota().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstanceMigration().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstance().Cache()),
		virtualmachineimage.NewValidator(
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage().Cache(),
			clients.Core.PersistentVolumeClaim().Cache(),
			clients.K8s.AuthorizationV1().SelfSubjectAccessReviews(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineTemplateVersion().Cache(),
			clients.Core.Secret().Cache(),
			clients.StorageFactory.Storage().V1().StorageClass().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache()),
		upgrade.NewValidator(
			clients.CloudweavFactory.Cloudweavhci().V1beta1().Upgrade().Cache(),
			clients.Core.Node().Cache(),
			clients.LonghornFactory.Longhorn().V1beta2().Volume().Cache(),
			clients.ClusterFactory.Cluster().V1beta1().Cluster().Cache(),
			clients.ClusterFactory.Cluster().V1beta1().Machine().Cache(),
			clients.RancherManagementFactory.Management().V3().ManagedChart().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().Version().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().ScheduleVMBackup().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstance().Cache(),
			&http.Client{
				Transport: transport,
				Timeout:   time.Second * 20,
			},
			string(bearToken),
		),
		virtualmachinebackup.NewValidator(
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().Setting().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineRestore().Cache(),
			clients.CoreFactory.Core().V1().PersistentVolumeClaim().Cache(),
			clients.LonghornFactory.Longhorn().V1beta2().Engine().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().ResourceQuota().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstanceMigration().Cache(),
		),
		virtualmachinerestore.NewValidator(
			clients.Core.Namespace().Cache(),
			clients.Core.Pod().Cache(),
			clients.CloudweavCoreFactory.Core().V1().ResourceQuota().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().Setting().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineRestore().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().ScheduleVMBackup().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstanceMigration().Cache(),
			clients.SnapshotFactory.Snapshot().V1().VolumeSnapshotClass().Cache(),
			clients.CNIFactory.K8s().V1().NetworkAttachmentDefinition().Cache(),
		),
		setting.NewValidator(
			clients.CloudweavFactory.Cloudweavhci().V1beta1().Setting().Cache(),
			clients.Core.Node().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineBackup().Cache(),
			clients.SnapshotFactory.Snapshot().V1().VolumeSnapshotClass().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineRestore().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachineInstance().Cache(),
			clients.RancherManagementFactory.Management().V3().Feature().Cache(),
			clients.LonghornFactory.Longhorn().V1beta2().Volume().Cache(),
			clients.CoreFactory.Core().V1().PersistentVolumeClaim().Cache(),
			clients.CloudweavNetworkFactory.Network().V1beta1().ClusterNetwork().Cache(),
			clients.CloudweavNetworkFactory.Network().V1beta1().VlanConfig().Cache(),
			clients.CloudweavNetworkFactory.Network().V1beta1().VlanStatus().Cache(),
		),
		templateversion.NewValidator(
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineTemplate().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineTemplateVersion().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().KeyPair().Cache()),
		managedchart.NewValidator(),
		bundle.NewValidator(),
		bundledeployment.NewValidator(
			clients.FleetFactory.Fleet().V1alpha1().Cluster().Cache(),
		),
		storageclass.NewValidator(
			clients.StorageFactory.Storage().V1().StorageClass().Cache(),
			clients.Core.Secret().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage().Cache()),
		namespace.NewValidator(clients.CloudweavCoreFactory.Core().V1().ResourceQuota().Cache()),
		addon.NewValidator(clients.CloudweavFactory.Cloudweavhci().V1beta1().Addon().Cache()),
		version.NewValidator(),
		volumesnapshot.NewValidator(
			clients.CoreFactory.Core().V1().PersistentVolumeClaim().Cache(),
			clients.LonghornFactory.Longhorn().V1beta2().Engine().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().ResourceQuota().Cache(),
			clients.KubevirtFactory.Kubevirt().V1().VirtualMachine().Cache(),
		),
		resourcequota.NewValidator(),
		schedulevmbackup.NewValidator(
			clients.CloudweavFactory.Cloudweavhci().V1beta1().Setting().Cache(),
			clients.Core.Secret().Cache(),
			clients.CloudweavFactory.Cloudweavhci().V1beta1().ScheduleVMBackup().Cache(),
		),
		secret.NewValidator(clients.StorageFactory.Storage().V1().StorageClass().Cache()),
	}

	router := webhook.NewRouter()
	for _, v := range validators {
		addHandler(router, types.AdmissionTypeValidation, types.NewValidatorAdapter(v), options)
		resources = append(resources, v.Resource())
	}

	return router, resources, nil
}
