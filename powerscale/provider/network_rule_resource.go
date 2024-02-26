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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// NetworkRuleResource creates a new resource.
type NetworkRuleResource struct {
	client *client.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NetworkRuleResource{}
	_ resource.ResourceWithConfigure   = &NetworkRuleResource{}
	_ resource.ResourceWithImportState = &NetworkRuleResource{}
)

// NewNetworkRuleResource is a helper function to simplify the provider implementation.
func NewNetworkRuleResource() resource.Resource {
	return &NetworkRuleResource{}
}

// Metadata describes the resource arguments.
func (r NetworkRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_rule"
}

// Schema describes the resource arguments.
func (r *NetworkRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "V3PoolsPoolRulesRule struct for V3PoolsPoolRulesRule",
		Description:         "V3PoolsPoolRulesRule struct for V3PoolsPoolRulesRule",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description:         "Description for the provisioning rule.",
				MarkdownDescription: "Description for the provisioning rule.",
				Optional:            true,
				Computed:            true,
			},
			"groupnet": schema.StringAttribute{
				Description:         "Name of the groupnet this rule belongs to",
				MarkdownDescription: "Name of the groupnet this rule belongs to",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Unique rule ID.",
				MarkdownDescription: "Unique rule ID.",
				Computed:            true,
			},
			"iface": schema.StringAttribute{
				Description:         "Interface name the provisioning rule applies to.",
				MarkdownDescription: "Interface name the provisioning rule applies to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the provisioning rule.",
				MarkdownDescription: "Name of the provisioning rule.",
				Required:            true,
			},
			"node_type": schema.StringAttribute{
				Description:         "Node type the provisioning rule applies to.",
				MarkdownDescription: "Node type the provisioning rule applies to.",
				Optional:            true,
				Computed:            true,
			},
			"pool": schema.StringAttribute{
				Description:         "Name of the pool this rule belongs to.",
				MarkdownDescription: "Name of the pool this rule belongs to.",
				Required:            true,
			},
			"subnet": schema.StringAttribute{
				Description:         "Name of the subnet this rule belongs to.",
				MarkdownDescription: "Name of the subnet this rule belongs to.",
				Required:            true,
			},
		},
	}
}

// Configure - defines configuration for network rule resource.
func (r *NetworkRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (r NetworkRuleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "creating network rule")

	var rulePlan models.V3PoolsPoolRulesRule
	diags := request.Plan.Get(ctx, &rulePlan)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ruleToCreate := powerscale.V3PoolsPoolRule{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, rulePlan, &ruleToCreate)
	if err != nil {
		response.Diagnostics.AddError("Error creating network rule",
			fmt.Sprintf("Could not read network rule : %s with error: %s", rulePlan.Name.ValueString(), err.Error()),
		)
		return
	}
	networkRuleID, err := helper.CreateNetworkRule(ctx, r.client, rulePlan.Groupnet.ValueString(), rulePlan.Subnet.ValueString(), rulePlan.Pool.ValueString(), ruleToCreate)
	if err != nil {
		errStr := constants.CreateRuleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating rule ",
			message)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("network rule %s created", networkRuleID.Id), map[string]interface{}{
		"networkRuleResponse": networkRuleID,
	})

	rule, err := helper.GetNetworkRule(ctx, r.client, networkRuleID.Id, rulePlan.Groupnet.ValueString(), rulePlan.Subnet.ValueString(), rulePlan.Pool.ValueString())
	if err != nil {
		errStr := constants.GetRuleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error getting rule ",
			message)
		return
	}

	// update resource state according to response
	err = helper.CopyFields(ctx, rule, &rulePlan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of rule resource",
			err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, rulePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create rule completed")
}

// Read reads the resource state.
func (r NetworkRuleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "reading rule")
	var ruleState models.V3PoolsPoolRulesRule
	diags := request.State.Get(ctx, &ruleState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling get rule by name", map[string]interface{}{
		"rule name":     ruleState.Name.ValueString(),
		"groupnet name": ruleState.Groupnet.ValueString(),
	})

	rule, err := helper.GetNetworkRule(ctx, r.client, ruleState.Name.ValueString(), ruleState.Groupnet.ValueString(), ruleState.Subnet.ValueString(), ruleState.Pool.ValueString())
	if err != nil {
		errStr := constants.GetRuleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading rule ",
			message)
		return
	}

	err = helper.CopyFields(ctx, rule, &ruleState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of rule resource",
			err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, ruleState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read rule completed")
}

// Update updates the resource state.
func (r NetworkRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating rule")
	var rulePlan models.V3PoolsPoolRulesRule
	diags := request.Plan.Get(ctx, &rulePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var ruleState models.V3PoolsPoolRulesRule
	diags = response.State.Get(ctx, &ruleState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update network rule", map[string]interface{}{
		"rulePlan":  rulePlan,
		"ruleState": ruleState,
	})

	var ruleToUpdate powerscale.V3GroupnetsGroupnetSubnetsSubnetPoolsPoolRule
	// Get param from tf input
	err := helper.ReadFromState(ctx, rulePlan, &ruleToUpdate)
	if err != nil {
		response.Diagnostics.AddError(
			"Error update rule",
			fmt.Sprintf("Could not read rule struct %s with error: %s", ruleState.Name.ValueString(), err.Error()),
		)
		return
	}
	err = helper.UpdateNetworkRule(ctx, r.client, ruleState.Name.ValueString(), rulePlan.Groupnet.ValueString(), rulePlan.Subnet.ValueString(), rulePlan.Pool.ValueString(), ruleToUpdate)
	if err != nil {
		errStr := constants.UpdateRuleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating rule ",
			message)
		return
	}
	// After update, rule name and groupnet name may change per rule plan
	tflog.Debug(ctx, "calling get rule by name", map[string]interface{}{
		"rule name":     rulePlan.Name,
		"groupnet name": rulePlan.Groupnet,
		"subnet name":   rulePlan.Subnet,
		"pool name":     rulePlan.Pool,
	})

	rule, err := helper.GetNetworkRule(ctx, r.client, rulePlan.Name.ValueString(), rulePlan.Groupnet.ValueString(), rulePlan.Subnet.ValueString(), rulePlan.Pool.ValueString())
	if err != nil {
		errStr := constants.GetSubnetErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading rule ",
			message)
		return
	}

	err = helper.CopyFields(ctx, rule, &ruleState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error copying fields of rule resource",
			err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, ruleState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update rule completed")
}

// Delete deletes the resource.
func (r NetworkRuleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting rule")
	var ruleState models.V3PoolsPoolRulesRule
	diags := request.State.Get(ctx, &ruleState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling delete rule on pscale client", map[string]interface{}{
		"rule name":     ruleState.Name,
		"groupnet name": ruleState.Groupnet,
		"subnet name":   ruleState.Subnet,
		"pool name":     ruleState.Pool,
	})
	err := helper.DeleteNetworkRule(ctx, r.client, ruleState.Name.ValueString(), ruleState.Groupnet.ValueString(), ruleState.Subnet.ValueString(), ruleState.Pool.ValueString())
	if err != nil {
		errStr := constants.DeleteRuleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error deleting rule ",
			message)
		return
	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete rule completed")
}

// ImportState imports the resource state.
func (r NetworkRuleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing network rule")
	idParts := strings.Split(request.ID, ".")

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: groupnet_name.subnetnet_name.pool_name.rule_name Got: %q", request.ID),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("groupnet"), idParts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("subnet"), idParts[1])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("pool"), idParts[2])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("name"), idParts[3])...)
}
