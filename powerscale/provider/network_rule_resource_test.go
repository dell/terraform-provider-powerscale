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
	"context"
	"fmt"
	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

var ruleMocker *mockey.Mocker

func TestAccRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RuleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "name", "tfacc_rule"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "iface", "ext-2"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "node_type", "any"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_network_rule.rule",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "groupnet0.subnet0.pool0.tfacc_rule", states[0].Attributes["id"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + RuleResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "name", "tfacc_rule_rename"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "iface", "ext-2"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "node_type", "any"),
				),
			},
		},
	})
}

func TestAccRuleResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RuleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "name", "tfacc_rule"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "subnet", "subnet0"),
				),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_network_rule.rule",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNetworkRule).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing wrong ID
			{
				ResourceName:  "powerscale_network_rule.rule",
				ImportState:   true,
				ImportStateId: "ruleId",
				ExpectError:   regexp.MustCompile("Unexpected Import Identifier"),
			},
		},
	})
}

func TestAccRuleResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + RuleResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CreateNetworkRule).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccRuleResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RuleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "name", "tfacc_rule"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "iface", "ext-2"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "node_type", "any"),
				),
			},
			// Update get error
			{
				Config: ProviderConfig + RuleResourceUpdateConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateNetworkRule).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + RuleResourceUpdateConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetNetworkRule).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + RuleResourceUpdateConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccRuleResourceErrorDelete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + RuleResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "name", "tfacc_rule"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "iface", "ext-2"),
					resource.TestCheckResourceAttr("powerscale_network_rule.rule", "node_type", "any"),
				),
			},
			{
				PreConfig: func() {
					if ruleMocker != nil {
						ruleMocker.UnPatch()
					}
					ruleMocker = mockey.Mock(helper.DeleteNetworkRule).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, ruleName string, groupnet string, subnet string, pool string) bool {
							return ruleMocker.Times() == 1
						})
				},
				Config:      ProviderConfig + RuleResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccRuleResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RuleResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var RuleResourceConfig = `
resource "powerscale_network_rule" "rule" {
  name = "tfacc_rule"
  groupnet = "groupnet0"
  subnet = "subnet0"
  pool = "pool0"
  description = "tfacc_rule"
  iface = "ext-2"
  node_type = "any"
}
`

var RuleResourceUpdateConfig = `
resource "powerscale_network_rule" "rule" {
  name = "tfacc_rule_rename"
  groupnet = "groupnet0"
  subnet = "subnet0"
  pool = "pool0"
  description = "tfacc_rule_rename"
  iface = "ext-2"
  node_type = "any"
}
`
