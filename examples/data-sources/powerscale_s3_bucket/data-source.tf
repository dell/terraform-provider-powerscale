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

# This Terraform DataSource is used to query the details of existing S3 Bucket from PowerScale array.

# Returns a list of PowerScale s3 bucket based on filter block
data "powerscale_s3_bucket" "example_s3_buckets" {
  filter {
    # Used for query parameter, supported by PowerScale Platform API

    # Only list s3 bucket in this zone.
    # zone = "System"

    # Only list s3 bucket owned by this.
    # owner = "root"
  }
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_s3_bucket.example_s3_buckets
output "powerscale_s3_bucket" {
  value = data.powerscale_s3_bucket.example_s3_buckets
}

# Returns all of the PowerScale S3 Bucket in default zone
data "powerscale_s3_bucket" "all" {
}

# Output value of above block by executing 'terraform output' command
# The user can use the fetched information by the variable data.powerscale_s3_bucket.all
output "powerscale_s3_bucket_all" {
  value = data.powerscale_s3_bucket.all
}
