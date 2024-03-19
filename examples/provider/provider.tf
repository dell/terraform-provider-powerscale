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

terraform {
  required_providers {
    powerscale = {
      source  = "registry.terraform.io/dell/powerscale"
      version = "1.1.0"
    }
  }
}

provider "powerscale" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = var.insecure
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