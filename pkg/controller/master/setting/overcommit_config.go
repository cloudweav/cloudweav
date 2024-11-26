package setting

import (
	"encoding/json"
	"fmt"
	"strconv"

	longhorn "github.com/longhorn/longhorn-manager/types"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/settings"
	"github.com/cloudweav/cloudweav/pkg/util"
)

func (h *Handler) syncOvercommitConfig(setting *cloudweavv1.Setting) error {
	overcommit := &settings.Overcommit{}
	var value string
	if setting.Value != "" {
		value = setting.Value
	} else {
		value = setting.Default
	}
	if err := json.Unmarshal([]byte(value), overcommit); err != nil {
		return fmt.Errorf("Invalid JSON `%s`: %s", setting.Value, err.Error())
	}

	// Longhorn storage overcommit
	storage, err := h.longhornSettingCache.Get(util.LonghornSystemNamespaceName, string(longhorn.SettingNameStorageOverProvisioningPercentage))
	if err != nil {
		return err
	}
	storageCpy := storage.DeepCopy()
	percentage := strconv.Itoa(overcommit.Storage)
	if storageCpy.Value != percentage {
		storageCpy.Value = percentage
		if _, err := h.longhornSettings.Update(storageCpy); err != nil {
			return err
		}
	}

	return nil
}
