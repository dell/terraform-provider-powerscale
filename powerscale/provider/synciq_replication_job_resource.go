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

// NewSyncIQReplicationJobResource is a helper function to simplify the provider implementation.
func NewSyncIQReplicationJobResource() resource.Resource {
	return &SyncIQReplicationJobResource{}
}

// SyncIQReplicationJobResource is the resource implementation.
type SyncIQReplicationJobResource struct {
	client *client.Client
	isDelete bool
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
		MarkdownDescription: "Resource for managing SyncIQReplicationJobResource on OpenManage Enterprise.",
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
				Required:            true,
				Description:         "Action for the job",
				MarkdownDescription: "Action for the job",
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
		// need to discuss with team if we should remove resource from state or not
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			r.isDelete = true
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"SyncIQ Replication Job not found: Cleaning up state",
				"Use SyncIQ Reports to get latest sync status.",
			)
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
		isPause := "running"
		if plan.IsPaused.ValueBool() {
			isPause = "paused"
		}
		updateJob := powerscale.V1SyncJobExtendedExtended{
			State: isPause,
		}
		err := helper.UpdateSyncIQReplicationJob(ctx, r.client, state.Id.ValueString(), updateJob)
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
	if r.isDelete {
		return // already deleted
	}
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
