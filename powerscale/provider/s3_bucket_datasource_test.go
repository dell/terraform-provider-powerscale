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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccS3BucketDatasource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + S3BucketDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_s3_bucket.bucket_datasource_test", "s3_buckets.#", "1"),
				),
			},
		},
	})
}

func TestAccS3BucketsourceGetAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				Config: ProviderConfig + S3BucketAllDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_s3_bucket.bucket_datasource_test_all", "filter.#", "0"),
				),
			},
		},
	})
}

func TestAccS3BucketDatasourceGetError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.ListS3Buckets).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + S3BucketAllDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.UnPatch()
}

func TestAccS3BucketDatasourcePagination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					resume := "1"
					buckets := powerscale.V12S3Buckets{
						Resume:  &resume,
						Buckets: nil,
						Total:   nil,
					}
					FunctionMocker = mockey.Mock(mockey.GetMethod(powerscale.ApiListProtocolsv12S3BucketsRequest{}, "Execute")).Return(&buckets, nil, nil).Build()
					FunctionMocker.When(func() bool {
						buckets := powerscale.V12S3Buckets{
							Resume:  nil,
							Buckets: []powerscale.V12S3Bucket{{Id: bucketName}},
							Total:   nil,
						}
						if FunctionMocker.MockTimes() > 0 {
							FunctionMocker.UnPatch()
							FunctionMocker = mockey.Mock(mockey.GetMethod(powerscale.ApiListProtocolsv12S3BucketsRequest{}, "Execute")).Return(&buckets, nil, nil).Build()
						}
						return FunctionMocker.MockTimes() == 0
					})
				},
				Config: ProviderConfig + S3BucketAllDatasourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerscale_s3_bucket.bucket_datasource_test_all", "filter.#", "0"),
				),
			},
		},
	})
	FunctionMocker.UnPatch()
}

func TestAccS3BucketDatasourceErrorCopyFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Read testing
			{
				PreConfig: func() {
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + S3BucketAllDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
			//Read testing with names
			{
				PreConfig: func() {
					FunctionMocker.UnPatch()
					FunctionMocker = mockey.Mock(helper.CopyFields).Return(fmt.Errorf("mock error")).Build()
				},
				Config:      ProviderConfig + S3BucketDatasourceConfig,
				ExpectError: regexp.MustCompile("mock error"),
			},
		},
	})
	FunctionMocker.UnPatch()
}

var S3BucketAllDatasourceConfig = FileSystemResourceConfigCommon3 + fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_resource_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	name = "%s"
	path = "/ifs/%s"
	create_path = true
	acl = [{
		grantee = {
			name = "Everyone"
			type = "wellknown"
		}
		permission = "FULL_CONTROL"
	}]
}

data "powerscale_s3_bucket" "bucket_datasource_test_all" {}
`, bucketName, bucketName)

var S3BucketDatasourceConfig = FileSystemResourceConfigCommon3 + fmt.Sprintf(`
resource "powerscale_s3_bucket" "bucket_resource_test" {
	depends_on = [powerscale_filesystem.file_system_test] 
	name = "%s"
	path = "/ifs/%s"
	create_path = true
	acl = [{
		grantee = {
			name = "Everyone"
			type = "wellknown"
		}
		permission = "FULL_CONTROL"
	}]
	zone  = "System"
	owner = "root"
}

data "powerscale_s3_bucket" "bucket_datasource_test" {
	filter {
		zone  = "System"
		owner = "root"
	}
  	depends_on = [
    	powerscale_s3_bucket.bucket_resource_test
  	]
}
`, bucketName, bucketName)
