load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "example.go",
        "interface.go",
    ],
    importpath = "k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/informers/externalversions/cr/v1",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/k8s.io/apiextensions-apiserver/examples/client-go/pkg/apis/cr/v1:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/clientset/versioned:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/informers/externalversions/internalinterfaces:go_default_library",
        "//vendor/k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/listers/cr/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
