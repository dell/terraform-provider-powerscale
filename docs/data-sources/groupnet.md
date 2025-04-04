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

title: "powerscale_groupnet data source"
linkTitle: "powerscale_groupnet"
page_title: "powerscale_groupnet Data Source - terraform-provider-powerscale"
subcategory: ""
description: |-
  This datasource is used to query the existing Groupnets from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.
---

# powerscale_groupnet (Data Source)

This datasource is used to query the existing Groupnets from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.

## Example Usage

```terraform
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

# PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.

# Returns a list of PowerScale Groupnets based on names filter block. 
data "powerscale_groupnet" "example_groupnet" {
  filter {
    # Optional list of names to filter upon
    names = ["groupnet_name"]
  }
}

# Returns a list of PowerScale Groupnets in order based on the filters in the filter block. 
data "powerscale_groupnet" "example_groupnet" {
  filter {
    sort = "name"
    dir  = "DESC"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_groupnet.example_groupnet
output "powerscale_groupnet_filter" {
  value = data.powerscale_groupnet.example_groupnet
}


# Returns all of the PowerScale Groupnets
data "powerscale_groupnet" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_groupnet.all
output "powerscale_groupnet_all" {
  value = data.powerscale_groupnet.all
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `groupnets` (Attributes List) List of groupnets. (see [below for nested schema](#nestedatt--groupnets))
- `id` (String) Unique identifier of the groupnet instance.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `dir` (String) The direction of the sort.
- `limit` (Number) Return no more than this many results.
- `names` (Set of String) Only list groupnet matching this name.
- `sort` (String) The field that will be used for sorting.


<a id="nestedatt--groupnets"></a>
### Nested Schema for `groupnets`

Read-Only:

- `allow_wildcard_subdomains` (Boolean) If enabled, SmartConnect treats subdomains of known dns zones as the known dns zone. This is required for S3 Virtual Host domains.
- `description` (String) A description of the groupnet.
- `dns_cache_enabled` (Boolean) DNS caching is enabled or disabled.
- `dns_resolver_rotate` (Boolean) Enable or disable DNS resolver rotate.
- `dns_search` (List of String) List of DNS search suffixes.
- `dns_servers` (List of String) List of Domain Name Server IP addresses.
- `id` (String) Unique Interface ID.
- `name` (String) The name of the groupnet.
- `server_side_dns_search` (Boolean) Enable or disable appending nodes DNS search list to client DNS inquiries directed at SmartConnect service IP.
- `subnets` (List of String) Name of the subnets in the groupnet.