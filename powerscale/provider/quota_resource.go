/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"regexp"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &QuotaResource{}
	_ resource.ResourceWithImportState = &QuotaResource{}
)

// NewQuotaResource returns the Quota resource object.
func NewQuotaResource() resource.Resource {
	return &QuotaResource{}
}

// QuotaResource defines the resource implementation.
type QuotaResource struct {
	client *client.Client
}

// Configure configures the resource.
func (r *QuotaResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Metadata describes the resource arguments.
func (r *QuotaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quota"
}

// Schema describes the resource arguments.
func (r *QuotaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the Quota entity of PowerScale Array. " +
			"Quota module monitors and enforces administrator-defined storage limits. " +
			"We can Create, Update and Delete the Quota using this resource. We can also import an existing Quota from PowerScale array.",
		Description: "This resource is used to manage the Quota entity of PowerScale Array. " +
			"Quota module monitors and enforces administrator-defined storage limits. " +
			"We can Create, Update and Delete the Quota using this resource. We can also import an existing Quota from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			// Read-only attributes
			"id": schema.StringAttribute{
				Description:         "The system ID given to the quota.",
				MarkdownDescription: "The system ID given to the quota.",
				Computed:            true,
			},
			"efficiency_ratio": schema.NumberAttribute{
				Description:         "Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.",
				MarkdownDescription: "Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.",
				Computed:            true,
			},
			"notifications": schema.StringAttribute{
				Description:         "Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.",
				MarkdownDescription: "Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.",
				Computed:            true,
			},
			"ready": schema.BoolAttribute{
				Description:         "True if the default resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
				MarkdownDescription: "True if the default resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
				Computed:            true,
			},
			"reduction_ratio": schema.NumberAttribute{
				Description:         "Represents the ratio of logical space provided to physical data space used. This accounts for compression and data deduplication effects.",
				MarkdownDescription: "Represents the ratio of logical space provided to physical data space used. This accounts for compression and data deduplication effects.",
				Computed:            true,
			},
			"usage": schema.SingleNestedAttribute{
				Description:         "The usage of quota",
				MarkdownDescription: "The usage of quota",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"applogical": schema.Int64Attribute{
						Description:         "Bytes used by governed data apparent to application.",
						MarkdownDescription: "Bytes used by governed data apparent to application.",
						Computed:            true,
					},
					"applogical_ready": schema.BoolAttribute{
						Description:         "True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"fslogical": schema.Int64Attribute{
						Description:         "Bytes used by governed data apparent to filesystem.",
						MarkdownDescription: "Bytes used by governed data apparent to filesystem.",
						Computed:            true,
					},
					"fslogical_ready": schema.BoolAttribute{
						Description:         "True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"fsphysical": schema.Int64Attribute{
						Description:         "Physical data usage adjusted to account for shadow store efficiency",
						MarkdownDescription: "Physical data usage adjusted to account for shadow store efficiency",
						Computed:            true,
					},
					"fsphysical_ready": schema.BoolAttribute{
						Description:         "True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"inodes": schema.Int64Attribute{
						Description:         "Number of inodes (filesystem entities) used by governed data.",
						MarkdownDescription: "Number of inodes (filesystem entities) used by governed data.",
						Computed:            true,
					},
					"inodes_ready": schema.BoolAttribute{
						Description:         "True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"physical": schema.Int64Attribute{
						Description:         "Bytes used for governed data and filesystem overhead.",
						MarkdownDescription: "Bytes used for governed data and filesystem overhead.",
						Computed:            true,
					},
					"physical_data": schema.Int64Attribute{
						Description:         "Number of physical blocks for file data",
						MarkdownDescription: "Number of physical blocks for file data",
						Computed:            true,
					},
					"physical_data_ready": schema.BoolAttribute{
						Description:         "True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"physical_protection": schema.Int64Attribute{
						Description:         "Number of physical blocks for file protection",
						MarkdownDescription: "Number of physical blocks for file protection",
						Computed:            true,
					},
					"physical_protection_ready": schema.BoolAttribute{
						Description:         "True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"physical_ready": schema.BoolAttribute{
						Description:         "True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
					"shadow_refs": schema.Int64Attribute{
						Description:         "Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.",
						MarkdownDescription: "Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.",
						Computed:            true,
					},
					"shadow_refs_ready": schema.BoolAttribute{
						Description:         "True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						MarkdownDescription: "True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.",
						Computed:            true,
					},
				},
			},

			// Required and cannot be updated
			"path": schema.StringAttribute{
				Description:         "The ifs path governed.",
				MarkdownDescription: "The ifs path governed.",
				Required:            true,
				Validators: []validator.String{stringvalidator.RegexMatches(
					regexp.MustCompile(`^/ifs$|^/ifs/`), "must begin with /ifs",
				)},
			},
			"type": schema.StringAttribute{
				Description:         "The type of quota.",
				MarkdownDescription: "The type of quota.",
				Required:            true,
			},
			"include_snapshots": schema.BoolAttribute{
				Description:         "If true, quota governs snapshot data as well as head data.",
				MarkdownDescription: "If true, quota governs snapshot data as well as head data.",
				Required:            true,
			},

			// Optional and cannot be updated
			"zone": schema.StringAttribute{
				Description:         "Optional named zone to use for user and group resolution.",
				MarkdownDescription: "Optional named zone to use for user and group resolution.",
				Optional:            true,
			},
			"persona": schema.SingleNestedAttribute{
				Description:         "Specifies the persona of the file group. persona is required for user and group type.",
				MarkdownDescription: "Specifies the persona of the file group. persona is required for user and group type.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
						MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
						Optional:            true,
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Specifies the persona name, which must be combined with a type.",
						MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
						Optional:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Specifies the type of persona, which must be combined with a name.",
						MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
						Optional:            true,
					},
				},
			},

			// Optional after creation
			"linked": schema.BoolAttribute{
				Description:         "For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked. Computed by PowerScale, do not set Linked while creating.",
				MarkdownDescription: "For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked. Computed by PowerScale, do not set Linked while creating.",
				Optional:            true,
				Computed:            true,
			},

			// Optional, not computed
			"ignore_limit_checks": schema.BoolAttribute{
				Description:         "If true, skip child quota's threshold comparison with parent quota path.",
				MarkdownDescription: "If true, skip child quota's threshold comparison with parent quota path.",
				Optional:            true,
			},
			"force": schema.BoolAttribute{
				Description:         "Force creation of quotas on the root of /ifs or percent based quotas.",
				MarkdownDescription: "Force creation of quotas on the root of /ifs or percent based quotas.",
				Optional:            true,
			},

			// Computed and can be updated
			"container": schema.BoolAttribute{
				Description:         "If true, quotas using the quota directory see the quota thresholds as share size.",
				MarkdownDescription: "If true, quotas using the quota directory see the quota thresholds as share size.",
				Optional:            true,
				Computed:            true,
			},
			"enforced": schema.BoolAttribute{
				Description:         "True if the quota provides enforcement, otherwise an accounting quota.",
				MarkdownDescription: "True if the quota provides enforcement, otherwise an accounting quota.",
				Optional:            true,
				Computed:            true,
			},
			"thresholds": schema.SingleNestedAttribute{
				Description:         "The thresholds of quota",
				MarkdownDescription: "The thresholds of quota",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"advisory": schema.Int64Attribute{
						Description:         "Usage bytes at which notifications will be sent but writes will not be denied.",
						MarkdownDescription: "Usage bytes at which notifications will be sent but writes will not be denied.",
						Optional:            true,
						Computed:            true,
					},
					"advisory_exceeded": schema.BoolAttribute{
						Description:         "True if the advisory threshold has been hit.",
						MarkdownDescription: "True if the advisory threshold has been hit.",
						Computed:            true,
					},
					"advisory_last_exceeded": schema.Int64Attribute{
						Description:         "Time at which advisory threshold was hit.",
						MarkdownDescription: "Time at which advisory threshold was hit.",
						Computed:            true,
					},
					"hard": schema.Int64Attribute{
						Description:         "Usage bytes at which further writes will be denied.",
						MarkdownDescription: "Usage bytes at which further writes will be denied.",
						Optional:            true,
						Computed:            true,
					},
					"hard_exceeded": schema.BoolAttribute{
						Description:         "True if the hard threshold has been hit.",
						MarkdownDescription: "True if the hard threshold has been hit.",
						Computed:            true,
					},
					"hard_last_exceeded": schema.Int64Attribute{
						Description:         "Time at which hard threshold was hit.",
						MarkdownDescription: "Time at which hard threshold was hit.",
						Computed:            true,
					},
					"percent_advisory": schema.NumberAttribute{
						Description:         "Advisory threshold as percent of hard threshold. Usage bytes at which notifications will be sent but writes will not be denied. Must be >= 0.01 <= 99.99, precision 2",
						MarkdownDescription: "Advisory threshold as percent of hard threshold. Usage bytes at which notifications will be sent but writes will not be denied. Must be >= 0.01 <= 99.99, precision 2",
						Optional:            true,
						Computed:            true,
					},
					"percent_soft": schema.NumberAttribute{
						Description:         "Soft threshold as percent of hard threshold. Usage bytes at which notifications will be sent and soft grace time will be started. Must be >= 0.01 <= 99.99, precision 2",
						MarkdownDescription: "Soft threshold as percent of hard threshold. Usage bytes at which notifications will be sent and soft grace time will be started. Must be >= 0.01 <= 99.99, precision 2",
						Optional:            true,
						Computed:            true,
					},
					"soft": schema.Int64Attribute{
						Description:         "Usage bytes at which notifications will be sent and soft grace time will be started.",
						MarkdownDescription: "Usage bytes at which notifications will be sent and soft grace time will be started.",
						Optional:            true,
						Computed:            true,
					},
					"soft_exceeded": schema.BoolAttribute{
						Description:         "True if the soft threshold has been hit.",
						MarkdownDescription: "True if the soft threshold has been hit.",
						Computed:            true,
					},
					"soft_grace": schema.Int64Attribute{
						Description:         "Time in seconds after which the soft threshold has been hit before writes will be denied.",
						MarkdownDescription: "Time in seconds after which the soft threshold has been hit before writes will be denied.",
						Optional:            true,
						Computed:            true,
					},
					"soft_last_exceeded": schema.Int64Attribute{
						Description:         "Time at which soft threshold was hit",
						MarkdownDescription: "Time at which soft threshold was hit",
						Computed:            true,
					},
				},
			},
			"thresholds_on": schema.StringAttribute{
				Description:         "Thresholds apply on quota accounting metric.",
				MarkdownDescription: "Thresholds apply on quota accounting metric.",
				Optional:            true,
				Computed:            true,
			},
		},
	}

}

