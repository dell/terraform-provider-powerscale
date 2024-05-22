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
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"
)

func TestAccNamespaceAclResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.2", "dir_gen_execute"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "1"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.0", "container_inherit"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "powerscale_namespace_acl.namespace_acl_test",
				ImportStateId:                        namespace,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "namespace",
				ImportStateVerifyIgnore:              []string{"acl_custom", "nsaccess"},
			},
			// Update
			{
				Config: ProviderConfig + NamespaceACLUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "5"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.4", "delete_child"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "root"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
		},
	})
}

func TestAccNamespaceAclResourceEmptyConfig1(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NamespaceACLResourceEmptyConfig1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "5"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.4", "delete_child"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "root"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
		},
	})
}

func TestAccNamespaceAclResourceEmptyConfig2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NamespaceACLResourceEmptyConfig2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "0"),
				),
			},
		},
	})
}

func TestAccNamespaceAclResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.2", "dir_gen_execute"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "1"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.0", "container_inherit"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
			// ImportState testing get error
			{
				ResourceName:  "powerscale_namespace_acl.namespace_acl_test",
				ImportStateId: namespace,
				ImportState:   true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetNamespaceACL).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNamespaceAclResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.2", "dir_gen_execute"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "1"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.0", "container_inherit"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
			{
				Config:      ProviderConfig + NamespaceACLInvalidResourceConfig,
				ExpectError: regexp.MustCompile("Namespace acl param invalid"),
			},
			// Update param read error
			{
				Config: ProviderConfig + NamespaceACLUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config: ProviderConfig + NamespaceACLUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.CheckNamespaceACLParam).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + NamespaceACLUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.UpdateNamespaceACL).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNamespaceACL).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, model models.NamespaceACLResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNamespaceAclResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CheckNamespaceACLParam).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.UpdateNamespaceACL).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNamespaceAclResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NamespaceACLResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.2", "dir_gen_execute"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "1"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.0", "container_inherit"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:10"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
			{
				ResourceName:  "powerscale_namespace_acl.namespace_acl_test",
				ImportStateId: namespace,
				ImportState:   true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NamespaceACLResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + NamespaceACLUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNamespaceAclResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NamespaceACLResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Change the settings back to default
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + NamespaceACLUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "namespace", namespace),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "owner.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "group.id", "GID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.#", "3"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.#", "5"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.0", "dir_gen_read"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accessrights.4", "delete_child"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.accesstype", "allow"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.inherit_flags.#", "0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.id", "UID:0"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.name", "root"),
					resource.TestCheckResourceAttr("powerscale_namespace_acl.namespace_acl_test", "acl.2.trustee.type", "user"),
				),
			},
		},
	})
}

var namespace = "ifs/home"

var NamespaceACLResourceConfig = fmt.Sprintf(`
resource "powerscale_namespace_acl" "namespace_acl_test" {
	namespace = "%s"
	nsaccess = true
	owner = { id = "UID:10"}
	group = { id = "GID:10"}
	acl_custom = [
	{
		accessrights = ["dir_gen_read","dir_gen_write","dir_gen_execute"]
		accesstype = "allow"
		inherit_flags = ["container_inherit"]
		trustee = {
			id = "UID:10"
		}
	},
	{
		accessrights = ["dir_gen_read","dir_gen_write","dir_gen_execute"]
		accesstype = "allow"
		inherit_flags = ["container_inherit"]
		trustee = {
			name = "Isilon Users",
			type = "group"
		}
	},
	{
		accessrights = ["dir_gen_read"]
		accesstype = "allow"
		inherit_flags = ["container_inherit"]
		"trustee": {
			"id": "SID:S-1-1-0"
		}
	},
	]
}
`, namespace)

var NamespaceACLUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_namespace_acl" "namespace_acl_test" {
	namespace = "%s"
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
		accessrights = ["dir_gen_read","dir_gen_write"]
		accesstype = "allow"
		inherit_flags = []
		"trustee": {
			"id": "SID:S-1-1-0"
		}
	},
	]
}
`, namespace)

var NamespaceACLInvalidResourceConfig = fmt.Sprintf(`
resource "powerscale_namespace_acl" "namespace_acl_test" {
	namespace = "%s"
	nsaccess = true
	owner = { id = "UID:0"}
	group = { id = "GID:0"}
}
`, namespace)

var NamespaceACLResourceEmptyConfig1 = fmt.Sprintf(`
resource "powerscale_namespace_acl" "namespace_acl_test" {
	namespace = "%s"
}
`, namespace)

var NamespaceACLResourceEmptyConfig2 = fmt.Sprintf(`
resource "powerscale_namespace_acl" "namespace_acl_test" {
	namespace = "%s"
	nsaccess = true
	owner = { id = "UID:0"}
	group = { id = "GID:0"}
	acl_custom = []
}
`, namespace)
