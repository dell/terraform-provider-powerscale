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

func TestAccNfsZoneSettingsResourceCreate(t *testing.T) {
	resourceName := "powerscale_nfs_zone_settings.example"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "nfsv4_no_names"),
					resource.TestCheckResourceAttrSet(resourceName, "nfsv4_replace_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "nfsv4_allow_numeric_ids"),
					resource.TestCheckResourceAttrSet(resourceName, "nfsv4_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "nfsv4_no_domain"),
					resource.TestCheckResourceAttrSet(resourceName, "nfsv4_no_domain_uids"),
				),
			},
		},
	})
}

func TestAccNfsZoneSettingsResourceImport(t *testing.T) {
	resourceName := "powerscale_nfs_zone_settings.example"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigBasic,
			},
			// import testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNfsZoneSettingsResourceUpdate(t *testing.T) {
	resourceName := "powerscale_nfs_zone_settings.example"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigBasic,
			},
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigNewDomain,
			},
			// update and read testing
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigUpdatedDomain,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nfsv4_no_names", "false"),
					resource.TestCheckResourceAttr(resourceName, "nfsv4_replace_domain", "true"),
					resource.TestCheckResourceAttr(resourceName, "nfsv4_allow_numeric_ids", "true"),
					resource.TestCheckResourceAttr(resourceName, "nfsv4_domain", "localdomain_Updated"),
					resource.TestCheckResourceAttr(resourceName, "nfsv4_no_domain", "false"),
					resource.TestCheckResourceAttr(resourceName, "nfsv4_no_domain_uids", "true"),
				),
			},
		},
	})
}

func TestAccNfsZoneSettingsCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsZoneSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetNfsZoneSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsZoneSettingsReadMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigBasic,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNfsZoneSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsZoneSettingsUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigBasic,
			},
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigNewDomain,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigUpdatedDomain,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsZoneSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigUpdatedDomain,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetNfsZoneSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigUpdatedDomain,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsZoneSettingsResourceConfigUpdatedDomain,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsZoneSettingsImportMockErr(t *testing.T) {
	resourceName := "powerscale_nfs_zone_settings.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsZoneSettingsResourceConfigBasic,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNfsZoneSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + nfsZoneSettingsResourceConfigBasic,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var nfsZoneSettingsResourceConfigBasic = `
resource "powerscale_nfs_zone_settings" "example" {
	zone = "tfaccAccessZone"
}
`

var nfsZoneSettingsResourceConfigNewDomain = `
resource "powerscale_nfs_zone_settings" "example" {
	zone = "tfaccAccessZone"
	nfsv4_domain = "localdomain_New"
}
`

var nfsZoneSettingsResourceConfigUpdatedDomain = `
resource "powerscale_nfs_zone_settings" "example" {
	zone = "tfaccAccessZone"
	nfsv4_no_names = false
	nfsv4_replace_domain = true
	nfsv4_allow_numeric_ids = true
	nfsv4_domain = "localdomain_Updated"
	nfsv4_no_domain = false
	nfsv4_no_domain_uids = true
}
`
