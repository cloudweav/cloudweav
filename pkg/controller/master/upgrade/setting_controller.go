package upgrade

import (
	"github.com/sirupsen/logrus"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
)

// settingHandler do version syncs on server-version setting changes
type settingHandler struct {
	versionSyncer *versionSyncer
}

func (h *settingHandler) OnChanged(_ string, setting *cloudweavv1.Setting) (*cloudweavv1.Setting, error) {
	if setting == nil || setting.DeletionTimestamp != nil || setting.Name != "server-version" {
		return setting, nil
	}
	if err := h.versionSyncer.sync(); err != nil {
		logrus.Errorf("failed syncing version metadata: %v", err)
	}
	return setting, nil
}
