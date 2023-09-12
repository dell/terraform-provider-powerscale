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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/sync/errgroup"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

// NewUserResource creates a new resource.
func NewUserResource() resource.Resource {
	return &UserResource{}
}

// UserResource defines the resource implementation.
type UserResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema describes the resource arguments.
func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Users in PowerScale cluster. PowerScale User allows you to authenticate through a local authentication provider. Remote users are restricted to read-only operations.",
		Description:         "Resource for managing Users in PowerScale cluster. PowerScale User allows you to authenticate through a local authentication provider. Remote users are restricted to read-only operations.",

		Attributes: map[string]schema.Attribute{
			"query_force": schema.BoolAttribute{
				Description:         "If true, skip validation checks when creating user. Need to be true, when changing user UID.",
				MarkdownDescription: "If true, skip validation checks when creating user. Need to be true, when changing user UID.",
				Optional:            true,
			},
			"query_zone": schema.StringAttribute{
				Description:         "Specifies the zone that the object belongs to.",
				MarkdownDescription: "Specifies the zone that the object belongs to.",
				Optional:            true,
			},
			"query_provider": schema.StringAttribute{
				Description:         "Specifies the provider type.",
				MarkdownDescription: "Specifies the provider type.",
				Optional:            true,
			},
			"provider_name": schema.StringAttribute{
				Description:         "Specifies the authentication provider that the object belongs to.",
				MarkdownDescription: "Specifies the authentication provider that the object belongs to.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Specifies the user ID.",
				MarkdownDescription: "Specifies the user ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Specifies a user name.",
				MarkdownDescription: "Specifies a user name.",
				Required:            true,
			},
			"uid": schema.Int64Attribute{
				Description:         "Specifies a numeric user identifier. (Update Supported)",
				MarkdownDescription: "Specifies a numeric user identifier. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"email": schema.StringAttribute{
				Description:         "Specifies an email address. (Update Supported)",
				MarkdownDescription: "Specifies an email address. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "If true, the authenticated user is enabled. (Update Supported)",
				MarkdownDescription: "If true, the authenticated user is enabled. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"expiry": schema.Int64Attribute{
				Description:         "Specifies the Unix Epoch time at which the authenticated user will expire. (Update Supported)",
				MarkdownDescription: "Specifies the Unix Epoch time at which the authenticated user will expire. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"gecos": schema.StringAttribute{
				Description:         "Specifies the GECOS value, which is usually the full name. (Update Supported)",
				MarkdownDescription: "Specifies the GECOS value, which is usually the full name. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"home_directory": schema.StringAttribute{
				Description:         "Specifies a home directory for the user. (Update Supported)",
				MarkdownDescription: "Specifies a home directory for the user. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				Description:         "Sets or Changes the password for the user. (Update Supported)",
				MarkdownDescription: "Sets or Changes the password for the user. (Update Supported)",
				Optional:            true,
				Sensitive:           true,
			},
			"password_expires": schema.BoolAttribute{
				Description:         "If true, the password is allowed to expire. (Update Supported)",
				MarkdownDescription: "If true, the password is allowed to expire. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"prompt_password_change": schema.BoolAttribute{
				Description:         "If true, Prompts the user to change their password at the next login. (Update Supported)",
				MarkdownDescription: "If true, Prompts the user to change their password at the next login. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"shell": schema.StringAttribute{
				Description:         "Specifies a path to the shell for the user. (Update Supported)",
				MarkdownDescription: "Specifies a path to the shell for the user. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"unlock": schema.BoolAttribute{
				Description:         "If true, the user account should be unlocked. (Update Supported)",
				MarkdownDescription: "If true, the user account should be unlocked. (Update Supported)",
				Optional:            true,
			},
			"roles": schema.ListAttribute{
				Description:         "List of roles, the user is assigned. (Update Supported)",
				MarkdownDescription: "List of roles, the user is assigned. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"primary_group": schema.StringAttribute{
				Description:         "Specifies the name of the primary group. (Update Supported)",
				MarkdownDescription: "Specifies the name of the primary group. (Update Supported)",
				Computed:            true,
				Optional:            true,
			},
			"dn": schema.StringAttribute{
				Description:         "Specifies a principal name for the user.",
				MarkdownDescription: "Specifies a principal name for the user.",
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
				Optional:            true,
				Computed:            true,
			},
			"expired": schema.BoolAttribute{
				Description:         "If true, the authenticated user has expired.",
				MarkdownDescription: "If true, the authenticated user has expired.",
				Computed:            true,
			},
			"gid": schema.ObjectAttribute{
				Description:         "Specifies a group identifier.",
				MarkdownDescription: "Specifies a group identifier.",
				Computed:            true,
				AttributeTypes: map[string]attr.Type{
					"name": types.StringType,
					"id":   types.StringType,
					"type": types.StringType,
				},
			},
			"primary_group_sid": schema.ObjectAttribute{
				Description:         "Specifies the persona of the primary group.",
				MarkdownDescription: "Specifies the persona of the primary group.",
				Computed:            true,
				AttributeTypes: map[string]attr.Type{
					"name": types.StringType,
					"id":   types.StringType,
					"type": types.StringType,
				},
			},
			"sam_account_name": schema.StringAttribute{
				Description:         "Specifies a user name.",
				MarkdownDescription: "Specifies a user name.",
				Computed:            true,
			},
			"sid": schema.StringAttribute{
				Description:         "Specifies a security identifier. (Update Supported)",
				MarkdownDescription: "Specifies a security identifier. (Update Supported)",
				Computed:            true,
				Optional:            true,
			},
			"type": schema.StringAttribute{
				Description:         "Specifies the object type.",
				MarkdownDescription: "Specifies the object type.",
				Computed:            true,
			},
			"upn": schema.StringAttribute{
				Description:         "Specifies a principal name for the user.",
				MarkdownDescription: "Specifies a principal name for the user.",
				Computed:            true,
			},
			"generated_gid": schema.BoolAttribute{
				Description:         "If true, the GID was generated.",
				MarkdownDescription: "If true, the GID was generated.",
				Computed:            true,
			},
			"generated_uid": schema.BoolAttribute{
				Description:         "If true, the UID was generated.",
				MarkdownDescription: "If true, the UID was generated.",
				Computed:            true,
			},
			"generated_upn": schema.BoolAttribute{
				Description:         "If true, the UPN was generated.",
				MarkdownDescription: "If true, the UPN was generated.",
				Computed:            true,
			},
			"locked": schema.BoolAttribute{
				Description:         "If true, indicates that the account is locked.",
				MarkdownDescription: "If true, indicates that the account is locked.",
				Computed:            true,
			},
			"password_expired": schema.BoolAttribute{
				Description:         "If true, the password has expired.",
				MarkdownDescription: "If true, the password has expired.",
				Computed:            true,
			},
			"user_can_change_password": schema.BoolAttribute{
				Description:         "Specifies whether the password for the user can be changed.",
				MarkdownDescription: "Specifies whether the password for the user can be changed.",
				Computed:            true,
			},
			"max_password_age": schema.Int64Attribute{
				Description:         "Specifies the maximum time in seconds allowed before the password expires.",
				MarkdownDescription: "Specifies the maximum time in seconds allowed before the password expires.",
				Computed:            true,
			},
			"password_expiry": schema.Int64Attribute{
				Description:         "Specifies the time in Unix Epoch seconds that the password will expire.",
				MarkdownDescription: "Specifies the time in Unix Epoch seconds that the password will expire.",
				Computed:            true,
			},
			"password_last_set": schema.Int64Attribute{
				Description:         "Specifies the last time the password was set.",
				MarkdownDescription: "Specifies the last time the password was set.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating User resource...")
	var plan models.UserResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roleList []string
	if !plan.Roles.IsNull() && !plan.Roles.IsUnknown() {
		diags := plan.Roles.ElementsAs(ctx, &roleList, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	userName := plan.Name.ValueString()
	err := helper.CreateUser(ctx, r.client, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error creating the User - %s", userName),
			err.Error(),
		)
		return
	}

	if diags := helper.UpdateUserRoles(ctx, r.client, &models.UserResourceModel{}, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	result, err := helper.GetUserWithZone(ctx, r.client, userName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the User - %s", userName),
			err.Error(),
		)
		return
	}

	// parse user response to state user model
	helper.UpdateUserResourceState(&plan, result.Users[0], nil)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create User resource")
}

// Read reads the resource state.
func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading User resource")
	var plan models.UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userName := plan.Name.ValueString()
	result, err := helper.GetUserWithZone(ctx, r.client, userName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the User - %s", userName),
			err.Error(),
		)
		return
	}

	// parse user response to state user model
	helper.UpdateUserResourceState(&plan, result.Users[0], nil)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Read User resource")
}

