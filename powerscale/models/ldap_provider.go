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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// LdapProviderModel describes the resource data model.
type LdapProviderModel struct {

	// Query param when creating and updating
	// Ignore unresolvable server URIs.
	IgnoreUnresolvableServerURIs types.Bool `tfsdk:"ignore_unresolvable_server_urls"`

	// Get params
	// Specifies the ID of the LDAP provider.
	ID types.String `tfsdk:"id"`
	// Specifies the name of the access zone in which this provider was created.
	ZoneName types.String `tfsdk:"zone_name"`

	// Create params
	// Groupnet identifier.
	Groupnet types.String `tfsdk:"groupnet"`

	// Create and Update params
	// Specifies the attribute name used when searching for alternate security identities.
	AlternateSecurityIdentitiesAttribute types.String `tfsdk:"alternate_security_identities_attribute"`
	// If true, enables authentication and identity management through the authentication provider.
	Authentication types.Bool `tfsdk:"authentication"`
	// If true, connects the provider to a random server.
	BalanceServers types.Bool `tfsdk:"balance_servers"`
	// Specifies the root of the tree in which to search identities.
	BaseDn types.String `tfsdk:"base_dn"`
	// Specifies the distinguished name for binding to the LDAP server.
	BindDn types.String `tfsdk:"bind_dn"`
	// Specifies which bind mechanism to use when connecting to an LDAP server. The only supported option is the 'simple' value.
	BindMechanism types.String `tfsdk:"bind_mechanism"`
	// Specifies the timeout in seconds when binding to an LDAP server.
	BindTimeout types.Int64 `tfsdk:"bind_timeout"`
	// Specifies the path to the root certificates file.
	CertificateAuthorityFile types.String `tfsdk:"certificate_authority_file"`
	// Specifies the time in seconds between provider online checks.
	CheckOnlineInterval types.Int64 `tfsdk:"check_online_interval"`
	// Specifies the canonical name.
	CnAttribute types.String `tfsdk:"cn_attribute"`
	// Automatically create the home directory on the first login.
	CreateHomeDirectory types.Bool `tfsdk:"create_home_directory"`
	// Specifies the hashed password value.
	CryptPasswordAttribute types.String `tfsdk:"crypt_password_attribute"`
	// Specifies the LDAP Email attribute.
	EmailAttribute types.String `tfsdk:"email_attribute"`
	// If true, enables the LDAP provider.
	Enabled types.Bool `tfsdk:"enabled"`
	// If true, allows the provider to enumerate groups.
	EnumerateGroups types.Bool `tfsdk:"enumerate_groups"`
	// If true, allows the provider to enumerate users.
	EnumerateUsers types.Bool `tfsdk:"enumerate_users"`
	// Specifies the list of groups that can be resolved.
	FindableGroups types.List `tfsdk:"findable_groups"`
	// Specifies the list of users that can be resolved.
	FindableUsers types.List `tfsdk:"findable_users"`
	// Specifies the LDAP GECOS attribute.
	GecosAttribute types.String `tfsdk:"gecos_attribute"`
	// Specifies the LDAP GID attribute.
	GIDAttribute types.String `tfsdk:"gid_attribute"`
	// Specifies the distinguished name of the entry where LDAP searches for groups are started.
	GroupBaseDn types.String `tfsdk:"group_base_dn"`
	// Specifies the domain for this provider through which groups are qualified.
	GroupDomain types.String `tfsdk:"group_domain"`
	// Specifies the LDAP filter for group objects.
	GroupFilter types.String `tfsdk:"group_filter"`
	// Specifies the LDAP Group Members attribute.
	GroupMembersAttribute types.String `tfsdk:"group_members_attribute"`
	// Specifies the depth from the base DN to perform LDAP searches.
	GroupSearchScope types.String `tfsdk:"group_search_scope"`
	// Specifies the path to the home directory template.
	HomeDirectoryTemplate types.String `tfsdk:"home_directory_template"`
	// Specifies the LDAP Homedir attribute.
	HomedirAttribute types.String `tfsdk:"homedir_attribute"`
	// If true, continues over secure connections even if identity checks fail.
	IgnoreTLSErrors types.Bool `tfsdk:"ignore_tls_errors"`
	// Specifies the groups that can be viewed in the provider.
	ListableGroups types.List `tfsdk:"listable_groups"`
	// Specifies the users that can be viewed in the provider.
	ListableUsers types.List `tfsdk:"listable_users"`
	// Specifies the login shell path.
	LoginShell types.String `tfsdk:"login_shell"`
	// Sets the method by which group member lookups are performed. Use caution when changing this option directly.
	MemberLookupMethod types.String `tfsdk:"member_lookup_method"`
	// Specifies the LDAP Query Member Of attribute, which performs reverse membership queries.
	MemberOfAttribute types.String `tfsdk:"member_of_attribute"`
	// Specifies the name of the LDAP provider.
	Name types.String `tfsdk:"name"`
	// Specifies the LDAP UID attribute, which is used as the login name.
	NameAttribute types.String `tfsdk:"name_attribute"`
	// Specifies the distinguished name of the entry where LDAP searches for netgroups are started.
	NetgroupBaseDn types.String `tfsdk:"netgroup_base_dn"`
	// Specifies the LDAP filter for netgroup objects.
	NetgroupFilter types.String `tfsdk:"netgroup_filter"`
	// Specifies the LDAP Netgroup Members attribute.
	NetgroupMembersAttribute types.String `tfsdk:"netgroup_members_attribute"`
	// Specifies the depth from the base DN to perform LDAP searches.
	NetgroupSearchScope types.String `tfsdk:"netgroup_search_scope"`
	// Specifies the LDAP Netgroup Triple attribute.
	NetgroupTripleAttribute types.String `tfsdk:"netgroup_triple_attribute"`
	// Normalizes group names to lowercase before look up.
	NormalizeGroups types.Bool `tfsdk:"normalize_groups"`
	// Normalizes user names to lowercase before look up.
	NormalizeUsers types.Bool `tfsdk:"normalize_users"`
	// Specifies the LDAP NT Password attribute.
	NtPasswordAttribute types.String `tfsdk:"nt_password_attribute"`
	// Specifies which NTLM versions to support for users with NTLM-compatible credentials.
	NtlmSupport types.String `tfsdk:"ntlm_support"`
	// Specifies the provider domain.
	ProviderDomain types.String `tfsdk:"provider_domain"`
	// Determines whether to continue over a non-TLS connection.
	RequireSecureConnection types.Bool `tfsdk:"require_secure_connection"`
	// If true, checks the provider for filtered lists of findable and unfindable users and groups.
	RestrictFindable types.Bool `tfsdk:"restrict_findable"`
	// If true, checks the provider for filtered lists of listable and unlistable users and groups.
	RestrictListable types.Bool `tfsdk:"restrict_listable"`
	// Specifies the default depth from the base DN to perform LDAP searches.
	SearchScope types.String `tfsdk:"search_scope"`
	// Specifies the search timeout period in seconds.
	SearchTimeout types.Int64 `tfsdk:"search_timeout"`
	// Specifies the server URIs.
	ServerUris types.List `tfsdk:"server_uris"`
	// Sets the attribute name that indicates the absolute date to expire the account.
	ShadowExpireAttribute types.String `tfsdk:"shadow_expire_attribute"`
	// Sets the attribute name that indicates the section of the shadow map that is used to store the flag value.
	ShadowFlagAttribute types.String `tfsdk:"shadow_flag_attribute"`
	// Sets the attribute name that indicates the number of days of inactivity that is allowed for the user.
	ShadowInactiveAttribute types.String `tfsdk:"shadow_inactive_attribute"`
	// Sets the attribute name that indicates the last change of the shadow information.
	ShadowLastChangeAttribute types.String `tfsdk:"shadow_last_change_attribute"`
	// Sets the attribute name that indicates the maximum number of days a password can be valid.
	ShadowMaxAttribute types.String `tfsdk:"shadow_max_attribute"`
	// Sets the attribute name that indicates the minimum number of days between shadow changes.
	ShadowMinAttribute types.String `tfsdk:"shadow_min_attribute"`
	// Sets LDAP filter for shadow user objects.
	ShadowUserFilter types.String `tfsdk:"shadow_user_filter"`
	// Sets the attribute name that indicates the number of days before the password expires to warn the user.
	ShadowWarningAttribute types.String `tfsdk:"shadow_warning_attribute"`
	// Specifies the LDAP Shell attribute.
	ShellAttribute types.String `tfsdk:"shell_attribute"`
	// Sets the attribute name that indicates the SSH Public Key for the user.
	SSHPublicKeyAttribute types.String `tfsdk:"ssh_public_key_attribute"`
	// Specifies the status of the provider.
	Status types.String `tfsdk:"status"`
	// If true, indicates that this provider instance was created by OneFS and cannot be removed.
	System types.Bool `tfsdk:"system"`
	// Specifies the minimum TLS protocol version.
	TLSProtocolMin types.String `tfsdk:"tls_protocol_min"`
	// Specifies the LDAP UID Number attribute.
	UIDAttribute types.String `tfsdk:"uid_attribute"`
	// Specifies the groups that cannot be resolved by the provider.
	UnfindableGroups types.List `tfsdk:"unfindable_groups"`
	// Specifies users that cannot be resolved by the provider.
	UnfindableUsers types.List `tfsdk:"unfindable_users"`
	// Sets the LDAP Unique Group Members attribute.
	UniqueGroupMembersAttribute types.String `tfsdk:"unique_group_members_attribute"`
	// Specifies a group that cannot be listed by the provider.
	UnlistableGroups types.List `tfsdk:"unlistable_groups"`
	// Specifies a user that cannot be listed by the provider.
	UnlistableUsers types.List `tfsdk:"unlistable_users"`
	// Specifies the distinguished name of the entry at which to start LDAP searches for users.
	UserBaseDn types.String `tfsdk:"user_base_dn"`
	// Specifies the domain for this provider through which users are qualified.
	UserDomain types.String `tfsdk:"user_domain"`
	// Specifies the LDAP filter for user objects.
	UserFilter types.String `tfsdk:"user_filter"`
	// Specifies the depth from the base DN to perform LDAP searches.
	UserSearchScope types.String `tfsdk:"user_search_scope"`

	// create and updating params -  Only available for PowerScale 9.5 and above.
	// This setting controls the behavior of the certificate revocation checking algorithm when the LDAP provider is presented with a digital certificate by an LDAP server.
	TLSRevocationCheckLevel types.String `tfsdk:"tls_revocation_check_level"`
	// Specifies the OCSP server URIs.
	OcspServerUris types.List `tfsdk:"ocsp_server_uris"`
}
