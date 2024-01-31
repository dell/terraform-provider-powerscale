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
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var filePoolPolicyDsMocker *mockey.Mocker

func TestAccFilePoolPolicyDataSource(t *testing.T) {
	var policyTerraformName = "data.powerscale_filepool_policy.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + filePoolPolicyResourceConfig + filePoolPolicyAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(policyTerraformName, "file_pool_policies.#"),
				),
			},
			// filter with default policy read testing
			{
				Config: ProviderConfig + filePoolPolicyResourceConfig + filePoolPolicyFilterDefaultDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(policyTerraformName, "file_pool_policies.#"),
				),
			},
			// filter with names read testing
			{
				Config: ProviderConfig + filePoolPolicyResourceConfig + filePoolPolicyFilterNameDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyTerraformName, "file_pool_policies.0.name", "tfacc_filePoolPolicy"),
					resource.TestCheckResourceAttr(policyTerraformName, "file_pool_policies.0.description", "tfacc_filePoolPolicy description"),
					resource.TestCheckResourceAttr(policyTerraformName, "file_pool_policies.0.apply_order", "1"),
					resource.TestCheckResourceAttr(policyTerraformName, "file_pool_policies.0.actions.#", "7"),
					resource.TestCheckResourceAttr(policyTerraformName, "file_pool_policies.0.file_matching_pattern.or_criteria.#", "2"),
				),
			},
		},
	})
}

func TestAccFilePoolPolicyDataSourceInvalidNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid names read testing
			{
				Config:      ProviderConfig + filePoolPolicyInvalidFilterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error one or more of the filtered File Pool Policy names is not a valid powerscale File Pool Policy.*.`),
			},
		},
	})
}

func TestAccFilePoolPolicyDatasourceErrorGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyDsMocker != nil {
						filePoolPolicyDsMocker.UnPatch()
					}
					filePoolPolicyDsMocker = mockey.Mock((*powerscale.FilepoolApiService).ListFilepoolv12FilepoolPoliciesExecute).Return(nil, nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceConfig + filePoolPolicyAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccFilePoolPolicyDatasourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyDsMocker != nil {
						filePoolPolicyDsMocker.UnPatch()
					}
					filePoolPolicyDsMocker = mockey.Mock(helper.UpdateFilePoolPolicyDataSourceState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceConfig + filePoolPolicyAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccFilePoolPolicyDatasourceReleaseMock(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyDsMocker != nil {
						filePoolPolicyDsMocker.Release()
					}
				},
				Config: ProviderConfig + filePoolPolicyResourceConfig + filePoolPolicyAllDataSourceConfig,
			},
		},
	})
}

var filePoolPolicyAllDataSourceConfig = `

data "powerscale_filepool_policy" "test" {
    depends_on = [powerscale_filepool_policy.policy_test]
}
`

var filePoolPolicyFilterDefaultDataSourceConfig = `

data "powerscale_filepool_policy" "test" {
	filter {
    names = ["Default policy"]
  }
  depends_on = [powerscale_filepool_policy.policy_test]
}
`

var filePoolPolicyFilterNameDataSourceConfig = `

data "powerscale_filepool_policy" "test" {
	filter {
    names = [powerscale_filepool_policy.policy_test.name]
  }
  depends_on = [powerscale_filepool_policy.policy_test]
}
`

var filePoolPolicyInvalidFilterDataSourceConfig = `
data "powerscale_filepool_policy" "test" {
	filter {
    names = ["invalidName"]
  }
}
`
