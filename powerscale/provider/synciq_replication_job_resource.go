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
	"net/http"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &synciqPolicyResource{}
	_ resource.ResourceWithConfigure   = &synciqPolicyResource{}
	_ resource.ResourceWithImportState = &synciqPolicyResource{}
)

const (
	paused  = "paused"
	running = "running"
)

// NewSyncIQReplicationJobResource is a helper function to simplify the provider implementation.
func NewSyncIQReplicationJobResource() resource.Resource {
	return &SyncIQReplicationJobResource{}
}

// SyncIQReplicationJobResource is the resource implementation.
type SyncIQReplicationJobResource struct {
	client *client.Client
}

// Configure implements resource.ResourceWithConfigure.
func (s *SyncIQReplicationJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	s.client = pscaleClient
}

// Metadata returns the resource type name.
func (r *SyncIQReplicationJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_replication_job"
}

// Schema defines the schema for the resource.
func (r *SyncIQReplicationJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing SyncIQ ReplicationJob on PowerScale. This resource can be used to manually trigger the replication job to replicate data from source powerscale cluster to a target powerscale cluster.",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "ID/Name of the policy",
				MarkdownDescription: "ID/Name of the policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"action": schema.StringAttribute{
				Required: true,
				Description: `Action for the job 
				 run - to start the replication job using synciq policy
				 test - to test the replication job using synciq policy, 
				 resync_prep - Resync_prep is a preparation step in PowerScale SyncIQ replication jobs that helps ensure a successful replication operation by performing a series of checks and verifications on the source and target volumes before starting the replication process., 
				 allow_write - allow_write determines whether the replication job allows writes to the target volume during the replication process. When configured, the target volume is writable, and any changes made to the target volume will be replicated to the source volume. This is useful in scenarios where you need to make changes to the target volume, such as updating files or creating new files, while the replication job is running.,
				 allow_write_revert - allow_write_revert determines whether the replication job allows writes to the target volume when reverting a replication job. When configure, the target volume is writable during the revert process, allowing changes made to the target volume during the revert process to be replicated to the source volume.`,
				MarkdownDescription: `Action for the job 
				 run - to start the replication job using synciq policy
				 test - to test the replication job using synciq policy, 
				 resync_prep - Resync_prep is a preparation step in PowerScale SyncIQ replication jobs that helps ensure a successful replication operation by performing a series of checks and verifications on the source and target volumes before starting the replication process., 
				 allow_write - allow_write determines whether the replication job allows writes to the target volume during the replication process. When configured, the target volume is writable, and any changes made to the target volume will be replicated to the source volume. This is useful in scenarios where you need to make changes to the target volume, such as updating files or creating new files, while the replication job is running.,
				 allow_write_revert - allow_write_revert determines whether the replication job allows writes to the target volume when reverting a replication job. When configure, the target volume is writable during the revert process, allowing changes made to the target volume during the revert process to be replicated to the source volume.`,
				Validators: []validator.String{
					stringvalidator.OneOf("run", "test", "resync_prep", "allow_write", "allow_write_revert"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"is_paused": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "change job state to running or paused.",
				MarkdownDescription: "change job state to running or paused.",
				Default:             booldefault.StaticBool(false),
			},
			"wait_time": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Wait Time for the job",
				MarkdownDescription: "Wait Time for the job",
				Default:             int64default.StaticInt64(5),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *SyncIQReplicationJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource create : Started")
	//Get Plan Data
	var plan models.SyncIQReplicationJobResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.IsPaused.ValueBool() {
		resp.Diagnostics.AddError("Config Error", "SyncIQ Replication Job cannot be paused befor job creation.")
	}

	var createJob powerscale.V1SyncJob
	// Get param from tf input
	err := helper.ReadFromState(ctx, &plan, &createJob)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading create plan",
			err.Error(),
		)
		return
	}
	_, err = helper.CreateSyncIQReplicationJob(ctx, r.client, createJob)
	if err != nil {
		errStr := "Could not create syncIQ Replication Job with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating syncIQ Replication Job",
			message,
		)
		return
	}
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource create: updating state finished, saving ...")
	// Save into State
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource create: finish")
}

// Read refreshes the Terraform state with the latest data.
func (r *SyncIQReplicationJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource read: started")
	var state models.SyncIQReplicationJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	time.Sleep(time.Duration(state.WaitTime.ValueInt64()) * time.Second)
	tflog.Debug(ctx, "calling get syncIQ Replication Job on powerscale client")
	readState, httpResp, err := helper.GetSyncIQReplicationJob(ctx, r.client, state.Id.ValueString())
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			diags = resp.State.Set(ctx, &state)
			resp.Diagnostics.Append(diags...)
			tflog.Trace(ctx, "resource_SyncIQReplicationJobResource read: finished")
			return
		}
		errStr := "Could not read syncIQ Replication Job with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading syncIQ Replication Job",
			message,
		)
		return
	}
	if len(readState.Jobs) > 0 {
		job := readState.Jobs[0]
		state.Id = types.StringValue(job.PolicyName)
		if job.State == "running" {
			state.IsPaused = types.BoolValue(false)
		} else if job.State == "paused" {
			state.IsPaused = types.BoolValue(true)
		}
	}

	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource read: finished reading state")
	//Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource read: finished")
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SyncIQReplicationJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource update: started")
	var state, plan models.SyncIQReplicationJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.IsPaused.Equal(state.IsPaused) {
		isPause := running
		if plan.IsPaused.ValueBool() {
			isPause = paused
		}
		updateJob := powerscale.V1SyncJobExtendedExtended{
			State: isPause,
		}
		_, err := helper.UpdateSyncIQReplicationJob(ctx, r.client, state.Id.ValueString(), updateJob)
		if err != nil {
			errStr := "Could not update syncIQ Replication Job with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error updating syncIQ Replication Job",
				message,
			)
			return
		}
		diags = resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
	}
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource update: finished")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SyncIQReplicationJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource delete: started")
	var state models.SyncIQReplicationJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := helper.DeleteSyncIQReplicationJob(ctx, r.client, state.Id.ValueString())
	if err != nil {
		errStr := "Could not delete syncIQ Replication Job with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting syncIQ Replication Job",
			message,
		)
	}
	time.Sleep(time.Duration(state.WaitTime.ValueInt64()) * time.Second)
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_SyncIQReplicationJobResource delete: finished")
}
