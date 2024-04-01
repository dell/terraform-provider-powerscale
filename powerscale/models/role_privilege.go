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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// RolePrivilegeDataSourceModel describes the data source data model.
type RolePrivilegeDataSourceModel struct {
	ID             types.String               `tfsdk:"id"`
	RolePrivileges []RolePrivilegeDetailModel `tfsdk:"role_privileges_details"`

	// Filters
	RolePrivilegeFilter *RolePrivilegeFilterType `tfsdk:"filter"`
}

// RolePrivilegeDetailModel Specifies the properties for a role privilege.
type RolePrivilegeDetailModel struct {
	URI            types.String `tfsdk:"uri"`
	Category       types.String `tfsdk:"category"`
	Description    types.String `tfsdk:"description"`
	Permission     types.String `tfsdk:"permission"`
	PrivilegeLevel types.String `tfsdk:"privilegelevel"`
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ParentID       types.String `tfsdk:"parent_id"`
}

// RolePrivilegeFilterType describes the filter data model.
type RolePrivilegeFilterType struct {
	Names []types.String `tfsdk:"names"`
}
