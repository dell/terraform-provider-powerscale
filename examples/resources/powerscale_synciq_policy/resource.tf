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
# If resource arguments are omitted, `terraform apply` will load User Mapping Rules from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load User Mapping Rules (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting User Mapping Rules from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale Sync IQ Policies can be used to replicate files or directories from one cluster to another. 
resource "powerscale_synciq_policy" "policy" {
  # Required
  name             = "policy1"
  action           = "sync"
  source_root_path = "/ifs"
  target_host      = "10.10.10.10"
  target_path      = "/ifs/policy1Sink"

  # Optional
  description = "Policy 1 description"
  enabled     = true
  file_matching_pattern = {
    or_criteria = [
      {
        and_criteria = [
          {
            type     = "name"
            value    = "tfacc"
            operator = "=="
          }
        ]
      }
    ]
  }
}

# After the execution of above resource block, Sync IQ Policies would have been cached in terraform state file, or
# Sync IQ Policies would have been updated on PowerScale.
# For more information, Please check the terraform state file.