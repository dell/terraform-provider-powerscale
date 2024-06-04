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

package provider

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"strings"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &SmbServerSettingsResource{}
	_ resource.ResourceWithConfigure   = &SmbServerSettingsResource{}
	_ resource.ResourceWithImportState = &SmbServerSettingsResource{}
)

// NewSmbServerSettingsResource is a helper function to simplify the provider implementation.
func NewSmbServerSettingsResource() resource.Resource {
	return &SmbServerSettingsResource{}
}

// SmbServerSettingsResource is the resource implementation.
type SmbServerSettingsResource struct {
	client *client.Client
}

// Metadata defines the resource type name.
func (r *SmbServerSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_server_settings"
}

// Schema defines the schema for the resource.
func (r *SmbServerSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `This resource is used to manage the SMB Server Settings of PowerScale Array. We can Create, Update and Delete the SMB Server Settings using this resource.  
		Note that, SMB Server Settings is the native functionality of PowerScale. When creating the resource, we actually load SMB Server Settings from PowerScale to the resource.`,
		Description: `This resource is used to manage the SMB Server Settings of PowerScale Array. We can Create, Update and Delete the SMB Server Settings using this resource.  
		Note that, SMB Server Settings is the native functionality of PowerScale. When creating the resource, we actually load SMB Server Settings from PowerScale to the resource.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of SMB Server Settings.",
				MarkdownDescription: "ID of SMB Server Settings.",
			},
			"scope": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If specified as \"effective\" or not specified, all fields are returned.  If specified as \"user\", only fields with non-default values are shown.  If specified as \"default\", the original values are returned.",
				MarkdownDescription: "If specified as \"effective\" or not specified, all fields are returned.  If specified as \"user\", only fields with non-default values are shown.  If specified as \"default\", the original values are returned.",
			},
			"support_smb2": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Support the SMB2 protocol on the server.",
				MarkdownDescription: "Support the SMB2 protocol on the server.",
			},
			"support_smb3_encryption": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Support the SMB3 encryption on the server.",
				MarkdownDescription: "Support the SMB3 encryption on the server.",
			},
			"audit_logon": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specify the level of logon audit events to log.",
				MarkdownDescription: "Specify the level of logon audit events to log.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"all",
						"success",
						"failure",
						"none",
					),
				},
			},
			"srv_cpu_multiplier": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specify the number of SRV service worker threads per CPU.",
				MarkdownDescription: "Specify the number of SRV service worker threads per CPU.",
			},
			"server_string": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Provides a description of the server.",
				MarkdownDescription: "Provides a description of the server.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 511),
				},
			},
			"service": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specify whether service is enabled.",
				MarkdownDescription: "Specify whether service is enabled.",
			},
			"support_multichannel": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Support multichannel.",
				MarkdownDescription: "Support multichannel.",
			},
			"dot_snap_visible_root": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Show the .snapshot directory in the root of a share.",
				MarkdownDescription: "Show the .snapshot directory in the root of a share.",
			},
			"onefs_num_workers": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Set the maximum number of OneFS driver worker threads.",
				MarkdownDescription: "Set the maximum number of OneFS driver worker threads.",
			},
			"srv_num_workers": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Set the maximum number of SRV service worker threads.",
				MarkdownDescription: "Set the maximum number of SRV service worker threads.",
			},
			"enable_security_signatures": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Indicates whether the server supports signed SMB packets.",
				MarkdownDescription: "Indicates whether the server supports signed SMB packets.",
			},
			"guest_user": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the fully-qualified user to use for guest access.",
				MarkdownDescription: "Specifies the fully-qualified user to use for guest access.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 511),
				},
			},
			"require_security_signatures": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Indicates whether the server requires signed SMB packets.",
				MarkdownDescription: "Indicates whether the server requires signed SMB packets.",
			},
			"onefs_cpu_multiplier": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specify the number of OneFS driver worker threads per CPU.",
				MarkdownDescription: "Specify the number of OneFS driver worker threads per CPU.",
			},
			"dot_snap_accessible_child": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Allow access to .snapshot directories in share subdirectories.",
				MarkdownDescription: "Allow access to .snapshot directories in share subdirectories.",
			},
			"access_based_share_enum": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Only enumerate files and folders the requesting user has access to.",
				MarkdownDescription: "Only enumerate files and folders the requesting user has access to.",
			},
			"audit_fileshare": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specify level of file share audit events to log.",
				MarkdownDescription: "Specify level of file share audit events to log.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"all",
						"success",
						"failure",
						"none",
					),
				},
			},
			"dot_snap_visible_child": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Show .snapshot directories in share subdirectories.",
				MarkdownDescription: "Show .snapshot directories in share subdirectories.",
			},
			"dot_snap_accessible_root": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Allow access to the .snapshot directory in the root of the share.",
				MarkdownDescription: "Allow access to the .snapshot directory in the root of the share.",
			},
			"server_side_copy": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Enable Server Side Copy.",
				MarkdownDescription: "Enable Server Side Copy.",
			},
			"reject_unencrypted_access": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If SMB3 encryption is enabled, reject unencrypted access from clients.",
				MarkdownDescription: "If SMB3 encryption is enabled, reject unencrypted access from clients.",
			},
			"support_netbios": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Support NetBIOS.",
				MarkdownDescription: "Support NetBIOS.",
			},
			"ignore_eas": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specify whether to ignore EAs on files.",
				MarkdownDescription: "Specify whether to ignore EAs on files.",
			},
		},
	}
}

// Configure configures the resource.
func (r *SmbServerSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = pscaleClient
}

// Create creates the resource.
func (r *SmbServerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Started creating smb server settings")

	var plan models.SmbServerSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := plan.Scope.ValueString()
	if scope == "" {
		scope = "effective"
	}

	var toUpdate powerscale.V6SmbSettingsGlobalSettings
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating smb server settings",
			fmt.Sprintf("Could not read smb server settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateSmbServerSettings(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating smb server settings",
			message,
		)
		return
	}

	settings, err := helper.GetSmbServerSettings(ctx, r.client, scope)
	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading smb server settings",
			message,
		)
		return
	}

	var state models.SmbServerSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of smb server settings resource",
			err.Error(),
		)
		return
	}
	state.Scope = types.StringValue(scope)
	state.ID = types.StringValue("smb_server_settings_" + scope)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed creating smb server settings")
}

// Read reads the resource.
func (r *SmbServerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Started reading smb server settings")

	var state models.SmbServerSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := state.Scope.ValueString()

	settings, err := helper.GetSmbServerSettings(ctx, r.client, scope)
	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading smb server settings",
			message,
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of smb server settings resource",
			err.Error(),
		)
		return
	}
	state.Scope = types.StringValue(scope)
	state.ID = types.StringValue("smb_server_settings_" + scope)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed reading smb server settings")
}

// Update updates the resource.
func (r *SmbServerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Started updating smb server settings")

	var plan models.SmbServerSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := plan.Scope.ValueString()
	if scope == "" {
		scope = "effective"
	}

	var toUpdate powerscale.V6SmbSettingsGlobalSettings
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating smb server settings",
			fmt.Sprintf("Could not read smb server settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateSmbServerSettings(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating smb server settings",
			message,
		)
		return
	}

	settings, err := helper.GetSmbServerSettings(ctx, r.client, scope)
	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading smb server settings",
			message,
		)
		return
	}

	var state models.SmbServerSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of smb server settings resource",
			err.Error(),
		)
		return
	}
	state.Scope = types.StringValue(scope)
	state.ID = types.StringValue("smb_server_settings_" + scope)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed updating smb server settings")
}

// Delete deletes the resource.
func (r *SmbServerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Started deleting smb server settings")

	// Read Terraform prior state data into the model
	var state models.SmbServerSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// SMB Server settings is the native functionality that cannot be deleted, so just remove state
	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, "Completed deleting smb server settings")
}

// ImportState imports the resource.
func (r *SmbServerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Started importing smb server settings")

	reqID := req.ID
	scope := strings.TrimSpace(reqID)

	if scope == "" {
		scope = "effective"
	}

	settings, err := helper.GetSmbServerSettings(ctx, r.client, scope)
	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading smb server settings",
			message,
		)
		return
	}

	var state models.SmbServerSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of smb server settings resource",
			err.Error(),
		)
		return
	}
	state.Scope = types.StringValue(scope)
	state.ID = types.StringValue("smb_server_settings_" + scope)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed importing smb server settings")
}
