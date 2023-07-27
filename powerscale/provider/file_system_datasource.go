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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FileSystemDataSource{}

// NewFileSystemDataSource creates a new data source.
func NewFileSystemDataSource() datasource.DataSource {
	return &FileSystemDataSource{}
}

// FileSystemDataSource defines the data source implementation.
type FileSystemDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *FileSystemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem"
}

// Schema describes the data source arguments.
func (d *FileSystemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "FileSystem data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "FileSystem identifier",
				Computed:            true,
				Optional:            true,
			},
			"directory_path": schema.StringAttribute{
				MarkdownDescription: "FileSystem identifier",
				Required:            true,
			},
			"file_systems_details": schema.SingleNestedAttribute{
				Description:         "Details of the Filesystem",
				MarkdownDescription: "Details of the Filesystem",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"file_system_attributes": schema.ListNestedAttribute{
						Description:         "FileSystems Attributes",
						MarkdownDescription: "FileSystems Attributes",
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         "Attribute name",
									MarkdownDescription: "Attribute name",
									Computed:            true,
									Optional:            true,
								},
								"namespace": schema.StringAttribute{
									Description:         "Attribute namespace",
									MarkdownDescription: "Attribute namespace",
									Computed:            true,
									Optional:            true,
								},
								"value": schema.StringAttribute{
									Description:         "Attribute value",
									MarkdownDescription: "Attribute value",
									Computed:            true,
									Optional:            true,
								},
							},
						},
					},
					"file_system_quotas": schema.ListNestedAttribute{
						Description:         "Filesystem quotas",
						MarkdownDescription: "Filesystem quotas",
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "Quota Id",
									MarkdownDescription: "Quota Id",
									Computed:            true,
									Optional:            true,
								},
								"enforced": schema.BoolAttribute{
									Description:         "True if the quota provides enforcement, otherwise a accounting quota.",
									MarkdownDescription: "True if the quota provides enforcement, otherwise a accounting quota.",
									Computed:            true,
									Optional:            true,
								},
								"container": schema.BoolAttribute{
									Description:         "If true, SMB shares using the quota directory see the quota thresholds as share size.",
									MarkdownDescription: "If true, SMB shares using the quota directory see the quota thresholds as share size.",
									Computed:            true,
									Optional:            true,
								},
								"type": schema.StringAttribute{
									Description:         "The type of quota.",
									MarkdownDescription: "The type of quota.",
									Computed:            true,
									Optional:            true,
								},
								"path": schema.StringAttribute{
									Description:         "The path of quota.",
									MarkdownDescription: "The path of quota.",
									Computed:            true,
									Optional:            true,
								},
								"usage": schema.SingleNestedAttribute{
									Description:         "Usage",
									MarkdownDescription: "Usage",
									Computed:            true,
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"applogical": schema.Int64Attribute{
											Description:         "Bytes used by governed data apparent to application",
											MarkdownDescription: "Bytes used by governed data apparent to application",
											Computed:            true,
											Optional:            true,
										},
										"applogical_ready": schema.BoolAttribute{
											Description:         "True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"fslogical": schema.Int64Attribute{
											Description:         "Bytes used by governed data apparent to filesystem.",
											MarkdownDescription: "Bytes used by governed data apparent to filesystem.",
											Computed:            true,
											Optional:            true,
										},
										"fslogical_ready": schema.BoolAttribute{
											Description:         "True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"fsphysical": schema.Int64Attribute{
											Description:         "Bytes used by governed data apparent to filesystem.",
											MarkdownDescription: "Bytes used by governed data apparent to filesystem.",
											Computed:            true,
											Optional:            true,
										},
										"fsphysical_ready": schema.BoolAttribute{
											Description:         "True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"inodes": schema.Int64Attribute{
											Description:         "Number of inodes (filesystem entities) used by governed data.",
											MarkdownDescription: "Number of inodes (filesystem entities) used by governed data.",
											Computed:            true,
											Optional:            true,
										},
										"inodes_ready": schema.BoolAttribute{
											Description:         "True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"physical": schema.Int64Attribute{
											Description:         "Bytes used for governed data and filesystem overhead.",
											MarkdownDescription: "Bytes used for governed data and filesystem overhead.",
											Computed:            true,
											Optional:            true,
										},
										"physical_data": schema.Int64Attribute{
											Description:         "Number of physical blocks for file data",
											MarkdownDescription: "Number of physical blocks for file data",
											Computed:            true,
											Optional:            true,
										},
										"physical_data_ready": schema.BoolAttribute{
											Description:         "True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"physical_protection": schema.Int64Attribute{
											Description:         "Number of physical blocks for file protection",
											MarkdownDescription: "Number of physical blocks for file protection",
											Computed:            true,
											Optional:            true,
										},
										"physical_protection_ready": schema.BoolAttribute{
											Description:         "True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"physical_ready": schema.BoolAttribute{
											Description:         "True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
										"shadow_refs": schema.Int64Attribute{
											Description:         "Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.",
											MarkdownDescription: "Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.",
											Computed:            true,
											Optional:            true,
										},
										"shadow_refs_ready": schema.BoolAttribute{
											Description:         "True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											MarkdownDescription: "True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
											Computed:            true,
											Optional:            true,
										},
									},
								},
							},
						},
					},
					"file_system_namespace_acl": schema.SingleNestedAttribute{
						Description:         "Filesystem acl",
						MarkdownDescription: "Filesystem acl",
						Computed:            true,
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"acl": schema.ListNestedAttribute{
								Description:         "Filesystem Access Control List",
								MarkdownDescription: "Filesystem Access Control List",
								Computed:            true,
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"accessrights": schema.ListAttribute{
											Description:         "Access rights",
											MarkdownDescription: "Access rights",
											Computed:            true,
											Optional:            true,
											ElementType:         types.StringType,
										},
										"accesstype": schema.StringAttribute{
											Description:         "Access type",
											MarkdownDescription: "Access type",
											Computed:            true,
											Optional:            true,
										},
										"inherit_flags": schema.ListAttribute{
											Description:         "Inherit flags",
											MarkdownDescription: "Inherit flags",
											Computed:            true,
											Optional:            true,
											ElementType:         types.StringType,
										},
										"op": schema.StringAttribute{
											Description:         "Op",
											MarkdownDescription: "Op",
											Computed:            true,
											Optional:            true,
										},
										"trustee": schema.SingleNestedAttribute{
											Description:         "Trustee",
											MarkdownDescription: "Trustee",
											Computed:            true,
											Optional:            true,
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description:         "Trustee identifier",
													MarkdownDescription: "Trustee identifier",
													Computed:            true,
													Optional:            true,
												},
												"name": schema.StringAttribute{
													Description:         "Trustee name",
													MarkdownDescription: "Trustee name",
													Computed:            true,
													Optional:            true,
												},
												"type": schema.StringAttribute{
													Description:         "Trustee type",
													MarkdownDescription: "Trustee type",
													Computed:            true,
													Optional:            true,
												},
											},
										},
									},
								},
							},
							"action": schema.StringAttribute{
								Description:         "Acl action",
								MarkdownDescription: "Acl action",
								Computed:            true,
								Optional:            true,
							},
							"authoritative": schema.StringAttribute{
								Description:         "Acl authoritative",
								MarkdownDescription: "Acl authoritative",
								Computed:            true,
								Optional:            true,
							},
							"group": schema.SingleNestedAttribute{
								Description:         "ACL group",
								MarkdownDescription: "ACL group",
								Computed:            true,
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Group identifier",
										MarkdownDescription: "Group identifier",
										Computed:            true,
										Optional:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Group name",
										MarkdownDescription: "Group name",
										Computed:            true,
										Optional:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Group type",
										MarkdownDescription: "Group type",
										Computed:            true,
										Optional:            true,
									},
								},
							},
							"mode": schema.StringAttribute{
								Description:         "Acl mode",
								MarkdownDescription: "Acl mode",
								Computed:            true,
								Optional:            true,
							},
							"owner": schema.SingleNestedAttribute{
								Description:         "ACL owner",
								MarkdownDescription: "ACL owner",
								Computed:            true,
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Owner identifier",
										MarkdownDescription: "Owner identifier",
										Computed:            true,
										Optional:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Owner name",
										MarkdownDescription: "Owner name",
										Computed:            true,
										Optional:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Owner type",
										MarkdownDescription: "Owner type",
										Computed:            true,
										Optional:            true,
									},
								},
							},
						},
					},
					"file_system_snapshots": schema.ListNestedAttribute{
						Description:         "Filesystem snapshots",
						MarkdownDescription: "Filesystem shots",
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"alias": schema.StringAttribute{
									Description:         "The name of the alias, none for real snapshots.",
									MarkdownDescription: "The name of the alias, none for real snapshots.",
									Computed:            true,
									Optional:            true,
								},
								"created": schema.Int64Attribute{
									Description:         "The Unix Epoch time the snapshot was created.",
									MarkdownDescription: "The Unix Epoch time the snapshot was created.",
									Computed:            true,
									Optional:            true,
								},
								"expires": schema.Int64Attribute{
									Description:         "The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.",
									MarkdownDescription: "The Unix Epoch time the snapshot will expire and be eligible for automatic deletion.",
									Computed:            true,
									Optional:            true,
								},
								"has_locks": schema.BoolAttribute{
									Description:         "True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of locks.",
									MarkdownDescription: "True if the snapshot has one or more locks present see, see the locks subresource of a snapshot for a list of locks.",
									Computed:            true,
									Optional:            true,
								},
								"id": schema.Int64Attribute{
									Description:         "The system ID given to the snapshot. This is useful for tracking the status of delete pending snapshots.",
									MarkdownDescription: "The system ID given to the snapshot. This is useful for tracking the status of delete pending snapshots.",
									Computed:            true,
									Optional:            true,
								},
								"name": schema.StringAttribute{
									Description:         "The user or system supplied snapshot name. This will be null for snapshots pending delete.",
									MarkdownDescription: "The user or system supplied snapshot name. This will be null for snapshots pending delete.",
									Computed:            true,
									Optional:            true,
								},
								"path": schema.StringAttribute{
									Description:         "The /ifs path snapshotted.",
									MarkdownDescription: "The /ifs path snapshotted.",
									Computed:            true,
									Optional:            true,
								},
								"pct_filesystem": schema.NumberAttribute{
									Description:         "Percentage of /ifs used for storing this snapshot.",
									MarkdownDescription: "Percentage of /ifs used for storing this snapshot.",
									Computed:            true,
									Optional:            true,
								},
								"pct_reserve": schema.NumberAttribute{
									Description:         "Percentage of configured snapshot reserved used for storing this snapshot.",
									MarkdownDescription: "Percentage of configured snapshot reserved used for storing this snapshot.",
									Computed:            true,
									Optional:            true,
								},
								"schedule": schema.StringAttribute{
									Description:         "The name of the schedule used to create this snapshot, if applicable.",
									MarkdownDescription: "The name of the schedule used to create this snapshot, if applicable.",
									Computed:            true,
									Optional:            true,
								},
								"shadow_bytes": schema.Int64Attribute{
									Description:         "The amount of shadow bytes referred to by this snapshot.",
									MarkdownDescription: "The amount of shadow bytes referred to by this snapshot.",
									Computed:            true,
									Optional:            true,
								},
								"size": schema.Int64Attribute{
									Description:         "The amount of storage in bytes used to store this snapshot.",
									MarkdownDescription: "The amount of storage in bytes used to store this snapshot.",
									Computed:            true,
									Optional:            true,
								},
								"state": schema.StringAttribute{
									Description:         "Snapshot state.",
									MarkdownDescription: "Snapshot state.",
									Computed:            true,
									Optional:            true,
								},
								"target_id": schema.Int64Attribute{
									Description:         "The ID of the snapshot pointed to if this is an alias. 18446744073709551615 (max uint64) is returned for an alias to the live filesystem.",
									MarkdownDescription: "The ID of the snapshot pointed to if this is an alias. 18446744073709551615 (max uint64) is returned for an alias to the live filesystem.",
									Computed:            true,
									Optional:            true,
								},
								"target_name": schema.StringAttribute{
									Description:         "The name of the snapshot pointed to if this is an alias.",
									MarkdownDescription: "The name of the snapshot pointed to if this is an alias.",
									Computed:            true,
									Optional:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *FileSystemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *FileSystemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.FileSystemDataSourceModel

	// Read Terraform configuration data into the model.
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	// Remove the "/" infront of the beginning of the directory path for the API calls
	usablePath := data.DirectoryPath.ValueString()
	usablePath = usablePath[1:]
	meta, err := helper.GetDirectoryMetadata(ctx, d.client, usablePath)

	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the metadata for the filesystem",
			message,
		)
		return
	}

	acl, err := helper.GetDirectoryACL(ctx, d.client, usablePath)
	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the acl for the filesystem",
			message,
		)
		return
	}

	quota, err := helper.GetDirectoryQuota(ctx, d.client, usablePath)
	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the quota for the filesystem",
			message,
		)
		return
	}

	snapshots, err := helper.GetDirectorySnapshots(ctx, d.client)
	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the snapshots for the filesystem",
			message,
		)
		return
	}

	filteredSnapshots := helper.FilterPowerScaleSnapshots(snapshots, usablePath)

	err = helper.BuildFilesystemDatasource(ctx, &data, filteredSnapshots, quota, acl, meta)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error readingsnapshots for the filesystem",
			err.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue("FileSystem-id")

	tflog.Trace(ctx, "read a filesystem data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
