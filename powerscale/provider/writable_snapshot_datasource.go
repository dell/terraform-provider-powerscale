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
	"regexp"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WritableSnapshotDataSource{}

// NewWritableSnapshotDataSource creates a new data source.
func NewWritableSnapshotDataSource() datasource.DataSource {
	return &WritableSnapshotDataSource{}
}

// WritableSnapshotDataSource defines the data source implementation.
type WritableSnapshotDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *WritableSnapshotDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_writable_snapshot"
}

// Schema describes the data source arguments.
func (d *WritableSnapshotDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = WritablesnapshotDatasourceSchema(ctx)
}

// Configure configures the data source.
func (d *WritableSnapshotDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// WritablesnapshotDatasourceSchema describes the data source arguments.
func WritablesnapshotDatasourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
			},
			"writable": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The system ID given to the writable snapshot. This is useful for debugging.",
							MarkdownDescription: "The system ID given to the writable snapshot. This is useful for debugging.",
						},
						"phys_size": schema.Int64Attribute{
							Computed:            true,
							Description:         "The amount of storage in bytes used to store this writable snapshot.",
							MarkdownDescription: "The amount of storage in bytes used to store this writable snapshot.",
						},
						"created": schema.Int64Attribute{
							Computed:            true,
							Description:         "The Unix Epoch time the writable snapshot was created.",
							MarkdownDescription: "The Unix Epoch time the writable snapshot was created.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "Writable Snapshot state.",
							MarkdownDescription: "Writable Snapshot state.",
						},
						"src_id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The system ID of the user supplied source snapshot. This is useful for debugging.",
							MarkdownDescription: "The system ID of the user supplied source snapshot. This is useful for debugging.",
						},
						"src_path": schema.StringAttribute{
							Computed:            true,
							Description:         "The /ifs path of user supplied source snapshot. This will be null for writable snapshots pending delete.",
							MarkdownDescription: "The /ifs path of user supplied source snapshot. This will be null for writable snapshots pending delete.",
						},
						"src_snap": schema.StringAttribute{
							Computed:            true,
							Description:         "The user supplied source snapshot name or ID. This will be null for writable snapshots pending delete.",
							MarkdownDescription: "The user supplied source snapshot name or ID. This will be null for writable snapshots pending delete.",
						},
						"dst_path": schema.StringAttribute{
							Computed:            true,
							Description:         "The user supplied /ifs path of writable snapshot.",
							MarkdownDescription: "The user supplied /ifs path of writable snapshot.",
						},
						"log_size": schema.Int64Attribute{
							Computed:            true,
							Description:         "The sum in bytes of logical size of files in this writable snapshot.",
							MarkdownDescription: "The sum in bytes of logical size of files in this writable snapshot.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"path": schema.StringAttribute{
						Optional:            true,
						Description:         "Only list writable snapshots matching this path.",
						MarkdownDescription: "Only list writable snapshots matching this path.",
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile("^/ifs$|^/ifs/"), "must start with /ifs or /ifs/"),
							stringvalidator.LengthBetween(4, 4096),
							stringvalidator.ConflictsWith(path.MatchRoot("filter").AtName("sort"), path.MatchRoot("filter").AtName("state"), path.MatchRoot("filter").AtName("limit"), path.MatchRoot("filter").AtName("dir"), path.MatchRoot("filter").AtName("resume")),
						},
					},
					"sort": schema.StringAttribute{
						Optional:            true,
						Description:         "The field that will be used for sorting.  Choices are path, src name, src path, created, size and state. Default is created.",
						MarkdownDescription: "The field that will be used for sorting.  Choices are path, src name, src path, created, size and state. Default is created.",
					},
					"state": schema.StringAttribute{
						Optional:            true,
						Description:         "Only list writable snapshots matching this state.",
						MarkdownDescription: "Only list writable snapshots matching this state.",
					},
					"limit": schema.Int32Attribute{
						Optional:            true,
						Description:         "Return no more than this many results at once (see resume).",
						MarkdownDescription: "Return no more than this many results at once (see resume).",
					},
					"dir": schema.StringAttribute{
						Optional:            true,
						Description:         "The direction of the sort.",
						MarkdownDescription: "The direction of the sort.",
					},
					"resume": schema.StringAttribute{
						Optional:            true,
						Description:         "Continue returning results from previous call using this token (token should come from the previous call, resume cannot be used with other options).",
						MarkdownDescription: "Continue returning results from previous call using this token (token should come from the previous call, resume cannot be used with other options).",
					},
				},
			},
		},
	}
}

// Read reads data from the data source.
func (d *WritableSnapshotDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read Terraform configuration data into the model
	var data models.WritablesnapshotModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *models.WritablesnapshotModel
	var errD error
	if data.WritableSnapshotFilter != nil && data.WritableSnapshotFilter.Path.ValueString() != "" {
		path := data.WritableSnapshotFilter.Path.ValueString()
		config, err := helper.GetWritableSnapshot(ctx, d.client, path)
		if err != nil {
			errStr := constants.ListWritableSnapshotMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading writable snapshots", message)
			return
		}
		state, errD = helper.NewWritableSnapshotDataSource(ctx, config.Writable)
	} else {
		config, err := helper.GetAllWritableSnapshots(ctx, d.client, &data)
		if err != nil {
			errStr := constants.ListWritableSnapshotMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading writable snapshots", message)
			return
		}
		state, errD = helper.NewWritableSnapshotDataSource(ctx, config.Writable)
	}
	if errD != nil {
		resp.Diagnostics.AddError("Failed to map writable snapshots fields", errD.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
