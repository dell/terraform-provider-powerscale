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

# Available actions: Create, Update, Delete and Import.
# After `terraform apply` of this example file it will create a new groupnet with the name set in `name` attribute on the PowerScale.

# PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.
resource "powerscale_groupnet" "example_groupnet" {
  # Required name for creating and updating. (Update Supported)
  name = "testGroupnetResourceSample"

  # Optional fields when creating and updating.

  # A description of the groupnet. (Update Supported)
  # description = "description of the groupnet"

  # DNS caching is enabled or disabled. Defaults to True. (Update Supported)
  # dns_cache_enabled = true

  # If enabled, SmartConnect treats subdomains of known dns zones as the known dns zone. This is required for S3 Virtual Host domains. Defaults to True. (Update Supported)
  # allow_wildcard_subdomains = true

  # Enable or disable appending nodes DNS search list to client DNS inquiries directed at SmartConnect service IP. Defaults to True. (Update Supported)
  # server_side_dns_search = true

  # Enable or disable DNS resolver rotate. Defaults to False. (Update Supported)
  # dns_resolver_rotate = false

  # List of DNS search suffixes. (Update Supported)
  # dns_search = ["DNS search suffixes"]

  # List of Domain Name Server IP addresses. (Update Supported)
  # dns_servers = ["10.230.44.169"]
}

# After the execution of above resource block, groupnet would have been created on the PowerScale array. 
# For more information, Please check the terraform state file. 