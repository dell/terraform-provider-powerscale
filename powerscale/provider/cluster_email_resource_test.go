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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccClusterEmailResourceImport(t *testing.T) {
	var clusterEmailResourceName = "powerscale_cluster_email.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterEmailResourceConfig,
			},
			// Import testing
			{
				Config:       ProviderConfig + clusterEmailResourceConfig,
				ResourceName: clusterEmailResourceName,
				ImportState:  true,
				ExpectError:  nil,
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "id")
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "settings")
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "settings.batch_mode")
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "settings.smtp_auth_passwd_set")
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "settings.smtp_auth_security")
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "settings.smtp_port")
					resource.TestCheckResourceAttrSet(clusterEmailResourceName, "settings.use_smtp_auth")
					return nil
				},
			},
		},
	})
}

func TestAccClusterEmailResourceNullableField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterEmailResourceConfig,
			},
			{
				Config: ProviderConfig + clusterEmailUserTemplateSetResourceConfig,
			},
			{
				Config: ProviderConfig + clusterEmailUserTemplateEmptyResourceConfig,
			},
		},
	})
}

func TestAccClusterEmailResourceUpdate(t *testing.T) {
	var clusterEmailResourceName = "powerscale_cluster_email.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterEmailResourceConfig,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + clusterEmailUpdateResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterEmailResourceName, "settings.batch_mode", "all"),
					resource.TestCheckResourceAttr(clusterEmailResourceName, "settings.smtp_auth_security", "none"),
					resource.TestCheckResourceAttr(clusterEmailResourceName, "settings.smtp_port", "6225"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + clusterEmailUpdateRevertResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterEmailResourceName, "settings.batch_mode", "none"),
					resource.TestCheckResourceAttr(clusterEmailResourceName, "settings.smtp_auth_security", "starttls"),
					resource.TestCheckResourceAttr(clusterEmailResourceName, "settings.smtp_port", "25"),
				),
			},
		},
	})
}

func TestAccClusterEmailResourceCreateMockErr(t *testing.T) {
	var FunctionMocker2 *mockey.Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterEmailResourceConfig,
			},
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
					FunctionMocker2 = mockey.Mock(helper.GetV21ClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					if FunctionMocker2 != nil {
						FunctionMocker2.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterEmail).Return(fmt.Errorf("mock error")).Build()
					FunctionMocker2 = mockey.Mock(helper.UpdateV21ClusterEmail).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailUpdateResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
	FunctionMocker2.Release()
}

func TestAccClusterEmailResourceUpdateMockErr(t *testing.T) {
	var FunctionMocker2 *mockey.Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
					FunctionMocker2 = mockey.Mock(helper.GetV21ClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					if FunctionMocker2 != nil {
						FunctionMocker2.Release()
					}
					FunctionMocker = mockey.Mock(helper.UpdateClusterEmail).Return(fmt.Errorf("mock error")).Build()
					FunctionMocker2 = mockey.Mock(helper.UpdateV21ClusterEmail).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + clusterEmailResourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
	FunctionMocker2.Release()
}

func TestAccClusterEmailResourceImportMockErr(t *testing.T) {
	var FunctionMocker2 *mockey.Mocker
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + clusterEmailResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if FunctionMocker != nil {
						FunctionMocker.Release()
					}
					if FunctionMocker2 != nil {
						FunctionMocker2.Release()
					}
					FunctionMocker = mockey.Mock(helper.GetClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
					FunctionMocker2 = mockey.Mock(helper.GetV21ClusterEmail).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:            ProviderConfig + clusterEmailResourceConfig,
				ResourceName:      "powerscale_cluster_email.test",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
	FunctionMocker2.Release()
}

var clusterEmailResourceConfig = `
resource "powerscale_cluster_email" "test" {
	settings = {
	}
}
`

var clusterEmailUserTemplateSetResourceConfig = `
resource "powerscale_cluster_email" "test" {
	settings = {
		user_template = "/ifs/README.txt"
	}
}
`

var clusterEmailUserTemplateEmptyResourceConfig = `
resource "powerscale_cluster_email" "test" {
	settings = {
		user_template = ""
	}
}
`

var clusterEmailUpdateResourceConfig = `
resource "powerscale_cluster_email" "test" {
	settings = {
		batch_mode = "all"
		smtp_auth_security = "none"
		smtp_port = 6225
	}
	
}
`

var clusterEmailUpdateRevertResourceConfig = `
resource "powerscale_cluster_email" "test" {
	settings = {
		batch_mode = "none"
		smtp_auth_security = "starttls"
		smtp_port = 25
	}
}
`
