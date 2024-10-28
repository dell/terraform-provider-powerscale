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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSnapshotDataSource(t *testing.T) {
	var snapshotTerraformName = "data.powerscale_snapshot.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + SnapshotDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(snapshotTerraformName, "snapshots_details.0.expires"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "snapshots_details.0.has_locks", "false"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "snapshots_details.0.pct_reserve", "0"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "snapshots_details.0.path", "/ifs/tfacc_file_system_test"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "snapshots_details.0.shadow_bytes", "0"),
					resource.TestCheckResourceAttr(snapshotTerraformName, "snapshots_details.0.state", "active"),
				),
			},
		},
	})
}

func TestAccSnapshotDataSourceAll(t *testing.T) {
	var azTerraformName = "data.powerscale_snapshot.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + SnapshotAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(azTerraformName, "snapshots_details.#"),
				),
			},
		},
	})
}

func TestAccSnapshotDataSourceGetErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetAllSnapshots).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccSnapshotDataSourceMapErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.SnapshotDetailMapper).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SnapshotAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var SnapshotDataSourceConfig = `
data "powerscale_snapshot" "test" {
  filter {
    path = "/ifs/tfacc_file_system_test"
	name = "Snapshot: 2024Apr15, 4:34 PM"
	state = "active"
  }
}
output "powerscale_snapshot" {
	value = data.powerscale_snapshot.test
}
`

var SnapshotAllDataSourceConfig = `
data "powerscale_snapshot" "all" {
}
`
