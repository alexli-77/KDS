/*
Copyright 2016 The Kubernetes Authors.

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

package priorities

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/kubernetes/plugin/pkg/scheduler/algorithm"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
)

type PriorityMetadataFactory struct {
	serviceLister     algorithm.ServiceLister
	controllerLister  algorithm.ControllerLister
	replicaSetLister  algorithm.ReplicaSetLister
	statefulSetLister algorithm.StatefulSetLister
}

func NewPriorityMetadataFactory(serviceLister algorithm.ServiceLister, controllerLister algorithm.ControllerLister, replicaSetLister algorithm.ReplicaSetLister, statefulSetLister algorithm.StatefulSetLister) algorithm.MetadataProducer {
	factory := &PriorityMetadataFactory{
		serviceLister:     serviceLister,
		controllerLister:  controllerLister,
		replicaSetLister:  replicaSetLister,
		statefulSetLister: statefulSetLister,
	}
	return factory.PriorityMetadata
}

// priorityMetadata is a type that is passed as metadata for priority functions
type priorityMetadata struct {
	nonZeroRequest *schedulercache.Resource
	podTolerations []v1.Toleration
	affinity       *v1.Affinity
	podSelectors   []labels.Selector
}

// PriorityMetadata is a MetadataProducer.  Node info can be nil.
func (pmf *PriorityMetadataFactory) PriorityMetadata(pod *v1.Pod, nodeNameToInfo map[string]*schedulercache.NodeInfo) interface{} {
	// If we cannot compute metadata, just return nil
	if pod == nil {
		return nil
	}
	tolerationsPreferNoSchedule := getAllTolerationPreferNoSchedule(pod.Spec.Tolerations)
	podSelectors := getSelectors(pod, pmf.serviceLister, pmf.controllerLister, pmf.replicaSetLister, pmf.statefulSetLister)
	return &priorityMetadata{
		nonZeroRequest: getNonZeroRequests(pod),
		podTolerations: tolerationsPreferNoSchedule,
		affinity:       pod.Spec.Affinity,
		podSelectors:   podSelectors,
	}
}

// getSelectors returns selectors of services, RCs and RSs matching the given pod.
func getSelectors(pod *v1.Pod, sl algorithm.ServiceLister, cl algorithm.ControllerLister, rsl algorithm.ReplicaSetLister, ssl algorithm.StatefulSetLister) []labels.Selector {
	var selectors []labels.Selector
	if services, err := sl.GetPodServices(pod); err == nil {
		for _, service := range services {
			selectors = append(selectors, labels.SelectorFromSet(service.Spec.Selector))
		}
	}
	if rcs, err := cl.GetPodControllers(pod); err == nil {
		for _, rc := range rcs {
			selectors = append(selectors, labels.SelectorFromSet(rc.Spec.Selector))
		}
	}
	if rss, err := rsl.GetPodReplicaSets(pod); err == nil {
		for _, rs := range rss {
			if selector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector); err == nil {
				selectors = append(selectors, selector)
			}
		}
	}
	if sss, err := ssl.GetPodStatefulSets(pod); err == nil {
		for _, ss := range sss {
			if selector, err := metav1.LabelSelectorAsSelector(ss.Spec.Selector); err == nil {
				selectors = append(selectors, selector)
			}
		}
	}
	return selectors
}
