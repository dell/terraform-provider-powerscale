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