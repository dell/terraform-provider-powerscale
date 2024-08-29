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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &synciqPolicyResource{}
	_ resource.ResourceWithConfigure   = &synciqPolicyResource{}
	_ resource.ResourceWithImportState = &synciqPolicyResource{}
)

// NewSynciqPolicyResource creates a new resource.
func NewSynciqPolicyResource() resource.Resource {
	return &synciqPolicyResource{}
}

// synciqPolicyResource defines the resource implementation.
type synciqPolicyResource struct {
	client *client.Client
}

// Configure implements resource.ResourceWithConfigure.
func (r *synciqPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Schema implements resource.Resource.
func (r *synciqPolicyResource) Schema(ctx context.Context, res resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SynciqpolicyResourceSchema(ctx)
}

// Metadata describes the resource arguments.
func (r *synciqPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_policy"
}

// The function to be called when a resource is created.
func (s *synciqPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan into the model
	var plan models.SynciqpolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V14SyncPolicy
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading create plan",
			err.Error(),
		)
		return
	}

	id, err := helper.CreateSyncIQPolicy(ctx, s.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating syncIQ Policy",
			message,
		)
		return
	}

	state, dgs := s.GetStateById(ctx, id)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster Email resource state")
}

func (s *synciqPolicyResource) GetStateById(ctx context.Context, id string) (models.SynciqpolicyResourceModel, diag.Diagnostics) {
	var dgs diag.Diagnostics
	resp, err := helper.GetSyncIQPolicyByID(ctx, s.client, id)
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		dgs.AddError(
			"Error reading syncIQ Policy",
			message,
		)
		return models.SynciqpolicyResourceModel{}, dgs
	}
	state, diags := helper.NewSynciqpolicyResourceModel(ctx, resp)
	dgs.Append(diags...)
	return state, dgs
}

// The function to be called when a resource is read.
func (s *synciqPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform plan into the model
	var oldState models.SynciqpolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, dgs := s.GetStateById(ctx, oldState.Id.ValueString())
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// The function to be called when a resource is updated.
func (s *synciqPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read Terraform plan into the model
	var plan, OldState models.SynciqpolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &OldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get param from tf input
	var toUpdate powerscale.V14SyncPolicyExtendedExtended
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading update plan",
			err.Error(),
		)
		return
	}

	// if file matching pattern is null then set it to empty
	if plan.FileMatchingPattern.IsNull() {
		toUpdate.FileMatchingPattern = &powerscale.V1SyncJobPolicyFileMatchingPattern{
			OrCriteria: make([]powerscale.V1SyncJobPolicyFileMatchingPatternOrCriteriaItem, 0),
		}
	}

	err = helper.UpdateSyncIQPolicy(ctx, s.client, OldState.Id.ValueString(), toUpdate)
	if err != nil {
		errStr := "Could not update syncIQ Policy with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating syncIQ Policy",
			message,
		)
		return
	}

	state, dgs := s.GetStateById(ctx, OldState.Id.ValueString())
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// The function to be called when a resource is deleted.
func (s *synciqPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read Terraform plan into the model
	var state models.SynciqpolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.DeleteSyncIQPolicy(ctx, s.client, state.Id.ValueString())
	if err != nil {
		errStr := "Could not delete syncIQ Policy with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting syncIQ Policy",
			message,
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState implements resource.ResourceWithImportState.
func (r *synciqPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := helper.GetSyncIQPolicyIDByName(ctx, r.client, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error getting SyncIQ Policy", err.Error())
		return
	}
	req.ID = id
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
