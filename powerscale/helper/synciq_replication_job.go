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
func UpdateSyncIQReplicationJob(ctx context.Context, client *client.Client, jobID string, job powerscale.V1SyncJobExtendedExtended) (*http.Response,error) {
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
