package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "register.go",
        "types.go",
        "zz_generated.deepcopy.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/apis/batch",
    deps = [
        "//pkg/apis/core:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
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
        "//pkg/apis/batch/fuzzer:all-srcs",
        "//pkg/apis/batch/install:all-srcs",
        "//pkg/apis/batch/v1:all-srcs",
        "//pkg/apis/batch/v1beta1:all-srcs",
        "//pkg/apis/batch/v2alpha1:all-srcs",
        "//pkg/apis/batch/validation:all-srcs",
    ],
    tags = ["automanaged"],
)
