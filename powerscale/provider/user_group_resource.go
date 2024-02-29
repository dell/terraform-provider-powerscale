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
	"regexp"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/sync/errgroup"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &UserGroupResource{}
	_ resource.ResourceWithConfigure   = &UserGroupResource{}
	_ resource.ResourceWithImportState = &UserGroupResource{}
)

// NewUserGroupResource creates a new resource.
func NewUserGroupResource() resource.Resource {
	return &UserGroupResource{}
}

// UserGroupResource defines the resource implementation.
type UserGroupResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *UserGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

// Schema describes the resource arguments.
func (r *UserGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the User Group entity of PowerScale Array. We can Create, Update and Delete the User Group using this resource. We can also import an existing User Group from PowerScale array. PowerScale User Group allows you to do operations on a set of users, groups and well-knowns.",
		Description:         "This resource is used to manage the User Group entity of PowerScale Array. We can Create, Update and Delete the User Group using this resource. We can also import an existing User Group from PowerScale array. PowerScale User Group allows you to do operations on a set of users, groups and well-knowns.",
		Attributes: map[string]schema.Attribute{
			"query_force": schema.BoolAttribute{
				Description:         "If true, skip validation checks when creating user group. Need to be true, when changing group GID.",
				MarkdownDescription: "If true, skip validation checks when creating user group. Need to be true, when changing group GID.",
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
			"name": schema.StringAttribute{
				Description:         "Specifies a user group name.",
				MarkdownDescription: "Specifies a user group name.",
				Required:            true,
			},
			"gid": schema.Int64Attribute{
				Description:         "Specifies a numeric user group identifier. (Update Supported)",
				MarkdownDescription: "Specifies a numeric user group identifier. (Update Supported)",
				Optional:            true,
				Computed:            true,
			},
			"sid": schema.StringAttribute{
				Description:         "Specifies a security identifier.",
				MarkdownDescription: "Specifies a security identifier.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^SID`), "must start with 'SID:'",
					),
				},
			},
			"roles": schema.ListAttribute{
				Description:         "List of roles, the user is assigned. (Update Supported)",
				MarkdownDescription: "List of roles, the user is assigned. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"users": schema.ListAttribute{
				Description:         "Specifies list members of user within the group. (Update Supported)",
				MarkdownDescription: "Specifies list members of user within the group. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"groups": schema.ListAttribute{
				Description:         "Specifies list members of group within the group. (Update Supported)",
				MarkdownDescription: "Specifies list members of group within the group. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"well_knowns": schema.ListAttribute{
				Description:         "Specifies list members of well_known within the group. (Update Supported)",
				MarkdownDescription: "Specifies list members of well_known within the group. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"provider_name": schema.StringAttribute{
				Description:         "Specifies the authentication provider that the object belongs to.",
				MarkdownDescription: "Specifies the authentication provider that the object belongs to.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Specifies the user group ID.",
				MarkdownDescription: "Specifies the user group ID.",
				Computed:            true,
			},
			"dn": schema.StringAttribute{
				Description:         "Specifies a principal name for the user group.",
				MarkdownDescription: "Specifies a principal name for the user group.",
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
			"sam_account_name": schema.StringAttribute{
				Description:         "Specifies a user group name.",
				MarkdownDescription: "Specifies a user group name.",
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
		},
	}
}

// Configure configures the resource.
func (r *UserGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating User Group resource...")
	var plan models.UserGroupResourceModel

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

	userGroupName := plan.Name.ValueString()
	err := helper.CreateUserGroup(ctx, r.client, &plan)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error creating the User Group - %s", userGroupName), err.Error())
		return
	}

	// add user group members
	if diags := helper.UpdateUserGroupMembers(ctx, r.client, &models.UserGroupResourceModel{}, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	if diags := helper.UpdateUserGroupRoles(ctx, r.client, &models.UserGroupResourceModel{}, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	result, err := helper.GetUserGroupWithZone(ctx, r.client, userGroupName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting the User Group - %s", userGroupName), err.Error())
		return
	}

	// parse user response to state user group model
	helper.UpdateUserGroupResourceState(&plan, result.Groups[0], nil, nil)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create User Group resource")
}

// Read reads the resource state.
func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading User Group resource")
	var plan models.UserGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupName := plan.Name.ValueString()
	result, err := helper.GetUserGroupWithZone(ctx, r.client, groupName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting the User Group - %s", groupName), err.Error())
		return
	}
	roles, roleErr := helper.GetAllRolesWithZone(ctx, r.client, plan.QueryZone.ValueString())
	if roleErr != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Roles", roleErr.Error())
	}
	members, err := helper.GetAllGroupMembersWithZone(ctx, r.client, groupName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting the list of PowerScale Group Members of %s", groupName), err.Error())
	}
	// parse user response to state user group model
	helper.UpdateUserGroupResourceState(&plan, result.Groups[0], members, roles)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Read User Group resource")
}

// Update updates the resource state.
func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating User Group resource...")
	// Read Terraform plan into the model
	var plan models.UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userGroupName := state.Name.ValueString()
	if err := helper.UpdateUserGroup(ctx, r.client, &state, &plan); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error updating the User Group - %s", userGroupName), err.Error())
		return
	}
	// update user group members
	if diags := helper.UpdateUserGroupMembers(ctx, r.client, &state, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	// update user group roles
	if diags := helper.UpdateUserGroupRoles(ctx, r.client, &state, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	result, err := helper.GetUserGroupWithZone(ctx, r.client, userGroupName, plan.QueryZone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting the User Group - %s", userGroupName), err.Error())
		return
	}

	// parse user response to state user group model
	helper.UpdateUserGroupResourceState(&plan, result.Groups[0], nil, nil)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Update User Group resource")
}

// Delete deletes the resource.
func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting User Group resource")
	var state models.UserGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// remove user group from roles
	var roleList []string
	state.Roles.ElementsAs(ctx, &roleList, false)

	for _, role := range roleList {
		_ = helper.RemoveUserGroupRoleWithZone(ctx, r.client, role, state.GID.ValueInt64(), state.QueryZone.ValueString())
	}

	if err := helper.DeleteUserGroupWithZone(ctx, r.client, state.Name.ValueString(), state.QueryZone.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting the User Group - %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete User Group resource")
}

// ImportState imports the resource state.
func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing User Group resource")
	var state models.UserGroupResourceModel

	var zoneID string
	groupName := req.ID
	//requestID format is zoneID:groupName
	if strings.Contains(groupName, ":") {
		params := strings.Split(groupName, ":")
		groupName = strings.Trim(params[1], " ")
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

	result, err := helper.GetUserGroupWithZone(ctx, r.client, groupName, zoneID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting the User Group - %s", groupName), err.Error())
		return
	}

	members, err := helper.GetAllGroupMembersWithZone(ctx, r.client, groupName, zoneID)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting the list of PowerScale Group Members of %s", groupName), err.Error())
	}

	if err := eg.Wait(); err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Roles", err.Error())
	}

	// parse user response to state user group model
	helper.UpdateUserGroupResourceState(&state, result.Groups[0], members, roles)
	if zoneID != "" {
		state.QueryZone = types.StringValue(zoneID)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import User Group resource")
}
