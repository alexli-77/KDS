# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

all: push

# Keep this one version ahead, so no one accidentally blows away the latest published version.
TAG = 1.1

explorer: explorer.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' ./explorer.go

container: explorer
	docker build --pull -t gcr.io/google_containers/explorer:$(TAG) .

push: container
	gcloud docker -- push gcr.io/google_containers/explorer:$(TAG)

clean:
	rm -f explorer
