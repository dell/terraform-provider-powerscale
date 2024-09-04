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

func TestAccClusterIdentityResourceImport(t *testing.T) {
	var clusterIdentityResourceName = "powerscale_cluster_identity.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Success
			{
				Config: ProviderConfig + clusterIdentityResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterIdentityResourceName, "name", "cluster1"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:       ProviderConfig + clusterIdentityResourceConfig,
				ResourceName: clusterIdentityResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
			// Import testing
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:       ProviderConfig + clusterIdentityResourceConfig,
				ResourceName: clusterIdentityResourceName,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(clusterIdentityResourceName, "id")
					resource.TestCheckResourceAttrSet(clusterIdentityResourceName, "name")
					resource.TestCheckResourceAttrSet(clusterIdentityResourceName, "description")
					return nil
				},
			},
		},
	})
}

// TestAccClusterIdentityResource - Tests the creation of a cluster Identity resource.
func TestAccClusterIdentityResource(t *testing.T) {
	var clusterIdentityResourceName = "powerscale_cluster_identity.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterIdentityResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterIdentityResourceName, "name", "cluster1"),
				),
			},
		},
	})
}

// TestAccClusterIdentityResource_Update - Tests the update of a cluster Identity resource along with error mocking.
func TestAccClusterIdentityResource_Update(t *testing.T) {
	var clusterIdentityResourceName = "powerscale_cluster_identity.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterIdentityResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterIdentityResourceName, "id"),
					resource.TestCheckResourceAttr(clusterIdentityResourceName, "name", "cluster1"),
				),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfig,
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
				Config:      ProviderConfig + clusterIdentityResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterIdentity).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + clusterIdentityResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + clusterIdentityResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterIdentityResourceName, "name", "cluster2"),
				),
			},
			{
				Config: ProviderConfig + clusterIdentityResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(clusterIdentityResourceName, "id"),
					resource.TestCheckResourceAttr(clusterIdentityResourceName, "name", "cluster1"),
				),
			},
		},
	})
}

// TestAccClusterIdentityResource_Create - Tests the mock errors during the create operation of the cluster Identity resource.
func TestAccClusterIdentityResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterIdentity).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterIdentityState).Return(nil).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfig,
				ExpectError: regexp.MustCompile(`.*Value Conversion*.`),
			},
		},
	})
}

// TestAccClusterIdentityResource_Update - Tests the mock errors during the update operation of the cluster Identity resource.
func TestAccClusterIdentityResourceUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterIdentity).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterIdentityResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

// TestAccClusterIdentityResource_Import - Tests the mock errors during the import of the cluster Identity resource.
func TestAccClusterIdentityResourceImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterIdentityResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterIdentity).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + clusterIdentityResourceConfig,
				ResourceName:      "powerscale_cluster_identity.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterIdentityState).Return(nil).Build()
				},
				Config:            ProviderConfig + clusterIdentityResourceConfig,
				ResourceName:      "powerscale_cluster_identity.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*Value Conversion*.`),
				ImportStateVerify: true,
			},
		},
	})
}

var clusterIdentityResourceConfig = `
resource "powerscale_cluster_identity" "test" {
	name = "cluster1"
	description = "cluster name"
	logon= {motd = "motd",motd_header = "motd_header"}
}
`

var clusterIdentityResourceConfigUpdate = `
resource "powerscale_cluster_identity" "test" {
	name = "cluster2"
}
`
