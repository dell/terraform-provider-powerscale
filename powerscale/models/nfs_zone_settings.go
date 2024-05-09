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

// NfsZoneSettingsResourceModel defines the resource implementation.
type NfsZoneSettingsResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Zone                 types.String `tfsdk:"zone"`
	Nfsv4NoNames         types.Bool   `tfsdk:"nfsv4_no_names"`
	Nfsv4ReplaceDomain   types.Bool   `tfsdk:"nfsv4_replace_domain"`
	Nfsv4AllowNumericIds types.Bool   `tfsdk:"nfsv4_allow_numeric_ids"`
	Nfsv4Domain          types.String `tfsdk:"nfsv4_domain"`
	Nfsv4NoDomain        types.Bool   `tfsdk:"nfsv4_no_domain"`
	Nfsv4NoDomainUids    types.Bool   `tfsdk:"nfsv4_no_domain_uids"`
}
