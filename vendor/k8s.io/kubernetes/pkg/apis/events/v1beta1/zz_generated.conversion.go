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
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/events/v1beta1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/kubernetes/pkg/apis/core"
	unsafe "unsafe"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1beta1_Event_To_core_Event,
		Convert_core_Event_To_v1beta1_Event,
		Convert_v1beta1_EventList_To_core_EventList,
		Convert_core_EventList_To_v1beta1_EventList,
		Convert_v1beta1_EventSeries_To_core_EventSeries,
		Convert_core_EventSeries_To_v1beta1_EventSeries,
	)
}

func autoConvert_v1beta1_Event_To_core_Event(in *v1beta1.Event, out *core.Event, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.EventTime = in.EventTime
	out.Series = (*core.EventSeries)(unsafe.Pointer(in.Series))
	out.ReportingController = in.ReportingController
	out.ReportingInstance = in.ReportingInstance
	out.Action = in.Action
	out.Reason = in.Reason
	// WARNING: in.Regarding requires manual conversion: does not exist in peer-type
	out.Related = (*core.ObjectReference)(unsafe.Pointer(in.Related))
	// WARNING: in.Note requires manual conversion: does not exist in peer-type
	out.Type = in.Type
	// WARNING: in.DeprecatedSource requires manual conversion: does not exist in peer-type
	// WARNING: in.DeprecatedFirstTimestamp requires manual conversion: does not exist in peer-type
	// WARNING: in.DeprecatedLastTimestamp requires manual conversion: does not exist in peer-type
	// WARNING: in.DeprecatedCount requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_core_Event_To_v1beta1_Event(in *core.Event, out *v1beta1.Event, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	// WARNING: in.InvolvedObject requires manual conversion: does not exist in peer-type
	out.Reason = in.Reason
	// WARNING: in.Message requires manual conversion: does not exist in peer-type
	// WARNING: in.Source requires manual conversion: does not exist in peer-type
	// WARNING: in.FirstTimestamp requires manual conversion: does not exist in peer-type
	// WARNING: in.LastTimestamp requires manual conversion: does not exist in peer-type
	// WARNING: in.Count requires manual conversion: does not exist in peer-type
	out.Type = in.Type
	out.EventTime = in.EventTime
	out.Series = (*v1beta1.EventSeries)(unsafe.Pointer(in.Series))
	out.Action = in.Action
	out.Related = (*v1.ObjectReference)(unsafe.Pointer(in.Related))
	out.ReportingController = in.ReportingController
	out.ReportingInstance = in.ReportingInstance
	return nil
}

func autoConvert_v1beta1_EventList_To_core_EventList(in *v1beta1.EventList, out *core.EventList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]core.Event, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_Event_To_core_Event(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_EventList_To_core_EventList is an autogenerated conversion function.
func Convert_v1beta1_EventList_To_core_EventList(in *v1beta1.EventList, out *core.EventList, s conversion.Scope) error {
	return autoConvert_v1beta1_EventList_To_core_EventList(in, out, s)
}

func autoConvert_core_EventList_To_v1beta1_EventList(in *core.EventList, out *v1beta1.EventList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1beta1.Event, len(*in))
		for i := range *in {
			if err := Convert_core_Event_To_v1beta1_Event(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_core_EventList_To_v1beta1_EventList is an autogenerated conversion function.
func Convert_core_EventList_To_v1beta1_EventList(in *core.EventList, out *v1beta1.EventList, s conversion.Scope) error {
	return autoConvert_core_EventList_To_v1beta1_EventList(in, out, s)
}

func autoConvert_v1beta1_EventSeries_To_core_EventSeries(in *v1beta1.EventSeries, out *core.EventSeries, s conversion.Scope) error {
	out.Count = in.Count
	out.LastObservedTime = in.LastObservedTime
	out.State = core.EventSeriesState(in.State)
	return nil
}

// Convert_v1beta1_EventSeries_To_core_EventSeries is an autogenerated conversion function.
func Convert_v1beta1_EventSeries_To_core_EventSeries(in *v1beta1.EventSeries, out *core.EventSeries, s conversion.Scope) error {
	return autoConvert_v1beta1_EventSeries_To_core_EventSeries(in, out, s)
}

func autoConvert_core_EventSeries_To_v1beta1_EventSeries(in *core.EventSeries, out *v1beta1.EventSeries, s conversion.Scope) error {
	out.Count = in.Count
	out.LastObservedTime = in.LastObservedTime
	out.State = v1beta1.EventSeriesState(in.State)
	return nil
}

// Convert_core_EventSeries_To_v1beta1_EventSeries is an autogenerated conversion function.
func Convert_core_EventSeries_To_v1beta1_EventSeries(in *core.EventSeries, out *v1beta1.EventSeries, s conversion.Scope) error {
	return autoConvert_core_EventSeries_To_v1beta1_EventSeries(in, out, s)
}
