package data

import (
	"github.com/rancher/wrangler/v3/pkg/apply"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	publicNamespace = "cloudweav-public"
)

func addPublicNamespace(apply apply.Apply) error {
	// All authenticated users are readable in the public namespace
	return apply.
		WithDynamicLookup().
		WithSetID("cloudweav-public").
		ApplyObjects(
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: publicNamespace},
			})
}
