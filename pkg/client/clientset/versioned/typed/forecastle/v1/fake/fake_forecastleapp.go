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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	forecastlev1 "github.com/stakater/Forecastle/pkg/apis/forecastle/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeForecastleApps implements ForecastleAppInterface
type FakeForecastleApps struct {
	Fake *FakeForecastleV1
	ns   string
}

var forecastleappsResource = schema.GroupVersionResource{Group: "forecastle.stakater.com", Version: "v1", Resource: "forecastleapps"}

var forecastleappsKind = schema.GroupVersionKind{Group: "forecastle.stakater.com", Version: "v1", Kind: "ForecastleApp"}

// Get takes name of the forecastleApp, and returns the corresponding forecastleApp object, and an error if there is any.
func (c *FakeForecastleApps) Get(name string, options v1.GetOptions) (result *forecastlev1.ForecastleApp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(forecastleappsResource, c.ns, name), &forecastlev1.ForecastleApp{})

	if obj == nil {
		return nil, err
	}
	return obj.(*forecastlev1.ForecastleApp), err
}

// List takes label and field selectors, and returns the list of ForecastleApps that match those selectors.
func (c *FakeForecastleApps) List(opts v1.ListOptions) (result *forecastlev1.ForecastleAppList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(forecastleappsResource, forecastleappsKind, c.ns, opts), &forecastlev1.ForecastleAppList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &forecastlev1.ForecastleAppList{ListMeta: obj.(*forecastlev1.ForecastleAppList).ListMeta}
	for _, item := range obj.(*forecastlev1.ForecastleAppList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested forecastleApps.
func (c *FakeForecastleApps) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(forecastleappsResource, c.ns, opts))

}

// Create takes the representation of a forecastleApp and creates it.  Returns the server's representation of the forecastleApp, and an error, if there is any.
func (c *FakeForecastleApps) Create(forecastleApp *forecastlev1.ForecastleApp) (result *forecastlev1.ForecastleApp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(forecastleappsResource, c.ns, forecastleApp), &forecastlev1.ForecastleApp{})

	if obj == nil {
		return nil, err
	}
	return obj.(*forecastlev1.ForecastleApp), err
}

// Update takes the representation of a forecastleApp and updates it. Returns the server's representation of the forecastleApp, and an error, if there is any.
func (c *FakeForecastleApps) Update(forecastleApp *forecastlev1.ForecastleApp) (result *forecastlev1.ForecastleApp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(forecastleappsResource, c.ns, forecastleApp), &forecastlev1.ForecastleApp{})

	if obj == nil {
		return nil, err
	}
	return obj.(*forecastlev1.ForecastleApp), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeForecastleApps) UpdateStatus(forecastleApp *forecastlev1.ForecastleApp) (*forecastlev1.ForecastleApp, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(forecastleappsResource, "status", c.ns, forecastleApp), &forecastlev1.ForecastleApp{})

	if obj == nil {
		return nil, err
	}
	return obj.(*forecastlev1.ForecastleApp), err
}

// Delete takes name of the forecastleApp and deletes it. Returns an error if one occurs.
func (c *FakeForecastleApps) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(forecastleappsResource, c.ns, name), &forecastlev1.ForecastleApp{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeForecastleApps) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(forecastleappsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &forecastlev1.ForecastleAppList{})
	return err
}

// Patch applies the patch and returns the patched forecastleApp.
func (c *FakeForecastleApps) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *forecastlev1.ForecastleApp, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(forecastleappsResource, c.ns, name, pt, data, subresources...), &forecastlev1.ForecastleApp{})

	if obj == nil {
		return nil, err
	}
	return obj.(*forecastlev1.ForecastleApp), err
}