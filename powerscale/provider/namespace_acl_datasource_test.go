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
	"context"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNamespaceAclDataSource(t *testing.T) {
	var namespaceACLTerraformName = "data.powerscale_namespace_acl.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NamespaceACLDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(namespaceACLTerraformName, "owner.id", "UID:0"),
					resource.TestCheckResourceAttr(namespaceACLTerraformName, "group.id", "GID:0"),
					resource.TestCheckResourceAttr(namespaceACLTerraformName, "authoritative", "acl"),
					resource.TestCheckResourceAttr(namespaceACLTerraformName, "mode", "0775"),
					resource.TestCheckResourceAttr(namespaceACLTerraformName, "acl.#", "3"),
				),
			},
		},
	})
}

func TestAccNamespaceAclDataSourceMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNamespaceACLDatasource).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NamespaceACLDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + NamespaceACLDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNamespaceAclDataSourceParamErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:      ProviderConfig + NamespaceACLDataSourceConfigInvalid,
				ExpectError: regexp.MustCompile(`.*Error getting the namespace acl*.`),
			},
			{
				Config:      ProviderConfig + NamespaceACLDataSourceConfigEmpty,
				ExpectError: regexp.MustCompile(`.*Missing required argument*.`),
			},
		},
	})
}

var NamespaceACLDataSourceConfig = `
resource "powerscale_namespace_acl" "namespace_acl_test" {
	namespace = "ifs/home"
	nsaccess = true
	owner = { id = "UID:0"}
	group = { id = "GID:0"}
	acl_custom = [
	{
		accessrights = ["dir_gen_read","dir_gen_write","dir_gen_execute","std_write_dac","delete_child"]
		accesstype = "allow"
		inherit_flags = []
		trustee = {
			id = "UID:0"
		}
	},
	{
		accessrights = ["dir_gen_read","dir_gen_write","dir_gen_execute","delete_child"]
		accesstype = "allow"
		inherit_flags = []
		trustee = {
			id = "GID:0"
		}
	},
	{
		accessrights = ["dir_gen_read","dir_gen_execute"]
		accesstype = "allow"
		inherit_flags = []
		"trustee": {
			"id": "SID:S-1-1-0"
		}
	},
	]
}

data "powerscale_namespace_acl" "test" {
	filter {
		namespace = "ifs/home"
		nsaccess = true
	}
    depends_on = [
        powerscale_namespace_acl.namespace_acl_test
    ]
}
`

var NamespaceACLDataSourceConfigInvalid = `
data "powerscale_namespace_acl" "test" {
	filter {
		namespace = "invalid"
		nsaccess = true
	}
}
`

var NamespaceACLDataSourceConfigEmpty = `
data "powerscale_namespace_acl" "test" {
	filter {
		nsaccess = true
	}
}
`
