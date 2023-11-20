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
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &GroupnetResource{}
	_ resource.ResourceWithConfigure   = &GroupnetResource{}
	_ resource.ResourceWithImportState = &GroupnetResource{}
)

// NewGroupnetResource creates a new resource.
func NewGroupnetResource() resource.Resource {
	return &GroupnetResource{}
}

// GroupnetResource defines the resource implementation.
type GroupnetResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *GroupnetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groupnet"
}

// Schema describes the resource arguments.
func (r *GroupnetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the Groupnet entity of PowerScale Array. We can Create, Update and Delete the Groupnet using this resource. We can also import an existing Groupnet from PowerScale array. PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.",
		Description:         "This resource is used to manage the Groupnet entity of PowerScale Array. We can Create, Update and Delete the Groupnet using this resource. We can also import an existing Groupnet from PowerScale array. PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "The name of the groupnet. (Update Supported)",
				MarkdownDescription: "The name of the groupnet. (Update Supported)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9a-zA-Z_-]*$`), "must follow the pattern '^[0-9a-zA-Z_-]*$'",
					),
				},
			},
			"id": schema.StringAttribute{
				Description:         "Unique Interface ID.",
				MarkdownDescription: "Unique Interface ID.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A description of the groupnet. (Update Supported)",
				MarkdownDescription: "A description of the groupnet. (Update Supported)",
				Optional:            true,
				Validators:          []validator.String{stringvalidator.LengthBetween(1, 128)},
			},
			"allow_wildcard_subdomains": schema.BoolAttribute{
				Description:         "If enabled, SmartConnect treats subdomains of known dns zones as the known dns zone. This is required for S3 Virtual Host domains. Defaults to True. (Update Supported)",
				MarkdownDescription: "If enabled, SmartConnect treats subdomains of known dns zones as the known dns zone. This is required for S3 Virtual Host domains. Defaults to True. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"dns_cache_enabled": schema.BoolAttribute{
				Description:         "DNS caching is enabled or disabled. Defaults to True. (Update Supported)",
				MarkdownDescription: "DNS caching is enabled or disabled. Defaults to True. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"server_side_dns_search": schema.BoolAttribute{
				Description:         "Enable or disable appending nodes DNS search list to client DNS inquiries directed at SmartConnect service IP. Defaults to True. (Update Supported)",
				MarkdownDescription: "Enable or disable appending nodes DNS search list to client DNS inquiries directed at SmartConnect service IP. Defaults to True. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"dns_resolver_rotate": schema.BoolAttribute{
				Description:         "Enable or disable DNS resolver rotate. Defaults to False. (Update Supported)",
				MarkdownDescription: "Enable or disable DNS resolver rotate. Defaults to False. (Update Supported)",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"dns_search": schema.ListAttribute{
				Description:         "List of DNS search suffixes. (Update Supported)",
				MarkdownDescription: "List of DNS search suffixes. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 2048)),
					listvalidator.SizeBetween(0, 6),
				},
			},
			"dns_servers": schema.ListAttribute{
				Description:         "List of Domain Name Server IP addresses. (Update Supported)",
				MarkdownDescription: "List of Domain Name Server IP addresses. (Update Supported)",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 40)),
					listvalidator.SizeBetween(0, 3),
				},
			},
			"subnets": schema.ListAttribute{
				Description:         "Name of the subnets in the groupnet.",
				MarkdownDescription: "Name of the subnets in the groupnet.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *GroupnetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *GroupnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Groupnet resource...")
	var plan models.GroupnetModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupnetName := plan.Name.ValueString()
	if diags := helper.CreateGroupnet(ctx, r.client, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	groupnetResponse, err := helper.GetGroupnet(ctx, r.client, groupnetName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the groupnet - %s", groupnetName),
			err.Error(),
		)
		// if err, revert create
		_ = helper.DeleteGroupnet(ctx, r.client, groupnetName)
		return
	}

	// parse groupnet response to state groupnet model

	if err := helper.UpdateGroupnetResourceState(ctx, &plan, groupnetResponse); err != nil {
		resp.Diagnostics.AddError("Error creating groupnet Resource",
			fmt.Sprintf("Error parsing groupnet resource state: %s", err.Error()))
		// if err, revert create
		_ = helper.DeleteGroupnet(ctx, r.client, groupnetName)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create Groupnet resource")
}

// Read reads the resource state.
func (r *GroupnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Groupnet resource")
	var state models.GroupnetModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupnetName := state.Name.ValueString()
	groupnetResponse, err := helper.GetGroupnet(ctx, r.client, groupnetName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the Groupnet - %s", groupnetName),
			err.Error(),
		)
		return
	}

	// parse groupnet response to state groupnet model
	if err := helper.UpdateGroupnetResourceState(ctx, &state, groupnetResponse); err != nil {
		resp.Diagnostics.AddError("Error reading groupnet Resource",
			fmt.Sprintf("Error parsing groupnet resource state: %s", err.Error()))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Groupnet resource")
}

// Update updates the resource state.
func (r *GroupnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating Groupnet resource...")
	// Read Terraform plan into the model
	var plan models.GroupnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.GroupnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if diags := helper.UpdateGroupnet(ctx, r.client, &state, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	groupnetResponse, err := helper.GetGroupnet(ctx, r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the Groupnet - %s", plan.Name.ValueString()),
			err.Error(),
		)
		return
	}

	// parse groupnet response to state groupnet model
	if err := helper.UpdateGroupnetResourceState(ctx, &plan, groupnetResponse); err != nil {
		resp.Diagnostics.AddError("Error updating groupnet Resource",
			fmt.Sprintf("Error parsing groupnet resource state: %s", err.Error()))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Update Groupnet resource")

}

// Delete deletes the resource.
func (r *GroupnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Groupnet resource")
	var state models.GroupnetModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := helper.DeleteGroupnet(ctx, r.client, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting the Groupnet - %s", state.Name.ValueString()),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Groupnet resource")
}

// ImportState imports the resource state.
func (r *GroupnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Groupnet resource")
	var state models.GroupnetModel

	groupnetName := req.ID

	groupnetResponse, err := helper.GetGroupnet(ctx, r.client, groupnetName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error getting the Groupnet - %s", groupnetName),
			err.Error(),
		)
		return
	}

	// parse groupnet response to state groupnet model
	if err := helper.UpdateGroupnetImportState(ctx, &state, groupnetResponse); err != nil {
		resp.Diagnostics.AddError("Error importing groupnet Resource",
			fmt.Sprintf("Error parsing groupnet resource state: %s", err.Error()))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import Groupnet resource")
}
