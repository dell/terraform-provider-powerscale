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

# PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules.

# Returns a list of PowerScale User Mapping Rules based on names and zone filter block. 
data "powerscale_user_mapping_rules" "testUserMappingRules" {
  filter {
    # Optional Names filter for source user name or target user name.
    names = ["admin", "Guest"]
    # Optional Operators filter for user mapping rules.
    operators = ["append", "union"]
    # Optional zone filter. The zone to which the user mapping applies. Defaults to System.
    zone = "System"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_user_mapping_rules.testUserMappingRules
output "powerscale_user_mapping_rules_filter" {
  value = data.powerscale_user_mapping_rules.testUserMappingRules
}


# Returns all of the PowerScale User Mapping Rules.
data "powerscale_user_mapping_rules" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_user_mapping_rules.all
output "powerscale_user_mapping_rules_all" {
  value = data.powerscale_user_mapping_rules.all
}
