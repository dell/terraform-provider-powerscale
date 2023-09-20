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
# After `terraform apply` of this example file it will create a new user with the name set in `name` attribute on the PowerScale.

# PowerScale User allows you to authenticate through a local authentication provider. Remote users are restricted to read-only operations.
resource "powerscale_user" "testUser" {
  # Required name for creating
  name = "testUserResourceSample"

  # Optional query_force. If true, skip validation checks when creating user. The force option is required for user ID changes.
  # query_force = false

  # Optional query_zone, will return user according to zone. Specifies the zone that the user will belong to when creating. Once user is created, its zone cannot be changed.
  # query_zone = "testZone"

  # Optional query_provider, will return user according to provider. Specifies the provider that the user will belong to when creating. Once user is created, its provider cannot be changed.
  # query_provider = "testProvider"

  # Optional parameters when creating and updating. 
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

# After the execution of above resource block, user would have been created on the PowerScale array. 
# For more information, Please check the terraform state file. 