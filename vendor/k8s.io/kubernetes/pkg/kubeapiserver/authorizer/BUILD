package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_test(
    name = "go_default_test",
    srcs = ["config_test.go"],
    data = [
        "//pkg/auth/authorizer/abac:example_policy",
    ],
    importpath = "k8s.io/kubernetes/pkg/kubeapiserver/authorizer",
    library = ":go_default_library",
    deps = ["//pkg/kubeapiserver/authorizer/modes:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["config.go"],
    importpath = "k8s.io/kubernetes/pkg/kubeapiserver/authorizer",
    deps = [
        "//pkg/auth/authorizer/abac:go_default_library",
        "//pkg/auth/nodeidentifier:go_default_library",
        "//pkg/client/informers/informers_generated/internalversion:go_default_library",
        "//pkg/kubeapiserver/authorizer/modes:go_default_library",
        "//plugin/pkg/auth/authorizer/node:go_default_library",
        "//plugin/pkg/auth/authorizer/rbac:go_default_library",
        "//plugin/pkg/auth/authorizer/rbac/bootstrappolicy:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/authorization/authorizer:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/authorization/authorizerfactory:go_default_library",
        "//vendor/k8s.io/apiserver/pkg/authorization/union:go_default_library",
        "//vendor/k8s.io/apiserver/plugin/pkg/authorizer/webhook:go_default_library",
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
        "//pkg/kubeapiserver/authorizer/modes:all-srcs",
    ],
    tags = ["automanaged"],
)
