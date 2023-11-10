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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Description:         "Id of SmartPools settings. Readonly. Fixed value of \"smartpools_settings\"",
				MarkdownDescription: "Id of SmartPools settings. Readonly. Fixed value of \"smartpools_settings\"",
				Optional:            false,
				Required:            false,
				Computed:            true,
			},
			"manage_io_optimization": schema.BoolAttribute{
				Description:         "Manage I/O optimization settings.",
				MarkdownDescription: "Manage I/O optimization settings.",
				Computed:            true,
			},
			"manage_io_optimization_apply_to_files": schema.BoolAttribute{
				Description:         "Apply to files with manually-managed I/O optimization settings.",
				MarkdownDescription: "Apply to files with manually-managed I/O optimization settings.",
				Computed:            true,
			},
			"manage_protection": schema.BoolAttribute{
				Description:         "Manage protection settings.",
				MarkdownDescription: "Manage protection settings.",
				Computed:            true,
			},
			"manage_protection_apply_to_files": schema.BoolAttribute{
				Description:         "Apply to files with manually-managed protection.",
				MarkdownDescription: "Apply to files with manually-managed protection.",
				Computed:            true,
			},
			"global_namespace_acceleration_enabled": schema.BoolAttribute{
				Description:         "Enable global namespace acceleration.",
				MarkdownDescription: "Enable global namespace acceleration.",
				Computed:            true,
			},
			"global_namespace_acceleration_state": schema.StringAttribute{
				Description:         "Whether or not namespace operation optimizations are currently in effect.",
				MarkdownDescription: "Whether or not namespace operation optimizations are currently in effect.",
				Computed:            true,
			},
			"protect_directories_one_level_higher": schema.BoolAttribute{
				Description:         "Increase directory protection to a higher requested protection than its contents.",
				MarkdownDescription: "Increase directory protection to a higher requested protection than its contents.",
				Computed:            true,
			},
			"spillover_enabled": schema.BoolAttribute{
				Description:         "Enable global spillover.",
				MarkdownDescription: "Enable global spillover.",
				Computed:            true,
			},
			"spillover_target": schema.SingleNestedAttribute{
				Description:         "Spillover data target.",
				MarkdownDescription: "Spillover data target.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description:         "Target pool ID if target specified, otherwise null.",
						MarkdownDescription: "Target pool ID if target specified, otherwise null.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Target pool name if target specified, otherwise null.",
						MarkdownDescription: "Target pool name if target specified, otherwise null.",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Type of target pool.",
						MarkdownDescription: "Type of target pool.",
						Computed:            true,
					},
				},
			},
			"ssd_l3_cache_default_enabled": schema.BoolAttribute{
				Description:         "Use SSDs as L3 cache by default for new node pools",
				MarkdownDescription: "Use SSDs as L3 cache by default for new node pools.",
				Computed:            true,
			},
			"ssd_qab_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of QAB blocks to place on SSDs.",
				MarkdownDescription: "Controls number of mirrors of QAB blocks to place on SSDs.",
				Computed:            true,
			},
			"ssd_system_btree_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of system B-tree blocks to place on SSDs.",
				MarkdownDescription: "Controls number of mirrors of system B-tree blocks to place on SSDs.",
				Computed:            true,
			},
			"ssd_system_delta_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of system delta blocks to place on SSDs.",
				MarkdownDescription: "Controls number of mirrors of system delta blocks to place on SSDs.",
				Computed:            true,
			},
			"virtual_hot_spare_deny_writes": schema.BoolAttribute{
				Description:         "Deny data writes to reserved disk space",
				MarkdownDescription: "Deny data writes to reserved disk space",
				Computed:            true,
			},
			"virtual_hot_spare_hide_spare": schema.BoolAttribute{
				Description:         "Subtract the space reserved for the virtual hot spare when calculating available free space",
				MarkdownDescription: "Subtract the space reserved for the virtual hot spare when calculating available free space",
				Computed:            true,
			},
			"virtual_hot_spare_limit_drives": schema.Int64Attribute{
				Description:         "The number of drives to reserve for the virtual hot spare, from 0-4.",
				MarkdownDescription: "The number of drives to reserve for the virtual hot spare, from 0-4.",
				Computed:            true,
			},
			"virtual_hot_spare_limit_percent": schema.Int64Attribute{
				Description:         "The percent space to reserve for the virtual hot spare, from 0-20.",
				MarkdownDescription: "The percent space to reserve for the virtual hot spare, from 0-20.",
				Computed:            true,
			},
			"default_transfer_limit_state": schema.StringAttribute{
				Description:         "How the default transfer limit value is applied. Only available for PowerScale 9.5 and above.",
				MarkdownDescription: "How the default transfer limit value is applied.Only available for PowerScale 9.5 and above.",
				Computed:            true,
				Optional:            true,
			},
			"default_transfer_limit_pct": schema.NumberAttribute{
				Description:         "Applies to all storagepools that fall back on the default transfer limit. Stop moving files to this pool when this limit is met. The value must be between 0 and 100. Only available for PowerScale 9.5 and above.",
				MarkdownDescription: "Applies to all storagepools that fall back on the default transfer limit. Stop moving files to this pool when this limit is met. The value must be between 0 and 100. Only available for PowerScale 9.5 and above.",
				Computed:            true,
				Optional:            true,
			},
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

	summary, detail := helper.UpdateSmartPoolSettingsDatasourceModel(ctx, settings, &state)
	if summary != "" && detail != "" {
		resp.Diagnostics.AddError(summary, detail)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read SmartPoolSettings data source ")
}
