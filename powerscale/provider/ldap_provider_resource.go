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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &LdapProviderResource{}
	_ resource.ResourceWithConfigure   = &LdapProviderResource{}
	_ resource.ResourceWithImportState = &LdapProviderResource{}
)

// NewLdapProviderResource creates a new resource.
func NewLdapProviderResource() resource.Resource {
	return &LdapProviderResource{}
}

// LdapProviderResource defines the resource implementation.
type LdapProviderResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *LdapProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_provider"
}

// Schema describes the resource arguments.
func (r *LdapProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the LDAP provider entity of PowerScale Array. We can Create, Update and Delete the LDAP provider using this resource. We can also import an existing LDAP provider from PowerScale array. PowerScale LDAP provider enables you to define, query, and modify directory services and resources.",
		Description:         "This resource is used to manage the LDAP provider entity of PowerScale Array. We can Create, Update and Delete the LDAP provider using this resource. We can also import an existing LDAP provider from PowerScale array. PowerScale LDAP provider enables you to define, query, and modify directory services and resources.",

		Attributes: map[string]schema.Attribute{
			// Query param when creating and updating
			"ignore_unresolvable_server_urls": schema.BoolAttribute{
				Description:         "Ignore unresolvable server URIs when creating and updating. (Update Supported)",
				MarkdownDescription: "Ignore unresolvable server URIs when creating and updating. (Update Supported)",
				Optional:            true,
			},
			// Get params
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
			// Create params
			"groupnet": schema.StringAttribute{
				Description:         "Groupnet identifier. ",
				MarkdownDescription: "Groupnet identifier. ",
				Optional:            true,
				Computed:            true,
			},
			// Required params
			"name": schema.StringAttribute{
				Description:         "Specifies the name of the LDAP provider. (Update Supported)",
				MarkdownDescription: "Specifies the name of the LDAP provider. (Update Supported)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"base_dn": schema.StringAttribute{
				Description:         "Specifies the root of the tree in which to search identities. (Update Supported)",
				MarkdownDescription: "Specifies the root of the tree in which to search identities. (Update Supported)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"server_uris": schema.ListAttribute{
				Description:         "Specifies the server URIs. (Update Supported)",
				MarkdownDescription: "Specifies the server URIs. (Update Supported)",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(4, 2048)),
					listvalidator.SizeBetween(1, 4294967295),
				},
			},
			// Create and Update params - only available for PowerScale 9.5 and above
			"tls_revocation_check_level": schema.StringAttribute{
				Description:         "This setting controls the behavior of the certificate revocation checking algorithm when the LDAP provider is presented with a digital certificate by an LDAP server. Acceptable values: \"none\", \"allowNoData\", \"allowNoSrc\", \"strict\". Only available for PowerScale 9.5 and above. (Update Supported)",
				MarkdownDescription: "This setting controls the behavior of the certificate revocation checking algorithm when the LDAP provider is presented with a digital certificate by an LDAP server. Acceptable values: \"none\", \"allowNoData\", \"allowNoSrc\", \"strict\". Only available for PowerScale 9.5 and above. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"ocsp_server_uris": schema.ListAttribute{
				Description:         "Specifies the OCSP server URIs. Only available for PowerScale 9.5 and above. (Update Supported)",
				MarkdownDescription: "Specifies the OCSP server URIs. Only available for PowerScale 9.5 and above. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(4, 2048)),
					listvalidator.SizeBetween(0, 10),
				},
			},
			// Create and Update params
			"authentication": schema.BoolAttribute{
				Description:         "If true, enables authentication and identity management through the authentication provider. (Update Supported)",
				MarkdownDescription: "If true, enables authentication and identity management through the authentication provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"balance_servers": schema.BoolAttribute{
				Description:         "If true, connects the provider to a random server. (Update Supported)",
				MarkdownDescription: "If true, connects the provider to a random server. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"create_home_directory": schema.BoolAttribute{
				Description:         "Automatically create the home directory on the first login. (Update Supported)",
				MarkdownDescription: "Automatically create the home directory on the first login. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "If true, enables the LDAP provider. (Update Supported)",
				MarkdownDescription: "If true, enables the LDAP provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"enumerate_groups": schema.BoolAttribute{
				Description:         "If true, allows the provider to enumerate groups. (Update Supported)",
				MarkdownDescription: "If true, allows the provider to enumerate groups. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"enumerate_users": schema.BoolAttribute{
				Description:         "If true, allows the provider to enumerate users. (Update Supported)",
				MarkdownDescription: "If true, allows the provider to enumerate users. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"ignore_tls_errors": schema.BoolAttribute{
				Description:         "If true, continues over secure connections even if identity checks fail. (Update Supported)",
				MarkdownDescription: "If true, continues over secure connections even if identity checks fail. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"normalize_groups": schema.BoolAttribute{
				Description:         "Normalizes group names to lowercase before look up. (Update Supported)",
				MarkdownDescription: "Normalizes group names to lowercase before look up. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"normalize_users": schema.BoolAttribute{
				Description:         "Normalizes user names to lowercase before look up. (Update Supported)",
				MarkdownDescription: "Normalizes user names to lowercase before look up. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"require_secure_connection": schema.BoolAttribute{
				Description:         "Determines whether to continue over a non-TLS connection. (Update Supported)",
				MarkdownDescription: "Determines whether to continue over a non-TLS connection. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"restrict_findable": schema.BoolAttribute{
				Description:         "If true, checks the provider for filtered lists of findable and unfindable users and groups. (Update Supported)",
				MarkdownDescription: "If true, checks the provider for filtered lists of findable and unfindable users and groups. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"restrict_listable": schema.BoolAttribute{
				Description:         "If true, checks the provider for filtered lists of listable and unlistable users and groups. (Update Supported)",
				MarkdownDescription: "If true, checks the provider for filtered lists of listable and unlistable users and groups. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"system": schema.BoolAttribute{
				Description:         "If true, indicates that this provider instance was created by OneFS and cannot be removed. (Update Supported)",
				MarkdownDescription: "If true, indicates that this provider instance was created by OneFS and cannot be removed. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"bind_timeout": schema.Int64Attribute{
				Description:         "Specifies the timeout in seconds when binding to an LDAP server. (Update Supported)",
				MarkdownDescription: "Specifies the timeout in seconds when binding to an LDAP server. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 3600),
				},
			},
			"check_online_interval": schema.Int64Attribute{
				Description:         "Specifies the time in seconds between provider online checks. (Update Supported)",
				MarkdownDescription: "Specifies the time in seconds between provider online checks. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 3600),
				},
			},
			"search_timeout": schema.Int64Attribute{
				Description:         "Specifies the search timeout period in seconds. (Update Supported)",
				MarkdownDescription: "Specifies the search timeout period in seconds. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(10, 3600),
				},
			},
			"findable_groups": schema.ListAttribute{
				Description:         "Specifies the list of groups that can be resolved. (Update Supported)",
				MarkdownDescription: "Specifies the list of groups that can be resolved. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"findable_users": schema.ListAttribute{
				Description:         "Specifies the list of users that can be resolved. (Update Supported)",
				MarkdownDescription: "Specifies the list of users that can be resolved. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"listable_groups": schema.ListAttribute{
				Description:         "Specifies the groups that can be viewed in the provider. (Update Supported)",
				MarkdownDescription: "Specifies the groups that can be viewed in the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"listable_users": schema.ListAttribute{
				Description:         "Specifies the users that can be viewed in the provider. (Update Supported)",
				MarkdownDescription: "Specifies the users that can be viewed in the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"unfindable_groups": schema.ListAttribute{
				Description:         "Specifies the groups that cannot be resolved by the provider. (Update Supported)",
				MarkdownDescription: "Specifies the groups that cannot be resolved by the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"unfindable_users": schema.ListAttribute{
				Description:         "Specifies users that cannot be resolved by the provider. (Update Supported)",
				MarkdownDescription: "Specifies users that cannot be resolved by the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"unlistable_groups": schema.ListAttribute{
				Description:         "Specifies a group that cannot be listed by the provider. (Update Supported)",
				MarkdownDescription: "Specifies a group that cannot be listed by the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"unlistable_users": schema.ListAttribute{
				Description:         "Specifies a user that cannot be listed by the provider. (Update Supported)",
				MarkdownDescription: "Specifies a user that cannot be listed by the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 255)),
					listvalidator.SizeBetween(0, 4294967295),
				},
			},
			"alternate_security_identities_attribute": schema.StringAttribute{
				Description:         "Specifies the attribute name used when searching for alternate security identities. (Update Supported)",
				MarkdownDescription: "Specifies the attribute name used when searching for alternate security identities. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"bind_dn": schema.StringAttribute{
				Description:         "Specifies the distinguished name for binding to the LDAP server. (Update Supported)",
				MarkdownDescription: "Specifies the distinguished name for binding to the LDAP server. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"bind_mechanism": schema.StringAttribute{
				Description:         "Specifies which bind mechanism to use when connecting to an LDAP server. The only supported option is the 'simple' value. (Update Supported)",
				MarkdownDescription: "Specifies which bind mechanism to use when connecting to an LDAP server. The only supported option is the 'simple' value. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"certificate_authority_file": schema.StringAttribute{
				Description:         "Specifies the path to the root certificates file. (Update Supported)",
				MarkdownDescription: "Specifies the path to the root certificates file. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"cn_attribute": schema.StringAttribute{
				Description:         "Specifies the canonical name. (Update Supported)",
				MarkdownDescription: "Specifies the canonical name. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"crypt_password_attribute": schema.StringAttribute{
				Description:         "Specifies the hashed password value. (Update Supported)",
				MarkdownDescription: "Specifies the hashed password value. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"email_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Email attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Email attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 64),
				},
			},
			"gecos_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP GECOS attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP GECOS attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"gid_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP GID attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP GID attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"group_base_dn": schema.StringAttribute{
				Description:         "Specifies the distinguished name of the entry where LDAP searches for groups are started. (Update Supported)",
				MarkdownDescription: "Specifies the distinguished name of the entry where LDAP searches for groups are started. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"group_domain": schema.StringAttribute{
				Description:         "Specifies the domain for this provider through which groups are qualified. (Update Supported)",
				MarkdownDescription: "Specifies the domain for this provider through which groups are qualified. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"group_filter": schema.StringAttribute{
				Description:         "Specifies the LDAP filter for group objects. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP filter for group objects. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"group_members_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Group Members attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Group Members attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"group_search_scope": schema.StringAttribute{
				Description:         "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				MarkdownDescription: "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"home_directory_template": schema.StringAttribute{
				Description:         "Specifies the path to the home directory template. (Update Supported)",
				MarkdownDescription: "Specifies the path to the home directory template. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"homedir_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Homedir attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Homedir attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"login_shell": schema.StringAttribute{
				Description:         "Specifies the login shell path. (Update Supported)",
				MarkdownDescription: "Specifies the login shell path. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"member_lookup_method": schema.StringAttribute{
				Description:         "Sets the method by which group member lookups are performed. Use caution when changing this option directly. Acceptable values: \"default\", \"rfc2307bis\". (Update Supported)",
				MarkdownDescription: "Sets the method by which group member lookups are performed. Use caution when changing this option directly. Acceptable values: \"default\", \"rfc2307bis\". (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"member_of_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Query Member Of attribute, which performs reverse membership queries. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Query Member Of attribute, which performs reverse membership queries. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"name_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP UID attribute, which is used as the login name. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP UID attribute, which is used as the login name. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"netgroup_base_dn": schema.StringAttribute{
				Description:         "Specifies the distinguished name of the entry where LDAP searches for netgroups are started. (Update Supported)",
				MarkdownDescription: "Specifies the distinguished name of the entry where LDAP searches for netgroups are started. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"netgroup_filter": schema.StringAttribute{
				Description:         "Specifies the LDAP filter for netgroup objects. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP filter for netgroup objects. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"netgroup_members_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Netgroup Members attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Netgroup Members attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"netgroup_search_scope": schema.StringAttribute{
				Description:         "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				MarkdownDescription: "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"netgroup_triple_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Netgroup Triple attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Netgroup Triple attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"nt_password_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP NT Password attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP NT Password attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"ntlm_support": schema.StringAttribute{
				Description:         "Specifies which NTLM versions to support for users with NTLM-compatible credentials. Acceptable values: \"all\", \"v2only\", \"none\". (Update Supported)",
				MarkdownDescription: "Specifies which NTLM versions to support for users with NTLM-compatible credentials. Acceptable values: \"all\", \"v2only\", \"none\". (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"provider_domain": schema.StringAttribute{
				Description:         "Specifies the provider domain. (Update Supported)",
				MarkdownDescription: "Specifies the provider domain. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"search_scope": schema.StringAttribute{
				Description:         "Specifies the default depth from the base DN to perform LDAP searches. Acceptable values: \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				MarkdownDescription: "Specifies the default depth from the base DN to perform LDAP searches. Acceptable values: \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_expire_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the absolute date to expire the account. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the absolute date to expire the account. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_flag_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the section of the shadow map that is used to store the flag value. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the section of the shadow map that is used to store the flag value. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_inactive_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the number of days of inactivity that is allowed for the user. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the number of days of inactivity that is allowed for the user. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_last_change_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the last change of the shadow information. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the last change of the shadow information. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_max_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the maximum number of days a password can be valid. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the maximum number of days a password can be valid. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_min_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the minimum number of days between shadow changes. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the minimum number of days between shadow changes. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_user_filter": schema.StringAttribute{
				Description:         "Sets LDAP filter for shadow user objects. (Update Supported)",
				MarkdownDescription: "Sets LDAP filter for shadow user objects. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shadow_warning_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the number of days before the password expires to warn the user. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the number of days before the password expires to warn the user. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"shell_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP Shell attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP Shell attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"ssh_public_key_attribute": schema.StringAttribute{
				Description:         "Sets the attribute name that indicates the SSH Public Key for the user. (Update Supported)",
				MarkdownDescription: "Sets the attribute name that indicates the SSH Public Key for the user. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.StringAttribute{
				Description:         "Specifies the status of the provider. (Update Supported)",
				MarkdownDescription: "Specifies the status of the provider. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"tls_protocol_min": schema.StringAttribute{
				Description:         "Specifies the minimum TLS protocol version. (Update Supported)",
				MarkdownDescription: "Specifies the minimum TLS protocol version. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"uid_attribute": schema.StringAttribute{
				Description:         "Specifies the LDAP UID Number attribute. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP UID Number attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"unique_group_members_attribute": schema.StringAttribute{
				Description:         "Sets the LDAP Unique Group Members attribute. (Update Supported)",
				MarkdownDescription: "Sets the LDAP Unique Group Members attribute. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"user_base_dn": schema.StringAttribute{
				Description:         "Specifies the distinguished name of the entry at which to start LDAP searches for users. (Update Supported)",
				MarkdownDescription: "Specifies the distinguished name of the entry at which to start LDAP searches for users. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"user_domain": schema.StringAttribute{
				Description:         "Specifies the domain for this provider through which users are qualified. (Update Supported)",
				MarkdownDescription: "Specifies the domain for this provider through which users are qualified. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"user_filter": schema.StringAttribute{
				Description:         "Specifies the LDAP filter for user objects. (Update Supported)",
				MarkdownDescription: "Specifies the LDAP filter for user objects. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"user_search_scope": schema.StringAttribute{
				Description:         "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				MarkdownDescription: "Specifies the depth from the base DN to perform LDAP searches. Acceptable values: \"default\", \"base\", \"onelevel\", \"subtree\", \"children\". (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
		}}
}

// Configure configures the resource.
func (r *LdapProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *LdapProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating LdapProvider resource...")
	var plan models.LdapProviderModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ldapName := plan.Name.ValueString()
	if err := helper.CreateLdapProvider(ctx, r.client, &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating ldap provider - %s", ldapName),
			err.Error(),
		)
		return
	}

	ldapResponse, err := helper.GetLdapProvider(ctx, r.client, plan.Name.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting ldap provider after creation",
			err.Error(),
		)
		// if err, revert create
		_ = helper.DeleteLdapProvider(ctx, r.client, ldapName)
		return
	}

	if err := helper.UpdateLdapProviderResourceState(ctx, &plan, ldapResponse); err != nil {
		resp.Diagnostics.AddError("Error creating LdapProvider Resource",
			fmt.Sprintf("Error parsing LdapProvider resource state: %s", err.Error()))
		// if err, revert create
		_ = helper.DeleteLdapProvider(ctx, r.client, ldapName)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create LdapProvider resource")
}

// Read reads the resource state.
func (r *LdapProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading LdapProvider resource")
	var state models.LdapProviderModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ldapResponse, err := helper.GetLdapProvider(ctx, r.client, state.Name.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the LdapProvider - %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	// parse ldapProvider response to state ldapProvider model
	if err := helper.UpdateLdapProviderResourceState(ctx, &state, ldapResponse); err != nil {
		resp.Diagnostics.AddError("Error reading LdapProvider Resource",
			fmt.Sprintf("Error parsing LdapProvider resource state: %s", err.Error()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read LdapProvider resource")
}

// Update updates the resource state.
func (r *LdapProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating LdapProvider resource...")
	// Read Terraform plan into the model
	var plan models.LdapProviderModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.LdapProviderModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.UpdateLdapProvider(ctx, r.client, &state, &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating the LdapProvider resource - %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	ldapResponse, err := helper.GetLdapProvider(ctx, r.client, plan.Name.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the LdapProvider - %s", plan.Name.ValueString()),
			err.Error(),
		)
		return
	}

	if err := helper.UpdateLdapProviderResourceState(ctx, &plan, ldapResponse); err != nil {
		resp.Diagnostics.AddError("Error updating LdapProvider Resource",
			fmt.Sprintf("Error parsing LdapProvider resource state: %s", err.Error()))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Update LdapProvider resource")
}

// Delete deletes the resource.
func (r *LdapProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting LdapProvider resource")
	var state models.LdapProviderModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.DeleteLdapProvider(ctx, r.client, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting the LdapProvider - %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete LdapProvider resource")
}

// ImportState imports the resource state.
func (r *LdapProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing LdapProvider resource")
	var state models.LdapProviderModel

	ldapName := req.ID
	ldapResponse, err := helper.GetLdapProvider(ctx, r.client, ldapName, "")
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the LdapProvider - %s", ldapName),
			err.Error(),
		)
		return
	}

	// parse ldapProvider response to state ldapProvider model
	if err := helper.UpdateLdapProviderResourceState(ctx, &state, ldapResponse); err != nil {
		resp.Diagnostics.AddError("Error reading LdapProvider Resource",
			fmt.Sprintf("Error parsing LdapProvider resource state: %s", err.Error()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import LdapProvider resource")
}
