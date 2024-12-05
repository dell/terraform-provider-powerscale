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
# After `terraform apply` of this example file will create the storage pool tier on the PowerScale array with the attributes set in the config.
# For update, name, children, transfer_limit_pct and transfer_limit_state are supported. transfer_limit_pct and transfer_limit_state are mutually exclusive

resource "powerscale_storagepool_tier" "example" {
  # Required field both for creating and updating
  name = "Sample_terraform_tier_7"

  # Optional parameters
  children = [
    "x410_34tb_1.6tb-ssd_64gb"
  ]
  transfer_limit_pct = 40
}