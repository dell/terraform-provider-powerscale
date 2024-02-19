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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file for the first time, you will create a network rule on the PowerScale

# PowerScale network rule allows you to manage the network rule on the Powerscale array
resource "powerscale_network_rule" "rule" {
  # Required. Name of the provisioning rule.
  name = "tfacc_rule"
  # Required.
  groupnet = "groupnet0"
  # Required.
  subnet = "subnet0"
  # Required.
  pool = "pool0"
  # Required. Interface name the provisioning rule applies to.
  iface = "ext-2"
  # Optional. Description for the provisioning rule.
  description = "tfacc_rule"
  # Optional. Node type the provisioning rule applies to.
  node_type = "any"
}

# After the execution of above resource block, network rule would have been created on the PowerScale array.
# For more information, Please check the terraform state file.