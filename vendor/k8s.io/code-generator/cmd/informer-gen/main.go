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

package main

import (
	"path/filepath"

	"k8s.io/code-generator/cmd/informer-gen/generators"
	"k8s.io/gengo/args"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

func main() {
	arguments := args.Default()

	// Custom arguments.
	customArgs := &generators.CustomArgs{
		VersionedClientSetPackage: "k8s.io/kubernetes/pkg/client/clientset_generated/clientset",
		InternalClientSetPackage:  "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset",
		ListersPackage:            "k8s.io/kubernetes/pkg/client/listers",
		SingleDirectory:           false,
	}
	customArgs.AddFlags(pflag.CommandLine)

	// Override defaults.
	arguments.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), "k8s.io/kubernetes/hack/boilerplate/boilerplate.go.txt")
	arguments.OutputPackagePath = "k8s.io/kubernetes/pkg/client/informers/informers_generated"
	arguments.CustomArgs = customArgs

	// Run it.
	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		glog.Fatalf("Error: %v", err)
	}
	glog.V(2).Info("Completed successfully.")
}
