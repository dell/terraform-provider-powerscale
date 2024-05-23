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

// GetSmbShareSettings get smb share settings.
func GetSmbShareSettings(ctx context.Context, client *client.Client, scope string, zone string) (*powerscale.V7SmbSettingsShare, error) {
	param := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv7SmbSettingsShare(ctx)

	if zone != "" {
		param = param.Zone(zone)
	}

	if scope != "" {
		param = param.Scope(scope)
	}
	response, _, err := param.Execute()
	return response, err
}

// UpdateSmbShareSettings update smb share settings.
func UpdateSmbShareSettings(ctx context.Context, client *client.Client, v7SmbSettingsShare powerscale.V7SmbSettingsShareSettings, zone string) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv7SmbSettingsShare(ctx).V7SmbSettingsShare(v7SmbSettingsShare)

	if zone != "" {
		updateParam = updateParam.Zone(zone)
	}
	_, err := updateParam.Execute()
	return err
}
