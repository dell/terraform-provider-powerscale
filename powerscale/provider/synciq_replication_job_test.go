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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var host = ""
var user = ""
var password = ""

func TestAccSyncIQReplicationJobResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + createReplicationJobConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_synciq_replication_job.job1", "id", "TerraformPolicy"),
				),
			},
			{
				Config: ProviderConfig + updateReplicationJobConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_synciq_replication_job.job1", "id", "TerraformPolicy"),
				),
			},
		},
	})
}

var createReplicationJobConfig = fmt.Sprintf(`
resource "terraform_data" "large_file" {
  provisioner "remote-exec" {
    inline = [
      "mkdir -p /ifs/terraform/source",
      "head -c 1000000000 /dev/urandom > /ifs/terraform/source/large_file.dat",
      "mkdir -p /ifs/terraform/target"
    ]
    connection {
      host     = "` + host + `"
      user     = "` + user + `"
      password = "` + password + `"
      type     = "ssh"
    }
  }

  provisioner "remote-exec" {
    when = destroy
    inline = [
      "rm -rf /ifs/terraform",
      "echo 'yes' | isi sync rules delete bw-0",
    ]
    connection {
      host     = "` + host + `"
      user     = "` + user + `"
      password = "` + password + `"
      type     = "ssh"
    }
  }
}

resource "powerscale_synciq_policy" "policy" {
  name             = "TerraformPolicy"
  action           = "sync"
  source_root_path = "/ifs/terraform/source"
  target_host      = "` + host + `"
  target_path      = "/ifs/terraform/target"
  depends_on       = [terraform_data.large_file]
}

resource "powerscale_synciq_rules" "kb-10" {
  bandwidth_rules = [
    {
      limit = 10
      schedule = {
        begin        = "00:00"
        days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
        end          = "23:59"
      }
    },
  ]
}

resource "powerscale_synciq_replication_job" "job1" {
  action = "run"
  id     = "TerraformPolicy"
  is_paused = false
  depends_on = [powerscale_synciq_policy.policy]
}
`)

var updateReplicationJobConfig = `
resource "powerscale_synciq_replication_job" "job1" {
  action = "run"
  id     = "TerraformPolicy"
  is_paused = true
}
`
