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
# After `terraform apply` of this example file it will create NFS export on specified paths on the PowerScale

# PowerScale provides an NFS server so you can share files on your cluster
resource "powerscale_nfs_export" "example_export" {
  # Required path for creating
  paths = ["/ifs/example_export"]

  # Computed identifier for export
  # id = 1

  # Optional parameter when creating
  # ignore_bad_auth = true
  # ignore_bad_paths= true
  # ignore_conflicts = true
  # ignore_unresolvable_hosts = true
  # zone = "System"

  # Optional parameter when creating and updating. Will ignore the conflict when set true
  # force = true

  # Optional query. Will return the information according to scope
  # scope = "default"

  # Computed attributes, can be updated
  # all_dirs = false
  # block_size = 8192
  # can_set_time = true
  # case_insensitive = true
  # case_preserving = false
  # chown_restricted = false
  # clients = ["client1"]
  # commit_asynchronous = false
  # conflicting_paths = ["/ifs/conflicting_path"]
  # description = "Example path"
  # directory_transfer_size = 131072
  # encoding = "DEFAULT"
  # link_max = 32767
  # map_all = {
  #   enabled = false,
  #   primary_group = {
  #       id = "GROUP:nobody"
  #   }
  #   secondary_groups = [
  #     {
  #       id   = "GROUP:Users"
  #     }
  #   ]
  #   user = {
  #       id = "USER:nobody"
  #   }
  # }
  # map_failure = {}
  # map_full = true
  # map_lookup_uid = false
  # map_non_root = {}
  # map_retry = true
  # map_root = {}
  # max_file_size = 9223372036854775807
  # name_max_size = 255
  # no_truncate = false
  # read_only = false
  # read_only_clients = []
  # read_transfer_max_size = 1048576
  # read_transfer_multiple = 4194304
  # read_transfer_size = 131072
  # read_write_clients = []
  # readdirplus = true
  # readdirplus_prefetch = 10
  # return_32bit_file_ids = false
  # root_clients = []
  # security_flavors = ["unix"]
  # setattr_asynchronous = false
  # snapshot = "-"
  # symlinks = true
  # time_delta = 0.0000000009999999717180685
  # unresolved_clients = []
  # write_datasync_action = "DATASYNC"
  # write_datasync_reply = "DATASYNC"
  # write_filesync_action = "FILESYNC"
  # write_filesync_reply = "FILESYNC"
  # write_transfer_max_size = 1048576
  # write_transfer_multiple = 512
  # write_transfer_size = 524288
  # write_unstable_action = "UNSTABLE"
  # write_unstable_reply = "UNSTABLE"
}

# After the execution of above resource block, an NFS export would have been created on the PowerScale array.
# For more information, Please check the terraform state file.