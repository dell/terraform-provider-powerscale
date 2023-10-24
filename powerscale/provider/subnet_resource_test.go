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
	"context"
	"fmt"
	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"
)

var subnetMocker *mockey.Mocker

func TestAccSubnetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SubnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "name", "tfacc_subnet"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "addr_family", "ipv4"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "prefixlen", "21"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_subnet.subnet",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "groupnet0.tfacc_subnet", states[0].Attributes["id"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + SubnetResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "name", "tfacc_subnet_rename"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "addr_family", "ipv4"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "prefixlen", "21"),
				),
			},
		},
	})
}

func TestAccSubnetResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SubnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "name", "tfacc_subnet"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "addr_family", "ipv4"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "prefixlen", "21"),
				),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_subnet.subnet",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetSubnet).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing wrong ID
			{
				ResourceName:  "powerscale_subnet.subnet",
				ImportState:   true,
				ImportStateId: "subnetId",
				ExpectError:   regexp.MustCompile("Unexpected Import Identifier"),
			},
		},
	})
}

func TestAccSubnetResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SubnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "name", "tfacc_subnet"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "addr_family", "ipv4"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "prefixlen", "21"),
				),
			},
			// Update get error
			{
				Config: ProviderConfig + SubnetResourceUpdateConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateSubnet).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + SubnetResourceUpdateConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetSubnet).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + SubnetResourceUpdateConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSubnetResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + SubnetInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccSubnetResourceErrorDelete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SubnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "name", "tfacc_subnet"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "addr_family", "ipv4"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "prefixlen", "21"),
				),
			},
			{
				PreConfig: func() {
					if subnetMocker != nil {
						subnetMocker.UnPatch()
					}
					subnetMocker = mockey.Mock(helper.DeleteSubnet).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, subnetName, groupnet string) bool {
							return subnetMocker.Times() == 1
						})
				},
				Config:      ProviderConfig + SubnetResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSubnetResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + SubnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "name", "tfacc_subnet"),
					resource.TestCheckResourceAttr("powerscale_subnet.subnet", "groupnet", "groupnet0"),
				),
			},
			{
				ResourceName: "powerscale_subnet.subnet",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SubnetResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SubnetResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccSubnetResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					subnetMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + SubnetResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var SubnetResourceConfig = `
resource "powerscale_subnet" "subnet" {
  name = "tfacc_subnet"
  groupnet = "groupnet0"
  addr_family = "ipv4"
  prefixlen = 21
}
`

var SubnetResourceUpdateConfig = `
resource "powerscale_subnet" "subnet" {
  name = "tfacc_subnet_rename"
  groupnet = "groupnet0"
  addr_family = "ipv4"
  prefixlen = 21
}
`
var SubnetInvalidResourceConfig = `
resource "powerscale_subnet" "subnet" {
  name = "tfacc_subnet_invalid"
  groupnet = "groupnet0"
  addr_family = "invalid"
  prefixlen = 21
}
`
