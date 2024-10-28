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

# PowerScale SyncIQ Replication Job datasource allows you to get a list of SyncIQ replication jobs.

# Returns a list of PowerScale SyncIQ Replication Jobs 
data "powerscale_synciq_replication_job" "all_jobs" {
}

# Returns a list of PowerScale SyncIQ Replication Jobs with running state
data "powerscale_synciq_replication_job" "jobByState" {
  filter {
    state = "running"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_synciq_replication_job.all_jobs.synciq_jobs
output "powerscale_synciq_all_replication_jobs" {
  value = data.powerscale_synciq_replication_job.all_jobs.synciq_jobs
}

# The user can use the fetched information by the variable data.powerscale_synciq_replication_job.all_jobs.synciq_jobs
output "replicationJobByState" {
  value = data.powerscale_synciq_replication_job.jobByState.synciq_jobs
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
