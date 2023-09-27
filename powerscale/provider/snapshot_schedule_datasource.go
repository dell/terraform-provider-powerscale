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
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &SnapshotScheduleDataSource{}
	_ datasource.DataSourceWithConfigure = &SnapshotScheduleDataSource{}
)

// NewSnapshotScheduleDataSource returns the SnapshotSchedule data source object.
func NewSnapshotScheduleDataSource() datasource.DataSource {
	return &SnapshotScheduleDataSource{}
}

// SnapshotScheduleDataSource defines the data source implementation.
type SnapshotScheduleDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SnapshotScheduleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_schedule"
}

// Schema describes the data source arguments.
func (d *SnapshotScheduleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing Snapshot Schedules from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details / for further processing in resource block. Uses are able to see information like duration, path, schedule, name etc. for the existing snapshot schedules",
		Description: "This datasource is used to query the existing Snapshot Schedules from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details / for further processing in resource block. Uses are able to see information like duration, path, schedule, name etc. for the existing snapshot schedules",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"schedules": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of snapshot schedules",
				MarkdownDescription: "List of snapshot schedules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							Description:         "Alias name to create for each snapshot.",
							MarkdownDescription: "Alias name to create for each snapshot.",
							Computed:            true,
							Optional:            true,
						},
						"duration": schema.Int64Attribute{
							Description:         "Time in seconds added to creation time to construction expiration time.",
							MarkdownDescription: "Time in seconds added to creation time to construction expiration time.",
							Computed:            true,
							Optional:            true,
						},
						"id": schema.Int64Attribute{
							Description:         "The system ID given to the schedule.",
							MarkdownDescription: "The system ID given to the schedule.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "The schedule name.",
							MarkdownDescription: "The schedule name.",
							Computed:            true,
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
							Computed:            true,
						},
						"pattern": schema.StringAttribute{
							Description:         "Pattern expanded with strftime to create snapshot names.",
							MarkdownDescription: "Pattern expanded with strftime to create snapshot names.",
							Computed:            true,
							Optional:            true,
						},
						"schedule": schema.StringAttribute{
							Description:         "The isidate compatible natural language description of the schedule.",
							MarkdownDescription: "The isidate compatible natural language description of the schedule.",
							Computed:            true,
							Optional:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Description:         "Names to filter snapshot schedules.",
						MarkdownDescription: "Names to filter snapshot schedules.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"sort": schema.StringAttribute{
						Description:         "The field that will be used for sorting. Choices are id, name, path, pattern, schedule, duration, alias, next_run, and next_snapshot. Default is id.",
						MarkdownDescription: "The field that will be used for sorting. Choices are id, name, path, pattern, schedule, duration, alias, next_run, and next_snapshot. Default is id.",
						Optional:            true,
					},
					"limit": schema.Int64Attribute{
						Description:         "Return no more than this many results at once.",
						MarkdownDescription: "Return no more than this many results at once.",
						Optional:            true,
					},
					"dir": schema.StringAttribute{
						Description:         "The direction of the sort.Supported Values:ASC , DESC",
						MarkdownDescription: "The direction of the sort.Supported Values:ASC , DESC",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (d *SnapshotScheduleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *SnapshotScheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Snapshot Schedule data source ")
	var ssPlan models.SnapshotScheduleDataSourceModel
	var ssState models.SnapshotScheduleDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &ssPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}
	snapshotSchedules, err := helper.ListSnapshotSchedules(ctx, d.client, ssPlan.SnapshotScheduleFilter)
	if err != nil {
		errStr := constants.ListSnapshotSchedulesMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading snapshot schedules ",
			message)
		return
	}

	for _, s := range snapshotSchedules {
		currentSchedule := s
		if ssPlan.SnapshotScheduleFilter == nil || len(ssPlan.SnapshotScheduleFilter.Names) == 0 {
			// If there is no filter, include all schedules
			entity := models.SnapshotScheduleEntity{}
			err := helper.CopyFields(ctx, &currentSchedule, &entity)
			if err != nil {
				resp.Diagnostics.AddError("Error reading snapshot schedule datasource plan",
					fmt.Sprintf("Could not list snapshot schedule with error: %s", err.Error()))
				return
			}
			ssState.SnapshotSchedules = append(ssState.SnapshotSchedules, entity)
		} else {
			for _, name := range ssPlan.SnapshotScheduleFilter.Names {
				if *currentSchedule.Name == name.ValueString() {
					// If the schedule name matches the filter, include it
					entity := models.SnapshotScheduleEntity{}
					err := helper.CopyFields(ctx, &currentSchedule, &entity)
					if err != nil {
						resp.Diagnostics.AddError("Error reading snapshot schedule datasource plan",
							fmt.Sprintf("Could not list snapshot schedule with error: %s", err.Error()))
						return
					}
					ssState.SnapshotSchedules = append(ssState.SnapshotSchedules, entity)
					break // No need to check further names for this schedule
				}
			}
		}
	}
	if ssPlan.SnapshotScheduleFilter != nil && len(ssPlan.SnapshotScheduleFilter.Names) > 0 &&
		(ssState.SnapshotSchedules == nil || len(ssState.SnapshotSchedules) == 0) {
		resp.Diagnostics.AddError(
			"Error reading snapshot schedule datasource plan",
			"No matching snapshot schedules found. Names might be invalid.",
		)
		return
	}
	ssState.ID = types.StringValue("1")
	ssState.SnapshotScheduleFilter = ssPlan.SnapshotScheduleFilter

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &ssState)...)
}
