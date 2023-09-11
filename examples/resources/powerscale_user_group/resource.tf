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

# Available actions: Create, Update, Delete and Import.
# After `terraform apply` of this example file it will create a new user group with the name set in `name` attribute on the PowerScale.

# PowerScale User Group allows you can do operations on a set of users, groups and well-knowns.
resource "powerscale_user_group" "testUserGroup" {
  # Required name for creating
  name = "testUserGroupResourceSample"

  # Optional query_force. If true, skip validation checks when creating user group. The force option is required for user group ID changes.
  # query_force = false

  # Optional query parameters when creating and updating. Will return the information according to zone and provider. 
  # query_zone = "testZone"
  # query_provider = "testProvider"

  # Optional parameters when creating
  # sid = "SID:XXXX"

  # Optional parameters when creating and updating. 
  # gid      = 11000
  # roles    = ["SystemAdmin"]
  # users    = ["MemberOfUser"]
  # groups   = ["MemberOfGroup"]
  # well_knowns    = ["MemberOfWellKnown"]
}