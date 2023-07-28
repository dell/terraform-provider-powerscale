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
	"context"
	"terraform-provider-powerscale/client"

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
	client *client.Client
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Data describes the provider data model.
type Data struct {
	Endpoint                types.String `tfsdk:"endpoint"`
	Username                types.String `tfsdk:"username"`
	Group                   types.String `tfsdk:"group"`
	Password                types.String `tfsdk:"password"`
	Insecure                types.Bool   `tfsdk:"insecure"`
	VolumePath              types.String `tfsdk:"volume_path"`
	VolumePathPermissions   types.String `tfsdk:"volume_path_permissions"`
	AuthType                types.Int64  `tfsdk:"auth_type"`
	VerboseLogging          types.Int64  `tfsdk:"verbose_logging"`
	IgnoreUnresolvableHosts types.Bool   `tfsdk:"ignore_unresolvable_hosts"`
}

// Metadata describes the provider arguments.
func (p *PscaleProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerscale"
	resp.Version = p.version
}

// Schema describes the provider arguments.
func (p *PscaleProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
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
			"group": schema.StringAttribute{
				MarkdownDescription: "The user's group",
				Description:         "The user's group",
				Optional:            true,
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
			"volume_path": schema.StringAttribute{
				MarkdownDescription: "which base path to use when looking for volume directories",
				Description:         "which base path to use when looking for volume directories",
				Optional:            true,
			},
			"volume_path_permissions": schema.StringAttribute{
				MarkdownDescription: "permissions for new volume directory",
				Description:         "permissions for new volume directory",
				Optional:            true,
			},
			"auth_type": schema.Int64Attribute{
				MarkdownDescription: "what should be the auth type, 0 for basic and 1 for session-based",
				Description:         "what should be the auth type, 0 for basic and 1 for session-based",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"verbose_logging": schema.Int64Attribute{
				MarkdownDescription: "what verbose level should be used for logging, high(0), medium(1) or low(2)",
				Description:         "what verbose level should be used for logging, high(0), medium(1) or low(2)",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1, 2),
				},
			},
			"ignore_unresolvable_hosts": schema.BoolAttribute{
				MarkdownDescription: "whether to ignore unresolvable hosts",
				Description:         "whether to ignore unresolvable hosts",
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

	// Configuration values are now available.
	pscaleClient, err := client.NewClient(
		data.Endpoint.ValueString(),
		data.Insecure.ValueBool(),
		uint(data.VerboseLogging.ValueInt64()),
		data.Username.ValueString(),
		data.Group.ValueString(),
		data.Password.ValueString(),
		data.VolumePath.ValueString(),
		data.VolumePathPermissions.ValueString(),
		data.IgnoreUnresolvableHosts.ValueBool(),
		uint8(data.AuthType.ValueInt64()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create powerscale client",
			err.Error(),
		)
		return
	}

	// client configuration for data sources and resources
	p.client = pscaleClient
	resp.DataSourceData = pscaleClient
	resp.ResourceData = pscaleClient
}

// Resources describes the provider resources.
func (p *PscaleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccessZoneResource,
	}
}

// DataSources describes the provider data sources.
func (p *PscaleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccessZoneDataSource,
		NewClusterDataSource,
		NewUserDataSource,
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
