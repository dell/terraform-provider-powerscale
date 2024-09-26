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

package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SyncIQRuleDataSourceSchema defines the schema for the data source.
func SyncIQRuleDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing SyncIQ Replication Rules from PowerScale array." +
			" The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Description: "This datasource is used to query the existing SyncIQ Replication Rules from PowerScale array." +
			" The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unique identifier of the performance rule.",
				MarkdownDescription: "Unique identifier of the performance rule.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "The type of system resource this rule limits.",
							MarkdownDescription: "The type of system resource this rule limits.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"bandwidth",
									"file_count",
									"cpu",
									"worker",
								),
							},
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "User-entered description of this performance rule.",
							MarkdownDescription: "User-entered description of this performance rule.",
						},
						"enabled": schema.BoolAttribute{
							Computed:            true,
							Description:         "Whether this performance rule is currently in effect during its specified intervals.",
							MarkdownDescription: "Whether this performance rule is currently in effect during its specified intervals.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "The system ID given to this performance rule.",
							MarkdownDescription: "The system ID given to this performance rule.",
						},
						"limit": schema.Int64Attribute{
							Computed:            true,
							Description:         "Amount the specified system resource type is limited by this rule.  Units are kb/s for bandwidth, files/s for file-count, processing percentage used for cpu, or percentage of maximum available workers.",
							MarkdownDescription: "Amount the specified system resource type is limited by this rule.  Units are kb/s for bandwidth, files/s for file-count, processing percentage used for cpu, or percentage of maximum available workers.",
						},
						"schedule": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "A schedule defining when during a week this performance rule is in effect.  If unspecified or null, the schedule will always be in effect.",
							MarkdownDescription: "A schedule defining when during a week this performance rule is in effect.  If unspecified or null, the schedule will always be in effect.",
							Attributes: map[string]schema.Attribute{
								"tuesday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Tuesday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Tuesday.  If false, or unspecified, it is not.",
								},
								"monday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Monday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Monday.  If false, or unspecified, it is not.",
								},
								"wednesday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Wednesday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Wednesday.  If false, or unspecified, it is not.",
								},
								"saturday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Saturday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Saturday.  If false, or unspecified, it is not.",
								},
								"friday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Friday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Friday.  If false, or unspecified, it is not.",
								},
								"end": schema.StringAttribute{
									Computed:            true,
									Description:         "End time (inclusive) for this schedule, during its specified days.  Format is \"hh:mm\" (three-letter weekday name abbreviation, 24h format hour, and minute).  A null value indicates the end of the day (\"23:59\").",
									MarkdownDescription: "End time (inclusive) for this schedule, during its specified days.  Format is \"hh:mm\" (three-letter weekday name abbreviation, 24h format hour, and minute).  A null value indicates the end of the day (\"23:59\").",
								},
								"sunday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Sunday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Sunday.  If false, or unspecified, it is not.",
								},
								"begin": schema.StringAttribute{
									Computed:            true,
									Description:         "Start time (inclusive) for this schedule, during its specified days.  Format is \"hh:mm\" (24h format hour, and minute).  A null value indicates the beginning of the day (\"00:00\").",
									MarkdownDescription: "Start time (inclusive) for this schedule, during its specified days.  Format is \"hh:mm\" (24h format hour, and minute).  A null value indicates the beginning of the day (\"00:00\").",
								},
								"thursday": schema.BoolAttribute{
									Computed:            true,
									Description:         "If true, this rule is in effect on Thursday.  If false, or unspecified, it is not.",
									MarkdownDescription: "If true, this rule is in effect on Thursday.  If false, or unspecified, it is not.",
								},
							},
						},
					},
				},
			},
		},
	}
}
