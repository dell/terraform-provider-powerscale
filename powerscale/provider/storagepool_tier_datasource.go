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
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoragepoolTierDataSource{}

// NewStoragepoolTierDataSource creates a new data source.
func NewStoragepoolTierDataSource() datasource.DataSource {
	return &StoragepoolTierDataSource{}
}

// StoragepoolTierDataSource defines the data source implementation.
type StoragepoolTierDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d StoragepoolTierDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagepool_tier"
}

// Schema describes the data source arguments.
func (d StoragepoolTierDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the Storagepool tiers from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale Storagepool tiers provide the ability to configure Storagepool tiers on the cluster.",
		Description: "This datasource is used to query the Storagepool tiers from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale Storagepool tiers provide the ability to configure Storagepool tiers on the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Id of Storagepool tiers. Readonly. Fixed value of \"storagepool_tiers\"",
				MarkdownDescription: "Id of Storagepool tiers. Readonly. Fixed value of \"storagepool_tiers\"",
				Computed:            true,
			},
			"storagepool_tiers": schema.ListNestedAttribute{
				Description:         "List of Storagepool tiers",
				MarkdownDescription: "List of Storagepool tiers",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Id of storagepool tier.",
							MarkdownDescription: "Id of storagepool tier.",
							Computed:            true,
						},
						"children": schema.ListAttribute{
							Description:         "The names or IDs of the tier's children.",
							MarkdownDescription: "The names or IDs of the tier's children.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"lnns": schema.ListAttribute{
							Description:         "The nodes that are part of this tier.",
							MarkdownDescription: "The nodes that are part of this tier.",
							Computed:            true,
							ElementType:         types.Int32Type,
						},
						"name": schema.StringAttribute{
							Description:         "Name of storagepool tier.",
							MarkdownDescription: "Name of storagepool tier.",
							Computed:            true,
						},
						"transfer_limit_pct": schema.Int32Attribute{
							Description:         "Stop moving files to this tier when this limit is met.",
							MarkdownDescription: "Stop moving files to this tier when this limit is met.",
							Computed:            true,
						},
						"transfer_limit_state": schema.StringAttribute{
							Description:         "How the transfer limit value is being applied.",
							MarkdownDescription: "How the transfer limit value is being applied.",
							Computed:            true,
						},
						"usage": schema.SingleNestedAttribute{
							Description:         "Usage.",
							MarkdownDescription: "Usage.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"avail_bytes": schema.StringAttribute{
									Description:         "Available free bytes remaining in the pool when virtual hot spare is taken into account.",
									MarkdownDescription: "Available free bytes remaining in the pool when virtual hot spare is taken into account.",
									Computed:            true,
								},
								"avail_hdd_bytes": schema.StringAttribute{
									Description:         "Available free bytes remaining in the pool on HDD drives when virtual hot spare is taken into account.",
									MarkdownDescription: "Available free bytes remaining in the pool on HDD drives when virtual hot spare is taken into account.",
									Computed:            true,
								},
								"avail_ssd_bytes": schema.StringAttribute{
									Description:         "Available free bytes remaining in the pool on SSD drives when virtual hot spare is taken into account.",
									MarkdownDescription: "Available free bytes remaining in the pool on SSD drives when virtual hot spare is taken into account.",
									Computed:            true,
								},
								"balanced": schema.BoolAttribute{
									Description:         "Whether or not the pool usage is currently balanced.",
									MarkdownDescription: "Whether or not the pool usage is currently balanced.",
									Computed:            true,
								},
								"free_bytes": schema.StringAttribute{
									Description:         "Free bytes remaining in the pool.",
									MarkdownDescription: "Free bytes remaining in the pool.",
									Computed:            true,
								},
								"free_hdd_bytes": schema.StringAttribute{
									Description:         "Free bytes remaining in the pool on HDD drives.",
									MarkdownDescription: "Free bytes remaining in the pool on HDD drives.",
									Computed:            true,
								},
								"free_ssd_bytes": schema.StringAttribute{
									Description:         "Free bytes remaining in the pool on SSD drives.",
									MarkdownDescription: "Free bytes remaining in the pool on SSD drives.",
									Computed:            true,
								},
								"pct_used": schema.StringAttribute{
									Description:         "Percentage of usable space in the pool which is used.",
									MarkdownDescription: "Percentage of usable space in the pool which is used.",
									Computed:            true,
								},
								"pct_used_hdd": schema.StringAttribute{
									Description:         "Percentage of usable space on HDD drives in the pool which is used.",
									MarkdownDescription: "Percentage of usable space on HDD drives in the pool which is used.",
									Computed:            true,
								},
								"pct_used_ssd": schema.StringAttribute{
									Description:         "Percentage of usable space on SSD drives in the pool which is used.",
									MarkdownDescription: "Percentage of usable space on SSD drives in the pool which is used.",
									Computed:            true,
								},
								"total_bytes": schema.StringAttribute{
									Description:         "Total bytes in the pool.",
									MarkdownDescription: "Total bytes in the pool.",
									Computed:            true,
								},
								"total_hdd_bytes": schema.StringAttribute{
									Description:         "Total bytes in the pool on HDD drives.",
									MarkdownDescription: "Total bytes in the pool on HDD drives.",
									Computed:            true,
								},
								"total_ssd_bytes": schema.StringAttribute{
									Description:         "Total bytes in the pool on SSD drives.",
									MarkdownDescription: "Total bytes in the pool on SSD drives.",
									Computed:            true,
								},
								"usable_bytes": schema.StringAttribute{
									Description:         "Total bytes in the pool drives when virtual hot spare is taken into account.",
									MarkdownDescription: "Total bytes in the pool drives when virtual hot spare is taken into account.",
									Computed:            true,
								},
								"usable_hdd_bytes": schema.StringAttribute{
									Description:         "Total bytes in the pool on HDD drives when virtual hot spare is taken into account.",
									MarkdownDescription: "Total bytes in the pool on HDD drives when virtual hot spare is taken into account.",
									Computed:            true,
								},
								"usable_ssd_bytes": schema.StringAttribute{
									Description:         "Total bytes in the pool on SSD drives when virtual hot spare is taken into account.",
									MarkdownDescription: "Total bytes in the pool on SSD drives when virtual hot spare is taken into account.",
									Computed:            true,
								},
								"used_bytes": schema.StringAttribute{
									Description:         "Used bytes in the pool.",
									MarkdownDescription: "Used bytes in the pool.",
									Computed:            true,
								},
								"used_hdd_bytes": schema.StringAttribute{
									Description:         "Used bytes in the pool on HDD drives.",
									MarkdownDescription: "Used bytes in the pool on HDD drives.",
									Computed:            true,
								},
								"used_ssd_bytes": schema.StringAttribute{
									Description:         "Used bytes in the pool on SSD drives.",
									MarkdownDescription: "Used bytes in the pool on SSD drives.",
									Computed:            true,
								},
								"virtual_hot_spare_bytes": schema.StringAttribute{
									Description:         "Bytes reserved for virtual hot spare in the pool.",
									MarkdownDescription: "Bytes reserved for virtual hot spare in the pool.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *StoragepoolTierDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d StoragepoolTierDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.StoragepoolTierDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read and map Storagepool tier state
	storagepoolTiers, err := helper.GetAllStoragepoolTiers(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Storagepool tiers", err.Error())
		return
	}
	allTiers := make([]models.StoragepoolTierModel, len(storagepoolTiers))
	for i := range len(allTiers) {
		err = helper.CopyFields(ctx, storagepoolTiers[i], &allTiers[i])
		if err != nil {
			resp.Diagnostics.AddError("Error copying fields of storagepool tiers datasource", err.Error())
			return
		}
		allTiers[i].Id = types.StringValue(strconv.Itoa(int(storagepoolTiers[i].GetId())))
		if storagepoolTiers[i].HasTransferLimitPct() {
			allTiers[i].TransferLimitPct = types.Int32Value(int32(storagepoolTiers[i].GetTransferLimitPct()))
		}
	}
	idValue := "storagepool_tier_"
	var state models.StoragepoolTierDataSourceModel
	state.ID = types.StringValue(idValue)
	state.StoragepoolTier = allTiers

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Storagepool tier data source ")
}
