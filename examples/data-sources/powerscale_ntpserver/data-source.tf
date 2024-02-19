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

# This Terraform DataSource is used to query the details of existing NTP servers from PowerScale array.

# Returns a list of PowerScale NTP servers based on names specified in the filter block.
data "powerscale_ntpserver" "test" {
  filter {
    names = ["ntp_server_name_example"]
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_ntpserver.test
output "powerscale_ntpserver" {
  value = data.powerscale_ntpserver.test
}

# Returns all PowerScale NTP servers on PowerScale array
data "powerscale_ntpserver" "all" {
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_ntpserver.all
output "powerscale_ntpserver_data_all" {
  value = data.powerscale_ntpserver.all
}