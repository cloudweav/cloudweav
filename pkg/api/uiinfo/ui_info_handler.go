package uiinfo

import (
	"net/http"
	"os"

	"github.com/cloudweav/cloudweav/pkg/config"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/settings"
	"github.com/cloudweav/cloudweav/pkg/util"
)

type Handler struct {
	settingsCache ctlcloudweavv1.SettingCache
}

func NewUIInfoHandler(scaled *config.Scaled, _ config.Options) *Handler {
	return &Handler{
		settingsCache: scaled.CloudweavFactory.Cloudweavhci().V1beta1().Setting().Cache(),
	}
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	uiSource := settings.UISource.Get()
	if uiSource == "auto" {
		if !settings.IsRelease() {
			uiSource = "external"
		} else {
			uiSource = "bundled"
		}
	}
	util.ResponseOKWithBody(rw, map[string]string{
		settings.UISourceSettingName:               uiSource,
		settings.UIIndexSettingName:                settings.UIIndex.Get(),
		settings.UIPluginIndexSettingName:          settings.UIPluginIndex.Get(),
		settings.UIPluginBundledVersionSettingName: os.Getenv(settings.GetEnvKey(settings.UIPluginBundledVersionSettingName)),
	})
}
