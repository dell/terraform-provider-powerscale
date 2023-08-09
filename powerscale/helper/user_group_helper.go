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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// GetAllGroupMembers returns all group members
func GetAllGroupMembers(ctx context.Context, client *client.Client, groupName string) (members []powerscale.V1AuthAccessAccessItemFileGroup, err error) {
	memberParams := client.PscaleOpenAPIClient.AuthGroupsApi.ListAuthGroupsv1GroupMembers(ctx, groupName)
	result, _, err := memberParams.Execute()
	if err != nil {
		errStr := constants.ReadUserGroupMemberErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user group members: %s", message)
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
			return nil, fmt.Errorf("error getting user group members with resume: %s", message)
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
		if !filter.ResolveNames.IsNull() {
			groupParams = groupParams.ResolveNames(filter.ResolveNames.ValueBool())
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
					continue
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
