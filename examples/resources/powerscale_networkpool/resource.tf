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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file for the first time, you will create a network pool on the PowerScale

# PowerScale network pool allows you to add network interfaces to network pools to associate address ranges with a node or a group of nodes.
resource "powerscale_networkpool" "pool_test" {
  #   Required (subnet and groupnet cannot be modified once designated)
  #   Recommend using powerscale_groupnet.groupnet_example.name and powerscale_subnet.subnet_example.name to manage network pool together with groupnet and subnet
  name     = "pool_test"
  subnet   = "subnet0"
  groupnet = "groupnet0"

  #   Optional fields both for creating and updating
  #   access_zone = "System"
  #   aggregation_mode = "lacp"
  #   alloc_method = "static"
  #   description = "Test"
  #   ifaces = [
  #     {
  #       iface = "testIface",
  #       lnn = 0
  #     }
  #   ]
  #   nfsv3_rroce_only = false
  #   ranges = [
  #     {
  #       high = "testIPAddress1",
  #       low = "testIPAddress2"
  #     }
  #   ]
  #   rebalance_policy = "auto"
  #   sc_auto_unsuspend_delay = 0
  #   sc_connect_policy = "round_robin"
  #   sc_dns_zone = "testZoneName"
  #   sc_dns_zone_aliases = ["testZoneNameAlias"]
  #   sc_failover_policy = "round_robin"
  #   sc_subnet = "testSubnetName"
  #   sc_ttl = 0
  #   static_routes = [
  #     {
  #       gateway = "testGatewayAddress",
  #       prefixlen = 10,
  #       subnet = "testSubnetAddress"
  #     }
  #   ]
}

# After the execution of above resource block, network pool would have been created on the PowerScale array.
# For more information, Please check the terraform state file.