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
	"encoding/json"
	"fmt"
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	// FilePoolPolicyActionSetDataAccessPatternType specifies action type for set_data_access_pattern.
	FilePoolPolicyActionSetDataAccessPatternType = "set_data_access_pattern"
	// FilePoolPolicyActionApplyDataStoragePolicyType specifies action type for apply_data_storage_policy.
	FilePoolPolicyActionApplyDataStoragePolicyType = "apply_data_storage_policy"
	// FilePoolPolicyActionApplySnapshotStoragePolicyType specifies action type for apply_snapshot_storage_policy.
	FilePoolPolicyActionApplySnapshotStoragePolicyType = "apply_snapshot_storage_policy"
	// FilePoolPolicyActionEnableCoalescerType specifies action type for enable_coalescer.
	FilePoolPolicyActionEnableCoalescerType = "enable_coalescer"
	// FilePoolPolicyActionEnablePackingType specifies action type for enable_packing.
	FilePoolPolicyActionEnablePackingType = "enable_packing"
	// FilePoolPolicyActionSetRequestedProtectionType specifies action type for set_requested_protection.
	FilePoolPolicyActionSetRequestedProtectionType = "set_requested_protection"
	// FilePoolPolicyActionSetCloudPoolPolicyType specifies action type for set_cloudpool_policy.
	FilePoolPolicyActionSetCloudPoolPolicyType = "set_cloudpool_policy"
)

// UpdateFilePoolPolicyResourceState updates resource state.
func UpdateFilePoolPolicyResourceState(ctx context.Context, policyModel *models.FilePoolPolicyModel, policyResponse *powerscale.V12FilepoolPolicyExtended) {
	if policyModel.ApplyOrder.IsUnknown() {
		policyModel.ApplyOrder = types.Int64Value(*policyResponse.ApplyOrder)
	}
	policyModel.Description = types.StringValue(*policyResponse.Description)
	policyModel.BirthClusterID = types.StringValue(*policyResponse.BirthClusterId)
	policyModel.ID = types.StringValue(*policyResponse.Id)
	policyModel.Name = types.StringValue(*policyResponse.Name)
	policyModel.State = types.StringValue(*policyResponse.State)
	policyModel.StateDetails = types.StringValue(*policyResponse.StateDetails)
}

// UpdateFilePoolPolicyImportState updates resource import state.
func UpdateFilePoolPolicyImportState(ctx context.Context, policyModel *models.FilePoolPolicyModel, policyResponse *powerscale.V12FilepoolPolicyExtended) (err error) {
	policyModel.ApplyOrder = types.Int64Value(*policyResponse.ApplyOrder)
	policyModel.Description = types.StringValue(*policyResponse.Description)
	policyModel.BirthClusterID = types.StringValue(*policyResponse.BirthClusterId)
	policyModel.ID = types.StringValue(*policyResponse.Id)
	policyModel.Name = types.StringValue(*policyResponse.Name)
	policyModel.State = types.StringValue(*policyResponse.State)
	policyModel.StateDetails = types.StringValue(*policyResponse.StateDetails)
	if err = parseActionParams(ctx, policyModel, policyResponse.Actions); err != nil {
		return
	}
	if err = parseFileMatchPattern(ctx, policyModel, policyResponse.FileMatchingPattern); err != nil {
		return
	}
	return
}

// UpdateFilePoolDefaultPolicyState updates resource state.
func UpdateFilePoolDefaultPolicyState(ctx context.Context, defaultPolicyModel *models.FilePoolPolicyModel, policyResponse *powerscale.V4FilepoolDefaultPolicyDefaultPolicy) (err error) {
	defaultPolicyModel.ID = types.StringValue("filepool_defaultpolicy")
	defaultPolicyModel.Description = types.StringValue("This policy applies to all files not selected by higher-priority policies.")
	defaultPolicyModel.IsDefaultPolicy = types.BoolValue(true)
	defaultPolicyModel.Name = types.StringValue("Default policy")
	defaultPolicyModel.ApplyOrder = types.Int64Null()
	defaultPolicyModel.BirthClusterID = types.StringNull()
	defaultPolicyModel.State = types.StringNull()
	defaultPolicyModel.StateDetails = types.StringNull()

	if len(defaultPolicyModel.Actions) != 0 {
		return
	}
	if err = parseActionParams(ctx, defaultPolicyModel, policyResponse.Actions); err != nil {
		return
	}
	return
}

