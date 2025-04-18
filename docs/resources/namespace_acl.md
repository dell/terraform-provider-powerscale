---
# Copyright (c) 2024-2025 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://mozilla.org/MPL/2.0/
#
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

title: "powerscale_namespace_acl resource"
linkTitle: "powerscale_namespace_acl"
page_title: "powerscale_namespace_acl Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  This resource is used to manage the Namespace ACL on PowerScale Array. We can Create, Update and Delete the Namespace ACL using this resource. We can also import the existing Namespace ACL from PowerScale array. Note that, when creating the resource, we actually load Namespace ACL from PowerScale to the resource state.
---

# powerscale_namespace_acl (Resource)

This resource is used to manage the Namespace ACL on PowerScale Array. We can Create, Update and Delete the Namespace ACL using this resource. We can also import the existing Namespace ACL from PowerScale array. Note that, when creating the resource, we actually load Namespace ACL from PowerScale to the resource state.


## Example Usage

```terraform
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
  #   zone = "system"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `namespace` (String) Indicate the namespace to set/get acl.

### Optional

- `acl_custom` (Attributes List) Customer's raw configuration of the JSON array of access rights. (see [below for nested schema](#nestedatt--acl_custom))
- `group` (Attributes) Provides the JSON object for the group persona of the owner. (see [below for nested schema](#nestedatt--group))
- `nsaccess` (Boolean) Indicates that the operation is on the access point instead of the store path.
- `owner` (Attributes) Provides the JSON object for the group persona of the owner. (see [below for nested schema](#nestedatt--owner))
- `zone` (String) Indicates the zone of the namespace.

### Read-Only

- `acl` (Attributes List) Array effective configuration of the JSON array of access rights. (see [below for nested schema](#nestedatt--acl))
- `authoritative` (String) If the directory has access rights set, then this field is returned as acl. If the directory has POSIX permissions set, then this field is returned as mode.
- `mode` (String) Provides the POSIX mode.

<a id="nestedatt--acl_custom"></a>
### Nested Schema for `acl_custom`

Required:

- `accesstype` (String) Grants or denies access control permissions. Options: allow, deny
- `trustee` (Attributes) Provides the JSON object for the group persona of the owner. (see [below for nested schema](#nestedatt--acl_custom--trustee))

Optional:

- `accessrights` (List of String) Specifies the access control permissions for a specific user or group. Options: std_delete, std_read_dac, std_write_dac, std_write_owner, std_synchronize, std_required, generic_all, generic_read, generic_write, generic_exec, dir_gen_all, dir_gen_read, dir_gen_write, dir_gen_execute, file_gen_all, file_gen_read, file_gen_write, file_gen_execute, modify, file_read, file_write, append, execute, file_read_attr, file_write_attr, file_read_ext_attr, file_write_ext_attr, delete_child, list, add_file, add_subdir, traverse, dir_read_attr, dir_write_attr, dir_read_ext_attr, dir_write_ext_attr
- `inherit_flags` (List of String) Grants or denies access control permissions. Options: object_inherit, container_inherit, inherit_only, no_prop_inherit, inherited_ace
- `op` (String) Operations for updating access control permissions. Unnecessary for access right replacing scenario

<a id="nestedatt--acl_custom--trustee"></a>
### Nested Schema for `acl_custom.trustee`

Optional:

- `id` (String) Specifies the serialized form of a persona, which can be 'UID:0' or 'GID:0'
- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.



<a id="nestedatt--group"></a>
### Nested Schema for `group`

Optional:

- `id` (String) Specifies the serialized form of a persona, which can be 'GID:0'
- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.


<a id="nestedatt--owner"></a>
### Nested Schema for `owner`

Optional:

- `id` (String) Specifies the serialized form of a persona, which can be 'UID:0'
- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.


<a id="nestedatt--acl"></a>
### Nested Schema for `acl`

Read-Only:

- `accessrights` (List of String) Specifies the access control permissions for a specific user or group.
- `accesstype` (String) Grants or denies access control permissions.
- `inherit_flags` (List of String) Grants or denies access control permissions.
- `op` (String) Operations for updating access control permissions. Unnecessary for access right replacing scenario
- `trustee` (Attributes) Provides the JSON object for the group persona of the owner. (see [below for nested schema](#nestedatt--acl--trustee))

<a id="nestedatt--acl--trustee"></a>
### Nested Schema for `acl.trustee`

Read-Only:

- `id` (String) Specifies the serialized form of a persona, which can be 'UID:0' or 'GID:0'
- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.

Unless specified otherwise, all fields of this resource can be updated.

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The command is
# terraform import powerscale_namespace_acl.namespace_acl_test <namespace_path>
# Example:
terraform import powerscale_namespace_acl.namespace_acl_test namespace_path
# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
```