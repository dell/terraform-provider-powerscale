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

func TestAccNfsZoneSettingsDataSourceReadWithoutFilter(t *testing.T) {
	dataSourceName := "data.powerscale_nfs_zone_settings.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + nfsZoneSettingsDataSourceConfigWithoutFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_domain_uids"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_allow_numeric_ids"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_names"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_replace_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_domain"),
				),
			},
		},
	})
}

func TestAccNfsZoneSettingsDataSourceReadWithFilter(t *testing.T) {
	dataSourceName := "data.powerscale_nfs_zone_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + nfsZoneSettingsDataSourceConfigWithFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_domain_uids"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_allow_numeric_ids"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_names"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_replace_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_domain"),
				),
			},
		},
	})
}

func TestAccNfsZoneSettingsDataSourceReadMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.FilterNfsZoneSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsDataSourceConfigWithFilter,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsDataSourceConfigWithFilter,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}
func TestAccNfsZoneSettingsDataSourceReadEmptyFilter(t *testing.T) {
	dataSourceName := "data.powerscale_nfs_zone_settings.empty"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + nfsZoneSettingsDataSourceConfigWithEmptyFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_domain_uids"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_allow_numeric_ids"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_no_names"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_replace_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nfs_zone_settings.nfsv4_domain"),
				),
			},
		},
	})
}

var nfsZoneSettingsDataSourceConfigWithoutFilter = `
data "powerscale_nfs_zone_settings" "all" {
}
`

var nfsZoneSettingsDataSourceConfigWithFilter = `
data "powerscale_nfs_zone_settings" "test" {
	filter {
		zone = "System"
	}
}
`

var nfsZoneSettingsDataSourceConfigWithEmptyFilter = `
data "powerscale_nfs_zone_settings" "empty" {
	filter {
	}
}
`
