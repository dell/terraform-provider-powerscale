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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSmbShareSettingsDataSourceReadWithoutFilter(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithoutFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_share_settings.HostACL"),
					/*resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.enable_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_netbios"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.access_based_share_enum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.ignore_eas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_num_workers"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.require_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_side_copy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb3_encryption"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.guest_user"),*/
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadWithFilter(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.scopeandzone"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_share_settings.HostACL"),
					/*resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.enable_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_netbios"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.access_based_share_enum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.ignore_eas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_num_workers"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.require_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_side_copy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb3_encryption"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.guest_user"),
					*/
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadWithEmptyFilter(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.empty"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithEmptyFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_share_settings.HostACL"),
					/*resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.enable_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_netbios"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.access_based_share_enum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.ignore_eas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_num_workers"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.require_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_side_copy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb3_encryption"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.guest_user"),
					*/
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadWithFilterDefaultScope(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.default"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithFilterDefaultScope,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_share_settings.HostACL"),
					/*resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.enable_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_netbios"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.access_based_share_enum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.ignore_eas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_num_workers"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.require_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_side_copy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb3_encryption"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.guest_user"),
					*/
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadWithFilterUserScope(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.user"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithFilterUserScope,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadWithFilterSystemZone(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.system"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithFilterSystemZone,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
					/*resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.enable_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_netbios"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.access_based_share_enum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.ignore_eas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_num_workers"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.require_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_side_copy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb3_encryption"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.guest_user"),
					*/
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadWithFiltereffectiveScope(t *testing.T) {
	dataSourceName := "data.powerscale_smb_share_settings.effective"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbShareSettingsDataSourceConfigWithFilterEffectiveScope,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_share_settings_System"),
					/*resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.enable_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_netbios"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.access_based_share_enum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_root"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_accessible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.ignore_eas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_num_workers"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.dot_snap_visible_child"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.require_security_signatures"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_side_copy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.server_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_smb3_encryption"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.reject_unencrypted_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.onefs_cpu_multiplier"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.guest_user"),
					*/
				),
			},
		},
	})
}

func TestAccSmbShareSettingsDataSourceReadMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.FilterSmbShareSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbShareSettingsDataSourceConfigWithFilter,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbShareSettingsDataSourceConfigWithFilter,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var smbShareSettingsDataSourceConfigWithoutFilter = `
data "powerscale_smb_share_settings" "all" {
}
`

var smbShareSettingsDataSourceConfigWithEmptyFilter = `
data "powerscale_smb_share_settings" "empty" {
	filter {
	}
}
`

var smbShareSettingsDataSourceConfigWithFilterSystemZone = `
data "powerscale_smb_share_settings" "system" {
	filter {
		zone = "System"
	}
}
`

var smbShareSettingsDataSourceConfigWithFilterDefaultScope = `
data "powerscale_smb_share_settings" "default" {
	filter {
		scope = "default"
	}
}
`

var smbShareSettingsDataSourceConfigWithFilterUserScope = `
data "powerscale_smb_share_settings" "user" {
	filter {
		scope = "user"
	}
}
`

var smbShareSettingsDataSourceConfigWithFilterEffectiveScope = `
data "powerscale_smb_share_settings" "effective" {
	filter {
		scope = "effective"
	}
}
`

var smbShareSettingsDataSourceConfigWithFilter = `
data "powerscale_smb_share_settings" "scopeandzone" {
	filter {
		scope = "effective"
		zone = "System"
	}
}
`
