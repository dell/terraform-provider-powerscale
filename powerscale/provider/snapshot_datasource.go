/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SnapshotDataSource{}

// NewSnapshotDataSource creates a new data source.
func NewSnapshotDataSource() datasource.DataSource {
	return &SnapshotDataSource{}
}

// SnapshotDataSource defines the data source implementation.
type SnapshotDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SnapshotDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

// Schema describes the data source arguments.
func (d *SnapshotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing Snapshots from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Snapshots is a logical pointer to data that is stored on a cluster at a specific point in time.",
		Description:         "This datasource is used to query the existing Snapshots from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Snapshot is a logical pointer to data that is stored on a cluster at a specific point in time.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"snapshots_details": schema.ListNestedAttribute{
				Description:         "List of Snapshots",
				MarkdownDescription: "List of Snapshots",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"alias": schema.StringAttribute{
							Description:         "The name of the alias, none for real snapshots.",
							MarkdownDescription: "The name of the alias, none for real snapshots.",
							Computed:            true,
						},
						"created": schema.Int64Attribute{
							Description:         "The Unix Epoch time the snapshot was created.",
							MarkdownDescription: "The Unix Epoch time the snapshot was created.",
							Computed:            true,
						},
						"expires": schema.Int64Attribute{
							Description:         "The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.",
							MarkdownDescription: "The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.",
							Computed:            true,
						},
						"set_expires": schema.StringAttribute{
							Description:         "The amount of time from creation before the snapshot will expire and be eligible for automatic deletion.",
							MarkdownDescription: "The amount of time from creation before the snapshot will expire and be eligible for automatic deletion.",
							Optional:            true,
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
						"name": schema.StringAttribute{
							Description:         "The user or system supplied snapshot name. This will be null for snapshots pending delete.",
							MarkdownDescription: "The user or system supplied snapshot name. This will be null for snapshots pending delete.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							Description:         "The /ifs path snapshotted.",
							MarkdownDescription: "The /ifs path snapshotted.",
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
							Description:         "The ID of the snapshot pointed to if this is an alias. An alias to the live filesystem is represented by the value -1.",
							MarkdownDescription: "The ID of the snapshot pointed to if this is an alias. An alias to the live filesystem is represented by the value -1.",
							Computed:            true,
						},
						"target_name": schema.StringAttribute{
							Description:         "The name of the snapshot pointed to if this is an alias.",
							MarkdownDescription: "The name of the snapshot pointed to if this is an alias.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRoot("filter").AtName("limit")),
						},
					},

					"name": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRoot("filter").AtName("sort"), path.MatchRoot("filter").AtName("state"), path.MatchRoot("filter").AtName("limit"), path.MatchRoot("filter").AtName("dir"), path.MatchRoot("filter").AtName("path"), path.MatchRoot("filter").AtName("schedule"), path.MatchRoot("filter").AtName("type")),
						},
					},

					"sort": schema.StringAttribute{
						Optional:            true,
						Description:         "The field that will be used for sorting.",
						MarkdownDescription: "The field that will be used for sorting.",
					},

					"limit": schema.Int64Attribute{
						Optional:            true,
						Description:         "Return no more than this many results at once (see resume).",
						MarkdownDescription: "Return no more than this many results at once (see resume).",
					},

					"dir": schema.StringAttribute{
						Optional:            true,
						Description:         "The direction of the sort.",
						MarkdownDescription: "The direction of the sort.",
					},

					"state": schema.StringAttribute{
						Optional:            true,
						Description:         "The state of the snapshot.",
						MarkdownDescription: "The state of the snapshot.",
					},

					"type": schema.StringAttribute{
						Optional:            true,
						Description:         "The type of the snapshot.",
						MarkdownDescription: "The type of the snapshot.",
					},

					"schedule": schema.StringAttribute{
						Optional:            true,
						Description:         "The schedule of the snapshot.",
						MarkdownDescription: "The schedule of the snapshot.",
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *SnapshotDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *SnapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.SnapshotDataSourceModel
	var plan models.SnapshotDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := helper.GetAllSnapshots(ctx, d.client, &plan)
	if err != nil {
		errStr := constants.ReadSnapshotErrorMessage + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of snapshots",
			message,
		)
		return
	}

	// Do the TF Mapping
	fulldetail := []models.SnapshotDetailModel{}
	for _, vsse := range result {
		detail, err := helper.SnapshotDetailMapper(ctx, vsse)
		if err != nil {
			errStr := constants.ReadSnapshotErrorMessage + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error getting the list of snapshots",
				message,
			)
			return
		}
		fulldetail = append(fulldetail, detail)
	}

	// Apply the Name Filter if it is set
	if plan.SnapshotFilter != nil && plan.SnapshotFilter.Name.ValueString() != "" {
		for _, sdm := range fulldetail {
			if plan.SnapshotFilter.Name.ValueString() == sdm.Name.ValueString() {
				state.Snapshots = append(state.Snapshots, sdm)
			}
		}
		// If after the filter the length is still zero then that filter is invalid
		if len(state.Snapshots) == 0 {
			resp.Diagnostics.AddError(
				"Error getting snapshots",
				fmt.Sprintf("Unable to find snapshot with name `%s`", plan.SnapshotFilter.Name.ValueString()),
			)
			return
		}
	}
	// Apply the Path Filter if it is set
	if plan.SnapshotFilter != nil && plan.SnapshotFilter.Path.ValueString() != "" {

		if !plan.SnapshotFilter.Limit.IsNull() {
			resp.Diagnostics.AddError(
				"Error getting snapshots",
				"Path filter cannot be applied along with limit",
			)
			return
		}

		for _, sdm := range fulldetail {
			if plan.SnapshotFilter.Path.ValueString() == sdm.Path.ValueString() {
				state.Snapshots = append(state.Snapshots, sdm)
			}
		}
		// If after the filter the length is still zero then that filter is invalid
		if len(state.Snapshots) == 0 {
			resp.Diagnostics.AddError(
				"Error getting snapshots",
				fmt.Sprintf("Unable to find snapshot with path `%s`", plan.SnapshotFilter.Path.ValueString()),
			)
			return
		}
	} else {
		state.Snapshots = append(state.Snapshots, fulldetail...)
	}

	// save into the Terraform state.
	state.ID = types.StringValue("snapshot_datasource")

	state.SnapshotFilter = plan.SnapshotFilter

	tflog.Trace(ctx, "read the snapshot datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
