package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "base.go",
        "doc.go",
        "version.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/version",
    deps = ["//vendor/k8s.io/apimachinery/pkg/version:go_default_library"],
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
        "//pkg/version/prometheus:all-srcs",
        "//pkg/version/verflag:all-srcs",
    ],
    tags = ["automanaged"],
)
