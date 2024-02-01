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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ClusterEmailResource{}
	_ resource.ResourceWithConfigure   = &ClusterEmailResource{}
	_ resource.ResourceWithImportState = &ClusterEmailResource{}
)

// NewClusterEmailResource creates a new resource.
func NewClusterEmailResource() resource.Resource {
	return &ClusterEmailResource{}
}

// ClusterEmailResource defines the resource implementation.
type ClusterEmailResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *ClusterEmailResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_email"
}

// Schema describes the resource arguments.
func (r *ClusterEmailResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the Cluster Email Settings entity of PowerScale Array. " +
			"PowerScale Cluster Email Settings provide the ability to configure email settings on the cluster." +
			"We can Create, Update and Delete the Cluster Email Settings using this resource. We can also import existing Cluster Email Settings from PowerScale array. " +
			"Note that, Cluster Email Settings is the native functionality of PowerScale. When creating the resource, we actually load Cluster Email Settings from PowerScale to the resource state. ",
		Description: "This resource is used to manage the Cluster Email Settings entity of PowerScale Array. " +
			"PowerScale Cluster Email Settings provide the ability to configure email settings on the cluster." +
			"We can Create, Update and Delete the Cluster Email Settings using this resource. We can also import existing Cluster Email Settings from PowerScale array. " +
			"Note that, Cluster Email Settings is the native functionality of PowerScale. When creating the resource, we actually load Cluster Email Settings from PowerScale to the resource state. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Cluster Email Settings.",
				MarkdownDescription: "ID of the Cluster Email Settings.",
				Computed:            true,
			},
			"settings": schema.SingleNestedAttribute{
				Description:         "Cluster email notification settings.",
				MarkdownDescription: "Cluster email notification settings.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"batch_mode": schema.StringAttribute{
						Description:         "This setting determines how notifications will be batched together to be sent by email.  'none' means each notification will be sent separately.  'severity' means notifications of the same severity will be sent together.  'category' means notifications of the same category will be sent together.  'all' means all notifications will be batched together and sent in a single email.",
						MarkdownDescription: "This setting determines how notifications will be batched together to be sent by email.  'none' means each notification will be sent separately.  'severity' means notifications of the same severity will be sent together.  'category' means notifications of the same category will be sent together.  'all' means all notifications will be batched together and sent in a single email.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("all", "severity", "category", "none"),
						},
					},
					"mail_relay": schema.StringAttribute{
						Description:         "The address of the SMTP server to be used for relaying the notification messages.  An SMTP server is required in order to send notifications.  If this string is empty, no emails will be sent.",
						MarkdownDescription: "The address of the SMTP server to be used for relaying the notification messages.  An SMTP server is required in order to send notifications.  If this string is empty, no emails will be sent.",
						Optional:            true,
						Computed:            true,
					},
					"mail_sender": schema.StringAttribute{
						Description:         "The full email address that will appear as the sender of notification messages.",
						MarkdownDescription: "The full email address that will appear as the sender of notification messages.",
						Optional:            true,
						Computed:            true,
					},
					"mail_subject": schema.StringAttribute{
						Description:         "The subject line for notification messages from this cluster.",
						MarkdownDescription: "The subject line for notification messages from this cluster.",
						Optional:            true,
						Computed:            true,
					},
					"smtp_auth_passwd_set": schema.BoolAttribute{
						Description:         "Indicates if an SMTP authentication password is set.",
						MarkdownDescription: "Indicates if an SMTP authentication password is set.",
						Computed:            true,
					},
					"smtp_auth_passwd": schema.StringAttribute{
						Description:         "Password to authenticate with if SMTP authentication is being used.",
						MarkdownDescription: "Password to authenticate with if SMTP authentication is being used.",
						Optional:            true,
						Computed:            true,
					},
					"smtp_auth_security": schema.StringAttribute{
						Description:         "The type of secure communication protocol to use if SMTP is being used.  If 'none', plain text will be used, if 'starttls', the encrypted STARTTLS protocol will be used.",
						MarkdownDescription: "The type of secure communication protocol to use if SMTP is being used.  If 'none', plain text will be used, if 'starttls', the encrypted STARTTLS protocol will be used.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "starttls"),
						},
					},
					"smtp_auth_username": schema.StringAttribute{
						Description:         "Username to authenticate with if SMTP authentication is being used.",
						MarkdownDescription: "Username to authenticate with if SMTP authentication is being used.",
						Optional:            true,
						Computed:            true,
					},
					"smtp_port": schema.Int64Attribute{
						Description:         "The port on the SMTP server to be used for relaying the notification messages.",
						MarkdownDescription: "The port on the SMTP server to be used for relaying the notification messages.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.Int64{
							int64validator.Between(0, 65535),
						},
					},
					"use_smtp_auth": schema.BoolAttribute{
						Description:         "If true, this cluster will send SMTP authentication credentials to the SMTP relay server in order to send its notification emails.  If false, the cluster will attempt to send its notification emails without authentication.",
						MarkdownDescription: "If true, this cluster will send SMTP authentication credentials to the SMTP relay server in order to send its notification emails.  If false, the cluster will attempt to send its notification emails without authentication.",
						Optional:            true,
						Computed:            true,
					},
					"user_template": schema.StringAttribute{
						Description:         "Location of a custom template file that can be used to specify the layout of the notification emails.  If this string is empty, the default template will be used.",
						MarkdownDescription: "Location of a custom template file that can be used to specify the layout of the notification emails.  If this string is empty, the default template will be used.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *ClusterEmailResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

// Create allocates the resource.
func (r *ClusterEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Cluster Email Settings resource state")
	// Read Terraform plan into the model
	var plan models.ClusterEmail
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V1ClusterEmailExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan.Settings, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			fmt.Sprintf("Could not read cluster email param with error: %s", message),
		)
		return
	}
	// if the field is set to empty, set it to null to update back to default
	// if not set, unset the field to not update it
	if plan.Settings.UserTemplate.IsUnknown() {
		toUpdate.UserTemplate.Unset()
	} else if plan.Settings.UserTemplate.ValueString() == "" {
		toUpdate.UserTemplate.Set(nil)
	}
	err = helper.UpdateClusterEmail(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}

	clusterEmail, err := helper.GetClusterEmail(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}
	var state models.ClusterEmail
	err = helper.CopyFields(ctx, clusterEmail, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("cluster_email")
	state.Settings.SMTPAuthPasswd = plan.Settings.SMTPAuthPasswd
	if state.Settings.SMTPAuthPasswd.IsUnknown() {
		state.Settings.SMTPAuthPasswd = types.StringValue("")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster Email resource state")
}

// Read reads the resource state.
func (r *ClusterEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Cluster Email resource")
	var state models.ClusterEmail
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterEmail, err := helper.GetClusterEmail(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}
	err = helper.CopyFields(ctx, clusterEmail, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("cluster_email")
	if state.Settings.SMTPAuthPasswd.IsUnknown() {
		state.Settings.SMTPAuthPasswd = types.StringValue("")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Email resource")
}

// Update updates the resource state.
func (r *ClusterEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Cluster Email resource...")
	var plan models.ClusterEmail
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V1ClusterEmailExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan.Settings, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			fmt.Sprintf("Could not read cluster email param with error: %s", message),
		)
		return
	}
	if plan.Settings.UserTemplate.IsUnknown() {
		toUpdate.UserTemplate.Unset()
	} else if plan.Settings.UserTemplate.ValueString() == "" {
		toUpdate.UserTemplate.Set(nil)
	}
	err = helper.UpdateClusterEmail(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}

	clusterEmail, err := helper.GetClusterEmail(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster email",
			message,
		)
		return
	}
	var state models.ClusterEmail
	err = helper.CopyFields(ctx, clusterEmail, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("cluster_email")
	state.Settings.SMTPAuthPasswd = plan.Settings.SMTPAuthPasswd
	if state.Settings.SMTPAuthPasswd.IsUnknown() {
		state.Settings.SMTPAuthPasswd = types.StringValue("")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update Cluster Email resource")

}

// Delete deletes the resource.
func (r *ClusterEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Cluster Email resource state")
	var state models.ClusterEmail

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Cluster Email resource state")
}

// ImportState imports the resource state.
func (r *ClusterEmailResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Cluster Email resource")
	var state models.ClusterEmail
	clusterEmail, err := helper.GetClusterEmail(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading cluster email",
			message,
		)
		return
	}
	err = helper.CopyFields(ctx, clusterEmail, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("cluster_email")
	if state.Settings.SMTPAuthPasswd.IsUnknown() {
		state.Settings.SMTPAuthPasswd = types.StringValue("")
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Email resource")
}
