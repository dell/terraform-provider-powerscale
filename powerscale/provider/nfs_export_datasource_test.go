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
					FunctionMocker = mockey.Mock(helper.ListNFSExports).Return(exports, nil).Build()
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

func TestAccNfsExportDatasourceGetWithQueryParam(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + NfsExportDatasourceGetWithQueryParam,
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
					FunctionMocker = mockey.Mock(helper.ListNFSExports).Return(nil, fmt.Errorf("mock error")).Build()
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
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
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

var NfsExportDatasourceGetWithQueryParam = FileSystemResourceConfigCommon2 + `
resource "powerscale_nfs_export" "test_export" {
	depends_on = [powerscale_filesystem.file_system_test]
	paths = ["/ifs/tfacc_nfs_export"]
}

data "powerscale_nfs_export" "export_datasource_test" {
	filter {
        check = true
        dir   = "ASC"
		limit = 10
		scope = "effective"
        sort  = "id"
	}
  	depends_on = [
    	powerscale_nfs_export.test_export
  	]
}
`

var NfsExportDatasourceGetAllConfig = FileSystemResourceConfigCommon2 + `
resource "powerscale_nfs_export" "test_export" {
	depends_on = [powerscale_filesystem.file_system_test]
	paths = ["/ifs/tfacc_nfs_export"]
}

data "powerscale_nfs_export" "export_datasource_test" {
  	depends_on = [
    	powerscale_nfs_export.test_export
  	]
}
`
