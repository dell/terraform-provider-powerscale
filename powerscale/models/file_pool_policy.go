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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// FilePoolPolicyModel describes the resource data model.
type FilePoolPolicyModel struct {
	// Specifies if the policy is default policy. Default policy applies to all files not selected by higher-priority policies.
	IsDefaultPolicy types.Bool `tfsdk:"is_default_policy"`
	// A list of actions to be taken for matching files
	Actions []V1FilepoolDefaultPolicyAction `tfsdk:"actions"`
	// The order in which this policy should be applied (relative to other policies)
	ApplyOrder types.Int64 `tfsdk:"apply_order"`
	// The guid assigned to the cluster on which the policy was created
	BirthClusterID types.String `tfsdk:"birth_cluster_id"`
	// A description for this policy
	Description types.String `tfsdk:"description"`
	// The file matching rules for this policy
	FileMatchingPattern *V1FilepoolPolicyFileMatchingPattern `tfsdk:"file_matching_pattern"`
	// A unique name for this policy
	ID types.String `tfsdk:"id"`
	// A unique name for this policy
	Name types.String `tfsdk:"name"`
	// Indicates whether this policy is in a good state (\"OK\") or disabled (\"disabled\")
	State types.String `tfsdk:"state"`
	// Gives further information to describe the state of this policy
	StateDetails types.String `tfsdk:"state_details"`
}

// V1FilepoolDefaultPolicyAction An action to apply to a file matching the policy
type V1FilepoolDefaultPolicyAction struct {
	// ActionParam *V12ActionParam `tfsdk:"action_param"`
	// action for set_data_access_pattern type
	DataAccessPatternAction types.String `tfsdk:"data_access_pattern_action"`
	// action for apply_data_storage_policy	type
	DataStoragePolicyAction *V12StoragePolicyActionParams `tfsdk:"data_storage_policy_action"`
	// action for apply_snapshot_storage_policy type
	SnapshotStoragePolicyAction *V12StoragePolicyActionParams `tfsdk:"snapshot_storage_policy_action"`
	// action for enable_coalescer type
	EnableCoalescerAction types.Bool `tfsdk:"enable_coalescer_action"`
	// action for enable_packing type
	EnablePackingAction types.Bool `tfsdk:"enable_packing_action"`
	// action for set_requested_protection type
	RequestedProtectionAction types.String `tfsdk:"requested_protection_action"`
	// action for set_cloudpool_policy type
	CloudPoolPolicyAction *V12CloudPolicyActionParams `tfsdk:"cloudpool_policy_action"`
	// action_type Acceptable values: set_requested_protection, set_data_access_pattern, enable_coalescer, apply_data_storage_policy, apply_snapshot_storage_policy, set_cloudpool_policy, enable_packing.
	ActionType types.String `tfsdk:"action_type"`
}

// V12StoragePolicyActionParams - Action for apply_data_storage_policy and apply_snapshot_storage_policy type
type V12StoragePolicyActionParams struct {
	SSDStrategy types.String `tfsdk:"ssd_strategy"`
	StoragePool types.String `tfsdk:"storagepool"`
}

// V12StoragePolicyActionParamsJSONModel - Json model for V12StoragePolicyActionParam
type V12StoragePolicyActionParamsJSONModel struct {
	SSDStrategy string `json:"ssd_strategy"`
	StoragePool string `json:"storagepool"`
}

// V12CloudPolicyArchiveParams - Archive Params for CloudPolicyAction
type V12CloudPolicyArchiveParams struct {
	CloudPolicyActionParams *V12CloudPolicyActionParams `tfsdk:"archive_parameters"`
}

// V12CloudPolicyArchiveParamsJSONModel - Json model for V12CloudPolicyArchiveParams
type V12CloudPolicyArchiveParamsJSONModel struct {
	CloudPolicyActionParams *V12CloudPolicyActionParamsJSONModel `json:"archive_parameters"`
}

// V12CloudPolicyActionParams The filepool policy values for cloudpools.
type V12CloudPolicyActionParams struct {
	// Specifies the cloudPool storage target.
	Pool types.String `tfsdk:"pool"`
	// Specifies if files with snapshots should be archived.
	ArchiveSnapshotFiles types.Bool                       `tfsdk:"archive_snapshot_files"`
	Cache                *V12CloudPolicyActionCacheParams `tfsdk:"cache"`
	// Specifies if files should be compressed.
	Compression types.Bool `tfsdk:"compression"`
	// Specifies the minimum amount of time archived data will be retained in the cloud after deletion.
	DataRetention types.Int64 `tfsdk:"data_retention"`
	// Specifies if files should be encrypted.
	Encryption types.Bool `tfsdk:"encryption"`
	// (Used with NDMP backups only.  Not applicable to SyncIQ.)  The minimum amount of time cloud files will be retained after the creation of a full NDMP backup.
	FullBackupRetention types.Int64 `tfsdk:"full_backup_retention"`
	// (Used with SyncIQ and NDMP backups.)  The minimum amount of time cloud files will be retained after the creation of a SyncIQ backup or an incremental NDMP backup.
	IncrementalBackupRetention types.Int64 `tfsdk:"incremental_backup_retention"`
	// The minimum amount of time to wait before updating cloud data with local changes.
	WritebackFrequency types.Int64 `tfsdk:"writeback_frequency"`
}

