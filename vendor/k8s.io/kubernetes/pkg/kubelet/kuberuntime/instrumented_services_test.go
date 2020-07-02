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

package kuberuntime

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
)

func TestRecordOperation(t *testing.T) {
	prometheus.MustRegister(metrics.RuntimeOperations)
	prometheus.MustRegister(metrics.RuntimeOperationsLatency)
	prometheus.MustRegister(metrics.RuntimeOperationsErrors)

	temporalServer := "127.0.0.1:1234"
	l, err := net.Listen("tcp", temporalServer)
	assert.NoError(t, err)
	defer l.Close()

	prometheusUrl := "http://" + temporalServer + "/metrics"
	mux := http.NewServeMux()
	mux.Handle("/metrics", prometheus.Handler())
	server := &http.Server{
		Addr:    temporalServer,
		Handler: mux,
	}
	go func() {
		server.Serve(l)
	}()

	recordOperation("create_container", time.Now())
	runtimeOperationsCounterExpected := "kubelet_runtime_operations{operation_type=\"create_container\"} 1"
	runtimeOperationsLatencyExpected := "kubelet_runtime_operations_latency_microseconds_count{operation_type=\"create_container\"} 1"

	assert.HTTPBodyContains(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}), "GET", prometheusUrl, nil, runtimeOperationsCounterExpected)

	assert.HTTPBodyContains(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}), "GET", prometheusUrl, nil, runtimeOperationsLatencyExpected)
}

func TestInstrumentedVersion(t *testing.T) {
	fakeRuntime, _, _, _ := createTestRuntimeManager()
	irs := newInstrumentedRuntimeService(fakeRuntime)
	vr, err := irs.Version("1")
	assert.NoError(t, err)
	assert.Equal(t, kubeRuntimeAPIVersion, vr.Version)
}

func TestStatus(t *testing.T) {
	fakeRuntime, _, _, _ := createTestRuntimeManager()
	fakeRuntime.FakeStatus = &runtimeapi.RuntimeStatus{
		Conditions: []*runtimeapi.RuntimeCondition{
			{Type: runtimeapi.RuntimeReady, Status: false},
			{Type: runtimeapi.NetworkReady, Status: true},
		},
	}
	irs := newInstrumentedRuntimeService(fakeRuntime)
	actural, err := irs.Status()
	assert.NoError(t, err)
	expected := &runtimeapi.RuntimeStatus{
		Conditions: []*runtimeapi.RuntimeCondition{
			{Type: runtimeapi.RuntimeReady, Status: false},
			{Type: runtimeapi.NetworkReady, Status: true},
		},
	}
	assert.Equal(t, expected, actural)
}
