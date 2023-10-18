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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateGroupnetDataSourceState updates datasource state.
func UpdateGroupnetDataSourceState(ctx context.Context, groupnetState *models.GroupnetDataSourceModel, groupnetResponse []powerscale.V10NetworkGroupnetExtended) (err error) {
	for _, groupnet := range groupnetResponse {
		var model models.GroupnetModel

		if err = CopyFields(ctx, groupnet, &model); err != nil {
			return
		}

		model.DNSResolverRotate = types.BoolValue(false)
		if groupnet.HasDnsOptions() {
			for _, option := range groupnet.DnsOptions {
				if option == "rotate" {
					model.DNSResolverRotate = types.BoolValue(true)
					break
				}
			}
		}

		groupnetState.Groupnets = append(groupnetState.Groupnets, model)
	}
	return
}

// GetAllGroupnets returns all groupnets.
func GetAllGroupnets(ctx context.Context, client *client.Client) (groupnets []powerscale.V10NetworkGroupnetExtended, err error) {

	groupnetParams := client.PscaleOpenAPIClient.NetworkApi.ListNetworkv10NetworkGroupnets(ctx)
	result, _, err := groupnetParams.Execute()
	if err != nil {
		errStr := constants.ReadGroupnetErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting groupnets: %s", message)
	}

	for {
		groupnets = append(groupnets, result.Groupnets...)
		if result.Resume == nil || *result.Resume == "" {
			break
		}

		groupnetParams = client.PscaleOpenAPIClient.NetworkApi.ListNetworkv10NetworkGroupnets(ctx).Resume(*result.Resume)
		if result, _, err = groupnetParams.Execute(); err != nil {
			errStr := constants.ReadGroupnetErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting groupnets with resume: %s", message)
		}
	}
	return
}
