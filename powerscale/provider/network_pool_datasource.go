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
var _ datasource.DataSource = &NetworkPoolDataSource{}

// NewNetworkPoolDataSource creates a new data source.
func NewNetworkPoolDataSource() datasource.DataSource {
	return &NetworkPoolDataSource{}
}

// NetworkPoolDataSource defines the data source implementation.
type NetworkPoolDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NetworkPoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networkpool"
}

// Schema describes the data source arguments.
func (d *NetworkPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing network pools from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can add network interfaces to network pools to associate address ranges with a node or a group of nodes.",
		Description:         "This datasource is used to query the existing network pools from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can add network interfaces to network pools to associate address ranges with a node or a group of nodes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the network pool instance.",
				MarkdownDescription: "Unique identifier of the network pool instance.",
				Computed:            true,
			},
			"network_pools_details": schema.ListNestedAttribute{
				Description:         "List of Network Pools.",
				MarkdownDescription: "List of Network Pools.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_zone": schema.StringAttribute{
							Description:         "Name of a valid access zone to map IP address pool to the zone.",
							MarkdownDescription: "Name of a valid access zone to map IP address pool to the zone.",
							Computed:            true,
						},
						"addr_family": schema.StringAttribute{
							Description:         "IP address format.",
							MarkdownDescription: "IP address format.",
							Computed:            true,
						},
						"aggregation_mode": schema.StringAttribute{
							Description:         "OneFS supports the following NIC aggregation modes.",
							MarkdownDescription: "OneFS supports the following NIC aggregation modes.",
							Computed:            true,
						},
						"alloc_method": schema.StringAttribute{
							Description:         "Specifies how IP address allocation is done among pool members.",
							MarkdownDescription: "Specifies how IP address allocation is done among pool members.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "A description of the pool.",
							MarkdownDescription: "A description of the pool.",
							Computed:            true,
						},
						"groupnet": schema.StringAttribute{
							Description:         "Name of the groupnet this pool belongs to.",
							MarkdownDescription: "Name of the groupnet this pool belongs to.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Unique Pool ID.",
							MarkdownDescription: "Unique Pool ID.",
							Computed:            true,
						},
						"ifaces": schema.ListNestedAttribute{
							Description:         "List of interface members in this pool.",
							MarkdownDescription: "List of interface members in this pool.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"iface": schema.StringAttribute{
										Description:         "A string that defines an interface name.",
										MarkdownDescription: "A string that defines an interface name.",
										Computed:            true,
									},
									"lnn": schema.Int64Attribute{
										Description:         "Logical Node Number (LNN) of a node.",
										MarkdownDescription: "Logical Node Number (LNN) of a node.",
										Computed:            true,
									},
								},
							},
						},
						"name": schema.StringAttribute{
							Description:         "The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.",
							MarkdownDescription: "The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.",
							Computed:            true,
						},
						"nfsv3_rroce_only": schema.BoolAttribute{
							Description:         "Indicates that pool contains only RDMA RRoCE capable interfaces.",
							MarkdownDescription: "Indicates that pool contains only RDMA RRoCE capable interfaces.",
							Computed:            true,
						},
						"ranges": schema.ListNestedAttribute{
							Description:         "List of IP address ranges in this pool.",
							MarkdownDescription: "List of IP address ranges in this pool.",
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
						"rebalance_policy": schema.StringAttribute{
							Description:         "Rebalance policy..",
							MarkdownDescription: "Rebalance policy..",
							Computed:            true,
						},
						"rules": schema.ListAttribute{
							Description:         "Names of the rules in this pool.",
							MarkdownDescription: "Names of the rules in this pool.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"sc_auto_unsuspend_delay": schema.Int64Attribute{
							Description:         "Time delay in seconds before a node which has been automatically unsuspended becomes usable in SmartConnect responses for pool zones.",
							MarkdownDescription: "Time delay in seconds before a node which has been automatically unsuspended becomes usable in SmartConnect responses for pool zones.",
							Computed:            true,
						},
						"sc_connect_policy": schema.StringAttribute{
							Description:         "SmartConnect client connection balancing policy.",
							MarkdownDescription: "SmartConnect client connection balancing policy.",
							Computed:            true,
						},
						"sc_dns_zone": schema.StringAttribute{
							Description:         "SmartConnect zone name for the pool.",
							MarkdownDescription: "SmartConnect zone name for the pool.",
							Computed:            true,
						},
						"sc_dns_zone_aliases": schema.ListAttribute{
							Description:         "List of SmartConnect zone aliases (DNS names) to the pool.",
							MarkdownDescription: "List of SmartConnect zone aliases (DNS names) to the pool.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"sc_failover_policy": schema.StringAttribute{
							Description:         "SmartConnect IP failover policy.",
							MarkdownDescription: "SmartConnect IP failover policy.",
							Computed:            true,
						},
						"sc_subnet": schema.StringAttribute{
							Description:         "Name of SmartConnect service subnet for this pool.",
							MarkdownDescription: "Name of SmartConnect service subnet for this pool.",
							Computed:            true,
						},
						"sc_suspended_nodes": schema.ListAttribute{
							Description:         "List of LNNs showing currently suspended nodes in SmartConnect.",
							MarkdownDescription: "List of LNNs showing currently suspended nodes in SmartConnect.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"sc_ttl": schema.Int64Attribute{
							Description:         "Time to live value for SmartConnect DNS query responses in seconds.",
							MarkdownDescription: "Time to live value for SmartConnect DNS query responses in seconds.",
							Computed:            true,
						},
						"static_routes": schema.ListNestedAttribute{
							Description:         "List of interface members in this pool.",
							MarkdownDescription: "List of interface members in this pool.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"gateway": schema.StringAttribute{
										Description:         "Address of the gateway in the format: yyy.yyy.yyy.yyy",
										MarkdownDescription: "Address of the gateway in the format: yyy.yyy.yyy.yyy",
										Computed:            true,
									},
									"prefixlen": schema.Int64Attribute{
										Description:         "Prefix length in the format: nn.",
										MarkdownDescription: "Prefix length in the format: nn.",
										Computed:            true,
									},
									"subnet": schema.StringAttribute{
										Description:         "Network address in the format: xxx.xxx.xxx.xxx",
										MarkdownDescription: "Network address in the format: xxx.xxx.xxx.xxx",
										Computed:            true,
									},
								},
							},
						},
						"subnet": schema.StringAttribute{
							Description:         "The name of the subnet.",
							MarkdownDescription: "The name of the subnet.",
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
						Description:         "Filter network pools by names.",
						MarkdownDescription: "Filter network pools by names.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"subnet": schema.StringAttribute{
						Description:         "If specified, only pools for this subnet will be returned.",
						MarkdownDescription: "If specified, only pools for this subnet will be returned.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"groupnet": schema.StringAttribute{
						Description:         "If specified, only pools for this groupnet will be returned.",
						MarkdownDescription: "If specified, only pools for this groupnet will be returned.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"access_zone": schema.StringAttribute{
						Description:         "If specified, only pools with this zone name will be returned.",
						MarkdownDescription: "If specified, only pools with this zone name will be returned.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"alloc_method": schema.StringAttribute{
						Description:         "If specified, only pools with this allocation type will be returned.",
						MarkdownDescription: "If specified, only pools with this allocation type will be returned.",
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
func (d *NetworkPoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NetworkPoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading network pool data source")

	var state models.NetworkPoolDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	networkPoolList, err := helper.GetNetworkPools(ctx, d.client, state)

	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of network pools",
			message,
		)
		return
	}

	var networkPools []models.NetworkPoolDetailModel
	for _, poolItem := range networkPoolList.Pools {
		val := poolItem
		networkPool, err := helper.NetworkPoolDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error mapping the list of network pools",
				message,
			)
			return
		}
		networkPools = append(networkPools, networkPool)
	}

	state.NetworkPools = networkPools

	// filter network pools by names
	if state.NetworkPoolFilter != nil && len(state.NetworkPoolFilter.Names) > 0 {
		var validNetworkPools []string
		var filteredNetworkPools []models.NetworkPoolDetailModel

		for _, pool := range state.NetworkPools {
			for _, name := range state.NetworkPoolFilter.Names {
				if !name.IsNull() && pool.Name.Equal(name) {
					filteredNetworkPools = append(filteredNetworkPools, pool)
					validNetworkPools = append(validNetworkPools, fmt.Sprintf("Name: %s", pool.Name))
					continue
				}
			}
		}

		state.NetworkPools = filteredNetworkPools

		if len(state.NetworkPools) != len(state.NetworkPoolFilter.Names) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered network pool names is not a valid powerscale network pool.",
				fmt.Sprintf("Valid network pools: [%v], filtered list: [%v]", strings.Join(validNetworkPools, " ; "), state.NetworkPoolFilter.Names),
			)
		}
	}

	// save into the Terraform state.
	state.ID = types.StringValue("network_pool_datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading network pool data source ")
}
