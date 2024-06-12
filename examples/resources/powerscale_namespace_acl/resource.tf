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
# If only namespace is provided, `terraform apply` will load Namespace ACL from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Namespace ACL (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting Namespace ACL from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale Namespace ACL allow you to manage the access control list for a namespace.
resource "powerscale_namespace_acl" "example_namespace_acl" {
  # Required and immutable once set
  namespace = "ifs/example"

  #   # Optional query parameters
  #   nsaccess = true
  #
  #   # Optional fields both for creating and updating
  #   # For owner and group, please provide either the UID/GID or the name+type
  #   owner = {
  #     id = "UID:0"
  #   }
  #   group = {
  #     name = "Isilon Users",
  #     type = "group"
  #   }
  #
  #   # acl_custom is required for updating. It can be set to [] to remove all acl.
  #   # While creating, if owner or group is provided, acl_custom must be specified as well. If none of
  #   # the three parameters are provided, Terraform will load the settings from the array directly.
  #   # For trustee, please provide either the UID/GID or the name+type
  #   # Please notice, the field acl_custom is the raw configuration, PowerScale will identify and calculate the accessrights
  #   # and inherit_flags provided and return its effective settings, which will be represented in the field acl in state.
  #   acl_custom = [
  #     {
  #       accessrights  = ["dir_gen_all"]
  #       accesstype    = "allow"
  #       inherit_flags = ["container_inherit"]
  #       trustee = {
  #         id = "UID:0"
  #       }
  #     },
  #     {
  #       accessrights  = ["dir_gen_write", "dir_gen_read", "dir_gen_execute", "std_read_dac"]
  #       accesstype    = "allow"
  #       inherit_flags = ["container_inherit"]
  #       trustee = {
  #         name = "Isilon Users",
  #         type = "group"
  #       }
  #     },
  #   ]
}

# After the execution of above resource block, Namespace ACL would have been cached in terraform state file, or
# Namespace ACL would have been updated on PowerScale.
# For more information, Please check the terraform state file.

# Example: A comparison of accessrights between acl_custom (user raw configuration) and acl (array effective settings)
# |                  acl_custom                   |        acl        |    note   |
# |-----------------------------------------------|-------------------|-----------|
# |          "dir_gen_all","dir_gen_read"         |   "dir_gen_all"   |     1     |
# |                 "generic_read"                |   "dir_gen_read"  |     2     |
# | "traverse", "std_read_dac", "std_synchronize" | "dir_gen_execute" |     3     |
#
# Note:
# 1. "dir_gen_read" will be merged into "dir_gen_all"
# 2. "generic_read" will be converted to "dir_gen_read" for a directory
# 3. "traverse", "std_read_dac", "std_synchronize" will be combined to "dir_gen_execute"
# These examples are only typical cases for your reference. Not all cases are listed.