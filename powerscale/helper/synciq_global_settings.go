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
)

// UpdateSyncIQGlobalSettings updates the SyncIQ global settings
func UpdateSyncIQGlobalSettings(ctx context.Context, client *client.Client, edit powerscale.V16SyncSettingsExtended) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv16SyncSettings(ctx).V16SyncSettings(edit).Execute()
	return err
}

// GetSyncIQGlobalSettings fetches the SyncIQ global settings
func GetSyncIQGlobalSettings(ctx context.Context, client *client.Client) (*powerscale.V16SyncSettings, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv16SyncSettings(ctx).Execute()
	return resp, err
}
