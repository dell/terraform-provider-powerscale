---
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerscale_user_mapping_rules resource"
linkTitle: "powerscale_user_mapping_rules"
page_title: "powerscale_user_mapping_rules Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  This resource is used to manage the User Mapping Rules entity of PowerScale Array. PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules.We can Create, Update and Delete the User Mapping Rules using this resource. We can also import an existing User Mapping Rules from PowerScale array. Note that, User Mapping Rules is the native functionality of PowerScale. When creating the resource, we actually load User Mapping Rules from PowerScale to the resource state.
---

# powerscale_user_mapping_rules (Resource)

This resource is used to manage the User Mapping Rules entity of PowerScale Array. PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules.We can Create, Update and Delete the User Mapping Rules using this resource. We can also import an existing User Mapping Rules from PowerScale array. Note that, User Mapping Rules is the native functionality of PowerScale. When creating the resource, we actually load User Mapping Rules from PowerScale to the resource state.


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
# If resource arguments are omitted, `terraform apply` will load User Mapping Rules from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load User Mapping Rules (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting User Mapping Rules from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules.
resource "powerscale_user_mapping_rules" "testUserMappingRules" {

  # Optional params for updating.

  # The zone to which the user mapping applies. Defaults to System
  zone = "System"

  # Specifies the parameters for user mapping rules.
  parameters = {
    # Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information.
    # When default_unix_user is not null: Designate the user as the default UNIX user.
    # When default_unix_user is null: Allow Dell Technologies to generate a primary UID and GID.
    # When default_unix_user.user is " ": Deny access with the following error: no such user.
    default_unix_user = {
      # Specifies the domain of the user that is being mapped.
      domain = "domain",
      # Specifies the name of the user that is being mapped.
      user = "username"
    }
  }


  # Specifies the list of user mapping rules.
  rules = [
    {
      # Specifies the operator to make rules on specified users or groups. Acceptable values: append, insert, replace, trim, union.
      operator = "append",
      # Specifies the mapping options for this user mapping rule. 
      options = {
        # If true, and the rule was applied successfully, stop processing further.
        break = true,
        # Specifies the default user information that can be applied if the final credentials do not have valid UID and GID information. 
        default_user = {
          # Specifies the domain of the user that is being mapped.
          domain = "domain",
          # Specifies the name of the user that is being mapped.
          user = "Guest"
        },
        # If true, the primary GID and primary group SID should be copied to the existing credential.
        group = true,
        # If true, all additional identifiers should be copied to the existing credential.
        groups = true,
        # If true, the primary UID and primary user SID should be copied to the existing credential.
        user = true
      },
      # Specifies the target user information that the rule can be applied to.
      target_user = {
        domain = "domain",
        user   = "testMappingRule"
      },
      # Specifies the source user information that the rule can be applied from.
      source_user = {
        domain = "domain",
        user   = "Guest"
      }
    },
    {
      # Operator 'trim' only accepts 'break' option and only accepts a single user.
      operator = "trim",
      options = {
        break = true,
      },
      target_user = {
        domain = "domain",
        user   = "testMappingRule"
      }
    },
    {
      # Operator 'union' only accepts 'break' option.
      operator = "union",
      options = {
        break = true,
        default_user = {
          domain = "domain",
          user   = "Guest"
        },
      },
      target_user = {
        user = "tfaccUserMappungRuleUser"
      },
      source_user = {
        user = "admin"
      }
    },
    {
      # Operator 'replace' only accepts 'break' option.
      operator = "replace",
      options = {
        break = true,
        default_user = {
          domain = "domain",
          user   = "Guest"
        },
      },
      target_user = {
        domain = "domain",
        user   = "tfaccUserMappungRuleUser"
      },
      source_user = {
        domain = "domain",
        user   = "admin"
      }
    },
    {
      operator = "insert",
      options = {
        break = true,
        default_user = {
          domain = "domain",
          user   = "Guest"
        },
        group  = true,
        groups = true,
        user   = true
      },
      target_user = {
        domain = "domain",
        user   = "tfaccUserMappungRuleUser"
      },
      source_user = {
        domain = "domain",
        user   = "admin"
      }
    },
  ]

  # List of user identity for mapping test.
  test_mapping_users = [
    {
      # Specifies a user name.
      name = "root"
      # Specifies a numeric user identifier.
      uid = 0
    },
    {
      name = "admin"
    }
  ]
}

