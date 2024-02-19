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
	"context"
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

func TestAccAclSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + ACLSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "create_over_smb", aclCreateOverSmb),
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "dos_attr", aclReadOnlyDosAttr),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_aclsettings.acl_settings_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, aclCreateOverSmb, states[0].Attributes["create_over_smb"])
					assert.Equal(t, aclReadOnlyDosAttr, states[0].Attributes["dos_attr"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + ACLSettingsUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "create_over_smb", "dis"+aclCreateOverSmb),
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "dos_attr", aclReadOnlyDosAttr+"_and_nfs"),
				),
			},
		},
	})
}

func TestAccAclSettingsResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + ACLSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "create_over_smb", aclCreateOverSmb),
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "dos_attr", aclReadOnlyDosAttr),
				),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_aclsettings.acl_settings_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetACLSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccAclSettingsResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + ACLSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "create_over_smb", aclCreateOverSmb),
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "dos_attr", aclReadOnlyDosAttr),
				),
			},
			// Update param read error
			{
				Config: ProviderConfig + ACLSettingsUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + ACLSettingsUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.UpdateACLSettings).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + ACLSettingsResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetACLSettings).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccAclSettingsResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + ACLSettingsResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.UpdateACLSettings).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccAclSettingsResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + ACLSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "create_over_smb", aclCreateOverSmb),
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "dos_attr", aclReadOnlyDosAttr),
				),
			},
			{
				ResourceName: "powerscale_aclsettings.acl_settings_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + ACLSettingsResourceConfig,
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
				Config:      ProviderConfig + ACLSettingsUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccAclSettingsResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + ACLSettingsResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Change the settings back to default
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + ACLSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "create_over_smb", aclCreateOverSmb),
					resource.TestCheckResourceAttr("powerscale_aclsettings.acl_settings_test", "dos_attr", aclReadOnlyDosAttr),
				),
			},
		},
	})
}

var aclCreateOverSmb = "allow"
var aclReadOnlyDosAttr = "deny_smb"

var ACLSettingsResourceConfig = fmt.Sprintf(`
resource "powerscale_aclsettings" "acl_settings_test" {
	create_over_smb = "%s"
	dos_attr = "%s"
}
`, aclCreateOverSmb, aclReadOnlyDosAttr)

var ACLSettingsUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_aclsettings" "acl_settings_test" {
	create_over_smb = "%s"
	dos_attr = "%s"
}
`, "dis"+aclCreateOverSmb, aclReadOnlyDosAttr+"_and_nfs")
