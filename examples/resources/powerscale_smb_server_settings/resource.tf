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
# If resource arguments are omitted, `terraform apply` will load SMB server settings from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load SMB server settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting SMB server settings from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale SMB server settings allow you to configure SMB server settings on PowerScale.
resource "powerscale_smb_server_settings" "example" {

  # Optional fields both for creating and updating
  scope = "effective"
  # access_based_share_enum = false
  # dot_snap_accessible_child = true
  # dot_snap_accessible_root = true
  # dot_snap_visible_child = false
  # dot_snap_visible_root = true
  # enable_security_signatures = false
  # guest_user = "nobody"
  # ignore_eas = false
  # onefs_cpu_multiplier = 4
  # onefs_num_workers = 0
  # reject_unencrypted_access = true
  # require_security_signatures = false
  # server_side_copy = true
  # server_string = "PowerScale Server"
  # service = true
  # support_multichannel = true
  # support_netbios = false
  # support_smb2 = true
  # support_smb3_encryption = false
}

# After the execution of above resource block, SMB server settings would have been cached in terraform state file, or
# SMB server settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.