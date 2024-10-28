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
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"fmt"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNfsAliasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NfsAliasResourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
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
				Config:      ProviderConfig + NfsAliasResourceConfigUpdateErr2,
				ExpectError: regexp.MustCompile(`.*Error updating nfs alias*.`),
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
				Config:       ProviderConfig + NfsAliasResourceConfig,
				ResourceName: "powerscale_nfs_alias.example",
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var NfsAliasResourceConfig = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias"
   path = "/ifs"
   zone = "System"
  }
`

var NfsAliasResourceConfigCreateErr = `
resource "powerscale_nfs_alias" "example" {
   name = "NfsAlias"
   path = "/ifs"
   zone = "Invalid"
  }
`

var NfsAliasResourceConfigUpdateErr = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias"
   path = "/ifs"
   zone = "Update"
  }
`

var NfsAliasResourceConfigUpdateErr2 = `
resource "powerscale_nfs_alias" "example" {
   name = "/NfsAlias"
   path = "/Invalid"
   zone = "System"
  }
`
