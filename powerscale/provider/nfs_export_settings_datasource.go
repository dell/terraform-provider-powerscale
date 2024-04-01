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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &NfsExportSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &NfsExportSettingsDataSource{}
)

// NewNfsExportSettingsDataSource creates a new cluster email settings data source.
func NewNfsExportSettingsDataSource() datasource.DataSource {
	return &NfsExportSettingsDataSource{}
}

// NfsExportSettingsDataSource defines the data source implementation.
type NfsExportSettingsDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NfsExportSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_export_settings"
}

// Schema describes the data source arguments.
func (d *NfsExportSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the NFS Export Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Description:         "This datasource is used to query the NFS Export Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"nfs_export_settings": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "NFS Export Settings",
				MarkdownDescription: "NFS Export Settings",
				Attributes: map[string]schema.Attribute{
					"all_dirs": schema.BoolAttribute{
						Description:         "True if all directories under the specified paths are mountable.",
						MarkdownDescription: "True if all directories under the specified paths are mountable.",
						Computed:            true,
					},
					"block_size": schema.Int64Attribute{
						Description:         "Specifies the block size returned by the NFS statfs procedure.",
						MarkdownDescription: "Specifies the block size returned by the NFS statfs procedure.",
						Computed:            true,
					},
					"can_set_time": schema.BoolAttribute{
						Description:         "True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"case_insensitive": schema.BoolAttribute{
						Description:         "True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"case_preserving": schema.BoolAttribute{
						Description:         "True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"chown_restricted": schema.BoolAttribute{
						Description:         "True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"commit_asynchronous": schema.BoolAttribute{
						Description:         "True if NFS  commit  requests execute asynchronously.",
						MarkdownDescription: "True if NFS  commit  requests execute asynchronously.",
						Computed:            true,
					},
					"directory_transfer_size": schema.Int64Attribute{
						Description:         "Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"encoding": schema.StringAttribute{
						Description:         "Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.",
						MarkdownDescription: "Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.",
						Computed:            true,
					},
					"link_max": schema.Int64Attribute{
						Description:         "Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"map_all": schema.SingleNestedAttribute{
						Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
						MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         "True if the user mapping is applied.",
								MarkdownDescription: "True if the user mapping is applied.",
								Computed:            true,
							},
							"primary_group": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
							"secondary_groups": schema.ListNestedAttribute{
								Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Computed:            true,
										},
									},
								},
							},
							"user": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
						},
					},
					"map_failure": schema.SingleNestedAttribute{
						Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
						MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         "True if the user mapping is applied.",
								MarkdownDescription: "True if the user mapping is applied.",
								Computed:            true,
							},
							"primary_group": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
							"secondary_groups": schema.ListNestedAttribute{
								Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Computed:            true,
										},
									},
								},
							},
							"user": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
						},
					},
					"map_full": schema.BoolAttribute{
						Description:         "True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.",
						MarkdownDescription: "True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.",
						Computed:            true,
					},
					"map_lookup_uid": schema.BoolAttribute{
						Description:         "True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.",
						MarkdownDescription: "True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.",
						Computed:            true,
					},
					"map_non_root": schema.SingleNestedAttribute{
						Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
						MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         "True if the user mapping is applied.",
								MarkdownDescription: "True if the user mapping is applied.",
								Computed:            true,
							},
							"primary_group": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
							"secondary_groups": schema.ListNestedAttribute{
								Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Computed:            true,
										},
									},
								},
							},
							"user": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
						},
					},
					"map_retry": schema.BoolAttribute{
						Description:         "Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.",
						MarkdownDescription: "Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.",
						Computed:            true,
					},
					"map_root": schema.SingleNestedAttribute{
						Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
						MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description:         "True if the user mapping is applied.",
								MarkdownDescription: "True if the user mapping is applied.",
								Computed:            true,
							},
							"primary_group": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
							"secondary_groups": schema.ListNestedAttribute{
								Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Computed:            true,
										},
									},
								},
							},
							"user": schema.SingleNestedAttribute{
								Description:         "Specifies the persona of the file group.",
								MarkdownDescription: "Specifies the persona of the file group.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
						},
					},
					"max_file_size": schema.Int64Attribute{
						Description:         "Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"name_max_size": schema.Int64Attribute{
						Description:         "Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"no_truncate": schema.BoolAttribute{
						Description:         "True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						MarkdownDescription: "True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
						Computed:            true,
					},
					"read_only": schema.BoolAttribute{
						Description:         "True if the export is set to read-only.",
						MarkdownDescription: "True if the export is set to read-only.",
						Computed:            true,
					},
					"read_transfer_max_size": schema.Int64Attribute{
						Description:         "Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"read_transfer_multiple": schema.Int64Attribute{
						Description:         "Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"read_transfer_size": schema.Int64Attribute{
						Description:         "Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"readdirplus": schema.BoolAttribute{
						Description:         "True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.",
						MarkdownDescription: "True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.",
						Computed:            true,
					},
					"readdirplus_prefetch": schema.Int64Attribute{
						Description:         "Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)",
						MarkdownDescription: "Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)",
						Computed:            true,
					},
					"return_32bit_file_ids": schema.BoolAttribute{
						Description:         "Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).",
						MarkdownDescription: "Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).",
						Computed:            true,
					},
					"security_flavors": schema.ListAttribute{
						Description:         "Specifies the authentication types that are supported for this export.",
						MarkdownDescription: "Specifies the authentication types that are supported for this export.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"setattr_asynchronous": schema.BoolAttribute{
						Description:         "True if set attribute operations execute asynchronously.",
						MarkdownDescription: "True if set attribute operations execute asynchronously.",
						Computed:            true,
					},
					"snapshot": schema.StringAttribute{
						Description:         "Specifies the snapshot for all mounts.",
						MarkdownDescription: "Specifies the snapshot for all mounts.",
						Computed:            true,
					},
					"symlinks": schema.BoolAttribute{
						Description:         "True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"time_delta": schema.NumberAttribute{
						Description:         "Specifies the resolution of all time values that are returned to the clients",
						MarkdownDescription: "Specifies the resolution of all time values that are returned to the clients",
						Computed:            true,
					},
					"write_datasync_action": schema.StringAttribute{
						Description:         "Specifies the action to be taken when an NFSv3+ datasync write is requested.",
						MarkdownDescription: "Specifies the action to be taken when an NFSv3+ datasync write is requested.",
						Computed:            true,
					},
					"write_datasync_reply": schema.StringAttribute{
						Description:         "Specifies the stability disposition returned when an NFSv3+ datasync write is processed.",
						MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ datasync write is processed.",
						Computed:            true,
					},
					"write_filesync_action": schema.StringAttribute{
						Description:         "Specifies the action to be taken when an NFSv3+ filesync write is requested.",
						MarkdownDescription: "Specifies the action to be taken when an NFSv3+ filesync write is requested.",
						Computed:            true,
					},
					"write_filesync_reply": schema.StringAttribute{
						Description:         "Specifies the stability disposition returned when an NFSv3+ filesync write is processed.",
						MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ filesync write is processed.",
						Computed:            true,
					},
					"write_transfer_max_size": schema.Int64Attribute{
						Description:         "Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"write_transfer_multiple": schema.Int64Attribute{
						Description:         "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"write_transfer_size": schema.Int64Attribute{
						Description:         "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						MarkdownDescription: "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
						Computed:            true,
					},
					"write_unstable_action": schema.StringAttribute{
						Description:         "Specifies the action to be taken when an NFSv3+ unstable write is requested.",
						MarkdownDescription: "Specifies the action to be taken when an NFSv3+ unstable write is requested.",
						Computed:            true,
					},
					"write_unstable_reply": schema.StringAttribute{
						Description:         "Specifies the stability disposition returned when an NFSv3+ unstable write is processed.",
						MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ unstable write is processed.",
						Computed:            true,
					},
					"zone": schema.StringAttribute{
						Description:         "Specifies the zone in which the export is valid.",
						MarkdownDescription: "Specifies the zone in which the export is valid.",
						Computed:            true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"zone": schema.StringAttribute{
						Description:         "Specifies which access zone to use.",
						MarkdownDescription: "Specifies which access zone to use.",
						Optional:            true,
					},
					"scope": schema.StringAttribute{
						Description: "If specified as \"effective\" or not specified, all fields are returned. " +
							"If specified as \"user\", only fields with non-default values are shown. If specified as \"default\", the original values are returned.",
						MarkdownDescription: "If specified as \"effective\" or not specified, all fields are returned. " +
							"If specified as \"user\", only fields with non-default values are shown. If specified as \"default\", the original values are returned.",
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *NfsExportSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NfsExportSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Nfs Export Settings data source ")

	var settingsPlan models.NfsSettingsExportDatasource
	var settingsState models.NfsSettingsExportDatasource
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &settingsPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nfsExportSettings, err := helper.FilterNfsExportSettings(ctx, d.client, settingsPlan.NfsSettingsExportFilter)

	if err != nil {
		errStr := constants.ReadNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs export settings",
			message,
		)
		return
	}

	var settings models.V2NfsSettingsExportSettings
	err = helper.CopyFields(ctx, nfsExportSettings.GetSettings(), &settings)
	if err != nil {
		resp.Diagnostics.AddError("Error copying fields of nfs export settings datasource", err.Error())
		return
	}

	settingsState.NfsSettingsExport = &settings
	settingsState.ID = types.StringValue("nfs_export_settings")
	settingsState.NfsSettingsExportFilter = settingsPlan.NfsSettingsExportFilter

	resp.Diagnostics.Append(resp.State.Set(ctx, &settingsState)...)
	tflog.Info(ctx, "Done with Read Nfs Export Settings data source ")
}
