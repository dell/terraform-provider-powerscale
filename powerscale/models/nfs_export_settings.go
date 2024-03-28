package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// NfsexportsettingsModel defines the resource implementation.
type NfsexportsettingsModel struct {
	ID                    types.String `tfsdk:"id"`
	Symlinks              types.Bool   `tfsdk:"symlinks"`
	MapNonRoot            types.Object `tfsdk:"map_non_root"`
	TimeDelta             types.Number `tfsdk:"time_delta"`
	Return32bitFileIds    types.Bool   `tfsdk:"return_32bit_file_ids"`
	CasePreserving        types.Bool   `tfsdk:"case_preserving"`
	LinkMax               types.Int64  `tfsdk:"link_max"`
	MapFailure            types.Object `tfsdk:"map_failure"`
	WriteUnstableReply    types.String `tfsdk:"write_unstable_reply"`
	Zone                  types.String `tfsdk:"zone"`
	WriteDatasyncAction   types.String `tfsdk:"write_datasync_action"`
	ReadOnly              types.Bool   `tfsdk:"read_only"`
	AllDirs               types.Bool   `tfsdk:"all_dirs"`
	Readdirplus           types.Bool   `tfsdk:"readdirplus"`
	MapRetry              types.Bool   `tfsdk:"map_retry"`
	ReadTransferMaxSize   types.Int64  `tfsdk:"read_transfer_max_size"`
	WriteTransferSize     types.Int64  `tfsdk:"write_transfer_size"`
	MapLookupUID          types.Bool   `tfsdk:"map_lookup_uid"`
	ReaddirplusPrefetch   types.Int64  `tfsdk:"readdirplus_prefetch"`
	CaseInsensitive       types.Bool   `tfsdk:"case_insensitive"`
	MapAll                types.Object `tfsdk:"map_all"`
	SecurityFlavors       types.List   `tfsdk:"security_flavors"`
	ChownRestricted       types.Bool   `tfsdk:"chown_restricted"`
	CanSetTime            types.Bool   `tfsdk:"can_set_time"`
	WriteTransferMaxSize  types.Int64  `tfsdk:"write_transfer_max_size"`
	CommitAsynchronous    types.Bool   `tfsdk:"commit_asynchronous"`
	Encoding              types.String `tfsdk:"encoding"`
	WriteFilesyncReply    types.String `tfsdk:"write_filesync_reply"`
	WriteDatasyncReply    types.String `tfsdk:"write_datasync_reply"`
	MapFull               types.Bool   `tfsdk:"map_full"`
	Snapshot              types.String `tfsdk:"snapshot"`
	MaxFileSize           types.Int64  `tfsdk:"max_file_size"`
	ReadTransferMultiple  types.Int64  `tfsdk:"read_transfer_multiple"`
	WriteUnstableAction   types.String `tfsdk:"write_unstable_action"`
	WriteTransferMultiple types.Int64  `tfsdk:"write_transfer_multiple"`
	DirectoryTransferSize types.Int64  `tfsdk:"directory_transfer_size"`
	ReadTransferSize      types.Int64  `tfsdk:"read_transfer_size"`
	WriteFilesyncAction   types.String `tfsdk:"write_filesync_action"`
	BlockSize             types.Int64  `tfsdk:"block_size"`
	MapRoot               types.Object `tfsdk:"map_root"`
	NameMaxSize           types.Int64  `tfsdk:"name_max_size"`
	NoTruncate            types.Bool   `tfsdk:"no_truncate"`
	SetattrAsynchronous   types.Bool   `tfsdk:"setattr_asynchronous"`
}

// NfsSettingsExportDatasource represents the data source implementation.
type NfsSettingsExportDatasource struct {
	ID                      types.String                       `tfsdk:"id"`
	NfsSettingsExport       *V2NfsSettingsExportSettings       `tfsdk:"nfs_export_settings"`
	NfsSettingsExportFilter *NfsSettingsExportDatasourceFilter `tfsdk:"filter"`
}

// NfsSettingsExportDatasourceFilter holds filter conditions.
type NfsSettingsExportDatasourceFilter struct {
	Zone  types.String `tfsdk:"zone"`
	Scope types.String `tfsdk:"scope"`
}

// V2NfsSettingsExportSettings Specifies configuration values for NFS exports.
type V2NfsSettingsExportSettings struct {
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
	// True if NFS  commit  requests execute asynchronously.
	CommitAsynchronous types.Bool `tfsdk:"commit_asynchronous"`
	// Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.
	DirectoryTransferSize types.Int64 `tfsdk:"directory_transfer_size"`
	// Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.
	Encoding types.String `tfsdk:"encoding"`
	// Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	LinkMax types.Int64 `tfsdk:"link_max"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapAll *V2NfsExportMapAll `tfsdk:"map_all"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapFailure *V2NfsExportMapAll `tfsdk:"map_failure"`
	// True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.
	MapFull types.Bool `tfsdk:"map_full"`
	// True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.
	MapLookupUID types.Bool `tfsdk:"map_lookup_uid"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapNonRoot *V2NfsExportMapAll `tfsdk:"map_non_root"`
	// Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.
	MapRetry types.Bool `tfsdk:"map_retry"`
	// Specifies the users and groups to which non-root and root clients are mapped.
	MapRoot *V2NfsExportMapAll `tfsdk:"map_root"`
	// Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	MaxFileSize types.Int64 `tfsdk:"max_file_size"`
	// Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NameMaxSize types.Int64 `tfsdk:"name_max_size"`
	// True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NoTruncate types.Bool `tfsdk:"no_truncate"`
	// True if the export is set to read-only.
	ReadOnly types.Bool `tfsdk:"read_only"`
	// Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMaxSize types.Int64 `tfsdk:"read_transfer_max_size"`
	// Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMultiple types.Int64 `tfsdk:"read_transfer_multiple"`
	// Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferSize types.Int64 `tfsdk:"read_transfer_size"`
	// True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.
	Readdirplus types.Bool `tfsdk:"readdirplus"`
	// Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)
	ReaddirplusPrefetch types.Int64 `tfsdk:"readdirplus_prefetch"`
	// Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).
	Return32bitFileIds types.Bool `tfsdk:"return_32bit_file_ids"`
	// Specifies the authentication types that are supported for this export.
	SecurityFlavors types.List `tfsdk:"security_flavors"`
	// True if set attribute operations execute asynchronously.
	SetattrAsynchronous types.Bool `tfsdk:"setattr_asynchronous"`
	// Specifies the snapshot for all mounts.
	Snapshot types.String `tfsdk:"snapshot"`
	// True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.
	Symlinks types.Bool `tfsdk:"symlinks"`
	// Specifies the resolution of all time values that are returned to the clients
	TimeDelta *float32 `tfsdk:"time_delta"`
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
