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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NamespaceACLResource{}
var _ resource.ResourceWithConfigure = &NamespaceACLResource{}
var _ resource.ResourceWithImportState = &NamespaceACLResource{}

// NewNamespaceACLResource creates a new resource.
func NewNamespaceACLResource() resource.Resource {
	return &NamespaceACLResource{}
}

// NamespaceACLResource defines the resource implementation.
type NamespaceACLResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *NamespaceACLResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_acl"
}

// Schema describes the resource arguments.
func (r *NamespaceACLResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the Namespace ACL on PowerScale Array. We can Create, Update and Delete the Namespace ACL using this resource. " +
			"We can also import the existing Namespace ACL from PowerScale array. Note that, when creating the resource, we actually load Namespace ACL from PowerScale to the resource state.",
		Description: "This resource is used to manage the Namespace ACL on PowerScale Array. We can Create, Update and Delete the Namespace ACL using this resource. " +
			"We can also import the existing Namespace ACL from PowerScale array. Note that, when creating the resource, we actually load Namespace ACL from PowerScale to the resource state.",
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
				Optional:            true,
				Computed:            true,
				Description:         "Provides the JSON object for the group persona of the owner.",
				MarkdownDescription: "Provides the JSON object for the group persona of the owner.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies the serialized form of a persona, which can be 'UID:0'",
						MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0'",
					},
					"name": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies the persona name, which must be combined with a type.",
						MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
					},
					"type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies the type of persona, which must be combined with a name.",
						MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
					},
				},
			},
			"group": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Provides the JSON object for the group persona of the owner.",
				MarkdownDescription: "Provides the JSON object for the group persona of the owner.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies the persona name, which must be combined with a type.",
						MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
					},
					"type": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies the type of persona, which must be combined with a name.",
						MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
					},
					"id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Specifies the serialized form of a persona, which can be 'GID:0'",
						MarkdownDescription: "Specifies the serialized form of a persona, which can be 'GID:0'",
					},
				},
			},
			"acl_custom": schema.ListNestedAttribute{
				Optional:            true,
				Description:         "Customer's raw configuration of the JSON array of access rights.",
				MarkdownDescription: "Customer's raw configuration of the JSON array of access rights.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"accesstype": schema.StringAttribute{
							Required:            true,
							Description:         "Grants or denies access control permissions. Options: allow, deny",
							MarkdownDescription: "Grants or denies access control permissions. Options: allow, deny",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"allow",
									"deny",
								),
							},
						},
						"op": schema.StringAttribute{
							Optional:            true,
							Description:         "Operations for updating access control permissions. Unnecessary for access right replacing scenario",
							MarkdownDescription: "Operations for updating access control permissions. Unnecessary for access right replacing scenario",
						},
						"inherit_flags": schema.ListAttribute{
							Optional:            true,
							Description:         "Grants or denies access control permissions. Options: object_inherit, container_inherit, inherit_only, no_prop_inherit, inherited_ace",
							MarkdownDescription: "Grants or denies access control permissions. Options: object_inherit, container_inherit, inherit_only, no_prop_inherit, inherited_ace",
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.UniqueValues(),
								listvalidator.ValueStringsAre(
									stringvalidator.OneOf(
										"object_inherit",
										"container_inherit",
										"inherit_only",
										"no_prop_inherit",
										"inherited_ace",
									),
								),
							},
						},
						"trustee": schema.SingleNestedAttribute{
							Required:            true,
							Description:         "Provides the JSON object for the group persona of the owner.",
							MarkdownDescription: "Provides the JSON object for the group persona of the owner.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Optional:            true,
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
								},
								"id": schema.StringAttribute{
									Optional:            true,
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0' or 'GID:0'",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0' or 'GID:0'",
								},
								"name": schema.StringAttribute{
									Optional:            true,
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
								},
							},
						},
						"accessrights": schema.ListAttribute{
							Optional:            true,
							Description:         "Specifies the access control permissions for a specific user or group. Options: std_delete, std_read_dac, std_write_dac, std_write_owner, std_synchronize, std_required, generic_all, generic_read, generic_write, generic_exec, dir_gen_all, dir_gen_read, dir_gen_write, dir_gen_execute, file_gen_all, file_gen_read, file_gen_write, file_gen_execute, modify, file_read, file_write, append, execute, file_read_attr, file_write_attr, file_read_ext_attr, file_write_ext_attr, delete_child, list, add_file, add_subdir, traverse, dir_read_attr, dir_write_attr, dir_read_ext_attr, dir_write_ext_attr",
							MarkdownDescription: "Specifies the access control permissions for a specific user or group. Options: std_delete, std_read_dac, std_write_dac, std_write_owner, std_synchronize, std_required, generic_all, generic_read, generic_write, generic_exec, dir_gen_all, dir_gen_read, dir_gen_write, dir_gen_execute, file_gen_all, file_gen_read, file_gen_write, file_gen_execute, modify, file_read, file_write, append, execute, file_read_attr, file_write_attr, file_read_ext_attr, file_write_ext_attr, delete_child, list, add_file, add_subdir, traverse, dir_read_attr, dir_write_attr, dir_read_ext_attr, dir_write_ext_attr",
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.UniqueValues(),
								listvalidator.ValueStringsAre(
									stringvalidator.OneOf(
										"std_delete", "std_read_dac", "std_write_dac", "std_write_owner",
										"std_synchronize", "std_required", "generic_all", "generic_read", "generic_write",
										"generic_exec", "dir_gen_all", "dir_gen_read", "dir_gen_write",
										"dir_gen_execute", "file_gen_all", "file_gen_read", "file_gen_write",
										"file_gen_execute", "modify", "file_read", "file_write", "append", "execute",
										"file_read_attr", "file_write_attr", "file_read_ext_attr", "file_write_ext_attr",
										"delete_child", "list", "add_file", "add_subdir", "traverse", "dir_read_attr",
										"dir_write_attr", "dir_read_ext_attr", "dir_write_ext_attr",
									),
								),
							},
						},
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
	}
}

