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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var createMockerLocal *mockey.Mocker
var setACLMockerLocal *mockey.Mocker
var metadataMocker *mockey.Mocker

func TestAccFileSystemResourceWithUIDChange(t *testing.T) {
	var fileSystemResourceName = "powerscale_filesystem.file_system_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + FileSystemResourceConfigWithUserChange,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fileSystemResourceName, "name", "tfaccDirTf"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "id", "ifs/tfaccDirTf"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "owner.name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "owner.id", "UID:20000"),
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
			// Update testing
			{
				Config: ProviderConfig + FileSystemResourceConfigWithUserChangeUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fileSystemResourceName, "name", "tfaccDirTf"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "id", "ifs/tfaccDirTf"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "owner.name", "tfaccUserCreation"),
					resource.TestCheckResourceAttr(fileSystemResourceName, "owner.id", "UID:20001"),
				),
			},
		},
	})
}

func TestAccFileSystemResource(t *testing.T) {
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
			// Import failure GetMetadata
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					metadataMocker = mockey.Mock((*powerscale.NamespaceApiService).GetDirectoryMetadataExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()

				},
				ResourceName: fileSystemResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
			// Import but GetACL failed and throw warning
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					setACLMockerLocal = mockey.Mock((*powerscale.NamespaceApiService).GetAclExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				ResourceName: fileSystemResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccFileSystemResourceUpdate(t *testing.T) {
	var fileSystemResourceName = "powerscale_filesystem.file_system_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
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
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccFileSystemResourceUpdateMetadataError(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl, then Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					setACLMockerLocal = mockey.Mock(helper.UpdateFileSystemOwnerAndGroup).Return(nil).Build()
					FunctionMocker = mockey.Mock(helper.UpdateFileSystemAccessControl).Return(nil).Build()
					metadataMocker = mockey.Mock((*powerscale.NamespaceApiService).GetDirectoryMetadataExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
				Destroy:     true,
			},
		},
	})
}

func TestAccFileSystemResourceUpdateGetAclError(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with ACL error and Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock(helper.GetDirectoryACL).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemResourceUpdConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl, then Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createMockerLocal = mockey.Mock(helper.UpdateFileSystemAccessControl).Return(nil).Build()
					FunctionMocker = mockey.Mock(helper.GetDirectoryACL).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
				Destroy:     true,
			},
		},
	})
}

func TestAccFileSystemResourceUpdateUserErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + FileSystemResourceUpdConfig,
			},
			//Update owner/group/accessControl failed, but throw warning
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock((*powerscale.NamespaceApiService).SetAclExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}
func TestAccFileSystemResourceUpdateFail(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
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

func TestAccFileSystemResourceGetMetaErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createMockerLocal = mockey.Mock(helper.ExecuteCreate).Return(nil, nil, nil).Build()
					FunctionMocker = mockey.Mock(helper.GetDirectoryMetadata).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccFileSystemResourceCreateFSErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					FunctionMocker = mockey.Mock((*powerscale.NamespaceApiService).CreateDirectoryExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccFileSystemResourceDeleteMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + FileSystemResourceConfig,
			},
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.UnPatch()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.UnPatch()
					}
					if metadataMocker != nil {
						metadataMocker.UnPatch()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
					createMockerLocal = mockey.Mock((*powerscale.NamespaceApiService).DeleteDirectoryExecute).When(func(r powerscale.ApiDeleteDirectoryRequest) bool {
						return createMockerLocal.Times() == 1
					}).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + FileSystemResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccFileSystemResourceReleaseMock(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if createMockerLocal != nil {
						createMockerLocal.Release()
					}
					if setACLMockerLocal != nil {
						setACLMockerLocal.Release()
					}
					if metadataMocker != nil {
						metadataMocker.Release()
					}
					if FunctionMocker != nil {
						FunctionMocker.UnPatch()
					}
				},
				Config: ProviderConfig + FileSystemResourceConfig,
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
var FileSystemResourceUpdConfig = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTfUpd"
  
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

var FileSystemUpdateResourceConfig = `
resource "powerscale_filesystem" "file_system_test" {
	# Default set to '/ifs'
	# directory_path         = "/ifs"
  
	# Required
	name = "tfaccDirTfUpd"
  
	recursive = true
	overwrite = true
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
	overwrite = true
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
	overwrite = true
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

var FileSystemResourceConfigWithUserChange = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"

	uid = 20000
	email = "test@dell.com"
	primary_group = "Administrators"
	roles = ["SystemAdmin"]
  }

resource "powerscale_filesystem" "file_system_test" {
	depends_on = [powerscale_user.test]

	name = "tfaccDirTf"
  
	recursive = true
	overwrite = true
	group = {
	  id   = "GID:0"
	}
	owner = {
	  name = "tfaccUserCreation",
	}
  }
`

var FileSystemResourceConfigWithUserChangeUpdate = `
resource "powerscale_user" "test" {
	name = "tfaccUserCreation"
  
	uid = 20001
	query_force = true
	enabled = true
	email = "newTest@dell.com"
	primary_group = "Administrators"
	roles = ["tfaccUserRole"]
  }

resource "powerscale_user_group" "testDep" {
	name = "tfaccUserGroupDatasource"
	gid = 10000
}

resource "powerscale_filesystem" "file_system_test" {
	depends_on = [powerscale_user.test, powerscale_user_group.testDep]

	name = "tfaccDirTf"
	recursive = true
	overwrite = true
	group = {
		name = "tfaccUserGroupDatasource",
	}
	owner = {
	  name = "tfaccUserCreation",
	}
  }
`
