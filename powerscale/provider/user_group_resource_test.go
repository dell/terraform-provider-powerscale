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
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var userGroupMocker *mockey.Mocker
var userGroupCreateMocker *mockey.Mocker

func TestAccUserGroupResourceCreate(t *testing.T) {
	var userGroupResourceName = "powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + userGroupResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userGroupResourceName, "name", "tfaccUserGroupCreation"),
					resource.TestCheckResourceAttr(userGroupResourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(userGroupResourceName, "roles.0", "SystemAdmin"),
					resource.TestCheckResourceAttr(userGroupResourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(userGroupResourceName, "users.0", "tfaccMemberUser"),
					resource.TestCheckResourceAttr(userGroupResourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(userGroupResourceName, "groups.0", "wheel"),
					resource.TestCheckResourceAttr(userGroupResourceName, "well_knowns.#", "1"),
					resource.TestCheckResourceAttr(userGroupResourceName, "well_knowns.0", "Everyone"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + userGroupUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userGroupResourceName, "name", "tfaccUserGroupCreation"),
					resource.TestCheckResourceAttr(userGroupResourceName, "gid", "20001"),
					resource.TestCheckResourceAttr(userGroupResourceName, "roles.#", "2"),
					resource.TestCheckResourceAttr(userGroupResourceName, "users.#", "2"),
					resource.TestCheckResourceAttr(userGroupResourceName, "groups.#", "1"),
					resource.TestCheckResourceAttr(userGroupResourceName, "groups.0", "admin"),
					resource.TestCheckResourceAttr(userGroupResourceName, "well_knowns.#", "1"),
					resource.TestCheckResourceAttr(userGroupResourceName, "well_knowns.0", "NT AUTHORITY\\NETWORK"),
				),
			},
			// Update member Error testing
			{
				Config:      ProviderConfig + userGroupInvalidMemberResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error add User*.`),
			},
			// Update role Error testing
			{
				Config:      ProviderConfig + userGroupInvalidRoleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error assign User Group to Role*.`),
			},
		},
	})
}

func TestAccUserGroupResourceCreateErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error testing
			{
				Config:      ProviderConfig + userGroupErrCreateResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not add member to user group*.`),
			},
		},
	})
}

func TestAccUserGroupResourceAddRoleErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      ProviderConfig + userGroupInvalidRoleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error assign User Group to Role*.`),
			},
		},
	})
}

func TestAccUserGroupResourceImport(t *testing.T) {
	var userGroupResourceName = "powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + userGroupResourceConfig,
			},
			{
				ResourceName:      userGroupResourceName,
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "tfaccUserGroupCreation", s[0].Attributes["name"])
					assert.Equal(t, "1", s[0].Attributes["roles.#"])
					assert.Equal(t, "1", s[0].Attributes["users.#"])
					assert.Equal(t, "1", s[0].Attributes["groups.#"])
					assert.Equal(t, "wheel", s[0].Attributes["groups.0"])
					assert.Equal(t, "1", s[0].Attributes["well_knowns.#"])
					assert.Equal(t, "Everyone", s[0].Attributes["well_knowns.0"])
					return nil
				},
			},
		},
	})
}

func TestAccUserGroupResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Error testing
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupCreateMocker = mockey.Mock(helper.CreateUserGroup).Return(nil).Build()
					userGroupMocker = mockey.Mock(helper.GetUserGroupWithZone).Return(nil, fmt.Errorf("user group read mock error")).Build()
				},
				Config:      ProviderConfig + userGroupBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*user group read mock error*.`),
			},
		},
	})
}

func TestAccUserGroupRolesResourceImportRolesErr(t *testing.T) {
	var userGroupResourceName = "powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create User Group and Read Roles Error testing
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userGroupBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupMocker = mockey.Mock(helper.GetAllRolesWithZone).Return(nil, fmt.Errorf("roles read mock error")).Build()
					userGroupCreateMocker = mockey.Mock(helper.GetAllGroupMembersWithZone).Return(nil, fmt.Errorf("members read mock error")).Build()
				},
				ResourceName:      userGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*roles read mock error*.`),
			},
		},
	})
}

