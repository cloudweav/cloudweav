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

// AuthProvidersGetter has a method to return a AuthProviderInterface.
// A group's client should implement this interface.
type AuthProvidersGetter interface {
	AuthProviders() AuthProviderInterface
}

// AuthProviderInterface has methods to work with AuthProvider resources.
type AuthProviderInterface interface {
	Create(ctx context.Context, authProvider *v3.AuthProvider, opts v1.CreateOptions) (*v3.AuthProvider, error)
	Update(ctx context.Context, authProvider *v3.AuthProvider, opts v1.UpdateOptions) (*v3.AuthProvider, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v3.AuthProvider, error)
	List(ctx context.Context, opts v1.ListOptions) (*v3.AuthProviderList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.AuthProvider, err error)
	AuthProviderExpansion
}

// authProviders implements AuthProviderInterface
type authProviders struct {
	client rest.Interface
}

// newAuthProviders returns a AuthProviders
func newAuthProviders(c *ManagementV3Client) *authProviders {
	return &authProviders{
		client: c.RESTClient(),
	}
}

// Get takes name of the authProvider, and returns the corresponding authProvider object, and an error if there is any.
func (c *authProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v3.AuthProvider, err error) {
	result = &v3.AuthProvider{}
	err = c.client.Get().
		Resource("authproviders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AuthProviders that match those selectors.
func (c *authProviders) List(ctx context.Context, opts v1.ListOptions) (result *v3.AuthProviderList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v3.AuthProviderList{}
	err = c.client.Get().
		Resource("authproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested authProviders.
func (c *authProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("authproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a authProvider and creates it.  Returns the server's representation of the authProvider, and an error, if there is any.
func (c *authProviders) Create(ctx context.Context, authProvider *v3.AuthProvider, opts v1.CreateOptions) (result *v3.AuthProvider, err error) {
	result = &v3.AuthProvider{}
	err = c.client.Post().
		Resource("authproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(authProvider).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a authProvider and updates it. Returns the server's representation of the authProvider, and an error, if there is any.
func (c *authProviders) Update(ctx context.Context, authProvider *v3.AuthProvider, opts v1.UpdateOptions) (result *v3.AuthProvider, err error) {
	result = &v3.AuthProvider{}
	err = c.client.Put().
		Resource("authproviders").
		Name(authProvider.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(authProvider).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the authProvider and deletes it. Returns an error if one occurs.
func (c *authProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("authproviders").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *authProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("authproviders").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched authProvider.
func (c *authProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.AuthProvider, err error) {
	result = &v3.AuthProvider{}
	err = c.client.Patch(pt).
		Resource("authproviders").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
