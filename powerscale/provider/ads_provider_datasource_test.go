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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdsProviderDataSourceNames(t *testing.T) {
	var adsTerraformName = "data.powerscale_adsprovider.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Filter by names
			{
				Config: ProviderConfig + AdsDataSourceNamesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.#", "1"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.name", "PIE.LAB.EMC.COM"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.allocate_gids", "true"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.check_online_interval", "300"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.site", "Default-First-Site-Name"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.status", "online"),
					resource.TestCheckResourceAttrSet(adsTerraformName, "ads_providers_details.0.ignored_trusted_domains.#"),
				),
			},
		},
	})
}

func TestAccAdsProviderDataSourceFilter(t *testing.T) {
	var adsTerraformName = "data.powerscale_adsprovider.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter by api filter
			{
				Config: ProviderConfig + AdsDataSourceFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.#", "1"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.name", "PIE.LAB.EMC.COM"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.site", "Default-First-Site-Name"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.status", "online"),
					resource.TestCheckResourceAttr(adsTerraformName, "ads_providers_details.0.zone_name", "System"),
				),
			},
		},
	})
}

func TestAccAdsProviderDataSourceAll(t *testing.T) {
	var adsTerraformName = "data.powerscale_adsprovider.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + AdsAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(adsTerraformName, "ads_providers_details.#"),
				),
			},
		},
	})
}

func TestAccAdsProviderDataSourceNamesErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + AdsDataSourceNameConfigErr,
				ExpectError: regexp.MustCompile(`.*not a valid powerscale ads provider*.`),
			},
		},
	})
}

func TestAccAdsProviderDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + AdsDataSourceFilterConfigErr,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccAdsProviderDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.AdsProviderDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + AdsAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var AdsDataSourceNamesConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "PIE.LAB.EMC.COM"
	user = "administrator"
	password = "Password123!"
}

data "powerscale_adsprovider" "test" {
	filter {
		names = ["PIE.LAB.EMC.COM"]
	}
	depends_on = [
		powerscale_adsprovider.ads_test
	]
}
`

var AdsDataSourceFilterConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "PIE.LAB.EMC.COM"
	user = "administrator"
	password = "Password123!"
}

data "powerscale_adsprovider" "test" {
	filter {
		scope = "user"
	}
	depends_on = [
		powerscale_adsprovider.ads_test
	]
}
`

var AdsAllDataSourceConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "PIE.LAB.EMC.COM"
	user = "administrator"
	password = "Password123!"
}

data "powerscale_adsprovider" "all" {
	depends_on = [
		powerscale_adsprovider.ads_test
	]
}
`

var AdsDataSourceNameConfigErr = `
resource "powerscale_adsprovider" "ads_test" {
	name = "PIE.LAB.EMC.COM"
	user = "administrator"
	password = "Password123!"
}

data "powerscale_adsprovider" "test" {
	filter {
		names = ["BadName"]
	}
	depends_on = [
		powerscale_adsprovider.ads_test
	]
}
`

var AdsDataSourceFilterConfigErr = `
resource "powerscale_adsprovider" "ads_test" {
	name = "PIE.LAB.EMC.COM"
	user = "administrator"
	password = "Password123!"
}

data "powerscale_adsprovider" "test" {
	filter {
		invalidFilter = "badFilter"
	}
	depends_on = [
		powerscale_adsprovider.ads_test
	]
}
`
