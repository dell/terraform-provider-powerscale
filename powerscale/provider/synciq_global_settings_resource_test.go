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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccSyncIQGlobalSettingsResource - Tests the creation of a syncIQ global settings resource.
func TestAccSyncIQGlobalSettingsResource(t *testing.T) {
	var syncIQGlobalSettingsResource = "powerscale_synciq_global_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(syncIQGlobalSettingsResource, "preferred_rpo_alert", "3"),
				),
			},
		},
	})
}

// TestAccSyncIQGlobalSettingsResource_Update - Tests the update of a synciq global settings resource along with error mocking.
func TestAccSyncIQGlobalSettingsResource_Update(t *testing.T) {
	var syncIQGlobalSettingsResource = "powerscale_synciq_global_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(syncIQGlobalSettingsResource, "preferred_rpo_alert", "3"),
				),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSyncIQGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSyncIQGlobalSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSyncIQGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(syncIQGlobalSettingsResource, "preferred_rpo_alert", "5"),
				),
			},
			{
				Config: ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(syncIQGlobalSettingsResource, "preferred_rpo_alert", "3"),
				),
			},
		},
	})
}

// TestAccSyncIQGlobalSettingsResourceCreate - Tests the mock errors during the create operation of the SyncIQ global settings resource.
func TestAccSyncIQGlobalSettingsResourceCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSyncIQGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSyncIQGlobalSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

// TestAccSyncIQGlobalSettingsResourceUpdateMockErr - Tests the mock errors during the update operation of the syncIQ Global Settings resource.
func TestAccSyncIQGlobalSettingsResourceUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSyncIQGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateSyncIQGlobalSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

// TestAccSyncIQGlobalSettingsResourceImportMockErr - Tests the mock errors during the import of the SyncIQ Global Settings resource.
func TestAccSyncIQGlobalSettingsResourceImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + SyncIQGlobalSettingsResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetSyncIQGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ResourceName:      "powerscale_synciq_global_settings.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:       ProviderConfig + SyncIQGlobalSettingsResourceConfig,
				ResourceName: "powerscale_synciq_global_settings.test",
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var SyncIQGlobalSettingsResourceConfig = `
resource "powerscale_synciq_global_settings" "test" {
	preferred_rpo_alert = 3
	report_email = ["mail1@dell.com", "mail2@dell.com"]
   restrict_target_network = true
   rpo_alerts           = true
   service              = "paused"
}
`

var SyncIQGlobalSettingsResourceConfigUpdate = `
resource "powerscale_synciq_global_settings" "test" {
	preferred_rpo_alert = 5
	report_email = ["mail2@dell.com", "mail3@dell.com"]
   restrict_target_network = true
   rpo_alerts           = true
   service              = "paused"
}
`
