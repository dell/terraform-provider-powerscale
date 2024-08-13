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
