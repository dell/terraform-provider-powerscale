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

// WritableSnapshot defines the writable snapshot.
type WritableSnapshot struct {
	// The system ID given to the snapshot. This is useful for tracking the status of delete pending snapshots.
	ID types.Int32 `tfsdk:"id"`

	// The destination path of the snapshot.
	DstPath types.String `tfsdk:"dst_path"`

	// The source snapshot ID.
	SrcSnap types.String `tfsdk:"snap_id"`

	// The source snapshot name.
	SnapName types.String `tfsdk:"snap_name"`

	// The source path.
	SrcPath types.String `tfsdk:"src_path"`

	// Snapshot state.
	State types.String `tfsdk:"state"`
}

// WritableSnapshotDataSource defines the writable snapshot data source.
type WritableSnapshotDataSource struct {
	ID       types.Int64  `tfsdk:"id"`
	PhysSize types.Int64  `tfsdk:"phys_size"`
	Created  types.Int64  `tfsdk:"created"`
	State    types.String `tfsdk:"state"`
	SrcID    types.Int64  `tfsdk:"src_id"`
	SrcPath  types.String `tfsdk:"src_path"`
	SrcSnap  types.String `tfsdk:"src_snap"`
	DstPath  types.String `tfsdk:"dst_path"`
	LogSize  types.Int64  `tfsdk:"log_size"`
}

// WritableSnapshotFilter defines the writable snapshot filters.
type WritableSnapshotFilter struct {
	Path   types.String `tfsdk:"path"`
	Sort   types.String `tfsdk:"sort"`
	Resume types.String `tfsdk:"resume"`
	State  types.String `tfsdk:"state"`
	Limit  types.Int32  `tfsdk:"limit"`
	Dir    types.String `tfsdk:"dir"`
}

// WritablesnapshotModel defines the writable snapshot model for data source.
type WritablesnapshotModel struct {
	ID                     types.String                 `tfsdk:"id"`
	Writable               []WritableSnapshotDataSource `tfsdk:"writable"`
	WritableSnapshotFilter *WritableSnapshotFilter      `tfsdk:"filter"`
}
