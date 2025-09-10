/*
Copyright (c) 2024-2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"fmt"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccNfsAliasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NfsAliasResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_nfs_alias.example", "name", "/NfsAlias"),
				),
			},
			// ImportState testing
			{
				ResourceName:  "powerscale_nfs_alias.example",
				ImportState:   true,
				ImportStateId: "dev-tcz:NfsAlias",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, aliasName, states[0].Attributes["name"])
					assert.Equal(t, zone, states[0].Attributes["zone"])
					return nil
				},
			},
			// Update testing
			{
				Config: ProviderConfig + NfsAliasResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_nfs_alias.example", "name", "/NfsAlias_Update"),
				),
			},
		},
	})
}

func TestAccNfsAliasResourceCreateErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + NfsAliasResourceConfigCreateErr,
				ExpectError: regexp.MustCompile(`.*Error creating nfs alias*.`),
			},
		},
	})
}

func TestAccNfsAliasResourceModifyErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NfsAliasResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config:      ProviderConfig + NfsAliasResourceConfigUpdateErr,
				ExpectError: regexp.MustCompile(`.*Error updating nfs alias*.`),
			},
			{
				Config:      ProviderConfig + NfsAliasResourceConfigUpdatePathErr,
				ExpectError: regexp.MustCompile(`.*Error updating nfs alias*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsAlias).Return(diags).Build()
				},
				Config:      ProviderConfig + NfsAliasResourceConfigUpdateErr2,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsAliasResourceMockErr(t *testing.T) {
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
					FunctionMocker = mockey.Mock(helper.CreateNfsAlias).Return(diags).Build()
				},
				Config:      ProviderConfig + NfsAliasResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadNfsAlias).Return(diags).Build()
				},
				Config:      ProviderConfig + NfsAliasResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},

			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NfsAliasResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsAliasResourceImportMockErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NfsAliasResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadNfsAlias).Return(diags).Build()
				},
				Config:            ProviderConfig + NfsAliasResourceConfig,
				ResourceName:      "powerscale_nfs_alias.example",
				ImportState:       true,
				ImportStateId:     "dev-tcz:NfsAlias",
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
				Config:        ProviderConfig + NfsAliasResourceConfig,
				ResourceName:  "powerscale_nfs_alias.example",
				ImportStateId: "dev-tcz:NfsAlias",
				ImportState:   true,
				ExpectError:   regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

// alias exist in other than system zone and you are not providing req id with zone:alias_name
// by default it will look into system zone
func TestAccNfsAliasResourceImportErr(t *testing.T) {
	var diags diag.Diagnostics
	diags.AddError("mock error", "mock error")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import and read Error testing
			{
				Config:            ProviderConfig + NfsAliasResourceConfig,
				ResourceName:      "powerscale_nfs_alias.example",
				ImportState:       true,
				ImportStateId:     "NfsAlias",
				ExpectError:       regexp.MustCompile(`Error reading nfs alias*.`),
				ImportStateVerify: true,
			},
		},
	})
}

var aliasName = "/NfsAlias"
var zone = "dev-tcz"

var NfsAliasResourceConfig = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias"
   path = "/ifs/data"
   zone = "dev-tcz"
  }
`

var NfsAliasResourceConfigUpdate = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias_Update"
   path = "/ifs/data"
   zone = "dev-tcz"
  }
`

var NfsAliasResourceConfigCreateErr = `
resource "powerscale_nfs_alias" "example" {
   name = "NfsAlias"
   path = "/ifs/data"
   zone = "Invalid"
  }
`

var NfsAliasResourceConfigUpdateErr = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias"
   path = "/ifs/data"
   zone = "Update"
  }
`
var NfsAliasResourceConfigUpdatePathErr = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias"
   path = "/ifs"
   zone = "dev-tcz"
  }
`

var NfsAliasResourceConfigUpdateErr2 = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias2"
   path = "/ifs/data"
   zone = "dev-tcz"
  }
`
