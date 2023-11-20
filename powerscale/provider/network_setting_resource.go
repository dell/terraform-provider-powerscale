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
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NetworkSettingResource{}
	_ resource.ResourceWithConfigure   = &NetworkSettingResource{}
	_ resource.ResourceWithImportState = &NetworkSettingResource{}
)

// NewNetworkSettingResource creates a new resource.
func NewNetworkSettingResource() resource.Resource {
	return &NetworkSettingResource{}
}

// NetworkSettingResource defines the resource implementation.
type NetworkSettingResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *NetworkSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_settings"
}

// Schema describes the resource arguments.
func (r *NetworkSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the Network Settings entity of PowerScale Array. " +
			"PowerScale Network Settings provide the ability to configure external network configuration on the cluster." +
			"We can Create, Update and Delete the Network Settings using this resource. We can also import an existing Network Settings from PowerScale array. " +
			"Note that, Network Settings is the native functionality of PowerScale. When creating the resource, we actually load Network Settings from PowerScale to the resource state. ",
		Description: "This resource is used to manage the Network Settings entity of PowerScale Array. " +
			"PowerScale Network Settings provide the ability to configure external network configuration on the cluster." +
			"We can Create, Update and Delete the Network Settings using this resource. We can also import an existing Network Settings from PowerScale array. " +
			"Note that, Network Settings is the native functionality of PowerScale. When creating the resource, we actually load Network Settings from PowerScale to the resource state. ",

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
				Description:         "Enable or disable Source Based Routing. (Update Supported)",
				MarkdownDescription: "Enable or disable Source Based Routing. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"sc_rebalance_delay": schema.Int64Attribute{
				Description:         "Delay in seconds for IP rebalance. (Update Supported)",
				MarkdownDescription: "Delay in seconds for IP rebalance. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.Int64{int64validator.Between(0, 10)},
			},
			"tcp_ports": schema.ListAttribute{
				Description:         "List of client TCP ports. (Update Supported)",
				MarkdownDescription: "List of client TCP ports. (Update Supported)",
				ElementType:         types.Int64Type,
				Optional:            true,
				Computed:            true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueInt64sAre(int64validator.Between(0, 65535)),
					listvalidator.SizeBetween(0, 65535),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *NetworkSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *NetworkSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Network Settings resource state")
	// Read Terraform plan into the model
	var plan models.NetworkSettingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.NetworkSettingModel
	settingResponse, err := helper.GetNetworkSetting(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the network settings", err.Error())
		return
	}

	// parse network settings response to state model
	helper.UpdateNetworkSettingState(ctx, &state, settingResponse)

	if diags := helper.UpdateNetworkSetting(ctx, r.client, &state, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	settingResponse, err = helper.GetNetworkSetting(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the network settings", err.Error())
		return
	}

	// parse network settings response to state model
	helper.UpdateNetworkSettingState(ctx, &plan, settingResponse)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create Network Settings resource state")
}

// Read reads the resource state.
func (r *NetworkSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Network Settings resource")
	var state models.NetworkSettingModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settingResponse, err := helper.GetNetworkSetting(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the network settings", err.Error())
		return
	}

	// parse network settings response to state model
	helper.UpdateNetworkSettingState(ctx, &state, settingResponse)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Network Settings resource")
}

// Update updates the resource state.
func (r *NetworkSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Network Settings resource...")
	// Read Terraform plan into the model
	var plan models.NetworkSettingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.NetworkSettingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diags := helper.UpdateNetworkSetting(ctx, r.client, &state, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	settingResponse, err := helper.GetNetworkSetting(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the network settings", err.Error())
		return
	}

	// parse network settings response to state model
	helper.UpdateNetworkSettingState(ctx, &plan, settingResponse)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Update Network Settings resource")

}

// Delete deletes the resource.
func (r *NetworkSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Network Settings resource state")
	var state models.NetworkSettingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Network Settings resource state")
}

// ImportState imports the resource state.
func (r *NetworkSettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Network Settings resource")
	var state models.NetworkSettingModel

	settingResponse, err := helper.GetNetworkSetting(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the network settings", err.Error())
		return
	}

	// parse network settings response to state model
	helper.UpdateNetworkSettingState(ctx, &state, settingResponse)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Network Settings resource")
}
