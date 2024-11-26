package supportbundle

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corefake "k8s.io/client-go/kubernetes/fake"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/controller/master/supportbundle/types"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/fake"
	"github.com/cloudweav/cloudweav/pkg/util/fakeclients"
)

func Test_checkExistTime(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	coreclientset := corefake.NewSimpleClientset()
	namespace := "test-support-bundle"

	if _, err := coreclientset.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{}); err != nil {
		assert.Nil(t, err, "failed to create namespace", namespace)
	}

	handler := Handler{
		supportBundles:          fakeclients.SupportBundleClient(clientset.CloudweavhciV1beta1().SupportBundles),
		supportBundleController: fakeclients.SupportBundleClient(clientset.CloudweavhciV1beta1().SupportBundles),
	}

	tests := []struct {
		name             string
		getSupportBundle func() *cloudweavv1.SupportBundle
		expected         func(*cloudweavv1.SupportBundle, error, string)
	}{
		{
			name: "ready state",
			getSupportBundle: func() *cloudweavv1.SupportBundle {
				sb := &cloudweavv1.SupportBundle{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test1",
						Namespace: namespace,
					},
					Status: cloudweavv1.SupportBundleStatus{
						State: types.StateReady,
					},
				}
				cloudweavv1.SupportBundleInitialized.True(sb)
				cloudweavv1.SupportBundleInitialized.LastUpdated(sb, time.Now().Add(-35*time.Minute).Format(time.RFC3339))
				return sb
			},
			expected: func(sb *cloudweavv1.SupportBundle, err error, name string) {
				assert.Nil(t, sb, name)
				assert.True(t, apierrors.IsNotFound(err), name)
			},
		},
		{
			name: "error state",
			getSupportBundle: func() *cloudweavv1.SupportBundle {
				sb := &cloudweavv1.SupportBundle{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test2",
						Namespace: namespace,
					},
					Status: cloudweavv1.SupportBundleStatus{
						State: types.StateError,
					},
				}
				cloudweavv1.SupportBundleInitialized.False(sb)
				cloudweavv1.SupportBundleInitialized.Message(sb, "custom error")
				cloudweavv1.SupportBundleInitialized.LastUpdated(sb, time.Now().Add(-35*time.Minute).Format(time.RFC3339))
				return sb
			},
			expected: func(sb *cloudweavv1.SupportBundle, err error, name string) {
				assert.Nil(t, sb, name)
				assert.True(t, apierrors.IsNotFound(err), name)
			},
		},
		{
			name: "non-final state should not be deleted",
			getSupportBundle: func() *cloudweavv1.SupportBundle {
				sb := &cloudweavv1.SupportBundle{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test3",
						Namespace: namespace,
					},
					Status: cloudweavv1.SupportBundleStatus{
						State: "other state",
					},
				}
				return sb
			},
			expected: func(sb *cloudweavv1.SupportBundle, err error, name string) {
				assert.Equal(t, sb.Name, "test3", name)
				assert.Equal(t, sb.Status.State, "other state", name)
				assert.Nil(t, err, name)
			},
		},
	}

	for _, tc := range tests {
		sb := tc.getSupportBundle()
		_, err := clientset.CloudweavhciV1beta1().SupportBundles(namespace).Create(context.Background(), sb, metav1.CreateOptions{})
		assert.Nil(t, err, tc.name)

		_, err = handler.OnSupportBundleChanged("", sb)
		assert.Nil(t, err, tc.name)

		sb, err = clientset.CloudweavhciV1beta1().SupportBundles(namespace).Get(context.Background(), sb.Name, metav1.GetOptions{})
		tc.expected(sb, err, tc.name)
	}
}
