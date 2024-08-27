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
	"encoding/json"
	"fmt"
	"os"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// ImportState implements resource.ResourceWithImportState.
func (r *synciqPolicyResource) ImportState(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	panic("unimplemented")
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
	tflog.Info(ctx, "Creating Cluster Email Settings resource state")
	// Read Terraform plan into the model
	var plan SynciqpolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V14SyncPolicy
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			fmt.Sprintf("Could not read cluster email param with error: %s", message),
		)
		return
	}

	respC, _, err := s.client.PscaleOpenAPIClient.SyncApi.CreateSyncv14SyncPolicy(context.Background()).V14SyncPolicy(toUpdate).Execute()
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}

	// clusterEmail, err := helper.GetClusterEmail(ctx, r.client)
	respR, _, err := s.client.PscaleOpenAPIClient.SyncApi.GetSyncv14SyncPolicy(context.Background(), respC.Id).Execute()
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}
	state, dgs := s.GetState(ctx, respR)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster Email resource state")
}

// The function to be called when a resource is read.
func (s *synciqPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Creating Cluster Email Settings resource state")
	// Read Terraform plan into the model
	var oldState SynciqpolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	respR, _, err := s.client.PscaleOpenAPIClient.SyncApi.GetSyncv14SyncPolicy(context.Background(), oldState.Id.ValueString()).Execute()
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}
	state, dgs := s.GetState(ctx, respR)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster Email resource state")
}

// The function to be called when a resource is updated.
func (s *synciqPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Creating Cluster Email Settings resource state")
	// Read Terraform plan into the model
	var plan, OldState SynciqpolicyModel
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
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			fmt.Sprintf("Could not read cluster email param with error: %s", message),
		)
		return
	}
	reqJson, _ := json.MarshalIndent(toUpdate, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `SyncApi.GetSyncv14SyncPolicy` after Update: %s\n", string(reqJson))

	_, err = s.client.PscaleOpenAPIClient.SyncApi.UpdateSyncv14SyncPolicy(context.Background(), OldState.Id.ValueString()).V14SyncPolicy(toUpdate).Execute()
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}

	// clusterEmail, err := helper.GetClusterEmail(ctx, r.client)
	respR, _, err := s.client.PscaleOpenAPIClient.SyncApi.GetSyncv14SyncPolicy(context.Background(), OldState.Id.ValueString()).Execute()
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}
	state, dgs := s.GetState(ctx, respR)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster Email resource state")
}

func (s *synciqPolicyResource) GetState(ctx context.Context, respR *powerscale.V14SyncPoliciesExtended) (SynciqpolicyModel, diag.Diagnostics) {
	var state SynciqpolicyModel
	var dgs diag.Diagnostics
	source := respR.Policies[0]
	err := helper.CopyFieldsToNonNestedModel(ctx, source, &state)
	if err != nil {
		dgs.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return state, dgs
	}
	return state, nil
}

// The function to be called when a resource is deleted.
func (s *synciqPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Creating Cluster Email Settings resource state")
	// Read Terraform plan into the model
	var state SynciqpolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := s.client.PscaleOpenAPIClient.SyncApi.DeleteSyncv14SyncPolicy(context.Background(), state.Id.ValueString()).Execute()
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting cluster email",
			message,
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
