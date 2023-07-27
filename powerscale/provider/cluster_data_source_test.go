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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: ProviderConfig + testAccClusterDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_cluster.test", "id", "cluster-data-source"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "config.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "config.devices.#"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "config.onefs_version.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "config.timezone.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "identity.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "identity.logon.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "internal_networks.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.#"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.drives.#"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.drives.0.firmware.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.hardware.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.partitions.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.sensors.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.state.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "nodes.nodes.0.status.%"),
					resource.TestCheckResourceAttrSet("data.powerscale_cluster.test", "acs.%"),
				),
			},
		},
	})
}

func TestAccClusterConfigError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetClusterConfig).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Client Error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccClusterIdentityError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Client Error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccClusterNodesError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetClusterNodes).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Client Error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccClusterInternalNetworksError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetClusterInternalNetworks).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Client Error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

func TestAccClusterAcsError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ListClusterAcs).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Client Error*.`),
			},
		},
		CheckDestroy: func(_ *terraform.State) error {
			if FunctionMocker != nil {
				FunctionMocker.UnPatch()
			}
			return nil
		},
	})
}

var testAccClusterDataSourceConfig = `
data "powerscale_cluster" "test" {
}
`
