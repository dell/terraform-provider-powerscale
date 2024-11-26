---
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
page_title: "SyncIQ Workflow",
title: "SyncIQ Workflow"
linkTitle: "SyncIQ Workflow"
---


SyncIQ Workflow Explanation

SyncIQ is a data management solution that helps manage data movement between different locations. The workflow involves enabling the SyncIQ service, creating a SyncIQ policy, creating a replication job from the policy, and monitoring the job's status using a replication report.

Step 1: Enable SyncIQ Service

To begin, the SyncIQ service needs to be enabled. This is done using the powerscale_synciq_global_settings resource in Terraform. The service parameter is set to "on" to enable the service.

```
### Make sure SyncIQ service is enabled

resource "powerscale_synciq_global_settings" "enable_synciq" {
  service                 = "on"
}

```
Step 2: Create SyncIQ Policy

Next, a SyncIQ policy needs to be created. A policy defines the rules for data movement, such as the source and target locations, and the action to take (e.g., sync). In this example, a policy named policy1 is created with the following settings:

- action is set to "sync" to synchronize data between the source and target locations.
- source_root_path is set to "/ifs" to specify the source location.
- target_host is set to "10.10.10.10" to specify the target location.
- target_path is set to "/ifs/policy1Sink" to specify the target path.

```
### Create SyncIQ policy with action sync

resource "powerscale_synciq_policy" "policy1" {
  name             = "policy1"
  action           = "sync" # action can be sync or copy
  source_root_path = "/ifs"
  target_host      = "10.10.10.10"
  target_path      = "/ifs/policy1Sink"
}
```

Step 3: Create Replication Job

A replication job is created from the SyncIQ policy using the powerscale_synciq_replication_job resource. The job is configured to run the policy (identified by the id parameter) and is not paused (i.e., is_paused is set to false).
```
### Create replication job from SyncIQ policy

resource "powerscale_synciq_replication_job" "job1" {
  action    = "run" # action can be run, test, resync_prep, allow_write or allow_write_revert
  id        = "policy1"
  is_paused = false
}
```
Step 4: Monitor Replication Job Status

To monitor the status of the replication job, a replication report can be used. The powerscale_synciq_replication_report resource is used to filter the report to show only the replication job with the name Policy1.

```
### Use replication report to view the status of the job

data "powerscale_synciq_replication_report" "filtering" {
  filter {
    policy_name        = "Policy1"
  }
}
```
By following these steps, SyncIQ can be used to manage data movement between different locations and monitor the status of the replication jobs.








