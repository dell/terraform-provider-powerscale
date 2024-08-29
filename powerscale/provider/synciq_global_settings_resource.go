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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource              = &SyncIQGlobalSettingsResource{}
	_ resource.ResourceWithConfigure = &SyncIQGlobalSettingsResource{}
	//_ resource.ResourceWithImportState = &SyncIQGlobalSettingsResource{}
)

// NewSyncIQGlobalSettingsResource creates a new resource.
func NewSyncIQGlobalSettingsResource() resource.Resource {
	return &SyncIQGlobalSettingsResource{}
}

// SyncIQGlobalSettingsResource defines the resource implementation.
type SyncIQGlobalSettingsResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *SyncIQGlobalSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_global_settings"
}

// Schema describes the resource arguments.
func (r *SyncIQGlobalSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the SyncIQ Global Settings entity of PowerScale Array. " +
			"We can Update the SyncIQ Global Settings using this resource. We can also import existing SyncIQ Global Settings from PowerScale array. ",
		Description: "This resource is used to manage the SyncIQ Global Settings entity of PowerScale Array. " +
			"We can Update the SyncIQ Global Settings using this resource. We can also import existing SyncIQ Global Settings from PowerScale array. ",
		Attributes: map[string]schema.Attribute{
			"preferred_rpo_alert": schema.Int64Attribute{
				Description:         "If specified, display as default RPO Alert value for new policy creation via WebUI.",
				MarkdownDescription: "If specified, display as default RPO Alert value for new policy creation via WebUI.",
				Optional:            true,
				Computed:            true,
			},
			"renegotiation_period": schema.Int64Attribute{
				Description:         "If specified, the duration to persist encrypted connection before forcing a renegotiation.",
				MarkdownDescription: "If specified, the duration to persist encrypted connection before forcing a renegotiation.",
				Optional:            true,
				Computed:            true,
			},
			"report_email": schema.SetAttribute{
				Description:         "Email sync reports to these addresses.",
				MarkdownDescription: "Email sync reports to these addresses.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"report_max_age": schema.Int64Attribute{
				Description:         "ID of the Cluster Email Settings.",
				MarkdownDescription: "ID of the Cluster Email Settings.",
				Optional:            true,
				Computed:            true,
			},
			"report_max_count": schema.Int64Attribute{
				Description:         "The default length of time (in seconds) a policy report will be stored.",
				MarkdownDescription: "The default length of time (in seconds) a policy report will be stored.",
				Optional:            true,
				Computed:            true,
			},
			"restrict_target_network": schema.BoolAttribute{
				Description:         "Default for the \"restrict_target_network\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.",
				MarkdownDescription: "Default for the \"restrict_target_network\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.",
				Optional:            true,
				Computed:            true,
			},
			"rpo_alerts": schema.BoolAttribute{
				Description:         "If disabled, no RPO alerts will be generated.",
				MarkdownDescription: "If disabled, no RPO alerts will be generated.",
				Optional:            true,
				Computed:            true,
			},
			"service": schema.StringAttribute{
				Description:         "Specifies if the SyncIQ service currently on, paused, or off.  If paused, all sync jobs will be paused.  If turned off, all jobs will be canceled.",
				MarkdownDescription: "Specifies if the SyncIQ service currently on, paused, or off.  If paused, all sync jobs will be paused.  If turned off, all jobs will be canceled.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"on",
						"off",
						"paused",
					),
				},
			},
			"bandwidth_reservation_reserve_percentage": schema.Int64Attribute{
				Description:         "The percentage of SyncIQ bandwidth to reserve for policies that did not specify a bandwidth reservation.",
				MarkdownDescription: "The percentage of SyncIQ bandwidth to reserve for policies that did not specify a bandwidth reservation.",
				Optional:            true,
				Computed:            true,
			},
			"cluster_certificate_id": schema.StringAttribute{
				Description:         "The ID of this cluster's certificate being used for encryption.",
				MarkdownDescription: "The ID of this cluster's certificate being used for encryption.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"encryption_cipher_list": schema.StringAttribute{
				Description:         "The cipher list being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.",
				MarkdownDescription: "The cipher list being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"encryption_required": schema.BoolAttribute{
				Description:         "If true, requires all SyncIQ policies to utilize encrypted communications.",
				MarkdownDescription: "If true, requires all SyncIQ policies to utilize encrypted communications.",
				Optional:            true,
				Computed:            true,
			},
			"force_interface": schema.BoolAttribute{
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  Default for the \"force_interface\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  Default for the \"force_interface\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.",
				Optional:            true,
				Computed:            true,
			},
			"max_concurrent_jobs": schema.Int64Attribute{
				Description:         "The max concurrent jobs that SyncIQ can support. This number is based on the size of the current cluster and the current SyncIQ worker throttle rule.",
				MarkdownDescription: "The max concurrent jobs that SyncIQ can support. This number is based on the size of the current cluster and the current SyncIQ worker throttle rule.",
				Optional:            true,
				Computed:            true,
			},
			"ocsp_address": schema.StringAttribute{
				Description:         "The address of the OCSP responder to which to connect.",
				MarkdownDescription: "The address of the OCSP responder to which to connect.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"ocsp_issuer_certificate_id": schema.StringAttribute{
				Description:         "The ID of the certificate authority that issued the certificate whose revocation status is being checked.",
				MarkdownDescription: "The ID of the certificate authority that issued the certificate whose revocation status is being checked.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"service_history_max_age": schema.Int64Attribute{
				Description:         "Maximum age of service information to maintain, in seconds.",
				MarkdownDescription: "Maximum age of service information to maintain, in seconds.",
				Optional:            true,
				Computed:            true,
			},
			"service_history_max_count": schema.Int64Attribute{
				Description:         "Maximum number of historical service information records to maintain.",
				MarkdownDescription: "Maximum number of historical service information records to maintain.",
				Optional:            true,
				Computed:            true,
			},
			"use_workers_per_node": schema.BoolAttribute{
				Description:         "If enabled, SyncIQ will use the deprecated workers_per_node field with worker pools functionality and limit workers accordingly.",
				MarkdownDescription: "If enabled, SyncIQ will use the deprecated workers_per_node field with worker pools functionality and limit workers accordingly.",
				Optional:            true,
				Computed:            true,
			},
			"bandwidth_reservation_reserve_absolute": schema.Int64Attribute{
				Description:         "The amount of SyncIQ bandwidth to reserve in kb/s for policies that did not specify a bandwidth reservation. This field takes precedence over bandwidth_reservation_reserve_percentage.",
				MarkdownDescription: "The amount of SyncIQ bandwidth to reserve in kb/s for policies that did not specify a bandwidth reservation. This field takes precedence over bandwidth_reservation_reserve_percentage.",
				Optional:            true,
				Computed:            true,
			},
			"source_network": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Restricts replication policies on the local cluster to running on the specified subnet and pool.",
				MarkdownDescription: "Restricts replication policies on the local cluster to running on the specified subnet and pool.",
				Attributes: map[string]schema.Attribute{
					"subnet": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "The subnet to restrict replication policies to.",
						MarkdownDescription: "The subnet to restrict replication policies to.",
					},
					"pool": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "The pool to restrict replication policies to.",
						MarkdownDescription: "The pool to restrict replication policies to.",
					},
				},
			},
			"tw_chkpt_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The interval (in seconds) in which treewalk syncs are forced to checkpoint.",
				MarkdownDescription: "The interval (in seconds) in which treewalk syncs are forced to checkpoint.",
			},
		},
	}
}

