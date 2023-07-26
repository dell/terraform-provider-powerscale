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
	powerscale "dell/powerscale-go-client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateUserDataSourceState updates datasource state.
func UpdateUserDataSourceState(userState *models.UserDataSourceModel, userResponse []powerscale.V1MappingUsersLookupMappingItemUser, roles []powerscale.V1AuthRoleExtended) {
	for _, user := range userResponse {
		var model models.UserModel
		UpdateUserState(&model, user)

		var roleAttrs []attr.Value
		for _, r := range roles {
			for _, m := range r.Members {
				if *m.Id == *user.Uid.Id {
					roleAttrs = append(roleAttrs, types.StringValue(r.Name))
				}
			}
		}
		model.Roles, _ = types.ListValue(types.StringType, roleAttrs)
		userState.Users = append(userState.Users, model)
	}
}

// UpdateUserState updates user state.
func UpdateUserState(model *models.UserModel, user powerscale.V1MappingUsersLookupMappingItemUser) {

	model.Dn = types.StringValue(user.Dn)
	model.DNSDomain = types.StringValue(user.DnsDomain)
	model.Domain = types.StringValue(user.Domain)
	model.Email = types.StringValue(user.Email)
	model.Gecos = types.StringValue(user.Gecos)
	model.HomeDirectory = types.StringValue(user.HomeDirectory)
	model.ID = types.StringValue(user.Id)
	model.Name = types.StringValue(user.Name)
	model.Provider = types.StringValue(user.Provider)
	model.SamAccountName = types.StringValue(user.SamAccountName)
	model.Shell = types.StringValue(user.Shell)
	model.Type = types.StringValue(user.Type)
	model.Upn = types.StringValue(user.Upn)
	if user.Gid.Id != nil {
		model.GID = types.StringValue(*user.Gid.Id)
	}
	if user.PrimaryGroupSid.Id != nil {
		model.PrimaryGroupSID = types.StringValue(*user.PrimaryGroupSid.Id)
	}
	if user.Sid.Id != nil {
		model.SID = types.StringValue(*user.Sid.Id)
	}
	if user.Uid.Id != nil {
		model.UID = types.StringValue(*user.Uid.Id)
	}

	model.Enabled = types.BoolValue(user.Enabled)
	model.Expired = types.BoolValue(user.Expired)
	model.GeneratedGID = types.BoolValue(user.GeneratedGid)
	model.GeneratedUID = types.BoolValue(user.GeneratedUid)
	model.GeneratedUpn = types.BoolValue(user.GeneratedUpn)
	model.Locked = types.BoolValue(user.Locked)
	model.PasswordExpired = types.BoolValue(user.PasswordExpired)
	model.PasswordExpires = types.BoolValue(user.PasswordExpires)
	model.PromptPasswordChange = types.BoolValue(user.PromptPasswordChange)
	model.UserCanChangePassword = types.BoolValue(user.UserCanChangePassword)

	model.Expiry = types.Int64Value(int64(user.Expiry))
	model.MaxPasswordAge = types.Int64Value(int64(user.MaxPasswordAge))
	model.PasswordExpiry = types.Int64Value(int64(user.PasswordExpiry))
	model.PasswordLastSet = types.Int64Value(int64(user.PasswordLastSet))
}