// GetFilePoolPolicy Returns the file pool policy by name.
func GetFilePoolPolicy(ctx context.Context, client *client.Client, policyName string) (*powerscale.V12FilepoolPolicyExtended, error) {
	result, _, err := client.PscaleOpenAPIClient.FilepoolApi.GetFilepoolv12FilepoolPolicy(ctx, policyName).Execute()
	if err != nil {
		errStr := constants.ReadFilePoolPolicyErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting file pool policy: %s", message)
	}
	if len(result.Policies) <= 0 {
		message := constants.ReadFilePoolPolicyErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty file pool policy: %s", message)
	}
	return &result.Policies[0], err
}

// CreateFilePoolPolicy Creates a FilePoolPolicy.
func CreateFilePoolPolicy(ctx context.Context, client *client.Client, plan *models.FilePoolPolicyModel) (err error) {
	policyToCreate := powerscale.V12FilepoolPolicy{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueString(),
	}
	if !plan.ApplyOrder.IsUnknown() {
		policyToCreate.ApplyOrder = plan.ApplyOrder.ValueInt64Pointer()
	}
	actionList, err := buildActionParams(ctx, plan.Actions)
	if err != nil {
		return
	}
	policyToCreate.Actions = append(policyToCreate.Actions, actionList...)

	matchPattern, err := buildFileMatchPattern(ctx, plan)
	if err != nil {
		return
	}
	policyToCreate.FileMatchingPattern = *matchPattern

	if _, _, err := client.PscaleOpenAPIClient.FilepoolApi.CreateFilepoolv12FilepoolPolicy(ctx).V12FilepoolPolicy(policyToCreate).Execute(); err != nil {
		errStr := constants.CreateFilePoolPolicyErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error creating file pool policy: %s", message)
	}
	return
}

// UpdateFilePoolPolicy Updates a FilePoolPolicy parameters.
func UpdateFilePoolPolicy(ctx context.Context, client *client.Client, state *models.FilePoolPolicyModel, plan *models.FilePoolPolicyModel) (err error) {
	policyToUpdate := powerscale.V12FilepoolPolicyExtendedExtended{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
	}
	if !plan.ApplyOrder.IsUnknown() && plan.ApplyOrder.ValueInt64() != state.ApplyOrder.ValueInt64() {
		policyToUpdate.ApplyOrder = plan.ApplyOrder.ValueInt64Pointer()
	}
	actionList, err := buildActionParams(ctx, plan.Actions)
	if err != nil {
		return
	}
	policyToUpdate.Actions = append(policyToUpdate.Actions, actionList...)

	matchPattern, err := buildFileMatchPattern(ctx, plan)
	if err != nil {
		return
	}
	policyToUpdate.FileMatchingPattern = matchPattern
	if _, err := client.PscaleOpenAPIClient.FilepoolApi.UpdateFilepoolv12FilepoolPolicy(ctx, state.Name.ValueString()).V12FilepoolPolicy(policyToUpdate).Execute(); err != nil {
		errStr := constants.UpdateFilePoolPolicyErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error updating file pool policy: %s", message)
	}

	return
}

// DeleteFilePoolPolicy Deletes a FilePoolPolicy.
func DeleteFilePoolPolicy(ctx context.Context, client *client.Client, policyName string) error {
	if _, err := client.PscaleOpenAPIClient.FilepoolApi.DeleteFilepoolv12FilepoolPolicy(ctx, policyName).Execute(); err != nil {
		errStr := constants.DeleteFilePoolPolicyErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting file pool policy - %s : %s", policyName, message)
	}
	return nil
}

