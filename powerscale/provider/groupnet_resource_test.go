/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var groupnetGetMocker *Mocker
var groupnetMocker *Mocker

func TestAccGroupnetResourceCreate(t *testing.T) {
	var groupnetResourceName = "powerscale_groupnet.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + groupnetBasicResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetResourceName, "name", "tfaccGroupnetCreation"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_resolver_rotate", "false"),
					resource.TestCheckResourceAttr(groupnetResourceName, "allow_wildcard_subdomains", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "server_side_dns_search", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_cache_enabled", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "id", "tfaccGroupnetCreation"),
					resource.TestCheckNoResourceAttr(groupnetResourceName, "description"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_search.#", "0"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_servers.#", "0"),
					resource.TestCheckResourceAttr(groupnetResourceName, "subnets.#", "0"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + groupnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetResourceName, "name", "tfaccGroupnetUpdate"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_resolver_rotate", "false"),
					resource.TestCheckResourceAttr(groupnetResourceName, "allow_wildcard_subdomains", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "server_side_dns_search", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_cache_enabled", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "id", "tfaccGroupnetUpdate"),
					resource.TestCheckResourceAttr(groupnetResourceName, "description", "terraform groupnet resource"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_search.#", "1"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_search.0", "pie.lab.emc.com"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_servers.#", "1"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_servers.0", "10.230.44.169"),
					resource.TestCheckResourceAttr(groupnetResourceName, "subnets.#", "0"),
				),
			},
			{
				// Update and Read testing
				Config: ProviderConfig + groupnetUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetResourceName, "name", "tfaccGroupnetUpdate2"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_resolver_rotate", "true"),
					resource.TestCheckResourceAttr(groupnetResourceName, "allow_wildcard_subdomains", "false"),
					resource.TestCheckResourceAttr(groupnetResourceName, "server_side_dns_search", "false"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_cache_enabled", "false"),
					resource.TestCheckResourceAttr(groupnetResourceName, "id", "tfaccGroupnetUpdate2"),
					resource.TestCheckNoResourceAttr(groupnetResourceName, "description"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_search.#", "0"),
					resource.TestCheckResourceAttr(groupnetResourceName, "dns_servers.#", "0"),
					resource.TestCheckResourceAttr(groupnetResourceName, "subnets.#", "0"),
				),
			},
		},
	})
}

func TestAccGroupnetResourceErr(t *testing.T) {
	var groupnetResourceName = "powerscale_groupnet.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error - invalid name
			{
				Config:      ProviderConfig + groupnetInvalidNameResourceConfig,
				ExpectError: regexp.MustCompile(`.*must follow the pattern*.`),
			},
			// Create Error - invalid server
			{
				Config:      ProviderConfig + groupnetInvalidServerResourceConfig,
				ExpectError: regexp.MustCompile(`.*error creating groupnet*.`),
			},
			// Create Error - invalid search
			{
				Config:      ProviderConfig + groupnetInvalidSearchResourceConfig,
				ExpectError: regexp.MustCompile(`.*error creating groupnet*.`),
			},
			// Create and Read testing
			{
				Config: ProviderConfig + groupnetBasicResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetResourceName, "name", "tfaccGroupnetCreation"),
				),
			},
			// Update Error - invalid name
			{
				Config:      ProviderConfig + groupnetInvalidNameResourceConfig,
				ExpectError: regexp.MustCompile(`.*must follow the pattern*.`),
			},
			// Update Error - invalid server
			{
				Config:      ProviderConfig + groupnetInvalidServerResourceConfig,
				ExpectError: regexp.MustCompile(`.*error updating groupnet*.`),
			},
			// Update Error - invalid search
			{
				Config:      ProviderConfig + groupnetInvalidSearchResourceConfig,
				ExpectError: regexp.MustCompile(`.*error updating groupnet*.`),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + groupnetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetResourceName, "name", "tfaccGroupnetUpdate"),
				),
			},
		},
	})
}

func TestAccGroupnetResourceImport(t *testing.T) {
	var groupnetResourceName = "powerscale_groupnet.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + groupnetResourceConfig,
			},
			{
				ResourceName:      groupnetResourceName,
				ImportState:       true,
				ExpectError:       nil,
				ImportStateVerify: true,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					assert.Equal(t, "tfaccGroupnetUpdate", s[0].Attributes["name"])
					assert.Equal(t, "1", s[0].Attributes["dns_search.#"])
					assert.Equal(t, "pie.lab.emc.com", s[0].Attributes["dns_search.0"])
					assert.Equal(t, "1", s[0].Attributes["dns_servers.#"])
					assert.Equal(t, "10.230.44.169", s[0].Attributes["dns_servers.0"])
					return nil
				},
			},
		},
	})
}

func TestAccGroupnetResourceCreateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Error testing
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
					groupnetMocker = Mock(helper.CreateGroupnet).Return(nil).Build()
					groupnetGetMocker = Mock(helper.GetGroupnet).Return(nil, fmt.Errorf("groupnet read mock error")).Build()
				},
				Config:      ProviderConfig + groupnetBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*groupnet read mock error*.`),
			},
			// Create and Parse state Error testing
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
					groupnetGetMocker = Mock(helper.UpdateGroupnetResourceState).Return(fmt.Errorf("groupnet parse state mock error")).Build()
				},
				Config:      ProviderConfig + groupnetBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*groupnet parse state mock error*.`),
			},
		},
	})
}

func TestAccGroupnetResourceUpdateMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
				},
				Config: ProviderConfig + groupnetBasicResourceConfig,
			},
			// Update and Read Error testing
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
					groupnetMocker = Mock(helper.UpdateGroupnet).Return(nil).Build()
					groupnetGetMocker = Mock(helper.GetGroupnet).Return(nil, fmt.Errorf("groupnet read mock error")).Build()
				},
				Config:      ProviderConfig + groupnetBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*groupnet read mock error*.`),
			},
			// Update and Parse state Error testing
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
					groupnetGetMocker = Mock(helper.UpdateGroupnetResourceState).Return(fmt.Errorf("groupnet parse state mock error")).Build()
				},
				Config:      ProviderConfig + groupnetBasicResourceConfig,
				ExpectError: regexp.MustCompile(`.*groupnet parse state mock error*.`),
			},
		},
	})
}

func TestAccGroupnetResourceDeleteMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
				},
				Config: ProviderConfig + groupnetBasicResourceConfig,
			},
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.UnPatch()
					}
					if groupnetMocker != nil {
						groupnetMocker.UnPatch()
					}
					groupnetMocker = Mock(helper.DeleteGroupnet).Return(fmt.Errorf("groupnet delete mock error")).Build().
						When(func(ctx context.Context, client *client.Client, groupnetName string) bool {
							return groupnetMocker.Times() == 1
						})
				},
				Config:      ProviderConfig + groupnetBasicResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("groupnet delete mock error"),
			},
		},
	})
}

func TestAccGroupnetReleaseMockResource(t *testing.T) {
	var groupnetResourceName = "powerscale_groupnet.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if groupnetGetMocker != nil {
						groupnetGetMocker.Release()
					}
					if groupnetMocker != nil {
						groupnetMocker.Release()
					}
				},
				Config: ProviderConfig + groupnetBasicResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(groupnetResourceName, "name", "tfaccGroupnetCreation"),
				)},
		},
	})
}

var groupnetBasicResourceConfig = `
resource "powerscale_groupnet" "test" {
	name = "tfaccGroupnetCreation"
  }
`

var groupnetResourceConfig = `
resource "powerscale_groupnet" "test" {
	name = "tfaccGroupnetUpdate"
	description = "terraform groupnet resource"
	dns_cache_enabled = true
	allow_wildcard_subdomains = true
	server_side_dns_search = true
	dns_resolver_rotate = false
	dns_search = ["pie.lab.emc.com"]
	dns_servers = ["10.230.44.169"]
  }
`

var groupnetUpdateResourceConfig = `
resource "powerscale_groupnet" "test" {
	name = "tfaccGroupnetUpdate2"
	dns_cache_enabled = false
	allow_wildcard_subdomains = false
	server_side_dns_search = false
	dns_resolver_rotate = true
  }
`

var groupnetInvalidNameResourceConfig = `
resource "powerscale_groupnet" "test" {
	name = "&&invalidName"
  }
`

var groupnetInvalidServerResourceConfig = `
resource "powerscale_groupnet" "test" {
	name = "tfaccGroupnetCreation"
	dns_servers = ["invalidServer"]
  }
`

var groupnetInvalidSearchResourceConfig = `
resource "powerscale_groupnet" "test" {
	name = "tfaccGroupnetCreation"
	dns_search = ["_invalid_"]
  }
`
