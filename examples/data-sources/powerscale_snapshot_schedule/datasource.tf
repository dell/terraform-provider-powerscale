/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# PowerScale SMB shares provide clients network access to file system resources on the cluster

# Returns a list of all the PowerScale Snapshot schedules
data "powerscale_snapshot_schedule" "example_snapshot_schedule_all" {
  filter {

    # Used for query parameter, supported by PowerScale Platform API
    #The direction of the sort.Supported Values:ASC , DESC
    # dir = "ASC"
    # Return no more than this many results at once.
    # limit = 1
    # The field that will be used for sorting. Choices are id, name, path, pattern, schedule, duration, alias, next_run, and next_snapshot. Default is id.
    # sort = "name"
  }
}
output "powerscale_snapshot_schedule_all" {
  value = data.powerscale_snapshot_schedule.example_snapshot_schedule_all
}

# Returns a list of PowerScale Snapshot schedules based on path filter block
data "powerscale_snapshot_schedule" "example_snapshot_schedule" {
  filter {
    # Used to specify names of snapshot schedules
    names = ["Snapshot schedule 370395356"]

    # Used for query parameter, supported by PowerScale Platform API
    #The direction of the sort.Supported Values:ASC , DESC
    # dir = "ASC"
    # Return no more than this many results at once.
    # limit = 1
    # The field that will be used for sorting. Choices are id, name, path, pattern, schedule, duration, alias, next_run, and next_snapshot. Default is id.
    # sort = "name"

  }
}
output "powerscale_snapshot_schedule" {
  value = data.powerscale_snapshot_schedule.example_snapshot_schedule
}