# After the execution of above resource block, User Mapping Rules would have been cached in terraform state file, or
# User Mapping Rules would have been updated on PowerScale.
# For more information, Please check the terraform state file.
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `parameters` (Attributes) Specifies the parameters for user mapping rules. (Update Supported) (see [below for nested schema](#nestedatt--parameters))
- `rules` (Attributes List) Specifies the list of user mapping rules. (Update Supported) (see [below for nested schema](#nestedatt--rules))
- `test_mapping_users` (Attributes List) List of user identity for mapping test. (Update Supported) (see [below for nested schema](#nestedatt--test_mapping_users))
- `zone` (String) The zone to which the user mapping applies. (Update Supported)

### Read-Only

- `id` (String) User Mapping Rules ID.
- `mapping_users` (Attributes List) List of test mapping user result. (see [below for nested schema](#nestedatt--mapping_users))

<a id="nestedatt--parameters"></a>
### Nested Schema for `parameters`

Optional:

- `default_unix_user` (Attributes) Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information. (Update Supported) (see [below for nested schema](#nestedatt--parameters--default_unix_user))

<a id="nestedatt--parameters--default_unix_user"></a>
### Nested Schema for `parameters.default_unix_user`

Required:

- `user` (String) Specifies the name of the user that is being mapped. (Update Supported)

Optional:

- `domain` (String) Specifies the domain of the user that is being mapped. (Update Supported)



<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Required:

- `operator` (String) Specifies the operator to make rules on specified users or groups. (Update Supported)
- `target_user` (Attributes) Specifies the target user information that the rule can be applied to. (Update Supported) (see [below for nested schema](#nestedatt--rules--target_user))

Optional:

- `options` (Attributes) Specifies the mapping options for this user mapping rule. (Update Supported) (see [below for nested schema](#nestedatt--rules--options))
- `source_user` (Attributes) Specifies the source user information that the rule can be applied from. (Update Supported) (see [below for nested schema](#nestedatt--rules--source_user))

<a id="nestedatt--rules--target_user"></a>
### Nested Schema for `rules.target_user`

Required:

- `user` (String) Specifies the name of the user that is being mapped. (Update Supported)

Optional:

- `domain` (String) Specifies the domain of the user that is being mapped. (Update Supported)


<a id="nestedatt--rules--options"></a>
### Nested Schema for `rules.options`

Optional:

- `break` (Boolean) If true, and the rule was applied successfully, stop processing further. (Update Supported)
- `default_user` (Attributes) Specifies the default user information that can be applied if the final credentials do not have valid UID and GID information. (Update Supported) (see [below for nested schema](#nestedatt--rules--options--default_user))
- `group` (Boolean) If true, the primary GID and primary group SID should be copied to the existing credential. (Update Supported)
- `groups` (Boolean) If true, all additional identifiers should be copied to the existing credential. (Update Supported)
- `user` (Boolean) If true, the primary UID and primary user SID should be copied to the existing credential. (Update Supported)

<a id="nestedatt--rules--options--default_user"></a>
### Nested Schema for `rules.options.default_user`

Required:

- `user` (String) Specifies the name of the user that is being mapped. (Update Supported)

Optional:

- `domain` (String) Specifies the domain of the user that is being mapped. (Update Supported)



<a id="nestedatt--rules--source_user"></a>
### Nested Schema for `rules.source_user`

Required:

- `user` (String) Specifies the name of the user that is being mapped. (Update Supported)

Optional:

- `domain` (String) Specifies the domain of the user that is being mapped. (Update Supported)



<a id="nestedatt--test_mapping_users"></a>
### Nested Schema for `test_mapping_users`

Optional:

- `name` (String) Specifies a user name. (Update Supported)
- `uid` (Number) Specifies a numeric user identifier. (Update Supported)


<a id="nestedatt--mapping_users"></a>
### Nested Schema for `mapping_users`

Read-Only:

- `privileges` (Attributes List) Specifies the system-defined privilege that may be granted to users. (see [below for nested schema](#nestedatt--mapping_users--privileges))
- `supplemental_identities` (Attributes List) Specifies the configuration properties for a user. (see [below for nested schema](#nestedatt--mapping_users--supplemental_identities))
- `user` (Attributes) Specifies the configuration properties for a user. (see [below for nested schema](#nestedatt--mapping_users--user))
- `zid` (Number) Numeric ID of the access zone which contains this user.
- `zone` (String) Name of the access zone which contains this user.

<a id="nestedatt--mapping_users--privileges"></a>
### Nested Schema for `mapping_users.privileges`

Read-Only:

- `id` (String) Specifies the ID of the privilege.
- `name` (String) Specifies the name of the privilege.
- `read_only` (Boolean) True, if the privilege is read-only.


<a id="nestedatt--mapping_users--supplemental_identities"></a>
### Nested Schema for `mapping_users.supplemental_identities`

Read-Only:

- `gid` (String) Specifies a user or group GID.
- `name` (String) Specifies a user or group name.
- `sid` (String) Specifies a user or group SID.


<a id="nestedatt--mapping_users--user"></a>
### Nested Schema for `mapping_users.user`

Read-Only:

- `name` (String) Specifies the user name.
- `on_disk_user_identity` (String) Specifies the user identity on disk.
- `primary_group_name` (String) Specifies the primary group name.
- `primary_group_sid` (String) Specifies the primary group SID.
- `sid` (String) Specifies a user or group SID.
- `uid` (String) Specifies the user UID.

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
# terraform import powerscale_user_mapping_rules.testUserMappingRules <zoneName>
# Example:
terraform import powerscale_user_mapping_rules.testUserMappingRules System
# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
```