// GetFilePoolDefaultPolicy returns default FilePoolPolicy.
func GetFilePoolDefaultPolicy(ctx context.Context, client *client.Client) (defaultPolicy *powerscale.V4FilepoolDefaultPolicyDefaultPolicy, err error) {
	result, _, err := client.PscaleOpenAPIClient.FilepoolApi.GetFilepoolv4FilepoolDefaultPolicy(ctx).Execute()
	if err != nil {
		errStr := constants.ReadFilePoolPolicyErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting default file pool policy: %s", message)
	}
	return result.DefaultPolicy, err
}

// UpdateFilePoolDefaultPolicy Updates a default FilePoolPolicy params.
func UpdateFilePoolDefaultPolicy(ctx context.Context, client *client.Client, plan *models.FilePoolPolicyModel) (err error) {
	policyToUpdate := powerscale.V1FilepoolDefaultPolicyExtended{}
	actionList, err := buildActionParams(ctx, plan.Actions)
	if err != nil {
		return
	}
	policyToUpdate.Actions = append(policyToUpdate.Actions, actionList...)
	if _, err := client.PscaleOpenAPIClient.FilepoolApi.UpdateFilepoolv4FilepoolDefaultPolicy(ctx).V4FilepoolDefaultPolicy(policyToUpdate).Execute(); err != nil {
		errStr := constants.UpdateFilePoolPolicyErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error updating default file pool policy: %s", message)
	}
	return
}

// parseActionParams parses action params response to model according different action type.
func parseActionParams(ctx context.Context, policyModel *models.FilePoolPolicyModel, actionsResponse []powerscale.V1FilepoolDefaultPolicyAction) (err error) {
	actions := make([]models.V1FilepoolDefaultPolicyAction, 0)
	for _, action := range actionsResponse {
		actionModel := models.V1FilepoolDefaultPolicyAction{ActionType: types.StringValue(action.ActionType)}
		correctType := false
		switch action.ActionType {
		case FilePoolPolicyActionSetRequestedProtectionType:
			if value, ok := action.ActionParam.(string); ok {
				correctType = true
				actionModel.RequestedProtectionAction = types.StringValue(value)
			}
		case FilePoolPolicyActionSetDataAccessPatternType:
			if value, ok := action.ActionParam.(string); ok {
				correctType = true
				actionModel.DataAccessPatternAction = types.StringValue(value)
			}
		case FilePoolPolicyActionEnableCoalescerType:
			if value, ok := action.ActionParam.(bool); ok {
				correctType = true
				actionModel.EnableCoalescerAction = types.BoolValue(value)
			}
		case FilePoolPolicyActionEnablePackingType:
			if value, ok := action.ActionParam.(bool); ok {
				correctType = true
				actionModel.EnablePackingAction = types.BoolValue(value)
			}
		case FilePoolPolicyActionApplyDataStoragePolicyType:
			paramBytes, err := json.Marshal(action.ActionParam)
			if err != nil {
				return err
			}
			jsonModel := &models.V12StoragePolicyActionParamsJSONModel{}
			if err = json.Unmarshal(paramBytes, jsonModel); err != nil {
				return err
			}
			correctType = true
			actionModel.DataStoragePolicyAction = &models.V12StoragePolicyActionParams{
				SSDStrategy: types.StringValue(jsonModel.SSDStrategy),
				StoragePool: types.StringValue(jsonModel.StoragePool),
			}
		case FilePoolPolicyActionApplySnapshotStoragePolicyType:
			paramBytes, err := json.Marshal(action.ActionParam)
			if err != nil {
				return err
			}
			jsonModel := &models.V12StoragePolicyActionParamsJSONModel{}
			if err = json.Unmarshal(paramBytes, jsonModel); err != nil {
				return err
			}
			correctType = true
			actionModel.SnapshotStoragePolicyAction = &models.V12StoragePolicyActionParams{
				SSDStrategy: types.StringValue(jsonModel.SSDStrategy),
				StoragePool: types.StringValue(jsonModel.StoragePool),
			}
		case FilePoolPolicyActionSetCloudPoolPolicyType:
			paramBytes, err := json.Marshal(action.ActionParam)
			if err != nil {
				return err
			}
			jsonModel := &models.V12CloudPolicyArchiveParamsJSONModel{}
			if err = json.Unmarshal(paramBytes, jsonModel); err != nil {
				return err
			}
			cloudPoolPolicyAction := models.V12CloudPolicyArchiveParams{}
			if err = CopyFields(ctx, jsonModel, &cloudPoolPolicyAction); err == nil {
				correctType = true
				actionModel.CloudPoolPolicyAction = cloudPoolPolicyAction.CloudPolicyActionParams
			}
		default:
			return fmt.Errorf("unexpected action type: %s", action.ActionType)
		}
		if err != nil {
			return fmt.Errorf("failed to parse action param: %v for this action type: %s, Error: %s", action.ActionParam, action.ActionType, err.Error())
		}
		if !correctType {

			return fmt.Errorf("unexpected action param: %v for this action type: %s", action.ActionParam, action.ActionType)
		}
		actions = append(actions, actionModel)
	}
	policyModel.Actions = actions
	return
}

