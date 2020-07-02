/*
Copyright 2017 The Kubernetes Authors.

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

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	core "k8s.io/kubernetes/pkg/apis/core"
)

// FakeServices implements ServiceInterface
type FakeServices struct {
	Fake *FakeCore
	ns   string
}

var servicesResource = schema.GroupVersionResource{Group: "", Version: "", Resource: "services"}

var servicesKind = schema.GroupVersionKind{Group: "", Version: "", Kind: "Service"}

// Get takes name of the service, and returns the corresponding service object, and an error if there is any.
func (c *FakeServices) Get(name string, options v1.GetOptions) (result *core.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(servicesResource, c.ns, name), &core.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*core.Service), err
}

// List takes label and field selectors, and returns the list of Services that match those selectors.
func (c *FakeServices) List(opts v1.ListOptions) (result *core.ServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(servicesResource, servicesKind, c.ns, opts), &core.ServiceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &core.ServiceList{}
	for _, item := range obj.(*core.ServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested services.
func (c *FakeServices) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(servicesResource, c.ns, opts))

}

// Create takes the representation of a service and creates it.  Returns the server's representation of the service, and an error, if there is any.
func (c *FakeServices) Create(service *core.Service) (result *core.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(servicesResource, c.ns, service), &core.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*core.Service), err
}

// Update takes the representation of a service and updates it. Returns the server's representation of the service, and an error, if there is any.
func (c *FakeServices) Update(service *core.Service) (result *core.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(servicesResource, c.ns, service), &core.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*core.Service), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeServices) UpdateStatus(service *core.Service) (*core.Service, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(servicesResource, "status", c.ns, service), &core.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*core.Service), err
}

// Delete takes name of the service and deletes it. Returns an error if one occurs.
func (c *FakeServices) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(servicesResource, c.ns, name), &core.Service{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServices) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(servicesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &core.ServiceList{})
	return err
}

// Patch applies the patch and returns the patched service.
func (c *FakeServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *core.Service, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(servicesResource, c.ns, name, data, subresources...), &core.Service{})

	if obj == nil {
		return nil, err
	}
	return obj.(*core.Service), err
}
