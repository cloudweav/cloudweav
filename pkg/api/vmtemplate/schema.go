package vmtemplate

import (
	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/schema"
	"github.com/rancher/steve/pkg/server"

	"github.com/cloudweav/cloudweav/pkg/config"
)

const (
	templateSchemaID        = "cloudweavhci.io.virtualmachinetemplate"
	templateVersionSchemaID = "cloudweavhci.io.virtualmachinetemplateversion"
)

func RegisterSchema(scaled *config.Scaled, server *server.Server, _ config.Options) error {
	templateVersionCache := scaled.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineTemplateVersion().Cache()
	th := &templateLinkHandler{
		templateVersionCache: templateVersionCache,
	}

	t := []schema.Template{
		{
			ID:        templateSchemaID,
			Formatter: formatter,
			Customize: func(apiSchema *types.APISchema) {
				apiSchema.ByIDHandler = th.byIDHandler
			},
		},
		{
			ID:        templateVersionSchemaID,
			Formatter: versionFormatter,
		},
	}

	server.SchemaFactory.AddTemplate(t...)
	return nil
}
