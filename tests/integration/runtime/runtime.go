package runtime

import (
	"fmt"
	"os"

	"github.com/cloudweav/cloudweav/pkg/config"
	"github.com/cloudweav/cloudweav/pkg/settings"
	"github.com/cloudweav/cloudweav/tests/framework/fuzz"
)

const (
	testChartDir            = "../../../deploy/charts/cloudweav"
	testCRDChartDir         = "../../../deploy/charts/cloudweav-crd"
	testCloudweavNamespace  = "cloudweav-system"
	testLonghornNamespace   = "longhorn-system"
	testCattleNamespace     = "cattle-system"
	testChartReleaseName    = "cloudweav"
	testCRDChartReleaseName = "cloudweav-crd"
)

var (
	testDeploymentManifest = []string{
		"virt-operator",
		"virt-api",
		"virt-controller",
	}
	testDaemonSetManifest = []string{
		"virt-handler",
	}
	longhornDeploymentManifest = []string{
		"csi-attacher",
		"csi-snapshotter",
		"csi-provisioner",
		"csi-resizer",
		"longhorn-driver-deployer",
	}
	longhornDaemonSetManifest = []string{
		"longhorn-manager",
		"engine-image-ei-2938e020",
		"longhorn-csi-plugin",
	}
)

// SetConfig configures the public variables exported in github.com/cloudweav/cloudweav/pkg/config package.
func SetConfig() (config.Options, error) {
	var options config.Options

	// generate two random ports
	ports, err := fuzz.FreePorts(2)
	if err != nil {
		return options, fmt.Errorf("failed to get listening ports of cloudweav server, %v", err)
	}

	// config http and https
	options.HTTPListenPort = ports[0]
	options.HTTPSListenPort = ports[1]
	options.Namespace = testCloudweavNamespace
	options.RancherEmbedded = false

	// inject the preset envs, this is used for testing setting.
	err = os.Setenv(settings.GetEnvKey(settings.APIUIVersion.Name), settings.APIUIVersion.Default)
	if err != nil {
		return options, fmt.Errorf("failed to preset ENVs of cloudweav server, %w", err)
	}
	return options, nil
}
