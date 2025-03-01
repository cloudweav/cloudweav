package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	lhv1beta2 "github.com/longhorn/longhorn-manager/k8s/pkg/apis/longhorn/v1beta2"
	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/config"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	ctlkubevirtv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/kubevirt.io/v1"
	ctllonghornv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/longhorn.io/v1beta2"
	"github.com/cloudweav/cloudweav/pkg/settings"
	"github.com/cloudweav/cloudweav/pkg/util"
)

const (
	backupTargetControllerName = "cloudweav-backup-target-controller"

	longhornBackupTargetSettingName       = "backup-target"
	longhornBackupTargetSecretSettingName = "backup-target-credential-secret"
)

// RegisterBackupTarget register the setting controller and reconsile longhorn setting when backup target changed
func RegisterBackupTarget(ctx context.Context, management *config.Management, _ config.Options) error {
	settings := management.CloudweavFactory.Cloudweavhci().V1beta1().Setting()
	secrets := management.CoreFactory.Core().V1().Secret()
	longhornSettings := management.LonghornFactory.Longhorn().V1beta2().Setting()
	vms := management.VirtFactory.Kubevirt().V1().VirtualMachine()

	backupTargetController := &TargetHandler{
		ctx:                  ctx,
		longhornSettings:     longhornSettings,
		longhornSettingCache: longhornSettings.Cache(),
		secrets:              secrets,
		secretCache:          secrets.Cache(),
		vms:                  vms,
		settings:             settings,
	}

	settings.OnChange(ctx, backupTargetControllerName, backupTargetController.OnBackupTargetChange)
	return nil
}

type TargetHandler struct {
	ctx                  context.Context
	longhornSettings     ctllonghornv1.SettingClient
	longhornSettingCache ctllonghornv1.SettingCache
	secrets              ctlcorev1.SecretClient
	secretCache          ctlcorev1.SecretCache
	vms                  ctlkubevirtv1.VirtualMachineController
	settings             ctlcloudweavv1.SettingClient
}

// OnBackupTargetChange handles backupTarget setting object on change
func (h *TargetHandler) OnBackupTargetChange(_ string, setting *cloudweavv1.Setting) (*cloudweavv1.Setting, error) {
	if setting == nil || setting.DeletionTimestamp != nil ||
		setting.Name != settings.BackupTargetSettingName {
		return setting, nil
	}

	target, err := settings.DecodeBackupTarget(setting.Value)
	if err != nil {
		return h.setConfiguredCondition(setting, "", err)
	}

	logrus.Debugf("backup target change:%s:%s", target.Type, target.Endpoint)

	switch target.Type {
	case settings.S3BackupType:
		// Since S3 access key id and secret access key are stripped after S3 backup target has been verified
		// in reUpdateBackupTargetSettingSecret
		// stop the controller to reconcile it
		if target.SecretAccessKey == "" && target.AccessKeyID == "" {
			break
		}

		if err = h.updateLonghornTarget(target); err != nil {
			return h.setConfiguredCondition(setting, "", err)
		}

		if err = h.updateBackupTargetSecret(target); err != nil {
			return h.setConfiguredCondition(setting, "", err)
		}

		return h.reUpdateBackupTargetSettingSecret(setting, target)

	case settings.NFSBackupType:
		if err = h.updateLonghornTarget(target); err != nil {
			return h.setConfiguredCondition(setting, "", err)
		}

		// reset previous secret of S3
		if err = h.resetBackupTargetSecret(); err != nil {
			return h.setConfiguredCondition(setting, "", err)
		}

	default:
		// reset backup target to default, then delete/update related settings
		if target.IsDefaultBackupTarget() {
			if err = h.updateLonghornTarget(target); err != nil {
				return h.setConfiguredCondition(setting, "", err)
			}

			// reset previous secret of S3
			if err = h.resetBackupTargetSecret(); err != nil {
				return h.setConfiguredCondition(setting, "", err)
			}

			settingCpy := setting.DeepCopy()
			cloudweavv1.SettingConfigured.False(settingCpy)
			cloudweavv1.SettingConfigured.Message(settingCpy, "")
			cloudweavv1.SettingConfigured.Reason(settingCpy, "")
			return h.settings.Update(settingCpy)
		}

		return h.setConfiguredCondition(setting, "", fmt.Errorf("Invalid backup target type:%s or parameter", target.Type))
	}

	if len(setting.Status.Conditions) == 0 || cloudweavv1.SettingConfigured.IsFalse(setting) {
		return h.setConfiguredCondition(setting, "", nil)
	}
	return setting, nil
}

func (h *TargetHandler) reUpdateBackupTargetSettingSecret(setting *cloudweavv1.Setting, target *settings.BackupTarget) (*cloudweavv1.Setting, error) {
	// only do a second update when s3 with credentials
	if target.Type != settings.S3BackupType {
		return nil, nil
	}

	// reset the s3 credentials to prevent controller reconcile and not to expose secret key
	target.SecretAccessKey = ""
	target.AccessKeyID = ""
	targetBytes, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	settingCpy := setting.DeepCopy()
	settingCpy.Value = string(targetBytes)

	return h.settings.Update(settingCpy)
}

