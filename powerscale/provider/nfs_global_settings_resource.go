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
	_ resource.Resource              = &NfsGlobalSettingsResource{}
	_ resource.ResourceWithConfigure = &NfsGlobalSettingsResource{}
)

// NewNfsGlobalSettingsResource creates a new resource.
func NewNfsGlobalSettingsResource() resource.Resource {
	return &NfsGlobalSettingsResource{}
}

// NfsGlobalSettingsResource defines the resource implementation.
type NfsGlobalSettingsResource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (r *NfsGlobalSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_global_settings"
}

// Schema describes the data source arguments.
func (r *NfsGlobalSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `This resource is used to manage the NFS Global Settings of PowerScale Array. We can Create, Update and Delete the NFS Global Settings using this resource.  
Note that, NFS Global Settings is the native functionality of PowerScale. When creating the resource, we actually load NFS Global Settings from PowerScale to the resource.`,
		Description: `This resource is used to manage the NFS Global Settings of PowerScale Array. We can Create, Update and Delete the NFS Global Settings using this resource.  
Note that, NFS Global Settings is the native functionality of PowerScale. When creating the resource, we actually load NFS Global Settings from PowerScale to the resource.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Id of NFS Global settings. Readonly. ",
				MarkdownDescription: "Id of NFS Global settings. Readonly. ",
			},
			"nfsv3_enabled": schema.BoolAttribute{
				Description:         "True if NFSv3 is enabled.",
				MarkdownDescription: "True if NFSv3 is enabled.",
				Optional:            true,
				Computed:            true,
			},
			"nfsv3_rdma_enabled": schema.BoolAttribute{
				Description:         "True if the RDMA is enabled for NFSv3.",
				MarkdownDescription: "True if the RDMA is enabled for NFSv3.",
				Optional:            true,
				Computed:            true,
			},
			"nfsv4_enabled": schema.BoolAttribute{
				Description:         "True if NFSv4 is enabled.",
				MarkdownDescription: "True if NFSv4 is enabled.",
				Optional:            true,
				Computed:            true,
			},
			"rpc_maxthreads": schema.Int64Attribute{
				Description:         "Specifies the maximum number of threads in the nfsd thread pool.",
				MarkdownDescription: "Specifies the maximum number of threads in the nfsd thread pool.",
				Optional:            true,
				Computed:            true,
			},
			"rpc_minthreads": schema.Int64Attribute{
				Description:         "Specifies the minimum number of threads in the nfsd thread pool.",
				MarkdownDescription: "Specifies the minimum number of threads in the nfsd thread pool.",
				Optional:            true,
				Computed:            true,
			},
			"rquota_enabled": schema.BoolAttribute{
				Description:         "True if the rquota protocol is enabled.",
				MarkdownDescription: "True if the rquota protocol is enabled.",
				Optional:            true,
				Computed:            true,
			},
			"service": schema.BoolAttribute{
				Description:         "True if the NFS service is enabled. When set to false, the NFS service is disabled.",
				MarkdownDescription: "True if the NFS service is enabled. When set to false, the NFS service is disabled.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NfsGlobalSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NfsGlobalSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating NFS Global Settings resource...")

	var plan models.NfsGlobalSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V12NfsSettingsGlobalSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs global settings",
			fmt.Sprintf("Could not read nfs global settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateNfsGlobalSettings(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs global settings",
			message,
		)
		return
	}

	settings, err := helper.GetNfsGlobalSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs global settings", message)
		return
	}

	var state models.NfsGlobalSettingsModel
	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs global settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_global_settings")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create nfs global settings resource")
}

// Read reads the resource state.
func (r *NfsGlobalSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Nfs Global Settings resource")

	var state models.NfsGlobalSettingsModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := helper.GetNfsGlobalSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs global settings", message)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs global settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_global_settings")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read nfs global settings resource")
}

// Update updates the resource state.
func (r *NfsGlobalSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating NFS Global Settings resource...")

	var plan models.NfsGlobalSettingsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.NfsGlobalSettingsModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var toUpdate powerscale.V12NfsSettingsGlobalSettings
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs global settings",
			fmt.Sprintf("Could not read nfs global settings param with error: %s", message),
		)
		return
	}

	err = helper.UpdateNfsGlobalSettings(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating nfs global settings",
			message,
		)
		return
	}

	settings, err := helper.GetNfsGlobalSettings(ctx, r.client)
	if err != nil {
		errStr := constants.ReadNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error reading nfs global settings", message)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, settings.GetSettings(), &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of nfs global settings resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("nfs_global_settings")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update nfs global settings resource")
}

// Delete deletes the resource.
func (r *NfsGlobalSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Nfs Global Settings resource")
	var state models.NfsGlobalSettingsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Nfs global settings is the native functionality that cannot be deleted, so just remove state
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete nfs global settings resource")
}

// ImportState imports the resource state.
func (r *NfsGlobalSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Nfs Global Settings resource")

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
