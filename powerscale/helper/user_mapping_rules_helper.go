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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateUserMappingRulesState updates resource state.
func UpdateUserMappingRulesState(ctx context.Context, rulesState *models.UserMappingRulesResourceModel, rulesResponse *powerscale.V1MappingUsersRulesRules) (diags diag.Diagnostics) {
	mappingUsersRulesRuleUser2Type := map[string]attr.Type{
		"domain": types.StringType,
		"user":   types.StringType,
	}
	defaultUnixUserType := map[string]attr.Type{
		"default_unix_user": types.ObjectType{AttrTypes: mappingUsersRulesRuleUser2Type},
	}
	mappingUsersRulesRuleOptionsType := map[string]attr.Type{
		"default_user": types.ObjectType{AttrTypes: mappingUsersRulesRuleUser2Type},
		"break":        types.BoolType,
		"group":        types.BoolType,
		"groups":       types.BoolType,
		"user":         types.BoolType,
	}
	mappingUsersRulesRuleType := map[string]attr.Type{
		"operator":    types.StringType,
		"options":     types.ObjectType{AttrTypes: mappingUsersRulesRuleOptionsType},
		"target_user": types.ObjectType{AttrTypes: mappingUsersRulesRuleUser2Type},
		"source_user": types.ObjectType{AttrTypes: mappingUsersRulesRuleUser2Type},
	}

	if rulesResponse.Parameters != nil {
		paramMap := map[string]attr.Value{
			"default_unix_user": types.ObjectNull(mappingUsersRulesRuleUser2Type),
		}
		if rulesResponse.Parameters.DefaultUnixUser != nil {
			defaultUserMap := map[string]attr.Value{
				"user":   types.StringValue(rulesResponse.Parameters.DefaultUnixUser.User),
				"domain": types.StringNull(),
			}
			if rulesResponse.Parameters.DefaultUnixUser.Domain != nil {
				defaultUserMap["domain"] = types.StringValue(*rulesResponse.Parameters.DefaultUnixUser.Domain)
			}
			defaultUserObject, diags := types.ObjectValue(mappingUsersRulesRuleUser2Type, defaultUserMap)
			if diags.HasError() {
				return diags
			}
			paramMap["default_unix_user"] = defaultUserObject
		}
		rulesState.Parameters, diags = types.ObjectValue(defaultUnixUserType, paramMap)
		if diags.HasError() {
			return
		}
	}

	var rulesAttrs []attr.Value
	for _, ruleResp := range rulesResponse.Rules {
		optionObj := types.ObjectNull(mappingUsersRulesRuleOptionsType)
		if ruleResp.Options != nil {
			optionMap := map[string]attr.Value{
				"group":        types.BoolNull(),
				"groups":       types.BoolNull(),
				"user":         types.BoolNull(),
				"break":        types.BoolNull(),
				"default_user": types.ObjectNull(mappingUsersRulesRuleUser2Type),
			}
			onlySupportBreak := false
			operator := ruleResp.GetOperator()
			if operator == "replace" || operator == "trim" || operator == "union" {
				onlySupportBreak = true
			}
			if ruleResp.Options.Group != nil && !onlySupportBreak {
				optionMap["group"] = types.BoolValue(*ruleResp.Options.Group)
			}
			if ruleResp.Options.Groups != nil && !onlySupportBreak {
				optionMap["groups"] = types.BoolValue(*ruleResp.Options.Groups)
			}
			if ruleResp.Options.User != nil && !onlySupportBreak {
				optionMap["user"] = types.BoolValue(*ruleResp.Options.User)
			}
			if ruleResp.Options.Break != nil {
				optionMap["break"] = types.BoolValue(*ruleResp.Options.Break)
			}
			if ruleResp.Options.DefaultUser != nil {
				defaultUserMap := map[string]attr.Value{
					"domain": types.StringNull(),
					"user":   types.StringValue(ruleResp.Options.DefaultUser.User),
				}
				if ruleResp.Options.DefaultUser.Domain != nil {
					defaultUserMap["domain"] = types.StringValue(*ruleResp.Options.DefaultUser.Domain)
				}
				optionMap["default_user"], diags = types.ObjectValue(mappingUsersRulesRuleUser2Type, defaultUserMap)
				if diags.HasError() {
					return
				}
			}
			optionObj, diags = types.ObjectValue(mappingUsersRulesRuleOptionsType, optionMap)
			if diags.HasError() {
				return
			}
		}

		targetUserMap := map[string]attr.Value{
			"domain": types.StringNull(),
			"user":   types.StringValue(ruleResp.User1.User),
		}
		if ruleResp.User1.Domain != nil {
			targetUserMap["domain"] = types.StringValue(*ruleResp.User1.Domain)
		}
		targetUserObj, diag := types.ObjectValue(mappingUsersRulesRuleUser2Type, targetUserMap)
		if diag.HasError() {
			return diag
		}
		sourceUserObj := types.ObjectNull(mappingUsersRulesRuleUser2Type)
		if ruleResp.User2 != nil {
			sourceUserMap := map[string]attr.Value{
				"domain": types.StringNull(),
				"user":   types.StringValue(ruleResp.User2.User),
			}
			if ruleResp.User2.Domain != nil {
				sourceUserMap["domain"] = types.StringValue(*ruleResp.User2.Domain)
			}
			sourceUserObj, diags = types.ObjectValue(mappingUsersRulesRuleUser2Type, sourceUserMap)
			if diags.HasError() {
				return
			}
		}

		ruleItemMap := map[string]attr.Value{
			"operator":    types.StringNull(),
			"target_user": targetUserObj,
			"source_user": sourceUserObj,
			"options":     optionObj,
		}
		if ruleResp.Operator != nil {
			ruleItemMap["operator"] = types.StringValue(*ruleResp.Operator)
		}
		ruleItemObj, diags := types.ObjectValue(mappingUsersRulesRuleType, ruleItemMap)
		if diags.HasError() {
			return diags
		}
		rulesAttrs = append(rulesAttrs, ruleItemObj)
	}
	rulesState.Rules, diags = types.ListValue(types.ObjectType{AttrTypes: mappingUsersRulesRuleType}, rulesAttrs)
	if diags.HasError() {
		return
	}

	rulesState.ID = types.StringValue("user_mapping_rules")
	return
}

