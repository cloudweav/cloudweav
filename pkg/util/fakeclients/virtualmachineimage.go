package fakeclients

import (
	"context"
	"fmt"
	"time"

	"github.com/rancher/wrangler/v3/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	harv1type "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/util"
	"github.com/cloudweav/cloudweav/pkg/webhook/indexeres"
	"github.com/cloudweav/cloudweav/tests/framework/fuzz"
)

type VirtualMachineImageClient func(string) harv1type.VirtualMachineImageInterface

func (c VirtualMachineImageClient) Informer() cache.SharedIndexInformer {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) GroupVersionKind() schema.GroupVersionKind {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) AddGenericHandler(_ context.Context, _ string, _ generic.Handler) {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) AddGenericRemoveHandler(_ context.Context, _ string, _ generic.Handler) {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) Updater() generic.Updater {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) OnChange(_ context.Context, _ string, _ generic.ObjectHandler[*cloudweavv1.VirtualMachineImage]) {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) OnRemove(_ context.Context, _ string, _ generic.ObjectHandler[*cloudweavv1.VirtualMachineImage]) {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) Enqueue(_, _ string) {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) EnqueueAfter(_, _ string, _ time.Duration) {
	// do nothing
}

func (c VirtualMachineImageClient) Cache() generic.CacheInterface[*cloudweavv1.VirtualMachineImage] {
	//TODO implement me
	panic("implement me")
}

func (c VirtualMachineImageClient) Update(virtualMachineImage *cloudweavv1.VirtualMachineImage) (*cloudweavv1.VirtualMachineImage, error) {
	return c(virtualMachineImage.Namespace).Update(context.TODO(), virtualMachineImage, metav1.UpdateOptions{})
}
func (c VirtualMachineImageClient) Get(namespace, name string, options metav1.GetOptions) (*cloudweavv1.VirtualMachineImage, error) {
	return c(namespace).Get(context.TODO(), name, options)
}
func (c VirtualMachineImageClient) Create(virtualMachineImage *cloudweavv1.VirtualMachineImage) (*cloudweavv1.VirtualMachineImage, error) {
	if virtualMachineImage.GenerateName != "" {
		virtualMachineImage.Name = fmt.Sprintf("%s%s", virtualMachineImage.GenerateName, fuzz.String(5))
	}
	return c(virtualMachineImage.Namespace).Create(context.TODO(), virtualMachineImage, metav1.CreateOptions{})
}
func (c VirtualMachineImageClient) Delete(_, _ string, _ *metav1.DeleteOptions) error {
	panic("implement me")
}
func (c VirtualMachineImageClient) List(_ string, _ metav1.ListOptions) (*cloudweavv1.VirtualMachineImageList, error) {
	panic("implement me")
}
func (c VirtualMachineImageClient) UpdateStatus(*cloudweavv1.VirtualMachineImage) (*cloudweavv1.VirtualMachineImage, error) {
	panic("implement me")
}
func (c VirtualMachineImageClient) Watch(_ string, _ metav1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}
func (c VirtualMachineImageClient) Patch(_, _ string, _ types.PatchType, _ []byte, _ ...string) (result *cloudweavv1.VirtualMachineImage, err error) {
	panic("implement me")
}
func (c VirtualMachineImageClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1.VirtualMachineImage, *cloudweavv1.VirtualMachineImageList], error) {
	panic("implement me")
}

type VirtualMachineImageCache func(string) harv1type.VirtualMachineImageInterface

func (c VirtualMachineImageCache) Get(namespace, name string) (*cloudweavv1.VirtualMachineImage, error) {
	return c(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
func (c VirtualMachineImageCache) List(namespace string, selector labels.Selector) ([]*cloudweavv1.VirtualMachineImage, error) {
	list, err := c(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	result := make([]*cloudweavv1.VirtualMachineImage, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, &list.Items[i])
	}
	return result, err
}
func (c VirtualMachineImageCache) AddIndexer(_ string, _ generic.Indexer[*cloudweavv1.VirtualMachineImage]) {
	panic("implement me")
}
func (c VirtualMachineImageCache) GetByIndex(key, scName string) ([]*cloudweavv1.VirtualMachineImage, error) {
	var vmimages []*cloudweavv1.VirtualMachineImage

	// TODO:
	// Need to figure out how to better test this.
	// Otherwise, we should add more testing NS here.
	testingNS := []string{"default"}

	switch key {
	case indexeres.ImageByStorageClass:
		for _, ns := range testingNS {
			vmList, err := c(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			for _, vm := range vmList.Items {
				vm := vm
				sc, ok := vm.Annotations[util.AnnotationStorageClassName]
				if !ok {
					continue
				}
				if sc == scName {
					vmimages = append(vmimages, &vm)
				}
			}
		}
	default:
		panic(fmt.Sprintf("unimplemented indexer: %s", key))
	}

	return vmimages, nil
}
