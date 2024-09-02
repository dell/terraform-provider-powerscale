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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"
)

// UpdateSyncIQGlobalSettings updates the SyncIQ global settings
func UpdateSyncIQGlobalSettings(ctx context.Context, client *client.Client, edit powerscale.V16SyncSettingsExtended) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv16SyncSettings(ctx).V16SyncSettings(edit).Execute()
	return err
}

// GetSyncIQGlobalSettings fetches the SyncIQ global settings
func GetSyncIQGlobalSettings(ctx context.Context, client *client.Client) (*powerscale.V16SyncSettings, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv16SyncSettings(ctx).Execute()
	return resp, err
}

// ManageSyncIQGlobalSettings does all the update functionality for SyncIQ Global settings
func ManageSyncIQGlobalSettings(ctx context.Context, plan models.SyncIQGlobalSettingsModel, state *models.SyncIQGlobalSettingsModel, client *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	var toUpdate powerscale.V16SyncSettingsExtended

	err := ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateSyncIQGlobalSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error updating synciq global settings",
			fmt.Sprintf("Could not read synciq global setting param with error: %s", message),
		)
		return diags
	}

	if !plan.ReportEmail.IsUnknown() {
		diags.Append(plan.ReportEmail.ElementsAs(ctx, &toUpdate.ReportEmail, true)...)
	}

	if toUpdate.SourceNetwork != nil {
		if (toUpdate.SourceNetwork.Subnet == "" && toUpdate.SourceNetwork.Pool != "") || (toUpdate.SourceNetwork.Subnet != "" && toUpdate.SourceNetwork.Pool == "") {
			diags.AddError(
				"Valid value for both subnet and pool needs to be provided",
				"Please either provide value for both subnet and pool or else keep both as empty",
			)
			return diags
		}
	}

	err = UpdateSyncIQGlobalSettings(ctx, client, toUpdate)
	if err != nil {
		errStr := constants.UpdateSyncIQGlobalSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error updating synciq global settings",
			message,
		)
		return diags
	}

	globalSetting, err := GetSyncIQGlobalSettings(ctx, client)
	if err != nil {
		errStr := constants.ReadSyncIQGlobalSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading synciq global settings",
			message,
		)
		return diags
	}

	err = CopyFields(ctx, globalSetting.Settings, state)
	if err != nil {
		diags.AddError(
			"Error copying fields of synciq global settings resource",
			err.Error(),
		)
		return diags
	}

	sourceNetwork, emailObj, diags2 := GetSourceNetworkAndEmail(globalSetting)
	if diags2.HasError() {
		diags.AddError(
			"Error copying fields of synciq global settings resource",
			"Error occurred while copying source network or report email",
		)
		return diags
	}
	state.SourceNetwork = sourceNetwork
	state.ReportEmail = emailObj

	return diags
}

// ManageReadSyncIQGlobalSettings does the read functionality for SyncIQ global settings
func ManageReadSyncIQGlobalSettings(ctx context.Context, state *models.SyncIQGlobalSettingsModel, client *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	globalSettings, err := GetSyncIQGlobalSettings(ctx, client)
	if err != nil {
		errStr := constants.ReadSyncIQGlobalSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading synciq global settings",
			message,
		)
		return diags
	}
	err = CopyFields(ctx, globalSettings.Settings, state)
	if err != nil {
		diags.AddError(
			"Error copying fields of synciq global settings resource",
			err.Error(),
		)
		return diags
	}

	sourceNetwork, emailObj, diags := GetSourceNetworkAndEmail(globalSettings)
	if diags.HasError() {
		diags.AddError(
			"Error copying fields of synciq global settings resource",
			"Error occurred while copying source network or report email",
		)
		return diags
	}
	state.SourceNetwork = sourceNetwork
	state.ReportEmail = emailObj

	return diags
}

// ManageReadDataSourceSyncIQGlobalSettings does the read functionality for SyncIQ global settings datasource
func ManageReadDataSourceSyncIQGlobalSettings(ctx context.Context, state *models.SyncIQGlobalSettingsDataSourceModel, client *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	globalSettings, err := GetSyncIQGlobalSettings(ctx, client)
	if err != nil {
		errStr := constants.ReadSyncIQGlobalSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading synciq global settings",
			message,
		)
		return diags
	}
	err = CopyFields(ctx, globalSettings.Settings, state)
	if err != nil {
		diags.AddError(
			"Error copying fields of synciq global settings resource",
			err.Error(),
		)
		return diags
	}

	sourceNetwork, emailObj, diags := GetSourceNetworkAndEmail(globalSettings)
	if diags.HasError() {
		diags.AddError(
			"Error copying fields of synciq global settings resource",
			"Error occurred while copying source network or report email",
		)
		return diags
	}
	state.SourceNetwork = sourceNetwork
	state.ReportEmail = emailObj

	return diags
}

// GetSourceNetworkAndEmail returns the object and set value for source network and report email attributes.
func GetSourceNetworkAndEmail(globalSetting *powerscale.V16SyncSettings) (basetypes.ObjectValue, basetypes.SetValue, diag.Diagnostics) {
	var sourceNetworkObject basetypes.ObjectValue
	var diags diag.Diagnostics
	var emailObj basetypes.SetValue
	if globalSetting.Settings.SourceNetwork != nil {
		sourceNetworkObjectType := map[string]attr.Type{
			"subnet": types.StringType,
			"pool":   types.StringType,
		}

		sourceNetworkElemMap := map[string]attr.Value{
			"subnet": types.StringValue(globalSetting.Settings.SourceNetwork.Subnet),
			"pool":   types.StringValue(globalSetting.Settings.SourceNetwork.Pool),
		}

		sourceNetworkObject, diags = types.ObjectValue(sourceNetworkObjectType, sourceNetworkElemMap)
		if diags.HasError() {
			diags.Append(diags...)
		}
	} else {
		sourceNetworkObjectType := map[string]attr.Type{
			"subnet": types.StringType,
			"pool":   types.StringType,
		}

		SourceNetworkElemMap := map[string]attr.Value{
			"subnet": types.StringValue(""),
			"pool":   types.StringValue(""),
		}

		sourceNetworkObject, diags = types.ObjectValue(sourceNetworkObjectType, SourceNetworkElemMap)
		if diags.HasError() {
			diags.Append(diags...)
		}
	}
	emails := []attr.Value{}
	for _, email := range globalSetting.Settings.ReportEmail {
		emails = append(emails, types.StringValue(email))
	}
	emailObj, diags = types.SetValue(types.StringType, emails)
	if diags.HasError() {
		diags.Append(diags...)
	}

	return sourceNetworkObject, emailObj, diags
}
