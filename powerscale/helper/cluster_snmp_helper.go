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

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetClusterEmail retrieve cluster email.
func GetClusterSNMP(ctx context.Context, client *client.Client) (*powerscale.V16SnmpSettings, error) {
	clusterSNMP, _, err := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv16SnmpSettings(ctx).Execute()
	return clusterSNMP, err
}

// UpdateClusterEmail update cluster email.
func UpdateClusterSNMP(ctx context.Context, client *client.Client, v1ClusterSNMP powerscale.V16SnmpSettingsExtended) error {
	_, err := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv16SnmpSettings(ctx).V16SnmpSettings(v1ClusterSNMP).Execute()
	return err
}

func UpdateclusterSNMPResourceState(ctx context.Context, plan, clusterSNMPModel *models.ClusterSNMPModel, clusterSNMPResponse *powerscale.V16SnmpSettingsSettings) {
	clusterSNMPModel.ID = types.StringValue("cluster_snmp")
	clusterSNMPModel.Service = types.BoolValue(*clusterSNMPResponse.Service)
	clusterSNMPModel.ReadOnlyCommunity = types.StringValue(*clusterSNMPResponse.ReadOnlyCommunity)
	clusterSNMPModel.SnmpV1V2cAccess = types.BoolValue(*clusterSNMPResponse.SnmpV1V2cAccess)
	clusterSNMPModel.SnmpV3Access = types.BoolValue(*clusterSNMPResponse.SnmpV3Access)
	clusterSNMPModel.SnmpV3AuthProtocol = types.StringValue(*clusterSNMPResponse.SnmpV3AuthProtocol)
	clusterSNMPModel.SnmpV3PrivProtocol = types.StringValue(*clusterSNMPResponse.SnmpV3PrivProtocol)
	clusterSNMPModel.SnmpV3ReadOnlyUser = types.StringValue(*clusterSNMPResponse.SnmpV3ReadOnlyUser)
	clusterSNMPModel.SnmpV3SecurityLevel = types.StringValue(*clusterSNMPResponse.SnmpV3SecurityLevel)
	clusterSNMPModel.SystemContact = types.StringValue(*clusterSNMPResponse.SystemContact)
	clusterSNMPModel.SystemLocation = types.StringValue(*clusterSNMPResponse.SystemLocation)
	clusterSNMPModel.SnmpV3PrivPassword = types.StringValue(plan.SnmpV3PrivPassword.ValueString())
	clusterSNMPModel.SnmpV3Password = types.StringValue(plan.SnmpV3Password.ValueString())

}
