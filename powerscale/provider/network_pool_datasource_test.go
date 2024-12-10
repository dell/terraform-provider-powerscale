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

func TestAccNetworkPoolDataSourceNames(t *testing.T) {
	var poolTerraformName = "data.powerscale_networkpool.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Filter by names
			{
				Config: ProviderConfig + PoolDataSourceNamesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.#", "1"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.access_zone", "System"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.subnet", "subnet0"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.sc_ttl", "0"),
					resource.TestCheckResourceAttrSet(poolTerraformName, "network_pools_details.0.ifaces.#"),
				),
			},
		},
	})
}

func TestAccNetworkPoolDataSourceFilter(t *testing.T) {
	var poolTerraformName = "data.powerscale_networkpool.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter by api filter
			{
				Config: ProviderConfig + PoolDataSourceFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.#", "1"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.access_zone", "System"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.subnet", "subnet0"),
					resource.TestCheckResourceAttr(poolTerraformName, "network_pools_details.0.sc_ttl", "0"),
					resource.TestCheckResourceAttrSet(poolTerraformName, "network_pools_details.0.ifaces.#"),
				),
			},
		},
	})
}

func TestAccNetworkPoolDataSourceAll(t *testing.T) {
	var poolTerraformName = "data.powerscale_networkpool.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all network pools
			{
				Config: ProviderConfig + PoolAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(poolTerraformName, "network_pools_details.#"),
				),
			},
		},
	})
}

func TestAccNetworkPoolDataSourceNamesErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + PoolDataSourceNameConfigErr,
				ExpectError: regexp.MustCompile(`.*not a valid powerscale network pool*.`),
			},
		},
	})
}

func TestAccNetworkPoolDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + PoolDataSourceFilterConfigErr,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccNetworkPoolDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNetworkPools).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + PoolAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNetworkPoolDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.NetworkPoolDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + PoolAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var PoolDataSourceNamesConfig = `
data "powerscale_networkpool" "test" {
	filter {
		names = ["pool0"]
	}
}
`

var PoolDataSourceFilterConfig = `
data "powerscale_networkpool" "test" {
	filter {
		groupnet = "groupnet0"
		subnet = "subnet0"
		names = ["pool0"]
	}
}
`

var PoolAllDataSourceConfig = `
data "powerscale_networkpool" "all" {
}
`

var PoolDataSourceNameConfigErr = `
data "powerscale_networkpool" "test" {
	filter {
		names = ["BadName"]
	}
}
`

var PoolDataSourceFilterConfigErr = `
data "powerscale_networkpool" "test" {
	filter {
		invalidFilter = "badFilter"
	}
}
`
