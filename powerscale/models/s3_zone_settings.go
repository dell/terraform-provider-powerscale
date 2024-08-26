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
