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
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessZoneDataSource(t *testing.T) {
	var azTerraformName = "data.powerscale_accesszone.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + AzDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.#", "1"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.alternate_system_provider", "lsa-file-provider:MinimumRequired"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.id", "tfaccAccessZone"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.negative_cache_entry_expiry", "60"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.path", "/ifs"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.skeleton_directory", "/usr/share/skel"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.system", "false"),
					resource.TestCheckResourceAttr(azTerraformName, "access_zones_details.0.system_provider", "lsa-file-provider:System"),
				),
			},
		},
	})
}

func TestAccAccessZoneDataSourceAll(t *testing.T) {
	var azTerraformName = "data.powerscale_accesszone.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + AzAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(azTerraformName, "access_zones_details.#"),
				),
			},
		},
	})
}
func TestAccAccessZoneDataSourceErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + AzDataSourceConfigErr,
				ExpectError: regexp.MustCompile(`.*not a valid powerscale access zone*.`),
			},
		},
	})
}

func TestAccAccessZoneDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + AzDataSourceConfigErrFilter,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccAccessZoneDataSourceGetErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetAllAccessZones).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + AzAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}
func TestQueryZoneNameByID(t *testing.T) {
	var azs []powerscale.V3ZoneExtended
	for i := 0; i < 10; i++ {
		idInt := int32(i)
		zoneName := fmt.Sprintf("Zone%d", i)
		azs = append(azs, powerscale.V3ZoneExtended{
			ZoneId: &idInt,
			Name:   &zoneName,
		})
	}

	zones := &powerscale.V3Zones{Zones: azs}
	if FunctionMocker != nil {
		FunctionMocker.UnPatch()
	}
	FunctionMocker = mockey.Mock(helper.GetAllAccessZones).Return(zones, nil).Build()
	defer FunctionMocker.UnPatch()
	name, err := helper.QueryZoneNameByID(context.Background(), nil, 0)
	assert.Nil(t, err)
	assert.Equal(t, "Zone0", name)
}

var AzDataSourceConfig = tfaccAccessZoneConfig + `
data "powerscale_accesszone" "test" {
  filter {
    names = ["tfaccAccessZone"]
  }

  depends_on = [powerscale_accesszone.tfaccAccessZone]
}
output "powerscale_accesszone" {
	value = data.powerscale_accesszone.test
}
`

var AzDataSourceConfigErr = tfaccAccessZoneConfig + `
data "powerscale_accesszone" "test" {
  filter {
    names = ["BadName"]
  }

  depends_on = [powerscale_accesszone.tfaccAccessZone]
}
output "powerscale_accesszone" {
	value = data.powerscale_accesszone.test
}
`

var AzDataSourceConfigErrFilter = tfaccAccessZoneConfig + `
data "powerscale_accesszone" "test" {
  filter {
    names = ["BadName"]
	invalidFilter = "badFilter"
  }

  depends_on = [powerscale_accesszone.tfaccAccessZone]
}
output "powerscale_accesszone" {
	value = data.powerscale_accesszone.test
}
`

var AzAllDataSourceConfig = `
data "powerscale_accesszone" "all" {
}
`
