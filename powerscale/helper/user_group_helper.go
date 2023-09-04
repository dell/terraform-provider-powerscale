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

// UpdateUserGroupDataSourceState updates datasource state.
func UpdateUserGroupDataSourceState(model *models.UserGroupModel, groupResponse powerscale.V1AuthGroupExtended,
	groupMembers []powerscale.V1AuthAccessAccessItemFileGroup, roles []powerscale.V1AuthRoleExtended) {
	model.Dn = types.StringValue(groupResponse.Dn)
	model.Domain = types.StringValue(groupResponse.Domain)
	model.DNSDomain = types.StringValue(groupResponse.DnsDomain)
	model.ID = types.StringValue(groupResponse.Id)
	model.Name = types.StringValue(groupResponse.Name)
	model.Provider = types.StringValue(groupResponse.Provider)
	model.SamAccountName = types.StringValue(groupResponse.SamAccountName)
	model.Type = types.StringValue(groupResponse.Type)
	if groupResponse.Gid.Id != nil {
		model.GID = types.StringValue(*groupResponse.Gid.Id)
	}
	if groupResponse.Sid.Id != nil {
		model.SID = types.StringValue(*groupResponse.Sid.Id)
	}
	model.GeneratedGID = types.BoolValue(groupResponse.GeneratedGid)

	//parse roles
	var roleAttrs []attr.Value
	for _, r := range roles {
		for _, m := range r.Members {
			if *m.Id == *groupResponse.Gid.Id {
				roleAttrs = append(roleAttrs, types.StringValue(r.Name))
			}
		}
	}
	model.Roles, _ = types.ListValue(types.StringType, roleAttrs)

	// parse group members
	var members []models.V1AuthAccessAccessItemFileGroup
	for _, m := range groupMembers {
		members = append(members, models.V1AuthAccessAccessItemFileGroup{
			Name: types.StringValue(*m.Name),
			ID:   types.StringValue(*m.Id),
			Type: types.StringValue(*m.Type),
		})
	}
	model.Members = members
}

// GetAllGroupMembers returns all group members.
func GetAllGroupMembers(ctx context.Context, client *client.Client, groupName string) (members []powerscale.V1AuthAccessAccessItemFileGroup, err error) {
	members = make([]powerscale.V1AuthAccessAccessItemFileGroup, 0)
	emptyMembers := make([]powerscale.V1AuthAccessAccessItemFileGroup, 0)
	memberParams := client.PscaleOpenAPIClient.AuthGroupsApi.ListAuthGroupsv1GroupMembers(ctx, groupName)
	result, _, err := memberParams.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return emptyMembers, fmt.Errorf("error getting user group members: %s", message)
	}

	for {
		members = append(members, result.Members...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}
		memberParams = client.PscaleOpenAPIClient.AuthGroupsApi.ListAuthGroupsv1GroupMembers(ctx, groupName).Resume(*result.Resume)
		if result, _, err = memberParams.Execute(); err != nil {
			errStr := constants.ReadUserGroupMemberErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return emptyMembers, fmt.Errorf("error getting user group members with resume: %s", message)
		}
	}

	return
}

// GetUserGroupsWithFilter returns all filtered groups.
func GetUserGroupsWithFilter(ctx context.Context, client *client.Client, filter *models.UserGroupFilterType) (groups []powerscale.V1AuthGroupExtended, err error) {
	groupParams := client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthGroups(ctx)
	if filter != nil {
		if !filter.NamePrefix.IsNull() {
			groupParams = groupParams.Filter(filter.NamePrefix.ValueString())
		}
		if !filter.Domain.IsNull() {
			groupParams = groupParams.Domain(filter.Domain.ValueString())
		}
		if !filter.Zone.IsNull() {
			groupParams = groupParams.Zone(filter.Zone.ValueString())
		}
		if !filter.Provider.IsNull() {
			groupParams = groupParams.Provider(filter.Provider.ValueString())
		}
		if !filter.Cached.IsNull() {
			groupParams = groupParams.Cached(filter.Cached.ValueBool())
		}
	}

	result, _, err := groupParams.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user groups: %s", message)
	}

	for {
		groups = append(groups, result.Groups...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}
		groupParams = client.PscaleOpenAPIClient.AuthApi.ListAuthv1AuthGroups(ctx).Resume(*result.Resume)
		if result, _, err = groupParams.Execute(); err != nil {
			errStr := constants.ReadUserGroupErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting user groups with resume: %s", message)
		}
	}

	if filter != nil && len(filter.Names) > 0 {
		var validUserGroups []string
		var filteredUserGroups []powerscale.V1AuthGroupExtended

		for _, group := range groups {
			for _, name := range filter.Names {
				if (!name.Name.IsNull() && group.Name == name.Name.ValueString()) ||
					(!name.GID.IsNull() && fmt.Sprintf("GID:%d", name.GID.ValueInt64()) == *group.Gid.Id) {
					filteredUserGroups = append(filteredUserGroups, group)
					validUserGroups = append(validUserGroups, fmt.Sprintf("Name: %s, GID: %s", group.Name, *group.Gid.Id))
					break
				}
			}
		}

		if len(filteredUserGroups) != len(filter.Names) {
			return nil, fmt.Errorf(
				"error one or more of the filtered user group names is not a valid powerscale user group. Valid user groups: [%v], filtered list: [%v]",
				validUserGroups, filter.Names)
		}

		groups = filteredUserGroups
	}
	return
}

