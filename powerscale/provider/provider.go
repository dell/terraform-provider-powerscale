/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure PscaleProvider satisfies various provider interfaces.
var _ provider.Provider = &PscaleProvider{}

// PscaleProvider defines the provider implementation.
type PscaleProvider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Data describes the provider data model.
type Data struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
	AuthType types.Int64  `tfsdk:"auth_type"`
	Timeout  types.Int64  `tfsdk:"timeout"`
}

// Metadata describes the provider arguments.
func (p *PscaleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerscale"
	resp.Version = p.version
}

// Schema describes the provider arguments.
func (p *PscaleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Terraform provider for Dell PowerScale can be used to interact with a Dell PowerScale array in order to manage the array resources.",
		Description:         "The Terraform provider for Dell PowerScale can be used to interact with a Dell PowerScale array in order to manage the array resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The API endpoint, ex. https://172.17.177.230:8080",
				Description:         "The API endpoint, ex. https://172.17.177.230:8080",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username",
				Description:         "The username",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password",
				Description:         "The password",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "whether to skip SSL validation",
				Description:         "whether to skip SSL validation",
				Required:            true,
			},
			"auth_type": schema.Int64Attribute{
				MarkdownDescription: "what should be the auth type, 0 for basic and 1 for session-based",
				Description:         "what should be the auth type, 0 for basic and 1 for session-based",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "specifies a time limit for requests",
				Description:         "specifies a time limit for requests",
				Optional:            true,
			},
		},
	}
}

// Configure configures the provider.
func (p *PscaleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Data

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// if timeout is not set. use default value 2000
	if data.Timeout.IsNull() || data.Timeout.IsUnknown() {
		data.Timeout = types.Int64Value(2000)
	}
	// If auth type is not set, use session based auth by default
	if data.AuthType.IsNull() || data.AuthType.IsUnknown() {
		data.AuthType = types.Int64Value(1)
	}
	// Configuration values are now available.
	pscaleClient, err := client.NewClient(
		data.Endpoint.ValueString(),
		data.Insecure.ValueBool(),
		data.Username.ValueString(),
		data.Password.ValueString(),
		data.AuthType.ValueInt64(),
		data.Timeout.ValueInt64(),
	)

	if err != nil {
		message := helper.GetErrorString(err, "")
		resp.Diagnostics.AddError(
			"Unable to create powerscale client",
			message,
		)
		return
	}

	// client configuration for data sources and resources
	resp.DataSourceData = pscaleClient
	resp.ResourceData = pscaleClient
}

// Resources describes the provider resources.
func (p *PscaleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccessZoneResource,
		NewSmbShareResource,
		NewUserResource,
		NewNfsExportResource,
		NewUserGroupResource,
		NewAdsProviderResource,
		NewNetworkPoolResource,
		NewNtpServerResource,
		NewNtpSettingsResource,
		NewACLSettingsResource,
		NewRoleResource,
		NewFileSystemResource,
		NewSnapshotResource,
		NewSnapshotScheduleResource,
		NewGroupnetResource,
		NewQuotaResource,
		NewSubnetResource,
		NewSmartPoolSettingResource,
		NewNetworkSettingResource,
		NewNetworkRuleResource,
		NewLdapProviderResource,
		NewClusterEmailResource,
		NewFilePoolPolicyResource,
		NewNfsExportSettingsResource,
		NewNfsZoneSettingsResource,
		NewNfsGlobalSettingsResource,
		NewUserMappingRulesResource,
		NewS3BucketResource,
	}
}

// DataSources describes the provider data sources.
func (p *PscaleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFileSystemDataSource,
		NewAccessZoneDataSource,
		NewClusterDataSource,
		NewUserDataSource,
		NewSmbShareDataSource,
		NewAdsProviderDataSource,
		NewNetworkPoolDataSource,
		NewNtpServerDataSource,
		NewNtpSettingsDataSource,
		NewACLSettingsDataSource,
		NewRoleDataSource,
		NewRolePrivilegeDataSource,
		NewUserGroupDataSource,
		NewNfsExportDataSource,
		NewSnapshotDataSource,
		NewSnapshotScheduleDataSource,
		NewQuotaDataSource,
		NewSubnetDataSource,
		NewGroupnetDataSource,
		NewSmartPoolSettingDataSource,
		NewNetworkSettingDataSource,
		NewNetworkRuleDataSource,
		NewLdapProviderDataSource,
		NewClusterEmailDataSource,
		NewFilePoolPolicyDataSource,
		NewNfsExportSettingsDataSource,
		NewNfsGlobalSettingsDataSource,
		NewUserMappingRulesDataSource,
		NewS3BucketDataSource,
		NewNfsZoneSettingsDataSource,
		NewSmbServerSettingsDataSource,
	}
}

// New returns a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PscaleProvider{
			version: version,
		}
	}
}
