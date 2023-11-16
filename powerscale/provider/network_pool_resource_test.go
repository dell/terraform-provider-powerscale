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
	powerscale "dell/powerscale-go-client"
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"
)

func TestAccNetworkPoolResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NetworkPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "name", poolName),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "subnet", "subnet0"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "access_zone", "System"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "sc_ttl", "0"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "nfsv3_rroce_only", "false"),
				),
			},
			{
				ResourceName:  "powerscale_networkpool.pool_test",
				ImportStateId: "groupnet0,subnet0,pool_test",
				ImportState:   true,
				ExpectError:   regexp.MustCompile("Unexpected Import Identifier"),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_networkpool.pool_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, poolName, states[0].Attributes["name"])
					assert.Equal(t, "groupnet0", states[0].Attributes["groupnet"])
					assert.Equal(t, "subnet0", states[0].Attributes["subnet"])
					assert.Equal(t, "System", states[0].Attributes["access_zone"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + NetworkPoolUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "alloc_method", "dynamic"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "sc_ttl", "1"),
				),
			},
		},
	})
}

func TestAccNetworkPoolResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NetworkPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "name", poolName),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "subnet", "subnet0"),
				),
			},
			// ImportState testing get none pool
			{
				ResourceName: "powerscale_networkpool.pool_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetNetworkPool).Return(&powerscale.V12GroupnetsGroupnetSubnetsSubnetPools{}, nil).Build()
				},
				ExpectError: regexp.MustCompile("not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_networkpool.pool_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNetworkPool).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNetworkPoolResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NetworkPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "name", poolName),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "subnet", "subnet0"),
				),
			},
			// Update param read error
			{
				Config: ProviderConfig + NetworkPoolUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get error
			{
				Config: ProviderConfig + NetworkPoolUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.UpdateNetworkPool).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none pool
			{
				Config: ProviderConfig + NetworkPoolUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNetworkPool).Return(&powerscale.V12GroupnetsGroupnetSubnetsSubnetPools{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, poolModel models.NetworkPoolResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".*Could not read updated network pool*."),
			},
			// Update get error
			{
				Config: ProviderConfig + NetworkPoolUpdatedResourceConfig2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetNetworkPool).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, poolModel models.NetworkPoolResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			//Update Invalid Config
			{
				PreConfig: func() {
					FunctionMocker.Release()
				},
				Config:      ProviderConfig + NetworkPoolInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccNetworkPoolResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + NetworkPoolInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccNetworkPoolResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + NetworkPoolResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "name", poolName),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "groupnet", "groupnet0"),
					resource.TestCheckResourceAttr("powerscale_networkpool.pool_test", "subnet", "subnet0"),
				),
			},
			{
				ResourceName: "powerscale_networkpool.pool_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NetworkPoolResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + NetworkPoolUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccNetworkPoolResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + NetworkPoolResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var poolName = "pool_at_test"

var NetworkPoolResourceConfig = fmt.Sprintf(`
resource "powerscale_networkpool" "pool_test" {
	name = "%s"
	groupnet = "groupnet0"
	subnet = "subnet0"
	ranges = [
		{
			high = "10.225.108.142",
			low = "10.225.108.142"
		}
	]
}
`, poolName)

var NetworkPoolInvalidResourceConfig = fmt.Sprintf(`
resource "powerscale_networkpool" "pool_test" {
	name = "%s"
	groupnet = "groupnet0"
	subnet = "subnet0"
	ranges = [
		{
			high = "10.225.108.142",
			low = "10.225.108.142"
		}
	]
	alloc_method = "invalid"
}
`, poolName)

var NetworkPoolUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_networkpool" "pool_test" {
	name = "%s"
	groupnet = "groupnet0"
	subnet = "subnet0"
	ranges = [
		{
			high = "10.225.108.142",
			low = "10.225.108.142"
		}
	]
	alloc_method = "dynamic"
	sc_ttl = 1
}
`, poolName)

var NetworkPoolUpdatedResourceConfig2 = fmt.Sprintf(`
resource "powerscale_networkpool" "pool_test" {
	name = "%s"
	groupnet = "groupnet0"
	subnet = "subnet0"
	ranges = [
		{
			high = "10.225.108.142",
			low = "10.225.108.142"
		}
	]
	description = "network pool test"
}
`, poolName)
