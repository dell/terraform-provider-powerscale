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
)

// GetSmbServerSettings retrieve smb server settings.
func GetSmbServerSettings(ctx context.Context, client *client.Client, scope string) (*powerscale.V6SmbSettingsGlobal, error) {
	getParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv7SmbSettingsGlobal(ctx)
	getParam = getParam.Scope(scope)
	smbServerSettings, _, err := getParam.Execute()
	return smbServerSettings, err
}

// UpdateSmbServerSettings update smb server settings.
func UpdateSmbServerSettings(ctx context.Context, client *client.Client, smbServerSettings powerscale.V6SmbSettingsGlobalSettings) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv7SmbSettingsGlobal(ctx)
	updateParam = updateParam.V7SmbSettingsGlobal(smbServerSettings)
	_, err := updateParam.Execute()
	return err
}

// FilterSmbServerSettings filter smb server settings.
func FilterSmbServerSettings(ctx context.Context, client *client.Client, filter *models.SmbServerSettingsFilter) (*powerscale.V6SmbSettingsGlobal, error) {
	filterParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv7SmbSettingsGlobal(ctx)

	if filter != nil {
		if scopeStr := filter.Scope.ValueString(); scopeStr != "" {
			filterParam = filterParam.Scope(scopeStr)
		}
	}

	smbServerSettings, _, err := filterParam.Execute()
	return smbServerSettings, err
}
