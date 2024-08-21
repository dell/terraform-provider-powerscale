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
# If resource arguments are omitted, `terraform apply` will load Cluster Owner Details from PowerScale, and save to
# terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Cluster Owner Details (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting Cluster Owner Details from PowerScale.
# For more information, Please check the terraform state file.

resource "powerscale_cluster_owner" "test" {
  company          = "company_name"
  location         = "location"
  primary_email    = "primary_email@example.com"
  primary_name     = "primary_name"
  primary_phone1   = "+91-12345-67890"
  primary_phone2   = "+1 123-456-7890" # primary alternate phone number
  secondary_email  = "secondary_email@example.com"
  secondary_name   = "secondary_name"
  secondary_phone1 = "+44 (20) 1234 5678"
  secondary_phone2 = "+1 (800) 555-5555" # secondary alternate phone number
}

# After the execution of above resource block, Cluster Owner Settings would have been cached in terraform state file, and
# Cluster Owner Settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.