// buildActionParams builds action params according different action type.
func buildActionParams(ctx context.Context, actions []models.V1FilepoolDefaultPolicyAction) (actionsListBody []powerscale.V1FilepoolDefaultPolicyAction, err error) {
	actionMap := map[string]powerscale.V1FilepoolDefaultPolicyAction{
		FilePoolPolicyActionSetRequestedProtectionType:     {ActionType: FilePoolPolicyActionSetRequestedProtectionType, ActionParam: nil},
		FilePoolPolicyActionSetDataAccessPatternType:       {ActionType: FilePoolPolicyActionSetDataAccessPatternType, ActionParam: nil},
		FilePoolPolicyActionEnableCoalescerType:            {ActionType: FilePoolPolicyActionEnableCoalescerType, ActionParam: nil},
		FilePoolPolicyActionEnablePackingType:              {ActionType: FilePoolPolicyActionEnablePackingType, ActionParam: nil},
		FilePoolPolicyActionApplyDataStoragePolicyType:     {ActionType: FilePoolPolicyActionApplyDataStoragePolicyType, ActionParam: nil},
		FilePoolPolicyActionApplySnapshotStoragePolicyType: {ActionType: FilePoolPolicyActionApplySnapshotStoragePolicyType, ActionParam: nil},
		FilePoolPolicyActionSetCloudPoolPolicyType:         {ActionType: FilePoolPolicyActionSetCloudPoolPolicyType, ActionParam: nil},
	}
	for _, action := range actions {
		switch action.ActionType.ValueString() {
		case FilePoolPolicyActionSetRequestedProtectionType:
			actionBody := actionMap[FilePoolPolicyActionSetRequestedProtectionType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionSetRequestedProtectionType)
			}
			actionBody.ActionParam = action.RequestedProtectionAction.ValueString()
			actionMap[FilePoolPolicyActionSetRequestedProtectionType] = actionBody
		case FilePoolPolicyActionSetDataAccessPatternType:
			actionBody := actionMap[FilePoolPolicyActionSetDataAccessPatternType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionSetDataAccessPatternType)
			}
			actionBody.ActionParam = action.DataAccessPatternAction.ValueString()
			actionMap[FilePoolPolicyActionSetDataAccessPatternType] = actionBody
		case FilePoolPolicyActionEnableCoalescerType:
			actionBody := actionMap[FilePoolPolicyActionEnableCoalescerType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionEnableCoalescerType)
			}
			actionBody.ActionParam = action.EnableCoalescerAction.ValueBool()
			actionMap[FilePoolPolicyActionEnableCoalescerType] = actionBody
		case FilePoolPolicyActionEnablePackingType:
			actionBody := actionMap[FilePoolPolicyActionEnablePackingType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionEnablePackingType)
			}
			actionBody.ActionParam = action.EnablePackingAction.ValueBool()
			actionMap[FilePoolPolicyActionEnablePackingType] = actionBody
		case FilePoolPolicyActionApplyDataStoragePolicyType:
			actionBody := actionMap[FilePoolPolicyActionApplyDataStoragePolicyType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionApplyDataStoragePolicyType)
			}
			dataStoragePolicyAction := models.V12StoragePolicyActionParamsJSONModel{}
			if err = ReadFromState(ctx, action.DataStoragePolicyAction, &dataStoragePolicyAction); err != nil {
				return
			}
			actionBody.ActionParam = dataStoragePolicyAction
			actionMap[FilePoolPolicyActionApplyDataStoragePolicyType] = actionBody
		case FilePoolPolicyActionApplySnapshotStoragePolicyType:
			actionBody := actionMap[FilePoolPolicyActionApplySnapshotStoragePolicyType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionApplySnapshotStoragePolicyType)
			}
			snapshotStoragePolicyAction := models.V12StoragePolicyActionParamsJSONModel{}
			if err = ReadFromState(ctx, action.SnapshotStoragePolicyAction, &snapshotStoragePolicyAction); err != nil {
				return
			}
			actionBody.ActionParam = snapshotStoragePolicyAction
			actionMap[FilePoolPolicyActionApplySnapshotStoragePolicyType] = actionBody
		case FilePoolPolicyActionSetCloudPoolPolicyType:
			actionBody := actionMap[FilePoolPolicyActionSetCloudPoolPolicyType]
			if actionBody.ActionParam != nil {
				return actionsListBody, fmt.Errorf("duplicated action type: %s", FilePoolPolicyActionSetCloudPoolPolicyType)
			}
			cloudPoolPolicyAction := models.V12CloudPolicyActionParamsJSONModel{}
			if err = ReadFromState(ctx, action.CloudPoolPolicyAction, &cloudPoolPolicyAction); err != nil {
				return
			}
			if action.CloudPoolPolicyAction != nil && action.CloudPoolPolicyAction.Cache != nil {
				cloudPoolPolicyCacheAction := models.V12CloudPolicyActionCacheParamsJSONModel{}
				if err = ReadFromState(ctx, action.CloudPoolPolicyAction.Cache, &cloudPoolPolicyCacheAction); err != nil {
					return
				}
				cloudPoolPolicyAction.Cache = &cloudPoolPolicyCacheAction
			}
			actionBody.ActionParam = models.V12CloudPolicyArchiveParamsJSONModel{CloudPolicyActionParams: &cloudPoolPolicyAction}
			actionMap[FilePoolPolicyActionSetCloudPoolPolicyType] = actionBody
		default:
			return actionsListBody, fmt.Errorf("unexpected action type: %s", action.ActionType.ValueString())
		}
	}
	for _, value := range actionMap {
		actionsListBody = append(actionsListBody, value)
	}
	return
}

