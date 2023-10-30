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

// SmartPoolSettingsDataSource smartpool settings struct.
type SmartPoolSettingsDataSource struct {
	ID       types.String       `tfsdk:"id"`
	Settings *SmartPoolSettings `tfsdk:"settings"`
}

// SmartPoolSettings smartpool settings detail.
type SmartPoolSettings struct {
	// ID               types.String             `tfsdk:"id"`
	// Automatically manage IO optimization settings on files.
	ManageIoOptimization             types.Bool `tfsdk:"manage_io_optimization"`
	ManageIoOptimizationApplyToFiles types.Bool `tfsdk:"manage_io_optimization_apply_to_files"`
	// Automatically manage protection settings on files.
	ManageProtection             types.Bool `tfsdk:"manage_protection"`
	ManageProtectionApplyToFiles types.Bool `tfsdk:"manage_protection_apply_to_files"`
	// Optimize namespace operations by storing metadata on SSDs.
	GlobalNamespaceAccelerationEnabled types.Bool `tfsdk:"global_namespace_acceleration_enabled"`
	// Whether or not namespace operation optimizations are currently in effect.
	GlobalNamespaceAccelerationState types.String `tfsdk:"global_namespace_acceleration_state"`
	// Automatically add additional protection level to all directories.
	ProtectDirectoriesOneLevelHigher types.Bool `tfsdk:"protect_directories_one_level_higher"`
	// Spill writes into other pools as needed.
	SpilloverEnabled types.Bool                          `tfsdk:"spillover_enabled"`
	SpilloverTarget  *StoragepoolSettingsSpilloverTarget `tfsdk:"spillover_target"`
	// The L3 Cache default enabled state. This specifies whether L3 Cache should be enabled on new node pools.
	SsdL3CacheDefaultEnabled types.Bool `tfsdk:"ssd_l3_cache_default_enabled"`
	// Controls number of mirrors of QAB blocks to place on SSDs.
	SsdQabMirrors types.String `tfsdk:"ssd_qab_mirrors"`
	// Controls number of mirrors of system B-tree blocks to place on SSDs.
	SsdSystemBtreeMirrors types.String `tfsdk:"ssd_system_btree_mirrors"`
	// Controls number of mirrors of system delta blocks to place on SSDs.
	SsdSystemDeltaMirrors types.String `tfsdk:"ssd_system_delta_mirrors"`
	// Deny writes into reserved virtual hot spare space.
	VirtualHotSpareDenyWrites types.Bool `tfsdk:"virtual_hot_spare_deny_writes"`
	// Hide reserved virtual hot spare space from free space counts.
	VirtualHotSpareHideSpare types.Bool `tfsdk:"virtual_hot_spare_hide_spare"`
	// The number of drives to reserve for the virtual hot spare, from 0-4.
	VirtualHotSpareLimitDrives types.Int64 `tfsdk:"virtual_hot_spare_limit_drives"`
	// The percent space to reserve for the virtual hot spare, from 0-20.
	VirtualHotSpareLimitPercent types.Int64  `tfsdk:"virtual_hot_spare_limit_percent"`
	DefaultTransferLimitState   types.String `tfsdk:"default_transfer_limit_state"`
	DefaultTransferLimitPct     types.Number `tfsdk:"default_transfer_limit_pct"`
}

// StoragepoolSettingsSpilloverTarget Target pool for spilled writes.
type StoragepoolSettingsSpilloverTarget struct {
	// Target pool ID if target specified, otherwise null.
	ID types.Int64 `tfsdk:"id"`
	// Target pool name if target specified, otherwise null.
	Name types.String `tfsdk:"name"`
	// Type of target pool.
	Type types.String `tfsdk:"type"`
}
