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

func TestAccSmbServerSettingsResourceCreate(t *testing.T) {
	resourceName := "powerscale_smb_server_settings.example"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// create and read testing
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "smb_server_settings_effective"),
					resource.TestCheckResourceAttr(resourceName, "scope", "effective"),
					resource.TestCheckResourceAttrSet(resourceName, "access_based_share_enum"),
					resource.TestCheckResourceAttrSet(resourceName, "dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(resourceName, "dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(resourceName, "dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(resourceName, "dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_security_signatures"),
					resource.TestCheckResourceAttrSet(resourceName, "guest_user"),
					resource.TestCheckResourceAttrSet(resourceName, "ignore_eas"),
					resource.TestCheckResourceAttrSet(resourceName, "onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(resourceName, "onefs_num_workers"),
					resource.TestCheckResourceAttrSet(resourceName, "reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(resourceName, "require_security_signatures"),
					resource.TestCheckResourceAttrSet(resourceName, "server_side_copy"),
					resource.TestCheckResourceAttrSet(resourceName, "server_string"),
					resource.TestCheckResourceAttrSet(resourceName, "service"),
					resource.TestCheckResourceAttrSet(resourceName, "support_multichannel"),
					resource.TestCheckResourceAttrSet(resourceName, "support_netbios"),
					resource.TestCheckResourceAttrSet(resourceName, "support_smb2"),
					resource.TestCheckResourceAttrSet(resourceName, "support_smb3_encryption"),
				),
			},
		},
	})
}

func TestAccSmbServerSettingsResourceImport(t *testing.T) {
	resourceName := "powerscale_smb_server_settings.example"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigBasic,
			},
			// import testing
			{
				ResourceName:      resourceName,
				ImportStateId:     "effective",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSmbServerSettingsResourceUpdate(t *testing.T) {
	resourceName := "powerscale_smb_server_settings.example"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigBasic,
			},
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigDefault,
			},
			// update and read testing
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "smb_server_settings_effective"),
					resource.TestCheckResourceAttr(resourceName, "scope", "effective"),
					resource.TestCheckResourceAttr(resourceName, "access_based_share_enum", "true"),
					resource.TestCheckResourceAttr(resourceName, "dot_snap_accessible_child", "false"),
					resource.TestCheckResourceAttr(resourceName, "dot_snap_accessible_root", "false"),
					resource.TestCheckResourceAttr(resourceName, "dot_snap_visible_child", "true"),
					resource.TestCheckResourceAttr(resourceName, "dot_snap_visible_root", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_security_signatures", "true"),
					resource.TestCheckResourceAttr(resourceName, "guest_user", "everybody"),
					resource.TestCheckResourceAttr(resourceName, "ignore_eas", "true"),
					resource.TestCheckResourceAttr(resourceName, "onefs_cpu_multiplier", "1"),
					resource.TestCheckResourceAttr(resourceName, "onefs_num_workers", "4"),
					resource.TestCheckResourceAttr(resourceName, "reject_unencrypted_access", "false"),
					resource.TestCheckResourceAttr(resourceName, "require_security_signatures", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_side_copy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_string", "PowerScale Server Updated"),
					resource.TestCheckResourceAttr(resourceName, "service", "false"),
					resource.TestCheckResourceAttr(resourceName, "support_multichannel", "false"),
					resource.TestCheckResourceAttr(resourceName, "support_netbios", "true"),
					resource.TestCheckResourceAttr(resourceName, "support_smb2", "false"),
					resource.TestCheckResourceAttr(resourceName, "support_smb3_encryption", "true"),
				),
			},
		},
	})
}

func TestAccSmbServerSettingsCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSmbServerSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSmbServerSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSmbServerSettingsReadMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigBasic,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSmbServerSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigBasic,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSmbServerSettingsUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigBasic,
			},
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigDefault,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigUpdated,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSmbServerSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigUpdated,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSmbServerSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigUpdated,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsResourceConfigUpdated,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSmbServerSettingsImportMockErr(t *testing.T) {
	resourceName := "powerscale_smb_server_settings.example"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + smbServerSettingsResourceConfigBasic,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSmbServerSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + smbServerSettingsResourceConfigBasic,
				ResourceName:      resourceName,
				ImportStateId:     "effective",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + smbServerSettingsResourceConfigBasic,
				ResourceName:      resourceName,
				ImportStateId:     "effective",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var smbServerSettingsResourceConfigBasic = `
resource "powerscale_smb_server_settings" "example" {
	scope = "effective"
}
`

var smbServerSettingsResourceConfigUpdated = `
resource "powerscale_smb_server_settings" "example" {
	scope = "effective"
	access_based_share_enum = true
    dot_snap_accessible_child = false
    dot_snap_accessible_root = false
    dot_snap_visible_child = true
    dot_snap_visible_root = false
    enable_security_signatures = true
    guest_user = "everybody"
    ignore_eas = true
    onefs_cpu_multiplier = 1
    onefs_num_workers = 4
    reject_unencrypted_access = false
    require_security_signatures = true
    server_side_copy = false
    server_string = "PowerScale Server Updated"
    service = false
    support_multichannel = false
    support_netbios = true
    support_smb2 = false
    support_smb3_encryption = true
}
`

var smbServerSettingsResourceConfigDefault = `
resource "powerscale_smb_server_settings" "example" {
	scope = "effective"
	access_based_share_enum = false
    dot_snap_accessible_child = true
    dot_snap_accessible_root = true
    dot_snap_visible_child = false
    dot_snap_visible_root = true
    enable_security_signatures = false
    guest_user = "nobody"
    ignore_eas = false
    onefs_cpu_multiplier = 4
    onefs_num_workers = 0
    reject_unencrypted_access = true
    require_security_signatures = false
    server_side_copy = true
    server_string = "PowerScale Server"
    service = true
    support_multichannel = true
    support_netbios = false
    support_smb2 = true
    support_smb3_encryption = false
}
`
