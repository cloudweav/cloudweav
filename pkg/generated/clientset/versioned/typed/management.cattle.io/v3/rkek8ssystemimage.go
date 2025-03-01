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

// RkeK8sSystemImagesGetter has a method to return a RkeK8sSystemImageInterface.
// A group's client should implement this interface.
type RkeK8sSystemImagesGetter interface {
	RkeK8sSystemImages(namespace string) RkeK8sSystemImageInterface
}

// RkeK8sSystemImageInterface has methods to work with RkeK8sSystemImage resources.
type RkeK8sSystemImageInterface interface {
	Create(ctx context.Context, rkeK8sSystemImage *v3.RkeK8sSystemImage, opts v1.CreateOptions) (*v3.RkeK8sSystemImage, error)
	Update(ctx context.Context, rkeK8sSystemImage *v3.RkeK8sSystemImage, opts v1.UpdateOptions) (*v3.RkeK8sSystemImage, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v3.RkeK8sSystemImage, error)
	List(ctx context.Context, opts v1.ListOptions) (*v3.RkeK8sSystemImageList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.RkeK8sSystemImage, err error)
	RkeK8sSystemImageExpansion
}

// rkeK8sSystemImages implements RkeK8sSystemImageInterface
type rkeK8sSystemImages struct {
	client rest.Interface
	ns     string
}

// newRkeK8sSystemImages returns a RkeK8sSystemImages
func newRkeK8sSystemImages(c *ManagementV3Client, namespace string) *rkeK8sSystemImages {
	return &rkeK8sSystemImages{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the rkeK8sSystemImage, and returns the corresponding rkeK8sSystemImage object, and an error if there is any.
func (c *rkeK8sSystemImages) Get(ctx context.Context, name string, options v1.GetOptions) (result *v3.RkeK8sSystemImage, err error) {
	result = &v3.RkeK8sSystemImage{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RkeK8sSystemImages that match those selectors.
func (c *rkeK8sSystemImages) List(ctx context.Context, opts v1.ListOptions) (result *v3.RkeK8sSystemImageList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v3.RkeK8sSystemImageList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested rkeK8sSystemImages.
func (c *rkeK8sSystemImages) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a rkeK8sSystemImage and creates it.  Returns the server's representation of the rkeK8sSystemImage, and an error, if there is any.
func (c *rkeK8sSystemImages) Create(ctx context.Context, rkeK8sSystemImage *v3.RkeK8sSystemImage, opts v1.CreateOptions) (result *v3.RkeK8sSystemImage, err error) {
	result = &v3.RkeK8sSystemImage{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(rkeK8sSystemImage).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a rkeK8sSystemImage and updates it. Returns the server's representation of the rkeK8sSystemImage, and an error, if there is any.
func (c *rkeK8sSystemImages) Update(ctx context.Context, rkeK8sSystemImage *v3.RkeK8sSystemImage, opts v1.UpdateOptions) (result *v3.RkeK8sSystemImage, err error) {
	result = &v3.RkeK8sSystemImage{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		Name(rkeK8sSystemImage.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(rkeK8sSystemImage).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the rkeK8sSystemImage and deletes it. Returns an error if one occurs.
func (c *rkeK8sSystemImages) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rkeK8sSystemImages) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched rkeK8sSystemImage.
func (c *rkeK8sSystemImages) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v3.RkeK8sSystemImage, err error) {
	result = &v3.RkeK8sSystemImage{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("rkek8ssystemimages").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
