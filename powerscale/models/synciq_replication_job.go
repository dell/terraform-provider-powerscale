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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SyncIQReplicationJobResourceModel describes the SyncIQ Replication Job resource data model.
type SyncIQReplicationJobResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Action   types.String `tfsdk:"action"`
	IsPaused types.Bool   `tfsdk:"is_paused"`
	WaitTime types.Int64  `tfsdk:"wait_time"`
}

// SyncIQReplicationJobDataSourceModel describes the SyncIQ Replication Job datasource data model.
type SyncIQReplicationJobDataSourceModel struct {
	SyncIQReplicationJobs []SyncIQReplicationJobModel `tfsdk:"synciq_jobs"`
	ID                    types.String                `tfsdk:"id"`
	SyncIQJobFilter       *SyncIQJobFilterModel       `tfsdk:"filter"`
}

// SyncIQJobFilterModel describes the filters supported by api.
type SyncIQJobFilterModel struct {
	// The field that will be used for sorting.
	Sort types.String `tfsdk:"sort"`
	// Return no more than this many results at once.
	Limit types.Int32 `tfsdk:"limit"`
	// The direction of the sort. Supported Values: ASC, DESC.
	Dir types.String `tfsdk:"dir"`
	// Filter on the state of the SyncIQ job.
	State types.String `tfsdk:"state"`
}

