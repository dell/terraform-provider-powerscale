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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var userMocker *mockey.Mocker
var userCreateMocker *mockey.Mocker

func TestAccUserResourceCreate(t *testing.T) {
	var userResourceName = "powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + userResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(userResourceName, "uid", "20000"),
					resource.TestCheckResourceAttr(userResourceName, "email", "test@dell.com"),
					resource.TestCheckResourceAttr(userResourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "roles.0", "SystemAdmin"),
					resource.TestCheckResourceAttr(userResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(userResourceName, "primary_group", "Administrators"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + userUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(userResourceName, "uid", "20001"),
					resource.TestCheckResourceAttr(userResourceName, "email", "newTest@dell.com"),
					resource.TestCheckResourceAttr(userResourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(userResourceName, "roles.0", "tfaccUserRole"),
					resource.TestCheckResourceAttr(userResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(userResourceName, "primary_group", "Administrators"),
				),
			},
			// Update role Error testing
			{
				Config:      ProviderConfig + userInvalidRoleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error assign User to Role*.`),
			},
		},
	})
}

func TestAccUserResourceResetPassword(t *testing.T) {
	var userResourceName = "powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + userBasicResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(userResourceName, "prompt_password_change", "false"),
				),
			},
			// Update and Read testing - enable prompt_password_change
			{
				Config: ProviderConfig + userResourceConfigPasswordPrompt,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(userResourceName, "prompt_password_change", "true"),
				),
			},
			// Update Error testing - disable prompt_password_change separately
			{
				Config:      ProviderConfig + userResourceConfigPasswordPromptError,
				ExpectError: regexp.MustCompile(`.*Error updating the User*.`),
			},
			// Update and Read testing - update password and disable prompt_password_change
			{
				Config: ProviderConfig + userResourceConfigPasswordReset,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(userResourceName, "prompt_password_change", "false"),
				),
			},
		},
	})
}

func TestAccUserResourceCreateErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error testing
			{
				Config:      ProviderConfig + userErrCreateResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating the User*.`),
			},
		},
	})
}

func TestAccUserResourceAddRoleErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      ProviderConfig + userInvalidRoleResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error assign User to Role*.`),
			},
		},
	})
}

func TestAccUserResourceImport(t *testing.T) {
	var userResourceName = "powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + userResourceConfig,
			},
			{
				ResourceName:      userResourceName,
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "tfaccUserCreation", s[0].Attributes["name"])
					assert.Equal(t, "test@dell.com", s[0].Attributes["email"])
					assert.Equal(t, "1", s[0].Attributes["roles.#"])
					assert.Equal(t, "SystemAdmin", s[0].Attributes["roles.0"])
					assert.Equal(t, "Administrators", s[0].Attributes["primary_group"])
					return nil
				},
			},
		},
	})
}

func TestAccUserResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Error testing
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
					userCreateMocker = mockey.Mock(helper.CreateUser).Return(nil).Build()
					userMocker = mockey.Mock(helper.GetUserWithZone).Return(nil, fmt.Errorf("user read mock error")).Build()
				},
				Config:      ProviderConfig + userBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*user read mock error*.`),
			},
		},
	})
}

func TestAccUserRolesResourceImportRolesErr(t *testing.T) {
	var userResourceName = "powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create User and Read Roles Error testing
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
					userMocker = mockey.Mock(helper.GetAllRolesWithZone).Return(nil, fmt.Errorf("roles read mock error")).Build()
				},
				ResourceName:      userResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*roles read mock error*.`),
			},
		},
	})
}

func TestAccUserRolesResourceImportGetErr(t *testing.T) {
	var userResourceName = "powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create User and Read Roles Error testing
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
					userMocker = mockey.Mock(helper.GetUserWithZone).Return(nil, fmt.Errorf("user read mock error")).Build()
				},
				ResourceName:      userResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*user read mock error*.`),
			},
		},
	})
}

func TestAccUserRolesResourceDeleteMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.UnPatch()
					}
					if userCreateMocker != nil {
						userCreateMocker.UnPatch()
					}
					userMocker = mockey.Mock((*powerscale.AuthApiService).DeleteAuthv1AuthUserExecute).When(func(r powerscale.ApiDeleteAuthv1AuthUserRequest) bool {
						return userMocker.Times() == 1
					}).Return(nil, fmt.Errorf("user mock error")).Build()
				},
				Config:      ProviderConfig + userBasicResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("user mock error"),
			},
		},
	})
}

func TestAccUserReleaseMockResource(t *testing.T) {
	var userResourceName = "powerscale_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create User and Read Roles Error testing
			{
				PreConfig: func() {
					if userMocker != nil {
						userMocker.Release()
					}
					if userCreateMocker != nil {
						userCreateMocker.Release()
					}
				},
				Config: ProviderConfig + userBasicResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userResourceName, "name", "tfaccUserCreation"),
				)},
		},
	})
}

var userBasicResourceConfig = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
	roles = ["tfaccUserRole"]
  }
`

var userResourceConfig = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"

	uid = 20000
	email = "test@dell.com"
	primary_group = "Administrators"
	roles = ["SystemAdmin"]
  }
`

var userUpdateResourceConfig = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
  
	uid = 20001
	query_force = true
	enabled = true
	email = "newTest@dell.com"
	primary_group = "Administrators"
	roles = ["tfaccUserRole"]
  }
`

var userErrCreateResourceConfig = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"

	uid = 20000
	email = "test@dell.com"
	primary_group = "InvalidGroup"
	roles = ["SystemAdmin"]
  }
`

var userInvalidRoleResourceConfig = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
  
	uid = 20001
	query_force = true
	enabled = true
	email = "newTest@dell.com"
	primary_group = "Administrators"
	roles = ["invalidRole"]
  }
`

var userResourceConfigPasswordPrompt = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
	email = "PasswordPrompt@dell.com"
	prompt_password_change = true
	roles = ["tfaccUserRole"]
  }
`

var userResourceConfigPasswordPromptError = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
	email = "PasswordPromptError@dell.com"
	prompt_password_change = false
	roles = ["tfaccUserRole"]
  }
`

var userResourceConfigPasswordReset = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
	email = "PasswordReset@dell.com"
	password = "testPasswordReset"
	prompt_password_change = false
	roles = ["tfaccUserRole"]
  }
`
