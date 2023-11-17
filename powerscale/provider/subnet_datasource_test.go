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
	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

func TestAccSubnetDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + SubnetDatasourceGetAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_subnet.subnet_datasource_test", "subnets.#"),
				),
			},
		},
	})
}

func TestAccSubnetDatasourceGetFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + SubnetDatasourceGetFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_subnet.subnet_datasource_test", "subnets.#"),
				),
			},
			// Read testing
			{
				Config: ProviderConfig + SubnetDatasourceGetFilterNameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_subnet.subnet_datasource_test", "subnets.#"),
				),
			},
		},
	})
}

func TestAccSubnetDatasourceGetFilterError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      ProviderConfig + SubnetDatasourceGetFilterErrorConfig,
				ExpectError: regexp.MustCompile("Error reading subnets"),
			},
		},
	})
}

func TestAccSubnetDatasourceGetPaginationError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ResumeSubnets).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SubnetDatasourceGetAllConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var SubnetDatasourceGetAllConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
}
`

var SubnetDatasourceGetFilterConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
  filter{
    groupnet_name="groupnet0"
  }
}
`

var SubnetDatasourceGetFilterNameConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
  filter{
    names=["subnet0"]
  }
}
`

var SubnetDatasourceGetFilterErrorConfig = `
data "powerscale_subnet" "subnet_datasource_test" {
  filter{
    groupnet_name="groupnet-non-existent"
  }
}
`
