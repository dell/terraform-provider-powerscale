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
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &SnapshotRestoreResource{}
)

// NewSnapshotRestoreResource returns the snapshot restore resource object.
func NewSnapshotRestoreResource() resource.Resource {
	return &SnapshotRestoreResource{}
}

// SnapshotRestoreResource defines the resource implementation.
type SnapshotRestoreResource struct {
	client *client.Client
}

// Configure configures the resource.
func (r *SnapshotRestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

// Metadata describes the resource arguments.
func (r *SnapshotRestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_restore"
}

// Schema defines the schema for the resource.
func (r *SnapshotRestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to restore the data from the snapshot of PowerScale Array. The restore is done using copy/clone/snaprevert job. We can Create, Update and Delete using this resource.",
		Description:         "This resource is used to restore the data from the snapshot of PowerScale Array. The restore is done using copy/clone/snaprevert job. We can Create, Update and Delete using this resource.",
		Attributes:          SnapshotRestoreResourceSchema(),
	}
}

// SnapshotRestoreResourceSchema defines the schema for the resource.
func SnapshotRestoreResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder ID",
			MarkdownDescription: "Placeholder ID",
			Computed:            true,
		},
		"snaprevert_params": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"allow_dup": schema.BoolAttribute{
					Optional:            true,
					Description:         "Whether or not to queue the job if one of the same type is already running or queued.",
					MarkdownDescription: "Whether or not to queue the job if one of the same type is already running or queued.",
				},
				"snapshot_id": schema.Int32Attribute{
					Required:            true,
					Description:         "Snapshot ID.",
					MarkdownDescription: "Snapshot ID.",
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *SnapshotRestoreResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "Creating snapshot restore resource state")
	var plan models.SnapshotRestoreModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageSnapshotRestore(ctx, r.client, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with creating snapshot restore resource state")
}

// Read refreshes the Terraform state with the latest value.
func (r *SnapshotRestoreResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "Reading snapshot restore resource state")
	var state models.SnapshotRestoreModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading snapshot restore resource state")
}

// Update updates the resource and sets the updated Terraform state.
func (r *SnapshotRestoreResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating snapshot restore resource state")
	var plan models.SnapshotRestoreModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var state models.SnapshotRestoreModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageSnapshotRestore(ctx, r.client, plan)
	response.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with updating snapshot restore resource state")
}

// Delete deletes the resource.
func (r *SnapshotRestoreResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting snapshot restore resource state")
	var state models.SnapshotRestoreModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Delete snaprevert domain
	if !state.SnapRevertParams.IsNull() {
		diags := helper.DeleteSnaprevertDomain(ctx, r.client, state)
		response.Diagnostics.Append(diags...)
	}

	if response.Diagnostics.HasError() {
		return
	}

	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with deleting snapshot restore resource state")
}
