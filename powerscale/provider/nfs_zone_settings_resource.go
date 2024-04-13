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
	"strings"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NfsZoneSettingsResource{}
	_ resource.ResourceWithConfigure   = &NfsZoneSettingsResource{}
	_ resource.ResourceWithImportState = &NfsZoneSettingsResource{}
)

// NewNfsZoneSettingsResource is a helper function to simplify the provider implementation.
func NewNfsZoneSettingsResource() resource.Resource {
	return &NfsZoneSettingsResource{}
}

// NfsZoneSettingsResource is the resource implementation.
type NfsZoneSettingsResource struct {
	client *client.Client
}

// Metadata defines the resource type name.
func (r *NfsZoneSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_zone_settings"
}

// Schema defines the schema for the resource.
func (r *NfsZoneSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `This resource is used to manage the NFS Zone Settings of PowerScale Array. We can Create, Update and Delete the NFS Zone Settings using this resource.  
		Note that, NFS Zone Settings is the native functionality of PowerScale. When creating the resource, we actually load NFS Zone Settings from PowerScale to the resource.`,
		Description: `This resource is used to manage the NFS Zone Settings of PowerScale Array. We can Create, Update and Delete the NFS Zone Settings using this resource.  
		Note that, NFS Zone Settings is the native functionality of PowerScale. When creating the resource, we actually load NFS Zone Settings from PowerScale to the resource.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of NFS Zone Settings. Value of ID will be same as the access zone.",
				MarkdownDescription: "ID of NFS Zone Settings. Value of ID will be same as the access zone.",
			},
			"zone": schema.StringAttribute{
				Required:            true,
				Description:         "Access zone",
				MarkdownDescription: "Access zone",
			},
			"nfsv4_no_names": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If true, sends owners and groups as UIDs and GIDs.",
				MarkdownDescription: "If true, sends owners and groups as UIDs and GIDs.",
			},
			"nfsv4_replace_domain": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If true, replaces the owner or group domain with an NFS domain name.",
				MarkdownDescription: "If true, replaces the owner or group domain with an NFS domain name.",
			},
			"nfsv4_allow_numeric_ids": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If true, sends owners and groups as UIDs and GIDs when look up fails or if the 'nfsv4_no_name' property is set to 1.",
				MarkdownDescription: "If true, sends owners and groups as UIDs and GIDs when look up fails or if the 'nfsv4_no_name' property is set to 1.",
			},
			"nfsv4_domain": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "Specifies the domain or realm through which users and groups are associated.",
				MarkdownDescription: "Specifies the domain or realm through which users and groups are associated.",
			},
			"nfsv4_no_domain": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If true, sends owners and groups without a domain name.",
				MarkdownDescription: "If true, sends owners and groups without a domain name.",
			},
			"nfsv4_no_domain_uids": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Description:         "If true, sends UIDs and GIDs without a domain name.",
				MarkdownDescription: "If true, sends UIDs and GIDs without a domain name.",
			},
		},
	}
}

// Configure configures the resource.
func (r *NfsZoneSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = pscaleClient
}

// Create creates the resource.
func (r *NfsZoneSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Started creating nfs zone settings")

	var plan models.NfsZoneSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone := plan.Zone.ValueString()

	var toUpdate powerscale.V2NfsSettingsZoneSettings
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs zone settings",
			fmt.Sprintf("Could not read nfs zone settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateNfsZoneSettings(ctx, r.client, toUpdate, zone)
	if err != nil {
		errStr := constants.UpdateNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs zone settings",
			message,
		)
		return
	}

	settings, err := helper.GetNfsZoneSettings(ctx, r.client, zone)
	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs zone settings",
			message,
		)
		return
	}

	var state models.NfsZoneSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs zone settings resource",
			err.Error(),
		)
		return
	}
	state.Zone = plan.Zone
	state.ID = plan.Zone

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed creating nfs zone settings")
}

// Read reads the resource.
func (r *NfsZoneSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Started reading nfs zone settings")

	var state models.NfsZoneSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone := state.Zone.ValueString()

	settings, err := helper.GetNfsZoneSettings(ctx, r.client, zone)
	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs zone settings",
			message,
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs zone settings resource",
			err.Error(),
		)
		return
	}
	state.Zone = types.StringValue(zone)
	state.ID = types.StringValue(zone)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed reading nfs zone settings")
}

// Update updates the resource.
func (r *NfsZoneSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Started updating nfs zone settings")

	var plan models.NfsZoneSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone := plan.Zone.ValueString()

	var toUpdate powerscale.V2NfsSettingsZoneSettings
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs zone settings",
			fmt.Sprintf("Could not read nfs zone settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateNfsZoneSettings(ctx, r.client, toUpdate, zone)
	if err != nil {
		errStr := constants.UpdateNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs zone settings",
			message,
		)
		return
	}

	settings, err := helper.GetNfsZoneSettings(ctx, r.client, zone)
	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs zone settings",
			message,
		)
		return
	}

	var state models.NfsZoneSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs zone settings resource",
			err.Error(),
		)
		return
	}
	state.Zone = plan.Zone
	state.ID = plan.Zone

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed updating nfs zone settings")
}

// Delete deletes the resource.
func (r *NfsZoneSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Started deleting nfs zone settings")

	// Read Terraform prior state data into the model
	var state models.NfsZoneSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Nfs zone settings is the native functionality that cannot be deleted, so just remove state
	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, "Completed deleting nfs zone settings")
}

// ImportState imports the resource.
func (r *NfsZoneSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Started importing nfs zone settings")

	reqID := req.ID
	zone := strings.TrimSpace(reqID)

	settings, err := helper.GetNfsZoneSettings(ctx, r.client, zone)
	if err != nil {
		errStr := constants.ReadNfsZoneSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs zone settings",
			message,
		)
		return
	}

	var state models.NfsZoneSettingsResourceModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs zone settings resource",
			err.Error(),
		)
		return
	}
	state.Zone = types.StringValue(zone)
	state.ID = types.StringValue(zone)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Completed importing nfs zone settings")
}
