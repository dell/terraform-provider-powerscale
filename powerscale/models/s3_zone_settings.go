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

// S3ZoneSettingsResource defines the resource arguments.
type S3ZoneSettingsResource struct {
	Zone                      types.String `tfsdk:"zone"`
	BaseDomain                types.String `tfsdk:"base_domain"`
	BucketDirectoryCreateMode types.Int64  `tfsdk:"bucket_directory_create_mode"`
	ObjectACLPolicy           types.String `tfsdk:"object_acl_policy"`
	RootPath                  types.String `tfsdk:"root_path"`
}
