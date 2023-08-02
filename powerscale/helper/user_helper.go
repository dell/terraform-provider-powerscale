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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateUserDataSourceState updates datasource state.
func UpdateUserDataSourceState(userState *models.UserDataSourceModel, userResponse []powerscale.V1MappingUsersLookupMappingItemUser, roles []powerscale.V1AuthRoleExtended) {
	for _, user := range userResponse {
		var model models.UserModel
		updateUserState(&model, user)

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

// updateUserState updates user state.
func updateUserState(model *models.UserModel, user powerscale.V1MappingUsersLookupMappingItemUser) {

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

// UpdateUserResourceState updates resource state.
func UpdateUserResourceState(model *models.UserReourceModel, user powerscale.V1MappingUsersLookupMappingItemUser, roles []powerscale.V1AuthRoleExtended) {
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
	model.GID = AuthAccessKeyObjectMapper(user.Gid)
	model.PrimaryGroupSID = AuthAccessKeyObjectMapper(user.PrimaryGroupSid)

	if user.PrimaryGroupSid.Name != nil {
		model.PrimaryGroup = types.StringValue(*user.PrimaryGroupSid.Name)
	}
	if user.Sid.Id != nil {
		model.SID = types.StringValue(*user.Sid.Id)
	}
	if user.Uid.Id != nil && strings.HasPrefix(*user.Uid.Id, "UID:") {
		uidInt, _ := strconv.Atoi(strings.Trim(*user.Uid.Id, "UID:"))
		model.UID = types.Int64Value(int64(uidInt))
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

	if roles != nil {
		var roleAttrs []attr.Value
		for _, r := range roles {
			for _, m := range r.Members {
				if *m.Id == *user.Uid.Id {
					roleAttrs = append(roleAttrs, types.StringValue(r.Name))
				}
			}
		}

		model.Roles, _ = types.ListValue(types.StringType, roleAttrs)
	}
}

// GetAllRoles returns all roles
func GetAllRoles(ctx context.Context, client *client.Client) (roles []powerscale.V1AuthRoleExtended, err error) {
	roleParams := client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthRoles(ctx)
	result, _, err := roleParams.Execute()
	if err != nil {
		return
	}

	for {
		roles = append(roles, result.Roles...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}

		roleParams = client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthRoles(ctx).Resume(*result.Resume)
		if result, _, err = roleParams.Execute(); err != nil {
			return nil, err
		}
	}

	return
}

// GetUser Returns the User by userName.
func GetUser(ctx context.Context, client *client.Client, userName string) (*powerscale.V1AuthUsersExtended, error) {
	authID := fmt.Sprintf("USER:%s", userName)
	getParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthUser(ctx, userName)
	result, _, err := getParam.Execute()
	if err != nil {
		return nil, err
	}

	if len(result.Users) < 1 {
		return nil, fmt.Errorf("error getting the User - %s", authID)
	}
	return result, err
}

// CreateUser Creates an User
func CreateUser(ctx context.Context, client *client.Client, plan *models.UserReourceModel) error {

	createParam := client.PscaleOpenAPIClient.AuthApi.CreateAuthv1AuthUser(ctx)
	if !plan.QueryForce.IsNull() {
		createParam = createParam.Force(plan.QueryForce.ValueBool())
	}
	if !plan.QueryZone.IsNull() {
		createParam = createParam.Zone(plan.QueryZone.ValueString())
	}
	if !plan.QueryProvider.IsNull() {
		createParam = createParam.Provider(plan.QueryProvider.ValueString())
	}

	body := &powerscale.V1AuthUser{
		Name:                 plan.Name.ValueString(),
		Password:             plan.Password.ValueStringPointer(),
		Enabled:              plan.Enabled.ValueBoolPointer(),
		PromptPasswordChange: plan.PromptPasswordChange.ValueBoolPointer(),
		PasswordExpires:      plan.PasswordExpires.ValueBoolPointer(),
	}
	if !plan.Expiry.IsNull() && plan.Expiry.ValueInt64() > 0 {
		body.Expiry = plan.Expiry.ValueInt64Pointer()
	}
	if !plan.UID.IsNull() && plan.UID.ValueInt64() > 0 {
		body.Uid = plan.UID.ValueInt64Pointer()
	}
	if !plan.UnLock.IsNull() {
		body.Unlock = plan.UnLock.ValueBoolPointer()
	}
	if !plan.Email.IsNull() && plan.Email.ValueString() != "" {
		body.Email = plan.Email.ValueStringPointer()
	}
	if !plan.Gecos.IsNull() && plan.Gecos.ValueString() != "" {
		body.Gecos = plan.Gecos.ValueStringPointer()
	}
	if !plan.HomeDirectory.IsNull() && plan.HomeDirectory.ValueString() != "" {
		body.HomeDirectory = plan.HomeDirectory.ValueStringPointer()
	}
	if !plan.Shell.IsNull() && plan.Shell.ValueString() != "" {
		body.Shell = plan.Shell.ValueStringPointer()
	}
	if !plan.SID.IsNull() && plan.SID.ValueString() != "" {
		body.Sid = plan.SID.ValueStringPointer()
	}
	if !plan.PrimaryGroup.IsNull() && plan.PrimaryGroup.ValueString() != "" {
		primaryGroupName := plan.PrimaryGroup.ValueString()
		body.PrimaryGroup = &powerscale.V1AuthAccessAccessItemFileGroup{Name: &primaryGroupName}
	}

	createParam = createParam.V1AuthUser(*body)
	_, _, err := createParam.Execute()

	return err
}

// UpdateUser Updates an User parameters,
// including uid, sid, roles, password, enabled, home_directory, primary_group, unlock, email, expiry, gecos, shell, prompt_password_change, password_expires.
func UpdateUser(ctx context.Context, client *client.Client, state *models.UserReourceModel, plan *models.UserReourceModel) error {
	authID := fmt.Sprintf("USER:%s", plan.Name.ValueString())
	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv1AuthUser(ctx, authID)

	if !plan.QueryForce.IsNull() {
		updateParam = updateParam.Force(plan.QueryForce.ValueBool())
	}
	if !plan.QueryZone.IsNull() {
		updateParam = updateParam.Zone(plan.QueryZone.ValueString())
	}
	if !plan.QueryProvider.IsNull() {
		updateParam = updateParam.Provider(plan.QueryProvider.ValueString())
	}

	body := &powerscale.V1AuthUserExtended{}
	if !state.Name.Equal(plan.Name) {
		return fmt.Errorf("may not change user's name")
	}
	if !state.UID.Equal(plan.UID) && plan.UID.ValueInt64() > 0 {
		if !plan.QueryForce.ValueBool() {
			return fmt.Errorf("may not change user's uid without using the force option")
		}
		body.Uid = plan.UID.ValueInt64Pointer()
	}
	if !state.Expiry.Equal(plan.Expiry) && plan.Expiry.ValueInt64() > 0 {
		body.Expiry = plan.Expiry.ValueInt64Pointer()
	}
	if !state.Email.Equal(plan.Email) {
		body.Email = plan.Email.ValueStringPointer()
	}
	if !state.Enabled.Equal(plan.Enabled) {
		body.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if !state.PasswordExpires.Equal(plan.PasswordExpires) {
		body.PasswordExpires = plan.PasswordExpires.ValueBoolPointer()
	}
	if !state.PromptPasswordChange.Equal(plan.PromptPasswordChange) {
		body.PromptPasswordChange = plan.PromptPasswordChange.ValueBoolPointer()
	}
	if !state.Gecos.Equal(plan.Gecos) && plan.Gecos.ValueString() != "" {
		body.Gecos = plan.Gecos.ValueStringPointer()
	}
	if !state.HomeDirectory.Equal(plan.HomeDirectory) && plan.HomeDirectory.ValueString() != "" {
		body.HomeDirectory = plan.HomeDirectory.ValueStringPointer()
	}
	if !state.Password.Equal(plan.Password) {
		body.Password = plan.Password.ValueStringPointer()
	}
	if !state.Shell.Equal(plan.Shell) && plan.Shell.ValueString() != "" {
		body.Shell = plan.Shell.ValueStringPointer()
	}
	if !state.UnLock.Equal(plan.UnLock) {
		body.Unlock = plan.UnLock.ValueBoolPointer()
	}
	if !state.SID.Equal(plan.SID) && plan.SID.ValueString() != "" {
		body.Sid = plan.SID.ValueStringPointer()
	}
	if !state.PrimaryGroup.Equal(plan.PrimaryGroup) && plan.PrimaryGroup.ValueString() != "" {
		body.Sid = plan.SID.ValueStringPointer()
	}

	updateParam = updateParam.V1AuthUser(*body)
	_, err := updateParam.Execute()

	return err
}

// UpdateUserRoles Updates an User roles
func UpdateUserRoles(ctx context.Context, client *client.Client, state *models.UserReourceModel, plan *models.UserReourceModel) (diags diag.Diagnostics) {

	var duplicatedRoles []attr.Value
	for _, i := range state.Roles.Elements() {
		for _, j := range plan.Roles.Elements() {
			if i.Equal(j) {
				duplicatedRoles = append(duplicatedRoles, i)
			}
		}
	}

	// remove roles from user
	for _, i := range state.Roles.Elements() {
		for _, role := range duplicatedRoles {
			if role.Equal(i) {
				continue
			}
		}
		roleID := strings.Trim(i.String(), "\"")
		if err := RemoveUserRole(ctx, client, roleID, state.UID.ValueInt64()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove User from Role - %s", roleID), err.Error())
		}
	}

	// assign roles to user
	for _, i := range plan.Roles.Elements() {
		for _, role := range duplicatedRoles {
			if role.Equal(i) {
				continue
			}
		}
		roleID := strings.Trim(i.String(), "\"")
		if err := AddUserRole(ctx, client, roleID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error assign User to Role - %s", roleID), err.Error())
		}
	}

	return
}

// AddUserRole Assigns role to user
func AddUserRole(ctx context.Context, client *client.Client, roleID, userName string) error {
	authID := userName
	roleParam := client.PscaleOpenAPIClient.AuthRolesApi.CreateAuthRolesv1RoleMember(ctx, roleID).
		V1RoleMember(powerscale.V1AuthAccessAccessItemFileGroup{Name: &authID})
	if _, _, err := roleParam.Execute(); err != nil {
		return err
	}
	return nil
}

// RemoveUserRole Removes role from user
func RemoveUserRole(ctx context.Context, client *client.Client, roleID string, uid int64) error {
	authID := fmt.Sprintf("UID:%d", uid)
	roleParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1RolesRoleMember(ctx, authID, roleID)
	if _, err := roleParam.Execute(); err != nil {
		return err
	}
	return nil
}
