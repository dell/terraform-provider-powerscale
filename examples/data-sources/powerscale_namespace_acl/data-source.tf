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

# This Terraform DataSource is used to query the details of the Namespace ACL from PowerScale array.

# Returns the PowerScale Namespace ACL on PowerScale array
data "powerscale_namespace_acl" "example" {
  # Note: namespace must be specified
  filter {
    namespace = "ifs/example"
    nsaccess  = true
  }
}

# Output value of above block by executing 'terraform output' command
# You can use the the fetched information by the variable data.powerscale_namespace_acl.example
output "powerscale_namespace_acl_example" {
  value = data.powerscale_namespace_acl.example
}