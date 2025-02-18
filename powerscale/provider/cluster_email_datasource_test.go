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

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterEmailDataSource(t *testing.T) {
	var clusterEmailName = "data.powerscale_cluster_email.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + clusterEmailDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterEmailName, "id"),
					resource.TestCheckResourceAttrSet(clusterEmailName, "settings.%"),
					resource.TestCheckResourceAttrSet(clusterEmailName, "settings.batch_mode"),
					resource.TestCheckResourceAttrSet(clusterEmailName, "settings.smtp_auth_passwd_set"),
					resource.TestCheckResourceAttrSet(clusterEmailName, "settings.smtp_auth_security"),
					resource.TestCheckResourceAttrSet(clusterEmailName, "settings.smtp_port"),
					resource.TestCheckResourceAttrSet(clusterEmailName, "settings.use_smtp_auth"),
				),
			},
		},
	})
}

func TestAccClusterEmailDatasourceErrorGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetClusterVersion).Return("", fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var clusterEmailDataSourceConfig = `
data "powerscale_cluster_email" "test" {
}
`
