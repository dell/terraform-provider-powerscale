/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// GetNetworkPools Get a list of Network Pools.
func GetNetworkPools(ctx context.Context, client *client.Client, state models.NetworkPoolDataSourceModel) (*[]powerscale.V12NetworkPool, error) {
	networkPoolParams := client.PscaleOpenAPIClient.NetworkApi.GetNetworkv12NetworkPools(ctx)

	if state.NetworkPoolFilter != nil && !state.NetworkPoolFilter.Subnet.IsNull() {
		networkPoolParams = networkPoolParams.Subnet(state.NetworkPoolFilter.Subnet.ValueString())
	}
	if state.NetworkPoolFilter != nil && !state.NetworkPoolFilter.Groupnet.IsNull() {
		networkPoolParams = networkPoolParams.Groupnet(state.NetworkPoolFilter.Groupnet.ValueString())
	}
	if state.NetworkPoolFilter != nil && !state.NetworkPoolFilter.AccessZone.IsNull() {
		networkPoolParams = networkPoolParams.AccessZone(state.NetworkPoolFilter.AccessZone.ValueString())
	}
	if state.NetworkPoolFilter != nil && !state.NetworkPoolFilter.AllocMethod.IsNull() {
		networkPoolParams = networkPoolParams.AllocMethod(state.NetworkPoolFilter.AllocMethod.ValueString())
	}

	networkPools, _, err := networkPoolParams.Execute()
	return &networkPools.Pools, err
}

// NetworkPoolDetailMapper Does the mapping from response to model.
//
//go:noinline
func NetworkPoolDetailMapper(ctx context.Context, networkPool *powerscale.V12NetworkPool) (models.NetworkPoolDetailModel, error) {
	model := models.NetworkPoolDetailModel{}
	err := CopyFields(ctx, networkPool, &model)
	return model, err
}
