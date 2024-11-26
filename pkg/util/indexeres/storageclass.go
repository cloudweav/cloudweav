package indexeres

import (
	storagev1 "k8s.io/api/storage/v1"

	"github.com/cloudweav/cloudweav/pkg/ref"
	"github.com/cloudweav/cloudweav/pkg/util"
)

const (
	StorageClassBySecretIndex = "cloudweavhci.io/storageclass-by-secret"
)

func StorageClassBySecret(sc *storagev1.StorageClass) ([]string, error) {
	encrypted, ok := sc.Parameters[util.LonghornOptionEncrypted]
	if !ok || encrypted != "true" {
		return []string{}, nil
	}

	secretName := sc.Parameters[util.CSINodePublishSecretNameKey]
	secretNS := sc.Parameters[util.CSINodePublishSecretNamespaceKey]

	return []string{
		ref.Construct(secretNS, secretName),
	}, nil
}
