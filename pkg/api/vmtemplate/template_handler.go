package vmtemplate

import (
	"fmt"
	"net/http"

	"github.com/rancher/apiserver/pkg/apierror"
	"github.com/rancher/apiserver/pkg/handlers"
	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/wrangler/v3/pkg/schemas/validation"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/cloudweav/cloudweav/pkg/controller/master/template"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
)

type templateLinkHandler struct {
	templateVersionCache ctlcloudweavv1.VirtualMachineTemplateVersionCache
}

func (h *templateLinkHandler) byIDHandler(request *types.APIRequest) (types.APIObject, error) {
	if request.Link == "versions" {
		versions, err := h.getVersions(request.Namespace, request.Name)
		if err != nil {
			return types.APIObject{}, err
		}

		request.ResponseWriter.WriteList(request, http.StatusOK, versions)
	}

	return handlers.ByIDHandler(request)
}

func (h *templateLinkHandler) getVersions(templateNs, templateName string) (types.APIObjectList, error) {
	sets := labels.Set{
		template.TemplateLabel: templateName,
	}
	versions, err := h.templateVersionCache.List(templateNs, sets.AsSelector())
	if err != nil {
		return types.APIObjectList{}, apierror.NewAPIError(validation.ServerError, err.Error())
	}

	result := make([]types.APIObject, 0, len(versions))
	for _, vtr := range versions {
		id := fmt.Sprintf("%s/%s", vtr.Namespace, vtr.Name)
		result = append(result, types.APIObject{
			Type:   templateVersionSchemaID,
			ID:     id,
			Object: vtr,
		})
	}

	return types.APIObjectList{
		Objects: result,
	}, nil
}
