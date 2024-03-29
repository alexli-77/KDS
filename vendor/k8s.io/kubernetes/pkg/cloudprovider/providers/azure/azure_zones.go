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

package azure

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

const instanceInfoURL = "http://169.254.169.254/metadata/v1/InstanceInfo"

var faultMutex = &sync.Mutex{}
var faultDomain *string

type instanceInfo struct {
	ID           string `json:"ID"`
	UpdateDomain string `json:"UD"`
	FaultDomain  string `json:"FD"`
}

// GetZone returns the Zone containing the current failure zone and locality region that the program is running in
func (az *Cloud) GetZone() (cloudprovider.Zone, error) {
	faultMutex.Lock()
	defer faultMutex.Unlock()
	if faultDomain == nil {
		var err error
		faultDomain, err = fetchFaultDomain()
		if err != nil {
			return cloudprovider.Zone{}, err
		}
	}
	zone := cloudprovider.Zone{
		FailureDomain: *faultDomain,
		Region:        az.Location,
	}
	return zone, nil
}

// GetZoneByProviderID implements Zones.GetZoneByProviderID
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
func (az *Cloud) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	nodeName, err := splitProviderID(providerID)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	return az.GetZoneByNodeName(nodeName)
}

// GetZoneByNodeName implements Zones.GetZoneByNodeName
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
func (az *Cloud) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {

	vm, err := az.VirtualMachinesClient.Get(az.ResourceGroup, string(nodeName), compute.InstanceView)

	if err != nil {
		return cloudprovider.Zone{}, err
	}

	failureDomain := strconv.Itoa(int(*vm.VirtualMachineProperties.InstanceView.PlatformFaultDomain))

	zone := cloudprovider.Zone{
		FailureDomain: failureDomain,
		Region:        *(vm.Location),
	}
	return zone, nil
}

func fetchFaultDomain() (*string, error) {
	resp, err := http.Get(instanceInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return readFaultDomain(resp.Body)
}

func readFaultDomain(reader io.Reader) (*string, error) {
	var instanceInfo instanceInfo
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &instanceInfo)
	if err != nil {
		return nil, err
	}
	return &instanceInfo.FaultDomain, nil
}
