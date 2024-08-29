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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SyncIQPolicyDataSource returns the overall cluster details.
type SyncIQPolicyDataSource struct {
	ID       types.String                 `tfsdk:"id"`
	Policies []V14SyncPolicyExtendedModel `tfsdk:"policies"`
}

// V14SyncPolicyExtendedModel is the tfsdk model of V14SyncPolicyExtended.
type V14SyncPolicyExtendedModel struct {
	AcceleratedFailback               types.Bool                              `tfsdk:"accelerated_failback"`
	Action                            types.String                            `tfsdk:"action"`
	AllowCopyFb                       types.Bool                              `tfsdk:"allow_copy_fb"`
	BandwidthReservation              types.Int64                             `tfsdk:"bandwidth_reservation"`
	Changelist                        types.Bool                              `tfsdk:"changelist"`
	CheckIntegrity                    types.Bool                              `tfsdk:"check_integrity"`
	CloudDeepCopy                     types.String                            `tfsdk:"cloud_deep_copy"`
	Conflicted                        types.Bool                              `tfsdk:"conflicted"`
	DatabaseMirrored                  types.Bool                              `tfsdk:"database_mirrored"`
	DeleteQuotas                      types.Bool                              `tfsdk:"delete_quotas"`
	Description                       types.String                            `tfsdk:"description"`
	DisableFileSplit                  types.Bool                              `tfsdk:"disable_file_split"`
	DisableFofb                       types.Bool                              `tfsdk:"disable_fofb"`
	DisableQuotaTmpDir                types.Bool                              `tfsdk:"disable_quota_tmp_dir"`
	DisableStf                        types.Bool                              `tfsdk:"disable_stf"`
	EnableHashTmpdir                  types.Bool                              `tfsdk:"enable_hash_tmpdir"`
	Enabled                           types.Bool                              `tfsdk:"enabled"`
	Encrypted                         types.Bool                              `tfsdk:"encrypted"`
	EncryptionCipherList              types.String                            `tfsdk:"encryption_cipher_list"`
	ExpectedDataloss                  types.Bool                              `tfsdk:"expected_dataloss"`
	FileMatchingPattern               V1SyncJobPolicyFileMatchingPatternModel `tfsdk:"file_matching_pattern"`
	ForceInterface                    types.Bool                              `tfsdk:"force_interface"`
	HasSyncState                      types.Bool                              `tfsdk:"has_sync_state"`
	ID                                types.String                            `tfsdk:"id"`
	IgnoreRecursiveQuota              types.Bool                              `tfsdk:"ignore_recursive_quota"`
	JobDelay                          types.Int64                             `tfsdk:"job_delay"`
	LastJobState                      types.String                            `tfsdk:"last_job_state"`
	LastStarted                       types.Int64                             `tfsdk:"last_started"`
	LastSuccess                       types.Int64                             `tfsdk:"last_success"`
	LinkedServicePolicies             []types.String                          `tfsdk:"linked_service_policies"`
	LogLevel                          types.String                            `tfsdk:"log_level"`
	LogRemovedFiles                   types.Bool                              `tfsdk:"log_removed_files"`
	Name                              types.String                            `tfsdk:"name"`
	NextRun                           types.Int64                             `tfsdk:"next_run"`
	OcspAddress                       types.String                            `tfsdk:"ocsp_address"`
	OcspIssuerCertificateID           types.String                            `tfsdk:"ocsp_issuer_certificate_id"`
	PasswordSet                       types.Bool                              `tfsdk:"password_set"`
	Priority                          types.Int64                             `tfsdk:"priority"`
	ReportMaxAge                      types.Int64                             `tfsdk:"report_max_age"`
	ReportMaxCount                    types.Int64                             `tfsdk:"report_max_count"`
	RestrictTargetNetwork             types.Bool                              `tfsdk:"restrict_target_network"`
	RpoAlert                          types.Int64                             `tfsdk:"rpo_alert"`
	Schedule                          types.String                            `tfsdk:"schedule"`
	ServicePolicy                     types.Bool                              `tfsdk:"service_policy"`
	SkipLookup                        types.Bool                              `tfsdk:"skip_lookup"`
	SkipWhenSourceUnmodified          types.Bool                              `tfsdk:"skip_when_source_unmodified"`
	SnapshotSyncExisting              types.Bool                              `tfsdk:"snapshot_sync_existing"`
	SnapshotSyncPattern               types.String                            `tfsdk:"snapshot_sync_pattern"`
	SourceCertificateID               types.String                            `tfsdk:"source_certificate_id"`
	SourceDomainMarked                types.Bool                              `tfsdk:"source_domain_marked"`
	SourceExcludeDirectories          []types.String                          `tfsdk:"source_exclude_directories"`
	SourceIncludeDirectories          []types.String                          `tfsdk:"source_include_directories"`
	SourceNetwork                     V1SyncPolicySourceNetworkModel          `tfsdk:"source_network"`
	SourceRootPath                    types.String                            `tfsdk:"source_root_path"`
	SourceSnapshotArchive             types.Bool                              `tfsdk:"source_snapshot_archive"`
	SourceSnapshotExpiration          types.Int64                             `tfsdk:"source_snapshot_expiration"`
	SourceSnapshotPattern             types.String                            `tfsdk:"source_snapshot_pattern"`
	SyncExistingSnapshotExpiration    types.Bool                              `tfsdk:"sync_existing_snapshot_expiration"`
	SyncExistingTargetSnapshotPattern types.String                            `tfsdk:"sync_existing_target_snapshot_pattern"`
	TargetCertificateID               types.String                            `tfsdk:"target_certificate_id"`
	TargetCompareInitialSync          types.Bool                              `tfsdk:"target_compare_initial_sync"`
	TargetDetectModifications         types.Bool                              `tfsdk:"target_detect_modifications"`
	TargetHost                        types.String                            `tfsdk:"target_host"`
	TargetPath                        types.String                            `tfsdk:"target_path"`
	TargetSnapshotAlias               types.String                            `tfsdk:"target_snapshot_alias"`
	TargetSnapshotArchive             types.Bool                              `tfsdk:"target_snapshot_archive"`
	TargetSnapshotExpiration          types.Int64                             `tfsdk:"target_snapshot_expiration"`
	TargetSnapshotPattern             types.String                            `tfsdk:"target_snapshot_pattern"`
	WorkersPerNode                    types.Int64                             `tfsdk:"workers_per_node"`
}

