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

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccS3BucketResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3BucketResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "name", bucketName),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "id", bucketName),
				),
			},
			// ImportState testing
			{
				ResourceName: "powerscale_s3_bucket.bucket_test",
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, bucketName, states[0].Attributes["id"])
					assert.Equal(t, bucketName, states[0].Attributes["name"])
					return nil
				},
			},
			// ImportState testing
			{
				ResourceName:  "powerscale_s3_bucket.bucket_test",
				ImportState:   true,
				ImportStateId: fmt.Sprintf("System:%s", bucketName),
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, bucketName, states[0].Attributes["id"])
					assert.Equal(t, bucketName, states[0].Attributes["name"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + S3BucketUpdatedResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "id", bucketName),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "name", bucketName),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "description", "Updated Description"),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "object_acl_policy", "deny"),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "acl.#", "0"),
				),
			},
		},
	})
}

func TestAccS3BucketResourceErrorRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3BucketResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "name", bucketName),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "id", bucketName),
				),
			},
			// ImportState testing get none bucket
			{
				ResourceName: "powerscale_s3_bucket.bucket_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.GetS3Bucket).Return(&powerscale.V12S3BucketsExtended{}, nil).Build()
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_s3_bucket.bucket_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetS3Bucket).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// ImportState testing get none bucket
			{
				ResourceName: "powerscale_s3_bucket.bucket_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetS3Bucket).Return(&powerscale.V12S3BucketsExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, bucketID string, zone string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile(".not found"),
			},
			// ImportState testing get error
			{
				ResourceName: "powerscale_s3_bucket.bucket_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetS3Bucket).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.UnPatch()
}

func TestAccS3BucketResourceErrorUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3BucketResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "name", bucketName),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "id", bucketName),
				),
			},
			// Update
			{
				Config:      ProviderConfig + S3BucketResourceConfigUpdateZone,
				ExpectError: regexp.MustCompile(".Zone"),
			},
			{
				Config:      ProviderConfig + S3BucketResourceConfigUpdateName,
				ExpectError: regexp.MustCompile(".Name"),
			},
			{
				Config:      ProviderConfig + S3BucketResourceConfigUpdatePath,
				ExpectError: regexp.MustCompile(".Path"),
			},
			// Update get error
			{
				Config: ProviderConfig + S3BucketUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.UpdateS3Bucket).Return(fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
			// Update get none share
			{
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetS3Bucket).Return(&powerscale.V12S3BucketsExtended{}, nil).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				Config:      ProviderConfig + S3BucketUpdatedResourceConfig,
				ExpectError: regexp.MustCompile(".not found"),
			},
			// Update get error
			{
				Config: ProviderConfig + S3BucketUpdatedResourceConfig,
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.GetS3Bucket).Return(nil, fmt.Errorf("mock error")).Build().
						When(func(ctx context.Context, client *client.Client, shareID string, zone string) bool {
							return FunctionMocker.Times() == 2
						})
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.UnPatch()
}

func TestAccS3BucketResourceErrorCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      ProviderConfig + S3BucketInvalidResourceConfig,
				ExpectError: regexp.MustCompile(".*Bad Request*."),
			},
			{
				Config: ProviderConfig + S3BucketResourceConfig,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CreateS3Bucket).Return(nil, fmt.Errorf("mock error")).Build()
				},
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
}

func TestAccS3BucketResourceErrorCopyField(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + S3BucketResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "name", bucketName),
					resource.TestCheckResourceAttr("powerscale_s3_bucket.bucket_test", "id", bucketName),
				),
			},
			{
				ResourceName: "powerscale_s3_bucket.bucket_test",
				ImportState:  true,
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + S3BucketResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			{
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.CopyFieldsToNonNestedModel).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + S3BucketUpdatedResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.UnPatch()
}

func TestAccS3BucketResourceErrorReadState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ReadFromState).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + S3BucketResourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.UnPatch()
}

var bucketName = "tfacc-test-s3-bucket"

var S3BucketResourceConfig = fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_test" {
	name = "%s"
	path = "/ifs/%s"
	create_path = true
	zone = "System"
	acl = [{
		grantee = {
			name = "Everyone"
			type = "wellknown"
		}
		permission = "FULL_CONTROL"
	}]
}
`, bucketName, bucketName)

var S3BucketInvalidResourceConfig = fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_test" {
	name = "%s"
	path = "/ifs/%s"
	create_path = true
	zone = "System"
	acl = [{
		grantee = {
			name = "invalid"
			type = "invalid"
		}
		permission = "FULL_CONTROL"
	}]
}
`, bucketName, bucketName)

var S3BucketUpdatedResourceConfig = fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_test" {
	name = "%s"
	path = "/ifs/%s"
	create_path = true
	zone = "System"
	description = "Updated Description"
	object_acl_policy = "deny"
}
`, bucketName, bucketName)

var S3BucketResourceConfigUpdateName = fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_test" {
	name = "%s-update"
	path = "/ifs/%s"
	create_path = true
	zone = "System"
	acl = [{
		grantee = {
			name = "Everyone"
			type = "wellknown"
		}
		permission = "FULL_CONTROL"
	}]
}
`, bucketName, bucketName)

var S3BucketResourceConfigUpdatePath = fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_test" {
	name = "%s"
	path = "/ifs/%s-update"
	create_path = true
	zone = "System"
	acl = [{
		grantee = {
			name = "Everyone"
			type = "wellknown"
		}
		permission = "FULL_CONTROL"
	}]
}
`, bucketName, bucketName)

var S3BucketResourceConfigUpdateZone = fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_test" {
	name = "%s"
	path = "/ifs/%s"
	create_path = true
	zone = "System-update"
	acl = [{
		grantee = {
			name = "Everyone"
			type = "wellknown"
		}
		permission = "FULL_CONTROL"
	}]
}
`, bucketName, bucketName)
