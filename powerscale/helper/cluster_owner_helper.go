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
)

// GetClusterOwner retrieve cluster owner.
func GetClusterOwner(ctx context.Context, client *client.Client) (*powerscale.V1ClusterOwner, error) {
	ClusterOwner, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv1ClusterOwner(ctx).Execute()
	return ClusterOwner, err
}

// UpdateClusterOwner update cluster owner.
func UpdateClusterOwner(ctx context.Context, client *client.Client, v1ClusterOwner powerscale.V1ClusterOwner) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv1ClusterOwner(ctx).V1ClusterOwner(v1ClusterOwner).Execute()
	return err
}

// ManageClusterOwner for common operation for the Create and Update.
func ManageClusterOwner(ctx context.Context, client *client.Client, plan models.ClusterOwner) (state models.ClusterOwner, resp diag.Diagnostics) {
	var toUpdate powerscale.V1ClusterOwner
	// Get param from tf input
	err := ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterOwnerSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster owner",
			fmt.Sprintf("Could not read cluster owner param with error: %s", message),
		)
		return state, resp
	}

	err = UpdateClusterOwner(ctx, client, toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterOwnerSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster owner",
			message,
		)
		return state, resp
	}

	clusterOwner, err := GetClusterOwner(ctx, client)
	if err != nil {
		errStr := constants.ReadClusterOwnerSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster owner",
			message,
		)
		return state, resp
	}

	err = CopyFields(ctx, clusterOwner, &state)
	if err != nil {
		resp.AddError(
			"Error copying fields of cluster owner resource",
			err.Error(),
		)
		return state, resp
	}

	return state, nil
}
