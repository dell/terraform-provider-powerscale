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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// S3BucketDatasource holds s3 bucket datasource schema attribute details.
type S3BucketDatasource struct {
	ID             types.String               `tfsdk:"id"`
	S3Buckets      []S3BucketDatasourceEntity `tfsdk:"s3_buckets"`
	S3BucketFilter *S3BucketDatasourceFilter  `tfsdk:"filter"`
}

// V10S3BucketACLItem specifies properties for an S3 Access Control Entry.
type V10S3BucketACLItem struct {
	// Specifies the persona of the file group.
	Grantee V1AuthAccessAccessItemFileGroup `tfsdk:"grantee"`
	// Specifies the S3 rights being allowed.
	Permission types.String `tfsdk:"permission"`
}

// S3BucketDatasourceEntity struct for S3 Bucket data source model.
type S3BucketDatasourceEntity struct {
	// Specifies an ordered list of S3 permissions.
	ACL []V10S3BucketACLItem `tfsdk:"acl"`
	// Description for this S3 bucket.
	Description types.String `tfsdk:"description"`
	// Bucket ID.
	ID types.String `tfsdk:"id"`
	// Bucket name.
	Name types.String `tfsdk:"name"`
	// Set behavior of modifying object acls.
	ObjectACLPolicy types.String `tfsdk:"object_acl_policy"`
	// Specifies the name of the owner.
	Owner types.String `tfsdk:"owner"`
	// Path of bucket within /ifs.
	Path types.String `tfsdk:"path"`
	// Zone ID.
	Zid types.Int64 `tfsdk:"zid"`
}

// S3BucketDatasourceFilter describes the filter data model.
type S3BucketDatasourceFilter struct {
	Zone  types.String `tfsdk:"zone"`
	Owner types.String `tfsdk:"owner"`
}

// S3BucketResource describes the resource data model.
type S3BucketResource struct {
	// Specifies an ordered list of S3 permissions.
	ACL types.List `tfsdk:"acl"`
	// Create path if does not exist.
	CreatePath types.Bool `tfsdk:"create_path"`
	// Description for this S3 bucket.
	Description types.String `tfsdk:"description"`
	// Bucket ID.
	ID types.String `tfsdk:"id"`
	// Bucket name.
	Name types.String `tfsdk:"name"`
	// Set behavior of modifying object acls.
	ObjectACLPolicy types.String `tfsdk:"object_acl_policy"`
	// Specifies the name of the owner.
	Owner types.String `tfsdk:"owner"`
	// Path of bucket within /ifs.
	Path types.String `tfsdk:"path"`
	// Numeric ID of the access zone to use.
	Zid types.Int64 `tfsdk:"zid"`
	// Name of the access zone to use.
	Zone types.String `tfsdk:"zone"`
}
