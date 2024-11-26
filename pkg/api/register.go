package api

import (
	"context"

	"github.com/rancher/steve/pkg/server"

	"github.com/cloudweav/cloudweav/pkg/api/cluster"
	"github.com/cloudweav/cloudweav/pkg/api/image"
	"github.com/cloudweav/cloudweav/pkg/api/keypair"
	"github.com/cloudweav/cloudweav/pkg/api/namespace"
	"github.com/cloudweav/cloudweav/pkg/api/node"
	"github.com/cloudweav/cloudweav/pkg/api/upgradelog"
	"github.com/cloudweav/cloudweav/pkg/api/vm"
	"github.com/cloudweav/cloudweav/pkg/api/vmtemplate"
	"github.com/cloudweav/cloudweav/pkg/api/volume"
	"github.com/cloudweav/cloudweav/pkg/api/volumesnapshot"
	"github.com/cloudweav/cloudweav/pkg/config"
)

type registerSchema func(scaled *config.Scaled, server *server.Server, options config.Options) error

func registerSchemas(scaled *config.Scaled, server *server.Server, options config.Options, registers ...registerSchema) error {
	for _, register := range registers {
		if err := register(scaled, server, options); err != nil {
			return err
		}
	}
	return nil
}

func Setup(ctx context.Context, server *server.Server, _ *server.Controllers, options config.Options) error {
	scaled := config.ScaledWithContext(ctx)
	return registerSchemas(scaled, server, options,
		image.RegisterSchema,
		keypair.RegisterSchema,
		vmtemplate.RegisterSchema,
		vm.RegisterSchema,
		node.RegisterSchema,
		upgradelog.RegisterSchema,
		volume.RegisterSchema,
		volumesnapshot.RegisterSchema,
		cluster.RegisterSchema,
		namespace.RegisterSchema,
	)
}
