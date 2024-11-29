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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStoragepoolTierResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + StoragepoolTierResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_storagepool_tier.example", "name", "Sample_terraform_tier_1"),
					resource.TestCheckResourceAttr("powerscale_storagepool_tier.example", "transfer_limit_pct", "20"),
				),
			},
			{
				Config: ProviderConfig + StoragepoolTierResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_storagepool_tier.example", "name", "Sample_terraform_tier_2"),
				),
			},
		},
	})
}

func TestAccStoragepoolTierResourceModifyErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + StoragepoolTierResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateStoragepoolTier).Return(diags).Build()
				},
				Config:      ProviderConfig + StoragepoolTierResourceConfigUpdateErr2,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccStoragepoolTierResourceMockErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CreateStoragepoolTier).Return(diags).Build()
				},
				Config:      ProviderConfig + StoragepoolTierResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadStoragepoolTier).Return(diags).Build()
				},
				Config:      ProviderConfig + StoragepoolTierResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + StoragepoolTierResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccStoragepoolTierResourceImportMockErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + StoragepoolTierResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadStoragepoolTier).Return(diags).Build()
				},
				Config:            ProviderConfig + StoragepoolTierResourceConfig,
				ResourceName:      "powerscale_storagepool_tier.example",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:       ProviderConfig + StoragepoolTierResourceConfig,
				ResourceName: "powerscale_storagepool_tier.example",
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var StoragepoolTierResourceConfig = `
resource "powerscale_storagepool_tier" "example" {
    children = [
        "x410_34tb_1.6tb-ssd_64gb"
    ]
    name = "Sample_terraform_tier_1"
    transfer_limit_pct = 20
  }
`

var StoragepoolTierResourceConfigUpdate = `
resource "powerscale_storagepool_tier" "example" {
    children = []
    name = "Sample_terraform_tier_2"
    transfer_limit_state = "default"
  }
`

var StoragepoolTierResourceConfigCreateErr = `
resource "powerscale_storagepool_tier" "example" {
    children = [
        "x410_34tb_1.6tb-ssd_64gb"
    ]
	name = "Sample_terraform_tier_1"
    transfer_limit_pct = 110
  }
`

var StoragepoolTierResourceConfigUpdateErr = `
resource "powerscale_storagepool_tier" "example" {
	name = "Sample_terraform_tier_1"
    transfer_limit_pct = 110
  }
`

var StoragepoolTierResourceConfigUpdateErr2 = `
resource "powerscale_storagepool_tier" "example" {
    children = [
        "x410_34tb_1.6tb-ssd_64gb"
    ]
	name = "Sample_terraform_tier_1"
    transfer_limit_pct = 110
  }
`
