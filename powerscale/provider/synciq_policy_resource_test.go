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

func TestAccSynciqPolicyResource(t *testing.T) {
	var policyTerraformName = "powerscale_synciq_policy.policy"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink"

					file_matching_pattern = {
						or_criteria = [
							{
								and_criteria = [
									{
										type = "name"
										value = "tfacc"
										operator = "=="
									}
								]
							}
						]
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyTerraformName, "name", "tfaccPolicy"),
					resource.TestCheckResourceAttr(policyTerraformName, "action", "sync"),
					resource.TestCheckResourceAttr(policyTerraformName, "source_root_path", "/ifs"),
					resource.TestCheckResourceAttr(policyTerraformName, "target_path", "/ifs/tfaccSink"),
				),
			},
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyTerraformName, "name", "tfaccPolicy2"),
					resource.TestCheckResourceAttr(policyTerraformName, "action", "sync"),
					resource.TestCheckResourceAttr(policyTerraformName, "source_root_path", "/ifs"),
					resource.TestCheckResourceAttr(policyTerraformName, "target_path", "/ifs/tfaccSink2"),
				),
			},
		},
	})
}
