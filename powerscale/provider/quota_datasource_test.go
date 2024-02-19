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
	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

func TestAccQuotaDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + QuotaDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerscale_quota.quota_datasource_test", "quotas.#"),
				),
			},
		},
	})
}

func TestAccQuotaDatasourceErrorListQuota(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ListQuotas).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + QuotaDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccQuotaDatasourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + QuotaDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var QuotaDatasourceConfig = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "directory"
	include_snapshots = false
}

data "powerscale_quota" "quota_datasource_test" {
  filter {
    enforced = false
    exceeded = false
    include_snapshots = false
    path = "/ifs/tfacc_quota_test"
    recurse_path_children = true
    recurse_path_parents  = true
    type = "directory"
    zone  = "System"
  }
  depends_on = [	
    powerscale_quota.quota_test
  ]
}
`
