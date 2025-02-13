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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NetworkPoolResource{}
var _ resource.ResourceWithConfigure = &NetworkPoolResource{}
var _ resource.ResourceWithImportState = &NetworkPoolResource{}

// NewNetworkPoolResource creates a new resource.
func NewNetworkPoolResource() resource.Resource {
	return &NetworkPoolResource{}
}

// NetworkPoolResource defines the resource implementation.
type NetworkPoolResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *NetworkPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networkpool"
}

// Schema describes the resource arguments.
func (r *NetworkPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the network pool entity of PowerScale Array. We can Create, Update and Delete the network pool using this resource. We can also import an existing network pool from PowerScale array.",
		Description:         "This resource is used to manage the network pool entity of PowerScale Array. We can Create, Update and Delete the network pool using this resource. We can also import an existing network pool from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"access_zone": schema.StringAttribute{
				Description:         "Name of a valid access zone to map IP address pool to the zone.",
				MarkdownDescription: "Name of a valid access zone to map IP address pool to the zone.",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"alloc_method": schema.StringAttribute{
				Description:         "Specifies how IP address allocation is done among pool members.",
				MarkdownDescription: "Specifies how IP address allocation is done among pool members.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description of the pool.",
				MarkdownDescription: "A description of the pool.",
				Optional:            true,
				Computed:            true,
			},
			"groupnet": schema.StringAttribute{
				Description:         "Name of the groupnet this pool belongs to. Cannot be modified once designated",
				MarkdownDescription: "Name of the groupnet this pool belongs to. Cannot be modified once designated",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Unique Pool ID.",
				MarkdownDescription: "Unique Pool ID.",
				Computed:            true,
			},
			"ifaces": schema.ListNestedAttribute{
				Description:         "List of interface members in this pool.",
				MarkdownDescription: "List of interface members in this pool.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"iface": schema.StringAttribute{
							Description:         "A string that defines an interface name.",
							MarkdownDescription: "A string that defines an interface name.",
							Optional:            true,
							Computed:            true,
						},
						"lnn": schema.Int64Attribute{
							Description:         "Logical Node Number (LNN) of a node.",
							MarkdownDescription: "Logical Node Number (LNN) of a node.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.",
				MarkdownDescription: "The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.",
				Required:            true,
			},
			"nfsv3_rroce_only": schema.BoolAttribute{
				Description:         "Indicates that pool contains only RDMA RRoCE capable interfaces.",
				MarkdownDescription: "Indicates that pool contains only RDMA RRoCE capable interfaces.",
				Optional:            true,
				Computed:            true,
			},
			"ranges": schema.ListNestedAttribute{
				Description:         "List of IP address ranges in this pool.",
				MarkdownDescription: "List of IP address ranges in this pool.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"high": schema.StringAttribute{
							Description:         "High IP",
							MarkdownDescription: "High IP",
							Optional:            true,
							Computed:            true,
						},
						"low": schema.StringAttribute{
							Description:         "Low IP",
							MarkdownDescription: "Low IP",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"rebalance_policy": schema.StringAttribute{
				Description:         "Rebalance policy..",
				MarkdownDescription: "Rebalance policy..",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"sc_connect_policy": schema.StringAttribute{
				Description:         "SmartConnect client connection balancing policy.",
				MarkdownDescription: "SmartConnect client connection balancing policy.",
				Optional:            true,
				Computed:            true,
			},
			"sc_dns_zone": schema.StringAttribute{
				Description:         "SmartConnect zone name for the pool.",
				MarkdownDescription: "SmartConnect zone name for the pool.",
				Optional:            true,
				Computed:            true,
			},
			"sc_dns_zone_aliases": schema.ListAttribute{
				Description:         "List of SmartConnect zone aliases (DNS names) to the pool.",
				MarkdownDescription: "List of SmartConnect zone aliases (DNS names) to the pool.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sc_failover_policy": schema.StringAttribute{
				Description:         "SmartConnect IP failover policy.",
				MarkdownDescription: "SmartConnect IP failover policy.",
				Optional:            true,
				Computed:            true,
			},
			"sc_subnet": schema.StringAttribute{
				Description:         "Name of SmartConnect service subnet for this pool.",
				MarkdownDescription: "Name of SmartConnect service subnet for this pool.",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"static_routes": schema.ListNestedAttribute{
				Description:         "List of interface members in this pool.",
				MarkdownDescription: "List of interface members in this pool.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gateway": schema.StringAttribute{
							Description:         "Address of the gateway in the format: yyy.yyy.yyy.yyy",
							MarkdownDescription: "Address of the gateway in the format: yyy.yyy.yyy.yyy",
							Optional:            true,
							Computed:            true,
						},
						"prefixlen": schema.Int64Attribute{
							Description:         "Prefix length in the format: nn.",
							MarkdownDescription: "Prefix length in the format: nn.",
							Optional:            true,
							Computed:            true,
						},
						"subnet": schema.StringAttribute{
							Description:         "Network address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "Network address in the format: xxx.xxx.xxx.xxx",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"subnet": schema.StringAttribute{
				Description:         "The name of the subnet. Cannot be modified once designated",
				MarkdownDescription: "The name of the subnet. Cannot be modified once designated",
				Required:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NetworkPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	powerscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = powerscaleClient
}

// Create allocates the resource.
func (r *NetworkPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating network pool")

	var plan models.NetworkPoolResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planBackup models.NetworkPoolResourceModel
	diagsB := req.Plan.Get(ctx, &planBackup)

	resp.Diagnostics.Append(diagsB...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolToCreate := powerscale.V12SubnetsSubnetPool{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &poolToCreate)
	if err != nil {
		errStr := constants.CreateNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating network pool",
			fmt.Sprintf("Could not read pool param with error: %s", message),
		)
		return
	}
	poolID, err := helper.CreateNetworkPool(ctx, r.client, plan.Groupnet.ValueString(), plan.Subnet.ValueString(), poolToCreate)
	if err != nil {
		errStr := constants.CreateNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating network pool",
			message,
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("network pool %s created", poolID.Id), map[string]interface{}{
		"networkPoolResponse": poolID,
	})

	plan.ID = types.StringValue(poolID.Id)
	getPoolResponse, err := helper.GetNetworkPool(ctx, r.client, plan)
	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating network pool",
			message,
		)
		return
	}

	// update resource state according to response
	if len(getPoolResponse.Pools) <= 0 {
		resp.Diagnostics.AddError(
			"Error creating network pool",
			fmt.Sprintf("Could not read created network pool %s", poolID),
		)
		return
	}

	createdPool := getPoolResponse.Pools[0]
	err = helper.CopyFields(ctx, createdPool, &plan)
	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating network pool",
			fmt.Sprintf("Could not read network pool struct %s with error: %s", poolID, message),
		)
		return
	}
	helper.NetworkPoolListsDiff(ctx, planBackup, &plan)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create network pool completed")
}

// Read reads the resource state.
func (r *NetworkPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading network pool")

	var poolState models.NetworkPoolResourceModel
	diags := req.State.Get(ctx, &poolState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var poolStateBackup models.NetworkPoolResourceModel
	diagsB := req.State.Get(ctx, &poolStateBackup)
	resp.Diagnostics.Append(diagsB...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := poolState.ID
	tflog.Debug(ctx, "calling get network pool by ID", map[string]interface{}{
		"networkPoolID": poolID,
	})
	poolResponse, err := helper.GetNetworkPool(ctx, r.client, poolState)
	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading network pool",
			message,
		)
		return
	}

	if len(poolResponse.Pools) <= 0 {
		resp.Diagnostics.AddError(
			"Error reading network pool",
			fmt.Sprintf("Could not read network pool %s from powerscale with error: network pool not found", poolID),
		)
		return
	}
	tflog.Debug(ctx, "updating read network pool state", map[string]interface{}{
		"networkPoolResponse": poolResponse,
		"networkPoolState":    poolState,
	})
	err = helper.CopyFields(ctx, poolResponse.Pools[0], &poolState)
	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading network pool",
			fmt.Sprintf("Could not read network pool struct %s with error: %s", poolID, message),
		)
		return
	}
	helper.NetworkPoolListsDiff(ctx, poolStateBackup, &poolState)
	diags = resp.State.Set(ctx, poolState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read network pool completed")
}

// Update updates the resource state.
func (r *NetworkPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating network pool")

	var poolPlan models.NetworkPoolResourceModel
	diags := req.Plan.Get(ctx, &poolPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planBackup models.NetworkPoolResourceModel
	diagsB := req.Plan.Get(ctx, &planBackup)

	resp.Diagnostics.Append(diagsB...)
	if resp.Diagnostics.HasError() {
		return
	}

	var poolState models.NetworkPoolResourceModel
	diags = resp.State.Get(ctx, &poolState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update network pool", map[string]interface{}{
		"poolPlan":  poolPlan,
		"poolState": poolState,
	})

	poolID := poolState.ID.ValueString()
	poolPlan.ID = poolState.ID
	var poolToUpdate powerscale.V12GroupnetsGroupnetSubnetsSubnetPool
	// Get param from tf input
	err := helper.ReadFromState(ctx, poolPlan, &poolToUpdate)
	if err != nil {
		errStr := constants.UpdateNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating network pool",
			fmt.Sprintf("Could not read pool param with error: %s", message),
		)
		return
	}
	err = helper.UpdateNetworkPool(ctx, r.client, poolState.Name.ValueString(), poolPlan.Groupnet.ValueString(), poolPlan.Subnet.ValueString(), poolToUpdate)
	if err != nil {
		errStr := constants.UpdateNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating network pool",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get network pool by ID on powerscale client", map[string]interface{}{
		"networkPoolID": poolID,
	})
	updatedPool, err := helper.GetNetworkPool(ctx, r.client, poolPlan)
	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating network pool",
			message,
		)
		return
	}

	if len(updatedPool.Pools) <= 0 {
		resp.Diagnostics.AddError(
			"Error updating network pool",
			fmt.Sprintf("Could not read updated network pool %s", poolID),
		)
		return
	}

	err = helper.CopyFields(ctx, updatedPool.Pools[0], &poolPlan)
	if err != nil {
		errStr := constants.ReadNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating network pool",
			fmt.Sprintf("Could not read network pool struct %s with error: %s", poolID, message),
		)
		return
	}
	helper.NetworkPoolListsDiff(ctx, planBackup, &poolPlan)
	diags = resp.State.Set(ctx, poolPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update network pool completed")
}

// Delete deletes the resource.
func (r *NetworkPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting network pool")

	var poolState models.NetworkPoolResourceModel
	diags := req.State.Get(ctx, &poolState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := poolState.ID.ValueString()
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete network pool on powerscale client", map[string]interface{}{
		"networkPoolID": poolID,
	})
	err := helper.DeleteNetworkPool(ctx, r.client, poolState.Name.ValueString(), poolState.Groupnet.ValueString(), poolState.Subnet.ValueString())
	if err != nil {
		errStr := constants.DeleteNetworkPoolErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting network pool",
			message,
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete network pool completed")
}

// ImportState imports the resource state.
func (r *NetworkPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: groupnet_name.subnetnet_name.pool_name Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("groupnet"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subnet"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
}
