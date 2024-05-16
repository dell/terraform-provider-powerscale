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

func TestAccSmbServerSettingsDataSourceReadWithoutFilter(t *testing.T) {
	dataSourceName := "data.powerscale_smb_server_settings.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbServerSettingsDataSourceConfigWithoutFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_server_settings_effective"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
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
				),
			},
		},
	})
}

func TestAccSmbServerSettingsDataSourceReadWithFilter(t *testing.T) {
	dataSourceName := "data.powerscale_smb_server_settings.effective"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbServerSettingsDataSourceConfigWithFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_server_settings_effective"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
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
				),
			},
		},
	})
}

func TestAccSmbServerSettingsDataSourceReadWithEmptyFilter(t *testing.T) {
	dataSourceName := "data.powerscale_smb_server_settings.empty"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbServerSettingsDataSourceConfigWithEmptyFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_server_settings_effective"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
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
				),
			},
		},
	})
}

func TestAccSmbServerSettingsDataSourceReadWithFilterDefaultScope(t *testing.T) {
	dataSourceName := "data.powerscale_smb_server_settings.default"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbServerSettingsDataSourceConfigWithFilterDefaultScope,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_server_settings_default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "smb_server_settings.support_multichannel"),
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
				),
			},
		},
	})
}

func TestAccSmbServerSettingsDataSourceReadWithFilterUserScope(t *testing.T) {
	dataSourceName := "data.powerscale_smb_server_settings.user"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
			{
				Config: ProviderConfig + smbServerSettingsDataSourceConfigWithFilterUserScope,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "smb_server_settings_user"),
				),
			},
		},
	})
}

func TestAccSmbServerSettingsDataSourceReadMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.FilterSmbServerSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsDataSourceConfigWithFilter,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + smbServerSettingsDataSourceConfigWithFilter,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var smbServerSettingsDataSourceConfigWithoutFilter = `
data "powerscale_smb_server_settings" "all" {
}
`

var smbServerSettingsDataSourceConfigWithEmptyFilter = `
data "powerscale_smb_server_settings" "empty" {
	filter {
	}
}
`

var smbServerSettingsDataSourceConfigWithFilter = `
data "powerscale_smb_server_settings" "effective" {
	filter {
		scope = "effective"
	}
}
`

var smbServerSettingsDataSourceConfigWithFilterDefaultScope = `
data "powerscale_smb_server_settings" "default" {
	filter {
		scope = "default"
	}
}
`

var smbServerSettingsDataSourceConfigWithFilterUserScope = `
data "powerscale_smb_server_settings" "user" {
	filter {
		scope = "user"
	}
}
`
