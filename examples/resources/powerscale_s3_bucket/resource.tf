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
# After `terraform apply` of this example file it will create S3 Bucket on specified paths on the PowerScale Array.
# For more information, Please check the terraform state file.

# PowerScale S3 Bucket enables access to file-based data that is stored on OneFS clusters as objects.

resource "powerscale_s3_bucket" "s3_bucket_example" {
  # Required attributes and update not supported
  name = "s3-bucket-example"
  path = "/ifs/s3_bucket_example"

  # Optional attributes and update not supported, 
  # Their default value shows as below if not provided during creation 
  # create_path = false
  # owner = "root"
  # zone = "System"

  # Optional attributes, can be updated
  #
  # By default acl is an empty list. To add an acl item, both grantee and permission are required.
  # Accepted values for permission are: READ, WRITE, READ_ACP, WRITE_ACP, FULL_CONTROL 
  # acl = [{
  #   grantee = {
  #     name = "root"
  #     type = "user"
  #   }
  #   permission = "FULL_CONTROL"
  # }]
  #
  # By default description is empty
  # description = ""
  #
  # Accepted values for object_acl_policy are: replace, deny.
  # The default value would be replace if unset.
  # object_acl_policy = "replace"
}

# After the execution of above resource block, a S3 Bucket would have been created on the PowerScale array.
# For more information, Please check the terraform state file.
