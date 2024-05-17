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

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNfsGlobalSettingsDataSource(t *testing.T) {
	var nfsGlobalSettings = "data.powerscale_nfs_global_settings.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read all testing
			{
				Config: ProviderConfig + globalSettingsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "id"),
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "nfsv3_enabled"),
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "nfsv3_rdma_enabled"),
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "nfsv4_enabled"),
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "rquota_enabled"),
					resource.TestCheckResourceAttrSet(nfsGlobalSettings, "service"),
				),
			},
		},
	})
}

func TestAccNfsGlobalSettingsDataSourceErrorGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(helper.GetNfsGlobalSettings).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + globalSettingsDataSourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

var globalSettingsDataSourceConfig = `
data "powerscale_nfs_global_settings" "test" {
}
`
