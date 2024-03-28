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

// GetNfsExportSettings retrieve nfs export settings.
func GetNfsExportSettings(ctx context.Context, client *client.Client) (*powerscale.V2NfsSettingsExport, error) {
	nfsExportSettings, _, err := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsSettingsExport(ctx).Execute()
	return nfsExportSettings, err
}

// UpdateNfsExportSettings update nfs export settings.
func UpdateNfsExportSettings(ctx context.Context, client *client.Client, nfsExportSettings powerscale.V2NfsSettingsExportSettings) error {
	_, err := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv2NfsSettingsExport(ctx).V2NfsSettingsExport(nfsExportSettings).Execute()
	return err
}

// FilterNfsExportSettings retrieve nfs export settings by zone and scope.
func FilterNfsExportSettings(ctx context.Context, client *client.Client, filter *models.NfsSettingsExportDatasourceFilter) (*powerscale.V2NfsSettingsExport, error) {
	exportRequest := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsSettingsExport(ctx)
	if filter != nil {
		if zoneStr := filter.Zone.ValueString(); zoneStr != "" {
			exportRequest = exportRequest.Zone(zoneStr)
		}
		if scopeStr := filter.Scope.ValueString(); scopeStr != "" {
			exportRequest = exportRequest.Scope(scopeStr)
		}
	}
	nfsExportSettings, _, err := exportRequest.Execute()
	return nfsExportSettings, err
}
