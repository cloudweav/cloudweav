package upgrade

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	"github.com/cloudweav/cloudweav/pkg/config"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/scheme"
)

const (
	upgradeControllerName = "cloudweav-upgrade-controller"
	planControllerName    = "cloudweav-plan-controller"
	jobControllerName     = "cloudweav-upgrade-job-controller"
	podControllerName     = "cloudweav-upgrade-pod-controller"
	settingControllerName = "cloudweav-version-setting-controller"
	vmImageControllerName = "cloudweav-upgrade-vm-image-controller"
	secretControllerName  = "cloudweav-upgrade-secret-controller"
	nodeControllerName    = "cloudweav-upgrade-node-controller"
)

func Register(ctx context.Context, management *config.Management, options config.Options) error {
	if !options.HCIMode {
		return nil
	}

	upgrades := management.CloudweavFactory.Cloudweavhci().V1beta1().Upgrade()
	upgradeLogs := management.CloudweavFactory.Cloudweavhci().V1beta1().UpgradeLog()
	versions := management.CloudweavFactory.Cloudweavhci().V1beta1().Version()
	settings := management.CloudweavFactory.Cloudweavhci().V1beta1().Setting()
	plans := management.UpgradeFactory.Upgrade().V1().Plan()
	nodes := management.CoreFactory.Core().V1().Node()
	jobs := management.BatchFactory.Batch().V1().Job()
	pods := management.CoreFactory.Core().V1().Pod()
	vmImages := management.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage()
	vms := management.VirtFactory.Kubevirt().V1().VirtualMachine()
	services := management.CoreFactory.Core().V1().Service()
	namespaces := management.CoreFactory.Core().V1().Namespace()
	clusters := management.ProvisioningFactory.Provisioning().V1().Cluster()
	machines := management.ClusterFactory.Cluster().V1beta1().Machine()
	secrets := management.CoreFactory.Core().V1().Secret()
	pvcs := management.CoreFactory.Core().V1().PersistentVolumeClaim()
	lhSettings := management.LonghornFactory.Longhorn().V1beta2().Setting()
	kubeVirt := management.VirtFactory.Kubevirt().V1().KubeVirt()

	virtSubsrcConfig := rest.CopyConfig(management.RestConfig)
	virtSubsrcConfig.GroupVersion = &schema.GroupVersion{Group: "subresources.kubevirt.io", Version: "v1"}
	virtSubsrcConfig.APIPath = "/apis"
	virtSubsrcConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	virtSubresourceClient, err := rest.RESTClientFor(virtSubsrcConfig)
	if err != nil {
		return err
	}

	controller := &upgradeHandler{
		ctx:               ctx,
		jobClient:         jobs,
		jobCache:          jobs.Cache(),
		nodeCache:         nodes.Cache(),
		namespace:         options.Namespace,
		upgradeClient:     upgrades,
		upgradeCache:      upgrades.Cache(),
		upgradeController: upgrades,
		upgradeLogClient:  upgradeLogs,
		upgradeLogCache:   upgradeLogs.Cache(),
		versionCache:      versions.Cache(),
		planClient:        plans,
		planCache:         plans.Cache(),
		vmImageClient:     vmImages,
		vmImageCache:      vmImages.Cache(),
		vmClient:          vms,
		vmCache:           vms.Cache(),
		serviceClient:     services,
		pvcClient:         pvcs,
		clusterClient:     clusters,
		clusterCache:      clusters.Cache(),
		lhSettingClient:   lhSettings,
		lhSettingCache:    lhSettings.Cache(),
		kubeVirtCache:     kubeVirt.Cache(),
		vmRestClient:      virtSubresourceClient,
	}
	upgrades.OnChange(ctx, upgradeControllerName, controller.OnChanged)
	upgrades.OnRemove(ctx, upgradeControllerName, controller.OnRemove)

	planHandler := &planHandler{
		namespace:     options.Namespace,
		upgradeClient: upgrades,
		upgradeCache:  upgrades.Cache(),
		nodeCache:     nodes.Cache(),
		planClient:    plans,
	}
	plans.OnChange(ctx, planControllerName, planHandler.OnChanged)

	jobHandler := &jobHandler{
		namespace:     options.Namespace,
		planCache:     plans.Cache(),
		upgradeClient: upgrades,
		upgradeCache:  upgrades.Cache(),
		machineCache:  machines.Cache(),
		secretClient:  secrets,
		nodeClient:    nodes,
		nodeCache:     nodes.Cache(),
	}
	jobs.OnChange(ctx, jobControllerName, jobHandler.OnChanged)

	podHandler := &podHandler{
		namespace:     options.Namespace,
		planCache:     plans.Cache(),
		upgradeClient: upgrades,
		upgradeCache:  upgrades.Cache(),
	}
	pods.OnChange(ctx, podControllerName, podHandler.OnChanged)

	vmImageHandler := &vmImageHandler{
		namespace:     options.Namespace,
		upgradeClient: upgrades,
		upgradeCache:  upgrades.Cache(),
	}
	vmImages.OnChange(ctx, vmImageControllerName, vmImageHandler.OnChanged)

	secretHandler := &secretHandler{
		namespace:     options.Namespace,
		upgradeClient: upgrades,
		upgradeCache:  upgrades.Cache(),
		jobClient:     jobs,
		jobCache:      jobs.Cache(),
		machineCache:  machines.Cache(),
	}
	secrets.OnChange(ctx, secretControllerName, secretHandler.OnChanged)

	nodeHandler := &nodeHandler{
		namespace:     options.Namespace,
		nodeClient:    nodes,
		nodeCache:     nodes.Cache(),
		upgradeClient: upgrades,
		upgradeCache:  upgrades.Cache(),
		secretClient:  secrets,
	}
	nodes.OnChange(ctx, nodeControllerName, nodeHandler.OnChanged)

	versionSyncer := newVersionSyncer(ctx, options.Namespace, versions, nodes, namespaces)

	settingHandler := settingHandler{
		versionSyncer: versionSyncer,
	}
	settings.OnChange(ctx, settingControllerName, settingHandler.OnChanged)

	go versionSyncer.start()

	return nil
}
