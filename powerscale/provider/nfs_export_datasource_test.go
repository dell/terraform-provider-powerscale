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
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

func TestAccNfsExportDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					id := int64(1)
					exports := &[]powerscale.V2NfsExportExtended{{
						Id:    &id,
						Paths: []string{"/ifs/primary"},
					}}
					FunctionMocker = Mock(helper.ListNFSExports).Return(exports, nil).Build()
				},
				Config: ProviderConfig + NfsExportDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_nfs_export.export_datasource_test", "nfs_exports.#", "1"),
					resource.TestCheckResourceAttr("data.powerscale_nfs_export.export_datasource_test", "nfs_exports.0.paths.0", "/ifs/primary"),
				),
			},
		},
	})
}

func TestAccNfsExportDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + NfsExportDatasourceGetAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_nfs_export.export_datasource_test", "nfs_exports.#"),
				),
			},
		},
	})
}

func TestAccNfsExportDatasourceErrorList(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ListNFSExports).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NfsExportDatasourceGetAllConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNfsExportDatasourceErrorCopy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NfsExportDatasourceGetAllConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var NfsExportDatasourceConfig = `
data "powerscale_nfs_export" "export_datasource_test" {
	filter {
		ids = [1]
		paths = ["/ifs/primary", "/ifs/secondary"]
	}
}
`

var NfsExportDatasourceGetAllConfig = `
data "powerscale_nfs_export" "export_datasource_test" {}
`
