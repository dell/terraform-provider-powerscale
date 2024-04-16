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
	"strings"
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
var _ datasource.DataSource = &RoleDataSource{}

// NewRoleDataSource creates a new data source.
func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

// RoleDataSource defines the data source implementation.
type RoleDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *RoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema describes the data source arguments.
func (d *RoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing roles from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can permit and limit access to administrative areas of your cluster on a per-user basis through roles.",
		Description:         "This datasource is used to query the existing roles from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can permit and limit access to administrative areas of your cluster on a per-user basis through roles.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the network pool instance.",
				MarkdownDescription: "Unique identifier of the network pool instance.",
				Computed:            true,
			},
			"roles_details": schema.ListNestedAttribute{
				Description:         "List of Roles.",
				MarkdownDescription: "List of Roles.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "Specifies the ID of the role.",
							MarkdownDescription: "Specifies the ID of the role.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Specifies the description of the role.",
							MarkdownDescription: "Specifies the description of the role.",
							Computed:            true,
						},
						"members": schema.ListNestedAttribute{
							Description:         "Specifies the users or groups that have this role.",
							MarkdownDescription: "Specifies the users or groups that have this role.",
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
						"name": schema.StringAttribute{
							Description:         "Specifies the name of the role.",
							MarkdownDescription: "Specifies the name of the role.",
							Computed:            true,
						},
						"privileges": schema.ListNestedAttribute{
							Description:         "Specifies the privileges granted by this role.",
							MarkdownDescription: "Specifies the privileges granted by this role.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the ID of the privilege.",
										MarkdownDescription: "Specifies the ID of the privilege.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the name of the privilege.",
										MarkdownDescription: "Specifies the name of the privilege.",
										Computed:            true,
									},
									"permission": schema.StringAttribute{
										Description:         "permission of the privilege, 'r' = read , 'x' = read-execute, 'w' = read-execute-write, '-' = no permission",
										MarkdownDescription: "permission of the privilege, 'r' = read , 'x' = read-execute, 'w' = read-execute-write, '-' = no permission",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Description:         "Filter roles by names.",
						MarkdownDescription: "Filter roles by names.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"zone": schema.StringAttribute{
						Description:         "Specifies which access zone to use.",
						MarkdownDescription: "Specifies which access zone to use.",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading role data source")

	var state models.RoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleList, err := helper.GetRoles(ctx, d.client, state)

	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of roles",
			message,
		)
		return
	}

	var roles []models.RoleDetailModel
	for _, roleItem := range roleList.Roles {
		val := roleItem
		role, err := helper.RoleDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadRoleErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error mapping the list of roles",
				message,
			)
			return
		}
		roles = append(roles, role)
	}

	state.Roles = roles

	// filter roles by names
	if state.RoleFilter != nil && len(state.RoleFilter.Names) > 0 {
		var validRoles []string
		var filteredRoles []models.RoleDetailModel

		for _, role := range state.Roles {
			for _, name := range state.RoleFilter.Names {
				if !name.IsNull() && role.Name.Equal(name) {
					filteredRoles = append(filteredRoles, role)
					validRoles = append(validRoles, fmt.Sprintf("Name: %s", role.Name))
					continue
				}
			}
		}

		state.Roles = filteredRoles

		if len(state.Roles) != len(state.RoleFilter.Names) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered role names is not a valid powerscale role.",
				fmt.Sprintf("Valid roles: [%v], filtered list: [%v]", strings.Join(validRoles, " ; "), state.RoleFilter.Names),
			)
		}
	}

	// save into the Terraform state.
	state.ID = types.StringValue("role_datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading role data source ")
}
