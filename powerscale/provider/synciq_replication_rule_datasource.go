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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SyncIQRuleDataSource{}

// NewSyncIQRuleDataSource creates a new data source.
func NewSyncIQRuleDataSource() datasource.DataSource {
	return &SyncIQRuleDataSource{}
}

// SyncIQRuleDataSource defines the data source implementation.
type SyncIQRuleDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SyncIQRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_rule"
}

// Schema describes the data source arguments.
func (d *SyncIQRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = helper.SyncIQRuleDataSourceSchema(ctx)
}

// Configure configures the data source.
func (d *SyncIQRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *SyncIQRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read Terraform configuration data into the model
	var data models.SyncIQRuleDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *models.SyncIQRuleDataSource
	var errD error
	if id := data.ID.ValueString(); id == "" {
		config, err := helper.GetAllSyncIQRules(ctx, d.client)
		if err != nil {
			errStr := constants.ListSynciqRulesMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading syncIQ rules", message)
			return
		}
		state, errD = helper.NewSyncIQRuleDataSource(ctx, config.GetRules())
	} else {
		config, err := helper.GetSyncIQRuleByID(ctx, d.client, id)
		if err != nil {
			errStr := constants.ListSynciqRulesMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading syncIQ rules", message)
			return
		}
		state, errD = helper.NewSyncIQRuleDataSource(ctx, config.GetRules())
	}

	if errD != nil {
		resp.Diagnostics.AddError("Failed to map sync rule fields", errD.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
