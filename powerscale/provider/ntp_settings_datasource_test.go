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

package provider

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNtpSettingsDataSourceAll(t *testing.T) {
	var ntpServerTerraformName = "data.powerscale_ntpsettings.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + NtpSettingsAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ntpServerTerraformName, "key_file", "/ifs/"),
				),
			},
		},
	})
}

func TestAccNtpSettingsDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNtpSettings).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + NtpSettingsAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNtpSettingsDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.NtpSettingsDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NtpSettingsAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var NtpSettingsAllDataSourceConfig = `
resource "powerscale_ntpsettings" "ntp_settings_test" {
	key_file = "/ifs/"
}

data "powerscale_ntpsettings" "all" {
	depends_on = [
		powerscale_ntpsettings.ntp_settings_test
	]
}
`
