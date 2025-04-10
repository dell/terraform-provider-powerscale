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
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &SubnetDataSource{}
	_ datasource.DataSourceWithConfigure = &SubnetDataSource{}
)

// NewSubnetDataSource creates a new user data source.
func NewSubnetDataSource() datasource.DataSource {
	return &SubnetDataSource{}
}

// SubnetDataSource defines the data source implementation.
type SubnetDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SubnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

// Schema describes the data source arguments.
func (d *SubnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing Subnets from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. ",
		Description:         "This datasource is used to query the existing Subnets from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. ",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier",
				Description:         "Unique identifier",
			},
			"subnets": schema.ListNestedAttribute{
				Description:         "List of subnets",
				MarkdownDescription: "List of subnets",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"addr_family": schema.StringAttribute{
							Description:         "IP address format.",
							MarkdownDescription: "IP address format.",
							Computed:            true,
						},
						"base_addr": schema.StringAttribute{
							Description:         "The base IP address.",
							MarkdownDescription: "The base IP address.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "A description of the subnet.",
							MarkdownDescription: "A description of the subnet.",
							Computed:            true,
						},
						"dsr_addrs": schema.ListAttribute{
							Description:         "List of Direct Server Return addresses.",
							MarkdownDescription: "List of Direct Server Return addresses.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"gateway": schema.StringAttribute{
							Description:         "Gateway IP address.",
							MarkdownDescription: "Gateway IP address.",
							Computed:            true,
						},
						"gateway_priority": schema.Int64Attribute{
							Description:         "Gateway priority.",
							MarkdownDescription: "Gateway priority.",
							Computed:            true,
						},
						"groupnet": schema.StringAttribute{
							Description:         "Name of the groupnet this subnet belongs to.",
							MarkdownDescription: "Name of the groupnet this subnet belongs to.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Unique Subnet ID.",
							MarkdownDescription: "Unique Subnet ID.",
							Computed:            true,
						},
						"mtu": schema.Int64Attribute{
							Description:         "MTU of the subnet.",
							MarkdownDescription: "MTU of the subnet.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "The name of the subnet.",
							MarkdownDescription: "The name of the subnet.",
							Computed:            true,
						},
						"pools": schema.ListAttribute{
							Description:         "Name of the pools in the subnet.",
							MarkdownDescription: "Name of the pools in the subnet.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"prefixlen": schema.Int64Attribute{
							Description:         "Subnet Prefix Length.",
							MarkdownDescription: "Subnet Prefix Length.",
							Computed:            true,
						},
						"sc_service_addrs": schema.ListNestedAttribute{
							Description:         "List of IP addresses that SmartConnect listens for DNS requests.",
							MarkdownDescription: "List of IP addresses that SmartConnect listens for DNS requests.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"high": schema.StringAttribute{
										Description:         "High IP",
										MarkdownDescription: "High IP",
										Computed:            true,
									},
									"low": schema.StringAttribute{
										Description:         "Low IP",
										MarkdownDescription: "Low IP",
										Computed:            true,
									},
								},
							},
						},
						"sc_service_name": schema.StringAttribute{
							Description:         "Domain Name corresponding to the SmartConnect Service Address.",
							MarkdownDescription: "Domain Name corresponding to the SmartConnect Service Address.",
							Computed:            true,
						},
						"vlan_enabled": schema.BoolAttribute{
							Description:         "VLAN tagging enabled or disabled.",
							MarkdownDescription: "VLAN tagging enabled or disabled.",
							Computed:            true,
						},
						"vlan_id": schema.Int64Attribute{
							Description:         "VLAN ID for all interfaces in the subnet.",
							MarkdownDescription: "VLAN ID for all interfaces in the subnet.",
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
						Description:         "List of subnet name.",
						MarkdownDescription: "List of subnet name.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"groupnet_name": schema.StringAttribute{
						Description:         "Specifies which groupnet to query.",
						MarkdownDescription: "Specifies which groupnet to query.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *SubnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *SubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Subnet data source ")
	var subnetPlan models.SubnetDs
	var subnetState models.SubnetDs
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &subnetPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	subnets, err := helper.ListSubnets(ctx, d.client, subnetPlan.SubnetFilter)
	if err != nil {
		errStr := constants.ListSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading subnets",
			message)
		return
	}

	var validSubnets []string
	for _, subnet := range *subnets {
		entity := models.SubnetDsItem{}
		err := helper.CopyFields(ctx, subnet, &entity)
		if err != nil {
			resp.Diagnostics.AddError("Error copying fields of subnet datasource",
				fmt.Sprintf("Could not list subnets with error: %s", err.Error()))
			return
		}
		validSubnets = append(validSubnets, entity.Name.ValueString())
		subnetState.Subnets = append(subnetState.Subnets, entity)
	}

	if subnetPlan.SubnetFilter != nil && len(subnetState.Subnets) < len(subnetPlan.SubnetFilter.Names) {
		resp.Diagnostics.AddError(
			"Error one or more of the filtered subnet name is not a valid powerscale Subnet.",
			fmt.Sprintf("Valid Subnets: [%v], filtered list: [%v]", strings.Join(validSubnets, " , "), subnetPlan.SubnetFilter.Names),
		)
	}

	subnetState.ID = types.StringValue("Subnet-id")
	subnetState.SubnetFilter = subnetPlan.SubnetFilter
	resp.Diagnostics.Append(resp.State.Set(ctx, &subnetState)...)
}