// GetUserMappingRulesByZone returns user mapping rules detail in specific zone.
func GetUserMappingRulesByZone(ctx context.Context, client *client.Client, zone string) (*powerscale.V1MappingUsersRulesRules, error) {

	getParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv1MappingUsersRules(ctx)
	if zone != "" {
		getParam = getParam.Zone(zone)
	}
	result, _, err := getParam.Execute()
	if err != nil {
		errStr := constants.ReadUserMappingRulesErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user mapping rules: %s", message)
	}

	return result.Rules, nil
}

// UpdateUserMappingRules Updates user mapping rules.
func UpdateUserMappingRules(ctx context.Context, client *client.Client, state *models.UserMappingRulesResourceModel, plan *models.UserMappingRulesResourceModel) (diags diag.Diagnostics) {

	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv1MappingUsersRules(ctx)
	if !plan.Zone.IsNull() {
		updateParam = updateParam.Zone(plan.Zone.ValueString())
	}
	params := &powerscale.V1MappingUsersRulesRulesParameters{}
	if !plan.Parameters.IsUnknown() {
		if err := assignObjectToField(ctx, plan.Parameters, params); err != nil {
			diags.AddError("error updating user mapping rules: parse Rule Parameters failed", err.Error())
			return
		}
	}

	rulesBody := make([]powerscale.V1MappingUsersRulesRule, 0)
	if !plan.Rules.IsNull() {
		var rulesObjectList []models.V1MappingUsersRulesRule
		if !plan.Rules.IsUnknown() {
			if diags = plan.Rules.ElementsAs(ctx, &rulesObjectList, false); diags.HasError() {
				return
			}
		}

		for _, rulesObj := range rulesObjectList {
			ruleInput, err := buildUserMappingRuleInput(ctx, rulesObj)
			if err != nil {
				diags.AddError("error updating user mapping rules: parse Rules failed", err.Error())
				return
			}
			rulesBody = append(rulesBody, *ruleInput)
		}
	}

	body := &powerscale.V1MappingUsersRulesRules{Parameters: params, Rules: rulesBody}
	updateParam = updateParam.V1MappingUsersRules(*body)
	if _, err := updateParam.Execute(); err != nil {
		errStr := constants.UpdateUserMappingRulesErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError("error updating user mapping rules", message)
		return
	}

	return
}

