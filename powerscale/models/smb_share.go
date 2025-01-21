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

import "github.com/hashicorp/terraform-plugin-framework/types"

// SmbShareResource smb share schema attribute details.
type SmbShareResource struct {
	// ID of the smb share
	ID types.String `tfsdk:"id"`
	// Only enumerate files and folders the requesting user has access to.
	AccessBasedEnumeration types.Bool `tfsdk:"access_based_enumeration"`
	// Access-based enumeration on only the root directory of the share.
	AccessBasedEnumerationRootOnly types.Bool `tfsdk:"access_based_enumeration_root_only"`
	// Allow deletion of read-only files in the share.
	AllowDeleteReadonly types.Bool `tfsdk:"allow_delete_readonly"`
	// Allows users to execute files they have read rights for.
	AllowExecuteAlways types.Bool `tfsdk:"allow_execute_always"`
	// Allow automatic expansion of variables for home directories.
	AllowVariableExpansion types.Bool `tfsdk:"allow_variable_expansion"`
	// Automatically create home directories.
	AutoCreateDirectory types.Bool `tfsdk:"auto_create_directory"`
	// Share is visible in net view and the browse list.
	Browsable types.Bool `tfsdk:"browsable"`
	// Persistent open timeout for the share.
	CaTimeout types.Int64 `tfsdk:"ca_timeout"`
	// Specify the level of write-integrity on continuously available shares.
	CaWriteIntegrity types.String `tfsdk:"ca_write_integrity"`
	// Level of change notification alerts on the share.
	ChangeNotify types.String `tfsdk:"change_notify"`
	// Specify if persistent opens are allowed on the share.
	ContinuouslyAvailable types.Bool `tfsdk:"continuously_available"`
	// Create path if does not exist.
	CreatePath types.Bool `tfsdk:"create_path"`
	// Create permissions for new files and directories in share.
	CreatePermissions types.String `tfsdk:"create_permissions"`
	// Client-side caching policy for the shares.
	CscPolicy types.String `tfsdk:"csc_policy"`
	// Description for this SMB share.
	Description types.String `tfsdk:"description"`
	// Directory create mask bits.
	DirectoryCreateMask types.Int64 `tfsdk:"directory_create_mask"`
	// Directory create mode bits.
	DirectoryCreateMode types.Int64 `tfsdk:"directory_create_mode"`
	// File create mask bits.
	FileCreateMask types.Int64 `tfsdk:"file_create_mask"`
	// File create mode bits.
	FileCreateMode types.Int64 `tfsdk:"file_create_mode"`
	// Specifies the list of file extensions.
	FileFilterExtensions types.List `tfsdk:"file_filter_extensions"`
	// Specifies if filter list is for deny or allow. Default is deny.
	FileFilterType types.String `tfsdk:"file_filter_type"`
	// Enables file filtering on this zone.
	FileFilteringEnabled types.Bool `tfsdk:"file_filtering_enabled"`
	// Hide files and directories that begin with a period '.'.
	HideDotFiles types.Bool `tfsdk:"hide_dot_files"`
	// An ACL expressing which hosts are allowed access. A deny clause must be the final entry.
	HostACL types.List `tfsdk:"host_acl"`
	// Specify the condition in which user access is done as the guest account.
	ImpersonateGuest types.String `tfsdk:"impersonate_guest"`
	// User account to be used as guest account.
	ImpersonateUser types.String `tfsdk:"impersonate_user"`
	// Set the inheritable ACL on the share path.
	InheritablePathACL types.Bool `tfsdk:"inheritable_path_acl"`
	// Specifies the wchar_t starting point for automatic byte mangling.
	MangleByteStart types.Int64 `tfsdk:"mangle_byte_start"`
	// Character mangle map.
	MangleMap types.List `tfsdk:"mangle_map"`
	// Share name.
	Name types.String `tfsdk:"name"`
	// Support NTFS ACLs on files and directories.
	NtfsACLSupport types.Bool `tfsdk:"ntfs_acl_support"`
	// Support oplocks.
	Oplocks types.Bool `tfsdk:"oplocks"`
	// Path of share within /ifs.
	Path types.String `tfsdk:"path"`
	// Specifies an ordered list of permission modifications.
	Permissions types.List `tfsdk:"permissions"`
	// Allow account to run as root.
	RunAsRoot types.List `tfsdk:"run_as_root"`
	// Enables SMB3 encryption for the share.
	Smb3EncryptionEnabled types.Bool `tfsdk:"smb3_encryption_enabled"`
	// Enables sparse file.
	SparseFile types.Bool `tfsdk:"sparse_file"`
	// Specifies if persistent opens would do strict lockout on the share.
	StrictCaLockout types.Bool `tfsdk:"strict_ca_lockout"`
	// Handle SMB flush operations.
	StrictFlush types.Bool `tfsdk:"strict_flush"`
	// Specifies whether byte range locks contend against SMB I/O.
	StrictLocking types.Bool `tfsdk:"strict_locking"`
	// Name of the access zone to which to move this SMB share.
	Zone types.String `tfsdk:"zone"`
	// Numeric ID of the access zone which contains this SMB share.
	Zid types.Int64 `tfsdk:"zid"`
}

// V1SmbSharePermission Specifies properties for an Access Control Entry.
type V1SmbSharePermission struct {
	// Specifies the file system rights that are allowed or denied.
	Permission types.String `tfsdk:"permission"`
	// Determines whether the permission is allowed or denied.
	PermissionType types.String `tfsdk:"permission_type"`
	//
	Trustee V1AuthAccessAccessItemFileGroup `tfsdk:"trustee"`
}

