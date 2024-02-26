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

# The SmartQuotas module is an optional quota-management tool that monitors and enforces administrator-defined storage limits.

# Returns a list of PowerScale quotas based on filter block
data "powerscale_quota" "example_quotas" {
  filter {
    # Used for query parameter, supported by PowerScale Platform API

    # Only list quotas with this enforcement (non-accounting).
    # enforced = false

    # Set to true to only list quotas which have exceeded one or more of their thresholds.
    # exceeded = false

    # Only list quotas with this setting for include_snapshots.
    # include_snapshots = true

    # Only list quotas with this setting for include_snapshots.
    # path = "/ifs/tfacc_file_system_test"

    # Only list user or group quotas matching this persona (must be used with the corresponding type argument).
    # Format is <PERSONA_TYPE>:<string/integer>, where PERSONA_TYPE is one of USER, GROUP, SID, ID, or GID.
    # persona= "SID:0"

    # If used with the path argument, match all quotas at that path or any descendent sub-directory.
    # recurse_path_children = true

    # If used with the path argument, match all quotas at that path or any parent directory.
    # recurse_path_parents  = true

    # Use the named report as a source rather than the live quotas. See the /quota/reports resource for a list of valid reports.
    # report_id  = "id"

    # Only list quotas matching this type.
    # Allowed values: directory, user, group, default-directory, default-user, default-group
    # type = "directory"

    # Optional named zone to use for user and group resolution.
    # zone  = "System"
  }
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_quota.example_quotas
output "powerscale_quota" {
  value = data.powerscale_quota.example_quotas
}

# Returns all of the PowerScale Quotas in default zone
data "powerscale_quota" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_quota.all
output "powerscale_quota_all" {
  value = data.powerscale_quota.all
}