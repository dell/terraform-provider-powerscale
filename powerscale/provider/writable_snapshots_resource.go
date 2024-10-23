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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// NewWriteableSnapshotResource creates a new resource.
func NewWriteableSnapshotResource() resource.Resource {
	return &WritableSnapshotResource{}
}

// WritableSnapshotResource is the structure for the resource.
type WritableSnapshotResource struct {
	client *client.Client
}

// Configure sets the client for the resource.
func (r *WritableSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = pscaleClient
}

// Metadata sets the type name for the resource.
func (r *WritableSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_writable_snapshot"
}

// Schema returns the schema for the resource.
func (r *WritableSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Description:         "Unique identifier of the writable snapshot.",
				MarkdownDescription: "Unique identifier of the writable snapshot.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"dst_path": schema.StringAttribute{
				Description:         "The destination path for the writable snapshot.",
				MarkdownDescription: "The destination path for the writable snapshot.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"snap_id": schema.StringAttribute{
				Description:         "The ID of the source snapshot for the writable snapshot.",
				MarkdownDescription: "The ID of the source snapshot for the writable snapshot.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(\s*|\d+)$`),
						"must contain only numbers",
					),
				},
			},
			"snap_name": schema.StringAttribute{
				Description:         "The name of the source snapshot for the writable snapshot.",
				MarkdownDescription: "The name of the source snapshot for the writable snapshot.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"src_path": schema.StringAttribute{
				Description:         "The source path of the writable snapshot.",
				MarkdownDescription: "The source path of the writable snapshot.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description:         "The state of the writable snapshot.",
				MarkdownDescription: "The state of the writable snapshot.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Read reads the resource.
func (r *WritableSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read reads the resource.
	tflog.Info(ctx, "Reading Writable Snapshot resource state")

	// Read Terraform prior state data into the model
	var state models.WritableSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// fetch writable snapshot settings
	writableSnapshotResponse, err := helper.GetWritableSnapshot(ctx, r.client, state.DstPath.ValueString())
	if err != nil {
		errStr := constants.ReadWritableSnapshotErrorMsg + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error Reading writable snapshot",
			message,
		)
		return
	}

	// Save updated data into Terraform state
	if len(writableSnapshotResponse.Writable) > 0 {
		helper.UpdateWritableSnapshotState(&state, &writableSnapshotResponse.Writable[0])
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Writable Snapshot resource state")
}

// Create creates the resource.
func (r *WritableSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Writable Snapshot resource state")

	// Read Terraform plan into the model
	var plan models.WritableSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update writable snapshot settings
	toUpdate := powerscale.V14SnapshotWritableItem{
		DstPath: plan.DstPath.ValueString(),
		SrcSnap: plan.SrcSnap.ValueString(),
	}

	// update writable snapshot settings
	writableSnapshotResponse, err := helper.UpdateWritableSnapshot(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateWritableSnapshotErrorMsg + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating writable snapshot",
			message,
		)
		return
	}

	var state models.WritableSnapshot
	helper.UpdateWritableSnapshotState(&state, writableSnapshotResponse)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Writable Snapshot resource state")
}

// Update updates the resource.
func (r *WritableSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// update is not supported, it will perform read-refresh
	resp.Diagnostics.AddError(
		"Error updating Writable Snapshot resource.",
		"An update plan of Writable Snapshot resource should never be invoked. This resource is supposed to be replaced on update.",
	)
}

// Delete deletes the resource.
func (r *WritableSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Writable Snapshot resource state")

	// Read Terraform prior state data into the model
	var state models.WritableSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the writable snapshot
	err := helper.DeleteWritableSnapshot(ctx, r.client, state.DstPath.ValueString())
	if err != nil {
		errStr := constants.DeleteWritableSnapshotErrorMsg + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting writable snapshot",
			message,
		)
		return
	}

	// Remove the resource from the state
	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, "Done with Delete Writable Snapshot resource state")
}

// ImportState imports the resource state.
func (r *WritableSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Writable Snapshot resource state")

	// Read the resource state
	var state models.WritableSnapshot
	// Update the writable snapshot settings
	writableSnapshotResponse, err := helper.GetWritableSnapshot(ctx, r.client, state.DstPath.ValueString())
	if err != nil {
		errStr := constants.ReadWritableSnapshotErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading writable snapshot",
			message,
		)
		return
	}

	if len(writableSnapshotResponse.Writable) > 0 {
		helper.UpdateWritableSnapshotState(&state, &writableSnapshotResponse.Writable[0])
	}

	// Save the updated resource state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import Writable Snapshot resource state")
}
