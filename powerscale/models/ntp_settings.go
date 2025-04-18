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

// NtpSettingsDataSourceModel describes the datasource data model.
type NtpSettingsDataSourceModel struct {
	// Number of nodes that will contact the NTP servers.
	Chimers types.Int64 `tfsdk:"chimers"`
	// Node number (LNN) for nodes excluded from chimer duty.
	Excluded types.List `tfsdk:"excluded"`
	// Path to NTP key file within /ifs.
	KeyFile types.String `tfsdk:"key_file"`
}

// NtpSettingsResourceModel describes the resource data model.
type NtpSettingsResourceModel struct {
	// Number of nodes that will contact the NTP servers.
	Chimers types.Int64 `tfsdk:"chimers"`
	// Node number (LNN) for nodes excluded from chimer duty.
	Excluded types.Set `tfsdk:"excluded"`
	// Path to NTP key file within /ifs.
	KeyFile types.String `tfsdk:"key_file"`
}
