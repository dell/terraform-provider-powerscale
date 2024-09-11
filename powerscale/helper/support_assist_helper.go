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
	"terraform-provider-powerscale/powerscale/models"
	"time"

	"terraform-provider-powerscale/powerscale/constants"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetSupportAssistTerms retrieves the support assist terms
func GetSupportAssistTerms(ctx context.Context, client *client.Client) (*powerscale.V16SupportassistTerms, error) {
	terms, _, err := client.PscaleOpenAPIClient.SupportassistApi.GetSupportassistv16SupportassistTerms(ctx).Execute()
	return terms, err
}

// UpdateSupportAssistTerms updates the support assist terms
func UpdateSupportAssistTerms(ctx context.Context, client *client.Client, V16SupportAssistTermsExtended powerscale.V16SupportassistTermsExtended) error {
	_, err := client.PscaleOpenAPIClient.SupportassistApi.UpdateSupportassistv16SupportassistTerms(ctx).V16SupportassistTerms(V16SupportAssistTermsExtended).Execute()
	return err
}

// GetSupportAssist retrieves the support assist settings
func GetSupportAssist(ctx context.Context, client *client.Client) (*powerscale.V16SupportassistSettings, error) {
	supportAssistSettings, _, err := client.PscaleOpenAPIClient.SupportassistApi.GetSupportassistv16SupportassistSettings(ctx).Execute()
	return supportAssistSettings, err
}

// UpdateSupportAssistStatus updates the status of the support assist
func UpdateSupportAssistStatus(ctx context.Context, client *client.Client, V16SupportAssistStatusExtended powerscale.V16SupportassistStatusExtended) error {
	_, err := client.PscaleOpenAPIClient.SupportassistApi.UpdateSupportassistv16SupportassistStatus(ctx).V16SupportassistStatus(V16SupportAssistStatusExtended).Execute()
	return err
}

// UpdateSupportAssistSettings updates the support assist settings
func UpdateSupportAssistSettings(ctx context.Context, client *client.Client, V16SupportAssistSettingsExtended powerscale.V16SupportassistSettingsExtended) error {
	_, err := client.PscaleOpenAPIClient.SupportassistApi.UpdateSupportassistv16SupportassistSettings(ctx).V16SupportassistSettings(V16SupportAssistSettingsExtended).Execute()
	return err
}

// CreateSupportAssistv16Task creates a support assist task
func CreateSupportAssistv16Task(ctx context.Context, client *client.Client, V16SupportAssistTask powerscale.V16SupportassistTaskItem) (*powerscale.CreateTaskResponse, error) {
	response, _, err := client.PscaleOpenAPIClient.SupportassistApi.CreateSupportassistv16SupportassistTaskItem(ctx).V16SupportassistTaskItem(V16SupportAssistTask).Execute()
	return response, err
}

// GetSupportAssistv16Task retrieves a support assist task by its ID
func GetSupportAssistv16Task(ctx context.Context, client *client.Client, id string) (*powerscale.V16SupportassistTaskId, error) {
	response, _, err := client.PscaleOpenAPIClient.SupportassistApi.GetSupportassistv16SupportassistTaskById(ctx, id).Execute()
	return response, err
}

// CreateSupportAssistv17Task creates a support assist task
func CreateSupportAssistv17Task(ctx context.Context, client *client.Client, V17SupportAssistTask powerscale.V16SupportassistTaskItem) (*powerscale.CreateTaskResponse, error) {
	response, _, err := client.PscaleOpenAPIClient.SupportassistApi.CreateSupportassistv17SupportassistTaskItem(ctx).V17SupportassistTaskItem(V17SupportAssistTask).Execute()
	return response, err
}

// GetSupportAssistv17Task retrieves a support assist task by its ID
func GetSupportAssistv17Task(ctx context.Context, client *client.Client, id string) (*powerscale.V16SupportassistTaskId, error) {
	response, _, err := client.PscaleOpenAPIClient.SupportassistApi.GetSupportassistv17SupportassistTaskById(ctx, id).Execute()
	return response, err
}

// GetClusterVersion retrieves the cluster version
func GetClusterVersion(ctx context.Context, client *client.Client) (*powerscale.V3ClusterVersion, error) {
	clusterVersion, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterVersion(ctx).Execute()
	return clusterVersion, err
}

