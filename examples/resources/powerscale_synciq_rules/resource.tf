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

# Available actions: Create and Update updates the syncIQ  replication rules on the PowerScale array.
# Delete will delete the state file. Import action is also available.
# After `terraform apply` of this example file will update the settings according to the attributes set in the config

# PowerScale SynIQ replication rules allows you to manage the replication rules on the Powerscale array
resource "powerscale_synciq_rules" "all_rules" {
  bandwidth_rules = [
    {
      limit       = 10000
      description = "Bandwidth limit for Weekend"
      schedule = {
        begin        = "00:00"
        days_of_week = ["saturday", "sunday"]
        end          = "23:59"
      }
    },
    {
      limit       = 2000
      description = "Bandwidth limit for Weekdays"
      schedule = {
        days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      }
    },
  ]

  cpu_rules = [
    {
      limit       = 16
      description = "CPU limit"
    }
  ]

  file_count_rules = [
    {
      limit       = 50
      description = "File limit"
    }
  ]

  worker_rules = [
    {
      limit       = 10
      description = "Worker limit"
      enabled     = false
    }
  ]
}

# After the execution of above resource block, replication rules would have been updated on the PowerScale array. For more information, Please check the terraform state file.
