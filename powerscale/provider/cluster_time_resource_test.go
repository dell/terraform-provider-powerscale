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

var clusterTimeResourceConfigInvalid = `
resource "powerscale_cluster_time" "test" {

}
`

var clusterTimeResourceConfigInvalid2 = `
resource "powerscale_cluster_time" "test" {
	date = "2024"
}
`

var clusterTimeResourceConfigInvalid3 = `
resource "powerscale_cluster_time" "test" {
	time = "2024"
}
`

var clusterTimeResourceConfigInvalid4 = `
resource "powerscale_cluster_time" "test" {
	path = "invalid"
}
`

var clusterTimeResourceConfig = `
resource "powerscale_cluster_time" "test" {
	date = "01/12/2024"
  	time = "00:32"
}
`

var clusterTimeResourceConfigUpdate = `
resource "powerscale_cluster_time" "test" {
	date = "12/12/2024"
  	time = "10:32"
}
`

func TestAccClusterTimeResource(t *testing.T) {
	clusterOwnerResourceName := "powerscale_cluster_time.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Case 1 - Create the resource with invalid fields.
			{
				Config:      ProviderConfig + clusterTimeResourceConfigInvalid,
				ExpectError: regexp.MustCompile(`.*Please provide at least one of the following*`),
			},
			// Case 2 - Create the resource with invalid value.
			{
				Config:      ProviderConfig + clusterTimeResourceConfigInvalid2,
				ExpectError: regexp.MustCompile(`.* Please follow the format MM/dd/yyyy*`),
			},
			// Case 3 - Create the resource with invalid value.
			{
				Config:      ProviderConfig + clusterTimeResourceConfigInvalid3,
				ExpectError: regexp.MustCompile(`.*Please follow the format HH:mm*`),
			},
			// Case 4 - Create the resource with invalid value.
			{
				Config:      ProviderConfig + clusterTimeResourceConfigInvalid4,
				ExpectError: regexp.MustCompile(`.*Africa/Abidjan*`),
			},
			// Case 5 - Create the resource with valid value.
			{
				Config:      ProviderConfig + clusterTimeResourceConfig,
				ExpectError: nil,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "date", "01/12/2024"),
				),
			},
			// Case 6 - Read Refresh
			{
				Config:      ProviderConfig + clusterTimeResourceConfig,
				ExpectError: nil,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "date", "01/12/2024"),
				),
			},
			// Case 7 - Update
			{
				Config:      ProviderConfig + clusterTimeResourceConfigUpdate,
				ExpectError: nil,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "date", "12/12/2024"),
				),
			},
		},
	})
}

func TestAccClusterTimeResourceMock(t *testing.T) {
	clusterOwnerResourceName := "powerscale_cluster_time.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Create
			{
				Config:      ProviderConfig + clusterTimeResourceConfig,
				ExpectError: nil,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterOwnerResourceName, "date", "01/12/2024"),
				),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterTime).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterTimeResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Read testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterTimeZone).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterTimeResourceConfig,
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
				Config:      ProviderConfig + clusterTimeResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterTimeZone).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + clusterTimeResourceConfig,
				ResourceName:      "powerscale_cluster_time.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}
