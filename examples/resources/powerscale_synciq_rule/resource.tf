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
# After `terraform apply` of this example file it will create a new SyncIQ replication rule on PowerScale.

# PowerScale SyncIQ performance rule allows you to limit the resources used for SyncIQ replication.
resource "powerscale_synciq_rule" "test" {
  // required

  // type and value of the replication rule limit
  // Units are kb/s for bandwidth, files/s for file-count, processing percentage used for cpu, or percentage of maximum available workers
  type  = "bandwidth"
  limit = 20000

  // optional
  description = "tfacc updated"
  enabled     = true
  schedule = {
    days_of_week = ["monday", "wednesday", "thursday"]
    // time in 24h format
    begin = "01:00",
    end   = "22:59",
  }
}

# After the execution of above resource block, syncIQ performance rule would have been created on the PowerScale array. 
# For more information, Please check the terraform state file. 