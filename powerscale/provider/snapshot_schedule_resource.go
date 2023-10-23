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
	"regexp"
	"strconv"
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

// SnapshotScheduleResource creates a new resource.
type SnapshotScheduleResource struct {
	client *client.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &SnapshotScheduleResource{}
	_ resource.ResourceWithConfigure   = &SnapshotScheduleResource{}
	_ resource.ResourceWithImportState = &SnapshotScheduleResource{}
)

// NewSnapshotScheduleResource is a helper function to simplify the provider implementation.
func NewSnapshotScheduleResource() resource.Resource {
	return &SnapshotScheduleResource{}
}

// Metadata describes the resource arguments.
func (r SnapshotScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_schedule"
}

// Schema describes the resource arguments.
func (r *SnapshotScheduleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the Snapshot Schedule entity on PowerScale array. " +
			"We can Create, Update and Delete the Snapshot Schedules using this resource. We can also import an existing Snapshot Schedule from PowerScale array.",
		Description: "This resource is used to manage the Snapshot Schedule entity on PowerScale array. " +
			"We can Create, Update and Delete the Snapshot Schedules using this resource. We can also import an existing Snapshot Schedule from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"alias": schema.StringAttribute{
				Description:         "Alias name to create for each snapshot.",
				MarkdownDescription: "Alias name to create for each snapshot.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(255),
				},
			},
			"duration": schema.Int64Attribute{
				Description:         "Time in seconds added to creation time to construction expiration time.",
				MarkdownDescription: "Time in seconds added to creation time to construction expiration time.",
				Computed:            true,
				Optional:            true,
			},
			"retention_time": schema.StringAttribute{
				Description:         "Time value in String for which snapshots created by this snapshot schedule should be retained.Values supported are of format : " + "Never Expires, x Seconds(s), x Minute(s), x Hour(s), x Week(s), x Day(s), x Month(s), x Year(s) where x can be any integer value",
				MarkdownDescription: "Time value in String for which snapshots created by this snapshot schedule should be retained.Values supported are of format : " + "Never Expires, x Seconds(s), x Minute(s), x Hour(s), x Week(s), x Day(s), x Month(s), x Year(s) where x can be any integer value",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("1 Week(s)"),
				Validators: []validator.String{stringvalidator.RegexMatches(
					regexp.MustCompile(`^(\d+)\s+(Second\(s\)|Minute\(s\)|Hour\(s\)|Day\(s\)|Week\(s\)|Year\(s\)|Never Expires)$`), "must be in proper format'",
				)},
			},
			"id": schema.StringAttribute{
				Description:         "The system ID given to the schedule.",
				MarkdownDescription: "The system ID given to the schedule.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "The schedule name.",
				MarkdownDescription: "The schedule name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(255),
				},
			},
			"next_run": schema.Int64Attribute{
				Description:         "Unix Epoch time of next snapshot to be created.",
				MarkdownDescription: "Unix Epoch time of next snapshot to be created.",
				Computed:            true,
				Optional:            true,
			},
			"next_snapshot": schema.StringAttribute{
				Description:         "Formatted name (see pattern) of next snapshot to be created.",
				MarkdownDescription: "Formatted name (see pattern) of next snapshot to be created",
				Computed:            true,
				Optional:            true,
			},
			"path": schema.StringAttribute{
				Description:         "The /ifs path snapshotted.",
				MarkdownDescription: "The /ifs path snapshotted.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("/ifs"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"pattern": schema.StringAttribute{
				Description:         "Pattern expanded with strftime to create snapshot names.Some sample values for pattern are: 'snap-%F' would yield snap-1984-03-20 , 'backup-%FT%T' would yield backup-1984-03-20T22:30:00",
				MarkdownDescription: "Pattern expanded with strftime to create snapshot names.Some sample values for pattern are: 'snap-%F' would yield snap-1984-03-20 , 'backup-%FT%T' would yield backup-1984-03-20T22:30:00",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ScheduleName_duration_%Y-%m-%d_%H:%M"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"schedule": schema.StringAttribute{
				Description: `The isidate-compatible natural language description of the schedule. It specifies the frequency of the schedule. You can specify this as combination of <interval> and <frequency> where each of them can be defined as: 
				<interval>:
					*Every [ ( other | <integer> ) ] ( weekday | day | week [ on <day>] | month [ on the <integer> ] | <day>[, ...] [ of every [ ( other | <integer> ) ] week ] | The last (day | weekday | <day>) of every [ (other | <integer>) ] month | The <integer> (weekday | <day>) of every [ (other | <integer>) ] month | The <integer> of every [ (other | <integer>) ] month | Yearly on <month> <integer> | Yearly on the (last | <integer>) [ weekday | <day> ] of <month>
				<frequency>:
					*at <hh>[:<mm>] [ (AM | PM) ] | every [ <integer> ] (hours | minutes) [ between <hh>[:<mm>] [ (AM | PM) ] and <hh>[:<mm>] [ (AM | PM) ] | every [ <integer> ] (hours | minutes) [ from <hh>[:<mm>] [ (AM | PM) ] to <hh>[:<mm>] [ (AM | PM) ]
				Additionally:
					<integer> can include "st," "th," or "rd," e.g., "Every 1st month."
					<day> can be a day of the week or a three-letter abbreviation, e.g., "saturday" or "sat."
					<month> must be the name of the month or its abbreviation, e.g., "July" or "Jul."
				Some sample values:  "Every 2 days.", "Every 3rd weekday at 11 PM.", "Every month on the 15th at 1:30 AM."`,
				MarkdownDescription: `The isidate-compatible natural language description of the schedule. It specifies the frequency of the schedule. You can specify this as combination of <interval> and <frequency> where each of them can be defined as: 
				<interval>:
					*Every [ ( other | <integer> ) ] ( weekday | day | week [ on <day>] | month [ on the <integer> ] | <day>[, ...] [ of every [ ( other | <integer> ) ] week ] | The last (day | weekday | <day>) of every [ (other | <integer>) ] month | The <integer> (weekday | <day>) of every [ (other | <integer>) ] month | The <integer> of every [ (other | <integer>) ] month | Yearly on <month> <integer> | Yearly on the (last | <integer>) [ weekday | <day> ] of <month>
				<frequency>:
					*at <hh>[:<mm>] [ (AM | PM) ] | every [ <integer> ] (hours | minutes) [ between <hh>[:<mm>] [ (AM | PM) ] and <hh>[:<mm>] [ (AM | PM) ] | every [ <integer> ] (hours | minutes) [ from <hh>[:<mm>] [ (AM | PM) ] to <hh>[:<mm>] [ (AM | PM) ]
				Additionally:
					<integer> can include "st," "th," or "rd," e.g., "Every 1st month."
					<day> can be a day of the week or a three-letter abbreviation, e.g., "saturday" or "sat."
					<month> must be the name of the month or its abbreviation, e.g., "July" or "Jul."
				Some sample values:  "Every 2 days.", "Every 3rd weekday at 11 PM.", "Every month on the 15th at 1:30 AM."`,
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("every 1 days at 12:00 AM"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// Configure - defines configuration for snapshot schedule resource.
func (r *SnapshotScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (r SnapshotScheduleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Snapshot Schedule")

	var plan *models.SnapshotScheduleResource
	diags := request.Plan.Get(ctx, &plan)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	createResp, err := helper.CreateSnapshotSchedule(ctx, r.client, plan)
	if err != nil {
		errStr := constants.CreateSnapshotScheduleErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error creating snapshot schedule",
			message,
		)
		return
	}
	result, err := helper.GetSpecificSnapshotSchedule(ctx, r.client, strconv.FormatInt(int64(createResp.Id), 10))
	if err != nil {
		errStr := constants.CreateSnapshotScheduleErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error creating snapshot schedule",
			message,
		)
		return
	}
	state := models.SnapshotScheduleResource{}
	err = helper.SnapshotScheduleMapper(ctx, result, &state)
	if err != nil {
		response.Diagnostics.AddError(constants.CreateSnapshotScheduleErrorMessage+" with error: ", err.Error())
		return
	}
	state.RetentionTime = plan.RetentionTime
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)

	tflog.Info(ctx, "Create Snapshot Schedule completed")
}

