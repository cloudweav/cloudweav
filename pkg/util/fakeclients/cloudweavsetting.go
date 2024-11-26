package fakeclients

import (
	"context"

	"github.com/rancher/wrangler/v3/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	"github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	cloudweavtype "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/cloudweavhci.io/v1beta1"
)

type CloudweavSettingClient func() cloudweavtype.SettingInterface

func (c CloudweavSettingClient) Create(s *v1beta1.Setting) (*v1beta1.Setting, error) {
	return c().Create(context.TODO(), s, metav1.CreateOptions{})
}

func (c CloudweavSettingClient) Update(s *v1beta1.Setting) (*v1beta1.Setting, error) {
	return c().Update(context.TODO(), s, metav1.UpdateOptions{})
}

func (c CloudweavSettingClient) UpdateStatus(_ *v1beta1.Setting) (*v1beta1.Setting, error) {
	panic("implement me")
}

func (c CloudweavSettingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c().Delete(context.TODO(), name, *options)
}

func (c CloudweavSettingClient) Get(name string, options metav1.GetOptions) (*v1beta1.Setting, error) {
	return c().Get(context.TODO(), name, options)
}

func (c CloudweavSettingClient) List(opts metav1.ListOptions) (*v1beta1.SettingList, error) {
	return c().List(context.TODO(), opts)
}

func (c CloudweavSettingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c().Watch(context.TODO(), opts)
}

func (c CloudweavSettingClient) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Setting, err error) {
	return c().Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

func (c CloudweavSettingClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.NonNamespacedClientInterface[*v1beta1.Setting, *v1beta1.SettingList], error) {
	panic("implement me")
}

type CloudweavSettingCache func() cloudweavtype.SettingInterface

func (c CloudweavSettingCache) Get(name string) (*v1beta1.Setting, error) {
	return c().Get(context.TODO(), name, metav1.GetOptions{})
}

func (c CloudweavSettingCache) List(_ labels.Selector) ([]*v1beta1.Setting, error) {
	panic("implement me")
}

func (c CloudweavSettingCache) AddIndexer(_ string, _ generic.Indexer[*v1beta1.Setting]) {
	panic("implement me")
}

func (c CloudweavSettingCache) GetByIndex(_, _ string) ([]*v1beta1.Setting, error) {
	panic("implement me")
}
