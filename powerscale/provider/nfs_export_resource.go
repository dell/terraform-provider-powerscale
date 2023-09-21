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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NfsExportResource{}
var _ resource.ResourceWithImportState = &NfsExportResource{}

// NewNfsExportResource creates a new resource.
func NewNfsExportResource() resource.Resource {
	return &NfsExportResource{}
}

// NfsExportResource defines the resource implementation.
type NfsExportResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r NfsExportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_export"
}

// Schema describes the resource arguments.
func (r *NfsExportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the NFS export entity of PowerScale Array. " +
			"PowerScale provides an NFS server so you can share files on your cluster. " +
			"We can Create, Update and Delete the NFS export using this resource. We can also import an existing NFS export from PowerScale array.",
		Description: "This resource is used to manage the NFS export entity of PowerScale Array. " +
			"PowerScale provides an NFS server so you can share files on your cluster. " +
			"We can Create, Update and Delete the NFS export using this resource. We can also import an existing NFS export from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"scope": schema.StringAttribute{
				Description:         "When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.",
				MarkdownDescription: "When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.",
				Optional:            true,
			},
			"force": schema.BoolAttribute{
				Description:         "If true, the export will be created even if it conflicts with another export.",
				MarkdownDescription: "If true, the export will be created even if it conflicts with another export.",
				Optional:            true,
			},
			"ignore_unresolvable_hosts": schema.BoolAttribute{
				Description:         "Ignore unresolvable hosts.",
				MarkdownDescription: "Ignore unresolvable hosts.",
				Optional:            true,
			},
			"ignore_conflicts": schema.BoolAttribute{
				Description:         "Ignore conflicts with existing exports.",
				MarkdownDescription: "Ignore conflicts with existing exports.",
				Optional:            true,
			},
			"ignore_bad_paths": schema.BoolAttribute{
				Description:         "Ignore nonexistent or otherwise bad paths.",
				MarkdownDescription: "Ignore nonexistent or otherwise bad paths.",
				Optional:            true,
			},
			"ignore_bad_auth": schema.BoolAttribute{
				Description:         "Ignore invalid users.",
				MarkdownDescription: "Ignore invalid users.",
				Optional:            true,
			},

			"all_dirs": schema.BoolAttribute{
				Description:         "True if all directories under the specified paths are mountable.",
				MarkdownDescription: "True if all directories under the specified paths are mountable.",
				Optional:            true,
				Computed:            true,
			},
			"block_size": schema.Int64Attribute{
				Description:         "Specifies the block size returned by the NFS statfs procedure.",
				MarkdownDescription: "Specifies the block size returned by the NFS statfs procedure.",
				Optional:            true,
				Computed:            true,
			},
			"can_set_time": schema.BoolAttribute{
				Description:         "True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				Optional:            true,
				Computed:            true,
			},
			"case_insensitive": schema.BoolAttribute{
				Description:         "True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				Optional:            true,
				Computed:            true,
			},
			"case_preserving": schema.BoolAttribute{
				Description:         "True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				Optional:            true,
				Computed:            true,
			},
			"chown_restricted": schema.BoolAttribute{
				Description:         "True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				Optional:            true,
				Computed:            true,
			},
			"clients": schema.ListAttribute{
				Description:         "Specifies the clients with root access to the export.",
				MarkdownDescription: "Specifies the clients with root access to the export.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"commit_asynchronous": schema.BoolAttribute{
				Description:         "True if NFS  commit  requests execute asynchronously.",
				MarkdownDescription: "True if NFS  commit  requests execute asynchronously.",
				Optional:            true,
				Computed:            true,
			},
			"conflicting_paths": schema.ListAttribute{
				Description:         "Reports the paths that conflict with another export.",
				MarkdownDescription: "Reports the paths that conflict with another export.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				Description:         "Specifies the user-defined string that is used to identify the export.",
				MarkdownDescription: "Specifies the user-defined string that is used to identify the export.",
				Optional:            true,
				Computed:            true,
			},
			"directory_transfer_size": schema.Int64Attribute{
				Description:         "Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"encoding": schema.StringAttribute{
				Description:         "Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.",
				MarkdownDescription: "Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"map_all": schema.SingleNestedAttribute{
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Optional:            true,
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
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"map_lookup_uid": schema.BoolAttribute{
				Description:         "True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.",
				MarkdownDescription: "True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.",
				Optional:            true,
				Computed:            true,
			},
			"map_non_root": schema.SingleNestedAttribute{
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"map_root": schema.SingleNestedAttribute{
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Optional:            true,
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
				Optional:            true,
				Computed:            true,
			},
			"name_max_size": schema.Int64Attribute{
				Description:         "Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				Optional:            true,
				Computed:            true,
			},
			"no_truncate": schema.BoolAttribute{
				Description:         "True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				Optional:            true,
				Computed:            true,
			},
			"paths": schema.ListAttribute{
				Description:         "Specifies the paths under /ifs that are exported.",
				MarkdownDescription: "Specifies the paths under /ifs that are exported.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"read_only": schema.BoolAttribute{
				Description:         "True if the export is set to read-only.",
				MarkdownDescription: "True if the export is set to read-only.",
				Optional:            true,
				Computed:            true,
			},
			"read_only_clients": schema.ListAttribute{
				Description:         "Specifies the clients with read-only access to the export.",
				MarkdownDescription: "Specifies the clients with read-only access to the export.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"read_transfer_max_size": schema.Int64Attribute{
				Description:         "Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"read_transfer_multiple": schema.Int64Attribute{
				Description:         "Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"read_transfer_size": schema.Int64Attribute{
				Description:         "Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"read_write_clients": schema.ListAttribute{
				Description:         "Specifies the clients with both read and write access to the export, even when the export is set to read-only.",
				MarkdownDescription: "Specifies the clients with both read and write access to the export, even when the export is set to read-only.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"readdirplus": schema.BoolAttribute{
				Description:         "True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.",
				MarkdownDescription: "True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.",
				Optional:            true,
				Computed:            true,
			},
			"readdirplus_prefetch": schema.Int64Attribute{
				Description:         "Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)",
				MarkdownDescription: "Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)",
				Optional:            true,
				Computed:            true,
			},
			"return_32bit_file_ids": schema.BoolAttribute{
				Description:         "Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).",
				MarkdownDescription: "Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).",
				Optional:            true,
				Computed:            true,
			},
			"root_clients": schema.ListAttribute{
				Description:         "Clients that have root access to the export.",
				MarkdownDescription: "Clients that have root access to the export.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"security_flavors": schema.ListAttribute{
				Description:         "Specifies the authentication types that are supported for this export.",
				MarkdownDescription: "Specifies the authentication types that are supported for this export.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"setattr_asynchronous": schema.BoolAttribute{
				Description:         "True if set attribute operations execute asynchronously.",
				MarkdownDescription: "True if set attribute operations execute asynchronously.",
				Optional:            true,
				Computed:            true,
			},
			"snapshot": schema.StringAttribute{
				Description:         "Specifies the snapshot for all mounts.",
				MarkdownDescription: "Specifies the snapshot for all mounts.",
				Optional:            true,
				Computed:            true,
			},
			"symlinks": schema.BoolAttribute{
				Description:         "True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"time_delta": schema.NumberAttribute{
				Description:         "Specifies the resolution of all time values that are returned to the clients",
				MarkdownDescription: "Specifies the resolution of all time values that are returned to the clients",
				Optional:            true,
				Computed:            true,
			},
			"unresolved_clients": schema.ListAttribute{
				Description:         "Reports clients that cannot be resolved.",
				MarkdownDescription: "Reports clients that cannot be resolved.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"write_datasync_action": schema.StringAttribute{
				Description:         "Specifies the action to be taken when an NFSv3+ datasync write is requested.",
				MarkdownDescription: "Specifies the action to be taken when an NFSv3+ datasync write is requested.",
				Optional:            true,
				Computed:            true,
			},
			"write_datasync_reply": schema.StringAttribute{
				Description:         "Specifies the stability disposition returned when an NFSv3+ datasync write is processed.",
				MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ datasync write is processed.",
				Optional:            true,
				Computed:            true,
			},
			"write_filesync_action": schema.StringAttribute{
				Description:         "Specifies the action to be taken when an NFSv3+ filesync write is requested.",
				MarkdownDescription: "Specifies the action to be taken when an NFSv3+ filesync write is requested.",
				Optional:            true,
				Computed:            true,
			},
			"write_filesync_reply": schema.StringAttribute{
				Description:         "Specifies the stability disposition returned when an NFSv3+ filesync write is processed.",
				MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ filesync write is processed.",
				Optional:            true,
				Computed:            true,
			},
			"write_transfer_max_size": schema.Int64Attribute{
				Description:         "Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"write_transfer_multiple": schema.Int64Attribute{
				Description:         "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"write_transfer_size": schema.Int64Attribute{
				Description:         "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				Optional:            true,
				Computed:            true,
			},
			"write_unstable_action": schema.StringAttribute{
				Description:         "Specifies the action to be taken when an NFSv3+ unstable write is requested.",
				MarkdownDescription: "Specifies the action to be taken when an NFSv3+ unstable write is requested.",
				Optional:            true,
				Computed:            true,
			},
			"write_unstable_reply": schema.StringAttribute{
				Description:         "Specifies the stability disposition returned when an NFSv3+ unstable write is processed.",
				MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ unstable write is processed.",
				Optional:            true,
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "Specifies the zone in which the export is valid. Cannot be changed once set",
				MarkdownDescription: "Specifies the zone in which the export is valid. Cannot be changed once set",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure - defines configuration for nfs export resource.
func (r *NfsExportResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Create allocates the resource.
func (r NfsExportResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "creating nfs export")

	var exportPlan models.NfsExportResource
	diags := request.Plan.Get(ctx, &exportPlan)
	response.Diagnostics.Append(diags...)
	var exportPlanBackUp models.NfsExportResource
	diags = request.Plan.Get(ctx, &exportPlanBackUp)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	createResp, err := helper.CreateNFSExport(ctx, r.client, exportPlan)
	if err != nil {
		errStr := constants.CreateNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating nfs export ",
			message)
		return
	}
	exportID := int64(createResp.Id)
	tflog.Debug(ctx, fmt.Sprintf("nfs export %s created", strconv.FormatInt(exportID, 10)), map[string]interface{}{
		"nfsExportResponse": exportID,
	})

	exportPlan.ID = types.Int64Value(exportID)
	getExportResponse, err := helper.GetNFSExport(ctx, r.client, exportPlan)
	if err != nil {
		errStr := constants.GetNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating nfs export ",
			message)
		return
	}

	// update resource state according to response
	if len(getExportResponse.Exports) <= 0 {
		response.Diagnostics.AddError(
			"Error creating nfs export",
			fmt.Sprintf("Could not get created nfs export state %d with error: nfs export not found", exportID),
		)
		return
	}
	createdExport := getExportResponse.Exports[0]
	err = helper.CopyFieldsToNonNestedModel(ctx, createdExport, &exportPlan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating nfs export",
			fmt.Sprintf("Could not read nfs export %d with error: %s", exportID, err.Error()),
		)
		return
	}

	helper.ResolvePersonaDiff(ctx, exportPlanBackUp, &exportPlan)
	diags = response.State.Set(ctx, exportPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create nfs export completed")
}

// Read reads the resource state.
func (r NfsExportResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "reading nfs export")
	var exportState models.NfsExportResource
	diags := request.State.Get(ctx, &exportState)
	response.Diagnostics.Append(diags...)
	var exportStateBackUp models.NfsExportResource
	diags = request.State.Get(ctx, &exportStateBackUp)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	exportID := exportState.ID.ValueInt64()
	tflog.Debug(ctx, "calling get nfs export by ID", map[string]interface{}{
		"nfsExportID": exportID,
	})
	exportResponse, err := helper.GetNFSExport(ctx, r.client, exportState)
	if err != nil {
		errStr := constants.GetNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading nfs export ",
			message)
		return
	}

	if len(exportResponse.Exports) <= 0 {
		response.Diagnostics.AddError(
			"Error reading nfs export",
			fmt.Sprintf("Could not read nfs export %d from pscale with error: nfs export not found", exportID),
		)
		return
	}
	tflog.Debug(ctx, "updating read nfs export state", map[string]interface{}{
		"nfsExportResponse": exportResponse,
		"nfsExportState":    exportState,
	})
	err = helper.CopyFieldsToNonNestedModel(ctx, exportResponse.Exports[0], &exportState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error read nfs export",
			fmt.Sprintf("Could not read nfs export struct %d with error: %s", exportID, err.Error()),
		)
		return
	}

	helper.ResolvePersonaDiff(ctx, exportStateBackUp, &exportState)
	diags = response.State.Set(ctx, exportState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read nfs export completed")
}

