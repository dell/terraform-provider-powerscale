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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var networkSettingGetMocker *Mocker
var networkSettingMocker *Mocker

func TestAccNetworkSettingResourceImport(t *testing.T) {
	var networkSettingResourceName = "powerscale_network_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + networkSettingBasicResourceConfig,
			},
			// Import testing
			{
				Config:            ProviderConfig + networkSettingBasicResourceConfig,
				ResourceName:      networkSettingResourceName,
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "groupnet0", s[0].Attributes["default_groupnet"])
					return nil
				},
			},
			// Update and Read testing
			{
				Config: ProviderConfig + networkSettingUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkSettingResourceName, "source_based_routing_enabled", "true"),
					resource.TestCheckResourceAttr(networkSettingResourceName, "sc_rebalance_delay", "10"),
					resource.TestCheckResourceAttr(networkSettingResourceName, "default_groupnet", "groupnet0"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + networkSettingUpdateRevertResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(networkSettingResourceName, "source_based_routing_enabled", "false"),
					resource.TestCheckResourceAttr(networkSettingResourceName, "sc_rebalance_delay", "0"),
					resource.TestCheckResourceAttr(networkSettingResourceName, "default_groupnet", "groupnet0"),
				),
			},
		},
	})
}

func TestAccNetworkSettingResourceErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + networkSettingBasicResourceConfig,
			},
			// Update Error - invalid sc_rebalance_delay
			{
				Config:      ProviderConfig + networkSettingInvalidDelayResourceConfig,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value*.`),
			},
			// Update Error - invalid tcp_ports
			{
				Config:      ProviderConfig + networkSettingInvalidPortResourceConfig,
				ExpectError: regexp.MustCompile(`.*Invalid Attribute Value*.`),
			},
			{
				Config: ProviderConfig + networkSettingBasicResourceConfig,
			},
		},
	})
}

func TestAccNetworkSettingResourceImportMockErr(t *testing.T) {
	var networkSettingResourceName = "powerscale_network_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.UnPatch()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.UnPatch()
					}
				},
				Config: ProviderConfig + networkSettingBasicResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.UnPatch()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.UnPatch()
					}
					networkSettingGetMocker = Mock(helper.GetNetworkSetting).Return(nil, fmt.Errorf("networkSettings read mock error")).Build()
				},
				Config:            ProviderConfig + networkSettingBasicResourceConfig,
				ResourceName:      networkSettingResourceName,
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*networkSettings read mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkSettingResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read Error testing
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.UnPatch()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.UnPatch()
					}
					networkSettingGetMocker = Mock(helper.GetNetworkSetting).Return(nil, fmt.Errorf("networkSettings read mock error")).Build()
				},
				Config:      ProviderConfig + networkSettingBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*networkSettings read mock error*.`),
			},
			// Create and Read Error testing
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.UnPatch()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.UnPatch()
					}
					networkSettingGetMocker = Mock(helper.GetNetworkSetting).To(func(ctx context.Context, powerscaleClient *client.Client) (*powerscale.V12NetworkExternalSettings, error) {
						if networkSettingGetMocker.MockTimes() > 0 {
							return nil, fmt.Errorf("networkSettings read mock error")
						}
						return mockNetworkSettingsOrigin, nil
					}).Build()
				},
				Config:      ProviderConfig + networkSettingBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*networkSettings read mock error*.`),
			},
		},
	})
}

func TestAccNetworkSettingResourceUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import testing
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.UnPatch()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.UnPatch()
					}
				},
				Config: ProviderConfig + networkSettingBasicResourceConfig,
			},
			// Update and Read Error testing
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.UnPatch()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.UnPatch()
					}
					networkSettingMocker = Mock(helper.UpdateNetworkSetting).Return(nil).Build()
					networkSettingGetMocker = Mock(helper.GetNetworkSetting).To(func(ctx context.Context, powerscaleClient *client.Client) (*powerscale.V12NetworkExternalSettings, error) {
						if networkSettingGetMocker.MockTimes() > 0 {
							return nil, fmt.Errorf("networkSettings update mock error")
						}
						return mockNetworkSettingsOrigin, nil
					}).Build()
				},
				Config:      ProviderConfig + networkSettingUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*networkSettings update mock error*.`),
			},
		},
	})
}

func TestAccNetworkSettingReleaseMockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if networkSettingGetMocker != nil {
						networkSettingGetMocker.Release()
					}
					if networkSettingMocker != nil {
						networkSettingMocker.Release()
					}
				},
				Config: ProviderConfig + networkSettingBasicResourceConfig,
			},
		},
	})
}

var networkSettingBasicResourceConfig = `
resource "powerscale_network_settings" "test" {
  }
`

var networkSettingUpdateResourceConfig = `
resource "powerscale_network_settings" "test" {
	source_based_routing_enabled = true
	sc_rebalance_delay = 10
  }
`

var networkSettingUpdateRevertResourceConfig = `
resource "powerscale_network_settings" "test" {
	source_based_routing_enabled = false
	sc_rebalance_delay = 0
  }
`

var networkSettingInvalidDelayResourceConfig = `
resource "powerscale_network_settings" "test" {
	sc_rebalance_delay = 1000
  }
`

var networkSettingInvalidPortResourceConfig = `
resource "powerscale_network_settings" "test" {
	tcp_ports = [1000000]
  }
`

var mockNetworkSettingsOrigin = &powerscale.V12NetworkExternalSettings{
	Sbr:              false,
	ScRebalanceDelay: 0,
	DefaultGroupnet:  "groupnet0",
	TcpPorts:         []int64{20, 21, 80, 445, 2049},
}
