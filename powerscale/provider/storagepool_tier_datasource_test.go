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

func TestAccStoragepoolTierDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + StoragepoolTierDatasourceAllWithResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_storagepool_tier.all_test", "storagepool_tiers.#"),
				),
			},
		},
	})
}

func TestAccStoragepoolTierDatasourceErrorList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + StoragepoolTierDatasourceAllConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetAllStoragepoolTiers).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.Release()
}

func TestAccStoragepoolTierDatasourceErrorCopy(t *testing.T) {
	FunctionMocker2 = mockey.Mock(helper.GetAllStoragepoolTiers).Return(nil, fmt.Errorf("mock error")).Build()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + StoragepoolTierDatasourceAllConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker2.Release()
}

var StoragepoolTierDatasourceAllWithResourceConfig = StoragepoolTierResourceConfigForDatasource + `

data "powerscale_storagepool_tier" "all_test" {
	depends_on = [powerscale_storagepool_tier.example]
}
`

var StoragepoolTierResourceConfigForDatasource = `
resource "powerscale_storagepool_tier" "example" {
    name = "Sample_terraform_tier_1"
    transfer_limit_pct = 20
}
`

var StoragepoolTierDatasourceAllConfig = `
data "powerscale_storagepool_tier" "all_test" {
	
}
`
