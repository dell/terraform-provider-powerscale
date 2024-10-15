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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccwritableSnapshotResourceImport - Tests the import of the writable snapshot resource.
func TestAccwritableSnapshotResourceImport(t *testing.T) {
	var writableSnapshotResourceName = "powerscale_writable_snapshot.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Success
			{
				Config: ProviderConfig + writableSnapshotResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(writableSnapshotResourceName, "dst_path", "/ifs/abcd"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:       ProviderConfig + writableSnapshotResourceConfig,
				ResourceName: writableSnapshotResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
			// Import testing
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:       ProviderConfig + writableSnapshotResourceConfig,
				ResourceName: writableSnapshotResourceName,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(writableSnapshotResourceName, "dst_path")
					resource.TestCheckResourceAttrSet(writableSnapshotResourceName, "src_path")
					resource.TestCheckResourceAttrSet(writableSnapshotResourceName, "snap_id")
					return nil
				},
			},
		},
	})
}

// TestAccwritableSnapshotResourceImportMockErr - Tests the mock errors during the import of the writable snapshot resource.
func TestAccwritableSnapshotResourceImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + writableSnapshotResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + writableSnapshotResourceConfig,
				ResourceName:      "powerscale_writable_snapshot.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccwritableSnapshotResource - Tests the creation of a writable snapshot resource.
func TestAccwritableSnapshotResource(t *testing.T) {
	var writableSnapshotResourceName = "powerscale_writable_snapshot.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + writableSnapshotResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(writableSnapshotResourceName, "dst_path", "/ifs/abcd"),
				),
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.DeleteWritableSnapshot).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + writableSnapshotResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// mock refresh error
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				RefreshState: true,
				ExpectError:  regexp.MustCompile("mock error"),
			},
		},
	})
}

// TestAccwritableSnapshotResource_Update - Tests the update of a writable snapshot resource along with error mocking.
func TestAccwritableSnapshotResource_Update(t *testing.T) {
	var writableSnapshotResourceName = "powerscale_writable_snapshot.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + writableSnapshotResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(writableSnapshotResourceName, "dst_path", "/ifs/abcd"),
				),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + writableSnapshotResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},

			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, path string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + writableSnapshotResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + writableSnapshotResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(writableSnapshotResourceName, "dst_path", "/ifs/abcd1"),
				),
			},
			{
				Config: ProviderConfig + writableSnapshotResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(writableSnapshotResourceName, "dst_path", "/ifs/abcd"),
				),
			},
		},
	})
}

// TestAccwritableSnapshotResourceCreateMockErr - Tests the mock errors during the create operation of the writable snapshot resource.
func TestAccwritableSnapshotResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + writableSnapshotResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + writableSnapshotResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.DeleteWritableSnapshot).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + writableSnapshotResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + writableSnapshotResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var snapshotPrereqConfig = `
resource "powerscale_snapshot" "snap" {
	path = powerscale_filesystem.file_system_test.full_path
	name = "snap_restore_snap"
  }
`
var writableSnapshotResourceConfig = FileSystemResourceConfig + snapshotPrereqConfig + `
resource "powerscale_writable_snapshot" "test" {
	snap_id = powerscale_snapshot.snap.id
	dst_path = "/ifs/abcd"
}
`

var writableSnapshotResourceConfigUpdate = FileSystemResourceConfig + snapshotPrereqConfig + `
resource "powerscale_writable_snapshot" "test" {
	snap_id = powerscale_snapshot.snap.id
	dst_path = "/ifs/abcd1"
}
`
