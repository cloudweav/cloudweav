package fakeclients

import (
	"context"

	"github.com/rancher/wrangler/v3/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	cloudweavv1beta1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	cloudweavtype "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/indexeres"
	"github.com/cloudweav/cloudweav/pkg/ref"
)

type VMBackupClient func(string) cloudweavtype.VirtualMachineBackupInterface

func (c VMBackupClient) Create(vmBackup *cloudweavv1beta1.VirtualMachineBackup) (*cloudweavv1beta1.VirtualMachineBackup, error) {
	return c(vmBackup.Namespace).Create(context.TODO(), vmBackup, metav1.CreateOptions{})
}

func (c VMBackupClient) Update(volume *cloudweavv1beta1.VirtualMachineBackup) (*cloudweavv1beta1.VirtualMachineBackup, error) {
	return c(volume.Namespace).Update(context.TODO(), volume, metav1.UpdateOptions{})
}

func (c VMBackupClient) UpdateStatus(_ *cloudweavv1beta1.VirtualMachineBackup) (*cloudweavv1beta1.VirtualMachineBackup, error) {
	panic("implement me")
}

func (c VMBackupClient) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c(namespace).Delete(context.TODO(), name, *options)
}

func (c VMBackupClient) Get(namespace, name string, options metav1.GetOptions) (*cloudweavv1beta1.VirtualMachineBackup, error) {
	return c(namespace).Get(context.TODO(), name, options)
}

func (c VMBackupClient) List(namespace string, opts metav1.ListOptions) (*cloudweavv1beta1.VirtualMachineBackupList, error) {
	return c(namespace).List(context.TODO(), opts)
}

func (c VMBackupClient) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c(namespace).Watch(context.TODO(), opts)
}

func (c VMBackupClient) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *cloudweavv1beta1.VirtualMachineBackup, err error) {
	return c(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

func (c VMBackupClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1beta1.VirtualMachineBackup, *cloudweavv1beta1.VirtualMachineBackupList], error) {
	panic("implement me")
}

type VMBackupCache func(string) cloudweavtype.VirtualMachineBackupInterface

func (c VMBackupCache) Get(namespace, name string) (*cloudweavv1beta1.VirtualMachineBackup, error) {
	return c(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c VMBackupCache) List(namespace string, selector labels.Selector) ([]*cloudweavv1beta1.VirtualMachineBackup, error) {
	list, err := c(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	result := make([]*cloudweavv1beta1.VirtualMachineBackup, 0, len(list.Items))
	for i := range list.Items {
		result = append(result, &list.Items[i])
	}
	return result, err
}

func (c VMBackupCache) AddIndexer(_ string, _ generic.Indexer[*cloudweavv1beta1.VirtualMachineBackup]) {
	panic("implement me")
}

func (c VMBackupCache) GetByIndex(indexName, key string) ([]*cloudweavv1beta1.VirtualMachineBackup, error) {
	switch indexName {
	case indexeres.VMBackupBySourceVMNameIndex:
		vmNamespace, _ := ref.Parse(key)
		backupList, err := c(vmNamespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		var backups []*cloudweavv1beta1.VirtualMachineBackup
		for i := range backupList.Items {
			b := backupList.Items[i]
			if b.Spec.Source.Name == key {
				backups = append(backups, &b)
			}
		}
		return backups, nil
	default:
		return nil, nil
	}
}
