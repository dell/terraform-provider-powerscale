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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NtpSettingsDataSource{}

// NewNtpSettingsDataSource creates a new data source.
func NewNtpSettingsDataSource() datasource.DataSource {
	return &NtpSettingsDataSource{}
}

// NtpSettingsDataSource defines the data source implementation.
type NtpSettingsDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NtpSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ntpsettings"
}

// Schema describes the data source arguments.
func (d *NtpSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the NTP Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use NTP Settings to change the settings of NTP Servers",
		Description:         "This datasource is used to query the NTP Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use NTP Settings to change the settings of NTP Servers",
		Attributes: map[string]schema.Attribute{
			"chimers": schema.Int64Attribute{
				Description:         "Number of nodes that will contact the NTP servers.",
				MarkdownDescription: "Number of nodes that will contact the NTP servers.",
				Computed:            true,
			},
			"excluded": schema.ListAttribute{
				Description:         "Node number (LNN) for nodes excluded from chimer duty.",
				MarkdownDescription: "Node number (LNN) for nodes excluded from chimer duty.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"key_file": schema.StringAttribute{
				Description:         "Path to NTP key file within /ifs.",
				MarkdownDescription: "Path to NTP key file within /ifs.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *NtpSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NtpSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading ntp settings data source")

	var state models.NtpSettingsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ntpSettingsResp, err := helper.GetNtpSettings(ctx, d.client)

	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the ntp settings",
			message,
		)
		return
	}

	ntpSettings, err := helper.NtpSettingsDetailMapper(ctx, ntpSettingsResp.Settings)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error mapping the list of ntp settings",
			message,
		)
		return
	}

	state = ntpSettings

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading ntp settings data source ")
}
