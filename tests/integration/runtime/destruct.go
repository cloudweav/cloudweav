package runtime

import (
	"context"
	"fmt"

	restclient "k8s.io/client-go/rest"

	"github.com/cloudweav/cloudweav/tests/framework/env"
	"github.com/cloudweav/cloudweav/tests/framework/helm"
	"github.com/cloudweav/cloudweav/tests/framework/ready"
)

// Destruct releases the runtime if "SKIP_CLOUDWEAV_INSTALLATION" is not "true".
func Destruct(ctx context.Context, kubeConfig *restclient.Config) error {
	if env.IsKeepingCloudweavInstallation() || env.IsSkipCloudweavInstallation() {
		return nil
	}

	// uninstall cloudweav chart
	err := uninstallCloudweavCharts(ctx, kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

// uninstallCloudweavCharts uninstalls the basic components of cloudweav.
func uninstallCloudweavCharts(ctx context.Context, kubeConfig *restclient.Config) error {
	// uninstall chart
	_, err := helm.UninstallChart(testChartReleaseName, testCloudweavNamespace)
	if err != nil {
		return fmt.Errorf("failed to uninstall cloudweav chart: %v", err)
	}

	// verifies chart uninstallation
	namespaceReadyCondition, err := ready.NewNamespaceCondition(kubeConfig, testCloudweavNamespace)
	if err != nil {
		return fmt.Errorf("faield to create namespace ready condition from kubernetes config: %w", err)
	}
	namespaceReadyCondition.AddDeploymentsClean(testDeploymentManifest...)
	namespaceReadyCondition.AddDaemonSetsClean(testDaemonSetManifest...)

	return namespaceReadyCondition.Wait(ctx)
}
