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

// NtpSettingsDetailMapper Does the mapping from response to model.
//
//go:noinline
func NtpSettingsDetailMapper(ctx context.Context, ntpSettings *powerscale.V3NtpSettingsSettings) (models.NtpSettingsDataSourceModel, error) {
	model := models.NtpSettingsDataSourceModel{}
	err := CopyFields(ctx, ntpSettings, &model)
	return model, err
}

// GetNtpSettings retrieve NTP Settings information.
func GetNtpSettings(ctx context.Context, client *client.Client) (*powerscale.V3NtpSettings, error) {
	queryParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv3NtpSettings(ctx)
	ntpSettingsRes, _, err := queryParam.Execute()
	return ntpSettingsRes, err
}

// UpdateNtpSettings Update NTP Settings.
func UpdateNtpSettings(ctx context.Context, client *client.Client, ntpSettingsToUpdate powerscale.V3NtpSettingsSettings) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv3NtpSettings(ctx)
	_, err := updateParam.V3NtpSettings(ntpSettingsToUpdate).Execute()
	return err
}

// ResolveSettingsDiff implement state
// For nfs export settings info, response may only contain UID while type/username is given
// Need to manually copy plan info to state, or state would keep the type/username as null, which is inconsistent.
func ResolveSettingsDiff(ctx context.Context, plan models.NfsexportsettingsModel, state *models.NfsexportsettingsModel) {
	state.MapAll = assignKnownObjectToUnknown(ctx, plan.MapAll, state.MapAll)
	state.MapFailure = assignKnownObjectToUnknown(ctx, plan.MapFailure, state.MapFailure)
	state.MapNonRoot = assignKnownObjectToUnknown(ctx, plan.MapNonRoot, state.MapNonRoot)
	state.MapRoot = assignKnownObjectToUnknown(ctx, plan.MapRoot, state.MapRoot)
	state.SecurityFlavors = ListCheck(plan.SecurityFlavors, plan.SecurityFlavors.ElementType(ctx))
}