// GetLookupMappingUsers return lookup mapping users detail.
func GetLookupMappingUsers(ctx context.Context, client *client.Client, zone string, user models.UserMemberItem) ([]powerscale.V1MappingUsersLookupMappingItem, error) {
	getParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv1MappingUsersLookup(ctx)
	if !user.Name.IsNull() && user.Name.ValueString() != "" {
		getParam = getParam.User(user.Name.ValueString())
	}
	if !user.UID.IsNull() {
		getParam = getParam.Uid(int32(user.UID.ValueInt64()))
	}
	if zone != "" {
		getParam = getParam.Zone(zone)
	}
	result, _, err := getParam.Execute()
	if err != nil {
		errStr := constants.ReadTestUserMappingErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting user mapping test: %s", message)
	}

	return result.Mapping, nil
}

// UpdateLookupMappingUsersState update all lookup mapping users state.
func UpdateLookupMappingUsersState(ctx context.Context, client *client.Client, plan *models.UserMappingRulesResourceModel) (diags diag.Diagnostics) {

	var lookupUsers []powerscale.V1MappingUsersLookupMappingItem
	for _, user := range plan.TestMappingUsers {
		userResponse, err := GetLookupMappingUsers(ctx, client, plan.Zone.ValueString(), user)
		if err != nil {
			diags.AddError(fmt.Sprintf("error getting the mapping test user:%s", user.Name.ValueString()), err.Error())
			return
		}
		lookupUsers = append(lookupUsers, userResponse...)
	}

	var lookupUsersObjects []attr.Value
	groupItemType := map[string]attr.Type{
		"name": types.StringType,
		"gid":  types.StringType,
		"sid":  types.StringType,
	}
	privilegeItemType := map[string]attr.Type{
		"name":      types.StringType,
		"id":        types.StringType,
		"read_only": types.BoolType,
	}
	userItemType := map[string]attr.Type{
		"name":                  types.StringType,
		"on_disk_user_identity": types.StringType,
		"primary_group_sid":     types.StringType,
		"primary_group_name":    types.StringType,
		"sid":                   types.StringType,
		"uid":                   types.StringType,
	}
	lookupUserItemType := map[string]attr.Type{
		"zid":                     types.Int64Type,
		"zone":                    types.StringType,
		"user":                    types.ObjectType{AttrTypes: userItemType},
		"privileges":              types.ListType{ElemType: types.ObjectType{AttrTypes: privilegeItemType}},
		"supplemental_identities": types.ListType{ElemType: types.ObjectType{AttrTypes: groupItemType}},
	}
	for _, user := range lookupUsers {
		var groupsObjects []attr.Value
		for _, group := range user.Groups {
			groupItemMap := map[string]attr.Value{
				"name": types.StringValue(group.Name),
				"gid":  types.StringNull(),
				"sid":  types.StringNull(),
			}
			if group.Gid != nil && group.Gid.Id != nil {
				groupItemMap["gid"] = types.StringValue(*group.Gid.Id)
			}
			if group.Sid.Id != nil {
				groupItemMap["sid"] = types.StringValue(*group.Sid.Id)
			}
			groupObject, diags := types.ObjectValue(groupItemType, groupItemMap)
			if diags.HasError() {
				return diags
			}
			groupsObjects = append(groupsObjects, groupObject)
		}
		groupsListState := types.ListNull(types.ObjectType{AttrTypes: groupItemType})
		if len(groupsObjects) > 0 {
			groupsListState, diags = types.ListValue(types.ObjectType{AttrTypes: groupItemType}, groupsObjects)
			if diags.HasError() {
				return
			}
		}

		var privilegeObjects []attr.Value
		for _, privilege := range user.Privileges {
			privilegeItemMap := map[string]attr.Value{
				"name":      types.StringNull(),
				"id":        types.StringValue(privilege.Id),
				"read_only": types.BoolNull(),
			}
			if privilege.Name != nil {
				privilegeItemMap["name"] = types.StringValue(*privilege.Name)
			}
			if privilege.ReadOnly != nil {
				privilegeItemMap["read_only"] = types.BoolValue(*privilege.ReadOnly)
			}
			privilegeObject, diags := types.ObjectValue(privilegeItemType, privilegeItemMap)
			if diags.HasError() {
				return diags
			}
			privilegeObjects = append(privilegeObjects, privilegeObject)
		}
		privilegeListState := types.ListNull(types.ObjectType{AttrTypes: privilegeItemType})
		if len(privilegeObjects) > 0 {
			privilegeListState, diags = types.ListValue(types.ObjectType{AttrTypes: privilegeItemType}, privilegeObjects)
			if diags.HasError() {
				return
			}
		}

		var userItemObject attr.Value
		if user.User != nil {
			userItemMap := map[string]attr.Value{
				"name":                  types.StringValue(user.User.Name),
				"on_disk_user_identity": types.StringNull(),
				"primary_group_sid":     types.StringNull(),
				"primary_group_name":    types.StringNull(),
				"sid":                   types.StringNull(),
				"uid":                   types.StringNull(),
			}
			if user.User.Sid.Id != nil {
				userItemMap["sid"] = types.StringValue(*user.User.Sid.Id)
			}
			if user.User.Uid.Id != nil {
				userItemMap["uid"] = types.StringValue(*user.User.Uid.Id)
			}
			if user.User.OnDiskUserIdentity.Id != nil {
				userItemMap["on_disk_user_identity"] = types.StringValue(*user.User.OnDiskUserIdentity.Id)
			}
			if user.User.PrimaryGroupSid.Id != nil {
				userItemMap["primary_group_sid"] = types.StringValue(*user.User.PrimaryGroupSid.Id)
			}
			if user.User.PrimaryGroupSid.Name != nil {
				userItemMap["primary_group_name"] = types.StringValue(*user.User.PrimaryGroupSid.Name)
			}
			userItemObject, diags = types.ObjectValue(userItemType, userItemMap)
			if diags.HasError() {
				return
			}
		}

		lookupUserItemMap := map[string]attr.Value{
			"zid":                     types.Int64Null(),
			"zone":                    types.StringNull(),
			"user":                    userItemObject,
			"privileges":              privilegeListState,
			"supplemental_identities": groupsListState,
		}
		if user.Zid != nil {
			lookupUserItemMap["zid"] = types.Int64Value(int64(*user.Zid))
		}
		if user.Zone != nil {
			lookupUserItemMap["zone"] = types.StringValue(*user.Zone)
		}

		lookupUsersObject, diags := types.ObjectValue(lookupUserItemType, lookupUserItemMap)
		if diags.HasError() {
			return diags
		}
		lookupUsersObjects = append(lookupUsersObjects, lookupUsersObject)
	}

	plan.TestMappingUserResults = types.ListNull(types.ObjectType{AttrTypes: lookupUserItemType})
	if len(lookupUsersObjects) > 0 {
		plan.TestMappingUserResults, diags = types.ListValue(types.ObjectType{AttrTypes: lookupUserItemType}, lookupUsersObjects)
		if diags.HasError() {
			return
		}
	}

	return
}

// buildUserMappingRuleInput builds rule model into rule body.
func buildUserMappingRuleInput(ctx context.Context, rule models.V1MappingUsersRulesRule) (*powerscale.V1MappingUsersRulesRule, error) {
	ruleBody := &powerscale.V1MappingUsersRulesRule{}
	if err := ReadFromState(ctx, rule, ruleBody); err != nil {
		return nil, err
	}

	targetUser := &powerscale.V1MappingUsersRulesRuleUser2{}
	sourceUser := &powerscale.V1MappingUsersRulesRuleUser2{}
	if err := assignObjectToField(ctx, rule.User1, targetUser); err != nil {
		return nil, err
	}
	if err := assignObjectToField(ctx, rule.User2, sourceUser); err != nil {
		return nil, err
	}
	ruleBody.User1 = *targetUser
	ruleBody.User2 = sourceUser

	return ruleBody, nil
}
