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
	"fmt"
	"os"
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
				MarkdownDescription: "The API endpoint, ex. https://172.17.177.230:8080. This can also be set using the environment variable POWERSCALE_ENDPOINT",
				Description:         "The API endpoint, ex. https://172.17.177.230:8080. This can also be set using the environment variable POWERSCALE_ENDPOINT",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username. This can also be set using the environment variable POWERSCALE_USERNAME",
				Description:         "The username. This can also be set using the environment variable POWERSCALE_USERNAME",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password. This can also be set using the environment variable POWERSCALE_PASSWORD",
				Description:         "The password. This can also be set using the environment variable POWERSCALE_PASSWORD",
				// This should remain optional so user can use environment variables if they choose.
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "whether to skip SSL validation. This can also be set using the environment variable POWERSCALE_INSECURE",
				Description:         "whether to skip SSL validation. This can also be set using the environment variable POWERSCALE_INSECURE",
				// This should remain optional so user can use environment variables if they choose.
				Optional: true,
			},
			"auth_type": schema.Int64Attribute{
				MarkdownDescription: "what should be the auth type, 0 for basic and 1 for session-based. This can also be set using the environment variable POWERSCALE_AUTH_TYPE",
				Description:         "what should be the auth type, 0 for basic and 1 for session-based. This can also be set using the environment variable POWERSCALE_AUTH_TYPE",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "specifies a time limit for requests. This can also be set using the environment variable POWERSCALE_TIMEOUT",
				Description:         "specifies a time limit for requests. This can also be set using the environment variable POWERSCALE_TIMEOUT",
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
		timeoutEnv, errTimeout := strconv.ParseInt(os.Getenv("POWERSCALE_TIMEOUT"), 10, 64)
		if errTimeout == nil {
			data.Timeout = types.Int64Value(timeoutEnv)
		} else {
			data.Timeout = types.Int64Value(2000)
		}
	}
	// If auth type is not set, use session based auth by default
	if data.AuthType.IsNull() || data.AuthType.IsUnknown() {
		authTypeEnv, errAuthType := strconv.ParseInt(os.Getenv("POWERSCALE_AUTH_TYPE"), 10, 64)
		if errAuthType == nil {
			data.AuthType = types.Int64Value(authTypeEnv)
		} else {
			data.AuthType = types.Int64Value(1)
		}
	}
	// If Insecure is not set, set to false by default
	if data.Insecure.IsNull() || data.Insecure.IsUnknown() {
		insecureEnv, errInsecure := strconv.ParseBool(os.Getenv("POWERSCALE_INSECURE"))
		if errInsecure == nil {
			data.Insecure = types.BoolValue(insecureEnv)
		} else {
			data.Insecure = types.BoolValue(false)
		}
	}

	if data.Username.IsUnknown() || data.Username.ValueString() == "" {
		usernameEnv := os.Getenv("POWERSCALE_USERNAME")
		if usernameEnv != "" {
			data.Username = types.StringValue(usernameEnv)
		} else {
			resp.Diagnostics.AddError(
				"Unable to find username",
				"Username cannot be an empty/unknown string",
			)
			return
		}
	}

	if data.Password.IsUnknown() || data.Password.ValueString() == "" {
		passEnv := os.Getenv("POWERSCALE_PASSWORD")
		if passEnv != "" {
			data.Password = types.StringValue(passEnv)
		} else {
			// Cannot connect to client with an unknown value
			resp.Diagnostics.AddWarning(
				"Unable to create client",
				"Password cannot be an empty/unknown string",
			)
			return
		}
	}

	if data.Endpoint.IsUnknown() || data.Endpoint.ValueString() == "" {
		endpointEnv := os.Getenv("POWERSCALE_ENDPOINT")
		if endpointEnv != "" {
			data.Endpoint = types.StringValue(endpointEnv)
		} else {
			// Cannot connect to client with an unknown value
			resp.Diagnostics.AddWarning(
				"Unable to create client",
				"Password cannot be an empty/unknown string",
			)
			return
		}
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
		NewNamespaceACLResource,
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
		NewSmbShareSettingsResource,
		NewSmbServerSettingsResource,
		NewClusterSnmpResource,
		NewS3KeyResource,
		NewClusterOwnerResource,
		NewSynciqPolicyResource,
		NewSyncIQGlobalSettingsResource,
		NewS3GlobalSettingResource,
		NewS3ZoneSettingsResource,
		NewClusterIdentityResource,
		NewClusterTimeResource,
		NewSyncIQRuleResource,
		NewSyncIQPeerCertificateResource,
		NewSupportAssistResource,
		NewWriteableSnapshotResource,
		NewSnapshotRestoreResource,
		NewNfsAliasResource,
		NewSyncIQReplicationJobResource,
		NewStoragepoolTierResource,
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
		NewNamespaceACLDataSource,
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
		NewStoragepoolTierDataSource,
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
		NewSmbShareSettingsDataSource,
		NewSyncIQPolicyDataSource,
		NewSyncIQRuleDataSource,
		NewSyncIQGlobalSettingsDataSource,
		NewSyncIQPeerCertificateDataSource,
		NewReplicationReportDataSource,
		NewNfsAliasDataSource,
		NewWritableSnapshotDataSource,
		NewSyncIQReplicationJobDataSource,
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

type commonResourceConfigurer struct {
	client *client.Client
	name   string
}

// Configure configures the resource.
func (d *commonResourceConfigurer) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Metadata describes the resource arguments.
func (d *commonResourceConfigurer) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.name
}

// ImportState implements resource.ResourceWithImportState.
func (d *commonResourceConfigurer) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