// Create allocates the resource.
func (r *QuotaResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "Creating quota")

	var quotaPlan models.QuotaResource
	diags := request.Plan.Get(ctx, &quotaPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	var quotaPlanBackup models.QuotaResource
	diags = request.Plan.Get(ctx, &quotaPlanBackup)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if !quotaPlan.Linked.IsUnknown() {
		response.Diagnostics.AddError("Error creating quota", "Do not set attribute Linked while creating")
		return
	}
	if err := helper.IsQuotaParamInvalid(quotaPlan); err != nil {
		response.Diagnostics.AddError("Error creating quota", err.Error())
		return
	}
	quotaToCreate := powerscale.V12QuotaQuota{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, quotaPlan, &quotaToCreate)
	if err != nil {
		response.Diagnostics.AddError("Error creating quota",
			fmt.Sprintf("Could not read quota param of Path: %s with error: %s", quotaPlan.Path.ValueString(), err.Error()),
		)
		return
	}
	quotaID, err := helper.CreateQuota(ctx, r.client, quotaToCreate, quotaPlan.Zone.ValueString())
	if err != nil {
		errStr := constants.CreateQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating quota ",
			message)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("quota %s created", quotaID.Id), map[string]interface{}{
		"QuotaResponse": quotaID,
	})

	getQuotaResponse, err := helper.GetQuota(ctx, r.client, quotaID.Id, quotaPlan.Zone.ValueString())
	if err != nil {
		errStr := constants.ReadQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating quota ", message)
		return
	}

	// update resource state according to response
	if len(getQuotaResponse.Quotas) <= 0 {
		response.Diagnostics.AddError(
			"Error creating quota",
			fmt.Sprintf("Could not get created quota state %s with error: quota not found", quotaID),
		)
		return
	}
	createdQuota := getQuotaResponse.Quotas[0]

	err = helper.CopyFieldsToNonNestedModel(ctx, createdQuota, &quotaPlan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating quota",
			fmt.Sprintf("Could not read quota %s with error: %s", quotaID, err.Error()),
		)
		return
	}

	quotaPlan.IgnoreLimitChecks = quotaPlanBackup.IgnoreLimitChecks
	quotaPlan.Force = quotaPlanBackup.Force
	if quotaPlanBackup.Persona.IsNull() || createdQuota.Type == "directory" {
		quotaPlan.Persona = types.ObjectNull(quotaPlan.Persona.AttributeTypes(ctx))
	}
	diags = response.State.Set(ctx, quotaPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create quota completed")
}

