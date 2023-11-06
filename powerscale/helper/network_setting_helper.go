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

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateNetworkSettingDataSourceState updates datasource state.
func UpdateNetworkSettingDataSourceState(ctx context.Context, settingState *models.NetworkSettingModel, settingResponse *powerscale.V12NetworkExternalSettings) (err error) {

	if err = CopyFields(ctx, settingResponse, settingState); err != nil {
		return
	}

	settingState.SBREnabled = types.BoolValue(settingResponse.Sbr)
	settingState.ID = types.StringValue("network_setting_datasource")

	return
}

// GetNetworkSetting returns network setting detail.
func GetNetworkSetting(ctx context.Context, client *client.Client) (*powerscale.V12NetworkExternalSettings, error) {

	getParams := client.PscaleOpenAPIClient.NetworkApi.GetNetworkv12NetworkExternal(ctx)
	result, _, err := getParams.Execute()
	if err != nil {
		errStr := constants.ReadNetworkSettingErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting network setting: %s", message)
	}

	return result.Settings, nil
}
