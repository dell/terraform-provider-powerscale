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

# Returns PowerScale NFS export settings based on filter
data "powerscale_nfs_export_settings" "test" {
  filter = {
    # Used for query parameter, supported by PowerScale Platform API
    # scope = "user"
    # zone  = "System"
  }
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_nfs_export_settings.test
output "powerscale_nfs_export_settings" {
  value = data.powerscale_nfs_export_settings.test
}

# Returns NFS export settings
data "powerscale_nfs_export_settings" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_nfs_export_settings.all
output "powerscale_nfs_export_settings" {
  value = data.powerscale_nfs_export_settings.all
}