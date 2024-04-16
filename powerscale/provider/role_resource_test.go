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
	powerscale "dell/powerscale-go-client"
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"
)

func TestAccRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.#", "2"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.0.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.#", "1"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.name", "Support"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.permission", "r"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_role.role_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, roleName, states[0].Attributes["name"])
					assert.Equal(t, roleDescription, states[0].Attributes["description"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + RoleUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription+"_modified"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.#", "1"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.0.name", "root"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.#", "2"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.name", "Shutdown"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.permission", "r"),
				),
			},
		},
	})
}

func TestAccRoleResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.#", "2"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.0.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.#", "1"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.name", "Support"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.permission", "r"),
				),
			},
			// ImportState testing get none role
			{
				ResourceName: "powerscale_role.role_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetRole).Return(&powerscale.V14AuthRolesExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile("not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_role.role_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetRole).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				ResourceName: "powerscale_role.role_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetRole).Return(&powerscale.V14AuthRolesExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("not found"),
			},
			{
				ResourceName: "powerscale_role.role_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetRole).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccRoleResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + RoleInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccRoleResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.#", "2"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.0.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.#", "1"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.name", "Support"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.permission", "r"),
				),
			},
			{
				ResourceName: "powerscale_role.role_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RoleResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + RoleUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccRoleResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RoleResourceConfig,
				ExpectError: regexp.MustCompile("Could not read role param with error"),
			},
		},
	})
}

func TestAccRoleResourceReorderMemberError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + RoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription),
				),
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReorderRoleMembers).Return(types.List{}, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RoleUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("Could not reorder role members"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.ReorderRoleMembers).Return(types.List{}, fmt.Errorf("mock error")).Build().
						When(func(localMembers types.List, remoteMembers types.List) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + RoleUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("Could not reorder role members"),
			},
		},
	})
}

func TestAccRoleResourceReorderPrivilegeError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + RoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription),
				),
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReorderRolePrivileges).Return(types.List{}, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RoleUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("Could not reorder role privileges"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.ReorderRolePrivileges).Return(types.List{}, fmt.Errorf("mock error")).Build().
						When(func(localMembers types.List, remoteMembers types.List) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + RoleUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("Could not reorder role privileges"),
			},
		},
	})
}

func TestAccRoleResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RoleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_role.role_test", "name", roleName),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "description", roleDescription),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.#", "2"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "members.0.name", "admin"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.#", "1"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.name", "Support"),
					resource.TestCheckResourceAttr("powerscale_role.role_test", "privileges.0.permission", "r"),
				),
			},
			// Update param read error
			{
				Config: ProviderConfig + RoleUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("Could not read role param with error"),
			},
			// Update get error
			{
				Config: ProviderConfig + RoleUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.UpdateRole).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none role
			{
				Config: ProviderConfig + RoleUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetRole).Return(&powerscale.V14AuthRolesExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".*Could not read updated role*."),
			},
			// Update get error
			{
				Config: ProviderConfig + RoleUpdatedResourceConfig2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetRole).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, roleModel models.RoleResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			//Update Invalid Config
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:      ProviderConfig + RoleInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

var roleName = "role_at_test"
var roleDescription = "role_description"

var RoleResourceConfig = fmt.Sprintf(`
resource "powerscale_role" "role_test" {
	name = "%s"
	description = "%s"
	members = [
		{
			name = "admin"
		},
		{
			id = "UID:0"
		}
	]
	privileges = [
		{
			id = "ISI_PRIV_SYS_SUPPORT",
			permission = "r"
		}
	]
}
`, roleName, roleDescription)

var RoleInvalidResourceConfig = fmt.Sprintf(`
resource "powerscale_role" "role_test" {
	name = "%s"
	description = "%s"
	members = [
		{
			name = "admin"
		},
		{
			id = "UID:0"
		}
	]
	privileges = [
		{
			id = "INVALID_PRIVILEGE",
			permission = "r"
		}
	]
}
`, roleName, roleDescription)

var RoleUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_role" "role_test" {
	name = "%s"
	description = "%s"
	members = [
		{
			id = "UID:0"
		}
	]
	privileges = [
		{
			id = "ISI_PRIV_SYS_SHUTDOWN",
			permission = "r"
		},
		{
			id = "ISI_PRIV_SYS_SUPPORT",
			permission = "r"
		}
	]
}
`, roleName, roleDescription+"_modified")

var RoleUpdatedResourceConfig2 = fmt.Sprintf(`
resource "powerscale_role" "role_test" {
	name = "%s"
	description = "%s"
	members = [
		{
			id = "UID:0"
		}
	]
	privileges = [
		{
			id = "ISI_PRIV_SYS_TIME",
			permission = "w"
		}
	]
}
`, roleName, roleDescription+"_modified")
