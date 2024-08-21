package provider

import (
	"fmt"
	"regexp"
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
					FunctionMocker = mockey.Mock(SetGlobalSetting).Return(nil, fmt.Errorf("create error")).Build()
				},
				ExpectError: regexp.MustCompile("create error"),
			},
			{
				Config: ProviderConfig + testAccS3GlobalSettingUpdateConfig(),
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(GetGlobalSetting).Return(fmt.Errorf("read error")).Build()
				},
				ExpectError: regexp.MustCompile("read error"),
			},
			{
				PreConfig: func() { FunctionMocker.UnPatch() },
				Config:    ProviderConfig + testAccS3GlobalSettingConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "service", "true"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_only", "false"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "http_port", "9097"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_port", "9098"),
				),
			},
			{
				Config: ProviderConfig + testAccS3GlobalSettingUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "service", "false"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_only", "true"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "http_port", "9099"),
					resource.TestCheckResourceAttr("powerscale_s3_global_settings.s3_global_setting", "https_port", "9100"),
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
					FunctionMocker = mockey.Mock(GetGlobalSetting).Return(fmt.Errorf("mock error")).Build()
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
	service = true
	https_only = false
	http_port = 9097
	https_port = 9098
}`
}

func testAccS3GlobalSettingUpdateConfig() string {
	return `
resource "powerscale_s3_global_settings" "s3_global_setting" {
	service = false
	https_only = true
	http_port = 9099
	https_port = 9100
}`
}

func testAccS3GlobalSettingImportError() string {
	return `
	resource "powerscale_s3_global_settings" "s3_import" {}
	`
}