// SyncIQReplicationJobModel defines the model of a SyncIQ job.
type SyncIQReplicationJobModel struct {
	CharSpecsReplicated        types.Int64  `tfsdk:"char_specs_replicated"`
	TargetBytesRecovered       types.Int64  `tfsdk:"target_bytes_recovered"`
	DirsNew                    types.Int64  `tfsdk:"dirs_new"`
	PolicyID                   types.String `tfsdk:"policy_id"`
	UpdatedFilesReplicated     types.Int64  `tfsdk:"updated_files_replicated"`
	FilesWithAdsReplicated     types.Int64  `tfsdk:"files_with_ads_replicated"`
	EndTime                    types.Int64  `tfsdk:"end_time"`
	BytesRecoverable           types.Int64  `tfsdk:"bytes_recoverable"`
	ID                         types.String `tfsdk:"id"`
	UnchangedDataBytes         types.Int64  `tfsdk:"unchanged_data_bytes"`
	SourceBytesRecovered       types.Int64  `tfsdk:"source_bytes_recovered"`
	UpToDateFilesSkipped       types.Int64  `tfsdk:"up_to_date_files_skipped"`
	SourceFilesUnlinked        types.Int64  `tfsdk:"source_files_unlinked"`
	FifosReplicated            types.Int64  `tfsdk:"fifos_replicated"`
	Error                      types.String `tfsdk:"error"`
	JobID                      types.Int64  `tfsdk:"job_id"`
	FilesLinked                types.Int64  `tfsdk:"files_linked"`
	SourceDirectoriesLinked    types.Int64  `tfsdk:"source_directories_linked"`
	RetransmittedFiles         types.List   `tfsdk:"retransmitted_files"`
	DirsDeleted                types.Int64  `tfsdk:"dirs_deleted"`
	PolicyAction               types.String `tfsdk:"policy_action"`
	FlippedLins                types.Int64  `tfsdk:"flipped_lins"`
	QuotasDeleted              types.Int64  `tfsdk:"quotas_deleted"`
	TargetFilesDeleted         types.Int64  `tfsdk:"target_files_deleted"`
	Policy                     types.Object `tfsdk:"policy"`
	LinsTotal                  types.Int64  `tfsdk:"lins_total"`
	TotalDataBytes             types.Int64  `tfsdk:"total_data_bytes"`
	FileDataBytes              types.Int64  `tfsdk:"file_data_bytes"`
	TargetFilesLinked          types.Int64  `tfsdk:"target_files_linked"`
	TargetDirectoriesLinked    types.Int64  `tfsdk:"target_directories_linked"`
	Workers                    types.List   `tfsdk:"workers"`
	TotalFiles                 types.Int64  `tfsdk:"total_files"`
	CommittedFiles             types.Int64  `tfsdk:"committed_files"`
	Retry                      types.Int64  `tfsdk:"retry"`
	UserConflictFilesSkipped   types.Int64  `tfsdk:"user_conflict_files_skipped"`
	FilesNew                   types.Int64  `tfsdk:"files_new"`
	RunningChunks              types.Int64  `tfsdk:"running_chunks"`
	RegularFilesReplicated     types.Int64  `tfsdk:"regular_files_replicated"`
	State                      types.String `tfsdk:"state"`
	Action                     types.String `tfsdk:"action"`
	TotalPhases                types.Int64  `tfsdk:"total_phases"`
	NumRetransmittedFiles      types.Int64  `tfsdk:"num_retransmitted_files"`
	HardLinksReplicated        types.Int64  `tfsdk:"hard_links_replicated"`
	TargetSnapshots            types.List   `tfsdk:"target_snapshots"`
	DirsMoved                  types.Int64  `tfsdk:"dirs_moved"`
	TotalNetworkBytes          types.Int64  `tfsdk:"total_network_bytes"`
	SparseDataBytes            types.Int64  `tfsdk:"sparse_data_bytes"`
	Phases                     types.List   `tfsdk:"phases"`
	SocketsReplicated          types.Int64  `tfsdk:"sockets_replicated"`
	SourceDirectoriesCreated   types.Int64  `tfsdk:"source_directories_created"`
	ResyncedLins               types.Int64  `tfsdk:"resynced_lins"`
	SourceDirectoriesDeleted   types.Int64  `tfsdk:"source_directories_deleted"`
	TargetFilesUnlinked        types.Int64  `tfsdk:"target_files_unlinked"`
	NewFilesReplicated         types.Int64  `tfsdk:"new_files_replicated"`
	Duration                   types.Int64  `tfsdk:"duration"`
	PolicyName                 types.String `tfsdk:"policy_name"`
	FilesTransferred           types.Int64  `tfsdk:"files_transferred"`
	TargetDirectoriesCreated   types.Int64  `tfsdk:"target_directories_created"`
	DirectoriesReplicated      types.Int64  `tfsdk:"directories_replicated"`
	ServiceReport              types.List   `tfsdk:"service_report"`
	SucceededChunks            types.Int64  `tfsdk:"succeeded_chunks"`
	CorrectedLins              types.Int64  `tfsdk:"corrected_lins"`
	TargetDirectoriesDeleted   types.Int64  `tfsdk:"target_directories_deleted"`
	ErrorChecksumFilesSkipped  types.Int64  `tfsdk:"error_checksum_files_skipped"`
	FailedChunks               types.Int64  `tfsdk:"failed_chunks"`
	TargetDirectoriesUnlinked  types.Int64  `tfsdk:"target_directories_unlinked"`
	NetworkBytesToTarget       types.Int64  `tfsdk:"network_bytes_to_target"`
	FilesSelected              types.Int64  `tfsdk:"files_selected"`
	SourceFilesDeleted         types.Int64  `tfsdk:"source_files_deleted"`
	FilesChanged               types.Int64  `tfsdk:"files_changed"`
	TotalExportedServices      types.Int64  `tfsdk:"total_exported_services"`
	Encrypted                  types.Bool   `tfsdk:"encrypted"`
	BlockSpecsReplicated       types.Int64  `tfsdk:"block_specs_replicated"`
	StartTime                  types.Int64  `tfsdk:"start_time"`
	WormCommittedFileConflicts types.Int64  `tfsdk:"worm_committed_file_conflicts"`
	FilesUnlinked              types.Int64  `tfsdk:"files_unlinked"`
	ErrorIoFilesSkipped        types.Int64  `tfsdk:"error_io_files_skipped"`
	TotalChunks                types.Int64  `tfsdk:"total_chunks"`
	ErrorNetFilesSkipped       types.Int64  `tfsdk:"error_net_files_skipped"`
	SourceDirectoriesVisited   types.Int64  `tfsdk:"source_directories_visited"`
	HashExceptionsFixed        types.Int64  `tfsdk:"hash_exceptions_fixed"`
	SourceDirectoriesUnlinked  types.Int64  `tfsdk:"source_directories_unlinked"`
	AdsStreamsReplicated       types.Int64  `tfsdk:"ads_streams_replicated"`
	SymlinksReplicated         types.Int64  `tfsdk:"symlinks_replicated"`
	DeadNode                   types.Bool   `tfsdk:"dead_node"`
	SyncType                   types.String `tfsdk:"sync_type"`
	HashExceptionsFound        types.Int64  `tfsdk:"hash_exceptions_found"`
	Warnings                   types.List   `tfsdk:"warnings"`
	BytesTransferred           types.Int64  `tfsdk:"bytes_transferred"`
	Errors                     types.List   `tfsdk:"errors"`
	SourceFilesLinked          types.Int64  `tfsdk:"source_files_linked"`
	NetworkBytesToSource       types.Int64  `tfsdk:"network_bytes_to_source"`
	DirsChanged                types.Int64  `tfsdk:"dirs_changed"`
}
