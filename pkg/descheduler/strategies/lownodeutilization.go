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

package strategies

import (
	"sort"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	clientset "k8s.io/client-go/kubernetes"
	helper "k8s.io/kubernetes/pkg/api/v1/resource"

	//"github.com/descheduler/cmd/descheduler/app/options"
	"github.com/descheduler-controller/pkg/api"
	nodeutil "github.com/descheduler-controller/pkg/descheduler/node"
	podutil "github.com/descheduler-controller/pkg/descheduler/pod"
)

type NodeUsageMap struct {
	Node             *v1.Node
	Usage            api.ResourceThresholds
	allPods          []*v1.Pod
	nonRemovablePods []*v1.Pod
	bePods           []*v1.Pod
	bPods            []*v1.Pod
	gPods            []*v1.Pod
	SIsNodeUschedulable bool
	SIsShare bool
}

type NodePodsMap map[*v1.Node][]*v1.Pod


func validateThresholds(thresholds api.ResourceThresholds) bool {
	if thresholds == nil || len(thresholds) == 0 {
		glog.V(1).Infof("no resource threshold is configured")
		return false
	}
	for name := range thresholds {
		switch name {
		case v1.ResourceCPU:
			continue
		case v1.ResourceMemory:
			continue
		case v1.ResourcePods:
			continue
		default:
			glog.Errorf("only cpu, memory, or pods thresholds can be specified")
			return false
		}
	}
	return true
}

//This function could be merged into above once we are clear.
func validateTargetThresholds(targetThresholds api.ResourceThresholds) bool {
	if targetThresholds == nil {
		glog.V(1).Infof("no target resource threshold is configured")
		return false
	} else if _, ok := targetThresholds[v1.ResourcePods]; !ok {
		glog.V(1).Infof("no target resource threshold for pods is configured")
		return false
	}
	return true
}

func GetNodeUsage(npm NodePodsMap, evictLocalStoragePods bool) []NodeUsageMap {
	NodesStatus := []NodeUsageMap{}
	var sIsNodeUschedulable bool
	var sIsShare bool
	for node, pods := range npm {
		usage, allPods, nonRemovablePods, bePods, bPods, gPods := NodeUtilization(node, pods, evictLocalStoragePods)
		if nodeutil.IsNodeUschedulable(node) {
			sIsNodeUschedulable = true
		} else {
			sIsNodeUschedulable = false
		}
		if node.Labels["datatype"] == "share"{
			sIsShare = true
		}else {
			sIsShare = false
		}
		nuMap := NodeUsageMap{node, usage, allPods, nonRemovablePods, bePods, bPods, gPods,sIsNodeUschedulable,sIsShare}
		
		NodesStatus = append(NodesStatus, nuMap)
	}
	return NodesStatus
}

func SortNodesByUsage(nodes []NodeUsageMap,resourcename v1.ResourceName) {
	sort.Slice(nodes, func(i, j int) bool {
		var ti, tj api.Percentage
		for name, value := range nodes[i].Usage {
			if name == v1.ResourcePods {
				ti += value
			}
		}
		for name, value := range nodes[j].Usage {
			if name == v1.ResourcePods {
				tj += value
			}
		}
		// To return sorted in descending order
		return ti > tj
	})
}

// createNodePodsMap returns nodepodsmap with evictable pods on node.
func CreateNodePodsMap(client clientset.Interface, nodes []*v1.Node) NodePodsMap {
	npm := NodePodsMap{}
	for _, node := range nodes {
		pods, err := podutil.ListPodsOnANode(client, node)
		if err != nil {
			glog.Warningf("node %s will not be processed, error in accessing its pods (%#v)", node.Name, err)
		} else {
			npm[node] = pods
		}
	}
	return npm
}

func IsNodeAboveTargetUtilization(nodeThresholds api.ResourceThresholds, thresholds api.ResourceThresholds) bool {
	for name, nodeValue := range nodeThresholds {
		if name == v1.ResourceCPU || name == v1.ResourceMemory || name == v1.ResourcePods {
			if value, ok := thresholds[name]; !ok {
				continue
			} else if nodeValue > value {
				return true
			}
		}
	}
	return false
}

func IsNodeWithLowUtilization(nodeThresholds api.ResourceThresholds, thresholds api.ResourceThresholds) bool {
	for name, nodeValue := range nodeThresholds {
		if name == v1.ResourceCPU || name == v1.ResourceMemory || name == v1.ResourcePods {
			if value, ok := thresholds[name]; !ok {
				continue
			} else if nodeValue > value {
				return false
			}
		}
	}
	return true
}

// Nodeutilization returns the current usage of node.
func NodeUtilization(node *v1.Node, pods []*v1.Pod, evictLocalStoragePods bool) (api.ResourceThresholds, []*v1.Pod, []*v1.Pod, []*v1.Pod, []*v1.Pod, []*v1.Pod) {
	bePods := []*v1.Pod{}
	nonRemovablePods := []*v1.Pod{}
	bPods := []*v1.Pod{}
	gPods := []*v1.Pod{}
	totalReqs := map[v1.ResourceName]resource.Quantity{}
	for _, pod := range pods {
		// We need to compute the usage of nonRemovablePods unless it is a best effort pod. So, cannot use podutil.ListEvictablePodsOnNode
		if !podutil.IsEvictable(pod, evictLocalStoragePods) {
			nonRemovablePods = append(nonRemovablePods, pod)
			if podutil.IsBestEffortPod(pod) {
				continue
			}
		} else if podutil.IsBestEffortPod(pod) {
			bePods = append(bePods, pod)
			continue
		} else if podutil.IsBurstablePod(pod) {
			bPods = append(bPods, pod)
		} else {
			gPods = append(gPods, pod)
		}

		req, _ := helper.PodRequestsAndLimits(pod)
		for name, quantity := range req {
			if name == v1.ResourceCPU || name == v1.ResourceMemory {
				if value, ok := totalReqs[name]; !ok {
					totalReqs[name] = *quantity.Copy()
				} else {
					value.Add(quantity)
					totalReqs[name] = value
				}
			}
		}
	}

	nodeCapacity := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		nodeCapacity = node.Status.Allocatable
	}

	usage := api.ResourceThresholds{}
	totalCPUReq := totalReqs[v1.ResourceCPU]
	totalMemReq := totalReqs[v1.ResourceMemory]
	totalPods := len(pods)
	usage[v1.ResourceCPU] = api.Percentage((float64(totalCPUReq.MilliValue()) * 100) / float64(nodeCapacity.Cpu().MilliValue()))
	usage[v1.ResourceMemory] = api.Percentage(float64(totalMemReq.Value()) / float64(nodeCapacity.Memory().Value()) * 100)
	usage[v1.ResourcePods] = api.Percentage((float64(totalPods) * 100) / float64(nodeCapacity.Pods().Value()))
	return usage, pods, nonRemovablePods, bePods, bPods, gPods
}
