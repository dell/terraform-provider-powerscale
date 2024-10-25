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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
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
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", powerscaleAdsproviderName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", powerscaleAdsproviderUsername),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "lookup_users", "true"),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "check_online_interval", "300"),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, powerscaleAdsproviderName, states[0].Attributes["name"])
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
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", powerscaleAdsproviderName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", powerscaleAdsproviderUsername),
				),
			},
			// ImportState testing get none ads
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetAdsProvider).Return(&powerscale.V14ProvidersAdsExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetAdsProvider).Return(nil, fmt.Errorf("mock error")).Build()
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
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", powerscaleAdsproviderName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", powerscaleAdsproviderUsername),
				),
			},
			// Update get error
			{
				Config: ProviderConfig + AdsProviderUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateAdsProvider).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none ads
			{
				Config: ProviderConfig + AdsProviderUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetAdsProvider).Return(&powerscale.V14ProvidersAdsExtended{}, nil).Build().
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
					FunctionMocker = mockey.Mock(helper.GetAdsProvider).Return(nil, fmt.Errorf("mock error")).Build().
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
			{
				Config:      ProviderConfig + AdsProviderUpdatePreCheckConfig,
				ExpectError: regexp.MustCompile(".*Should not provide parameters for creating*."),
			},
			{
				Config:      ProviderConfig + AdsProviderUpdateGroupnetConfig,
				ExpectError: regexp.MustCompile(".*Should not use a different groupnet*."),
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
			{
				Config:      ProviderConfig + AdsProviderCreatePreCheckConfig,
				ExpectError: regexp.MustCompile(".*Should not provide parameters for updating*."),
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
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "name", powerscaleAdsproviderName),
					resource.TestCheckResourceAttr("powerscale_adsprovider.ads_test", "user", powerscaleAdsproviderUsername),
				),
			},
			{
				ResourceName: "powerscale_adsprovider.ads_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				// Config:      ProviderConfig + AdsProviderResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build().
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
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + AdsProviderResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var AdsProviderResourceConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
}
`

var AdsProviderInvalidResourceConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
	sfu_support = "invalid"
}
`

var AdsProviderCreatePreCheckConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
	reset_schannel = true
}
`

var AdsProviderUpdatedResourceConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
	lookup_users = false
	check_online_interval = 310
}
`

var AdsProviderUpdatedResourceConfig2 = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
	check_online_interval = 290
}
`

var AdsProviderUpdatePreCheckConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
	kerberos_hdfs_spn = true
}
`

var AdsProviderUpdateGroupnetConfig = `
resource "powerscale_adsprovider" "ads_test" {
	name = "%s"
	user = "%s"
	password = "%s"
	groupnet = "groupnet_x"
}
`

func initAdsProviderConfig() {
	// resource config
	AdsProviderResourceConfig = fmt.Sprintf(AdsProviderResourceConfig, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)
	AdsProviderInvalidResourceConfig = fmt.Sprintf(AdsProviderInvalidResourceConfig, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)
	AdsProviderCreatePreCheckConfig = fmt.Sprintf(AdsProviderCreatePreCheckConfig, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)
	AdsProviderUpdatedResourceConfig = fmt.Sprintf(AdsProviderUpdatedResourceConfig, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)
	AdsProviderUpdatedResourceConfig2 = fmt.Sprintf(AdsProviderUpdatedResourceConfig2, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)
	AdsProviderUpdatePreCheckConfig = fmt.Sprintf(AdsProviderUpdatePreCheckConfig, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)
	AdsProviderUpdateGroupnetConfig = fmt.Sprintf(AdsProviderUpdateGroupnetConfig, powerscaleAdsproviderName, powerscaleAdsproviderUsername, powerscaleAdsproviderPassword)

	// data source config
	// All datasources are appended with AdsProviderResourceConfig as pre-requirement
	AdsDataSourceNamesConfig = AdsProviderResourceConfig + fmt.Sprintf(AdsDataSourceNamesConfig, powerscaleAdsproviderName)
	AdsDataSourceFilterConfig = AdsProviderResourceConfig + AdsDataSourceFilterConfig
	AdsAllDataSourceConfig = AdsProviderResourceConfig + AdsAllDataSourceConfig
	AdsDataSourceNameConfigErr = AdsProviderResourceConfig + AdsDataSourceNameConfigErr
	AdsDataSourceFilterConfigErr = AdsProviderResourceConfig + AdsDataSourceFilterConfigErr

}
