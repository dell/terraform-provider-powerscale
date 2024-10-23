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
# If resource arguments are omitted, `terraform apply` will load Writable Snapshot Details from PowerScale, and save to
# terraform state file.
# If any resource arguments are specified, `terraform apply` will try to load Writable Snapshot Details (if not loaded).
# `terraform destroy` will delete the Writable Snapshot resource from terraform state file as well as from PowerScale.
# For more information, Please check the terraform state file.

resource "powerscale_writable_snapshot" "writablesnap" {
  # dst_path is the path of the writable snapshot.
  dst_path = "/ifs/abcd"

  # snap_id is the source snapshot of the writable snapshot.
  snap_id = "5709"
}

# After the execution of above resource block, Writable Snapshot Settings would have been cached in terraform state file, and
# A new Writable Snapshot would have been created on PowerScale.
# For more information, Please check the terraform state file.