// V12CloudPolicyActionParamsJSONModel - Json model for V12CloudPolicyActionParams.
type V12CloudPolicyActionParamsJSONModel struct {
	Pool                       *string                                   `json:"pool"`
	ArchiveSnapshotFiles       *bool                                     `json:"archive_snapshot_files,omitempty"`
	Cache                      *V12CloudPolicyActionCacheParamsJSONModel `json:"cache,omitempty"`
	Compression                *bool                                     `json:"compression,omitempty"`
	DataRetention              *int64                                    `json:"data_retention,omitempty"`
	Encryption                 *bool                                     `json:"encryption,omitempty"`
	FullBackupRetention        *int64                                    `json:"full_backup_retention,omitempty"`
	IncrementalBackupRetention *int64                                    `json:"incremental_backup_retention,omitempty"`
	WritebackFrequency         *int64                                    `json:"writeback_frequency,omitempty"`
}

// V12CloudPolicyActionCacheParams Specifies default cloudpool cache settings for new filepool policies.
type V12CloudPolicyActionCacheParams struct {
	// Specifies cache expiration.
	Expiration types.Int64 `tfsdk:"expiration"`
	// Specifies cache read ahead type.
	ReadAhead types.String `tfsdk:"read_ahead"`
	// Specifies cache type.
	Type types.String `tfsdk:"type"`
}

// V12CloudPolicyActionCacheParamsJSONModel - Json model for V12CloudPolicyActionCacheParams.
type V12CloudPolicyActionCacheParamsJSONModel struct {
	Expiration *int64  `json:"expiration,omitempty"`
	ReadAhead  *string `json:"read_ahead,omitempty"`
	Type       *string `json:"type,omitempty"`
}

// V1FilepoolPolicyFileMatchingPattern The file matching rules for this policy
type V1FilepoolPolicyFileMatchingPattern struct {
	OrCriteria []V1FilepoolPolicyFileMatchingPatternOrCriteriaItem `tfsdk:"or_criteria"`
}

// V1FilepoolPolicyFileMatchingPatternOrCriteriaItem struct for V1FilepoolPolicyFileMatchingPatternOrCriteriaItem
type V1FilepoolPolicyFileMatchingPatternOrCriteriaItem struct {
	AndCriteria []V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem `tfsdk:"and_criteria"`
}

// V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem struct for V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem
type V1FilepoolPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItem struct {
	// Indicates whether the existence of an attribute indicates a match (valid only with 'type' = 'custom_attribute')
	AttributeExists types.Bool `tfsdk:"attribute_exists"`
	// True to match the path exactly, False to match any subtree. (valid only with 'type' = 'path')
	BeginsWith types.Bool `tfsdk:"begins_with"`
	// True to indicate case sensitivity when comparing file attributes (valid only with 'type' = 'name' or 'type' = 'path')
	CaseSensitive types.Bool `tfsdk:"case_sensitive"`
	// File attribute field name to be compared in a custom comparison (valid only with 'type' = 'custom_attribute')
	Field types.String `tfsdk:"field"`
	// The comparison operator to use while comparing an attribute with its value
	Operator types.String `tfsdk:"operator"`
	// The file attribute to be compared to a given value.
	// Acceptable values are: name, path, link_count, accessed_time, birth_time, changed_time, metadata_changed_time, size, file_type, custom_attribute.
	Type types.String `tfsdk:"type"`
	// Size unit value. One of 'B','KB','MB','GB','TB','PB','EB' (valid only with 'type' = 'size')
	Units types.String `tfsdk:"units"`
	// Whether time units refer to a calendar date and time (e.g., Jun 3, 2009) or a relative duration (e.g., 2 weeks) (valid only with 'type' in {accessed_time, birth_time, changed_time or metadata_changed_time}
	UseRelativeTime types.Bool `tfsdk:"use_relative_time"`
	// The value to be compared against a file attribute.
	Value types.String `tfsdk:"value"`
}
