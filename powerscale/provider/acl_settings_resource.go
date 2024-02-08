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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ACLSettingsResource{}
var _ resource.ResourceWithConfigure = &ACLSettingsResource{}
var _ resource.ResourceWithImportState = &ACLSettingsResource{}

// NewACLSettingsResource creates a new resource.
func NewACLSettingsResource() resource.Resource {
	return &ACLSettingsResource{}
}

// ACLSettingsResource defines the resource implementation.
type ACLSettingsResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *ACLSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aclsettings"
}

// Schema describes the resource arguments.
func (r *ACLSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the ACL Settings entity of PowerScale Array. We can Create, Update and Delete the ACL Settings using this resource. " +
			"We can also import the existing ACL Settings from PowerScale array. Note that, ACL Settings is the native functionality of PowerScale. When creating the resource, we actually load ACL Settings from PowerScale to the resource state.",
		Description: "This resource is used to manage the ACL Settings entity of PowerScale Array. We can Create, Update and Delete the ACL Settings using this resource. " +
			"We can also import the existing ACL Settings from PowerScale array. Note that, ACL Settings is the native functionality of PowerScale. When creating the resource, we actually load ACL Settings from PowerScale to the resource state.",
		Attributes: map[string]schema.Attribute{
			"access": schema.StringAttribute{
				Description:         "Access checks (chmod, chown). Options: unix, windows",
				MarkdownDescription: "Access checks (chmod, chown). Options: unix, windows",
				Optional:            true,
				Computed:            true,
			},
			"calcmode": schema.StringAttribute{
				Description:         "Displayed mode bits. Options: approx, 777",
				MarkdownDescription: "Displayed mode bits. Options: approx, 777",
				Optional:            true,
				Computed:            true,
			},
			"calcmode_group": schema.StringAttribute{
				Description:         "Approximate group mode bits when ACL exists. Options: group_aces, group_only",
				MarkdownDescription: "Approximate group mode bits when ACL exists. Options: group_aces, group_only",
				Optional:            true,
				Computed:            true,
			},
			"calcmode_owner": schema.StringAttribute{
				Description:         "Approximate owner mode bits when ACL exists. Options: owner_aces, owner_only",
				MarkdownDescription: "Approximate owner mode bits when ACL exists. Options: owner_aces, owner_only",
				Optional:            true,
				Computed:            true,
			},
			"calcmode_traverse": schema.StringAttribute{
				Description:         "Require traverse rights in order to traverse directories with existing ACLs. Options: require, ignore",
				MarkdownDescription: "Require traverse rights in order to traverse directories with existing ACLs. Options: require, ignore",
				Optional:            true,
				Computed:            true,
			},
			"chmod": schema.StringAttribute{
				Description:         "chmod on files with existing ACLs. Options: remove, replace, replace_users_and_groups, merge_with_ugo_priority, merge, deny, ignore",
				MarkdownDescription: "chmod on files with existing ACLs. Options: remove, replace, replace_users_and_groups, merge_with_ugo_priority, merge, deny, ignore",
				Optional:            true,
				Computed:            true,
			},
			"chmod_007": schema.StringAttribute{
				Description:         "chmod (007) on files with existing ACLs. Options: default, remove",
				MarkdownDescription: "chmod (007) on files with existing ACLs. Options: default, remove",
				Optional:            true,
				Computed:            true,
			},
			"chmod_inheritable": schema.StringAttribute{
				Description:         "ACLs created on directories by UNIX chmod. Options: yes, no",
				MarkdownDescription: "ACLs created on directories by UNIX chmod. Options: yes, no",
				Optional:            true,
				Computed:            true,
			},
			"chown": schema.StringAttribute{
				Description:         "chown/chgrp on files with existing ACLs. Options: owner_group_and_acl, owner_group_only, ignore",
				MarkdownDescription: "chown/chgrp on files with existing ACLs. Options: owner_group_and_acl, owner_group_only, ignore",
				Optional:            true,
				Computed:            true,
			},
			"create_over_smb": schema.StringAttribute{
				Description:         "ACL creation over SMB. Options: allow, disallow",
				MarkdownDescription: "ACL creation over SMB. Options: allow, disallow",
				Optional:            true,
				Computed:            true,
			},
			"dos_attr": schema.StringAttribute{
				Description:         " Read only DOS attribute. Options: deny_smb, deny_smb_and_nfs",
				MarkdownDescription: " Read only DOS attribute. Options: deny_smb, deny_smb_and_nfs",
				Optional:            true,
				Computed:            true,
			},
			"group_owner_inheritance": schema.StringAttribute{
				Description:         "Group owner inheritance. Options: native, parent, creator",
				MarkdownDescription: "Group owner inheritance. Options: native, parent, creator",
				Optional:            true,
				Computed:            true,
			},
			"rwx": schema.StringAttribute{
				Description:         "Treatment of 'rwx' permissions. Options: retain, full_control",
				MarkdownDescription: "Treatment of 'rwx' permissions. Options: retain, full_control",
				Optional:            true,
				Computed:            true,
			},
			"synthetic_denies": schema.StringAttribute{
				Description:         "Synthetic 'deny' ACEs. Options: none, remove",
				MarkdownDescription: "Synthetic 'deny' ACEs. Options: none, remove",
				Optional:            true,
				Computed:            true,
			},
			"utimes": schema.StringAttribute{
				Description:         "Access check (utimes). Options: only_owner, owner_and_write",
				MarkdownDescription: "Access check (utimes). Options: only_owner, owner_and_write",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *ACLSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ACLSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating acl settings")

	var plan models.ACLSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	aclSettingsToCreate := powerscale.V11SettingsAclsAclPolicySettings{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &aclSettingsToCreate)
	if err != nil {
		errStr := constants.CreateACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating acl settings",
			fmt.Sprintf("Could not read acl settings param with error: %s", message),
		)
		return
	}
	err = helper.UpdateACLSettings(ctx, r.client, aclSettingsToCreate)
	if err != nil {
		errStr := constants.CreateACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating acl settings",
			message,
		)
		return
	}
	tflog.Debug(ctx, "acl settings initialized")

	getACLSettingsResponse, err := helper.GetACLSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating acl settings",
			message,
		)
		return
	}

	createdACLSettings := getACLSettingsResponse.AclPolicySettings
	err = helper.CopyFields(ctx, createdACLSettings, &plan)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating acl settings",
			fmt.Sprintf("Could not read acl settings struct with error: %s", message),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create acl settings completed")
}

