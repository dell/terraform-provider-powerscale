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

// GetNfsExportSettingsByZone retrieve nfs export settings by zone.
func GetNfsExportSettingsByZone(ctx context.Context, client *client.Client, zone string) (*powerscale.V2NfsSettingsExport, error) {
	nfsExportSettings, _, err := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsSettingsExport(ctx).Zone(zone).Execute()
	return nfsExportSettings, err
}

// UpdateNfsExportSettings update nfs export settings.
func UpdateNfsExportSettings(ctx context.Context, client *client.Client, nfsExportSettings powerscale.V2NfsSettingsExportSettings, zone string) error {
	_, err := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv2NfsSettingsExport(ctx).V2NfsSettingsExport(nfsExportSettings).Zone(zone).Execute()
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

// SetDefaultValues set default values for nfs export settings
func SetDefaultValues(nfsExportSettingsModel *models.NfsexportsettingsModel, nfsExportSettings *powerscale.V2NfsSettingsExportSettings) error {
	if nfsExportSettings == nil {
		return nil
	}
	if nfsExportSettingsModel.Snapshot.ValueString() == "-" {
		value := "@DEFAULT"
		nfsExportSettings.Snapshot = &value
	}
	return nil
}

// GetSpecifiedZone retrieve zone from plan, or state if it is not defined in plan
func GetSpecifiedZone(plan *models.NfsexportsettingsModel, state *models.NfsexportsettingsModel) string {
	if plan.Zone.ValueString() != "" {
		return plan.Zone.ValueString()
	}
	return state.Zone.ValueString()
}
