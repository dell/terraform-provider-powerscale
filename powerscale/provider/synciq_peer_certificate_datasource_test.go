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
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSyncIQPeerCertificateDatasource(t *testing.T) {
	certRs := getPeerCertProvisionerConfig() + `
	resource "powerscale_synciq_peer_certificate" "test" {
		path = terraform_data.certificate.output.cert
		name = "tfaccTest"
		description = "Tfacc Test"
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// read invalid ID
			{
				Config: ProviderConfig + `
				data "powerscale_synciq_peer_certificate" "test" {
					id = "invalid"
				}
				`,
				ExpectError: regexp.MustCompile(`.*not found.*`),
			},
			// Read valid ID
			{
				Config: ProviderConfig + certRs + `
				data "powerscale_synciq_peer_certificate" "test" {
					id = powerscale_synciq_peer_certificate.test.id
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_synciq_peer_certificate.test", "certificates.#"),
					resource.TestCheckResourceAttr("data.powerscale_synciq_peer_certificate.test", "certificates.#", "1"),
				),
			},
			// read all mock error
			{
				Config: ProviderConfig + certRs + `
				data "powerscale_synciq_peer_certificate" "test" {
				}
				`,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ListPeerCerts).Return(nil, fmt.Errorf("mock network error")).Build()
				},
				ExpectError: regexp.MustCompile("mock network error"),
			},
			// read all valid
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + certRs + `
				data "powerscale_synciq_peer_certificate" "test" {
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_synciq_peer_certificate.test", "certificates.#"),
				),
			},
		},
	})
}
