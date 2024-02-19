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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"testing"
)

func TestAccRuleDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + RuleDatasourceGetAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_network_rule.rule", "network_rules.#"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.id", "groupnet0.subnet0.pool0.rule0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.subnet", "subnet0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.pool", "pool0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.name", "rule0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.iface", "ext-1"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.node_type", "any"),
				),
			},
		},
	})
}

func TestAccRuleDatasourceGetFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + RuleDatasourceGetFilterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_network_rule.rule", "network_rules.#"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.id", "groupnet0.subnet0.pool0.rule0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.subnet", "subnet0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.pool", "pool0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.name", "rule0"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.iface", "ext-1"),
					resource.TestCheckResourceAttr("data.powerscale_network_rule.rule", "network_rules.0.node_type", "any"),
				),
			},
		},
	})
}

func TestAccRuleDatasourceGetFilterError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config:      ProviderConfig + RuleDatasourceGetFilterErrorConfig,
				ExpectError: regexp.MustCompile("Error getting the list"),
			},
		},
	})
}

var RuleDatasourceGetAllConfig = `
data "powerscale_network_rule" "rule" {
}
`

var RuleDatasourceGetFilterConfig = `
data "powerscale_network_rule" "rule" {
  filter{
	names=["rule0"]
	groupnet="groupnet0"
	subnet="subnet0"
	pool="pool0"
  }
}
`

var RuleDatasourceGetFilterErrorConfig = `
data "powerscale_network_rule" "rule" {
  filter{
	groupnet="groupnet-non-existent"
	subnet="subnet0"
	pool="pool0"
  }
}
`
