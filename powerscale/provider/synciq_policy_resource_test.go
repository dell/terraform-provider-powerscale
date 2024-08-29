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

func TestAccSynciqPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// mock create plan reading error test
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock create plan read error")).Build()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink"
				}`,
				ExpectError: regexp.MustCompile("mock create plan read error"),
			},
			// create test positive
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink"
				}
				`,
				ExpectNonEmptyPlan: true,
			},
			// import negative
			{
				ResourceName:  "powerscale_synciq_policy.policy",
				ImportState:   true,
				ImportStateId: "non-existing",
				ExpectError:   regexp.MustCompile("not found"),
			},
			// import positive
			{
				ResourceName:  "powerscale_synciq_policy.policy",
				ImportState:   true,
				ImportStateId: "tfaccPolicy",
			},
			// mock update plan reading error test
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock update plan reading error")).Build()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink"
				}`,
				ExpectError: regexp.MustCompile("mock update plan reading error"),
			},
			// update positive
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"

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
				ExpectNonEmptyPlan: true,
			},
			// create negative - Existing name
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"

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

				resource "powerscale_synciq_policy" "policy2" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"
				}
				`,
				ExpectError: regexp.MustCompile("Error creating syncIQ Policy"),
			},
			// update negative - Invalid source root path
			// root path needs to start with /ifs
			{
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/invalid"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"

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
				ExpectError: regexp.MustCompile(".*Could not update syncIQ Policy.*"),
			},
			// mock destroy error
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.DeleteSyncIQPolicy).Return(fmt.Errorf("mock delete error")).Build()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"
				}
				`,
				Destroy:     true,
				ExpectError: regexp.MustCompile("mock delete error"),
			},
			// mock refresh error
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSyncIQPolicyByID).Return(nil, fmt.Errorf("mock read error")).Build()
				},
				RefreshState: true,
				ExpectError:  regexp.MustCompile("mock read error"),
			},
			// update positive - Remove all file matching patterns
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + `
				resource "powerscale_synciq_policy" "policy" {
					name = "tfaccPolicy2"
					action = "sync"
					source_root_path = "/ifs"
					target_host = "10.10.10.10"
					target_path = "/ifs/tfaccSink2"
				}
				`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
