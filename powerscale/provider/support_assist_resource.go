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
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &SupportAssistResource{}
)

// NewSupportAssistResource returns the Support Assist resource object.
func NewSupportAssistResource() resource.Resource {
	return &SupportAssistResource{}
}

// SupportAssistResource defines the resource implementation.
type SupportAssistResource struct {
	client *client.Client
}

// Configure configures the resource.
func (r *SupportAssistResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (r *SupportAssistResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_support_assist"
}

func (r *SupportAssistResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("supportassist_enabled"),
			path.MatchRoot("enable_download"),
			path.MatchRoot("contact"),
			path.MatchRoot("telemetry"),
			path.MatchRoot("automatic_case_creation"),
			path.MatchRoot("connections"),
			path.MatchRoot("enable_remote_support"),
			path.MatchRoot("accepted_terms"),
			path.MatchRoot("access_key"),
			path.MatchRoot("pin"),
		),
		resourcevalidator.RequiredTogether(
			path.MatchRoot("access_key"),
			path.MatchRoot("pin"),
		),
	}
}

func (r *SupportAssistResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the Support Assist settings of PowerScale Array. We can Create, Update and Delete the Support Assist settings using this resource. Note that, Support Assist settings is the native functionality of PowerScale.",
		Description:         "This resource is used to manage the Support Assist settings of PowerScale Array. We can Create, Update and Delete the Support Assist settings using this resource. Note that, Support Assist settings is the native functionality of PowerScale.",
		Attributes:          SupportAssistResourceSchema(),
	}
}

func SupportAssistResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder ID",
			MarkdownDescription: "Placeholder ID",
			Computed:            true,
		},
		"supportassist_enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:         "Whether SupportAssist is enabled",
			MarkdownDescription: "Whether SupportAssist is enabled",
		},
		"enable_download": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:         "True indicates downloads are enabled",
			MarkdownDescription: "True indicates downloads are enabled",
		},
		"contact": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"primary": schema.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"phone": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's phone number.",
							MarkdownDescription: "Contact's phone number.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("([\\.\\-\\+\\/\\sxX]*([0-9]+|[\\(\\d+\\)])+)+"), "must be a valid phone number"),
								stringvalidator.LengthBetween(0, 40),
							},
						},
						"email": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's email address.",
							MarkdownDescription: "Contact's email address.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9._%-]+@([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9]+$"), "must be a valid email"),
								stringvalidator.LengthBetween(0, 320),
							},
						},
						"first_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's first name.",
							MarkdownDescription: "Contact's first name.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("[\\p{L}\\p{M}*\\-\\.\\' ]*"), "must be a valid first name"),
								stringvalidator.LengthBetween(0, 50),
							},
						},
						"language": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"Cs",
									"Da",
									"De",
									"El",
									"En",
									"Es",
									"es-LA",
									"Fi",
									"fr-CA",
									"He",
									"It",
									"Ja",
									"Ko",
									"Nl",
									"No",
									"Pl",
									"Pt",
									"pt-BR",
									"Ru",
									"Sk",
									"Sv",
									"Th",
									"Tr",
									"zh-CN",
									"zh-TW",
								),
								stringvalidator.LengthBetween(0, 255),
							},
						},
						"last_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's last name.",
							MarkdownDescription: "Contact's last name.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("[\\p{L}\\p{M}*\\-\\.\\' ]*"), "must be a valid last name"),
								stringvalidator.LengthBetween(0, 50),
							},
						},
					},
				},
				"secondary": schema.SingleNestedAttribute{
					Optional: true,
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"phone": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's phone number.",
							MarkdownDescription: "Contact's phone number.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("([\\.\\-\\+\\/\\sxX]*([0-9]+|[\\(\\d+\\)])+)+"), "must be a valid phone number"),
								stringvalidator.LengthBetween(0, 40),
							},
						},
						"email": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's email address.",
							MarkdownDescription: "Contact's email address.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9._%-]+@([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9]+$"), "must be a valid email"),
								stringvalidator.LengthBetween(0, 320),
							},
						},
						"first_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's first name.",
							MarkdownDescription: "Contact's first name.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("[\\p{L}\\p{M}*\\-\\.\\' ]*"), "must be a valid first name"),
								stringvalidator.LengthBetween(0, 50),
							},
						},
						"language": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"Cs",
									"Da",
									"De",
									"El",
									"En",
									"Es",
									"es-LA",
									"Fi",
									"fr-CA",
									"He",
									"It",
									"Ja",
									"Ko",
									"Nl",
									"No",
									"Pl",
									"Pt",
									"pt-BR",
									"Ru",
									"Sk",
									"Sv",
									"Th",
									"Tr",
									"zh-CN",
									"zh-TW",
								),
								stringvalidator.LengthBetween(0, 255),
							},
						},
						"last_name": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Description:         "Contact's last name.",
							MarkdownDescription: "Contact's last name.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("[\\p{L}\\p{M}*\\-\\.\\' ]*"), "must be a valid last name"),
								stringvalidator.LengthBetween(0, 50),
							},
						},
					},
				},
			},
		},
		"telemetry": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"telemetry_threads": schema.Int64Attribute{
					Optional:            true,
					Computed:            true,
					Description:         "Change the number of threads for telemetry gathers",
					MarkdownDescription: "Change the number of threads for telemetry gathers",
				},
				"offline_collection_period": schema.Int64Attribute{
					Optional:            true,
					Computed:            true,
					Description:         "Change the offline collection period for when the connection to gateway is down",
					MarkdownDescription: "Change the offline collection period for when the connection to gateway is down",
				},
				"telemetry_enabled": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Change the status of telemetry",
					MarkdownDescription: "Change the status of telemetry",
				},
				"telemetry_persist": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Change if files are kept after upload",
					MarkdownDescription: "Change if files are kept after upload",
				},
			},
		},
		"automatic_case_creation": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:         "True indicates automatic case creation is enabled",
			MarkdownDescription: "True indicates automatic case creation is enabled",
		},
		"connections": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Connection Mode for SupportAssist: can be direct or via gateway",
					MarkdownDescription: "Connection Mode for SupportAssist: can be direct or via gateway",
					Validators: []validator.String{
						stringvalidator.OneOf("direct", "gateway"),
					},
				},
				"network_pools": schema.ListAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Network pools for gateway use",
					MarkdownDescription: "Network pools for gateway use",
					ElementType:         types.StringType,
				},
				"gateway_endpoints": schema.ListNestedAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Gateway details",
					MarkdownDescription: "Gateway details",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"port": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								Description:         "Gateway port",
								MarkdownDescription: "Gateway port",
							},
							"validate_ssl": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Whether to validate SSL for this gateway",
								MarkdownDescription: "Whether to validate SSL for this gateway",
							},
							"priority": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								Description:         "Gateway's priority",
								MarkdownDescription: "Gateway's priority",
							},
							"use_proxy": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Whether to use Proxy for this gateway",
								MarkdownDescription: "Whether to use Proxy for this gateway",
							},
							"host": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Gateway hostname or IPv4 address",
								MarkdownDescription: "Gateway hostname or IPv4 address",
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile("(^$|^((([a-zA-Z0-9_][a-zA-Z0-9-]{0,61})?[a-zA-Z0-9])(\\.([a-zA-Z0-9_][a-zA-Z0-9-]{0,61})?[a-zA-Z0-9])*)$|^([01]?[0-9]?[0-9]|2[0-4][0-9]|25[0-5])(\\.([01]?[0-9]?[0-9]|2[0-4][0-9]|25[0-5])){3}$)"), "must be a valid hostname or IPv4 address"),
									stringvalidator.LengthBetween(0, 255),
								},
							},
							"enabled": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								Description:         "Whether this gateway is enabled/disabled",
								MarkdownDescription: "Whether this gateway is enabled/disabled",
							},
						},
					},
				},
			},
		},
		"enable_remote_support": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:         "Allow remote support.",
			MarkdownDescription: "Allow remote support.",
		},
		"accepted_terms": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			Description:         "Set T&C accepted or rejected status",
			MarkdownDescription: "Set T&C accepted or rejected status",
		},
		"access_key": schema.StringAttribute{
			Description:         "SupportAssist access key",
			MarkdownDescription: "SupportAssist access key",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"pin": schema.StringAttribute{
			Description:         "SupportAssist pin",
			MarkdownDescription: "SupportAssist pin",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}
}

func (r *SupportAssistResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "Creating support assist resource state")
	var plan models.SupportAssistModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageSupportAssist(ctx, r.client, plan)
	response.Diagnostics.Append(diags...)

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with creating support assist resource state")
}

func (r *SupportAssistResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "Reading support assist resource state")
	var state models.SupportAssistModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	state, dig := helper.ReadSupportAssistDetails(ctx, r.client, state)
	response.Diagnostics.Append(dig...)

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading support assist resource state")

}

func (r *SupportAssistResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating support assist resource state")
	var plan models.SupportAssistModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var state models.SupportAssistModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	state, diags := helper.ManageSupportAssist(ctx, r.client, plan)
	response.Diagnostics.Append(diags...)
	// if response.Diagnostics.HasError() {
	// 	return
	// }

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with updating support assist resource state")
}

func (r SupportAssistResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting support assist resource state")
	var state models.SupportAssistModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "Done with deleting support assist resource state")
}

func (r SupportAssistResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
