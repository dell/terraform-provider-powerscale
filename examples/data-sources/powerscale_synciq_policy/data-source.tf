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

# PowerScale SyncIQ Policy allows you to get a list of SyncIQ Policies or a policy by its ID.

# Returns a list of PowerScale SyncIQ Policies 
data "powerscale_synciq_policy" "all_policies" {
}

# Returns a the PowerScale SyncIQ Policy with given ID
data "powerscale_synciq_policy" "all_policies" {
  id = "g23j9a1f83h12n5j4"
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_synciq_policy.all_policies.policies
output "powerscale_synciq_all_policies" {
  value = data.powerscale_synciq_policy.all_policies.policies
}

# The user can use the fetched policy by ID by the variable data.powerscale_synciq_policy.one.policies[0]
output "policyByID" {
  value = data.powerscale_synciq_policy.one.policies[0]
}

# Get syncIQ policy by name
# Step 1: We shall use the datasource to get all the policies as shown above
# Step 2: We index them by name
locals {
  policiesByName = { for policy in data.powerscale_synciq_policy.all_policies.policies : policy.name => policy }
}

# Syep 3: The user can use the fetched policy by name by the variable local.policiesByName["<anme>"]
output "policiesByName" {
  value = {
    "aut" = local.policiesByName["aut_synciq_report_policy_01"]
    "but" = local.policiesByName["but_synciq_report_policy_01"]
    "cut" = local.policiesByName["cut_synciq_report_policy_01"]
  }
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
