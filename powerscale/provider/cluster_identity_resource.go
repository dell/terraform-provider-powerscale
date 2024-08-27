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

	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ClusterIdentityResource{}
	_ resource.ResourceWithConfigure   = &ClusterIdentityResource{}
	_ resource.ResourceWithImportState = &ClusterIdentityResource{}
)

func NewClusterIdentityResource() resource.Resource {
	return &ClusterIdentityResource{}
}

// ClusterIdentityResource represents a resource in the provider.
type ClusterIdentityResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *ClusterIdentityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_identity"
}

// Schema returns the schema for the resource.
func (r *ClusterIdentityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the Cluster Identity settings of PowerScale Array. " +
			"We can Create, Update and Delete the Cluster Identity using this resource. We can also import an existing Cluster Identity from PowerScale array.",
		Description: "This resource is used to manage the Cluster Identity settings of PowerScale Array. " +
			"We can Create, Update and Delete the Cluster Identity using this resource. We can also import an existing Cluster Identity from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The id for this cluster.",
				MarkdownDescription: "The id for this cluster.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "A unique name for this cluster.",
				MarkdownDescription: "A unique name for this cluster.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description of the cluster.",
				MarkdownDescription: "A description of the cluster.",
				Optional:            true,
				Computed:            true,
			},
			"logon": schema.SingleNestedAttribute{
				Description:         "The information displayed when a user logs in to the cluster.",
				MarkdownDescription: "The information displayed when a user logs in to the cluster.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"motd": schema.StringAttribute{
						Description:         "The message of the day.",
						MarkdownDescription: "The message of the day.",
						Optional:            true,
						Computed:            true,
					},
					"motd_header": schema.StringAttribute{
						Description:         "The header to the message of the day.",
						MarkdownDescription: "The header to the message of the day.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *ClusterIdentityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.
func (r *ClusterIdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Info(ctx, "Creating Cluster Identity Settings resource state")
	var plan, state models.ClusterIdentityResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ManageClusterIdentity(ctx, plan, &state, r.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Identity resource")

}

// Read reads the resource.
func (r *ClusterIdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Cluster Identity resource")

	var state models.ClusterIdentityResource

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterIdentity, err := helper.GetClusterIdentity(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterIdentitySettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading cluster identity",
			message,
		)
		return
	}

	// Save updated data into Terraform state
	diags := helper.UpdateClusterIdentityState(&state, clusterIdentity)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Identity resource")
}

// Update updates the resource.
func (r *ClusterIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Cluster Identity Settings resource state")
	var plan, state models.ClusterIdentityResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ManageClusterIdentity(ctx, plan, &state, r.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Identity resource")
}

// Delete deletes the resource.
func (r *ClusterIdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Cluster Identity resource state")
	var state models.ClusterIdentityResource

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Cluster Identity resource state")
}

// ImportState imports the resource.
func (r *ClusterIdentityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Cluster Identity resource state")

	var state models.ClusterIdentityResource
	clusterIdentity, err := helper.GetClusterIdentity(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving Cluster Identity resource",
			err.Error(),
		)
		return
	}

	diags := helper.UpdateClusterIdentityState(&state, clusterIdentity)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Identity resource")
}
