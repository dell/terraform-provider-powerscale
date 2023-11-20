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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/sync/errgroup"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

// NewUserDataSource creates a new user data source.
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema describes the data source arguments.
func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing Users from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale User allows you to authenticate through a local authentication provider. Remote users are restricted to read-only operations.",
		Description:         "This datasource is used to query the existing Users from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale User allows you to authenticate through a local authentication provider. Remote users are restricted to read-only operations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the user instance.",
				Description:         "Unique identifier of the user instance.",
			},
			"users": schema.ListNestedAttribute{
				MarkdownDescription: "List of users.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Specifies a user name.",
							MarkdownDescription: "Specifies a user name.",
							Optional:            true,
						},
						"uid": schema.StringAttribute{
							Description:         "Specifies a user identifier.",
							MarkdownDescription: "Specifies a user identifier.",
							Optional:            true,
						},
						"dn": schema.StringAttribute{
							Description:         "Specifies a principal name for the user.",
							MarkdownDescription: "Specifies a principal name for the user.",
							Optional:            true,
						},
						"dns_domain": schema.StringAttribute{
							Description:         "Specifies the DNS domain.",
							MarkdownDescription: "Specifies the DNS domain.",
							Optional:            true,
						},
						"domain": schema.StringAttribute{
							Description:         "Specifies the domain that the object is part of.",
							MarkdownDescription: "Specifies the domain that the object is part of.",
							Optional:            true,
						},
						"email": schema.StringAttribute{
							Description:         "Specifies an email address.",
							MarkdownDescription: "Specifies an email address.",
							Optional:            true,
						},
						"gecos": schema.StringAttribute{
							Description:         "Specifies the GECOS value, which is usually the full name.",
							MarkdownDescription: "Specifies the GECOS value, which is usually the full name.",
							Optional:            true,
						},
						"gid": schema.StringAttribute{
							Description:         "Specifies a group identifier.",
							MarkdownDescription: "Specifies a group identifier.",
							Optional:            true,
						},
						"home_directory": schema.StringAttribute{
							Description:         "Specifies a home directory for the user.",
							MarkdownDescription: "Specifies a home directory for the user.",
							Optional:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Specifies the user ID.",
							MarkdownDescription: "Specifies the user ID.",
							Optional:            true,
						},
						"primary_group_sid": schema.StringAttribute{
							Description:         "Specifies the persona of the primary group.",
							MarkdownDescription: "Specifies the persona of the primary group.",
							Optional:            true,
						},
						"provider": schema.StringAttribute{
							Description:         "Specifies the authentication provider that the object belongs to.",
							MarkdownDescription: "Specifies the authentication provider that the object belongs to.",
							Optional:            true,
						},
						"sam_account_name": schema.StringAttribute{
							Description:         "Specifies a user name.",
							MarkdownDescription: "Specifies a user name.",
							Optional:            true,
						},
						"shell": schema.StringAttribute{
							Description:         "Specifies a path to the shell for the user.",
							MarkdownDescription: "Specifies a path to the shell for the user.",
							Optional:            true,
						},
						"sid": schema.StringAttribute{
							Description:         "Specifies a security identifier.",
							MarkdownDescription: "Specifies a security identifier.",
							Optional:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Specifies the object type.",
							MarkdownDescription: "Specifies the object type.",
							Optional:            true,
						},
						"upn": schema.StringAttribute{
							Description:         "Specifies a principal name for the user.",
							MarkdownDescription: "Specifies a principal name for the user.",
							Optional:            true,
						},
						"enabled": schema.BoolAttribute{
							Description:         "If true, the authenticated user is enabled.",
							MarkdownDescription: "If true, the authenticated user is enabled.",
							Optional:            true,
						},
						"expired": schema.BoolAttribute{
							Description:         "If true, the authenticated user has expired.",
							MarkdownDescription: "If true, the authenticated user has expired.",
							Optional:            true,
						},
						"generated_gid": schema.BoolAttribute{
							Description:         "If true, the GID was generated.",
							MarkdownDescription: "If true, the GID was generated.",
							Optional:            true,
						},
						"generated_uid": schema.BoolAttribute{
							Description:         "If true, the UID was generated.",
							MarkdownDescription: "If true, the UID was generated.",
							Optional:            true,
						},
						"generated_upn": schema.BoolAttribute{
							Description:         "If true, the UPN was generated.",
							MarkdownDescription: "If true, the UPN was generated.",
							Optional:            true,
						},
						"locked": schema.BoolAttribute{
							Description:         "If true, indicates that the account is locked.",
							MarkdownDescription: "If true, indicates that the account is locked.",
							Optional:            true,
						},
						"password_expired": schema.BoolAttribute{
							Description:         "If true, the password has expired.",
							MarkdownDescription: "If true, the password has expired.",
							Optional:            true,
						},
						"password_expires": schema.BoolAttribute{
							Description:         "If true, the password is allowed to expire.",
							MarkdownDescription: "If true, the password is allowed to expire.",
							Optional:            true,
						},
						"prompt_password_change": schema.BoolAttribute{
							Description:         "If true, Prompts the user to change their password at the next login.",
							MarkdownDescription: "If true, Prompts the user to change their password at the next login.",
							Optional:            true,
						},
						"user_can_change_password": schema.BoolAttribute{
							Description:         "Specifies whether the password for the user can be changed.",
							MarkdownDescription: "Specifies whether the password for the user can be changed.",
							Optional:            true,
						},
						"expiry": schema.Int64Attribute{
							Description:         "Specifies the Unix Epoch time at which the authenticated user will expire.",
							MarkdownDescription: "Specifies the Unix Epoch time at which the authenticated user will expire.",
							Optional:            true,
						},
						"max_password_age": schema.Int64Attribute{
							Description:         "Specifies the maximum time in seconds allowed before the password expires.",
							MarkdownDescription: "Specifies the maximum time in seconds allowed before the password expires.",
							Optional:            true,
						},
						"password_expiry": schema.Int64Attribute{
							Description:         "Specifies the time in Unix Epoch seconds that the password will expire.",
							MarkdownDescription: "Specifies the time in Unix Epoch seconds that the password will expire.",
							Optional:            true,
						},
						"password_last_set": schema.Int64Attribute{
							Description:         "Specifies the last time the password was set.",
							MarkdownDescription: "Specifies the last time the password was set.",
							Optional:            true,
						},
						"roles": schema.ListAttribute{
							Description:         "List of roles.",
							MarkdownDescription: "List of roles.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.ListNestedAttribute{
						Description:         "List of user identity.",
						MarkdownDescription: "List of user identity.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         "Specifies a user name.",
									MarkdownDescription: "Specifies a user name.",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"uid": schema.Int64Attribute{
									Description:         "Specifies a numeric user identifier.",
									MarkdownDescription: "Specifies a numeric user identifier.",
									Optional:            true,
								},
							},
						},
					},
					"name_prefix": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter users by name prefix.",
						MarkdownDescription: "Filter users by name prefix.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"domain": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter users by domain.",
						MarkdownDescription: "Filter users by domain.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"zone": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter users by zone.",
						MarkdownDescription: "Filter users by zone.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"provider": schema.StringAttribute{
						Optional:            true,
						Description:         "Filter users by provider.",
						MarkdownDescription: "Filter users by provider.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"cached": schema.BoolAttribute{
						Optional:            true,
						Description:         "If true, only return cached objects.",
						MarkdownDescription: "If true, only return cached objects.",
					},
					"member_of": schema.BoolAttribute{
						Optional:            true,
						Description:         "Enumerate all users that a group is a member of.",
						MarkdownDescription: "Enumerate all users that a group is a member of.",
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading User data source ")

	var state models.UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// start goroutine to cache all roles
	var eg errgroup.Group
	var roles []powerscale.V1AuthRoleExtended
	var roleErr error
	var zoneID string
	if state.Filter != nil && !state.Filter.Zone.IsNull() {
		zoneID = state.Filter.Zone.ValueString()
	}
	eg.Go(func() error {
		roles, roleErr = helper.GetAllRolesWithZone(ctx, d.client, zoneID)
		return roleErr
	})

	users, err := helper.GetUsersWithFilter(ctx, d.client, state.Filter)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Users.", err.Error())
		return
	}

	if err := eg.Wait(); err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Roles", err.Error())
	}

	// parse user response to state user model
	helper.UpdateUserDataSourceState(&state, users, roles)

	// filter users by names
	if state.Filter != nil && len(state.Filter.Names) > 0 {
		var validUsers []string
		var filteredUsers []models.UserModel

		for _, user := range state.Users {
			for _, name := range state.Filter.Names {
				if (!name.Name.IsNull() && user.Name.Equal(name.Name)) ||
					(!name.UID.IsNull() && fmt.Sprintf("UID:%d", name.UID.ValueInt64()) == user.UID.ValueString()) {
					filteredUsers = append(filteredUsers, user)
					validUsers = append(validUsers, fmt.Sprintf("Name: %s, UID: %s", user.Name, user.UID))
					break
				}
			}
		}

		state.Users = filteredUsers

		if len(state.Users) != len(state.Filter.Names) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered user names is not a valid powerscale user.",
				fmt.Sprintf("Valid users: [%v], filtered list: [%v]", strings.Join(validUsers, " ; "), state.Filter.Names),
			)
		}
	}

	state.ID = types.StringValue("user_datasource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read User data source ")
}
