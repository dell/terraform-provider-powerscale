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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource              = &NfsExportSettingsResource{}
	_ resource.ResourceWithConfigure = &NfsExportSettingsResource{}
)

// NewNfsExportSettingsResource creates a new resource.
func NewNfsExportSettingsResource() resource.Resource {
	return &NfsExportSettingsResource{}
}

// NfsExportSettingsResource defines the resource implementation.
type NfsExportSettingsResource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (r *NfsExportSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_export_settings"
}

// Schema describes the data source arguments.
func (r *NfsExportSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `This resource is used to manage the NFS Export Settings of PowerScale Array. We can Create, Update and Delete the NFS Export Settings using this resource.  
Note that, NFS Export Settings is the native functionality of PowerScale. When creating the resource, we actually load NFS Export Settings from PowerScale to the resource.`,
		Description: `This resource is used to manage the NFS Export Settings of PowerScale Array. We can Create, Update and Delete the NFS Export Settings using this resource.  
Note that, NFS Export Settings is the native functionality of PowerScale. When creating the resource, we actually load NFS Export Settings from PowerScale to the resource.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Id of NFS Export settings. Readonly. ",
				MarkdownDescription: "Id of NFS Export settings. Readonly. ",
			},
			"symlinks": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"map_non_root": schema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Attributes: map[string]schema.Attribute{
					"user": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
						},
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "True if the user mapping is applied.",
						MarkdownDescription: "True if the user mapping is applied.",
					},
					"primary_group": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
						},
					},
					"secondary_groups": schema.ListNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"user",
											"group",
											"wellknown",
										),
									},
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 261),
									},
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 255),
									},
								},
							},
						},
					},
				},
			},
			"time_delta": schema.NumberAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the resolution of all time values that are returned to the clients",
				MarkdownDescription: "Specifies the resolution of all time values that are returned to the clients",
			},
			"return_32bit_file_ids": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).",
				MarkdownDescription: "Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).",
			},
			"case_preserving": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"link_max": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"map_failure": schema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Attributes: map[string]schema.Attribute{
					"primary_group": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"secondary_groups": schema.ListNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"user",
											"group",
											"wellknown",
										),
									},
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 261),
									},
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 255),
									},
								},
							},
						},
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "True if the user mapping is applied.",
						MarkdownDescription: "True if the user mapping is applied.",
					},
					"user": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
				},
			},
			"write_unstable_reply": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the stability disposition returned when an NFSv3+ unstable write is processed.",
				MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ unstable write is processed.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DATASYNC",
						"FILESYNC",
						"UNSTABLE",
					),
				},
			},
			"zone": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the zone in which the export is valid.",
				MarkdownDescription: "Specifies the zone in which the export is valid.",
			},
			"write_datasync_action": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the action to be taken when an NFSv3+ datasync write is requested.",
				MarkdownDescription: "Specifies the action to be taken when an NFSv3+ datasync write is requested.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DATASYNC",
						"FILESYNC",
						"UNSTABLE",
					),
				},
			},
			"read_only": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if the export is set to read-only.",
				MarkdownDescription: "True if the export is set to read-only.",
			},
			"all_dirs": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if all directories under the specified paths are mountable.",
				MarkdownDescription: "True if all directories under the specified paths are mountable.",
			},
			"readdirplus": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.",
				MarkdownDescription: "True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.",
			},
			"map_retry": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.",
				MarkdownDescription: "Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.",
			},
			"read_transfer_max_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"write_transfer_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"map_lookup_uid": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.",
				MarkdownDescription: "True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.",
			},
			"readdirplus_prefetch": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)",
				MarkdownDescription: "Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)",
			},
			"case_insensitive": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"map_all": schema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Attributes: map[string]schema.Attribute{
					"user": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "True if the user mapping is applied.",
						MarkdownDescription: "True if the user mapping is applied.",
					},
					"primary_group": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"secondary_groups": schema.ListNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"user",
											"group",
											"wellknown",
										),
									},
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 261),
									},
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 255),
									},
								},
							},
						},
					},
				},
			},
			"security_flavors": schema.ListAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the authentication types that are supported for this export.",
				MarkdownDescription: "Specifies the authentication types that are supported for this export.",
				ElementType:         types.StringType,
			},
			"chown_restricted": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"can_set_time": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"write_transfer_max_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"commit_asynchronous": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if NFS  commit  requests execute asynchronously.",
				MarkdownDescription: "True if NFS  commit  requests execute asynchronously.",
			},
			"encoding": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.",
				MarkdownDescription: "Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.",
			},
			"write_filesync_reply": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the stability disposition returned when an NFSv3+ filesync write is processed.",
				MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ filesync write is processed.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DATASYNC",
						"FILESYNC",
						"UNSTABLE",
					),
				},
			},
			"write_datasync_reply": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the stability disposition returned when an NFSv3+ datasync write is processed.",
				MarkdownDescription: "Specifies the stability disposition returned when an NFSv3+ datasync write is processed.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DATASYNC",
						"FILESYNC",
						"UNSTABLE",
					),
				},
			},
			"map_full": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.",
				MarkdownDescription: "True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.",
			},
			"snapshot": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the snapshot for all mounts.",
				MarkdownDescription: "Specifies the snapshot for all mounts.",
			},
			"max_file_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"read_transfer_multiple": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"write_unstable_action": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the action to be taken when an NFSv3+ unstable write is requested.",
				MarkdownDescription: "Specifies the action to be taken when an NFSv3+ unstable write is requested.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DATASYNC",
						"FILESYNC",
						"UNSTABLE",
					),
				},
			},
			"write_transfer_multiple": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"directory_transfer_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"read_transfer_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
				MarkdownDescription: "Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.",
			},
			"write_filesync_action": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the action to be taken when an NFSv3+ filesync write is requested.",
				MarkdownDescription: "Specifies the action to be taken when an NFSv3+ filesync write is requested.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"DATASYNC",
						"FILESYNC",
						"UNSTABLE",
					),
				},
			},
			"block_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the block size returned by the NFS statfs procedure.",
				MarkdownDescription: "Specifies the block size returned by the NFS statfs procedure.",
			},
			"map_root": schema.SingleNestedAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the users and groups to which non-root and root clients are mapped.",
				MarkdownDescription: "Specifies the users and groups to which non-root and root clients are mapped.",
				Attributes: map[string]schema.Attribute{
					"user": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
						},
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "True if the user mapping is applied.",
						MarkdownDescription: "True if the user mapping is applied.",
					},
					"primary_group": schema.SingleNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies the persona of the file group.",
						MarkdownDescription: "Specifies the persona of the file group.",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the type of persona, which must be combined with a name.",
								MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								Validators: []validator.String{
									stringvalidator.OneOf(
										"user",
										"group",
										"wellknown",
									),
								},
							},
							"id": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 261),
								},
							},
							"name": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								Description:         "Specifies the persona name, which must be combined with a type.",
								MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(0, 255),
								},
							},
						},
					},
					"secondary_groups": schema.ListNestedAttribute{
						Computed:            true,
						Optional:            true,
						Description:         "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						MarkdownDescription: "Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 261),
									},
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Validators: []validator.String{
										stringvalidator.LengthBetween(0, 255),
									},
								},
								"type": schema.StringAttribute{
									Computed:            true,
									Optional:            true,
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"user",
											"group",
											"wellknown",
										),
									},
								},
							},
						},
					},
				},
			},
			"name_max_size": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"no_truncate": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
				MarkdownDescription: "True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.",
			},
			"setattr_asynchronous": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "True if set attribute operations execute asynchronously.",
				MarkdownDescription: "True if set attribute operations execute asynchronously.",
			},
		},
	}
}

// Configure configures the resource.
func (r *NfsExportSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *NfsExportSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating NFS Export Settings resource...")

	var plan models.NfsexportsettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V2NfsSettingsExportSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs export settings",
			fmt.Sprintf("Could not read nfs export settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateNfsExportSettings(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs export settings",
			message,
		)
		return
	}

	settings, err := helper.GetNfsExportSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs export settings", message)
		return
	}

	var state models.NfsexportsettingsModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs export settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_export_settings")
	helper.ResolveSettingsDiff(ctx, plan, &state)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create nfs export settings resource")
}

// Read reads the resource state.
func (r *NfsExportSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Nfs Export Settings resource")

	var state models.NfsexportsettingsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var stateBackup models.NfsexportsettingsModel
	diags = req.State.Get(ctx, &stateBackup)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := helper.GetNfsExportSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs export settings", message)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs export settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_export_settings")
	helper.ResolveSettingsDiff(ctx, stateBackup, &state)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read nfs export settings resource")
}

// Update updates the resource state.
func (r *NfsExportSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating NFS Export Settings resource...")

	var plan models.NfsexportsettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V2NfsSettingsExportSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs export settings",
			fmt.Sprintf("Could not read nfs export settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateNfsExportSettings(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs export settings",
			message,
		)
		return
	}

	settings, err := helper.GetNfsExportSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs export settings", message)
		return
	}

	var state models.NfsexportsettingsModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs export settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_export_settings")
	helper.ResolveSettingsDiff(ctx, plan, &state)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update nfs export settings resource")
}

// Delete deletes the resource.
func (r *NfsExportSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Nfs Export Settings resource")
	var state models.NfsexportsettingsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Nfs export settings is the native functionality that cannot be deleted, so just remove state
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete nfs export settings resource")
}

// ImportState imports the resource state.
func (r *NfsExportSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Nfs Export Settings resource")

	var state models.NfsexportsettingsModel
	settings, err := helper.GetNfsExportSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsExportSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs export settings", message)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs export settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_export_settings")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import nfs export settings resource")
}
