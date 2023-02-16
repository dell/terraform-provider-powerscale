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
	"testing"
)

func TestAccExampleResource(t *testing.T) {
	//resource.Test(t, resource.TestCase{
	//	PreCheck:                 func() { testAccPreCheck(t) },
	//	ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
	//	Steps: []resource.TestStep{
	//		// Create and Read testing
	//		{
	//			Config: testAccExampleResourceConfig("one"),
	//			Check: resource.ComposeAggregateTestCheckFunc(
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "configurable_attribute", "one"),
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "defaulted", "example value when not configured"),
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "id", "example-id"),
	//			),
	//		},
	//		// ImportState testing
	//		{
	//			ResourceName:      "scaffolding_example.test",
	//			ImportState:       true,
	//			ImportStateVerify: true,
	//			// This is not normally necessary, but is here because this
	//			// example code does not have an actual upstream service.
	//			// Once the Read method is able to refresh information from
	//			// the upstream service, this can be removed.
	//			ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
	//		},
	//		// Update and Read testing
	//		{
	//			Config: testAccExampleResourceConfig("two"),
	//			Check: resource.ComposeAggregateTestCheckFunc(
	//				resource.TestCheckResourceAttr("scaffolding_example.test", "configurable_attribute", "two"),
	//			),
	//		},
	//		// Delete testing automatically occurs in TestCase
	//	},
	//})
}
