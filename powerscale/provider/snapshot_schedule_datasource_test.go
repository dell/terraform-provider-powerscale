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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnapshotScheduleDataSource(t *testing.T) {
	var snapshotTerraformName = "data.powerscale_snapshot_schedule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + SnapshotScheduleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotTerraformName, "schedules.0.duration", "604800"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "schedules.0.name", "Snapshot schedule 370395356"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "schedules.0.alias", "tfacc_Snapshot schedule 370395356_Alias"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "schedules.0.path", "/ifs/tfacc_test_dirNew"),
				),
			},
		},
	})
}

func TestAccSnapshotScheduleDataSourceAll(t *testing.T) {
	var ssTerraformName = "data.powerscale_snapshot_schedule.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + SnapshotScheduleAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(ssTerraformName, "schedules.#"),
				),
			},
		},
	})
}

func TestAccSnapshotScheduleDataSourceOtherFilters(t *testing.T) {
	var filterTerraformName = "data.powerscale_snapshot_schedule.filterTest"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + SnapshotScheduleOtherFiltersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(filterTerraformName, "schedules.#", "1"),
					resource.TestCheckResourceAttr(filterTerraformName, "schedules.0.path", "/ifs"),
				),
			},
		},
	})
}

func TestAccSnapshotScheduleDataSourceGetErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ListSnapshotSchedules).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotScheduleDataSourceGetErrCopyAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotScheduleDataSourceGetErrCopy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var SnapshotScheduleDataSourceConfig = `
data "powerscale_snapshot_schedule" "test" {
  filter {
    names = ["Snapshot schedule 370395356"]
  }
}
output "powerscale_snapshot_schedule" {
	value = data.powerscale_snapshot_schedule.test
}
`

var SnapshotScheduleAllDataSourceConfig = `
data "powerscale_snapshot_schedule" "all" {
}
`

var SnapshotScheduleOtherFiltersDataSourceConfig = `
data "powerscale_snapshot_schedule" "filterTest" {
	filter {	   
		dir = "ASC"
		limit = 1
		sort = "name"
	  }
}
`
