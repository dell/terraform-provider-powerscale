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
	powerscale "dell/powerscale-go-client"
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

var ldapMocker *Mocker
var ldapV16Mocker *Mocker
var ldapV11Mocker *Mocker

func TestAccLdapProviderResource(t *testing.T) {
	var ldapResourceName = "powerscale_ldap_provider.ldap_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + ldapProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapResourceName, "name", "tfacc_ldap"),
					resource.TestCheckResourceAttr(ldapResourceName, "zone_name", "System"),
					resource.TestCheckResourceAttr(ldapResourceName, "base_dn", "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com"),
					resource.TestCheckResourceAttr(ldapResourceName, "server_uris.0", powerscaleLdapHost),
				),
			},
			// ImportState testing
			{
				ResourceName: ldapResourceName,
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "tfacc_ldap", states[0].Attributes["name"])
					assert.Equal(t, "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com", states[0].Attributes["base_dn"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + ldapProviderResourceRenameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapResourceName, "name", "tfacc_ldap_update"),
				),
			},
		},
	})
}

func TestAccLdapProviderResourceErr(t *testing.T) {
	var ldapResourceName = "powerscale_ldap_provider.ldap_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error - invalid server
			{
				Config:      ProviderConfig + ldapProviderInvalidResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating ldap provider*.`),
			},
			// Create and Read testing
			{
				Config: ProviderConfig + ldapProviderResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapResourceName, "name", "tfacc_ldap"),
				),
			},
			// Update Error - invalid server
			{
				Config:      ProviderConfig + ldapProviderInvalidResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error updating the LdapProvider resource*.`),
			},
			// Update
			{
				Config: ProviderConfig + ldapProviderResourceRenameConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ldapResourceName, "name", "tfacc_ldap_update"),
				),
			},
		},
	})
}

func TestAccLdapProviderResourceMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapV16Mocker = Mock((*powerscale.AuthApiService).CreateAuthv16ProvidersLdapItemExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
					ldapV11Mocker = Mock((*powerscale.AuthApiService).CreateAuthv11ProvidersLdapItemExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Create and GetVersion Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock((*client.Client).GetOnefsVersion).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Create and Read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(helper.GetLdapProvider).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Create and Parse state Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapV16Mocker = Mock(helper.UpdateLdapProviderResourceState).Return(fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + ldapProviderResourceConfig,
			},
			// Read and GetVersion Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock((*client.Client).GetOnefsVersion).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(helper.GetLdapProvider).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(helper.UpdateLdapProviderResourceState).Return(fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Update and GetVersion Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock((*client.Client).GetOnefsVersion).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Update Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapV16Mocker = Mock((*powerscale.AuthApiService).UpdateAuthv16ProvidersLdapByIdExecute).Return(nil, fmt.Errorf("ldap mock error")).Build()
					ldapV11Mocker = Mock((*powerscale.AuthApiService).UpdateAuthv11ProvidersLdapByIdExecute).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Update and Read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(helper.GetLdapProvider).When(func(ctx context.Context, client *client.Client, ldapProviderName, scope string) bool {
						return ldapMocker.Times() > 1
					}).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + ldapProviderResourceDisableConfig,
			},
		},
	})
}

func TestAccLdapProviderResourceDeleteMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + ldapProviderResourceConfig,
			},
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock((*powerscale.AuthApiService).DeleteAuthv11ProvidersLdapByIdExecute).When(func(r powerscale.ApiDeleteAuthv11ProvidersLdapByIdRequest) bool {
						return ldapMocker.Times() == 1
					}).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("ldap mock error"),
			},
		},
	})
}

