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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSupportAssistResource(t *testing.T) {
	supportAssistResourceName := "powerscale_support_assist.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + supportAssistResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(supportAssistResourceName, "supportassist_enabled", "true"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "enable_download", "false"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "automatic_case_creation", "true"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "enable_remote_support", "true"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "accepted_terms", "true"),
				),
			},
			{
				Config: ProviderConfig + supportAssistResourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(supportAssistResourceName, "supportassist_enabled", "false"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "enable_download", "false"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "automatic_case_creation", "false"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "enable_remote_support", "true"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "accepted_terms", "false"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "telemetry.offline_collection_period", "7800"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "telemetry.telemetry_enabled", "false"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "telemetry.telemetry_persist", "true"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "telemetry.telemetry_threads", "5"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "contact.primary.first_name", "terraform_first"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "contact.secondary.email", "xyz@gmail.com"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "connections.mode", "gateway"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "connections.network_pools.0", "subnet0:pool0"),
					resource.TestCheckResourceAttr(supportAssistResourceName, "connections.gateway_endpoints.0.host", "1.2.3.4"),
				),
			},
			{
				Config:       ProviderConfig + supportAssistResourceConfigUpdate,
				ResourceName: supportAssistResourceName,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(supportAssistResourceName, "id")
					resource.TestCheckResourceAttrSet(supportAssistResourceName, "supportassist_enabled")
					resource.TestCheckResourceAttrSet(supportAssistResourceName, "automatic_case_creation")
					resource.TestCheckResourceAttrSet(supportAssistResourceName, "enable_remote_support")
					resource.TestCheckResourceAttrSet(supportAssistResourceName, "accepted_terms")
					return nil
				},
			},
		},
	})
}

func TestAccSupportAssistResourceMockError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSupportAssistSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + supportAssistResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + supportAssistResourceConfigUpdate,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSupportAssist).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + supportAssistResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSupportAssistSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + supportAssistResourceConfigUpdate1,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var supportAssistResourceConfig = `
resource "powerscale_support_assist" "test" {
	supportassist_enabled   = true
  	enable_download         = false
  	automatic_case_creation = true
  	enable_remote_support   = true
  	accepted_terms          = true
}
`

var supportAssistResourceConfigUpdate = `
resource "powerscale_support_assist" "test" {
	supportassist_enabled   = false
  	enable_download         = false
  	automatic_case_creation = false
  	enable_remote_support   = true
  	accepted_terms                = false
	telemetry = {
		offline_collection_period = 7800,
		telemetry_enabled         = false,
		telemetry_persist		  = true,
		telemetry_threads         = 5
	}
	contact = {
		primary = {
		  	email      = "abc@gmail.com",
		  	first_name = "terraform_first",
		  	language   = "No",
		  	last_name  = "terraform_last",
		  	phone      = "1234567890"
		},
		secondary = {
		  	email 	   = "xyz@gmail.com",
			first_name = "terraform_sec_first",
			language   = "En",
			last_name  = "terraform_sec_last",
			phone      = "1234567980"
		}
	}
	connections = {
		mode = "gateway"
		gateway_endpoints = [
		  	{
		  		enabled = false,
		  		host = "1.2.3.4",
		  		port = 9443,
		  		priority = 1,
		  		use_proxy = true,
		  		validate_ssl = true
			},
		],
		network_pools = ["subnet0:pool0"]
	}
}
`

var supportAssistResourceConfigUpdate1 = `
resource "powerscale_support_assist" "test" {
	supportassist_enabled   = false
  	enable_download         = false
  	automatic_case_creation = true
  	enable_remote_support   = true
  	accepted_terms                = false
	telemetry = {
		offline_collection_period = 7800,
		telemetry_enabled         = false,
		telemetry_persist		  = true,
		telemetry_threads         = 5
	}
	contact = {
		primary = {
		  	email      = "abc@gmail.com",
		  	first_name = "terraform_first",
		  	language   = "No",
		  	last_name  = "terraform_last",
		  	phone      = "1234567890"
		},
		secondary = {
		  	email 	   = "xyz@gmail.com",
			first_name = "terraform_sec_first",
			language   = "En",
			last_name  = "terraform_sec_last",
			phone      = "1234567980"
		}
	}
	connections = {
		mode = "gateway"
		gateway_endpoints = [
		  	{
		  		enabled = false,
		  		host = "1.2.3.4",
		  		port = 9443,
		  		priority = 1,
		  		use_proxy = true,
		  		validate_ssl = true
			},
		],
		network_pools = ["subnet0:pool0"]
	}
}
`
