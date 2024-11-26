package data

import (
	"context"

	nodev1 "github.com/cloudweav/node-manager/pkg/apis/node.cloudweavhci.io/v1beta1"
	loggingv1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	lhv1beta2 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	fleetv1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	rancherv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	provisioningv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	upgradev1 "github.com/rancher/system-upgrade-controller/pkg/apis/upgrade.cattle.io/v1"
	"k8s.io/client-go/rest"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/util/crd"
)

func createCRDs(ctx context.Context, restConfig *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(ctx, restConfig)
	if err != nil {
		return err
	}
	return factory.
		BatchCreateCRDsIfNotExisted(
			crd.NonNamespacedFromGV(cloudweavv1.SchemeGroupVersion, "Setting", cloudweavv1.Setting{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "APIService", rancherv3.APIService{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "Setting", rancherv3.Setting{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "User", rancherv3.User{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "Group", rancherv3.Group{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "GroupMember", rancherv3.GroupMember{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "Token", rancherv3.Token{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "NodeDriver", rancherv3.NodeDriver{}),
			crd.NonNamespacedFromGV(rancherv3.SchemeGroupVersion, "Feature", rancherv3.Feature{}),
			crd.NonNamespacedFromGV(upgradev1.SchemeGroupVersion, "Plan", upgradev1.Plan{}),
			crd.NonNamespacedFromGV(loggingv1.GroupVersion, "Logging", loggingv1.Logging{}),
		).
		BatchCreateCRDsIfNotExisted(
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "KeyPair", cloudweavv1.KeyPair{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "Upgrade", cloudweavv1.Upgrade{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "UpgradeLog", cloudweavv1.UpgradeLog{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "Version", cloudweavv1.Version{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "VirtualMachineImage", cloudweavv1.VirtualMachineImage{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "VirtualMachineTemplate", cloudweavv1.VirtualMachineTemplate{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "VirtualMachineTemplateVersion", cloudweavv1.VirtualMachineTemplateVersion{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "VirtualMachineBackup", cloudweavv1.VirtualMachineBackup{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "VirtualMachineRestore", cloudweavv1.VirtualMachineRestore{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "Preference", cloudweavv1.Preference{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "SupportBundle", cloudweavv1.SupportBundle{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "ResourceQuota", cloudweavv1.ResourceQuota{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "ScheduleVMBackup", cloudweavv1.ScheduleVMBackup{}),
			// The BackingImage struct is not compatible with wrangler schemas generation, pass nil as the workaround.
			// The expected CRD will be applied by Longhorn chart.
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "BackingImage", nil),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "BackingImageDataSource", lhv1beta2.BackingImageDataSource{}),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "Backup", lhv1beta2.Backup{}),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "BackupBackingImage", lhv1beta2.BackupBackingImage{}),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "Engine", lhv1beta2.Engine{}),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "Snapshot", lhv1beta2.Snapshot{}),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "Volume", lhv1beta2.Volume{}),
			crd.FromGV(lhv1beta2.SchemeGroupVersion, "Setting", lhv1beta2.Setting{}),
			crd.FromGV(provisioningv1.SchemeGroupVersion, "Cluster", provisioningv1.Cluster{}),
			crd.FromGV(fleetv1alpha1.SchemeGroupVersion, "Cluster", fleetv1alpha1.Cluster{}),
			crd.FromGV(clusterv1.GroupVersion, "Cluster", clusterv1.Cluster{}),
			crd.FromGV(clusterv1.GroupVersion, "Machine", clusterv1.Machine{}),
			crd.FromGV(cloudweavv1.SchemeGroupVersion, "Addon", cloudweavv1.Addon{}).WithStatus(),
			crd.FromGV(monitoringv1.SchemeGroupVersion, "Prometheus", monitoringv1.Prometheus{}),
			crd.FromGV(monitoringv1.SchemeGroupVersion, "Alertmanager", monitoringv1.Alertmanager{}),
			crd.FromGV(loggingv1.GroupVersion, "ClusterFlow", loggingv1.ClusterFlow{}),
			crd.FromGV(loggingv1.GroupVersion, "ClusterOutput", loggingv1.ClusterOutput{}),
			crd.FromGV(nodev1.SchemeGroupVersion, "NodeConfig", nodev1.NodeConfig{}),
			crd.FromGV(rkev1.SchemeGroupVersion, "RKEControlPlane", rkev1.RKEControlPlane{}),
		).
		BatchWait()
}
