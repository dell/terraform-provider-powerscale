/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# PowerScale User Group allows you can do operations on a set of users, groups and well-knowns.

# Returns a list of PowerScale User Groups based on gid or name in names filter block. 
data "powerscale_user_group" "test_user_group" {
  filter {
    # Optional list of names to filter upon
    names = [
      # {
      #   gid = 0
      # },
      # {
      #   name = "Administrators"
      # },
      {
        name = "tfaccUserGroupDatasource"
        gid  = 10000
      }
    ]

    # Optional query parameters.
    cached      = false
    name_prefix = "tfacc"
    # domain = "testDomain"
    # zone = "testZone"
    # provider = "testProvider"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_user_group.test_user_group
output "powerscale_user_group_filter" {
  value = data.powerscale_user_group.test_user_group
}