// buildFileMatchPattern builds criteria into FileMatchPattern.
func buildFileMatchPattern(ctx context.Context, plan *models.FilePoolPolicyModel) (matchPattern *powerscale.V1FilepoolPolicyFileMatchingPattern, err error) {
	orCriteriaListBody := make([]powerscale.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem, 0)
	for _, orCriteria := range plan.FileMatchingPattern.OrCriteria {
		andCriteriaListBody := make([]powerscale.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem, 0)
		for _, andCriteria := range orCriteria.AndCriteria {
			andCriteriaBody := powerscale.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem{}
			if err = ReadFromState(ctx, andCriteria, &andCriteriaBody); err != nil {
				return nil, err
			}
			switch andCriteria.Type.ValueString() {
			case "name", "custom_attribute", "file_type", "path":
				andCriteriaBody.Value = andCriteria.Value.ValueString()
			case "link_count", "accessed_time", "birth_time", "changed_time", "metadata_changed_time", "size":
				v, err := strconv.ParseInt(andCriteria.Value.ValueString(), 10, 64)
				if err != nil {
					return nil, fmt.Errorf("input value should be an integral string. unexpected value:%s for criteria type: %s", andCriteria.Value.ValueString(), andCriteria.Type.ValueString())
				}
				andCriteriaBody.Value = v
			default:
				return nil, fmt.Errorf("unsupported criteria type: %s", andCriteria.Type.ValueString())
			}
			andCriteriaListBody = append(andCriteriaListBody, andCriteriaBody)
		}
		orCriteriaBody := powerscale.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem{AndCriteria: andCriteriaListBody}
		orCriteriaListBody = append(orCriteriaListBody, orCriteriaBody)
	}
	matchPattern = &powerscale.V1FilepoolPolicyFileMatchingPattern{OrCriteria: orCriteriaListBody}
	return
}

