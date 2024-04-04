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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var userMappingRulesMocker *mockey.Mocker

func TestAccUserMappingRuleResource(t *testing.T) {
	var mappingRuleResourceName = "powerscale_user_mapping_rules.mapping_rule_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userMappingRuleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mappingRuleResourceName, "parameters.default_unix_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.#", "3"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.operator", "append"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.options.break", "true"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.options.default_user.user", "Guest"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.target_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.target_user.domain", "domain"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.source_user.user", "admin"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.source_user.domain", "domain"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "mapping_users.#", "2"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "mapping_users.0.user.uid", "UID:0"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "mapping_users.0.user.name", "root"),
				),
			},
			// ImportState testing
			{
				ResourceName:      mappingRuleResourceName,
				ImportState:       true,
				ImportStateId:     "System",
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "3", s[0].Attributes["rules.#"])
					assert.Equal(t, "System", s[0].Attributes["zone"])
					assert.Equal(t, "tfaccUserMappungRuleUser", s[0].Attributes["parameters.default_unix_user.user"])
					assert.Equal(t, "admin", s[0].Attributes["rules.0.source_user.user"])
					return nil
				},
				ImportStateVerifyIgnore: []string{"mapping_users", "test_mapping_users"},
			},
			// Update
			{
				Config: ProviderConfig + userMappingRuleResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mappingRuleResourceName, "parameters.default_unix_user.user", " "),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.operator", "replace"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.options.break", "false"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.target_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.target_user.domain", "domain"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.source_user.user", "admin"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.source_user.domain", "domain"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.1.options.break", "true"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.1.options.default_user.user", "Guest"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.1.target_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "mapping_users.#", "1"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "mapping_users.0.user.name", "Guest"),
				),
			},
		},
	})
}

func TestAccUserMappingRuleResourceEmpty(t *testing.T) {
	var mappingRuleResourceName = "powerscale_user_mapping_rules.mapping_rule_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with Empty config and Read testing
			{
				Config: ProviderConfig + userMappingRuleEmptyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(mappingRuleResourceName, "parameters.default_unix_user.user", " "),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.#", "2"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.operator", "replace"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.options.break", "false"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.target_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.target_user.domain", "domain"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.source_user.user", "admin"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.0.source_user.domain", "domain"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.1.options.break", "true"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.1.options.default_user.user", "Guest"),
					resource.TestCheckResourceAttr(mappingRuleResourceName, "rules.1.target_user.user", "tfaccUserMappungRuleUser"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccUserMappingRuleResourceErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error - invalid rule
			{
				Config:      ProviderConfig + userMappingRuleResourceInvalidRule,
				ExpectError: regexp.MustCompile(`.*error updating user mapping rules*.`),
			},
			// Create Error - invalid zone
			{
				Config:      ProviderConfig + userMappingRuleResourceInvalidZone,
				ExpectError: regexp.MustCompile(`.*error getting user mapping rules*.`),
			},
			// Create Error - invalid lookup user
			{
				Config:      ProviderConfig + userMappingRuleResourceInvalidLookupUser,
				ExpectError: regexp.MustCompile(`.*error getting the mapping test user*.`),
			},
			{
				Config: ProviderConfig + userMappingRuleResourceConfig,
			},
			// Update Error - invalid rule
			{
				Config:      ProviderConfig + userMappingRuleResourceInvalidRule,
				ExpectError: regexp.MustCompile(`.*error updating user mapping rules*.`),
			},
			{
				Config: ProviderConfig + userMappingRuleResourceUpdateConfig,
			},
		},
	})
}

