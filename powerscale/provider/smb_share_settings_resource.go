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
	_ resource.Resource                = &SmbShareSettingsResource{}
	_ resource.ResourceWithConfigure   = &SmbShareSettingsResource{}
	_ resource.ResourceWithImportState = &SmbShareSettingsResource{}
)

// NewSmbShareSettingsResource is a helper function to simplify the provider implementation.
func NewSmbShareSettingsResource() resource.Resource {
	return &SmbShareSettingsResource{}
}

// SmbShareSettingsResource creates a new resource.
type SmbShareSettingsResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *SmbShareSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_share_settings"
}

// Schema describes the resource arguments.
func (r *SmbShareSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the SMB share Settings entity of PowerScale Array. We can Create, Update and Delete the SMB share Settings using this resource. " +
			"We can also import the existing SMB share Settings from PowerScale array. Note that, SMB share Settings is the native functionality of PowerScale. When creating the resource, we actually load SMB share Settings from PowerScale to the resource state.",
		Description: "This resource is used to manage the SMB share Settings entity of PowerScale Array. We can Create, Update and Delete the SMB share Settings using this resource. " +
			"We can also import the existing SMB share Settings from PowerScale array. Note that, SMB share Settings is the native functionality of PowerScale. When creating the resource, we actually load SMB share Settings from PowerScale to the resource state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of SMB share Settings. Value of ID will be same as the access zone.",
				MarkdownDescription: "ID of SMB share Settings. Value of ID will be same as the access zone.",
			},
			"hide_dot_files": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Hide files and directories that begin with a period '.'.",
				MarkdownDescription: "Hide files and directories that begin with a period '.'.",
			},
			"allow_execute_always": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Allows users to execute files they have read rights for.",
				MarkdownDescription: "Allows users to execute files they have read rights for.",
			},
			"host_acl": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
				MarkdownDescription: "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
				ElementType:         types.StringType,
			},
			"directory_create_mask": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unix umask or mode bits.",
				MarkdownDescription: "Unix umask or mode bits.",
			},
			"impersonate_user": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User account to be used as guest account.",
				MarkdownDescription: "User account to be used as guest account.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 511),
				},
			},
			"file_filter_extensions": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies the list of file extensions.",
				MarkdownDescription: "Specifies the list of file extensions.",
				ElementType:         types.StringType,
			},
			"file_create_mode": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unix umask or mode bits.",
				MarkdownDescription: "Unix umask or mode bits.",
			},
			"ntfs_acl_support": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Support NTFS ACLs on files and directories.",
				MarkdownDescription: "Support NTFS ACLs on files and directories.",
			},
			"access_based_enumeration_root_only": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Access-based enumeration on only the root directory of the share.",
				MarkdownDescription: "Access-based enumeration on only the root directory of the share.",
			},
			"directory_create_mode": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unix umask or mode bits.",
				MarkdownDescription: "Unix umask or mode bits.",
			},
			"allow_delete_readonly": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Allow deletion of read-only files in the share.",
				MarkdownDescription: "Allow deletion of read-only files in the share.",
			},
			"ca_write_integrity": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specify the level of write-integrity on continuously available shares. Acceptable values: none, write-read-coherent, full",
				MarkdownDescription: "Specify the level of write-integrity on continuously available shares. Acceptable values: none, write-read-coherent, full",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"none",
						"write-read-coherent",
						"full",
					),
				},
			},
			"strict_flush": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Handle SMB flush operations.",
				MarkdownDescription: "Handle SMB flush operations.",
			},
			"zone": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the access zone in which to update settings",
				MarkdownDescription: "Name of the access zone in which to update settings",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 511),
				},
			},
			"smb3_encryption_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Enables SMB3 encryption for the share.",
				MarkdownDescription: "Enables SMB3 encryption for the share.",
			},
			"mangle_byte_start": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies the wchar_t starting point for automatic byte mangling.",
				MarkdownDescription: "Specifies the wchar_t starting point for automatic byte mangling.",
			},
			"access_based_enumeration": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Only enumerate files and folders the requesting user has access to.",
				MarkdownDescription: "Only enumerate files and folders the requesting user has access to.",
			},
			"file_filtering_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Enables file filtering on the share.",
				MarkdownDescription: "Enables file filtering on the share.",
			},
			"sparse_file": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Enables sparse file.",
				MarkdownDescription: "Enables sparse file.",
			},
			"change_notify": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specify level of change notification alerts on the share. Acceptable values: all, norecurse, none",
				MarkdownDescription: "Specify level of change notification alerts on the share. Acceptable values: all, norecurse, none",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"all",
						"norecurse",
						"none",
					),
				},
			},
			"mangle_map": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Character mangle map.",
				MarkdownDescription: "Character mangle map.",
				ElementType:         types.StringType,
			},
			"file_create_mask": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unix umask or mode bits.",
				MarkdownDescription: "Unix umask or mode bits.",
			},
			"impersonate_guest": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specify the condition in which user access is done as the guest account. Acceptable values: always, bad user, never",
				MarkdownDescription: "Specify the condition in which user access is done as the guest account. Acceptable values: always, bad user, never",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"always",
						"bad user",
						"never",
					),
				},
			},
			"strict_ca_lockout": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies if persistent opens would do strict lockout on the share.",
				MarkdownDescription: "Specifies if persistent opens would do strict lockout on the share.",
			},
			"file_filter_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies if filter list is for deny or allow. Default is deny.",
				MarkdownDescription: "Specifies if filter list is for deny or allow. Default is deny.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"deny",
						"allow",
					),
				},
			},
			"create_permissions": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Set the create permissions for new files and directories in share. Acceptable values: default acl, inherit mode bits, use create mask and mode",
				MarkdownDescription: "Set the create permissions for new files and directories in share. Acceptable values: default acl, inherit mode bits, use create mask and mode",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default acl",
						"inherit mode bits",
						"use create mask and mode",
					),
				},
			},
			"ca_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Persistent open timeout for the share.",
				MarkdownDescription: "Persistent open timeout for the share.",
			},
			"csc_policy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Client-side caching policy for the shares. Acceptable values: manual, documents, programs, none",
				MarkdownDescription: "Client-side caching policy for the shares. Acceptable values: manual, documents, programs, none",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"manual",
						"documents",
						"programs",
						"none",
					),
				},
			},
			"oplocks": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Allow oplock requests.",
				MarkdownDescription: "Allow oplock requests.",
			},
			"strict_locking": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies whether byte range locks contend against SMB I/O.",
				MarkdownDescription: "Specifies whether byte range locks contend against SMB I/O.",
			},
			"continuously_available": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specify if persistent opens are allowed on the share.",
				MarkdownDescription: "Specify if persistent opens are allowed on the share.",
			},
			"scope": schema.StringAttribute{
				Description:         "When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.",
				MarkdownDescription: "When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.",
				Optional:            true,
			},
		},
	}
}