// parseFileMatchPattern parses FileMatchPattern response to terraform model.
func parseFileMatchPattern(ctx context.Context, policyModel *models.FilePoolPolicyModel, matchPattern *powerscale.V1FilepoolPolicyFileMatchingPattern) (err error) {
	orCriteriaList := make([]models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem, 0)
	for _, orCriteria := range matchPattern.OrCriteria {
		andCriteriaList := make([]models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem, 0)
		for _, andCriteria := range orCriteria.AndCriteria {
			andCriteriaModel := models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem{}
			if err = CopyFields(ctx, andCriteria, &andCriteriaModel); err != nil {
				return
			}
			correctType := false
			switch andCriteria.Type {
			case "name", "custom_attribute", "file_type", "path":
				if value, ok := andCriteria.Value.(string); ok {
					correctType = true
					andCriteriaModel.Value = types.StringValue(value)
				}
			case "link_count", "accessed_time", "birth_time", "changed_time", "metadata_changed_time", "size":
				if value, ok := andCriteria.Value.(float64); ok {
					correctType = true
					andCriteriaModel.Value = types.StringValue(strconv.FormatInt(int64(value), 10))
				}
			default:
				return fmt.Errorf("unsupported criteria type: %s", andCriteria.Type)
			}
			if !correctType {
				return fmt.Errorf("unexpected value for criteria type: %s", andCriteria.Type)
			}
			andCriteriaList = append(andCriteriaList, andCriteriaModel)
		}
		orCriteriaList = append(orCriteriaList, models.V1FilepoolPolicyFileMatchingPatternOrCriteriaItem{AndCriteria: andCriteriaList})
	}
	policyModel.FileMatchingPattern = &models.V1FilepoolPolicyFileMatchingPattern{OrCriteria: orCriteriaList}
	return
}

// IsPolicyParamInvalid Verify if policy params is valid for default policy and common policy.
func IsPolicyParamInvalid(plan models.FilePoolPolicyModel) error {
	if plan.IsDefaultPolicy.ValueBool() {
		if plan.Name.ValueString() != "Default policy" {
			return fmt.Errorf("the default policy name should be \"Default policy\"")
		}
		if plan.FileMatchingPattern != nil || !plan.ApplyOrder.IsUnknown() || !plan.Description.IsUnknown() {
			return fmt.Errorf("may not specify file_matching_pattern, apply_order and description for default policy")
		}
		if plan.Actions != nil {
			for _, action := range plan.Actions {
				if action.ActionType.ValueString() == FilePoolPolicyActionSetCloudPoolPolicyType {
					return fmt.Errorf("default policy not support \"set_cloudpool_policy\" action")
				}
			}
		}
		return nil
	}

	if plan.FileMatchingPattern == nil {
		return fmt.Errorf("file_matching_pattern is required")
	}
	return nil
}
