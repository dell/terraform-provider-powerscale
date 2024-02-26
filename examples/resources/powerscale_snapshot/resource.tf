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

# Available actions: Create, Update (name, expire), Delete and Import
# After `terraform apply` of this example file it will create a new snapshot for the path name set in `path` attribute on the PowerScale

# Note: This resource can be used to create a manual snapshot at a single point in time. Howerver, if the user wants to take snapshots on a regular cadence, they should use the snapshot_schedules resource. 

# PowerScale Snapshots is a logical pointer to data that is stored on a cluster at a specific point in time.
resource "powerscale_snapshot" "snap" {

  # Required path to the filesystem to which the snapshot will be taken of
  # This cannot be changed after create
  path = "/ifs/tfacc_file_system_test"

  # Optional name of the new snapshot. If unset uses the current date and time for the name attribute (Can be modified)
  name = "tfacc_snapshot_1"

  # Optional set_expires The amount of time from creation before the snapshot will expire and be eligible for automatic deletion.  (Can be modified)
  # Options: Never(default if unset), 1 Day, 1 Week, 1 Month.
  set_expires = "1 Day"
}

# After the execution of above resource block, snapshot would have been created on the PowerScale array. For more information, Please check the terraform state file. 