func (h *TargetHandler) updateLonghornTarget(backupTarget *settings.BackupTarget) error {
	target, err := h.longhornSettingCache.Get(util.LonghornSystemNamespaceName, longhornBackupTargetSettingName)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		if _, err := h.longhornSettings.Create(&lhv1beta2.Setting{
			ObjectMeta: metav1.ObjectMeta{
				Name:      longhornBackupTargetSettingName,
				Namespace: util.LonghornSystemNamespaceName,
			},
			Value: util.ConstructEndpoint(backupTarget),
		}); err != nil {
			return err
		}
		return nil
	}

	targetCpy := target.DeepCopy()
	targetCpy.Value = util.ConstructEndpoint(backupTarget)

	if !reflect.DeepEqual(target, targetCpy) {
		_, err := h.longhornSettings.Update(targetCpy)
		return err
	}
	return nil
}

func getBackupSecretData(target *settings.BackupTarget) (map[string]string, error) {
	data := map[string]string{
		util.AWSAccessKey:       target.AccessKeyID,
		util.AWSSecretKey:       target.SecretAccessKey,
		util.AWSEndpoints:       target.Endpoint,
		util.AWSCERT:            target.Cert,
		util.VirtualHostedStyle: strconv.FormatBool(target.VirtualHostedStyle),
	}
	if settings.AdditionalCA.Get() != "" {
		data[util.AWSCERT] = settings.AdditionalCA.Get()
	}

	var httpProxyConfig util.HTTPProxyConfig
	if err := json.Unmarshal([]byte(settings.HTTPProxy.Get()), &httpProxyConfig); err != nil {
		return nil, err
	}
	data[util.HTTPProxyEnv] = httpProxyConfig.HTTPProxy
	data[util.HTTPSProxyEnv] = httpProxyConfig.HTTPSProxy
	data[util.NoProxyEnv] = util.AddBuiltInNoProxy(httpProxyConfig.NoProxy)

	// trim spaces for all values in order to pass the Longhorn webhook, refer to https://github.com/longhorn/longhorn-manager/pull/970
	for k, v := range data {
		data[k] = strings.TrimSpace(v)
	}

	return data, nil
}

func (h *TargetHandler) updateBackupTargetSecret(target *settings.BackupTarget) error {
	backupSecretData, err := getBackupSecretData(target)
	if err != nil {
		return err
	}
	secret, err := h.secretCache.Get(util.LonghornSystemNamespaceName, util.BackupTargetSecretName)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		newSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      util.BackupTargetSecretName,
				Namespace: util.LonghornSystemNamespaceName,
			},
		}
		newSecret.StringData = backupSecretData
		if _, err = h.secrets.Create(newSecret); err != nil {
			return err
		}
	} else {
		secretCpy := secret.DeepCopy()
		secretCpy.StringData = backupSecretData
		if !reflect.DeepEqual(secret.StringData, secretCpy.StringData) {
			if _, err := h.secrets.Update(secretCpy); err != nil {
				return err
			}
		}
	}

	return h.updateLonghornBackupTargetSecretSetting()
}

func (h *TargetHandler) resetBackupTargetSecret() error {
	secret, err := h.secretCache.Get(util.LonghornSystemNamespaceName, util.BackupTargetSecretName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	if secret != nil && len(secret.Data) != 0 {
		secretCpy := secret.DeepCopy()
		secretCpy.Data = map[string][]byte{}
		if _, err = h.secrets.Update(secretCpy); err != nil {
			return err
		}
	}

	setting, err := h.longhornSettingCache.Get(util.LonghornSystemNamespaceName, longhornBackupTargetSecretSettingName)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	if setting != nil && setting.Value != "" {
		settingCpy := setting.DeepCopy()
		settingCpy.Value = ""
		if _, err := h.longhornSettings.Update(settingCpy); err != nil {
			return err
		}
	}

	return nil
}

func (h *TargetHandler) updateLonghornBackupTargetSecretSetting() error {
	targetSecret, err := h.longhornSettingCache.Get(util.LonghornSystemNamespaceName, longhornBackupTargetSecretSettingName)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		if _, err := h.longhornSettings.Create(&lhv1beta2.Setting{
			ObjectMeta: metav1.ObjectMeta{
				Name:      longhornBackupTargetSecretSettingName,
				Namespace: util.LonghornSystemNamespaceName,
			},
			Value: util.BackupTargetSecretName,
		}); err != nil {
			return err
		}
		return nil
	}

	targetSecCpy := targetSecret.DeepCopy()
	targetSecCpy.Value = util.BackupTargetSecretName

	if targetSecret.Value != targetSecCpy.Value {
		if _, err := h.longhornSettings.Update(targetSecCpy); err != nil {
			return err
		}
	}

	return nil
}

func (h *TargetHandler) setConfiguredCondition(setting *cloudweavv1.Setting, reason string, err error) (*cloudweavv1.Setting, error) {
	settingCpy := setting.DeepCopy()
	// SetError with nil error will cleanup message in condition and set the status to true
	cloudweavv1.SettingConfigured.SetError(settingCpy, reason, err)
	return h.settings.Update(settingCpy)
}
