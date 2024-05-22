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

func TestAccSmbShareSettingsResourceCreate(t *testing.T) {
	resourceName := "powerscale_smb_share_settings.share_settings_test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + SmbShareSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_based_enumeration", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_based_enumeration_root_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "allow_delete_readonly", "false"),
					resource.TestCheckResourceAttr(resourceName, "ca_timeout", "120"),
				),
			},
		},
	})
}

func TestAccSmbShareSettingsResourceImport(t *testing.T) {
	resourceName := "powerscale_smb_share_settings.share_settings_test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + SmbShareSettingsResourceConfig,
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

func TestAccSmbShareSettingsResourceUpdate(t *testing.T) {
	resourceName := "powerscale_smb_share_settings.share_settings_test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + SmbShareSettingsResourceConfig,
			},
			// update and read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + SmbShareSettingsUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_based_enumeration", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_based_enumeration_root_only", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_delete_readonly", "false"),
					resource.TestCheckResourceAttr(resourceName, "ca_timeout", "60"),
				),
			},
		},
	})
}

func TestAccSmbShareSettingsResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSmbShareSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSmbShareSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}
func TestAccSmbShareSettingsResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmbShareSettingsResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSmbShareSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSmbShareSettingsResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmbShareSettingsResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsUpdatedResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSmbShareSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsUpdatedResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSmbShareSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsUpdatedResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareSettingsUpdatedResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSmbShareSettingsResourceErrorImport(t *testing.T) {
	resourceName := "powerscale_smb_share_settings.share_settings_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + SmbShareSettingsResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSmbShareSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + SmbShareSettingsResourceConfig,
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
				Config:            ProviderConfig + SmbShareSettingsResourceConfig,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var SmbShareSettingsResourceConfig = `
resource "powerscale_smb_share_settings" "share_settings_test" {
access_based_enumeration           = false
		access_based_enumeration_root_only = false
		allow_delete_readonly              = false
		ca_timeout                         = 120
		zone                               = "System"
}
`
var SmbShareSettingsUpdatedResourceConfig = `
resource "powerscale_smb_share_settings" "share_settings_test" {
	access_based_enumeration           = true
		access_based_enumeration_root_only = true
		allow_delete_readonly              = false
		ca_timeout                         = 60
		zone                               = "tfaccAccessZone"
}
`
