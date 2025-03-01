package setting

import (
	"encoding/json"
	"reflect"

	"github.com/mitchellh/mapstructure"
	catalogv1api "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"

	"github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/settings"
)

const (
	SupportBundleRepository = "support-bundle-kit"
)

func UpdateSupportBundleImage(settingClient v1beta1.SettingClient, settingCache v1beta1.SettingCache, app *catalogv1api.App) error {
	if app.Spec.Chart == nil {
		logrus.WithFields(logrus.Fields{
			"namespace": app.Namespace,
			"name":      app.Name,
		}).Warning("app has empty chart, skip to update the support-bundle-image setting")
		return nil
	}
	// merge template and chart
	values, err := chartutil.CoalesceValues(
		&chart.Chart{
			// In latest version of chartutil, it read chart.metadata.name, so we put a value here.
			// ref: https://github.com/helm/helm/blob/b8d3535991dd5089d58bc88c46a5ffe2721ae830/pkg/chartutil/coalesce.go#L160
			Metadata: &chart.Metadata{Name: "merge-templates-and-values"},
			Values:   app.Spec.Chart.Values,
		},
		app.Spec.Values,
	)
	if err != nil {
		return err
	}

	var supportBundleYaml map[string]interface{}
	v, ok := values[SupportBundleRepository]
	if !ok {
		logrus.WithFields(logrus.Fields{
			"namespace": app.Namespace,
			"name":      app.Name,
		}).Warningf("app chart values cant't find %s, skip to update the support-bundle-image setting", SupportBundleRepository)
		return nil
	}
	if supportBundleYaml, ok = v.(map[string]interface{}); !ok {
		logrus.WithFields(logrus.Fields{
			"namespace": app.Namespace,
			"name":      app.Name,
		}).Warningf("unknown %s yaml struct %+v, skip to update the support-bundle-image setting", SupportBundleRepository, v)
		return nil
	}
	if len(supportBundleYaml) == 0 {
		logrus.WithFields(logrus.Fields{
			"namespace": app.Namespace,
			"name":      app.Name,
		}).Warning("supportBundleYaml map is empty, skip to convert to support-bundle-image setting")
		return nil
	}

	var supportBundle struct {
		Image settings.Image
	}
	if err := mapstructure.Decode(supportBundleYaml, &supportBundle); err != nil {
		return err
	}
	if supportBundle.Image.ImageName() == "" {
		logrus.WithFields(logrus.Fields{
			"namespace": app.Namespace,
			"name":      app.Name,
		}).Warning("the converted imagename is empty, skip to update the support-bundle-image setting")
		return nil
	}
	imageBytes, err := json.Marshal(&supportBundle.Image)
	if err != nil {
		return err
	}
	supportBundleImage, err := settingCache.Get(settings.SupportBundleImageName)
	if err != nil {
		return err
	}

	supportBundleImageCpy := supportBundleImage.DeepCopy()
	supportBundleImageCpy.Default = string(imageBytes)

	if reflect.DeepEqual(supportBundleImage, supportBundleImageCpy) {
		return nil
	}
	_, err = settingClient.Update(supportBundleImageCpy)
	return err
}