// Configure configures the resource.
func (r *NamespaceACLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NamespaceACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating namespace acl")

	var plan models.NamespaceACLResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	namespaceACLToCreate := powerscale.NamespaceAcl{}
	plan.ACL = plan.CustomACL
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &namespaceACLToCreate)
	if err != nil {
		errStr := constants.CreateNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating namespace acl",
			fmt.Sprintf("Could not read namespace acl param with error: %s", message),
		)
		return
	}
	err = helper.CheckNamespaceACLParam(&namespaceACLToCreate)
	if err != nil {
		errStr := constants.CreateNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating namespace acl",
			fmt.Sprintf("Namespace acl param invalid with error: %s", message),
		)
		return
	}

	if namespaceACLToCreate.HasAcl() {
		err = helper.UpdateNamespaceACL(ctx, r.client, plan, namespaceACLToCreate)
		if err != nil {
			errStr := constants.CreateNamespaceACLErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error creating namespace acl",
				message,
			)
			return
		}
	}

	getNamespaceACLResponse, err := helper.GetNamespaceACL(ctx, r.client, plan)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating namespace acl",
			message,
		)
		return
	}

	originalPlan := plan
	err = helper.CopyFieldsToNonNestedModel(ctx, getNamespaceACLResponse, &plan)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating namespace acl",
			fmt.Sprintf("Could not read namespace acl struct with error: %s", message),
		)
		return
	}

	plan.CustomACL = originalPlan.CustomACL

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create namespace acl completed")
}

