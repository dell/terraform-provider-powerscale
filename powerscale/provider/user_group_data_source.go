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
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &UserGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &UserGroupDataSource{}
)

// NewUserGroupDataSource creates a new user data source.
func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

// UserGroupDataSource defines the data source implementation.
type UserGroupDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *UserGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

// Schema describes the data source arguments.
func (d *UserGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing User Groups from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale User Group allows you to do operations on a set of users, groups and well-knowns.",
		Description:         "This datasource is used to query the existing User Groups from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale User Group allows you to do operations on a set of users, groups and well-knowns.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the user group instance.",
				Description:         "Unique identifier of the user group instance.",
			},
			"user_groups": schema.ListNestedAttribute{
				MarkdownDescription: "List of user groups.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Specifies a user group name.",
							MarkdownDescription: "Specifies a user group name.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Specifies the user group ID.",
							MarkdownDescription: "Specifies the user group ID.",
							Computed:            true,
						},
						"dn": schema.StringAttribute{
							Description:         "Specifies the distinguished name for the user group.",
							MarkdownDescription: "Specifies the distinguished name for the user group.",
							Computed:            true,
						},
						"dns_domain": schema.StringAttribute{
							Description:         "Specifies the DNS domain.",
							MarkdownDescription: "Specifies the DNS domain.",
							Computed:            true,
						},
						"domain": schema.StringAttribute{
							Description:         "Specifies the domain that the object is part of.",
							MarkdownDescription: "Specifies the domain that the object is part of.",
							Computed:            true,
						},
						"gid": schema.StringAttribute{
							Description:         "Specifies a user group identifier.",
							MarkdownDescription: "Specifies a user group identifier.",
							Computed:            true,
						},
						"provider": schema.StringAttribute{
							Description:         "Specifies the authentication provider that the object belongs to.",
							MarkdownDescription: "Specifies the authentication provider that the object belongs to.",
							Computed:            true,
						},
						"sam_account_name": schema.StringAttribute{
							Description:         "Specifies a user group name.",
							MarkdownDescription: "Specifies a user group name.",
							Computed:            true,
						},
						"sid": schema.StringAttribute{
							Description:         "Specifies a security identifier.",
							MarkdownDescription: "Specifies a security identifier.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Specifies the object type.",
							MarkdownDescription: "Specifies the object type.",
							Computed:            true,
						},
						"generated_gid": schema.BoolAttribute{
							Description:         "If true, the GID was generated.",
							MarkdownDescription: "If true, the GID was generated.",
							Computed:            true,
						},
						"roles": schema.ListAttribute{
							Description:         "List of roles.",
							MarkdownDescription: "List of roles.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"members": schema.ListNestedAttribute{
							Description:         "List of members of group. Group Member can be user or group.",
							MarkdownDescription: "List of members of group. Group Member can be user or group.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description:         "Specifies a user or group name.",
										MarkdownDescription: "Specifies a user or group name.",
										Computed:            true,
									},
									"id": schema.StringAttribute{
										Description:         "Specifies a user or group id.",
										MarkdownDescription: "Specifies a user or group id.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the object type.",
										MarkdownDescription: "Specifies the object type.",
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
					"names": schema.ListNestedAttribute{
						Description:         "List of user group identity.",
						MarkdownDescription: "List of user group identity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         "Specifies a user group name.",
									MarkdownDescription: "Specifies a user group name.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"gid": schema.Int64Attribute{
									Description:         "Specifies a numeric user group identifier.",
									MarkdownDescription: "Specifies a numeric user group identifier.",
									Optional:            true,
								},
							},
						},
					},
					"name_prefix": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter user groups by name prefix.",
						MarkdownDescription: "Filter user groups by name prefix.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"domain": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter user groups by domain.",
						MarkdownDescription: "Filter user groups by domain.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"zone": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter user groups by zone.",
						MarkdownDescription: "Filter user groups by zone.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"provider": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter user groups by provider.",
						MarkdownDescription: "Filter user groups by provider.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"cached": schema.BoolAttribute{
						Optional:            true,
						Description:         "If true, only return cached objects.",
						MarkdownDescription: "If true, only return cached objects.",
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *UserGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading User Group data source ")

	var state models.UserGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var roles []powerscale.V1AuthRoleExtended
	var roleErr error
	var zoneID string
	if state.Filter != nil && !state.Filter.Zone.IsNull() {
		zoneID = state.Filter.Zone.ValueString()
	}
	roles, roleErr = helper.GetAllRolesWithZone(ctx, d.client, zoneID)
	if roleErr != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Roles.", roleErr.Error())
		return
	}

	groupsResponse, err := helper.GetUserGroupsWithFilter(ctx, d.client, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale User Groups.", err.Error())
		return
	}

	var groups []models.UserGroupModel
	for _, group := range groupsResponse {
		model := models.UserGroupModel{}
		members, err := helper.GetAllGroupMembersWithZone(ctx, d.client, group.Name, zoneID)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error getting the list of PowerScale Group Members of %s", group.Name), err.Error())
		}
		helper.UpdateUserGroupDataSourceState(&model, group, members, roles)
		groups = append(groups, model)
	}

	state.UserGroups = groups

	state.ID = types.StringValue("user_group_datasource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read User Group data source ")
}
