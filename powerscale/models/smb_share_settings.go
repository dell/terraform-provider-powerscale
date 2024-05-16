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

// SmbShareSettingsResourceModel defines the resource implementation.
type SmbShareSettingsResourceModel struct {
	ID                             types.String `tfsdk:"id"`
	HideDotFiles                   types.Bool   `tfsdk:"hide_dot_files"`
	AllowExecuteAlways             types.Bool   `tfsdk:"allow_execute_always"`
	HostACL                        types.List   `tfsdk:"host_acl"`
	DirectoryCreateMask            types.Int64  `tfsdk:"directory_create_mask"`
	ImpersonateUser                types.String `tfsdk:"impersonate_user"`
	FileFilterExtensions           types.List   `tfsdk:"file_filter_extensions"`
	FileCreateMode                 types.Int64  `tfsdk:"file_create_mode"`
	NtfsACLSupport                 types.Bool   `tfsdk:"ntfs_acl_support"`
	AccessBasedEnumerationRootOnly types.Bool   `tfsdk:"access_based_enumeration_root_only"`
	DirectoryCreateMode            types.Int64  `tfsdk:"directory_create_mode"`
	AllowDeleteReadonly            types.Bool   `tfsdk:"allow_delete_readonly"`
	CaWriteIntegrity               types.String `tfsdk:"ca_write_integrity"`
	StrictFlush                    types.Bool   `tfsdk:"strict_flush"`
	Zone                           types.String `tfsdk:"zone"`
	Smb3EncryptionEnabled          types.Bool   `tfsdk:"smb3_encryption_enabled"`
	MangleByteStart                types.Int64  `tfsdk:"mangle_byte_start"`
	AccessBasedEnumeration         types.Bool   `tfsdk:"access_based_enumeration"`
	FileFilteringEnabled           types.Bool   `tfsdk:"file_filtering_enabled"`
	SparseFile                     types.Bool   `tfsdk:"sparse_file"`
	ChangeNotify                   types.String `tfsdk:"change_notify"`
	MangleMap                      types.List   `tfsdk:"mangle_map"`
	FileCreateMask                 types.Int64  `tfsdk:"file_create_mask"`
	ImpersonateGuest               types.String `tfsdk:"impersonate_guest"`
	StrictCaLockout                types.Bool   `tfsdk:"strict_ca_lockout"`
	FileFilterType                 types.String `tfsdk:"file_filter_type"`
	CreatePermissions              types.String `tfsdk:"create_permissions"`
	CaTimeout                      types.Int64  `tfsdk:"ca_timeout"`
	CscPolicy                      types.String `tfsdk:"csc_policy"`
	Oplocks                        types.Bool   `tfsdk:"oplocks"`
	StrictLocking                  types.Bool   `tfsdk:"strict_locking"`
	ContinuouslyAvailable          types.Bool   `tfsdk:"continuously_available"`
}
