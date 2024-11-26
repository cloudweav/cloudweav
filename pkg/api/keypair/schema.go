package keypair

import (
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/schema"
	"github.com/rancher/steve/pkg/server"
	"github.com/rancher/wrangler/v3/pkg/schemas"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/config"
)

const (
	keygen = "keygen"
)

func RegisterSchema(scaled *config.Scaled, server *server.Server, _ config.Options) error {
	server.BaseSchemas.MustImportAndCustomize(cloudweavv1.KeyGenInput{}, nil)
	t := schema.Template{
		ID: "cloudweavhci.io.keypair",
		Customize: func(s *types.APISchema) {
			s.CollectionFormatter = CollectionFormatter
			s.CollectionActions = map[string]schemas.Action{
				keygen: {
					Input: "keyGenInput",
				},
			}
			s.Formatter = Formatter
			s.ActionHandlers = map[string]http.Handler{
				keygen: KeyGenActionHandler{
					KeyPairs:     scaled.CloudweavFactory.Cloudweavhci().V1beta1().KeyPair(),
					KeyPairCache: scaled.CloudweavFactory.Cloudweavhci().V1beta1().KeyPair().Cache(),
				},
			}
		},
	}
	server.SchemaFactory.AddTemplate(t)
	return nil
}
