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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SnapshotResource{}
var _ resource.ResourceWithImportState = &SnapshotResource{}

// NewSnapshotResource creates a new resource.
func NewSnapshotResource() resource.Resource {
	return &SnapshotResource{}
}

// SnapshotResource defines the resource implementation.
type SnapshotResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *SnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

// Schema describes the resource arguments.
func (r *SnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the Snapshot entity of PowerScale Array. We can Create, Update and Delete the Snapshot using this resource. We can also import an existing Snapshot from PowerScale array. PowerScale Snapshots is a logical pointer to data that is stored on a cluster at a specific point in time.",
		Description:         "This resource is used to manage the Snapshot entity of PowerScale Array. We can Create, Update and Delete the Snapshot using this resource. We can also import an existing Snapshot from PowerScale array. PowerScale Snapshots is a logical pointer to data that is stored on a cluster at a specific point in time.",

		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Description:         "The /ifs path snapshotted.",
				MarkdownDescription: "The /ifs path snapshotted.",
				Required:            true,
			},
			"alias": schema.StringAttribute{
				Description:         "The name of the alias, none for real snapshots.",
				MarkdownDescription: "The name of the alias, none for real snapshots.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The user or system supplied snapshot name. This will be null for snapshots pending delete. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed. (Update Supported)",
				MarkdownDescription: "The user or system supplied snapshot name. This will be null for snapshots pending delete. Only alphanumeric characters, underscores ( _ ), and hyphens (-) are allowed. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_\-$@#&]*$`),
						"must contain only alphanumeric characters and _-$@#&",
					),
				},
			},
			"set_expires": schema.StringAttribute{
				Description:         "The amount of time from creation before the snapshot will expire and be eligible for automatic deletion. Resets each time this value is updated (Update Supported)",
				MarkdownDescription: "The amount of time from creation before the snapshot will expire and be eligible for automatic deletion. Resets each time this value is updated (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("Never"),
				Validators: []validator.String{
					stringvalidator.OneOf("Never", "1 Day", "1 Week", "1 Month"),
				},
			},
			"expires": schema.Int64Attribute{
				Description:         "The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.",
				MarkdownDescription: "The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.",
				Computed:            true,
			},
			"created": schema.Int64Attribute{
				Description:         "The Unix Epoch time the snapshot was created.",
				MarkdownDescription: "The Unix Epoch time the snapshot was created.",
				Computed:            true,
			},
			"has_locks": schema.BoolAttribute{
				Description:         "True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of lock.",
				MarkdownDescription: "True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of lock.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "The system ID given to the snapshot. This is useful for tracking the status of delete pending snapshots.",
				MarkdownDescription: "The system ID given to the snapshot. This is useful for tracking the status of delete pending snapshots.",
				Computed:            true,
			},

			"pct_filesystem": schema.NumberAttribute{
				Description:         "Percentage of /ifs used for storing this snapshot.",
				MarkdownDescription: "Percentage of /ifs used for storing this snapshot.",
				Computed:            true,
			},
			"pct_reserve": schema.NumberAttribute{
				Description:         "Percentage of configured snapshot reserved used for storing this snapshot.",
				MarkdownDescription: "Percentage of configured snapshot reserved used for storing this snapshot.",
				Computed:            true,
			},
			"shadow_bytes": schema.Int64Attribute{
				Description:         "The amount of shadow bytes referred to by this snapshot.",
				MarkdownDescription: "The amount of shadow bytes referred to by this snapshot.",
				Computed:            true,
			},
			"schedule": schema.StringAttribute{
				Description:         "The name of the schedule used to create this snapshot, if applicable.",
				MarkdownDescription: "The name of the schedule used to create this snapshot, if applicable.",
				Computed:            true,
			},
			"size": schema.Int64Attribute{
				Description:         "The amount of storage in bytes used to store this snapshot.",
				MarkdownDescription: "The amount of storage in bytes used to store this snapshot.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				Description:         "Snapshot state.",
				MarkdownDescription: "Snapshot state.",
				Computed:            true,
			},
			"target_id": schema.Int64Attribute{
				Description:         "The ID of the snapshot pointed to if this is an alias. 18446744073709551615 (max uint64) is returned for an alias to the live filesystem.",
				MarkdownDescription: "The ID of the snapshot pointed to if this is an alias. 18446744073709551615 (max uint64) is returned for an alias to the live filesystem.",
				Computed:            true,
			},
			"target_name": schema.StringAttribute{
				Description:         "The name of the snapshot pointed to if this is an alias.",
				MarkdownDescription: "The name of the snapshot pointed to if this is an alias.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *SnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating snapshot")
	var plan *models.SnapshotDetailModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {

		return
	}
	result, err := helper.CreateSnapshot(ctx, r.client, plan)
	if err != nil {
		errStr := constants.CreateSnapshotErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			message,
		)
		return
	}

	state, err := helper.SnapshotResourceDetailMapper(ctx, *result)
	if err != nil {
		resp.Diagnostics.AddError(constants.CreateSnapshotErrorMessage+" with error: ", err.Error())
		return
	}
	state.SetExpires = plan.SetExpires
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read reads the resource state.
func (r *SnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading snapshot")
	var plan *models.SnapshotDetailModel

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := helper.GetSpecificSnapshot(ctx, r.client, plan.ID.ValueString())
	if err != nil {
		errStr := constants.ReadSnapshotErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of snapshots",
			message,
		)
		return
	}

	state, err := helper.SnapshotResourceDetailMapper(ctx, result)
	if err != nil {
		resp.Diagnostics.AddError(constants.ReadSnapshotErrorMessage+" with error: ", err.Error())
		return
	}
	state.SetExpires = plan.SetExpires
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource state Path, Name, AuthProviders.
func (r *SnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating snapshot")
	var plan *models.SnapshotDetailModel
	var state *models.SnapshotDetailModel

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

	if plan.Path != state.Path {
		resp.Diagnostics.AddError(
			"Error editing snapshot",
			"Path variable is not able to be updated after creation",
		)
		return
	}

	// Populate the Edit Parameters
	editVal := powerscale.V1SnapshotSnapshotExtendedExtended{}
	if plan.Name != state.Name {
		editVal.Name = plan.Name.ValueStringPointer()
		state.Name = plan.Name
	}
	if plan.Expires != state.Expires {
		expire := helper.CalclulateExpire(plan.SetExpires.ValueString())
		editVal.Expires = &expire
	}
	err := helper.ModifySnapshot(ctx, r.client, state.ID.ValueString(), editVal)
	if err != nil {
		errStr := constants.UpdateSnapshotErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error editing snapshot",
			message,
		)
		return
	}
	state.SetExpires = plan.SetExpires
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete deletes the resource.
func (r *SnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting snapshot")
	var data *models.SnapshotDetailModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deleteParam := r.client.PscaleOpenAPIClient.SnapshotApi.DeleteSnapshotv1SnapshotSnapshot(ctx, data.ID.ValueString())
	_, err := deleteParam.Execute()
	if err != nil {
		errStr := constants.DeleteSnapshotErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting snapshot",
			message,
		)
	}
}

// ImportState imports the resource state.
func (r *SnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
