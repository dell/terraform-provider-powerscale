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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NtpServerResource{}
var _ resource.ResourceWithConfigure = &NtpServerResource{}
var _ resource.ResourceWithImportState = &NtpServerResource{}

// NewNtpServerResource creates a new resource.
func NewNtpServerResource() resource.Resource {
	return &NtpServerResource{}
}

// NtpServerResource defines the resource implementation.
type NtpServerResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *NtpServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ntpserver"
}

// Schema describes the resource arguments.
func (r *NtpServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the NTP Server entity of PowerScale Array. We can Create, Update and Delete the NTP Server using this resource. We can also import an existing NTP Server from PowerScale array.",
		Description:         "This resource is used to manage the NTP Server entity of PowerScale Array. We can Create, Update and Delete the NTP Server using this resource. We can also import an existing NTP Server from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Description:         "Key value from key_file that maps to this server.",
				MarkdownDescription: "Key value from key_file that maps to this server.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "NTP server name.",
				MarkdownDescription: "NTP server name.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Field ID.",
				MarkdownDescription: "Field ID.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NtpServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NtpServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating ntp server")

	var plan models.NtpServerResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ntpServerToCreate := powerscale.V3NtpServer{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &ntpServerToCreate)
	if err != nil {
		errStr := constants.CreateNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp server",
			fmt.Sprintf("Could not read ntp server param with error: %s", message),
		)
		return
	}
	ntpServerID, err := helper.CreateNtpServer(ctx, r.client, ntpServerToCreate)
	if err != nil {
		errStr := constants.CreateNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp server",
			message,
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("ntp server %s created", ntpServerID), map[string]interface{}{
		"ntpServerResponse": ntpServerID,
	})

	plan.ID = types.StringValue(ntpServerID)
	getNtpServerResponse, err := helper.GetNtpServer(ctx, r.client, plan)
	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp server",
			message,
		)
		return
	}

	// update resource state according to response
	if len(getNtpServerResponse.Servers) <= 0 {
		resp.Diagnostics.AddError(
			"Error creating ntp server",
			fmt.Sprintf("Could not read created ntp server %s", ntpServerID),
		)
		return
	}

	createdNtpServer := getNtpServerResponse.Servers[0]
	err = helper.CopyFields(ctx, createdNtpServer, &plan)
	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp server",
			fmt.Sprintf("Could not read ntp server struct %s with error: %s", ntpServerID, message),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create ntp server completed")
}

// Read reads the resource state.
func (r *NtpServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading ntp server")

	var ntpServerState models.NtpServerResourceModel
	diags := req.State.Get(ctx, &ntpServerState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ntpServerID := ntpServerState.ID
	tflog.Debug(ctx, "calling get ntp server by ID", map[string]interface{}{
		"ntpServerID": ntpServerID,
	})
	ntpServerResponse, err := helper.GetNtpServer(ctx, r.client, ntpServerState)
	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading ntp server",
			message,
		)
		return
	}

	if len(ntpServerResponse.Servers) <= 0 {
		resp.Diagnostics.AddError(
			"Error reading ntp server",
			fmt.Sprintf("Could not read ntp server %s from powerscale with error: ntp server not found", ntpServerID),
		)
		return
	}
	tflog.Debug(ctx, "updating read ntp server state", map[string]interface{}{
		"ntpServerResponse": ntpServerResponse,
		"ntpServerState":    ntpServerState,
	})
	err = helper.CopyFields(ctx, ntpServerResponse.Servers[0], &ntpServerState)
	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading ntp server",
			fmt.Sprintf("Could not read ntp server struct %s with error: %s", ntpServerID, message),
		)
		return
	}

	diags = resp.State.Set(ctx, ntpServerState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read ntp server completed")
}

// Update updates the resource state.
func (r *NtpServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating ntp server")

	var ntpServerPlan models.NtpServerResourceModel
	diags := req.Plan.Get(ctx, &ntpServerPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ntpServerState models.NtpServerResourceModel
	diags = resp.State.Get(ctx, &ntpServerState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update ntp server", map[string]interface{}{
		"ntpServerPlan":  ntpServerPlan,
		"ntpServerState": ntpServerState,
	})

	if helper.IsUpdateNtpServerParamInvalid(ntpServerPlan, ntpServerState) {
		resp.Diagnostics.AddError(
			"Error updating ntp server",
			"Should not provide parameters for creating",
		)
		return
	}

	ntpServerID := ntpServerState.ID.ValueString()
	ntpServerPlan.ID = ntpServerState.ID
	var ntpServerToUpdate powerscale.V3NtpServerExtendedExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, ntpServerPlan, &ntpServerToUpdate)
	if err != nil {
		errStr := constants.UpdateNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp server",
			fmt.Sprintf("Could not read ntp server param with error: %s", message),
		)
		return
	}
	err = helper.UpdateNtpServer(ctx, r.client, ntpServerState.ID.ValueString(), ntpServerToUpdate)
	if err != nil {
		errStr := constants.UpdateNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp server",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get ntp server by ID on powerscale client", map[string]interface{}{
		"ntpServerID": ntpServerID,
	})
	updatedNtpServer, err := helper.GetNtpServer(ctx, r.client, ntpServerPlan)
	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp server",
			message,
		)
		return
	}

	if len(updatedNtpServer.Servers) <= 0 {
		resp.Diagnostics.AddError(
			"Error updating ntp server",
			fmt.Sprintf("Could not read updated ntp server %s", ntpServerID),
		)
		return
	}

	err = helper.CopyFields(ctx, updatedNtpServer.Servers[0], &ntpServerPlan)
	if err != nil {
		errStr := constants.ReadNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp server",
			fmt.Sprintf("Could not read ntp server struct %s with error: %s", ntpServerID, message),
		)
		return
	}
	diags = resp.State.Set(ctx, ntpServerPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update ntp server completed")
}

// Delete deletes the resource.
func (r *NtpServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting ntp server")

	var ntpServerState models.NtpServerResourceModel
	diags := req.State.Get(ctx, &ntpServerState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ntpServerID := ntpServerState.ID.ValueString()
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete ntp server on powerscale client", map[string]interface{}{
		"ntpServerID": ntpServerID,
	})
	err := helper.DeleteNtpServer(ctx, r.client, ntpServerState.ID.ValueString())
	if err != nil {
		errStr := constants.DeleteNtpServerErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting ntp server",
			message,
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete ntp server completed")
}

// ImportState imports the resource state.
func (r *NtpServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
