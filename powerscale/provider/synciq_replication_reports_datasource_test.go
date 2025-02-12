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
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReplicationReportsDataSourceAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + SetupHostIP() + RRDataSourceConfig + SyncIQSettingsConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccReplicationReportsDataSourceFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			// read all
			{
				Config: ProviderConfig + SetupHostIP() + RRDataSourceConfigFilter + SyncIQSettingsConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(),
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
				Config:      ProviderConfig + RRDataSourceFilterConfigErr + SyncIQSettingsConfig,
				ExpectError: regexp.MustCompile(`.*Unsupported argument*.`),
			},
		},
	})
}

func TestAccReplicationReportsDataSourceGettingErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetReplicationReports).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SetupHostIP() + RRDataSourceConfig + SyncIQSettingsConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func SetupHostIP() string {

	if powerScaleSSHIP == "localhost" || powerScaleSSHIP == "0.0.0.0" {
		return `
		locals {
			host_ip = "10.10.10.10"
		}
	`
	}
	var result = fmt.Sprintf(`
		locals {
			host_ip = "%s"
		}
	`, powerScaleSSHIP)

	return result

}

var JobConfig = `
resource "powerscale_synciq_policy" "policy1" {
	name             = "tfaccPolicy"
	action           = "sync"
	source_root_path = "/ifs"
	target_host      = local.host_ip
	target_path      = "/ifs/tfaccSink"
}

resource "powerscale_synciq_replication_job" "job1" {
  action = "run"
  id     = powerscale_synciq_policy.policy1.id
  is_paused = false
  wait_time = 5
  depends_on = [
    powerscale_synciq_policy.policy1
  ]
}

resource "time_sleep" "wait_60_seconds" {
  create_duration = "90s"

  depends_on = [powerscale_synciq_replication_job.job1]
}
`

var RRDataSourceConfig = JobConfig + `
data "powerscale_synciq_replication_report" "all" {
	depends_on = [time_sleep.wait_60_seconds]
}
`

var RRDataSourceConfigFilter = JobConfig + `
data "powerscale_synciq_replication_report" "filtering" {
	filter {
		reports_per_policy = 1
	}
	depends_on = [time_sleep.wait_60_seconds]
}
`

var RRDataSourceNameConfigErr = `
data "powerscale_synciq_replication_report" "test" {
	filter {
		policy_name = "InvalidName"
	}
}
`

var RRDataSourceFilterConfigErr = `
data "powerscale_synciq_replication_report" "test" {
	filter {
		invalidFilter = "Invalid"
	}
}
`
