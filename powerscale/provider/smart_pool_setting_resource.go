/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource              = &SmartPoolSettingResource{}
	_ resource.ResourceWithConfigure = &SmartPoolSettingResource{}
)

// NewSmartPoolSettingResource creates a new resource.
func NewSmartPoolSettingResource() resource.Resource {
	return &SmartPoolSettingResource{}
}

// SmartPoolSettingResource defines the resource implementation.
type SmartPoolSettingResource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (r *SmartPoolSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smartpool_settings"
}

// Schema describes the data source arguments.
func (r *SmartPoolSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `This resource is used to manage the SmartPools Settings of PowerScale Array. We can Create, Update and Delete the SmartPools Settings using this resource.  
Note that, SmartPools Settings is the native functionality of PowerScale. When creating the resource, we actually load SmartPools Settings from PowerScale to the resource.`,
		Description: `This resource is used to manage the SmartPools Settings of PowerScale Array. We can Create, Update and Delete the SmartPools Settings using this resource.  
Note that, SmartPools Settings is the native functionality of PowerScale. When creating the resource, we actually load SmartPools Settings from PowerScale to the resource.`,
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
				Optional:            true,
			},
			"manage_io_optimization_apply_to_files": schema.BoolAttribute{
				Description:         "Apply to files with manually-managed I/O optimization settings.",
				MarkdownDescription: "Apply to files with manually-managed I/O optimization settings.",
				Computed:            true,
				Optional:            true,
			},
			"manage_protection": schema.BoolAttribute{
				Description:         "Manage protection settings.",
				MarkdownDescription: "Manage protection settings.",
				Computed:            true,
				Optional:            true,
			},
			"manage_protection_apply_to_files": schema.BoolAttribute{
				Description:         "Apply to files with manually-managed protection.",
				MarkdownDescription: "Apply to files with manually-managed protection.",
				Computed:            true,
				Optional:            true,
			},
			"global_namespace_acceleration_enabled": schema.BoolAttribute{
				Description:         "Enable global namespace acceleration.",
				MarkdownDescription: "Enable global namespace acceleration.",
				Computed:            true,
				Optional:            true,
			},
			"global_namespace_acceleration_state": schema.StringAttribute{
				Description:         "Whether or not namespace operation optimizations are currently in effect.",
				MarkdownDescription: "Whether or not namespace operation optimizations are currently in effect.",
				Computed:            true,
				Optional:            true,
			},
			"protect_directories_one_level_higher": schema.BoolAttribute{
				Description:         "Increase directory protection to a higher requested protection than its contents.",
				MarkdownDescription: "Increase directory protection to a higher requested protection than its contents.",
				Computed:            true,
				Optional:            true,
			},
			"spillover_enabled": schema.BoolAttribute{
				Description:         "Enable global spillover.",
				MarkdownDescription: "Enable global spillover.",
				Computed:            true,
				Optional:            true,
			},
			"spillover_target": schema.SingleNestedAttribute{
				Description:         "Spillover data target.",
				MarkdownDescription: "Spillover data target.",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "Target pool name if target specified as storagepool, otherwise empty string.",
						MarkdownDescription: "Target pool name if target specified as storagepool, otherwise empty string.",
						Computed:            true,
						Optional:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Type of target pool. Acceptable values: storagepool, anywhere",
						MarkdownDescription: "Type of target pool. Acceptable values: storagepool, anywhere",
						Computed:            true,
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("storagepool", "anywhere"),
						},
					},
				},
			},
			"ssd_l3_cache_default_enabled": schema.BoolAttribute{
				Description:         "Use SSDs as L3 cache by default for new node pools",
				MarkdownDescription: "Use SSDs as L3 cache by default for new node pools.",
				Computed:            true,
				Optional:            true,
			},
			"ssd_qab_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of QAB blocks to place on SSDs. Acceptable values: one, all",
				MarkdownDescription: "Controls number of mirrors of QAB blocks to place on SSDs. Acceptable values: one, all",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("one", "all"),
				},
			},
			"ssd_system_btree_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of system B-tree blocks to place on SSDs. Acceptable values: one, all",
				MarkdownDescription: "Controls number of mirrors of system B-tree blocks to place on SSDs. Acceptable values: one, all",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("one", "all"),
				},
			},
			"ssd_system_delta_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of system delta blocks to place on SSDs. Acceptable values: one, all",
				MarkdownDescription: "Controls number of mirrors of system delta blocks to place on SSDs. Acceptable values: one, all",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("one", "all"),
				},
			},
			"virtual_hot_spare_deny_writes": schema.BoolAttribute{
				Description:         "Deny data writes to reserved disk space",
				MarkdownDescription: "Deny data writes to reserved disk space",
				Computed:            true,
				Optional:            true,
			},
			"virtual_hot_spare_hide_spare": schema.BoolAttribute{
				Description:         "Subtract the space reserved for the virtual hot spare when calculating available free space",
				MarkdownDescription: "Subtract the space reserved for the virtual hot spare when calculating available free space",
				Computed:            true,
				Optional:            true,
			},
			"virtual_hot_spare_limit_drives": schema.Int64Attribute{
				Description:         "The number of drives to reserve for the virtual hot spare, from 0-4.",
				MarkdownDescription: "The number of drives to reserve for the virtual hot spare, from 0-4.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
					int64validator.AtMost(20),
				},
			},
			"virtual_hot_spare_limit_percent": schema.Int64Attribute{
				Description:         "The percent space to reserve for the virtual hot spare, from 0-20.",
				MarkdownDescription: "The percent space to reserve for the virtual hot spare, from 0-20.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
					int64validator.AtMost(20),
				},
			},
			"default_transfer_limit_state": schema.StringAttribute{
				Description:         "How the default transfer limit value is applied. Only available for PowerScale 9.5 and above.",
				MarkdownDescription: "How the default transfer limit value is applied. Only available for PowerScale 9.5 and above.",
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
func (r *SmartPoolSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *SmartPoolSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating SmartPoolSettings resource...")

	var plan models.SmartPoolSettingsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	var planBackUp models.SmartPoolSettingsResource
	diags = req.Plan.Get(ctx, &planBackUp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read and map SmartPool setting state
	settings, err := helper.GetSmartPoolSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading SmartPool settings", message)
		return
	}

	summary, detail := helper.UpdateSmartPoolSettingsResourceModel(ctx, settings, &plan)
	if summary != "" && detail != "" {
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	summary, detail = r.updateResource(ctx, &planBackUp)
	if len(summary) > 0 && len(detail) > 0 {
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planBackUp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Done with Create SmartPoolSettings data source ")
}

// Read reads the resource state.
func (r *SmartPoolSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading SmartPoolSettings resource")

	var state models.SmartPoolSettingsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := helper.GetSmartPoolSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading SmartPool settings", message)
		return
	}

	summary, detail := helper.UpdateSmartPoolSettingsResourceModel(ctx, settings, &state)
	if summary != "" && detail != "" {
		resp.Diagnostics.AddError(summary, detail)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Done with Read SmartPool Settings resource")
}

func (r *SmartPoolSettingResource) updateResource(ctx context.Context, plan *models.SmartPoolSettingsResource) (string, string) {
	if (!plan.ManageProtectionApplyToFiles.IsUnknown()) &&
		!plan.ManageProtectionApplyToFiles.IsNull() &&
		plan.ManageProtectionApplyToFiles.ValueBool() {
		if plan.ManageProtection.IsNull() || plan.ManageProtection.IsUnknown() || !plan.ManageProtection.ValueBool() {
			return "Input validation failed.",
				"manage_protection should be set to true when manage_protection_apply_to_files is true."
		}
	}

	if (!plan.ManageIoOptimizationApplyToFiles.IsUnknown()) &&
		!plan.ManageIoOptimizationApplyToFiles.IsNull() &&
		plan.ManageIoOptimizationApplyToFiles.ValueBool() {
		if plan.ManageIoOptimization.IsNull() || plan.ManageIoOptimization.IsUnknown() || !plan.ManageIoOptimization.ValueBool() {
			return "Input validation failed.",
				"manage_io_optimization should be set to true when manage_io_optimization_apply_to_files is true."
		}
	}

	// update smartpool settings
	if err := helper.UpdateSmartPoolSettings(ctx, r.client, plan); err != nil {
		errStr := constants.UpdateSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		return "Error updating SmartPool settings", message
	}

	// Read and map SmartPool setting state
	settings, err := helper.GetSmartPoolSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		return "Error reading SmartPool settings", message
	}

	summary, detail := helper.UpdateSmartPoolSettingsResourceModel(ctx, settings, plan)
	if summary != "" && detail != "" {
		return summary, detail
	}

	return "", ""
}

// Update updates the resource state.
func (r *SmartPoolSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating SmartPoolSettings resource...")
	// Read Terraform plan into the model
	var plan models.SmartPoolSettingsResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	summary, detail := r.updateResource(ctx, &plan)
	if len(summary) > 0 && len(detail) > 0 {
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Done with Update SmartPool Settings resource")
}

// Delete deletes the resource.
func (r *SmartPoolSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting SmartPoolSettings resource")
	var state models.SmartPoolSettingsResource

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// SmartPool settings is the native functionality that cannot be deleted, so just remove state
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete SmartPool SmartPoolSettings")
}

// ImportState imports the resource state.
func (r *SmartPoolSettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing SmartPoolSettings resource")
	var state models.SmartPoolSettingsResource

	// Read and map SmartPool setting state
	settings, err := helper.GetSmartPoolSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading SmartPool settings", message)
		return
	}

	summary, detail := helper.UpdateSmartPoolSettingsResourceModel(ctx, settings, &state)
	if summary != "" && detail != "" {
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	tflog.Info(ctx, "Done with Import SmartPoolSettings resource")
}
