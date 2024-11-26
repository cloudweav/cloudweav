package setting

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ctlnodev1 "github.com/cloudweav/node-manager/pkg/generated/controllers/node.cloudweavhci.io/v1beta1"
	ctlhelmv1 "github.com/k3s-io/helm-controller/pkg/generated/controllers/helm.cattle.io/v1"
	catalogv1api "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	catalogv1 "github.com/rancher/rancher/pkg/generated/controllers/catalog.cattle.io/v1"
	ctlmgmtv3 "github.com/rancher/rancher/pkg/generated/controllers/management.cattle.io/v3"
	provisioningv1 "github.com/rancher/rancher/pkg/generated/controllers/provisioning.cattle.io/v1"
	ctlrkev1 "github.com/rancher/rancher/pkg/generated/controllers/rke.cattle.io/v1"
	"github.com/rancher/wrangler/v3/pkg/apply"
	v1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/apps/v1"
	ctlcorev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/v3/pkg/slice"
	"k8s.io/apimachinery/pkg/api/errors"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	kubevirtv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/kubevirt.io/v1"
	ctllhv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/longhorn.io/v1beta2"
	networkingv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/networking.k8s.io/v1"
	"github.com/cloudweav/cloudweav/pkg/settings"
	"github.com/cloudweav/cloudweav/pkg/util"
)

type syncerFunc func(*cloudweavv1.Setting) error

var (
	syncers map[string]syncerFunc
	// bootstrapSettings are the setting that syncs on bootstrap
	bootstrapSettings = []string{
		settings.SSLCertificatesSettingName,
		settings.KubeconfigDefaultTokenTTLMinutesSettingName,
		// The Longhorn storage over-provisioning percentage is set to 100, whereas Cloudweav uses 200.
		// This needs to be synchronized when Cloudweav starts.
		settings.OvercommitConfigSettingName,
		// always run this when Cloudweav POD starts
		settings.AdditionalGuestMemoryOverheadRatioName,
	}
	skipHashCheckSettings = []string{
		settings.AutoRotateRKE2CertsSettingName,
		settings.LogLevelSettingName,
		settings.KubeconfigDefaultTokenTTLMinutesSettingName,
		settings.AdditionalGuestMemoryOverheadRatioName,
	}
)

type Handler struct {
	namespace            string
	httpClient           http.Client
	apply                apply.Apply
	clusterCache         provisioningv1.ClusterCache
	clusters             provisioningv1.ClusterClient
	settings             v1beta1.SettingClient
	settingCache         v1beta1.SettingCache
	settingController    v1beta1.SettingController
	secrets              ctlcorev1.SecretClient
	secretCache          ctlcorev1.SecretCache
	deployments          v1.DeploymentClient
	deploymentCache      v1.DeploymentCache
	ingresses            networkingv1.IngressClient
	ingressCache         networkingv1.IngressCache
	longhornSettings     ctllhv1.SettingClient
	longhornSettingCache ctllhv1.SettingCache
	configmaps           ctlcorev1.ConfigMapClient
	configmapCache       ctlcorev1.ConfigMapCache
	serviceCache         ctlcorev1.ServiceCache
	apps                 catalogv1.AppClient
	managedCharts        ctlmgmtv3.ManagedChartClient
	managedChartCache    ctlmgmtv3.ManagedChartCache
	helmChartConfigs     ctlhelmv1.HelmChartConfigClient
	helmChartConfigCache ctlhelmv1.HelmChartConfigCache
	nodeClient           ctlcorev1.NodeController
	nodeCache            ctlcorev1.NodeCache
	nodeConfigs          ctlnodev1.NodeConfigClient
	nodeConfigsCache     ctlnodev1.NodeConfigCache
	rkeControlPlaneCache ctlrkev1.RKEControlPlaneCache
	rancherSettings      ctlmgmtv3.SettingClient
	rancherSettingsCache ctlmgmtv3.SettingCache
	kubeVirtConfig       kubevirtv1.KubeVirtClient
	kubeVirtConfigCache  kubevirtv1.KubeVirtCache
}