// Read reads the resource state.
func (r SnapshotScheduleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Snapshot Schedule")
	var plan *models.SnapshotScheduleResource

	// Read Terraform prior state plan into the model
	response.Diagnostics.Append(request.State.Get(ctx, &plan)...)

	if response.Diagnostics.HasError() {
		return
	}

	result, err := helper.GetSpecificSnapshotSchedule(ctx, r.client, plan.ID.ValueString())
	if err != nil {
		errStr := constants.ReadSnapshotScheduleErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error getting the snapshot schedule",
			message,
		)
		return
	}
	state := models.SnapshotScheduleResource{}
	err = helper.SnapshotScheduleMapper(ctx, result, &state)
	if err != nil {
		response.Diagnostics.AddError(constants.ReadSnapshotScheduleErrorMessage+" with error: ", err.Error())
		return
	}
	state.RetentionTime = plan.RetentionTime
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Read Snapshot Schedule  completed")
}

// Update updates the resource state.
func (r SnapshotScheduleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Snapshot Schedule ")
	var plan models.SnapshotScheduleResource
	var state models.SnapshotScheduleResource

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Read the state
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := helper.UpdateSnapshotSchedule(ctx, r.client, &plan, &state)
	if err != nil {
		errStr := constants.UpdateSnapshotScheduleErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error Updating the snapshot schedule",
			message,
		)
		return
	}

	result, err := helper.GetSpecificSnapshotSchedule(ctx, r.client, state.ID.ValueString())
	if err != nil {
		errStr := constants.UpdateSnapshotScheduleErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error Updating snapshot schedule",
			message,
		)
		return
	}

	errMap := helper.SnapshotScheduleMapper(ctx, result, &state)
	if errMap != nil {
		response.Diagnostics.AddError(constants.UpdateSnapshotScheduleErrorMessage+" with error: ", err.Error())
		return
	}

	state.RetentionTime = plan.RetentionTime
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Update Snapshot Schedule completed")
}

// Delete deletes the resource.
func (r SnapshotScheduleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Snapshot Schedule ")
	var state *models.SnapshotScheduleResource

	// Read Terraform prior state into the model
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	deleteParam := r.client.PscaleOpenAPIClient.SnapshotApi.DeleteSnapshotv1SnapshotSchedule(ctx, state.ID.ValueString())
	_, err := deleteParam.Execute()
	if err != nil {
		errStr := constants.DeleteSnapshotScheduleErrorMessage + " with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError(
			"Error deleting snapshot schedule",
			message,
		)
	}
	tflog.Info(ctx, "Deleting Snapshot Schedule completed")
}

// ImportState imports the resource state.
func (r SnapshotScheduleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
