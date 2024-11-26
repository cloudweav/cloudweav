package runtime

import (
	"context"
	"fmt"

	restclient "k8s.io/client-go/rest"

	"github.com/cloudweav/cloudweav/tests/framework/client"
	"github.com/cloudweav/cloudweav/tests/framework/env"
	"github.com/cloudweav/cloudweav/tests/framework/helm"
	"github.com/cloudweav/cloudweav/tests/framework/ready"
)

// Construct prepares runtime if "SKIP_CLOUDWEAV_INSTALLATION" is not "true".
func Construct(ctx context.Context, kubeConfig *restclient.Config) error {
	if env.IsSkipCloudweavInstallation() {
		return nil
	}

	// create namespaces
	var err error
	namespaces := []string{testCloudweavNamespace, testLonghornNamespace, testCattleNamespace}
	for _, namespace := range namespaces {
		err = client.CreateNamespace(kubeConfig, namespace)
		if err != nil {
			return fmt.Errorf("failed to create target namespace %s, %v", namespace, err)
		}
	}

	err = createCRDs(ctx, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create CRDs, %v", err)
	}

	// install cloudweav chart
	err = installCloudweavChart(ctx, kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to install cloudweav chart, %w", err)
	}

	return nil
}

// installCloudweavChart installs the basic components of cloudweav.
func installCloudweavChart(ctx context.Context, kubeConfig *restclient.Config) error {
	// chart values patches
	patches := map[string]interface{}{
		"replicas":                             0,
		"cloudweav-network-controller.enabled": true,
	}

	// webhook
	patches["webhook.controllerUser"] = "kubernetes-admin"
	patches["webhook.image.imagePullPolicy"] = "Never"
	repo, tag := env.GetWebhookImage()
	if repo != "" {
		patches["webhook.image.repository"] = repo
		patches["webhook.image.tag"] = tag
		patches["webhook.debug"] = true
	}

	if !env.IsE2ETestsEnabled() {
		patches["longhorn.enabled"] = "false"
	}

	if env.IsUsingEmulation() {
		patches["kubevirt.spec.configuration.developerConfiguration.useEmulation"] = "true"
	}

	// install crd chart
	_, err := helm.InstallChart(testCRDChartReleaseName, testCloudweavNamespace, testCRDChartDir, nil)
	if err != nil {
		return fmt.Errorf("failed to install cloudweav-crd chart: %w", err)
	}

	// install chart
	_, err = helm.InstallChart(testChartReleaseName, testCloudweavNamespace, testChartDir, patches)
	if err != nil {
		return fmt.Errorf("failed to install cloudweav chart: %w", err)
	}

	// verifies chart installation
	cloudweavReadyCondition, err := ready.NewNamespaceCondition(kubeConfig, testCloudweavNamespace)
	if err != nil {
		return fmt.Errorf("faield to create namespace ready condition from kubernetes config: %w", err)
	}
	cloudweavReadyCondition.AddDeploymentsReady(testDeploymentManifest...)
	cloudweavReadyCondition.AddDaemonSetsReady(testDaemonSetManifest...)

	if env.IsE2ETestsEnabled() {
		longhornReadyCondition, err := ready.NewNamespaceCondition(kubeConfig, testLonghornNamespace)
		if err != nil {
			return fmt.Errorf("faield to create namespace ready condition from kubernetes config: %w", err)
		}
		longhornReadyCondition.AddDeploymentsReady(longhornDeploymentManifest...)
		longhornReadyCondition.AddDaemonSetsReady(longhornDaemonSetManifest...)

		if err := longhornReadyCondition.Wait(ctx); err != nil {
			return err
		}
	}

	if err := cloudweavReadyCondition.Wait(ctx); err != nil {
		return err
	}

	return nil
}