func (h *Handler) settingOnChanged(_ string, setting *cloudweavv1.Setting) (*cloudweavv1.Setting, error) {
	if setting == nil || setting.DeletionTimestamp != nil {
		return nil, nil
	}

	// The setting value hash is stored in the annotation when a setting syncer completes.
	// So that we only proceed when value is changed.
	if setting.Value == "" && setting.Annotations[util.AnnotationHash] == "" &&
		!slice.ContainsString(bootstrapSettings, setting.Name) {
		return nil, nil
	}

	toMeasure := io.MultiReader(
		strings.NewReader(setting.Value),
		strings.NewReader(setting.Annotations[util.AnnotationUpgradePatched]),
	)

	hash := sha256.New224()
	io.Copy(hash, toMeasure)
	currentHash := fmt.Sprintf("%x", hash.Sum(nil))
	if !slice.ContainsString(skipHashCheckSettings, setting.Name) && currentHash == setting.Annotations[util.AnnotationHash] {
		return nil, nil
	}

	toUpdate := setting.DeepCopy()
	if toUpdate.Annotations == nil {
		toUpdate.Annotations = make(map[string]string)
	}

	var err error
	if syncer, ok := syncers[setting.Name]; ok {
		err = syncer(setting)
		if err == nil {
			toUpdate.Annotations[util.AnnotationHash] = currentHash
		}
		if updateErr := h.setConfiguredCondition(toUpdate, err); updateErr != nil {
			return setting, updateErr
		}
	}

	return setting, err
}

func (h *Handler) setConfiguredCondition(settingCopy *cloudweavv1.Setting, err error) error {
	if err != nil && (!cloudweavv1.SettingConfigured.IsFalse(settingCopy) ||
		cloudweavv1.SettingConfigured.GetMessage(settingCopy) != err.Error()) {
		cloudweavv1.SettingConfigured.False(settingCopy)
		cloudweavv1.SettingConfigured.Message(settingCopy, err.Error())
		if _, err := h.settings.Update(settingCopy); err != nil {
			return err
		}
	} else if err == nil {
		if settingCopy.Value == "" {
			cloudweavv1.SettingConfigured.False(settingCopy)
		} else {
			cloudweavv1.SettingConfigured.True(settingCopy)
		}
		cloudweavv1.SettingConfigured.Message(settingCopy, "")
		if _, err := h.settings.Update(settingCopy); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) updateBackupSecret(data map[string]string) error {
	secret, err := h.secretCache.Get(util.LonghornSystemNamespaceName, util.BackupTargetSecretName)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	toUpdate := secret.DeepCopy()
	if toUpdate.Data == nil {
		toUpdate.Data = make(map[string][]byte)
	}
	for key, value := range data {
		toUpdate.Data[key] = []byte(value)
	}
	_, err = h.secrets.Update(toUpdate)
	return err
}

func (h *Handler) redeployDeployment(namespace, name string) error {
	deployment, err := h.deploymentCache.Get(namespace, name)
	if err != nil {
		return err
	}
	toUpdate := deployment.DeepCopy()
	if deployment.Spec.Template.Annotations == nil {
		toUpdate.Spec.Template.Annotations = make(map[string]string)
	}
	toUpdate.Spec.Template.Annotations[util.AnnotationTimestamp] = time.Now().Format(time.RFC3339)

	_, err = h.deployments.Update(toUpdate)
	return err
}

func (h *Handler) appOnChanged(_ string, app *catalogv1api.App) (*catalogv1api.App, error) {
	if app == nil || app.DeletionTimestamp != nil {
		return nil, nil
	}

	cloudweavManagedChart, err := h.managedChartCache.Get(ManagedChartNamespace, CloudweavManagedChartName)
	if err != nil {
		return nil, err
	}

	if app.Namespace != cloudweavManagedChart.Spec.DefaultNamespace || app.Name != cloudweavManagedChart.Spec.ReleaseName {
		return nil, nil
	}

	return nil, UpdateSupportBundleImage(h.settings, h.settingCache, app)
}
