/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
)

// GetStoragepoolTier gets storage pool tier.
func GetStoragepoolTier(ctx context.Context, client *client.Client, path string) (*powerscale.V16StoragepoolTiersExtended, error) {
	StoragepoolTier, _, err := client.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv16StoragepoolTier(ctx, path).Execute()
	return StoragepoolTier, err
}

// GetAllStoragepoolTiers returns the full list of storage pool tiers.
func GetAllStoragepoolTiers(ctx context.Context, client *client.Client) ([]powerscale.V16StoragepoolTierExtended, error) {
	StoragepoolTierParams := client.PscaleOpenAPIClient.StoragepoolApi.ListStoragepoolv16StoragepoolTiers(ctx)

	StoragepoolTiers, _, err := StoragepoolTierParams.Execute()
	if err != nil {
		errStr := constants.ReadStoragepoolTiersErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting storagepool tiers: %s", message)
	}
	return StoragepoolTiers.Tiers, nil
}
