package keypair

import (
	"fmt"

	"golang.org/x/crypto/ssh"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	ctlcloudweavv1 "github.com/cloudweav/cloudweav/pkg/generated/controllers/cloudweavhci.io/v1beta1"
)

// Handler computes key pairs' fingerprints
type Handler struct {
	keyPairClient ctlcloudweavv1.KeyPairClient
}

func (h *Handler) OnKeyPairChanged(_ string, keyPair *cloudweavv1.KeyPair) (*cloudweavv1.KeyPair, error) {
	if keyPair == nil || keyPair.DeletionTimestamp != nil {
		return keyPair, nil
	}

	if keyPair.Spec.PublicKey == "" || keyPair.Status.FingerPrint != "" {
		return keyPair, nil
	}

	toUpdate := keyPair.DeepCopy()
	publicKey := []byte(keyPair.Spec.PublicKey)
	pk, _, _, _, err := ssh.ParseAuthorizedKey(publicKey)
	if err != nil {
		cloudweavv1.KeyPairValidated.False(toUpdate)
		cloudweavv1.KeyPairValidated.Reason(toUpdate, fmt.Sprintf("failed to parse the public key, error: %v", err))
	} else {
		fingerPrint := ssh.FingerprintLegacyMD5(pk)
		toUpdate.Status.FingerPrint = fingerPrint
		cloudweavv1.KeyPairValidated.True(toUpdate)
	}
	return h.keyPairClient.Update(toUpdate)
}
