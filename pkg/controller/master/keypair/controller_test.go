package keypair

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	"github.com/rancher/wrangler/v3/pkg/generic"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/fake"
	typeharv1 "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/cloudweavhci.io/v1beta1"
)

func TestHandler_OnKeyPairChanged(t *testing.T) {
	type input struct {
		key     string
		keyPair *cloudweavv1.KeyPair
	}
	type output struct {
		keyPair *cloudweavv1.KeyPair
		err     error
	}

	var testPublicKey, testPublicKeyFingerprint, err = generateSSHPublicKey()
	assert.Nil(t, err, "mock SSH public key should be created")
	var testCases = []struct {
		name     string
		given    input
		expected output
	}{
		{
			name: "nil resource",
			given: input{
				key:     "",
				keyPair: nil,
			},
			expected: output{
				keyPair: nil,
				err:     nil,
			},
		},
		{
			name: "deleted resource",
			given: input{
				key: "default/test",
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "test",
						DeletionTimestamp: &metav1.Time{},
					},
				},
			},
			expected: output{
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "test",
						DeletionTimestamp: &metav1.Time{},
					},
				},
				err: nil,
			},
		},
		{
			name: "blank public key",
			given: input{
				key: "default/test",
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: "",
					},
				},
			},
			expected: output{
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: "",
					},
				},
				err: nil,
			},
		},
		{
			name: "not blank fingerprint",
			given: input{
				key: "default/test",
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: "FAKE_PUBLIC_KEY",
					},
					Status: cloudweavv1.KeyPairStatus{
						FingerPrint: "FAKE_FINGER_PRINT",
					},
				},
			},
			expected: output{
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: "FAKE_PUBLIC_KEY",
					},
					Status: cloudweavv1.KeyPairStatus{
						FingerPrint: "FAKE_FINGER_PRINT",
					},
				},
				err: nil,
			},
		},
		{
			name: "illegal public key",
			given: input{
				key: "default/test",
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: "FAKE_PUBLIC_KEY",
					},
				},
			},
			expected: output{
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: "FAKE_PUBLIC_KEY",
					},
					Status: cloudweavv1.KeyPairStatus{
						Conditions: []cloudweavv1.Condition{
							{
								Type:   cloudweavv1.KeyPairValidated,
								Status: corev1.ConditionFalse,
								Reason: "failed to parse the public key, error: ssh: no key found",
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "generate fingerprint for legal public key",
			given: input{
				key: "default/test",
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: testPublicKey,
					},
				},
			},
			expected: output{
				keyPair: &cloudweavv1.KeyPair{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.KeyPairSpec{
						PublicKey: testPublicKey,
					},
					Status: cloudweavv1.KeyPairStatus{
						Conditions: []cloudweavv1.Condition{
							{
								Type:   cloudweavv1.KeyPairValidated,
								Status: corev1.ConditionTrue,
							},
						},
						FingerPrint: testPublicKeyFingerprint,
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		var clientset = fake.NewSimpleClientset()
		if tc.given.keyPair != nil {
			var err = clientset.Tracker().Add(tc.given.keyPair)
			assert.Nil(t, err, "mock resource should add into fake controller tracker")
		}

		var handler = &Handler{
			keyPairClient: fakeKeyPairClient(clientset.CloudweavhciV1beta1().KeyPairs),
		}
		var actual output
		actual.keyPair, actual.err = handler.OnKeyPairChanged(tc.given.key, tc.given.keyPair)
		// NB(thxCode) we don't need to compare the `lastUpdateTime` and `lastTransitionTime` of conditions.
		if actual.keyPair != nil {
			for i := range actual.keyPair.Status.Conditions {
				actual.keyPair.Status.Conditions[i].LastUpdateTime = ""
				actual.keyPair.Status.Conditions[i].LastTransitionTime = ""
			}
		}

		assert.Equal(t, tc.expected, actual, "case %q", tc.name)
	}
}

func generateSSHPublicKey() (pk string, fingerprint string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key, %v", err)
	}
	pubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create SSH public key, %v", err)
	}
	pk = string(ssh.MarshalAuthorizedKey(pubKey))
	fingerprint = ssh.FingerprintLegacyMD5(pubKey)
	return pk, fingerprint, nil
}

type fakeKeyPairClient func(string) typeharv1.KeyPairInterface

func (c fakeKeyPairClient) Create(pair *cloudweavv1.KeyPair) (*cloudweavv1.KeyPair, error) {
	return c(pair.Namespace).Create(context.TODO(), pair, metav1.CreateOptions{})
}

func (c fakeKeyPairClient) Update(pair *cloudweavv1.KeyPair) (*cloudweavv1.KeyPair, error) {
	return c(pair.Namespace).Update(context.TODO(), pair, metav1.UpdateOptions{})
}

func (c fakeKeyPairClient) UpdateStatus(pair *cloudweavv1.KeyPair) (*cloudweavv1.KeyPair, error) {
	return c(pair.Namespace).UpdateStatus(context.TODO(), pair, metav1.UpdateOptions{})
}

func (c fakeKeyPairClient) Delete(namespace, name string, opts *metav1.DeleteOptions) error {
	return c(namespace).Delete(context.TODO(), name, *opts)
}

func (c fakeKeyPairClient) Get(namespace, name string, opts metav1.GetOptions) (*cloudweavv1.KeyPair, error) {
	return c(namespace).Get(context.TODO(), name, opts)
}

func (c fakeKeyPairClient) List(namespace string, opts metav1.ListOptions) (*cloudweavv1.KeyPairList, error) {
	return c(namespace).List(context.TODO(), opts)
}

func (c fakeKeyPairClient) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c(namespace).Watch(context.TODO(), opts)
}

func (c fakeKeyPairClient) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *cloudweavv1.KeyPair, err error) {
	return c(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

func (c fakeKeyPairClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1.KeyPair, *cloudweavv1.KeyPairList], error) {
	panic("implement me")
}
