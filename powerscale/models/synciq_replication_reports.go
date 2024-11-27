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

import "github.com/hashicorp/terraform-plugin-framework/types"

// ReplicationReportsDatasourceModel describes the struct for replication reports datasource model
type ReplicationReportsDatasourceModel struct {
	ID                      types.String                 `tfsdk:"id"`
	Reports                 []ReplicationReportsDetail   `tfsdk:"replication_reports"`
	ReplicationReportFilter []ReplicationReportFilterType `tfsdk:"filter"`
}

type ReplicationReportFilterType struct {
	Name types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
// ReplicationReportFilterType describes the struct for filter block
// type ReplicationReportFilterType struct {
// 	Sort             types.String `tfsdk:"sort"`
// 	NewerThan        types.Int64  `tfsdk:"newer_than"`
// 	PolicyName       types.String `tfsdk:"policy_name"`
// 	State            types.String `tfsdk:"state"`
// 	Limit            types.Int64  `tfsdk:"limit"`
// 	ReportsPerPolicy types.Int64  `tfsdk:"reports_per_policy"`
// 	Summary          types.Bool   `tfsdk:"summary"`
// 	Dir              types.String `tfsdk:"dir"`
// }

// ReplicationReportsDetail describes the struct for replication report
type ReplicationReportsDetail struct {
	Action                     types.String   `tfsdk:"action"`
	AdsStreamsReplicated       types.Int64    `tfsdk:"ads_streams_replicated"`
	BlockSpecsReplicated       types.Int64    `tfsdk:"block_specs_replicated"`
	BytesRecoverable           types.Int64    `tfsdk:"bytes_recoverable"`
	BytesTransferred           types.Int64    `tfsdk:"bytes_transferred"`
	CharSpecsReplicated        types.Int64    `tfsdk:"char_specs_replicated"`
	CommittedFiles             types.Int64    `tfsdk:"committed_files"`
	CorrectedLins              types.Int64    `tfsdk:"corrected_lins"`
	DeadNode                   types.Bool     `tfsdk:"dead_node"`
	DirectoriesReplicated      types.Int64    `tfsdk:"directories_replicated"`
	DirsChanged                types.Int64    `tfsdk:"dirs_changed"`
	DirsDeleted                types.Int64    `tfsdk:"dirs_deleted"`
	DirsMoved                  types.Int64    `tfsdk:"dirs_moved"`
	DirsNew                    types.Int64    `tfsdk:"dirs_new"`
	Duration                   types.Int64    `tfsdk:"duration"`
	Encrypted                  types.Bool     `tfsdk:"encrypted"`
	EndTime                    types.Int64    `tfsdk:"end_time"`
	Error                      types.String   `tfsdk:"error"`
	ErrorChecksumFilesSkipped  types.Int64    `tfsdk:"error_checksum_files_skipped"`
	ErrorIoFilesSkipped        types.Int64    `tfsdk:"error_io_files_skipped"`
	ErrorNetFilesSkipped       types.Int64    `tfsdk:"error_net_files_skipped"`
	Errors                     types.List     `tfsdk:"errors"`
	FailedChunks               types.Int64    `tfsdk:"failed_chunks"`
	FifosReplicated            types.Int64    `tfsdk:"fifos_replicated"`
	FileDataBytes              types.Int64    `tfsdk:"file_data_bytes"`
	FilesChanged               types.Int64    `tfsdk:"files_changed"`
	FilesLinked                types.Int64    `tfsdk:"files_linked"`
	FilesNew                   types.Int64    `tfsdk:"files_new"`
	FilesSelected              types.Int64    `tfsdk:"files_selected"`
	FilesTransferred           types.Int64    `tfsdk:"files_transferred"`
	FilesUnlinked              types.Int64    `tfsdk:"files_unlinked"`
	FilesWithAdsReplicated     types.Int64    `tfsdk:"files_with_ads_replicated"`
	FlippedLins                types.Int64    `tfsdk:"flipped_lins"`
	HardLinksReplicated        types.Int64    `tfsdk:"hard_links_replicated"`
	HashExceptionsFixed        types.Int64    `tfsdk:"hash_exceptions_fixed"`
	HashExceptionsFound        types.Int64    `tfsdk:"hash_exceptions_found"`
	ID                         types.String   `tfsdk:"id"`
	JobID                      types.Int64    `tfsdk:"job_id"`
	LinsTotal                  types.Int64    `tfsdk:"lins_total"`
	NetworkBytesToSource       types.Int64    `tfsdk:"network_bytes_to_source"`
	NetworkBytesToTarget       types.Int64    `tfsdk:"network_bytes_to_target"`
	NewFilesReplicated         types.Int64    `tfsdk:"new_files_replicated"`
	NumRetransmittedFiles      types.Int64    `tfsdk:"num_retransmitted_files"`
	Phases                     []PhasesDetail `tfsdk:"phases"`
	Policy                     PolicyDetail   `tfsdk:"policy"`
	PolicyAction               types.String   `tfsdk:"policy_action"`
	PolicyID                   types.String   `tfsdk:"policy_id"`
	PolicyName                 types.String   `tfsdk:"policy_name"`
	QuotasDeleted              types.Int64    `tfsdk:"quotas_deleted"`
	RegularFilesReplicated     types.Int64    `tfsdk:"regular_files_replicated"`
	ResyncedLins               types.Int64    `tfsdk:"resynced_lins"`
	RetransmittedFiles         types.List     `tfsdk:"retransmitted_files"`
	Retry                      types.Int64    `tfsdk:"retry"`
	RunningChunks              types.Int64    `tfsdk:"running_chunks"`
	SocketsReplicated          types.Int64    `tfsdk:"sockets_replicated"`
	SourceBytesRecovered       types.Int64    `tfsdk:"source_bytes_recovered"`
	SourceDirectoriesCreated   types.Int64    `tfsdk:"source_directories_created"`
	SourceDirectoriesDeleted   types.Int64    `tfsdk:"source_directories_deleted"`
	SourceDirectoriesLinked    types.Int64    `tfsdk:"source_directories_linked"`
	SourceDirectoriesUnlinked  types.Int64    `tfsdk:"source_directories_unlinked"`
	SourceDirectoriesVisited   types.Int64    `tfsdk:"source_directories_visited"`
	SourceFilesDeleted         types.Int64    `tfsdk:"source_files_deleted"`
	SourceFilesLinked          types.Int64    `tfsdk:"source_files_linked"`
	SourceFilesUnlinked        types.Int64    `tfsdk:"source_files_unlinked"`
	SparseDataBytes            types.Int64    `tfsdk:"sparse_data_bytes"`
	StartTime                  types.Int64    `tfsdk:"start_time"`
	State                      types.String   `tfsdk:"state"`
	SubreportCount             types.Int64    `tfsdk:"subreport_count"`
	SucceededChunks            types.Int64    `tfsdk:"succeeded_chunks"`
	SymlinksReplicated         types.Int64    `tfsdk:"symlinks_replicated"`
	SyncType                   types.String   `tfsdk:"sync_type"`
	TargetBytesRecovered       types.Int64    `tfsdk:"target_bytes_recovered"`
	TargetDirectoriesCreated   types.Int64    `tfsdk:"target_directories_created"`
	TargetDirectoriesDeleted   types.Int64    `tfsdk:"target_directories_deleted"`
	TargetDirectoriesLinked    types.Int64    `tfsdk:"target_directories_linked"`
	TargetDirectoriesUnlinked  types.Int64    `tfsdk:"target_directories_unlinked"`
	TargetFilesDeleted         types.Int64    `tfsdk:"target_files_deleted"`
	TargetFilesLinked          types.Int64    `tfsdk:"target_files_linked"`
	TargetFilesUnlinked        types.Int64    `tfsdk:"target_files_unlinked"`
	TargetSnapshots            types.List     `tfsdk:"target_snapshots"`
	Throughput                 types.String   `tfsdk:"throughput"`
	TotalChunks                types.Int64    `tfsdk:"total_chunks"`
	TotalDataBytes             types.Int64    `tfsdk:"total_data_bytes"`
	TotalFiles                 types.Int64    `tfsdk:"total_files"`
	TotalNetworkBytes          types.Int64    `tfsdk:"total_network_bytes"`
	TotalPhases                types.Int64    `tfsdk:"total_phases"`
	UnchangedDataBytes         types.Int64    `tfsdk:"unchanged_data_bytes"`
	UpToDateFilesSkipped       types.Int64    `tfsdk:"up_to_date_files_skipped"`
	UpdatedFilesReplicated     types.Int64    `tfsdk:"updated_files_replicated"`
	UserConflictFilesSkipped   types.Int64    `tfsdk:"user_conflict_files_skipped"`
	Warnings                   types.List     `tfsdk:"warnings"`
	WormCommittedFileConflicts types.Int64    `tfsdk:"worm_committed_file_conflicts"`
}

// PolicyDetail describes the struct for policy detail
type PolicyDetail struct {
	Action                   types.String              `tfsdk:"action"`
	FileMatchingPattern      FileMatchingPatternDetail `tfsdk:"file_matching_pattern"`
	Name                     types.String              `tfsdk:"name"`
	SourceExcludeDirectories types.List                `tfsdk:"source_exclude_directories"`
	SourceIncludeDirectories types.List                `tfsdk:"source_include_directories"`
	SourceRootPath           types.String              `tfsdk:"source_root_path"`
	TargetHost               types.String              `tfsdk:"target_host"`
	TargetPath               types.String              `tfsdk:"target_path"`
}

// OrCriteriaDetail describes the struct for 'or criteria' in the policy
type OrCriteriaDetail struct {
	AndCriteria []AndCriteriaDetail `tfsdk:"and_criteria"`
}

// AndCriteriaDetail describes the struct for 'and criteria' in the policy
type AndCriteriaDetail struct {
	AttributeExists types.Bool   `tfsdk:"attribute_exists"`
	CaseSensitive   types.Bool   `tfsdk:"case_sensitive"`
	Field           types.String `tfsdk:"field"`
	Operator        types.String `tfsdk:"operator"`
	Type            types.String `tfsdk:"type"`
	Value           types.String `tfsdk:"value"`
	WholeWord       types.Bool   `tfsdk:"whole_word"`
}

// StatisticsDetail describes the struct for statistics
type StatisticsDetail struct {
	ComplianceDirLinks types.String `tfsdk:"compliance_dir_links"`
	CorrectedLins      types.String `tfsdk:"corrected_lins"`
	DeletedDirs        types.String `tfsdk:"deleted_dirs"`
	Dirs               types.String `tfsdk:"dirs"`
	Files              types.String `tfsdk:"files"`
	FlippedLins        types.String `tfsdk:"flipped_lins"`
	HashExceptions     types.String `tfsdk:"hash_exceptions"`
	LinkedDirs         types.String `tfsdk:"linked_dirs"`
	LinkedFiles        types.String `tfsdk:"linked_files"`
	MarkedDirectories  types.String `tfsdk:"marked_directories"`
	MarkedFiles        types.String `tfsdk:"marked_files"`
	ModifiedDirs       types.String `tfsdk:"modified_dirs"`
	ModifiedFiles      types.String `tfsdk:"modified_files"`
	ModifiedLins       types.String `tfsdk:"modified_lins"`
	NewComplianceDirs  types.String `tfsdk:"new_compliance_dirs"`
	NewDirs            types.String `tfsdk:"new_dirs"`
	NewFiles           types.String `tfsdk:"new_files"`
	NewResyncedFiles   types.String `tfsdk:"new_resynced_files"`
	ResyncedFileLinks  types.String `tfsdk:"resynced_file_links"`
	ResyncedLins       types.String `tfsdk:"resynced_lins"`
	UnlinkedFiles      types.String `tfsdk:"unlinked_files"`
}

// PhasesDetail describes the struct for phase
type PhasesDetail struct {
	EndTime    types.Int64      `tfsdk:"end_time"`
	Phase      types.String     `tfsdk:"phase"`
	StartTime  types.Int64      `tfsdk:"start_time"`
	Statistics StatisticsDetail `tfsdk:"statistics"`
}

// FileMatchingPatternDetail describes the struct for file matching pattern
type FileMatchingPatternDetail struct {
	OrCriteria []OrCriteriaDetail `tfsdk:"or_criteria"`
}
