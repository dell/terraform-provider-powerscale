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
