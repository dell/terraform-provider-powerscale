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


# Available actions: Create, Read, Update, Delete and Import.

# Step 1: Upload a certificate to PowerScale.
# This can be done using the file provisioner.
resource "terraform_data" "cert" {
  provisioner "file" {
    source      = "/root/certs/pscale.crt" # certificate that user has been given
    destination = "/ifs/peerCert1.crt"
  }

  connection {
    type     = "ssh"
    user     = var.username
    password = var.password
    host     = var.powerscaleIP
    port     = var.PowerscaleSSHPort
  }

  provisioner "remote-exec" {
    when   = destroy
    inline = ["rm -f /ifs/peerCert1.crt"]
  }
}

# PowerScale Sync IQ Peer Certificates can be used to establish trust to the peer/target cluster where files are to be replicated to. 
resource "powerscale_synciq_certificate" "certificate" {
  depends_on = [terraform_data.cert]
  // required
  // Cannot be updated, requires-replace
  path = "/ifs/peerCert1.crt"
  // optional
  name         = "cert1"
  descriptuion = "cert1"
}

# After the execution of above resource block, Sync IQ Certificate would have been cached in terraform state file, or
# Sync IQ Certificate would have been created on PowerScale.
# For more information, Please check the terraform state file.