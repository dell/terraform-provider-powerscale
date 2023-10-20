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

// NetworkPoolDataSourceModel describes the data source data model.
type NetworkPoolDataSourceModel struct {
	ID           types.String             `tfsdk:"id"`
	NetworkPools []NetworkPoolDetailModel `tfsdk:"network_pools_details"`

	// Filters
	NetworkPoolFilter *NetworkPoolFilterType `tfsdk:"filter"`
}

// NetworkPoolDetailModel Specifies the properties for a network pool.
type NetworkPoolDetailModel struct {
	// Name of a valid access zone to map IP address pool to the zone.
	AccessZone types.String `tfsdk:"access_zone"`
	// IP address format.
	AddrFamily types.String `tfsdk:"addr_family"`
	// OneFS supports the following NIC aggregation modes.
	AggregationMode types.String `tfsdk:"aggregation_mode"`
	// Specifies how IP address allocation is done among pool members.
	AllocMethod types.String `tfsdk:"alloc_method"`
	// A description of the pool.
	Description types.String `tfsdk:"description"`
	// Name of the groupnet this pool belongs to.
	Groupnet types.String `tfsdk:"groupnet"`
	// Unique Pool ID.
	ID types.String `tfsdk:"id"`
	// List of interface members in this pool.
	Ifaces []V12SubnetsSubnetPoolIface `tfsdk:"ifaces"`
	// The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.
	Name types.String `tfsdk:"name"`
	// Indicates that pool contains only RDMA RRoCE capable interfaces.
	Nfsv3RroceOnly types.Bool `tfsdk:"nfsv3_rroce_only"`
	// List of IP address ranges in this pool.
	Ranges []V12GroupnetSubnetScServiceAddr `tfsdk:"ranges"`
	// Rebalance policy..
	RebalancePolicy types.String `tfsdk:"rebalance_policy"`
	// Names of the rules in this pool.
	Rules types.List `tfsdk:"rules"`
	// Time delay in seconds before a node which has been automatically unsuspended becomes usable in SmartConnect responses for pool zones.
	ScAutoUnsuspendDelay types.Int64 `tfsdk:"sc_auto_unsuspend_delay"`
	// SmartConnect client connection balancing policy.
	ScConnectPolicy types.String `tfsdk:"sc_connect_policy"`
	// SmartConnect zone name for the pool.
	ScDNSZone types.String `tfsdk:"sc_dns_zone"`
	// List of SmartConnect zone aliases (DNS names) to the pool.
	ScDNSZoneAliases types.List `tfsdk:"sc_dns_zone_aliases"`
	// SmartConnect IP failover policy.
	ScFailoverPolicy types.String `tfsdk:"sc_failover_policy"`
	// Name of SmartConnect service subnet for this pool.
	ScSubnet types.String `tfsdk:"sc_subnet"`
	// List of LNNs showing currently suspended nodes in SmartConnect.
	ScSuspendedNodes []types.Int64 `tfsdk:"sc_suspended_nodes"`
	// Time to live value for SmartConnect DNS query responses in seconds.
	ScTTL types.Int64 `tfsdk:"sc_ttl"`
	// List of interface members in this pool.
	StaticRoutes []V12SubnetsSubnetPoolStaticRoute `tfsdk:"static_routes"`
	// The name of the subnet.
	Subnet types.String `tfsdk:"subnet"`
}

// V12SubnetsSubnetPoolIface struct for pool interface.
type V12SubnetsSubnetPoolIface struct {
	// A string that defines an interface name.
	Iface types.String `tfsdk:"iface"`
	// Logical Node Number (LNN) of a node.
	Lnn types.Int64 `tfsdk:"lnn"`
}

// V12SubnetsSubnetPoolStaticRoute struct for static route.
type V12SubnetsSubnetPoolStaticRoute struct {
	// Address of the gateway in the format: yyy.yyy.yyy.yyy
	Gateway types.String `tfsdk:"gateway"`
	// Prefix length in the format: nn.
	Prefixlen types.Int64 `tfsdk:"prefixlen"`
	// Network address in the format: xxx.xxx.xxx.xxx
	Subnet types.String `tfsdk:"subnet"`
}

// NetworkPoolFilterType describes the filter data model.
type NetworkPoolFilterType struct {
	Names       []types.String `tfsdk:"names"`
	Subnet      types.String   `tfsdk:"subnet"`
	Groupnet    types.String   `tfsdk:"groupnet"`
	AccessZone  types.String   `tfsdk:"access_zone"`
	AllocMethod types.String   `tfsdk:"alloc_method"`
}

// NetworkPoolResourceModel describes the resource data model.
type NetworkPoolResourceModel struct {
	// Name of a valid access zone to map IP address pool to the zone.
	AccessZone types.String `tfsdk:"access_zone"`
	// IP address format.
	AddrFamily types.String `tfsdk:"addr_family"`
	// OneFS supports the following NIC aggregation modes.
	AggregationMode types.String `tfsdk:"aggregation_mode"`
	// Specifies how IP address allocation is done among pool members.
	AllocMethod types.String `tfsdk:"alloc_method"`
	// A description of the pool.
	Description types.String `tfsdk:"description"`
	// Name of the groupnet this pool belongs to.
	Groupnet types.String `tfsdk:"groupnet"`
	// Unique Pool ID.
	ID types.String `tfsdk:"id"`
	// List of interface members in this pool.
	Ifaces types.List `tfsdk:"ifaces"`
	// The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.
	Name types.String `tfsdk:"name"`
	// Indicates that pool contains only RDMA RRoCE capable interfaces.
	Nfsv3RroceOnly types.Bool `tfsdk:"nfsv3_rroce_only"`
	// List of IP address ranges in this pool.
	Ranges types.List `tfsdk:"ranges"`
	// Rebalance policy..
	RebalancePolicy types.String `tfsdk:"rebalance_policy"`
	// Names of the rules in this pool.
	Rules types.List `tfsdk:"rules"`
	// Time delay in seconds before a node which has been automatically unsuspended becomes usable in SmartConnect responses for pool zones.
	ScAutoUnsuspendDelay types.Int64 `tfsdk:"sc_auto_unsuspend_delay"`
	// SmartConnect client connection balancing policy.
	ScConnectPolicy types.String `tfsdk:"sc_connect_policy"`
	// SmartConnect zone name for the pool.
	ScDNSZone types.String `tfsdk:"sc_dns_zone"`
	// List of SmartConnect zone aliases (DNS names) to the pool.
	ScDNSZoneAliases types.List `tfsdk:"sc_dns_zone_aliases"`
	// SmartConnect IP failover policy.
	ScFailoverPolicy types.String `tfsdk:"sc_failover_policy"`
	// Name of SmartConnect service subnet for this pool.
	ScSubnet types.String `tfsdk:"sc_subnet"`
	// List of LNNs showing currently suspended nodes in SmartConnect.
	ScSuspendedNodes types.List `tfsdk:"sc_suspended_nodes"`
	// Time to live value for SmartConnect DNS query responses in seconds.
	ScTTL types.Int64 `tfsdk:"sc_ttl"`
	// List of interface members in this pool.
	StaticRoutes types.List `tfsdk:"static_routes"`
	// The name of the subnet.
	Subnet types.String `tfsdk:"subnet"`
}
