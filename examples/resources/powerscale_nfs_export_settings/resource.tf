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

# Available actions: Create, Update, Delete and Import.
# If resource arguments are omitted, `terraform apply` will load NFS export settings from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load NFS export settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting NFS export settings from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale NFS Export Settings allow you to configure NFS export settings on PowerScale.
resource "powerscale_nfs_export_settings" "example" {
  # Optional fields both for creating and updating
  #  all_dirs = false
  #  case_insensitive = false
  #  case_preserving = true
  #  commit_asynchronous = false
  #  no_truncate = false
  #  write_datasync_action = "DATASYNC"
  #  security_flavors = ["unix"]
  #  map_non_root = {
  #    enabled = false
  #    primary_group = {}
  #    secondary_groups = []
  #    user = {
  #      id = "USER:nobody"
  #    }
  #  }
  #
  # Specifies the zone in which the export is valid. Notice that update this field will change the resource you manage.
  #  zone = "System"
}

# After the execution of above resource block, NFS export settings would have been cached in terraform state file, or
# NFS export settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.