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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SmbShareResource creates a new resource.
type SmbShareResource struct {
	client *client.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &SmbShareResource{}
	_ resource.ResourceWithConfigure   = &SmbShareResource{}
	_ resource.ResourceWithImportState = &SmbShareResource{}
)

// NewSmbShareResource is a helper function to simplify the provider implementation.
func NewSmbShareResource() resource.Resource {
	return &SmbShareResource{}
}

// Metadata describes the resource arguments.
func (r SmbShareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_share"
}

// Schema describes the resource arguments.
func (r *SmbShareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the SMB share entity on PowerScale array. " +
			"PowerScale SMB shares provide clients network access to file system resources on the cluster. " +
			"We can Create, Update and Delete the SMB share using this resource. We can also import an existing SMB Share from PowerScale array.",
		Description: "This resource is used to manage the SMB share entity on PowerScale array. " +
			"PowerScale SMB shares provide clients network access to file system resources on the cluster. " +
			"We can Create, Update and Delete the SMB share using this resource. We can also import an existing SMB Share from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the smb share.",
				MarkdownDescription: "The ID of the smb share.",
				Computed:            true,
			},
			"access_based_enumeration": schema.BoolAttribute{
				Description:         "Only enumerate files and folders the requesting user has access to.",
				MarkdownDescription: "Only enumerate files and folders the requesting user has access to.",
				Optional:            true,
				Computed:            true,
			},
			"access_based_enumeration_root_only": schema.BoolAttribute{
				Description:         "Access-based enumeration on only the root directory of the share.",
				MarkdownDescription: "Access-based enumeration on only the root directory of the share.",
				Optional:            true,
				Computed:            true,
			},
			"allow_delete_readonly": schema.BoolAttribute{
				Description:         "Allow deletion of read-only files in the share.",
				MarkdownDescription: "Allow deletion of read-only files in the share.",
				Optional:            true,
				Computed:            true,
			},
			"allow_execute_always": schema.BoolAttribute{
				Description:         "Allows users to execute files they have read rights for.",
				MarkdownDescription: "Allows users to execute files they have read rights for.",
				Optional:            true,
				Computed:            true,
			},
			"allow_variable_expansion": schema.BoolAttribute{
				Description:         "Allow automatic expansion of variables for home directories.",
				MarkdownDescription: "Allow automatic expansion of variables for home directories.",
				Optional:            true,
				Computed:            true,
			},
			"auto_create_directory": schema.BoolAttribute{
				Description:         "Automatically create home directories.",
				MarkdownDescription: "Automatically create home directories.",
				Optional:            true,
				Computed:            true,
			},
			"browsable": schema.BoolAttribute{
				Description:         "Share is visible in net view and the browse list.",
				MarkdownDescription: "Share is visible in net view and the browse list.",
				Optional:            true,
				Computed:            true,
			},
			"ca_timeout": schema.Int64Attribute{
				Description:         "Persistent open timeout for the share.",
				MarkdownDescription: "Persistent open timeout for the share.",
				Optional:            true,
				Computed:            true,
			},
			"ca_write_integrity": schema.StringAttribute{
				Description:         "Specify the level of write-integrity on continuously available shares.",
				MarkdownDescription: "Specify the level of write-integrity on continuously available shares.",
				Optional:            true,
				Computed:            true,
			},
			"change_notify": schema.StringAttribute{
				Description:         "Level of change notification alerts on the share.",
				MarkdownDescription: "Level of change notification alerts on the share.",
				Optional:            true,
				Computed:            true,
			},
			"continuously_available": schema.BoolAttribute{
				Description:         "Specify if persistent opens are allowed on the share.",
				MarkdownDescription: "Specify if persistent opens are allowed on the share.",
				Computed:            true,
			},
			"create_path": schema.BoolAttribute{
				Description:         "Create path if does not exist.",
				MarkdownDescription: "Create path if does not exist.",
				Optional:            true,
			},
			"create_permissions": schema.StringAttribute{
				Description:         "Create permissions for new files and directories in share.",
				MarkdownDescription: "Create permissions for new files and directories in share.",
				Optional:            true,
				Computed:            true,
			},
			"csc_policy": schema.StringAttribute{
				Description:         "Client-side caching policy for the shares.",
				MarkdownDescription: "Client-side caching policy for the shares.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description for this SMB share.",
				MarkdownDescription: "Description for this SMB share.",
				Optional:            true,
				Computed:            true,
			},
			"directory_create_mask": schema.Int64Attribute{
				Description:         "Directory create mask bits.",
				MarkdownDescription: "Directory create mask bits.",
				Optional:            true,
				Computed:            true,
			},
			"directory_create_mode": schema.Int64Attribute{
				Description:         "Directory create mode bits.",
				MarkdownDescription: "Directory create mode bits.",
				Optional:            true,
				Computed:            true,
			},
			"file_create_mask": schema.Int64Attribute{
				Description:         "File create mask bits.",
				MarkdownDescription: "File create mask bits.",
				Optional:            true,
				Computed:            true,
			},
			"file_create_mode": schema.Int64Attribute{
				Description:         "File create mode bits.",
				MarkdownDescription: "File create mode bits.",
				Optional:            true,
				Computed:            true,
			},
			"file_filter_extensions": schema.ListAttribute{
				Description:         "Specifies the list of file extensions.",
				MarkdownDescription: "Specifies the list of file extensions.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"file_filter_type": schema.StringAttribute{
				Description:         "Specifies if filter list is for deny or allow. Default is deny.",
				MarkdownDescription: "Specifies if filter list is for deny or allow. Default is deny.",
				Optional:            true,
				Computed:            true,
			},
			"file_filtering_enabled": schema.BoolAttribute{
				Description:         "Enables file filtering on this zone.",
				MarkdownDescription: "Enables file filtering on this zone.",
				Optional:            true,
				Computed:            true,
			},
			"hide_dot_files": schema.BoolAttribute{
				Description:         "Hide files and directories that begin with a period '.'.",
				MarkdownDescription: "Hide files and directories that begin with a period '.'.",
				Optional:            true,
				Computed:            true,
			},
			"host_acl": schema.ListAttribute{
				Description:         "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
				MarkdownDescription: "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"impersonate_guest": schema.StringAttribute{
				Description:         "Specify the condition in which user access is done as the guest account.",
				MarkdownDescription: "Specify the condition in which user access is done as the guest account.",
				Optional:            true,
				Computed:            true,
			},
			"impersonate_user": schema.StringAttribute{
				Description:         "User account to be used as guest account.",
				MarkdownDescription: "User account to be used as guest account.",
				Optional:            true,
				Computed:            true,
			},
			"inheritable_path_acl": schema.BoolAttribute{
				Description:         "Set the inheritable ACL on the share path.",
				MarkdownDescription: "Set the inheritable ACL on the share path.",
				Optional:            true,
				Computed:            true,
			},
			"mangle_byte_start": schema.Int64Attribute{
				Description:         "Specifies the wchar_t starting point for automatic byte mangling.",
				MarkdownDescription: "Specifies the wchar_t starting point for automatic byte mangling.",
				Optional:            true,
				Computed:            true,
			},
			"mangle_map": schema.ListAttribute{
				Description:         "Character mangle map.",
				MarkdownDescription: "Character mangle map.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"name": schema.StringAttribute{
				Description:         "Share name.",
				MarkdownDescription: "Share name.",
				Required:            true,
			},
			"ntfs_acl_support": schema.BoolAttribute{
				Description:         "Support NTFS ACLs on files and directories.",
				MarkdownDescription: "Support NTFS ACLs on files and directories.",
				Optional:            true,
				Computed:            true,
			},
			"oplocks": schema.BoolAttribute{
				Description:         "Support oplocks.",
				MarkdownDescription: "Support oplocks.",
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				Description:         "Path of share within /ifs.",
				MarkdownDescription: "Path of share within /ifs.",
				Required:            true,
			},
			"permissions": schema.ListNestedAttribute{
				Description:         "Specifies an ordered list of permission modifications.",
				MarkdownDescription: "Specifies an ordered list of permission modifications.",
				Required:            true,
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission": schema.StringAttribute{
							Description:         "Specifies the file system rights that are allowed or denied.",
							MarkdownDescription: "Specifies the file system rights that are allowed or denied.",
							Required:            true,
						},
						"permission_type": schema.StringAttribute{
							Description:         "Determines whether the permission is allowed or denied.",
							MarkdownDescription: "Determines whether the permission is allowed or denied.",
							Required:            true,
						},
						"trustee": schema.SingleNestedAttribute{
							Description:         "Specifies the persona of the file group.",
							MarkdownDescription: "Specifies the persona of the file group.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "Specifies the serialized form of a persona using security identifier, which can be 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona using security identifier, which can be 'SID:S-1-1'.",
									Optional:            true,
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Optional:            true,
									Computed:            true,
								},
								"type": schema.StringAttribute{
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Optional:            true,
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"run_as_root": schema.ListNestedAttribute{
				Description:         "Allow account to run as root.",
				MarkdownDescription: "Allow account to run as root.",
				Optional:            true,
				Computed:            true,
				Default: listdefault.StaticValue(types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
					"id": types.StringType, "name": types.StringType, "type": types.StringType}})),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Specifies the serialized form of a persona using security identifier, which can be 'SID:S-1-1'.",
							MarkdownDescription: "Specifies the serialized form of a persona using security identifier, which can be 'SID:S-1-1'.",
							Optional:            true,
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Specifies the persona name, which must be combined with a type.",
							MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
							Optional:            true,
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Specifies the type of persona, which must be combined with a name.",
							MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"smb3_encryption_enabled": schema.BoolAttribute{
				Description:         "Enables SMB3 encryption for the share.",
				MarkdownDescription: "Enables SMB3 encryption for the share.",
				Optional:            true,
				Computed:            true,
			},
			"sparse_file": schema.BoolAttribute{
				Description:         "Enables sparse file.",
				MarkdownDescription: "Enables sparse file.",
				Optional:            true,
				Computed:            true,
			},
			"strict_ca_lockout": schema.BoolAttribute{
				Description:         "Specifies if persistent opens would do strict lockout on the share.",
				MarkdownDescription: "Specifies if persistent opens would do strict lockout on the share.",
				Optional:            true,
				Computed:            true,
			},
			"strict_flush": schema.BoolAttribute{
				Description:         "Handle SMB flush operations.",
				MarkdownDescription: "Handle SMB flush operations.",
				Optional:            true,
				Computed:            true,
			},
			"strict_locking": schema.BoolAttribute{
				Description:         "Specifies whether byte range locks contend against SMB I/O.",
				MarkdownDescription: "Specifies whether byte range locks contend against SMB I/O.",
				Optional:            true,
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "Name of the access zone to which to move this SMB share.",
				MarkdownDescription: "Name of the access zone to which to move this SMB share.",
				Optional:            true,
			},
			"zid": schema.Int64Attribute{
				Description:         "Numeric ID of the access zone which contains this SMB share.",
				MarkdownDescription: "Numeric ID of the access zone which contains this SMB share.",
				Computed:            true,
			},
		},
	}
}

