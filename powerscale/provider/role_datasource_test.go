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

func TestAccRoleDataSourceNames(t *testing.T) {
	var roleTerraformName = "data.powerscale_role.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Filter by names
			{
				Config: ProviderConfig + RoleDataSourceNamesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.name", "SystemAdmin"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.id", "SystemAdmin"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.members.0.name", "admin"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.members.0.id", "UID:10"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.members.0.type", "user"),
					resource.TestCheckResourceAttrSet(roleTerraformName, "roles_details.0.privileges.#"),
				),
			},
		},
	})
}

func TestAccRoleDataSourceFilter(t *testing.T) {
	var roleTerraformName = "data.powerscale_role.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter by api filter
			{
				Config: ProviderConfig + RoleDataSourceFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.name", "SystemAdmin"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.id", "SystemAdmin"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.members.0.name", "admin"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.members.0.id", "UID:10"),
					resource.TestCheckResourceAttr(roleTerraformName, "roles_details.0.members.0.type", "user"),
					resource.TestCheckResourceAttrSet(roleTerraformName, "roles_details.0.privileges.#"),
				),
			},
		},
	})
}

func TestAccRoleDataSourceAll(t *testing.T) {
	var roleTerraformName = "data.powerscale_role.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + RoleAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(roleTerraformName, "roles_details.#"),
				),
			},
		},
	})
}

func TestAccRoleDataSourceNamesErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + RoleDataSourceNameConfigErr,
				ExpectError: regexp.MustCompile(`.*not a valid powerscale role*.`),
			},
		},
	})
}

func TestAccRoleDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + RoleDataSourceFilterConfigErr,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccRoleDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetRoles).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RoleAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccRoleDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.RoleDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RoleAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var RoleDataSourceNamesConfig = `
data "powerscale_role" "test" {
	filter {
		names = ["SystemAdmin"]
	}
}
`

var RoleDataSourceFilterConfig = `
data "powerscale_role" "test" {
	filter {
		zone = "System"
		names = ["SystemAdmin"]
	}
}
`

var RoleAllDataSourceConfig = `
data "powerscale_role" "all" {
}
`

var RoleDataSourceNameConfigErr = `
data "powerscale_role" "test" {
	filter {
		names = ["BadName"]
	}
}
`

var RoleDataSourceFilterConfigErr = `
data "powerscale_role" "test" {
	filter {
		invalidFilter = "badFilter"
	}
}
`