// Update updates the resource state.
func (r *QuotaResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "updating quota")
	var quotaPlan models.QuotaResource
	diags := request.Plan.Get(ctx, &quotaPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var quotaState models.QuotaResource
	diags = response.State.Get(ctx, &quotaState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update quota", map[string]interface{}{
		"quotaPlan":  quotaPlan,
		"quotaState": quotaState,
	})

	quotaID := quotaState.ID.ValueString()

	// validate update params
	if err := helper.ValidateQuotaUpdate(quotaPlan, quotaState); err != nil {
		response.Diagnostics.AddError(
			"Error updating quota",
			fmt.Sprintf("Could not update Quota %s with error: %s", quotaID, err.Error()),
		)
		return
	}

	var quotaToUpdate powerscale.V12QuotaQuotaExtendedExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, quotaPlan, &quotaToUpdate)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating quota",
			fmt.Sprintf("Could not read quota struct %s with error: %s", quotaID, err.Error()),
		)
		return
	}
	err = helper.UpdateQuota(ctx, r.client, quotaID, quotaToUpdate)
	if err != nil {
		errStr := constants.UpdateQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating quota ",
			message)
		return
	}
	quotaID = quotaState.ID.ValueString()
	tflog.Debug(ctx, "calling get quota by ID on pscale client", map[string]interface{}{
		"smbQuotaID": quotaID,
	})
	updatedQuota, err := helper.GetQuota(ctx, r.client, quotaID, quotaPlan.Zone.ValueString())
	if err != nil {
		errStr := constants.UpdateQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating quota ",
			message)
		return
	}

	if len(updatedQuota.Quotas) <= 0 {
		response.Diagnostics.AddError(
			"Error reading quota",
			fmt.Sprintf("Could not read quota %s from pscale with error: quota not found", quotaID),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, updatedQuota.Quotas[0], &quotaState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading quota",
			fmt.Sprintf("Could not read quota struct %s with error: %s", quotaID, err.Error()),
		)
		return
	}
	quotaState.IgnoreLimitChecks = quotaPlan.IgnoreLimitChecks
	quotaState.Force = quotaPlan.Force
	if quotaPlan.Zone.ValueString() != "" {
		quotaState.Zone = quotaPlan.Zone
	}
	if quotaPlan.Persona.IsNull() || quotaPlan.Type.ValueString() == "directory" {
		quotaState.Persona = types.ObjectNull(quotaState.Persona.AttributeTypes(ctx))
	}
	diags = response.State.Set(ctx, quotaState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update quota completed")
}