// Configure - defines configuration for smb share resource.
func (r *SmbShareResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Create allocates the resource.
func (r SmbShareResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "creating smb share")

	var sharePlan models.SmbShareResource
	diags := request.Plan.Get(ctx, &sharePlan)
	//cachedPermission := sharePlan.Permissions

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	shareToCreate := powerscale.V7SmbShare{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, sharePlan, &shareToCreate)
	if err != nil {
		response.Diagnostics.AddError("Error creating smb share",
			fmt.Sprintf("Could not read smb share param of Path: %s with error: %s", sharePlan.Path.ValueString(), err.Error()),
		)
		return
	}
	shareID, err := helper.CreateSmbShare(ctx, r.client, shareToCreate)
	if err != nil {
		errStr := constants.CreateSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating smb share ",
			message)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("smb share %s created", shareID.Id), map[string]interface{}{
		"smbShareResponse": shareID,
	})

	getShareResponse, err := helper.GetSmbShare(ctx, r.client, shareID.Id, shareToCreate.Zone)
	if err != nil {
		errStr := constants.GetSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating smb share ",
			message)
		return
	}

	// update resource state according to response
	if len(getShareResponse.Shares) <= 0 {
		response.Diagnostics.AddError(
			"Error creating smb share",
			fmt.Sprintf("Could not get created smb share state %s with error: smb share not found", shareID),
		)
		return
	}
	createdShare := getShareResponse.Shares[0]
	err = helper.CopyFieldsToNonNestedModel(ctx, createdShare, &sharePlan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating smb share",
			fmt.Sprintf("Could not read smb share %s with error: %s", shareID, err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, sharePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create smb share completed")
}

// Read reads the resource state.
func (r SmbShareResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "reading smb share")
	var shareState models.SmbShareResource
	diags := request.State.Get(ctx, &shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	shareID := shareState.ID
	tflog.Debug(ctx, "calling get smb share by ID", map[string]interface{}{
		"smbShareID": shareID,
	})
	shareResponse, err := helper.GetSmbShare(ctx, r.client, shareID.ValueString(), shareState.Zone.ValueStringPointer())
	if err != nil {
		errStr := constants.GetSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading smb share ",
			message)
		return
	}

	if len(shareResponse.Shares) <= 0 {
		response.Diagnostics.AddError(
			"Error reading smb share",
			fmt.Sprintf("Could not read smb share %s from pscale with error: smb share not found", shareID),
		)
		return
	}
	tflog.Debug(ctx, "reading read smb share state", map[string]interface{}{
		"smbShareResponse": shareResponse,
		"smbShareState":    shareState,
	})
	err = helper.CopyFieldsToNonNestedModel(ctx, shareResponse.Shares[0], &shareState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading smb share",
			fmt.Sprintf("Could not read smb share struct %s with error: %s", shareID, err.Error()),
		)
		return
	}
	diags = response.State.Set(ctx, shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read smb share completed")
}

