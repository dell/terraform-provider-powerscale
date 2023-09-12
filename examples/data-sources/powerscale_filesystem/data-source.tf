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

# This terraform DataSource is used to query the existing FileSystem(Namespace Directory) from PowerScale array.It allows you to get information which includes all metadata , access control , quotas and snapshots related information for the directory.

# Returns the information related to the specified PowerScale FileSystem(Namespace Directory) based on the directory path. If directory path is not set it will give details regarding the default "/ifs" directory.
data "powerscale_filesystem" "system" {
  # Required parameter, path of the directory filesystem datasource, defaults to "/ifs" if not set
  directory_path = "/ifs/tfacc_file_system_test"
}

output "powerscale_filesystem_1" {
  value = data.powerscale_filesystem.system
}
# After the successful execution of above block, We can see the output value by executing 'terraform output' command.
# Also, we can use the fetched information by the variable data.powerscale_filesystem.system"
