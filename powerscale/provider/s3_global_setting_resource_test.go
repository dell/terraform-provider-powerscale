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

func TestAccS3GlobalSettingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error in setGlobalSetting
			{
				Config: ProviderConfig + testAccS3GlobalSettingConfig(),
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.SetGlobalSetting).Return(nil, fmt.Errorf("create error")).Build()
				},
				ExpectError: regexp.MustCompile("create error"),
			},
			{
				PreConfig: func() { FunctionMocker.UnPatch() },
				Config:    ProviderConfig + testAccS3GlobalSettingConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "service", "false"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_only", "false"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "http_port", "9097"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_port", "9098"),
				),
			},
			{
				Config: ProviderConfig + testAccS3GlobalSettingUpdateConfig(),
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetGlobalSetting).Return(fmt.Errorf("read error")).Build()
				},
				ExpectError: regexp.MustCompile("read error"),
			},
			{
				PreConfig: func() { FunctionMocker.UnPatch() },
				Config:    ProviderConfig + testAccS3GlobalSettingUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "service", "true"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_only", "true"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "http_port", "9020"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_port", "9021"),
				),
			},

			{
				Config:        ProviderConfig + testAccS3GlobalSettingImportError(),
				ResourceName:  "powerscale_s3_global_settings.s3_import",
				ImportState:   true,
				ImportStateId: "1",
			},

			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetGlobalSetting).Return(fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + testAccS3GlobalSettingImportError(),
				ResourceName:      "powerscale_s3_global_settings.s3_import",
				ImportState:       true,
				ImportStateId:     "1",
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func testAccS3GlobalSettingConfig() string {
	return `
resource "powerscale_s3_global_settings" "s3_global_setting" {
	service = false
	https_only = false
	http_port = 9097
	https_port = 9098
}`
}

func testAccS3GlobalSettingUpdateConfig() string {
	return `
resource "powerscale_s3_global_settings" "s3_global_setting" {
	service = true
	https_only = true
	http_port = 9020
	https_port = 9021
}`
}

func testAccS3GlobalSettingImportError() string {
	return `
	resource "powerscale_s3_global_settings" "s3_import" {}
	`
}
