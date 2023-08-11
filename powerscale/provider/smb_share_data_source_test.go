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
					FunctionMocker = Mock(GetMethod(powerscale.ApiListProtocolsv7SmbSharesRequest{}, "Execute")).Return(&shares, nil, nil).Build()
					FunctionMocker.When(func() bool {
						shares := powerscale.V7SmbShares{
							Digest: nil,
							Resume: nil,
							Shares: []powerscale.V7SmbShareExtended{{Id: shareName}},
							Total:  nil,
						}
						if FunctionMocker.MockTimes() > 0 {
							FunctionMocker.UnPatch()
							FunctionMocker = Mock(GetMethod(powerscale.ApiListProtocolsv7SmbSharesRequest{}, "Execute")).Return(&shares, nil, nil).Build()
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
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SmbShareAllDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var SmbShareAllDatasourceConfig = fmt.Sprintf(`
resource "powerscale_smb_share" "share_resource_test" {
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

var SmbShareDatasourceConfig = fmt.Sprintf(`
resource "powerscale_smb_share" "share_resource_test" {
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
		resolve_names = true
		names = ["%s"]
		limit = 1
	}
  	depends_on = [
    	powerscale_smb_share.share_resource_test
  	]
}
`, shareName, shareName, shareName)
