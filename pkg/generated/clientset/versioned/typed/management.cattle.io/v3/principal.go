/*
Copyright 2024 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	scheme "github.com/cloudweav/cloudweav/pkg/generated/clientset/versioned/scheme"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PrincipalsGetter has a method to return a PrincipalInterface.
// A group's client should implement this interface.
type PrincipalsGetter interface {
	Principals() PrincipalInterface
}

// PrincipalInterface has methods to work with Principal resources.
type PrincipalInterface interface {
	Create(ctx context.Context, principal *v3.Principal, opts v1.CreateOptions) (*v3.Principal, error)
	Update(ctx context.Context, principal *v3.Principal, opts v1.UpdateOptions) (*v3.Principal, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v3.Principal, error)
	List(ctx context.Context, opts v1.ListOptions) (*v3.PrincipalList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.Principal, err error)
	PrincipalExpansion
}

// principals implements PrincipalInterface
type principals struct {
	client rest.Interface
}

// newPrincipals returns a Principals
func newPrincipals(c *ManagementV3Client) *principals {
	return &principals{
		client: c.RESTClient(),
	}
}

// Get takes name of the principal, and returns the corresponding principal object, and an error if there is any.
func (c *principals) Get(ctx context.Context, name string, options v1.GetOptions) (result *v3.Principal, err error) {
	result = &v3.Principal{}
	err = c.client.Get().
		Resource("principals").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Principals that match those selectors.
func (c *principals) List(ctx context.Context, opts v1.ListOptions) (result *v3.PrincipalList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v3.PrincipalList{}
	err = c.client.Get().
		Resource("principals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested principals.
func (c *principals) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("principals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a principal and creates it.  Returns the server's representation of the principal, and an error, if there is any.
func (c *principals) Create(ctx context.Context, principal *v3.Principal, opts v1.CreateOptions) (result *v3.Principal, err error) {
	result = &v3.Principal{}
	err = c.client.Post().
		Resource("principals").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(principal).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a principal and updates it. Returns the server's representation of the principal, and an error, if there is any.
func (c *principals) Update(ctx context.Context, principal *v3.Principal, opts v1.UpdateOptions) (result *v3.Principal, err error) {
	result = &v3.Principal{}
	err = c.client.Put().
		Resource("principals").
		Name(principal.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(principal).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the principal and deletes it. Returns an error if one occurs.
func (c *principals) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("principals").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *principals) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("principals").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched principal.
func (c *principals) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.Principal, err error) {
	result = &v3.Principal{}
	err = c.client.Patch(pt).
		Resource("principals").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