// Update updates the resource state.
func (r NfsExportResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating nfs export")
	var exportPlan models.NfsExportResource
	diags := request.Plan.Get(ctx, &exportPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var exportState models.NfsExportResource
	diags = response.State.Get(ctx, &exportState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update nfs export", map[string]interface{}{
		"exportPlan":  exportPlan,
		"exportState": exportState,
	})

	if !exportPlan.Zone.IsUnknown() && exportPlan.Zone.ValueString() != exportState.Zone.ValueString() {
		response.Diagnostics.AddError("Error updating nfs export", "Do not change access zone once set")
		return
	}

	exportID := exportState.ID.ValueInt64()
	exportPlan.ID = exportState.ID
	err := helper.UpdateNFSExport(ctx, r.client, exportPlan)
	if err != nil {
		errStr := constants.UpdateNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating nfs export ",
			message)
		return
	}

	tflog.Debug(ctx, "calling get nfs export by ID on pscale client", map[string]interface{}{
		"nfsExportID": exportID,
	})
	// Use export plan to query updated export
	updatedShare, err := helper.GetNFSExport(ctx, r.client, exportPlan)
	if err != nil {
		errStr := constants.UpdateNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error getting nfs export ",
			message)
		return
	}

	if len(updatedShare.Exports) <= 0 {
		response.Diagnostics.AddError(
			"Error reading nfs export",
			fmt.Sprintf("Could not read nfs export %d from pscale with error: nfs export not found", exportID),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, updatedShare.Exports[0], &exportState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error read nfs export",
			fmt.Sprintf("Could not read nfs export struct %d with error: %s", exportID, err.Error()),
		)
		return
	}
	helper.ResolvePersonaDiff(ctx, exportPlan, &exportState)
	diags = response.State.Set(ctx, exportState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update nfs export completed")
}