// Read reads the resource state.
func (r *NamespaceACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading namespace acl")

	var namespaceACLState models.NamespaceACLResourceModel
	diags := req.State.Get(ctx, &namespaceACLState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling get namespace acl")
	namespaceACLResponse, err := helper.GetNamespaceACL(ctx, r.client, namespaceACLState)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading namespace acl",
			message,
		)
		return
	}

	tflog.Debug(ctx, "updating read namespace acl state", map[string]interface{}{
		"namespaceACLResponse": namespaceACLResponse,
		"namespaceACLState":    namespaceACLState,
	})

	originalState := namespaceACLState
	err = helper.CopyFieldsToNonNestedModel(ctx, namespaceACLResponse, &namespaceACLState)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading namespace acl",
			fmt.Sprintf("Could not read namespace acl struct with error: %s", message),
		)
		return
	}

	namespaceACLState.CustomACL = originalState.CustomACL

	diags = resp.State.Set(ctx, namespaceACLState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read namespace acl completed")
}

// Update updates the resource state.
func (r *NamespaceACLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating namespace acl")

	var namespaceACLPlan models.NamespaceACLResourceModel
	diags := req.Plan.Get(ctx, &namespaceACLPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var namespaceACLState models.NamespaceACLResourceModel
	diags = resp.State.Get(ctx, &namespaceACLState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update namespace acl", map[string]interface{}{
		"namespaceACLPlan":  namespaceACLPlan,
		"namespaceACLState": namespaceACLState,
	})

	var namespaceACLToUpdate powerscale.NamespaceAcl
	namespaceACLPlan.ACL = namespaceACLPlan.CustomACL
	// Get param from tf input
	err := helper.ReadFromState(ctx, namespaceACLPlan, &namespaceACLToUpdate)
	if err != nil {
		errStr := constants.UpdateNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating namespace acl",
			fmt.Sprintf("Could not read namespace acl param with error: %s", message),
		)
		return
	}
	err = helper.CheckNamespaceACLParam(&namespaceACLToUpdate)
	if err != nil {
		errStr := constants.CreateNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating namespace acl",
			fmt.Sprintf("Namespace acl param invalid with error: %s", message),
		)
		return
	}
	err = helper.UpdateNamespaceACL(ctx, r.client, namespaceACLPlan, namespaceACLToUpdate)
	if err != nil {
		errStr := constants.UpdateNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating namespace acl",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get namespace acl on powerscale client")
	updatedNamespaceACL, err := helper.GetNamespaceACL(ctx, r.client, namespaceACLPlan)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating namespace acl",
			message,
		)
		return
	}

	originalPlan := namespaceACLPlan
	err = helper.CopyFieldsToNonNestedModel(ctx, updatedNamespaceACL, &namespaceACLPlan)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating namespace acl",
			fmt.Sprintf("Could not read namespace acl struct with error: %s", message),
		)
		return
	}

	namespaceACLPlan.CustomACL = originalPlan.CustomACL

	diags = resp.State.Set(ctx, namespaceACLPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update namespace acl completed")
}

// Delete deletes the resource.
func (r *NamespaceACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting namespace acl")

	var namespaceACLState models.NamespaceACLResourceModel
	diags := req.State.Get(ctx, &namespaceACLState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete namespace acl completed")
}

// ImportState imports the resource state.
func (r *NamespaceACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var namespaceACLModel models.NamespaceACLResourceModel
	namespaceACLModel.Namespace = types.StringValue(req.ID)

	tflog.Debug(ctx, "calling get namespace acl")
	namespaceACLResponse, err := helper.GetNamespaceACL(ctx, r.client, namespaceACLModel)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing namespace acl",
			message,
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, namespaceACLResponse, &namespaceACLModel)
	if err != nil {
		errStr := constants.ReadNamespaceACLErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing namespace acl",
			fmt.Sprintf("Could not read namespace acl struct with error: %s", message),
		)
		return
	}

	namespaceACLModel.CustomACL = namespaceACLModel.ACL
	diags := resp.State.Set(ctx, namespaceACLModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import namespace acl completed")
}
