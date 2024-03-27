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
var _ datasource.DataSource = &RolePrivilegeDataSource{}

// NewRolePrivilegeDataSource creates a new data source.
func NewRolePrivilegeDataSource() datasource.DataSource {
	return &RolePrivilegeDataSource{}
}

// RolePrivilegeDataSource defines the data source implementation.
type RolePrivilegeDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *RolePrivilegeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roleprivilege"
}

// Schema describes the data source arguments.
func (d *RolePrivilegeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing Role Privileges from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can designate certain privileges as no permission, read, execute, or write when adding the privilege to a role.",
		Description:         "This datasource is used to query the existing Role Privileges from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can designate certain privileges as no permission, read, execute, or write when adding the privilege to a role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the Role Privilege instance.",
				MarkdownDescription: "Unique identifier of the Role Privilege instance.",
				Computed:            true,
			},
			"role_privileges_details": schema.ListNestedAttribute{
				Description:         "List of Role Privileges.",
				MarkdownDescription: "List of Role Privileges.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uri": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the associated uri for the privilege.",
							MarkdownDescription: "Specifies the associated uri for the privilege.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the general categorization of the privilege.",
							MarkdownDescription: "Specifies the general categorization of the privilege.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies a short description of the privilege.",
							MarkdownDescription: "Specifies a short description of the privilege.",
						},
						"permission": schema.StringAttribute{
							Computed:            true,
							Description:         "Permissions the privilege has r=read , x=read-execute, w=read-execute-write.",
							MarkdownDescription: "Permissions the privilege has r=read , x=read-execute, w=read-execute-write.",
						},
						"privilegelevel": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the level of the privilege.",
							MarkdownDescription: "Specifies the level of the privilege.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the ID of the privilege.",
							MarkdownDescription: "Specifies the ID of the privilege.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the name of the privilege.",
							MarkdownDescription: "Specifies the name of the privilege.",
						},
						"parent_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the parent ID of the privilege.",
							MarkdownDescription: "Specifies the parent ID of the privilege.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Description:         "Filter Role Privileges by names.",
						MarkdownDescription: "Filter Role Privileges by names.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *RolePrivilegeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *RolePrivilegeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading role privilege data source")

	var state models.RolePrivilegeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rolePrivilegeList, err := helper.GetRolePrivileges(ctx, d.client)

	if err != nil {
		errStr := constants.ReadRolePrivilegeErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of role privileges",
			message,
		)
		return
	}

	var rolePrivileges []models.RolePrivilegeDetailModel
	for _, rolePrivilegeItem := range rolePrivilegeList.Privileges {
		val := rolePrivilegeItem
		rolePrivilege, err := helper.RolePrivilegeDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadRolePrivilegeErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error mapping the list of role privileges",
				message,
			)
			return
		}
		rolePrivileges = append(rolePrivileges, rolePrivilege)
	}

	state.RolePrivileges = rolePrivileges

	// filter role privileges by names
	if state.RolePrivilegeFilter != nil && len(state.RolePrivilegeFilter.Names) > 0 {
		var filteredRolePrivileges []models.RolePrivilegeDetailModel

		for _, rolePrivilege := range state.RolePrivileges {
			for _, name := range state.RolePrivilegeFilter.Names {
				if !name.IsNull() && strings.Contains(strings.ToLower(rolePrivilege.Name.ValueString()), strings.ToLower(name.ValueString())) {
					filteredRolePrivileges = append(filteredRolePrivileges, rolePrivilege)
					continue
				}
			}
		}

		state.RolePrivileges = filteredRolePrivileges

		if len(state.RolePrivileges) == 0 {
			resp.Diagnostics.AddError(
				"Error getting the list of role privileges",
				"No relevant role privileges are found",
			)
		}
	}

	// save into the Terraform state.
	state.ID = types.StringValue("role_privilege_datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading role privilege data source ")
}
