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

// CatalogsGetter has a method to return a CatalogInterface.
// A group's client should implement this interface.
type CatalogsGetter interface {
	Catalogs() CatalogInterface
}

// CatalogInterface has methods to work with Catalog resources.
type CatalogInterface interface {
	Create(ctx context.Context, catalog *v3.Catalog, opts v1.CreateOptions) (*v3.Catalog, error)
	Update(ctx context.Context, catalog *v3.Catalog, opts v1.UpdateOptions) (*v3.Catalog, error)
	UpdateStatus(ctx context.Context, catalog *v3.Catalog, opts v1.UpdateOptions) (*v3.Catalog, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v3.Catalog, error)
	List(ctx context.Context, opts v1.ListOptions) (*v3.CatalogList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.Catalog, err error)
	CatalogExpansion
}

// catalogs implements CatalogInterface
type catalogs struct {
	client rest.Interface
}

// newCatalogs returns a Catalogs
func newCatalogs(c *ManagementV3Client) *catalogs {
	return &catalogs{
		client: c.RESTClient(),
	}
}

// Get takes name of the catalog, and returns the corresponding catalog object, and an error if there is any.
func (c *catalogs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v3.Catalog, err error) {
	result = &v3.Catalog{}
	err = c.client.Get().
		Resource("catalogs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Catalogs that match those selectors.
func (c *catalogs) List(ctx context.Context, opts v1.ListOptions) (result *v3.CatalogList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v3.CatalogList{}
	err = c.client.Get().
		Resource("catalogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested catalogs.
func (c *catalogs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("catalogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a catalog and creates it.  Returns the server's representation of the catalog, and an error, if there is any.
func (c *catalogs) Create(ctx context.Context, catalog *v3.Catalog, opts v1.CreateOptions) (result *v3.Catalog, err error) {
	result = &v3.Catalog{}
	err = c.client.Post().
		Resource("catalogs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(catalog).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a catalog and updates it. Returns the server's representation of the catalog, and an error, if there is any.
func (c *catalogs) Update(ctx context.Context, catalog *v3.Catalog, opts v1.UpdateOptions) (result *v3.Catalog, err error) {
	result = &v3.Catalog{}
	err = c.client.Put().
		Resource("catalogs").
		Name(catalog.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(catalog).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *catalogs) UpdateStatus(ctx context.Context, catalog *v3.Catalog, opts v1.UpdateOptions) (result *v3.Catalog, err error) {
	result = &v3.Catalog{}
	err = c.client.Put().
		Resource("catalogs").
		Name(catalog.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(catalog).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the catalog and deletes it. Returns an error if one occurs.
func (c *catalogs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("catalogs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *catalogs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("catalogs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched catalog.
func (c *catalogs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.Catalog, err error) {
	result = &v3.Catalog{}
	err = c.client.Patch(pt).
		Resource("catalogs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
