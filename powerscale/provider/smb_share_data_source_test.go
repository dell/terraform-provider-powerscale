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

func TestAccSmbShareDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + SmbShareDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_smb_share.share_datasource_test", "smb_shares.#", "1"),
				),
			},
		},
	})
}

func TestAccSmbShareDatasourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + SmbShareAllDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_smb_share.share_datasource_test_all", "filter.#", "0"),
				),
			},
		},
	})
}

func TestAccSmbShareDatasourceGetError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ListSmbShares).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareAllDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSmbShareDatasourcePagination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					resume := "1"
					shares := powerscale.V7SmbShares{
						Digest: nil,
						Resume: &resume,
						Shares: nil,
						Total:  nil,
					}
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(mockey.GetMethod(powerscale.ApiListProtocolsv7SmbSharesRequest{}, "Execute")).Return(&shares, nil, nil).Build()
					FunctionMocker.When(func() bool {
						shares := powerscale.V7SmbShares{
							Digest: nil,
							Resume: nil,
							Shares: []powerscale.V7SmbShareExtended{{Id: shareName}},
							Total:  nil,
						}
						if FunctionMocker.MockTimes() > 0 {
							FunctionMocker.UnPatch()
							FunctionMocker = mockey.Mock(mockey.GetMethod(powerscale.ApiListProtocolsv7SmbSharesRequest{}, "Execute")).Return(&shares, nil, nil).Build()
						}
						return FunctionMocker.MockTimes() == 0
					})
				},
				Config: ProviderConfig + SmbShareAllDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_smb_share.share_datasource_test_all", "filter.#", "0"),
				),
			},
		},
	})
}

func TestAccSmbShareDatasourceErrorCopyFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareAllDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			//Read testing with names
			{
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var SmbShareAllDatasourceConfig = FileSystemResourceConfigCommon4 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_resource_test" {
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
}

data "powerscale_smb_share" "share_datasource_test_all" {}
`, shareName, shareName)

var SmbShareDatasourceConfig = FileSystemResourceConfigCommon4 + fmt.Sprintf(`
resource "powerscale_smb_share" "share_resource_test" {
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
}

data "powerscale_smb_share" "share_datasource_test" {
	filter {
		names = ["%s"]
		zone  = "System"
		scope = "effective"
		sort  = "id"
		offset= 0
		dir = "ASC"
	}
  	depends_on = [
    	powerscale_smb_share.share_resource_test
  	]
}
`, shareName, shareName, shareName)
