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

// This package contains hand-coded set implementations that should be similar
// to the autogenerated ones in pkg/util/sets.
// We can't simply use net.IPNet as a map-key in Go (because it contains a
// []byte).
// We could use the same workaround we use here (a string representation as the
// key) to autogenerate sets.  If we do that, or decide on an alternate
// approach, we should replace the implementations in this package with the
// autogenerated versions.
// It is expected that callers will alias this import as "netsets" i.e. import
// netsets "k8s.io/kubernetes/pkg/util/net/sets"

package sets
