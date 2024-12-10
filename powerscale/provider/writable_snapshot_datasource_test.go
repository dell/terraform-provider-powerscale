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

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccWritableSnapshotDatasourceGetAll tests the writable snapshot datasource.
func TestAccWritableSnapshotDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing error
			{
				Config: ProviderConfig + `
				data "powerscale_writable_snapshot" "test" {
				}
				`,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetAllWritableSnapshots).Return(nil, fmt.Errorf("mock network error")).Build()
				},
				ExpectError: regexp.MustCompile("mock network error"),
			},
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + `
				data "powerscale_writable_snapshot" "test" {
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_writable_snapshot.test", "writable.#"),
				),
			},
		},
	})
}

// TestAccWritableSnapshotDatasourceID tests the writable snapshot datasource by ID.
func TestAccWritableSnapshotDatasourceID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing error
			{
				Config: ProviderConfig + `
				data "powerscale_writable_snapshot" "test" {
					filter {
						path = ""
					}
				}
				`,
				ExpectError: regexp.MustCompile(`.*string length must be between 4 and 4096*|.*id must start with*`),
			},
			{
				Config: ProviderConfig + `
				data "powerscale_writable_snapshot" "test" {
					filter {
						path = "invalid"
					}
				}
				`,
				ExpectError: regexp.MustCompile(`.*path must start with*`),
			},
			// Read testing
			{
				Config: ProviderConfig + writableSnapshotDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_writable_snapshot.test", "writable.#"),
					resource.TestCheckResourceAttr("data.powerscale_writable_snapshot.test", "writable.#", "1"),
				),
			},
			{
				Config: ProviderConfig + writableSnapshotDatasourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetWritableSnapshot).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config: ProviderConfig + writableSnapshotDatasourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.NewWritableSnapshotDataSource).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var writableSnapshotDatasourceConfig = writableSnapshotResourceConfig + `
data "powerscale_writable_snapshot" "preq" {
	depends_on = [powerscale_writable_snapshot.test]
}
data "powerscale_writable_snapshot" "test" {
	filter {
		path = data.powerscale_writable_snapshot.preq.writable[0].dst_path
	}
}
`
