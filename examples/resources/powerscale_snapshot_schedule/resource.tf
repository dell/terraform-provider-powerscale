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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will create NFS export on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale SMB shares provide clients network access to file system resources on the cluster
resource "powerscale_snapshot_schedule" "snap_schedule" {
  # Required name of snapshot schedule for creating
  name = "test_snap_schedule"

  # Optional fields

  # The path snapshotted. Default set to : "/ifs"
  # path = "/ifs/tfacc_test"

  # Time value in String for which snapshots created by this snapshot schedule should be retained. Values supported are of format : 
  # "Never Expires, x Seconds(s), x Minute(s), x Hour(s), x Week(s), x Day(s), x Month(s), x Year(s) where x can be any integer value.
  # Default set to : "1 Week(s)"
  # retention_time = "3 Hour(s)" 

  # Alias name to create for each snapshot.
  # alias = "snap_alias"

  # Pattern expanded with strftime to create snapshot names. Some sample values for pattern are: 'snap-%F' would yield snap-1984-03-20 , 'backup-%FT%T' would yield backup-1984-03-20T22:30:00".
  # Default set to : "ScheduleName_duration_%Y-%m-%d_%H:%M"
  # pattern = "ScheduleName_duration_%Y-%m-%d_%H:%M"

  # The isidate compatible natural language description of the schedule.It specifies the frequency of the schedule.You can specify this as combination of <interval> and <frequency> where each of them can be defined as:  
  #  <interval>:
  # 	  *Every [ ( other | <integer> ) ] ( weekday | day | week [ on <day>] | month [ on the <integer> ] | <day>[, ...] [ of every [ ( other | <integer> ) ] week ] | The last (day | weekday | <day>) of every [ (other | <integer>) ] month | The <integer> (weekday | <day>) of every [ (other | <integer>) ] month | The <integer> of every [ (other | <integer>) ] month | Yearly on <month> <integer> | Yearly on the (last | <integer>) [ weekday | <day> ] of <month>
  #  <frequency>:
  # 			*at <hh>[:<mm>] [ (AM | PM) ] | every [ <integer> ] (hours | minutes) [ between <hh>[:<mm>] [ (AM | PM) ] and <hh>[:<mm>] [ (AM | PM) ] | every [ <integer> ] (hours | minutes) [ from <hh>[:<mm>] [ (AM | PM) ] to <hh>[:<mm>] [ (AM | PM) ]
  #  Additionally:
  # 			<integer> can include "st," "th," or "rd," e.g., "Every 1st month."
  # 			<day> can be a day of the week or a three-letter abbreviation, e.g., "saturday" or "sat."
  # 			<month> must be the name of the month or its abbreviation, e.g., "July" or "Jul."
  # 	Some sample values:  "Every 2 days", "Every 3rd weekday at 11 PM", "Every month on the 15th at 1:30 AM"`
  # Default set to : "every 1 days at 12:00 AM"
  # schedule = "every 1 days at 12:00 AM"

}

# After the execution of above resource block, an SMB share would have been created on the PowerScale array.
# For more information, Please check the terraform state file.