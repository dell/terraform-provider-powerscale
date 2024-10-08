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

title: "powerscale_accesszone data source"
linkTitle: "powerscale_accesszone"
page_title: "powerscale_accesszone Data Source - terraform-provider-powerscale"
subcategory: ""
description: |-
  Access Zone Datasource. This datasource is used to query the existing Access Zone from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Access Zones allow you to isolate data and control who can access data in each zone.
---

# powerscale_accesszone (Data Source)

Access Zone Datasource. This datasource is used to query the existing Access Zone from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Access Zones allow you to isolate data and control who can access data in each zone.

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

# Returns all of the PowerScale access zones and their details
# PowerScale access zones allow you to isolate data and control who can access data in each zone.
data "powerscale_accesszone" "all" {
}

output "powerscale_accesszone_data_all" {
  value = data.powerscale_accesszone.all
}

# Returns a subset of the PowerScale access zones based on the names provided in the `names` filter block and their details
data "powerscale_accesszone" "test" {
  # Optional list of names to filter upon
  filter {
    names = ["tfaccAccessZone"]
  }
}

output "powerscale_accesszone" {
  value = data.powerscale_accesszone.test
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
# Also, we can use the fetched information by the variable data.powerscale_accesszone.all
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block, Optional) (see [below for nested schema](#nestedblock--filter))

### Read-Only

- `access_zones_details` (Attributes List) List of AccessZones (see [below for nested schema](#nestedatt--access_zones_details))
- `id` (String) Identifier

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Optional:

- `names` (Set of String)


<a id="nestedatt--access_zones_details"></a>
### Nested Schema for `access_zones_details`

Read-Only:

- `alternate_system_provider` (String) Specifies an alternate system provider.
- `auth_providers` (List of String) Specifies the list of authentication providers available on this access zone.
- `cache_entry_expiry` (Number) Specifies amount of time in seconds to cache a user/group.
- `create_path` (Boolean) Determines if a path is created when a path does not exist.
- `groupnet` (String) Groupnet identifier
- `home_directory_umask` (Number) Specifies the permissions set on automatically created user home directories.
- `id` (String) Specifies the system-assigned ID for the access zone. This value is returned when an access zone is created through the POST method
- `ifs_restricted` (Attributes List) Specifies a list of users and groups that have read and write access to /ifs. (see [below for nested schema](#nestedatt--access_zones_details--ifs_restricted))
- `map_untrusted` (String) Maps untrusted domains to this NetBIOS domain during authentication.
- `name` (String) Specifies the access zone name.
- `negative_cache_entry_expiry` (Number) Specifies number of seconds the negative cache entry is valid.
- `netbios_name` (String) Specifies the NetBIOS name.
- `path` (String) Specifies the access zone base directory path.
- `skeleton_directory` (String) Specifies the skeleton directory that is used for user home directories.
- `system` (Boolean) True if the access zone is built-in.
- `system_provider` (String) Specifies the system provider for the access zone.
- `user_mapping_rules` (List of String) Specifies the current ID mapping rules.
- `zone_id` (Number) Specifies the access zone ID on the system.

<a id="nestedatt--access_zones_details--ifs_restricted"></a>
### Nested Schema for `access_zones_details.ifs_restricted`

Read-Only:

- `id` (String) Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.
- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.