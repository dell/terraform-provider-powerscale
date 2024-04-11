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

# This Terraform DataSource is used to query the details of existing Role from PowerScale array.

# Returns a list of PowerScale Role based on the filters specified in the filter block.
data "powerscale_role" "test" {
  filter {
    names = ["SystemAdmin"]
    zone  = "System"
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_role.test
output "powerscale_role" {
  value = data.powerscale_role.test
}

# Returns all PowerScale Role on PowerScale array
data "powerscale_role" "all" {
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_role.all
output "powerscale_role_data_all" {
  value = data.powerscale_role.all
}