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

title: "powerscale_quota resource"
linkTitle: "powerscale_quota"
page_title: "powerscale_quota Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  This resource is used to manage the Quota entity of PowerScale Array. Quota module monitors and enforces administrator-defined storage limits. We can Create, Update and Delete the Quota using this resource. We can also import an existing Quota from PowerScale array.
---

# powerscale_quota (Resource)

This resource is used to manage the Quota entity of PowerScale Array. Quota module monitors and enforces administrator-defined storage limits. We can Create, Update and Delete the Quota using this resource. We can also import an existing Quota from PowerScale array.


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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will create Quota on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale Quota module monitors and enforces administrator-defined storage limits.

resource "powerscale_quota" "quota_test" {
  # Required and update not supported
  path              = "/ifs/example_quota"
  type              = "user"
  include_snapshots = "false"

  # Optional and update not supported
  # zone = "System"
  # persona = {
  # id = "UID:1501"
  # name = "Guest"
  # type = "user"
  # }

  # Optional and update supported
  # container = true
  # enforced = false
  # force = false
  # thresholds_on = "applogicalsize"
  # ignore_limit_checks = true
  # thresholds = {
  # advisory = 1000
  # soft = 2000
  # hard = 4000
  # soft_grace = 120
  # }
}

# After the execution of above resource block, a Quota would have been created on the PowerScale array.
# For more information, Please check the terraform state file.
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `include_snapshots` (Boolean) If true, quota governs snapshot data as well as head data.
- `path` (String) The ifs path governed.
- `type` (String) The type of quota.

### Optional

- `container` (Boolean) If true, quotas using the quota directory see the quota thresholds as share size.
- `enforced` (Boolean) True if the quota provides enforcement, otherwise an accounting quota.
- `force` (Boolean) Force creation of quotas on the root of /ifs or percent based quotas.
- `ignore_limit_checks` (Boolean) If true, skip child quota's threshold comparison with parent quota path.
- `linked` (Boolean) For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked. Set linked as true or false to link or unlink quota
- `persona` (Attributes) Specifies the persona of the file group. persona is required for user and group type. (see [below for nested schema](#nestedatt--persona))
- `thresholds` (Attributes) The thresholds of quota (see [below for nested schema](#nestedatt--thresholds))
- `thresholds_on` (String) Thresholds apply on quota accounting metric.
- `zone` (String) Optional named zone to use for user and group resolution.

### Read-Only

- `efficiency_ratio` (Number) Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.
- `id` (String) The system ID given to the quota.
- `notifications` (String) Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.
- `ready` (Boolean) True if the default resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `reduction_ratio` (Number) Represents the ratio of logical space provided to physical data space used. This accounts for compression and data deduplication effects.
- `usage` (Attributes) The usage of quota (see [below for nested schema](#nestedatt--usage))

<a id="nestedatt--persona"></a>
### Nested Schema for `persona`

Optional:

- `id` (String) Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.
- `name` (String) Specifies the persona name, which must be combined with a type.
- `type` (String) Specifies the type of persona, which must be combined with a name.


<a id="nestedatt--thresholds"></a>
### Nested Schema for `thresholds`

Optional:

- `advisory` (Number) Usage bytes at which notifications will be sent but writes will not be denied.
- `hard` (Number) Usage bytes at which further writes will be denied.
- `percent_advisory` (Number) Advisory threshold as percent of hard threshold. Usage bytes at which notifications will be sent but writes will not be denied.
- `percent_soft` (Number) Soft threshold as percent of hard threshold. Usage bytes at which notifications will be sent and soft grace time will be started.
- `soft` (Number) Usage bytes at which notifications will be sent and soft grace time will be started.
- `soft_grace` (Number) Time in seconds after which the soft threshold has been hit before writes will be denied.

Read-Only:

- `advisory_exceeded` (Boolean) True if the advisory threshold has been hit.
- `advisory_last_exceeded` (Number) Time at which advisory threshold was hit.
- `hard_exceeded` (Boolean) True if the hard threshold has been hit.
- `hard_last_exceeded` (Number) Time at which hard threshold was hit.
- `soft_exceeded` (Boolean) True if the soft threshold has been hit.
- `soft_last_exceeded` (Number) Time at which soft threshold was hit


<a id="nestedatt--usage"></a>
### Nested Schema for `usage`

Read-Only:

- `applogical` (Number) Bytes used by governed data apparent to application.
- `applogical_ready` (Boolean) True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `fslogical` (Number) Bytes used by governed data apparent to filesystem.
- `fslogical_ready` (Boolean) True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `fsphysical` (Number) Physical data usage adjusted to account for shadow store efficiency
- `fsphysical_ready` (Boolean) True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `inodes` (Number) Number of inodes (filesystem entities) used by governed data.
- `inodes_ready` (Boolean) True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `physical` (Number) Bytes used for governed data and filesystem overhead.
- `physical_data` (Number) Number of physical blocks for file data
- `physical_data_ready` (Boolean) True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `physical_protection` (Number) Number of physical blocks for file protection
- `physical_protection_ready` (Boolean) True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `physical_ready` (Boolean) True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
- `shadow_refs` (Number) Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.
- `shadow_refs_ready` (Boolean) True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.

Unless specified otherwise, all fields of this resource can be updated.

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The command is
# terraform import powerscale_quota.quota_example [<zoneID>]:<id>
# Example 1: <zoneID> is Optional, defaults to System:
terraform import powerscale_quota.quota_example example_quota_id
# Example 2:
terraform import powerscale_quota.quota_example zone_id:example_quota_id
# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
```