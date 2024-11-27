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
	powerscale "dell/powerscale-go-client"
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// GetReplicationReports gets a list of replication reports.
func GetReplicationReports(ctx context.Context, client *client.Client, state models.ReplicationReportsDatasourceModel) (*[]powerscale.V15SyncReport, error) {
	listRRParam := client.PscaleOpenAPIClient.SyncApi.GetSyncv15SyncReports(ctx)
	isLimit := false
	for _, filter := range state.ReplicationReportFilter {
		if !filter.Name.IsNull() && !filter.Value.IsNull() {
			filterName := filter.Name.ValueString()
			filterValue := filter.Value.ValueString()
			if filterName == "Sort" {
				listRRParam = listRRParam.Sort(filterValue)
			}
			if filterName == "NewerThan" {
				intVal, _ := strconv.Atoi(filterValue)
				listRRParam = listRRParam.NewerThan(int32(intVal))
			}
			if filterName == "PolicyName" {
				listRRParam = listRRParam.PolicyName(filterValue)
			}
			if filterName == "State" {
				listRRParam = listRRParam.State(filterValue)
			}
			if filterName == "Limit" {
				intVal, _ := strconv.Atoi(filterValue)
				listRRParam = listRRParam.Limit(int32(intVal))
				isLimit = true
			}
			if filterName == "ReportsPerPolicy" {
				intVal, _ := strconv.Atoi(filterValue)
				listRRParam = listRRParam.ReportsPerPolicy(int32(intVal))
			}
			if filterName == "Summary" {
				boolVal, _ := strconv.ParseBool(filterValue)
				listRRParam = listRRParam.Summary(boolVal)
			}
			if filterName == "Dir" {
				listRRParam = listRRParam.Dir(filterValue)
			}
		}
	}

	// if state.ReplicationReportFilter != nil {
	// 	if !state.ReplicationReportFilter.Sort.IsNull() {
	// 		listRRParam = listRRParam.Sort(state.ReplicationReportFilter.Sort.ValueString())
	// 	}
	// 	if !state.ReplicationReportFilter.NewerThan.IsNull() {
	// 		listRRParam = listRRParam.NewerThan(int32(state.ReplicationReportFilter.NewerThan.ValueInt64()))
	// 	}
	// 	if !state.ReplicationReportFilter.PolicyName.IsNull() {
	// 		listRRParam = listRRParam.PolicyName(state.ReplicationReportFilter.PolicyName.ValueString())
	// 	}
	// 	if !state.ReplicationReportFilter.State.IsNull() {
	// 		listRRParam = listRRParam.State(state.ReplicationReportFilter.State.ValueString())
	// 	}
	// 	if !state.ReplicationReportFilter.Limit.IsNull() {
	// 		listRRParam = listRRParam.Limit(int32(state.ReplicationReportFilter.Limit.ValueInt64()))
	// 	}
	// 	if !state.ReplicationReportFilter.ReportsPerPolicy.IsNull() {
	// 		listRRParam = listRRParam.ReportsPerPolicy(int32(state.ReplicationReportFilter.ReportsPerPolicy.ValueInt64()))
	// 	}
	// 	if !state.ReplicationReportFilter.Dir.IsNull() {
	// 		listRRParam = listRRParam.Dir(state.ReplicationReportFilter.Dir.ValueString())
	// 	}
	// 	if !state.ReplicationReportFilter.Summary.IsNull() {
	// 		listRRParam = listRRParam.Summary(state.ReplicationReportFilter.Summary.ValueBool())
	// 	}
	// }
	resp, _, err := listRRParam.Execute()
	if err != nil {
		return nil, err
	}
	totalReplicationReports := resp.Reports
	for resp.Resume != nil && isLimit {
		resumeReplicationReportParam := client.PscaleOpenAPIClient.SyncApi.GetSyncv15SyncReports(ctx).Resume(*resp.Resume)
		resp, _, err = resumeReplicationReportParam.Execute()
		if err != nil {
			return &totalReplicationReports, err
		}
		totalReplicationReports = append(totalReplicationReports, resp.Reports...)
	}
	return &totalReplicationReports, nil
}

// ReplicationReportDetailMapper maps the tfsdk struct to model.
func ReplicationReportDetailMapper(ctx context.Context, rr *powerscale.V15SyncReport) (models.ReplicationReportsDetail, error) {
	model := models.ReplicationReportsDetail{}
	err := CopyFields(ctx, rr, &model)
	return model, err
}
