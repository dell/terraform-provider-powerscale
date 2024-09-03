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
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSynciqRuleResource(t *testing.T) {
	var resourceName = "powerscale_synciq_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// invalid create
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = -200
					description = "tfacc created"
				}
				`,
				ExpectError: regexp.MustCompile(`.*Error creating syncIQ rule*.`),
			},
			// Create
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = 10000
					description = "tfacc created"
				}
				`,
			},
			// check that import is creating correct state
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			// import with invalid ID
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: "invalid",
				ExpectError:   regexp.MustCompile(`.*Error reading syncIQ rule*.`),
			},
			// Update
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = 20000
					description = "tfacc updated"
					enabled = true
					schedule = {
						begin = "01:00",
						end = "22:59",
					}
				}
				`,
			},
			// invalid update with wrong schedule
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = 20000
					description = "tfacc updated"
					enabled = true
					schedule = {
						begin = "invalid"
						end = "22:59",
					}
				}
				`,
				ExpectError: regexp.MustCompile(`.*Error updating syncIQ rule*.`),
			},
			// Add days of week
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = 20000
					description = "tfacc updated"
					enabled = true
					schedule = {
						begin = "01:00",
						end = "22:59",
						days_of_week = ["monday", "wednesday", "thursday"]
					}
				}
				`,
			},
			// mock delete error
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = 20000
					description = "tfacc updated"
					enabled = true
					schedule = {
						begin = "01:00",
						end = "22:59",
						days_of_week = ["monday", "wednesday", "thursday"]
					}
				}
				`,
				Destroy: true,
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.DeleteSyncIQRule).Return(fmt.Errorf("mock delete error")).Build()
				},
				ExpectError: regexp.MustCompile(`.*Error deleting syncIQ rule*.`),
			},
			// remove days of week
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_rule" "test" {
					type = "bandwidth"
					limit = 20000
					description = "tfacc updated"
					enabled = false
					schedule = {
						begin = "01:00",
						end = "22:59",
						days_of_week = ["monday"]
					}
				}
				`,
			},
		},
	})
}