// Update updates the resource state.
func (r SmbShareResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating smb share")
	var sharePlan models.SmbShareResource
	diags := request.Plan.Get(ctx, &sharePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var shareState models.SmbShareResource
	diags = response.State.Get(ctx, &shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update smb share", map[string]interface{}{
		"sharePlan":  sharePlan,
		"shareState": shareState,
	})

	shareID := shareState.ID.ValueString()
	var shareToUpdate powerscale.V7SmbShareExtendedExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, sharePlan, &shareToUpdate)
	if err != nil {
		response.Diagnostics.AddError(
			"Error update smb share",
			fmt.Sprintf("Could not read smb share struct %s with error: %s", shareID, err.Error()),
		)
		return
	}
	zoneName := shareState.Zone.ValueString()
	// if share name is updated, query original zone
	if !sharePlan.Zone.Equal(shareState.Zone) {
		zoneName, err = helper.QueryZoneNameByID(ctx, r.client, int64(shareState.Zid.ValueInt64()))
		if err != nil {
			response.Diagnostics.AddError(
				"Error update smb share",
				fmt.Sprintf("Could not read zone %d for share %s with error: %s",
					shareState.Zid.ValueInt64(), shareID, err.Error()),
			)
			return
		}
	}
	err = helper.UpdateSmbShare(ctx, r.client, shareID, &zoneName, shareToUpdate)
	if err != nil {
		errStr := constants.UpdateSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating smb share ",
			message)
		return
	}
	// Share plan must have the field name, update if updated
	shareID = sharePlan.Name.ValueString()
	tflog.Debug(ctx, "calling get smb share by ID on pscale client", map[string]interface{}{
		"smbShareID": shareID,
	})
	updatedShare, err := helper.GetSmbShare(ctx, r.client, shareID, shareToUpdate.Zone)
	if err != nil {
		errStr := constants.GetSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating smb share ",
			message)
		return
	}

	if len(updatedShare.Shares) <= 0 {
		response.Diagnostics.AddError(
			"Error reading smb share",
			fmt.Sprintf("Could not read smb share %s from pscale with error: smb share not found", shareID),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, updatedShare.Shares[0], &shareState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading smb share",
			fmt.Sprintf("Could not read smb share struct %s with error: %s", shareID, err.Error()),
		)
		return
	}
	// Zone need to be manually set
	shareState.Zone = sharePlan.Zone
	diags = response.State.Set(ctx, shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update smb share completed")
}

