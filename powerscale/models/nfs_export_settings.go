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
