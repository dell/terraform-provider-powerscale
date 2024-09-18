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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &S3ZoneSettingsResource{}
	_ resource.ResourceWithImportState = &S3ZoneSettingsResource{}
)

// NewS3ZoneSettingsResource is a helper function to simplify the provider implementation.
func NewS3ZoneSettingsResource() resource.Resource {
	return &S3ZoneSettingsResource{}
}

// S3ZoneSettingsResource is the resource implementation.
type S3ZoneSettingsResource struct {
	client *client.Client
}

// Configure implements resource.ResourceWithConfigure.
func (r *S3ZoneSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Metadata returns the resource type name.
func (r *S3ZoneSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_zone_settings"
}

// Schema defines the schema for the resource.
func (r *S3ZoneSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing S3ZoneSettings on PowerScale.",
		Description:         "This resource is used to manage the S3ZoneSettings of PowerScale Array. PowerScale S3 zone settings are used to configure the S3 zone of the S3 bucket.",
		Attributes:          S3ZoneSettingsSchema(),
	}
}

// S3ZoneSettingsSchema is a function that returns the schema for S3ZoneSettingsResource.
func S3ZoneSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"zone": schema.StringAttribute{
			MarkdownDescription: "The name of the access zone you want to update settings for s3 service",
			Description:         "The name of the access zone you want to update settings for s3 service",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIfConfigured(),
			},
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"base_domain": schema.StringAttribute{
			MarkdownDescription: "Base Domain for S3 zone",
			Description:         "Base Domain for S3 zone",
			Optional:            true,
			Computed:            true,
		},
		"bucket_directory_create_mode": schema.Int64Attribute{
			MarkdownDescription: " The permission mode for creating bucket directories.",
			Description:         " The permission mode for creating bucket directories.",
			Optional:            true,
			Computed:            true,
		},
		"object_acl_policy": schema.StringAttribute{
			MarkdownDescription: "The default policy for object access control lists (ACLs), which can be either “replace” or “deny”",
			Description:         "The default policy for object access control lists (ACLs), which can be either “replace” or “deny”",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("replace", "deny"),
			},
		},
		"root_path": schema.StringAttribute{
			MarkdownDescription: " The root path for the S3 bucket.",
			Description:         " The root path for the S3 bucket.",
			Optional:            true,
			Computed:            true,
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *S3ZoneSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_S3ZoneSettings create : Started")
	//Get Plan Data
	var plan, state models.S3ZoneSettingsResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := helper.SetZoneSetting(ctx, r.client, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating s3 zone settings ", err.Error())
		return
	}

	tflog.Trace(ctx, "resource_S3ZoneSettings create: updating state finished, saving ...")
	// Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_S3ZoneSettings create: finish")
}

// Read refreshes the Terraform state with the latest data.
func (r *S3ZoneSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_S3ZoneSettings read: started")
	var state models.S3ZoneSettingsResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.GetZoneSetting(ctx, r.client, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error reading s3 zone settings ", err.Error())
		return
	}

	tflog.Trace(ctx, "resource_S3ZoneSettings read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_S3ZoneSettings read: finished")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *S3ZoneSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Get state Data
	tflog.Trace(ctx, "resource_S3ZoneSettings update: started")
	var state, plan models.S3ZoneSettingsResource

	// Get plan Data
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := helper.SetZoneSetting(ctx, r.client, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating s3 zone settings ", err.Error())
		return
	}

	tflog.Trace(ctx, "resource_S3ZoneSettings update: finished state update")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "resource_S3ZoneSettings update: finished")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *S3ZoneSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_S3ZoneSettings delete: started")
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_S3ZoneSettings delete: finished")
}

// ImportState import state for existing S3ZoneSettings.
func (r S3ZoneSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Cannot import S3 Zone Settings with empty zone name.", "S3 Zone Settings do not have empty zone name.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("zone"), req, resp)
}
