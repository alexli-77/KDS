package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "fake_authorization_client.go",
        "fake_generated_expansion.go",
        "fake_localsubjectaccessreview.go",
        "fake_localsubjectaccessreview_expansion.go",
        "fake_selfsubjectaccessreview.go",
        "fake_selfsubjectaccessreview_expansion.go",
        "fake_selfsubjectrulesreview.go",
        "fake_selfsubjectrulesreview_expansion.go",
        "fake_subjectaccessreview.go",
        "fake_subjectaccessreview_expansion.go",
    ],
    importpath = "k8s.io/client-go/kubernetes/typed/authorization/v1beta1/fake",
    deps = [
        "//vendor/k8s.io/api/authorization/v1beta1:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/typed/authorization/v1beta1:go_default_library",
        "//vendor/k8s.io/client-go/rest:go_default_library",
        "//vendor/k8s.io/client-go/testing:go_default_library",
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
