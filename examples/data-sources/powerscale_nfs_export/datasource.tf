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

data "powerscale_nfs_export" "test" {
  filter = {
    ids   = [1, 2, 3]
    paths = ["/ifs/primary", "/ifs/secondary"]
    # check = true
    # dir   = ASC
    # limit = 10
    # offset= 1
    # path  = "/ifs/primary"
    # resume= "resume_token"
    # scope = "user"
    # sort  = "id"
    # zone  = "System"
  }
}

output "powerscale_nfs_export" {
  value = data.powerscale_nfs_export.test
}

data "powerscale_nfs_export" "all" {
}

output "powerscale_nfs_export" {
  value = data.powerscale_nfs_export.all
}