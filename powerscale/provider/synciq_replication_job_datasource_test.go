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
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSyncIQReplicationJobDataSource(t *testing.T) {
	diags := diag.Diagnostics{}
	diags.AddError("mock err", "mock err")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + SyncIQReplicationJobDataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ManageDataSourceSyncIQReplicationJob).Return(nil, diags).Build()
				},
				Config:      ProviderConfig + SyncIQReplicationJobDataSourceConfigFilter,
				ExpectError: regexp.MustCompile(`.*mock err*.`),
			},
		},
	})
}

var SyncIQReplicationJobDataSourceConfig = `
	data "powerscale_synciq_replication_job" "all" {
	}
`

var SyncIQReplicationJobDataSourceConfigFilter = `
	data "powerscale_synciq_replication_job" "all" {
		filter {
			state = "running"
		  }
	}
`
