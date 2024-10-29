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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NfsAliasResource{}
var _ resource.ResourceWithImportState = &NfsAliasResource{}

// NewNfsAliasResource creates a new resource.
func NewNfsAliasResource() resource.Resource {
	return &NfsAliasResource{}
}

// NfsAliasResource defines the resource implementation.
type NfsAliasResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *NfsAliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_alias"
}

// Schema describes the resource arguments.
func (r *NfsAliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the NFS Alias entity of PowerScale Array. We can Create, Update and Delete the NFS Aliases using this resource. We can also import an existing NFS Alias from PowerScale array.",
		Description:         "This resource is used to manage the NFS Alias entity of PowerScale Array. We can Create, Update and Delete the NFS Aliases using this resource. We can also import an existing NFS Alias from PowerScale array.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "Specifies the name by which the alias can be referenced.",
				MarkdownDescription: "Specifies the name by which the alias can be referenced.",
				Required:            true,
			},
			"path": schema.StringAttribute{
				Description:         "Specifies the path to which the alias points.",
				MarkdownDescription: "Specifies the path to which the alias points.",
				Required:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "Specifies the zone in which the alias is valid.",
				MarkdownDescription: "Specifies the zone in which the alias is valid.",
				Optional:            true,
				Computed:            true,
			},
			"health": schema.StringAttribute{
				Description:         "Specifies whether the alias is usable.",
				MarkdownDescription: "Specifies whether the alias is usable.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Specifies a string which represents the unique location of the alias.",
				MarkdownDescription: "Specifies a string which represents the unique location of the alias.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *NfsAliasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NfsAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating nfs alias")
	var plan models.NfsAliasResourceModel
	var state models.NfsAliasResourceModel
	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.CreateNfsAlias(ctx, r.client, plan, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create NFS Alias resource state")
}

// Read reads the resource state.
func (r *NfsAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading nfs alias")
	var plan models.NfsAliasResourceModel
	var state models.NfsAliasResourceModel

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ReadNfsAlias(ctx, r.client, plan, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read NFS Alias resource state")
}

// Update updates the resource state Path, Name.
func (r *NfsAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating access zone")
	var plan *models.NfsAliasResourceModel
	var state *models.NfsAliasResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Populate the Edit Parameters
	editValues := powerscale.V2NfsAliasExtended{}
	if state.Zone != plan.Zone {
		resp.Diagnostics.AddError(
			"Error updating nfs alias",
			"Zone can't be updated",
		)
		return
	}
	if state.Name != plan.Name {
		editValues.Name = plan.Name.ValueStringPointer()
	}
	if state.Path != plan.Path {
		editValues.Path = plan.Path.ValueStringPointer()
	}

	diags := helper.UpdateNfsAlias(ctx, r.client, editValues, state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = helper.ReadNfsAlias(ctx, r.client, *plan, state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update NFS Alias resource state")
}

// Delete deletes the resource.
func (r *NfsAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting nfs alias")
	var data *models.NfsAliasResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteParam := r.client.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv2NfsAlias(ctx, data.ID.ValueString())
	_, err := deleteParam.Execute()
	if err != nil {
		errStr := constants.DeleteNfsAliasErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting access zone",
			message,
		)
	}
	tflog.Info(ctx, "Done with Delete NFS Alias resource state")
}

// ImportState imports the resource state.
func (r *NfsAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Please provide valid nfs alias ID", "Please provide valid nfs alias ID")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
