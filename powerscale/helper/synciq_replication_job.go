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
	"net/http"
	"terraform-provider-powerscale/client"
)

// GetSyncIQReplicationJob get syncIQ replication job.
func GetSyncIQReplicationJob(ctx context.Context, client *client.Client, jobID string) (*powerscale.V1SyncJobsExtended, *http.Response, error) {
	resp, httpResp, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv1SyncJob(ctx, jobID).Execute()
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, err
}

// CreateSyncIQReplicationJob create syncIQ replication job.
func CreateSyncIQReplicationJob(ctx context.Context, client *client.Client, job powerscale.V1SyncJob) (string, error) {
	if job.Action != nil && *job.Action == "run" {
		job.Action = nil
	}
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.CreateSyncv1SyncJob(ctx).V1SyncJob(job).Execute()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// UpdateSyncIQReplicationJob update syncIQ replication job.
func UpdateSyncIQReplicationJob(ctx context.Context, client *client.Client, jobID string, job powerscale.V1SyncJobExtendedExtended) (*http.Response, error) {
	resp, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv1SyncJob(ctx, jobID).V1SyncJob(job).Execute()
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DeleteSyncIQReplicationJob delete syncIQ replication job.
func DeleteSyncIQReplicationJob(ctx context.Context, client *client.Client, jobID string) error {
	deleteJob := powerscale.V1SyncJobExtendedExtended{
		State: "canceled",
	}
	resp, err := UpdateSyncIQReplicationJob(ctx, client, jobID, deleteJob)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil // already deleted
	}
	return err
}
