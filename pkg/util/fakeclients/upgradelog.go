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

type UpgradeLogClient func(string) harv1type.UpgradeLogInterface

func (c UpgradeLogClient) Create(upgradeLog *cloudweavv1.UpgradeLog) (*cloudweavv1.UpgradeLog, error) {
	return c(upgradeLog.Namespace).Create(context.TODO(), upgradeLog, metav1.CreateOptions{})
}
func (c UpgradeLogClient) Update(upgradeLog *cloudweavv1.UpgradeLog) (*cloudweavv1.UpgradeLog, error) {
	return c(upgradeLog.Namespace).Update(context.TODO(), upgradeLog, metav1.UpdateOptions{})
}
func (c UpgradeLogClient) UpdateStatus(*cloudweavv1.UpgradeLog) (*cloudweavv1.UpgradeLog, error) {
	panic("implement me")
}
func (c UpgradeLogClient) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c(namespace).Delete(context.TODO(), name, *options)
}
func (c UpgradeLogClient) Get(namespace, name string, options metav1.GetOptions) (*cloudweavv1.UpgradeLog, error) {
	return c(namespace).Get(context.TODO(), name, options)
}
func (c UpgradeLogClient) List(_ string, _ metav1.ListOptions) (*cloudweavv1.UpgradeLogList, error) {
	panic("implement me")
}
func (c UpgradeLogClient) Watch(_ string, _ metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}
func (c UpgradeLogClient) Patch(_, _ string, _ types.PatchType, _ []byte, _ ...string) (result *cloudweavv1.UpgradeLog, err error) {
	panic("implement me")
}

func (c UpgradeLogClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1.UpgradeLog, *cloudweavv1.UpgradeLogList], error) {
	panic("implement me")
}

type UpgradeLogCache func(string) harv1type.UpgradeLogInterface

func (c UpgradeLogCache) Get(namespace, name string) (*cloudweavv1.UpgradeLog, error) {
	return c(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
func (c UpgradeLogCache) List(_ string, _ labels.Selector) ([]*cloudweavv1.UpgradeLog, error) {
	panic("implement me")
}
func (c UpgradeLogCache) AddIndexer(_ string, _ generic.Indexer[*cloudweavv1.UpgradeLog]) {
	panic("implement me")
}
func (c UpgradeLogCache) GetByIndex(_, _ string) ([]*cloudweavv1.UpgradeLog, error) {
	panic("implement me")
}
