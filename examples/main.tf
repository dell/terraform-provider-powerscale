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

variable "username" {
  type = string
}

variable "password" {
  type = string
}

variable "endpoint" {
  type = string
}

variable "insecure" {
  type = bool
}

variable "group" {
  type = string
}

variable "volume_path" {
  type = string
}

variable "volume_path_permissions" {
  type = string
}

variable "ignore_unresolvable_hosts" {
  type = bool
}

variable "auth_type" {
  type = number
}

variable "verbose_logging" {
  type = number
}

provider "powerscale" {
  username                  = var.username
  password                  = var.password
  endpoint                  = var.endpoint
  insecure                  = var.insecure
  group                     = var.group
  volume_path               = var.volume_path
  volume_path_permissions   = var.volume_path_permissions
  ignore_unresolvable_hosts = var.ignore_unresolvable_hosts
  auth_type                 = var.auth_type
  verbose_logging           = var.verbose_logging
}
