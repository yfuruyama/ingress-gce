/*
Copyright 2019 The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1beta1 "k8s.io/ingress-gce/pkg/apis/frontendconfig/v1beta1"
	scheme "k8s.io/ingress-gce/pkg/frontendconfig/client/clientset/versioned/scheme"
)

// FrontendConfigsGetter has a method to return a FrontendConfigInterface.
// A group's client should implement this interface.
type FrontendConfigsGetter interface {
	FrontendConfigs(namespace string) FrontendConfigInterface
}

// FrontendConfigInterface has methods to work with FrontendConfig resources.
type FrontendConfigInterface interface {
	Create(*v1beta1.FrontendConfig) (*v1beta1.FrontendConfig, error)
	Update(*v1beta1.FrontendConfig) (*v1beta1.FrontendConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.FrontendConfig, error)
	List(opts v1.ListOptions) (*v1beta1.FrontendConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.FrontendConfig, err error)
	FrontendConfigExpansion
}

// frontendConfigs implements FrontendConfigInterface
type frontendConfigs struct {
	client rest.Interface
	ns     string
}

// newFrontendConfigs returns a FrontendConfigs
func newFrontendConfigs(c *NetworkingV1beta1Client, namespace string) *frontendConfigs {
	return &frontendConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the frontendConfig, and returns the corresponding frontendConfig object, and an error if there is any.
func (c *frontendConfigs) Get(name string, options v1.GetOptions) (result *v1beta1.FrontendConfig, err error) {
	result = &v1beta1.FrontendConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("frontendconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FrontendConfigs that match those selectors.
func (c *frontendConfigs) List(opts v1.ListOptions) (result *v1beta1.FrontendConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.FrontendConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("frontendconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested frontendConfigs.
func (c *frontendConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("frontendconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a frontendConfig and creates it.  Returns the server's representation of the frontendConfig, and an error, if there is any.
func (c *frontendConfigs) Create(frontendConfig *v1beta1.FrontendConfig) (result *v1beta1.FrontendConfig, err error) {
	result = &v1beta1.FrontendConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("frontendconfigs").
		Body(frontendConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a frontendConfig and updates it. Returns the server's representation of the frontendConfig, and an error, if there is any.
func (c *frontendConfigs) Update(frontendConfig *v1beta1.FrontendConfig) (result *v1beta1.FrontendConfig, err error) {
	result = &v1beta1.FrontendConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("frontendconfigs").
		Name(frontendConfig.Name).
		Body(frontendConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the frontendConfig and deletes it. Returns an error if one occurs.
func (c *frontendConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("frontendconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *frontendConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("frontendconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched frontendConfig.
func (c *frontendConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.FrontendConfig, err error) {
	result = &v1beta1.FrontendConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("frontendconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
