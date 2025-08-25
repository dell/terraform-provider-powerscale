---
# Copyright (c) 2023-2025 Dell Inc., or its subsidiaries. All Rights Reserved.
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

listIgnoreTitle: true
weight: 1
title: "powerscale provider"
linkTitle: "Provider"
page_title: "powerscale Provider"
subcategory: ""
description: |-
  The Terraform provider for Dell PowerScale can be used to interact with a Dell PowerScale array in order to manage the array resources.
---

# powerscale Provider

The Terraform provider for Dell PowerScale can be used to interact with a Dell PowerScale array in order to manage the array resources.

## Example Usage
The following abridged example demonstrates the usage of the provider to create groupnet, subnet, network pool, ads provider, access zone, quota, snapshot, snapshot schedule, user, user group, filesystem ,nfs export and smb share.

```terraform
/*
Copyright (c) 2023-2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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

terraform {
  required_providers {
    powerscale = {
      source  = "registry.terraform.io/dell/powerscale"
      version = "1.8.0"
    }
  }
}

provider "powerscale" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = var.insecure

  ## Provider can also be set using environment variables
  ## If environment variables are set it will override this configuration
  ## Example environment variables
  # POWERSCALE_USERNAME="username"
  # POWERSCALE_PASSWORD="password"
  # POWERSCALE_ENDPOINT="https://yourhost.host.com:8080"
  # POWERSCALE_INSECURE="false"
  # POWERSCALE_TIMEOUT="2000"
  # POWERSCALE_AUTH_TYPE="0"
}

# Specify 'alias' value in provider section to use multiple providers.
provider "powerscale" {
  alias    = "secondProvider"
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = var.insecure
}

# Specify 'provider' value to use the second provider in resource section.
resource "powerscale_groupnet" "example_multiple_provider_groupnet" {
  provider = powerscale.secondProvider
  name     = "example_groupnet"
}

resource "powerscale_groupnet" "example_groupnet" {
  name = "example_groupnet"
}

resource "powerscale_accesszone" "zone" {
  name     = "example_acczone"
  groupnet = powerscale_groupnet.example_groupnet.name
  path     = "/ifs"
}

resource "powerscale_adsprovider" "ads_test" {
  name     = "ADS.PROVIDER.EXAMPLE.COM"
  groupnet = powerscale_groupnet.example_groupnet.name
  user     = "admin"
  password = "password"
}

resource "powerscale_subnet" "subnet" {
  name     = "example_subnet"
  groupnet = powerscale_groupnet.example_groupnet.name
}

resource "powerscale_networkpool" "pool_test" {
  name        = "example_pool"
  subnet      = powerscale_subnet.subnet.name
  groupnet    = powerscale_groupnet.example_groupnet.name
  access_zone = powerscale_accesszone.zone.name
}

resource "powerscale_quota" "quota_test" {
  path              = powerscale_filesystem.example_file_system.full_path
  type              = "user"
  include_snapshots = "false"
  zone              = powerscale_accesszone.zone.name
  persona = {
    id   = format("%s:%s", "UID", powerscale_user.example_user.uid)
    name = powerscale_user.example_user.name
    type = "user"
  }
}

resource "powerscale_snapshot" "snap" {
  path        = powerscale_filesystem.example_file_system.full_path
  name        = "example_snapshot"
  set_expires = "1 Day"
}

resource "powerscale_snapshot_schedule" "snap_schedule" {
  name = "example_snap_schedule"
  path = powerscale_filesystem.example_file_system.full_path
}

resource "powerscale_user" "example_user" {
  name    = "example_user"
  enabled = true
}

resource "powerscale_user" "testUser1" {
  name       = "example_user1"
  query_zone = powerscale_accesszone.zone.name
}

resource "powerscale_user" "testUser2" {
  name       = "example_user2"
  query_zone = powerscale_accesszone.zone.name
}

resource "powerscale_user_group" "example_user_group" {
  name  = "example_user_group"
  users = [powerscale_user.example_user.name]
}

resource "powerscale_filesystem" "example_file_system" {
  directory_path = "/ifs/data"
  name           = "example_file_system"
  group = {
    id   = format("%s:%s", "GID", powerscale_user_group.example_user_group.gid)
    name = powerscale_user_group.example_user_group.name
    type = "group"
  }
  owner = {
    id   = format("%s:%s", "UID", powerscale_user.example_user.uid)
    name = powerscale_user.example_user.name,
    type = "user"
  }
  access_control = "public_read_write"
}

resource "powerscale_nfs_export" "example_export" {
  paths = [powerscale_filesystem.example_file_system.full_path]
  zone  = powerscale_accesszone.zone.name
  map_all = {
    enabled = true,
    primary_group = {
      id = format("%s:%s", "GROUP", powerscale_user_group.example_user_group.gid)
    }
    user = {
      id = format("%s:%s", "USER", powerscale_user.example_user.uid)
    }
  }
}

resource "powerscale_smb_share" "share_example" {
  name = "example_smb_share"
  path = powerscale_filesystem.example_file_system.full_path
  zone = powerscale_accesszone.zone.name
  permissions = [
    {
      permission      = "full"
      permission_type = "allow"
      trustee = {
        id   = powerscale_user.example_user.sid,
        name = powerscale_user.example_user.name,
        type = "user"
      }
    }
  ]
}

resource "powerscale_role" "role_test" {
  name        = "role_test"
  zone        = powerscale_accesszone.zone.name
  description = "tfacc for role test"
  members = [
    {
      id = format("%s:%s", "UID", powerscale_user.testUser1.uid)
    },
  ]
  privileges = [
    {
      id         = "ISI_PRIV_LOGIN_PAPI",
      permission = "r"
    }
  ]
}

resource "powerscale_nfs_export_settings" "example" {
  map_failure = {
    enabled       = false
    primary_group = {}
    user = {
      id = format("%s:%s", "UID", powerscale_user.testUser1.uid)
    }
  }
  map_root = {
    enabled       = true
    primary_group = {}
    user = {
      id = format("%s:%s", "UID", powerscale_user.testUser1.uid)
    }
  }
  zone = powerscale_accesszone.zone.name
}

resource "powerscale_user_mapping_rules" "testUserMappingRules" {
  zone = powerscale_accesszone.zone.name
  parameters = {
    default_unix_user = {
      user = "Guest"
    }
  }
  rules = [
    {
      operator = "insert",
      options = {
        break = true,
        default_user = {
          user = "Guest"
        },
        group  = true,
        groups = true,
        user   = true
      },
      target_user = {
        user = powerscale_user.testUser2.name
      },
      source_user = {
        user = powerscale_user.testUser1.name
      }
    },
  ]
  test_mapping_users = [
    {
      name = powerscale_user.testUser1.name
    },
    {
      name = powerscale_user.testUser2.name
    }
  ]
}

resource "powerscale_namespace_acl" "example_namespace_acl" {
  namespace = powerscale_filesystem.example_file_system.id
  acl_custom = [
    {
      accessrights = [
        "dir_gen_read",
        "dir_gen_write",
        "dir_gen_execute",
        "std_write_dac",
        "delete_child",
      ]
      accesstype    = "allow"
      inherit_flags = []
      trustee = {
        id = format("%s:%s", "UID", powerscale_user.example_user.uid)
      }
    },
    {
      accessrights = [
        "dir_gen_read",
        "dir_gen_write",
        "dir_gen_execute",
        "delete_child",
      ]
      accesstype    = "allow"
      inherit_flags = []
      trustee = {
        id = format("%s:%s", "GID", powerscale_user_group.example_user_group.gid)
      }
    },
  ]
}

resource "powerscale_s3_bucket" "s3_bucket_example" {
  name        = "s3-bucket-example"
  path        = powerscale_filesystem.example_file_system.full_path
  create_path = false
  owner       = "Guest"
  zone        = powerscale_accesszone.zone.name
  acl = [{
    grantee = {
      name = powerscale_user.example_user.name
      type = "user"
    }
    permission = "FULL_CONTROL"
  }]
  description       = "tfacc-s3-bucket-test creation"
  object_acl_policy = "replace"
}

resource "powerscale_smb_share_settings" "example" {
  access_based_enumeration           = true
  access_based_enumeration_root_only = true
  allow_delete_readonly              = false
  ca_timeout                         = 60
  zone                               = powerscale_accesszone.zone.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `auth_type` (Number) what should be the auth type, 0 for basic and 1 for session-based. This can also be set using the environment variable POWERSCALE_AUTH_TYPE
- `endpoint` (String) The API endpoint, ex. https://172.17.177.230:8080. This can also be set using the environment variable POWERSCALE_ENDPOINT
- `insecure` (Boolean) whether to skip SSL validation. This can also be set using the environment variable POWERSCALE_INSECURE
- `password` (String, Sensitive) The password. This can also be set using the environment variable POWERSCALE_PASSWORD
- `timeout` (Number) specifies a time limit for requests. This can also be set using the environment variable POWERSCALE_TIMEOUT
- `username` (String) The username. This can also be set using the environment variable POWERSCALE_USERNAME

## Best Practices
1. The parent resource attributes of a certain resource (e.g. groupnet field of subnet resource) can only be designated
   at creation. Once designated, they cannot be modified except for parent resource renaming.
2. The name of a resource is modifiable, but it is necessary to make sure its name referenced in the child resources
   is also updated (can be done manually or use reference resource_id.name).
3. Resources with child resources cannot be deleted independently. Use terraform destroy to delete all resources directly
   or delete all the child resources at the same time (depends_on is recommended to manage resources, serving as a
   precheck for delete operations).
