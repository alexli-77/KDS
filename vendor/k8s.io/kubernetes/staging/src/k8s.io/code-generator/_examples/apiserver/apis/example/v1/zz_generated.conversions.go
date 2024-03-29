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

package v1

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	example "k8s.io/code-generator/_examples/apiserver/apis/example"
	unsafe "unsafe"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1_TestType_To_example_TestType,
		Convert_example_TestType_To_v1_TestType,
		Convert_v1_TestTypeList_To_example_TestTypeList,
		Convert_example_TestTypeList_To_v1_TestTypeList,
		Convert_v1_TestTypeStatus_To_example_TestTypeStatus,
		Convert_example_TestTypeStatus_To_v1_TestTypeStatus,
	)
}

func autoConvert_v1_TestType_To_example_TestType(in *TestType, out *example.TestType, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1_TestTypeStatus_To_example_TestTypeStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1_TestType_To_example_TestType is an autogenerated conversion function.
func Convert_v1_TestType_To_example_TestType(in *TestType, out *example.TestType, s conversion.Scope) error {
	return autoConvert_v1_TestType_To_example_TestType(in, out, s)
}

func autoConvert_example_TestType_To_v1_TestType(in *example.TestType, out *TestType, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_example_TestTypeStatus_To_v1_TestTypeStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_example_TestType_To_v1_TestType is an autogenerated conversion function.
func Convert_example_TestType_To_v1_TestType(in *example.TestType, out *TestType, s conversion.Scope) error {
	return autoConvert_example_TestType_To_v1_TestType(in, out, s)
}

func autoConvert_v1_TestTypeList_To_example_TestTypeList(in *TestTypeList, out *example.TestTypeList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]example.TestType)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1_TestTypeList_To_example_TestTypeList is an autogenerated conversion function.
func Convert_v1_TestTypeList_To_example_TestTypeList(in *TestTypeList, out *example.TestTypeList, s conversion.Scope) error {
	return autoConvert_v1_TestTypeList_To_example_TestTypeList(in, out, s)
}

func autoConvert_example_TestTypeList_To_v1_TestTypeList(in *example.TestTypeList, out *TestTypeList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]TestType)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_example_TestTypeList_To_v1_TestTypeList is an autogenerated conversion function.
func Convert_example_TestTypeList_To_v1_TestTypeList(in *example.TestTypeList, out *TestTypeList, s conversion.Scope) error {
	return autoConvert_example_TestTypeList_To_v1_TestTypeList(in, out, s)
}

func autoConvert_v1_TestTypeStatus_To_example_TestTypeStatus(in *TestTypeStatus, out *example.TestTypeStatus, s conversion.Scope) error {
	out.Blah = in.Blah
	return nil
}

// Convert_v1_TestTypeStatus_To_example_TestTypeStatus is an autogenerated conversion function.
func Convert_v1_TestTypeStatus_To_example_TestTypeStatus(in *TestTypeStatus, out *example.TestTypeStatus, s conversion.Scope) error {
	return autoConvert_v1_TestTypeStatus_To_example_TestTypeStatus(in, out, s)
}

func autoConvert_example_TestTypeStatus_To_v1_TestTypeStatus(in *example.TestTypeStatus, out *TestTypeStatus, s conversion.Scope) error {
	out.Blah = in.Blah
	return nil
}

// Convert_example_TestTypeStatus_To_v1_TestTypeStatus is an autogenerated conversion function.
func Convert_example_TestTypeStatus_To_v1_TestTypeStatus(in *example.TestTypeStatus, out *TestTypeStatus, s conversion.Scope) error {
	return autoConvert_example_TestTypeStatus_To_v1_TestTypeStatus(in, out, s)
}
