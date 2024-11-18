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

# This Terraform DataSource is used to query the details of existing Replication Report from PowerScale array.

# Returns the entire list of PowerScale replication report.
data "powerscale_synciq_replication_report" "all" {
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_replication_report.all
output "powerscale_synciq_replication_report" {
  value = data.powerscale_replication_report.all
}

# Returns a list of PowerScale Replication Report based on the filters specified in the filter block.
data "powerscale_synciq_replication_report" "filtering" {
  filter {
    policy_name        = "Policy"
    reports_per_policy = 2
    sort               = "policy_name"
    dir                = "ASC"
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_replication_report.filtering
output "powerscale_replication_report_filter" {
  value = data.powerscale_replication_report.filtering
}
