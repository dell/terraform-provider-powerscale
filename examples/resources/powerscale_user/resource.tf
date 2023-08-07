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

resource "powerscale_user" "testUser" {
  # Required field
  name = "testUserResourceSample"

  # Optional Query
  # query_force = false
  # query_zone = "testZone"
  # query_provider = "testProvider"

  # Optional fields
  # uid      = 11000
  # password = "testPassword"
  # roles    = ["SystemAdmin"]
  # enabled = false
  # unlock = false
  # email = "testTerraform@dell.com"
  # home_directory = "/ifs/home/testUserResourceSample"
  # password_expires = true
  # primary_group = "testPrimaryGroup"
  # prompt_password_change = false
  # shell = "/bin/zsh"
  # sid = "SID:XXXX"
  # expiry = 123456
  # gecos = "testFullName"
}