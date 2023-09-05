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

# Returns all of the PowerScale access zones and their details
# PowerScale access zones allow you to isolate data and control who can access data in each zone.
data "powerscale_accesszone" "all" {
}

output "powerscale_accesszone_data_all" {
  value = data.powerscale_accesszone.all
}

# Returns a subset of the PowerScale access zones based on the names provided in the `names` filter block and their details
data "powerscale_accesszone" "test" {
  # Optional list of names to filter upon
  filter {
    names = ["tfaccAccessZone"]
  }
}

output "powerscale_accesszone" {
  value = data.powerscale_accesszone.test
}