func TestAccUserMappingRuleResourceMockErr(t *testing.T) {
	diags := diag.Diagnostics{}
	diags.AddError("mock err", "mock err")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock((*powerscale.AuthApiService).UpdateAuthv1MappingUsersRulesExecute).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Create and Read Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock((*powerscale.AuthApiService).GetAuthv1MappingUsersRulesExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Create and parse Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateUserMappingRulesState).Return(diags).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
			// Create and Lookup User Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock((*powerscale.AuthApiService).GetAuthv1MappingUsersLookupExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceConfig,
				ExpectError: regexp.MustCompile(`.*error getting user mapping test*.`),
			},
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userMappingRuleResourceConfig,
			},
			// Read Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock((*powerscale.AuthApiService).GetAuthv1MappingUsersRulesExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Read and parse Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateUserMappingRulesState).Return(diags).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
			// Read and lookupUser Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateLookupMappingUsersState).Return(diags).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
			// Update Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock((*powerscale.AuthApiService).UpdateAuthv1MappingUsersRulesExecute).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Read Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.GetUserMappingRulesByZone).When(func(ctx context.Context, client *client.Client, zone string) bool {
						return userMappingRulesMocker.Times() > 1
					}).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Parse Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateUserMappingRulesState).When(func(ctx context.Context, rulesState *models.UserMappingRulesResourceModel, rulesResponse *powerscale.V1MappingUsersRulesRules) bool {
						return userMappingRulesMocker.Times() > 1
					}).Return(diags).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
			// Update and Lookup User Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateLookupMappingUsersState).When(func(ctx context.Context, client *client.Client, plan *models.UserMappingRulesResourceModel) bool {
						return userMappingRulesMocker.Times() > 1
					}).Return(diags).Build()
				},
				Config:      ProviderConfig + userMappingRuleResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
				},
				Config: ProviderConfig + userMappingRuleResourceUpdateConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.GetUserMappingRulesByZone).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + userMappingRuleResourceUpdateConfig,
				ResourceName:      "powerscale_user_mapping_rules.mapping_rule_test",
				ImportStateId:     "System",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
			// Import and parse Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateUserMappingRulesState).Return(diags).Build()
				},
				Config:            ProviderConfig + userMappingRuleResourceConfig,
				ResourceName:      "powerscale_user_mapping_rules.mapping_rule_test",
				ImportStateId:     "System",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock err*.`),
				ImportStateVerify: true,
			},
			// Import and lookup user Error testing
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.UnPatch()
					}
					userMappingRulesMocker = mockey.Mock(helper.UpdateLookupMappingUsersState).Return(diags).Build()
				},
				Config:            ProviderConfig + userMappingRuleResourceConfig,
				ResourceName:      "powerscale_user_mapping_rules.mapping_rule_test",
				ImportStateId:     "System",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock err*.`),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUserMappingRuleReleaseMockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if userMappingRulesMocker != nil {
						userMappingRulesMocker.Release()
					}
				},
				Config: ProviderConfig + userMappingRuleResourceConfig,
			},
		},
	})
}

var userMappingRuleEmptyResourceConfig = `
resource "powerscale_user_mapping_rules" "mapping_rule_test" {
}
`

var userMappingRuleResourceConfig = `
resource "powerscale_user_mapping_rules" "mapping_rule_test" {	
	zone = "System"
	parameters = {
	  default_unix_user = {
		domain = "domain",
		user   = "tfaccUserMappungRuleUser"
	  }
	}

	rules = [
	  {
		operator = "append",
		options = {
		  break = true,
		  default_user = {
			domain = "domain",
			user   = "Guest"
		  },
		  group  = true,
		  groups = true,
		  user   = true
		},
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  domain = "domain",
		  user   = "admin"
		}
	  },
	  {
		operator = "union",
		target_user = {
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  user   = "admin"
		}
	  },
	  {
		operator = "trim",
		options = {
		  break = true,
		},
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		}
	  }
	]
  
	test_mapping_users = [
	  {
		uid = 0
	  },
	  {
		name = "admin"
	  }
	]
}
`

var userMappingRuleResourceUpdateConfig = `
resource "powerscale_user_mapping_rules" "mapping_rule_test" {	
	zone = "System"
	parameters = {
	  default_unix_user = {
		user   = " "
	  }
	}

	rules = [
	  {
		operator = "replace",
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  domain = "domain",
		  user   = "admin"
		}
	  },
	  {
		operator = "insert",
		options = {
			break = true,
			default_user = {
			  domain = "domain",
			  user   = "Guest"
			},
			group  = true,
			groups = true,
			user   = true
		},
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  domain = "domain",
		  user   = "admin"
		}
	  },
	]
  
	test_mapping_users = [
	  {
		name = "Guest"
	  }
	]
}
`

var userMappingRuleResourceInvalidRule = `
resource "powerscale_user_mapping_rules" "mapping_rule_test" {
	zone = "System"
	parameters = {
	  default_unix_user = {
		user   = " "
	  }
	}

	rules = [
	  {
		operator = "replace",
		options = {
			break = true,
			default_user = {
			  domain = "domain",
			  user   = "Guest"
			},
			group  = true,
			groups = true,
			user   = true
		},
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  domain = "domain",
		  user   = "admin"
		}
	  }
	]
}
`

var userMappingRuleResourceInvalidLookupUser = `
resource "powerscale_user_mapping_rules" "mapping_rule_test" {	
	zone = "System"
	parameters = {
	  default_unix_user = {
		user   = " "
	  }
	}

	rules = [
	  {
		operator = "replace",
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  domain = "domain",
		  user   = "admin"
		}
	  },
	  {
		operator = "insert",
		options = {
			break = true,
			default_user = {
			  domain = "domain",
			  user   = "Guest"
			},
			group  = true,
			groups = true,
			user   = true
		},
		target_user = {
		  domain = "domain",
		  user   = "tfaccUserMappungRuleUser"
		},
		source_user = {
		  domain = "domain",
		  user   = "admin"
		}
	  },
	]
  
	test_mapping_users = [
	  {
		name = "invalidLookupUser"
	  }
	]
}
`

var userMappingRuleResourceInvalidZone = `
resource "powerscale_user_mapping_rules" "mapping_rule_test" {	
	zone = "invalidMappingUserZone"
}
`
