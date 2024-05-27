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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will create NFS export on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale SMB share settings provide clients network access to file system resources on the cluster
resource "powerscale_smb_share_settings" "example" {

  # Required field both for creating and updating
  zone = "System"

  # Optional fields both for creating and updating
  # Please check the acceptable inputs for each setting in the documentation
  #access_based_enumeration           = true
  #access_based_enumeration_root_only = true
  #allow_delete_readonly              = true
  #allow_execute_always               = true
  #ca_timeout                         = 12
  #ca_write_integrity                 = "write-read-coherent"
  #change_notify                      = "none"
  #create_permissions                 = "default acl"
  #directory_create_mask              = 0
  #directory_create_mode              = 0
  #file_create_mask                   = 448
  #file_create_mode                   = 64
  #file_filter_extensions             = []
  #file_filter_type                   = "deny"
  #file_filtering_enabled             = true
  #hide_dot_files                     = true
  #host_acl                           = []
  #impersonate_guest                  = "never"
  #mangle_byte_start                  = 258
  #mangle_map                         = [
  #    "0x01-0x1F:-1",
  #    "0x22:-1",
  #    "0x2A:-1",
  #    "0x3A:-1",
  #    "0x3C:-1",
  #   "0x3E:-1",
  #    "0x3F:-1",
  #    "0x5C:-1",
  #]
  #ntfs_acl_support                   = true
  #oplocks                            = true
  #smb3_encryption_enabled            = true
  #sparse_file                        = true
  #strict_ca_lockout                  = true
  #strict_flush                       = true
  #strict_locking                     = true
  #scope                             = "effective"
}

# After the execution of above resource block, an SMB share would have been created on the PowerScale array.
# For more information, Please check the terraform state file.