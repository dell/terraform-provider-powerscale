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
resource "powerscale_filesystem" "file_system_test" {
  # Default set to '/ifs'
  # directory_path         = "/ifs"

  # Required
  name = "DirTf"
  group = {
    id   = "GID:0"
    name = "wheel"
    type = "group"
  }
  owner = {
    id   = "UID:1501",
    name = "Guest",
    type = "user"
  }

  # Optional. Default values set.
  recursive = true
  overwrite = false


  /* Optional : The ACL value for the directory. Users can either provide access rights input such as 'private_read' , 'private' ,
    'public_read', 'public_read_write', 'public' or permissions in POSIX format as '0550', '0770', '0775','0777' or 0700. The Default value is (0700). 
     Modification of ACL is only supported from POSIX to POSIX mode. 
  */

  # access_control = "0777"
}
