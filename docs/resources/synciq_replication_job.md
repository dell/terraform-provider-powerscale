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

title: "powerscale_synciq_replication_job resource"
linkTitle: "powerscale_synciq_replication_job"
page_title: "powerscale_synciq_replication_job Resource - terraform-provider-powerscale"
subcategory: ""
description: |-
  The PowerScale SyncIQ ReplicationJob resource provides a means of managing replication jobs on PowerScale clusters.
           This resource allows for the manual triggering of replication jobs to replicate data from a source PowerScale cluster to a target PowerScale cluster.
           Note: The replication job is an asynchronous operation, and this resource does not provide real-time monitoring of the job's status.
           To check the status of the job,please use the powerscalesynciqreplication_report datasource.
---

# powerscale_synciq_replication_job (Resource)

The PowerScale SyncIQ ReplicationJob resource provides a means of managing replication jobs on PowerScale clusters.
		 This resource allows for the manual triggering of replication jobs to replicate data from a source PowerScale cluster to a target PowerScale cluster. 
		 Note: The replication job is an asynchronous operation, and this resource does not provide real-time monitoring of the job's status. 
		 To check the status of the job,please use the powerscale_synciq_replication_report datasource.


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

# Available actions: Create and  Update updates the syncIQ Replication Job. Delete will delete job and clear the state file. 
# After `terraform apply` of this example file will perform action on the synciq replicaiton job according to the attributes set in the config

# PowerScale SynIQ Replication Job allows you to manage the SyncIQ Replication Jobs on the Powerscale array
resource "powerscale_synciq_replication_job" "job1" {
  action    = "run"             # action can be run, test, resync_prep, allow_write or allow_write_revert
  id        = "TerraformPolicy" # id/name of the synciq policy, use synciq policy resource to create policy.
  is_paused = false             # change job state to running or paused.
}

# There are other attributes values as well. Please refer the documentation.
# After the execution of above resource block, job would have been extecuted/updated on the PowerScale array. For more information, Please check the terraform state file.
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) Action for the job 
				 run - to start the replication job using synciq policy
				 test - to test the replication job using synciq policy, 
				 resync_prep - Resync_prep is a preparation step in PowerScale SyncIQ replication jobs that helps ensure a successful replication operation by performing a series of checks and verifications on the source and target volumes before starting the replication process., 
				 allow_write - allow_write determines whether the replication job allows writes to the target volume during the replication process. When configured, the target volume is writable, and any changes made to the target volume will be replicated to the source volume. This is useful in scenarios where you need to make changes to the target volume, such as updating files or creating new files, while the replication job is running.,
				 allow_write_revert - allow_write_revert determines whether the replication job allows writes to the target volume when reverting a replication job. When configure, the target volume is writable during the revert process, allowing changes made to the target volume during the revert process to be replicated to the source volume.
- `id` (String) ID/Name of the policy

### Optional

- `is_paused` (Boolean) change job state to running or paused.
- `wait_time` (Number) Wait Time for the job

Unless specified otherwise, all fields of this resource can be updated.

