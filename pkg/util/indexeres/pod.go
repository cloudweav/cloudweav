package indexeres

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/cloudweav/cloudweav/pkg/ref"
	"github.com/cloudweav/cloudweav/pkg/util"
)

const (
	PodByVMNameIndex = "cloudweavhci.io/pod-by-vmname"
)

func PodByVMName(obj *corev1.Pod) ([]string, error) {
	vmName, ok := obj.Labels[util.LabelVMName]
	if !ok {
		return nil, nil
	}
	return []string{ref.Construct(obj.Namespace, vmName)}, nil
}
