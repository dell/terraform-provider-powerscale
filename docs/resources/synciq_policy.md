---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
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

title: "powerscale_synciq_policy resource"
linkTitle: "powerscale_synciq_policy"
page_title: "powerscale_synciq_policy Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  This resource is used to manage the SyncIQ Replication Policy entity of PowerScale Array. We can Create, Read, Update and Delete the SyncIQ Replication Policy using this resource. We can also import existing SyncIQ Replication Policy from PowerScale array.
---

# powerscale_synciq_policy (Resource)

This resource is used to manage the SyncIQ Replication Policy entity of PowerScale Array. We can Create, Read, Update and Delete the SyncIQ Replication Policy using this resource. We can also import existing SyncIQ Replication Policy from PowerScale array.


## Example Usage

```terraform
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


# Available actions: Create, Update, Delete and Import.
# If resource arguments are omitted, `terraform apply` will load User Mapping Rules from PowerScale, and save to terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load User Mapping Rules (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting User Mapping Rules from PowerScale.
# For more information, Please check the terraform state file.

resource "powerscale_synciq_peer_certificate" "cert_10_10_10_10" {
  path        = "/ifs/peerCert_10_10_10_10.pem"
  name        = "peerCert_10_10_10_10"
  description = "Certificate for the replication peer cluster 10.10.10.10"
}

# PowerScale Sync IQ Policies can be used to replicate files or directories from one cluster to another. 
resource "powerscale_synciq_policy" "policy" {
  # Required
  name             = "policy1"
  action           = "sync"
  source_root_path = "/ifs"
  target_host      = "10.10.10.10"
  target_path      = "/ifs/policy1Sink"

  # Optional
  description = "Policy 1 description"
  enabled     = true
  file_matching_pattern = {
    or_criteria = [
      {
        and_criteria = [
          {
            type     = "name"
            value    = "tfacc"
            operator = "=="
          },
        ]
      },
      {
        and_criteria = [
          {
            type     = "size"
            value    = "200MB"
            operator = ">"
          }
        ]
      }
    ]
  }

  source_network = {
    pool   = "pool0"
    subnet = "subnet0"
  }

  accelerated_failback                  = false
  bandwidth_reservation                 = 100
  changelist                            = false
  check_integrity                       = true
  cloud_deep_copy                       = "allow"
  delete_quotas                         = true
  disable_quota_tmp_dir                 = true
  enable_hash_tmpdir                    = true
  encryption_cipher_list                = "aes128-ctr,aes192-ctr,aes256-ctr"
  ignore_recursive_quota                = true
  log_level                             = "fatal"
  priority                              = 1
  report_max_age                        = 3600
  report_max_count                      = 100
  rpo_alert                             = 10
  schedule                              = "Every Monday at 2AM"
  skip_when_source_unmodified           = true
  snapshot_sync_existing                = true
  snapshot_sync_pattern                 = "*"
  source_exclude_directories            = ["/ifs/SourceExcludeDir"]
  source_include_directories            = ["/ifs/SourceIncludeDir"]
  source_snapshot_archive               = true
  source_snapshot_expiration            = 8 * 60 * 60 * 24 # 8 days
  source_snapshot_pattern               = ".*"
  sync_existing_snapshot_expiration     = true
  sync_existing_target_snapshot_pattern = "%%{SnapName}-%%{SnapCreateTime}"
  target_certificate_id                 = powerscale_synciq_peer_certificate.cert_10_10_10_10.id
  target_compare_initial_sync           = true
  target_detect_modifications           = true
  target_snapshot_alias                 = "SIQ-%%{SrcCluster}-%%{PolicyName}-latest"
  target_snapshot_archive               = true
  target_snapshot_expiration            = 80 * 60 * 60 * 24 # 80 days
  target_snapshot_pattern               = "SIQ-%%{SrcCluster}-%%{PolicyName}-%Y-%m-%d_%H-%M-%S"
}

# After the execution of above resource block, Sync IQ Policies would have been cached in terraform state file, or
# Sync IQ Policies would have been updated on PowerScale.
# For more information, Please check the terraform state file.

# sheduling a policy when source is modified
resource "powerscale_synciq_policy" "policy_when_source_modified" {
  name             = "policy_when_source_modified"
  enabled          = true
  action           = "sync"
  source_root_path = "/ifs/Source"
  target_host      = "10.10.10.9"
  password         = "W0ulntUWannaKn0w"
  target_path      = "/ifs/Sink2"

  # scheduling
  schedule  = "when-source-modified"
  job_delay = 20 * 60 * 60

  # ocsp
  ocsp_address = "10.20.10.9"
  # Installed ca certificates can be listed using "isi certificate authority list"
  # which gives a list of all certificates with their short IDs.
  # Then full ID can be looked up using "isi certificate authority view <short ID>".
  ocsp_issuer_certificate_id = "16af57a9f676b0ab126095aa5ebadef22ab31119d644ac95cd4b93dbf3f26aeb"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) If 'copy', source files will be copied to the target cluster.  If 'sync', the target directory will be made an image of the source directory:  Files and directories that have been deleted on the source, have been moved within the target directory, or no longer match the selection criteria will be deleted from the target directory.
- `name` (String) User-assigned name of this sync policy.
- `source_root_path` (String) The root directory on the source cluster the files will be synced from.  Modifying this field will result in a full synchronization of all data.
- `target_host` (String) Hostname or IP address of sync target cluster.  Modifying the target cluster host can result in the policy being unrunnable if the new target does not match the current target association.
- `target_path` (String) Absolute filesystem path on the target cluster for the sync destination.

### Optional

- `accelerated_failback` (Boolean) If set to true, SyncIQ will perform failback configuration tasks during the next job run, rather than waiting to perform those tasks during the failback process. Performing these tasks ahead of time will increase the speed of failback operations.
- `allow_copy_fb` (Boolean) If set to true, SyncIQ will allow a policy with copy action failback which is not supported by default.
- `bandwidth_reservation` (Number) The desired bandwidth reservation for this policy in kb/s. This feature will not activate unless a SyncIQ bandwidth rule is in effect.
- `changelist` (Boolean) If true, retain previous source snapshot and incremental repstate, both of which are required for changelist creation.
- `check_integrity` (Boolean) If true, the sync target performs cyclic redundancy checks (CRC) on the data as it is received.
- `cloud_deep_copy` (String) If set to deny, replicates all CloudPools smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, the job will fail. If set to force, replicates all smartlinks to the target cluster as regular files. If set to allow, SyncIQ will attempt to replicate smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, SyncIQ will replicate the smartlinks as regular files.
- `conflicted` (Boolean) NOTE: This field should not be changed without the help of PowerScale support.  If true, the most recent run of this policy encountered an error and this policy will not start any more scheduled jobs until this field is manually set back to 'false'.
- `delete_quotas` (Boolean) If true, forcibly remove quotas on the target after they have been removed on the source.
- `description` (String) User-assigned description of this sync policy.
- `disable_file_split` (Boolean) NOTE: This field should not be changed without the help of PowerScale support.  If true, the 7.2+ file splitting capability will be disabled.
- `disable_fofb` (Boolean) NOTE: This field should not be changed without the help of PowerScale support.  Enable/disable sync failover/failback.
- `disable_quota_tmp_dir` (Boolean) If set to true, SyncIQ will not create temporary quota directories to aid in replication to target paths which contain quotas.
- `disable_stf` (Boolean) NOTE: This field should not be changed without the help of PowerScale support.  Enable/disable the 6.5+ STF based data transfer and uses only treewalk.
- `enable_hash_tmpdir` (Boolean) If true, syncs will use temporary working directory subdirectories to reduce lock contention.
- `enabled` (Boolean) If true, jobs will be automatically run based on this policy, according to its schedule.
- `encryption_cipher_list` (String) The cipher list (comma separated) being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.
- `expected_dataloss` (Boolean) NOTE: This field should not be changed without the help of PowerScale support.  Continue sending files even with the corrupted filesystem.
- `file_matching_pattern` (Attributes) A file matching pattern, organized as an OR'ed set of AND'ed file criteria, for example ((a AND b) OR (x AND y)) used to define a set of files with specific properties.  Policies of type 'sync' cannot use 'path' or time criteria in their matching patterns, but policies of type 'copy' can use all listed criteria. (see [below for nested schema](#nestedatt--file_matching_pattern))
- `force_interface` (Boolean) NOTE: This field should not be changed without the help of PowerScale support.  Determines whether data is sent only through the subnet and pool specified in the "source_network" field. This option can be useful if there are multiple interfaces for the given source subnet.  If you enable this option, the net.inet.ip.choose_ifa_by_ipsrc sysctl should be set.
- `ignore_recursive_quota` (Boolean) If set to true, SyncIQ will not check the recursive quota in target paths to aid in replication to target paths which contain no quota but target cluster has lots of quotas.
- `job_delay` (Number) If `schedule` is set to `when-source-modified`, the duration to wait after a modification is made before starting a job (default is 0 seconds).
- `log_level` (String) Severity an event must reach before it is logged. Accepted values are `fatal`, `error`, `notice`, `info`, `copy`, `debug`, `trace`.
- `log_removed_files` (Boolean) If true, the system will log any files or directories that are deleted due to a sync.
- `ocsp_address` (String) The address of the OCSP responder to which to connect. Set to empty string to disable OCSP.
- `ocsp_issuer_certificate_id` (String) The ID of the certificate authority that issued the certificate whose revocation status is being checked. Set to empty string to disable certificate verification.
- `password` (String) The password for the target cluster. This field is not readable.
- `priority` (Number) Determines the priority level of a policy. Policies with higher priority will have precedence to run over lower priority policies. Valid range is [0, 1]. Default is 0.
- `report_max_age` (Number) Length of time (in seconds) a policy report will be stored.
- `report_max_count` (Number) Maximum number of policy reports that will be stored on the system.
- `restrict_target_network` (Boolean) If you specify true, and you specify a SmartConnect zone in the "target_host" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.
- `rpo_alert` (Number) If `schedule` is set to a time/date, an alert is created if the specified RPO for this policy is exceeded. The default value is 0, which will not generate RPO alerts.
- `schedule` (String) The schedule on which new jobs will be run for this policy.
- `skip_lookup` (Boolean) Skip DNS lookup of target IPs.
- `skip_when_source_unmodified` (Boolean) If true and `schedule` is set to a time/date, the policy will not run if no changes have been made to the contents of the source directory since the last job successfully completed.
- `snapshot_sync_existing` (Boolean) If true, snapshot-triggered syncs will include snapshots taken before policy creation time (requires --schedule when-snapshot-taken).
- `snapshot_sync_pattern` (String) The naming pattern that a snapshot must match to trigger a sync when the schedule is when-snapshot-taken (default is "*").
- `source_exclude_directories` (List of String) Directories that will be excluded from the sync.  Modifying this field will result in a full synchronization of all data.
- `source_include_directories` (List of String) Directories that will be included in the sync.  Modifying this field will result in a full synchronization of all data.
- `source_network` (Attributes) Restricts replication policies on the local cluster to running on the specified subnet and pool. (see [below for nested schema](#nestedatt--source_network))
- `source_snapshot_archive` (Boolean) If true, archival snapshots of the source data will be taken on the source cluster before a sync.
- `source_snapshot_expiration` (Number) The length of time in seconds to keep snapshots on the source cluster.
- `source_snapshot_pattern` (String) The name pattern for snapshots taken on the source cluster before a sync.
- `sync_existing_snapshot_expiration` (Boolean) If set to true, the expire duration for target archival snapshot is the remaining expire duration of source snapshot, requires --sync-existing-snapshot=true
- `sync_existing_target_snapshot_pattern` (String) The naming pattern for snapshot on the destination cluster when --sync-existing-snapshot is true
- `target_certificate_id` (String) The ID of the target cluster certificate being used for encryption. Set to empty string to disable target certificate verification.
- `target_compare_initial_sync` (Boolean) If true, the target creates diffs against the original sync.
- `target_detect_modifications` (Boolean) If true, target cluster will detect if files have been changed on the target by legacy tree walk syncs.
- `target_snapshot_alias` (String) The alias of the snapshot taken on the target cluster after the sync completes. Do not use the value `DEFAULT`.
- `target_snapshot_archive` (Boolean) If true, archival snapshots of the target data will be taken on the target cluster after successful sync completions.
- `target_snapshot_expiration` (Number) The length of time in seconds to keep snapshots on the target cluster.
- `target_snapshot_pattern` (String) The name pattern for snapshots taken on the target cluster after the sync completes. Do not use the value `@DEFAULT`.
- `workers_per_node` (Number) The number of worker threads on a node performing a sync.

### Read-Only

- `id` (String) The system ID given to this sync policy.

<a id="nestedatt--file_matching_pattern"></a>
### Nested Schema for `file_matching_pattern`

Required:

- `or_criteria` (Attributes List) An array containing objects with "and_criteria" properties, each set of and_criteria will be logically OR'ed together to create the full file matching pattern. (see [below for nested schema](#nestedatt--file_matching_pattern--or_criteria))

<a id="nestedatt--file_matching_pattern--or_criteria"></a>
### Nested Schema for `file_matching_pattern.or_criteria`

Required:

- `and_criteria` (Attributes List) An array containing individual file criterion objects each describing one criterion.  These are logically AND'ed together to form a set of criteria. (see [below for nested schema](#nestedatt--file_matching_pattern--or_criteria--and_criteria))

<a id="nestedatt--file_matching_pattern--or_criteria--and_criteria"></a>
### Nested Schema for `file_matching_pattern.or_criteria.and_criteria`

Optional:

- `attribute_exists` (Boolean) For "custom_attribute" type criteria.  The file will match as long as the attribute named by "field" exists.  Default is true.
- `case_sensitive` (Boolean) If true, the value comparison will be case sensitive.  Default is true.
- `field` (String) The name of the file attribute to match on (only required if this is a custom_attribute type criterion).  Default is an empty string "".
- `operator` (String) How to compare the specified attribute of each file to the specified value.  Possible values are: `==`, `!=`, `>`, `>=`, `<`, `<=`, `!`.  Default is `==`.
- `type` (String) The type of this criterion, that is, which file attribute to match on. Accepted values are , `name`, `path`, `accessed_time`, `birth_time`, `changed_time`, `size`, `file_type`, `posix_regex_name`, `user_name`, `user_id`, `group_name`, `group_id`, `no_user`, `no_group`.
- `value` (String) The value to compare the specified attribute of each file to.
- `whole_word` (Boolean) If true, the attribute must match the entire word.  Default is true.




<a id="nestedatt--source_network"></a>
### Nested Schema for `source_network`

Optional:

- `pool` (String) The pool to restrict replication policies to.
- `subnet` (String) The subnet to restrict replication policies to.

## Import

Import is supported using the following syntax:

```shell
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
# terraform import powerscale_synciq_policy.policy <policy name>
# Example:
terraform import powerscale_synciq_policy.policy "policy1"
# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
```