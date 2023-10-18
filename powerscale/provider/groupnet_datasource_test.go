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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupnetDataSourceAll(t *testing.T) {
	var groupnetTerraformName = "data.powerscale_groupnet.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + groupnetAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(groupnetTerraformName, "groupnets.#"),
				),
			},
		},
	})
}

func TestAccGroupnetDataSourceFilterNames(t *testing.T) {
	var groupnetTerraformName = "data.powerscale_groupnet.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with names read testing
			{
				Config: ProviderConfig + groupnetFilterDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.#", "1"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.dns_resolver_rotate", "true"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.allow_wildcard_subdomains", "false"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.server_side_dns_search", "true"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.dns_cache_enabled", "true"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.name", "tfaccGroupnetDatasource"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.id", "tfaccGroupnetDatasource"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.description", "terraform groupnet datasource"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.dns_search.0", "pie.lab.emc.com"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.dns_servers.0", "10.230.44.169"),
					resource.TestCheckResourceAttr(groupnetTerraformName, "groupnets.0.subnets.#", "0"),
				),
			},
		},
	})
}

func TestAccGroupnetDataSourceInvalidNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid names read testing
			{
				Config:      ProviderConfig + groupnetInvalidFilterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error one or more of the filtered groupnet names is not a valid powerscale groupnet.*.`),
			},
		},
	})
}

func TestAccGroupnetDatasourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + groupnetAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccGroupnetDatasourceErrorGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetAllGroupnets).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + groupnetAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var groupnetFilterDataSourceConfig = `
data "powerscale_groupnet" "test" {
  filter {
    names = ["tfaccGroupnetDatasource"]
  }
}
`

var groupnetAllDataSourceConfig = `
data "powerscale_groupnet" "all" {
}
`
var groupnetInvalidFilterDataSourceConfig = `
data "powerscale_groupnet" "test" {
  filter {
    names = ["invalidName"]
  }
}
`
