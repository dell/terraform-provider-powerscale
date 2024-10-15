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
	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

func TestAccSmbShareResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmbShareResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "name", shareName),
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "ca_timeout", "120"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_smb_share.share_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, shareName, states[0].Attributes["id"])
					assert.Equal(t, "120", states[0].Attributes["ca_timeout"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + SmbShareNameUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "id", shareName+"_update"),
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "allow_delete_readonly", "true"),
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "ca_timeout", "30"),
				),
			},
		},
	})
}

func TestAccSmbShareResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmbShareResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "name", shareName),
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "ca_timeout", "120"),
				),
			},
			// ImportState testing get none share
			{
				ResourceName: "powerscale_smb_share.share_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSmbShare).Return(&powerscale.V7SmbSharesExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_smb_share.share_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSmbShare).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing get none share
			{
				ResourceName: "powerscale_smb_share.share_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSmbShare).Return(&powerscale.V7SmbSharesExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone *string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_smb_share.share_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSmbShare).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone *string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSmbShareResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmbShareResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "name", shareName),
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "ca_timeout", "120"),
				),
			},
			// Update get error
			{
				Config: ProviderConfig + SmbShareUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateSmbShare).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update zone get error
			{
				Config: ProviderConfig + SmbShareUpdatedZoneResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.QueryZoneNameByID).Return("", fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none share
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSmbShare).Return(&powerscale.V7SmbSharesExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone *string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + SmbShareUpdatedResourceConfig2,
				ExpectError: regexp.MustCompile(".not found"),
			},
			// Update get error
			{
				Config: ProviderConfig + SmbShareUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSmbShare).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone *string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSmbShareResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + SmbShareInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccSmbShareResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SmbShareResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "name", shareName),
					resource.TestCheckResourceAttr("powerscale_smb_share.share_test", "ca_timeout", "120"),
				),
			},
			{
				ResourceName: "powerscale_smb_share.share_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSmbShareResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var shareName = "tfacc_test_smb_share"

var FileSystemResourceConfigCommon6 = fmt.Sprintf(`
resource "powerscale_filesystem" "file_system_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	directory_path         = "/ifs"	
	name = "%s"	
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
`, shareName)

var SmbShareResourceConfig = FileSystemResourceConfigCommon6 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	auto_create_directory = true
	name = "%s"
	path = "/ifs/%s"
	permissions = [
		{
			permission = "full"
			permission_type = "allow"
			trustee = {
				id = "SID:S-1-1-0",
				name = "Everyone",
				type = "wellknown"
			}
		}
	]
	ca_timeout = 120
	zone = "System"
}
`, shareName, shareName)

var SmbShareInvalidResourceConfig = FileSystemResourceConfigCommon6 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	auto_create_directory = true
	name = "%s"
	path = "/ifs/%s"
	zone = "System"
	permissions = [
		{
			permission = "full"
			permission_type = "allow"
			trustee = {
				id = "SID:S-1-1-0",
				name = "invalid",
				type = "invalid"
			}
		}
	]
	ca_timeout = 120
}
`, shareName, shareName)

var SmbShareUpdatedResourceConfig = FileSystemResourceConfigCommon6 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	auto_create_directory = true
	name = "%s"
	path = "/ifs/%s"
	permissions = [
		{
			permission = "full"
			permission_type = "allow"
			trustee = {
				id = "SID:S-1-1-0",
				name = "Everyone",
				type = "wellknown"
			}
		}
	]
	allow_delete_readonly = true
	ca_timeout = 30
	zone = "System"
}
`, shareName, shareName)

var SmbShareUpdatedResourceConfig2 = FileSystemResourceConfigCommon6 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	auto_create_directory = true
	name = "%s"
	path = "/ifs/%s"
	permissions = [
		{
			permission = "full"
			permission_type = "allow"
			trustee = {
				id = "SID:S-1-1-0",
				name = "Everyone",
				type = "wellknown"
			}
		}
	]
	allow_delete_readonly = true
	ca_timeout = 60
	zone = "System"
}
`, shareName, shareName)

var SmbShareNameUpdatedResourceConfig = FileSystemResourceConfigCommon6 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	auto_create_directory = true
	name = "%s_update"
	path = "/ifs/%s_update_path"
	permissions = [
		{
			permission = "full"
			permission_type = "allow"
			trustee = {
				id = "SID:S-1-1-0",
				name = "Everyone",
				type = "wellknown"
			}
		}
	]
	allow_delete_readonly = true
	ca_timeout = 30
	zone = "System"
}
`, shareName, shareName)

var SmbShareUpdatedZoneResourceConfig = FileSystemResourceConfigCommon6 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	auto_create_directory = true
	name = "%s"
	path = "/ifs/%s"
	permissions = [
		{
			permission = "full"
			permission_type = "allow"
			trustee = {
				id = "SID:S-1-1-0",
				name = "Everyone",
				type = "wellknown"
			}
		}
	]
	allow_delete_readonly = true
	ca_timeout = 30
	zone = "System2"
}
`, shareName, shareName)
