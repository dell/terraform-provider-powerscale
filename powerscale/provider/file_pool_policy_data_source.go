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
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &FilePoolPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &FilePoolPolicyDataSource{}
)

// NewFilePoolPolicyDataSource creates a new data source.
func NewFilePoolPolicyDataSource() datasource.DataSource {
	return &FilePoolPolicyDataSource{}
}

// FilePoolPolicyDataSource defines the data source implementation.
type FilePoolPolicyDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *FilePoolPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filepool_policy"
}

// Schema describes the data source arguments.
func (d *FilePoolPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing File Pool Policies from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale File Pool Policy can identify logical groups of files and specify storage operations for these files.",
		Description:         "This datasource is used to query the existing File Pool Policies from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale File Pool Policy can identify logical groups of files and specify storage operations for these files.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the file pool policy instance.",
				Description:         "Unique identifier of the file pool policy instance.",
			},
			"file_pool_policies": schema.ListNestedAttribute{
				MarkdownDescription: "List of file pool policies.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "A unique name for this File Pool Policy. ",
							MarkdownDescription: "A unique name for this File Pool Policy. ",
							Computed:            true,
						},
						"file_matching_pattern": schema.SingleNestedAttribute{
							Description:         "Specifies the file matching rules for determining which files will be managed by this policy. ",
							MarkdownDescription: "Specifies the file matching rules for determining which files will be managed by this policy. ",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"or_criteria": schema.ListNestedAttribute{
									Description:         "List of or_criteria file matching rules for this policy.",
									MarkdownDescription: "List of or_criteria file matching rules for this policy.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"and_criteria": schema.ListNestedAttribute{
												Description:         "List of and_criteria file matching rules for this policy.",
												MarkdownDescription: "List of and_criteria file matching rules for this policy.",
												Computed:            true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"value": schema.StringAttribute{
															Description:         "The value to be compared against a file attribute.",
															MarkdownDescription: "The value to be compared against a file attribute.",
															Computed:            true,
														},
														"units": schema.StringAttribute{
															Description:         "Size unit value. One of 'B','KB','MB','GB','TB','PB','EB' (valid only with 'type' = 'size').",
															MarkdownDescription: "Size unit value. One of 'B','KB','MB','GB','TB','PB','EB' (valid only with 'type' = 'size').",
															Computed:            true,
														},
														"type": schema.StringAttribute{
															Description:         "The file attribute to be compared to a given value.",
															MarkdownDescription: "The file attribute to be compared to a given value.",
															Computed:            true,
														},
														"operator": schema.StringAttribute{
															Description:         "The comparison operator to use while comparing an attribute with its value.",
															MarkdownDescription: "The comparison operator to use while comparing an attribute with its value.",
															Computed:            true,
														},
														"field": schema.StringAttribute{
															Description:         "File attribute field name to be compared in a custom comparison (valid only with 'type' = 'custom_attribute').",
															MarkdownDescription: "File attribute field name to be compared in a custom comparison (valid only with 'type' = 'custom_attribute').",
															Computed:            true,
														},
														"use_relative_time": schema.BoolAttribute{
															Description:         "Whether time units refer to a calendar date and time (e.g., Jun 3, 2009) or a relative duration (e.g., 2 weeks) (valid only with 'type' in {accessed_time, birth_time, changed_time or metadata_changed_time}.",
															MarkdownDescription: "Whether time units refer to a calendar date and time (e.g., Jun 3, 2009) or a relative duration (e.g., 2 weeks) (valid only with 'type' in {accessed_time, birth_time, changed_time or metadata_changed_time}.",
															Computed:            true,
														},
														"case_sensitive": schema.BoolAttribute{
															Description:         "True to indicate case sensitivity when comparing file attributes (valid only with 'type' = 'name' or 'type' = 'path').",
															MarkdownDescription: "True to indicate case sensitivity when comparing file attributes (valid only with 'type' = 'name' or 'type' = 'path').",
															Computed:            true,
														},
														"begins_with": schema.BoolAttribute{
															Description:         "True to match the path exactly, False to match any subtree. (valid only with 'type' = 'path').",
															MarkdownDescription: "True to match the path exactly, False to match any subtree. (valid only with 'type' = 'path').",
															Computed:            true,
														},
														"attribute_exists": schema.BoolAttribute{
															Description:         "Indicates whether the existence of an attribute indicates a match (valid only with 'type' = 'custom_attribute').",
															MarkdownDescription: "Indicates whether the existence of an attribute indicates a match (valid only with 'type' = 'custom_attribute').",
															Computed:            true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"actions": schema.ListNestedAttribute{
							Description:         "A list of actions to be taken for matching files. ",
							MarkdownDescription: "A list of actions to be taken for matching files. ",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"enable_packing_action": schema.BoolAttribute{
										Description:         "Action for enable_packing type. True to enable enable_packing action.",
										MarkdownDescription: "Action for enable_packing type. True to enable enable_packing action.",
										Computed:            true,
									},
									"enable_coalescer_action": schema.BoolAttribute{
										Description:         "Action for enable_coalescer type. Set write performance optimization. True to enable SmartCache action.",
										MarkdownDescription: "Action for enable_coalescer type. Set write performance optimization. True to enable SmartCache action.",
										Computed:            true,
									},
									"data_access_pattern_action": schema.StringAttribute{
										Description:         "Action for set_data_access_pattern type. Set data access pattern optimization. Acceptable values: random, concurrency, streaming.",
										MarkdownDescription: "Action for set_data_access_pattern type. Set data access pattern optimization. Acceptable values: random, concurrency, streaming.",
										Computed:            true,
									},
									"requested_protection_action": schema.StringAttribute{
										Description:         "Action for set_requested_protection type. Acceptable values: default, +1n, +2d:1n, +2n, +3d:1n, +3d:1n1d, +3n, +4d:1n, +4d:2n, +4n, 2x, 3x, 4x, 5x, 6x, 7x, 8x.",
										MarkdownDescription: "Action for set_requested_protection type. Acceptable values: default, +1n, +2d:1n, +2n, +3d:1n, +3d:1n1d, +3n, +4d:1n, +4d:2n, +4n, 2x, 3x, 4x, 5x, 6x, 7x, 8x.",
										Computed:            true,
									},
									"data_storage_policy_action": schema.SingleNestedAttribute{
										Description:         "Action for apply_data_storage_policy.",
										MarkdownDescription: "Action for apply_data_storage_policy.",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"ssd_strategy": schema.StringAttribute{
												Description:         "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
												MarkdownDescription: "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
												Computed:            true,
											},
											"storagepool": schema.StringAttribute{
												Description:         "Specifies the storage target.",
												MarkdownDescription: "Specifies the storage target.",
												Computed:            true,
											},
										},
									},
									"snapshot_storage_policy_action": schema.SingleNestedAttribute{
										Description:         "Action for apply_snapshot_storage_policy.",
										MarkdownDescription: "Action for apply_snapshot_storage_policy.",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"ssd_strategy": schema.StringAttribute{
												Description:         "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
												MarkdownDescription: "Specifies the SSD strategy. Acceptable values: metadata, metadata-write, data, avoid.",
												Computed:            true,
											},
											"storagepool": schema.StringAttribute{
												Description:         "Specifies the snapshot storage target.",
												MarkdownDescription: "Specifies the snapshot storage target.",
												Computed:            true,
											},
										},
									},
									"cloudpool_policy_action": schema.SingleNestedAttribute{
										Description:         "Action for set_cloudpool_policy type.",
										MarkdownDescription: "Action for set_cloudpool_policy type.",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"pool": schema.StringAttribute{
												Description:         "Specifies the cloudPool storage target.",
												MarkdownDescription: "Specifies the cloudPool storage target.",
												Computed:            true,
											},
											"archive_snapshot_files": schema.BoolAttribute{
												Description:         "Specifies if files with snapshots should be archived.",
												MarkdownDescription: "Specifies if files with snapshots should be archived.",
												Computed:            true,
											},
											"compression": schema.BoolAttribute{
												Description:         "Specifies if files should be compressed.",
												MarkdownDescription: "Specifies if files should be compressed.",
												Computed:            true,
											},
											"encryption": schema.BoolAttribute{
												Description:         "Specifies if files should be encrypted.",
												MarkdownDescription: "Specifies if files should be encrypted.",
												Computed:            true,
											},
											"data_retention": schema.Int64Attribute{
												Description:         "Specifies the minimum amount of time archived data will be retained in the cloud after deletion.",
												MarkdownDescription: "Specifies the minimum amount of time archived data will be retained in the cloud after deletion.",
												Computed:            true,
											},
											"full_backup_retention": schema.Int64Attribute{
												Description:         "The minimum amount of time cloud files will be retained after the creation of a full NDMP backup. (Used with NDMP backups only.  Not applicable to SyncIQ.) ",
												MarkdownDescription: "The minimum amount of time cloud files will be retained after the creation of a full NDMP backup. (Used with NDMP backups only.  Not applicable to SyncIQ.) ",
												Computed:            true,
											},
											"incremental_backup_retention": schema.Int64Attribute{
												Description:         "The minimum amount of time cloud files will be retained after the creation of a SyncIQ backup or an incremental NDMP backup. (Used with SyncIQ and NDMP backups.) ",
												MarkdownDescription: "The minimum amount of time cloud files will be retained after the creation of a SyncIQ backup or an incremental NDMP backup. (Used with SyncIQ and NDMP backups.) ",
												Computed:            true,
											},
											"writeback_frequency": schema.Int64Attribute{
												Description:         "The minimum amount of time to wait before updating cloud data with local changes.",
												MarkdownDescription: "The minimum amount of time to wait before updating cloud data with local changes.",
												Computed:            true,
											},
											"cache": schema.SingleNestedAttribute{
												Description:         "Specifies default cloudpool cache settings for new filepool policies.",
												MarkdownDescription: "Specifies default cloudpool cache settings for new filepool policies.",
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"expiration": schema.Int64Attribute{
														Description:         "Specifies cache expiration.",
														MarkdownDescription: "Specifies cache expiration.",
														Computed:            true,
													},
													"read_ahead": schema.StringAttribute{
														Description:         "Specifies cache read ahead type. Acceptable values: partial, full.",
														MarkdownDescription: "Specifies cache read ahead type. Acceptable values: partial, full.",
														Computed:            true,
													},
													"type": schema.StringAttribute{
														Description:         "Specifies cache type. Acceptable values: cached, no-cache.",
														MarkdownDescription: "Specifies cache type. Acceptable values: cached, no-cache.",
														Computed:            true,
													},
												},
											},
										},
									},
									"action_type": schema.StringAttribute{
										Description:         "action_type Acceptable values: set_requested_protection, set_data_access_pattern, enable_coalescer, apply_data_storage_policy, apply_snapshot_storage_policy, set_cloudpool_policy, enable_packing.							",
										MarkdownDescription: "action_type Acceptable values: set_requested_protection, set_data_access_pattern, enable_coalescer, apply_data_storage_policy, apply_snapshot_storage_policy, set_cloudpool_policy, enable_packing.",
										Computed:            true,
									},
								},
							},
						},
						"apply_order": schema.Int64Attribute{
							Description:         "The order in which this policy should be applied (relative to other policies). ",
							MarkdownDescription: "The order in which this policy should be applied (relative to other policies). ",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "A description for this File Pool Policy. ",
							MarkdownDescription: "A description for this File Pool Policy. ",
							Computed:            true,
						},
						"birth_cluster_id": schema.StringAttribute{
							Description:         "The guid assigned to the cluster on which the policy was created.",
							MarkdownDescription: "The guid assigned to the cluster on which the policy was created.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "A unique name for this File Pool Policy.",
							MarkdownDescription: "A unique name for this File Pool Policy.",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							Description:         "Indicates whether this policy is in a good state (\"OK\") or disabled (\"disabled\").",
							MarkdownDescription: "Indicates whether this policy is in a good state (\"OK\") or disabled (\"disabled\").",
							Computed:            true,
						},
						"state_details": schema.StringAttribute{
							Description:         "Gives further information to describe the state of this policy.",
							MarkdownDescription: "Gives further information to describe the state of this policy.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *FilePoolPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *FilePoolPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading File Pool Policy data source ")

	var state models.FilePoolPolicyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := helper.GetAllFilePoolPolicies(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error listing of PowerScale File Pool Policies.", err.Error())
		return
	}

	defaultPolicyState, err := helper.GetFilePoolDefaultPolicyDataSourceState(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting Default File Pool Policy", err.Error())
		return
	}

	// parse response to state model
	if err := helper.UpdateFilePoolPolicyDataSourceState(ctx, &state, policies); err != nil {
		resp.Diagnostics.AddError("Error reading File Pool Policies datasource plan",
			fmt.Sprintf("Could not list File Pool Policies with error: %s", err.Error()))
		return
	}

	// filter by names
	if state.Filter != nil && len(state.Filter.Names) > 0 {
		var validPolicies []string
		var filteredPolicies []models.FilePoolPolicyDetailModel
		filterLen := len(state.Filter.Names)
		for _, name := range state.Filter.Names {
			if name.ValueString() == helper.FilePoolDefaultPolicyName {
				filteredPolicies = append(filteredPolicies, *defaultPolicyState)
				validPolicies = append(validPolicies, name.ValueString())
			}
			for _, policy := range state.FilePoolPolicies {
				if name.ValueString() == helper.FilePoolDefaultPolicyName && policy.Name.Equal(name) {
					filterLen++
				}
				if policy.Name.Equal(name) {
					filteredPolicies = append(filteredPolicies, policy)
					validPolicies = append(validPolicies, policy.Name.ValueString())
					break
				}
			}
		}
		state.FilePoolPolicies = filteredPolicies
		if filterLen != len(state.FilePoolPolicies) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered File Pool Policy names is not a valid powerscale File Pool Policy.",
				fmt.Sprintf("Valid File Pool Policies: [%v], filtered list: [%v]", strings.Join(validPolicies, " , "), state.Filter.Names),
			)
		}
	} else {
		state.FilePoolPolicies = append(state.FilePoolPolicies, *defaultPolicyState)
	}

	state.ID = types.StringValue("filepool_policy_datasource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read File Pool Policy data source ")
}
