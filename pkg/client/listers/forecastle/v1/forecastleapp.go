/*
Copyright The Kubernetes Authors.

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

package v1

import (
	v1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ForecastleAppLister helps list ForecastleApps.
type ForecastleAppLister interface {
	// List lists all ForecastleApps in the indexer.
	List(selector labels.Selector) (ret []*v1.ForecastleApp, err error)
	// ForecastleApps returns an object that can list and get ForecastleApps.
	ForecastleApps(namespace string) ForecastleAppNamespaceLister
	ForecastleAppListerExpansion
}

// forecastleAppLister implements the ForecastleAppLister interface.
type forecastleAppLister struct {
	indexer cache.Indexer
}

// NewForecastleAppLister returns a new ForecastleAppLister.
func NewForecastleAppLister(indexer cache.Indexer) ForecastleAppLister {
	return &forecastleAppLister{indexer: indexer}
}

// List lists all ForecastleApps in the indexer.
func (s *forecastleAppLister) List(selector labels.Selector) (ret []*v1.ForecastleApp, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ForecastleApp))
	})
	return ret, err
}

// ForecastleApps returns an object that can list and get ForecastleApps.
func (s *forecastleAppLister) ForecastleApps(namespace string) ForecastleAppNamespaceLister {
	return forecastleAppNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ForecastleAppNamespaceLister helps list and get ForecastleApps.
type ForecastleAppNamespaceLister interface {
	// List lists all ForecastleApps in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ForecastleApp, err error)
	// Get retrieves the ForecastleApp from the indexer for a given namespace and name.
	Get(name string) (*v1.ForecastleApp, error)
	ForecastleAppNamespaceListerExpansion
}

// forecastleAppNamespaceLister implements the ForecastleAppNamespaceLister
// interface.
type forecastleAppNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ForecastleApps in the indexer for a given namespace.
func (s forecastleAppNamespaceLister) List(selector labels.Selector) (ret []*v1.ForecastleApp, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ForecastleApp))
	})
	return ret, err
}

// Get retrieves the ForecastleApp from the indexer for a given namespace and name.
func (s forecastleAppNamespaceLister) Get(name string) (*v1.ForecastleApp, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("forecastleapp"), name)
	}
	return obj.(*v1.ForecastleApp), nil
}