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

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NfsExportResource Specifies configuration values for NFS exports.
type NfsExportResource struct {
	// query param
	// When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.
	Scope types.String `tfsdk:"scope"`

	// create and modify param
	// If true, the export will be created even if it conflicts with another export.
	Force types.Bool `tfsdk:"force"`
	// Ignore unresolvable hosts.
	IgnoreUnresolvableHosts types.Bool `tfsdk:"ignore_unresolvable_hosts"`
	// Ignore conflicts with existing exports.
	IgnoreConflicts types.Bool `tfsdk:"ignore_conflicts"`
	// Ignore nonexistent or otherwise bad paths.
	IgnoreBadPaths types.Bool `tfsdk:"ignore_bad_paths"`
	// Ignore invalid users.
	IgnoreBadAuth types.Bool `tfsdk:"ignore_bad_auth"`

	// True if all directories under the specified paths are mountable.
	AllDirs types.Bool `tfsdk:"all_dirs"`
	// Specifies the block size returned by the NFS statfs procedure.
	BlockSize types.Int64 `tfsdk:"block_size"`
	// True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CanSetTime types.Bool `tfsdk:"can_set_time"`
	// True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CaseInsensitive types.Bool `tfsdk:"case_insensitive"`
	// True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CasePreserving types.Bool `tfsdk:"case_preserving"`
	// True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	ChownRestricted types.Bool `tfsdk:"chown_restricted"`
	// Specifies the clients with root access to the export.
	Clients types.List `tfsdk:"clients"`
	// True if NFS  commit  requests execute asynchronously.
	CommitAsynchronous types.Bool `tfsdk:"commit_asynchronous"`
	// Reports the paths that conflict with another export.
	ConflictingPaths types.List `tfsdk:"conflicting_paths"`
	// Specifies the user-defined string that is used to identify the export.
	Description types.String `tfsdk:"description"`
	// Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.
	DirectoryTransferSize types.Int64 `tfsdk:"directory_transfer_size"`
	// Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.
	Encoding types.String `tfsdk:"encoding"`
	// Specifies the system-assigned ID for the export. This ID is returned when an export is created through the POST method.
	ID types.Int64 `tfsdk:"id"`
	// Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	LinkMax types.Int64 `tfsdk:"link_max"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapAll types.Object `tfsdk:"map_all"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapFailure types.Object `tfsdk:"map_failure"`
	// True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.
	MapFull types.Bool `tfsdk:"map_full"`
	// True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.
	MapLookupUID types.Bool `tfsdk:"map_lookup_uid"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapNonRoot types.Object `tfsdk:"map_non_root"`
	// Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.
	MapRetry types.Bool `tfsdk:"map_retry"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapRoot types.Object `tfsdk:"map_root"`
	// Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	MaxFileSize types.Int64 `tfsdk:"max_file_size"`
	// Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NameMaxSize types.Int64 `tfsdk:"name_max_size"`
	// True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NoTruncate types.Bool `tfsdk:"no_truncate"`
	// Specifies the paths under /ifs that are exported.
	Paths types.List `tfsdk:"paths"`
	// True if the export is set to read-only.
	ReadOnly types.Bool `tfsdk:"read_only"`
	// Specifies the clients with read-only access to the export.
	ReadOnlyClients types.List `tfsdk:"read_only_clients"`
	// Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMaxSize types.Int64 `tfsdk:"read_transfer_max_size"`
	// Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMultiple types.Int64 `tfsdk:"read_transfer_multiple"`
	// Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferSize types.Int64 `tfsdk:"read_transfer_size"`
	// Specifies the clients with both read and write access to the export, even when the export is set to read-only.
	ReadWriteClients types.List `tfsdk:"read_write_clients"`
	// True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.
	Readdirplus types.Bool `tfsdk:"readdirplus"`
	// Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)
	ReaddirplusPrefetch types.Int64 `tfsdk:"readdirplus_prefetch"`
	// Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).
	Return32bitFileIds types.Bool `tfsdk:"return_32bit_file_ids"`
	// Clients that have root access to the export.
	RootClients types.List `tfsdk:"root_clients"`
	// Specifies the authentication types that are supported for this export.
	SecurityFlavors types.List `tfsdk:"security_flavors"`
	// True if set attribute operations execute asynchronously.
	SetattrAsynchronous types.Bool `tfsdk:"setattr_asynchronous"`
	// Specifies the snapshot for all mounts.
	Snapshot types.String `tfsdk:"snapshot"`
	// True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.
	Symlinks types.Bool `tfsdk:"symlinks"`
	// Specifies the resolution of all time values that are returned to the clients
	TimeDelta types.Number `tfsdk:"time_delta"`
	// Reports clients that cannot be resolved.
	UnresolvedClients types.List `tfsdk:"unresolved_clients"`
	// Specifies the action to be taken when an NFSv3+ datasync write is requested.
	WriteDatasyncAction types.String `tfsdk:"write_datasync_action"`
	// Specifies the stability disposition returned when an NFSv3+ datasync write is processed.
	WriteDatasyncReply types.String `tfsdk:"write_datasync_reply"`
	// Specifies the action to be taken when an NFSv3+ filesync write is requested.
	WriteFilesyncAction types.String `tfsdk:"write_filesync_action"`
	// Specifies the stability disposition returned when an NFSv3+ filesync write is processed.
	WriteFilesyncReply types.String `tfsdk:"write_filesync_reply"`
	// Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferMaxSize types.Int64 `tfsdk:"write_transfer_max_size"`
	// Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferMultiple types.Int64 `tfsdk:"write_transfer_multiple"`
	// Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferSize types.Int64 `tfsdk:"write_transfer_size"`
	// Specifies the action to be taken when an NFSv3+ unstable write is requested.
	WriteUnstableAction types.String `tfsdk:"write_unstable_action"`
	// Specifies the stability disposition returned when an NFSv3+ unstable write is processed.
	WriteUnstableReply types.String `tfsdk:"write_unstable_reply"`
	// Specifies the zone in which the export is valid.
	Zone CaseInsensitiveStringValue `tfsdk:"zone"`
}

