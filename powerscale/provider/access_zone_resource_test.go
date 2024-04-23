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

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var mocker *mockey.Mocker
var createMocker *mockey.Mocker

func TestAccAccessZoneA(t *testing.T) {
	var accessZoneResourceName = "powerscale_accesszone.zone"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + AccessZoneResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.#", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            accessZoneResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_auth_providers"},
			},
			// Update name, path and auth providers, then Read testing
			{
				Config: ProviderConfig + AccessZoneUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone4"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone4"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs/home"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.#", "1"),
				),
			},
			// Update name, add auth providers, then Read testing
			{
				Config: ProviderConfig + AccessZoneResourceConfigAddProvider,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone5"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone5"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs/home"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.#", "3"),
				),
			},
			// Update name, reorder auth providers, then Read testing
			{
				Config: ProviderConfig + AccessZoneResourceConfigReorderProvider,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone5-1"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone5-1"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs/home"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.0", "lsa-file-provider:System"),
				),
			},
			// Update to error state
			{
				Config:      ProviderConfig + AccessZoneUpdateResourceConfigError,
				ExpectError: regexp.MustCompile(".*Error editing access zone*."),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAccessZoneError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      ProviderConfig + AccessZoneErrorResourceConfig,
				ExpectError: regexp.MustCompile(".*Error creating access zone*."),
			},
		},
	})
}

func TestAccAccessZoneResourceAfterCreateGetErr(t *testing.T) {
	var accessZoneResourceName = "powerscale_accesszone.zone"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + AccessZoneResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.#", "2"),
				),
			},
			{
				PreConfig: func() {
					if mocker != nil {
						mocker.UnPatch()
					}
					mocker = mockey.Mock(helper.GetAllAccessZones).Return(nil, fmt.Errorf("access zone read mock error")).Build()
				},
				Config:      ProviderConfig + AccessZoneResourceConfig,
				ExpectError: regexp.MustCompile(`.*access zone read mock error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if mocker != nil {
				mocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccAccessZoneResourceGetErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if mocker != nil {
						mocker.UnPatch()
					}
					if createMocker != nil {
						createMocker.UnPatch()
					}
					createMocker = mockey.Mock(helper.CreateAccessZones).Return(nil).Build()
					mocker = mockey.Mock(helper.GetAllAccessZones).Return(nil, fmt.Errorf("access zone read mock error")).Build()
				},
				Config:      ProviderConfig + AccessZoneResourceConfig,
				ExpectError: regexp.MustCompile(`.*access zone read mock error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if mocker != nil {
				mocker.UnPatch()
			}
			if createMocker != nil {
				createMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccAccessZoneResourceGetImportErr(t *testing.T) {
	var accessZoneResourceName = "powerscale_accesszone.zone"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if mocker != nil {
						mocker.UnPatch()
					}
					if createMocker != nil {
						createMocker.UnPatch()
					}
				},
				Config: ProviderConfig + AccessZoneResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.#", "2"),
				),
			},
			{
				PreConfig: func() {
					if mocker != nil {
						mocker.UnPatch()
					}
					mocker = mockey.Mock(helper.GetAllAccessZones).Return(nil, fmt.Errorf("access zone read mock error")).Build()
				},
				ResourceName:      accessZoneResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*access zone read mock error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if createMocker != nil {
				createMocker.UnPatch()
			}
			if mocker != nil {
				mocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccAccessZoneResourceGetImportSpecificErr(t *testing.T) {
	var accessZoneResourceName = "powerscale_accesszone.zone"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if mocker != nil {
						mocker.UnPatch()
					}
					if createMocker != nil {
						createMocker.UnPatch()
					}
				},
				Config: ProviderConfig + AccessZoneResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessZoneResourceName, "name", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "id", "tfaccTestAccessZone3"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "path", "/ifs"),
					resource.TestCheckResourceAttr(accessZoneResourceName, "auth_providers.#", "2"),
				),
			},
			{
				PreConfig: func() {
					mocker = mockey.Mock(helper.GetSpecificZone).Return(nil, fmt.Errorf("access zone read specific mock error")).Build()
				},
				ResourceName:            accessZoneResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_auth_providers"},
				ExpectError:             regexp.MustCompile(`.*access zone read specific mock error*.`),
			},
		},
	})
}

var AccessZoneResourceConfig = `
resource "powerscale_accesszone" "zone" {
	# Required fields
	name = "tfaccTestAccessZone3"
	groupnet = "groupnet0"
	path = "/ifs"
  
	# Optional to apply Auth Providers
	custom_auth_providers = ["System"]
  }
`

var AccessZoneUpdateResourceConfig = `
resource "powerscale_accesszone" "zone" {
	# Required fields
	name = "tfaccTestAccessZone4"
	groupnet = "groupnet0"
	path = "/ifs/home"
  
	# Optional to apply Auth Providers
	custom_auth_providers = []
  }
`

var AccessZoneResourceConfigAddProvider = `
resource "powerscale_accesszone" "zone" {
	# Required fields
	name = "tfaccTestAccessZone5"
	groupnet = "groupnet0"
	path = "/ifs/home"
  
	# Optional to apply Auth Providers
	custom_auth_providers = ["lsa-local-provider:System", "lsa-file-provider:System"]
  }
`

var AccessZoneResourceConfigReorderProvider = `
resource "powerscale_accesszone" "zone" {
	# Required fields
	name = "tfaccTestAccessZone5-1"
	groupnet = "groupnet0"
	path = "/ifs/home"
  
	# Optional to apply Auth Providers
	custom_auth_providers = ["lsa-file-provider:System", "lsa-local-provider:System"]
  }
`

var AccessZoneErrorResourceConfig = `
resource "powerscale_accesszone" "zone" {
	# Required fields
	name = "tfaccAccessZoneError"
	groupnet = "groupnet0"
	path = "/ifs"
  
	# Optional to apply Auth Providers
	custom_auth_providers = ["System"]
  }
`
var AccessZoneUpdateResourceConfigError = `
resource "powerscale_accesszone" "zone" {
	# Required fields
	name = "tfaccTestAccessZone6"
	groupnet = "groupnet0"
	path = "/some/bad/path/lol"
  
	# Optional to apply Auth Providers
	custom_auth_providers = ["fakeAuthProvider"]
  }
`
