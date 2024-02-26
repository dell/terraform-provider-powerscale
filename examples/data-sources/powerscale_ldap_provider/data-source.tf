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

# PowerScale LDAP provider enables you to define, query, and modify directory services and resources.

# Returns a list of PowerScale LDAP providers based on names and scope filter block. 
data "powerscale_ldap_provider" "example_ldap_provider" {
  filter {
    # Optional list of names to filter upon
    names = ["ldap_provider_name"]
    # If specified as "effective" or not specified, all fields are returned. If specified as "user", only fields with non-default values are shown. If specified as "default", the original values are returned.
    scope = "effective"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_ldap_provider.example_ldap_provider
output "powerscale_ldap_provider_filter" {
  value = data.powerscale_ldap_provider.example_ldap_provider
}


# Returns all of the PowerScale LDAP providers
data "powerscale_ldap_provider" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_ldap_provider.all
output "powerscale_ldap_provider_all" {
  value = data.powerscale_ldap_provider.all
}
