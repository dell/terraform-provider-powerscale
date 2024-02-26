/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/stretchr/testify/assert"
	"math/big"
	"regexp"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSmartPoolSettingsResourceCreate(t *testing.T) {
	var data = "powerscale_smartpool_settings.settings"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(data, "global_namespace_acceleration_enabled", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "global_namespace_acceleration_state", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "manage_io_optimization", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "manage_io_optimization_apply_to_files", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "manage_protection", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "manage_protection_apply_to_files", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "protect_directories_one_level_higher", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "spillover_enabled", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "spillover_target.name", regexp.MustCompile(`^\d+|\w*$`)),
					resource.TestMatchResourceAttr(data, "spillover_target.type", regexp.MustCompile(`^\w*$`)),
					resource.TestMatchResourceAttr(data, "ssd_l3_cache_default_enabled", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "ssd_qab_mirrors", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "ssd_system_btree_mirrors", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "ssd_system_delta_mirrors", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_deny_writes", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_hide_spare", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_limit_drives", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_limit_percent", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdate(t *testing.T) {
	var restV5UpdateFuncMocker *Mocker
	var restV16UpdateFuncMocker *Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create
			{
				Config: ProviderConfig + UpdatePoolSettingResourceConfig,
				PreConfig: func() {
					restV5UpdateFuncMocker = Mock(powerscale.ApiUpdateStoragepoolv5StoragepoolSettingsRequest.Execute).Return(nil, nil).Build()
					restV16UpdateFuncMocker = Mock(powerscale.ApiUpdateStoragepoolv16StoragepoolSettingsRequest.Execute).Return(nil, nil).Build()
					FunctionMocker = Mock(helper.GetSmartPoolSettings).To(func(ctx context.Context, powerscaleClient *client.Client) (any, error) {
						onefsVersion, _ := powerscaleClient.GetOnefsVersion()
						if restV5UpdateFuncMocker.MockTimes() > 0 {
							return mockV5StoragepoolSettingsAfterUpdate, nil
						}
						if restV16UpdateFuncMocker.MockTimes() > 0 {
							return mockV16StoragepoolSettingsAfterUpdate, nil
						}
						if onefsVersion.IsGreaterThan("9.4.0") {
							return mockV16StoragepoolSettingsBeforeUpdate, nil
						}
						return mockV5StoragepoolSettingsBeforeUpdate, nil
					}).Build()
				},
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			if restV5UpdateFuncMocker != nil {
				restV5UpdateFuncMocker.UnPatch()
			}
			if restV16UpdateFuncMocker != nil {
				restV16UpdateFuncMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdateIoOptimizationErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create
			{
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
			},
			// error updating manage_io_optimization
			{
				Config:      ProviderConfig + errUpdateManageIoOptimizationConfig,
				ExpectError: regexp.MustCompile(`.*Input validation failed*.`),
			},
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdateManageProtectionErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create
			{
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
			},
			// error updating manage_protection
			{
				Config:      ProviderConfig + errUpdateManageProtectionConfig,
				ExpectError: regexp.MustCompile(`.*Input validation failed*.`),
			},
		},
	})
}

func TestAccSmartPoolSettingsResourceCreateErrorRequest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetSmartPoolSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmartPoolSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error reading SmartPool settings*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdateErrorRequest(t *testing.T) {
	var updateFuncMocker *Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					updateFuncMocker = Mock(helper.UpdateSmartPoolSettings).Return(nil).Build()
					FunctionMocker = Mock(helper.GetSmartPoolSettings).Build().
						To(func(ctx context.Context, powerscaleClient *client.Client) (any, error) {
							if updateFuncMocker.MockTimes() > 0 {
								return nil, fmt.Errorf("mock error")
							}
							onefsVersion, _ := powerscaleClient.GetOnefsVersion()
							if onefsVersion.IsGreaterThan("9.4.0") {
								return mockV16StoragepoolSettingsBeforeUpdate, nil
							}
							return mockV5StoragepoolSettingsBeforeUpdate, nil
						})
				},
				Config:      ProviderConfig + UpdatePoolSettingResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error reading SmartPool settings*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			if updateFuncMocker != nil {
				updateFuncMocker.UnPatch()
			}
			return nil
		},
	})
}

