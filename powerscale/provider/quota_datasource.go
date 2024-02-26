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
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

var (
	_ datasource.DataSource              = &QuotaDataSource{}
	_ datasource.DataSourceWithConfigure = &QuotaDataSource{}
)

// NewQuotaDataSource returns the Quota data source object.
func NewQuotaDataSource() datasource.DataSource {
	return &QuotaDataSource{}
}

// QuotaDataSource defines the data source implementation.
type QuotaDataSource struct {
	client *client.Client
}

// Configure configures the resource.
func (d *QuotaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Metadata describes the data source arguments.
func (d *QuotaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quota"
}

// Schema describes the data source arguments.
func (d *QuotaDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing quotas from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"Quota module monitors and enforces administrator-defined storage limits",
		Description: "This datasource is used to query the existing quotas from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"Quota module monitors and enforces administrator-defined storage limits",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"quotas": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of Quotas",
				MarkdownDescription: "List of Quotas",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"container": schema.BoolAttribute{
							Description:         "If true, SMB shares using the quota directory see the quota thresholds as share size.",
							MarkdownDescription: "If true, SMB shares using the quota directory see the quota thresholds as share size.",
							Computed:            true,
						},
						"efficiency_ratio": schema.NumberAttribute{
							Description:         "Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.",
							MarkdownDescription: "Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.",
							Computed:            true,
						},
						"enforced": schema.BoolAttribute{
							Description:         "True if the quota provides enforcement, otherwise an accounting quota.",
							MarkdownDescription: "True if the quota provides enforcement, otherwise an accounting quota.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "The system ID given to the quota.",
							MarkdownDescription: "The system ID given to the quota.",
							Computed:            true,
						},
						"include_snapshots": schema.BoolAttribute{
							Description:         "If true, quota governs snapshot data as well as head data.",
							MarkdownDescription: "If true, quota governs snapshot data as well as head data.",
							Computed:            true,
						},
						"linked": schema.BoolAttribute{
							Description:         "For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked.",
							MarkdownDescription: "For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked.",
							Computed:            true,
						},
						"notifications": schema.StringAttribute{
							Description:         "Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.",
							MarkdownDescription: "Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							Description:         "The ifs path governed.",
							MarkdownDescription: "The ifs path governed.",
							Computed:            true,
						},
						"persona": schema.SingleNestedAttribute{
							Description:         "Specifies the persona of the file group.",
							MarkdownDescription: "Specifies the persona of the file group.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Computed:            true,
								},
							},
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
						"thresholds": schema.SingleNestedAttribute{
							Description:         "The thresholds of quota",
							MarkdownDescription: "The thresholds of quota",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"advisory": schema.Int64Attribute{
									Description:         "Usage bytes at which notifications will be sent but writes will not be denied.",
									MarkdownDescription: "Usage bytes at which notifications will be sent but writes will not be denied.",
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
									Description:         "Advisory threshold as percent of hard threshold. Usage bytes at which notifications will be sent but writes will not be denied.",
									MarkdownDescription: "Advisory threshold as percent of hard threshold. Usage bytes at which notifications will be sent but writes will not be denied.",
									Computed:            true,
								},
								"percent_soft": schema.NumberAttribute{
									Description:         "Soft threshold as percent of hard threshold. Usage bytes at which notifications will be sent and soft grace time will be started.",
									MarkdownDescription: "Soft threshold as percent of hard threshold. Usage bytes at which notifications will be sent and soft grace time will be started.",
									Computed:            true,
								},
								"soft": schema.Int64Attribute{
									Description:         "Usage bytes at which notifications will be sent and soft grace time will be started.",
									MarkdownDescription: "Usage bytes at which notifications will be sent and soft grace time will be started.",
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
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "The type of quota.",
							MarkdownDescription: "The type of quota.",
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
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"enforced": schema.BoolAttribute{
						Description:         "Only list quotas with this enforcement (non-accounting).",
						MarkdownDescription: "Only list quotas with this enforcement (non-accounting).",
						Optional:            true,
					},
					"exceeded": schema.BoolAttribute{
						Description:         "Set to true to only list quotas which have exceeded one or more of their thresholds.",
						MarkdownDescription: "Set to true to only list quotas which have exceeded one or more of their thresholds.",
						Optional:            true,
					},
					"include_snapshots": schema.BoolAttribute{
						Description:         "Only list quotas with this setting for include_snapshots.",
						MarkdownDescription: "Only list quotas with this setting for include_snapshots.",
						Optional:            true,
					},
					"path": schema.StringAttribute{
						Description:         "Only list quotas matching this path (see also recurse_path_*).",
						MarkdownDescription: "Only list quotas matching this path (see also recurse_path_*).",
						Optional:            true,
					},
					"persona": schema.StringAttribute{
						Description:         "Only list user or group quotas matching this persona (must be used with the corresponding type argument).",
						MarkdownDescription: "Only list user or group quotas matching this persona (must be used with the corresponding type argument).",
						Optional:            true,
					},
					"recurse_path_children": schema.BoolAttribute{
						Description:         "If used with the path argument, match all quotas at that path or any descendent sub-directory.",
						MarkdownDescription: "If used with the path argument, match all quotas at that path or any descendent sub-directory.",
						Optional:            true,
					},
					"recurse_path_parents": schema.BoolAttribute{
						Description:         "If used with the path argument, match all quotas at that path or any parent directory.",
						MarkdownDescription: "If used with the path argument, match all quotas at that path or any parent directory.",
						Optional:            true,
					},
					"report_id": schema.StringAttribute{
						Description:         "Use the named report as a source rather than the live quotas.",
						MarkdownDescription: "Use the named report as a source rather than the live quotas.",
						Optional:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Only list quotas matching this type.",
						MarkdownDescription: "Only list quotas matching this type.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("directory", "user", "group",
								"default-directory", "default-user", "default-group"),
						},
					},
					"zone": schema.StringAttribute{
						Description:         "Optional named zone to use for user and group resolution.",
						MarkdownDescription: "Optional named zone to use for user and group resolution.",
						Optional:            true,
					},
				},
			},
		},
	}

}

// Read reads data from the data source.
func (d *QuotaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Quota data source ")
	var plan models.QuotaDatasource
	var state models.QuotaDatasource
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	quotas, err := helper.ListQuotas(ctx, d.client, plan.QuotaFilter)
	var filteredQuotas []models.QuotaDatasourceEntity

	if err != nil {
		errStr := constants.ListQuotaErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of quotas",
			message,
		)
		return
	}

	// Convert from json to terraform model
	for _, quota := range quotas {
		entity := models.QuotaDatasourceEntity{}
		err := helper.CopyFields(ctx, quota, &entity)
		if err != nil {
			resp.Diagnostics.AddError("Error reading quota datasource plan",
				fmt.Sprintf("Could not list quotas with error: %s", err.Error()))
			return
		}
		filteredQuotas = append(filteredQuotas, entity)
	}

	state.ID = types.StringValue("quota_datasource")
	state.QuotaFilter = plan.QuotaFilter
	state.Quotas = filteredQuotas

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
