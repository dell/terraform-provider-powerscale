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

# Returns all of the PowerScale snapshots and their details
# PowerScale Snapshots is a logical pointer to data that is stored on a cluster at a specific point in time.
data "powerscale_snapshot" "all" {
}

output "powerscale_snapshot_data_all" {
  value = data.powerscale_snapshot.all
}

# Returns a subset of the PowerScale snapshots in order based on the filters provided in the filter block and their details
data "powerscale_snapshot" "test" {
  filter {
    path = "/ifs"
    sort = "name"
    dir  = "DESC"
  }
}

#Returns the snapshot with the given name
data "powerscale_snapshot" "test" {
  filter {
    name = "name"
  }
}

output "powerscale_snapshot" {
  value = data.powerscale_snapshot.test
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
# Also, we can use the fetched information by the variable data.powerscale_snapshot.all