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
resource "powerscale_smb_share" "share_example" {
  auto_create_directory = true
  name                  = "smb_share_example"
  path                  = "/ifs/smb_share_example"
  permissions = [
    {
      permission      = "full"
      permission_type = "allow"
      trustee = {
        id   = "SID:S-1-1-0",
        name = "Everyone",
        type = "wellknown"
      }
    }
  ]
}