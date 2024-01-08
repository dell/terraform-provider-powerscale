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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file for the first time, you will create a NTP Server on the PowerScale

# PowerScale NTP Server allows you to synchronize the system time
resource "powerscale_ntpserver" "ntp_server_test" {
  #   Required
  #   Name should be a qualified name of an existing NTP Server
  name = "ntp_server_example"

  #   Optional query parameters
  key = "ntp_server_key_example"
}

# After the execution of above resource block, NTP Server would have been created on the PowerScale array.
# For more information, Please check the terraform state file.