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

# Available actions: Create, Update and Delete
# After `terraform apply` of this example file it will generate the s3 key for the user.
# For more information, Please check the terraform state file.

# PowerScale S3 key to generate the keys for users to sign the requests you send to the S3 protocol.

resource "powerscale_s3_key" "skm" {
  user                     = "tf_user"
  zone                     = "System"
  existing_key_expiry_time = 10
}

output "key" {
  value = powerscale_s3_key.skm
}