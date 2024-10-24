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

// SnapshotRestoreModel represents snapshot restore resource model.
type SnapshotRestoreModel struct {
	ID               types.String `tfsdk:"id"`
	SnapRevertParams types.Object `tfsdk:"snaprevert_params"`
	CopyParams       types.Object `tfsdk:"copy_params"`
	CloneParams      types.Object `tfsdk:"clone_params"`
}

// SnapRevertParamsModel represents snapshot revert parameters model.
type SnapRevertParamsModel struct {
	AllowDup types.Bool  `tfsdk:"allow_dup"`
	SnapID   types.Int32 `tfsdk:"snapshot_id"`
}

// CopyParamsModel represents the copy parameters model.
type CopyParamsModel struct {
	Directory types.Object `tfsdk:"directory"`
	File      types.Object `tfsdk:"file"`
}

// DirectoryModel represents the directory model.
type DirectoryModel struct {
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
	Overwrite   types.Bool   `tfsdk:"overwrite"`
	Merge       types.Bool   `tfsdk:"merge"`
	Continue    types.Bool   `tfsdk:"continue"`
}

// FileModel represents the file model.
type FileModel struct {
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
	Overwrite   types.Bool   `tfsdk:"overwrite"`
}

// CloneParamsModel represents the clone parameters model.
type CloneParamsModel struct {
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
	Overwrite   types.Bool   `tfsdk:"overwrite"`
	SnapID      types.Int32  `tfsdk:"snapshot_id"`
}
