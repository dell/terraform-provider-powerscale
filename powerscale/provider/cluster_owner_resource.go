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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ClusterOwnerResource{}
	_ resource.ResourceWithConfigure   = &ClusterOwnerResource{}
	_ resource.ResourceWithImportState = &ClusterOwnerResource{}
)

// NewClusterOwnerResource creates a new resource.
func NewClusterOwnerResource() resource.Resource {
	return &ClusterOwnerResource{}
}

// ClusterOwnerResource defines the resource implementation.
type ClusterOwnerResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *ClusterOwnerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_owner"
}

// Schema describes the resource arguments.
func (r *ClusterOwnerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the Cluster Owner Settings entity of PowerScale Array. " +
			"PowerScale Cluster Owner Settings provide the ability to configure owner settings on the cluster." +
			"We can Create, Update and Delete the Cluster Owner Settings using this resource. We can also import existing Cluster Owner Settings from PowerScale array. " +
			"Note that, Cluster Owner Settings is the native functionality of PowerScale. When creating the resource, we actually load Cluster Owner Settings from PowerScale to the resource state. ",
		Description: "This resource is used to manage the Cluster Owner Settings entity of PowerScale Array. " +
			"PowerScale Cluster Owner Settings provide the ability to configure Owner settings on the cluster." +
			"We can Create, Update and Delete the Cluster Owner Settings using this resource. We can also import existing Cluster Owner Settings from PowerScale array. " +
			"Note that, Cluster Owner Settings is the native functionality of PowerScale. When creating the resource, we actually load Cluster Owner Settings from PowerScale to the resource state. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Cluster Owner Settings.",
				MarkdownDescription: "ID of the Cluster Owner Settings.",
				Computed:            true,
			},
			"company": schema.StringAttribute{
				Description:         "Company Name of the Cluster Owner Settings.",
				MarkdownDescription: "Company Name of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"location": schema.StringAttribute{
				Description:         "Location of the Cluster Owner Settings.",
				MarkdownDescription: "Location of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"primary_email": schema.StringAttribute{
				Description:         "Primary Email of the Cluster Owner Settings.",
				MarkdownDescription: "Primary Email of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`), "must be a valid email"),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"primary_name": schema.StringAttribute{
				Description:         "Primary Name of the Cluster Owner Settings.",
				MarkdownDescription: "Primary Name of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"primary_phone1": schema.StringAttribute{
				Description:         "Primary Phone of the Cluster Owner Settings.",
				MarkdownDescription: "Primary Phone of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"primary_phone2": schema.StringAttribute{
				Description:         "Primary Alternate Phone of the Cluster Owner Settings.",
				MarkdownDescription: "Primary Alternate Phone of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"secondary_email": schema.StringAttribute{
				Description:         "Secondary Email of the Cluster Owner Settings.",
				MarkdownDescription: "Secondary Email of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`), "must be a valid email"),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"secondary_name": schema.StringAttribute{
				Description:         "Secondary Name of the Cluster Owner Settings.",
				MarkdownDescription: "Secondary Name of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"secondary_phone1": schema.StringAttribute{
				Description:         "Secondary Phone of the Cluster Owner Settings.",
				MarkdownDescription: "Secondary Phone of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"secondary_phone2": schema.StringAttribute{
				Description:         "Secondary Alternate Phone of the Cluster Owner Settings.",
				MarkdownDescription: "Secondary Alternate Phone of the Cluster Owner Settings.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *ClusterOwnerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ValidateConfig validate config of the resource.
func (r *ClusterOwnerResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Retrieve values from plan
	var cfg models.ClusterOwner

	diags := req.Config.Get(ctx, &cfg)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.Company.IsNull() && cfg.Location.IsNull() && cfg.PrimaryEmail.IsNull() && cfg.PrimaryName.IsNull() && cfg.PrimaryPhone1.IsNull() && cfg.PrimaryPhone2.IsNull() && cfg.SecondaryEmail.IsNull() && cfg.SecondaryName.IsNull() && cfg.SecondaryPhone1.IsNull() && cfg.SecondaryPhone2.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("company"),
			"Please provide at least one of the following: company, location, primary_email, primary_name, primary_phone1, primary_phone2, secondary_email, secondary_name, secondary_phone1, secondary_phone2",
			"",
		)
	}

}

// Create allocates the resource.
func (r *ClusterOwnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Cluster Owner Settings resource state")
	var plan models.ClusterOwner
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageClusterOwner(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue("cluster_owner")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster Owner resource state")
}

// Read reads the resource state.
func (r *ClusterOwnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Cluster Owner resource")
	var state models.ClusterOwner
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterOwner, err := helper.GetClusterOwner(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterOwnerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating cluster owner",
			message,
		)
		return
	}
	err = helper.CopyFields(ctx, clusterOwner, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of cluster owner resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("cluster_owner")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Owner resource")
}

// Update updates the resource state.
func (r *ClusterOwnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Cluster Owner resource...")
	var plan models.ClusterOwner
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageClusterOwner(ctx, r.client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue("cluster_email")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Update Cluster Owner resource")
}

// Delete deletes the resource.
func (r *ClusterOwnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Cluster Owner resource state")
	var state models.ClusterOwner

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Cluster Owner resource state")
}

// ImportState imports the resource state.
func (r *ClusterOwnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Cluster Owner resource")
	var state models.ClusterOwner
	clusterOwner, err := helper.GetClusterOwner(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterOwnerSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading cluster owner",
			message,
		)
		return
	}
	err = helper.CopyFields(ctx, clusterOwner, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error copying fields of cluster owner resource",
			err.Error(),
		)
		return
	}
	state.ID = types.StringValue("cluster_owner")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Owner resource")
}
