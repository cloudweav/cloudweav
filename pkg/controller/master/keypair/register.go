package keypair

import (
	"context"

	"github.com/cloudweav/cloudweav/pkg/config"
)

const (
	controllerAgentName = "vm-keypair-controller"
)

func Register(ctx context.Context, management *config.Management, _ config.Options) error {
	keyPairs := management.CloudweavFactory.Cloudweavhci().V1beta1().KeyPair()
	controller := &Handler{
		keyPairClient: keyPairs,
	}

	keyPairs.OnChange(ctx, controllerAgentName, controller.OnKeyPairChanged)
	return nil
}
