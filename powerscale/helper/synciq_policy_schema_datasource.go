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

package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SynciqpolicyDatasourceSchema is a function that returns the schema for SyncIQPolicyDataSource.
func SynciqpolicyDatasourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: SyncIQPolicyDataSourceSchema(),
	}
}

// SyncIQPolicyDataSourceSchema is a function that returns the schema for SyncIQPolicyDataSource.
func SyncIQPolicyDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Optional:            true,
			Computed:            true,
		},
		"policies": schema.ListNestedAttribute{
			MarkdownDescription: "Policies",
			Description:         "Policies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: V14SyncPolicyExtendedModelSchema()},
		},
	}
}

// V14SyncPolicyExtendedModelSchema is a function that returns the schema for V14SyncPolicyExtendedModel.
func V14SyncPolicyExtendedModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"accelerated_failback": schema.BoolAttribute{
			MarkdownDescription: "Accelerated Failback",
			Description:         "Accelerated Failback",
			Computed:            true,
		},
		"action": schema.StringAttribute{
			MarkdownDescription: "Action",
			Description:         "Action",
			Computed:            true,
		},
		"allow_copy_fb": schema.BoolAttribute{
			MarkdownDescription: "Allow Copy Fb",
			Description:         "Allow Copy Fb",
			Computed:            true,
		},
		"bandwidth_reservation": schema.Int64Attribute{
			MarkdownDescription: "Bandwidth Reservation",
			Description:         "Bandwidth Reservation",
			Computed:            true,
		},
		"changelist": schema.BoolAttribute{
			MarkdownDescription: "Changelist",
			Description:         "Changelist",
			Computed:            true,
		},
		"check_integrity": schema.BoolAttribute{
			MarkdownDescription: "Check Integrity",
			Description:         "Check Integrity",
			Computed:            true,
		},
		"cloud_deep_copy": schema.StringAttribute{
			MarkdownDescription: "Cloud Deep Copy",
			Description:         "Cloud Deep Copy",
			Computed:            true,
		},
		"conflicted": schema.BoolAttribute{
			MarkdownDescription: "Conflicted",
			Description:         "Conflicted",
			Computed:            true,
		},
		"database_mirrored": schema.BoolAttribute{
			MarkdownDescription: "Database Mirrored",
			Description:         "Database Mirrored",
			Computed:            true,
		},
		"delete_quotas": schema.BoolAttribute{
			MarkdownDescription: "Delete Quotas",
			Description:         "Delete Quotas",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Description:         "Description",
			Computed:            true,
		},
		"disable_file_split": schema.BoolAttribute{
			MarkdownDescription: "Disable File Split",
			Description:         "Disable File Split",
			Computed:            true,
		},
		"disable_fofb": schema.BoolAttribute{
			MarkdownDescription: "Disable Fofb",
			Description:         "Disable Fofb",
			Computed:            true,
		},
		"disable_quota_tmp_dir": schema.BoolAttribute{
			MarkdownDescription: "Disable Quota Tmp Dir",
			Description:         "Disable Quota Tmp Dir",
			Computed:            true,
		},
		"disable_stf": schema.BoolAttribute{
			MarkdownDescription: "Disable Stf",
			Description:         "Disable Stf",
			Computed:            true,
		},
		"enable_hash_tmpdir": schema.BoolAttribute{
			MarkdownDescription: "Enable Hash Tmpdir",
			Description:         "Enable Hash Tmpdir",
			Computed:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Enabled",
			Description:         "Enabled",
			Computed:            true,
		},
		"encrypted": schema.BoolAttribute{
			MarkdownDescription: "Encrypted",
			Description:         "Encrypted",
			Computed:            true,
		},
		"encryption_cipher_list": schema.StringAttribute{
			MarkdownDescription: "Encryption Cipher List",
			Description:         "Encryption Cipher List",
			Computed:            true,
		},
		"expected_dataloss": schema.BoolAttribute{
			MarkdownDescription: "Expected Dataloss",
			Description:         "Expected Dataloss",
			Computed:            true,
		},
		"file_matching_pattern": schema.SingleNestedAttribute{
			MarkdownDescription: "File Matching Pattern",
			Description:         "File Matching Pattern",
			Computed:            true,
			Attributes:          V1SyncJobPolicyFileMatchingPatternModelSchema(),
		},
		"force_interface": schema.BoolAttribute{
			MarkdownDescription: "Force Interface",
			Description:         "Force Interface",
			Computed:            true,
		},
		"has_sync_state": schema.BoolAttribute{
			MarkdownDescription: "Has Sync State",
			Description:         "Has Sync State",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Id",
			Description:         "Id",
			Computed:            true,
		},
		"ignore_recursive_quota": schema.BoolAttribute{
			MarkdownDescription: "Ignore Recursive Quota",
			Description:         "Ignore Recursive Quota",
			Computed:            true,
		},
		"job_delay": schema.Int64Attribute{
			MarkdownDescription: "Job Delay",
			Description:         "Job Delay",
			Computed:            true,
		},
		"last_job_state": schema.StringAttribute{
			MarkdownDescription: "Last Job State",
			Description:         "Last Job State",
			Computed:            true,
		},
		"last_started": schema.Int64Attribute{
			MarkdownDescription: "Last Started",
			Description:         "Last Started",
			Computed:            true,
		},
		"last_success": schema.Int64Attribute{
			MarkdownDescription: "Last Success",
			Description:         "Last Success",
			Computed:            true,
		},
		"linked_service_policies": schema.ListAttribute{
			MarkdownDescription: "Linked Service Policies",
			Description:         "Linked Service Policies",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"log_level": schema.StringAttribute{
			MarkdownDescription: "Log Level",
			Description:         "Log Level",
			Computed:            true,
		},
		"log_removed_files": schema.BoolAttribute{
			MarkdownDescription: "Log Removed Files",
			Description:         "Log Removed Files",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",
			Computed:            true,
		},
		"next_run": schema.Int64Attribute{
			MarkdownDescription: "Next Run",
			Description:         "Next Run",
			Computed:            true,
		},
		"ocsp_address": schema.StringAttribute{
			MarkdownDescription: "Ocsp Address",
			Description:         "Ocsp Address",
			Computed:            true,
		},
		"ocsp_issuer_certificate_id": schema.StringAttribute{
			MarkdownDescription: "Ocsp Issuer Certificate Id",
			Description:         "Ocsp Issuer Certificate Id",
			Computed:            true,
		},
		"password_set": schema.BoolAttribute{
			MarkdownDescription: "Password Set",
			Description:         "Password Set",
			Computed:            true,
		},
		"priority": schema.Int64Attribute{
			MarkdownDescription: "Priority",
			Description:         "Priority",
			Computed:            true,
		},
		"report_max_age": schema.Int64Attribute{
			MarkdownDescription: "Report Max Age",
			Description:         "Report Max Age",
			Computed:            true,
		},
		"report_max_count": schema.Int64Attribute{
			MarkdownDescription: "Report Max Count",
			Description:         "Report Max Count",
			Computed:            true,
		},
		"restrict_target_network": schema.BoolAttribute{
			MarkdownDescription: "Restrict Target Network",
			Description:         "Restrict Target Network",
			Computed:            true,
		},
		"rpo_alert": schema.Int64Attribute{
			MarkdownDescription: "Rpo Alert",
			Description:         "Rpo Alert",
			Computed:            true,
		},
		"schedule": schema.StringAttribute{
			MarkdownDescription: "Schedule",
			Description:         "Schedule",
			Computed:            true,
		},
		"service_policy": schema.BoolAttribute{
			MarkdownDescription: "Service Policy",
			Description:         "Service Policy",
			Computed:            true,
		},
		"skip_lookup": schema.BoolAttribute{
			MarkdownDescription: "Skip Lookup",
			Description:         "Skip Lookup",
			Computed:            true,
		},
		"skip_when_source_unmodified": schema.BoolAttribute{
			MarkdownDescription: "Skip When Source Unmodified",
			Description:         "Skip When Source Unmodified",
			Computed:            true,
		},
		"snapshot_sync_existing": schema.BoolAttribute{
			MarkdownDescription: "Snapshot Sync Existing",
			Description:         "Snapshot Sync Existing",
			Computed:            true,
		},
		"snapshot_sync_pattern": schema.StringAttribute{
			MarkdownDescription: "Snapshot Sync Pattern",
			Description:         "Snapshot Sync Pattern",
			Computed:            true,
		},
		"source_certificate_id": schema.StringAttribute{
			MarkdownDescription: "Source Certificate Id",
			Description:         "Source Certificate Id",
			Computed:            true,
		},
		"source_domain_marked": schema.BoolAttribute{
			MarkdownDescription: "Source Domain Marked",
			Description:         "Source Domain Marked",
			Computed:            true,
		},
		"source_exclude_directories": schema.ListAttribute{
			MarkdownDescription: "Source Exclude Directories",
			Description:         "Source Exclude Directories",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"source_include_directories": schema.ListAttribute{
			MarkdownDescription: "Source Include Directories",
			Description:         "Source Include Directories",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"source_network": schema.SingleNestedAttribute{
			MarkdownDescription: "Source Network",
			Description:         "Source Network",
			Computed:            true,
			Attributes:          V1SyncPolicySourceNetworkModelSchema(),
		},
		"source_root_path": schema.StringAttribute{
			MarkdownDescription: "Source Root Path",
			Description:         "Source Root Path",
			Computed:            true,
		},
		"source_snapshot_archive": schema.BoolAttribute{
			MarkdownDescription: "Source Snapshot Archive",
			Description:         "Source Snapshot Archive",
			Computed:            true,
		},
		"source_snapshot_expiration": schema.Int64Attribute{
			MarkdownDescription: "Source Snapshot Expiration",
			Description:         "Source Snapshot Expiration",
			Computed:            true,
		},
		"source_snapshot_pattern": schema.StringAttribute{
			MarkdownDescription: "Source Snapshot Pattern",
			Description:         "Source Snapshot Pattern",
			Computed:            true,
		},
		"sync_existing_snapshot_expiration": schema.BoolAttribute{
			MarkdownDescription: "Sync Existing Snapshot Expiration",
			Description:         "Sync Existing Snapshot Expiration",
			Computed:            true,
		},
		"sync_existing_target_snapshot_pattern": schema.StringAttribute{
			MarkdownDescription: "Sync Existing Target Snapshot Pattern",
			Description:         "Sync Existing Target Snapshot Pattern",
			Computed:            true,
		},
		"target_certificate_id": schema.StringAttribute{
			MarkdownDescription: "Target Certificate Id",
			Description:         "Target Certificate Id",
			Computed:            true,
		},
		"target_compare_initial_sync": schema.BoolAttribute{
			MarkdownDescription: "Target Compare Initial Sync",
			Description:         "Target Compare Initial Sync",
			Computed:            true,
		},
		"target_detect_modifications": schema.BoolAttribute{
			MarkdownDescription: "Target Detect Modifications",
			Description:         "Target Detect Modifications",
			Computed:            true,
		},
		"target_host": schema.StringAttribute{
			MarkdownDescription: "Target Host",
			Description:         "Target Host",
			Computed:            true,
		},
		"target_path": schema.StringAttribute{
			MarkdownDescription: "Target Path",
			Description:         "Target Path",
			Computed:            true,
		},
		"target_snapshot_alias": schema.StringAttribute{
			MarkdownDescription: "Target Snapshot Alias",
			Description:         "Target Snapshot Alias",
			Computed:            true,
		},
		"target_snapshot_archive": schema.BoolAttribute{
			MarkdownDescription: "Target Snapshot Archive",
			Description:         "Target Snapshot Archive",
			Computed:            true,
		},
		"target_snapshot_expiration": schema.Int64Attribute{
			MarkdownDescription: "Target Snapshot Expiration",
			Description:         "Target Snapshot Expiration",
			Computed:            true,
		},
		"target_snapshot_pattern": schema.StringAttribute{
			MarkdownDescription: "Target Snapshot Pattern",
			Description:         "Target Snapshot Pattern",
			Computed:            true,
		},
		"workers_per_node": schema.Int64Attribute{
			MarkdownDescription: "Workers Per Node",
			Description:         "Workers Per Node",
			Computed:            true,
		},
	}
}

// V1SyncJobPolicyFileMatchingPatternModelSchema is a function that returns the schema for V1SyncJobPolicyFileMatchingPatternModel.
func V1SyncJobPolicyFileMatchingPatternModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"or_criteria": schema.ListNestedAttribute{
			MarkdownDescription: "Or Criteria",
			Description:         "Or Criteria",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModelSchema()},
		},
	}
}

