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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupDataSourceFilter(t *testing.T) {
	var userGroupTerraformName = "data.powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + userGroupFilterDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(userGroupTerraformName, "user_groups.#"),
				),
			},
		},
	})
}

func TestAccUserGroupDataSourceNames(t *testing.T) {
	var userGroupTerraformName = "data.powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter by names read testing
			{
				Config: ProviderConfig + userGroupNamesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userGroupTerraformName, "user_groups.#", "3"),
				),
			},
		},
	})
}

func TestAccUserGroupDataSourceFilterNames(t *testing.T) {
	var userTerraformName = "data.powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with names read testing
			{
				Config: ProviderConfig + userGroupFilterNamesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userTerraformName, "user_groups.#", "1"),
					resource.TestCheckResourceAttr(userTerraformName, "user_groups.0.gid", "GID:10000"),
					resource.TestCheckResourceAttr(userTerraformName, "user_groups.0.name", "tfaccUserGroupDatasource"),
					resource.TestCheckResourceAttr(userTerraformName, "user_groups.0.members.#", "0"),
					resource.TestCheckResourceAttr(userTerraformName, "user_groups.0.roles.#", "0"),
				),
			},
		},
	})
}

func TestAccUserGroupDataSourceInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid names read testing
			{
				Config:      ProviderConfig + userGroupInvalidNamesDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting the list of PowerScale User Groups.*.`),
			},
		},
	})
}

func TestAccUserGroupDataSourceAll(t *testing.T) {
	var userGroupTerraformName = "data.powerscale_user_group.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + userGroupAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(userGroupTerraformName, "user_groups.#"),
				),
			},
		},
	})
}

var userGroupFilterDataSourceConfig = `
data "powerscale_user_group" "test" {
  filter {
	cached = false
	# name_prefix = ""
	# domain = ""
	# zone = ""
	# provider = ""
  }
}
`

var userGroupFilterNamesDataSourceConfig = `
resource "powerscale_user_group" "testDep" {
	name = "tfaccUserGroupDatasource"
	gid = 10000
}

data "powerscale_user_group" "test" {
  filter {
    names = [
		{
			name = "tfaccUserGroupDatasource"
			gid = 10000
		}
	]
  }
  depends_on = [
	powerscale_user_group.testDep
  ]
}
`

var userGroupNamesDataSourceConfig = `
resource "powerscale_user_group" "testDep" {
	name = "tfaccUserGroupDatasource"
	gid = 10000
}

data "powerscale_user_group" "test" {
  filter {
    names = [
		{
			gid = 0
		},
		{
			name = "Administrators"
		},
		{
			name = "tfaccUserGroupDatasource"
			gid = 10000
		}
	]
  }
  depends_on = [
	powerscale_user_group.testDep
  ]
}
`

var userGroupInvalidNamesDataSourceConfig = `
data "powerscale_user_group" "test" {
  filter {
    names = [
		{
			gid = 0
		},
		{
			name = "invalidUserGroup"
		}
	]
  }
}
`

var userGroupAllDataSourceConfig = `
data "powerscale_user_group" "all" {
}
`
