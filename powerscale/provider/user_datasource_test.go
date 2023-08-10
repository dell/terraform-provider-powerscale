/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

func TestAccUserDataSourceFilter(t *testing.T) {
	var userTerraformName = "data.powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + userFilterDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(userTerraformName, "users.#"),
				),
			},
		},
	})
}

func TestAccUserDataSourceNames(t *testing.T) {
	var userTerraformName = "data.powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter by names read testing
			{
				Config: ProviderConfig + userNamesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userTerraformName, "users.#", "3"),
				),
			},
		},
	})
}

func TestAccUserDataSourceFilterNames(t *testing.T) {
	var userTerraformName = "data.powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with names read testing
			{
				Config: ProviderConfig + userFilterNamesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userTerraformName, "users.#", "1"),
					resource.TestCheckResourceAttr(userTerraformName, "users.0.uid", "UID:10000"),
					resource.TestCheckResourceAttr(userTerraformName, "users.0.name", "tfaccUserDatasource"),
					resource.TestCheckResourceAttr(userTerraformName, "users.0.roles.#", "0"),
				),
			},
		},
	})
}

func TestAccUserDataSourceInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid filter read testing
			{
				Config:      ProviderConfig + userInvalidFilterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*error getting users*.`),
			},
		},
	})
}

func TestAccUserDataSourceInvalidNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid names read testing
			{
				Config:      ProviderConfig + userInvalidNamesDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error one or more of the filtered user names is not a valid powerscale user*.`),
			},
		},
	})
}

func TestAccUserDataSourceAll(t *testing.T) {
	var userTerraformName = "data.powerscale_user.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + userAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(userTerraformName, "users.#"),
				),
			},
		},
	})
}

var userFilterDataSourceConfig = `
data "powerscale_user" "test" {
  filter {
	cached = false
	resolve_names = false
	member_of = false
	# domain = ""
	# zone = ""
	# provider = ""
  }
}
`

var userFilterNamesDataSourceConfig = `
data "powerscale_user" "test" {
  filter {
    names = [
		{
			name = "tfaccUserDatasource"
			uid = 10000
		}
	]
	cached = false
	name_prefix = "tfacc"
	resolve_names = false
	member_of = false
	# domain = "testDomain"
	# zone = "testZone"
	# provider = "testProvider"
  }
}
`

var userNamesDataSourceConfig = `
data "powerscale_user" "test" {
  filter {
    names = [
		{
			uid = 0
		},
		{
			name = "admin"
		},
		{
			name = "tfaccUserDatasource"
			uid = 10000
		}
	]
  }
}
`
var userInvalidFilterDataSourceConfig = `
data "powerscale_user" "test" {
	filter {
	  names = [
		  {
			  name = "tfaccUserDatasource"
			  uid = 10000
		  }
	  ]
	  cached = false
	  name_prefix = "tfacc"
	  resolve_names = false
	  member_of = false
	  domain = " "
	  zone = " "
	  provider = " "
	}
  }
`

var userInvalidNamesDataSourceConfig = `
data "powerscale_user" "test" {
  filter {
    names = [
		{
			uid = 0
		},
		{
			name = "invalidUser"
		}
	]
  }
}
`

var userAllDataSourceConfig = `
data "powerscale_user" "all" {
}
`
