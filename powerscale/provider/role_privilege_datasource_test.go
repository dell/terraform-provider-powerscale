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

func TestAccRolePrivilegeDataSourceNames(t *testing.T) {
	var ntpServerTerraformName = "data.powerscale_roleprivilege.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Filter by names
			{
				Config: ProviderConfig + RolePrivilegeDataSourceNamesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.#", "1"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.category", "System"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.description", "Shutdown the system"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.id", "ISI_PRIV_SYS_SHUTDOWN"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.name", "Shutdown"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.parent_id", "ISI_PRIV_ZERO"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.permission", "r"),
					resource.TestCheckResourceAttr(ntpServerTerraformName, "role_privileges_details.0.privilegelevel", "flag"),
				),
			},
		},
	})
}

func TestAccRolePrivilegeDataSourceAll(t *testing.T) {
	var ntpServerTerraformName = "data.powerscale_roleprivilege.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + RolePrivilegeAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(ntpServerTerraformName, "role_privileges_details.#"),
				),
			},
		},
	})
}

func TestAccRolePrivilegeDataSourceNamesErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + RolePrivilegeDataSourceNameConfigErr,
				ExpectError: regexp.MustCompile(`.*No relevant role privileges are found*.`),
			},
		},
	})
}

func TestAccRolePrivilegeDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + RolePrivilegeDataSourceFilterConfigErr,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccRolePrivilegeDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetRolePrivileges).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RolePrivilegeAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccRolePrivilegeDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.RolePrivilegeDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RolePrivilegeAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var RolePrivilegeDataSourceNamesConfig = `
data "powerscale_roleprivilege" "test" {
	filter {
		names = ["Shutdown"]
	}
}
`

var RolePrivilegeAllDataSourceConfig = `
data "powerscale_roleprivilege" "all" {
}
`

var RolePrivilegeDataSourceNameConfigErr = `
data "powerscale_roleprivilege" "test" {
	filter {
		names = ["BadName"]
	}
}
`

var RolePrivilegeDataSourceFilterConfigErr = `
data "powerscale_roleprivilege" "test" {
	filter {
		invalidFilter = "badFilter"
	}
}
`
