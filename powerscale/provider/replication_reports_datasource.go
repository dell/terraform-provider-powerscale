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
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ReplicationReportDataSource{}

// NewReplicationReportsDataSource creates a new data source.
func NewReplicationReportDataSource() datasource.DataSource {
	return &ReplicationReportDataSource{}
}

// ReplicationReportsDataSource defines the data source implementation.
type ReplicationReportDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *ReplicationReportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication_report"
}

// Schema describes the data source arguments.
func (d *ReplicationReportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the network pool instance.",
				MarkdownDescription: "Unique identifier of the network pool instance.",
				Computed:            true,
			},
			"replication_reports": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"total_phases": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of phases for this job.",
							MarkdownDescription: "The total number of phases for this job.",
						},
						"symlinks_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of symlinks replicated by this job.",
							MarkdownDescription: "The number of symlinks replicated by this job.",
						},
						"block_specs_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of block specs replicated by this job.",
							MarkdownDescription: "The number of block specs replicated by this job.",
						},
						"encrypted": schema.BoolAttribute{
							Computed:            true,
							Description:         "If true, syncs will be encrypted.",
							MarkdownDescription: "If true, syncs will be encrypted.",
						},
						"source_bytes_recovered": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes recovered on the source.",
							MarkdownDescription: "The number of bytes recovered on the source.",
						},
						"network_bytes_to_source": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes sent to the source by this job.",
							MarkdownDescription: "The total number of bytes sent to the source by this job.",
						},
						"error": schema.StringAttribute{
							Computed:            true,
							Description:         "The primary error message for this job.",
							MarkdownDescription: "The primary error message for this job.",
						},
						"lins_total": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs transferred by this job.",
							MarkdownDescription: "The number of LINs transferred by this job.",
						},
						"job_id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The ID of the job.",
							MarkdownDescription: "The ID of the job.",
						},
						"source_directories_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories deleted on the source.",
							MarkdownDescription: "The number of directories deleted on the source.",
						},
						"throughput": schema.StringAttribute{
							Computed:            true,
							Description:         "Throughput of a job",
							MarkdownDescription: "Throughput of a job",
						},
						"phases": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Data for each phase of this job.",
							MarkdownDescription: "Data for each phase of this job.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"statistics": schema.SingleNestedAttribute{
										Computed:            true,
										Description:         "Statistics for each phase of this job.",
										MarkdownDescription: "Statistics for each phase of this job.",
										Attributes: map[string]schema.Attribute{
											"compliance_dir_links": schema.StringAttribute{
												Computed:            true,
												Description:         "Compliance Dir Links",
												MarkdownDescription: "Compliance Dir Links",
											},
											"files": schema.StringAttribute{
												Computed:            true,
												Description:         "Files",
												MarkdownDescription: "Files",
											},
											"linked_dirs": schema.StringAttribute{
												Computed:            true,
												Description:         "Linked Dirs",
												MarkdownDescription: "Linked Dirs",
											},
											"marked_files": schema.StringAttribute{
												Computed:            true,
												Description:         "Marked Files",
												MarkdownDescription: "Marked Files",
											},
											"resynced_lins": schema.StringAttribute{
												Computed:            true,
												Description:         "Resynced LINs",
												MarkdownDescription: "Resynced LINs",
											},
											"linked_files": schema.StringAttribute{
												Computed:            true,
												Description:         "Linked Files",
												MarkdownDescription: "Linked Files",
											},
											"new_compliance_dirs": schema.StringAttribute{
												Computed:            true,
												Description:         "New Compliance Dirs",
												MarkdownDescription: "New Compliance Dirs",
											},
											"modified_dirs": schema.StringAttribute{
												Computed:            true,
												Description:         "Modified Dirs",
												MarkdownDescription: "Modified Dirs",
											},
											"deleted_dirs": schema.StringAttribute{
												Computed:            true,
												Description:         "Deleted Dirs",
												MarkdownDescription: "Deleted Dirs",
											},
											"flipped_lins": schema.StringAttribute{
												Computed:            true,
												Description:         "Flipped LINs",
												MarkdownDescription: "Flipped LINs",
											},
											"hash_exceptions": schema.StringAttribute{
												Computed:            true,
												Description:         "Hash Exceptions",
												MarkdownDescription: "Hash Exceptions",
											},
											"new_files": schema.StringAttribute{
												Computed:            true,
												Description:         "New Files",
												MarkdownDescription: "New Files",
											},
											"new_resynced_files": schema.StringAttribute{
												Computed:            true,
												Description:         "New Resynced Files",
												MarkdownDescription: "New Resynced Files",
											},
											"resynced_file_links": schema.StringAttribute{
												Computed:            true,
												Description:         "Resynced File Links",
												MarkdownDescription: "Resynced File Links",
											},
											"unlinked_files": schema.StringAttribute{
												Computed:            true,
												Description:         "Unlinked Files",
												MarkdownDescription: "Unlinked Files",
											},
											"dirs": schema.StringAttribute{
												Computed:            true,
												Description:         "Dirs",
												MarkdownDescription: "Dirs",
											},
											"modified_files": schema.StringAttribute{
												Computed:            true,
												Description:         "Modified Files",
												MarkdownDescription: "Modified Files",
											},
											"corrected_lins": schema.StringAttribute{
												Computed:            true,
												Description:         "Corrected LINs",
												MarkdownDescription: "Corrected LINs",
											},
											"new_dirs": schema.StringAttribute{
												Computed:            true,
												Description:         "New Dirs",
												MarkdownDescription: "New Dirs",
											},
											"modified_lins": schema.StringAttribute{
												Computed:            true,
												Description:         "Modified LINs",
												MarkdownDescription: "Modified LINs",
											},
											"marked_directories": schema.StringAttribute{
												Computed:            true,
												Description:         "Marked Directories",
												MarkdownDescription: "Marked Directories",
											},
											"deleted_files": schema.StringAttribute{
												Computed:            true,
												Description:         "Deleted Files",
												MarkdownDescription: "Deleted Files",
											},
										},
									},
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
						"policy_id": schema.StringAttribute{
							Computed:            true,
							Description:         "The ID of the policy.",
							MarkdownDescription: "The ID of the policy.",
						},
						"target_directories_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories deleted on the target.",
							MarkdownDescription: "The number of directories deleted on the target.",
						},
						"files_changed": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files changed by this job.",
							MarkdownDescription: "The number of files changed by this job.",
						},
						"dirs_changed": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories changed by this job.",
							MarkdownDescription: "The number of directories changed by this job.",
						},
						"target_files_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files deleted on the target.",
							MarkdownDescription: "The number of files deleted on the target.",
						},
						"source_directories_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories unlinked on the source.",
							MarkdownDescription: "The number of directories unlinked on the source.",
						},
						"dirs_moved": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories moved by this job.",
							MarkdownDescription: "The number of directories moved by this job.",
						},
						"source_files_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files deleted on the source.",
							MarkdownDescription: "The number of files deleted on the source.",
						},
						"error_io_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with io errors skipped by this job.",
							MarkdownDescription: "The number of files with io errors skipped by this job.",
						},
						"total_network_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes sent over the network by this job.",
							MarkdownDescription: "The total number of bytes sent over the network by this job.",
						},
						"error_net_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with network errors skipped by this job.",
							MarkdownDescription: "The number of files with network errors skipped by this job.",
						},
						"warnings": schema.ListAttribute{
							Computed:            true,
							Description:         "A list of warning messages for this job.",
							MarkdownDescription: "A list of warning messages for this job.",
							ElementType:         types.StringType,
						},
						"char_specs_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of char specs replicated by this job.",
							MarkdownDescription: "The number of char specs replicated by this job.",
						},
						"total_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes transferred by this job.",
							MarkdownDescription: "The total number of bytes transferred by this job.",
						},
						"subreport_count": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of subreports that are available for this job report.",
							MarkdownDescription: "The number of subreports that are available for this job report.",
						},
						"sparse_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of sparse data bytes transferred by this job.",
							MarkdownDescription: "The number of sparse data bytes transferred by this job.",
						},
						"action": schema.StringAttribute{
							Computed:            true,
							Description:         "The action to be taken by this job.",
							MarkdownDescription: "The action to be taken by this job.",
						},
						"source_files_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files unlinked on the source.",
							MarkdownDescription: "The number of files unlinked on the source.",
						},
						"error_checksum_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with checksum errors skipped by this job.",
							MarkdownDescription: "The number of files with checksum errors skipped by this job.",
						},
						"up_to_date_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of up-to-date files skipped by this job.",
							MarkdownDescription: "The number of up-to-date files skipped by this job.",
						},
						"unchanged_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes unchanged by this job.",
							MarkdownDescription: "The number of bytes unchanged by this job.",
						},
						"hash_exceptions_fixed": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hash exceptions fixed by this job.",
							MarkdownDescription: "The number of hash exceptions fixed by this job.",
						},
						"target_files_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files unlinked on the target.",
							MarkdownDescription: "The number of files unlinked on the target.",
						},
						"new_files_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of new files replicated by this job.",
							MarkdownDescription: "The number of new files replicated by this job.",
						},
						"directories_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories replicated.",
							MarkdownDescription: "The number of directories replicated.",
						},
						"end_time": schema.Int64Attribute{
							Computed:            true,
							Description:         "The time the job ended in unix epoch seconds. The field is null if the job hasn't ended.",
							MarkdownDescription: "The time the job ended in unix epoch seconds. The field is null if the job hasn't ended.",
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
						"source_files_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files linked on the source.",
							MarkdownDescription: "The number of files linked on the source.",
						},
						"files_new": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files created by this job.",
							MarkdownDescription: "The number of files created by this job.",
						},
						"total_files": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files affected by this job.",
							MarkdownDescription: "The number of files affected by this job.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "A unique identifier for this object.",
							MarkdownDescription: "A unique identifier for this object.",
						},
						"dirs_new": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories created by this job.",
							MarkdownDescription: "The number of directories created by this job.",
						},
						"target_directories_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories linked on the target.",
							MarkdownDescription: "The number of directories linked on the target.",
						},
						"bytes_recoverable": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes recoverable by this job.",
							MarkdownDescription: "The number of bytes recoverable by this job.",
						},
						"corrected_lins": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs corrected by this job.",
							MarkdownDescription: "The number of LINs corrected by this job.",
						},
						"files_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files linked by this job.",
							MarkdownDescription: "The number of files linked by this job.",
						},
						"target_bytes_recovered": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes recovered on the target.",
							MarkdownDescription: "The number of bytes recovered on the target.",
						},
						"worm_committed_file_conflicts": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of WORM committed files which needed to be reverted. Since WORM committed files cannot be reverted, this is the number of files that were preserved in the compliance store.",
							MarkdownDescription: "The number of WORM committed files which needed to be reverted. Since WORM committed files cannot be reverted, this is the number of files that were preserved in the compliance store.",
						},
						"user_conflict_files_skipped": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with user conflicts skipped by this job.",
							MarkdownDescription: "The number of files with user conflicts skipped by this job.",
						},
						// 	// "service_report": schema.ListNestedAttribute{
						// 	// 	Computed:            true,
						// 	// 	Description:         "Data for each component exported as part of service replication.",
						// 	// 	MarkdownDescription: "Data for each component exported as part of service replication.",
						// 	// 	NestedObject: schema.NestedAttributeObject{
						// 	// 	Attributes: map[string]schema.Attribute{
						// 	// 		"handlers_transferred": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The number of handlers exported.",
						// 	// 			MarkdownDescription: "The number of handlers exported.",
						// 	// 		},
						// 	// 		"start_time": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The time the job began this component.",
						// 	// 			MarkdownDescription: "The time the job began this component.",
						// 	// 		},
						// 	// 		"component": schema.StringAttribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The component that was processed.",
						// 	// 			MarkdownDescription: "The component that was processed.",
						// 	// 		},
						// 	// 		"directory": schema.StringAttribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The directory of the service export.",
						// 	// 			MarkdownDescription: "The directory of the service export.",
						// 	// 		},
						// 	// 		"handlers_skipped": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The number of handlers skipped during export.",
						// 	// 			MarkdownDescription: "The number of handlers skipped during export.",
						// 	// 		},
						// 	// 		"filter": schema.ListAttribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "A list of path-based filters for exporting components.",
						// 	// 			MarkdownDescription: "A list of path-based filters for exporting components.",
						// 	// 			ElementType: types.StringType,
						// 	// 		},
						// 	// 		"handlers_failed": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The number of handlers failed during export.",
						// 	// 			MarkdownDescription: "The number of handlers failed during export.",
						// 	// 		},
						// 	// 		"records_failed": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The number of records failed during export.",
						// 	// 			MarkdownDescription: "The number of records failed during export.",
						// 	// 		},
						// 	// 		"end_time": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The time the job ended this component.",
						// 	// 			MarkdownDescription: "The time the job ended this component.",
						// 	// 		},
						// 	// 		"status": schema.StringAttribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The current status of export for this component.",
						// 	// 			MarkdownDescription: "The current status of export for this component.",
						// 	// 		},
						// 	// 		"error_msg": schema.ListAttribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "A list of error messages generated while exporting components.",
						// 	// 			MarkdownDescription: "A list of error messages generated while exporting components.",
						// 	// 			ElementType:         types.StringType,
						// 	// 		},
						// 	// 		"records_skipped": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The number of records skipped during export.",
						// 	// 			MarkdownDescription: "The number of records skipped during export.",
						// 	// 		},
						// 	// 		"records_transferred": schema.Int64Attribute{
						// 	// 			Computed:            true,
						// 	// 			Description:         "The number of records exported.",
						// 	// 			MarkdownDescription: "The number of records exported.",
						// 	// 		},
						// 	// 	},
						// 	// },
						// 	// },
						"target_directories_created": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories created on the target.",
							MarkdownDescription: "The number of directories created on the target.",
						},
						"file_data_bytes": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes transferred that belong to files.",
							MarkdownDescription: "The number of bytes transferred that belong to files.",
						},
						"files_transferred": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files transferred by this job.",
							MarkdownDescription: "The number of files transferred by this job.",
						},
						"hash_exceptions_found": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hash exceptions found by this job.",
							MarkdownDescription: "The number of hash exceptions found by this job.",
						},
						"resynced_lins": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs resynched by this job.",
							MarkdownDescription: "The number of LINs resynched by this job.",
						},
						"ads_streams_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of ads streams replicated by this job.",
							MarkdownDescription: "The number of ads streams replicated by this job.",
						},
						"network_bytes_to_target": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of bytes sent to the target by this job.",
							MarkdownDescription: "The total number of bytes sent to the target by this job.",
						},
						"retransmitted_files": schema.ListAttribute{
							Computed:            true,
							Description:         "The files that have been retransmitted by this job.",
							MarkdownDescription: "The files that have been retransmitted by this job.",
							ElementType:         types.StringType,
						},
						"policy_action": schema.StringAttribute{
							Computed:            true,
							Description:         "This is the action the policy is configured to perform.",
							MarkdownDescription: "This is the action the policy is configured to perform.",
						},
						"dirs_deleted": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories deleted by this job.",
							MarkdownDescription: "The number of directories deleted by this job.",
						},
						"fifos_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of fifos replicated by this job.",
							MarkdownDescription: "The number of fifos replicated by this job.",
						},
						"total_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "The total number of data chunks transmitted by this job.",
							MarkdownDescription: "The total number of data chunks transmitted by this job.",
						},
						"flipped_lins": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of LINs flipped by this job.",
							MarkdownDescription: "The number of LINs flipped by this job.",
						},
						"sockets_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of sockets replicated by this job.",
							MarkdownDescription: "The number of sockets replicated by this job.",
						},
						"dead_node": schema.BoolAttribute{
							Computed:            true,
							Description:         "This field is true if the node running this job is dead.",
							MarkdownDescription: "This field is true if the node running this job is dead.",
						},
						"files_selected": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files selected by this job.",
							MarkdownDescription: "The number of files selected by this job.",
						},
						"source_directories_created": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories created on the source.",
							MarkdownDescription: "The number of directories created on the source.",
						},
						"bytes_transferred": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of bytes that have been transferred by this job.",
							MarkdownDescription: "The number of bytes that have been transferred by this job.",
						},
						"succeeded_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of data chunks that have been transmitted successfully.",
							MarkdownDescription: "The number of data chunks that have been transmitted successfully.",
						},
						"retry": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of times the job has been retried.",
							MarkdownDescription: "The number of times the job has been retried.",
						},
						"updated_files_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of updated files replicated by this job.",
							MarkdownDescription: "The number of updated files replicated by this job.",
						},
						"source_directories_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories linked on the source.",
							MarkdownDescription: "The number of directories linked on the source.",
						},
						"committed_files": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of WORM committed files.",
							MarkdownDescription: "The number of WORM committed files.",
						},
						"errors": schema.ListAttribute{
							Computed:            true,
							Description:         "A list of error messages for this job.",
							MarkdownDescription: "A list of error messages for this job.",
							ElementType:         types.StringType,
						},
						"target_snapshots": schema.ListAttribute{
							Computed:            true,
							Description:         "The target snapshots created by this job.",
							MarkdownDescription: "The target snapshots created by this job.",
							ElementType:         types.StringType,
						},
						"source_directories_visited": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories visited on the source.",
							MarkdownDescription: "The number of directories visited on the source.",
						},
						"regular_files_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of regular files replicated by this job.",
							MarkdownDescription: "The number of regular files replicated by this job.",
						},
						"running_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of data chunks currently being transmitted.",
							MarkdownDescription: "The number of data chunks currently being transmitted.",
						},
						"start_time": schema.Int64Attribute{
							Computed:            true,
							Description:         "The time the job started in unix epoch seconds. The field is null if the job hasn't started.",
							MarkdownDescription: "The time the job started in unix epoch seconds. The field is null if the job hasn't started.",
						},
						"duration": schema.Int64Attribute{
							Computed:            true,
							Description:         "The amount of time in seconds between when the job was started and when it ended.  If the job has not yet ended, this is the amount of time since the job started.  This field is null if the job has not yet started.",
							MarkdownDescription: "The amount of time in seconds between when the job was started and when it ended.  If the job has not yet ended, this is the amount of time since the job started.  This field is null if the job has not yet started.",
						},
						"files_with_ads_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files with ads replicated by this job.",
							MarkdownDescription: "The number of files with ads replicated by this job.",
						},
						"files_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files unlinked by this job.",
							MarkdownDescription: "The number of files unlinked by this job.",
						},
						"hard_links_replicated": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of hard links replicated by this job.",
							MarkdownDescription: "The number of hard links replicated by this job.",
						},
						"num_retransmitted_files": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files that have been retransmitted by this job.",
							MarkdownDescription: "The number of files that have been retransmitted by this job.",
						},
						"target_directories_unlinked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of directories unlinked on the target.",
							MarkdownDescription: "The number of directories unlinked on the target.",
						},
						"target_files_linked": schema.Int64Attribute{
							Computed:            true,
							Description:         "The number of files linked on the target.",
							MarkdownDescription: "The number of files linked on the target.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "The state of the job.",
							MarkdownDescription: "The state of the job.",
						},
						"sync_type": schema.StringAttribute{
							Computed:            true,
							Description:         "The type of sync being performed by this job.",
							MarkdownDescription: "The type of sync being performed by this job.",
						},
						"failed_chunks": schema.Int64Attribute{
							Computed:            true,
							Description:         "Tyhe number of data chunks that failed transmission.",
							MarkdownDescription: "Tyhe number of data chunks that failed transmission.",
						},
						"policy": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "The policy associated with this job, or null if there is currently no policy associated with this job (this can happen if the job is newly created and not yet fully populated in the underlying database).",
							MarkdownDescription: "The policy associated with this job, or null if there is currently no policy associated with this job (this can happen if the job is newly created and not yet fully populated in the underlying database).",
							Attributes: map[string]schema.Attribute{
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
																"attribute_exists": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "For \"custom_attribute\" type criteria.  The file will match as long as the attribute named by \"field\" exists.  Default is true.",
																	MarkdownDescription: "For \"custom_attribute\" type criteria.  The file will match as long as the attribute named by \"field\" exists.  Default is true.",
																},
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
																"whole_word": schema.BoolAttribute{
																	Computed:            true,
																	Description:         "If true, the attribute must match the entire word.  Default is true.",
																	MarkdownDescription: "If true, the attribute must match the entire word.  Default is true.",
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
															},
														},
													},
												},
											},
										},
									},
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
								"action": schema.StringAttribute{
									Computed:            true,
									Description:         "The action to be taken by the job.",
									MarkdownDescription: "The action to be taken by the job.",
								},
								"target_path": schema.StringAttribute{
									Computed:            true,
									Description:         "Absolute filesystem path on the target cluster for the sync destination.",
									MarkdownDescription: "Absolute filesystem path on the target cluster for the sync destination.",
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
							},
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
					"resume": schema.StringAttribute{
						Optional:            true,
						Description:         "Continue returning results from previous call using this token (token should come from the previous call, resume cannot be used with other options).",
						MarkdownDescription: "Continue returning results from previous call using this token (token should come from the previous call, resume cannot be used with other options).",
					},
					"newer_than": schema.Int64Attribute{
						Optional:            true,
						Description:         "Filter the returned reports to include only those whose jobs started more recently than the specified number of days ago.",
						MarkdownDescription: "Filter the returned reports to include only those whose jobs started more recently than the specified number of days ago.",
					},
					"policy_name": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter the returned reports to include only those with this policy name.",
						MarkdownDescription: "Filter the returned reports to include only those with this policy name.",
					},
					"state": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter the returned reports to include only those whose jobs are in this state.",
						MarkdownDescription: "Filter the returned reports to include only those whose jobs are in this state.",
					},
					"limit": schema.Int64Attribute{
						Optional:            true,
						Description:         "Return no more than this many results at once (see resume).",
						MarkdownDescription: "Return no more than this many results at once (see resume).",
					},
					"reports_per_policy": schema.Int64Attribute{
						Optional:            true,
						Description:         "If specified, only the N most recent reports will be returned per policy.  If no other query args are present this argument defaults to 1. ",
						MarkdownDescription: "If specified, only the N most recent reports will be returned per policy.  If no other query args are present this argument defaults to 1. ",
					},
					"summary": schema.BoolAttribute{
						Optional:            true,
						Description:         "Return a summary rather than entire objects",
						MarkdownDescription: "Return a summary rather than entire objects",
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
func (d *ReplicationReportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ReplicationReportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading replication report data source")

	var state models.ReplicationReportsDatasourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	replicationReportList, err := helper.GetReplicationReports(ctx, d.client, state)

	if err != nil {
		errStr := constants.ReadReplicationReportsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of replication reports",
			message,
		)
		return
	}
	var rr []models.ReplicationReportsDetail
	for _, rrItem := range replicationReportList.Reports {
		entity := models.ReplicationReportsDetail{}
		err := helper.CopyFields(ctx, rrItem, &entity)
		if err != nil {
			resp.Diagnostics.AddError("Error reading replication report datasource plan",
				fmt.Sprintf("Could not list replication report with error: %s", err.Error()))
			return
		}
		rr = append(rr, entity)
	}
	state.Reports = rr
	state.ID = types.StringValue("replication_report_datasource")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading replication report data source ")
}
