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

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"terraform-provider-powerscale/powerscale/helper"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func getPeerCertProvisionerConfig() string {
	return fmt.Sprintf(`
	locals {
		cert_dir = "/ifs/testAccPeerCertDir"
		key = "${local.cert_dir}/testAccPeerCertKey.pem"
		cert = "${local.cert_dir}/testAccPeerCert.pem"
		subj = "/C=US/ST=California/L=The Cloud/O=Dell/OU=ISG/CN=test_acc_user"
	}
	
	resource "terraform_data" "certificate" {
		connection {
			type     = "ssh"
			user     = "%s"
			password = "%s"
			host     = "%s"
			port     = %s
		}
		input = {
			dir = local.cert_dir
			cert = local.cert
		}
		provisioner "remote-exec" {
			inline = [
				"rm -rf ${local.cert_dir}",
				"mkdir -m 777 ${local.cert_dir}",
				"openssl req -x509 -days 365 -newkey rsa:4096 -keyout ${local.key} -out ${local.cert} -nodes -subj \"${local.subj}\"",
			]
		}
		provisioner "remote-exec" {
			when = destroy
			inline = ["rm -rf ${self.output.dir}"]
		}
	}
	
	`,
		powerscaleUsername,
		powerscalePassword,
		powerScaleSSHIP,
		powerscaleSSHPort,
	)
}

// TestAccSyncIQPeerCertificateResource - Tests syncIQ peer certificate resource.
func TestAccSyncIQPeerCertificateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// invalid create
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = "/ifs/invalid.pem"
				}
				`,
				ExpectError: regexp.MustCompile(`.*Failed to create SyncIQ Peer Certificate.*`),
			},
			// create valid
			{
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
					name = "tfaccTest"
					description = "Tfacc Test"
				}
				`,
			},
			// import with name
			{
				ResourceName: "powerscale_synciq_peer_certificate.test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if states[0].Attributes["description"] != "Tfacc Test" {
						return fmt.Errorf("expected description %s, got %s", "Tfacc Test", states[0].Attributes["description"])
					}
					return nil
				},
				ImportStateId: "tfaccTest",
			},
			// invalid import - invalid name
			{
				ResourceName:  "powerscale_synciq_peer_certificate.test",
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile(`.*Could not find syncIQ peer certificate with name invalid.*`),
			},
			// invalid import - empty name
			{
				ResourceName: "powerscale_synciq_peer_certificate.test",
				ImportState:  true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return "", nil
				},
				ExpectError: regexp.MustCompile(`.*Cannot import syncIQ peer certificate with empty name.*`),
			},
			// mock import with name error
			{
				ResourceName:  "powerscale_synciq_peer_certificate.test",
				ImportState:   true,
				ImportStateId: "tfaccTest",
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ListPeerCerts).Return(nil, fmt.Errorf("mock import with name error")).Build()
				},
				ExpectError: regexp.MustCompile(`.*mock import with name error.*`),
			},
			// mock delete error
			{
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
					name = "tfaccTest"
					description = "Tfacc Test"
				}
				`,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.DeletePeerCert).Return(fmt.Errorf("mock delete error")).Build()
				},
				ExpectError: regexp.MustCompile(`.*mock delete error.*`),
				Destroy:     true,
			},
			// mock update error
			{
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
					name = "tfaccTest2"
					description = "Tfacc Test"
				}
				`,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.UpdatePeerCert).Return(fmt.Errorf("mock update error")).Build()
				},
				ExpectError: regexp.MustCompile(`.*mock update error.*`),
			},
			{
				// Update testing
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
					name = "tfaccTest2"
					description = "Tfacc Test 2"
				}
				`,
			},
		},
	})
}

// TestAccSyncIQPeerCertificateResourceMinimal - Tests syncIQ peer certificate resource with minimal config.
func TestAccSyncIQPeerCertificateResourceMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
				}
				`,
			},
			{
				// Add name and description
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
					name = "tfaccTest2"
					description = "Tfacc Test 2"
				}
				`,
			},
			// remove name and description
			{
				Config: ProviderConfig + getPeerCertProvisionerConfig() + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = terraform_data.certificate.output.cert
				}
				`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// check that changing the path creates a recreate plan
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_peer_certificate" "test" {
					path = "/ifs/invalid.pem"
				}
				`,
				ExpectError: regexp.MustCompile(`.*Failed to create SyncIQ Peer Certificate.*`),
			},
		},
	})
}
