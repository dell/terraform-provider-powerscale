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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ldapDsMocker *mockey.Mocker
var ldapGetDsMocker *mockey.Mocker

func TestAccLdapProviderDataSource(t *testing.T) {
	var ldapProviderTerraformName = "data.powerscale_ldap_provider.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + ldapProviderAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(ldapProviderTerraformName, "ldap_providers.#"),
				),
			},
			// filter with names read testing
			{
				Config: ProviderConfig + ldapProviderFilterNameDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.#", "1"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.name", "tfacc_ldap"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.server_uris.0", "ldap://10.225.108.54"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.base_dn", "dc=tthe,dc=testLdap,dc=com"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.authentication", "true"),
				),
			},
			// filter with scope read testing
			{
				Config: ProviderConfig + ldapProviderFilterScopeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.#", "1"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.name", "tfacc_ldap"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.server_uris.0", "ldap://10.225.108.54"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.base_dn", "dc=tthe,dc=testLdap,dc=com"),
					resource.TestCheckNoResourceAttr(ldapProviderTerraformName, "ldap_providers.0.authentication"),
				),
			},
		},
	})
}

func TestAccLdapProviderDataSourceInvalidNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// filter with invalid names read testing
			{
				Config:      ProviderConfig + ldapProviderInvalidFilterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error one or more of the filtered LdapProvider names is not a valid powerscale LdapProvider.*.`),
			},
		},
	})
}

func TestAccLdapProviderDatasourceErrorGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock((*client.Client).GetOnefsVersion).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock((*powerscale.AuthApiService).ListAuthv16ProvidersLdapExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
					ldapGetDsMocker = mockey.Mock((*powerscale.AuthApiService).ListAuthv11ProvidersLdapExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccLdapProviderDatasourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock(helper.UpdateLdapProviderDataSourceState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderAllDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccLdapProviderDatasourceHelperMockErr(t *testing.T) {
	var ldapProviderTerraformName = "data.powerscale_ldap_provider.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read v16 Error testing
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock(client.OnefsVersion.IsGreaterThan).Return(true).Build()
					ldapGetDsMocker = mockey.Mock((*powerscale.AuthApiService).ListAuthv16ProvidersLdapExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting the list of PowerScale LdapProviders*.`),
			},
			// Read v11 Error testing
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock(client.OnefsVersion.IsGreaterThan).Return(false).Build()
					ldapGetDsMocker = mockey.Mock((*powerscale.AuthApiService).ListAuthv11ProvidersLdapExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderAllDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting the list of PowerScale LdapProviders*.`),
			},
			// Read v16 testing
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock(client.OnefsVersion.IsGreaterThan).Return(true).Build()
					ldapGetDsMocker = mockey.Mock((*powerscale.AuthApiService).ListAuthv16ProvidersLdapExecute).Return(&mockV16LdapProviders, nil, nil).Build()
				},
				Config: ProviderConfig + ldapProviderAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.#", "1"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.name", "tfacc_ldap"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.server_uris.0", "ldap://10.225.108.54"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.base_dn", "dc=tthe,dc=testLdap,dc=com"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.tls_revocation_check_level", "none"),
				),
			},
			// Read v11 testing
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.UnPatch()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.UnPatch()
					}
					ldapDsMocker = mockey.Mock(client.OnefsVersion.IsGreaterThan).Return(false).Build()
					ldapGetDsMocker = mockey.Mock((*powerscale.AuthApiService).ListAuthv11ProvidersLdapExecute).Return(&mockV11LdapProviders, nil, nil).Build()
				},
				Config: ProviderConfig + ldapProviderAllDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.#", "1"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.name", "tfacc_ldap"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.server_uris.0", "ldap://10.225.108.54"),
					resource.TestCheckResourceAttr(ldapProviderTerraformName, "ldap_providers.0.base_dn", "dc=tthe,dc=testLdap,dc=com"),
					resource.TestCheckNoResourceAttr(ldapProviderTerraformName, "ldap_providers.0.tls_revocation_check_level"),
				),
			},
		},
	})
}
func TestAccLdapProviderDatasourceReleaseMock(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if ldapDsMocker != nil {
						ldapDsMocker.Release()
					}
					if ldapGetDsMocker != nil {
						ldapGetDsMocker.Release()
					}
				},
				Config: ProviderConfig + ldapProviderAllDataSourceConfig,
			},
		},
	})
}

var ldapProviderAllDataSourceConfig = `
data "powerscale_ldap_provider" "test" {
}
`

var ldapProviderFilterNameDataSourceConfig = `
resource "powerscale_ldap_provider" "test" {
	name = "tfacc_ldap"
	base_dn = "dc=tthe,dc=testLdap,dc=com"
	server_uris = ["ldap://10.225.108.54"]
}

data "powerscale_ldap_provider" "test" {
	filter {
    names = [powerscale_ldap_provider.test.name]
  }
  depends_on = [powerscale_ldap_provider.test]
}
`

var ldapProviderFilterScopeDataSourceConfig = `
resource "powerscale_ldap_provider" "test" {
	name = "tfacc_ldap"
	base_dn = "dc=tthe,dc=testLdap,dc=com"
	server_uris = ["ldap://10.225.108.54"]
}

data "powerscale_ldap_provider" "test" {
	filter {
		scope = "user"
		names = [powerscale_ldap_provider.test.name]
  }
  depends_on = [powerscale_ldap_provider.test]
}
`

var ldapProviderInvalidFilterDataSourceConfig = `
data "powerscale_ldap_provider" "test" {
	filter {
    names = ["invalidName"]
  }
}
`
