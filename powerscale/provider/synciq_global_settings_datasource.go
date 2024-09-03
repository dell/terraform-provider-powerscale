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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &SyncIQGlobalSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &SyncIQGlobalSettingsDataSource{}
)

// NewSyncIQGlobalSettingsDataSource creates a new syncIQ global settings data source.
func NewSyncIQGlobalSettingsDataSource() datasource.DataSource {
	return &SyncIQGlobalSettingsDataSource{}
}

// SyncIQGlobalSettingsDataSource defines the data source implementation.
type SyncIQGlobalSettingsDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SyncIQGlobalSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_global_settings"
}

// Schema describes the data source arguments.
func (d *SyncIQGlobalSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the SyncIQ Global Settings from PowerScale array. The information fetched from this datasource can be used for getting the details.",
		Description:         "This datasource is used to query the Cluster Email Settings from PowerScale array. The information fetched from this datasource can be used for getting the details.",

		Attributes: map[string]schema.Attribute{
			"force_interface": schema.BoolAttribute{
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  Default for the \"force_interface\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  Default for the \"force_interface\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.",
			},
			"preferred_rpo_alert": schema.Int64Attribute{
				Computed:            true,
				Description:         "If specified, display as default RPO Alert value for new policy creation via WebUI",
				MarkdownDescription: "If specified, display as default RPO Alert value for new policy creation via WebUI",
			},
			"use_workers_per_node": schema.BoolAttribute{
				Computed:            true,
				Description:         "If enabled, SyncIQ will use the deprecated workers_per_node field with worker pools functionality and limit workers accordingly.",
				MarkdownDescription: "If enabled, SyncIQ will use the deprecated workers_per_node field with worker pools functionality and limit workers accordingly.",
			},
			"ocsp_address": schema.StringAttribute{
				Computed:            true,
				Description:         "The address of the OCSP responder to which to connect.",
				MarkdownDescription: "The address of the OCSP responder to which to connect.",
			},
			"source_network": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Restricts replication policies on the local cluster to running on the specified subnet and pool.",
				MarkdownDescription: "Restricts replication policies on the local cluster to running on the specified subnet and pool.",
				Attributes: map[string]schema.Attribute{
					"subnet": schema.StringAttribute{
						Computed:            true,
						Description:         "The subnet to restrict replication policies to.",
						MarkdownDescription: "The subnet to restrict replication policies to.",
					},
					"pool": schema.StringAttribute{
						Computed:            true,
						Description:         "The pool to restrict replication policies to.",
						MarkdownDescription: "The pool to restrict replication policies to.",
					},
				},
			},
			"service": schema.StringAttribute{
				Computed:            true,
				Description:         "Specifies if the SyncIQ service currently on, paused, or off.  If paused, all sync jobs will be paused.  If turned off, all jobs will be canceled.",
				MarkdownDescription: "Specifies if the SyncIQ service currently on, paused, or off.  If paused, all sync jobs will be paused.  If turned off, all jobs will be canceled.",
			},
			"cluster_certificate_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of this cluster's certificate being used for encryption.",
				MarkdownDescription: "The ID of this cluster's certificate being used for encryption.",
			},
			"ocsp_issuer_certificate_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the certificate authority that issued the certificate whose revocation status is being checked.",
				MarkdownDescription: "The ID of the certificate authority that issued the certificate whose revocation status is being checked.",
			},
			"report_max_count": schema.Int64Attribute{
				Computed:            true,
				Description:         "The default maximum number of reports to retain for a policy.",
				MarkdownDescription: "The default maximum number of reports to retain for a policy.",
			},
			"rpo_alerts": schema.BoolAttribute{
				Computed:            true,
				Description:         "If disabled, no RPO alerts will be generated.",
				MarkdownDescription: "If disabled, no RPO alerts will be generated.",
			},
			"max_concurrent_jobs": schema.Int64Attribute{
				Computed:            true,
				Description:         "The max concurrent jobs that SyncIQ can support. This number is based on the size of the current cluster and the current SyncIQ worker throttle rule.",
				MarkdownDescription: "The max concurrent jobs that SyncIQ can support. This number is based on the size of the current cluster and the current SyncIQ worker throttle rule.",
			},
			"tw_chkpt_interval": schema.Int64Attribute{
				Computed:            true,
				Description:         "The interval (in seconds) in which treewalk syncs are forced to checkpoint.",
				MarkdownDescription: "The interval (in seconds) in which treewalk syncs are forced to checkpoint.",
			},
			"encryption_required": schema.BoolAttribute{
				Computed:            true,
				Description:         "If true, requires all SyncIQ policies to utilize encrypted communications.",
				MarkdownDescription: "If true, requires all SyncIQ policies to utilize encrypted communications.",
			},
			"bandwidth_reservation_reserve_percentage": schema.Int64Attribute{
				Computed:            true,
				Description:         "The percentage of SyncIQ bandwidth to reserve for policies that did not specify a bandwidth reservation.",
				MarkdownDescription: "The percentage of SyncIQ bandwidth to reserve for policies that did not specify a bandwidth reservation.",
			},
			"report_email": schema.SetAttribute{
				Computed:            true,
				Description:         "Email sync reports to these addresses.",
				MarkdownDescription: "Email sync reports to these addresses.",
				ElementType:         types.StringType,
			},
			"service_history_max_age": schema.Int64Attribute{
				Computed:            true,
				Description:         "Maximum age of service information to maintain, in seconds.",
				MarkdownDescription: "Maximum age of service information to maintain, in seconds.",
			},
			"renegotiation_period": schema.Int64Attribute{
				Computed:            true,
				Description:         "If specified, the duration to persist encrypted connection before forcing a renegotiation.",
				MarkdownDescription: "If specified, the duration to persist encrypted connection before forcing a renegotiation.",
			},
			"bandwidth_reservation_reserve_absolute": schema.Int64Attribute{
				Computed:            true,
				Description:         "The amount of SyncIQ bandwidth to reserve in kb/s for policies that did not specify a bandwidth reservation. This field takes precedence over bandwidth_reservation_reserve_percentage.",
				MarkdownDescription: "The amount of SyncIQ bandwidth to reserve in kb/s for policies that did not specify a bandwidth reservation. This field takes precedence over bandwidth_reservation_reserve_percentage.",
			},
			"report_max_age": schema.Int64Attribute{
				Computed:            true,
				Description:         "The default length of time (in seconds) a policy report will be stored.",
				MarkdownDescription: "The default length of time (in seconds) a policy report will be stored.",
			},
			"restrict_target_network": schema.BoolAttribute{
				Computed:            true,
				Description:         "Default for the \"restrict_target_network\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.",
				MarkdownDescription: "Default for the \"restrict_target_network\" property that will be applied to each new sync policy unless otherwise specified at the time of policy creation.  If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.",
			},
			"encryption_cipher_list": schema.StringAttribute{
				Computed:            true,
				Description:         "The cipher list being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.",
				MarkdownDescription: "The cipher list being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.",
			},
			"service_history_max_count": schema.Int64Attribute{
				Computed:            true,
				Description:         "Maximum number of historical service information records to maintain.",
				MarkdownDescription: "Maximum number of historical service information records to maintain.",
			},
			"password_set": schema.BoolAttribute{
				Computed:            true,
				Description:         "Indicates if a password is set for authentication. Password value is not shown with GET.",
				MarkdownDescription: "Indicates if a password is set for authentication. Password value is not shown with GET.",
			},
		},
	}
}

// Configure configures the data source.
func (d *SyncIQGlobalSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *SyncIQGlobalSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading SyncIQ Global Settings data source ")

	var state models.SyncIQGlobalSettingsDataSourceModel

	if resp.Diagnostics.HasError() {
		return
	}

	diags := helper.ManageReadDataSourceSyncIQGlobalSettings(ctx, &state, d.client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read SyncIQ Gloabl Settings data source ")
}