// NfsExportDatasource holds nfs exports datasource schema attribute details.
type NfsExportDatasource struct {
	ID               types.String                `tfsdk:"id"`
	NfsExports       []NfsExportDatasourceEntity `tfsdk:"nfs_exports"`
	NfsExportsFilter *NfsExportDatasourceFilter  `tfsdk:"filter"`
}

// NfsExportDatasourceFilter holds filter conditions.
type NfsExportDatasourceFilter struct {
	// supported by api
	Sort   types.String `tfsdk:"sort"`
	Zone   types.String `tfsdk:"zone"`
	Resume types.String `tfsdk:"resume"`
	Scope  types.String `tfsdk:"scope"`
	Limit  types.Int32  `tfsdk:"limit"`
	Offset types.Int32  `tfsdk:"offset"`
	Path   types.String `tfsdk:"path"`
	Check  types.Bool   `tfsdk:"check"`
	Dir    types.String `tfsdk:"dir"`
	// custom id & path list
	IDs   []types.Int64  `tfsdk:"ids"`
	Paths []types.String `tfsdk:"paths"`
}

// NfsExportDatasourceEntity Specifies entity values for NFS exports.
type NfsExportDatasourceEntity struct {
	// True if all directories under the specified paths are mountable.
	AllDirs types.Bool `tfsdk:"all_dirs"`
	// Specifies the block size returned by the NFS statfs procedure.
	BlockSize types.Int64 `tfsdk:"block_size"`
	// True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CanSetTime types.Bool `tfsdk:"can_set_time"`
	// True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CaseInsensitive types.Bool `tfsdk:"case_insensitive"`
	// True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CasePreserving types.Bool `tfsdk:"case_preserving"`
	// True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	ChownRestricted types.Bool `tfsdk:"chown_restricted"`
	// Specifies the clients with root access to the export.
	Clients types.List `tfsdk:"clients"`
	// True if NFS  commit  requests execute asynchronously.
	CommitAsynchronous types.Bool `tfsdk:"commit_asynchronous"`
	// Reports the paths that conflict with another export.
	ConflictingPaths types.List `tfsdk:"conflicting_paths"`
	// Specifies the user-defined string that is used to identify the export.
	Description types.String `tfsdk:"description"`
	// Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.
	DirectoryTransferSize types.Int64 `tfsdk:"directory_transfer_size"`
	// Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.
	Encoding types.String `tfsdk:"encoding"`
	// Specifies the system-assigned ID for the export. This ID is returned when an export is created through the POST method.
	ID types.Int64 `tfsdk:"id"`
	// Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	LinkMax types.Int64 `tfsdk:"link_max"`
	//
	MapAll *V2NfsExportMapAll `tfsdk:"map_all"`
	//
	MapFailure *V2NfsExportMapAll `tfsdk:"map_failure"`
	// True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.
	MapFull types.Bool `tfsdk:"map_full"`
	// True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.
	MapLookupUID types.Bool `tfsdk:"map_lookup_uid"`
	//
	MapNonRoot *V2NfsExportMapAll `tfsdk:"map_non_root"`
	// Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.
	MapRetry types.Bool `tfsdk:"map_retry"`
	//
	MapRoot *V2NfsExportMapAll `tfsdk:"map_root"`
	// Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	MaxFileSize types.Int64 `tfsdk:"max_file_size"`
	// Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NameMaxSize types.Int64 `tfsdk:"name_max_size"`
	// True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NoTruncate types.Bool `tfsdk:"no_truncate"`
	// Specifies the paths under /ifs that are exported.
	Paths types.List `tfsdk:"paths"`
	// True if the export is set to read-only.
	ReadOnly types.Bool `tfsdk:"read_only"`
	// Specifies the clients with read-only access to the export.
	ReadOnlyClients types.List `tfsdk:"read_only_clients"`
	// Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMaxSize types.Int64 `tfsdk:"read_transfer_max_size"`
	// Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMultiple types.Int64 `tfsdk:"read_transfer_multiple"`
	// Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferSize types.Int64 `tfsdk:"read_transfer_size"`
	// Specifies the clients with both read and write access to the export, even when the export is set to read-only.
	ReadWriteClients types.List `tfsdk:"read_write_clients"`
	// True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.
	Readdirplus types.Bool `tfsdk:"readdirplus"`
	// Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)
	ReaddirplusPrefetch types.Int64 `tfsdk:"readdirplus_prefetch"`
	// Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).
	Return32bitFileIds types.Bool `tfsdk:"return_32bit_file_ids"`
	// Clients that have root access to the export.
	RootClients types.List `tfsdk:"root_clients"`
	// Specifies the authentication types that are supported for this export.
	SecurityFlavors types.List `tfsdk:"security_flavors"`
	// True if set attribute operations execute asynchronously.
	SetattrAsynchronous types.Bool `tfsdk:"setattr_asynchronous"`
	// Specifies the snapshot for all mounts.
	Snapshot types.String `tfsdk:"snapshot"`
	// True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.
	Symlinks types.Bool `tfsdk:"symlinks"`
	// Specifies the resolution of all time values that are returned to the clients
	TimeDelta types.Number `tfsdk:"time_delta"`
	// Reports clients that cannot be resolved.
	UnresolvedClients types.List `tfsdk:"unresolved_clients"`
	// Specifies the action to be taken when an NFSv3+ datasync write is requested.
	WriteDatasyncAction types.String `tfsdk:"write_datasync_action"`
	// Specifies the stability disposition returned when an NFSv3+ datasync write is processed.
	WriteDatasyncReply types.String `tfsdk:"write_datasync_reply"`
	// Specifies the action to be taken when an NFSv3+ filesync write is requested.
	WriteFilesyncAction types.String `tfsdk:"write_filesync_action"`
	// Specifies the stability disposition returned when an NFSv3+ filesync write is processed.
	WriteFilesyncReply types.String `tfsdk:"write_filesync_reply"`
	// Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferMaxSize types.Int64 `tfsdk:"write_transfer_max_size"`
	// Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferMultiple types.Int64 `tfsdk:"write_transfer_multiple"`
	// Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferSize types.Int64 `tfsdk:"write_transfer_size"`
	// Specifies the action to be taken when an NFSv3+ unstable write is requested.
	WriteUnstableAction types.String `tfsdk:"write_unstable_action"`
	// Specifies the stability disposition returned when an NFSv3+ unstable write is processed.
	WriteUnstableReply types.String `tfsdk:"write_unstable_reply"`
	// Specifies the zone in which the export is valid.
	Zone types.String `tfsdk:"zone"`
}

// V2NfsExportMapAll Specifies the users and groups to which non-root and root clients are mapped.
type V2NfsExportMapAll struct {
	// True if the user mapping is applied.
	Enabled types.Bool `tfsdk:"enabled"`
	//
	PrimaryGroup *V1AuthAccessAccessItemFileGroup `tfsdk:"primary_group"`
	// Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.
	SecondaryGroups []V11NfsExportMapAllSecondaryGroupsInner `tfsdk:"secondary_groups"`
	//
	User *V1AuthAccessAccessItemFileGroup `tfsdk:"user"`
}

// V11NfsExportMapAllSecondaryGroupsInner Specifies properties for a persona, which consists of either a 'type' and a 'name' or an 'ID'.
type V11NfsExportMapAllSecondaryGroupsInner struct {
	// Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.
	ID types.String `tfsdk:"id"`
	// Specifies the persona name, which must be combined with a type.
	Name types.String `tfsdk:"name"`
	// Specifies the type of persona, which must be combined with a name.
	Type types.String `tfsdk:"type"`
}
