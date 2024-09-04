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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateClusterIdentity calls the API to update the cluster identity.
func UpdateClusterIdentity(ctx context.Context, client *client.Client, clusterIdentity powerscale.V3ClusterIdentity) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv3ClusterIdentity(ctx).V3ClusterIdentity(clusterIdentity).Execute()
	return err
}

// ManageClusterIdentity reads and updates the cluster identity.
func ManageClusterIdentity(ctx context.Context, plan models.ClusterIdentityResource, state *models.ClusterIdentityResource, client *client.Client) diag.Diagnostics {
	var toUpdate powerscale.V3ClusterIdentity
	var diags diag.Diagnostics
	err := ReadFromState(ctx, &plan, &toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterIdentitySettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error updating cluster identity",
			fmt.Sprintf("Could not read cluster identity param with error: %s", message),
		)
		return diags
	}

	err = UpdateClusterIdentity(ctx, client, toUpdate)
	if err != nil {
		errStr := constants.UpdateClusterIdentitySettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error updating cluster identity",
			message,
		)
		return diags
	}

	clusterIdentity, err := GetClusterIdentity(ctx, client)
	if err != nil {
		errStr := constants.ReadClusterIdentitySettingsErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error fetching updated cluster identity settings",
			message,
		)
		return diags
	}
	UpdateClusterIdentityState(state, clusterIdentity)

	return diags
}

// UpdateClusterIdentityState is a function to update cluster identity state.
func UpdateClusterIdentityState(state *models.ClusterIdentityResource, clusterIdentity *powerscale.V1ClusterIdentity) diag.Diagnostics {
	mappingLogonObject2Type := map[string]attr.Type{
		"motd":        types.StringType,
		"motd_header": types.StringType,
	}

	LogonElemMap := map[string]attr.Value{
		"motd":        types.StringValue(clusterIdentity.Logon.Motd),
		"motd_header": types.StringValue(clusterIdentity.Logon.MotdHeader),
	}

	lookuplogonObject, diags := types.ObjectValue(mappingLogonObject2Type, LogonElemMap)
	if diags.HasError() {
		diags.AddError(
			"Unable to update cluster identity state",
			"Unmarshalling logon object failed",
		)
		return diags
	}
	state.Logon = lookuplogonObject
	state.Name = types.StringValue(clusterIdentity.Name)
	state.Description = types.StringValue(clusterIdentity.Description)
	state.ID = types.StringValue("cluster_identity")
	return nil
}
