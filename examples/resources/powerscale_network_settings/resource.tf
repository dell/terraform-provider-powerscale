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

# Available actions: Create, Update, Delete and Import.
# If resource arguments are omitted, `terraform apply` will load Network Settings from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Network Settings (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting Network Settings from PowerScale.
# For more information, Please check the terraform state file.

# PowerScale Network Settings provide the ability to configure external network configuration on the cluster.
resource "powerscale_network_settings" "example_network_settings" {

  # Optional fields when updating.

  # Enable or disable Source Based Routing. (Update Supported)
  # source_based_routing_enabled = false

  # Delay in seconds for IP rebalance. (Update Supported)
  # sc_rebalance_delay = 0

  # List of client TCP ports. (Update Supported)
  # tcp_ports = [20, 21, 80, 445, 2049]
}

# After the execution of above resource block, Network Settings would have been cached in terraform state file, or
# Network Settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.