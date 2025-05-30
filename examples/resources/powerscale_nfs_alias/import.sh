# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
#terraform import powerscale_nfs_alias.example global_setting name_of_nfs_alias
# Example:
terraform import powerscale_nfs_alias.example global_setting "alias"
# after running this command, populate parameters in the config file to start managing this resource.
# Note: running "terraform show" after importing shows the current config/state of the resource.

