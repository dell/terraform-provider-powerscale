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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// SnapshotScheduleDataSourceModel Defines the Snapshot schedule datasource model.
type SnapshotScheduleDataSourceModel struct {
	ID                types.String             `tfsdk:"id"`
	SnapshotSchedules []SnapshotScheduleEntity `tfsdk:"schedules"`

	//filter
	SnapshotScheduleFilter *SnapshotScheduleFilter `tfsdk:"filter"`
}

// SnapshotScheduleEntity Defines the Snapshot schedule entity.
type SnapshotScheduleEntity struct {
	// Alias name to create for each snapshot.
	Alias types.String `tfsdk:"alias"`
	// Time in seconds added to creation time to construction expiration time.
	Duration types.Int64 `tfsdk:"duration"`
	// The system ID given to the schedule.
	ID types.Int64 `tfsdk:"id"`
	// The schedule name.
	Name types.String `tfsdk:"name"`
	// Unix Epoch time of next snapshot to be created.
	NextRun types.Int64 `tfsdk:"next_run"`
	// Formatted name (see pattern) of next snapshot to be created.
	NextSnapshot types.String `tfsdk:"next_snapshot"`
	// The /ifs path snapshotted.
	Path types.String `tfsdk:"path"`
	// Pattern expanded with strftime to create snapshot names.
	Pattern types.String `tfsdk:"pattern"`
	// The isidate compatible natural language description of the schedule.
	Schedule types.String `tfsdk:"schedule"`
}

// SnapshotScheduleFilter defines the snapshot schedule filter.
type SnapshotScheduleFilter struct {
	// filters supported by api

	// The field that will be used for sorting. Choices are id, name, path, pattern, schedule, duration, alias, next_run, and next_snapshot. Default is id.
	Sort types.String `tfsdk:"sort"`
	// Return no more than this many results at once.
	Limit types.Int64 `tfsdk:"limit"`
	// The direction of the sort.Supported Values:ASC , DESC
	Dir types.String `tfsdk:"dir"`

	// custom names filter to filter on snapshot schedule names
	Names []types.String `tfsdk:"names"`
}

// SnapshotScheduleResource Defines the Snapshot schedule entity.
type SnapshotScheduleResource struct {
	// Alias name to create for each snapshot.
	Alias types.String `tfsdk:"alias"`
	// Time in seconds added to creation time to construction expiration time.
	Duration types.Int64 `tfsdk:"duration"`
	// The system ID given to the schedule.
	ID types.String `tfsdk:"id"`
	// The schedule name.
	Name types.String `tfsdk:"name"`
	// Unix Epoch time of next snapshot to be created.
	NextRun types.Int64 `tfsdk:"next_run"`
	// Formatted name (see pattern) of next snapshot to be created.
	NextSnapshot types.String `tfsdk:"next_snapshot"`
	// The /ifs path snapshotted.
	Path types.String `tfsdk:"path"`
	// Pattern expanded with strftime to create snapshot names.
	Pattern types.String `tfsdk:"pattern"`
	// The isidate compatible natural language description of the schedule.
	Schedule types.String `tfsdk:"schedule"`
	//Time value in String for which snapshots created by this snapshot schedule should be retained. Values supported are of format : "Never Expires, x Seconds(s), x Minute(s), x Hour(s), x Week(s), x Day(s), x Month(s), x Year(s) where x can be any integer value.
	RetentionTime types.String `tfsdk:"retention_time"`
}
