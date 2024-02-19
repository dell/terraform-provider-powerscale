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

# This Terraform DataSource is used to query the details of existing ADS providers from PowerScale array.

# Returns a list of PowerScale ADS providers based on names and scope specified in the filter block.
data "powerscale_adsprovider" "test" {
  filter {
    names = ["PIE.LAB.EMC.COM"]
    scope = "effective"
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_adsprovider.test
output "powerscale_adsprovider" {
  value = data.powerscale_adsprovider.test
}

# Returns all PowerScale ADS providers on PowerScale array
data "powerscale_adsprovider" "all" {
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_adsprovider.all
output "powerscale_adsprovider_data_all" {
  value = data.powerscale_adsprovider.all
}