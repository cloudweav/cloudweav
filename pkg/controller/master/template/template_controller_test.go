package template

import (
	"context"
	"testing"

	"github.com/rancher/wrangler/v3/pkg/generic"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"

	cloudweavv1 "github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/fake"
	typeharv1 "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/typed/cloudweavhci.io/v1beta1"
)

func TestTemplateHandler_OnChanged(t *testing.T) {
	type input struct {
		key             string
		template        *cloudweavv1.VirtualMachineTemplate
		templateVersion *cloudweavv1.VirtualMachineTemplateVersion
	}
	type output struct {
		template *cloudweavv1.VirtualMachineTemplate
		err      error
	}

	var testCases = []struct {
		name     string
		given    input
		expected output
	}{
		{
			name: "nil resource",
			given: input{
				key:      "",
				template: nil,
			},
			expected: output{
				template: nil,
				err:      nil,
			},
		},
		{
			name: "deleted resource",
			given: input{
				key: "default/test",
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "test",
						DeletionTimestamp: &metav1.Time{},
					},
				},
			},
			expected: output{
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "test",
						DeletionTimestamp: &metav1.Time{},
					},
				},
				err: nil,
			},
		},
		{
			name: "blank default version ID",
			given: input{
				key: "default/test",
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "",
					},
				},
			},
			expected: output{
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "",
					},
				},
				err: nil,
			},
		},
		{
			name: "not corresponding version template",
			given: input{
				key: "default/test",
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "default/test",
					},
				},
				templateVersion: &cloudweavv1.VirtualMachineTemplateVersion{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "fake",
					},
					Spec: cloudweavv1.VirtualMachineTemplateVersionSpec{
						Description: "fake_description",
						TemplateID:  "default/test",
						ImageID:     "fake_image_id",
						VM:          cloudweavv1.VirtualMachineSourceSpec{},
					},
					Status: cloudweavv1.VirtualMachineTemplateVersionStatus{
						Version: 1,
						Conditions: []cloudweavv1.Condition{
							{
								Type:   cloudweavv1.VersionAssigned,
								Status: v1.ConditionTrue,
							},
						},
					},
				},
			},
			expected: output{
				template: nil,
				err:      errors.NewNotFound(schema.GroupResource{Group: "cloudweavhci.io", Resource: "virtualmachinetemplateversions"}, "test"),
			},
		},
		{
			name: "directly return as the template version is the same",
			given: input{
				key: "default/test",
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "default/test",
					},
					Status: cloudweavv1.VirtualMachineTemplateStatus{
						DefaultVersion: 1,
					},
				},
				templateVersion: &cloudweavv1.VirtualMachineTemplateVersion{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateVersionSpec{
						Description: "fake_description",
						TemplateID:  "fake_template_id",
						ImageID:     "fake_image_id",
						VM:          cloudweavv1.VirtualMachineSourceSpec{},
					},
					Status: cloudweavv1.VirtualMachineTemplateVersionStatus{
						Version: 1,
						Conditions: []cloudweavv1.Condition{
							{
								Type:   cloudweavv1.VersionAssigned,
								Status: v1.ConditionTrue,
							},
						},
					},
				},
			},
			expected: output{
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "default/test",
					},
					Status: cloudweavv1.VirtualMachineTemplateStatus{
						DefaultVersion: 1,
					},
				},
				err: nil,
			},
		},
		{
			name: "update template version",
			given: input{
				key: "default/test",
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "default/test",
					},
					Status: cloudweavv1.VirtualMachineTemplateStatus{
						DefaultVersion: 1,
						LatestVersion:  1,
					},
				},
				templateVersion: &cloudweavv1.VirtualMachineTemplateVersion{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateVersionSpec{
						Description: "fake_description",
						TemplateID:  "default/test",
						ImageID:     "fake_image_id",
						VM:          cloudweavv1.VirtualMachineSourceSpec{},
					},
					Status: cloudweavv1.VirtualMachineTemplateVersionStatus{
						Version: 2,
						Conditions: []cloudweavv1.Condition{
							{
								Type:   cloudweavv1.VersionAssigned,
								Status: v1.ConditionTrue,
							},
						},
					},
				},
			},
			expected: output{
				template: &cloudweavv1.VirtualMachineTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "test",
					},
					Spec: cloudweavv1.VirtualMachineTemplateSpec{
						DefaultVersionID: "default/test",
					},
					Status: cloudweavv1.VirtualMachineTemplateStatus{
						DefaultVersion: 2,
						LatestVersion:  2,
					},
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		var clientset = fake.NewSimpleClientset()
		if tc.given.template != nil {
			var err = clientset.Tracker().Add(tc.given.template)
			assert.Nil(t, err, "mock resource should add into fake controller tracker")
		}
		if tc.given.templateVersion != nil {
			var err = clientset.Tracker().Add(tc.given.templateVersion)
			assert.Nil(t, err, "mock resource should add into fake controller tracker")
		}

		var handler = &templateHandler{
			templates:            fakeTemplateClient(clientset.CloudweavhciV1beta1().VirtualMachineTemplates),
			templateVersions:     fakeTemplateVersionClient(clientset.CloudweavhciV1beta1().VirtualMachineTemplateVersions),
			templateVersionCache: fakeTemplateVersionCache(clientset.CloudweavhciV1beta1().VirtualMachineTemplateVersions),
		}
		var actual output
		actual.template, actual.err = handler.OnChanged(tc.given.key, tc.given.template)
		assert.Equal(t, tc.expected, actual, "case %q", tc.name)
	}
}

