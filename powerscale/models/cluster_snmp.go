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

// ClusterSNMPModel is the model for the Cluster SNMP.
type ClusterSNMPModel struct {
	ID                  types.String `tfsdk:"id"`
	Service             types.Bool   `tfsdk:"enabled"`
	ReadOnlyCommunity   types.String `tfsdk:"read_only_community"`
	SnmpV1V2cAccess     types.Bool   `tfsdk:"snmp_v1_v2c_access"`
	SnmpV3Access        types.Bool   `tfsdk:"snmp_v3_access"`
	SnmpV3Password      types.String `tfsdk:"snmp_v3_password"`
	SnmpV3AuthProtocol  types.String `tfsdk:"snmp_v3_auth_protocol"`
	SnmpV3PrivProtocol  types.String `tfsdk:"snmp_v3_priv_protocol"`
	SnmpV3PrivPassword  types.String `tfsdk:"snmp_v3_priv_password"`
	SnmpV3ReadOnlyUser  types.String `tfsdk:"snmp_v3_read_only_user"`
	SnmpV3SecurityLevel types.String `tfsdk:"snmp_v3_security_level"`
	SystemContact       types.String `tfsdk:"system_contact"`
	SystemLocation      types.String `tfsdk:"system_location"`
}
