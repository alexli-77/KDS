package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cronjob.go",
        "expansion_generated.go",
        "job.go",
        "job_expansion.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/client/listers/batch/internalversion",
    deps = [
        "//pkg/apis/batch:go_default_library",
        "//pkg/apis/core:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["job_test.go"],
    importpath = "k8s.io/kubernetes/pkg/client/listers/batch/internalversion",
    library = ":go_default_library",
    deps = [
        "//pkg/apis/batch:go_default_library",
        "//pkg/apis/core:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/sets:go_default_library",
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
)
