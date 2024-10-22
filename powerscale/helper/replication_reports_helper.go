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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// GetRoles Get a list of replication reports.
func GetReplicationReports(ctx context.Context, client *client.Client, state models.ReplicationReportsDatasourceModel) (*powerscale.V15SyncReports, error) {
	listRRParam := client.PscaleOpenAPIClient.SyncApi.GetSyncv15SyncReports(ctx)
	if state.ReplicationReportFilter != nil {
		if !state.ReplicationReportFilter.Sort.IsNull() {
			listRRParam = listRRParam.Sort(state.ReplicationReportFilter.Sort.ValueString())
		}
		if !state.ReplicationReportFilter.Resume.IsNull() {
			listRRParam = listRRParam.Resume(state.ReplicationReportFilter.Resume.ValueString())
		}
		if !state.ReplicationReportFilter.NewerThan.IsNull() {
			listRRParam = listRRParam.NewerThan(int32(state.ReplicationReportFilter.NewerThan.ValueInt64()))
		}
		if !state.ReplicationReportFilter.PolicyName.IsNull() {
			listRRParam = listRRParam.PolicyName(state.ReplicationReportFilter.PolicyName.ValueString())
		}
		if !state.ReplicationReportFilter.State.IsNull() {
			listRRParam = listRRParam.State(state.ReplicationReportFilter.State.ValueString())
		}
		if !state.ReplicationReportFilter.Limit.IsNull() {
			listRRParam = listRRParam.Limit(int32(state.ReplicationReportFilter.Limit.ValueInt64()))
		}
		if !state.ReplicationReportFilter.ReportsPerPolicy.IsNull() {
			listRRParam = listRRParam.ReportsPerPolicy(int32(state.ReplicationReportFilter.ReportsPerPolicy.ValueInt64()))
		}
		if !state.ReplicationReportFilter.Dir.IsNull() {
			listRRParam = listRRParam.Dir(state.ReplicationReportFilter.Dir.ValueString())
		}
		if !state.ReplicationReportFilter.Summary.IsNull() {
			listRRParam = listRRParam.Summary(state.ReplicationReportFilter.Summary.ValueBool())
		}

	}
	resp, _, err := listRRParam.Execute()
	return resp, err
}

func ReplicationReportDetailMapper(ctx context.Context, rr *powerscale.V15SyncReport) (models.ReplicationReportsDetail, error) {
	model := models.ReplicationReportsDetail{}
	err := CopyFields(ctx, rr, &model)
	return model, err
}
