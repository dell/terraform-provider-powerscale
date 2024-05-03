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

// GetNfsZoneSettings retrieve nfs zone settings.
func GetNfsZoneSettings(ctx context.Context, client *client.Client, zone string) (*powerscale.V2NfsSettingsZone, error) {
	getParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsSettingsZone(ctx)
	getParam = getParam.Zone(zone)
	nfsZoneSettings, _, err := getParam.Execute()
	return nfsZoneSettings, err
}

// UpdateNfsZoneSettings update nfs zone settings.
func UpdateNfsZoneSettings(ctx context.Context, client *client.Client, nfsZoneSettings powerscale.V2NfsSettingsZoneSettings, zone string) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv2NfsSettingsZone(ctx)
	updateParam = updateParam.V2NfsSettingsZone(nfsZoneSettings)
	updateParam = updateParam.Zone(zone)
	_, err := updateParam.Execute()
	return err
}

// FilterNfsZoneSettings filter nfs zone settings.
func FilterNfsZoneSettings(ctx context.Context, client *client.Client, filter *models.NfsZoneSettingsFilter) (*powerscale.V2NfsSettingsZone, error) {
	filterParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsSettingsZone(ctx)

	if filter != nil {
		if zoneStr := filter.Zone.ValueString(); zoneStr != "" {
			filterParam = filterParam.Zone(zoneStr)
		}
	}

	nfsZoneSettings, _, err := filterParam.Execute()
	return nfsZoneSettings, err
}
