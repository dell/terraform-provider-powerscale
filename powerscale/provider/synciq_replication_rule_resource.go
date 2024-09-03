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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SyncIQRuleResource{}
	_ resource.ResourceWithConfigure   = &SyncIQRuleResource{}
	_ resource.ResourceWithImportState = &SyncIQRuleResource{}
)

// NewSyncIQRuleResource creates a new data source.
func NewSyncIQRuleResource() resource.Resource {
	return &SyncIQRuleResource{}
}

// SyncIQRuleResource defines the resource implementation.
type SyncIQRuleResource struct {
	client *client.Client
}

// ImportState implements resource.ResourceWithImportState.
func (d *SyncIQRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata describes the resource arguments.
func (d *SyncIQRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_rule"
}

// Schema describes the data source arguments.
func (d *SyncIQRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = helper.SyncIQRuleResourceSchema(ctx)
}

// Configure configures the data source.
func (d *SyncIQRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Create allocates the resource.
func (d *SyncIQRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SyncIQRuleResource
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create
	apiReq := helper.GetRequestFromSynciqRuleResource(ctx, plan)
	id, err := helper.CreateSyncIQRule(ctx, d.client, apiReq)
	if err != nil {
		errStr := constants.ListSynciqRulesMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error creating syncIQ rule", message)
		return
	}

	// read state
	state, dgs := d.Get(ctx, id)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read reads data for the resource.
func (d *SyncIQRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform configuration data into the model
	var data models.SyncIQRuleResource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	state, dgs := d.Get(ctx, id)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Get fetches the resource state
func (d *SyncIQRuleResource) Get(ctx context.Context, id string) (models.SyncIQRuleResource, diag.Diagnostics) {
	// Read Terraform configuration data into the model
	var ret models.SyncIQRuleResource
	var dgs diag.Diagnostics
	config, err := helper.GetSyncIQRuleByID(ctx, d.client, id)
	if err != nil {
		errStr := constants.ListSynciqRulesMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		dgs.AddError("Error reading syncIQ rules", message)
		return ret, dgs
	}
	state, diags := helper.NewSyncIQRuleResource(ctx, config.GetRules()[0])
	dgs.Append(diags...)
	return state, dgs
}

// Update updates the resource.
func (d *SyncIQRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state models.SyncIQRuleResource
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update
	apiReq := helper.GetRequestFromSynciqRuleResource(ctx, plan)
	err := helper.UpdateSyncIQRule(ctx, d.client, state.ID.ValueString(), apiReq)
	if err != nil {
		errStr := constants.ListSynciqRulesMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error updating syncIQ rule", message)
		return
	}

	// read state
	state, dgs := d.Get(ctx, state.ID.ValueString())
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete implements resource.Resource.
func (d *SyncIQRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SyncIQRuleResource
	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.DeleteSyncIQRule(ctx, d.client, state.ID.ValueString())
	if err != nil {
		errStr := constants.ListSynciqRulesMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error deleting syncIQ rule", message)
		return
	}

	resp.State.RemoveResource(ctx)
}
