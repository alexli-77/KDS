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

package proxy

const (
	// KubeProxyConfigMap18 is the proxy ConfigMap manifest for Kubernetes version 1.8
	KubeProxyConfigMap18 = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-proxy
  namespace: kube-system
  labels:
    app: kube-proxy
data:
  kubeconfig.conf: |
    apiVersion: v1
    kind: Config
    clusters:
    - cluster:
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        server: {{ .MasterEndpoint }}
      name: default
    contexts:
    - context:
        cluster: default
        namespace: default
        user: default
      name: default
    current-context: default
    users:
    - name: default
      user:
        tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
`

	// KubeProxyConfigMap19 is the proxy ConfigMap manifest for Kubernetes 1.9 and above
	KubeProxyConfigMap19 = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-proxy
  namespace: kube-system
  labels:
    app: kube-proxy
data:
  kubeconfig.conf: |-
    apiVersion: v1
    kind: Config
    clusters:
    - cluster:
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        server: {{ .MasterEndpoint }}
      name: default
    contexts:
    - context:
        cluster: default
        namespace: default
        user: default
      name: default
    current-context: default
    users:
    - name: default
      user:
        tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
  config.conf: |-
{{ .ProxyConfig}}
`
	// KubeProxyDaemonSet18 is the proxy DaemonSet manifest for Kubernetes version 1.8
	KubeProxyDaemonSet18 = `
apiVersion: apps/v1beta2
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-proxy
  name: kube-proxy
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: kube-proxy
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: kube-proxy
    spec:
      containers:
      - name: kube-proxy
        image: {{ if .ImageOverride }}{{ .ImageOverride }}{{ else }}{{ .ImageRepository }}/kube-proxy-{{ .Arch }}:{{ .Version }}{{ end }}
        imagePullPolicy: IfNotPresent
        command:
        - /usr/local/bin/kube-proxy
        - --kubeconfig=/var/lib/kube-proxy/kubeconfig.conf
        {{ .ClusterCIDR }}
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /var/lib/kube-proxy
          name: kube-proxy
        - mountPath: /run/xtables.lock
          name: xtables-lock
          readOnly: false
        - mountPath: /lib/modules
          name: lib-modules
          readOnly: true
      hostNetwork: true
      serviceAccountName: kube-proxy
      tolerations:
      - key: {{ .MasterTaintKey }}
        effect: NoSchedule
      - key: {{ .CloudTaintKey }}
        value: "true"
        effect: NoSchedule
      volumes:
      - name: kube-proxy
        configMap:
          name: kube-proxy
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
      - name: lib-modules
        hostPath:
          path: /lib/modules
`

	// KubeProxyDaemonSet19 is the proxy DaemonSet manifest for Kubernetes 1.9 and above
	KubeProxyDaemonSet19 = `
apiVersion: apps/v1beta2
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-proxy
  name: kube-proxy
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: kube-proxy
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: kube-proxy
    spec:
      containers:
      - name: kube-proxy
        image: {{ if .ImageOverride }}{{ .ImageOverride }}{{ else }}{{ .ImageRepository }}/kube-proxy-{{ .Arch }}:{{ .Version }}{{ end }}
        imagePullPolicy: IfNotPresent
        command:
        - /usr/local/bin/kube-proxy
        - --config=/var/lib/kube-proxy/config.conf
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /var/lib/kube-proxy
          name: kube-proxy
        - mountPath: /run/xtables.lock
          name: xtables-lock
          readOnly: false
        - mountPath: /lib/modules
          name: lib-modules
          readOnly: true
      hostNetwork: true
      serviceAccountName: kube-proxy
      tolerations:
      - key: {{ .MasterTaintKey }}
        effect: NoSchedule
      - key: {{ .CloudTaintKey }}
        value: "true"
        effect: NoSchedule
      volumes:
      - name: kube-proxy
        configMap:
          name: kube-proxy
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
      - name: lib-modules
        hostPath:
          path: /lib/modules
`
)
