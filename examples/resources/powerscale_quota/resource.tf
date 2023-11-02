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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will create Quota on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale Quota module monitors and enforces administrator-defined storage limits.

resource "powerscale_quota" "quota_test" {
  # Required and update not supported
  path              = "/ifs/example_quota"
  type              = "user"
  include_snapshots = "false"

  # Optional and update not supported
  # zone = "System"
  # persona = {
  # id = "UID:1501"
  # name = "Guest"
  # type = "user"
  # }

  # Optional and update supported
  # container = true
  # enforced = false
  # force = false
  # thresholds_on = "applogicalsize"
  # ignore_limit_checks = true
  # thresholds = {
  # advisory = 1000
  # soft = 2000
  # hard = 4000
  # soft_grace = 120
  # }
}

# After the execution of above resource block, a Quota would have been created on the PowerScale array.
# For more information, Please check the terraform state file.