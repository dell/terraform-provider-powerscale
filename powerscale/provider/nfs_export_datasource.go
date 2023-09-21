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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

var (
	_ datasource.DataSource              = &NfsExportDataSource{}
	_ datasource.DataSourceWithConfigure = &NfsExportDataSource{}
)

// NewNfsExportDataSource returns the NfsExport data source object.
func NewNfsExportDataSource() datasource.DataSource {
	return &NfsExportDataSource{}
}

// NfsExportDataSource defines the data source implementation.
type NfsExportDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NfsExportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_export"
}

// Schema describes the data source arguments.
func (d *NfsExportDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing NFS exports from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details / for further processing in resource block. " +
			"PowerScale provides an NFS server so you can share files on your cluster",
		Description: "This datasource is used to query the existing NFS exports from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details / for further processing in resource block. " +
			"PowerScale provides an NFS server so you can share files on your cluster",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"nfs_exports": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of nfs exports",
				MarkdownDescription: "List of nfs exports",
				NestedObject: schema.NestedAttributeObject{
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
						"clients": schema.ListAttribute{
							Description:         "Specifies the clients with root access to the export.",
							MarkdownDescription: "Specifies the clients with root access to the export.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"commit_asynchronous": schema.BoolAttribute{
							Description:         "True if NFS  commit  requests execute asynchronously.",
							MarkdownDescription: "True if NFS  commit  requests execute asynchronously.",
							Computed:            true,
						},
						"conflicting_paths": schema.ListAttribute{
							Description:         "Reports the paths that conflict with another export.",
							MarkdownDescription: "Reports the paths that conflict with another export.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"description": schema.StringAttribute{
							Description:         "Specifies the user-defined string that is used to identify the export.",
							MarkdownDescription: "Specifies the user-defined string that is used to identify the export.",
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
						"id": schema.Int64Attribute{
							Description:         "Specifies the system-assigned ID for the export. This ID is returned when an export is created through the POST method.",
							MarkdownDescription: "Specifies the system-assigned ID for the export. This ID is returned when an export is created through the POST method.",
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
									Optional:            true,
									Computed:            true,
								},
								"primary_group": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
											Computed:            true,
										},
									},
								},
								"secondary_groups": schema.ListNestedAttribute{
									Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									Optional:            true,
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												Optional:            true,
												Computed:            true,
											},
											"name": schema.StringAttribute{
												Description:         "Specifies the persona name, which must be combined with a type.",
												MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
												Optional:            true,
												Computed:            true,
											},
											"type": schema.StringAttribute{
												Description:         "Specifies the type of persona, which must be combined with a name.",
												MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
												Optional:            true,
												Computed:            true,
											},
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
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
									Optional:            true,
									Computed:            true,
								},
								"primary_group": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
											Computed:            true,
										},
									},
								},
								"secondary_groups": schema.ListNestedAttribute{
									Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									Optional:            true,
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												Optional:            true,
												Computed:            true,
											},
											"name": schema.StringAttribute{
												Description:         "Specifies the persona name, which must be combined with a type.",
												MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
												Optional:            true,
												Computed:            true,
											},
											"type": schema.StringAttribute{
												Description:         "Specifies the type of persona, which must be combined with a name.",
												MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
												Optional:            true,
												Computed:            true,
											},
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
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
									Optional:            true,
									Computed:            true,
								},
								"primary_group": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
											Computed:            true,
										},
									},
								},
								"secondary_groups": schema.ListNestedAttribute{
									Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									Optional:            true,
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												Optional:            true,
												Computed:            true,
											},
											"name": schema.StringAttribute{
												Description:         "Specifies the persona name, which must be combined with a type.",
												MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
												Optional:            true,
												Computed:            true,
											},
											"type": schema.StringAttribute{
												Description:         "Specifies the type of persona, which must be combined with a name.",
												MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
												Optional:            true,
												Computed:            true,
											},
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
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
									Optional:            true,
									Computed:            true,
								},
								"primary_group": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
											Computed:            true,
										},
									},
								},
								"secondary_groups": schema.ListNestedAttribute{
									Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
									Optional:            true,
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
												Optional:            true,
												Computed:            true,
											},
											"name": schema.StringAttribute{
												Description:         "Specifies the persona name, which must be combined with a type.",
												MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
												Optional:            true,
												Computed:            true,
											},
											"type": schema.StringAttribute{
												Description:         "Specifies the type of persona, which must be combined with a name.",
												MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
												Optional:            true,
												Computed:            true,
											},
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description:         "Specifies the persona of the file group.",
									MarkdownDescription: "Specifies the persona of the file group.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
											Optional:            true,
											Computed:            true,
										},
										"name": schema.StringAttribute{
											Description:         "Specifies the persona name, which must be combined with a type.",
											MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
											Optional:            true,
											Computed:            true,
										},
										"type": schema.StringAttribute{
											Description:         "Specifies the type of persona, which must be combined with a name.",
											MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
											Optional:            true,
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
						"paths": schema.ListAttribute{
							Description:         "Specifies the paths under /ifs that are exported.",
							MarkdownDescription: "Specifies the paths under /ifs that are exported.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"read_only": schema.BoolAttribute{
							Description:         "True if the export is set to read-only.",
							MarkdownDescription: "True if the export is set to read-only.",
							Computed:            true,
						},
						"read_only_clients": schema.ListAttribute{
							Description:         "Specifies the clients with read-only access to the export.",
							MarkdownDescription: "Specifies the clients with read-only access to the export.",
							Computed:            true,
							ElementType:         types.StringType,
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
						"read_write_clients": schema.ListAttribute{
							Description:         "Specifies the clients with both read and write access to the export, even when the export is set to read-only.",
							MarkdownDescription: "Specifies the clients with both read and write access to the export, even when the export is set to read-only.",
							Computed:            true,
							ElementType:         types.StringType,
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
						"root_clients": schema.ListAttribute{
							Description:         "Clients that have root access to the export.",
							MarkdownDescription: "Clients that have root access to the export.",
							Computed:            true,
							ElementType:         types.StringType,
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
						"unresolved_clients": schema.ListAttribute{
							Description:         "Reports clients that cannot be resolved.",
							MarkdownDescription: "Reports clients that cannot be resolved.",
							Computed:            true,
							ElementType:         types.StringType,
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
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"ids": schema.SetAttribute{
						Description:         "IDs to filter nfs exports.",
						MarkdownDescription: "IDs to filter nfs exports.",
						Optional:            true,
						ElementType:         types.Int64Type,
					},
					"paths": schema.SetAttribute{
						Description:         "Paths to filter nfs exports.",
						MarkdownDescription: "Paths to filter nfs exports.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"sort": schema.StringAttribute{
						Description:         "The field that will be used for sorting.",
						MarkdownDescription: "The field that will be used for sorting.",
						Optional:            true,
					},
					"zone": schema.StringAttribute{
						Description:         "Specifies which access zone to use.",
						MarkdownDescription: "Specifies which access zone to use.",
						Optional:            true,
					},
					"resume": schema.StringAttribute{
						Description: "Continue returning results from previous call using this token " +
							"(token should come from the previous call, resume cannot be used with other options).",
						MarkdownDescription: "Continue returning results from previous call using this token " +
							"(token should come from the previous call, resume cannot be used with other options).",
						Optional: true,
					},
					"limit": schema.Int64Attribute{
						Description:         "Return no more than this many results at once (see resume).",
						MarkdownDescription: "Return no more than this many results at once (see resume).",
						Optional:            true,
					},
					"offset": schema.Int64Attribute{
						Description:         "The position of the first item returned for a paginated query within the full result set.",
						MarkdownDescription: "The position of the first item returned for a paginated query within the full result set.",
						Optional:            true,
					},
					"scope": schema.StringAttribute{
						Description: "If specified as \"effective\" or not specified, all fields are returned. " +
							"If specified as \"user\", only fields with non-default values are shown. If specified as \"default\", the original values are returned.",
						MarkdownDescription: "If specified as \"effective\" or not specified, all fields are returned. " +
							"If specified as \"user\", only fields with non-default values are shown. If specified as \"default\", the original values are returned.",
						Optional: true,
					},
					"dir": schema.StringAttribute{
						Description:         "The direction of the sort.",
						MarkdownDescription: "The direction of the sort.",
						Optional:            true,
					},
					"path": schema.StringAttribute{
						Description:         "If specified, only exports that explicitly reference at least one of the given paths will be returned.",
						MarkdownDescription: "If specified, only exports that explicitly reference at least one of the given paths will be returned.",
						Optional:            true,
					},
					"check": schema.BoolAttribute{
						Description:         "Check for conflicts when listing exports.",
						MarkdownDescription: "Check for conflicts when listing exports.",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (d *NfsExportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *NfsExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Nfs Export data source ")
	var exportsPlan models.NfsExportDatasource
	var exportsState models.NfsExportDatasource
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &exportsPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	totalNfsExports, err := helper.ListNFSExports(ctx, d.client, exportsPlan.NfsExportsFilter)
	if err != nil {
		resp.Diagnostics.AddError("Error reading nfs export datasource plan",
			fmt.Sprintf("Could not list nfs exports with error: %s", err.Error()))
		return
	}
	var exportsIDs []types.Int64
	var paths []types.String
	if exportsPlan.NfsExportsFilter != nil {
		exportsIDs = exportsPlan.NfsExportsFilter.IDs
		paths = exportsPlan.NfsExportsFilter.Paths
	}
	filteredExports, err := helper.FilterExports(paths, exportsIDs, *totalNfsExports)
	if err != nil {
		errStr := constants.ListNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error filtering nfs export",
			message)
		return
	}

	for _, export := range filteredExports {
		entity := models.NfsExportDatasourceEntity{}
		err := helper.CopyFields(ctx, export, &entity)
		if err != nil {
			resp.Diagnostics.AddError("Error reading nfs export datasource plan",
				fmt.Sprintf("Could not list nfs exports with error: %s", err.Error()))
			return
		}
		exportsState.NfsExports = append(exportsState.NfsExports, entity)
	}
	exportsState.ID = types.StringValue("1")
	exportsState.NfsExportsFilter = exportsPlan.NfsExportsFilter
	resp.Diagnostics.Append(resp.State.Set(ctx, &exportsState)...)
}