// Update updates the resource state.
func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating User resource...")
	// Read Terraform plan into the model
	var plan models.UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userName := state.Name.ValueString()
	if err := helper.UpdateUser(ctx, r.client, &state, &plan); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error updating the User - %s", userName),
			err.Error(),
		)
		return
	}

	if diags := helper.UpdateUserRoles(ctx, r.client, &state, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	result, err := helper.GetUserWithZone(ctx, r.client, userName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the User - %s", userName),
			err.Error(),
		)
		return
	}

	// parse user response to state user model
	helper.UpdateUserResourceState(&plan, result.Users[0], nil)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting User resource")
	var state models.UserResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// remove user from roles
	var roleList []string
	state.Roles.ElementsAs(ctx, &roleList, false)
	for _, role := range roleList {
		_ = helper.RemoveUserRoleWithZone(ctx, r.client, role, state.UID.ValueInt64(), state.QueryZone.ValueString())
	}

	if err := helper.DeleteUserWithZone(ctx, r.client, state.Name.ValueString(), state.QueryZone.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting the User - %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete User resource")
}

// ImportState imports the resource state.
func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing User resource")
	var state models.UserResourceModel

	var zoneID string
	userName := req.ID
	//requestID format is zoneID:userName
	if strings.Contains(userName, ":") {
		params := strings.Split(userName, ":")
		userName = strings.Trim(params[1], " ")
		zoneID = strings.Trim(params[0], " ")
	}

	// start goroutine to cache all roles
	var eg errgroup.Group
	var roles []powerscale.V1AuthRoleExtended
	var roleErr error
	eg.Go(func() error {
		roles, roleErr = helper.GetAllRolesWithZone(ctx, r.client, zoneID)
		return roleErr
	})

	result, err := helper.GetUserWithZone(ctx, r.client, userName, zoneID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the User - %s", userName),
			err.Error(),
		)
		return
	}

	if err := eg.Wait(); err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Roles", err.Error())
	}

	// parse user response to state user model
	helper.UpdateUserResourceState(&state, result.Users[0], roles)
	if zoneID != "" {
		state.QueryZone = types.StringValue(zoneID)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import User resource")
}
