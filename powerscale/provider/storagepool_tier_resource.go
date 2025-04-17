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
	"strconv"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &StoragepoolTierResource{}
	_ resource.ResourceWithImportState = &StoragepoolTierResource{}
)

// NewStoragepoolTierResource creates a new resource.
func NewStoragepoolTierResource() resource.Resource {
	return &StoragepoolTierResource{}
}

// StoragepoolTierResource defines the resource implementation.
type StoragepoolTierResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *StoragepoolTierResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagepool_tier"
}

// Schema describes the resource arguments.
func (r *StoragepoolTierResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the storagepool tier entity of PowerScale Array. We can Create, Update and Delete the storagepool tiers using this resource. We can also import an existing storagepool tier from PowerScale array.",
		Description:         "This resource is used to manage the storagepool tier entity of PowerScale Array. We can Create, Update and Delete the storagepool tiers using this resource. We can also import an existing storagepool tier from PowerScale array.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "Specifies the storagepool tier name.",
				MarkdownDescription: "Specifies the storagepool tier name.",
				Required:            true,
			},
			"children": schema.SetAttribute{
				Description:         "An optional parameter which adds new nodepools to the storagepool tier.",
				MarkdownDescription: "An optional parameter which adds new nodepools to the storagepool tier.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"transfer_limit_pct": schema.Int32Attribute{
				Description:         "Stop moving files to this tier when this limit is met",
				MarkdownDescription: "Stop moving files to this tier when this limit is met",
				Computed:            true,
				Optional:            true,
			},
			"transfer_limit_state": schema.StringAttribute{
				Description:         "How the transfer limit value is being applied",
				MarkdownDescription: "How the transfer limit value is being applied",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.Int64Attribute{
				Description:         "Specifies a string which represents the unique identifier of storagepool tier",
				MarkdownDescription: "Specifies a string which represents the unique identifier of storagepool tier",
				Computed:            true,
			},
			"lnns": schema.SetAttribute{
				Description:         "The nodes that are part of this tier.",
				MarkdownDescription: "The nodes that are part of this tier.",
				Computed:            true,
				ElementType:         types.Int32Type,
			},
		},
	}
}

// Configure configures the resource.
func (r *StoragepoolTierResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	powerscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = powerscaleClient
}

// Create allocates the resource.
func (r *StoragepoolTierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating storagepool tier")
	var plan, state models.StoragepoolTierResourceModel
	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.CreateStoragepoolTier(ctx, r.client, plan, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Info(ctx, "Done with Create Storagepool tier resource state")
}

// Read reads the resource state.
func (r *StoragepoolTierResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading storagepool tier resource")
	var state models.StoragepoolTierResourceModel

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ReadStoragepoolTier(ctx, r.client, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read storagepool tier resource state")
}

// Update updates the resource
func (r *StoragepoolTierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating storagepool settings resource state")
	var plan *models.StoragepoolTierResourceModel
	var state models.StoragepoolTierResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Populate the Edit Parameters
	editValues := powerscale.V16StoragepoolTierExtendedExtended{}
	if state.Name != plan.Name {
		editValues.Name = plan.Name.ValueStringPointer()
	}

	if !plan.Children.IsUnknown() && !state.Children.Equal(plan.Children) {
		var ChildrenList []string
		if len(plan.Children.Elements()) > 0 {
			diags := plan.Children.ElementsAs(ctx, &ChildrenList, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}

			editValues.Children = append(make([]string, 0), ChildrenList...)
		} else {
			editValues.Children = make([]string, 0)
		}
	}

	updatedThroughCondition := false
	if !plan.TransferLimitState.IsNull() && plan.TransferLimitState.ValueString() != "" && state.TransferLimitState != plan.TransferLimitState {
		editValues.TransferLimitState = plan.TransferLimitState.ValueStringPointer()
		if plan.TransferLimitPct.IsNull() || plan.TransferLimitPct.ValueInt32() == 0 {
			editValues.TransferLimitPct = nil
			updatedThroughCondition = true
		}
	}

	if !updatedThroughCondition && state.TransferLimitPct != plan.TransferLimitPct {
		value := plan.TransferLimitPct.ValueInt32()
		editValues.TransferLimitPct = &value
		if plan.TransferLimitState.IsNull() || plan.TransferLimitState.IsUnknown() {
			editValues.TransferLimitState = nil
		}
	}

	diags := helper.UpdateStoragepoolTier(ctx, r.client, editValues, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = helper.ReadStoragepoolTier(ctx, r.client, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update Storagepool Tier resource state")
}

// Delete deletes the resource.
func (r *StoragepoolTierResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting Storagepool tier")
	var data *models.StoragepoolTierResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteParam := r.client.PscaleOpenAPIClient.StoragepoolApi.DeleteStoragepoolv16StoragepoolTier(ctx, strconv.FormatInt(data.Id.ValueInt64(), 10))
	_, err := deleteParam.Execute()
	if err != nil {
		errStr := constants.DeleteStoragepoolTierErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting storage pool",
			message,
		)
	}
	tflog.Info(ctx, "Done with Delete Storagepool Tier resource state")
}

// ImportState imports the resource state.
func (r *StoragepoolTierResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing Storagepool tier resource state")
	var state models.StoragepoolTierResourceModel

	diags := helper.GetStoragepoolTierByID(ctx, r.client, req.ID, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if len(state.Children.Elements()) == 0 {
		state.Children = types.SetNull(types.StringType)
	}
	if len(state.Lnns.Elements()) == 0 {
		state.Lnns = types.SetNull(types.Int32Type)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update Storagepool Tier resource state")

}
