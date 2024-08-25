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
			// Error in setGlobalSetting
			{
				Config: ProviderConfig + testAccS3ZoneSettingConfig(),
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.SetGlobalSetting).Return(nil, fmt.Errorf("create error")).Build()
				},
				ExpectError: regexp.MustCompile("create error"),
			},
			{
				Config: ProviderConfig + testAccS3ZoneSettingUpdateConfig(),
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetGlobalSetting).Return(fmt.Errorf("read error")).Build()
				},
				ExpectError: regexp.MustCompile("read error"),
			},
			{
				PreConfig: func() { FunctionMocker.UnPatch() },
				Config:    ProviderConfig + testAccS3ZoneSettingConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "service", "true"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "https_only", "false"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "http_port", "9097"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "https_port", "9098"),
				),
			},
			{
				Config: ProviderConfig + testAccS3ZoneSettingUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "service", "false"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "https_only", "true"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "http_port", "9099"),
					resource.TestCheckResourceAttr("powerscale_s3_zone_settings.s3_global_setting", "https_port", "9100"),
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
					FunctionMocker = mockey.Mock(helper.GetGlobalSetting).Return(fmt.Errorf("mock error")).Build()
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
resource "powerscale_s3_zone_settings" "s3_global_setting" {
  zone                        = "System"
  root_path                   = "/data/s3_buckets"
  base_domain                 = "dell.com"
  bucket_directory_create_mode = 511
  object_acl_policy           = "replace"
}`
}

func testAccS3ZoneSettingUpdateConfig() string {
	return `
resource "powerscale_s3_zone_settings" "s3_global_setting" {
  zone                        = "System"
  root_path                   = "/data/s3_buckets/zone"
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
