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

# A Peer Certificate can be imported either by using name (if non-empty) or ID.

# The command to import by ID is
# terraform import powerscale_synciq_peer_certificate.certificate <certificate ID>
# Example:
terraform import powerscale_synciq_peer_certificate.certificate "08834c2212932cdeb75d61f9ff661892c97afe500152b70e50a8248188fe7e27"

# The command to import by name is
# terraform import powerscale_synciq_peer_certificate.certificate2 <name:certificate name>
# Example:
terraform import powerscale_synciq_peer_certificate.certificate2 "name:peer_certificate_01"

# after running any of these commands, populate the path field with the value "/dummy" to start managing this resource. Add other fields as required.
# Note: running "terraform show" after importing shows the current config/state of the resource. You can copy/paste that config to make it easier to manage the resource.