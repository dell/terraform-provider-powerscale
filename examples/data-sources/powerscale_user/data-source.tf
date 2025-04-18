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

# PowerScale User allows you to authenticate through a local authentication provider. Remote users are restricted to read-only operations.

# Returns a list of PowerScale Users based on uid or name in names filter block. 
data "powerscale_user" "test_user" {
  filter {
    # Optional list of names to filter upon
    names = [
      # {
      #   uid = 0
      # },
      # {
      #   name = "admin"
      # },
      #{
      #     sid = "S-1-5-21-3219966720-1480896164-796802738-501"
      #},
      #{
      #    name = "admin"
      #    sid = "S-1-5-21-3219966720-1480896164-796802738-501"
      #},
      {
        name = "tfaccUserDatasource"
        uid  = 10000
      }
    ]

    # Optional query parameters.
    cached      = false
    name_prefix = "tfacc"
    member_of   = false
    # domain = "testDomain"
    # zone = "testZone"
    # provider = "testProvider"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_user.test_user
output "powerscale_user_filter" {
  value = data.powerscale_user.test_user
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