type fakeTemplateClient func(string) typeharv1.VirtualMachineTemplateInterface

func (c fakeTemplateClient) Create(template *cloudweavv1.VirtualMachineTemplate) (*cloudweavv1.VirtualMachineTemplate, error) {
	return c(template.Namespace).Create(context.TODO(), template, metav1.CreateOptions{})
}

func (c fakeTemplateClient) Update(template *cloudweavv1.VirtualMachineTemplate) (*cloudweavv1.VirtualMachineTemplate, error) {
	return c(template.Namespace).Update(context.TODO(), template, metav1.UpdateOptions{})
}

func (c fakeTemplateClient) UpdateStatus(template *cloudweavv1.VirtualMachineTemplate) (*cloudweavv1.VirtualMachineTemplate, error) {
	return c(template.Namespace).UpdateStatus(context.TODO(), template, metav1.UpdateOptions{})
}

func (c fakeTemplateClient) Delete(namespace, name string, opts *metav1.DeleteOptions) error {
	return c(namespace).Delete(context.TODO(), name, *opts)
}

func (c fakeTemplateClient) Get(namespace, name string, opts metav1.GetOptions) (*cloudweavv1.VirtualMachineTemplate, error) {
	return c(namespace).Get(context.TODO(), name, opts)
}

func (c fakeTemplateClient) List(namespace string, opts metav1.ListOptions) (*cloudweavv1.VirtualMachineTemplateList, error) {
	return c(namespace).List(context.TODO(), opts)
}

func (c fakeTemplateClient) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c(namespace).Watch(context.TODO(), opts)
}

func (c fakeTemplateClient) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *cloudweavv1.VirtualMachineTemplate, err error) {
	return c(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

func (c fakeTemplateClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1.VirtualMachineTemplate, *cloudweavv1.VirtualMachineTemplateList], error) {
	panic("implement me")
}

type fakeTemplateVersionCache func(string) typeharv1.VirtualMachineTemplateVersionInterface

func (c fakeTemplateVersionCache) Get(namespace, name string) (*cloudweavv1.VirtualMachineTemplateVersion, error) {
	return c(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c fakeTemplateVersionCache) List(_ string, _ labels.Selector) ([]*cloudweavv1.VirtualMachineTemplateVersion, error) {
	panic("implement me")
}

func (c fakeTemplateVersionCache) AddIndexer(_ string, _ generic.Indexer[*cloudweavv1.VirtualMachineTemplateVersion]) {
	panic("implement me")
}

func (c fakeTemplateVersionCache) GetByIndex(_, _ string) ([]*cloudweavv1.VirtualMachineTemplateVersion, error) {
	panic("implement me")
}

type fakeTemplateVersionClient func(string) typeharv1.VirtualMachineTemplateVersionInterface

func (c fakeTemplateVersionClient) Create(templateVersion *cloudweavv1.VirtualMachineTemplateVersion) (*cloudweavv1.VirtualMachineTemplateVersion, error) {
	return c(templateVersion.Namespace).Create(context.TODO(), templateVersion, metav1.CreateOptions{})
}

func (c fakeTemplateVersionClient) UpdateStatus(templateVersion *cloudweavv1.VirtualMachineTemplateVersion) (*cloudweavv1.VirtualMachineTemplateVersion, error) {
	return c(templateVersion.Namespace).UpdateStatus(context.TODO(), templateVersion, metav1.UpdateOptions{})
}

func (c fakeTemplateVersionClient) Update(templateVersion *cloudweavv1.VirtualMachineTemplateVersion) (*cloudweavv1.VirtualMachineTemplateVersion, error) {
	return c(templateVersion.Namespace).Update(context.TODO(), templateVersion, metav1.UpdateOptions{})
}

func (c fakeTemplateVersionClient) Delete(namespace, name string, opts *metav1.DeleteOptions) error {
	return c(namespace).Delete(context.TODO(), name, *opts)
}

func (c fakeTemplateVersionClient) Get(namespace, name string, opts metav1.GetOptions) (*cloudweavv1.VirtualMachineTemplateVersion, error) {
	return c(namespace).Get(context.TODO(), name, opts)
}

func (c fakeTemplateVersionClient) List(namespace string, opts metav1.ListOptions) (*cloudweavv1.VirtualMachineTemplateVersionList, error) {
	return c(namespace).List(context.TODO(), opts)
}

func (c fakeTemplateVersionClient) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c(namespace).Watch(context.TODO(), opts)
}

func (c fakeTemplateVersionClient) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *cloudweavv1.VirtualMachineTemplateVersion, err error) {
	return c(namespace).Patch(context.TODO(), name, pt, data, metav1.PatchOptions{}, subresources...)
}

func (c fakeTemplateVersionClient) WithImpersonation(_ rest.ImpersonationConfig) (generic.ClientInterface[*cloudweavv1.VirtualMachineTemplateVersion, *cloudweavv1.VirtualMachineTemplateVersionList], error) {
	panic("implement me")
}
