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

# Restore snapshot using snaprevert job
resource "powerscale_snapshot_restore" "test" {
  snaprevert_params = {
    allow_dup   = true
    snapshot_id = "snapshot_id"
  }
}

# terraform destroy will delete the snaprevert domain if restore is done using snapshot revert.

# Restore snapshot using copy operation
resource "powerscale_snapshot_restore" "test" {
  copy_params = {
    directory = {
      source      = "Path of the snapshot to copy" # e.g. /namespace/ifs/.snapshot/snapshot_name/directory
      destination = "Path of the destination" # '/' is not required at the start e.g. ifs/dest
      overwrite   = true
    }
  }
}

# Restore snapshot using clone operation
resource "powerscale_snapshot_restore" "test" {
  clone_params = {
    source      = "Path of the snapshot to copy" # e.g. /namespace/ifs/.snapshot/snapshot_name/directory/file
    destination = "Path of the destination" # '/' is not required at the start e.g. ifs/dest/test.txt
    snapshot_id = "Snapshot ID"
  }
}
