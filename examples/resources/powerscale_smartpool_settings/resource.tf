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
# If resource arguments are omitted, `terraform apply` will load SmartPools Settings from PowerScale, and save to
# terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load SmartPools Settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting SmartPools Settings from PowerScale.
# For more information, Please check the terraform state file.

resource "powerscale_smartpool_settings" "settings" {
  #    global_namespace_acceleration_enabled = false
  #    manage_io_optimization                = true
  #    manage_io_optimization_apply_to_files = false
  #    manage_protection                     = true
  #    manage_protection_apply_to_files      = false
  #    protect_directories_one_level_higher  = true
  #    spillover_enabled                     = true
  #    spillover_target                      = {
  #      name    = "sample_storagepool"
  #      type    = "storagepool" // anywhere or storagepool. name should be empty string when type is anywhere
  #    }
  #    ssd_l3_cache_default_enabled          = true
  #    ssd_qab_mirrors                       = "all"
  #    ssd_system_btree_mirrors              = "all"
  #    ssd_system_delta_mirrors              = "all"
  #    virtual_hot_spare_deny_writes         = true
  #    virtual_hot_spare_hide_spare          = true
  #    virtual_hot_spare_limit_drives        = 4
  #    virtual_hot_spare_limit_percent       = 4
  #  # Note that, default_transfer_limit_state and default_transfer_limit_pct are mutually exclusive and only one can be specified.
  #    default_transfer_limit_state          = "disabled" // available for PowerScale 9.5 and above
  #    default_transfer_limit_pct            = 90 // available for PowerScale 9.5 and above
}

# After the execution of above resource block, SmartPools Settings would have been cached in terraform state file, or
# SmartPools Settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.