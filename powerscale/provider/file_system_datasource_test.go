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

func TestAccFileSystemDataSource(t *testing.T) {
	var fsTerraform = "data.powerscale_filesystem.system"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter read testing
			{
				Config: ProviderConfig + FileSystemDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fsTerraform, "directory_path", "/ifs/tfacc_file_system_test"),
					resource.TestCheckResourceAttrSet(fsTerraform, "file_systems_details.file_system_attributes.#"),
					resource.TestCheckResourceAttrSet(fsTerraform, "file_systems_details.file_system_namespace_acl.acl.0.access_rights.#"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.authoritative", "mode"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.group.id", "GID:0"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.group.name", "wheel"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.group.type", "group"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.owner.id", "UID:0"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.owner.name", "root"),
					resource.TestCheckResourceAttr(fsTerraform, "file_systems_details.file_system_namespace_acl.owner.type", "user"),
				),
			},
		},
	})
}

func TestAccFileSystemDataSourceFilterDefault(t *testing.T) {
	var fsTerraform = "data.powerscale_filesystem.system"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + FileSystemDataSourceDefaultConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fsTerraform, "directory_path", "/ifs"),
				),
			},
		},
	})
}

func TestAccFileSystemDataSourceGetAclErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetDirectoryACL).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccFileSystemDataSourceGetQuotaErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetDirectoryQuota).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccFileSystemDataSourceGetSnapErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetDirectorySnapshots).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccFileSystemDataSourceGetMetaErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetDirectoryMetadata).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccFileSystemDataSourceReleaseMock(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + FileSystemDataSourceConfig,
			},
		},
	})
}

var FileSystemResourceConfigCommon = `
resource "powerscale_filesystem" "file_system_test" {
	directory_path         = "/ifs"	
	name = "tfacc_file_system_test"	
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

var FileSystemDataSourceConfig = FileSystemResourceConfigCommon + `
data "powerscale_filesystem" "system" {
	# Required parameter, path of the directory filesystem you would like to create a datasource out of 

	depends_on = [powerscale_filesystem.file_system_test]
	directory_path = "/ifs/tfacc_file_system_test"
  }
`

var FileSystemDataSourceDefaultConfig = `
data "powerscale_filesystem" "system" {
	# No Directory_path should use default
  }
`
