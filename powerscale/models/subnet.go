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

// V12GroupnetSubnets struct for V12GroupnetSubnets.
type V12GroupnetSubnets struct {
	ID           types.String                `tfsdk:"id"`
	Subnets      []V12GroupnetSubnetExtended `tfsdk:"subnets"`
	SubnetFilter *SubnetFilterType           `tfsdk:"filter"`
}

// SubnetFilterType describes the filter data model.
type SubnetFilterType struct {
	Names        []types.String `tfsdk:"names"`
	GroupnetName types.String   `tfsdk:"groupnet_name"`
}

// V12GroupnetSubnetExtended struct for V12GroupnetSubnetExtended.
type V12GroupnetSubnetExtended struct {
	// IP address format.
	AddrFamily types.String `tfsdk:"addr_family"`
	// The base IP address.
	BaseAddr types.String `tfsdk:"base_addr"`
	// A description of the subnet.
	Description types.String `tfsdk:"description"`
	// List of Direct Server Return addresses.
	DsrAddrs types.List `tfsdk:"dsr_addrs"`
	// Gateway IP address.
	Gateway types.String `tfsdk:"gateway"`
	// Gateway priority.
	GatewayPriority types.Int64 `tfsdk:"gateway_priority"`
	// Name of the groupnet this subnet belongs to.
	Groupnet types.String `tfsdk:"groupnet"`
	// Unique Subnet ID.
	ID types.String `tfsdk:"id"`
	// MTU of the subnet.
	Mtu types.Int64 `tfsdk:"mtu"`
	// The name of the subnet.
	Name types.String `tfsdk:"name"`
	// Name of the pools in the subnet.
	Pools types.List `tfsdk:"pools"`
	// Subnet Prefix Length.
	Prefixlen types.Int64 `tfsdk:"prefixlen"`
	// List of IP addresses that SmartConnect listens for DNS requests.
	ScServiceAddrs []V12GroupnetSubnetScServiceAddr `tfsdk:"sc_service_addrs"`
	// Domain Name corresponding to the SmartConnect Service Address.
	ScServiceName types.String `tfsdk:"sc_service_name"`
	// VLAN tagging enabled or disabled.
	VlanEnabled types.Bool `tfsdk:"vlan_enabled"`
	// VLAN ID for all interfaces in the subnet.
	VlanID types.Int64 `tfsdk:"vlan_id"`
}

// V12GroupnetSubnetScServiceAddr struct for V12GroupnetSubnetScServiceAddr.
type V12GroupnetSubnetScServiceAddr struct {
	// High IP
	High types.String `tfsdk:"high"`
	// Low IP
	Low types.String `tfsdk:"low"`
}
