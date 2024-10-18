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

func TestAccS3ZoneSettingResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Error in setzoneSetting
			{
				Config: ProviderConfig + testAccS3ZoneSettingConfig(),
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.SetZoneSetting).Return(nil, fmt.Errorf("create error")).Build()
				},
				ExpectError: regexp.MustCompile("create error"),
			},
			{
				PreConfig: func() { FunctionMocker.UnPatch() },
				Config:    ProviderConfig + testAccS3ZoneSettingConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "zone", "System"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "root_path", "/ifs"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "base_domain", "dell"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "object_acl_policy", "deny"),
				),
			},
			{
				Config: ProviderConfig + testAccS3ZoneSettingUpdateConfig(),
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetZoneSetting).Return(fmt.Errorf("read error")).Build()
				},
				ExpectError: regexp.MustCompile("read error"),
			},
			{
				PreConfig: func() { FunctionMocker.UnPatch() },
				Config:    ProviderConfig + testAccS3ZoneSettingUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "zone", "System"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "root_path", "/ifs"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "base_domain", "dell.com"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_zone_setting", "object_acl_policy", "replace"),
				),
			},
			{
				Config:        ProviderConfig + testAccS3ZoneSettingImportError(),
				ResourceName:  "powerscale_s3_zone_settings.s3_import",
				ImportState:   true,
				ImportStateId: "System",
			},

			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetZoneSetting).Return(fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + testAccS3ZoneSettingImportError(),
				ResourceName:      "powerscale_s3_zone_settings.s3_import",
				ImportState:       true,
				ImportStateId:     "System",
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func testAccS3ZoneSettingConfig() string {
	return `
resource "powerscale_s3_zone_settings" "s3_zone_setting" {
  zone                        = "System"
  root_path                   = "/ifs"
  base_domain                 = "dell"
  bucket_directory_create_mode = 511
  object_acl_policy           = "deny"
}`
}

func testAccS3ZoneSettingUpdateConfig() string {
	return `
resource "powerscale_s3_zone_settings" "s3_zone_setting" {
  zone                        = "System"
  root_path                   = "/ifs"
  base_domain                 = "dell.com"
  bucket_directory_create_mode = 511
  object_acl_policy           = "replace"
}`
}

func testAccS3ZoneSettingImportError() string {
	return `
	resource "powerscale_s3_zone_settings" "s3_import" {}
	`
}
