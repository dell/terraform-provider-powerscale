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

# Available actions: Create, Update, Delete and Import.
# After `terraform apply` of this example file will create the nfs alias on the PowerScale array with the attributes set in the config.
# For update, only name and path is supported. Zone Cannot be updated.

resource "powerscale_nfs_alias" "example" {
  # Required field both for creating and updating
  name = "/exampleNfs"
  path = "/ifs"

  # Optional attribute for creating and updating
  zone = "System"
}