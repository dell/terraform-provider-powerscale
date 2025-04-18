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

# Available actions: Create, Update, Delete and Import
# If resource arguments are omitted, `terraform apply` will load Writable Snapshot Details from PowerScale, and save to
# terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Writable Snapshot Details (if not loaded).
# `terraform destroy` will delete the Writable Snapshot resource from terraform state file as well as from PowerScale.
# For more information, Please check the terraform state file.


# Scenario 1
# To create multiple writable snapshots using a snap id

# snapshot ID for which writable snapshot needs to be created
variable "snap_id" {
  description = "The ID of the source snapshot"
  type        = string
  # default     = "snap_id"
}

# Number of writable snapshots to be created from the snap_id
variable "num_writable_snapshots" {
  description = "The number of writable snapshots to create"
  type        = number
  default     = 3
}

# Example to create multiple writable snapshots from a single snapshot
resource "powerscale_writable_snapshot" "writablesnap_multiple1" {
  count = var.num_writable_snapshots

  dst_path = "/ifs/writable_snapshot_snap${count.index}"
  snap_id  = var.snap_id
}


# Scenario 2
# To create multiple writable snaphots using datasource

# Fetch snapshot data using snapshot datasource filters
data "powerscale_snapshot" "all" {
  filter {
    # sort = "created"
    # dir = "asc"
    # limit = 5
  }
}

# output command to verify the fetched snapshot
output "powerscale_snapshot_all_snaps" {
  value = data.powerscale_snapshot.all.snapshots_details
}

# create multiple writable snapshots with snap id fetched using datasource
locals {
  writable_snapshots = [
    for i, snap in data.powerscale_snapshot.all.snapshots_details : {
      dst_path = "/ifs/writable_snapshot${i}"
      snap_id  = snap.id
    }
  ]
}

resource "powerscale_writable_snapshot" "writablesnap_multiple2" {
  for_each = { for snap in local.writable_snapshots : snap.dst_path => snap }

  dst_path = each.value.dst_path
  snap_id  = each.value.snap_id
}


# Scenario 3
# To create a single writable snapshot

resource "powerscale_writable_snapshot" "writablesnap_single" {
  # dst_path is the path of the writable snapshot.
  dst_path = "/ifs/writable_snapshot"

  # snap_id is the source snapshot of the writable snapshot.
  snap_id = "snap_id"
}

# After the execution of above resource block, single/multiple Writable Snapshot(s) would have been cached in terraform state file, and
# Single/Multiple new Writable Snapshot(s) would have been created on PowerScale.
# For more information, Please check the terraform state file.