// GetUserGroup Returns the User Group by user group name.
func GetUserGroup(ctx context.Context, client *client.Client, groupName string) (*powerscale.V1AuthGroupsExtended, error) {
	authID := fmt.Sprintf("GROUP:%s", groupName)
	getParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthGroup(ctx, groupName)
	result, _, err := getParam.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user group - %s : %s", authID, message)
	}
	if len(result.Groups) < 1 {
		message := constants.ReadUserGroupErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty user group - %s : %s", authID, message)
	}

	return result, err
}

// UpdateUserGroupResourceState updates resource state.
func UpdateUserGroupResourceState(model *models.UserGroupResourceModel, group powerscale.V1AuthGroupExtended,
	groupMembers []powerscale.V1AuthAccessAccessItemFileGroup, roles []powerscale.V1AuthRoleExtended) {
	model.Dn = types.StringValue(group.Dn)
	model.Domain = types.StringValue(group.Domain)
	model.DNSDomain = types.StringValue(group.DnsDomain)
	model.ID = types.StringValue(group.Id)
	model.Name = types.StringValue(group.Name)
	model.Provider = types.StringValue(group.Provider)
	model.SamAccountName = types.StringValue(group.SamAccountName)
	model.Type = types.StringValue(group.Type)
	model.GeneratedGID = types.BoolValue(group.GeneratedGid)

	if group.Sid.Id != nil {
		model.SID = types.StringValue(*group.Sid.Id)
	}
	if group.Gid.Id != nil && strings.HasPrefix(*group.Gid.Id, "GID:") {
		gidInt, _ := strconv.Atoi(strings.Trim(*group.Gid.Id, "GID:"))
		model.GID = types.Int64Value(int64(gidInt))
	}

	if roles != nil {
		var roleAttrs []attr.Value
		for _, r := range roles {
			for _, m := range r.Members {
				if *m.Id == *group.Gid.Id {
					roleAttrs = append(roleAttrs, types.StringValue(r.Name))
				}
			}
		}

		model.Roles, _ = types.ListValue(types.StringType, roleAttrs)
	}
	if groupMembers != nil {
		var users []attr.Value
		var groups []attr.Value
		var wellKnowns []attr.Value
		for _, m := range groupMembers {
			switch *m.Type {
			case "user":
				users = append(users, types.StringValue(*m.Name))
			case "group":
				groups = append(groups, types.StringValue(*m.Name))
			case "wellknown":
				wellKnowns = append(wellKnowns, types.StringValue(*m.Name))
			}
		}
		model.Users, _ = types.ListValue(types.StringType, users)
		model.Groups, _ = types.ListValue(types.StringType, groups)
		model.WellKnowns, _ = types.ListValue(types.StringType, wellKnowns)
	}

}

// CreateUserGroup Creates a User Group.
func CreateUserGroup(ctx context.Context, client *client.Client, plan *models.UserGroupResourceModel) error {

	createParam := client.PscaleOpenAPIClient.AuthApi.CreateAuthv1AuthGroup(ctx)
	if !plan.QueryForce.IsNull() {
		createParam = createParam.Force(plan.QueryForce.ValueBool())
	}
	if !plan.QueryZone.IsNull() {
		createParam = createParam.Zone(plan.QueryZone.ValueString())
	}
	if !plan.QueryProvider.IsNull() {
		createParam = createParam.Provider(plan.QueryProvider.ValueString())
	}

	body := &powerscale.V1AuthGroup{Name: plan.Name.ValueString()}
	if !plan.GID.IsNull() && plan.GID.ValueInt64() > 0 {
		body.Gid = plan.GID.ValueInt64Pointer()
	}
	if !plan.SID.IsNull() && plan.SID.ValueString() != "" {
		body.Sid = plan.SID.ValueStringPointer()
	}

	createParam = createParam.V1AuthGroup(*body)
	if _, _, err := createParam.Execute(); err != nil {
		errStr := constants.CreateUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error creating user group: %s", message)
	}

	return nil
}

