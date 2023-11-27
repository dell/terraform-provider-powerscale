# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The command is
# terraform import powerscale_nfs_export.example_export [<zoneID>]:<name>
# Example 1:  <zoneID> is Optional, defaults to System:
terraform import powerscale_nfs_export.example_export example_export
# Example 2:
terraform import powerscale_nfs_export.example_export zone_id:example_export
# after running this command, populate the name field and other required parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.
