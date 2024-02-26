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

# Available actions: Create, Update (owner, group, access_control), Delete and Import existing FileSystem(Namespace directory) from Powerscale array.
# After `terraform apply` of this example file it will create a new FileSystem(Namespace directory) with the name set in `name` attribute in the directory path provided in `directory_path`on the PowerScale array

# PowerScale FileSystem Resource allows you to manage the Namespace Directory on the Powerscale array
resource "powerscale_filesystem" "file_system_test" {
  # Default set to '/ifs'
  # directory_path         = "/ifs"

  # Required attributes
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

  # Optional : query_zone, this will default to the default access zone if unset. However is needed if the user trying to be created is not in the default access zone.connection {
  # This should just be the access zone name. 
  # query_zone = "test_access_zone"

  # Optional attributes. Default values set.
  # Creates intermediate folders recursively, when set to true.
  recursive = true
  # Deletes and replaces the existing user attributes and ACLs of the directory with user-specified attributes and ACLS, when set to true.
  overwrite = false


  /* Optional : The ACL value for the directory. Users can either provide access rights input such as 'private_read' , 'private' ,
    'public_read', 'public_read_write', 'public' or permissions in POSIX format as '0550', '0770', '0775','0777' or 0700. The Default value is (0700). 
     Modification of ACL is only supported from POSIX to POSIX mode. 
  */

  # access_control = "0777"
}
# After the execution of above resource block, a PowerScale FileSystem(Namespace directory) would have been created at PowerScale array. You can also verify the changes made in terraform state file.
