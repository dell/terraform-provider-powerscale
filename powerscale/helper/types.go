/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

import "github.com/hashicorp/terraform-plugin-framework/types"

// GetKnownStringPointer returns a pointer to the string value if it is known, otherwise nil
func GetKnownStringPointer(in types.String) *string {
	if in.IsUnknown() {
		return nil
	}
	return in.ValueStringPointer()
}

// GetKnownBoolPointer returns a pointer to the bool value if it is known, otherwise nil
func GetKnownBoolPointer(in types.Bool) *bool {
	if in.IsUnknown() {
		return nil
	}
	return in.ValueBoolPointer()
}

// New returns a pointer to a copy of the given value.
func New[T any](in T) *T {
	return &in
}
