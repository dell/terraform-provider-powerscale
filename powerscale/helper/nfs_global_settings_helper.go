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

// GetNfsGlobalSettings retrieve nfs global settings.
func GetNfsGlobalSettings(ctx context.Context, client *client.Client) (*powerscale.V12NfsSettingsGlobal, error) {
	nfsGlobalSettings, _, err := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv12NfsSettingsGlobal(ctx).Execute()
	return nfsGlobalSettings, err
}

// UpdateNfsGlobalSettings update nfs global settings.
func UpdateNfsGlobalSettings(ctx context.Context, client *client.Client, v12NfsSettingsGlobal powerscale.V12NfsSettingsGlobalSettings) error {
	_, err := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv12NfsSettingsGlobal(ctx).V12NfsSettingsGlobal(v12NfsSettingsGlobal).Execute()
	return err
}