// SmbShareDatasource holds smb share datasource schema attribute details.
type SmbShareDatasource struct {
	ID              types.String               `tfsdk:"id"`
	SmbShares       []SmbShareDatasourceEntity `tfsdk:"smb_shares"`
	SmbSharesFilter *SmbShareDatasourceFilter  `tfsdk:"filter"`
}

// SmbShareDatasourceFilter holds filter conditions.
type SmbShareDatasourceFilter struct {
	// supported by api
	Sort   types.String `tfsdk:"sort"`
	Zone   types.String `tfsdk:"zone"`
	Resume types.String `tfsdk:"resume"`
	Limit  types.Int32  `tfsdk:"limit"`
	Offset types.Int32  `tfsdk:"offset"`
	Scope  types.String `tfsdk:"scope"`
	Dir    types.String `tfsdk:"dir"`
	// custom name list
	Names []types.String `tfsdk:"names"`
}

// SmbShareDatasourceEntity struct for SmbShareDatasource.
type SmbShareDatasourceEntity struct {
	// Only enumerate files and folders the requesting user has access to.
	AccessBasedEnumeration types.Bool `tfsdk:"access_based_enumeration"`
	// Access-based enumeration on only the root directory of the share.
	AccessBasedEnumerationRootOnly types.Bool `tfsdk:"access_based_enumeration_root_only"`
	// Allow deletion of read-only files in the share.
	AllowDeleteReadonly types.Bool `tfsdk:"allow_delete_readonly"`
	// Allows users to execute files they have read rights for.
	AllowExecuteAlways types.Bool `tfsdk:"allow_execute_always"`
	// Allow automatic expansion of variables for home directories.
	AllowVariableExpansion types.Bool `tfsdk:"allow_variable_expansion"`
	// Automatically create home directories.
	AutoCreateDirectory types.Bool `tfsdk:"auto_create_directory"`
	// Share is visible in net view and the browse list.
	Browsable types.Bool `tfsdk:"browsable"`
	// Persistent open timeout for the share.
	CaTimeout types.Int64 `tfsdk:"ca_timeout"`
	// Specify the level of write-integrity on continuously available shares.
	CaWriteIntegrity types.String `tfsdk:"ca_write_integrity"`
	// Level of change notification alerts on the share.
	ChangeNotify types.String `tfsdk:"change_notify"`
	// Specify if persistent opens are allowed on the share.
	ContinuouslyAvailable types.Bool `tfsdk:"continuously_available"`
	// Create permissions for new files and directories in share.
	CreatePermissions types.String `tfsdk:"create_permissions"`
	// Client-side caching policy for the shares.
	CscPolicy types.String `tfsdk:"csc_policy"`
	// Description for this SMB share.
	Description types.String `tfsdk:"description"`
	// Directory create mask bits.
	DirectoryCreateMask types.Int64 `tfsdk:"directory_create_mask"`
	// Directory create mode bits.
	DirectoryCreateMode types.Int64 `tfsdk:"directory_create_mode"`
	// File create mask bits.
	FileCreateMask types.Int64 `tfsdk:"file_create_mask"`
	// File create mode bits.
	FileCreateMode types.Int64 `tfsdk:"file_create_mode"`
	// Specifies the list of file extensions.
	FileFilterExtensions types.List `tfsdk:"file_filter_extensions"`
	// Specifies if filter list is for deny or allow. Default is deny.
	FileFilterType types.String `tfsdk:"file_filter_type"`
	// Enables file filtering on this zone.
	FileFilteringEnabled types.Bool `tfsdk:"file_filtering_enabled"`
	// Hide files and directories that begin with a period '.'.
	HideDotFiles types.Bool `tfsdk:"hide_dot_files"`
	// An ACL expressing which hosts are allowed access. A deny clause must be the final entry.
	HostACL types.List `tfsdk:"host_acl"`
	// Share ID.
	ID types.String `tfsdk:"id"`
	// Specify the condition in which user access is done as the guest account.
	ImpersonateGuest types.String `tfsdk:"impersonate_guest"`
	// User account to be used as guest account.
	ImpersonateUser types.String `tfsdk:"impersonate_user"`
	// Set the inheritable ACL on the share path.
	InheritablePathACL types.Bool `tfsdk:"inheritable_path_acl"`
	// Specifies the wchar_t starting point for automatic byte mangling.
	MangleByteStart types.Int64 `tfsdk:"mangle_byte_start"`
	// Character mangle map.
	MangleMap types.List `tfsdk:"mangle_map"`
	// Share name.
	Name types.String `tfsdk:"name"`
	// Support NTFS ACLs on files and directories.
	NtfsACLSupport types.Bool `tfsdk:"ntfs_acl_support"`
	// Support oplocks.
	Oplocks types.Bool `tfsdk:"oplocks"`
	// Path of share within /ifs.
	Path types.String `tfsdk:"path"`
	// Specifies an ordered list of permission modifications.
	Permissions []V1SmbSharePermission `tfsdk:"permissions"`
	// Allow account to run as root.
	RunAsRoot []V1AuthAccessAccessItemFileGroup `tfsdk:"run_as_root"`
	// Enables SMB3 encryption for the share.
	Smb3EncryptionEnabled types.Bool `tfsdk:"smb3_encryption_enabled"`
	// Enables sparse file.
	SparseFile types.Bool `tfsdk:"sparse_file"`
	// Specifies if persistent opens would do strict lockout on the share.
	StrictCaLockout types.Bool `tfsdk:"strict_ca_lockout"`
	// Handle SMB flush operations.
	StrictFlush types.Bool `tfsdk:"strict_flush"`
	// Specifies whether byte range locks contend against SMB I/O.
	StrictLocking types.Bool `tfsdk:"strict_locking"`
	// Numeric ID of the access zone which contains this SMB share
	Zid types.Int64 `tfsdk:"zid"`
}
