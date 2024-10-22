/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

func TestAccReplicationReportsDataSourceAll(t *testing.T) {
	var rrTerraformName = "data.powerscale_replication_report.all"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + RRDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(rrTerraformName, "replication_reports.#"),
				),
			},
		},
	})
}

func TestAccReplicationReportsDataSourceFilter(t *testing.T) {
	var rrTerraformName = "data.powerscale_replication_report.filtering"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + RRDataSourceConfigFilter,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(rrTerraformName, "replication_reports.#"),
				),
			},
		},
	})
}

func TestAccReplicationReportsDataSourceFilterErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + RRDataSourceFilterConfigErr,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccReplicationReportsDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetReplicationReports).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + RRDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

var RRDataSourceConfig = `
data "powerscale_replication_report" "all" {
}
`

var RRDataSourceConfigFilter = `
data "powerscale_replication_report" "filtering" {
	filter {
		reports_per_policy = 2
	}
}
`

var RRDataSourceNameConfigErr = `
data "powerscale_replication_report" "test" {
	filter {
		policy_name = "InvalidName"
	}
}
`

var RRDataSourceFilterConfigErr = `
data "powerscale_role" "test" {
	filter {
		invalidFilter = "Invalid"
	}
}
`