// Read reads data from the resource.
func (r *QuotaResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Quota resource")
	var quotaState models.QuotaResource
	diags := request.State.Get(ctx, &quotaState)
	response.Diagnostics.Append(diags...)
	var quotaStateBackup models.QuotaResource
	diags = request.State.Get(ctx, &quotaStateBackup)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	quotaID := quotaState.ID.ValueString()
	zone := quotaState.Zone.ValueString()
	tflog.Debug(ctx, "calling get quota by ID", map[string]interface{}{
		"QuotaID": quotaID,
		"Zone":    zone,
	})
	quotaResponse, err := helper.GetQuota(ctx, r.client, quotaID, zone)
	if err != nil {
		errStr := constants.ReadQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading quota ",
			message)
		return
	}

	if len(quotaResponse.Quotas) <= 0 {
		response.Diagnostics.AddError(
			"Error reading quota",
			fmt.Sprintf("Could not read quota %s from pscale with error: quota not found", quotaID),
		)
		return
	}
	tflog.Debug(ctx, "updating read quota state", map[string]interface{}{
		"QuotaResponse": quotaResponse,
		"QuotaState":    quotaState,
	})
	err = helper.CopyFieldsToNonNestedModel(ctx, quotaResponse.Quotas[0], &quotaState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error read quota",
			fmt.Sprintf("Could not read quota struct %s with error: %s", quotaID, err.Error()),
		)
		return
	}
	quotaState.IgnoreLimitChecks = quotaStateBackup.IgnoreLimitChecks
	quotaState.Force = quotaStateBackup.Force
	if quotaStateBackup.Persona.IsNull() || quotaStateBackup.Type.ValueString() == "directory" {
		quotaState.Persona = types.ObjectNull(quotaState.Persona.AttributeTypes(ctx))
	}
	diags = response.State.Set(ctx, quotaState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read quota completed")
}

