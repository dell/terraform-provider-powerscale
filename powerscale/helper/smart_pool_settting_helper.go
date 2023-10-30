/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"terraform-provider-powerscale/client"
)

// GetSmartPoolSettingsSchema Get cluster config schema.
func GetSmartPoolSettingsSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SmartPools Settings information",
		Description:         "SmartPools Settings information",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"manage_io_optimization": schema.BoolAttribute{
				Description:         "Manage I/O optimization settings.",
				MarkdownDescription: "Manage I/O optimization settings.",
				Computed:            true,
			},
			"manage_io_optimization_apply_to_files": schema.BoolAttribute{
				Description:         "Apply to files with manually-managed I/O optimization settings.",
				MarkdownDescription: "Apply to files with manually-managed I/O optimization settings.",
				Computed:            true,
			},
			"manage_protection": schema.BoolAttribute{
				Description:         "Manage protection settings.",
				MarkdownDescription: "Manage protection settings.",
				Computed:            true,
			},
			"manage_protection_apply_to_files": schema.BoolAttribute{
				Description:         "Apply to files with manually-managed protection.",
				MarkdownDescription: "Apply to files with manually-managed protection.",
				Computed:            true,
			},
			"global_namespace_acceleration_enabled": schema.BoolAttribute{
				Description:         "Enable global namespace acceleration.",
				MarkdownDescription: "Enable global namespace acceleration.",
				Computed:            true,
			},
			"global_namespace_acceleration_state": schema.StringAttribute{
				Description:         "Whether or not namespace operation optimizations are currently in effect.",
				MarkdownDescription: "Whether or not namespace operation optimizations are currently in effect.",
				Computed:            true,
			},
			"protect_directories_one_level_higher": schema.BoolAttribute{
				Description:         "Increase directory protection to a higher requested protection than its contents.",
				MarkdownDescription: "Increase directory protection to a higher requested protection than its contents.",
				Computed:            true,
			},
			"spillover_enabled": schema.BoolAttribute{
				Description:         "Enable global spillover.",
				MarkdownDescription: "Enable global spillover.",
				Computed:            true,
			},
			"spillover_target": schema.SingleNestedAttribute{
				Description:         "Spillover data target.",
				MarkdownDescription: "Spillover data target.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description:         "Target pool ID if target specified, otherwise null.",
						MarkdownDescription: "Target pool ID if target specified, otherwise null.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Target pool name if target specified, otherwise null.",
						MarkdownDescription: "Target pool name if target specified, otherwise null.",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Type of target pool.",
						MarkdownDescription: "Type of target pool.",
						Computed:            true,
					},
				},
			},
			"ssd_l3_cache_default_enabled": schema.BoolAttribute{
				Description:         "Use SSDs as L3 cache by default for new node pools",
				MarkdownDescription: "Use SSDs as L3 cache by default for new node pools.",
				Computed:            true,
			},
			"ssd_qab_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of QAB blocks to place on SSDs.",
				MarkdownDescription: "Controls number of mirrors of QAB blocks to place on SSDs.",
				Computed:            true,
			},
			"ssd_system_btree_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of system B-tree blocks to place on SSDs.",
				MarkdownDescription: "Controls number of mirrors of system B-tree blocks to place on SSDs.",
				Computed:            true,
			},
			"ssd_system_delta_mirrors": schema.StringAttribute{
				Description:         "Controls number of mirrors of system delta blocks to place on SSDs.",
				MarkdownDescription: "Controls number of mirrors of system delta blocks to place on SSDs.",
				Computed:            true,
			},
			"virtual_hot_spare_deny_writes": schema.BoolAttribute{
				Description:         "Deny data writes to reserved disk space",
				MarkdownDescription: "Deny data writes to reserved disk space",
				Computed:            true,
			},
			"virtual_hot_spare_hide_spare": schema.BoolAttribute{
				Description:         "Subtract the space reserved for the virtual hot spare when calculating available free space",
				MarkdownDescription: "Subtract the space reserved for the virtual hot spare when calculating available free space",
				Computed:            true,
			},
			"virtual_hot_spare_limit_drives": schema.Int64Attribute{
				Description:         "The number of drives to reserve for the virtual hot spare, from 0-4.",
				MarkdownDescription: "The number of drives to reserve for the virtual hot spare, from 0-4.",
				Computed:            true,
			},
			"virtual_hot_spare_limit_percent": schema.Int64Attribute{
				Description:         "The percent space to reserve for the virtual hot spare, from 0-20.",
				MarkdownDescription: "The percent space to reserve for the virtual hot spare, from 0-20.",
				Computed:            true,
			},
			"default_transfer_limit_state": schema.StringAttribute{
				Description:         "How the default transfer limit value is applied.",
				MarkdownDescription: "How the default transfer limit value is applied.",
				Optional:            true,
			},
			"default_transfer_limit_pct": schema.NumberAttribute{
				Description:         "Applies to all storagepools that fall back on the default transfer limit. Stop moving files to this pool when this limit is met. The value must be between 0 and 100.",
				MarkdownDescription: "Applies to all storagepools that fall back on the default transfer limit. Stop moving files to this pool when this limit is met. The value must be between 0 and 100.",
				Optional:            true,
			},
		},
	}
}

// GetSmartPoolSettings Get SmartPool settings based on Onefs version.
func GetSmartPoolSettings(ctx context.Context, powerscaleClient *client.Client) (any, error) {
	if powerscaleClient.OnefsVersion.IsGreaterThan("9.4.0") {
		settings, _, err := powerscaleClient.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv16StoragepoolSettings(ctx).Execute()
		return settings, err
	}
	settings, _, err := powerscaleClient.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv5StoragepoolSettings(ctx).Execute()
	return settings, err
}
