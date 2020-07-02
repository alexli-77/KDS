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

package descheduler

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"github.com/descheduler-controller/cmd/descheduler/app/options"
	"github.com/descheduler-controller/pkg/descheduler/client"
	nodeutil "github.com/descheduler-controller/pkg/descheduler/node"
	"github.com/descheduler-controller/pkg/descheduler/strategies"
	"github.com/descheduler-controller/pkg/api"
	"github.com/descheduler-controller/util"
	"os/exec"
)

func Run(rs *options.DeschedulerServer) error {

	evictLocalStoragePods := true
	rsclient, err := client.CreateClient(rs.KubeconfigFile)
	if err != nil {
		return err
	}
	rs.Client = rsclient
	policyPath,templatePath := util.GetDataFromYaml()
	glog.Warningf("path : %s ", policyPath)
	deschedulerPolicy, err := LoadPolicyConfig(policyPath)
	if err != nil {
		return err
	}
	if deschedulerPolicy == nil {
		return fmt.Errorf("deschedulerPolicy is nil")
	}

	stopChannel := make(chan struct{})
	nodes, err := nodeutil.ReadyNodes(rs.Client, rs.NodeSelector, stopChannel)
	if err != nil {
		return err
	}

	if len(nodes) <= 1 {
		glog.V(1).Infof("The cluster size is 0 or 1 meaning eviction causes service disruption or degradation. So aborting..")
		return nil
	}

	npm := strategies.CreateNodePodsMap(rs.Client, nodes)

	nodesStatus := strategies.GetNodeUsage(npm, evictLocalStoragePods)
	strategies.SortNodesByUsage(nodesStatus,v1.ResourcePods)

	var lowNode,highNode strategies.NodeUsageMap
	for i :=0;i < len(nodesStatus);i++ {
		if false == nodesStatus[i].SIsShare{
			continue
		} else {
			highNode = nodesStatus[i]
			break
		}
	}
	glog.Warningf("NO.1 : %s , cpu: %s, mem: %s, pod: %s", highNode.Node.Name, highNode.Usage["cpu"], highNode.Usage["memory"], highNode.Usage["pods"])
	for i :=len(nodesStatus)-1;i >= 0;i-- {
		if true == nodesStatus[i].SIsNodeUschedulable || false == nodesStatus[i].SIsShare{
			continue
		} else {
			lowNode = nodesStatus[i]
			break
		}
	}

	var cpu float64
	var memory float64
	var pods float64
	var i int
	i = 0

	for _, nodes := range nodesStatus {
//		glog.Warningf("lownode: %s , cpu: %s, mem: %s, pod: %s", nodes.Node.Name, nodes.Usage["cpu"], nodes.Usage["memory"], nodes.Usage["pods"],nodes.SIsNodeUschedulable)
		cpu += float64(nodes.Usage["cpu"])
		memory += float64(nodes.Usage["memory"])
		pods += float64(nodes.Usage["pods"])
		i++
	}

	avgCpu := float64(cpu/float64(i))
	avgMemory := float64(memory/float64(i))
	avgPods := float64(pods/float64(i))
	glog.Warningf("last : %s , cpu: %s, mem: %s, pod: %s, avgCpu: %s, avgMemory: %s, avgPods: %s", lowNode.Node.Name, lowNode.Usage["cpu"], lowNode.Usage["memory"], lowNode.Usage["pods"],avgCpu,avgMemory,avgPods)

	thresholds := make(api.ResourceThresholds)
	targetThresholds := make(api.ResourceThresholds)

	if rs.EvictMode == "avg" {
		targetThresholds[v1.ResourceCPU] = api.Percentage(avgCpu)
		targetThresholds[v1.ResourceMemory] = api.Percentage(avgMemory)
		targetThresholds[v1.ResourcePods] = api.Percentage(avgPods)
	}else {
		targetThresholds[v1.ResourceCPU] = highNode.Usage["cpu"]
		targetThresholds[v1.ResourceMemory] = highNode.Usage["memory"]
		targetThresholds[v1.ResourcePods] = highNode.Usage["pods"]
	}
	thresholds[v1.ResourceCPU] = lowNode.Usage["cpu"] + 0.001
	thresholds[v1.ResourceMemory] = lowNode.Usage["memory"] + 0.001
	thresholds[v1.ResourcePods] = lowNode.Usage["pods"] + 0.001

//		util.Modifyfile("./template.yaml",thresholds,targetThresholds)
	util.Modifyfile(templatePath,thresholds,targetThresholds)

	//If the difference between the maximum threshold and the average is less than 1, it is considered uniform
	if (float64(thresholds[v1.ResourcePods]) - float64(avgPods)) <= 1 {
		cmd := exec.Command("/bin/bash", "-c", "./descheduler --kubeconfig /root/.kube/config --policy-config-file ./config/policy.yaml --max-pods-to-evict-per-node 2 --evict-local-storage-pods --node-selector datatype!=inst")
		buf, err := cmd.Output()
		if err != nil{
			glog.Errorf(err.Error())
			return err
		}
		glog.Infof(string(buf))
	}


	return nil
}
