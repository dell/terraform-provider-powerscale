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

// UpdateFilePoolPolicyResourceState updates resource state.
func UpdateFilePoolPolicyResourceState(ctx context.Context, policyModel *models.FilePoolPolicyModel, policyResponse *powerscale.V12FilepoolPolicyExtended) {
	if policyResponse.HasApplyOrder() {
		policyModel.ApplyOrder = types.Int64Value(*policyResponse.ApplyOrder)
	} else if policyModel.ApplyOrder.IsUnknown() {
		policyModel.ApplyOrder = types.Int64Null()
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
	UpdateFilePoolPolicyResourceState(ctx, policyModel, policyResponse)
	if err = parseActionParams(ctx, policyModel, policyResponse.Actions); err != nil {
		return
	}
	if err = parseFileMatchPattern(ctx, policyModel, policyResponse.FileMatchingPattern); err != nil {
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
		ApplyOrder:  plan.ApplyOrder.ValueInt64Pointer(),
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueString(),
	}
	actionList, err := buildActionParams(ctx, plan)
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
		ApplyOrder:  plan.ApplyOrder.ValueInt64Pointer(),
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueStringPointer(),
	}
	actionList, err := buildActionParams(ctx, plan)
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

// parseActionParams parses action params response to model according different action type.
func parseActionParams(ctx context.Context, policyModel *models.FilePoolPolicyModel, actionsResponse []powerscale.V1FilepoolDefaultPolicyAction) (err error) {
	actions := make([]models.V1FilepoolDefaultPolicyAction, 0)
	for _, action := range actionsResponse {
		actionModel := models.V1FilepoolDefaultPolicyAction{ActionType: types.StringValue(action.ActionType)}
		correctType := false
		switch action.ActionType {
		case "set_requested_protection":
			if value, ok := action.ActionParam.(string); ok {
				correctType = true
				actionModel.RequestedProtectionAction = types.StringValue(value)
			}
		case "set_data_access_pattern":
			if value, ok := action.ActionParam.(string); ok {
				correctType = true
				actionModel.DataAccessPatternAction = types.StringValue(value)
			}
		case "enable_coalescer":
			if value, ok := action.ActionParam.(bool); ok {
				correctType = true
				actionModel.EnableCoalescerAction = types.BoolValue(value)
			}
		case "enable_packing":
			if value, ok := action.ActionParam.(bool); ok {
				correctType = true
				actionModel.EnablePackingAction = types.BoolValue(value)
			}
		case "apply_data_storage_policy":
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
		case "apply_snapshot_storage_policy":
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
		case "set_cloudpool_policy":
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
func buildActionParams(ctx context.Context, plan *models.FilePoolPolicyModel) (actionsListBody []powerscale.V1FilepoolDefaultPolicyAction, err error) {
	actionsListBody = make([]powerscale.V1FilepoolDefaultPolicyAction, 0)
	for _, action := range plan.Actions {
		actionBody := powerscale.V1FilepoolDefaultPolicyAction{ActionType: action.ActionType.ValueString()}
		switch action.ActionType.ValueString() {
		case "set_requested_protection":
			actionBody.ActionParam = action.RequestedProtectionAction.ValueString()
		case "set_data_access_pattern":
			actionBody.ActionParam = action.DataAccessPatternAction.ValueString()
		case "enable_coalescer":
			actionBody.ActionParam = action.EnableCoalescerAction.ValueBool()
		case "enable_packing":
			actionBody.ActionParam = action.EnablePackingAction.ValueBool()
		case "apply_data_storage_policy":
			dataStoragePolicyAction := models.V12StoragePolicyActionParamsJSONModel{}
			if err = ReadFromState(ctx, action.DataStoragePolicyAction, &dataStoragePolicyAction); err != nil {
				return
			}
			actionBody.ActionParam = dataStoragePolicyAction
		case "apply_snapshot_storage_policy":
			snapshotStoragePolicyAction := models.V12StoragePolicyActionParamsJSONModel{}
			if err = ReadFromState(ctx, action.SnapshotStoragePolicyAction, &snapshotStoragePolicyAction); err != nil {
				return
			}
			actionBody.ActionParam = snapshotStoragePolicyAction
		case "set_cloudpool_policy":
			cloudPoolPolicyAction := models.V12CloudPolicyActionParamsJSONModel{}
			if err = ReadFromState(ctx, action.CloudPoolPolicyAction, &cloudPoolPolicyAction); err != nil {
				return
			}
			actionBody.ActionParam = models.V12CloudPolicyArchiveParamsJSONModel{CloudPolicyActionParams: &cloudPoolPolicyAction}
		default:
			return actionsListBody, fmt.Errorf("unexpected action type: %s", action.ActionType.ValueString())
		}
		actionsListBody = append(actionsListBody, actionBody)
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
