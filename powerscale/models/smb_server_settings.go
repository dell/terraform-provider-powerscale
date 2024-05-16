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

// SmbServerSettingsDataSourceModel defines the data source implementation.
type SmbServerSettingsDataSourceModel struct {
	ID                      types.String             `tfsdk:"id"`
	SmbServerSettings       *SmbServerSettings       `tfsdk:"smb_server_settings"`
	SmbServerSettingsFilter *SmbServerSettingsFilter `tfsdk:"filter"`
}

// SmbServerSettings specifies the configuration values for SMB Server Settings.
type SmbServerSettings struct {
	SupportMultichannel       types.Bool   `tfsdk:"support_multichannel"`
	EnableSecuritySignatures  types.Bool   `tfsdk:"enable_security_signatures"`
	SupportNetbios            types.Bool   `tfsdk:"support_netbios"`
	DotSnapVisibleRoot        types.Bool   `tfsdk:"dot_snap_visible_root"`
	AccessBasedShareEnum      types.Bool   `tfsdk:"access_based_share_enum"`
	DotSnapAccessibleRoot     types.Bool   `tfsdk:"dot_snap_accessible_root"`
	SupportSmb2               types.Bool   `tfsdk:"support_smb2"`
	AuditLogon                types.String `tfsdk:"audit_logon"`
	DotSnapAccessibleChild    types.Bool   `tfsdk:"dot_snap_accessible_child"`
	SrvCPUMultiplier          types.Int64  `tfsdk:"srv_cpu_multiplier"`
	IgnoreEas                 types.Bool   `tfsdk:"ignore_eas"`
	AuditFileshare            types.String `tfsdk:"audit_fileshare"`
	OnefsNumWorkers           types.Int64  `tfsdk:"onefs_num_workers"`
	SrvNumWorkers             types.Int64  `tfsdk:"srv_num_workers"`
	DotSnapVisibleChild       types.Bool   `tfsdk:"dot_snap_visible_child"`
	RequireSecuritySignatures types.Bool   `tfsdk:"require_security_signatures"`
	ServerSideCopy            types.Bool   `tfsdk:"server_side_copy"`
	ServerString              types.String `tfsdk:"server_string"`
	Service                   types.Bool   `tfsdk:"service"`
	SupportSmb3Encryption     types.Bool   `tfsdk:"support_smb3_encryption"`
	RejectUnencryptedAccess   types.Bool   `tfsdk:"reject_unencrypted_access"`
	OnefsCPUMultiplier        types.Int64  `tfsdk:"onefs_cpu_multiplier"`
	GuestUser                 types.String `tfsdk:"guest_user"`
}

// SmbServerSettingsFilter holds the filter conditions.
type SmbServerSettingsFilter struct {
	Scope types.String `tfsdk:"scope"`
}
