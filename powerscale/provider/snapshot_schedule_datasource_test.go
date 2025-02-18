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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnapshotScheduleDataSourceA(t *testing.T) {
	var snapshotTerraformName = "data.powerscale_snapshot_schedule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + SnapshotScheduleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(snapshotTerraformName, "schedules.0.name", "tfacc_snap_schedule_test"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "schedules.0.alias", "test_alias"),
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
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ListSnapshotSchedules).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
	FunctionMocker.UnPatch()

}

func TestAccSnapshotScheduleDataSourceGetErrCopyAll(t *testing.T) {
	FunctionMockerList := mockey.Mock(helper.ListSnapshotSchedules).Return([]powerscale.V1SnapshotScheduleExtended{
		{}, {},
	}, nil).Build()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleAllDataSourceConfig2,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
	FunctionMockerList.UnPatch()
}

func NewValue[v any](in v) *v {
	return &in
}

func TestAccSnapshotScheduleDataSourceGetErrCopy(t *testing.T) {
	FunctionMockerList := mockey.Mock(helper.ListSnapshotSchedules).Return([]powerscale.V1SnapshotScheduleExtended{
		{
			Name: NewValue("tfacc_snap_schedule_test"),
		},
		{
			Name: NewValue("tfacc_snap_schedule_test1"),
		},
	}, nil).Build()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotScheduleDataSourceConfig2,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
	FunctionMockerList.UnPatch()
}

var SnapshotScheduleCommonConfig = FileSystemResourceConfigCommon + `
resource "powerscale_snapshot_schedule" "test" {
	# Required name of snapshot schedule
	depends_on = [powerscale_filesystem.file_system_test] 
	name = "tfacc_snap_schedule_test"
	alias = "test_alias"
	path = "/ifs/tfacc_file_system_test"
	retention_time = "3 Hour(s)"
}
`

var SnapshotScheduleDataSourceConfig = SnapshotScheduleCommonConfig + `
data "powerscale_snapshot_schedule" "test" {
  depends_on = [powerscale_snapshot_schedule.test]
  filter {
    names = ["tfacc_snap_schedule_test"]
  }
}
output "powerscale_snapshot_schedule" {
	value = data.powerscale_snapshot_schedule.test
}
`

var SnapshotScheduleAllDataSourceConfig = SnapshotScheduleCommonConfig + `
data "powerscale_snapshot_schedule" "all" {
depends_on = [powerscale_snapshot_schedule.test]
}
`

var SnapshotScheduleOtherFiltersDataSourceConfig = SnapshotScheduleCommonConfig + `
data "powerscale_snapshot_schedule" "filterTest" {
depends_on = [powerscale_snapshot_schedule.test]
	filter {	   
		dir = "ASC"
		limit = 1
		sort = "name"
	  }
}
`
var SnapshotScheduleDataSourceConfig2 = `
data "powerscale_snapshot_schedule" "test" {
	  filter {
		names = ["tfacc_snap_schedule_test"]
	  }
	}
`
var SnapshotScheduleAllDataSourceConfig2 = `
data "powerscale_snapshot_schedule" "all" {
}
`
