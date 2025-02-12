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

func TestAccNfsAliasDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + NfsAliasDatasourceGetAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_nfs_alias.alias_datasource_test", "nfs_aliases.#"),
				),
			},
		},
	})
}

func TestAccNfsAliasDatasourceGetWithQueryParam(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + NfsAliasDatasourceGetWithQueryParam,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_nfs_alias.export_datasource_test", "nfs_aliases.#"),
				),
			},
		},
	})
}

func TestAccNfsAliasDatasourceErrorList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ListNFSAliases).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NfsAliasDatasourceGetAllConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNfsAliasDatasourceErrorCopy(t *testing.T) {
	FunctionMocker2 = mockey.Mock(helper.ListNFSAliases).Return(nil, fmt.Errorf("mock error")).Build()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NfsAliasDatasourceGetErrorConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})

	FunctionMocker2.Release()
}

var NfsAliasDatasourceGetWithQueryParam = NfsAliasResourceConfig + `
data "powerscale_nfs_alias" "export_datasource_test" {
	depends_on = [powerscale_nfs_alias.example]
	filter {
        dir   = "ASC"
		limit = 10
        sort  = "name"
	}
}
`

var NfsAliasDatasourceGetAllConfig = NfsAliasResourceConfig + `
data "powerscale_nfs_alias" "alias_datasource_test" {
	depends_on = [powerscale_nfs_alias.example]
  	filter {
		limit = 1
	}
}
`
var NfsAliasDatasourceGetErrorConfig = `
data "powerscale_nfs_alias" "alias_datasource_test" {
	filter {
		limit = 1
	}
}`
