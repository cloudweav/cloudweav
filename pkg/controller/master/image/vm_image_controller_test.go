package image

import (
	"net/http"
	"testing"
	"time"

	longhorntypes "github.com/longhorn/longhorn-manager/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/fake"
	"github.com/cloudweav/cloudweav/pkg/util"
	"github.com/cloudweav/cloudweav/pkg/util/fakeclients"
)

func TestVMImageHandler_OnChanged(t *testing.T) {
	type input struct {
		image   *cloudweavv1.VirtualMachineImage
		objects []runtime.Object
	}
	var testCases = []struct {
		name     string
		given    input
		expected func(t *testing.T, handler *vmImageHandler, image *cloudweavv1.VirtualMachineImage, err error)
	}{
		{
			name: "Test case 1: Create Encrypted Image",
			given: input{
				objects: []runtime.Object{
					&cloudweavv1.VirtualMachineImage{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "source-image",
							Namespace: "default",
						},
						Spec: cloudweavv1.VirtualMachineImageSpec{
							SourceType:  "download",
							URL:         "https://dl-cdn.alpinelinux.org/alpine/v3.20/releases/x86_64/alpine-standard-3.20.2-x86_64.iso",
							DisplayName: "source-image",
						},
					},
				},
				image: &cloudweavv1.VirtualMachineImage{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "encrypted-test-image",
						Namespace: "default",
					},
					Spec: cloudweavv1.VirtualMachineImageSpec{
						SourceType: "clone",
						SecurityParameters: &cloudweavv1.VirtualMachineImageSecurityParameters{
							CryptoOperation:      "encrypt",
							SourceImageName:      "source-image",
							SourceImageNamespace: "default",
						},
						// After mutator, these parameters are from storage class of annotation
						// But, we didn't have mutator here, so we just put them here
						StorageClassParameters: map[string]string{
							util.LonghornOptionEncrypted:          "true",
							util.CSIProvisionerSecretNameKey:      "test-secret",
							util.CSIProvisionerSecretNamespaceKey: "default",
							util.CSINodeStageSecretNameKey:        "test-secret",
							util.CSINodeStageSecretNamespaceKey:   "default",
							util.CSINodePublishSecretNameKey:      "test-secret",
							util.CSINodePublishSecretNamespaceKey: "default",
							longhorntypes.OptionNumberOfReplicas:  "1",
						},
						DisplayName: "encrypted-test-image",
					},
				},
			},
			expected: func(t *testing.T, handler *vmImageHandler, _ *cloudweavv1.VirtualMachineImage, err error) {
				bis, _ := handler.backingImageCache.List("longhorn-system", labels.Everything())
				assert.Equal(t, 1, len(bis))
				assert.Equal(t, "default/encrypted-test-image", bis[0].Annotations[util.AnnotationImageID])
				assert.Equal(t, "test-secret", bis[0].Spec.SourceParameters["secret"])
				assert.Equal(t, "default", bis[0].Spec.SourceParameters["secret-namespace"])
				assert.Nil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset(append(tc.given.objects, tc.given.image)...)
			k8sclientset := k8sfake.NewSimpleClientset()

			handler := &vmImageHandler{
				backingImages:     fakeclients.BackingImageClient(clientset.LonghornV1beta2().BackingImages),
				backingImageCache: fakeclients.BackingImageCache(clientset.LonghornV1beta2().BackingImages),
				storageClasses:    fakeclients.StorageClassClient(k8sclientset.StorageV1().StorageClasses),
				storageClassCache: fakeclients.StorageClassCache(k8sclientset.StorageV1().StorageClasses),
				images:            fakeclients.VirtualMachineImageClient(clientset.CloudweavhciV1beta1().VirtualMachineImages),
				imageController:   fakeclients.VirtualMachineImageClient(clientset.CloudweavhciV1beta1().VirtualMachineImages),
				httpClient: http.Client{
					Timeout: 15 * time.Second,
				},
				pvcCache: fakeclients.PersistentVolumeClaimCache(k8sclientset.CoreV1().PersistentVolumeClaims),
			}

			image, err := handler.OnChanged("", tc.given.image)

			tc.expected(t, handler, image, err)
		})
	}
}
