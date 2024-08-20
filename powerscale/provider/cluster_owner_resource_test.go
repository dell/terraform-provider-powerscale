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

func TestAccClusterOwnerResourceImport(t *testing.T) {
	var clusterOwnerResourceName = "powerscale_cluster_owner.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with invalid fields
			{
				Config:      ProviderConfig + clusterOwnerResourceConfigInvalid,
				ExpectError: regexp.MustCompile(`.*Please provide at least one of the following.*`),
			},
			// Create Success
			{
				Config: ProviderConfig + clusterOwnerResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "company", "company_name"),
				),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:       ProviderConfig + clusterOwnerResourceConfig,
				ResourceName: clusterOwnerResourceName,
				ImportState:  true,
				ExpectError:  regexp.MustCompile(`.*mock error*.`),
			},
			// Import testing
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:       ProviderConfig + clusterOwnerResourceConfig,
				ResourceName: clusterOwnerResourceName,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(clusterOwnerResourceName, "id")
					resource.TestCheckResourceAttrSet(clusterOwnerResourceName, "company")
					resource.TestCheckResourceAttrSet(clusterOwnerResourceName, "primary_email")
					resource.TestCheckResourceAttrSet(clusterOwnerResourceName, "secondary_email")
					return nil
				},
			},
		},
	})
}

func TestAccClusterOwnerResourceUpdate(t *testing.T) {
	var clusterOwnerResourceName = "powerscale_cluster_owner.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterOwnerResourceConfig,
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterOwner).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerResourceConfig,
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
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterOwner).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterOwner).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update Mock Error
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source interface{}, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Update and Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
				},
				Config: ProviderConfig + clusterOwnerUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "company", "company_name_update"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + clusterOwnerUpdateRevertResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "company", "company_name"),
				),
			},
		},
	})
}

func TestAccClusterOwnerResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterOwner).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterOwner).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccClusterOwnerResourceUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterOwner).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterOwner).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterOwnerResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestAccClusterOwnerResourceImportMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterOwnerResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterOwner).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + clusterOwnerResourceConfig,
				ResourceName:      "powerscale_cluster_owner.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

var clusterOwnerResourceConfigInvalid = `
resource "powerscale_cluster_owner" "test" {
}
`

var clusterOwnerResourceConfig = `
resource "powerscale_cluster_owner" "test" {
	company = "company_name"
}
`

var clusterOwnerUserTemplateSetResourceConfig = `
resource "powerscale_cluster_owner" "test" {
	company = "company_name"
}
`

var clusterOwnerUserTemplateEmptyResourceConfig = `
resource "powerscale_cluster_owner" "test" {
	company = ""
}
`

var clusterOwnerUpdateResourceConfig = `
resource "powerscale_cluster_owner" "test" {
	company = "company_name_update"
	
}
`

var clusterOwnerUpdateRevertResourceConfig = `
resource "powerscale_cluster_owner" "test" {
	company = "company_name"
}
`