// V1SyncJobPolicyFileMatchingPatternModel is the tfsdk model of V1SyncJobPolicyFileMatchingPattern.
type V1SyncJobPolicyFileMatchingPatternModel struct {
	OrCriteria []V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModel `tfsdk:"or_criteria"`
}

// V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModel is the tfsdk model of V1SyncJobPolicyFileMatchingPatternOrCriteriaItem.
type V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModel struct {
	AndCriteria []V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModel `tfsdk:"and_criteria"`
}

// V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModel is the tfsdk model of V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem.
type V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModel struct {
	AttributeExists types.Bool   `tfsdk:"attribute_exists"`
	CaseSensitive   types.Bool   `tfsdk:"case_sensitive"`
	Field           types.String `tfsdk:"field"`
	Operator        types.String `tfsdk:"operator"`
	Type            types.String `tfsdk:"type"`
	Value           types.String `tfsdk:"value"`
	WholeWord       types.Bool   `tfsdk:"whole_word"`
}

// V1SyncPolicySourceNetworkModel is the tfsdk model of V1SyncPolicySourceNetwork.
type V1SyncPolicySourceNetworkModel struct {
	Pool   types.String `tfsdk:"pool"`
	Subnet types.String `tfsdk:"subnet"`
}

// Resource Model

type SynciqpolicyResourceModel struct {
	AcceleratedFailback               types.Bool   `tfsdk:"accelerated_failback"`
	Action                            types.String `tfsdk:"action"`
	AllowCopyFb                       types.Bool   `tfsdk:"allow_copy_fb"`
	BandwidthReservation              types.Int64  `tfsdk:"bandwidth_reservation"`
	Changelist                        types.Bool   `tfsdk:"changelist"`
	CheckIntegrity                    types.Bool   `tfsdk:"check_integrity"`
	CloudDeepCopy                     types.String `tfsdk:"cloud_deep_copy"`
	DeleteQuotas                      types.Bool   `tfsdk:"delete_quotas"`
	Description                       types.String `tfsdk:"description"`
	DisableFileSplit                  types.Bool   `tfsdk:"disable_file_split"`
	DisableFofb                       types.Bool   `tfsdk:"disable_fofb"`
	DisableQuotaTmpDir                types.Bool   `tfsdk:"disable_quota_tmp_dir"`
	DisableStf                        types.Bool   `tfsdk:"disable_stf"`
	EnableHashTmpdir                  types.Bool   `tfsdk:"enable_hash_tmpdir"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
	EncryptionCipherList              types.String `tfsdk:"encryption_cipher_list"`
	ExpectedDataloss                  types.Bool   `tfsdk:"expected_dataloss"`
	FileMatchingPattern               types.Object `tfsdk:"file_matching_pattern"`
	ForceInterface                    types.Bool   `tfsdk:"force_interface"`
	IgnoreRecursiveQuota              types.Bool   `tfsdk:"ignore_recursive_quota"`
	JobDelay                          types.Int64  `tfsdk:"job_delay"`
	LinkedServicePolicies             types.List   `tfsdk:"linked_service_policies"`
	LogLevel                          types.String `tfsdk:"log_level"`
	LogRemovedFiles                   types.Bool   `tfsdk:"log_removed_files"`
	Name                              types.String `tfsdk:"name"`
	OcspAddress                       types.String `tfsdk:"ocsp_address"`
	OcspIssuerCertificateId           types.String `tfsdk:"ocsp_issuer_certificate_id"`
	Password                          types.String `tfsdk:"password"`
	Priority                          types.Int64  `tfsdk:"priority"`
	ReportMaxAge                      types.Int64  `tfsdk:"report_max_age"`
	ReportMaxCount                    types.Int64  `tfsdk:"report_max_count"`
	RestrictTargetNetwork             types.Bool   `tfsdk:"restrict_target_network"`
	RpoAlert                          types.Int64  `tfsdk:"rpo_alert"`
	Schedule                          types.String `tfsdk:"schedule"`
	ServicePolicy                     types.Bool   `tfsdk:"service_policy"`
	SkipLookup                        types.Bool   `tfsdk:"skip_lookup"`
	SkipWhenSourceUnmodified          types.Bool   `tfsdk:"skip_when_source_unmodified"`
	SnapshotSyncExisting              types.Bool   `tfsdk:"snapshot_sync_existing"`
	SnapshotSyncPattern               types.String `tfsdk:"snapshot_sync_pattern"`
	SourceExcludeDirectories          types.List   `tfsdk:"source_exclude_directories"`
	SourceIncludeDirectories          types.List   `tfsdk:"source_include_directories"`
	SourceNetwork                     types.Object `tfsdk:"source_network"`
	SourceRootPath                    types.String `tfsdk:"source_root_path"`
	SourceSnapshotArchive             types.Bool   `tfsdk:"source_snapshot_archive"`
	SourceSnapshotExpiration          types.Int64  `tfsdk:"source_snapshot_expiration"`
	SourceSnapshotPattern             types.String `tfsdk:"source_snapshot_pattern"`
	SyncExistingSnapshotExpiration    types.Bool   `tfsdk:"sync_existing_snapshot_expiration"`
	SyncExistingTargetSnapshotPattern types.String `tfsdk:"sync_existing_target_snapshot_pattern"`
	TargetCertificateId               types.String `tfsdk:"target_certificate_id"`
	TargetCompareInitialSync          types.Bool   `tfsdk:"target_compare_initial_sync"`
	TargetDetectModifications         types.Bool   `tfsdk:"target_detect_modifications"`
	TargetHost                        types.String `tfsdk:"target_host"`
	TargetPath                        types.String `tfsdk:"target_path"`
	TargetSnapshotAlias               types.String `tfsdk:"target_snapshot_alias"`
	TargetSnapshotArchive             types.Bool   `tfsdk:"target_snapshot_archive"`
	TargetSnapshotExpiration          types.Int64  `tfsdk:"target_snapshot_expiration"`
	TargetSnapshotPattern             types.String `tfsdk:"target_snapshot_pattern"`
	WorkersPerNode                    types.Int64  `tfsdk:"workers_per_node"`
	Conflicted                        types.Bool   `tfsdk:"conflicted"`
	DatabaseMirrored                  types.Bool   `tfsdk:"database_mirrored"`
	Encrypted                         types.Bool   `tfsdk:"encrypted"`
	HasSyncState                      types.Bool   `tfsdk:"has_sync_state"`
	Id                                types.String `tfsdk:"id"`
	LastJobState                      types.String `tfsdk:"last_job_state"`
	LastStarted                       types.Int64  `tfsdk:"last_started"`
	LastSuccess                       types.Int64  `tfsdk:"last_success"`
	NextRun                           types.Int64  `tfsdk:"next_run"`
	PasswordSet                       types.Bool   `tfsdk:"password_set"`
	SourceCertificateId               types.String `tfsdk:"source_certificate_id"`
	SourceDomainMarked                types.Bool   `tfsdk:"source_domain_marked"`
}
