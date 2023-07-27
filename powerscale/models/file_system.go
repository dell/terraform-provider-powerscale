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

// FileSystemDataSourceModel describes the data source data model.
type FileSystemDataSourceModel struct {
	ID            types.String           `tfsdk:"id"`
	FileSystem    *FileSystemDetailModel `tfsdk:"file_systems_details"`
	DirectoryPath types.String           `tfsdk:"directory_path"`
}

// FileSystemDetailModel details of the Filesystem.
type FileSystemDetailModel struct {
	FileSystemAttribues []FileSystemAttribues `tfsdk:"file_system_attributes"`
	FileSystemQuota     []FileSystemQuota     `tfsdk:"file_system_quotas"`
	FileSystemACL       *FileSystemACL        `tfsdk:"file_system_namespace_acl"`
	FileSystemSnapshots []FileSystemSnaps     `tfsdk:"file_system_snapshots"`
}

// FileSystemSnaps details of the Filesystem.
type FileSystemSnaps struct {
	// Alias name to create for this snapshot. If null, remove any alias.
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

// FileSystemACL details of the Filesystem.
type FileSystemACL struct {
	// Provides the tfsdk array of access rights.
	ACL types.List `tfsdk:"acl"`
	// Action tells if we want to update the existing acl or delete and replace it with new acl defined. Default action is replace.
	Action types.String `tfsdk:"action"`
	// If the directory has access rights set, then this field is returned as acl. If the directory has POSIX permissions set, then this field is returned as mode.
	Authoritative types.String  `tfsdk:"authoritative"`
	Group         *MemberObject `tfsdk:"group"`
	// Provides the POSIX mode.
	Mode  types.String  `tfsdk:"mode"`
	Owner *MemberObject `tfsdk:"owner"`
}

// MemberObject Reusable object in FileSystem objects.
type MemberObject struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

// ACLObject Access Control List object.
type ACLObject struct {
	Accessrights types.ListType `tfsdk:"accessrights"`
	Accesstype   types.String   `tfsdk:"accesstype"`
	InheritFlags types.ListType `tfsdk:"inherit_flags"`
	Op           types.String   `tfsdk:"op"`
	Trustee      *MemberObject  `tfsdk:"trustee"`
}

// FileSystemQuota details of the Filesystem.
type FileSystemQuota struct {
	// If true, SMB shares using the quota directory see the quota thresholds as share size.
	Container types.Bool `tfsdk:"container"`
	// True if the quota provides enforcement, otherwise a accounting quota.
	Enforced types.Bool `tfsdk:"enforced"`
	// The system ID given to the quota.
	ID types.String `tfsdk:"id"`
	// The /ifs path governed.
	Path types.String `tfsdk:"path"`
	// The type of quota.
	Type  types.String         `tfsdk:"type"`
	Usage FileSystemQuotaUsage `tfsdk:"usage"`
}

// FileSystemQuotaUsage quota usage of the Filesystem.
type FileSystemQuotaUsage struct {
	// Bytes used by governed data apparent to application.
	Applogical types.Int64 `tfsdk:"applogical"`
	// True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	ApplogicalReady types.Bool `tfsdk:"applogical_ready"`
	// Bytes used by governed data apparent to filesystem.
	Fslogical types.Int64 `tfsdk:"fslogical"`
	// True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	FslogicalReady types.Bool `tfsdk:"fslogical_ready"`
	// Physical data usage adjusted to account for shadow store efficiency
	Fsphysical types.Int64 `tfsdk:"fsphysical"`
	// True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	FsphysicalReady types.Bool `tfsdk:"fsphysical_ready"`
	// Number of inodes (filesystem entities) used by governed data.
	Inodes types.Int64 `tfsdk:"inodes"`
	// True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	InodesReady types.Bool `tfsdk:"inodes_ready"`
	// Bytes used for governed data and filesystem overhead.
	Physical types.Int64 `tfsdk:"physical"`
	// Number of physical blocks for file data
	PhysicalData types.Int64 `tfsdk:"physical_data"`
	// True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	PhysicalDataReady types.Bool `tfsdk:"physical_data_ready"`
	// Number of physical blocks for file protection
	PhysicalProtection types.Int64 `tfsdk:"physical_protection"`
	// True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	PhysicalProtectionReady types.Bool `tfsdk:"physical_protection_ready"`
	// True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	PhysicalReady types.Bool `tfsdk:"physical_ready"`
	// Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.
	ShadowRefs types.Int64 `tfsdk:"shadow_refs"`
	// True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	ShadowRefsReady types.Bool `tfsdk:"shadow_refs_ready"`
}

// FileSystemAttribues attributes of the Filesystem.
type FileSystemAttribues struct {
	Name      types.String `tfsdk:"name"`
	Namespace types.String `tfsdk:"namespace"`
	Value     types.String `tfsdk:"value"`
}

// FileSystemFilterType describes the filter data model.
type FileSystemFilterType struct {
	Names []types.String `tfsdk:"names"`
}
