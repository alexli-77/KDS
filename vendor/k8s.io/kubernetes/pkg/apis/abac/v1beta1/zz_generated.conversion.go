// +build !ignore_autogenerated

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

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1beta1

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	abac "k8s.io/kubernetes/pkg/apis/abac"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1beta1_Policy_To_abac_Policy,
		Convert_abac_Policy_To_v1beta1_Policy,
		Convert_v1beta1_PolicySpec_To_abac_PolicySpec,
		Convert_abac_PolicySpec_To_v1beta1_PolicySpec,
	)
}

func autoConvert_v1beta1_Policy_To_abac_Policy(in *Policy, out *abac.Policy, s conversion.Scope) error {
	if err := Convert_v1beta1_PolicySpec_To_abac_PolicySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_Policy_To_abac_Policy is an autogenerated conversion function.
func Convert_v1beta1_Policy_To_abac_Policy(in *Policy, out *abac.Policy, s conversion.Scope) error {
	return autoConvert_v1beta1_Policy_To_abac_Policy(in, out, s)
}

func autoConvert_abac_Policy_To_v1beta1_Policy(in *abac.Policy, out *Policy, s conversion.Scope) error {
	if err := Convert_abac_PolicySpec_To_v1beta1_PolicySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_abac_Policy_To_v1beta1_Policy is an autogenerated conversion function.
func Convert_abac_Policy_To_v1beta1_Policy(in *abac.Policy, out *Policy, s conversion.Scope) error {
	return autoConvert_abac_Policy_To_v1beta1_Policy(in, out, s)
}

func autoConvert_v1beta1_PolicySpec_To_abac_PolicySpec(in *PolicySpec, out *abac.PolicySpec, s conversion.Scope) error {
	out.User = in.User
	out.Group = in.Group
	out.Readonly = in.Readonly
	out.APIGroup = in.APIGroup
	out.Resource = in.Resource
	out.Namespace = in.Namespace
	out.NonResourcePath = in.NonResourcePath
	return nil
}

// Convert_v1beta1_PolicySpec_To_abac_PolicySpec is an autogenerated conversion function.
func Convert_v1beta1_PolicySpec_To_abac_PolicySpec(in *PolicySpec, out *abac.PolicySpec, s conversion.Scope) error {
	return autoConvert_v1beta1_PolicySpec_To_abac_PolicySpec(in, out, s)
}

func autoConvert_abac_PolicySpec_To_v1beta1_PolicySpec(in *abac.PolicySpec, out *PolicySpec, s conversion.Scope) error {
	out.User = in.User
	out.Group = in.Group
	out.Readonly = in.Readonly
	out.APIGroup = in.APIGroup
	out.Resource = in.Resource
	out.Namespace = in.Namespace
	out.NonResourcePath = in.NonResourcePath
	return nil
}

// Convert_abac_PolicySpec_To_v1beta1_PolicySpec is an autogenerated conversion function.
func Convert_abac_PolicySpec_To_v1beta1_PolicySpec(in *abac.PolicySpec, out *PolicySpec, s conversion.Scope) error {
	return autoConvert_abac_PolicySpec_To_v1beta1_PolicySpec(in, out, s)
}
