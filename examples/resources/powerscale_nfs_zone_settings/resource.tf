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
# If resource arguments are omitted, `terraform apply` will load NFS zone settings from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load NFS zone settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting NFS zone settings from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale NFS zone settings allow you to configure NFS zone settings on PowerScale.
resource "powerscale_nfs_zone_settings" "example" {

  # Required field both for creating and updating
  zone = "tfaccAccessZone"

  # Optional fields both for creating and updating
  #  nfsv4_no_names = false
  #  nfsv4_replace_domain = true
  #  nfsv4_allow_numeric_ids = true
  #  nfsv4_domain = "localdomain_tfaccAZ"
  #  nfsv4_no_domain = false
  #  nfsv4_no_domain_uids = true
}

# After the execution of above resource block, NFS zone settings would have been cached in terraform state file, or
# NFS zone settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.