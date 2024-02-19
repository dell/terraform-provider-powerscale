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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will create NFS export on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale SMB shares provide clients network access to file system resources on the cluster
resource "powerscale_smb_share" "share_example" {
  # Required information for creating
  name = "smb_share_example"
  path = "/ifs/smb_share_example"
  permissions = [
    {
      permission      = "full"
      permission_type = "allow"
      trustee = {
        id   = "SID:S-1-1-0",
        name = "Everyone",
        type = "wellknown"
      }
    }
  ]

  # Zone is optional while creating and updating
  # zone = "System"

  # Optional attributes, can be updated
  # access_based_enumeration = false
  # access_based_enumeration_root_only = false
  # allow_delete_readonly = false
  # allow_execute_always = false
  # allow_variable_expansion = false
  # auto_create_directory = true
  # browsable = true
  # ca_timeout = 120
  # ca_write_integrity = "write-read-coherent"
  # change_notify = "norecurse"
  # create_path = false
  # create_permissions = "default acl"
  # csc_policy = "manual"
  # description = "description"
  # directory_create_mask = 448
  # directory_create_mode = 0
  # file_create_mask = 448
  # file_create_mode = 64
  # file_filter_extensions = ["ext"]
  # file_filter_type = "deny"
  # file_filtering_enabled = false
  # hide_dot_files = false
  # host_acl = ["example_host"]
  # impersonate_guest = "never"
  # impersonate_user = ""
  # inheritable_path_acl = false
  # mangle_byte_start = 60672
  # mangle_map = ["0x22:-1"]
  # ntfs_acl_support = true
  # oplocks = true
  # run_as_root = [{
  #   id   = "SID:S-1-1-0",
  #   name = "Everyone",
  #   type = "wellknown"
  # }]
  # smb3_encryption_enabled = false
  # sparse_file = false
  # strict_ca_lockout = true
  # strict_flush = true
  # strict_locking = false

  # zid should be computed according to zone
  # zid = 1
}

# After the execution of above resource block, an SMB share would have been created on the PowerScale array.
# For more information, Please check the terraform state file.