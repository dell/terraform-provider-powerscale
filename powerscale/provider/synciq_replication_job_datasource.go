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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &SyncIQReplicationJobDataSource{}
	_ datasource.DataSourceWithConfigure = &SyncIQReplicationJobDataSource{}
)

// NewSyncIQReplicationJobDataSource creates a new syncIQ job data source.
func NewSyncIQReplicationJobDataSource() datasource.DataSource {
	return &SyncIQReplicationJobDataSource{}
}

// SyncIQReplicationJobDataSource defines the data source implementation.
type SyncIQReplicationJobDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SyncIQReplicationJobDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_replication_job"
}

// Schema describes the data source arguments.
func (d *SyncIQReplicationJobDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the SyncIQ replication jobs from PowerScale array. The information fetched from this datasource can be used for getting the details.",
		Description:         "This datasource is used to query the SyncIQ replication jobs from PowerScale array. The information fetched from this datasource can be used for getting the details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier.",
				Description:         "Identifier.",
			},
			"synciq_jobs": schema.ListNestedAttribute{
				MarkdownDescription: "List of user groups.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							Computed:            true,
							Description:         "The action to be taken by this job.",
							MarkdownDescription: "The action to be taken by this job.",
						},
						"ads_streams_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of ads streams replicated by this job.",
							MarkdownDescription: "The number of ads streams replicated by this job.",
						},
						"block_specs_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of block specs replicated by this job.",
							MarkdownDescription: "The number of block specs replicated by this job.",
						},
						"bytes_recoverable": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes recoverable by this job.",
							MarkdownDescription: "The number of bytes recoverable by this job.",
						},
						"bytes_transferred": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes that have been transferred by this job.",
							MarkdownDescription: "The number of bytes that have been transferred by this job.",
						},
						"char_specs_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of char specs replicated by this job.",
							MarkdownDescription: "The number of char specs replicated by this job.",
						},
						"committed_files": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of WORM committed files.",
							MarkdownDescription: "The number of WORM committed files.",
						},
						"corrected_lins": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs corrected by this job.",
							MarkdownDescription: "The number of LINs corrected by this job.",
						},
						"dead_node": schema.BoolAttribute{
							Computed:            true,
							Description:         "This field is true if the node running this job is dead.",
							MarkdownDescription: "This field is true if the node running this job is dead.",
						},
						"directories_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories replicated.",
							MarkdownDescription: "The number of directories replicated.",
						},
						"dirs_changed": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories changed by this job.",
							MarkdownDescription: "The number of directories changed by this job.",
						},
						"dirs_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories deleted by this job.",
							MarkdownDescription: "The number of directories deleted by this job.",
						},
						"dirs_moved": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories moved by this job.",
							MarkdownDescription: "The number of directories moved by this job.",
						},
						"dirs_new": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories created by this job.",
							MarkdownDescription: "The number of directories created by this job.",
						},
						"duration": schema.Int64Attribute{
							Computed:            true,
							Description:         "The amount of time in seconds between when the job was started and when it ended.  If the job has not yet ended, this is the amount of time since the job started.  This field is null if the job has not yet started.",
							MarkdownDescription: "The amount of time in seconds between when the job was started and when it ended.  If the job has not yet ended, this is the amount of time since the job started.  This field is null if the job has not yet started.",
						},
						"encrypted": schema.BoolAttribute{
							Computed:            true,
							Description:         "If true, syncs will be encrypted.",
							MarkdownDescription: "If true, syncs will be encrypted.",
						},
						"end_time": schema.Int64Attribute{
							Computed:            true,
							Description:         "The time the job ended in unix epoch seconds. The field is null if the job hasn't ended.",
							MarkdownDescription: "The time the job ended in unix epoch seconds. The field is null if the job hasn't ended.",
						},
						"error": schema.StringAttribute{
							Computed:            true,
							Description:         "The primary error message for this job.",
							MarkdownDescription: "The primary error message for this job.",
						},
						"error_checksum_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with checksum errors skipped by this job.",
							MarkdownDescription: "The number of files with checksum errors skipped by this job.",
						},
						"error_io_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with io errors skipped by this job.",
							MarkdownDescription: "The number of files with io errors skipped by this job.",
						},
						"error_net_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with network errors skipped by this job.",
							MarkdownDescription: "The number of files with network errors skipped by this job.",
						},
						"errors": schema.ListAttribute{
							Computed:            true,
							Description:         "A list of error messages for this job.",
							MarkdownDescription: "A list of error messages for this job.",
							ElementType:         types.StringType,
						},
						"failed_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "They number of data chunks that failed transmission.",
							MarkdownDescription: "They number of data chunks that failed transmission.",
						},
						"fifos_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of fifos replicated by this job.",
							MarkdownDescription: "The number of fifos replicated by this job.",
						},
						"file_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes transferred that belong to files.",
							MarkdownDescription: "The number of bytes transferred that belong to files.",
						},
						"files_changed": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files changed by this job.",
							MarkdownDescription: "The number of files changed by this job.",
						},
						"files_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files linked by this job.",
							MarkdownDescription: "The number of files linked by this job.",
						},
						"files_new": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files created by this job.",
							MarkdownDescription: "The number of files created by this job.",
						},
						"files_selected": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files selected by this job.",
							MarkdownDescription: "The number of files selected by this job.",
						},
						"files_transferred": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files transferred by this job.",
							MarkdownDescription: "The number of files transferred by this job.",
						},
						"files_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files unlinked by this job.",
							MarkdownDescription: "The number of files unlinked by this job.",
						},
						"files_with_ads_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with ads replicated by this job.",
							MarkdownDescription: "The number of files with ads replicated by this job.",
						},
						"flipped_lins": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs flipped by this job.",
							MarkdownDescription: "The number of LINs flipped by this job.",
						},
						"hard_links_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hard links replicated by this job.",
							MarkdownDescription: "The number of hard links replicated by this job.",
						},
						"hash_exceptions_fixed": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hash exceptions fixed by this job.",
							MarkdownDescription: "The number of hash exceptions fixed by this job.",
						},
						"hash_exceptions_found": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hash exceptions found by this job.",
							MarkdownDescription: "The number of hash exceptions found by this job.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "A unique identifier for this object.",
							MarkdownDescription: "A unique identifier for this object.",
						},
						"job_id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The ID of the job.",
							MarkdownDescription: "The ID of the job.",
						},
						"lins_total": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs transferred by this job.",
							MarkdownDescription: "The number of LINs transferred by this job.",
						},
						"network_bytes_to_source": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes sent to the source by this job.",
							MarkdownDescription: "The total number of bytes sent to the source by this job.",
						},
						"network_bytes_to_target": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes sent to the target by this job.",
							MarkdownDescription: "The total number of bytes sent to the target by this job.",
						},
						"new_files_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of new files replicated by this job.",
							MarkdownDescription: "The number of new files replicated by this job.",
						},
						"num_retransmitted_files": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files that have been retransmitted by this job.",
							MarkdownDescription: "The number of files that have been retransmitted by this job.",
						},
						"phases": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Data for each phase of this job.",
							MarkdownDescription: "Data for each phase of this job.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"phase": schema.StringAttribute{
										Computed:            true,
										Description:         "The phase that the job was in.",
										MarkdownDescription: "The phase that the job was in.",
									},
									"start_time": schema.Int64Attribute{
										Computed:            true,
										Description:         "The time the job began this phase.",
										MarkdownDescription: "The time the job began this phase.",
									},
									"end_time": schema.Int64Attribute{
										Computed:            true,
										Description:         "The time the job ended this phase.",
										MarkdownDescription: "The time the job ended this phase.",
									},
								},
							},
						},
						"policy": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "The policy associated with this job, or null if there is currently no policy associated with this job (this can happen if the job is newly created and not yet fully populated in the underlying database).",
							MarkdownDescription: "The policy associated with this job, or null if there is currently no policy associated with this job (this can happen if the job is newly created and not yet fully populated in the underlying database).",
							Attributes: map[string]schema.Attribute{
								"action": schema.StringAttribute{
									Computed:            true,
									Description:         "The action to be taken by the job.",
									MarkdownDescription: "The action to be taken by the job.",
								},
								"file_matching_pattern": schema.SingleNestedAttribute{
									Computed:            true,
									Description:         "A file matching pattern, organized as an OR'ed set of AND'ed file criteria, for example ((a AND b) OR (x AND y)) used to define a set of files with specific properties.  Policies of type 'sync' cannot use 'path' or time criteria in their matching patterns, but policies of type 'copy' can use all listed criteria.",
									MarkdownDescription: "A file matching pattern, organized as an OR'ed set of AND'ed file criteria, for example ((a AND b) OR (x AND y)) used to define a set of files with specific properties.  Policies of type 'sync' cannot use 'path' or time criteria in their matching patterns, but policies of type 'copy' can use all listed criteria.",
									Attributes: map[string]schema.Attribute{
										"or_criteria": schema.ListNestedAttribute{
											Computed:            true,
											Description:         "An array containing objects with \"and_criteria\" properties, each set of and_criteria will be logically OR'ed together to create the full file matching pattern.",
											MarkdownDescription: "An array containing objects with \"and_criteria\" properties, each set of and_criteria will be logically OR'ed together to create the full file matching pattern.",
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"and_criteria": schema.ListNestedAttribute{
														Computed:            true,
														Description:         "An array containing individual file criterion objects each describing one criterion.  These are logically AND'ed together to form a set of criteria.",
														MarkdownDescription: "An array containing individual file criterion objects each describing one criterion.  These are logically AND'ed together to form a set of criteria.",
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"case_sensitive": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "If true, the value comparison will be case sensitive.  Default is true.",
																	MarkdownDescription: "If true, the value comparison will be case sensitive.  Default is true.",
																},
																"field": schema.StringAttribute{
																	Computed:            true,
																	Description:         "The name of the file attribute to match on (only required if this is a custom_attribute type criterion).  Default is an empty string \"\".",
																	MarkdownDescription: "The name of the file attribute to match on (only required if this is a custom_attribute type criterion).  Default is an empty string \"\".",
																},
																"operator": schema.StringAttribute{
																	Computed:            true,
																	Description:         "How to compare the specified attribute of each file to the specified value.",
																	MarkdownDescription: "How to compare the specified attribute of each file to the specified value.",
																},
																"type": schema.StringAttribute{
																	Computed:            true,
																	Description:         "The type of this criterion, that is, which file attribute to match on.",
																	MarkdownDescription: "The type of this criterion, that is, which file attribute to match on.",
																},
																"value": schema.StringAttribute{
																	Computed:            true,
																	Description:         "The value to compare the specified attribute of each file to.",
																	MarkdownDescription: "The value to compare the specified attribute of each file to.",
																},
																"whole_word": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "If true, the attribute must match the entire word.  Default is true.",
																	MarkdownDescription: "If true, the attribute must match the entire word.  Default is true.",
																},
																"attribute_exists": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "For \"custom_attribute\" type criteria.  The file will match as long as the attribute named by \"field\" exists.  Default is true.",
																	MarkdownDescription: "For \"custom_attribute\" type criteria.  The file will match as long as the attribute named by \"field\" exists.  Default is true.",
																},
															},
														},
													},
												},
											},
										},
									},
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "User-assigned name of this sync policy.",
									MarkdownDescription: "User-assigned name of this sync policy.",
								},
								"source_exclude_directories": schema.ListAttribute{
									Computed:            true,
									Description:         "Directories that will be excluded from the sync.  Modifying this field will result in a full synchronization of all data.",
									MarkdownDescription: "Directories that will be excluded from the sync.  Modifying this field will result in a full synchronization of all data.",
									ElementType:         types.StringType,
								},
								"source_include_directories": schema.ListAttribute{
									Computed:            true,
									Description:         "Directories that will be included in the sync.  Modifying this field will result in a full synchronization of all data.",
									MarkdownDescription: "Directories that will be included in the sync.  Modifying this field will result in a full synchronization of all data.",
									ElementType:         types.StringType,
								},
								"source_root_path": schema.StringAttribute{
									Computed:            true,
									Description:         "The root directory on the source cluster the files will be synced from.  Modifying this field will result in a full synchronization of all data.",
									MarkdownDescription: "The root directory on the source cluster the files will be synced from.  Modifying this field will result in a full synchronization of all data.",
								},
								"target_host": schema.StringAttribute{
									Computed:            true,
									Description:         "Hostname or IP address of sync target cluster.  Modifying the target cluster host can result in the policy being unrunnable if the new target does not match the current target association.",
									MarkdownDescription: "Hostname or IP address of sync target cluster.  Modifying the target cluster host can result in the policy being unrunnable if the new target does not match the current target association.",
								},
								"target_path": schema.StringAttribute{
									Computed:            true,
									Description:         "Absolute filesystem path on the target cluster for the sync destination.",
									MarkdownDescription: "Absolute filesystem path on the target cluster for the sync destination.",
								},
							},
						},
						"policy_action": schema.StringAttribute{
							Computed:            true,
							Description:         "This is the action the policy is configured to perform.",
							MarkdownDescription: "This is the action the policy is configured to perform.",
						},
						"policy_id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the policy.",
							MarkdownDescription: "The ID of the policy.",
						},
						"policy_name": schema.StringAttribute{
							Computed:            true,
							Description:         "The name of the policy.",
							MarkdownDescription: "The name of the policy.",
						},
						"quotas_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of quotas removed from the target.",
							MarkdownDescription: "The number of quotas removed from the target.",
						},
						"regular_files_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of regular files replicated by this job.",
							MarkdownDescription: "The number of regular files replicated by this job.",
						},
						"resynced_lins": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs resynched by this job.",
							MarkdownDescription: "The number of LINs resynched by this job.",
						},
						"retransmitted_files": schema.ListAttribute{
							Computed:            true,
							Description:         "The files that have been retransmitted by this job.",
							MarkdownDescription: "The files that have been retransmitted by this job.",
							ElementType:         types.StringType,
						},
						"retry": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of times the job has been retried.",
							MarkdownDescription: "The number of times the job has been retried.",
						},
						"running_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of data chunks currently being transmitted.",
							MarkdownDescription: "The number of data chunks currently being transmitted.",
						},
						"service_report": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Data for each component exported as part of service replication.",
							MarkdownDescription: "Data for each component exported as part of service replication.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"records_skipped": schema.Int64Attribute{
										Computed:            true,
										Description:         "The number of records skipped during export.",
										MarkdownDescription: "The number of records skipped during export.",
									},
									"start_time": schema.Int64Attribute{
										Computed:            true,
										Description:         "The time the job began this component.",
										MarkdownDescription: "The time the job began this component.",
									},
									"handlers_transferred": schema.Int64Attribute{
										Computed:            true,
										Description:         "The number of handlers exported.",
										MarkdownDescription: "The number of handlers exported.",
									},
									"records_transferred": schema.Int64Attribute{
										Computed:            true,
										Description:         "The number of records exported.",
										MarkdownDescription: "The number of records exported.",
									},
									"component": schema.StringAttribute{
										Computed:            true,
										Description:         "The component that was processed.",
										MarkdownDescription: "The component that was processed.",
									},
									"records_failed": schema.Int64Attribute{
										Computed:            true,
										Description:         "The number of records failed during export.",
										MarkdownDescription: "The number of records failed during export.",
									},
									"filter": schema.ListAttribute{
										Computed:            true,
										Description:         "A list of path-based filters for exporting components.",
										MarkdownDescription: "A list of path-based filters for exporting components.",
										ElementType:         types.StringType,
									},
									"handlers_failed": schema.Int64Attribute{
										Computed:            true,
										Description:         "The number of handlers failed during export.",
										MarkdownDescription: "The number of handlers failed during export.",
									},
									"handlers_skipped": schema.Int64Attribute{
										Computed:            true,
										Description:         "The number of handlers skipped during export.",
										MarkdownDescription: "The number of handlers skipped during export.",
									},
									"status": schema.StringAttribute{
										Computed:            true,
										Description:         "The current status of export for this component.",
										MarkdownDescription: "The current status of export for this component.",
									},
									"end_time": schema.Int64Attribute{
										Computed:            true,
										Description:         "The time the job ended this component.",
										MarkdownDescription: "The time the job ended this component.",
									},
									"directory": schema.StringAttribute{
										Computed:            true,
										Description:         "The directory of the service export.",
										MarkdownDescription: "The directory of the service export.",
									},
									"error_msg": schema.ListAttribute{
										Computed:            true,
										Description:         "A list of error messages generated while exporting components.",
										MarkdownDescription: "A list of error messages generated while exporting components.",
										ElementType:         types.StringType,
									},
								},
							},
						},
						"sockets_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of sockets replicated by this job.",
							MarkdownDescription: "The number of sockets replicated by this job.",
						},
						"source_bytes_recovered": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes recovered on the source.",
							MarkdownDescription: "The number of bytes recovered on the source.",
						},
						"source_directories_created": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories created on the source.",
							MarkdownDescription: "The number of directories created on the source.",
						},
						"source_directories_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories deleted on the source.",
							MarkdownDescription: "The number of directories deleted on the source.",
						},
						"source_directories_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories linked on the source.",
							MarkdownDescription: "The number of directories linked on the source.",
						},
						"source_directories_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories unlinked on the source.",
							MarkdownDescription: "The number of directories unlinked on the source.",
						},
						"source_directories_visited": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories visited on the source.",
							MarkdownDescription: "The number of directories visited on the source.",
						},
						"source_files_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files deleted on the source.",
							MarkdownDescription: "The number of files deleted on the source.",
						},
						"source_files_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files linked on the source.",
							MarkdownDescription: "The number of files linked on the source.",
						},
						"source_files_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files unlinked on the source.",
							MarkdownDescription: "The number of files unlinked on the source.",
						},
						"sparse_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of sparse data bytes transferred by this job.",
							MarkdownDescription: "The number of sparse data bytes transferred by this job.",
						},
						"start_time": schema.Int64Attribute{
							Computed:            true,
							Description:         "The time the job started in unix epoch seconds. The field is null if the job hasn't started.",
							MarkdownDescription: "The time the job started in unix epoch seconds. The field is null if the job hasn't started.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "The state of the job.",
							MarkdownDescription: "The state of the job.",
						},
						"succeeded_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of data chunks that have been transmitted successfully.",
							MarkdownDescription: "The number of data chunks that have been transmitted successfully.",
						},
						"symlinks_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of symlinks replicated by this job.",
							MarkdownDescription: "The number of symlinks replicated by this job.",
						},
						"sync_type": schema.StringAttribute{
							Computed:            true,
							Description:         "The type of sync being performed by this job.",
							MarkdownDescription: "The type of sync being performed by this job.",
						},
						"target_bytes_recovered": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes recovered on the target.",
							MarkdownDescription: "The number of bytes recovered on the target.",
						},
						"target_directories_created": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories created on the target.",
							MarkdownDescription: "The number of directories created on the target.",
						},
						"target_directories_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories deleted on the target.",
							MarkdownDescription: "The number of directories deleted on the target.",
						},
						"target_directories_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories linked on the target.",
							MarkdownDescription: "The number of directories linked on the target.",
						},
						"target_directories_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories unlinked on the target.",
							MarkdownDescription: "The number of directories unlinked on the target.",
						},
						"target_files_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files deleted on the target.",
							MarkdownDescription: "The number of files deleted on the target.",
						},
						"target_files_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files linked on the target.",
							MarkdownDescription: "The number of files linked on the target.",
						},
						"target_files_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files unlinked on the target.",
							MarkdownDescription: "The number of files unlinked on the target.",
						},
						"target_snapshots": schema.ListAttribute{
							Computed:            true,
							Description:         "The target snapshots created by this job.",
							MarkdownDescription: "The target snapshots created by this job.",
							ElementType:         types.StringType,
						},
						"total_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of data chunks transmitted by this job.",
							MarkdownDescription: "The total number of data chunks transmitted by this job.",
						},
						"total_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes transferred by this job.",
							MarkdownDescription: "The total number of bytes transferred by this job.",
						},
						"total_exported_services": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of components exported as part of service replication.",
							MarkdownDescription: "The total number of components exported as part of service replication.",
						},
						"total_files": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files affected by this job.",
							MarkdownDescription: "The number of files affected by this job.",
						},
						"total_network_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes sent over the network by this job.",
							MarkdownDescription: "The total number of bytes sent over the network by this job.",
						},
						"total_phases": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of phases for this job.",
							MarkdownDescription: "The total number of phases for this job.",
						},
						"unchanged_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes unchanged by this job.",
							MarkdownDescription: "The number of bytes unchanged by this job.",
						},
						"up_to_date_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of up-to-date files skipped by this job.",
							MarkdownDescription: "The number of up-to-date files skipped by this job.",
						},
						"updated_files_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of updated files replicated by this job.",
							MarkdownDescription: "The number of updated files replicated by this job.",
						},
						"user_conflict_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with user conflicts skipped by this job.",
							MarkdownDescription: "The number of files with user conflicts skipped by this job.",
						},
						"warnings": schema.ListAttribute{
							Computed:            true,
							Description:         "A list of warning messages for this job.",
							MarkdownDescription: "A list of warning messages for this job.",
							ElementType:         types.StringType,
						},
						"workers": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "A list of workers for this job.",
							MarkdownDescription: "A list of workers for this job.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"last_split": schema.Int64Attribute{
										Computed:            true,
										Description:         "The last time a network split occurred.",
										MarkdownDescription: "The last time a network split occurred.",
									},
									"lin": schema.Int64Attribute{
										Computed:            true,
										Description:         "The LIN being worked on.",
										MarkdownDescription: "The LIN being worked on.",
									},
									"lnn": schema.Int64Attribute{
										Computed:            true,
										Description:         "The lnn the worker is assigned to run on.",
										MarkdownDescription: "The lnn the worker is assigned to run on.",
									},
									"process_id": schema.Int64Attribute{
										Computed:            true,
										Description:         "The process ID of the worker.",
										MarkdownDescription: "The process ID of the worker.",
									},
									"target_host": schema.StringAttribute{
										Computed:            true,
										Description:         "The target host for this worker.",
										MarkdownDescription: "The target host for this worker.",
									},
									"source_host": schema.StringAttribute{
										Computed:            true,
										Description:         "The source host for this worker.",
										MarkdownDescription: "The source host for this worker.",
									},
									"last_work": schema.Int64Attribute{
										Computed:            true,
										Description:         "The last time the worker performed work.",
										MarkdownDescription: "The last time the worker performed work.",
									},
									"worker_id": schema.Int64Attribute{
										Computed:            true,
										Description:         "The ID of the worker.",
										MarkdownDescription: "The ID of the worker.",
									},
									"connected": schema.BoolAttribute{
										Computed:            true,
										Description:         "Whether there is a connection between the source and target.",
										MarkdownDescription: "Whether there is a connection between the source and target.",
									},
								},
							},
						},
						"worm_committed_file_conflicts": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of WORM committed files which needed to be reverted. Since WORM committed files cannot be reverted, this is the number of files that were preserved in the compliance store.",
							MarkdownDescription: "The number of WORM committed files which needed to be reverted. Since WORM committed files cannot be reverted, this is the number of files that were preserved in the compliance store.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"sort": schema.StringAttribute{
						Optional:            true,
						Description:         "The field that will be used for sorting.",
						MarkdownDescription: "The field that will be used for sorting.",
					},
					"state": schema.StringAttribute{
						Optional:            true,
						Description:         "Only list SyncIQ replication jobs matching this state.",
						MarkdownDescription: "Only list SyncIQ replication jobs matching this state.",
					},
					"limit": schema.Int64Attribute{
						Optional:            true,
						Description:         "Return no more than this many results at once.",
						MarkdownDescription: "Return no more than this many results at once.",
					},
					"dir": schema.StringAttribute{
						Optional:            true,
						Description:         "The direction of the sort.",
						MarkdownDescription: "The direction of the sort.",
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *SyncIQReplicationJobDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *SyncIQReplicationJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading SyncIQ replication job data source ")

	var plan models.SyncIQReplicationJobDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageDataSourceSyncIQReplicationJob(ctx, &plan, d.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read SyncIQ replication job data source ")
}
