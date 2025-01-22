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

package provider

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &SmbShareDataSource{}
	_ datasource.DataSourceWithConfigure = &SmbShareDataSource{}
)

// NewSmbShareDataSource returns the SmbShare data source object.
func NewSmbShareDataSource() datasource.DataSource {
	return &SmbShareDataSource{}
}

// SmbShareDataSource defines the data source implementation.
type SmbShareDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SmbShareDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smb_share"
}

// Schema describes the data source arguments.
func (d *SmbShareDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing SMB shares from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale SMB shares provide clients network access to file system resources on the cluster.",
		Description: "This datasource is used to query the existing SMB shares from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale SMB shares provide clients network access to file system resources on the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"smb_shares": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of smb shares",
				MarkdownDescription: "List of smb shares",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_based_enumeration": schema.BoolAttribute{
							Description:         "Only enumerate files and folders the requesting user has access to.",
							MarkdownDescription: "Only enumerate files and folders the requesting user has access to.",
							Computed:            true,
						},
						"access_based_enumeration_root_only": schema.BoolAttribute{
							Description:         "Access-based enumeration on only the root directory of the share.",
							MarkdownDescription: "Access-based enumeration on only the root directory of the share.",
							Computed:            true,
						},
						"allow_delete_readonly": schema.BoolAttribute{
							Description:         "Allow deletion of read-only files in the share.",
							MarkdownDescription: "Allow deletion of read-only files in the share.",
							Computed:            true,
						},
						"allow_execute_always": schema.BoolAttribute{
							Description:         "Allows users to execute files they have read rights for.",
							MarkdownDescription: "Allows users to execute files they have read rights for.",
							Computed:            true,
						},
						"allow_variable_expansion": schema.BoolAttribute{
							Description:         "Allow automatic expansion of variables for home directories.",
							MarkdownDescription: "Allow automatic expansion of variables for home directories.",
							Computed:            true,
						},
						"auto_create_directory": schema.BoolAttribute{
							Description:         "Automatically create home directories.",
							MarkdownDescription: "Automatically create home directories.",
							Computed:            true,
						},
						"browsable": schema.BoolAttribute{
							Description:         "Share is visible in net view and the browse list.",
							MarkdownDescription: "Share is visible in net view and the browse list.",
							Computed:            true,
						},
						"ca_timeout": schema.Int64Attribute{
							Description:         "Persistent open timeout for the share.",
							MarkdownDescription: "Persistent open timeout for the share.",
							Computed:            true,
						},
						"ca_write_integrity": schema.StringAttribute{
							Description:         "Specify the level of write-integrity on continuously available shares.",
							MarkdownDescription: "Specify the level of write-integrity on continuously available shares.",
							Computed:            true,
						},
						"change_notify": schema.StringAttribute{
							Description:         "Level of change notification alerts on the share.",
							MarkdownDescription: "Level of change notification alerts on the share.",
							Computed:            true,
						},
						"continuously_available": schema.BoolAttribute{
							Description:         "Specify if persistent opens are allowed on the share.",
							MarkdownDescription: "Specify if persistent opens are allowed on the share.",
							Computed:            true,
						},
						"create_permissions": schema.StringAttribute{
							Description:         "Create permissions for new files and directories in share.",
							MarkdownDescription: "Create permissions for new files and directories in share.",
							Computed:            true,
						},
						"csc_policy": schema.StringAttribute{
							Description:         "Client-side caching policy for the shares.",
							MarkdownDescription: "Client-side caching policy for the shares.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Description for this SMB share.",
							MarkdownDescription: "Description for this SMB share.",
							Computed:            true,
						},
						"directory_create_mask": schema.Int64Attribute{
							Description:         "Directory create mask bits.",
							MarkdownDescription: "Directory create mask bits.",
							Computed:            true,
						},
						"directory_create_mode": schema.Int64Attribute{
							Description:         "Directory create mode bits.",
							MarkdownDescription: "Directory create mode bits.",
							Computed:            true,
						},
						"file_create_mask": schema.Int64Attribute{
							Description:         "File create mask bits.",
							MarkdownDescription: "File create mask bits.",
							Computed:            true,
						},
						"file_create_mode": schema.Int64Attribute{
							Description:         "File create mode bits.",
							MarkdownDescription: "File create mode bits.",
							Computed:            true,
						},
						"file_filter_extensions": schema.ListAttribute{
							Description:         "Specifies the list of file extensions.",
							MarkdownDescription: "Specifies the list of file extensions.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"file_filter_type": schema.StringAttribute{
							Description:         "Specifies if filter list is for deny or allow. Default is deny.",
							MarkdownDescription: "Specifies if filter list is for deny or allow. Default is deny.",
							Computed:            true,
						},
						"file_filtering_enabled": schema.BoolAttribute{
							Description:         "Enables file filtering on this zone.",
							MarkdownDescription: "Enables file filtering on this zone.",
							Computed:            true,
						},
						"hide_dot_files": schema.BoolAttribute{
							Description:         "Hide files and directories that begin with a period '.'.",
							MarkdownDescription: "Hide files and directories that begin with a period '.'.",
							Computed:            true,
						},
						"host_acl": schema.ListAttribute{
							Description:         "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
							MarkdownDescription: "An ACL expressing which hosts are allowed access. A deny clause must be the final entry.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"id": schema.StringAttribute{
							Description:         "Share ID.",
							MarkdownDescription: "Share ID.",
							Computed:            true,
						},
						"impersonate_guest": schema.StringAttribute{
							Description:         "Specify the condition in which user access is done as the guest account.",
							MarkdownDescription: "Specify the condition in which user access is done as the guest account.",
							Computed:            true,
						},
						"impersonate_user": schema.StringAttribute{
							Description:         "User account to be used as guest account.",
							MarkdownDescription: "User account to be used as guest account.",
							Computed:            true,
						},
						"inheritable_path_acl": schema.BoolAttribute{
							Description:         "Set the inheritable ACL on the share path.",
							MarkdownDescription: "Set the inheritable ACL on the share path.",
							Computed:            true,
						},
						"mangle_byte_start": schema.Int64Attribute{
							Description:         "Specifies the wchar_t starting point for automatic byte mangling.",
							MarkdownDescription: "Specifies the wchar_t starting point for automatic byte mangling.",
							Computed:            true,
						},
						"mangle_map": schema.ListAttribute{
							Description:         "Character mangle map.",
							MarkdownDescription: "Character mangle map.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"name": schema.StringAttribute{
							Description:         "Share name.",
							MarkdownDescription: "Share name.",
							Computed:            true,
						},
						"ntfs_acl_support": schema.BoolAttribute{
							Description:         "Support NTFS ACLs on files and directories.",
							MarkdownDescription: "Support NTFS ACLs on files and directories.",
							Computed:            true,
						},
						"oplocks": schema.BoolAttribute{
							Description:         "Support oplocks.",
							MarkdownDescription: "Support oplocks.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							Description:         "Path of share within /ifs.",
							MarkdownDescription: "Path of share within /ifs.",
							Computed:            true,
						},
						"permissions": schema.ListNestedAttribute{
							Description:         "Specifies an ordered list of permission modifications.",
							MarkdownDescription: "Specifies an ordered list of permission modifications.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"permission": schema.StringAttribute{
										Description:         "Specifies the file system rights that are allowed or denied.",
										MarkdownDescription: "Specifies the file system rights that are allowed or denied.",
										Computed:            true,
									},
									"permission_type": schema.StringAttribute{
										Description:         "Determines whether the permission is allowed or denied.",
										MarkdownDescription: "Determines whether the permission is allowed or denied.",
										Computed:            true,
									},
									"trustee": schema.SingleNestedAttribute{
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
						},
						"run_as_root": schema.ListNestedAttribute{
							Description:         "Allow account to run as root.",
							MarkdownDescription: "Allow account to run as root.",
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
						"smb3_encryption_enabled": schema.BoolAttribute{
							Description:         "Enables SMB3 encryption for the share.",
							MarkdownDescription: "Enables SMB3 encryption for the share.",
							Computed:            true,
						},
						"sparse_file": schema.BoolAttribute{
							Description:         "Enables sparse file.",
							MarkdownDescription: "Enables sparse file.",
							Computed:            true,
						},
						"strict_ca_lockout": schema.BoolAttribute{
							Description:         "Specifies if persistent opens would do strict lockout on the share.",
							MarkdownDescription: "Specifies if persistent opens would do strict lockout on the share.",
							Computed:            true,
						},
						"strict_flush": schema.BoolAttribute{
							Description:         "Handle SMB flush operations.",
							MarkdownDescription: "Handle SMB flush operations.",
							Computed:            true,
						},
						"strict_locking": schema.BoolAttribute{
							Description:         "Specifies whether byte range locks contend against SMB I/O.",
							MarkdownDescription: "Specifies whether byte range locks contend against SMB I/O.",
							Computed:            true,
						},
						"zid": schema.Int64Attribute{
							Description:         "Numeric ID of the access zone which contains this SMB share",
							MarkdownDescription: "Numeric ID of the access zone which contains this SMB share",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Description:         "Names to filter smb shares.",
						MarkdownDescription: "Names to filter smb shares.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
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
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"resume": schema.StringAttribute{
						Description: "Continue returning results from previous call using this token " +
							"(token should come from the previous call, resume cannot be used with other options).",
						MarkdownDescription: "Continue returning results from previous call using this token " +
							"(token should come from the previous call, resume cannot be used with other options).",
						Optional: true,
					},
					"limit": schema.Int32Attribute{
						Description:         "Return no more than this many results at once (see resume).",
						MarkdownDescription: "Return no more than this many results at once (see resume).",
						Optional:            true,
					},
					"offset": schema.Int32Attribute{
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
				},
			},
		},
	}
}

