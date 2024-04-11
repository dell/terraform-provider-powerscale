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
	"github.com/bytedance/mockey"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNfsExportSettingsImport(t *testing.T) {
	var nfsExportSettings = "powerscale_nfs_export_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsExportSettingsResourceConfig,
			},
			// Import testing
			{
				ResourceName:  nfsExportSettings,
				ImportState:   true,
				ImportStateId: "System",
				ExpectError:   nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(nfsExportSettings, "id")
					resource.TestCheckResourceAttrSet(nfsExportSettings, "symlinks")
					resource.TestCheckResourceAttrSet(nfsExportSettings, "map_non_root")
					resource.TestCheckResourceAttrSet(nfsExportSettings, "map_root")
					resource.TestCheckResourceAttrSet(nfsExportSettings, "map_failure")
					return nil
				},
			},
		},
	})
}

func TestAccNfsExportSettingsUpdate(t *testing.T) {
	var nfsExportSettings = "powerscale_nfs_export_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsExportSettingsResourceConfig,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + nfsExportSettingsUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(nfsExportSettings, "all_dirs", "true"),
					resource.TestCheckResourceAttr(nfsExportSettings, "case_insensitive", "true"),
					resource.TestCheckResourceAttr(nfsExportSettings, "case_preserving", "false"),
					resource.TestCheckResourceAttr(nfsExportSettings, "commit_asynchronous", "true"),
					resource.TestCheckResourceAttr(nfsExportSettings, "no_truncate", "true"),
					resource.TestCheckResourceAttr(nfsExportSettings, "write_datasync_action", "UNSTABLE"),
					resource.TestCheckResourceAttr(nfsExportSettings, "security_flavors.#", "2"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + nfsExportSettingsUpdateRevertResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(nfsExportSettings, "all_dirs", "false"),
					resource.TestCheckResourceAttr(nfsExportSettings, "case_insensitive", "false"),
					resource.TestCheckResourceAttr(nfsExportSettings, "case_preserving", "true"),
					resource.TestCheckResourceAttr(nfsExportSettings, "commit_asynchronous", "false"),
					resource.TestCheckResourceAttr(nfsExportSettings, "no_truncate", "false"),
					resource.TestCheckResourceAttr(nfsExportSettings, "write_datasync_action", "DATASYNC"),
					resource.TestCheckResourceAttr(nfsExportSettings, "security_flavors.#", "1"),
				),
			},
		},
	})
}

func TestAccNfsExportSettingsSnapshot(t *testing.T) {
	var nfsExportSettings = "powerscale_nfs_export_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsExportSettingsSnapshotResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(nfsExportSettings, "snapshot", "-"),
				),
			},
		},
	})
}

func TestAccNfsExportSettingsCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNfsExportSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsExportSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsExportSettingsUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsExportSettingsResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNfsExportSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsExportSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsExportSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsExportSettingsImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsExportSettingsResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetNfsExportSettingsByZone).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + nfsExportSettingsResourceConfig,
				ResourceName:      "powerscale_nfs_export_settings.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

var nfsExportSettingsResourceConfig = `
resource "powerscale_nfs_export_settings" "test" {

}
`

var nfsExportSettingsSnapshotResourceConfig = `
resource "powerscale_nfs_export_settings" "test" {
	snapshot = "-"
}
`

var nfsExportSettingsUpdateResourceConfig = `
resource "powerscale_nfs_export_settings" "test" {
	all_dirs = true
	case_insensitive = true
	case_preserving = false
	commit_asynchronous = true
	no_truncate = true
	write_datasync_action = "UNSTABLE"
	security_flavors = ["unix", "krb5"]
	map_non_root = {
		enabled = true
		primary_group = {}
		secondary_groups = []
		user: {id: "USER:nobody"}
	}
}
`

var nfsExportSettingsUpdateRevertResourceConfig = `
resource "powerscale_nfs_export_settings" "test" {
	all_dirs = false
	case_insensitive = false
	case_preserving = true
	commit_asynchronous = false
	no_truncate = false
	write_datasync_action = "DATASYNC"
	security_flavors = ["unix"]
	map_non_root = {
		enabled = false
		primary_group = {}
		secondary_groups = []
		user = {
			id = "USER:nobody"
		}
	}
}
`
