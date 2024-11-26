package fakeclients

import (
	"context"

	"github.com/rancher/wrangler/v3/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	harv1type "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/cloudweavhci.io/v1beta1"
)

type UpgradeClient func(string) harv1type.UpgradeInterface

func (c UpgradeClient) Update(upgrade *cloudweavv1.Upgrade) (*cloudweavv1.Upgrade, error) {
	return c(upgrade.Namespace).Update(context.TODO(), upgrade, metav1.UpdateOptions{})
}
func (c UpgradeClient) Get(_, _ string, _ metav1.GetOptions) (*cloudweavv1.Upgrade, error) {
	panic("implement me")
}
func (c UpgradeClient) Create(*cloudweavv1.Upgrade) (*cloudweavv1.Upgrade, error) {
	panic("implement me")
}
func (c UpgradeClient) Delete(_, _ string, _ *metav1.DeleteOptions) error {
	panic("implement me")
}
func (c UpgradeClient) List(_ string, _ metav1.ListOptions) (*cloudweavv1.UpgradeList, error) {
	panic("implement me")
}
func (c UpgradeClient) UpdateStatus(*cloudweavv1.Upgrade) (*cloudweavv1.Upgrade, error) {
	panic("implement me")
}
func (c UpgradeClient) Watch(_ string, _ metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}
func (c UpgradeClient) Patch(_, _ string, _ types.PatchType, _ []byte, _ ...string) (result *cloudweavv1.Upgrade, err error) {
	panic("implement me")
}
func (c UpgradeClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1.Upgrade, *cloudweavv1.UpgradeList], error) {
	panic("implement me")
}

type UpgradeCache func(string) harv1type.UpgradeInterface

func (c UpgradeCache) Get(namespace, name string) (*cloudweavv1.Upgrade, error) {
	return c(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
func (c UpgradeCache) List(namespace string, selector labels.Selector) ([]*cloudweavv1.Upgrade, error) {
	list, err := c(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	result := make([]*cloudweavv1.Upgrade, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, &list.Items[i])
	}
	return result, err
}
func (c UpgradeCache) AddIndexer(_ string, _ generic.Indexer[*cloudweavv1.Upgrade]) {
	panic("implement me")
}
func (c UpgradeCache) GetByIndex(_, _ string) ([]*cloudweavv1.Upgrade, error) {
	panic("implement me")
}
