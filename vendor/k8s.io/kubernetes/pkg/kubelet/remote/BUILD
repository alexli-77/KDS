package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "remote_image.go",
        "remote_runtime.go",
        "utils.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubelet/remote",
    deps = [
        "//pkg/kubelet/apis/cri:go_default_library",
        "//pkg/kubelet/apis/cri/v1alpha1/runtime:go_default_library",
        "//pkg/kubelet/util:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/golang.org/x/net/context:go_default_library",
        "//vendor/google.golang.org/grpc:go_default_library",
        "//vendor/k8s.io/utils/exec:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//pkg/kubelet/remote/fake:all-srcs",
    ],
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["remote_runtime_test.go"],
    importpath = "k8s.io/kubernetes/pkg/kubelet/remote",
    library = ":go_default_library",
    tags = ["automanaged"],
    deps = [
        "//pkg/kubelet/apis/cri:go_default_library",
        "//pkg/kubelet/apis/cri/testing:go_default_library",
        "//pkg/kubelet/remote/fake:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)
