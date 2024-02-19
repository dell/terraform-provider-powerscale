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

# PowerScale File Pool Policy can identify logical groups of files and specify storage operations for these files.
# Returns a list of PowerScale File Pool Policies based on names filter block. 
data "powerscale_filepool_policy" "example_filepool_policy" {
  filter {
    # Optional list of names to filter upon
    names = ["filePoolPolicySample", "Default policy"]
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_filepool_policy.example_filepool_policy
output "powerscale_filepool_policy_filter" {
  value = data.powerscale_filepool_policy.example_filepool_policy
}


# Returns all of the PowerScale File Pool Policies including Default Policy.
data "powerscale_filepool_policy" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_filepool_policy.all
output "powerscale_filepool_policy_all" {
  value = data.powerscale_filepool_policy.all
}
