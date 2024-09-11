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

# Available actions: Create, Read, Update, Delete and Import. Delete will clear the the state file.
# After `terraform apply` of this example file will update the settings according to the attributes set in the config.
# All attributes are optional. This resource expects atleast one attribute to be present in the config.

resource "powerscale_support_assist" "test" {
  supportassist_enabled   = true
  enable_download         = true
  automatic_case_creation = false
  enable_remote_support   = true
  accepted_terms          = true
  access_key              = "key"
  pin                     = "pin"
  telemetry = {
    offline_collection_period = 7200,
    telemetry_enabled         = true,
    telemetry_persist         = true,
    telemetry_threads         = 6
  }
  contact = {
    primary = {
      email      = "abc@gmail.com",
      first_name = "terraform_first",
      language   = "En",
      last_name  = "terraform_last",
      phone      = "1234567890"
    },
    secondary = {
      email      = "xyz@gmail.com",
      first_name = "terraform_second",
      language   = "No",
      last_name  = "terraform",
      phone      = "1234567890"
    }
  }
  connections = {
    mode = "gateway"
    gateway_endpoints = [
      {
        enabled      = true,
        host         = "1.2.3.4",
        port         = 9443,
        priority     = 1,
        use_proxy    = true,
        validate_ssl = true
      },
    ],
    network_pools = ["subnet0:pool0"]
  }
}