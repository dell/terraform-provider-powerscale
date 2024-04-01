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

# This Terraform DataSource is used to query the details of existing role privileges from PowerScale array.

# Returns a list of PowerScale role privileges based on names in the filter block.
data "powerscale_roleprivilege" "test" {
  filter {
    #  Optional query parameters
    #  This will return the role privileges whose name contains the key word (case-insensitive)
    names = ["Shutdown"]
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_roleprivilege.test
output "powerscale_roleprivilege" {
  value = data.powerscale_roleprivilege.test
}

# Returns all PowerScale role privileges on PowerScale array
data "powerscale_roleprivilege" "all" {
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_roleprivilege.all
output "powerscale_roleprivilege_data_all" {
  value = data.powerscale_roleprivilege.all
}