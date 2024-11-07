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

# PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.

# Returns a list of PowerScale Groupnets based on names filter block. 
data "powerscale_groupnet" "example_groupnet" {
  filter {
    # Optional list of names to filter upon
    names = ["groupnet_name"]
  }
}

# Returns a list of PowerScale Groupnets in order based on the filters in the filter block. 
data "powerscale_groupnet" "example_groupnet" {
  filter {
    sort = "name"
    dir  = "DESC"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_groupnet.example_groupnet
output "powerscale_groupnet_filter" {
  value = data.powerscale_groupnet.example_groupnet
}


# Returns all of the PowerScale Groupnets
data "powerscale_groupnet" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_groupnet.all
output "powerscale_groupnet_all" {
  value = data.powerscale_groupnet.all
}
