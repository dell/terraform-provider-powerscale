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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithConfigure = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}

// NewRoleResource creates a new resource.
func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

// RoleResource defines the resource implementation.
type RoleResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema describes the resource arguments.
func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the role entity of PowerScale Array. We can Create, Update and Delete the role using this resource. We can also import an existing role from PowerScale array.",
		Description:         "This resource is used to manage the role entity of PowerScale Array. We can Create, Update and Delete the role using this resource. We can also import an existing role from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Optional:            true,
				Description:         "Specifies which access zone to use.",
				MarkdownDescription: "Specifies which access zone to use.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Specifies the name of the role.",
				MarkdownDescription: "Specifies the name of the role.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"members": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies the users or groups that have this role.",
				MarkdownDescription: "Specifies the users or groups that have this role.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
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
							Required:            true,
							Description:         "Specifies the serialized form of a persona, which can be 'UID:0'",
							MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0'",
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 261),
							},
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the persona name, which must be combined with a type.",
							MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 255),
							},
						},
					},
				},
			},
			"privileges": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies the privileges granted by this role.",
				MarkdownDescription: "Specifies the privileges granted by this role.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "permission of the privilege, 'r' = read , 'x' = read-execute, 'w' = read-execute-write, '-' = no permission",
							MarkdownDescription: "permission of the privilege, 'r' = read , 'x' = read-execute, 'w' = read-execute-write, '-' = no permission",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"r",
									"w",
									"x",
									"-",
								),
							},
						},
						"id": schema.StringAttribute{
							Required:            true,
							Description:         "Specifies the ID of the privilege.",
							MarkdownDescription: "Specifies the ID of the privilege.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 255),
							},
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Specifies the name of the privilege.",
							MarkdownDescription: "Specifies the name of the privilege.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 255),
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Specifies the description of the role.",
				MarkdownDescription: "Specifies the description of the role.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Specifies the ID of the role.",
				MarkdownDescription: "Specifies the ID of the role.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	powerscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = powerscaleClient
}

// Create allocates the resource.
func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating role")

	var plan models.RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleToCreate := powerscale.V14AuthRole{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &roleToCreate)
	if err != nil {
		errStr := constants.CreateRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Could not read role param with error: %s", message),
		)
		return
	}

	roleID, err := helper.CreateRole(ctx, r.client, roleToCreate, plan)
	if err != nil {
		errStr := constants.CreateRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating role",
			message,
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("role %s created", roleID.Id), map[string]interface{}{
		"roleResponse": roleID,
	})

	plan.ID = types.StringValue(roleID.Id)
	getRoleResponse, err := helper.GetRole(ctx, r.client, plan)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating role",
			message,
		)
		return
	}

	// update resource state according to response
	if len(getRoleResponse.Roles) <= 0 {
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Could not read created role %s", roleID),
		)
		return
	}

	createdRole := getRoleResponse.Roles[0]
	originalPlan := plan
	err = helper.CopyFields(ctx, createdRole, &plan)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Could not read role struct %s with error: %s", roleID, message),
		)
		return
	}

	orderedMemberList, err := helper.ReorderRoleMembers(originalPlan.Members, plan.Members)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Could not reorder role members for %s with error: %s", roleID, message),
		)
		return
	}
	plan.Members = orderedMemberList
	orderedPrivilegeList, err := helper.ReorderRolePrivileges(originalPlan.Privileges, plan.Privileges)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating role",
			fmt.Sprintf("Could not reorder role privileges for %s with error: %s", roleID, message),
		)
		return
	}
	plan.Privileges = orderedPrivilegeList

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create role completed")
}

// Read reads the resource state.
func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading role")

	var roleState models.RoleResourceModel
	diags := req.State.Get(ctx, &roleState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := roleState.ID
	tflog.Debug(ctx, "calling get role by ID", map[string]interface{}{
		"roleID": roleState.ID,
	})
	roleResponse, err := helper.GetRole(ctx, r.client, roleState)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading role",
			message,
		)
		return
	}

	if len(roleResponse.Roles) <= 0 {
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not read role %s from powerscale with error: role not found", roleID),
		)
		return
	}
	tflog.Debug(ctx, "updating role state", map[string]interface{}{
		"roleResponse": roleResponse,
		"roleState":    roleState,
	})

	originalState := roleState
	err = helper.CopyFields(ctx, roleResponse.Roles[0], &roleState)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not read role struct %s with error: %s", roleID, message),
		)
		return
	}

	orderedMemberList, err := helper.ReorderRoleMembers(originalState.Members, roleState.Members)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not reorder role members for %s with error: %s", roleID, message),
		)
		return
	}
	roleState.Members = orderedMemberList
	orderedPrivilegeList, err := helper.ReorderRolePrivileges(originalState.Privileges, roleState.Privileges)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not reorder role privileges for %s with error: %s", roleID, message),
		)
		return
	}
	roleState.Privileges = orderedPrivilegeList

	diags = resp.State.Set(ctx, roleState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read role completed")
}

