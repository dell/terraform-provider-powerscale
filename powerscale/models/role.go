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

// RoleResourceModel describes the resource data model.
type RoleResourceModel struct {
	// Query param
	// Specifies which access zone to use.
	Zone types.String `tfsdk:"zone"`

	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Members     types.List   `tfsdk:"members"`
	Privileges  types.List   `tfsdk:"privileges"`
	Description types.String `tfsdk:"description"`
}

// RoleDataSourceModel describes the data source data model.
type RoleDataSourceModel struct {
	ID    types.String      `tfsdk:"id"`
	Roles []RoleDetailModel `tfsdk:"roles_details"`

	// Filters
	RoleFilter *RoleFilterType `tfsdk:"filter"`
}

// RoleDetailModel Specifies the properties for a role.
type RoleDetailModel struct {
	// Specifies the ID of the role.
	ID types.String `tfsdk:"id"`
	// Specifies the description of the role.
	Description types.String `tfsdk:"description"`
	// Specifies the users or groups that have this role.
	Members []V1AuthAccessAccessItemFileGroup `tfsdk:"members"`
	// Specifies the name of the role.
	Name types.String `tfsdk:"name"`
	// Specifies the privileges granted by this role.
	Privileges []V14AuthIDNtokenPrivilegeItem `tfsdk:"privileges"`
}

// V14AuthIDNtokenPrivilegeItem Specifies the system-defined privilege that may be granted to users.
type V14AuthIDNtokenPrivilegeItem struct {
	// Specifies the ID of the privilege.
	ID types.String `tfsdk:"id"`
	// Specifies the name of the privilege.
	Name types.String `tfsdk:"name"`
	// permission of the privilege, 'r' = read , 'x' = read-execute, 'w' = read-execute-write, '-' = no permission
	Permission types.String `tfsdk:"permission"`
}

// RoleFilterType describes the filter data model.
type RoleFilterType struct {
	Names []types.String `tfsdk:"names"`
	Zone  types.String   `tfsdk:"zone"`
}