// UpdateUserGroup Updates a User Group GID.
func UpdateUserGroup(ctx context.Context, client *client.Client, state *models.UserGroupResourceModel, plan *models.UserGroupResourceModel) error {
	authID := fmt.Sprintf("GROUP:%s", plan.Name.ValueString())
	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv1AuthGroup(ctx, authID)

	if !plan.QueryForce.IsNull() {
		updateParam = updateParam.Force(plan.QueryForce.ValueBool())
	}
	if !plan.QueryZone.IsNull() {
		updateParam = updateParam.Zone(plan.QueryZone.ValueString())
	}
	if !plan.QueryProvider.IsNull() {
		updateParam = updateParam.Provider(plan.QueryProvider.ValueString())
	}

	body := &powerscale.V1AuthGroupExtendedExtended{}
	if !state.Name.Equal(plan.Name) || plan.SID.ValueString() != "" && !state.SID.Equal(plan.SID) {
		return fmt.Errorf("may not change user group's name or sid")
	}
	if !state.GID.Equal(plan.GID) && plan.GID.ValueInt64() > 0 {
		if !plan.QueryForce.ValueBool() {
			return fmt.Errorf("may not change user group's gid without using the force option")
		}
		body.Gid = plan.GID.ValueInt64Pointer()
	}

	updateParam = updateParam.V1AuthGroup(*body)
	if _, err := updateParam.Execute(); err != nil {
		errStr := constants.UpdateUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error updating user group - %s : %s", authID, message)
	}

	return nil
}

// UpdateUserGroupRoles Updates a User Group roles.
func UpdateUserGroupRoles(ctx context.Context, client *client.Client, state *models.UserGroupResourceModel, plan *models.UserGroupResourceModel) (diags diag.Diagnostics) {

	// get roles list changes
	toAdd, toRemove := GetElementsChanges(state.Roles.Elements(), plan.Roles.Elements())

	// if gid changed, should remove all roles firstly, then assign all roles.
	if !plan.GID.Equal(state.GID) {
		toAdd = plan.Roles.Elements()
		toRemove = state.Roles.Elements()
	}

	// remove roles from user group
	for _, i := range toRemove {
		roleID := strings.Trim(i.String(), "\"")
		if err := RemoveUserGroupRole(ctx, client, roleID, state.GID.ValueInt64()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove User Group from Role - %s", roleID), err.Error())
		}
	}

	// assign roles to user group
	for _, i := range toAdd {
		roleID := strings.Trim(i.String(), "\"")
		if err := AddUserGroupRole(ctx, client, roleID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error assign User Group to Role - %s", roleID), err.Error())
		}
	}

	return
}

// AddUserGroupRole Assigns role to user group.
func AddUserGroupRole(ctx context.Context, client *client.Client, roleID, userGroupName string) error {
	authID := userGroupName
	roleParam := client.PscaleOpenAPIClient.AuthRolesApi.CreateAuthRolesv1RoleMember(ctx, roleID).
		V1RoleMember(powerscale.V1AuthAccessAccessItemFileGroup{Name: &authID})
	if _, _, err := roleParam.Execute(); err != nil {
		errStr := constants.AddRoleMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error assign user group - %s to role - %s: %s", authID, roleID, message)
	}
	return nil
}

// RemoveUserGroupRole Removes role from user group.
func RemoveUserGroupRole(ctx context.Context, client *client.Client, roleID string, gid int64) error {
	authID := fmt.Sprintf("GID:%d", gid)
	roleParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1RolesRoleMember(ctx, authID, roleID)
	if _, err := roleParam.Execute(); err != nil {
		errStr := constants.DeleteRoleMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error remove user group - %s from role - %s: %s", authID, roleID, message)
	}
	return nil
}