// Update updates the resource state.
func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating role")

	var rolePlan models.RoleResourceModel
	diags := req.Plan.Get(ctx, &rolePlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var roleState models.RoleResourceModel
	diags = resp.State.Get(ctx, &roleState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update role", map[string]interface{}{
		"rolePlan":  rolePlan,
		"roleState": roleState,
	})

	roleID := roleState.ID.ValueString()
	rolePlan.ID = roleState.ID
	var roleToUpdate powerscale.V14AuthRoleExtendedExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, rolePlan, &roleToUpdate)
	if err != nil {
		errStr := constants.UpdateRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating role",
			fmt.Sprintf("Could not read role param with error: %s", message),
		)
		return
	}
	err = helper.UpdateRole(ctx, r.client, rolePlan, roleToUpdate)
	if err != nil {
		errStr := constants.UpdateRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating role",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get role by ID on powerscale client", map[string]interface{}{
		"roleID": roleID,
	})

	// Role ID and name should be consistent after the update
	rolePlan.ID = rolePlan.Name

	updatedRole, err := helper.GetRole(ctx, r.client, rolePlan)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating role",
			message,
		)
		return
	}

	if len(updatedRole.Roles) <= 0 {
		resp.Diagnostics.AddError(
			"Error updating role",
			fmt.Sprintf("Could not read updated role %s", roleID),
		)
		return
	}

	originalPlan := rolePlan
	err = helper.CopyFields(ctx, updatedRole.Roles[0], &rolePlan)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating role",
			fmt.Sprintf("Could not read role struct %s with error: %s", roleID, message),
		)
		return
	}

	orderedMemberList, err := helper.ReorderRoleMembers(originalPlan.Members, rolePlan.Members)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating role",
			fmt.Sprintf("Could not reorder role members for %s with error: %s", roleID, message),
		)
		return
	}
	rolePlan.Members = orderedMemberList
	orderedPrivilegeList, err := helper.ReorderRolePrivileges(originalPlan.Privileges, rolePlan.Privileges)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating role",
			fmt.Sprintf("Could not reorder role privileges for %s with error: %s", roleID, message),
		)
		return
	}
	rolePlan.Privileges = orderedPrivilegeList

	diags = resp.State.Set(ctx, rolePlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update role completed")
}

// Delete deletes the resource.
func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting role")

	var roleState models.RoleResourceModel
	diags := req.State.Get(ctx, &roleState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := roleState.ID.ValueString()
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete role on powerscale client", map[string]interface{}{
		"roleID": roleID,
	})
	err := helper.DeleteRole(ctx, r.client, roleState)
	if err != nil {
		errStr := constants.DeleteRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting role",
			message,
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete role completed")
}

// ImportState imports the resource state.
func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing role")
	var roleState models.RoleResourceModel

	roleID := req.ID
	var zoneID string
	if strings.Contains(req.ID, ":") {
		params := strings.Split(req.ID, ":")
		zoneID = strings.Trim(params[0], " ")
		roleID = strings.Trim(params[1], " ")
	}

	roleState.ID = types.StringValue(roleID)
	roleState.Zone = types.StringValue(zoneID)
	tflog.Debug(ctx, "calling get role by ID", map[string]interface{}{
		"roleID": roleState.ID,
	})
	roleResponse, err := helper.GetRole(ctx, r.client, roleState)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading role",
			message,
		)
		return
	}

	if len(roleResponse.Roles) <= 0 {
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not read role %s from powerscale with error: role not found", roleID),
		)
		return
	}
	tflog.Debug(ctx, "updating role state", map[string]interface{}{
		"roleResponse": roleResponse,
		"roleState":    roleState,
	})

	err = helper.CopyFields(ctx, roleResponse.Roles[0], &roleState)
	if err != nil {
		errStr := constants.ReadRoleErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading role",
			fmt.Sprintf("Could not read role struct %s with error: %s", roleID, message),
		)
		return
	}

	diags := resp.State.Set(ctx, roleState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read role completed")
}