// V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModelSchema is a function that returns the schema for V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModel.
func V1SyncJobPolicyFileMatchingPatternOrCriteriaItemModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"and_criteria": schema.ListNestedAttribute{
			MarkdownDescription: "And Criteria",
			Description:         "And Criteria",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModelSchema()},
		},
	}
}

// V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModelSchema is a function that returns the schema for V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModel.
func V1SyncJobPolicyFileMatchingPatternOrCriteriaItemAndCriteriaItemModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"attribute_exists": schema.BoolAttribute{
			MarkdownDescription: "Attribute Exists",
			Description:         "Attribute Exists",
			Computed:            true,
		},
		"case_sensitive": schema.BoolAttribute{
			MarkdownDescription: "Case Sensitive",
			Description:         "Case Sensitive",
			Computed:            true,
		},
		"field": schema.StringAttribute{
			MarkdownDescription: "Field",
			Description:         "Field",
			Computed:            true,
		},
		"operator": schema.StringAttribute{
			MarkdownDescription: "Operator",
			Description:         "Operator",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "Type",
			Description:         "Type",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "Value",
			Description:         "Value",
			Computed:            true,
		},
		"whole_word": schema.BoolAttribute{
			MarkdownDescription: "Whole Word",
			Description:         "Whole Word",
			Computed:            true,
		},
	}
}

// V1SyncPolicySourceNetworkModelSchema is a function that returns the schema for V1SyncPolicySourceNetworkModel.
func V1SyncPolicySourceNetworkModelSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"pool": schema.StringAttribute{
			MarkdownDescription: "Pool",
			Description:         "Pool",
			Computed:            true,
		},
		"subnet": schema.StringAttribute{
			MarkdownDescription: "Subnet",
			Description:         "Subnet",
			Computed:            true,
		},
	}
}
