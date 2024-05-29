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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &SmbShareSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &SmbShareSettingsDataSource{}
)

// NewSmbShareSettingsDataSource returns the SmbShareSettings data source object.
func NewSmbShareSettingsDataSource() datasource.DataSource {
	return &SmbShareSettingsDataSource{}
}

// SmbShareSettingsDataSource defines the data source implementation.
type SmbShareSettingsDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SmbShareSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_share_settings"
}

// Schema describes the data source arguments.
func (d *SmbShareSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing SMB shares settings from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale SMB shares settings provide clients network access to file system resources on the cluster.",
		Description: "This datasource is used to query the existing SMB shares settings from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale SMB shares settings provide clients network access to file system resources on the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"smb_share_settings": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ca_write_integrity": schema.StringAttribute{
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
					"ntfs_acl_support": schema.BoolAttribute{
						Computed:            true,
						Description:         "Support NTFS ACLs on files and directories.",
						MarkdownDescription: "Support NTFS ACLs on files and directories.",
					},
					"directory_create_mask": schema.Int64Attribute{
						Computed:            true,
						Description:         "Unix umask or mode bits.",
						MarkdownDescription: "Unix umask or mode bits.",
					},
					"mangle_map": schema.ListAttribute{
						Computed:            true,
						Description:         "Character mangle map.",
						MarkdownDescription: "Character mangle map.",
						ElementType:         types.StringType,
					},
					"ca_timeout": schema.Int64Attribute{
						Computed:            true,
						Description:         "Persistent open timeout for the share.",
						MarkdownDescription: "Persistent open timeout for the share.",
					},
					"change_notify": schema.StringAttribute{
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
					"strict_flush": schema.BoolAttribute{
						Computed:            true,
						Description:         "Handle SMB flush operations.",
						MarkdownDescription: "Handle SMB flush operations.",
					},
					"strict_ca_lockout": schema.BoolAttribute{
						Computed:            true,
						Description:         "Specifies if persistent opens would do strict lockout on the share.",
						MarkdownDescription: "Specifies if persistent opens would do strict lockout on the share.",
					},
					"host_acl": schema.ListAttribute{
						Computed:            true,
						Description:         "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
						MarkdownDescription: "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
						ElementType:         types.StringType,
					},
					"allow_delete_readonly": schema.BoolAttribute{
						Computed:            true,
						Description:         "Allow deletion of read-only files in the share.",
						MarkdownDescription: "Allow deletion of read-only files in the share.",
					},
					"create_permissions": schema.StringAttribute{
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
					"zone": schema.StringAttribute{
						Computed:            true,
						Description:         "Name of the access zone in which to update settings",
						MarkdownDescription: "Name of the access zone in which to update settings",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 511),
						},
					},
					"access_based_enumeration": schema.BoolAttribute{
						Computed:            true,
						Description:         "Only enumerate files and folders the requesting user has access to.",
						MarkdownDescription: "Only enumerate files and folders the requesting user has access to.",
					},
					"sparse_file": schema.BoolAttribute{
						Computed:            true,
						Description:         "Enables sparse file.",
						MarkdownDescription: "Enables sparse file.",
					},
					"file_create_mode": schema.Int64Attribute{
						Computed:            true,
						Description:         "Unix umask or mode bits.",
						MarkdownDescription: "Unix umask or mode bits.",
					},
					"file_filter_extensions": schema.ListAttribute{
						Computed:            true,
						Description:         "Specifies the list of file extensions.",
						MarkdownDescription: "Specifies the list of file extensions.",
						ElementType:         types.StringType,
					},
					"access_based_enumeration_root_only": schema.BoolAttribute{
						Computed:            true,
						Description:         "Access-based enumeration on only the root directory of the share.",
						MarkdownDescription: "Access-based enumeration on only the root directory of the share.",
					},
					"file_create_mask": schema.Int64Attribute{
						Computed:            true,
						Description:         "Unix umask or mode bits.",
						MarkdownDescription: "Unix umask or mode bits.",
					},
					"csc_policy": schema.StringAttribute{
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
					"impersonate_guest": schema.StringAttribute{
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
					"continuously_available": schema.BoolAttribute{
						Computed:            true,
						Description:         "Specify if persistent opens are allowed on the share.",
						MarkdownDescription: "Specify if persistent opens are allowed on the share.",
					},
					"strict_locking": schema.BoolAttribute{
						Computed:            true,
						Description:         "Specifies whether byte range locks contend against SMB I/O.",
						MarkdownDescription: "Specifies whether byte range locks contend against SMB I/O.",
					},
					"directory_create_mode": schema.Int64Attribute{
						Computed:            true,
						Description:         "Unix umask or mode bits.",
						MarkdownDescription: "Unix umask or mode bits.",
					},
					"allow_execute_always": schema.BoolAttribute{
						Computed:            true,
						Description:         "Allows users to execute files they have read rights for.",
						MarkdownDescription: "Allows users to execute files they have read rights for.",
					},
					"hide_dot_files": schema.BoolAttribute{
						Computed:            true,
						Description:         "Hide files and directories that begin with a period '.'.",
						MarkdownDescription: "Hide files and directories that begin with a period '.'.",
					},
					"mangle_byte_start": schema.Int64Attribute{
						Computed:            true,
						Description:         "Specifies the wchar_t starting point for automatic byte mangling.",
						MarkdownDescription: "Specifies the wchar_t starting point for automatic byte mangling.",
					},
					"smb3_encryption_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Enables SMB3 encryption for the share.",
						MarkdownDescription: "Enables SMB3 encryption for the share.",
					},
					"file_filter_type": schema.StringAttribute{
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
					"oplocks": schema.BoolAttribute{
						Computed:            true,
						Description:         "Allow oplock requests.",
						MarkdownDescription: "Allow oplock requests.",
					},
					"impersonate_user": schema.StringAttribute{
						Computed:            true,
						Description:         "User account to be used as guest account.",
						MarkdownDescription: "User account to be used as guest account.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(0, 511),
						},
					},
					"file_filtering_enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Enables file filtering on the share.",
						MarkdownDescription: "Enables file filtering on the share.",
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"scope": schema.StringAttribute{
						Optional:            true,
						Description:         "If specified as \"effective\" or not specified, all fields are returned.  If specified as \"user\", only fields with non-default values are shown.  If specified as \"default\", the original values are returned.",
						MarkdownDescription: "If specified as \"effective\" or not specified, all fields are returned.  If specified as \"user\", only fields with non-default values are shown.  If specified as \"default\", the original values are returned.",
					},
					"zone": schema.StringAttribute{
						Optional:            true,
						Description:         "Specifies which access zone to use.",
						MarkdownDescription: "Specifies which access zone to use.",
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (d *SmbShareSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *SmbShareSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Smb Share settings data source ")

	var config models.SmbShareSettingsDatasourceModel
	var state models.SmbShareSettingsDatasourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	smbShareSettings, err := helper.FilterSmbShareSettings(ctx, d.client, config.SmbShareSettingsFilter)

	if err != nil {
		errStr := constants.ReadSMBShareSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading smb share settings",
			message,
		)
		return
	}

	var settings models.SmbShareSettings
	err = helper.CopyFields(ctx, smbShareSettings.GetSettings(), &settings)
	if err != nil {
		resp.Diagnostics.AddError("Error copying fields of smb share settings datasource", err.Error())
		return
	}

	zoneStr := ""
	filter := config.SmbShareSettings
	if filter != nil {
		zoneStr = filter.Zone.ValueString()
	}
	if zoneStr == "" {
		zoneStr = "System"
	}
	idValue := "smb_share_settings_" + zoneStr

	state.SmbShareSettings = &settings
	state.ID = types.StringValue(idValue)

	state.SmbShareSettingsFilter = config.SmbShareSettingsFilter

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed reading smb share settings")
}
