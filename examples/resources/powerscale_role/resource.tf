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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file for the first time, you will create a role on the PowerScale

# PowerScale role allows you to permit and limit access to administrative areas of your cluster on a per-user basis through roles.
resource "powerscale_role" "role_test" {
  # Required
  name = "role_test"

  # Optional fields only for creating
  zone = "System"

  # Optional fields both for creating and updating
  description = "role_test_description"
  # To add members, please provide uid/gid or provide name and type
  members = [
    {
      name = "admin",
      type = "user"
    },
    {
      id = "UID:0"
    },
    {
      name = "guest",
      type = "group"
    }
  ]
  # To add privileges, the id is required. Please use role privilege datasource to look up the role privilege id needed.
  privileges = [
    {
      id         = "ISI_PRIV_SYS_SUPPORT",
      permission = "r"
    },
    {
      id         = "ISI_PRIV_SYS_SHUTDOWN",
      permission = "r"
    }
  ]
}

# After the execution of above resource block, role would have been created on the PowerScale array.
# For more information, Please check the terraform state file.