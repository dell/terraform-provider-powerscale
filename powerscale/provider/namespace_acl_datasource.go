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
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NamespaceACLDataSource{}

// NewNamespaceACLDataSource creates a new data source.
func NewNamespaceACLDataSource() datasource.DataSource {
	return &NamespaceACLDataSource{}
}

// NamespaceACLDataSource defines the data source implementation.
type NamespaceACLDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NamespaceACLDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_acl"
}

// Schema describes the data source arguments.
func (d *NamespaceACLDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the Namespace ACL from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use Namespace ACL to manage the access control list for a namespace.",
		Description:         "This datasource is used to query the Namespace ACL from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use Namespace ACL to manage the access control list for a namespace.",
		Attributes: map[string]schema.Attribute{
			"authoritative": schema.StringAttribute{
				Computed:            true,
				Description:         "If the directory has access rights set, then this field is returned as acl. If the directory has POSIX permissions set, then this field is returned as mode.",
				MarkdownDescription: "If the directory has access rights set, then this field is returned as acl. If the directory has POSIX permissions set, then this field is returned as mode.",
			},
			"mode": schema.StringAttribute{
				Computed:            true,
				Description:         "Provides the POSIX mode.",
				MarkdownDescription: "Provides the POSIX mode.",
			},
			"owner": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Provides the JSON object for the group persona of the owner.",
				MarkdownDescription: "Provides the JSON object for the group persona of the owner.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the serialized form of a persona, which can be 'UID:0'",
						MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0'",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the persona name, which must be combined with a type.",
						MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
					},
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the type of persona, which must be combined with a name.",
						MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
					},
				},
			},
			"group": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Provides the JSON object for the group persona of the owner.",
				MarkdownDescription: "Provides the JSON object for the group persona of the owner.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the persona name, which must be combined with a type.",
						MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
					},
					"type": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the type of persona, which must be combined with a name.",
						MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
					},
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Specifies the serialized form of a persona, which can be 'GID:0'",
						MarkdownDescription: "Specifies the serialized form of a persona, which can be 'GID:0'",
					},
				},
			},
			"acl": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "Array effective configuration of the JSON array of access rights.",
				MarkdownDescription: "Array effective configuration of the JSON array of access rights.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"accesstype": schema.StringAttribute{
							Computed:            true,
							Description:         "Grants or denies access control permissions.",
							MarkdownDescription: "Grants or denies access control permissions.",
						},
						"op": schema.StringAttribute{
							Computed:            true,
							Description:         "Operations for updating access control permissions. Unnecessary for access right replacing scenario",
							MarkdownDescription: "Operations for updating access control permissions. Unnecessary for access right replacing scenario",
						},
						"inherit_flags": schema.ListAttribute{
							Computed:            true,
							Description:         "Grants or denies access control permissions.",
							MarkdownDescription: "Grants or denies access control permissions.",
							ElementType:         types.StringType,
						},
						"trustee": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "Provides the JSON object for the group persona of the owner.",
							MarkdownDescription: "Provides the JSON object for the group persona of the owner.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:            true,
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								},
								"id": schema.StringAttribute{
									Computed:            true,
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0' or 'GID:0'",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0' or 'GID:0'",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								},
							},
						},
						"accessrights": schema.ListAttribute{
							Computed:            true,
							Description:         "Specifies the access control permissions for a specific user or group.",
							MarkdownDescription: "Specifies the access control permissions for a specific user or group.",
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"namespace": schema.StringAttribute{
						Required:            true,
						Description:         "Indicate the namespace to set/get acl.",
						MarkdownDescription: "Indicate the namespace to set/get acl.",
					},
					"nsaccess": schema.BoolAttribute{
						Optional:            true,
						Description:         "Indicates that the operation is on the access point instead of the store path.",
						MarkdownDescription: "Indicates that the operation is on the access point instead of the store path.",
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *NamespaceACLDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NamespaceACLDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "reading namespace acl data source")
	var state models.NamespaceACLDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	namespaceACLResp, err := helper.GetNamespaceACLDatasource(ctx, d.client, state)

	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the namespace acl",
			message,
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, namespaceACLResp, &state)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading namespace acl",
			fmt.Sprintf("Could not read namespace acl struct with error: %s", message),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading namespace acl data source")
}
