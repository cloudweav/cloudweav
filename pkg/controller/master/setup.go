package master

import (
	"context"

	"github.com/rancher/steve/pkg/server"
	"github.com/rancher/wrangler/v3/pkg/leader"

	"github.com/cloudweav/cloudweav/pkg/config"
	"github.com/cloudweav/cloudweav/pkg/controller/master/addon"
	"github.com/cloudweav/cloudweav/pkg/controller/master/backup"
	"github.com/cloudweav/cloudweav/pkg/controller/master/image"
	"github.com/cloudweav/cloudweav/pkg/controller/master/keypair"
	"github.com/cloudweav/cloudweav/pkg/controller/master/machine"
	"github.com/cloudweav/cloudweav/pkg/controller/master/mcmsettings"
	"github.com/cloudweav/cloudweav/pkg/controller/master/migration"
	"github.com/cloudweav/cloudweav/pkg/controller/master/node"
	"github.com/cloudweav/cloudweav/pkg/controller/master/nodedrain"
	"github.com/cloudweav/cloudweav/pkg/controller/master/rancher"
	"github.com/cloudweav/cloudweav/pkg/controller/master/schedulevmbackup"
	"github.com/cloudweav/cloudweav/pkg/controller/master/setting"
	"github.com/cloudweav/cloudweav/pkg/controller/master/storagenetwork"
	"github.com/cloudweav/cloudweav/pkg/controller/master/supportbundle"
	"github.com/cloudweav/cloudweav/pkg/controller/master/template"
	"github.com/cloudweav/cloudweav/pkg/controller/master/upgrade"
	"github.com/cloudweav/cloudweav/pkg/controller/master/upgradelog"
	"github.com/cloudweav/cloudweav/pkg/controller/master/virtualmachine"
)

type registerFunc func(context.Context, *config.Management, config.Options) error

var registerFuncs = []registerFunc{
	image.Register,
	keypair.Register,
	migration.Register,
	node.PromoteRegister,
	node.MaintainRegister,
	node.DownRegister,
	node.RemoveRegister,
	node.VolumeDetachRegister,
	node.CPUManagerRegister,
	machine.ControlPlaneRegister,
	setting.Register,
	template.Register,
	virtualmachine.Register,
	backup.RegisterBackup,
	backup.RegisterRestore,
	backup.RegisterBackupTarget,
	backup.RegisterBackupMetadata,
	backup.RegisterBackupBackingImage,
	supportbundle.Register,
	rancher.Register,
	upgrade.Register,
	upgradelog.Register,
	addon.Register,
	storagenetwork.Register,
	nodedrain.Register,
	mcmsettings.Register,
	schedulevmbackup.Register,
}

func register(ctx context.Context, management *config.Management, options config.Options) error {
	for _, f := range registerFuncs {
		if err := f(ctx, management, options); err != nil {
			return err
		}
	}

	return nil
}

func Setup(ctx context.Context, _ *server.Server, controllers *server.Controllers, options config.Options) error {
	scaled := config.ScaledWithContext(ctx)

	go leader.RunOrDie(ctx, "", "cloudweav-controllers", controllers.K8s, func(ctx context.Context) {
		if err := register(ctx, scaled.Management, options); err != nil {
			panic(err)
		}
		if err := scaled.Management.Start(options.Threadiness); err != nil {
			panic(err)
		}
		<-ctx.Done()
	})

	return nil
}
