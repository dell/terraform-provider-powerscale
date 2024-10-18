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

# PowerScale SyncIQ Certificate allows you to get a list of SyncIQ Peer Certificates or a peer certificate by its ID.

# Returns a list of PowerScale SyncIQ Certificates
data "powerscale_synciq_peer_certificate" "all_certificates" {
}

# Returns a the PowerScale SyncIQ Certificate with given ID
data "powerscale_synciq_peer_certificate" "one_certificate" {
  id = "g23j9a1f83h12n5j4"
}

# Returns the PowerScale SyncIQ Certificate with given name
data "powerscale_synciq_peer_certificate" "one_certificate_by_name" {
  filter {
    name = "tfaccTest"
  }
}

# Output value of above block by executing 'terraform output' command.
# The user can use the fetched information by the variable data.powerscale_synciq_peer_certificate.all_certificates.certificates
output "powerscale_synciq_all_certificates" {
  value = data.powerscale_synciq_peer_certificate.all_certificates.certificates
}

# The user can use the fetched certificate by ID by the variable data.powerscale_synciq_peer_certificate.one_certificate.certificates[0]
output "certificateByID" {
  value = data.powerscale_synciq_peer_certificate.one_certificate.certificates[0]
}

# The user can use the fetched certificate by name by the variable data.powerscale_synciq_peer_certificate.one_certificate_by_name.certificates[0]
output "certificateByName" {
  value = data.powerscale_synciq_peer_certificate.one_certificate_by_name.certificates[0]
}

# Get syncIQ certificate by status
# Step 1: We shall use the datasource to get all the certificates as shown above
# Step 2: We index them by status
output "certificatesByStatus" {
  value = {
    "valid"    = [for certificate in data.powerscale_synciq_peer_certificate.all_certificates.certificates : certificate if certificate.status == "valid"]
    "invalid"  = [for certificate in data.powerscale_synciq_peer_certificate.all_certificates.certificates : certificate if certificate.status == "invalid"]
    "expired"  = [for certificate in data.powerscale_synciq_peer_certificate.all_certificates.certificates : certificate if certificate.status == "expired"]
    "expiring" = [for certificate in data.powerscale_synciq_peer_certificate.all_certificates.certificates : certificate if certificate.status == "expiring"]
  }
}

# After the successful execution of above said block, We can see the output value by executing 'terraform output' command.
