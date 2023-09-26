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

// SnapshotDataSourceModel describes the data source data model.
type SnapshotDataSourceModel struct {
	ID        types.String          `tfsdk:"id"`
	Snapshots []SnapshotDetailModel `tfsdk:"snapshots_details"`
	// filter
	SnapshotFilter *SnapshotFilterType `tfsdk:"filter"`
}

// SnapshotFilterType describes the filter data model.
type SnapshotFilterType struct {
	Path types.String `tfsdk:"path"`
}

// SnapshotDetailModel details of the individual snapshot.
type SnapshotDetailModel struct {
	// The name of the alias, none for real snapshots.
	Alias types.String `tfsdk:"alias"`
	// The Unix Epoch time the snapshot was created.
	Created types.Int64 `tfsdk:"created"`
	// The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.
	Expires types.Int64 `tfsdk:"expires"`
	// True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of locks.
	HasLocks types.Bool `tfsdk:"has_locks"`
	// The system ID given to the snapshot. This is useful for tracking the status of delete pending snapshots.
	ID types.Int64 `tfsdk:"id"`
	// The user or system supplied snapshot name. This will be null for snapshots pending delete.
	Name types.String `tfsdk:"name"`
	// The /ifs path snapshotted.
	Path types.String `tfsdk:"path"`
	// Percentage of /ifs used for storing this snapshot.
	PctFilesystem types.Number `tfsdk:"pct_filesystem"`
	// Percentage of configured snapshot reserved used for storing this snapshot.
	PctReserve types.Number `tfsdk:"pct_reserve"`
	// The name of the schedule used to create this snapshot, if applicable.
	Schedule types.String `tfsdk:"schedule"`
	// The amount of shadow bytes referred to by this snapshot.
	ShadowBytes types.Int64 `tfsdk:"shadow_bytes"`
	// The amount of storage in bytes used to store this snapshot.
	Size types.Int64 `tfsdk:"size"`
	// Snapshot state.
	State types.String `tfsdk:"state"`
	// The ID of the snapshot pointed to if this is an alias. 18446744073709551615 (max uint64) is returned for an alias to the live filesystem.
	TargetID types.Int64 `tfsdk:"target_id"`
	// The name of the snapshot pointed to if this is an alias.
	TargetName types.String `tfsdk:"target_name"`
}
