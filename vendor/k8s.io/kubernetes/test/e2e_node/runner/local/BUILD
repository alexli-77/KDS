package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "local",
    importpath = "k8s.io/kubernetes/test/e2e_node/runner/local",
    library = ":go_default_library",
)

go_library(
    name = "go_default_library",
    srcs = ["run_local.go"],
    importpath = "k8s.io/kubernetes/test/e2e_node/runner/local",
    deps = [
        "//test/e2e_node/builder:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
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
)
