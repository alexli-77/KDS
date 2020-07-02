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

// This file was automatically generated by informer-gen

package internalversion

import (
	internalinterfaces "k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// InitializerConfigurations returns a InitializerConfigurationInformer.
	InitializerConfigurations() InitializerConfigurationInformer
	// MutatingWebhookConfigurations returns a MutatingWebhookConfigurationInformer.
	MutatingWebhookConfigurations() MutatingWebhookConfigurationInformer
	// ValidatingWebhookConfigurations returns a ValidatingWebhookConfigurationInformer.
	ValidatingWebhookConfigurations() ValidatingWebhookConfigurationInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// InitializerConfigurations returns a InitializerConfigurationInformer.
func (v *version) InitializerConfigurations() InitializerConfigurationInformer {
	return &initializerConfigurationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// MutatingWebhookConfigurations returns a MutatingWebhookConfigurationInformer.
func (v *version) MutatingWebhookConfigurations() MutatingWebhookConfigurationInformer {
	return &mutatingWebhookConfigurationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// ValidatingWebhookConfigurations returns a ValidatingWebhookConfigurationInformer.
func (v *version) ValidatingWebhookConfigurations() ValidatingWebhookConfigurationInformer {
	return &validatingWebhookConfigurationInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
