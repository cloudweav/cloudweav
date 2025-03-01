package fakeclients

import (
	"context"
	"fmt"

	mgmtv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/v3/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	mgmtv3type "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/management.cattle.io/v3"
	"github.com/cloudweav/cloudweav/tests/framework/fuzz"
)

type ManagedChartClient func(string) mgmtv3type.ManagedChartInterface

func (c ManagedChartClient) Update(managedChart *mgmtv3.ManagedChart) (*mgmtv3.ManagedChart, error) {
	return c(managedChart.Namespace).Update(context.TODO(), managedChart, metav1.UpdateOptions{})
}
func (c ManagedChartClient) Get(namespace, name string, options metav1.GetOptions) (*mgmtv3.ManagedChart, error) {
	return c(namespace).Get(context.TODO(), name, options)
}
func (c ManagedChartClient) Create(managedChart *mgmtv3.ManagedChart) (*mgmtv3.ManagedChart, error) {
	if managedChart.GenerateName != "" {
		managedChart.Name = fmt.Sprintf("%s%s", managedChart.GenerateName, fuzz.String(5))
	}
	return c(managedChart.Namespace).Create(context.TODO(), managedChart, metav1.CreateOptions{})
}
func (c ManagedChartClient) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c(namespace).Delete(context.TODO(), name, *options)
}
func (c ManagedChartClient) List(_ string, _ metav1.ListOptions) (*mgmtv3.ManagedChartList, error) {
	panic("implement me")
}
func (c ManagedChartClient) UpdateStatus(*mgmtv3.ManagedChart) (*mgmtv3.ManagedChart, error) {
	panic("implement me")
}
func (c ManagedChartClient) Watch(_ string, _ metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}
func (c ManagedChartClient) Patch(_, _ string, _ types.PatchType, _ []byte, _ ...string) (result *mgmtv3.ManagedChart, err error) {
	panic("implement me")
}
func (c ManagedChartClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*mgmtv3.ManagedChart, *mgmtv3.ManagedChartList], error) {
	panic("implement me")
}

type ManagedChartCache func(string) mgmtv3type.ManagedChartInterface

func (c ManagedChartCache) Get(namespace, name string) (*mgmtv3.ManagedChart, error) {
	return c(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
func (c ManagedChartCache) List(namespace string, selector labels.Selector) ([]*mgmtv3.ManagedChart, error) {
	list, err := c(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	result := make([]*mgmtv3.ManagedChart, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, &list.Items[i])
	}
	return result, err
}
func (c ManagedChartCache) AddIndexer(_ string, _ generic.Indexer[*mgmtv3.ManagedChart]) {
	panic("implement me")
}
func (c ManagedChartCache) GetByIndex(_, _ string) ([]*mgmtv3.ManagedChart, error) {
	panic("implement me")
}
