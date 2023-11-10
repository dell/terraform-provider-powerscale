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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// SubnetResource creates a new resource.
type SubnetResource struct {
	client *client.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &SubnetResource{}
	_ resource.ResourceWithConfigure   = &SubnetResource{}
	_ resource.ResourceWithImportState = &SubnetResource{}
)

// NewSubnetResource is a helper function to simplify the provider implementation.
func NewSubnetResource() resource.Resource {
	return &SubnetResource{}
}

// Metadata describes the resource arguments.
func (r SubnetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

// Schema describes the resource arguments.
func (r *SubnetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the Subnet entity on PowerScale array. " +
			"We can Create, Update and Delete the Subnet using this resource. We can also import an existing Subnet from PowerScale array.",
		Description: "This resource is used to manage the Subnet entity on PowerScale array. " +
			"We can Create, Update and Delete the Subnet using this resource. We can also import an existing Subnet from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"addr_family": schema.StringAttribute{
				Description:         "IP address format.",
				MarkdownDescription: "IP address format.",
				Required:            true,
			},
			"base_addr": schema.StringAttribute{
				Description:         "The base IP address.",
				MarkdownDescription: "The base IP address.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description of the subnet.",
				MarkdownDescription: "A description of the subnet.",
				Optional:            true,
				Computed:            true,
			},
			"dsr_addrs": schema.ListAttribute{
				Description:         "List of Direct Server Return addresses.",
				MarkdownDescription: "List of Direct Server Return addresses.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"gateway": schema.StringAttribute{
				Description:         "Gateway IP address.",
				MarkdownDescription: "Gateway IP address.",
				Optional:            true,
				Computed:            true,
			},
			"gateway_priority": schema.Int64Attribute{
				Description:         "Gateway priority.",
				MarkdownDescription: "Gateway priority.",
				Optional:            true,
				Computed:            true,
			},
			"groupnet": schema.StringAttribute{
				Description:         "Name of the groupnet this subnet belongs to. Updating is not allowed.",
				MarkdownDescription: "Name of the groupnet this subnet belongs to. Updating is not allowed.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Unique Subnet ID.",
				MarkdownDescription: "Unique Subnet ID.",
				Computed:            true,
			},
			"mtu": schema.Int64Attribute{
				Description:         "MTU of the subnet.",
				MarkdownDescription: "MTU of the subnet.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the subnet.",
				MarkdownDescription: "The name of the subnet.",
				Required:            true,
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
				Required:            true,
			},
			"sc_service_addrs": schema.ListNestedAttribute{
				Description:         "List of IP addresses that SmartConnect listens for DNS requests.",
				MarkdownDescription: "List of IP addresses that SmartConnect listens for DNS requests.",
				Computed:            true,
				Optional:            true,
				PlanModifiers:       []planmodifier.List{listplanmodifier.UseStateForUnknown()},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"high": schema.StringAttribute{
							Description:         "High IP",
							MarkdownDescription: "High IP",
							Computed:            true,
							Optional:            true,
						},
						"low": schema.StringAttribute{
							Description:         "Low IP",
							MarkdownDescription: "Low IP",
							Computed:            true,
							Optional:            true,
						},
					},
				},
			},
			"sc_service_name": schema.StringAttribute{
				Description:         "Domain Name corresponding to the SmartConnect Service Address.",
				MarkdownDescription: "Domain Name corresponding to the SmartConnect Service Address.",
				Optional:            true,
				Computed:            true,
			},
			"vlan_enabled": schema.BoolAttribute{
				Description:         "VLAN tagging enabled or disabled.",
				MarkdownDescription: "VLAN tagging enabled or disabled.",
				Optional:            true,
				Computed:            true,
			},
			"vlan_id": schema.Int64Attribute{
				Description:         "VLAN ID for all interfaces in the subnet.",
				MarkdownDescription: "VLAN ID for all interfaces in the subnet.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure - defines configuration for subnet resource.
func (r *SubnetResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Create allocates the resource.
func (r SubnetResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "creating subnet")

	var subnetPlan models.V12GroupnetSubnetExtended
	diags := request.Plan.Get(ctx, &subnetPlan)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	subnetToCreate := powerscale.V12GroupnetSubnet{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, subnetPlan, &subnetToCreate)
	if err != nil {
		response.Diagnostics.AddError("Error creating subnet",
			fmt.Sprintf("Could not read subnet : %s with error: %s", subnetPlan.Name.ValueString(), err.Error()),
		)
		return
	}
	subnetID, err := helper.CreateSubnet(ctx, r.client, subnetToCreate, subnetPlan.Groupnet.ValueString())
	if err != nil {
		errStr := constants.CreateSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating subnet ",
			message)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("subnet %s created", subnetID.Id), map[string]interface{}{
		"subnetResponse": subnetID,
	})

	subnet, err := helper.GetSubnet(ctx, r.client, subnetID.Id, subnetPlan.Groupnet.ValueString())
	if err != nil {
		errStr := constants.GetSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating subnet ",
			message)
		return
	}

	// update resource state according to response
	err = helper.CopyFields(ctx, subnet, &subnetPlan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of subnet resource",
			err.Error(),
		)
		return
	}

	if subnetPlan.VlanID.IsUnknown() {
		subnetPlan.VlanID = basetypes.NewInt64Null()
	}

	diags = response.State.Set(ctx, subnetPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create subnet completed")
}