// UpdateUserGroupMembers Updates a User Group members.
func UpdateUserGroupMembers(ctx context.Context, client *client.Client, state *models.UserGroupResourceModel, plan *models.UserGroupResourceModel) (diags diag.Diagnostics) {

	// update users in members
	toAdd, toRemove := GetElementsChanges(state.Users.Elements(), plan.Users.Elements())
	// remove users from user group by memberAuthID
	for _, i := range toRemove {
		memberAuthID := fmt.Sprintf("USER:%s", strings.Trim(i.String(), "\""))
		if err := RemoveUserGroupMember(ctx, client, memberAuthID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove User - %s from User Group.", memberAuthID), err.Error())
		}
	}
	// add users to user group by memberID
	for _, i := range toAdd {
		memberID := strings.Trim(i.String(), "\"")
		memberType := "user"
		memberIdentity := powerscale.V1AuthAccessAccessItemFileGroup{Name: &memberID, Type: &memberType}
		if err := AddUserGroupMember(ctx, client, memberIdentity, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error add User - %s to User Group.", memberID), err.Error())
		}
	}

	// update groups in members
	toAdd, toRemove = GetElementsChanges(state.Groups.Elements(), plan.Groups.Elements())
	// remove groups from user group by memberAuthID
	for _, i := range toRemove {
		memberAuthID := fmt.Sprintf("GROUP:%s", strings.Trim(i.String(), "\""))
		if err := RemoveUserGroupMember(ctx, client, memberAuthID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove Group - %s from User Group.", memberAuthID), err.Error())
		}
	}
	// add groups to user group by memberID
	for _, i := range toAdd {
		memberID := strings.Trim(i.String(), "\"")
		memberType := "group"
		memberIdentity := powerscale.V1AuthAccessAccessItemFileGroup{Name: &memberID, Type: &memberType}
		if err := AddUserGroupMember(ctx, client, memberIdentity, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error add Group - %s to User Group.", memberID), err.Error())
		}
	}

	// update well_konwns in members
	toAdd, toRemove = GetElementsChanges(state.WellKnowns.Elements(), plan.WellKnowns.Elements())
	// remove well-knowns from user group by wellKnownSID
	for _, i := range toRemove {
		wellKnownName := getWellKnownName(strings.Trim(i.String(), "\""))
		wellKnownSID := ""
		if wellKnowns, err := GetWellKnown(ctx, client, wellKnownName); err == nil {
			wellKnownSID = *wellKnowns.Wellknowns[0].Id
		} else {
			diags.AddError(fmt.Sprintf("Error remove Well-Known from User Group. Not found Well-Known - %s.", wellKnownName), err.Error())
			continue
		}
		if err := RemoveUserGroupMember(ctx, client, wellKnownSID, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error remove Well-Known - %s from User Group.", wellKnownName), err.Error())
		}
	}
	// add well-knowns to user group by wellKnownName
	for _, i := range toAdd {
		wellKnownName := getWellKnownName(strings.Trim(i.String(), "\""))
		memberType := "wellknown"
		memberIdentity := powerscale.V1AuthAccessAccessItemFileGroup{Name: &wellKnownName, Type: &memberType}
		if err := AddUserGroupMember(ctx, client, memberIdentity, plan.Name.ValueString()); err != nil {
			diags.AddError(fmt.Sprintf("Error add Well-Known - %s to User Group.", wellKnownName), err.Error())
		}
	}

	return
}

// RemoveUserGroupMember Removes member from user group by memberAuthID, like GROUP:groupName, USER:userName and SID:well-known-sid.
func RemoveUserGroupMember(ctx context.Context, client *client.Client, memberAuthID, userGroupName string) error {
	memberParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1GroupsGroupMember(ctx, memberAuthID, userGroupName)
	if _, err := memberParam.Execute(); err != nil {
		errStr := constants.DeleteUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error remove member - %s from user group - %s: %s", memberAuthID, userGroupName, message)
	}
	return nil
}

// AddUserGroupMember Adds member to user group by memberIdentity.
func AddUserGroupMember(ctx context.Context, client *client.Client, memberIdentity powerscale.V1AuthAccessAccessItemFileGroup, userGroupName string) error {
	memberParam := client.PscaleOpenAPIClient.AuthGroupsApi.CreateAuthGroupsv1GroupMember(ctx, userGroupName).V1GroupMember(memberIdentity)
	if _, _, err := memberParam.Execute(); err != nil {
		errStr := constants.AddUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error add member - %s:%s to user group - %s: %s", *memberIdentity.Type, *memberIdentity.Name, userGroupName, message)
	}
	return nil
}

// DeleteUserGroup Deletes a User Group.
func DeleteUserGroup(ctx context.Context, client *client.Client, groupName string) error {
	authID := fmt.Sprintf("GROUP:%s", groupName)
	deleteParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv1AuthGroup(ctx, authID)
	if _, err := deleteParam.Execute(); err != nil {
		errStr := constants.DeleteUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting user group - %s : %s", authID, message)
	}

	return nil
}

// GetWellKnown Returns the well-known by well-known name.
func GetWellKnown(ctx context.Context, client *client.Client, wellKnownName string) (*powerscale.V1AuthWellknowns, error) {
	result, _, err := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthWellknown(ctx, wellKnownName).Execute()
	if err != nil {
		errStr := constants.ReadWellKnownErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting well known - %s : %s", wellKnownName, message)
	}
	if len(result.Wellknowns) < 1 {
		message := constants.ReadWellKnownErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty well known - %s : %s", wellKnownName, message)
	}

	return result, err
}

// getWellKnownName returns Well-Known suffix name. 'NT AUTHORITY\\DIALUP' will return 'DIALUP'.
func getWellKnownName(name string) string {
	groups := strings.Split(name, "\\")
	return groups[len(groups)-1]
}
