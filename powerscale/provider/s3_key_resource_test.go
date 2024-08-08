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

func TestAccS3KeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3KeyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.key_test", "id", "123"),
				),
			},
			// Update
			{
				Config: ProviderConfig + S3KeyResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.key_test", "id", "123"),
				),
			},
		},
	})
}

func TestAccS3KeyResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3KeyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.key_test", "id", "123"),
				),
			},
		},
	})
}

func TestAccS3KeyResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3KeyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.key_test", "name", "123"),
					resource.TestCheckResourceAttr("powerscale_s3_key.key_test", "id", "123"),
				),
			},
			// Update
			{
				Config:      ProviderConfig + S3KeyResourceConfigUpdate,
				ExpectError: regexp.MustCompile(".key"),
			},
		},
	})
}

func TestAccS3KeyResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + S3KeyInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
		},
	})
}

func TestAccS3KeyResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3KeyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_key.key_test", "id", "123"),
				),
			},
		},
	})
}

var S3KeyResourceConfig = fmt.Sprintf(`
resource "powerscale_s3_key" "key_test" {
}
`)

var S3KeyInvalidResourceConfig = fmt.Sprintf(`
resource "powerscale_s3_key" "key_test" {
}
`)

var S3KeyResourceConfigUpdate = fmt.Sprintf(`
resource "powerscale_s3_key" "key_test" {
}
`)
