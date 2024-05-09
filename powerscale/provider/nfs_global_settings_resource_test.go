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

func TestAccNfsGlobalSettingsImport(t *testing.T) {
	var nfsGlobalSettings = "powerscale_nfs_global_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsGlobalSettingsResourceConfig,
			},
			// Import testing
			{
				ResourceName: nfsGlobalSettings,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "id")
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "nfsv3_enabled")
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "nfsv3_rdma_enabled")
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "nfsv4_enabled")
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "rquota_enabled")
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "service")
					return nil
				},
			},
		},
	})
}

func TestAccNfsGlobalSettingsUpdate(t *testing.T) {
	var nfsGlobalSettings = "powerscale_nfs_global_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsGlobalSettingsResourceConfig,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + nfsGlobalSettingsUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(nfsGlobalSettings, "nfsv3_enabled", "false"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "nfsv3_rdma_enabled", "true"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "nfsv4_enabled", "true"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "rpc_maxthreads", "32"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "rpc_minthreads", "32"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "rquota_enabled", "true"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "service", "false"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + nfsGlobalSettingsUpdateRevertResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(nfsGlobalSettings, "nfsv3_enabled", "true"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "nfsv3_rdma_enabled", "false"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "nfsv4_enabled", "false"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "rpc_maxthreads", "16"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "rpc_minthreads", "16"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "rquota_enabled", "false"),
					resource.TestCheckResourceAttr(nfsGlobalSettings, "service", "true"),
				),
			},
		},
	})
}

func TestAccNfsGlobalSettingsCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNfsGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsGlobalSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsGlobalSettingsUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsGlobalSettingsResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNfsGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateNfsGlobalSettings).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + nfsGlobalSettingsUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccNfsGlobalSettingsImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + nfsGlobalSettingsResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetNfsGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + nfsGlobalSettingsResourceConfig,
				ResourceName:      "powerscale_nfs_global_settings.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

var nfsGlobalSettingsResourceConfig = `
resource "powerscale_nfs_global_settings" "test" {

}
`

var nfsGlobalSettingsUpdateResourceConfig = `
resource "powerscale_nfs_global_settings" "test" {
	nfsv3_enabled = false
	nfsv3_rdma_enabled = true
	nfsv4_enabled = true
	rpc_maxthreads = 32
	rpc_minthreads = 32
	rquota_enabled = true
	service = false
}
`

var nfsGlobalSettingsUpdateRevertResourceConfig = `
resource "powerscale_nfs_global_settings" "test" {
	nfsv3_enabled = true
	nfsv3_rdma_enabled = false
	nfsv4_enabled = false
	rpc_maxthreads = 16
	rpc_minthreads = 16
	rquota_enabled = false
	service = true
}
`
