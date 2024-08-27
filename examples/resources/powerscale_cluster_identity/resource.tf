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
# If resource arguments are omitted, `terraform apply` will load Cluster identity Settings from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Cluster identity Settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting identity Settings from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale Cluster identity Settings allow you to change the settings of identity
resource "powerscale_cluster_identity" "example_identity_settings" {
  # Optional fields both for creating and updating
  name        = "cluster1"
  description = "cluster1 description"
  logon = {
    motd = "motd", motd_header = "motd header"
  }
}
# After the execution of above resource block, Cluster identity Settings would have been cached in terraform state file, or
# identity Settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.