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
	"time"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnapshotRestoreResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + snapRevertResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("powerscale_snapshot_restore.snap_restore", "snaprevert_params.snapshot_id", "powerscale_snapshot.snap", "id"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(30 * time.Second)
				},
				Config: ProviderConfig + snapRevertResourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("powerscale_snapshot_restore.snap_restore", "snaprevert_params.snapshot_id", "powerscale_snapshot.snap1", "id"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(30 * time.Second)
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyDirectory).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + copyDirectoryConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFile).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + copyFileConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CloneFile).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + cloneFileConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var syncIQPre = `
resource "powerscale_synciq_global_settings" "test" {
   service              = "on"
}
`

var snapshotPre = `
resource "powerscale_snapshot" "snap" {
	path = powerscale_filesystem.file_system_test.full_path
	name = "snap_restore_snap"
  }
`

var snapshotPre1 = `
resource "powerscale_snapshot" "snap1" {
	path = powerscale_filesystem.file_system_test.full_path
	name = "snap_restore_snap1"
  }
`

var snapRevertResourceConfig = syncIQPre + FileSystemResourceConfig + snapshotPre + `
resource "powerscale_snapshot_restore" "snap_restore" {
	snaprevert_params = {
    	snapshot_id = powerscale_snapshot.snap.id
  	}
}
`

var snapRevertResourceConfigUpdate = FileSystemResourceConfig + snapshotPre + snapshotPre1 + `
resource "powerscale_snapshot_restore" "snap_restore" {
	snaprevert_params = {
    	snapshot_id = powerscale_snapshot.snap1.id
  	}
}
`

var copyDirectoryConfig = `
resource "powerscale_snapshot_restore" "snap_restore_copy" {
	copy_params = {
		directory = {
		  source      = "/namespace/ifs/.snapshot/terraform_snap/terraform_test"
		  destination = "ifs/dest"                                 
		}
	}
}
`

var copyFileConfig = `
resource "powerscale_snapshot_restore" "snap_restore_copy" {
	copy_params = {
		file = {
		  source      = "/namespace/ifs/.snapshot/terraform_snap/terraform_test/test.txt"
		  destination = "ifs/dest/test.txt"
		  overwrite   = true                              
		}
	}
}
`

var cloneFileConfig = FileSystemResourceConfig + snapshotPre + `
resource "powerscale_snapshot_restore" "snap_restore_copy" {
	clone_params = {
		  source      = "/namespace/ifs/.snapshot/terraform_snap/terraform_test/test.txt"
		  destination = "ifs/dest/test.txt"
		  snapshot_id = powerscale_snapshot.snap.id
		  overwrite   = true                              
	}
}
`
