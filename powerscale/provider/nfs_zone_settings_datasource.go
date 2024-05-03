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

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &NfsZoneSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &NfsZoneSettingsDataSource{}
)

// NewNfsZoneSettingsDataSource is a helper function to simplify the provider implementation.
func NewNfsZoneSettingsDataSource() datasource.DataSource {
	return &NfsZoneSettingsDataSource{}
}

// NfsZoneSettingsDataSource is the data source implementation.
type NfsZoneSettingsDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *NfsZoneSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_zone_settings"
}

// Schema defines the schema for the data source.
func (d *NfsZoneSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This datasource is used to query the NFS Zone Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		MarkdownDescription: "This datasource is used to query the NFS Zone Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of NFS Zone Settings. Value of ID will be same as the access zone.",
				MarkdownDescription: "ID of NFS Zone Settings. Value of ID will be same as the access zone.",
			},
			"nfs_zone_settings": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "NFS Zone Settings",
				MarkdownDescription: "NFS Zone Settings",
				Attributes: map[string]schema.Attribute{
					"nfsv4_no_domain": schema.BoolAttribute{
						Computed:            true,
						Description:         "If true, sends owners and groups without a domain name.",
						MarkdownDescription: "If true, sends owners and groups without a domain name.",
					},
					"nfsv4_no_domain_uids": schema.BoolAttribute{
						Computed:            true,
						Description:         "If true, sends UIDs and GIDs without a domain name.",
						MarkdownDescription: "If true, sends UIDs and GIDs without a domain name.",
					},
					"nfsv4_allow_numeric_ids": schema.BoolAttribute{
						Computed:            true,
						Description:         "If true, sends owners and groups as UIDs and GIDs when look up fails or if the 'nfsv4_no_name' property is set to 1.",
						MarkdownDescription: "If true, sends owners and groups as UIDs and GIDs when look up fails or if the 'nfsv4_no_name' property is set to 1.",
					},
					"nfsv4_no_names": schema.BoolAttribute{
						Computed:            true,
						Description:         "If true, sends owners and groups as UIDs and GIDs.",
						MarkdownDescription: "If true, sends owners and groups as UIDs and GIDs.",
					},
					"nfsv4_replace_domain": schema.BoolAttribute{
						Computed:            true,
						Description:         "If true, replaces the owner or group domain with an NFS domain name.",
						MarkdownDescription: "If true, replaces the owner or group domain with an NFS domain name.",
					},
					"zone": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the access zones in which these settings apply.",
						MarkdownDescription: "Specifies the access zones in which these settings apply.",
					},
					"nfsv4_domain": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the domain or realm through which users and groups are associated.",
						MarkdownDescription: "Specifies the domain or realm through which users and groups are associated.",
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"zone": schema.StringAttribute{
						Optional:            true,
						Description:         "Access zone",
						MarkdownDescription: "Access zone",
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *NfsZoneSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NfsZoneSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started reading nfs zone settings")

	var config models.NfsZoneSettingsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nfsZoneSettings, err := helper.FilterNfsZoneSettings(ctx, d.client, config.NfsZoneSettingsFilter)

	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs zone settings",
			message,
		)
		return
	}

	var settings models.NfsZoneSettings
	err = helper.CopyFields(ctx, nfsZoneSettings.GetSettings(), &settings)
	if err != nil {
		resp.Diagnostics.AddError("Error copying fields of nfs zone settings datasource", err.Error())
		return
	}

	zoneStr := ""
	filter := config.NfsZoneSettingsFilter
	if filter != nil {
		zoneStr = filter.Zone.ValueString()
	}
	if zoneStr == "" {
		zoneStr = "System"
	}

	var state models.NfsZoneSettingsDataSourceModel
	state.ID = types.StringValue(zoneStr)
	state.NfsZoneSettings = &settings
	state.NfsZoneSettingsFilter = config.NfsZoneSettingsFilter

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed reading nfs zone settings")
}