// Delete deletes the resource.
func (r SmbShareResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting smb share")
	var shareState models.SmbShareResource
	diags := request.State.Get(ctx, &shareState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	shareID := shareState.ID.ValueString()
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete smb share on pscale client", map[string]interface{}{
		"smbShareID": shareID,
	})
	err := helper.DeleteSmbShare(ctx, r.client, shareID, shareState.Zone.ValueStringPointer())
	if err != nil {
		errStr := constants.DeleteSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error deleting smb share ",
			message)
		return
	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete smb share completed")
}

// ImportState imports the resource state.
func (r SmbShareResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var zoneName string
	shareID := request.ID
	// request.ID is form of zoneName:shareID
	if strings.Contains(request.ID, ":") {
		params := strings.Split(request.ID, ":")
		shareID = strings.Trim(params[1], " ")
		zoneName = strings.Trim(params[0], " ")
	}

	readSmbShare, err := helper.GetSmbShare(ctx, r.client, shareID, &zoneName)
	if err != nil {
		errStr := constants.GetSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error importing smb share ",
			message)
		return
	}

	if len(readSmbShare.Shares) <= 0 {
		response.Diagnostics.AddError(
			"Error importing smb share",
			fmt.Sprintf("Could not read smb share %s from pscale with error: smb share not found", shareID),
		)
		return
	}
	var model models.SmbShareResource
	err = helper.CopyFieldsToNonNestedModel(ctx, readSmbShare.Shares[0], &model)
	if err != nil {
		response.Diagnostics.AddError(
			"Error importing smb share",
			fmt.Sprintf("Could not read smb share struct %s with error: %s", shareID, err.Error()),
		)
		return
	}
	model.Zone = types.StringValue(zoneName)
	response.Diagnostics.Append(response.State.Set(ctx, model)...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import smb share completed")
}
