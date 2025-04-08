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

// StoragepoolTierDataSourceModel describes the data source data model.
type StoragepoolTierDataSourceModel struct {
	ID              types.String           `tfsdk:"id"`
	StoragepoolTier []StoragepoolTierModel `tfsdk:"storagepool_tiers"`
}

// StoragepoolTierDetailModel details of the individual storage pool tier.
type StoragepoolTierModel struct {
	Children           types.List                 `tfsdk:"children"`
	Id                 types.String               `tfsdk:"id"`
	Lnns               types.List                 `tfsdk:"lnns"`
	Name               types.String               `tfsdk:"name"`
	TransferLimitPct   types.Int32                `tfsdk:"transfer_limit_pct"`
	TransferLimitState types.String               `tfsdk:"transfer_limit_state"`
	Usage              *StoragepoolTierUsageModel `tfsdk:"usage"`
}

type StoragepoolTierUsageModel struct {
	AvailBytes           types.String `tfsdk:"avail_bytes"`
	AvailHddBytes        types.String `tfsdk:"avail_hdd_bytes"`
	AvailSsdBytes        types.String `tfsdk:"avail_ssd_bytes"`
	Balanced             types.Bool   `tfsdk:"balanced"`
	FreeBytes            types.String `tfsdk:"free_bytes"`
	FreeHddBytes         types.String `tfsdk:"free_hdd_bytes"`
	FreeSsdBytes         types.String `tfsdk:"free_ssd_bytes"`
	PctUsed              types.String `tfsdk:"pct_used"`
	PctUsedHdd           types.String `tfsdk:"pct_used_hdd"`
	PctUsedSsd           types.String `tfsdk:"pct_used_ssd"`
	TotalBytes           types.String `tfsdk:"total_bytes"`
	TotalHddBytes        types.String `tfsdk:"total_hdd_bytes"`
	TotalSsdBytes        types.String `tfsdk:"total_ssd_bytes"`
	UsableBytes          types.String `tfsdk:"usable_bytes"`
	UsableHddBytes       types.String `tfsdk:"usable_hdd_bytes"`
	UsableSsdBytes       types.String `tfsdk:"usable_ssd_bytes"`
	UsedBytes            types.String `tfsdk:"used_bytes"`
	UsedHddBytes         types.String `tfsdk:"used_hdd_bytes"`
	UsedSsdBytes         types.String `tfsdk:"used_ssd_bytes"`
	VirtualHotSpareBytes types.String `tfsdk:"virtual_hot_spare_bytes"`
}

type StoragepoolTierResourceModel struct {
	Children types.List `tfsdk:"children"`
	// Nodepools. The names or IDs of the tier's children.
	Id types.Int64 `tfsdk:"id"`
	// The unique identifier of the storagepool tier.
	Lnns types.List `tfsdk:"lnns"`
	// The nodes that are part of this tier.
	Name types.String `tfsdk:"name"`
	// The tier name.
	TransferLimitPct types.Int32 `tfsdk:"transfer_limit_pct"`
	// Stop moving files to this tier when this limit is met.
	TransferLimitState types.String `tfsdk:"transfer_limit_state"`
	// How the transfer limit value is being applied.
}
