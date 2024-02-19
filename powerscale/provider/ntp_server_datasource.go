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
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NtpServerDataSource{}

// NewNtpServerDataSource creates a new data source.
func NewNtpServerDataSource() datasource.DataSource {
	return &NtpServerDataSource{}
}

// NtpServerDataSource defines the data source implementation.
type NtpServerDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NtpServerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ntpserver"
}

// Schema describes the data source arguments.
func (d *NtpServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing NTP Servers from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use NTP Servers to synchronize the system time",
		Description:         "This datasource is used to query the existing NTP Servers from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use NTP Servers to synchronize the system time",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the NTP Server instance.",
				MarkdownDescription: "Unique identifier of the NTP Server instance.",
				Computed:            true,
			},
			"ntp_servers_details": schema.ListNestedAttribute{
				Description:         "List of NTP Servers.",
				MarkdownDescription: "List of NTP Servers.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description:         "Key value from key_file that maps to this server.",
							MarkdownDescription: "Key value from key_file that maps to this server.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "NTP server name.",
							MarkdownDescription: "NTP server name.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Field ID.",
							MarkdownDescription: "Field ID.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Description:         "Filter NTP Servers by names.",
						MarkdownDescription: "Filter NTP Servers by names.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *NtpServerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NtpServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading ntp server data source")

	var state models.NtpServerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ntpServerList, err := helper.GetNtpServers(ctx, d.client)

	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of ntp servers",
			message,
		)
		return
	}

	var ntpServers []models.NtpServerDetailModel
	for _, ntpServerItem := range ntpServerList.Servers {
		val := ntpServerItem
		ntpServer, err := helper.NtpServerDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadNtpServerErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error mapping the list of ntp servers",
				message,
			)
			return
		}
		ntpServers = append(ntpServers, ntpServer)
	}

	state.NtpServers = ntpServers

	// filter ntp servers by names
	if state.NtpServerFilter != nil && len(state.NtpServerFilter.Names) > 0 {
		var validNtpServers []string
		var filteredNtpServers []models.NtpServerDetailModel

		for _, ntpServer := range state.NtpServers {
			for _, name := range state.NtpServerFilter.Names {
				if !name.IsNull() && ntpServer.Name.Equal(name) {
					filteredNtpServers = append(filteredNtpServers, ntpServer)
					validNtpServers = append(validNtpServers, fmt.Sprintf("Name: %s", ntpServer.Name))
					continue
				}
			}
		}

		state.NtpServers = filteredNtpServers

		if len(state.NtpServers) != len(state.NtpServerFilter.Names) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered ntp server names is not a valid powerscale ntp server.",
				fmt.Sprintf("Valid ntp servers: [%v], filtered list: [%v]", strings.Join(validNtpServers, " ; "), state.NtpServerFilter.Names),
			)
		}
	}

	// save into the Terraform state.
	state.ID = types.StringValue("ntp_server_datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading ntp server data source ")
}
