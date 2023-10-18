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

package provider

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SmartPoolSettingDataSource{}

// NewSmartPoolSettingDataSource creates a new data source.
func NewSmartPoolSettingDataSource() datasource.DataSource {
	return &SmartPoolSettingDataSource{}
}

// SmartPoolSettingDataSource defines the data source implementation.
type SmartPoolSettingDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d SmartPoolSettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smartpool_settings"
}

// Schema describes the data source arguments.
func (d SmartPoolSettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the SmartPools settings from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details / for further processing in resource block. " +
			"PowerScale SmartPools settings provide the ability to configure SmartPools on the cluster.",
		Description: "This datasource is used to query the SmartPools settings from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details / for further processing in resource block. " +
			"PowerScale SmartPools settings provide the ability to configure SmartPools on the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Id of SmartPool settings.",
				MarkdownDescription: "Id of SmartPool settings.",
				Computed:            true,
			},
			"settings": helper.GetSmartPoolSettingsSchema(),
		},
	}
}

// Configure configures the data source.
func (d *SmartPoolSettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d SmartPoolSettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.SmartPoolSettingsDataSource

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read and map SmartPool setting state
	settings, err := helper.GetSmartPoolSettings(ctx, d.client)
	if err != nil {
		errStr := constants.ReadSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading SmartPool settings", message)
		return
	}

	var sts Settings
	switch v := settings.(type) {
	case *powerscale.V16StoragepoolSettings:
		sts, _ = settings.(*powerscale.V16StoragepoolSettings).GetSettingsOk()
		err := helper.CopyFields(ctx, settings, &state)
		if err != nil {
			resp.Diagnostics.AddError("Failed to map SmartPool settings fields", err.Error())
			return
		}
	case *powerscale.V5StoragepoolSettings:
		sts, _ = settings.(*powerscale.V5StoragepoolSettings).GetSettingsOk()
		err := helper.CopyFields(ctx, settings, &state)
		if err != nil {
			resp.Diagnostics.AddError("Failed to map SmartPool settings fields", err.Error())
			return
		}
	default:
		tflog.Error(ctx, fmt.Sprintf("Unknown type %s", v))
		resp.Diagnostics.AddError("Failed to parse StoragePool Settings.", fmt.Sprintf("Unknown type %s", v))
		return
	}

	state.ID = types.StringValue("smartpool_settings_datasource")
	// Compute the following two values that align with OneFS UI
	if sts.GetAutomaticallyManageIoOptimization() == "none" {
		state.Settings.ManageIoOptimization = types.BoolValue(false)
		state.Settings.ManageIoOptimizationApplyToFiles = types.BoolValue(false)
	} else if sts.GetAutomaticallyManageIoOptimization() == "files_at_default" {
		state.Settings.ManageIoOptimization = types.BoolValue(true)
		state.Settings.ManageIoOptimizationApplyToFiles = types.BoolValue(false)
	} else if sts.GetAutomaticallyManageIoOptimization() == "all" {
		state.Settings.ManageIoOptimization = types.BoolValue(true)
		state.Settings.ManageIoOptimizationApplyToFiles = types.BoolValue(true)
	}

	if sts.GetAutomaticallyManageProtection() == "none" {
		state.Settings.ManageProtection = types.BoolValue(false)
		state.Settings.ManageProtectionApplyToFiles = types.BoolValue(false)
	} else if sts.GetAutomaticallyManageProtection() == "files_at_default" {
		state.Settings.ManageProtection = types.BoolValue(true)
		state.Settings.ManageProtectionApplyToFiles = types.BoolValue(false)
	} else if sts.GetAutomaticallyManageProtection() == "all" {
		state.Settings.ManageProtection = types.BoolValue(true)
		state.Settings.ManageProtectionApplyToFiles = types.BoolValue(true)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read SmartPool Settings data source ")
}

// Settings The interface of various StoragepoolSettingsSettings.
type Settings interface {
	GetAutomaticallyManageIoOptimization() string
	GetAutomaticallyManageProtection() string
}
