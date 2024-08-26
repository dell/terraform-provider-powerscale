package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// S3GlobalSettingResource defines the resource arguments.
type S3GlobalSettingResource struct {
	// Specifies if the service is enabled.
	Service types.Bool `tfsdk:"service"`
	// Specifies if the service is HTTPS only.
	HTTPSOnly types.Bool `tfsdk:"https_only"`
	// Specifies the HTTP port.
	HTTPPort types.Int64 `tfsdk:"http_port"`
	// Specifies the HTTPS port.
	HTTPSPort types.Int64 `tfsdk:"https_port"`
}
