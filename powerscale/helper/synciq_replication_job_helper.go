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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSyncIQReplicationJobs gets the list of SyncIQ jobs.
func GetSyncIQReplicationJobs(ctx context.Context, client *client.Client, filter *models.SyncIQJobFilterModel) (*powerscale.V7SyncJobs, error) {
	jobParams := client.PscaleOpenAPIClient.SyncApi.ListSyncv7SyncJobs(context.Background())
	if filter != nil {
		if !filter.Sort.IsNull() {
			jobParams = jobParams.Sort(filter.Sort.ValueString())
		}

		if !filter.Dir.IsNull() {
			jobParams = jobParams.Dir(filter.Dir.ValueString())
		}

		if !filter.Limit.IsNull() {
			jobParams = jobParams.Limit(int32(filter.Limit.ValueInt64()))
		}

		if !filter.State.IsNull() {
			jobParams = jobParams.State(filter.State.ValueString())
		}
	}

	resp, _, err := jobParams.Execute()
	if err != nil {
		return nil, err
	}

	// Pagination
	if resp.Resume != nil && (filter == nil || filter.Limit.IsNull()) {
		jobParams = jobParams.Resume(*resp.Resume)
		newresp, _, errAdd := jobParams.Execute()
		if errAdd != nil {
			return resp, err
		}
		resp.Resume = newresp.Resume
		resp.Jobs = append(resp.Jobs, newresp.Jobs...)
	}
	return resp, err
}

// ManageDataSourceSyncIQReplicationJob gets the details of SyncIQ replication job and set the state.
func ManageDataSourceSyncIQReplicationJob(ctx context.Context, plan *models.SyncIQReplicationJobDataSourceModel, client *client.Client) (state *models.SyncIQReplicationJobDataSourceModel, resp diag.Diagnostics) {
	response, err := GetSyncIQReplicationJobs(ctx, client, plan.SyncIQJobFilter)
	if err != nil {
		errStr := constants.ReadSyncIQReplicationJobErrorMessage + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error getting the synciq jobs",
			message,
		)
		return state, resp
	}

	var jobs []models.SyncIQReplicationJobModel
	for _, job := range response.Jobs {
		var temp models.SyncIQReplicationJobModel

		err = CopyFieldsToNonNestedModel(ctx, job, &temp)
		if err != nil {
			resp.AddError(
				"Unable to update SyncIQ replication job state",
				fmt.Sprintf("Unable to update SyncIQ replication job state with error %s", err.Error()),
			)
			return state, resp
		}

		jobs = append(jobs, temp)
	}
	state = &models.SyncIQReplicationJobDataSourceModel{
		SyncIQReplicationJobs: jobs,
		ID:                    types.StringValue("synciq_job_datasource"),
	}
	return state, nil
}
