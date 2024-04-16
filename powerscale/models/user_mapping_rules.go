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

// UserMappingRulesDataSourceModel describes the data source data model.
type UserMappingRulesDataSourceModel struct {
	Parameters types.Object                `tfsdk:"user_mapping_rules_parameters"`
	Rules      types.List                  `tfsdk:"user_mapping_rules"`
	ID         types.String                `tfsdk:"id"`
	Filter     *UserMappingRulesFilterType `tfsdk:"filter"`
}

// UserMappingRulesFilterType holds filter attribute for User Mapping Rules.
type UserMappingRulesFilterType struct {
	// names filter for source and target users
	Names []types.String `tfsdk:"names"`
	// operators filter for user mapping rules.
	Operators []types.String `tfsdk:"operators"`
	// defaults to System
	Zone types.String `tfsdk:"zone"`
}

// UserMappingRulesResourceModel holds user mapping rules resource schema attribute details.
type UserMappingRulesResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Zone       types.String `tfsdk:"zone"`
	Parameters types.Object `tfsdk:"parameters"`
	// Specifies the list of user mapping rules.
	Rules types.List `tfsdk:"rules"`
	//additional attributes in terraform side
	TestMappingUsers       []UserMemberItem `tfsdk:"test_mapping_users"`
	TestMappingUserResults types.List       `tfsdk:"mapping_users"`
}

// V1MappingUsersRulesRule struct for V1MappingUsersRulesRule
type V1MappingUsersRulesRule struct {
	Operator types.String `tfsdk:"operator"`
	Options  types.Object `tfsdk:"options"`
	User1    types.Object `tfsdk:"target_user"`
	User2    types.Object `tfsdk:"source_user"`
}

// V1MappingUsersRulesRulesParameters Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information.
type V1MappingUsersRulesRulesParameters struct {
	DefaultUnixUser types.Object `tfsdk:"default_unix_user"`
}

// V1MappingUsersRulesRuleOptions Specifies the properties for user mapping rules.
type V1MappingUsersRulesRuleOptions struct {
	// If true, and the rule was applied successfully, stop processing further.
	Break       types.Bool   `tfsdk:"break"`
	DefaultUser types.Object `tfsdk:"default_user"`
	// If true, the primary GID and primary group SID should be copied to the existing credential.
	Group types.Bool `tfsdk:"group"`
	// If true, all additional identifiers should be copied to the existing credential.
	Groups types.Bool `tfsdk:"groups"`
	// If true, the primary UID and primary user SID should be copied to the existing credential.
	User types.Bool `tfsdk:"user"`
}

// V1MappingUsersRulesRuleUser2 struct for V1MappingUsersRulesRuleUser2
type V1MappingUsersRulesRuleUser2 struct {
	// Specifies the domain of the user that is being mapped.
	Domain types.String `tfsdk:"domain"`
	// Specifies the name of the user that is being mapped.
	User types.String `tfsdk:"user"`
}

// V1MappingUsersLookupMappingItem struct for V1MappingUsersLookupMappingItem
type V1MappingUsersLookupMappingItem struct {
	User                   types.Object `tfsdk:"user"`
	Privileges             types.List   `tfsdk:"privileges"`
	SupplementalIdentities types.List   `tfsdk:"supplemental_identities"`
	ZID                    types.Int64  `tfsdk:"zid"`
	Zone                   types.String `tfsdk:"zone"`
}

// V1MappingUsersLookupMappingItemGroup Specifies the configuration properties for a user.
type V1MappingUsersLookupMappingItemGroup struct {
	GID  types.String `tfsdk:"gid"`
	Name types.String `tfsdk:"name"`
	SID  types.String `tfsdk:"sid"`
}

// V1AuthIDNtokenPrivilegeItem Specifies the system-defined privilege that may be granted to users.
type V1AuthIDNtokenPrivilegeItem struct {
	// Specifies the ID of the privilege.
	ID types.String `tfsdk:"id"`
	// Specifies the name of the privilege.
	Name types.String `tfsdk:"name"`
	// True, if the privilege is read-only.
	ReadOnly types.Bool `tfsdk:"read_only"`
}

// V1MappingUsersLookupMappingItemUser Specifies the configuration properties for a user.
type V1MappingUsersLookupMappingItemUser struct {
	Name               types.String `tfsdk:"name"`
	OnDiskUserIdentity types.String `tfsdk:"on_disk_user_identity"`
	PrimaryGroupSID    types.String `tfsdk:"primary_group_sid"`
	PrimaryGroupName   types.String `tfsdk:"primary_group_name"`
	SID                types.String `tfsdk:"sid"`
	UID                types.String `tfsdk:"uid"`
}
