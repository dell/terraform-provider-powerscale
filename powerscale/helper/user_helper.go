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
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateUserDataSourceState updates datasource state.
func UpdateUserDataSourceState(userState *models.UserDataSourceModel, userResponse []powerscale.V1AuthUserExtended, roles []powerscale.V1AuthRoleExtended) {
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
func updateUserState(model *models.UserModel, user powerscale.V1AuthUserExtended) {

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

	model.Expiry = types.Int64Value(user.Expiry)
	model.MaxPasswordAge = types.Int64Value(int64(user.MaxPasswordAge))
	model.PasswordExpiry = types.Int64Value(int64(user.PasswordExpiry))
	model.PasswordLastSet = types.Int64Value(int64(user.PasswordLastSet))
}

// UpdateUserResourceState updates resource state.
func UpdateUserResourceState(model *models.UserResourceModel, user powerscale.V1AuthUserExtended, roles []powerscale.V1AuthRoleExtended) {
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

	model.Expiry = types.Int64Value(user.Expiry)
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

// GetUsersWithFilter returns all filtered users.
func GetUsersWithFilter(ctx context.Context, client *client.Client, filter *models.UserFilterType) (users []powerscale.V1AuthUserExtended, err error) {
	userParams := client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthUsers(ctx)

	if filter != nil {
		if !filter.NamePrefix.IsNull() {
			userParams = userParams.Filter(filter.NamePrefix.ValueString())
		}
		if !filter.Domain.IsNull() {
			userParams = userParams.Domain(filter.Domain.ValueString())
		}
		if !filter.Zone.IsNull() {
			userParams = userParams.Zone(filter.Zone.ValueString())
		}
		if !filter.Provider.IsNull() {
			userParams = userParams.Provider(filter.Provider.ValueString())
		}
		if !filter.Cached.IsNull() {
			userParams = userParams.Cached(filter.Cached.ValueBool())
		}
		if !filter.MemberOf.IsNull() {
			userParams = userParams.QueryMemberOf(filter.MemberOf.ValueBool())
		}
	}

	result, _, err := userParams.Execute()
	if err != nil {
		errStr := constants.ReadUserErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting users: %s", message)
	}

	for {
		users = append(users, result.Users...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}

		userParams = client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthUsers(ctx).Resume(*result.Resume)
		if result, _, err = userParams.Execute(); err != nil {
			errStr := constants.ReadUserErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting users with resume: %s", message)
		}
	}
	return
}

// GetAllRolesWithZone returns all roles in specific zone.
func GetAllRolesWithZone(ctx context.Context, client *client.Client, zone string) (roles []powerscale.V1AuthRoleExtended, err error) {
	roles = make([]powerscale.V1AuthRoleExtended, 0)
	emptyRoles := make([]powerscale.V1AuthRoleExtended, 0)

	roleParams := client.PscaleOpenAPIClient.AuthApi.ListAuthv7AuthRoles(ctx)
	if zone != "" {
		roleParams = roleParams.Zone(zone)
	}

	result, _, err := roleParams.Execute()
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return emptyRoles, fmt.Errorf("error getting roles : %s", message)
	}

	for {
		roles = append(roles, result.Roles...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}

		roleParams = client.PscaleOpenAPIClient.AuthApi.ListAuthv7AuthRoles(ctx).Resume(*result.Resume)
		if result, _, err = roleParams.Execute(); err != nil {
			errStr := constants.ReadRoleErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return emptyRoles, fmt.Errorf("error getting roles with resume: %s", message)
		}
	}

	return
}

// GetUserWithZone Returns the User by userName in specific zone.
func GetUserWithZone(ctx context.Context, client *client.Client, userName, zone string) (*powerscale.V1AuthUsersExtended, error) {
	authID := fmt.Sprintf("USER:%s", userName)
	getParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthUser(ctx, userName)
	if zone != "" {
		getParam = getParam.Zone(zone)
	}
	result, _, err := getParam.Execute()
	if err != nil {
		errStr := constants.ReadUserErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user - %s : %s", authID, message)
	}

	if len(result.Users) < 1 {
		message := constants.ReadUserErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty user - %s : %s", authID, message)
	}

	return result, err
}

// CreateUser Creates a User.
func CreateUser(ctx context.Context, client *client.Client, plan *models.UserResourceModel) error {

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
	if _, _, err := createParam.Execute(); err != nil {
		errStr := constants.CreateUserErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error creating user: %s", message)
	}

	return nil
}

// UpdateUser Updates a User parameters,
// including uid, sid, roles, password, enabled, home_directory, primary_group, unlock, email, expiry, gecos, shell, prompt_password_change, password_expires.
func UpdateUser(ctx context.Context, client *client.Client, state *models.UserResourceModel, plan *models.UserResourceModel) error {
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

	body := &powerscale.V1AuthUserExtendedExtended{}
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
		primaryGroupName := plan.PrimaryGroup.ValueString()
		body.PrimaryGroup = &powerscale.V1AuthAccessAccessItemFileGroup{Name: &primaryGroupName}
	}

	updateParam = updateParam.V1AuthUser(*body)
	if _, err := updateParam.Execute(); err != nil {
		errStr := constants.UpdateUserErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error updating user - %s : %s", authID, message)
	}

	return nil
}

