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

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNfsExportSettingsDataSource(t *testing.T) {
	var testName = "data.powerscale_nfs_export_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + settingsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testName, "id"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.%"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.all_dirs"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.case_insensitive"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.case_preserving"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.commit_asynchronous"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.no_truncate"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.write_datasync_action"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.security_flavors.#"),
				),
			},
		},
	})
}

func TestAccNfsExportSettingsDataSourceWithFilter(t *testing.T) {
	var testName = "data.powerscale_nfs_export_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + settingsDataSourceFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testName, "id"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.%"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.all_dirs"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.case_insensitive"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.case_preserving"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.commit_asynchronous"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.no_truncate"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.write_datasync_action"),
					resource.TestCheckResourceAttrSet(testName, "nfs_export_settings.security_flavors.#"),
				),
			},
		},
	})
}

func TestAccNfsExportSettingsDataSourceErrorGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.FilterNfsExportSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + settingsDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var settingsDataSourceConfig = `
data "powerscale_nfs_export_settings" "test" {
}
`

var settingsDataSourceFilterConfig = `
data "powerscale_nfs_export_settings" "test" {
  filter {
		scope = "effective"
	}
}
`
