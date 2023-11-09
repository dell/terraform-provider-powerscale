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
# After `terraform apply` of this example file it will create a new subnet with the name set in `name` attribute on the PowerScale

# PowerScale subnet allows you to manage the subnet on the Powerscale array
resource "powerscale_subnet" "subnet" {

  # Required. Name of the new subnet
  name = "subnet1"

  # Required. Name of the groupnet this subnet belongs to
  # Updating is not allowed
  # when managing resource together with groupnet, recommend:
  # groupnet = powerscale_groupnet.example_groupnet.name
  groupnet = "groupnet0"

  # Required. IP address format
  addr_family = "ipv4"

  # Required. Subnet Prefix Length
  prefixlen = 21

  # Optional fields both for creating and updating
  #  description = "terraform subnet"
  #  dsr_addrs = []
  #  gateway="0.0.0.0"
  #  gateway_priority=10
  #  mtu=1500
  #  sc_service_addrs=[{
  #    high="0.0.0.0"
  #    low="0.0.0.0"
  #  }]
  #  sc_service_name=""
  #  vlan_enabled=true
  #  vlan_id=1
}

# After the execution of above resource block, subnet would have been created on the PowerScale array. For more information, Please check the terraform state file.