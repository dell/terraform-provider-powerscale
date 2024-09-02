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

# Available actions: Create and  Update updates the syncIQ  global settings. Delete will delete the state file. Import action is also available
# After `terraform apply` of this example file will update the settings according to the attributes set in the config

# PowerScale SynIQ global settings allows you to manage the global settings on the Powerscale array
resource "powerscale_synciq_global_settings" "example" {
  preferred_rpo_alert = 3
  source_network = {
    subnet = "subnet0"
    pool   = "pool0"
  }
  report_email            = ["example1@mail.com", "example2@mail.com"]
  renegotiation_period    = 28800
  service                 = "paused"
  rpo_alerts              = true
  restrict_target_network = true
  report_max_count        = 2000
}

# There are other attributes as well. Please refer the documentation.
# After the execution of above resource block, settings would have been updated on the PowerScale array. For more information, Please check the terraform state file.