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
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetClusterTime retrieve cluster Time.
func GetClusterTime(ctx context.Context, client *client.Client) (*powerscale.V3ClusterTime, error) {
	ClusterTime, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterTime(ctx).Execute()
	return ClusterTime, err
}

// UpdateClusterTime update cluster Time.
func UpdateClusterTime(ctx context.Context, client *client.Client, V3ClusterTimeExtended powerscale.V3ClusterTimeExtended) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv3ClusterTime(ctx).V3ClusterTime(V3ClusterTimeExtended).Execute()
	return err
}

// GetClusterTimeZone retrieve cluster Timezone.
func GetClusterTimeZone(ctx context.Context, client *client.Client) (*powerscale.V3ClusterTimezone, error) {
	ClusterTimeZone, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterTimezone(ctx).Execute()
	return ClusterTimeZone, err
}

// UpdateClusterTimeZone update cluster Timezone.
func UpdateClusterTimeZone(ctx context.Context, client *client.Client, V3TimezoneRegionTimezone powerscale.V3TimezoneRegionTimezone) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv3TimezoneSettings(ctx).V3TimezoneSettings(V3TimezoneRegionTimezone).Execute()
	return err
}

// ManageClusterTime for common operation for the Create and Update.
func ManageClusterTime(ctx context.Context, client *client.Client, plan models.ClusterTime) (state models.ClusterTime, resp diag.Diagnostics) {
	var timeUpdate powerscale.V3ClusterTimeExtended
	var timezoneUpdate powerscale.V3TimezoneRegionTimezone

	// Get param from tf input
	err := ReadFromState(ctx, &plan, &timeUpdate)
	if err != nil {
		errStr := constants.UpdateClusterTimeSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster time",
			fmt.Sprintf("Could not read cluster time param with error: %s", message),
		)
		return state, resp
	}

	if !plan.Path.IsNull() && plan.Path.ValueString() != "" {
		timezoneUpdate.Path = plan.Path.ValueString()
		err = UpdateClusterTimeZone(ctx, client, timezoneUpdate)
		if err != nil {
			errStr := constants.UpdateClusterTimezoneSettingsErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error updating cluster time",
				message,
			)
			return state, resp
		}

	}

	if !plan.Date.IsNull() && plan.Date.ValueString() != "" && !plan.Time.IsNull() && plan.Time.ValueString() != "" {

		// Combine date and time into one string
		dateTimeStr := plan.Date.ValueString() + " " + plan.Time.ValueString()

		// Parse the combined string into a time.Time object
		layout := "01/02/2006 15:04"
		t, err := time.Parse(layout, dateTimeStr)
		if err != nil {
			errStr := constants.UpdateClusterTimeSettingsErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error updating cluster time",
				message,
			)
			return state, resp
		}

		timeUpdate.Time = int32(t.Unix())

		err = UpdateClusterTime(ctx, client, timeUpdate)
		if err != nil {
			errStr := constants.UpdateClusterTimeSettingsErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			resp.AddError(
				"Error updating cluster time",
				message,
			)
			return state, resp
		}
	}

	state, dig := ReadClusterTimeDetails(ctx, client, plan)
	if dig.HasError() {
		return state, dig
	}

	return state, nil
}

// ReadClusterTimeDetails read cluster Time details.
func ReadClusterTimeDetails(ctx context.Context, client *client.Client, plan models.ClusterTime) (state models.ClusterTime, resp diag.Diagnostics) {
	clusterTime, err := GetClusterTime(ctx, client)
	if err != nil {
		errStr := constants.ReadClusterTimeSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster time",
			message,
		)
		return state, resp
	}

	clusterTimezone, err := GetClusterTimeZone(ctx, client)
	if err != nil {
		errStr := constants.ReadClusterTimezoneSettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		resp.AddError(
			"Error updating cluster time",
			message,
		)
		return state, resp
	}

	state = plan

	var clusterConfigTimezone models.ClusterTimezoneSetting

	err = CopyFields(ctx, clusterTimezone, &clusterConfigTimezone)
	if err != nil {
		resp.AddError(
			"Error copying fields of cluster time resource",
			err.Error(),
		)
		return state, resp
	}

	state.Abbreviation = types.StringValue(clusterConfigTimezone.Settings.Abbreviation.ValueString())

	state.Path = types.StringValue(clusterConfigTimezone.Settings.Path.ValueString())

	state.TimeMillis = types.Int32Value(int32(*clusterTime.Nodes[0].Time))

	return state, resp
}
