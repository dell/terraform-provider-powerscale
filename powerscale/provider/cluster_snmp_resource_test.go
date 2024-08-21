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
	"context"
	"fmt"
	"regexp"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccClusterSnmpResource - Tests the creation of a cluster SNMP resource.
func TestAccClusterSnmpResource(t *testing.T) {
	var clusterSNMPResourceName = "powerscale_cluster_snmp.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterSnmpResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterSNMPResourceName, "enabled", "true"),
				),
			},
		},
	})
}

// TestAccClusterSnmpResource_Update - Tests the update of a cluster SNMP resource along with error mocking.
func TestAccClusterSnmpResource_Update(t *testing.T) {
	var clusterSNMPResourceName = "powerscale_cluster_snmp.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterSnmpResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterSNMPResourceName, "id"),
					resource.TestCheckResourceAttr(clusterSNMPResourceName, "enabled", "true"),
				),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterSNMP).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterSNMP).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterSNMP).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + clusterSnmpResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + clusterSnmpResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterSNMPResourceName, "enabled", "false"),
				),
			},
			{
				Config: ProviderConfig + clusterSnmpResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterSNMPResourceName, "id"),
					resource.TestCheckResourceAttr(clusterSNMPResourceName, "enabled", "true"),
				),
			},
		},
	})
}

// TestAccClusterSnmpResource_Create - Tests the mock errors during the create operation of the cluster SNMP resource.
func TestAccClusterSnmpResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterSNMP).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterSNMP).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

// TestAccClusterSnmpResource_Update - Tests the mock errors during the update operation of the cluster SNMP resource.
func TestAccClusterSnmpResourceUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterSNMP).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterSNMP).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterSnmpResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

// TestAccClusterSnmpResource_Import - Tests the mock errors during the import of the cluster SNMP resource.
func TestAccClusterSnmpResourceImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterSnmpResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterSNMP).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + clusterSnmpResourceConfig,
				ResourceName:      "powerscale_cluster_snmp.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccClusterSnmpResource_Import - Tests the import of the cluster SNMP resource.
func TestAccClusterSnmpResource_Import(t *testing.T) {
	var clusterSNMPResourceName = "powerscale_cluster_snmp.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with invalid fields
			{
				Config:      ProviderConfig + clusterSnmpResourceEmptyConfig,
				ExpectError: regexp.MustCompile(`.*Missing required argument.*`),
			},
			// Create Success
			{
				Config: ProviderConfig + clusterSnmpResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterSNMPResourceName, "enabled", "true"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterSNMP).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:       ProviderConfig + clusterSnmpResourceConfig,
				ResourceName: clusterSNMPResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
			// Import testing
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:       ProviderConfig + clusterSnmpResourceConfig,
				ResourceName: clusterSNMPResourceName,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(clusterSNMPResourceName, "id")
					resource.TestCheckResourceAttrSet(clusterSNMPResourceName, "enabled")
					resource.TestCheckResourceAttrSet(clusterSNMPResourceName, "snmp_v1_v2c_access")
					resource.TestCheckResourceAttrSet(clusterSNMPResourceName, "snmp_v3_access")
					return nil
				},
			},
		},
	})
}

var clusterSnmpResourceUpdateConfig = `
resource "powerscale_cluster_snmp" "test" {
	enabled = false
	snmp_v1_v2c_access = false
	read_only_community = "read_only_community"
}
`

var clusterSnmpResourceConfig = `
resource "powerscale_cluster_snmp" "test" {
	enabled = true
	snmp_v1_v2c_access = true
	read_only_community = "read_only_community"
}
`
var clusterSnmpResourceEmptyConfig = `
resource "powerscale_cluster_snmp" "test" {}
`
