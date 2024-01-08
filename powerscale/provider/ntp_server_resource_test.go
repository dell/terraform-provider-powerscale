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
	powerscale "dell/powerscale-go-client"
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"
)

func TestAccNtpServerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NtpServerResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "name", ntpServerName),
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "key", ntpServerKey),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_ntpserver.ntp_server_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, ntpServerName, states[0].Attributes["name"])
					assert.Equal(t, ntpServerKey, states[0].Attributes["key"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + NtpServerUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "name", ntpServerName),
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "key", ntpServerKey+"_updated_1"),
				),
			},
		},
	})
}

func TestAccNtpServerResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NtpServerResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "name", ntpServerName),
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "key", ntpServerKey),
				),
			},
			// ImportState testing get none ntp server
			{
				ResourceName: "powerscale_ntpserver.ntp_server_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetNtpServer).Return(&powerscale.V3NtpServersExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile("not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_ntpserver.ntp_server_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNtpServer).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNtpServerResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NtpServerResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "name", ntpServerName),
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "key", ntpServerKey),
				),
			},
			// Update param read error
			{
				Config: ProviderConfig + NtpServerUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + NtpServerUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.UpdateNtpServer).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none ntp server
			{
				Config: ProviderConfig + NtpServerUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNtpServer).Return(&powerscale.V3NtpServersExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, ntpServerModel models.NtpServerResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".*Could not read updated ntp server*."),
			},
			// Update get error
			{
				Config: ProviderConfig + NtpServerUpdatedResourceConfig2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNtpServer).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, ntpServerModel models.NtpServerResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config:      ProviderConfig + NtpServerUpdatePreCheckConfig,
				ExpectError: regexp.MustCompile(".*Should not provide parameters for creating*."),
			},
		},
	})
}

func TestAccNtpServerResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NtpServerResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CreateNtpServer).Return("", fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNtpServerResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NtpServerResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "name", ntpServerName),
					resource.TestCheckResourceAttr("powerscale_ntpserver.ntp_server_test", "key", ntpServerKey),
				),
			},
			{
				ResourceName: "powerscale_ntpserver.ntp_server_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NtpServerResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + NtpServerUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNtpServerResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NtpServerResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var ntpServerName = "ntp_server_at_test"
var ntpServerKey = "ntp_server_key"

var NtpServerResourceConfig = fmt.Sprintf(`
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "%s"
	key = "%s"
}
`, ntpServerName, ntpServerKey)

var NtpServerUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "%s"
	key = "%s"
}
`, ntpServerName, ntpServerKey+"_updated_1")

var NtpServerUpdatedResourceConfig2 = fmt.Sprintf(`
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "%s"
	key = "%s"
}
`, ntpServerName, ntpServerKey+"_updated_2")

var NtpServerUpdatePreCheckConfig = fmt.Sprintf(`
resource "powerscale_ntpserver" "ntp_server_test" {
	name = "%s"
	key = "%s"
}
`, ntpServerName+"_updated", ntpServerKey)