// Delete deletes the resource.
func (r QuotaResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting quota")
	var quotaState models.QuotaResource
	diags := request.State.Get(ctx, &quotaState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	quotaID := quotaState.ID.ValueString()
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete quota on pscale client", map[string]interface{}{
		"QuotaID": quotaID,
	})
	err := helper.DeleteQuota(ctx, r.client, quotaID)
	if err != nil {
		errStr := constants.DeleteQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error deleting quota ",
			message)
		return
	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete quota completed")
}

// ImportState imports the resource state.
func (r QuotaResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing Quota resource")
	var zoneName string
	quotaID := request.ID
	// request.ID is form of zoneName:exportID
	if strings.Contains(request.ID, ":") {
		params := strings.Split(request.ID, ":")
		quotaID = strings.Trim(params[1], " ")
		zoneName = strings.Trim(params[0], " ")
	}

	tflog.Debug(ctx, "calling get quota by ID", map[string]interface{}{
		"QuotaID": quotaID,
		"Zone":    zoneName,
	})
	quotaResponse, err := helper.GetQuota(ctx, r.client, quotaID, zoneName)
	if err != nil {
		errStr := constants.ReadQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error importing quota ",
			message)
		return
	}

	if len(quotaResponse.Quotas) <= 0 {
		response.Diagnostics.AddError(
			"Error importing quota",
			fmt.Sprintf("Could not read quota %s from pscale with error: quota not found", quotaID),
		)
		return
	}
	var quotaState models.QuotaResource
	tflog.Debug(ctx, "updating read quota state", map[string]interface{}{
		"QuotaResponse": quotaResponse,
		"QuotaState":    quotaState,
	})
	err = helper.CopyFieldsToNonNestedModel(ctx, quotaResponse.Quotas[0], &quotaState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error importing quota",
			fmt.Sprintf("Could not read quota struct %s with error: %s", quotaID, err.Error()),
		)
		return
	}
	if quotaState.Type.ValueString() == "directory" {
		quotaState.Persona = types.ObjectNull(quotaState.Persona.AttributeTypes(ctx))
	}
	quotaState.Zone = types.StringNull()
	if zoneName != "" {
		quotaState.Zone = types.StringValue(zoneName)
	}
	diags := response.State.Set(ctx, quotaState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "importing quota completed")
}
