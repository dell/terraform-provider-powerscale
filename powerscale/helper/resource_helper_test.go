/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type OpenapiStruct struct {
	BoolPtr           *bool                     `json:"bool_ptr,omitempty"`
	BoolVal           bool                      `json:"bool_val,omitempty"`
	StringPtr         *string                   `json:"string_ptr,omitempty"`
	StringVal         string                    `json:"string_val,omitempty"`
	NullableStringVal powerscale.NullableString `json:"nullable_string_val,omitempty"`
	Int64Ptr          *int64                    `json:"int_64_ptr,omitempty"`
	Int64Val          int64                     `json:"int_64_val,omitempty"`
	NestedSlice       []OpenapiChildStruct      `json:"nested_slice,omitempty"`
	NestedObject      OpenapiChildSingleStruct  `json:"nested_object,omitempty"`
	NestedObjectPtr   *OpenapiChildSingleStruct `json:"nested_object_ptr,omitempty"`
	EmptyObject       *OpenapiChildSingleStruct `json:"empty_obj,omitempty"`
	NestedSlicePtr    *[]OpenapiChildStruct     `json:"nested_slice_ptr,omitempty"`
	EmptySlicePtr     *[]OpenapiChildStruct     `json:"empty_nested_slice_ptr,omitempty"`
	EmptyBoolPtr      *bool                     `json:"empty_bool_ptr,omitempty"`
	EmptyStringPtr    *string                   `json:"empty_string_ptr,omitempty"`
	EmptyInt64Ptr     *int64                    `json:"empty_int_64_ptr,omitempty"`
}

type OpenapiChildStruct struct {
	Str               string                            `json:"str,omitempty"`
	NullableStringVal powerscale.NullableString         `json:"nullable_string_val,omitempty"`
	GrandSlice        []OpenapiGrandChildSingleStruct   `json:"grand_slice,omitempty"`
	GrandSlicePtr     *[]OpenapiGrandChildSingleStruct  `json:"grand_slice_ptr,omitempty"`
	Matrix            [][]OpenapiGrandChildSingleStruct `json:"matrix,omitempty"`
}

type OpenapiChildSingleStruct struct {
	Strings         []string                       `json:"strings,omitempty"`
	Integers        []int64                        `json:"integers,omitempty"`
	Structs         []OpenapiChildStruct           `json:"structs,omitempty"`
	SingleStruct    OpenapiGrandChildSingleStruct  `json:"single_struct,omitempty"`
	SingleStructPtr *OpenapiGrandChildSingleStruct `json:"single_struct_ptr,omitempty"`
}
type OpenapiGrandChildSingleStruct struct {
	String string `json:"str,omitempty"`
	IntVal int64  `json:"int_val,omitempty"`
}

type TfStruct struct {
	BoolPtr           types.Bool   `tfsdk:"bool_ptr"`
	BoolVal           types.Bool   `tfsdk:"bool_val"`
	StringPtr         types.String `tfsdk:"string_ptr"`
	StringVal         types.String `tfsdk:"string_val"`
	NullableStringVal types.String `tfsdk:"nullable_string_val"`
	Int64Ptr          types.Int64  `tfsdk:"int_64_ptr"`
	Int64Val          types.Int64  `tfsdk:"int_64_val"`
	NestedSlice       types.List   `tfsdk:"nested_slice"`
	NestedObject      types.Object `tfsdk:"nested_object"`
	NestedObjectPtr   types.Object `tfsdk:"nested_object_ptr"`
	NestedSlicePtr    types.List   `tfsdk:"nested_slice_ptr"`
	EmptySlicePtr     types.List   `tfsdk:"empty_nested_slice_ptr"`
	EmptyObject       types.Object `tfsdk:"empty_obj"`
	EmptyBoolPtr      types.Bool   `tfsdk:"empty_bool_ptr"`
	EmptyStringPtr    types.String `tfsdk:"empty_string_ptr"`
	EmptyInt64Ptr     types.Int64  `tfsdk:"empty_int_64_ptr"`
}

var fakeBool = true
var fakeString = "fake_string"
var fakeInt = int64(32)

var openapiStructObj = OpenapiStruct{
	BoolPtr:   &fakeBool,
	BoolVal:   fakeBool,
	StringPtr: &fakeString,
	StringVal: fakeString,
	Int64Ptr:  &fakeInt,
	Int64Val:  fakeInt,
	NestedSlice: []OpenapiChildStruct{{
		Str: "fake_child_1",
		Matrix: [][]OpenapiGrandChildSingleStruct{{{
			String: "matrix_1",
			IntVal: 0,
		}}},
	}, {
		Str: "fake_child_2",
		GrandSlicePtr: &[]OpenapiGrandChildSingleStruct{
			{
				String: "string",
				IntVal: int64(1),
			},
		},
	}},
	NestedSlicePtr: &[]OpenapiChildStruct{{
		NullableStringVal: *powerscale.NewNullableString(&fakeString),
		Str:               "fake_child_1",
	}, {
		Str: "fake_child_2",
	}},
	NestedObject: OpenapiChildSingleStruct{
		Strings:  []string{"1", "2", "3"},
		Integers: []int64{1, 2, 3},
		Structs: []OpenapiChildStruct{{
			Str: "single_child_1",
			GrandSlice: []OpenapiGrandChildSingleStruct{{
				String: "grand_1",
			}},
		}, {
			Str: "single_child_2",
		}},
	},
	NestedObjectPtr: &OpenapiChildSingleStruct{
		Integers: []int64{1, 2, 3},
	},
}

func Test_ResourceHelper(t *testing.T) {
	testCopyTfObj := TfStruct{}
	err := CopyFieldsToNonNestedModel(context.Background(), openapiStructObj, &testCopyTfObj)
	assert.Nil(t, err)
	var openapiTarget = OpenapiStruct{}
	err = ReadFromState(context.Background(), testCopyTfObj, &openapiTarget)
	assert.Nil(t, err)
	assert.Equal(t, openapiStructObj, openapiTarget)
}
