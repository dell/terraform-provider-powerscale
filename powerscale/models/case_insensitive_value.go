/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
// Ensure the implementation satisfies the expected interfaces
package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ basetypes.StringValuableWithSemanticEquals = (*CaseInsensitiveStringValue)(nil)
)

type CaseInsensitiveStringValue struct {
	basetypes.StringValue
}

// StringSemanticEquals checks the semantic equality of two CaseInsensitiveStringValue objects.
//
// ctx is the context.Context.
// newValuable is the basetypes.StringValuable to compare.
// Returns a boolean indicating if the values are equal and a diag.Diagnostics.
func (v CaseInsensitiveStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	other, ok := newValuable.(CaseInsensitiveStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}
	tflog.Debug(ctx, "CaseInsensitiveStringValue.StringSemanticEquals", map[string]interface{}{
		"other": other.StringValue.ValueString(), "current": v.StringValue.ValueString(), "equal": strings.EqualFold(v.StringValue.ValueString(), other.StringValue.ValueString())})

	results := false
	if strings.EqualFold(v.StringValue.ValueString(), other.StringValue.ValueString()) {
		results = true
		v = other
	}
	return results, diags
}

// Type returns the type of the CaseInsensitiveStringValue.
//
// ctx is the context.Context.
// Returns a CaseInsensitiveStringType.
func (v CaseInsensitiveStringValue) Type(ctx context.Context) attr.Type {
	return CaseInsensitiveStringType{}
}

// ValueString returns the value of the CaseInsensitiveStringValue.
//
// ctx is the context.Context.
// Returns a string.
func (v CaseInsensitiveStringValue) ValueString() string {
	return v.StringValue.ValueString()
}

// String returns the value of the CaseInsensitiveStringValue.
//
// ctx is the context.Context.
// Returns a string.
func (v CaseInsensitiveStringValue) String() string {
	return v.StringValue.ValueString()
}

// ValueFromTerraform converts a tftypes.Value to an attr.Value.
//
// ctx is the context.Context.
// in is the tftypes.Value to convert.
// Returns an attr.Value and an error.
func (v CaseInsensitiveStringValue) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	if in.IsNull() {
		return CaseInsensitiveStringValue{basetypes.NewStringNull()}, diags
	}

	var value string
	err := in.As(&value)
	if err != nil {
		diags.AddError("Error converting value", err.Error())
		return nil, diags
	}

	return CaseInsensitiveStringValue{basetypes.NewStringValue(value)}, diags
}

// Implement the ValueType method.

// ValueType returns the value type of the CaseInsensitiveStringValue.
func (v CaseInsensitiveStringValue) ValueType(_ context.Context) attr.Value {
	return CaseInsensitiveStringValue{}
}

// NewCaseInsensitiveStringNull creates an CaseInsensitiveStringValue with a null value. Determine whether the value is null via IsNull method.
func NewCaseInsensitiveStringNull() CaseInsensitiveStringValue {
	return CaseInsensitiveStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

// NewCaseInsensitiveStringUnknown creates an CaseInsensitiveStringValue with an unknown value. Determine whether the value is unknown via IsUnknown method.
func NewCaseInsensitiveStringUnknown() CaseInsensitiveStringValue {
	return CaseInsensitiveStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

// NewCaseInsensitiveStringValue creates an CaseInsensitiveStringValue with a known value. Access the value via ValueString method.
func NewCaseInsensitiveStringValue(value string) CaseInsensitiveStringValue {
	return CaseInsensitiveStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

// NewIPv4AddressPointerValue creates an CaseInsensitiveStringValue with a null value if nil or a known value. Access the value via ValueStringPointer method.
func NewCaseInsensitiveStringPointerValue(value *string) CaseInsensitiveStringValue {
	return CaseInsensitiveStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