// DeleteUserWithZone Deletes a User in specific zone.
func DeleteUserWithZone(ctx context.Context, client *client.Client, userName, zone string) error {
	authID := fmt.Sprintf("USER:%s", userName)
	deleteParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1AuthUser(ctx, authID)
	if zone != "" {
		deleteParam = deleteParam.Zone(zone)
	}
	if _, err := deleteParam.Execute(); err != nil {
		errStr := constants.DeleteUserErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting user - %s : %s", authID, message)
	}
	return nil
}

// UpdateUserRoles Updates a User roles.
func UpdateUserRoles(ctx context.Context, client *client.Client, state *models.UserResourceModel, plan *models.UserResourceModel) (diags diag.Diagnostics) {

	// get roles list changes
	toAdd, toRemove := GetElementsChanges(state.Roles.Elements(), plan.Roles.Elements())

	// if uid updated, should remove all roles firstly, then assign all roles.
	if !plan.UID.Equal(state.UID) {
		toAdd = plan.Roles.Elements()
		toRemove = state.Roles.Elements()
	}

	// remove roles from user
	for _, i := range toRemove {
		roleID := strings.Trim(i.String(), "\"")
		if err := RemoveUserRoleWithZone(ctx, client, roleID, state.UID.ValueInt64(), state.QueryZone.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove User from Role - %s", roleID), err.Error())
		}
	}

	// assign roles to user
	for _, i := range toAdd {
		roleID := strings.Trim(i.String(), "\"")
		if err := AddUserRoleWithZone(ctx, client, roleID, plan.Name.ValueString(), plan.QueryZone.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error assign User to Role - %s", roleID), err.Error())
		}
	}

	return
}

// AddUserRoleWithZone Assigns role to user in specific zone.
func AddUserRoleWithZone(ctx context.Context, client *client.Client, roleID, userName, zone string) error {
	authID := userName
	roleParam := client.PscaleOpenAPIClient.AuthRolesApi.CreateAuthRolesv7RoleMember(ctx, roleID).
		V7RoleMember(powerscale.V1AuthAccessAccessItemFileGroup{Name: &authID})
	if zone != "" {
		roleParam = roleParam.Zone(zone)
	}
	if _, _, err := roleParam.Execute(); err != nil {
		errStr := constants.AddRoleMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error assign user - %s to role - %s: %s", authID, roleID, message)
	}
	return nil
}

// RemoveUserRoleWithZone Removes role from user in specific zone.
func RemoveUserRoleWithZone(ctx context.Context, client *client.Client, roleID string, uid int64, zone string) error {
	authID := fmt.Sprintf("UID:%d", uid)
	roleParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv7RolesRoleMember(ctx, authID, roleID)
	if zone != "" {
		roleParam = roleParam.Zone(zone)
	}
	if _, err := roleParam.Execute(); err != nil {
		errStr := constants.DeleteRoleMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error remove user - %s from role - %s: %s", authID, roleID, message)
	}
	return nil
}
