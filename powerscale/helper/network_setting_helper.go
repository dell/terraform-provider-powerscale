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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateNetworkSettingState updates resource state.
func UpdateNetworkSettingState(ctx context.Context, settingState *models.NetworkSettingModel, settingResponse *powerscale.V12NetworkExternalSettings) {

	if len(settingState.TCPPorts.Elements()) != len(settingResponse.TcpPorts) {
		var portAttrs []attr.Value
		for _, port := range settingResponse.TcpPorts {
			portAttrs = append(portAttrs, types.Int64Value(port))
		}
		settingState.TCPPorts, _ = types.ListValue(types.Int64Type, portAttrs)
	}

	settingState.SBREnabled = types.BoolValue(settingResponse.Sbr)
	settingState.DefaultGroupnet = types.StringValue(settingResponse.DefaultGroupnet)
	settingState.SCRebalanceDelay = types.Int64Value(settingResponse.ScRebalanceDelay)
	settingState.ID = types.StringValue("network_settings")
}

// GetNetworkSetting returns network settings detail.
func GetNetworkSetting(ctx context.Context, client *client.Client) (*powerscale.V12NetworkExternalSettings, error) {

	getParams := client.PscaleOpenAPIClient.NetworkApi.GetNetworkv12NetworkExternal(ctx)
	result, _, err := getParams.Execute()
	if err != nil {
		errStr := constants.ReadNetworkSettingErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting network settings: %s", message)
	}

	return result.Settings, nil
}

// UpdateNetworkSetting Updates network setting.
func UpdateNetworkSetting(ctx context.Context, client *client.Client, state *models.NetworkSettingModel, plan *models.NetworkSettingModel) (diags diag.Diagnostics) {
	updateParam := client.PscaleOpenAPIClient.NetworkApi.UpdateNetworkv12NetworkExternal(ctx)

	body := &powerscale.V12NetworkExternalExtended{}

	if !plan.SBREnabled.IsNull() && !plan.SBREnabled.IsUnknown() && !state.SBREnabled.Equal(plan.SBREnabled) {
		body.Sbr = plan.SBREnabled.ValueBoolPointer()
	}
	if !plan.SCRebalanceDelay.IsNull() && !plan.SCRebalanceDelay.IsUnknown() && !state.SCRebalanceDelay.Equal(plan.SCRebalanceDelay) {
		body.ScRebalanceDelay = plan.SCRebalanceDelay.ValueInt64Pointer()
	}
	if !plan.TCPPorts.IsNull() && !plan.TCPPorts.IsUnknown() && !state.TCPPorts.Equal(plan.TCPPorts) {
		var ports []int64
		if diags = plan.TCPPorts.ElementsAs(ctx, &ports, false); diags.HasError() {
			return
		}
		body.TcpPorts = ports
	}

	updateParam = updateParam.V12NetworkExternal(*body)
	if _, err := updateParam.Execute(); err != nil {
		errStr := constants.UpdateNetworkSettingErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(fmt.Sprintf("error updating network settings: %s", message), err.Error())
	}

	return
}
