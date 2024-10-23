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

# PowerScale Writable Snapshots allows you to get a list of Writable Snapshots or a Writable Snapshot by its Path which can be added in filters.

# Returns a list of PowerScale Writable Snapshots 
data "powerscale_writable_snapshot" "all_snaps" {
}

# Returns a the PowerScale Writable Snapshot with given Path
data "powerscale_writable_snapshot" "test" {
  filter {
    path = "/ifs/path/to/writable/snap"
    limit = 5
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_writable_snapshot.all_snaps.writable
# output "powerscale_writable_snapshot" {
#   value = data.powerscale_writable_snapshot.all_snaps.writable
# }

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
