package image

import (
	"net/http"

	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/schema"
	"github.com/rancher/steve/pkg/server"
	"github.com/rancher/wrangler/v3/pkg/schemas"

	"github.com/cloudweav/cloudweav/pkg/config"
)

func RegisterSchema(scaled *config.Scaled, server *server.Server, _ config.Options) error {
	imgHandler := Handler{
		httpClient:                  http.Client{},
		Images:                      scaled.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage(),
		ImageCache:                  scaled.CloudweavFactory.Cloudweavhci().V1beta1().VirtualMachineImage().Cache(),
		BackingImageDataSources:     scaled.LonghornFactory.Longhorn().V1beta2().BackingImageDataSource(),
		BackingImageDataSourceCache: scaled.LonghornFactory.Longhorn().V1beta2().BackingImageDataSource().Cache(),
		BackingImageCache:           scaled.LonghornFactory.Longhorn().V1beta2().BackingImage().Cache(),
	}

	t := schema.Template{
		ID: "cloudweavhci.io.virtualmachineimage",
		Customize: func(s *types.APISchema) {
			s.Formatter = Formatter
			s.ResourceActions = map[string]schemas.Action{
				actionUpload: {},
			}
			/*
			 * ActionHandlers would let people define their own `POST` method.
			 * That would add to `actions` on API and would be filled as key-value
			 * pair in the current HTTP requests.
			 */
			s.ActionHandlers = map[string]http.Handler{
				actionUpload: imgHandler,
			}
			/*
			 * LinkHandlers would let people define their own `GET` method.
			 * That would add to `links` on API and would be filled as key-value
			 * pair in the current HTTP requests.
			 *
			 * Detail about `ActionHandlers` and `LinkHandlers` could be found
			 * with rancher/apiserver
			 */
			s.LinkHandlers = map[string]http.Handler{
				actionDownload: imgHandler,
			}
		},
	}
	server.SchemaFactory.AddTemplate(t)
	return nil
}
