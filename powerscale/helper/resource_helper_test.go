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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type OpenapiStruct struct {
	BoolPtr      *bool                    `json:"bool_ptr,omitempty"`
	BoolVal      bool                     `json:"bool_val,omitempty"`
	StringPtr    *string                  `json:"string_ptr,omitempty"`
	StringVal    string                   `json:"string_val,omitempty"`
	Int64Ptr     *int64                   `json:"int_64_ptr,omitempty"`
	Int64Val     int64                    `json:"int_64_val,omitempty"`
	NestedSlice  []OpenapiChildStruct     `json:"nested_slice,omitempty"`
	NestedObject OpenapiChildSingleStruct `json:"nested_object,omitempty"`
}

type OpenapiChildStruct struct {
	Str string `json:"str,omitempty"`
}

type OpenapiChildSingleStruct struct {
	Strings      []string                      `json:"strings,omitempty"`
	Integers     []int64                       `json:"integers,omitempty"`
	Structs      []OpenapiChildStruct          `json:"structs,omitempty"`
	SingleStruct OpenapiGrandChildSingleStruct `json:"single_struct,omitempty"`
}
type OpenapiGrandChildSingleStruct struct {
	String string `json:"str,omitempty"`
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
	}, {
		Str: "fake_child_2",
	}},
	NestedObject: OpenapiChildSingleStruct{
		Strings:  []string{"1", "2", "3"},
		Integers: []int64{1, 2, 3},
		Structs: []OpenapiChildStruct{{
			Str: "single_child_1",
		}, {
			Str: "single_child_2",
		}},
	},
}

type TfStruct struct {
	BoolPtr      types.Bool   `tfsdk:"bool_ptr"`
	BoolVal      types.Bool   `tfsdk:"bool_val"`
	StringPtr    types.String `tfsdk:"string_ptr"`
	StringVal    types.String `tfsdk:"string_val"`
	Int64Ptr     types.Int64  `tfsdk:"int_64_ptr"`
	Int64Val     types.Int64  `tfsdk:"int_64_val"`
	NestedSlice  types.List   `tfsdk:"nested_slice"`
	NestedObject types.Object `tfsdk:"nested_object"`
}

func Test_CopyFields(t *testing.T) {
	testCopyTfObj := TfStruct{}
	err := CopyFieldsToNonNestedModel(context.Background(), openapiStructObj, &testCopyTfObj)
	assert.Equal(t, fakeBool, testCopyTfObj.BoolPtr.ValueBool())
	assert.Equal(t, fakeString, testCopyTfObj.StringPtr.ValueString())
	assert.Equal(t, fakeInt, testCopyTfObj.Int64Val.ValueInt64())
	assert.Equal(t, 2, len(testCopyTfObj.NestedSlice.Elements()))
	assert.False(t, testCopyTfObj.NestedObject.IsNull())
	assert.Nil(t, err)
}

func Test_ReadFromState(t *testing.T) {
	nestedAttrMap := map[string]attr.Type{
		"strings": types.ListType{
			ElemType: types.StringType,
		},
	}
	nestedObj, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"state1, state2, state3"})
	nestedValueMap := map[string]attr.Value{
		"strings": nestedObj,
	}
	obj, _ := types.ObjectValue(nestedAttrMap, nestedValueMap)
	sliceAttrMap := map[string]attr.Type{
		"str": types.StringType,
	}
	sliceAttrVal := map[string]attr.Value{
		"str": types.StringValue("slice_1"),
	}
	sliceObj, _ := types.ObjectValue(sliceAttrMap, sliceAttrVal)
	sliceObjs, _ := types.ListValue(types.ObjectType{
		AttrTypes: sliceAttrMap,
	}, []attr.Value{sliceObj})
	tfStructObj := TfStruct{
		BoolPtr:      types.BoolValue(fakeBool),
		BoolVal:      types.BoolValue(fakeBool),
		StringPtr:    types.StringValue(fakeString),
		StringVal:    types.StringValue(fakeString),
		Int64Ptr:     types.Int64Value(fakeInt),
		Int64Val:     types.Int64Value(fakeInt),
		NestedObject: obj,
		NestedSlice:  sliceObjs,
	}
	target := OpenapiStruct{}
	err := ReadFromState(context.Background(), tfStructObj, &target)

	assert.Nil(t, err)
}
