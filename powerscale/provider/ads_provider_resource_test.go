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

func TestAccAdsProviderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + AdsProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", adsName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", "administrator"),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "lookup_users", "true"),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "check_online_interval", "300"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, adsName, states[0].Attributes["name"])
					assert.Equal(t, "300", states[0].Attributes["check_online_interval"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + AdsProviderUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "lookup_users", "false"),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "check_online_interval", "310"),
				),
			},
		},
	})
}

func TestAccAdsProviderResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + AdsProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", adsName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", "administrator"),
				),
			},
			// ImportState testing get none ads
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetAdsProvider).Return(&powerscale.V14ProvidersAdsExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetAdsProvider).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccAdsProviderResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + AdsProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", adsName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", "administrator"),
				),
			},
			// Update get error
			{
				Config: ProviderConfig + AdsProviderUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = Mock(helper.UpdateAdsProvider).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none ads
			{
				Config: ProviderConfig + AdsProviderUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetAdsProvider).Return(&powerscale.V14ProvidersAdsExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, adsModel models.AdsProviderResourceModel) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".*Could not read updated ads provider*."),
			},
			// Update get error
			{
				Config: ProviderConfig + AdsProviderUpdatedResourceConfig2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.GetAdsProvider).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, adsModel models.AdsProviderResourceModel) bool {
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
				Config:      ProviderConfig + AdsProviderInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccAdsProviderResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + AdsProviderInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccAdsProviderResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + AdsProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", adsName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", "administrator"),
				),
			},
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + AdsProviderResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, source, destination interface{}) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + AdsProviderUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccAdsProviderResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + AdsProviderResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var adsName = "PIE.LAB.EMC.COM"

var AdsProviderResourceConfig = fmt.Sprintf(`
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "administrator"
	password = "Password123!"
}
`, adsName)

var AdsProviderInvalidResourceConfig = fmt.Sprintf(`
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "administrator"
	password = "Password123!"
	sfu_support = "invalid"
}
`, adsName)

var AdsProviderUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "administrator"
	password = "Password123!"
	lookup_users = false
	check_online_interval = 310
}
`, adsName)

var AdsProviderUpdatedResourceConfig2 = fmt.Sprintf(`
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "administrator"
	password = "Password123!"
	check_online_interval = 290
}
`, adsName)
