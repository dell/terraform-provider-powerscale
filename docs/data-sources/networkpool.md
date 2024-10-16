---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerscale_networkpool data source"
linkTitle: "powerscale_networkpool"
page_title: "powerscale_networkpool Data Source - terraform-provider-powerscale"
subcategory: ""
description: |-
  This datasource is used to query the existing network pools from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can add network interfaces to network pools to associate address ranges with a node or a group of nodes.
---

# powerscale_networkpool (Data Source)

This datasource is used to query the existing network pools from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can add network interfaces to network pools to associate address ranges with a node or a group of nodes.

## Example Usage

```terraform
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

# This Terraform DataSource is used to query the details of existing network pools from PowerScale array.

# Returns a list of PowerScale network pools based on names and query parameters specified in the filter block.
data "powerscale_networkpool" "test" {
  filter {
    #   Optional query parameters
    #   Note: the following filters will be applied with AND logic
    names = ["pool0"]
    #     subnet = "subnet0"
    #     groupnet = "groupnet0"
    #     access_zone = "System"
    #     alloc_method = "static"
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_networkpool.test
output "powerscale_networkpool" {
  value = data.powerscale_networkpool.test
}

# Returns all PowerScale network pools on PowerScale array
data "powerscale_networkpool" "all" {
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_networkpool.all
output "powerscale_networkpool_data_all" {
  value = data.powerscale_networkpool.all
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `id` (String) Unique identifier of the network pool instance.
- `network_pools_details` (Attributes List) List of Network Pools. (see [below for nested schema](#nestedatt--network_pools_details))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `access_zone` (String) If specified, only pools with this zone name will be returned.
- `alloc_method` (String) If specified, only pools with this allocation type will be returned.
- `groupnet` (String) If specified, only pools for this groupnet will be returned.
- `names` (Set of String) Filter network pools by names.
- `subnet` (String) If specified, only pools for this subnet will be returned.


<a id="nestedatt--network_pools_details"></a>
### Nested Schema for `network_pools_details`

Read-Only:

- `access_zone` (String) Name of a valid access zone to map IP address pool to the zone.
- `addr_family` (String) IP address format.
- `aggregation_mode` (String) OneFS supports the following NIC aggregation modes.
- `alloc_method` (String) Specifies how IP address allocation is done among pool members.
- `description` (String) A description of the pool.
- `groupnet` (String) Name of the groupnet this pool belongs to.
- `id` (String) Unique Pool ID.
- `ifaces` (Attributes List) List of interface members in this pool. (see [below for nested schema](#nestedatt--network_pools_details--ifaces))
- `name` (String) The name of the pool. It must be unique throughout the given subnet.It's a required field with POST method.
- `nfsv3_rroce_only` (Boolean) Indicates that pool contains only RDMA RRoCE capable interfaces.
- `ranges` (Attributes List) List of IP address ranges in this pool. (see [below for nested schema](#nestedatt--network_pools_details--ranges))
- `rebalance_policy` (String) Rebalance policy..
- `rules` (List of String) Names of the rules in this pool.
- `sc_auto_unsuspend_delay` (Number) Time delay in seconds before a node which has been automatically unsuspended becomes usable in SmartConnect responses for pool zones.
- `sc_connect_policy` (String) SmartConnect client connection balancing policy.
- `sc_dns_zone` (String) SmartConnect zone name for the pool.
- `sc_dns_zone_aliases` (List of String) List of SmartConnect zone aliases (DNS names) to the pool.
- `sc_failover_policy` (String) SmartConnect IP failover policy.
- `sc_subnet` (String) Name of SmartConnect service subnet for this pool.
- `sc_suspended_nodes` (List of Number) List of LNNs showing currently suspended nodes in SmartConnect.
- `sc_ttl` (Number) Time to live value for SmartConnect DNS query responses in seconds.
- `static_routes` (Attributes List) List of interface members in this pool. (see [below for nested schema](#nestedatt--network_pools_details--static_routes))
- `subnet` (String) The name of the subnet.

<a id="nestedatt--network_pools_details--ifaces"></a>
### Nested Schema for `network_pools_details.ifaces`

Read-Only:

- `iface` (String) A string that defines an interface name.
- `lnn` (Number) Logical Node Number (LNN) of a node.


<a id="nestedatt--network_pools_details--ranges"></a>
### Nested Schema for `network_pools_details.ranges`

Read-Only:

- `high` (String) High IP
- `low` (String) Low IP


<a id="nestedatt--network_pools_details--static_routes"></a>
### Nested Schema for `network_pools_details.static_routes`

Read-Only:

- `gateway` (String) Address of the gateway in the format: yyy.yyy.yyy.yyy
- `prefixlen` (Number) Prefix length in the format: nn.
- `subnet` (String) Network address in the format: xxx.xxx.xxx.xxx