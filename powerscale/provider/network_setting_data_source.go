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
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &NetworkSettingDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkSettingDataSource{}
)

// NewNetworkSettingDataSource creates a new network settings data source.
func NewNetworkSettingDataSource() datasource.DataSource {
	return &NetworkSettingDataSource{}
}

// NetworkSettingDataSource defines the data source implementation.
type NetworkSettingDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NetworkSettingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_settings"
}

// Schema describes the data source arguments.
func (d *NetworkSettingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the Network Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Network Settings provide the ability to configure external network configuration on the cluster.",
		Description:         "This datasource is used to query the Network Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Network Settings provide the ability to configure external network configuration on the cluster.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Network Settings ID.",
				MarkdownDescription: "Network Settings ID.",
				Computed:            true,
			},
			"default_groupnet": schema.StringAttribute{
				Description:         "Default client-side DNS settings for non-multitenancy aware programs.",
				MarkdownDescription: "Default client-side DNS settings for non-multitenancy aware programs.",
				Computed:            true,
			},
			"source_based_routing_enabled": schema.BoolAttribute{
				Description:         "Enable or disable Source Based Routing.",
				MarkdownDescription: "Enable or disable Source Based Routing.",
				Computed:            true,
			},
			"sc_rebalance_delay": schema.Int64Attribute{
				Description:         "Delay in seconds for IP rebalance.",
				MarkdownDescription: "Delay in seconds for IP rebalance.",
				Computed:            true,
			},
			"tcp_ports": schema.ListAttribute{
				Description:         "List of client TCP ports.",
				MarkdownDescription: "List of client TCP ports.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *NetworkSettingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NetworkSettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Network Settings data source ")

	var state models.NetworkSettingModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	setting, err := helper.GetNetworkSetting(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting PowerScale Network Settings.", err.Error())
		return
	}

	// parse network settings response to state network settings model
	helper.UpdateNetworkSettingState(ctx, &state, setting)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Network Settings data source ")
}