func TestAccLdapProviderResourceImportMockErr(t *testing.T) {
	var ldapResourceName = "powerscale_ldap_provider.ldap_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + ldapProviderResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(helper.GetLdapProvider).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:            ProviderConfig + ldapProviderResourceConfig,
				ResourceName:      ldapResourceName,
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*ldap mock error*.`),
				ImportStateVerify: true,
			},
			// Import and parse Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(helper.UpdateLdapProviderResourceState).Return(fmt.Errorf("ldap mock error")).Build()
				},
				Config:            ProviderConfig + ldapProviderResourceConfig,
				ResourceName:      ldapResourceName,
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*ldap mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLdapProviderResourceHelperMockErr(t *testing.T) {
	mockV16LdapProviders, mockV11LdapProviders := getMockLdapProviderConfig()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create v16 Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(client.OnefsVersion.IsGreaterThan).Return(true).Build()
					ldapV16Mocker = Mock((*powerscale.AuthApiService).CreateAuthv16ProvidersLdapItemExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Create v11 Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(client.OnefsVersion.IsGreaterThan).Return(false).Build()
					ldapV16Mocker = Mock((*powerscale.AuthApiService).CreateAuthv11ProvidersLdapItemExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Create v16 and Read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(client.OnefsVersion.IsGreaterThan).Return(true).Build()
					ldapV11Mocker = Mock((*powerscale.AuthApiService).CreateAuthv16ProvidersLdapItemExecute).Return(nil, nil, nil).Build()
					ldapV16Mocker = Mock((*powerscale.AuthApiService).GetAuthv16ProvidersLdapByIdExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting ldap provider after creation*.`),
			},
			// Create v11 and Read Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(client.OnefsVersion.IsGreaterThan).Return(false).Build()
					ldapV11Mocker = Mock((*powerscale.AuthApiService).CreateAuthv11ProvidersLdapItemExecute).Return(nil, nil, nil).Build()
					ldapV16Mocker = Mock((*powerscale.AuthApiService).GetAuthv11ProvidersLdapByIdExecute).Return(nil, nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error getting ldap provider after creation*.`),
			},
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + ldapProviderResourceConfig,
			},
			// Update v16 Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(client.OnefsVersion.IsGreaterThan).Return(true).Build()
					ldapV11Mocker = Mock((*powerscale.AuthApiService).GetAuthv16ProvidersLdapByIdExecute).Return(&mockV16LdapProviders, nil, nil).Build()
					ldapV16Mocker = Mock((*powerscale.AuthApiService).UpdateAuthv16ProvidersLdapByIdExecute).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			// Update v11 Error testing
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
					ldapMocker = Mock(client.OnefsVersion.IsGreaterThan).Return(false).Build()
					ldapV11Mocker = Mock((*powerscale.AuthApiService).GetAuthv11ProvidersLdapByIdExecute).Return(&mockV11LdapProviders, nil, nil).Build()
					ldapV16Mocker = Mock((*powerscale.AuthApiService).UpdateAuthv11ProvidersLdapByIdExecute).Return(nil, fmt.Errorf("ldap mock error")).Build()
				},
				Config:      ProviderConfig + ldapProviderResourceDisableConfig,
				ExpectError: regexp.MustCompile(`.*ldap mock error*.`),
			},
			{
				PreConfig: func() {
					if ldapV16Mocker != nil {
						ldapV16Mocker.UnPatch()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.UnPatch()
					}
					if ldapMocker != nil {
						ldapMocker.UnPatch()
					}
				},
				Config: ProviderConfig + ldapProviderResourceDisableConfig,
			},
		},
	})
}

func TestAccLdapProviderReleaseMockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if ldapMocker != nil {
						ldapMocker.Release()
					}
					if ldapV16Mocker != nil {
						ldapV16Mocker.Release()
					}
					if ldapV11Mocker != nil {
						ldapV11Mocker.Release()
					}
				},
				Config: ProviderConfig + ldapProviderResourceConfig,
			},
		},
	})
}

var ldapProviderResourceConfig = `
resource "powerscale_ldap_provider" "ldap_test" {
	name = "tfacc_ldap"
	server_uris = ["%s"]
	base_dn = "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com"
}
`

var ldapProviderResourceRenameConfig = `
resource "powerscale_ldap_provider" "ldap_test" {
	name = "tfacc_ldap_update"
	server_uris = ["%s"]
	base_dn = "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com"
}
`

var ldapProviderResourceDisableConfig = `
resource "powerscale_ldap_provider" "ldap_test" {
	name = "tfacc_ldap"
	server_uris = ["%s"]
	base_dn = "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com"
	enabled = false
}
`

var ldapProviderInvalidResourceConfig = `
resource "powerscale_ldap_provider" "ldap_test" {
	name = "tfacc_ldap"
	server_uris = ["ldap://10.10.10.xx"]
	base_dn = "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com"
}
`

func initLdapVars() {
	// resource config
	ldapProviderResourceConfig = fmt.Sprintf(ldapProviderResourceConfig, powerscaleLdapHost)
	ldapProviderResourceRenameConfig = fmt.Sprintf(ldapProviderResourceRenameConfig, powerscaleLdapHost)
	ldapProviderResourceDisableConfig = fmt.Sprintf(ldapProviderResourceDisableConfig, powerscaleLdapHost)

	// datasource config
	ldapProviderFilterNameDataSourceConfig = ldapProviderResourceConfig + ldapProviderFilterNameDataSourceConfig
	ldapProviderFilterScopeDataSourceConfig = ldapProviderResourceConfig + ldapProviderFilterScopeDataSourceConfig
}

func getMockLdapProviderConfig() (powerscale.V16ProvidersLdap, powerscale.V11ProvidersLdap) {
	mockName, mockBaseDN, mockTLSRevocationCheckLevel := "tfacc_ldap", "dc=yulan,dc=pie,dc=lab,dc=emc,dc=com", "none"
	mockServerUris := []string{powerscaleLdapHost}
	return powerscale.V16ProvidersLdap{
			Ldap: []powerscale.V16ProvidersLdapLdapItem{
				{
					Name:                    &mockName,
					BaseDn:                  &mockBaseDN,
					ServerUris:              mockServerUris,
					TlsRevocationCheckLevel: &mockTLSRevocationCheckLevel,
				},
			},
		}, powerscale.V11ProvidersLdap{
			Ldap: []powerscale.V11ProvidersLdapLdapItem{
				{
					Name:       &mockName,
					BaseDn:     &mockBaseDN,
					ServerUris: mockServerUris,
				},
			},
		}
}
