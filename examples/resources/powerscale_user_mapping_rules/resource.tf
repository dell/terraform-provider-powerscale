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