func TestAccUserGroupRolesResourceImportGetErr(t *testing.T) {
	var userGroupResourceName = "powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create User Group and Read Roles Error testing
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userGroupBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupMocker = mockey.Mock(helper.GetUserGroupWithZone).Return(nil, fmt.Errorf("user group read mock error")).Build()
				},
				ResourceName:      userGroupResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*user group read mock error*.`),
			},
		},
	})
}

func TestAccUserGroupResourceDeleteMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userGroupBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupMocker = mockey.Mock(helper.DeleteUserGroupWithZone).Return(fmt.Errorf("user group delete mock error")).Build().
						When(func(ctx context.Context, client *client.Client, groupName, zone string) bool {
							return userGroupMocker.Times() == 1
						})
				},
				Config:      ProviderConfig + userGroupBasicResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("user group delete mock error"),
			},
		},
	})
}

func TestAccUserGroupResourceHelperMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupMocker = mockey.Mock((*powerscale.AuthApiService).GetAuthv1AuthGroupExecute).Return(nil, nil, fmt.Errorf("user group read mock error")).Build()
					userGroupCreateMocker = mockey.Mock(helper.CreateUserGroup).Return(nil).Build()
				},
				Config:      ProviderConfig + userGroupBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*user group read mock error*.`),
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupCreateMocker = mockey.Mock(helper.CreateUserGroup).Return(nil).Build()
					userGroupMocker = mockey.Mock((*powerscale.AuthApiService).GetAuthv1AuthGroupExecute).Return(&powerscale.V1AuthGroupsExtended{}, nil, nil).Build()
				},
				Config:      ProviderConfig + userGroupBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*got empty user group*.`),
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupMocker = mockey.Mock((*powerscale.AuthApiService).CreateAuthv1AuthGroupExecute).Return(nil, nil, fmt.Errorf("create mock error")).Build()
				},
				Config:      ProviderConfig + userGroupBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating the User Group*.`),
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupMocker = mockey.Mock((*powerscale.AuthRolesApiService).CreateAuthRolesv7RoleMemberExecute).Return(nil, nil, fmt.Errorf("role mock error")).Build()
					userGroupCreateMocker = mockey.Mock((*powerscale.AuthGroupsApiService).CreateAuthGroupsv1GroupMemberExecute).Return(nil, nil, fmt.Errorf("member mock error")).Build()
				},
				Config:      ProviderConfig + userGroupUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*role mock error*.`),
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userGroupResourceConfig,
			},
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.UnPatch()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.UnPatch()
					}
					userGroupCreateMocker = mockey.Mock((*powerscale.AuthApiService).DeleteAuthv1GroupsGroupMemberExecute).Return(nil, fmt.Errorf("member mock error")).Build()
				},
				Config:      ProviderConfig + userGroupBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*member mock error*.`),
			},
		},
	})
}

func TestAccUserGroupReleaseMockResource(t *testing.T) {
	var userGroupResourceName = "powerscale_user_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create User Group and Read Roles Error testing
			{
				PreConfig: func() {
					if userGroupMocker != nil {
						userGroupMocker.Release()
					}
					if userGroupCreateMocker != nil {
						userGroupCreateMocker.Release()
					}
				},
				Config: ProviderConfig + userGroupBasicResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userGroupResourceName, "name", "tfaccUserGroupCreation"),
				)},
		},
	})
}

var userGroupBasicResourceConfig = `
resource "powerscale_user_group" "test" {
    name = "tfaccUserGroupCreation"
	roles = ["SystemAdmin"]
	groups    = ["guest"]
  }
`

var userGroupResourceConfig = `
resource "powerscale_user" "testDepMember" {
	name = "tfaccMemberUser"
}

resource "powerscale_user" "testDepMember2" {
	name = "tfaccMemberUser2"
}

resource "powerscale_user_group" "test" {
    name = "tfaccUserGroupCreation"

    # query_force = true
    # query_zone = "testZone"
    # query_provider = "testProvider"

    # sid = "SID:testSID"
    roles = ["SystemAdmin"]
    users = ["tfaccMemberUser"]
    groups    = ["wheel"]
    well_knowns    = ["Everyone"]

	depends_on = [
		powerscale_user.testDepMember,
		powerscale_user.testDepMember2
	  ]
  }
`

var userGroupUpdateResourceConfig = `
resource "powerscale_user" "testDepMember" {
	name = "tfaccMemberUser"
}

resource "powerscale_user" "testDepMember2" {
	name = "tfaccMemberUser2"
}

resource "powerscale_user_group" "test" {
    name = "tfaccUserGroupCreation"

    query_force = true

    gid = 20001
    roles = ["SystemAdmin","tfaccUserRole"]
    users = ["tfaccMemberUser","tfaccMemberUser2"]
    groups    = ["admin"]
    well_knowns    = ["NT AUTHORITY\\NETWORK"]

	depends_on = [
		powerscale_user.testDepMember,
		powerscale_user.testDepMember2
	  ]
  }
`

var userGroupInvalidMemberResourceConfig = `
resource "powerscale_user_group" "test" {
    name = "tfaccUserGroupCreation"

    query_force = true
    gid = 20001
    users = ["invalidUserMember"]
  }
`

var userGroupErrCreateResourceConfig = `
resource "powerscale_user_group" "test" {
    name = "tfaccUserGroupCreation"

    users = ["invalidUserMember"]
  }
`

var userGroupInvalidRoleResourceConfig = `
resource "powerscale_user_group" "test" {
    name = "tfaccUserGroupCreation"

    query_force = true
    gid = 20001
    roles = ["invalidRole"]
  }
`
