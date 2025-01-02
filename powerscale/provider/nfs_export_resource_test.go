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
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNFSExport(t *testing.T) {
	var nfsExportResourceName = "powerscale_nfs_export.test_export"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NFSExportResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(nfsExportResourceName, "block_size"),
					resource.TestCheckResourceAttr(nfsExportResourceName, "paths.#", "1"),
					resource.TestCheckResourceAttr(nfsExportResourceName, "paths.0", "/ifs/tfacc_nfs_export"),
				),
			},
			// ImportState testing
			{
				ResourceName:      nfsExportResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"map_all", "force", "ignore_bad_auth", "ignore_bad_paths",
					"ignore_conflicts", "ignore_unresolvable_hosts"},
			},
			{
				Config: ProviderConfig + NFSExportUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(nfsExportResourceName, "block_size", "4096"),
					resource.TestCheckResourceAttr(nfsExportResourceName, "paths.#", "1"),
					resource.TestCheckResourceAttr(nfsExportResourceName, "paths.0", "/ifs/tfacc_nfs_export"),
					resource.TestCheckResourceAttr(nfsExportResourceName, "name_max_size", "127"),
				),
			},
		},
	})
}

func TestAccNFSExportErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NFSExportResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CreateNFSExport).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNFSExportErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NFSExportResourceConfig,
			},
			{
				Config: ProviderConfig + NFSExportUpdatedResourceConfig,
				PreConfig: func() {
					//FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.UpdateNFSExport).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config: ProviderConfig + NFSExportUpdatedResourceConfig2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetNFSExport).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				Config: ProviderConfig + NFSExportUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetNFSExport).Return(&powerscale.V2NfsExportsExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			{
				Config: ProviderConfig + NFSExportUpdatedResourceConfig2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNFSExportErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NFSExportResourceConfig,
			},
			// ImportState testing get none share
			{
				ResourceName: "powerscale_nfs_export.test_export",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetNFSExportByID).Return(&powerscale.V2NfsExportsExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_nfs_export.test_export",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetNFSExportByID).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_nfs_export.test_export",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing get none share
			{
				ResourceName: "powerscale_nfs_export.test_export",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetNFSExport).Return(&powerscale.V2NfsExportsExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_nfs_export.test_export",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetNFSExport).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_nfs_export.test_export",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNFSExportErrorDelete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + NFSExportResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.DeleteNFSExport).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) bool {
							return FunctionMocker.Times() == 1
						})
				},
				Config:      ProviderConfig + NFSExportResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var FileSystemResourceConfigCommon2 = `
resource "powerscale_filesystem" "file_system_test" {
	directory_path         = "/ifs"	
	name = "tfacc_nfs_export"
	
	  recursive = true
	  overwrite = true
	  group = {
		id   = "GID:0"
		name = "wheel"
		type = "group"
	  }
	  owner = {
		  id   = "UID:0",
		 name = "root",
		 type = "user"
	   }
	}
`
var NFSExportResourceConfig = FileSystemResourceConfigCommon2 + `
resource "powerscale_nfs_export" "test_export" {
	depends_on = [powerscale_filesystem.file_system_test]
	paths = ["/ifs/tfacc_nfs_export"]
	force = true
	map_all = {
		enabled = true
		primary_group = {
		  id = "GROUP:Users"
		}
		user = {
		  type = "user"
		  name = "Guest"	
		}
		secondary_groups = [
			{
			  id = "GROUP:Users"
			}
		]
	}
    ignore_bad_auth = true
    ignore_bad_paths= true
    ignore_conflicts = true
    ignore_unresolvable_hosts = true
    zone = "System"
}
`

var NFSExportUpdatedResourceConfig = FileSystemResourceConfigCommon2 + `
resource "powerscale_nfs_export" "test_export" {
	depends_on = [powerscale_filesystem.file_system_test]
	paths = ["/ifs/tfacc_nfs_export"]
	force = true
	map_all = {
		enabled = true
		primary_group = {
		  id = "GROUP:Users"
		}
		user = {
		  type = "user"
		  name = "Guest"	
		}
		secondary_groups = [
			{
			  id = "GROUP:Users"
			}
		]
	}
    ignore_bad_auth = true
    ignore_bad_paths= true
    ignore_conflicts = true
    ignore_unresolvable_hosts = true
    zone = "System"
	block_size = 4096
	name_max_size = 127
}
`

var NFSExportUpdatedResourceConfig2 = FileSystemResourceConfigCommon2 + `
resource "powerscale_nfs_export" "test_export" {
	paths = ["/ifs/tfacc_nfs_export"]
	force = true
	map_all = {
		enabled = true
		primary_group = {
		  id = "GROUP:Users"
		}
		user = {
		  type = "user"
		  name = "Guest"	
		}
		secondary_groups = [
			{
			  id = "GROUP:Users"
			}
		]
	}
    ignore_bad_auth = true
    ignore_bad_paths= true
    ignore_conflicts = true
    ignore_unresolvable_hosts = true
    zone = "System"
	block_size = 4096
	name_max_size = 255
}
`
