/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSyncIQGlobalSettingsDataSourceAll(t *testing.T) {
	var globalSettingsTerraformName = "data.powerscale_synciq_global_settings.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + SyncIQSettingsDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(globalSettingsTerraformName, "preferred_rpo_alert", "3"),
					resource.TestCheckResourceAttr(globalSettingsTerraformName, "service", "paused"),
					resource.TestCheckResourceAttr(globalSettingsTerraformName, "rpo_alerts", "true"),
					resource.TestCheckResourceAttr(globalSettingsTerraformName, "restrict_target_network", "true"),
				),
			},
		},
	})
}

func TestAccGlobalSettingsDataSourceReadingErr(t *testing.T) {
	diags := diag.Diagnostics{}
	diags.AddError("mock err", "mock err")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ManageReadDataSourceSyncIQGlobalSettings).Return(diags).Build()
				},
				Config:      ProviderConfig + SyncIQSettingsDatasourceConfig,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
		},
	})
}

var SyncIQSettingsDatasourceConfig = `
resource "powerscale_synciq_global_settings" "test" {
	preferred_rpo_alert = 3
	report_email = ["mail1@dell.com", "mail2@dell.com"]
   	restrict_target_network = true
   	rpo_alerts           = true
   	service              = "paused"
}

data "powerscale_synciq_global_settings" "all" {
    depends_on = [
        powerscale_synciq_global_settings.test
    ]
}
`
