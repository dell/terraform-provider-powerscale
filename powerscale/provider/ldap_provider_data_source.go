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
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &LdapProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &LdapProviderDataSource{}
)

// NewLdapProviderDataSource creates a new LDAP provider data source.
func NewLdapProviderDataSource() datasource.DataSource {
	return &LdapProviderDataSource{}
}

// LdapProviderDataSource defines the data source implementation.
type LdapProviderDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *LdapProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_provider"
}

// Schema describes the data source arguments.
func (d *LdapProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing LDAP providers from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale LDAP provider enables you to define, query, and modify directory services and resources.",
		Description:         "This datasource is used to query the existing LDAP providers from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale LDAP provider enables you to define, query, and modify directory services and resources.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the LDAP provider instance.",
				Description:         "Unique identifier of the LDAP provider instance.",
			},
			"ldap_providers": schema.ListNestedAttribute{
				MarkdownDescription: "List of LDAP providers.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Specifies the ID of the LDAP provider.",
							MarkdownDescription: "Specifies the ID of the LDAP provider.",
							Computed:            true,
						},
						"zone_name": schema.StringAttribute{
							Description:         "Specifies the name of the access zone in which this provider was created.",
							MarkdownDescription: "Specifies the name of the access zone in which this provider was created.",
							Computed:            true,
						},
						"groupnet": schema.StringAttribute{
							Description:         "Groupnet identifier. ",
							MarkdownDescription: "Groupnet identifier. ",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Specifies the name of the LDAP provider. ",
							MarkdownDescription: "Specifies the name of the LDAP provider. ",
							Computed:            true,
						},
						"base_dn": schema.StringAttribute{
							Description:         "Specifies the root of the tree in which to search identities. ",
							MarkdownDescription: "Specifies the root of the tree in which to search identities. ",
							Computed:            true,
						},
						"server_uris": schema.ListAttribute{
							Description:         "Specifies the server URIs. ",
							MarkdownDescription: "Specifies the server URIs. ",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"authentication": schema.BoolAttribute{
							Description:         "If true, enables authentication and identity management through the authentication provider. ",
							MarkdownDescription: "If true, enables authentication and identity management through the authentication provider. ",
							Computed:            true,
						},
						"balance_servers": schema.BoolAttribute{
							Description:         "If true, connects the provider to a random server. ",
							MarkdownDescription: "If true, connects the provider to a random server. ",
							Computed:            true,
						},
						"create_home_directory": schema.BoolAttribute{
							Description:         "Automatically create the home directory on the first login. ",
							MarkdownDescription: "Automatically create the home directory on the first login. ",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							Description:         "If true, enables the LDAP provider. ",
							MarkdownDescription: "If true, enables the LDAP provider. ",
							Computed:            true,
						},
						"enumerate_groups": schema.BoolAttribute{
							Description:         "If true, allows the provider to enumerate groups. ",
							MarkdownDescription: "If true, allows the provider to enumerate groups. ",
							Computed:            true,
						},
						"enumerate_users": schema.BoolAttribute{
							Description:         "If true, allows the provider to enumerate users. ",
							MarkdownDescription: "If true, allows the provider to enumerate users. ",
							Computed:            true,
						},
						"ignore_tls_errors": schema.BoolAttribute{
							Description:         "If true, continues over secure connections even if identity checks fail. ",
							MarkdownDescription: "If true, continues over secure connections even if identity checks fail. ",
							Computed:            true,
						},
						"normalize_groups": schema.BoolAttribute{
							Description:         "Normalizes group names to lowercase before look up. ",
							MarkdownDescription: "Normalizes group names to lowercase before look up. ",
							Computed:            true,
						},
						"normalize_users": schema.BoolAttribute{
							Description:         "Normalizes user names to lowercase before look up. ",
							MarkdownDescription: "Normalizes user names to lowercase before look up. ",
							Computed:            true,
						},
						"require_secure_connection": schema.BoolAttribute{
							Description:         "Determines whether to continue over a non-TLS connection. ",
							MarkdownDescription: "Determines whether to continue over a non-TLS connection. ",
							Computed:            true,
						},
						"restrict_findable": schema.BoolAttribute{
							Description:         "If true, checks the provider for filtered lists of findable and unfindable users and groups. ",
							MarkdownDescription: "If true, checks the provider for filtered lists of findable and unfindable users and groups. ",
							Computed:            true,
						},
						"restrict_listable": schema.BoolAttribute{
							Description:         "If true, checks the provider for filtered lists of listable and unlistable users and groups. ",
							MarkdownDescription: "If true, checks the provider for filtered lists of listable and unlistable users and groups. ",
							Computed:            true,
						},
						"system": schema.BoolAttribute{
							Description:         "If true, indicates that this provider instance was created by OneFS and cannot be removed. ",
							MarkdownDescription: "If true, indicates that this provider instance was created by OneFS and cannot be removed. ",
							Computed:            true,
						},
						"bind_timeout": schema.Int64Attribute{
							Description:         "Specifies the timeout in seconds when binding to an LDAP server. ",
							MarkdownDescription: "Specifies the timeout in seconds when binding to an LDAP server. ",
							Computed:            true,
						},
						"check_online_interval": schema.Int64Attribute{
							Description:         "Specifies the time in seconds between provider online checks. ",
							MarkdownDescription: "Specifies the time in seconds between provider online checks. ",
							Computed:            true,
						},
						"search_timeout": schema.Int64Attribute{
							Description:         "Specifies the search timeout period in seconds. ",
							MarkdownDescription: "Specifies the search timeout period in seconds. ",
							Computed:            true,
						},
						"findable_groups": schema.ListAttribute{
							Description:         "Specifies the list of groups that can be resolved. ",
							MarkdownDescription: "Specifies the list of groups that can be resolved. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"findable_users": schema.ListAttribute{
							Description:         "Specifies the list of users that can be resolved. ",
							MarkdownDescription: "Specifies the list of users that can be resolved. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"listable_groups": schema.ListAttribute{
							Description:         "Specifies the groups that can be viewed in the provider. ",
							MarkdownDescription: "Specifies the groups that can be viewed in the provider. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"listable_users": schema.ListAttribute{
							Description:         "Specifies the users that can be viewed in the provider. ",
							MarkdownDescription: "Specifies the users that can be viewed in the provider. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"unfindable_groups": schema.ListAttribute{
							Description:         "Specifies the groups that cannot be resolved by the provider. ",
							MarkdownDescription: "Specifies the groups that cannot be resolved by the provider. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"unfindable_users": schema.ListAttribute{
							Description:         "Specifies users that cannot be resolved by the provider. ",
							MarkdownDescription: "Specifies users that cannot be resolved by the provider. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"unlistable_groups": schema.ListAttribute{
							Description:         "Specifies a group that cannot be listed by the provider. ",
							MarkdownDescription: "Specifies a group that cannot be listed by the provider. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"unlistable_users": schema.ListAttribute{
							Description:         "Specifies a user that cannot be listed by the provider. ",
							MarkdownDescription: "Specifies a user that cannot be listed by the provider. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"alternate_security_identities_attribute": schema.StringAttribute{
							Description:         "Specifies the attribute name used when searching for alternate security identities. ",
							MarkdownDescription: "Specifies the attribute name used when searching for alternate security identities. ",
							Computed:            true,
						},
						"bind_dn": schema.StringAttribute{
							Description:         "Specifies the distinguished name for binding to the LDAP server. ",
							MarkdownDescription: "Specifies the distinguished name for binding to the LDAP server. ",
							Computed:            true,
						},
						"bind_mechanism": schema.StringAttribute{
							Description:         "Specifies which bind mechanism to use when connecting to an LDAP server. The only supported option is the 'simple' value. ",
							MarkdownDescription: "Specifies which bind mechanism to use when connecting to an LDAP server. The only supported option is the 'simple' value. ",
							Computed:            true,
						},
						"certificate_authority_file": schema.StringAttribute{
							Description:         "Specifies the path to the root certificates file. ",
							MarkdownDescription: "Specifies the path to the root certificates file. ",
							Computed:            true,
						},
						"cn_attribute": schema.StringAttribute{
							Description:         "Specifies the canonical name. ",
							MarkdownDescription: "Specifies the canonical name. ",
							Computed:            true,
						},
						"crypt_password_attribute": schema.StringAttribute{
							Description:         "Specifies the hashed password value. ",
							MarkdownDescription: "Specifies the hashed password value. ",
							Computed:            true,
						},
						"email_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Email attribute. ",
							MarkdownDescription: "Specifies the LDAP Email attribute. ",
							Computed:            true,
						},
						"gecos_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP GECOS attribute. ",
							MarkdownDescription: "Specifies the LDAP GECOS attribute. ",
							Computed:            true,
						},
						"gid_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP GID attribute. ",
							MarkdownDescription: "Specifies the LDAP GID attribute. ",
							Computed:            true,
						},
						"group_base_dn": schema.StringAttribute{
							Description:         "Specifies the distinguished name of the entry where LDAP searches for groups are started. ",
							MarkdownDescription: "Specifies the distinguished name of the entry where LDAP searches for groups are started. ",
							Computed:            true,
						},
						"group_domain": schema.StringAttribute{
							Description:         "Specifies the domain for this provider through which groups are qualified. ",
							MarkdownDescription: "Specifies the domain for this provider through which groups are qualified. ",
							Computed:            true,
						},
						"group_filter": schema.StringAttribute{
							Description:         "Specifies the LDAP filter for group objects. ",
							MarkdownDescription: "Specifies the LDAP filter for group objects. ",
							Computed:            true,
						},
						"group_members_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Group Members attribute. ",
							MarkdownDescription: "Specifies the LDAP Group Members attribute. ",
							Computed:            true,
						},
						"group_search_scope": schema.StringAttribute{
							Description:         "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							MarkdownDescription: "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							Computed:            true,
						},
						"home_directory_template": schema.StringAttribute{
							Description:         "Specifies the path to the home directory template. ",
							MarkdownDescription: "Specifies the path to the home directory template. ",
							Computed:            true,
						},
						"homedir_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Homedir attribute. ",
							MarkdownDescription: "Specifies the LDAP Homedir attribute. ",
							Computed:            true,
						},
						"login_shell": schema.StringAttribute{
							Description:         "Specifies the login shell path. ",
							MarkdownDescription: "Specifies the login shell path. ",
							Computed:            true,
						},
						"member_lookup_method": schema.StringAttribute{
							Description:         "Sets the method by which group member lookups are performed. Use caution when changing this option directly. Acceptable values: \"default\", \"rfc2307bis\". ",
							MarkdownDescription: "Sets the method by which group member lookups are performed. Use caution when changing this option directly. Acceptable values: \"default\", \"rfc2307bis\". ",
							Computed:            true,
						},
						"member_of_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Query Member Of attribute, which performs reverse membership queries. ",
							MarkdownDescription: "Specifies the LDAP Query Member Of attribute, which performs reverse membership queries. ",
							Computed:            true,
						},
						"name_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP UID attribute, which is used as the login name. ",
							MarkdownDescription: "Specifies the LDAP UID attribute, which is used as the login name. ",
							Computed:            true,
						},
						"netgroup_base_dn": schema.StringAttribute{
							Description:         "Specifies the distinguished name of the entry where LDAP searches for netgroups are started. ",
							MarkdownDescription: "Specifies the distinguished name of the entry where LDAP searches for netgroups are started. ",
							Computed:            true,
						},
						"netgroup_filter": schema.StringAttribute{
							Description:         "Specifies the LDAP filter for netgroup objects. ",
							MarkdownDescription: "Specifies the LDAP filter for netgroup objects. ",
							Computed:            true,
						},
						"netgroup_members_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Netgroup Members attribute. ",
							MarkdownDescription: "Specifies the LDAP Netgroup Members attribute. ",
							Computed:            true,
						},
						"netgroup_search_scope": schema.StringAttribute{
							Description:         "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							MarkdownDescription: "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							Computed:            true,
						},
						"netgroup_triple_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Netgroup Triple attribute. ",
							MarkdownDescription: "Specifies the LDAP Netgroup Triple attribute. ",
							Computed:            true,
						},
						"nt_password_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP NT Password attribute. ",
							MarkdownDescription: "Specifies the LDAP NT Password attribute. ",
							Computed:            true,
						},
						"ntlm_support": schema.StringAttribute{
							Description:         "Specifies which NTLM versions to support for users with NTLM-compatible credentials. Acceptable values: \"all\", \"v2only\", \"none\". ",
							MarkdownDescription: "Specifies which NTLM versions to support for users with NTLM-compatible credentials. Acceptable values: \"all\", \"v2only\", \"none\". ",
							Computed:            true,
						},
						"provider_domain": schema.StringAttribute{
							Description:         "Specifies the provider domain. ",
							MarkdownDescription: "Specifies the provider domain. ",
							Computed:            true,
						},
						"search_scope": schema.StringAttribute{
							Description:         "Specifies the default depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							MarkdownDescription: "Specifies the default depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							Computed:            true,
						},
						"shadow_expire_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the absolute date to expire the account. ",
							MarkdownDescription: "Sets the attribute name that indicates the absolute date to expire the account. ",
							Computed:            true,
						},
						"shadow_flag_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the section of the shadow map that is used to store the flag value. ",
							MarkdownDescription: "Sets the attribute name that indicates the section of the shadow map that is used to store the flag value. ",
							Computed:            true,
						},
						"shadow_inactive_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the number of days of inactivity that is allowed for the user. ",
							MarkdownDescription: "Sets the attribute name that indicates the number of days of inactivity that is allowed for the user. ",
							Computed:            true,
						},
						"shadow_last_change_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the last change of the shadow information. ",
							MarkdownDescription: "Sets the attribute name that indicates the last change of the shadow information. ",
							Computed:            true,
						},
						"shadow_max_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the maximum number of days a password can be valid. ",
							MarkdownDescription: "Sets the attribute name that indicates the maximum number of days a password can be valid. ",
							Computed:            true,
						},
						"shadow_min_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the minimum number of days between shadow changes. ",
							MarkdownDescription: "Sets the attribute name that indicates the minimum number of days between shadow changes. ",
							Computed:            true,
						},
						"shadow_user_filter": schema.StringAttribute{
							Description:         "Sets LDAP filter for shadow user objects. ",
							MarkdownDescription: "Sets LDAP filter for shadow user objects. ",
							Computed:            true,
						},
						"shadow_warning_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the number of days before the password expires to warn the user. ",
							MarkdownDescription: "Sets the attribute name that indicates the number of days before the password expires to warn the user. ",
							Computed:            true,
						},
						"shell_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP Shell attribute. ",
							MarkdownDescription: "Specifies the LDAP Shell attribute. ",
							Computed:            true,
						},
						"ssh_public_key_attribute": schema.StringAttribute{
							Description:         "Sets the attribute name that indicates the SSH Public Key for the user. ",
							MarkdownDescription: "Sets the attribute name that indicates the SSH Public Key for the user. ",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							Description:         "Specifies the status of the provider. ",
							MarkdownDescription: "Specifies the status of the provider. ",
							Computed:            true,
						},
						"tls_protocol_min": schema.StringAttribute{
							Description:         "Specifies the minimum TLS protocol version. ",
							MarkdownDescription: "Specifies the minimum TLS protocol version. ",
							Computed:            true,
						},
						"uid_attribute": schema.StringAttribute{
							Description:         "Specifies the LDAP UID Number attribute. ",
							MarkdownDescription: "Specifies the LDAP UID Number attribute. ",
							Computed:            true,
						},
						"unique_group_members_attribute": schema.StringAttribute{
							Description:         "Sets the LDAP Unique Group Members attribute. ",
							MarkdownDescription: "Sets the LDAP Unique Group Members attribute. ",
							Computed:            true,
						},
						"user_base_dn": schema.StringAttribute{
							Description:         "Specifies the distinguished name of the entry at which to start LDAP searches for users. ",
							MarkdownDescription: "Specifies the distinguished name of the entry at which to start LDAP searches for users. ",
							Computed:            true,
						},
						"user_domain": schema.StringAttribute{
							Description:         "Specifies the domain for this provider through which users are qualified. ",
							MarkdownDescription: "Specifies the domain for this provider through which users are qualified. ",
							Computed:            true,
						},
						"user_filter": schema.StringAttribute{
							Description:         "Specifies the LDAP filter for user objects. ",
							MarkdownDescription: "Specifies the LDAP filter for user objects. ",
							Computed:            true,
						},
						"user_search_scope": schema.StringAttribute{
							Description:         "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							MarkdownDescription: "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". ",
							Computed:            true,
						},
						// Only available for PowerScale 9.5 and above
						"tls_revocation_check_level": schema.StringAttribute{
							Description:         "This setting controls the behavior of the certificate revocation checking algorithm when the LDAP provider is presented with a digital certificate by an LDAP server. Acceptable values: \"none\", \"allowNoData\", \"allowNoSrc\", \"strict\". Only available for PowerScale 9.5 and above. ",
							MarkdownDescription: "This setting controls the behavior of the certificate revocation checking algorithm when the LDAP provider is presented with a digital certificate by an LDAP server. Acceptable values: \"none\", \"allowNoData\", \"allowNoSrc\", \"strict\". Only available for PowerScale 9.5 and above. ",
							Computed:            true,
						},
						"ocsp_server_uris": schema.ListAttribute{
							Description:         "Specifies the OCSP server URIs. Only available for PowerScale 9.5 and above. ",
							MarkdownDescription: "Specifies the OCSP server URIs. Only available for PowerScale 9.5 and above. ",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
					"scope": schema.StringAttribute{
						Description:         "If specified as \"effective\" or not specified, all fields are returned.  If specified as \"user\", only fields with non-default values are shown.  If specified as \"default\", the original values are returned. ",
						MarkdownDescription: "If specified as \"effective\" or not specified, all fields are returned.  If specified as \"user\", only fields with non-default values are shown.  If specified as \"default\", the original values are returned. ",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *LdapProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *LdapProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading LdapProvider data source ")

	var state models.LdapProviderDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ldapProviders, err := helper.GetAllLdapProvidersWithFilter(ctx, d.client, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale LdapProviders.", err.Error())
		return
	}

	// parse LdapProvider response to state LdapProvider model
	if err := helper.UpdateLdapProviderDataSourceState(ctx, &state, ldapProviders); err != nil {
		resp.Diagnostics.AddError("Error reading LdapProvider datasource plan",
			fmt.Sprintf("Could not list LdapProviders with error: %s", err.Error()))
		return
	}

	// filter LdapProvider by names
	if state.Filter != nil && len(state.Filter.Names) > 0 {
		// default scope cannot co-work with names filter
		if state.Filter.Scope.ValueString() == "default" {
			resp.Diagnostics.AddWarning("Returning all LDAP Provider with \"default\" scope", "filter.names is ignored when filter.scope is \"default\"")
		} else {
			var validLdapProviders []string
			var filteredLdapProviders []models.LdapProviderDetailModel

			for _, ldapProvider := range state.LdapProviders {
				for _, name := range state.Filter.Names {
					if ldapProvider.Name.Equal(name) {
						filteredLdapProviders = append(filteredLdapProviders, ldapProvider)
						validLdapProviders = append(validLdapProviders, ldapProvider.Name.ValueString())
						break
					}
				}
			}

			state.LdapProviders = filteredLdapProviders

			if len(state.LdapProviders) != len(state.Filter.Names) {
				resp.Diagnostics.AddError(
					"Error one or more of the filtered LdapProvider names is not a valid powerscale LdapProvider.",
					fmt.Sprintf("Valid LdapProviders: [%v], filtered list: [%v]", strings.Join(validLdapProviders, " , "), state.Filter.Names),
				)
			}
		}
	}

	state.ID = types.StringValue("ldap_provider_datasource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read LdapProvider data source ")
}
