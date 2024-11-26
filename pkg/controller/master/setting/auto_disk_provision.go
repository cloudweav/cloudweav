package setting

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
)

const (
	ManagedChartNamespace     = "fleet-local"
	CloudweavManagedChartName = "cloudweav"
	NDMName                   = "cloudweav-node-disk-manager"
)

func (h *Handler) syncNDMAutoProvisionPaths(setting *cloudweavv1.Setting) error {
	mChart, err := h.managedChartCache.Get(ManagedChartNamespace, CloudweavManagedChartName)
	if err != nil {
		return err
	}
	mChartCopy := mChart.DeepCopy()

	NDMValues, ok := mChartCopy.Spec.Values.Data[NDMName]
	if !ok {
		return fmt.Errorf("NDM chart value not found in ManagedChart")
	}

	NDMValuesMap, ok := NDMValues.(map[string]interface{})
	if !ok {
		return fmt.Errorf("NDM chart value is not a map[string]interface{}")
	}

	autoProvFilters := strings.Split(setting.Value, ",")
	for i, filter := range autoProvFilters {
		autoProvFilters[i] = strings.TrimSpace(filter)
	}

	NDMValuesMap["autoProvisionFilter"] = autoProvFilters
	mChartCopy.Spec.Values.Data[NDMName] = NDMValuesMap

	logrus.Debugf("NDM values to be updated to ManagedChart: %v", mChartCopy.Spec.Values.Data[NDMName])
	if _, err := h.managedCharts.Update(mChartCopy); err != nil {
		return err
	}

	return nil
}
