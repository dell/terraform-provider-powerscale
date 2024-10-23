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

// TestAccSynciqRulesResource tests the syncIQ rules resource.
func TestAccSynciqRulesResource(t *testing.T) {
	var resourceName = "powerscale_synciq_rules.test"
	withNoRules := `
				resource "powerscale_synciq_rules" "test" {
					bandwidth_rules = []
				}
				`
	withTwoRules := `
				resource "powerscale_synciq_rules" "test" {
					bandwidth_rules = [
						{
							limit       = 10000
							schedule = {
								begin = "00:00"
								days_of_week = ["friday", "monday"]
								end = "23:59"
							}
						},
						{
							description = "tfAcc2 description"
							limit       = 2000
							schedule = {
								begin = "01:00"
								days_of_week = ["monday", "saturday", "sunday"]
								end = "23:59"
							}
						},
					]
				}
				`
	with3RuleMinimal := `
				resource "powerscale_synciq_rules" "test" {
					bandwidth_rules = [
						{
							limit       = 10000
							schedule = {
								begin = "00:00"
								days_of_week = ["friday", "monday"]
								end = "23:59"
							}
						},    
						{
							limit       = 5000
							description = "tfAcc 1/1 description"
						}, 
						{
							description = "tfAcc2 description"
							limit       = 2000
							schedule = {
								begin = "01:00"
								days_of_week = ["monday", "saturday", "sunday"]
								end = "23:59"
							}
						},
					]
				}
				`
	with3RuleFull := `
				resource "powerscale_synciq_rules" "test" {
					bandwidth_rules = [
						{
							limit       = 10000
							schedule = {
								begin = "00:00"
								days_of_week = ["friday", "monday"]
								end = "23:59"
							}
						},    
						{
							limit       = 5000
							schedule = {
								begin = "00:40"
								days_of_week = ["friday", "monday", "tuesday"]
								end = "23:29"
							}
						}, 
						{
							description = "tfAcc2 description"
							limit       = 2000
							schedule = {
								begin = "01:00"
								days_of_week = ["monday", "saturday", "sunday"]
								end = "23:59"
							}
						},
					]
				}
				`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// error reading before create
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetAllSyncIQRules).Return(nil, fmt.Errorf("mock network error")).Build()
				},
				Config:      ProviderConfig + withNoRules,
				ExpectError: regexp.MustCompile(`.*mock network error*.`),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + withTwoRules,
			},
			// error during import
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetAllSyncIQRules).Return(nil, fmt.Errorf("mock network error")).Build()
				},
				ResourceName: resourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock network error*.`),
			},
			// check that import is creating correct state
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				ResourceName: resourceName,
				ImportState:  true,
			},
			// error creating new rule
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CreateSyncIQRule).Return("", fmt.Errorf("mock create error")).Build()
				},
				Config:      ProviderConfig + with3RuleMinimal,
				ExpectError: regexp.MustCompile(`.*mock create error*.`),
			},
			// Create new rule
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + with3RuleMinimal,
			},
			// error updating new rule with schedule
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.UpdateSyncIQRule).Return(fmt.Errorf("mock update error")).Build()
				},
				Config:      ProviderConfig + with3RuleFull,
				ExpectError: regexp.MustCompile(`.*mock update error*.`),
			},
			// Update new rule with schedule
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + with3RuleFull,
			},
			// error deleting new rule
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.DeleteSyncIQRule).Return(fmt.Errorf("mock delete error")).Build()
				},
				Config:      ProviderConfig + withTwoRules,
				ExpectError: regexp.MustCompile(`.*mock delete error*.`),
			},
			// delete new rule
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + withTwoRules,
			},
			// delete all rules
			{
				Config: ProviderConfig + withNoRules,
			},
		},
	})
}
