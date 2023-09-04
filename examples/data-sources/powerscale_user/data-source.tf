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

data "powerscale_user" "test_user" {
  filter {
    names = [
      # {
      #   uid = 0
      # },
      # {
      #   name = "admin"
      # },
      {
        name = "tfaccUserDatasource"
        uid  = 10000
      }
    ]
    cached      = false
    name_prefix = "tfacc"
    member_of   = false
    # domain = "testDomain"
    # zone = "testZone"
    # provider = "testProvider"
  }
}

output "powerscale_user_filter" {
  value = data.powerscale_user.test_user
}

data "powerscale_user" "test_all_user" {
}

output "powerscale_user_all" {
  value = data.powerscale_user.test_all_user
}