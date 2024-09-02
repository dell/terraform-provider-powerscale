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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSynciqRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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
			// remove days of week
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
						days_of_week = ["monday"]
					}
				}
				`,
			},
		},
	})
}
