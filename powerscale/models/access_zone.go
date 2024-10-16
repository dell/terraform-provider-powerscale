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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// AccessZoneDataSourceModel describes the data source data model.
type AccessZoneDataSourceModel struct {
	ID          types.String            `tfsdk:"id"`
	AccessZones []AccessZoneDetailModel `tfsdk:"access_zones_details"`

	//filter
	AccessZoneFilter *AccessZoneFilterType `tfsdk:"filter"`
}

// AccessZoneDetailModel details of the accessZone.
type AccessZoneDetailModel struct {
	// Specifies an alternate system provider.
	AlternateSystemProvider types.String `tfsdk:"alternate_system_provider"`
	// Specifies the list of authentication providers available on this access zone.
	AuthProviders types.List `tfsdk:"auth_providers"`
	// Specifies amount of time in seconds to cache a user/group.
	CacheEntryExpiry types.Int64 `tfsdk:"cache_entry_expiry"`
	// Determines if a path is created when a path does not exist.
	CreatePath types.Bool `tfsdk:"create_path"`
	// Groupnet identifier
	Groupnet types.String `tfsdk:"groupnet"`
	// Specifies the permissions set on automatically created user home directories.
	HomeDirectoryUmask types.Int64 `tfsdk:"home_directory_umask"`
	// Specifies the system-assigned ID for the access zone. This value is returned when an access zone is created through the POST method
	ID types.String `tfsdk:"id"`
	// Specifies a list of users and groups that have read and write access to /ifs.
	IfsRestricted types.List `tfsdk:"ifs_restricted"`
	// Maps untrusted domains to this NetBIOS domain during authentication.
	MapUntrusted types.String `tfsdk:"map_untrusted"`
	// Specifies the access zone name.
	Name types.String `tfsdk:"name"`
	// Specifies number of seconds the negative cache entry is valid.
	NegativeCacheEntryExpiry types.Int64 `tfsdk:"negative_cache_entry_expiry"`
	// Specifies the NetBIOS name.
	NetbiosName types.String `tfsdk:"netbios_name"`
	// Specifies the access zone base directory path.
	Path types.String `tfsdk:"path"`
	// Specifies the skeleton directory that is used for user home directories.
	SkeletonDirectory types.String `tfsdk:"skeleton_directory"`
	// True if the access zone is built-in.
	System types.Bool `tfsdk:"system"`
	// Specifies the system provider for the access zone.
	SystemProvider types.String `tfsdk:"system_provider"`
	// Specifies the current ID mapping rules.
	UserMappingRules types.List `tfsdk:"user_mapping_rules"`
	// Specifies the access zone ID on the system.
	ZoneID types.Int64 `tfsdk:"zone_id"`
}

// V1AuthAccessAccessItemFileGroup IfsRestricted object.
type V1AuthAccessAccessItemFileGroup struct {
	// Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.
	ID types.String `tfsdk:"id"`
	// Specifies the persona name, which must be combined with a type.
	Name types.String `tfsdk:"name"`
	// Specifies the type of persona, which must be combined with a name.
	Type types.String `tfsdk:"type"`
}

// AccessZoneFilterType describes the filter data model.
type AccessZoneFilterType struct {
	Name []types.String `tfsdk:"name"`
	AlternateSystemProvider []types.String `tfsdk:"alternate_system_provider"`
	CacheEntryExpiry []types.Int64 `tfsdk:"cache_entry_expiry"`
	Groupnet []types.String `tfsdk:"groupnet"`
	HomeDirectoryUmask []types.Int64 `tfsdk:"home_directory_umask"`
	ID []types.String `tfsdk:"id"`
	MapUntrusted []types.String `tfsdk:"map_untrusted"`
	NegativeCacheEntryExpiry []types.Int64 `tfsdk:"negative_cache_entry_expiry"`
	NetbiosName []types.String `tfsdk:"netbios_name"`
	Path []types.String `tfsdk:"path"`
	SkeletonDirectory []types.String `tfsdk:"skeleton_directory"`
	SystemProvider []types.String `tfsdk:"system_provider"`
	ZoneID []types.Int64 `tfsdk:"zone_id"`
	CreatePath types.Bool `tfsdk:"create_path"`
	System types.Bool `tfsdk:"system"`
}

// AccessZoneResourceModel describes the resource data model.
type AccessZoneResourceModel struct {
	ID types.String `tfsdk:"id"`
	// Specifies an alternate system provider.
	AlternateSystemProvider types.String `tfsdk:"alternate_system_provider"`
	// An Optional user modifiable list to add new custome auth providers
	CustomAuthProviders types.List `tfsdk:"custom_auth_providers"`
	// Specifies the list of authentication providers available on this access zone.
	AuthProviders types.List `tfsdk:"auth_providers"`
	// Specifies amount of time in seconds to cache a user/group.
	CacheEntryExpiry types.Int64 `tfsdk:"cache_entry_expiry"`
	// Determines if a path is created when a path does not exist.
	CreatePath types.Bool `tfsdk:"create_path"`
	// Groupnet identifier
	Groupnet types.String `tfsdk:"groupnet"`
	// Specifies the permissions set on automatically created user home directories.
	HomeDirectoryUmask types.Int64 `tfsdk:"home_directory_umask"`
	// Specifies a list of users and groups that have read and write access to /ifs.
	IfsRestricted types.List `tfsdk:"ifs_restricted"`
	// Maps untrusted domains to this NetBIOS domain during authentication.
	MapUntrusted types.String `tfsdk:"map_untrusted"`
	// Specifies the access zone name.
	Name types.String `tfsdk:"name"`
	// Specifies number of seconds the negative cache entry is valid.
	NegativeCacheEntryExpiry types.Int64 `tfsdk:"negative_cache_entry_expiry"`
	// Specifies the NetBIOS name.
	NetbiosName types.String `tfsdk:"netbios_name"`
	// Specifies the access zone base directory path.
	Path types.String `tfsdk:"path"`
	// Specifies the skeleton directory that is used for user home directories.
	SkeletonDirectory types.String `tfsdk:"skeleton_directory"`
	// True if the access zone is built-in.
	System types.Bool `tfsdk:"system"`
	// Specifies the system provider for the access zone.
	SystemProvider types.String `tfsdk:"system_provider"`
	// Specifies the current ID mapping rules.
	UserMappingRules types.List `tfsdk:"user_mapping_rules"`
	// Specifies the access zone ID on the system.
	ZoneID types.Int64 `tfsdk:"zone_id"`
}
