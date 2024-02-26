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

func TestAccNtpServerDataSourceNames(t *testing.T) {
	var ntpServerTerraformName = "data.powerscale_ntpserver.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Filter by names
			{
				Config: ProviderConfig + NtpServerDataSourceNamesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ntpServerTerraformName, "ntp_servers_details.#", "1"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "ntp_servers_details.0.id", "ntp_server_at_test"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "ntp_servers_details.0.name", "ntp_server_at_test"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "ntp_servers_details.0.key", "ntp_server_key"),
				),
			},
		},
	})
}

func TestAccNtpServerDataSourceAll(t *testing.T) {
	var ntpServerTerraformName = "data.powerscale_ntpserver.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + NtpServerAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(ntpServerTerraformName, "ntp_servers_details.#"),
				),
			},
		},
	})
}

func TestAccNtpServerDataSourceNamesErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + NtpServerDataSourceNameConfigErr,
				ExpectError: regexp.MustCompile(`.*not a valid powerscale ntp server*.`),
			},
		},
	})
}

func TestAccNtpServerDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + NtpServerDataSourceFilterConfigErr,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccNtpServerDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNtpServers).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NtpServerAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNtpServerDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.NtpServerDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NtpServerAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var NtpServerDataSourceNamesConfig = `
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "ntp_server_at_test"
	key = "ntp_server_key"
}

data "powerscale_ntpserver" "test" {
	filter {
		names = ["ntp_server_at_test"]
	}
	depends_on = [
		powerscale_ntpserver.ntp_server_test
	]
}
`

var NtpServerAllDataSourceConfig = `
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "ntp_server_at_test"
	key = "ntp_server_key"
}

data "powerscale_ntpserver" "all" {
	depends_on = [
		powerscale_ntpserver.ntp_server_test
	]
}
`

var NtpServerDataSourceNameConfigErr = `
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "ntp_server_at_test"
	key = "ntp_server_key"
}

data "powerscale_ntpserver" "test" {
	filter {
		names = ["BadName"]
	}
	depends_on = [
		powerscale_ntpserver.ntp_server_test
	]
}
`

var NtpServerDataSourceFilterConfigErr = `
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "ntp_server_at_test"
	key = "ntp_server_key"
}

data "powerscale_ntpserver" "test" {
	filter {
		invalidFilter = "badFilter"
	}
	depends_on = [
		powerscale_ntpserver.ntp_server_test
	]
}
`
