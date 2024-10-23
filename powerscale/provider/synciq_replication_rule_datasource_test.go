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
	"testing"

	"terraform-provider-powerscale/powerscale/helper"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSynciqRuleDatasourceGetAll tests the synciq rule datasource.
func TestAccSynciqRuleDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing error
			{
				Config: ProviderConfig + `
				data "powerscale_synciq_rule" "test" {
				}
				`,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetAllSyncIQRules).Return(nil, fmt.Errorf("mock network error")).Build()
				},
				ExpectError: regexp.MustCompile("mock network error"),
			},
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_rules" "preq" {
					bandwidth_rules = [
						{
							limit       = 10000
							schedule = {
								begin = "00:00"
								days_of_week = ["friday", "monday"]
								end = "23:59"
							}
						},
					]
				}
				data "powerscale_synciq_rule" "test" {
					depends_on = [powerscale_synciq_rules.preq]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_synciq_rule.test", "rules.#"),
				),
			},
		},
	})
}

// TestAccSynciqRuleDatasourceID tests the synciq rule datasource by ID.
func TestAccSynciqRuleDatasourceID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rules" "preq" {
					bandwidth_rules = [
						{
							limit       = 10000
							schedule = {
								begin = "00:00"
								days_of_week = ["friday", "monday"]
								end = "23:59"
							}
						},
					]
				}
				data "powerscale_synciq_rule" "test" {
					id = powerscale_synciq_rules.preq.bandwidth_rules[0].id
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_synciq_rule.test", "rules.#"),
					resource.TestCheckResourceAttr("data.powerscale_synciq_rule.test", "rules.#", "1"),
				),
			},
			// Read testing error
			{
				Config: ProviderConfig + `
				data "powerscale_synciq_rule" "test" {
					id = ""
				}
				`,
				ExpectError: regexp.MustCompile(`.*string length must be at least 1*`),
			},
			{
				Config: ProviderConfig + `
				data "powerscale_synciq_rule" "test" {
					id = "invalid"
				}
				`,
				ExpectError: regexp.MustCompile(`.*invalid rule ID was specified*`),
			},
		},
	})
}
