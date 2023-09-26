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
		MarkdownDescription: "This datasource is used to query the existing Snapshots from PowerScale array. The information fetched from this datasource can be used for getting the details / for further processing in resource block. PowerScale Snapshots is a logical pointer to data that is stored on a cluster at a specific point in time.",
		Description:         "This datasource is used to query the existing Snapshots from PowerScale array. The information fetched from this datasource can be used for getting the details / for further processing in resource block. PowerScale Snapshot is a logical pointer to data that is stored on a cluster at a specific point in time.",
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
						"has_locks": schema.BoolAttribute{
							Description:         "True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of lock.",
							MarkdownDescription: "True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of lock.",
							Computed:            true,
						},
						"id": schema.Int64Attribute{
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
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						Optional: true,
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

	result, err := helper.GetAllSnapshots(ctx, d.client)
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

	// Apply the Path Filter if it is set
	if plan.SnapshotFilter != nil && plan.SnapshotFilter.Path.ValueString() != "" {
		for _, sdm := range fulldetail {
			if plan.SnapshotFilter.Path.ValueString() == sdm.Path.ValueString() {
				state.Snapshots = append(state.Snapshots, sdm)
			}
		}
		// If after the filter the lenght is still zero then that filter is invalid
		if len(state.Snapshots) == 0 {
			resp.Diagnostics.AddError(
				"Error getting snapshots",
				fmt.Sprintf("Path `%s` is invalid, it has no snapshots ", plan.SnapshotFilter.Path.ValueString()),
			)
			return
		}
	} else {
		state.Snapshots = append(state.Snapshots, fulldetail...)
	}
	// save into the Terraform state.
	state.ID = types.StringValue("snapshot_datasource")

	tflog.Trace(ctx, "read the snapshot datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