// Configure configures the resource.
func (r *SyncIQGlobalSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SyncIQGlobalSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating SyncIQ Global Settings resource state")
	// Read Terraform plan into the model
	var plan models.SyncIQGlobalSettingsModel
	var state models.SyncIQGlobalSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ManageSyncIQGlobalSettings(ctx, plan, &state, r.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read reads the resource state.
func (r *SyncIQGlobalSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading global settings resource")
	var state models.SyncIQGlobalSettingsModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags := helper.ManageReadSyncIQGlobalSettings(ctx, &state, r.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Cluster Email resource")
}

// Update updates the resource state.
func (r *SyncIQGlobalSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating SyncIQ global settings resource...")
	var plan models.SyncIQGlobalSettingsModel
	var state models.SyncIQGlobalSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ManageSyncIQGlobalSettings(ctx, plan, &state, r.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete deletes the resource.
func (r *SyncIQGlobalSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting Global Settings resource state")
	var state models.SyncIQGlobalSettingsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with Delete Global Settings resource state")
}

// ImportState imports the resource.
func (r *SyncIQGlobalSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing Cluster Identity resource state")

	var state models.SyncIQGlobalSettingsModel

	diags := helper.ManageReadSyncIQGlobalSettings(ctx, &state, r.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
