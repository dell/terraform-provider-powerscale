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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserMappingRulesDataSource(t *testing.T) {
	var userMappingRulesDataSourceName = "data.powerscale_user_mapping_rules.mapping_rule_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + userMappingRuleResourceConfig + userMappingRulesDataSourceConfigAll,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules_parameters.default_unix_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.#", "3"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.operator", "append"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.options.break", "true"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.options.default_user.user", "Guest"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.target_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.target_user.domain", "domain"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.source_user.user", "admin"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.source_user.domain", "domain"),
				),
			},
			// read with filter testing
			{
				Config: ProviderConfig + userMappingRuleResourceConfig + userMappingRulesDataSourceConfigFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules_parameters.default_unix_user.user", "tfaccUserMappungRuleUser"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.#", "1"),
					resource.TestCheckResourceAttr(userMappingRulesDataSourceName, "user_mapping_rules.0.operator", "append"),
				),
			},
		},
	})
}

func TestAccUserMappingRulesDataSourceInvalidConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid names read testing
			{
				Config:      ProviderConfig + userMappingRulesDataSourceConfigInvalidName,
				ExpectError: regexp.MustCompile(`.*error one or more of the filtered user names is invalid.*.`),
			},
			// filter with invalid operators read testing
			{
				Config:      ProviderConfig + userMappingRulesDataSourceConfigInvalidOperators,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value Match*.`),
			},
			// filter with invalid zone read testing
			{
				Config:      ProviderConfig + userMappingRulesDataSourceConfigInvalidZone,
				ExpectError: regexp.MustCompile(`.*error getting PowerScale User Mapping Rules.*.`),
			},
		},
	})
}

var userMappingRulesDataSourceConfigAll = `
data "powerscale_user_mapping_rules" "mapping_rule_test" {
	depends_on = [powerscale_user_mapping_rules.mapping_rule_test]
}
`

var userMappingRulesDataSourceConfigFilter = `
data "powerscale_user_mapping_rules" "mapping_rule_test" {
	depends_on = [powerscale_user_mapping_rules.mapping_rule_test]

	filter {
		names = ["admin"]
		operators = ["append"]
		zone = "System"
	}
}
`

var userMappingRulesDataSourceConfigInvalidName = `
data "powerscale_user_mapping_rules" "mapping_rule_test" {
	filter {
		names = ["invalidMappingUserName"]
	}
}
`

var userMappingRulesDataSourceConfigInvalidOperators = `
data "powerscale_user_mapping_rules" "mapping_rule_test" {
	filter {
		operators = ["invalidOperators"]
	}
}
`

var userMappingRulesDataSourceConfigInvalidZone = `
data "powerscale_user_mapping_rules" "mapping_rule_test" {
	filter {
		zone = "invalidMappingUserZone"
	}
}
`
