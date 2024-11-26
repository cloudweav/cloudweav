package clients

import (
	"context"

	ctlfleetv1 "github.com/rancher/rancher/pkg/generated/controllers/fleet.cattle.io"
	rancherv3 "github.com/rancher/rancher/pkg/generated/controllers/management.cattle.io"
	"github.com/rancher/wrangler/v3/pkg/clients"
	ctrlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core"
	storagev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/storage"
	"github.com/rancher/wrangler/v3/pkg/schemes"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/client-go/rest"

	ctlclusterv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cluster.x-k8s.io"
	ctlcloudweavcorev1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/core"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io"
	ctlcniv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/k8s.cni.cncf.io"
	ctlkubevirtv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/kubevirt.io"
	ctllonghornv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/longhorn.io"
	ctlnetwork "github.com/cloudweav/cloudweav/pkg/generated/controllers/network.cloudweavhci.io"
	ctlsnapshotv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/snapshot.storage.k8s.io"
)

type Clients struct {
	clients.Clients

	CloudweavFactory         *ctlcloudweavv1.Factory
	CloudweavCoreFactory     *ctlcloudweavcorev1.Factory
	KubevirtFactory          *ctlkubevirtv1.Factory
	CNIFactory               *ctlcniv1.Factory
	SnapshotFactory          *ctlsnapshotv1.Factory
	FleetFactory             *ctlfleetv1.Factory
	StorageFactory           *storagev1.Factory
	LonghornFactory          *ctllonghornv1.Factory
	ClusterFactory           *ctlclusterv1.Factory
	RancherManagementFactory *rancherv3.Factory
	CoreFactory              *ctrlcorev1.Factory
	CloudweavNetworkFactory  *ctlnetwork.Factory
}

func New(ctx context.Context, rest *rest.Config, threadiness int) (*Clients, error) {
	clients, err := clients.NewFromConfig(rest, nil)
	if err != nil {
		return nil, err
	}

	if err := schemes.Register(v1.AddToScheme); err != nil {
		return nil, err
	}

	cloudweavFactory, err := ctlcloudweavv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = cloudweavFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	cloudweavCoreFactory, err := ctlcloudweavcorev1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = cloudweavCoreFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	kubevirtFactory, err := ctlkubevirtv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = kubevirtFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	cniFactory, err := ctlcniv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = cniFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	snapshotFactory, err := ctlsnapshotv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = snapshotFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	fleetFactory, err := ctlfleetv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = fleetFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	storageFactory, err := storagev1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	if err = storageFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	longhornFactory, err := ctllonghornv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	clusterFactory, err := ctlclusterv1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	rancherFactory, err := rancherv3.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	coreFactory, err := ctrlcorev1.NewFactoryFromConfigWithOptions(rest, clients.FactoryOptions)
	if err != nil {
		return nil, err
	}

	cloudweavNetworkFactory, err := ctlnetwork.NewFactoryFromConfigWithOptions(rest, (*ctlnetwork.FactoryOptions)(clients.FactoryOptions))
	if err != nil {
		return nil, err
	}

	if err = cloudweavNetworkFactory.Start(ctx, threadiness); err != nil {
		return nil, err
	}

	return &Clients{
		Clients:                  *clients,
		CloudweavFactory:         cloudweavFactory,
		CloudweavCoreFactory:     cloudweavCoreFactory,
		KubevirtFactory:          kubevirtFactory,
		CNIFactory:               cniFactory,
		SnapshotFactory:          snapshotFactory,
		FleetFactory:             fleetFactory,
		StorageFactory:           storageFactory,
		LonghornFactory:          longhornFactory,
		ClusterFactory:           clusterFactory,
		RancherManagementFactory: rancherFactory,
		CoreFactory:              coreFactory,
		CloudweavNetworkFactory:  cloudweavNetworkFactory,
	}, nil
}
