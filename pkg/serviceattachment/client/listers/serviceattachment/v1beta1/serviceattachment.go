/*
Copyright 2022 The Kubernetes Authors.

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
	v1beta1 "k8s.io/ingress-gce/pkg/apis/serviceattachment/v1beta1"
)

// ServiceAttachmentLister helps list ServiceAttachments.
// All objects returned here must be treated as read-only.
type ServiceAttachmentLister interface {
	// List lists all ServiceAttachments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ServiceAttachment, err error)
	// ServiceAttachments returns an object that can list and get ServiceAttachments.
	ServiceAttachments(namespace string) ServiceAttachmentNamespaceLister
	ServiceAttachmentListerExpansion
}

// serviceAttachmentLister implements the ServiceAttachmentLister interface.
type serviceAttachmentLister struct {
	indexer cache.Indexer
}

// NewServiceAttachmentLister returns a new ServiceAttachmentLister.
func NewServiceAttachmentLister(indexer cache.Indexer) ServiceAttachmentLister {
	return &serviceAttachmentLister{indexer: indexer}
}

// List lists all ServiceAttachments in the indexer.
func (s *serviceAttachmentLister) List(selector labels.Selector) (ret []*v1beta1.ServiceAttachment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ServiceAttachment))
	})
	return ret, err
}

// ServiceAttachments returns an object that can list and get ServiceAttachments.
func (s *serviceAttachmentLister) ServiceAttachments(namespace string) ServiceAttachmentNamespaceLister {
	return serviceAttachmentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServiceAttachmentNamespaceLister helps list and get ServiceAttachments.
// All objects returned here must be treated as read-only.
type ServiceAttachmentNamespaceLister interface {
	// List lists all ServiceAttachments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ServiceAttachment, err error)
	// Get retrieves the ServiceAttachment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.ServiceAttachment, error)
	ServiceAttachmentNamespaceListerExpansion
}

// serviceAttachmentNamespaceLister implements the ServiceAttachmentNamespaceLister
// interface.
type serviceAttachmentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ServiceAttachments in the indexer for a given namespace.
func (s serviceAttachmentNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.ServiceAttachment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ServiceAttachment))
	})
	return ret, err
}

// Get retrieves the ServiceAttachment from the indexer for a given namespace and name.
func (s serviceAttachmentNamespaceLister) Get(name string) (*v1beta1.ServiceAttachment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("serviceattachment"), name)
	}
	return obj.(*v1beta1.ServiceAttachment), nil
}
