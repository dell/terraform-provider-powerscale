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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = (*CaseInsensitiveStringType)(nil)

type CaseInsensitiveStringType struct {
	basetypes.StringType
}

// Equal checks if the given attribute type is equal to the current type.
//
// It takes an attribute type as a parameter.
// Returns a boolean indicating if the given type is equal to the current type.
func (t CaseInsensitiveStringType) Equal(o attr.Type) bool {
	other, ok := o.(CaseInsensitiveStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// String returns the string representation of the CaseInsensitiveStringType.
//
// It takes no parameters.
// Returns a string.
func (t CaseInsensitiveStringType) String() string {
	return "customtypes.CaseInsensitiveStringType"
}

// ValueFromString converts a basetypes.StringValue to a CaseInsensitiveStringValue.
//
// ctx is the context.Context.
// in is the input basetypes.StringValue.
// Returns a basetypes.StringValuable and a diag.Diagnostics.
func (t CaseInsensitiveStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	// CaseInsensitiveStringValue defined in the value type section
	value := CaseInsensitiveStringValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts a tftypes.Value to an attr.Value.
//
// ctx is the context.Context.
// in is the tftypes.Value to convert.
// Returns an attr.Value and an error.
func (t CaseInsensitiveStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// Validate validates the given tftypes.Value.
//
// ctx is the context.Context.
// in is the tftypes.Value to validate.
// path is the path.Path.
// Returns diag.Diagnostics.
func (t CaseInsensitiveStringType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if in.Type() == nil {
		return diags
	}

	return diags
}

// ValueType returns the value type of the CaseInsensitiveStringType.
//
// ctx is the context.Context.
// Returns a CaseInsensitiveStringValue.
func (t CaseInsensitiveStringType) ValueType(ctx context.Context) attr.Value {
	return CaseInsensitiveStringValue{}
}
