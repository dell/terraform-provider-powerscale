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

func TestAccSyncIQReplicationJobResource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + errorReplicationJob,
				ExpectError: regexp.MustCompile(`.*SyncIQ Replication Job cannot be paused befor job creation.*`),
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("Error reading create plan")).Build()
				},
				Config:      ProviderConfig + SetupReplication() + createReplicationJob,
				ExpectError: regexp.MustCompile("Error reading create plan"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CreateSyncIQReplicationJob).Return(nil, fmt.Errorf("Error creating syncIQ Replication Job")).Build()
				},
				Config:      ProviderConfig + SetupReplication() + createReplicationJob,
				ExpectError: regexp.MustCompile("Error creating syncIQ Replication Job"),
			},
			{
				// create synciq replication job positive test
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + SetupReplication() + createReplicationJob,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_synciq_replication_job.job1", "id", "TerraformPolicy"),
				),
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSyncIQReplicationJob).Return(nil, nil, fmt.Errorf("Error reading syncIQ Replication Job")).Build()
				},
				Config:      ProviderConfig + SetupReplication() + updateReplicationJob,
				ExpectError: regexp.MustCompile("Error reading syncIQ Replication Job"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.UpdateSyncIQReplicationJob).Return(nil, fmt.Errorf("Error updating syncIQ Replication Job")).Build()
				},
				Config:      ProviderConfig + SetupReplication() + updateReplicationJob,
				ExpectError: regexp.MustCompile("Error updating syncIQ Replication Job"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config: ProviderConfig + SetupReplication() + updateReplicationJob,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_synciq_replication_job.job1", "id", "TerraformPolicy"),
				),
			},
		},
	})
}

func SetupReplication() string {
	connection := fmt.Sprintf(`
  connection {
      host     = "%s"
	  port     = %s
      user     = "%s"
      password = "%s"
      type     = "ssh"
    }
  `, powerScaleSSHIP, powerscaleSSHPort, powerscaleUsername, powerscalePassword)

	createLargeFile := fmt.Sprintf(`
  resource "terraform_data" "large_file" {
    provisioner "remote-exec" {
      inline = [
        "mkdir -p /ifs/terraformAT/source",
        "head -c 10000000 /dev/urandom > /ifs/terraformAT/source/large_file.dat",
        "mkdir -p /ifs/terraformAT/target",
		" isi sync settings modify --service=on",
		"isi sync rules create bandwidth 00:00-23:59 X-S 129",
		"echo 'confirm create policy' | isi sync policies create --name=TerraformPolicy --source-root-path=/ifs/terraformAT/source/ --target-host=%s --target-path=/ifs/terraformAT/target/ --action=sync",
		"sleep 10",
      ]
      `+connection+`
    }

    provisioner "remote-exec" {
      when = destroy
      inline = [
        "rm -rf /ifs/terraform",
        "echo 'yes' | isi sync rules delete bw-0",
		"echo 'yes' | isi sync job cancel --all",
		"echo 'yes' | isi sync pol rm --all",
      ]
      `+connection+`
    }
  }`, powerScaleSSHIP)
	return createLargeFile
}

var createReplicationJob = `
resource "powerscale_synciq_replication_job" "job1" {
  action = "run"
  id     = "TerraformPolicy"
  is_paused = false
  depends_on = [terraform_data.large_file]
}
`

var updateReplicationJob = `
resource "powerscale_synciq_replication_job" "job1" {
  action = "run"
  id     = "TerraformPolicy"
  is_paused = true
  depends_on = [terraform_data.large_file]
}
`

var errorReplicationJob = `
resource "powerscale_synciq_replication_job" "errorJob" {
  action = "resync_prep"
  id     = "TerraformPolicy"
  is_paused = true
}`