// TestAccSmartPoolSettingsResourceV16 UT for SmartPoolSettingsResource using mock response from PowerScale 9.5.
func TestAccSmartPoolSettingsResourceV16(t *testing.T) {
	var data = "powerscale_smartpool_settings.settings"
	var restV5UpdateFuncMocker *Mocker
	var restV16UpdateFuncMocker *Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					restV5UpdateFuncMocker = Mock(powerscale.ApiUpdateStoragepoolv5StoragepoolSettingsRequest.Execute).Return(nil, nil).Build()
					restV16UpdateFuncMocker = Mock(powerscale.ApiUpdateStoragepoolv16StoragepoolSettingsRequest.Execute).Return(nil, nil).Build()
					FunctionMocker = Mock(helper.GetSmartPoolSettings).Build().
						When(func(ctx context.Context, powerscaleClient *client.Client) bool {
							onefsVersion, _ := powerscaleClient.GetOnefsVersion()
							if strings.Contains(powerscaleClient.PscaleOpenAPIClient.GetConfig().Servers[0].URL, "localhost") {
								powerscaleClient.SetOnefsVersion(9, 5, 0)
								return false
							}
							if !strings.Contains(powerscaleClient.PscaleOpenAPIClient.GetConfig().Servers[0].URL, "localhost") &&
								onefsVersion.IsGreaterThan("9.4.0") {
								return false
							}
							return true
						}).Return(mockV16StoragepoolSettings, nil)

				},
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(data, "default_transfer_limit_state", regexp.MustCompile(`^$|\w+`)),
					resource.TestMatchResourceAttr(data, "default_transfer_limit_pct", regexp.MustCompile(`^$|^\d+$`)),
					resource.TestMatchResourceAttr(data, "global_namespace_acceleration_enabled", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "global_namespace_acceleration_state", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "manage_io_optimization", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "manage_io_optimization_apply_to_files", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "manage_protection", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "manage_protection_apply_to_files", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "protect_directories_one_level_higher", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "spillover_enabled", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "spillover_target.name", regexp.MustCompile(`^\d+|\w*$`)),
					resource.TestMatchResourceAttr(data, "spillover_target.type", regexp.MustCompile(`^\w*$`)),
					resource.TestMatchResourceAttr(data, "ssd_l3_cache_default_enabled", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "ssd_qab_mirrors", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "ssd_system_btree_mirrors", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "ssd_system_delta_mirrors", regexp.MustCompile(`^\w+$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_deny_writes", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_hide_spare", regexp.MustCompile(`^false|true$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_limit_drives", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_limit_percent", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			if restV5UpdateFuncMocker != nil {
				restV5UpdateFuncMocker.UnPatch()
			}
			if restV16UpdateFuncMocker != nil {
				restV16UpdateFuncMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdateErrorPutRequest(t *testing.T) {
	var updateFuncMocker *Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create
			{
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
				PreConfig: func() {
					updateFuncMocker = Mock(helper.UpdateSmartPoolSettings).Return(fmt.Errorf("mock error")).Build()
					FunctionMocker = Mock(helper.GetSmartPoolSettings).To(func(ctx context.Context, powerscaleClient *client.Client) (any, error) {
						if updateFuncMocker.MockTimes() > 0 {
							return mockV5StoragepoolSettingsAfterUpdate, nil
						}
						return mockV5StoragepoolSettingsBeforeUpdate, nil
					}).Build()
				},
				ExpectError: regexp.MustCompile(`.*Error updating SmartPool settings*.`),
			},
			// Error updating
			{
				Config:      ProviderConfig + UpdatePoolSettingResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error updating SmartPool settings*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			if updateFuncMocker != nil {
				updateFuncMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdateErrorGetRequest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetSmartPoolSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmartPoolSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error reading SmartPool settings*.`),
			},
			// Error updating
			{
				Config:      ProviderConfig + UpdatePoolSettingResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error reading SmartPool settings*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccSmartPoolSettingsResourceUpdateErrorUpdatingModel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.UpdateSmartPoolSettingsResourceModel).Return("error summary", "error detail").Build()
				},
				Config:      ProviderConfig + SmartPoolSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*error summary*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccSmartPoolSettingsResourceImport(t *testing.T) {
	var data = "powerscale_smartpool_settings.settings"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
			},
			{
				ResourceName:      data,
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestMatchResourceAttr(data, "global_namespace_acceleration_enabled", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "global_namespace_acceleration_state", regexp.MustCompile(`^\w+$`))
					resource.TestMatchResourceAttr(data, "manage_io_optimization", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "manage_io_optimization_apply_to_files", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "manage_protection", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "manage_protection_apply_to_files", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "protect_directories_one_level_higher", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "spillover_enabled", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "spillover_target.name", regexp.MustCompile(`^\d+|\w*$`))
					resource.TestMatchResourceAttr(data, "spillover_target.type", regexp.MustCompile(`^\w*$`))
					resource.TestMatchResourceAttr(data, "ssd_l3_cache_default_enabled", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "ssd_qab_mirrors", regexp.MustCompile(`^\w+$`))
					resource.TestMatchResourceAttr(data, "ssd_system_btree_mirrors", regexp.MustCompile(`^\w+$`))
					resource.TestMatchResourceAttr(data, "ssd_system_delta_mirrors", regexp.MustCompile(`^\w+$`))
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_deny_writes", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_hide_spare", regexp.MustCompile(`^false|true$`))
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_limit_drives", regexp.MustCompile(`^\d+$`))
					resource.TestMatchResourceAttr(data, "virtual_hot_spare_limit_percent", regexp.MustCompile(`^\d+$`))
					return nil
				},
			},
		},
	})
}

func TestAccSmartPoolSettingsResourceImportErr(t *testing.T) {
	var data = "powerscale_smartpool_settings.settings"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Import SmartPool Settings resource Error testing
			{
				Config: ProviderConfig + SmartPoolSettingsResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetSmartPoolSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ResourceName:      data,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*Error reading SmartPool settings*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestBigFloatToInt32(t *testing.T) {
	testAccPreCheck(t)

	val, _ := helper.BigFloatToInt32(nil)
	assert.True(t, val == nil)

	bigFloat := big.NewFloat(92233720368547758071)
	_, err := helper.BigFloatToInt32(bigFloat)
	assert.True(t, err != nil)

	bigFloat = big.NewFloat(12345)
	val, _ = helper.BigFloatToInt32(bigFloat)
	assert.True(t, *val == 12345)
}

var SmartPoolSettingsResourceConfig = `
resource "powerscale_smartpool_settings" "settings" {
  
}
`

var UpdatePoolSettingResourceConfig = `
resource "powerscale_smartpool_settings" "settings" {
    manage_io_optimization                = true
    manage_io_optimization_apply_to_files = true
    manage_protection                     = true
    manage_protection_apply_to_files      = true
	protect_directories_one_level_higher  = true
    spillover_enabled                     = true
    spillover_target                      = {
      name    = ""
      type    = "anywhere"
    }
    ssd_l3_cache_default_enabled          = true
    ssd_qab_mirrors                       = "one"
    ssd_system_btree_mirrors              = "one"
    ssd_system_delta_mirrors              = "one"
    virtual_hot_spare_deny_writes         = true
    virtual_hot_spare_hide_spare          = true
    virtual_hot_spare_limit_drives        = 2
    virtual_hot_spare_limit_percent       = 2
  # Note that, default_transfer_limit_state and default_transfer_limit_pct are mutually exclusive and only one can be specified.
  # default_transfer_limit_state          = "disabled" // available for PowerScale 9.5 and above
  # default_transfer_limit_pct            = 90 // available for PowerScale 9.5 and above
}
`

var errUpdateManageProtectionConfig = `
resource "powerscale_smartpool_settings" "settings" {
    manage_protection                     = false
    manage_protection_apply_to_files      = true
}
`

var errUpdateManageIoOptimizationConfig = `
resource "powerscale_smartpool_settings" "settings" {
    manage_io_optimization                = false
    manage_io_optimization_apply_to_files = true
}
`

var mockV5StoragepoolSettingsBeforeUpdate = &powerscale.V5StoragepoolSettings{
	Settings: &powerscale.V5StoragepoolSettingsSettings{
		AutomaticallyManageIoOptimization:  "all",
		AutomaticallyManageProtection:      "all",
		GlobalNamespaceAccelerationEnabled: false,
		GlobalNamespaceAccelerationState:   "inactive",
		ProtectDirectoriesOneLevelHigher:   true,
		SpilloverEnabled:                   true,
		SpilloverTarget:                    powerscale.V1StoragepoolSettingsSettingsSpilloverTarget{Id: 0, Name: "", Type: "anywhere"},
		SsdL3CacheDefaultEnabled:           true,
		SsdQabMirrors:                      "one",
		SsdSystemBtreeMirrors:              "one",
		SsdSystemDeltaMirrors:              "one",
		VirtualHotSpareDenyWrites:          true,
		VirtualHotSpareHideSpare:           true,
		VirtualHotSpareLimitDrives:         2,
		VirtualHotSpareLimitPercent:        2,
	},
}

var mockV16StoragepoolSettingsBeforeUpdate = &powerscale.V16StoragepoolSettings{
	Settings: &powerscale.V16StoragepoolSettingsSettings{
		AutomaticallyManageIoOptimization:  "all",
		AutomaticallyManageProtection:      "all",
		GlobalNamespaceAccelerationEnabled: false,
		GlobalNamespaceAccelerationState:   "inactive",
		ProtectDirectoriesOneLevelHigher:   true,
		SpilloverEnabled:                   true,
		SpilloverTarget:                    powerscale.V1StoragepoolSettingsSettingsSpilloverTarget{Id: 0, Name: "", Type: "anywhere"},
		SsdL3CacheDefaultEnabled:           true,
		SsdQabMirrors:                      "one",
		SsdSystemBtreeMirrors:              "one",
		SsdSystemDeltaMirrors:              "one",
		VirtualHotSpareDenyWrites:          true,
		VirtualHotSpareHideSpare:           true,
		VirtualHotSpareLimitDrives:         2,
		VirtualHotSpareLimitPercent:        2,
		DefaultTransferLimitPct:            nil,
		DefaultTransferLimitState:          nil,
	},
}

var mockV5StoragepoolSettingsAfterUpdate = &powerscale.V5StoragepoolSettings{
	Settings: &powerscale.V5StoragepoolSettingsSettings{
		AutomaticallyManageIoOptimization:  "files_at_default",
		AutomaticallyManageProtection:      "files_at_default",
		GlobalNamespaceAccelerationEnabled: false,
		GlobalNamespaceAccelerationState:   "inactive",
		ProtectDirectoriesOneLevelHigher:   true,
		SpilloverEnabled:                   true,
		SpilloverTarget:                    powerscale.V1StoragepoolSettingsSettingsSpilloverTarget{Id: 0, Name: "", Type: "anywhere"},
		SsdL3CacheDefaultEnabled:           true,
		SsdQabMirrors:                      "one",
		SsdSystemBtreeMirrors:              "one",
		SsdSystemDeltaMirrors:              "one",
		VirtualHotSpareDenyWrites:          true,
		VirtualHotSpareHideSpare:           true,
		VirtualHotSpareLimitDrives:         2,
		VirtualHotSpareLimitPercent:        2,
	},
}

var mockV16StoragepoolSettingsAfterUpdate = &powerscale.V16StoragepoolSettings{
	Settings: &powerscale.V16StoragepoolSettingsSettings{
		AutomaticallyManageIoOptimization:  "files_at_default",
		AutomaticallyManageProtection:      "files_at_default",
		GlobalNamespaceAccelerationEnabled: false,
		GlobalNamespaceAccelerationState:   "inactive",
		ProtectDirectoriesOneLevelHigher:   true,
		SpilloverEnabled:                   true,
		SpilloverTarget:                    powerscale.V1StoragepoolSettingsSettingsSpilloverTarget{Id: 0, Name: "", Type: "anywhere"},
		SsdL3CacheDefaultEnabled:           true,
		SsdQabMirrors:                      "one",
		SsdSystemBtreeMirrors:              "one",
		SsdSystemDeltaMirrors:              "one",
		VirtualHotSpareDenyWrites:          true,
		VirtualHotSpareHideSpare:           true,
		VirtualHotSpareLimitDrives:         2,
		VirtualHotSpareLimitPercent:        2,
		DefaultTransferLimitPct:            nil,
		DefaultTransferLimitState:          nil,
	},
}
