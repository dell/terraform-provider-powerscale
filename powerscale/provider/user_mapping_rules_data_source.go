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

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &UserMappingRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &UserMappingRulesDataSource{}
)

// NewUserMappingRulesDataSource creates a new network settings data source.
func NewUserMappingRulesDataSource() datasource.DataSource {
	return &UserMappingRulesDataSource{}
}

// UserMappingRulesDataSource defines the data source implementation.
type UserMappingRulesDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *UserMappingRulesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_mapping_rules"
}

// Schema describes the data source arguments.
func (d *UserMappingRulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the User Mapping Rules from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules.",
		Description:         "This datasource is used to query the User Mapping Rules from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale User Mapping Rules combines user identities from different directory services into a single access token and then modifies it according to configured rules.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "User Mapping Rules ID.",
				MarkdownDescription: "User Mapping Rules ID.",
				Computed:            true,
			},
			"user_mapping_rules_parameters": schema.SingleNestedAttribute{
				Description:         "Specifies the parameters for user mapping rules.",
				MarkdownDescription: "Specifies the parameters for user mapping rules.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"default_unix_user": schema.SingleNestedAttribute{
						Description:         "Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information.",
						MarkdownDescription: "Specifies the default UNIX user information that can be applied if the final credentials do not have valid UID and GID information.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"domain": schema.StringAttribute{
								Description:         "Specifies the domain of the user that is being mapped.",
								MarkdownDescription: "Specifies the domain of the user that is being mapped.",
								Computed:            true,
							},
							"user": schema.StringAttribute{
								Description:         "Specifies the name of the user that is being mapped.",
								MarkdownDescription: "Specifies the name of the user that is being mapped.",
								Computed:            true,
							},
						},
					},
				},
			},
			"user_mapping_rules": schema.ListNestedAttribute{
				Description:         "Specifies the list of user mapping rules.",
				MarkdownDescription: "Specifies the list of user mapping rules.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operator": schema.StringAttribute{
							Description:         "Specifies the operator to make rules on specified users or groups.",
							MarkdownDescription: "Specifies the operator to make rules on specified users or groups.",
							Computed:            true,
						},
						"options": schema.SingleNestedAttribute{
							Description:         "Specifies the mapping options for this user mapping rule.",
							MarkdownDescription: "Specifies the mapping options for this user mapping rule.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"break": schema.BoolAttribute{
									Description:         "If true, and the rule was applied successfully, stop processing further.",
									MarkdownDescription: "If true, and the rule was applied successfully, stop processing further.",
									Computed:            true,
								},
								"user": schema.BoolAttribute{
									Description:         "If true, the primary UID and primary user SID should be copied to the existing credential.",
									MarkdownDescription: "If true, the primary UID and primary user SID should be copied to the existing credential.",
									Computed:            true,
								},
								"group": schema.BoolAttribute{
									Description:         "If true, the primary GID and primary group SID should be copied to the existing credential.",
									MarkdownDescription: "If true, the primary GID and primary group SID should be copied to the existing credential.",
									Computed:            true,
								},
								"groups": schema.BoolAttribute{
									Description:         "If true, all additional identifiers should be copied to the existing credential.",
									MarkdownDescription: "If true, all additional identifiers should be copied to the existing credential.",
									Computed:            true,
								},
								"default_user": schema.SingleNestedAttribute{
									Description:         "Specifies the default user information that can be applied if the final credentials do not have valid UID and GID information.",
									MarkdownDescription: "Specifies the default user information that can be applied if the final credentials do not have valid UID and GID information.",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"domain": schema.StringAttribute{
											Description:         "Specifies the domain of the user that is being mapped.",
											MarkdownDescription: "Specifies the domain of the user that is being mapped.",
											Computed:            true,
										},
										"user": schema.StringAttribute{
											Description:         "Specifies the name of the user that is being mapped.",
											MarkdownDescription: "Specifies the name of the user that is being mapped.",
											Computed:            true,
										},
									},
								},
							},
						},
						"target_user": schema.SingleNestedAttribute{
							Description:         "Specifies the target user information that the rule can be applied to.",
							MarkdownDescription: "Specifies the target user information that the rule can be applied to.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"domain": schema.StringAttribute{
									Description:         "Specifies the domain of the user that is being mapped.",
									MarkdownDescription: "Specifies the domain of the user that is being mapped.",
									Computed:            true,
								},
								"user": schema.StringAttribute{
									Description:         "Specifies the name of the user that is being mapped.",
									MarkdownDescription: "Specifies the name of the user that is being mapped.",
									Computed:            true,
								},
							},
						},
						"source_user": schema.SingleNestedAttribute{
							Description:         "Specifies the source user information that the rule can be applied from.",
							MarkdownDescription: "Specifies the source user information that the rule can be applied from.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"domain": schema.StringAttribute{
									Description:         "Specifies the domain of the user that is being mapped.",
									MarkdownDescription: "Specifies the domain of the user that is being mapped.",
									Computed:            true,
								},
								"user": schema.StringAttribute{
									Description:         "Specifies the name of the user that is being mapped.",
									MarkdownDescription: "Specifies the name of the user that is being mapped.",
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
					"names": schema.SetAttribute{
						Description:         "Names filter for source user name or target user name.",
						MarkdownDescription: "Names filter for source user name or target user name.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"operators": schema.SetAttribute{
						Description:         "Operators filter for user mapping rules.",
						MarkdownDescription: "Operators filter for user mapping rules.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.OneOf("append", "insert", "replace", "trim", "union"))},
					},
					"zone": schema.StringAttribute{
						Description:         "The zone to which the user mapping applies. Defaults to System.",
						MarkdownDescription: "The zone to which the user mapping applies. Defaults to System.",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *UserMappingRulesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *UserMappingRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading User Mapping Rules data source ")

	var plan models.UserMappingRulesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var zone string
	if plan.Filter != nil {
		zone = plan.Filter.Zone.ValueString()
	}
	result, err := helper.GetUserMappingRulesByZone(ctx, d.client, zone)
	if err != nil {
		resp.Diagnostics.AddError("error getting PowerScale User Mapping Rules.", err.Error())
		return
	}

	if diags := helper.UpdateUserMappingRulesDatasourceState(ctx, &plan, result); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Read User Mapping Rules data source ")
}
