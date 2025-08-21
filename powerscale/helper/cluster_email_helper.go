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

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// GetClusterEmail retrieve cluster email.
func GetClusterEmail(ctx context.Context, client *client.Client) (*powerscale.V1ClusterEmail, error) {
	clusterEmail, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv1ClusterEmail(ctx).Execute()
	return clusterEmail, err
}

// UpdateClusterEmail update cluster email.
func UpdateClusterEmail(ctx context.Context, client *client.Client, v1ClusterEmail powerscale.V1ClusterEmailExtended) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv1ClusterEmail(ctx).V1ClusterEmail(v1ClusterEmail).Execute()
	return err
}

// GetV21ClusterEmail retrieve cluster email.
func GetV21ClusterEmail(ctx context.Context, client *client.Client) (*powerscale.V1ClusterEmail, error) {
	clusterEmail, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv21ClusterEmail(ctx).Execute()
	return clusterEmail, err
}

// UpdateV21ClusterEmail update cluster email.
func UpdateV21ClusterEmail(ctx context.Context, client *client.Client, v21ClusterEmail powerscale.V1ClusterEmailExtended) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv21ClusterEmail(ctx).V21ClusterEmail(v21ClusterEmail).Execute()
	return err
}

// ManageClusterEmail manages the create and update of cluster email.
func ManageClusterEmail(ctx context.Context, client *client.Client, plan models.ClusterEmail) (state models.ClusterEmail, resp diag.Diagnostics) {

	var toUpdate powerscale.V1ClusterEmailExtended
	// Get param from tf input
	err := ReadFromState(ctx, &plan.Settings, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster email",
			fmt.Sprintf("Could not read cluster email param with error: %s", message),
		)
		return state, resp
	}

	// if the field is set to empty, set it to null to update back to default
	// if not set, unset the field to not update it
	if plan.Settings.UserTemplate.IsUnknown() {
		toUpdate.UserTemplate.Unset()
	} else if plan.Settings.UserTemplate.ValueString() == "" {
		toUpdate.UserTemplate.Set(nil)
	}

	clusterVersion, err := GetClusterVersion(ctx, client)
	if err != nil {
		errStr := constants.ReadClusterErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error reading cluster version",
			message,
		)
		return state, resp
	}

	if VersionGTE(clusterVersion, "9.10.0.0") {
		err = UpdateV21ClusterEmail(ctx, client, toUpdate)
	} else {
		err = UpdateClusterEmail(ctx, client, toUpdate)
	}

	if err != nil {
		errStr := constants.UpdateClusterEmailSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster email",
			message,
		)
		return state, resp
	}

	var clusterEmail *powerscale.V1ClusterEmail
	if VersionGTE(clusterVersion, "9.10.0.0") {
		clusterEmail, err = GetV21ClusterEmail(ctx, client)
	} else {
		clusterEmail, err = GetClusterEmail(ctx, client)
	}

	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster email",
			message,
		)
		return state, resp
	}

	err = CopyFields(ctx, clusterEmail, &state)
	if err != nil {
		resp.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return state, resp
	}
	state.ID = types.StringValue("cluster_email")
	state.Settings.SMTPAuthPasswd = plan.Settings.SMTPAuthPasswd
	if state.Settings.SMTPAuthPasswd.IsUnknown() {
		state.Settings.SMTPAuthPasswd = types.StringValue("")
	}
	return state, resp
}

// ReadClusterEmail manages read and import of cluster email.
func ReadClusterEmail(ctx context.Context, client *client.Client, state *models.ClusterEmail) (resp diag.Diagnostics) {
	clusterVersion, err := GetClusterVersion(ctx, client)
	if err != nil {
		errStr := constants.ReadClusterErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error reading cluster version",
			message,
		)
		return resp
	}

	var clusterEmail *powerscale.V1ClusterEmail
	if VersionGTE(clusterVersion, "9.10.0.0") {
		clusterEmail, err = GetV21ClusterEmail(ctx, client)
	} else {
		clusterEmail, err = GetClusterEmail(ctx, client)
	}

	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error reading cluster email",
			message,
		)
		return resp
	}
	err = CopyFields(ctx, clusterEmail, state)
	if err != nil {
		resp.AddError(
			"Error copying fields of cluster email resource",
			err.Error(),
		)
		return resp
	}
	state.ID = types.StringValue("cluster_email")
	if state.Settings.SMTPAuthPasswd.IsUnknown() {
		state.Settings.SMTPAuthPasswd = types.StringValue("")
	}
	return resp
}
