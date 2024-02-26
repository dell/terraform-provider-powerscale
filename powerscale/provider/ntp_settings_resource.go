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
var _ resource.Resource = &NtpSettingsResource{}
var _ resource.ResourceWithConfigure = &NtpSettingsResource{}
var _ resource.ResourceWithImportState = &NtpSettingsResource{}

// NewNtpSettingsResource creates a new resource.
func NewNtpSettingsResource() resource.Resource {
	return &NtpSettingsResource{}
}

// NtpSettingsResource defines the resource implementation.
type NtpSettingsResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *NtpSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ntpsettings"
}

// Schema describes the resource arguments.
func (r *NtpSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the NTP Settings entity of PowerScale Array. We can Create, Update and Delete the NTP Settings using this resource. " +
			"We can also import the existing NTP Settings from PowerScale array. Note that, NTP Settings is the native functionality of PowerScale. When creating the resource, we actually load NTP Settings from PowerScale to the resource state.",
		Description: "This resource is used to manage the NTP Settings entity of PowerScale Array. We can Create, Update and Delete the NTP Settings using this resource. " +
			"We can also import the existing NTP Settings from PowerScale array. Note that, NTP Settings is the native functionality of PowerScale. When creating the resource, we actually load NTP Settings from PowerScale to the resource state.",
		Attributes: map[string]schema.Attribute{
			"chimers": schema.Int64Attribute{
				Description:         "Number of nodes that will contact the NTP servers.",
				MarkdownDescription: "Number of nodes that will contact the NTP servers.",
				Optional:            true,
				Computed:            true,
			},
			"excluded": schema.ListAttribute{
				Description:         "Node number (LNN) for nodes excluded from chimer duty.",
				MarkdownDescription: "Node number (LNN) for nodes excluded from chimer duty.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"key_file": schema.StringAttribute{
				Description:         "Path to NTP key file within /ifs.",
				MarkdownDescription: "Path to NTP key file within /ifs.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NtpSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NtpSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating ntp settings")

	var plan models.NtpSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ntpSettingsToCreate := powerscale.V3NtpSettingsSettings{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &ntpSettingsToCreate)
	if err != nil {
		errStr := constants.CreateNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp settings",
			fmt.Sprintf("Could not read ntp settings param with error: %s", message),
		)
		return
	}
	err = helper.UpdateNtpSettings(ctx, r.client, ntpSettingsToCreate)
	if err != nil {
		errStr := constants.CreateNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp settings",
			message,
		)
		return
	}
	tflog.Debug(ctx, "ntp settings initialized")

	getNtpSettingsResponse, err := helper.GetNtpSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp settings",
			message,
		)
		return
	}

	createdNtpSettings := getNtpSettingsResponse.Settings
	err = helper.CopyFields(ctx, createdNtpSettings, &plan)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ntp settings",
			fmt.Sprintf("Could not read ntp settings struct with error: %s", message),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create ntp settings completed")
}

// Read reads the resource state.
func (r *NtpSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading ntp settings")

	var ntpSettingsState models.NtpSettingsResourceModel
	diags := req.State.Get(ctx, &ntpSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling get ntp settings")
	ntpSettingsResponse, err := helper.GetNtpSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading ntp settings",
			message,
		)
		return
	}

	tflog.Debug(ctx, "updating read ntp settings state", map[string]interface{}{
		"ntpSettingsResponse": ntpSettingsResponse,
		"ntpSettingsState":    ntpSettingsState,
	})
	err = helper.CopyFields(ctx, ntpSettingsResponse.Settings, &ntpSettingsState)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading ntp settings",
			fmt.Sprintf("Could not read ntp settings struct with error: %s", message),
		)
		return
	}

	diags = resp.State.Set(ctx, ntpSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read ntp settings completed")
}

// Update updates the resource state.
func (r *NtpSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating ntp settings")

	var ntpSettingsPlan models.NtpSettingsResourceModel
	diags := req.Plan.Get(ctx, &ntpSettingsPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ntpSettingsState models.NtpSettingsResourceModel
	diags = resp.State.Get(ctx, &ntpSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update ntp settings", map[string]interface{}{
		"ntpSettingsPlan":  ntpSettingsPlan,
		"ntpSettingsState": ntpSettingsState,
	})

	var ntpSettingsToUpdate powerscale.V3NtpSettingsSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, ntpSettingsPlan, &ntpSettingsToUpdate)
	if err != nil {
		errStr := constants.UpdateNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp settings",
			fmt.Sprintf("Could not read ntp settings param with error: %s", message),
		)
		return
	}
	err = helper.UpdateNtpSettings(ctx, r.client, ntpSettingsToUpdate)
	if err != nil {
		errStr := constants.UpdateNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp settings",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get ntp settings on powerscale client")
	updatedNtpSettings, err := helper.GetNtpSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp settings",
			message,
		)
		return
	}

	err = helper.CopyFields(ctx, updatedNtpSettings.Settings, &ntpSettingsPlan)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ntp settings",
			fmt.Sprintf("Could not read ntp settings struct with error: %s", message),
		)
		return
	}
	diags = resp.State.Set(ctx, ntpSettingsPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update ntp settings completed")
}

// Delete deletes the resource.
func (r *NtpSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting ntp settings")

	var ntpSettingsState models.NtpSettingsResourceModel
	diags := req.State.Get(ctx, &ntpSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete ntp settings completed")
}

// ImportState imports the resource state.
func (r *NtpSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var ntpSettingsState models.NtpSettingsResourceModel

	tflog.Debug(ctx, "calling get ntp settings")
	ntpSettingsResponse, err := helper.GetNtpSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing ntp settings",
			message,
		)
		return
	}

	err = helper.CopyFields(ctx, ntpSettingsResponse.Settings, &ntpSettingsState)
	if err != nil {
		errStr := constants.ReadNtpSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error importing ntp settings",
			fmt.Sprintf("Could not read ntp settings struct with error: %s", message),
		)
		return
	}

	diags := resp.State.Set(ctx, ntpSettingsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "import ntp settings completed")
}
