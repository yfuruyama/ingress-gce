/*
Copyright 2020 The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1beta1 "k8s.io/ingress-gce/pkg/apis/svcneg/v1beta1"
)

// ServiceNetworkEndpointGroupLister helps list ServiceNetworkEndpointGroups.
// All objects returned here must be treated as read-only.
type ServiceNetworkEndpointGroupLister interface {
	// List lists all ServiceNetworkEndpointGroups in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ServiceNetworkEndpointGroup, err error)
	// ServiceNetworkEndpointGroups returns an object that can list and get ServiceNetworkEndpointGroups.
	ServiceNetworkEndpointGroups(namespace string) ServiceNetworkEndpointGroupNamespaceLister
	ServiceNetworkEndpointGroupListerExpansion
}

// serviceNetworkEndpointGroupLister implements the ServiceNetworkEndpointGroupLister interface.
type serviceNetworkEndpointGroupLister struct {
	indexer cache.Indexer
}

// NewServiceNetworkEndpointGroupLister returns a new ServiceNetworkEndpointGroupLister.
func NewServiceNetworkEndpointGroupLister(indexer cache.Indexer) ServiceNetworkEndpointGroupLister {
	return &serviceNetworkEndpointGroupLister{indexer: indexer}
}

// List lists all ServiceNetworkEndpointGroups in the indexer.
func (s *serviceNetworkEndpointGroupLister) List(selector labels.Selector) (ret []*v1beta1.ServiceNetworkEndpointGroup, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ServiceNetworkEndpointGroup))
	})
	return ret, err
}

// ServiceNetworkEndpointGroups returns an object that can list and get ServiceNetworkEndpointGroups.
func (s *serviceNetworkEndpointGroupLister) ServiceNetworkEndpointGroups(namespace string) ServiceNetworkEndpointGroupNamespaceLister {
	return serviceNetworkEndpointGroupNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServiceNetworkEndpointGroupNamespaceLister helps list and get ServiceNetworkEndpointGroups.
// All objects returned here must be treated as read-only.
type ServiceNetworkEndpointGroupNamespaceLister interface {
	// List lists all ServiceNetworkEndpointGroups in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ServiceNetworkEndpointGroup, err error)
	// Get retrieves the ServiceNetworkEndpointGroup from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.ServiceNetworkEndpointGroup, error)
	ServiceNetworkEndpointGroupNamespaceListerExpansion
}

// serviceNetworkEndpointGroupNamespaceLister implements the ServiceNetworkEndpointGroupNamespaceLister
// interface.
type serviceNetworkEndpointGroupNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ServiceNetworkEndpointGroups in the indexer for a given namespace.
func (s serviceNetworkEndpointGroupNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.ServiceNetworkEndpointGroup, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ServiceNetworkEndpointGroup))
	})
	return ret, err
}

// Get retrieves the ServiceNetworkEndpointGroup from the indexer for a given namespace and name.
func (s serviceNetworkEndpointGroupNamespaceLister) Get(name string) (*v1beta1.ServiceNetworkEndpointGroup, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("servicenetworkendpointgroup"), name)
	}
	return obj.(*v1beta1.ServiceNetworkEndpointGroup), nil
}