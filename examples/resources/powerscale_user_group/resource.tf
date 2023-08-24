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

resource "powerscale_user_group" "testUserGroup" {
  # Required field
  name = "testUserGroupResourceSample"

  # Optional Query
  # query_force = false
  # query_zone = "testZone"
  # query_provider = "testProvider"

  # Optional fields
  # gid      = 11000
  # roles    = ["SystemAdmin"]
  # users    = ["MemberOfUser"]
  # groups   = ["MemberOfGroup"]
  # well_knowns    = ["MemberOfWellKnown"]
  # sid = "SID:XXXX"
}