// Read reads the resource state.
func (r *ACLSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading acl settings")

	var aclSettingsState models.ACLSettingsResourceModel
	diags := req.State.Get(ctx, &aclSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling get acl settings")
	aclSettingsResponse, err := helper.GetACLSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading acl settings",
			message,
		)
		return
	}

	tflog.Debug(ctx, "updating read acl settings state", map[string]interface{}{
		"aclSettingsResponse": aclSettingsResponse,
		"aclSettingsState":    aclSettingsState,
	})
	err = helper.CopyFields(ctx, aclSettingsResponse.AclPolicySettings, &aclSettingsState)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading acl settings",
			fmt.Sprintf("Could not read acl settings struct with error: %s", message),
		)
		return
	}

	diags = resp.State.Set(ctx, aclSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read acl settings completed")
}

// Update updates the resource state.
func (r *ACLSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating acl settings")

	var aclSettingsPlan models.ACLSettingsResourceModel
	diags := req.Plan.Get(ctx, &aclSettingsPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var aclSettingsState models.ACLSettingsResourceModel
	diags = resp.State.Get(ctx, &aclSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update acl settings", map[string]interface{}{
		"aclSettingsPlan":  aclSettingsPlan,
		"aclSettingsState": aclSettingsState,
	})

	var aclSettingsToUpdate powerscale.V11SettingsAclsAclPolicySettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, aclSettingsPlan, &aclSettingsToUpdate)
	if err != nil {
		errStr := constants.UpdateACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating acl settings",
			fmt.Sprintf("Could not read acl settings param with error: %s", message),
		)
		return
	}
	err = helper.UpdateACLSettings(ctx, r.client, aclSettingsToUpdate)
	if err != nil {
		errStr := constants.UpdateACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating acl settings",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get acl settings on powerscale client")
	updatedACLSettings, err := helper.GetACLSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating acl settings",
			message,
		)
		return
	}

	err = helper.CopyFields(ctx, updatedACLSettings.AclPolicySettings, &aclSettingsPlan)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating acl settings",
			fmt.Sprintf("Could not read acl settings struct with error: %s", message),
		)
		return
	}
	diags = resp.State.Set(ctx, aclSettingsPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update acl settings completed")
}

// Delete deletes the resource.
func (r *ACLSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting acl settings")

	var aclSettingsState models.ACLSettingsResourceModel
	diags := req.State.Get(ctx, &aclSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete acl settings completed")
}

// ImportState imports the resource state.
func (r *ACLSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var aclSettingsState models.ACLSettingsResourceModel

	tflog.Debug(ctx, "calling get acl settings")
	aclSettingsResponse, err := helper.GetACLSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing acl settings",
			message,
		)
		return
	}

	err = helper.CopyFields(ctx, aclSettingsResponse.AclPolicySettings, &aclSettingsState)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing acl settings",
			fmt.Sprintf("Could not read acl settings struct with error: %s", message),
		)
		return
	}

	diags := resp.State.Set(ctx, aclSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import acl settings completed")
}
