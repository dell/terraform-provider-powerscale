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

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ClusterSnmpResource{}
var _ resource.ResourceWithConfigure = &ClusterSnmpResource{}
var _ resource.ResourceWithImportState = &ClusterSnmpResource{}

// NewClusterSnmpResource is a resource that manages ClusterSnmp entities.
func NewClusterSnmpResource() resource.Resource {
	return &ClusterSnmpResource{}
}

// ClusterSnmpResource represents SNMP resource struct.
type ClusterSnmpResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *ClusterSnmpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_snmp"
}

// Schema returns the schema for the resource.
func (r *ClusterSnmpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Schema describes the resource arguments.
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the Cluster SNMP settings of PowerScale Array. We can Create, Update and Delete the Cluster SNMP using this resource. We can also import an existing Cluster SNMP from PowerScale array.",
		Description:         "This resource is used to manage the Cluster SNMP settings of PowerScale Array. We can Create, Update and Delete the Cluster SNMP using this resource. We can also import an existing Cluster SNMP from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the Cluster SNMP.",
				MarkdownDescription: "ID of the Cluster SNMP.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "True if the Cluster SNMP is enabled.",
				MarkdownDescription: "True if the Cluster SNMP is enabled.",
				Required:            true,
				Validators:          []validator.Bool{boolvalidator.AtLeastOneOf(path.MatchRoot("snmp_v1_v2c_access"), path.MatchRoot("snmp_v3_access"))},
			},
			"read_only_community": schema.StringAttribute{
				Description:         "The read-only community string for the Cluster SNMP.",
				MarkdownDescription: "The read-only community string for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
			},
			"snmp_v1_v2c_access": schema.BoolAttribute{
				Description:         "The SNMPv1/v2c access for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv1/v2c access for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
				Validators:          []validator.Bool{boolvalidator.AlsoRequires(path.MatchRoot("read_only_community"))},
			},
			"snmp_v3_access": schema.BoolAttribute{
				Description:         "The SNMPv3 access for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 access for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
				Validators:          []validator.Bool{boolvalidator.AlsoRequires(path.MatchRoot("snmp_v3_password"))},
			},
			"snmp_v3_password": schema.StringAttribute{
				Description:         "The SNMPv3 authentication password for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 authentication password for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
				},
			},
			"snmp_v3_auth_protocol": schema.StringAttribute{
				Description:         "The SNMPv3 authentication protocol for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 authentication protocol for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"MD5",
						"SHA",
					),
				},
			},
			"snmp_v3_priv_protocol": schema.StringAttribute{
				Description:         "The SNMPv3 privacy protocol for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 privacy protocol for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
			},
			"snmp_v3_priv_password": schema.StringAttribute{
				Description:         "The SNMPv3 privacy protocol password for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 privacy protocol password for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
			},
			"snmp_v3_read_only_user": schema.StringAttribute{
				Description:         "The SNMPv3 read-only user for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 read-only user for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
			},
			"snmp_v3_security_level": schema.StringAttribute{
				Description:         "The SNMPv3 security level for the Cluster SNMP.",
				MarkdownDescription: "The SNMPv3 security level for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
			},
			"system_contact": schema.StringAttribute{
				Description:         "The system contact for the Cluster SNMP.",
				MarkdownDescription: "The system contact for the Cluster SNMP.",
				Computed:            true,
				Optional:            true,
			},
			"system_location": schema.StringAttribute{
				Description:         "The system location for the Cluster SNMP.",
				MarkdownDescription: "The system location for the Cluster SNMP.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Configure configures the resource.
func (r *ClusterSnmpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ClusterSnmpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating Cluster SNMP Settings resource state")
	// Read Terraform plan into the model
	var plan, state models.ClusterSNMPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.UpdateClusterSNMP(ctx, plan, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster SNMP resource state")
}

// Read reads the resource state.
func (r *ClusterSnmpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Cluster SNMP Settings resource state")

	var state models.ClusterSNMPModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterSNMP, err := helper.GetClusterSNMP(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterSNMPSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading cluster SNMP",
			message,
		)
		return
	}

	helper.UpdateclusterSNMPResourceState(ctx, &state, &state, clusterSNMP.Settings)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Done with Read Cluster SNMP resource state")
}

// Update updates the resource state.
func (r *ClusterSnmpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Creating Cluster SNMP Settings resource state")
	// Read Terraform plan into the model
	var plan, state models.ClusterSNMPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags := r.UpdateClusterSNMP(ctx, plan, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Create Cluster SNMP resource state")
}

// Delete deletes the resource.
func (r *ClusterSnmpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Cluster SNMP resource state")
	var state models.ClusterSNMPModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Cluster SNMP resource state")
}

// ImportState imports the resource state.
func (r *ClusterSnmpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Cluster SNMP resource")
	var state models.ClusterSNMPModel
	clusterSNMP, err := helper.GetClusterSNMP(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterSNMPSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading cluster SNMP",
			message,
		)
		return
	}
	helper.UpdateclusterSNMPResourceState(ctx, &state, &state, clusterSNMP.Settings)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster SNMP resource")
}

func mapBoolValue(v bool) *bool {
	return &v
}

// UpdateClusterSNMP is a common function that both Create and Update functions can call to update cluster SNMP.
func (r *ClusterSnmpResource) UpdateClusterSNMP(ctx context.Context, plan models.ClusterSNMPModel, state *models.ClusterSNMPModel) diag.Diagnostics {
	tflog.Info(ctx, "Creating Cluster SNMP Settings resource state")
	var diags diag.Diagnostics
	var toUpdate powerscale.V16SnmpSettingsExtended

	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterSNMPSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		diags.AddError(
			"Error updating cluster SNMP",
			fmt.Sprintf("Could not read cluster SNMP param with error: %s", message),
		)
		return diags
	}

	toUpdate.Service = mapBoolValue(plan.Service.ValueBool())
	err = helper.UpdateClusterSNMP(ctx, r.client, toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterSNMPSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		diags.AddError(
			"Error updating cluster SNMP",
			message,
		)
		return diags
	}

	tflog.Debug(ctx, "calling get cluster SNMP settings on powerscale client")
	clusterSNMP, err := helper.GetClusterSNMP(ctx, r.client)
	if err != nil {
		errStr := constants.ReadClusterSNMPSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		diags.AddError(
			"Error fetching updated cluster SNMP settings",
			message,
		)
		return diags
	}

	helper.UpdateclusterSNMPResourceState(ctx, &plan, state, clusterSNMP.Settings)

	state.SnmpV3Password = types.StringValue(plan.SnmpV3Password.ValueString())
	state.SnmpV3PrivPassword = types.StringValue(plan.SnmpV3PrivPassword.ValueString())
	return diags
}
