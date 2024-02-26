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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAclSettingsDataSourceAll(t *testing.T) {
	var aclServerTerraformName = "data.powerscale_aclsettings.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + ACLSettingsAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(aclServerTerraformName, "create_over_smb", "allow"),
					resource.TestCheckResourceAttr(aclServerTerraformName, "dos_attr", "deny_smb"),
					resource.TestCheckResourceAttr(aclServerTerraformName, "chmod_inheritable", "no"),
					resource.TestCheckResourceAttr(aclServerTerraformName, "chmod_007", "default"),
				),
			},
		},
	})
}

func TestAccAclSettingsDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetACLSettings).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + ACLSettingsAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccAclSettingsDataSourceMappingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ACLSettingsDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + ACLSettingsAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var ACLSettingsAllDataSourceConfig = `
resource "powerscale_aclsettings" "acl_settings_test" {
    access                  = "windows"
    calcmode                = "approx"
    calcmode_group          = "group_aces"
    calcmode_owner          = "owner_aces"
    calcmode_traverse       = "ignore"
    chmod                   = "merge"
    chmod_007               = "default"
    chmod_inheritable       = "no"
    chown                   = "owner_group_and_acl"
    create_over_smb         = "allow"
    dos_attr                = "deny_smb"
    group_owner_inheritance = "creator"
    rwx                     = "retain"
    synthetic_denies        = "remove"
    utimes                  = "only_owner"
}

data "powerscale_aclsettings" "all" {
    depends_on = [
        powerscale_aclsettings.acl_settings_test
    ]
}
`
