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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestFileSystemResource(t *testing.T) {
	var fileSystemResourceName = "powerscale_filesystem.file_system_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + FileSystemResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fileSystemResourceName, "name", "tfaccDirTf"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "id", "ifs/tfaccDirTf"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "owner.name", "root"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "directory_path", "/ifs"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "group.name", "wheel"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "mode", "0700"),
				),
			},
			// ImportState testing
			{
				ResourceName: fileSystemResourceName,
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "tfaccDirTf", states[0].Attributes["name"])
					assert.Equal(t, "0700", states[0].Attributes["mode"])
					return nil
				},
			},
			// Update to error state
			{
				Config:      ProviderConfig + FileSystemUpdateResourceConfigError,
				ExpectError: regexp.MustCompile(".*Error updating the File system Resource*."),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestFileSystemResourceUpdate(t *testing.T) {
	var fileSystemResourceName = "powerscale_filesystem.file_system_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl, then Read testing
			{
				Config: ProviderConfig + FileSystemUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fileSystemResourceName, "name", "tfaccDirTfUpd"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "id", "ifs/tfaccDirTfUpd"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "owner.name", "Guest"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "directory_path", "/ifs"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "group.name", "wheel"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "mode", "0770"),
				),
			},
		},
	})
}

func TestFileSystemResourceUpdateUserErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl, then Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.UpdateFileSystem).Return(fmt.Errorf("Error updating user")).Build()
				},
				Config:      ProviderConfig + FileSystemUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error updating user*.`),
			},
		},
	})
}
func TestFileSystemResourceUpdateACLErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl, then Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.UpdateFileSystem).Return(fmt.Errorf("Error updating acl")).Build()
				},
				Config:      ProviderConfig + FileSystemUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error updating acl*.`),
			},
		},
	})
}
func TestFileSystemResourceUpdateFail(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl, then Read testing
			{
				Config:      ProviderConfig + FileSystemUpdateResourceConfigErr,
				ExpectError: regexp.MustCompile(`.*Renaming Directory is not supported*.`),
			},
		},
	})
}

var FileSystemResourceConfig = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTf"
  
	recursive = true
	overwrite = false
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
var FileSystemResourceUpdConfig = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTfUpd"
  
	recursive = true
	overwrite = false
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

var FileSystemUpdateResourceConfig = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTfUpd"
  
	recursive = true
	overwrite = false
	group = {
	  id   = "GID:0"
	  name = "wheel"
	  type = "group"
	}
	owner = {
	   id   = "UID:1501",
	  name = "Guest",
	  type = "user"
	}
	access_control = "0770"
  }
`
var FileSystemUpdateResourceConfigErr = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTfUpdErr"
  
	recursive = true
	overwrite = false
	group = {
	  id   = "GID:0"
	  name = "wheel"
	  type = "group"
	}
	owner = {
	   id   = "UID:1501",
	  name = "Guest",
	  type = "user"
	}
	access_control = "0770"
  }
`

var FileSystemUpdateResourceConfigError = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTf"
  
	recursive = true
	overwrite = false
	group = {
	  id   = "GID:0"
	  name = "wheel"
	  type = "group"
	}
	owner = {
		id   = "UID:1501",
	   name = "Guest",
	   type = "user"
	 }
	access_control = "private_read"
  }
`
