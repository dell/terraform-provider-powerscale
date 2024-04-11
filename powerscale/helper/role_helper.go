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

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// GetRoles Get a list of Roles.
func GetRoles(ctx context.Context, client *client.Client, state models.RoleDataSourceModel) (*powerscale.V14AuthRoles, error) {
	roleParams := client.PscaleOpenAPIClient.AuthApi.ListAuthv14AuthRoles(ctx)

	if state.RoleFilter != nil && !state.RoleFilter.Zone.IsNull() {
		roleParams = roleParams.Zone(state.RoleFilter.Zone.ValueString())
	}

	roles, _, err := roleParams.Execute()
	return roles, err
}

// RoleDetailMapper Does the mapping from response to model.
//
//go:noinline
func RoleDetailMapper(ctx context.Context, role *powerscale.V14AuthRoleExtended) (models.RoleDetailModel, error) {
	model := models.RoleDetailModel{}
	err := CopyFields(ctx, role, &model)
	return model, err
}

// CreateRole Create a Role.
func CreateRole(ctx context.Context, client *client.Client, role powerscale.V14AuthRole, roleModel models.RoleResourceModel) (*powerscale.CreateResponse, error) {
	createParam := client.PscaleOpenAPIClient.AuthApi.CreateAuthv14AuthRole(ctx).V14AuthRole(role)
	if !roleModel.Zone.IsNull() {
		createParam = createParam.Zone(roleModel.Zone.ValueString())
	}
	roleID, _, err := createParam.Execute()
	return roleID, err
}

// GetRole retrieve Role information.
func GetRole(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel) (*powerscale.V14AuthRolesExtended, error) {
	queryParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv14AuthRole(ctx, roleModel.ID.ValueString())
	if !roleModel.Zone.IsNull() {
		queryParam = queryParam.Zone(roleModel.Zone.ValueString())
	}
	roleRes, _, err := queryParam.Execute()
	return roleRes, err
}

// UpdateRole Update a Role.
func UpdateRole(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel, roleToUpdate powerscale.V14AuthRoleExtendedExtended) error {
	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv14AuthRole(ctx, roleModel.ID.ValueString())
	if !roleModel.Zone.IsNull() {
		updateParam = updateParam.Zone(roleModel.Zone.ValueString())
	}
	_, err := updateParam.V14AuthRole(roleToUpdate).Execute()
	return err
}

// DeleteRole Delete a Role.
func DeleteRole(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel) error {
	deleteParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv14AuthRole(ctx, roleModel.ID.ValueString())
	if !roleModel.Zone.IsNull() {
		deleteParam = deleteParam.Zone(roleModel.Zone.ValueString())
	}
	_, err := deleteParam.Execute()
	return err
}

// ReorderRoleMembers Reorder role members to ensure the consistency of state.
func ReorderRoleMembers(localMembers types.List, remoteMembers types.List) (types.List, error) {
	localMembersList := localMembers.Elements()
	remoteMembersList := remoteMembers.Elements()
	var orderedMembers []attr.Value

	addedRemoteMembers := make([]bool, len(remoteMembersList))
	for _, localMember := range localMembersList {
		localMemberObj, ok := localMember.(basetypes.ObjectValue)
		if !ok || localMemberObj.IsNull() || localMemberObj.IsUnknown() {
			return types.List{}, errors.New("failed to reorder role members")
		}
		localMemberName := localMemberObj.Attributes()["name"]
		localMemberID := localMemberObj.Attributes()["id"]
		localMemberType := localMemberObj.Attributes()["type"]

		for i, remoteMember := range remoteMembersList {
			remoteMemberObj, ok := remoteMember.(basetypes.ObjectValue)
			if !ok || remoteMemberObj.IsNull() || remoteMemberObj.IsUnknown() {
				return types.List{}, errors.New("failed to reorder role members")
			}
			remoteMemberName := remoteMemberObj.Attributes()["name"]
			remoteMemberID := remoteMemberObj.Attributes()["id"]
			remoteMemberType := remoteMemberObj.Attributes()["type"]

			if localMemberID == remoteMemberID || (localMemberType == remoteMemberType && localMemberName == remoteMemberName) || (localMemberType.IsUnknown() || localMemberType.IsNull()) && remoteMemberType.String() == "\"user\"" && localMemberName == remoteMemberName {
				orderedMembers = append(orderedMembers, remoteMember)
				addedRemoteMembers[i] = true
				break
			}
		}
	}

	for i, remoteMember := range remoteMembersList {
		if !addedRemoteMembers[i] {
			orderedMembers = append(orderedMembers, remoteMember)
		}
	}

	roleMembersType := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	}

	orderedMemberList, _ := types.ListValue(types.ObjectType{AttrTypes: roleMembersType}, orderedMembers)
	return orderedMemberList, nil
}

// ReorderRolePrivileges Reorder role privileges to ensure the consistency of state.
func ReorderRolePrivileges(localPrivileges types.List, remotePrivileges types.List) (types.List, error) {
	localPrivilegesList := localPrivileges.Elements()
	remotePrivilegesList := remotePrivileges.Elements()
	var orderedPrivileges []attr.Value

	addedRemotePrivileges := make([]bool, len(remotePrivilegesList))
	for _, localPrivilege := range localPrivilegesList {
		localPrivilegeObj, ok := localPrivilege.(basetypes.ObjectValue)
		if !ok || localPrivilegeObj.IsNull() || localPrivilegeObj.IsUnknown() {
			return types.List{}, errors.New("failed to reorder role privileges")
		}
		localPrivilegeID := localPrivilegeObj.Attributes()["id"]

		for i, remotePrivilege := range remotePrivilegesList {
			remotePrivilegeObj, ok := remotePrivilege.(basetypes.ObjectValue)
			if !ok || remotePrivilegeObj.IsNull() || remotePrivilegeObj.IsUnknown() {
				return types.List{}, errors.New("failed to reorder role privileges")
			}
			remotePrivilegeID := remotePrivilegeObj.Attributes()["id"]

			if localPrivilegeID == remotePrivilegeID {
				orderedPrivileges = append(orderedPrivileges, remotePrivilege)
				addedRemotePrivileges[i] = true
				break
			}
		}
	}

	for i, remotePrivilege := range remotePrivilegesList {
		if !addedRemotePrivileges[i] {
			orderedPrivileges = append(orderedPrivileges, remotePrivilege)
		}
	}

	rolePrivilegesType := map[string]attr.Type{
		"id":         types.StringType,
		"name":       types.StringType,
		"permission": types.StringType,
	}

	orderedPrivilegeList, _ := types.ListValue(types.ObjectType{AttrTypes: rolePrivilegesType}, orderedPrivileges)
	return orderedPrivilegeList, nil
}
