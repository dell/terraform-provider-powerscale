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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &UserMappingRulesResource{}
	_ resource.ResourceWithConfigure   = &UserMappingRulesResource{}
	_ resource.ResourceWithImportState = &UserMappingRulesResource{}
)

// NewUserMappingRulesResource creates a new resource.
func NewUserMappingRulesResource() resource.Resource {
	return &UserMappingRulesResource{}
}

// UserMappingRulesResource defines the resource implementation.
type UserMappingRulesResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *UserMappingRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_mapping_rules"
}

// Schema describes the resource arguments.
func (r *UserMappingRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the User Mapping Rules entity of PowerScale Array. " +
			"PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules." +
			"We can Create, Update and Delete the User Mapping Rules using this resource. We can also import an existing User Mapping Rules from PowerScale array. " +
			"Note that, User Mapping Rules is the native functionality of PowerScale. When creating the resource, we actually load User Mapping Rules from PowerScale to the resource state. ",
		Description: "This resource is used to manage the User Mapping Rules entity of PowerScale Array. " +
			"PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules." +
			"We can Create, Update and Delete the User Mapping Rules using this resource. We can also import an existing User Mapping Rules from PowerScale array. " +
			"Note that, User Mapping Rules is the native functionality of PowerScale. When creating the resource, we actually load User Mapping Rules from PowerScale to the resource state. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "User Mapping Rules ID.",
				MarkdownDescription: "User Mapping Rules ID.",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "The zone to which the user mapping applies.",
				MarkdownDescription: "The zone to which the user mapping applies.",
				Optional:            true,
			},
			"parameters": schema.SingleNestedAttribute{
				Description:         "Specifies the parameters for user mapping rules.",
				MarkdownDescription: "Specifies the parameters for user mapping rules.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"default_unix_user": schema.SingleNestedAttribute{
						Description:         "Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information.",
						MarkdownDescription: "Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"domain": schema.StringAttribute{
								Description:         "Specifies the domain of the user that is being mapped.",
								MarkdownDescription: "Specifies the domain of the user that is being mapped.",
								Optional:            true,
								Computed:            true,
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
							"user": schema.StringAttribute{
								Description:         "Specifies the name of the user that is being mapped.",
								MarkdownDescription: "Specifies the name of the user that is being mapped.",
								Required:            true,
								Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
							},
						},
					},
				},
			},
			"rules": schema.ListNestedAttribute{
				Description:         "Specifies the list of user mapping rules.",
				MarkdownDescription: "Specifies the list of user mapping rules.",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operator": schema.StringAttribute{
							Description:         "Specifies the operator to make rules on specified users or groups.",
							MarkdownDescription: "Specifies the operator to make rules on specified users or groups.",
							Required:            true,
							Validators:          []validator.String{stringvalidator.OneOf("append", "insert", "replace", "trim", "union")},
						},
						"options": schema.SingleNestedAttribute{
							Description:         "Specifies the mapping options for this user mapping rule.",
							MarkdownDescription: "Specifies the mapping options for this user mapping rule.",
							Optional:            true,
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"break": schema.BoolAttribute{
									Description:         "If true, and the rule was applied successfully, stop processing further.",
									MarkdownDescription: "If true, and the rule was applied successfully, stop processing further.",
									Optional:            true,
									Computed:            true,
								},
								"user": schema.BoolAttribute{
									Description:         "If true, the primary UID and primary user SID should be copied to the existing credential.",
									MarkdownDescription: "If true, the primary UID and primary user SID should be copied to the existing credential.",
									Optional:            true,
									Computed:            true,
								},
								"group": schema.BoolAttribute{
									Description:         "If true, the primary GID and primary group SID should be copied to the existing credential.",
									MarkdownDescription: "If true, the primary GID and primary group SID should be copied to the existing credential.",
									Optional:            true,
									Computed:            true,
								},
								"groups": schema.BoolAttribute{
									Description:         "If true, all additional identifiers should be copied to the existing credential.",
									MarkdownDescription: "If true, all additional identifiers should be copied to the existing credential.",
									Optional:            true,
									Computed:            true,
								},
								"default_user": schema.SingleNestedAttribute{
									Description:         "Specifies the default user information that can be applied if the final credentials do not have valid UID and GID information.",
									MarkdownDescription: "Specifies the default user information that can be applied if the final credentials do not have valid UID and GID information.",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"domain": schema.StringAttribute{
											Description:         "Specifies the domain of the user that is being mapped.",
											MarkdownDescription: "Specifies the domain of the user that is being mapped.",
											Optional:            true,
											Computed:            true,
											Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
										},
										"user": schema.StringAttribute{
											Description:         "Specifies the name of the user that is being mapped.",
											MarkdownDescription: "Specifies the name of the user that is being mapped.",
											Required:            true,
											Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
										},
									},
								},
							},
						},
						"target_user": schema.SingleNestedAttribute{
							Description:         "Specifies the target user information that the rule can be applied to.",
							MarkdownDescription: "Specifies the target user information that the rule can be applied to.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"domain": schema.StringAttribute{
									Description:         "Specifies the domain of the user that is being mapped.",
									MarkdownDescription: "Specifies the domain of the user that is being mapped.",
									Optional:            true,
									Computed:            true,
									Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
								},
								"user": schema.StringAttribute{
									Description:         "Specifies the name of the user that is being mapped.",
									MarkdownDescription: "Specifies the name of the user that is being mapped.",
									Required:            true,
									Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
								},
							},
						},
						"source_user": schema.SingleNestedAttribute{
							Description:         "Specifies the source user information that the rule can be applied from.",
							MarkdownDescription: "Specifies the source user information that the rule can be applied from.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"domain": schema.StringAttribute{
									Description:         "Specifies the domain of the user that is being mapped.",
									MarkdownDescription: "Specifies the domain of the user that is being mapped.",
									Optional:            true,
									Computed:            true,
									Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
								},
								"user": schema.StringAttribute{
									Description:         "Specifies the name of the user that is being mapped.",
									MarkdownDescription: "Specifies the name of the user that is being mapped.",
									Required:            true,
									Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
								},
							},
						},
					},
				},
			},
			"test_mapping_users": schema.ListNestedAttribute{
				Description:         "List of user identity for mapping test.",
				MarkdownDescription: "List of user identity for mapping test.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Specifies a user name.",
							MarkdownDescription: "Specifies a user name.",
							Optional:            true,
							Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
						},
						"uid": schema.Int32Attribute{
							Description:         "Specifies a numeric user identifier.",
							MarkdownDescription: "Specifies a numeric user identifier.",
							Optional:            true,
						},
					},
				},
			},
			"mapping_users": schema.ListNestedAttribute{
				Description:         "List of test mapping user result.",
				MarkdownDescription: "List of test mapping user result.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"zone": schema.StringAttribute{
							Description:         "Name of the access zone which contains this user. ",
							MarkdownDescription: "Name of the access zone which contains this user. ",
							Computed:            true,
						},
						"zid": schema.Int64Attribute{
							Description:         "Numeric ID of the access zone which contains this user. ",
							MarkdownDescription: "Numeric ID of the access zone which contains this user. ",
							Computed:            true,
						},
						"privileges": schema.ListNestedAttribute{
							Description:         "Specifies the system-defined privilege that may be granted to users. ",
							MarkdownDescription: "Specifies the system-defined privilege that may be granted to users. ",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description:         "Specifies the name of the privilege. ",
										MarkdownDescription: "Specifies the name of the privilege. ",
										Computed:            true,
									},
									"id": schema.StringAttribute{
										Description:         "Specifies the ID of the privilege. ",
										MarkdownDescription: "Specifies the ID of the privilege. ",
										Computed:            true,
									},
									"read_only": schema.BoolAttribute{
										Description:         "True, if the privilege is read-only. ",
										MarkdownDescription: "True, if the privilege is read-only. ",
										Computed:            true,
									},
								},
							},
						},
						"supplemental_identities": schema.ListNestedAttribute{
							Description:         "Specifies the configuration properties for a user.",
							MarkdownDescription: "Specifies the configuration properties for a user.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description:         "Specifies a user or group name.",
										MarkdownDescription: "Specifies a user or group name.",
										Computed:            true,
									},
									"gid": schema.StringAttribute{
										Description:         "Specifies a user or group GID.",
										MarkdownDescription: "Specifies a user or group GID.",
										Computed:            true,
									},
									"sid": schema.StringAttribute{
										Description:         "Specifies a user or group SID.",
										MarkdownDescription: "Specifies a user or group SID.",
										Computed:            true,
									},
								},
							},
						},
						"user": schema.SingleNestedAttribute{
							Description:         "Specifies the configuration properties for a user.",
							MarkdownDescription: "Specifies the configuration properties for a user.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"uid": schema.StringAttribute{
									Description:         "Specifies the user UID.",
									MarkdownDescription: "Specifies the user UID.",
									Computed:            true,
								},
								"sid": schema.StringAttribute{
									Description:         "Specifies a user or group SID.",
									MarkdownDescription: "Specifies a user or group SID.",
									Computed:            true,
								},
								"primary_group_sid": schema.StringAttribute{
									Description:         "Specifies the primary group SID.",
									MarkdownDescription: "Specifies the primary group SID.",
									Computed:            true,
								},
								"primary_group_name": schema.StringAttribute{
									Description:         "Specifies the primary group name.",
									MarkdownDescription: "Specifies the primary group name.",
									Computed:            true,
								},
								"on_disk_user_identity": schema.StringAttribute{
									Description:         "Specifies the user identity on disk.",
									MarkdownDescription: "Specifies the user identity on disk.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Specifies the user name.",
									MarkdownDescription: "Specifies the user name.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *UserMappingRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *UserMappingRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating User Mapping Rules resource state")
	// Read Terraform plan into the model
	var plan models.UserMappingRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.UserMappingRulesResourceModel
	if !plan.Rules.IsUnknown() || !plan.Parameters.IsUnknown() {
		if diags := helper.UpdateUserMappingRules(ctx, r.client, &state, &plan); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	rulesResponse, err := helper.GetUserMappingRulesByZone(ctx, r.client, plan.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error getting user mapping rules", err.Error())
		return
	}

	if diags := helper.UpdateUserMappingRulesState(ctx, &plan, rulesResponse); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if diags := helper.UpdateLookupMappingUsersState(ctx, r.client, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create User Mapping Rules resource state")
}

// Read reads the resource state.
func (r *UserMappingRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading User Mapping Rules resource state")

	var state models.UserMappingRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rulesResponse, err := helper.GetUserMappingRulesByZone(ctx, r.client, state.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error getting user mapping rules", err.Error())
		return
	}

	if diags := helper.UpdateUserMappingRulesState(ctx, &state, rulesResponse); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if diags := helper.UpdateLookupMappingUsersState(ctx, r.client, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read User Mapping Rules resource state")
}

// Update updates the resource state.
func (r *UserMappingRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating User Mapping Rules resource state")

	var plan models.UserMappingRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.UserMappingRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.Rules.IsUnknown() || !plan.Parameters.IsUnknown() {
		if diags := helper.UpdateUserMappingRules(ctx, r.client, &state, &plan); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	rulesResponse, err := helper.GetUserMappingRulesByZone(ctx, r.client, plan.Zone.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error getting user mapping rules", err.Error())
		return
	}

	if diags := helper.UpdateUserMappingRulesState(ctx, &plan, rulesResponse); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if diags := helper.UpdateLookupMappingUsersState(ctx, r.client, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Update User Mapping Rules resource state")
}

// Delete deletes the resource.
func (r *UserMappingRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting User Mapping Rules resource state")
	var state models.UserMappingRulesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete User Mapping Rules resource state")
}

// ImportState imports the resource state.
func (r *UserMappingRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing User Mapping Rules resource state")
	var state models.UserMappingRulesResourceModel

	zone := req.ID
	rulesResponse, err := helper.GetUserMappingRulesByZone(ctx, r.client, zone)
	if err != nil {
		resp.Diagnostics.AddError("error getting user mapping rules", err.Error())
		return
	}

	if diags := helper.UpdateUserMappingRulesState(ctx, &state, rulesResponse); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if diags := helper.UpdateLookupMappingUsersState(ctx, r.client, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.Zone = types.StringValue(zone)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import User Mapping Rules resource state")
}
