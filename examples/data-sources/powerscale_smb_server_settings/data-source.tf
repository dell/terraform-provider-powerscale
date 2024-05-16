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

# Returns PowerScale SMB Server Settings based on filter
data "powerscale_smb_server_settings" "test" {
  filter {
    # Used for query parameter, supported by PowerScale Platform API
    scope = "effective"
  }
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_smb_server_settings.test
output "powerscale_smb_server_settings_test" {
  value = data.powerscale_smb_server_settings.test
}

# Returns SMB Server Settings
data "powerscale_smb_server_settings" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_smb_server_settings.all
output "powerscale_smb_server_settings_all" {
  value = data.powerscale_smb_server_settings.all
}