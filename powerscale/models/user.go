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

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	Users []UserModel  `tfsdk:"users"`
	ID    types.String `tfsdk:"id"`

	//filter
	Filter *UserFilterType `tfsdk:"filter"`
}

// UserModel holds user data source schema attribute details.
type UserModel struct {
	Dn                    types.String `tfsdk:"dn"`
	DNSDomain             types.String `tfsdk:"dns_domain"`
	Domain                types.String `tfsdk:"domain"`
	Email                 types.String `tfsdk:"email"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Expired               types.Bool   `tfsdk:"expired"`
	Expiry                types.Int64  `tfsdk:"expiry"`
	Gecos                 types.String `tfsdk:"gecos"`
	GeneratedGID          types.Bool   `tfsdk:"generated_gid"`
	GeneratedUID          types.Bool   `tfsdk:"generated_uid"`
	GeneratedUpn          types.Bool   `tfsdk:"generated_upn"`
	GID                   types.String `tfsdk:"gid"`
	HomeDirectory         types.String `tfsdk:"home_directory"`
	ID                    types.String `tfsdk:"id"`
	Locked                types.Bool   `tfsdk:"locked"`
	MaxPasswordAge        types.Int64  `tfsdk:"max_password_age"`
	Name                  types.String `tfsdk:"name"`
	PasswordExpired       types.Bool   `tfsdk:"password_expired"`
	PasswordExpires       types.Bool   `tfsdk:"password_expires"`
	PasswordExpiry        types.Int64  `tfsdk:"password_expiry"`
	PasswordLastSet       types.Int64  `tfsdk:"password_last_set"`
	PrimaryGroupSID       types.String `tfsdk:"primary_group_sid"`
	PromptPasswordChange  types.Bool   `tfsdk:"prompt_password_change"`
	Provider              types.String `tfsdk:"provider"`
	SamAccountName        types.String `tfsdk:"sam_account_name"`
	Shell                 types.String `tfsdk:"shell"`
	SID                   types.String `tfsdk:"sid"`
	Type                  types.String `tfsdk:"type"`
	UID                   types.String `tfsdk:"uid"`
	Upn                   types.String `tfsdk:"upn"`
	UserCanChangePassword types.Bool   `tfsdk:"user_can_change_password"`
	Roles                 types.List   `tfsdk:"roles"`
}

// UserFilterType holds filter attribute for user.
type UserFilterType struct {
	Names        []UserMemberItem `tfsdk:"names"`
	NamePrefix   types.String     `tfsdk:"name_prefix"`
	Domain       types.String     `tfsdk:"domain"`
	Zone         types.String     `tfsdk:"zone"`
	Provider     types.String     `tfsdk:"provider"`
	Cached       types.Bool       `tfsdk:"cached"`
	ResolveNames types.Bool       `tfsdk:"resolve_names"`
	MemberOf     types.Bool       `tfsdk:"member_of"`
}

// UserMemberItem holds identity attribute for a auth member.
type UserMemberItem struct {
	Name types.String `tfsdk:"name"`
	UID  types.Int64  `tfsdk:"uid"`
}
