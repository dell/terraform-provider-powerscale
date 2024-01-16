/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
var _ datasource.DataSource = &NetworkRuleDataSource{}

// NewNetworkRuleDataSource creates a new data source.
func NewNetworkRuleDataSource() datasource.DataSource {
	return &NetworkRuleDataSource{}
}

// NetworkRuleDataSource defines the data source implementation.
type NetworkRuleDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NetworkRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_rule"
}

// Schema describes the data source arguments.
func (d *NetworkRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing network rules from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Description:         "This datasource is used to query the existing network rules from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the network rule.",
				MarkdownDescription: "Unique identifier of the network rule.",
				Computed:            true,
			},
			"network_rules": schema.ListNestedAttribute{
				Description:         "List of Network Rules.",
				MarkdownDescription: "List of Network Rules.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description:         "Description for the provisioning rule.",
							MarkdownDescription: "Description for the provisioning rule.",
							Computed:            true,
						},
						"groupnet": schema.StringAttribute{
							Description:         "Name of the groupnet this rule belongs to",
							MarkdownDescription: "Name of the groupnet this rule belongs to",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Unique rule ID.",
							MarkdownDescription: "Unique rule ID.",
							Computed:            true,
						},
						"iface": schema.StringAttribute{
							Description:         "Interface name the provisioning rule applies to.",
							MarkdownDescription: "Interface name the provisioning rule applies to.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Name of the provisioning rule.",
							MarkdownDescription: "Name of the provisioning rule.",
							Computed:            true,
						},
						"node_type": schema.StringAttribute{
							Description:         "Node type the provisioning rule applies to.",
							MarkdownDescription: "Node type the provisioning rule applies to.",
							Computed:            true,
						},
						"pool": schema.StringAttribute{
							Description:         "Name of the pool this rule belongs to.",
							MarkdownDescription: "Name of the pool this rule belongs to.",
							Computed:            true,
						},
						"subnet": schema.StringAttribute{
							Description:         "Name of the subnet this rule belongs to.",
							MarkdownDescription: "Name of the subnet this rule belongs to.",
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
						Description:         "Filter network rules by names.",
						MarkdownDescription: "Filter network rules by names.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"groupnet": schema.StringAttribute{
						Description:         "If specified, only rules for this groupnet will be returned.",
						MarkdownDescription: "If specified, only rules for this groupnet will be returned.",
						Optional:            true,
					},
					"subnet": schema.StringAttribute{
						Description:         "If specified, only rules for this subnet will be returned.",
						MarkdownDescription: "If specified, only rules for this subnet will be returned.",
						Optional:            true,
					},
					"pool": schema.StringAttribute{
						Description:         "If specified, only rules for this pool will be returned.",
						MarkdownDescription: "If specified, only rules for this pool will be returned.",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *NetworkRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NetworkRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading network rule data source")

	var state models.NetworkRuleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	networkRuleList, err := helper.ListNetworkRules(ctx, d.client, state.NetworkRuleFilter)

	if err != nil {
		errStr := constants.ListRuleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of network rules",
			message,
		)
		return
	}

	for _, rule := range networkRuleList {
		entity := models.V3PoolsPoolRulesRule{}
		err := helper.CopyFields(ctx, rule, &entity)
		if err != nil {
			resp.Diagnostics.AddError("Error copying fields of rule datasource",
				fmt.Sprintf("Could not list network rules with error: %s", err.Error()))
			return
		}
		state.NetworkRules = append(state.NetworkRules, entity)
	}

	state.ID = types.StringValue("network_rule_datasource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading network rule data source ")
}
