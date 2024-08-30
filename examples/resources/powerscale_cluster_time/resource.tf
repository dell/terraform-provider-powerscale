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
# If resource arguments are omitted, `terraform apply` will load Cluster Time Details from PowerScale, and save to
# terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Cluster Time Details (if not loaded) and update the settings.
# `terraform destroy` will delete the resource from terraform state file rather than deleting Cluster Time Details from PowerScale.
# For more information, Please check the terraform state file.

resource "powerscale_cluster_time" "test" {
  date = "01/12/2024"
  time = "00:32"
  path = "Asia/Kolkata"
}

# After the execution of above resource block, Cluster Time Settings would have been cached in terraform state file, and
# Cluster Time Settings would have been updated on PowerScale.
# For more information, Please check the terraform state file.
