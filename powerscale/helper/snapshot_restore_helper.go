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
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CreateSnapshotRestoreJob creates a job.
func CreateSnapshotRestoreJob(ctx context.Context, client *client.Client, jobPayload powerscale.V10JobJob) (*powerscale.Createv1JobJobResponse, error) {
	response, _, err := client.PscaleOpenAPIClient.JobApi.CreateJobv10JobJob(ctx).V10JobJob(jobPayload).Execute()
	return response, err
}

// GetSnapshotRestoreJob gets job details.
func GetSnapshotRestoreJob(ctx context.Context, client *client.Client, jobID string) (*powerscale.V10JobJobExtended, error) {
	response, _, err := client.PscaleOpenAPIClient.JobApi.GetJobv7JobJob(ctx, jobID).Execute()
	if err != nil {
		return nil, err
	}
	tflog.Info(ctx, "job response", map[string]interface{}{
		"fields": response,
	})
	return &response.Jobs[0], err
}

// ManageSnapshotRestore manages the snapshot restore.
func ManageSnapshotRestore(ctx context.Context, client *client.Client, plan models.SnapshotRestoreModel) (state models.SnapshotRestoreModel, resp diag.Diagnostics) {
	state = plan

	// Check if snaprevert params are provided
	if !plan.SnapRevertParams.IsNull() {
		var snapRevert models.SnapRevertParamsModel
		var path string

		diag := plan.SnapRevertParams.As(ctx, &snapRevert, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		if diag.HasError() {
			return state, diag
		}

		// Get Snapshot ID if name is provided
		if !snapRevert.SnapName.IsNull() {
			response, err := GetAllSnapshots(ctx, client)
			if err != nil {
				errStr := constants.ReadSnapshotErrorMessage + "with error: "
				message := GetErrorString(err, errStr)
				resp.AddError(
					"Error getting the list of snapshots",
					message,
				)
				return state, resp
			}

			for _, resp := range response {
				if resp.Name == snapRevert.SnapName.ValueString() {
					snapRevert.SnapID = types.Int32Value(resp.Id)
					path = resp.Path
					break
				}
			}

			if path == "" {
				resp.AddError(
					fmt.Sprintf("Error getting the snapshot with name %v", snapRevert.SnapName.String()),
					fmt.Sprintf("Error getting the snapshot with name %v", snapRevert.SnapName.String()),
				)
				return state, resp
			}
		} else {
			// Get Snapshot details if ID is provided
			response, err := GetSpecificSnapshot(ctx, client, snapRevert.SnapID.String())
			if err != nil {
				errStr := constants.ReadSnapshotErrorMessage + "with error: "
				message := GetErrorString(err, errStr)
				resp.AddError(
					fmt.Sprintf("Error getting the snapshot with id %v", snapRevert.SnapID.String()),
					message,
				)
				return state, resp
			}
			path = response.Path
		}

		// Populate the payload for creating snaprevert domain
		payload := powerscale.V10JobJob{
			Type:     "DomainMark",
			AllowDup: snapRevert.AllowDup.ValueBoolPointer(),
			DomainmarkParams: &powerscale.V1JobJobDomainmarkParams{
				Type: "SnapRevert",
				Root: path,
			},
		}
		createResponse, _ := CreateSnapshotRestoreJob(ctx, client, payload)
		strID := strconv.Itoa(int(createResponse.Id))
		response, err := GetSnapshotRestoreJob(ctx, client, strID)
		if err != nil {
			errStr := constants.ReadSnapshotRestoreJobErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error getting job",
				message,
			)
			return state, resp
		}
		_, diag = CheckJobStatus(ctx, client, strID, response)

		// Populate the payload for creating snaprevert job
		payload = powerscale.V10JobJob{
			Type:     "SnapRevert",
			AllowDup: snapRevert.AllowDup.ValueBoolPointer(),
			SnaprevertParams: &powerscale.V1JobJobSnaprevertParams{
				Snapid: snapRevert.SnapID.ValueInt32(),
			},
		}

		createResponse, err = CreateSnapshotRestoreJob(ctx, client, payload)
		if err != nil {
			errStr := constants.CreateSnapshotRestoreJobErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error creating job",
				message,
			)
			return state, resp
		}

		strID = strconv.Itoa(int(createResponse.Id))
		tflog.Info(ctx, fmt.Sprintf("SnapRestore job id: %v", createResponse.Id))
		response, err = GetSnapshotRestoreJob(ctx, client, strID)
		if err != nil {
			errStr := constants.ReadSnapshotRestoreJobErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error getting job",
				message,
			)
			return state, resp
		}

		response, diag = CheckJobStatus(ctx, client, strID, response)
		resp.Append(diag...)

		if response.State == "failed" {
			resp.AddError(
				"Error getting job report",
				"Please check if snaprevert domain is created",
			)
			return state, resp
		}
	}
	state.ID = types.StringValue("snapshot_restore")
	return state, nil
}

// CheckJobStatus checks the job status
func CheckJobStatus(ctx context.Context, client *client.Client, jobID string, response *powerscale.V10JobJobExtended) (res *powerscale.V10JobJobExtended, resp diag.Diagnostics) {
	var err error
	for !(response.State == "succeeded" || response.State == "failed") {
		time.Sleep(time.Second)
		response, err = GetSnapshotRestoreJob(ctx, client, jobID)
		if err != nil {
			errStr := constants.ReadSnapshotRestoreJobErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error getting job",
				message,
			)
			return nil, resp
		}
	}
	return response, nil
}
