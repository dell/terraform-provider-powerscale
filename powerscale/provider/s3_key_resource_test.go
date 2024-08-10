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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccS3KeyResourceErrorCreate(t *testing.T) {
	var S3KeyResourceConfigCreateError = tfConfig("tf_err_test", "invalid", "invalid", 80)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + S3KeyResourceConfigCreateError,
				ExpectError: regexp.MustCompile(".*Error creating s3 key*."),
			},
		},
	})
}

func TestAccS3KeyResource(t *testing.T) {
	var S3KeyResourceConfig = tfConfig("tf_test", "tf_user", "System", 40)
	var S3KeyResourceConfigUpdate = tfConfig("tf_test", "tf_user", "System", 80)
	var S3KeyResourceConfigUpdateError = tfConfig("tf_test", "tf_user", "System", -80)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{

				Config: ProviderConfig + S3KeyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.tf_test", "user", "tf_user"),
				),
			},
			// Update
			{
				Config: ProviderConfig + S3KeyResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.tf_test", "user", "tf_user"),
				),
			},
			// Update Error testing
			{
				Config:      ProviderConfig + S3KeyResourceConfigUpdateError,
				ExpectError: regexp.MustCompile(".*Error updating s3 key*."),
			},
		},
	})
}

func tfConfig(resource, user, zone string, expiry int) string {
	return fmt.Sprintf(`
resource "powerscale_s3_key" "%s" {
    user = "%s"
    zone = "%s"
    existing_key_expiry_time = %d
}
`, resource, user, zone, expiry)
}
