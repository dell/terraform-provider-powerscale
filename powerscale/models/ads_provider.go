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

// AdsProviderDataSourceModel describes the data source data model.
type AdsProviderDataSourceModel struct {
	ID           types.String             `tfsdk:"id"`
	AdsProviders []AdsProviderDetailModel `tfsdk:"ads_providers_details"`

	// Filters
	AdsProviderFilter *AdsProviderFilterType `tfsdk:"filter"`
}

// AdsProviderDetailModel Specifies the properties for an ADS authentication provider.
type AdsProviderDetailModel struct {
	// Allocates an ID for an unmapped Active Directory (ADS) group. ADS groups without GIDs can be proactively assigned a GID by the ID mapper. If the ID mapper option is disabled, GIDs are not proactively assigned, and when a primary group for a user does not include a GID, the system may allocate one.
	AllocateGids types.Bool `tfsdk:"allocate_gids"`
	// Allocates a user ID for an unmapped Active Directory (ADS) user. ADS users without UIDs can be proactively assigned a UID by the ID mapper. IF the ID mapper option is disabled, UIDs are not proactively assigned, and when an identify for a user does not include a UID, the system may allocate one.
	AllocateUids types.Bool `tfsdk:"allocate_uids"`
	// Enables lookup of unqualified user names in the primary domain.
	AssumeDefaultDomain types.Bool `tfsdk:"assume_default_domain"`
	// Enables authentication and identity management through the authentication provider.
	Authentication types.Bool `tfsdk:"authentication"`
	// Specifies the time in seconds between provider online checks.
	CheckOnlineInterval types.Int64 `tfsdk:"check_online_interval"`
	// Specifies the current time for the domain controllers.
	ControllerTime types.Int64 `tfsdk:"controller_time"`
	// Automatically creates a home directory on the first login.
	CreateHomeDirectory types.Bool `tfsdk:"create_home_directory"`
	// Sends an alert if the domain goes offline.
	DomainOfflineAlerts types.Bool `tfsdk:"domain_offline_alerts"`
	// Get duplicate SPNs in the provider domain
	DupSpns types.List `tfsdk:"dup_spns"`
	// List of additional SPNs to expect beyond what automatic checking routines might find
	ExtraExpectedSpns types.List `tfsdk:"extra_expected_spns"`
	// Sets list of groups that can be resolved.
	FindableGroups types.List `tfsdk:"findable_groups"`
	// Sets list of users that can be resolved.
	FindableUsers types.List `tfsdk:"findable_users"`
	// Specifies the Active Directory forest.
	Forest types.String `tfsdk:"forest"`
	// Groupnet identifier.
	Groupnet types.String `tfsdk:"groupnet"`
	// Specifies the path to the home directory template.
	HomeDirectoryTemplate types.String `tfsdk:"home_directory_template"`
	// Specifies the fully qualified hostname stored in the machine account.
	Hostname types.String `tfsdk:"hostname"`
	// Specifies the ID of the Active Directory provider instance.
	ID types.String `tfsdk:"id"`
	// If set to true, ignores all trusted domains.
	IgnoreAllTrusts types.Bool `tfsdk:"ignore_all_trusts"`
	// Includes trusted domains when 'ignore_all_trusts' is set to false.
	IgnoredTrustedDomains types.List `tfsdk:"ignored_trusted_domains"`
	// Includes trusted domains when 'ignore_all_trusts' is set to true.
	IncludeTrustedDomains types.List `tfsdk:"include_trusted_domains"`
	// Specifies Active Directory provider instance.
	Instance types.String `tfsdk:"instance"`
	// Enables encryption and signing on LDAP requests.
	LdapSignAndSeal types.Bool `tfsdk:"ldap_sign_and_seal"`
	// Specifies the login shell path.
	LoginShell types.String `tfsdk:"login_shell"`
	// Limits user and group lookups to the specified domains.
	LookupDomains types.List `tfsdk:"lookup_domains"`
	// Looks up AD groups in other providers before allocating a group ID.
	LookupGroups types.Bool `tfsdk:"lookup_groups"`
	// Normalizes AD group names to lowercase before look up.
	LookupNormalizeGroups types.Bool `tfsdk:"lookup_normalize_groups"`
	// Normalize AD user names to lowercase before look up.
	LookupNormalizeUsers types.Bool `tfsdk:"lookup_normalize_users"`
	// Looks up AD users in other providers before allocating a user ID.
	LookupUsers types.Bool `tfsdk:"lookup_users"`
	// Specifies the machine account name when creating a SAM account with Active Directory.
	MachineAccount types.String `tfsdk:"machine_account"`
	// Enables periodic changes of the machine password for security.
	MachinePasswordChanges types.Bool `tfsdk:"machine_password_changes"`
	// Sets maximum age of a password in seconds.
	MachinePasswordLifespan types.Int64 `tfsdk:"machine_password_lifespan"`
	// Specifies the Active Directory provider name.
	Name types.String `tfsdk:"name"`
	// Specifies the NetBIOS domain name associated with the machine account.
	NetbiosDomain types.String `tfsdk:"netbios_domain"`
	// Specifies the domain controller for which the node has affinity.
	NodeDcAffinity types.String `tfsdk:"node_dc_affinity"`
	// Specifies the timeout for the domain controller for which the local node has affinity.
	NodeDcAffinityTimeout types.Int64 `tfsdk:"node_dc_affinity_timeout"`
	// Enables the Active Directory provider to respond to 'getpwent' and 'getgrent' requests.
	NssEnumeration types.Bool `tfsdk:"nss_enumeration"`
	// Specifies the AD domain to which the provider is joined.
	PrimaryDomain types.String `tfsdk:"primary_domain"`
	// Check the provider for filtered lists of findable and unfindable users and groups.
	RestrictFindable types.Bool `tfsdk:"restrict_findable"`
	// The maximum amount of time (in seconds) an RPC call to Active Directory is allowed to take.
	RPCCallTimeout types.Int64 `tfsdk:"rpc_call_timeout"`
	// The number of retries attempted when a call to Active Directory fails due to network error.
	ServerRetryLimit types.Int64 `tfsdk:"server_retry_limit"`
	// Specifies whether to support RFC 2307 attributes on ADS domain controllers.
	SfuSupport types.String `tfsdk:"sfu_support"`
	// Specifies the site for the Active Directory.
	Site types.String `tfsdk:"site"`
	// Specifies the status of the provider.
	Status types.String `tfsdk:"status"`
	// Stores SFU mappings permanently in the ID mapper.
	StoreSfuMappings types.Bool `tfsdk:"store_sfu_mappings"`
	// If set to true, indicates that this provider instance was created by OneFS and cannot be removed.
	System types.Bool `tfsdk:"system"`
	// Specifies groups that cannot be resolved by the provider.
	UnfindableGroups types.List `tfsdk:"unfindable_groups"`
	// Specifies users that cannot be resolved by the provider.
	UnfindableUsers types.List `tfsdk:"unfindable_users"`
	// Specifies the name of the access zone in which this provider was created.
	ZoneName types.String `tfsdk:"zone_name"`
}

// AdsProviderFilterType describes the filter data model.
type AdsProviderFilterType struct {
	Names []types.String `tfsdk:"names"`
	Scope types.String   `tfsdk:"scope"`
}