// Configure - defines configuration for smb share resource.
func (r *SmbShareSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pscaleClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = pscaleClient
}

// Create allocates the resource.
func (r *SmbShareSettingsResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "creating smb share settings")

	var sharePlan models.SmbShareSettingsResourceModel
	diags := request.Plan.Get(ctx, &sharePlan)
	//cachedPermission := sharePlan.Permissions

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	zone := sharePlan.Zone.ValueString()
	scope := sharePlan.Scope.ValueString()

	var toUpdate powerscale.V7SmbSettingsShareSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, &sharePlan, &toUpdate)
	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating smb share settings",
			fmt.Sprintf("Could not read smb share settings param with error: %s", message),
		)
		return
	}
	err = helper.UpdateSmbShareSettings(ctx, r.client, toUpdate, zone)
	if err != nil {
		errStr := constants.UpdateSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating smb share settings",
			message)
		return
	}

	getShareSettingsResponse, err := helper.GetSmbShareSettings(ctx, r.client, scope, zone)
	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating smb share settings",
			message)
		return
	}

	var state models.SmbShareSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, getShareSettingsResponse.Settings, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of smb share settings resource", err.Error(),
		)
		return
	}
	state.Zone = sharePlan.Zone
	state.ID = types.StringValue("smb_share_settings_" + sharePlan.Zone.ValueString())
	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create smb share settings completed")
}

// Read reads the resource state.
func (r SmbShareSettingsResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "reading smb share settings")
	var shareState models.SmbShareSettingsResourceModel
	diags := request.State.Get(ctx, &shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	zone := shareState.Zone.ValueString()
	scope := shareState.Scope.ValueString()
	shareResponse, err := helper.GetSmbShareSettings(ctx, r.client, scope, zone)

	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading smb share settings ",
			message)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, shareResponse.Settings, &shareState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of smb share settings resource", err.Error(),
		)
		return
	}

	shareState.Zone = types.StringValue(zone)
	shareState.ID = types.StringValue("smb_share_settings_" + zone)
	diags = response.State.Set(ctx, &shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read smb share settings completed")
}

// Update updates the resource state.
func (r SmbShareSettingsResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating smb share settings")
	var sharePlan models.SmbShareSettingsResourceModel
	diags := request.Plan.Get(ctx, &sharePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	zone := sharePlan.Zone.ValueString()
	scope := sharePlan.Scope.ValueString()
	var shareToUpdate powerscale.V7SmbSettingsShareSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, sharePlan, &shareToUpdate)
	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error updating smb share settings",
			fmt.Sprintf("Could not read smb share settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateSmbShareSettings(ctx, r.client, shareToUpdate, zone)
	if err != nil {
		errStr := constants.UpdateSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating smb share settings ",
			message)
		return
	}

	updatedShare, err := helper.GetSmbShareSettings(ctx, r.client, scope, zone)
	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating smb share settings ",
			message)
		return
	}

	var shareState models.SmbShareSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, updatedShare.Settings, &shareState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of smb share settings resource", err.Error())

		return
	}
	// Zone need to be manually set
	shareState.Zone = sharePlan.Zone
	shareState.ID = types.StringValue("smb_share_settings_" + sharePlan.Zone.ValueString())
	diags = response.State.Set(ctx, shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update smb share settings completed")
}

// Delete deletes the resource.
func (r SmbShareSettingsResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "start deleting smb share settings")
	var shareState models.SmbShareSettingsResourceModel
	diags := request.State.Get(ctx, &shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// SMB share settings is the native functionality that cannot be deleted, so just remove state
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete smb share settings completed")
}

// ImportState imports the resource state.
func (r SmbShareSettingsResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {

	tflog.Info(ctx, "Started importing smb share settings")

	zone := request.ID

	readSmbShareSettings, err := helper.GetSmbShareSettings(ctx, r.client, "", zone)
	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error importing smb share settings",
			message)
		return
	}

	var state models.SmbShareSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, readSmbShareSettings.Settings, &state)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of smb share settings resource", err.Error(),
		)
		return
	}
	state.Zone = types.StringValue(zone)
	state.ID = types.StringValue("smb_share_settings_" + zone)
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import smb share settings completed")
}
