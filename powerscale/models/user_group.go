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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// UserGroupDataSourceModel describes the data source data model.
type UserGroupDataSourceModel struct {
	UserGroups []UserGroupModel `tfsdk:"user_groups"`
	ID         types.String     `tfsdk:"id"`

	//filter
	Filter *UserGroupFilterType `tfsdk:"filter"`
}

// UserGroupModel holds user group data source schema attribute details.
type UserGroupModel struct {
	Dn             types.String                      `tfsdk:"dn"`
	DNSDomain      types.String                      `tfsdk:"dns_domain"`
	Domain         types.String                      `tfsdk:"domain"`
	GeneratedGID   types.Bool                        `tfsdk:"generated_gid"`
	GID            types.String                      `tfsdk:"gid"`
	ID             types.String                      `tfsdk:"id"`
	Name           types.String                      `tfsdk:"name"`
	Provider       types.String                      `tfsdk:"provider"`
	SamAccountName types.String                      `tfsdk:"sam_account_name"`
	SID            types.String                      `tfsdk:"sid"`
	Type           types.String                      `tfsdk:"type"`
	Roles          types.List                        `tfsdk:"roles"`
	Members        []V1AuthAccessAccessItemFileGroup `tfsdk:"members"`
}

// UserGroupFilterType holds filter attribute for user group.
type UserGroupFilterType struct {
	Names        []UserGroupIdentityItem `tfsdk:"names"`
	NamePrefix   types.String            `tfsdk:"name_prefix"`
	Domain       types.String            `tfsdk:"domain"`
	Zone         types.String            `tfsdk:"zone"`
	Provider     types.String            `tfsdk:"provider"`
	Cached       types.Bool              `tfsdk:"cached"`
	ResolveNames types.Bool              `tfsdk:"resolve_names"`
}

// UserGroupIdentityItem holds identity attribute for a auth group.
type UserGroupIdentityItem struct {
	Name types.String `tfsdk:"name"`
	GID  types.Int64  `tfsdk:"gid"`
}

// UserGroupReourceModel describes the resource data model.
type UserGroupReourceModel struct {
	QueryForce     types.Bool   `tfsdk:"query_force"`
	QueryZone      types.String `tfsdk:"query_zone"`
	QueryProvider  types.String `tfsdk:"query_provider"`
	GID            types.Int64  `tfsdk:"gid"`
	Name           types.String `tfsdk:"name"`
	SID            types.String `tfsdk:"sid"`
	Roles          types.List   `tfsdk:"roles"`
	Users          types.List   `tfsdk:"users"`
	Dn             types.String `tfsdk:"dn"`
	DNSDomain      types.String `tfsdk:"dns_domain"`
	Domain         types.String `tfsdk:"domain"`
	GeneratedGID   types.Bool   `tfsdk:"generated_gid"`
	ID             types.String `tfsdk:"id"`
	Provider       types.String `tfsdk:"provider_name"`
	SamAccountName types.String `tfsdk:"sam_account_name"`
	Type           types.String `tfsdk:"type"`
}
