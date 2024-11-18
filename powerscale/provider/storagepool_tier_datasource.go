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
	"math"
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
var _ datasource.DataSource = &StoragepoolTierDataSource{}

// NewSmartPoolSettingDataSource creates a new data source.
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
		MarkdownDescription: "This datasource is used to query the StoragePool tiers from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale StoragePool tiers provide the ability to configure SmartPools on the cluster.",
		Description: "This datasource is used to query the StoragePool tiers from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale StoragePool tiers provide the ability to configure SmartPools on the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Id of StoragePool tiers. Readonly. Fixed value of \"storagepool_tiers\"",
				MarkdownDescription: "Id of StoragePool tiers. Readonly. Fixed value of \"storagepool_tiers\"",
				Optional:            false,
				Required:            false,
				Computed:            true,
			}, // Need to created nested attributes
			"storagepool_tiers": schema.ListNestedAttribute{
				Description:         "List of StoragePool tiers",
				MarkdownDescription: "List of StoragePool tiers",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Id of storagepool tier.",
							MarkdownDescription: "Id of storagepool tier.",
							Computed:            true,
						},
						"children": schema.ListAttribute{
							Description:         "Manage I/O optimization settings.",
							MarkdownDescription: "Manage I/O optimization settings.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"lnns": schema.ListAttribute{
							Description:         "Apply to files with manually-managed I/O optimization settings.",
							MarkdownDescription: "Apply to files with manually-managed I/O optimization settings.",
							Computed:            true,
							ElementType:         types.Int32Type,
						},
						"name": schema.StringAttribute{
							Description:         "Manage protection settings.",
							MarkdownDescription: "Manage protection settings.",
							Computed:            true,
						},
						// "node_type_ids": schema.ListAttribute{
						// 	Description:         "Apply to files with manually-managed protection.",
						// 	MarkdownDescription: "Apply to files with manually-managed protection.",
						// 	Computed:            true,
						// 	ElementType:         types.StringType,
						// },
						"transfer_limit_pct": schema.Int64Attribute{
							Description:         "Enable global namespace acceleration.",
							MarkdownDescription: "Enable global namespace acceleration.",
							Computed:            true,
						},
						"transfer_limit_state": schema.StringAttribute{
							Description:         "Whether or not namespace operation optimizations are currently in effect.",
							MarkdownDescription: "Whether or not namespace operation optimizations are currently in effect.",
							Computed:            true,
						},
						"usage": schema.SingleNestedAttribute{
							Description:         "Usage.",
							MarkdownDescription: "Usage.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"avail_bytes": schema.StringAttribute{
									Description:         "Avail bytes",
									MarkdownDescription: "Avail bytes",
									Computed:            true,
								},
								"avail_hdd_bytes": schema.StringAttribute{
									Description:         "Avail hdd bytes",
									MarkdownDescription: "Avail hdd bytes",
									Computed:            true,
								},
								"avail_ssd_bytes": schema.StringAttribute{
									Description:         "Avail ssd bytes",
									MarkdownDescription: "Avail ssd bytes",
									Computed:            true,
								},
								"balanced": schema.BoolAttribute{
									Description:         "Balanced",
									MarkdownDescription: "Balanced",
									Computed:            true,
								},
								"free_bytes": schema.StringAttribute{
									Description:         "Free bytes",
									MarkdownDescription: "Free bytes",
									Computed:            true,
								},
								"free_hdd_bytes": schema.StringAttribute{
									Description:         "Free hdd bytes",
									MarkdownDescription: "Free hdd bytes",
									Computed:            true,
								},
								"free_ssd_bytes": schema.StringAttribute{
									Description:         "Free ssd bytes",
									MarkdownDescription: "Free ssd bytes",
									Computed:            true,
								},
								"pct_used": schema.StringAttribute{
									Description:         "Pct used",
									MarkdownDescription: "Pct used",
									Computed:            true,
								},
								"pct_used_hdd": schema.StringAttribute{
									Description:         "Pct used hdd",
									MarkdownDescription: "Pct used hdd",
									Computed:            true,
								},
								"pct_used_ssd": schema.StringAttribute{
									Description:         "Pct used ssd",
									MarkdownDescription: "Pct used ssd",
									Computed:            true,
								},
								"total_bytes": schema.StringAttribute{
									Description:         "Total bytes",
									MarkdownDescription: "Total bytes",
									Computed:            true,
								},
								"total_hdd_bytes": schema.StringAttribute{
									Description:         "Total hdd bytes",
									MarkdownDescription: "Total hdd bytes",
									Computed:            true,
								},
								"total_ssd_bytes": schema.StringAttribute{
									Description:         "Total ssd bytes",
									MarkdownDescription: "Total ssd bytes",
									Computed:            true,
								},
								"usable_bytes": schema.StringAttribute{
									Description:         "Usable bytes",
									MarkdownDescription: "Usable bytes",
									Computed:            true,
								},
								"usable_hdd_bytes": schema.StringAttribute{
									Description:         "Usable hdd bytes",
									MarkdownDescription: "Usable hdd bytes",
									Computed:            true,
								},
								"usable_ssd_bytes": schema.StringAttribute{
									Description:         "Usable ssd bytes",
									MarkdownDescription: "Usable ssd bytes",
									Computed:            true,
								},
								"used_bytes": schema.StringAttribute{
									Description:         "Used bytes",
									MarkdownDescription: "Used bytes",
									Computed:            true,
								},
								"used_hdd_bytes": schema.StringAttribute{
									Description:         "Used hdd bytes",
									MarkdownDescription: "Used hdd bytes",
									Computed:            true,
								},
								"used_ssd_bytes": schema.StringAttribute{
									Description:         "Used ssd bytes",
									MarkdownDescription: "Used ssd bytes",
									Computed:            true,
								},
								"virtual_hot_spare_bytes": schema.StringAttribute{
									Description:         "Virtual hot spare bytes",
									MarkdownDescription: "Virtual hot spare bytes",
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

	// Read and map StoragePool tier state
	storagepoolTiers, err := helper.GetAllStoragepoolTiers(ctx, d.client)
	if err != nil {
		errStr := constants.ReadSmartPoolSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading Storagepool tiers", message)
		return
	}
	allTiers := make([]models.StoragepoolTierModel, len(storagepoolTiers.GetTiers()))
	eachTierValue := storagepoolTiers.GetTiers()
	for i := range len(allTiers) {
		err = helper.CopyFields(ctx, eachTierValue[i], &allTiers[i])
		if err != nil {
			resp.Diagnostics.AddError("Error copying fields of storagepool tiers datasource", err.Error())
			return
		}
		allTiers[i].Id = types.StringValue(strconv.Itoa(int(eachTierValue[i].GetId())))
		allTiers[i].TransferLimitPct = types.Int64Value(int64(math.Round(eachTierValue[i].GetTransferLimitPct())))
	}
	idValue := "storagepool_tier_"
	var state models.StoragepoolTierDataSourceModel
	state.ID = types.StringValue(idValue)
	state.StoragepoolTier = allTiers

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Storagepool tier data source ")
}