// Read reads the resource state.
func (r SubnetResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "reading subnet")
	var subnetState models.V12GroupnetSubnetExtended
	diags := request.State.Get(ctx, &subnetState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling get subnet by name and groupnet", map[string]interface{}{
		"subnet name":   subnetState.Name,
		"groupnet name": subnetState.Groupnet,
	})

	subnet, err := helper.GetSubnet(ctx, r.client, subnetState.Name.ValueString(), subnetState.Groupnet.ValueString())
	if err != nil {
		errStr := constants.GetSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading subnet ",
			message)
		return
	}

	err = helper.CopyFields(ctx, subnet, &subnetState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of subnet resource",
			err.Error(),
		)
		return
	}

	if subnetState.VlanID.IsUnknown() {
		subnetState.VlanID = basetypes.NewInt64Null()
	}

	diags = response.State.Set(ctx, subnetState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read subnet completed")
}

// Update updates the resource state.
func (r SubnetResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating subnet")
	var subnetPlan models.V12GroupnetSubnetExtended
	diags := request.Plan.Get(ctx, &subnetPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var subnetState models.V12GroupnetSubnetExtended
	diags = response.State.Get(ctx, &subnetState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update subnet", map[string]interface{}{
		"subnetPlan":  subnetPlan,
		"subnetState": subnetState,
	})

	var subnetToUpdate powerscale.V16GroupnetsGroupnetSubnet
	// Get param from tf input
	err := helper.ReadFromState(ctx, subnetPlan, &subnetToUpdate)
	if err != nil {
		response.Diagnostics.AddError(
			"Error update subnet",
			fmt.Sprintf("Could not read subnet struct %s with error: %s", subnetState.Name.ValueString(), err.Error()),
		)
		return
	}
	err = helper.UpdateSubnet(ctx, r.client, subnetState.Name.ValueString(), subnetPlan.Groupnet.ValueString(), subnetToUpdate)
	if err != nil {
		errStr := constants.UpdateSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating subnet ",
			message)
		return
	}
	// After update, subnet name and groupnet name may change per subnet plan
	tflog.Debug(ctx, "calling get subnet by name and groupnet", map[string]interface{}{
		"subnet name":   subnetPlan.Name,
		"groupnet name": subnetPlan.Groupnet,
	})

	subnet, err := helper.GetSubnet(ctx, r.client, subnetPlan.Name.ValueString(), subnetPlan.Groupnet.ValueString())
	if err != nil {
		errStr := constants.GetSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading subnet ",
			message)
		return
	}

	err = helper.CopyFields(ctx, subnet, &subnetState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of subnet resource",
			err.Error(),
		)
		return
	}

	if subnetState.VlanID.IsUnknown() {
		subnetState.VlanID = basetypes.NewInt64Null()
	}

	diags = response.State.Set(ctx, subnetState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update subnet completed")
}

// Delete deletes the resource.
func (r SubnetResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting subnet")
	var subnetState models.V12GroupnetSubnetExtended
	diags := request.State.Get(ctx, &subnetState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling delete subnet on pscale client", map[string]interface{}{
		"subnet name":   subnetState.Name,
		"groupnet name": subnetState.Groupnet,
	})
	err := helper.DeleteSubnet(ctx, r.client, subnetState.Name.ValueString(), subnetState.Groupnet.ValueString())
	if err != nil {
		errStr := constants.DeleteSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error deleting subnet ",
			message)
		return
	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete subnet completed")
}

// ImportState imports the resource state.
func (r SubnetResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing subnet")
	idParts := strings.Split(request.ID, ".")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: groupnet_name.subnet_name. Got: %q", request.ID),
		)
		return
	}
	var subnetState models.V12GroupnetSubnetExtended
	subnet, err := helper.GetSubnet(ctx, r.client, idParts[1], idParts[0])
	if err != nil {
		errStr := constants.GetSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error importing subnet ",
			message)
		return
	}
	err = helper.CopyFields(ctx, subnet, &subnetState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of subnet resource",
			err.Error(),
		)
		return
	}

	if subnetState.VlanID.IsUnknown() {
		subnetState.VlanID = basetypes.NewInt64Null()
	}

	response.Diagnostics.Append(response.State.Set(ctx, &subnetState)...)
	if response.Diagnostics.HasError() {
		return
	}
}