// ManageSupportAssist manages the support assist settings
func ManageSupportAssist(ctx context.Context, client *client.Client, plan models.SupportAssistModel) (state models.SupportAssistModel, resp diag.Diagnostics) {
	// Update support assist terms status
	if !plan.Accepted.IsNull() {
		terms := powerscale.V16SupportassistTermsExtended{
			Accepted: plan.Accepted.ValueBool(),
		}
		err := UpdateSupportAssistTerms(ctx, client, terms)
		if err != nil {
			errStr := constants.UpdateSupportAssistStatusTermsErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error updating support assist terms",
				message,
			)
		}
	}

	// Update support assist status
	if !plan.SupportassistEnabled.IsNull() {
		status := powerscale.V16SupportassistStatusExtended{
			Enabled: plan.SupportassistEnabled.ValueBoolPointer(),
		}

		err := UpdateSupportAssistStatus(ctx, client, status)
		if err != nil {
			errStr := constants.UpdateSupportAssistStatusErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error updating support assist status",
				message,
			)
		}
	}

	// Update support assist settings
	supportAssistSettings := models.UpdateSupportassistSettings{}
	err := ReadFromState(ctx, &plan, &supportAssistSettings)
	if err != nil {
		errStr := constants.UpdateClusterIdentitySettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating support assist",
			fmt.Sprintf("Could not read support assist param with error: %s", message),
		)
	}

	supportAssistSettingsExtended := powerscale.V16SupportassistSettingsExtended{
		AutomaticCaseCreation: supportAssistSettings.AutomaticCaseCreation,
		EnableDownload:        supportAssistSettings.EnableDownload,
		EnableRemoteSupport:   supportAssistSettings.EnableRemoteSupport,
		Connection:            supportAssistSettings.Connection,
		Telemetry:             supportAssistSettings.Telemetry,
		Contact:               supportAssistSettings.Contact,
	}
	err = UpdateSupportAssistSettings(ctx, client, supportAssistSettingsExtended)
	if err != nil {
		errStr := constants.UpdateSupportAssistSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating support assist settings",
			message,
		)
	}

	// Connect with support assist backend if access key and pin are provided
	if !plan.Access_key.IsNull() && !plan.Pin.IsNull() {
		taskSettings := powerscale.V16SupportassistTaskItem{
			Source: "CONFIG",
			TaskParams: &powerscale.V16SupportassistTaskItemTaskParams{
				SubTask:   "provision",
				AccessKey: plan.Access_key.ValueStringPointer(),
				Pin:       plan.Pin.ValueStringPointer(),
			},
		}

		clusterVersion, err := GetClusterVersion(ctx, client)
		if err != nil {
			errStr := constants.ReadClusterErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error reading cluster version",
				message,
			)
			state, _ := ReadSupportAssistDetails(ctx, client, plan)
			return state, resp
		}

		var (
			taskCreate *powerscale.CreateTaskResponse
			response   *powerscale.V16SupportassistTaskId
		)

		if clusterVersion.Nodes[0].Release == "9.5.0.0" {
			taskCreate, err = CreateSupportAssistv16Task(ctx, client, taskSettings)
		} else {
			taskCreate, err = CreateSupportAssistv17Task(ctx, client, taskSettings)
		}

		if err != nil {
			errStr := constants.CreateSupportAssistTaskErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error creating support assist task",
				message,
			)
			state, _ := ReadSupportAssistDetails(ctx, client, plan)
			return state, resp
		}

		if clusterVersion.Nodes[0].Release == "9.5.0.0" {
			response, err = GetSupportAssistv16Task(ctx, client, taskCreate.TaskId)
		} else {
			response, err = GetSupportAssistv17Task(ctx, client, taskCreate.TaskId)
		}

		if err != nil {
			errStr := constants.GetSupportAssistTaskErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error getting support assist task",
				message,
			)
			state, _ = ReadSupportAssistDetails(ctx, client, plan)
			return state, resp
		}

		jobState := "COMPLETED"
		for *response.Tasks.State != jobState {
			time.Sleep(time.Second)
			if clusterVersion.Nodes[0].Release == "9.5.0.0" {
				response, err = GetSupportAssistv16Task(ctx, client, taskCreate.TaskId)
			} else {
				response, err = GetSupportAssistv17Task(ctx, client, taskCreate.TaskId)
			}
			if err != nil {
				errStr := constants.GetSupportAssistTaskErrorMsg + "with error: "
				message := GetErrorString(err, errStr)
				resp.AddError(
					"Error getting support assist task",
					message,
				)
			}
		}

		if response.Tasks.ErrorMsg != "" {
			resp.AddError(
				"Error while support assist provisioning",
				response.Tasks.ErrorMsg,
			)
		}
	}

	state, dig := ReadSupportAssistDetails(ctx, client, plan)
	if dig.HasError() {
		resp.AddError(
			"Error reading support assist details",
			"Could not read support assist details",
		)
	}
	return state, resp
}

