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

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SyncIQRuleResourceSchema defines the schema for the syncIQ Replication Rule resource.
func SyncIQRuleResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Description:         "This resource is used to manage the SyncIQ Rule entity on PowerScale array.",
		MarkdownDescription: "This resource is used to manage the SyncIQ Rule entity on PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required: true,
				Description: "The type of system resource this rule limits." +
					" Acceptable values are: 'bandwidth', 'file_count', 'cpu', 'worker'.",
				MarkdownDescription: "The type of system resource this rule limits." +
					" Acceptable values are: `bandwidth`, `file_count`, `cpu`, `worker`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"bandwidth",
						"file_count",
						"cpu",
						"worker",
					),
				},
			},
			"limit": schema.Int64Attribute{
				Required:            true,
				Description:         "Amount the specified system resource type is limited by this rule.  Units are kb/s for bandwidth, files/s for file-count, processing percentage used for cpu, or percentage of maximum available workers.",
				MarkdownDescription: "Amount the specified system resource type is limited by this rule.  Units are kb/s for bandwidth, files/s for file-count, processing percentage used for cpu, or percentage of maximum available workers.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User-entered description of this performance rule.",
				MarkdownDescription: "User-entered description of this performance rule.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Whether this performance rule is currently in effect during its specified intervals.",
				MarkdownDescription: "Whether this performance rule is currently in effect during its specified intervals.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The system ID given to this performance rule.",
				MarkdownDescription: "The system ID given to this performance rule.",
			},
			"schedule": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "A schedule defining when during a week this performance rule is in effect.  If unspecified or null, the schedule will always be in effect.",
				MarkdownDescription: "A schedule defining when during a week this performance rule is in effect.  If unspecified or null, the schedule will always be in effect.",
				Attributes: map[string]schema.Attribute{
					"end": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "End time (inclusive) for this schedule, during its specified days.  Format is \"hh:mm\" (three-letter weekday name abbreviation, 24h format hour, and minute).  A null value indicates the end of the day (\"23:59\").",
						MarkdownDescription: "End time (inclusive) for this schedule, during its specified days.  Format is `hh:mm` (three-letter weekday name abbreviation, 24h format hour, and minute).  A null value indicates the end of the day (`23:59`).",
					},
					"begin": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Start time (inclusive) for this schedule, during its specified days.  Format is \"hh:mm\" (24h format hour, and minute).  A null value indicates the beginning of the day (\"00:00\").",
						MarkdownDescription: "Start time (inclusive) for this schedule, during its specified days.  Format is `hh:mm` (24h format hour, and minute).  A null value indicates the beginning of the day (`00:00`).",
					},
					"days_of_week": schema.SetAttribute{
						Description: "Set of days of the week during which this rule is in effect." +
							" Accepted values are 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'.",
						MarkdownDescription: "Set of days of the week during which this rule is in effect." +
							" Accepted values are `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`.",
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"monday",
									"tuesday",
									"wednesday",
									"thursday",
									"friday",
									"saturday",
									"sunday",
								),
							),
						},
					},
				},
			},
		},
	}
}
