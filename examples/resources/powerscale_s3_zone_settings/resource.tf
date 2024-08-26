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

# Available actions: Create, Update, Delete and Import
# After `terraform apply` of this example file it will modify S3 Zone Settings on  the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale S3 Zone settings allows you to configure S3 Zone settings on PowerScale.
resource "powerscale_s3_zone_settings" "s3_zone_setting" {
  zone                         = "tfacc-test"
  root_path                    = "/ifs/"
  base_domain                  = "dell.com"
  bucket_directory_create_mode = 511
  object_acl_policy            = "deny" # replace or deny
}