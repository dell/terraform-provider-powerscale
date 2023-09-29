/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var snapMocker *Mocker
var createSnapMocker *Mocker

func TestAccSnapshotResourceA(t *testing.T) {
	var snapshotResourceName = "powerscale_snapshot.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SnapshotResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotResourceName, "path", "/ifs/tfacc_file_system_test"),
					resource.TestCheckResourceAttr(snapshotResourceName, "expires", "0"),
					resource.TestCheckResourceAttr(snapshotResourceName, "set_expires", "Never"),
				),
			},
			// Update name, path and auth providers, then Read testing
			{
				Config: ProviderConfig + SnapshotResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotResourceName, "path", "/ifs/tfacc_file_system_test"),
					resource.TestCheckResourceAttr(snapshotResourceName, "set_expires", "1 Day"),
					resource.TestCheckResourceAttr(snapshotResourceName, "name", "tfacc_snapshot_1"),
				),
			},
			// Update to error state
			{
				Config:      ProviderConfig + SnapshotResourceConfigUpdateErrorPathModify,
				ExpectError: regexp.MustCompile(".*Path variable is not able to be updated after creation*."),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSnapshotResourceCreateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if snapMocker != nil {
						snapMocker.UnPatch()
					}
					if createSnapMocker != nil {
						createSnapMocker.UnPatch()
					}
					createSnapMocker = Mock(helper.CreateSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating snapshot*.`),
			},
		},
	})
}

func TestAccSnapshotResourceMapperError(t *testing.T) {
	create := powerscale.Createv1SnapshotSnapshotResponse{
		Id: 3,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if snapMocker != nil {
						snapMocker.UnPatch()
					}
					if createSnapMocker != nil {
						createSnapMocker.UnPatch()
					}
					createSnapMocker = Mock(helper.CreateSnapshot).Return(create, nil).Build()
					snapMocker = Mock(helper.SnapshotResourceDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotResourceConfig,
				ExpectError: regexp.MustCompile(`.*Could not create snapshot*.`),
			},
		},
	})
}

func TestAccSnapshotResourceReadAndUpdateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				PreConfig: func() {
					if snapMocker != nil {
						snapMocker.UnPatch()
					}
					if createSnapMocker != nil {
						createSnapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + SnapshotResourceConfig,
			},
			// Read Error
			{
				PreConfig: func() {
					snapMocker = Mock(helper.GetSpecificSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting the list of snapshots*.`),
			},
			// Update Error
			{
				PreConfig: func() {
					if snapMocker != nil {
						snapMocker.UnPatch()
					}
					if createSnapMocker != nil {
						createSnapMocker.UnPatch()
					}
					snapMocker = Mock(helper.ModifySnapshot).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*Error editing snapshot*.`),
			},
			// Update Error
			// Delete testing automatically occurs in TestCase
		},
	})
}

var SnapshotResourceConfig = `
resource "powerscale_snapshot" "test" {
  # Required path to the filesystem to which the snapshot will be taken of
  path = "/ifs/tfacc_file_system_test"
}
`

var SnapshotResourceConfigUpdate = `
resource "powerscale_snapshot" "test" {
  path = "/ifs/tfacc_file_system_test"
  name = "tfacc_snapshot_1"
  set_expires = "1 Day"
}
`

var SnapshotResourceConfigUpdateErrorPathModify = `
resource "powerscale_snapshot" "test" {
  path = "/ifs/tfacc_file_system_test/shouldnt/update"
  name = "tfacc_snapshot_1"
  set_expires = "1 Day"
}
`
