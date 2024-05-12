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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SmbServerSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &SmbServerSettingsDataSource{}
)

// NewSmbServerSettingsDataSource is a helper function to simplify the provider implementation.
func NewSmbServerSettingsDataSource() datasource.DataSource {
	return &SmbServerSettingsDataSource{}
}

// SmbServerSettingsDataSource is the data source implementation.
type SmbServerSettingsDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *SmbServerSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_server_settings"
}

// Schema defines the schema for the data source.
func (d *SmbServerSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This datasource is used to query the SMB Server Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		MarkdownDescription: "This datasource is used to query the SMB Server Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of SMB Server Settings.",
				MarkdownDescription: "ID of SMB Server Settings.",
			},
			"smb_server_settings": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "SMB Server Settings",
				MarkdownDescription: "SMB Server Settings",
				Attributes: map[string]schema.Attribute{
					"support_multichannel": schema.BoolAttribute{
						Computed:            true,
						Description:         "Support multichannel.",
						MarkdownDescription: "Support multichannel.",
					},
					"enable_security_signatures": schema.BoolAttribute{
						Computed:            true,
						Description:         "Indicates whether the server supports signed SMB packets.",
						MarkdownDescription: "Indicates whether the server supports signed SMB packets.",
					},
					"support_netbios": schema.BoolAttribute{
						Computed:            true,
						Description:         "Support NetBIOS.",
						MarkdownDescription: "Support NetBIOS.",
					},
					"dot_snap_visible_root": schema.BoolAttribute{
						Computed:            true,
						Description:         "Show the .snapshot directory in the root of a share.",
						MarkdownDescription: "Show the .snapshot directory in the root of a share.",
					},
					"access_based_share_enum": schema.BoolAttribute{
						Computed:            true,
						Description:         "Only enumerate files and folders the requesting user has access to.",
						MarkdownDescription: "Only enumerate files and folders the requesting user has access to.",
					},
					"dot_snap_accessible_root": schema.BoolAttribute{
						Computed:            true,
						Description:         "Allow access to the .snapshot directory in the root of the share.",
						MarkdownDescription: "Allow access to the .snapshot directory in the root of the share.",
					},
					"support_smb2": schema.BoolAttribute{
						Computed:            true,
						Description:         "Support the SMB2 protocol on the server.",
						MarkdownDescription: "Support the SMB2 protocol on the server.",
					},
					"audit_logon": schema.StringAttribute{
						Computed:            true,
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
					"dot_snap_accessible_child": schema.BoolAttribute{
						Computed:            true,
						Description:         "Allow access to .snapshot directories in share subdirectories.",
						MarkdownDescription: "Allow access to .snapshot directories in share subdirectories.",
					},
					"srv_cpu_multiplier": schema.Int64Attribute{
						Computed:            true,
						Description:         "Specify the number of SRV service worker threads per CPU.",
						MarkdownDescription: "Specify the number of SRV service worker threads per CPU.",
					},
					"ignore_eas": schema.BoolAttribute{
						Computed:            true,
						Description:         "Specify whether to ignore EAs on files.",
						MarkdownDescription: "Specify whether to ignore EAs on files.",
					},
					"audit_fileshare": schema.StringAttribute{
						Computed:            true,
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
					"onefs_num_workers": schema.Int64Attribute{
						Computed:            true,
						Description:         "Set the maximum number of OneFS driver worker threads.",
						MarkdownDescription: "Set the maximum number of OneFS driver worker threads.",
					},
					"srv_num_workers": schema.Int64Attribute{
						Computed:            true,
						Description:         "Set the maximum number of SRV service worker threads.",
						MarkdownDescription: "Set the maximum number of SRV service worker threads.",
					},
					"dot_snap_visible_child": schema.BoolAttribute{
						Computed:            true,
						Description:         "Show .snapshot directories in share subdirectories.",
						MarkdownDescription: "Show .snapshot directories in share subdirectories.",
					},
					"require_security_signatures": schema.BoolAttribute{
						Computed:            true,
						Description:         "Indicates whether the server requires signed SMB packets.",
						MarkdownDescription: "Indicates whether the server requires signed SMB packets.",
					},
					"server_side_copy": schema.BoolAttribute{
						Computed:            true,
						Description:         "Enable Server Side Copy.",
						MarkdownDescription: "Enable Server Side Copy.",
					},
					"server_string": schema.StringAttribute{
						Computed:            true,
						Description:         "Provides a description of the server.",
						MarkdownDescription: "Provides a description of the server.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 511),
						},
					},
					"service": schema.BoolAttribute{
						Computed:            true,
						Description:         "Specify whether service is enabled.",
						MarkdownDescription: "Specify whether service is enabled.",
					},
					"support_smb3_encryption": schema.BoolAttribute{
						Computed:            true,
						Description:         "Support the SMB3 encryption on the server.",
						MarkdownDescription: "Support the SMB3 encryption on the server.",
					},
					"reject_unencrypted_access": schema.BoolAttribute{
						Computed:            true,
						Description:         "If SMB3 encryption is enabled, reject unencrypted access from clients.",
						MarkdownDescription: "If SMB3 encryption is enabled, reject unencrypted access from clients.",
					},
					"onefs_cpu_multiplier": schema.Int64Attribute{
						Computed:            true,
						Description:         "Specify the number of OneFS driver worker threads per CPU.",
						MarkdownDescription: "Specify the number of OneFS driver worker threads per CPU.",
					},
					"guest_user": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the fully-qualified user to use for guest access.",
						MarkdownDescription: "Specifies the fully-qualified user to use for guest access.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 511),
						},
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *SmbServerSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *SmbServerSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started reading smb server settings")

	var config models.SmbServerSettingsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	smbServerSettings, err := helper.FilterSmbServerSettings(ctx, d.client, config.SmbServerSettingsFilter)

	if err != nil {
		errStr := constants.ReadSmbServerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading smb server settings",
			message,
		)
		return
	}

	var settings models.SmbServerSettings
	err = helper.CopyFields(ctx, smbServerSettings.GetSettings(), &settings)
	if err != nil {
		resp.Diagnostics.AddError("Error copying fields of smb server settings datasource", err.Error())
		return
	}

	scopeStr := ""
	filter := config.SmbServerSettingsFilter
	if filter != nil {
		scopeStr = filter.Scope.ValueString()
	}
	if scopeStr == "" {
		scopeStr = "effective"
	}
	idValue := "smb_server_settings_" + scopeStr

	var state models.SmbServerSettingsDataSourceModel
	state.ID = types.StringValue(idValue)
	state.SmbServerSettings = &settings
	state.SmbServerSettingsFilter = config.SmbServerSettingsFilter

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed reading smb server settings")
}
