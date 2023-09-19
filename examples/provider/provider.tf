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

terraform {
  required_providers {
    powerscale = {
      source = "registry.terraform.io/dell/powerscale"
    }
  }
}

provider "powerscale" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = var.insecure
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