// Configure configures the resource.
func (d *SmbShareDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *SmbShareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Smb Share data source ")
	var sharesPlan models.SmbShareDatasource
	var sharesState models.SmbShareDatasource
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &sharesPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var shareNames []types.String
	totalSmbShares, err := helper.ListSmbShares(ctx, d.client, sharesPlan.SmbSharesFilter)
	if err != nil {
		errStr := constants.ListSmbShareErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading smb shares ",
			message)
		return
	}

	if sharesPlan.SmbSharesFilter != nil {
		shareNames = sharesPlan.SmbSharesFilter.Names
	}
	// if names are specified filter locally
	var filteredShares []models.SmbShareDatasourceEntity
	if len(shareNames) > 0 {
		sharesMap := make(map[string]powerscale.V7SmbShareExtended)
		for _, s := range *totalSmbShares {
			sharesMap[s.Name] = s
		}
		for _, name := range shareNames {
			if specifiedShare, ok := sharesMap[name.ValueString()]; ok {
				entity := models.SmbShareDatasourceEntity{}
				err := helper.CopyFields(ctx, specifiedShare, &entity)
				if err != nil {
					resp.Diagnostics.AddError("Error reading smb share datasource plan",
						fmt.Sprintf("Could not list smb shares with error: %s", err.Error()))
					return
				}
				filteredShares = append(filteredShares, entity)
			}
		}

		if filteredShares == nil {
			resp.Diagnostics.AddError("Error reading smb share datasource plan", "No shares found with the specified name(s)")
		}

	} else {
		entity := models.SmbShareDatasourceEntity{}
		for _, share := range *totalSmbShares {
			err := helper.CopyFields(ctx, share, &entity)
			if err != nil {
				resp.Diagnostics.AddError("Error reading smb share datasource plan",
					fmt.Sprintf("Could not list smb shares with error: %s", err.Error()))
				return
			}
			filteredShares = append(filteredShares, entity)
		}
	}
	//check if there is any error while getting the port group
	sharesState.ID = types.StringValue("1")
	sharesState.SmbSharesFilter = sharesPlan.SmbSharesFilter
	sharesState.SmbShares = filteredShares

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &sharesState)...)
}