// ReadSupportAssistDetails reads the support assist details
func ReadSupportAssistDetails(ctx context.Context, client *client.Client, plan models.SupportAssistModel) (state models.SupportAssistModel, resp diag.Diagnostics) {
	supportAssist, err := GetSupportAssist(ctx, client)
	if err != nil {
		errStr := constants.ReadSupportAssistErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error reading support assist details",
			message,
		)
		return state, resp
	}

	terms, err := GetSupportAssistTerms(ctx, client)
	if err != nil {
		errStr := constants.ReadSupportAssistTermsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error reading support assist term details",
			message,
		)
		return state, resp
	}

	// Populate the network pools
	networkPools := make([]string, 0)
	var connectionStruct *models.V16SupportassistSettingsConnection
	if supportAssist.Connection != nil {
		for _, pool := range supportAssist.Connection.NetworkPools {
			temp := pool.Subnet + ":" + pool.Pool
			networkPools = append(networkPools, temp)
		}

		// Populate the connection details
		connectionStruct = &models.V16SupportassistSettingsConnection{
			Mode:             supportAssist.Connection.Mode,
			GatewayEndpoints: supportAssist.Connection.GatewayEndpoints,
			NetworkPools:     networkPools,
		}

	}

	copyStruct := &models.V16SupportassistSettingsCustomised{
		AutomaticCaseCreation: supportAssist.AutomaticCaseCreation,
		Connections:           connectionStruct,
		Contact:               supportAssist.Contact,
		EnableDownload:        supportAssist.EnableDownload,
		EnableRemoteSupport:   supportAssist.EnableRemoteSupport,
		SupportassistEnabled:  supportAssist.SupportassistEnabled,
		Telemetry:             supportAssist.Telemetry,
	}

	err = CopyFieldsToNonNestedModel(ctx, copyStruct, &state)
	if err != nil {
		resp.AddError(
			"Unable to update support assist state",
			fmt.Sprintf("Unable to update support assist state with error %s", err),
		)
		return state, resp
	}

	if copyStruct.Connections != nil && (len(copyStruct.Connections.GatewayEndpoints) == 0 || len(copyStruct.Connections.NetworkPools) == 0) {
		var settings models.Connection

		diag := state.Connection.As(ctx, &settings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		if diag.HasError() {
			return state, diag
		}
		// settings.Mode = types.StringValue(*copyStruct.Connections.Mode)
		if len(copyStruct.Connections.NetworkPools) == 0 {
			settings.NetworkPools, _ = types.ListValue(types.StringType, []attr.Value{})
		}

		if len(copyStruct.Connections.GatewayEndpoints) == 0 {
			gatewayElemType := types.ObjectType{
				AttrTypes: GetGatewayEndpointType(),
			}

			settings.GatewayEndpoints, diag = types.ListValue(gatewayElemType, []attr.Value{})
			if diag.HasError() {
				return state, diag
			}
		}
		obj := map[string]attr.Value{
			"mode":              settings.Mode,
			"network_pools":     settings.NetworkPools,
			"gateway_endpoints": settings.GatewayEndpoints,
		}

		// settings.NetworkPools, diag = types.ListValue(types.StringType, copyStruct.Connections.NetworkPools) //copyStruct.Connections.NetworkPools
		state.Connection, diag = types.ObjectValue(GetConnectionType(), obj)
		if diag.HasError() {
			return state, diag
		}
	}

	if !plan.Access_key.IsNull() && !plan.Pin.IsNull() {
		state.Access_key = plan.Access_key
		state.Pin = plan.Pin
	}
	state.ID = types.StringValue("support_assist")
	state.Accepted = types.BoolValue(terms.Terms.Accepted)
	return state, nil
}

// GetGatewayEndpointType returns the gateway type
func GetGatewayEndpointType() map[string]attr.Type {
	return map[string]attr.Type{
		"port":         types.Int64Type,
		"validate_ssl": types.BoolType,
		"priority":     types.Int64Type,
		"use_proxy":    types.BoolType,
		"host":         types.StringType,
		"enabled":      types.BoolType,
	}
}

// GetConnectionType returns the connection type
func GetConnectionType() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":              types.StringType,
		"network_pools":     types.ListType{ElemType: types.StringType},
		"gateway_endpoints": types.ListType{ElemType: types.ObjectType{AttrTypes: GetGatewayEndpointType()}},
	}
}
