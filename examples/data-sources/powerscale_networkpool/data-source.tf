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