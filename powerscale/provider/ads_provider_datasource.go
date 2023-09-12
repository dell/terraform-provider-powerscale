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
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AdsProviderDataSource{}

// NewAdsProviderDataSource creates a new data source.
func NewAdsProviderDataSource() datasource.DataSource {
	return &AdsProviderDataSource{}
}

// AdsProviderDataSource defines the data source implementation.
type AdsProviderDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *AdsProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adsprovider"
}

// Schema describes the data source arguments.
func (d *AdsProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "ADS Provider Datasource. Joining your cluster to an Active Directory domain allows you to perform user and group authentication. This Terraform DataSource is used to query the details of existing ADS providers from PowerScale array.",
		MarkdownDescription: "ADS Provider Datasource. Joining your cluster to an Active Directory domain allows you to perform user and group authentication. This Terraform DataSource is used to query the details of existing ADS providers from PowerScale array.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the ads provider instance.",
				MarkdownDescription: "Unique identifier of the ads provider instance.",
				Computed:            true,
			},
			"ads_providers_details": schema.ListNestedAttribute{
				Description:         "List of AdsProviders.",
				MarkdownDescription: "List of AdsProviders.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allocate_gids": schema.BoolAttribute{
							Description:         "Allocates an ID for an unmapped Active Directory (ADS) group. ADS groups without GIDs can be proactively assigned a GID by the ID mapper. If the ID mapper option is disabled, GIDs are not proactively assigned, and when a primary group for a user does not include a GID, the system may allocate one. ",
							MarkdownDescription: "Allocates an ID for an unmapped Active Directory (ADS) group. ADS groups without GIDs can be proactively assigned a GID by the ID mapper. If the ID mapper option is disabled, GIDs are not proactively assigned, and when a primary group for a user does not include a GID, the system may allocate one. ",
							Computed:            true,
						},
						"allocate_uids": schema.BoolAttribute{
							Description:         "Allocates a user ID for an unmapped Active Directory (ADS) user. ADS users without UIDs can be proactively assigned a UID by the ID mapper. IF the ID mapper option is disabled, UIDs are not proactively assigned, and when an identify for a user does not include a UID, the system may allocate one.",
							MarkdownDescription: "Allocates a user ID for an unmapped Active Directory (ADS) user. ADS users without UIDs can be proactively assigned a UID by the ID mapper. IF the ID mapper option is disabled, UIDs are not proactively assigned, and when an identify for a user does not include a UID, the system may allocate one.",
							Computed:            true,
						},
						"assume_default_domain": schema.BoolAttribute{
							Description:         "Enables lookup of unqualified user names in the primary domain.",
							MarkdownDescription: "Enables lookup of unqualified user names in the primary domain.",
							Computed:            true,
						},
						"authentication": schema.BoolAttribute{
							Description:         "Enables authentication and identity management through the authentication provider.",
							MarkdownDescription: "Enables authentication and identity management through the authentication provider.",
							Computed:            true,
						},
						"check_online_interval": schema.Int64Attribute{
							Description:         "Specifies the time in seconds between provider online checks.",
							MarkdownDescription: "Specifies the time in seconds between provider online checks.",
							Computed:            true,
						},
						"controller_time": schema.Int64Attribute{
							Description:         "Specifies the current time for the domain controllers.",
							MarkdownDescription: "Specifies the current time for the domain controllers.",
							Computed:            true,
						},
						"create_home_directory": schema.BoolAttribute{
							Description:         "Automatically creates a home directory on the first login.",
							MarkdownDescription: "Automatically creates a home directory on the first login.",
							Computed:            true,
						},
						"domain_offline_alerts": schema.BoolAttribute{
							Description:         "Sends an alert if the domain goes offline.",
							MarkdownDescription: "Sends an alert if the domain goes offline.",
							Computed:            true,
						},
						"dup_spns": schema.ListAttribute{
							Description:         "Get duplicate SPNs in the provider domain.",
							MarkdownDescription: "Get duplicate SPNs in the provider domain.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"extra_expected_spns": schema.ListAttribute{
							Description:         "List of additional SPNs to expect beyond what automatic checking routines might find.",
							MarkdownDescription: "List of additional SPNs to expect beyond what automatic checking routines might find.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"findable_groups": schema.ListAttribute{
							Description:         "Sets list of groups that can be resolved.",
							MarkdownDescription: "Sets list of groups that can be resolved.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"findable_users": schema.ListAttribute{
							Description:         "Sets list of users that can be resolved.",
							MarkdownDescription: "Sets list of users that can be resolved.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"forest": schema.StringAttribute{
							Description:         "Specifies the Active Directory forest.",
							MarkdownDescription: "Specifies the Active Directory forest.",
							Computed:            true,
						},
						"groupnet": schema.StringAttribute{
							Description:         "Groupnet identifier.",
							MarkdownDescription: "Groupnet identifier.",
							Computed:            true,
						},
						"home_directory_template": schema.StringAttribute{
							Description:         "Specifies the path to the home directory template.",
							MarkdownDescription: "Specifies the path to the home directory template.",
							Computed:            true,
						},
						"hostname": schema.StringAttribute{
							Description:         "Specifies the fully qualified hostname stored in the machine account.",
							MarkdownDescription: "Specifies the fully qualified hostname stored in the machine account.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Specifies the ID of the Active Directory provider instance.",
							MarkdownDescription: "Specifies the ID of the Active Directory provider instance.",
							Computed:            true,
						},
						"ignore_all_trusts": schema.BoolAttribute{
							Description:         "If set to true, ignores all trusted domains.",
							MarkdownDescription: "If set to true, ignores all trusted domains.",
							Computed:            true,
						},
						"ignored_trusted_domains": schema.ListAttribute{
							Description:         "Includes trusted domains when 'ignore_all_trusts' is set to false.",
							MarkdownDescription: "Includes trusted domains when 'ignore_all_trusts' is set to false.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"include_trusted_domains": schema.ListAttribute{
							Description:         "Includes trusted domains when 'ignore_all_trusts' is set to true.",
							MarkdownDescription: "Includes trusted domains when 'ignore_all_trusts' is set to true.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"instance": schema.StringAttribute{
							Description:         "Specifies Active Directory provider instance.",
							MarkdownDescription: "Specifies Active Directory provider instance.",
							Computed:            true,
						},
						"ldap_sign_and_seal": schema.BoolAttribute{
							Description:         "Enables encryption and signing on LDAP requests.",
							MarkdownDescription: "Enables encryption and signing on LDAP requests.",
							Computed:            true,
						},
						"login_shell": schema.StringAttribute{
							Description:         "Specifies the login shell path.",
							MarkdownDescription: "Specifies the login shell path.",
							Computed:            true,
						},
						"lookup_domains": schema.ListAttribute{
							Description:         "Limits user and group lookups to the specified domains.",
							MarkdownDescription: "Limits user and group lookups to the specified domains.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"lookup_groups": schema.BoolAttribute{
							Description:         "Looks up AD groups in other providers before allocating a group ID.",
							MarkdownDescription: "Looks up AD groups in other providers before allocating a group ID.",
							Computed:            true,
						},
						"lookup_normalize_groups": schema.BoolAttribute{
							Description:         "Normalizes AD group names to lowercase before look up.",
							MarkdownDescription: "Normalizes AD group names to lowercase before look up.",
							Computed:            true,
						},
						"lookup_normalize_users": schema.BoolAttribute{
							Description:         "Normalize AD user names to lowercase before look up.",
							MarkdownDescription: "Normalize AD user names to lowercase before look up.",
							Computed:            true,
						},
						"lookup_users": schema.BoolAttribute{
							Description:         "Looks up AD users in other providers before allocating a user ID.",
							MarkdownDescription: "Looks up AD users in other providers before allocating a user ID.",
							Computed:            true,
						},
						"machine_account": schema.StringAttribute{
							Description:         "Specifies the machine account name when creating a SAM account with Active Directory.",
							MarkdownDescription: "Specifies the machine account name when creating a SAM account with Active Directory.",
							Computed:            true,
						},
						"machine_password_changes": schema.BoolAttribute{
							Description:         "Enables periodic changes of the machine password for security.",
							MarkdownDescription: "Enables periodic changes of the machine password for security.",
							Computed:            true,
						},
						"machine_password_lifespan": schema.Int64Attribute{
							Description:         "Sets maximum age of a password in seconds.",
							MarkdownDescription: "Sets maximum age of a password in seconds.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Specifies the Active Directory provider name.",
							MarkdownDescription: "Specifies the Active Directory provider name.",
							Computed:            true,
						},
						"netbios_domain": schema.StringAttribute{
							Description:         "Specifies the NetBIOS domain name associated with the machine account.",
							MarkdownDescription: "Specifies the NetBIOS domain name associated with the machine account.",
							Computed:            true,
						},
						"node_dc_affinity": schema.StringAttribute{
							Description:         "Specifies the domain controller for which the node has affinity.",
							MarkdownDescription: "Specifies the domain controller for which the node has affinity.",
							Computed:            true,
						},
						"node_dc_affinity_timeout": schema.Int64Attribute{
							Description:         "Specifies the timeout for the domain controller for which the local node has affinity.",
							MarkdownDescription: "Specifies the timeout for the domain controller for which the local node has affinity.",
							Computed:            true,
						},
						"nss_enumeration": schema.BoolAttribute{
							Description:         "Enables the Active Directory provider to respond to 'getpwent' and 'getgrent' requests.",
							MarkdownDescription: "Enables the Active Directory provider to respond to 'getpwent' and 'getgrent' requests.",
							Computed:            true,
						},
						"primary_domain": schema.StringAttribute{
							Description:         "Specifies the AD domain to which the provider is joined.",
							MarkdownDescription: "Specifies the AD domain to which the provider is joined.",
							Computed:            true,
						},
						"restrict_findable": schema.BoolAttribute{
							Description:         "Check the provider for filtered lists of findable and unfindable users and groups.",
							MarkdownDescription: "Check the provider for filtered lists of findable and unfindable users and groups.",
							Computed:            true,
						},
						"rpc_call_timeout": schema.Int64Attribute{
							Description:         "The maximum amount of time (in seconds) an RPC call to Active Directory is allowed to take.",
							MarkdownDescription: "The maximum amount of time (in seconds) an RPC call to Active Directory is allowed to take.",
							Computed:            true,
						},
						"server_retry_limit": schema.Int64Attribute{
							Description:         "The number of retries attempted when a call to Active Directory fails due to network error.",
							MarkdownDescription: "The number of retries attempted when a call to Active Directory fails due to network error.",
							Computed:            true,
						},
						"sfu_support": schema.StringAttribute{
							Description:         "Specifies whether to support RFC 2307 attributes on ADS domain controllers.",
							MarkdownDescription: "Specifies whether to support RFC 2307 attributes on ADS domain controllers.",
							Computed:            true,
						},
						"site": schema.StringAttribute{
							Description:         "Specifies the site for the Active Directory.",
							MarkdownDescription: "Specifies the site for the Active Directory.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							Description:         "Specifies the status of the provider.",
							MarkdownDescription: "Specifies the status of the provider.",
							Computed:            true,
						},
						"store_sfu_mappings": schema.BoolAttribute{
							Description:         "Stores SFU mappings permanently in the ID mapper.",
							MarkdownDescription: "Stores SFU mappings permanently in the ID mapper.",
							Computed:            true,
						},
						"system": schema.BoolAttribute{
							Description:         "If set to true, indicates that this provider instance was created by OneFS and cannot be removed.",
							MarkdownDescription: "If set to true, indicates that this provider instance was created by OneFS and cannot be removed.",
							Computed:            true,
						},
						"unfindable_groups": schema.ListAttribute{
							Description:         "Specifies groups that cannot be resolved by the provider.",
							MarkdownDescription: "Specifies groups that cannot be resolved by the provider.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"unfindable_users": schema.ListAttribute{
							Description:         "Specifies users that cannot be resolved by the provider.",
							MarkdownDescription: "Specifies users that cannot be resolved by the provider.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"zone_name": schema.StringAttribute{
							Description:         "Specifies the name of the access zone in which this provider was created.",
							MarkdownDescription: "Specifies the name of the access zone in which this provider was created.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Description:         "Filter ads providers by names.",
						MarkdownDescription: "Filter ads providers by names.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"scope": schema.StringAttribute{
						Description:         "Filter ads providers by scope.",
						MarkdownDescription: "Filter ads providers by scope.",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *AdsProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *AdsProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading ads provider data source")

	var state models.AdsProviderDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	adsProviderParams := d.client.PscaleOpenAPIClient.AuthApi.ListAuthv14ProvidersAds(ctx)

	if state.AdsProviderFilter != nil && !state.AdsProviderFilter.Scope.IsNull() {
		adsProviderParams = adsProviderParams.Scope(state.AdsProviderFilter.Scope.ValueString())
	}

	result, _, err := adsProviderParams.Execute()

	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of ads providers",
			message,
		)
		return
	}

	var adsProviders []models.AdsProviderDetailModel
	for _, adsItem := range result.Ads {
		val := adsItem
		adsProvider, err := helper.AdsProviderDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadAdsProviderErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error getting the list of ads providers",
				message,
			)
			return
		}
		adsProviders = append(adsProviders, adsProvider)
	}

	state.AdsProviders = adsProviders

	// filter ads providers by names
	if state.AdsProviderFilter != nil && len(state.AdsProviderFilter.Names) > 0 {
		var validAdsProviders []string
		var filteredAdsProviders []models.AdsProviderDetailModel

		for _, ads := range state.AdsProviders {
			for _, name := range state.AdsProviderFilter.Names {
				if !name.IsNull() && ads.Name.Equal(name) {
					filteredAdsProviders = append(filteredAdsProviders, ads)
					validAdsProviders = append(validAdsProviders, fmt.Sprintf("Name: %s", ads.Name))
					continue
				}
			}
		}

		state.AdsProviders = filteredAdsProviders

		if len(state.AdsProviders) != len(state.AdsProviderFilter.Names) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered ads names is not a valid powerscale ads provider.",
				fmt.Sprintf("Valid ads providers: [%v], filtered list: [%v]", strings.Join(validAdsProviders, " ; "), state.AdsProviderFilter.Names),
			)
		}
	}

	// save into the Terraform state.
	state.ID = types.StringValue("ads_provider_datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading ads provider data source ")
}
