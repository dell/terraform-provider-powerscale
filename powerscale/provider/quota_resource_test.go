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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuotaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + QuotaResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_quota.quota_test", "linked", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "powerscale_quota.quota_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"persona", "ignore_limit_checks", "zone"},
			},
			// Update
			{
				Config: ProviderConfig + QuotaResourceConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_quota.quota_test", "thresholds.percent_advisory", "25.4"),
				),
			},
		},
	})
}

func TestAccQuotaResourceCreateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + QuotaResourceCreateWithLink,
				ExpectError: regexp.MustCompile(".set attribute Linked while creating"),
			},
			{
				Config:      ProviderConfig + QuotaResourceCreateDirectoryTypeWithPersona,
				ExpectError: regexp.MustCompile("\"persona\" is not needed"),
			},
			{
				Config:      ProviderConfig + QuotaResourceCreateUserTypeWithoutPersona,
				ExpectError: regexp.MustCompile("\"persona\" is required"),
			},
			// Update
			{
				Config: ProviderConfig + QuotaResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CreateQuota).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".error"),
			},
			{
				Config: ProviderConfig + QuotaResourceConfig,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".error"),
			},
		},
	})
}

func TestAccQuotaResourceUpdateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + QuotaResourceConfig,
			},
			// Update
			{
				Config:      ProviderConfig + QuotaResourceConfigNewZone,
				ExpectError: regexp.MustCompile(".Zone"),
			},
			{
				Config:      ProviderConfig + QuotaResourceConfigNewPersona,
				ExpectError: regexp.MustCompile(".Persona"),
			},
			{
				Config:      ProviderConfig + QuotaResourceConfigNewPath,
				ExpectError: regexp.MustCompile(".Path"),
			},
			{
				Config:      ProviderConfig + QuotaResourceConfigNewType,
				ExpectError: regexp.MustCompile(".Type"),
			},
			{
				Config:      ProviderConfig + QuotaResourceConfigDoesIncludeSnapshots,
				ExpectError: regexp.MustCompile(".IncludeSnapshots"),
			},
			{
				Config: ProviderConfig + QuotaResourceConfigUpdated,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateQuota).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".mock error"),
			},
			{
				Config: ProviderConfig + QuotaResourceConfigUpdated2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.ValidateQuotaUpdate).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".error"),
			},
			{
				Config: ProviderConfig + QuotaResourceConfigUpdated,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".error"),
			},
			{
				Config: ProviderConfig + QuotaResourceConfigUpdated2,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetQuota).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, quotaID string, zone string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".error"),
			},
			{
				Config: ProviderConfig + QuotaResourceConfigUpdated,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetQuota).Return(&powerscale.V12QuotaQuotasExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, quotaID string, zone string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			{
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + QuotaResourceConfigUpdated2,
				ExpectError: regexp.MustCompile(".error"),
			},
		},
	})
}

func TestAccQuotaResourceReadError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + QuotaResourceConfig,
			},
			{
				ResourceName: "powerscale_quota.quota_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetQuota).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".mock error"),
			},
			{
				ResourceName: "powerscale_quota.quota_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.GetQuota).Return(&powerscale.V12QuotaQuotasExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			{
				ResourceName: "powerscale_quota.quota_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.Release()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".error"),
			},
		},
	})
}

func TestAccQuotaResourceUpdateLinkError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + QuotaResourceUnLink,
			},
			{
				Config: ProviderConfig + QuotaResourceLink,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.LinkQuota).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile(".mock error"),
			},
		},
	})
}

var QuotaResourceConfig = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceConfigNewZone = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "AppZone"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceConfigDoesIncludeSnapshots = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = true
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceConfigNewType = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "directory"
	include_snapshots = false
	zone = "System"
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceConfigNewPath = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test2"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`
var QuotaResourceConfigNewPersona = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1502"
		name = "Guest2"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceCreateWithLink = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "directory"
	include_snapshots = false
	linked = true
}
`

var QuotaResourceConfigUpdated = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 25.4
        percent_soft = 50
		soft_grace = 120
	    hard = 4000
	}
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
	ignore_limit_checks = true
}
`

var QuotaResourceConfigUpdated2 = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 25.4
        percent_soft = 50
		soft_grace = 120
	    hard = 4000
	}
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
	ignore_limit_checks = true
}
`

var QuotaResourceCreateDirectoryTypeWithPersona = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "directory"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceCreateUserTypeWithoutPersona = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`

var QuotaResourceLink = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
	linked = true
}
`

var QuotaResourceUnLink = `
resource "powerscale_quota" "quota_test" {
	path = "/ifs/tfacc_quota_test"
	type = "user"
	include_snapshots = false
	zone = "System"
	persona = {
		id = "UID:1501"
		name = "Guest"
		type = "user"
	}
	thresholds = {
		percent_advisory = 10
        percent_soft = 20
		soft_grace = 120
	    hard = 4000
	}
	ignore_limit_checks = true
	container = true
	enforced = false
	thresholds_on = "applogicalsize"
}
`