// Delete deletes the resource.
func (r NfsExportResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting nfs export")
	var exportState models.NfsExportResource
	diags := request.State.Get(ctx, &exportState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	exportID := exportState.ID.ValueInt64()

	tflog.Debug(ctx, "calling delete nfs export on pscale client", map[string]interface{}{
		"nfsExportID": exportID,
	})
	err := helper.DeleteNFSExport(ctx, r.client, exportState)
	if err != nil {
		errStr := constants.DeleteNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error deleting nfs export ",
			message)
		return
	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete nfs export completed")
}

// ImportState imports the resource state.
func (r NfsExportResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	readNfsExport, err := helper.GetNFSExportByID(ctx, r.client, request.ID)
	if err != nil {
		errStr := constants.GetNfsExportErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading nfs export ",
			message)
		return
	}
	if len(readNfsExport.Exports) <= 0 {
		response.Diagnostics.AddError(
			"Error reading nfs export",
			fmt.Sprintf("Could not read nfs export %s from pscale with error: nfs export not found", request.ID),
		)
		return
	}
	var model models.NfsExportResource
	err = helper.CopyFieldsToNonNestedModel(ctx, readNfsExport.Exports[0], &model)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading nfs export",
			fmt.Sprintf("Could not set state for export %s with error: %s ",
				request.ID, err.Error()),
		)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}
}
