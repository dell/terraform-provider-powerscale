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

# Available actions: Create, Update, Delete and Import.
# If resource arguments are omitted, `terraform apply` will load ACL Settings from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load ACL Settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting ACL Settings from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale ACL Settings allow you to manage file and directory permissions, referred to as access rights.
resource "powerscale_aclsettings" "example_acl_settings" {
  # Optional fields both for creating and updating
  # Please check the acceptable inputs for each setting in the documentation
  #     access                  = "windows"
  #     calcmode                = "approx"
  #     calcmode_group          = "group_aces"
  #     calcmode_owner          = "owner_aces"
  #     calcmode_traverse       = "ignore"
  #     chmod                   = "merge"
  #     chmod_007               = "default"
  #     chmod_inheritable       = "no"
  #     chown                   = "owner_group_and_acl"
  #     create_over_smb         = "allow"
  #     dos_attr                = "deny_smb"
  #     group_owner_inheritance = "creator"
  #     rwx                     = "retain"
  #     synthetic_denies        = "remove"
  #     utimes                  = "only_owner"
}

# After the execution of above resource block, ACL Settings would have been cached in terraform state file, or
# ACL Settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.