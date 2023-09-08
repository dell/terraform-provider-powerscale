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

# PowerScale SMB shares provide clients network access to file system resources on the cluster

# Returns a list of PowerScale SMB shares based on name filter block
data "powerscale_smb_share" "example_smb_shares" {
  filter {
    # Used for specify names of SMB shares
    names = ["tfacc_smb_share"]

    # Used for query parameter, supported by PowerScale Platform API
    # dir = "ASC"
    # limit = 1
    # offset = "0"
    # resume = "resume-token"
    # scope = "default"
    # sort = "id"
    # zone = System
  }
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_smb_share.example_smb_shares
output "powerscale_smb_share" {
  value = data.powerscale_smb_share.example_smb_shares
}

# Returns all of the PowerScale SMB shares in default zone
data "powerscale_smb_share" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_smb_share.all_smb_shares
output "powerscale_smb_share_data_all" {
  value = data.powerscale_smb_share.all_smb_shares
}