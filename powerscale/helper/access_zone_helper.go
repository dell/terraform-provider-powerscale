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
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AccessZoneDetailMapper Does the mapping from response to model.
func AccessZoneDetailMapper(ctx context.Context, az *powerscale.V3ZoneExtended) (models.AccessZoneDetailModel, error) {
	model := models.AccessZoneDetailModel{}
	err := CopyFields(ctx, az, &model)
	model.IfsRestricted, _ = GetAuthAccessKeyObjects(az.IfsRestricted)
	// These need to be done manually because of the linter
	model.ZoneID = types.Int64Value(int64(*az.ZoneId))
	model.ID = types.StringValue(*az.Id)
	if err != nil {
		return model, err
	}
	return model, nil
}

// AccessZoneResourceDetailMapper detail mapper for access zone resource.
func AccessZoneResourceDetailMapper(ctx context.Context, az *powerscale.V3ZoneExtended) (models.AccessZoneResourceModel, error) {
	model := models.AccessZoneResourceModel{}
	err := CopyFields(ctx, az, &model)
	model.IfsRestricted, _ = GetAuthAccessKeyObjects(az.IfsRestricted)
	// These need to be done manually because of the linter
	model.ZoneID = types.Int64Value(int64(*az.ZoneId))
	model.ID = types.StringValue(*az.Id)
	if err != nil {
		return model, err
	}
	return model, nil
}

// GetAuthAccessKeyObjects returns auth hours key objects.
func GetAuthAccessKeyObjects(accessResponse []powerscale.V1AuthAccessAccessItemFileGroup) (types.List, diag.Diagnostics) {
	var accessKeyObjects []attr.Value
	accessType := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	}
	for _, access := range accessResponse {
		accessKeyObjects = append(accessKeyObjects, AuthAccessKeyObjectMapper(access))
	}

	return types.ListValue(types.ObjectType{AttrTypes: accessType}, accessKeyObjects)
}

// AuthAccessKeyObjectMapper parses V1AuthAccessAccessItemFileGroup to object.
func AuthAccessKeyObjectMapper(access powerscale.V1AuthAccessAccessItemFileGroup) types.Object {
	accessType := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	}
	accessMap := make(map[string]attr.Value)
	accessMap["id"] = types.StringValue(*access.Id)
	accessMap["name"] = types.StringValue(*access.Name)
	accessMap["type"] = types.StringValue(*access.Type)
	accessObject, _ := types.ObjectValue(accessType, accessMap)

	return accessObject
}

// GetAllAccessZones returns the full list of access zones.
func GetAllAccessZones(ctx context.Context, client *client.Client) (*powerscale.V3Zones, error) {
	result, _, err := client.PscaleOpenAPIClient.ZonesApi.ListZonesv3Zones(ctx).Execute()
	return result, err
}

// CreateAccessZones Creates an Access Zone.
func CreateAccessZones(ctx context.Context, client *client.Client, authProv []string, plan *models.AccessZoneResourceModel) error {
	forceOverlap := true
	createPath := false
	createParam := client.PscaleOpenAPIClient.ZonesApi.CreateZonesv3Zone(ctx)
	createParam = createParam.V3Zone(powerscale.V3Zone{
		AuthProviders: authProv,
		CreatePath:    &createPath,
		Groupnet:      plan.Groupnet.ValueStringPointer(),
		Name:          plan.Name.ValueString(),
		Path:          plan.Path.ValueStringPointer(),
		ForceOverlap:  &forceOverlap,
	})

	_, _, err := createParam.Execute()
	return err
}

// GetSpecificZone returns a specific zone or an error.
func GetSpecificZone(ctx context.Context, matchZone string, zoneList []powerscale.V3ZoneExtended) (models.AccessZoneResourceModel, error) {
	for _, vze := range zoneList {
		if *vze.Name == matchZone {
			zone := vze
			state, err := AccessZoneResourceDetailMapper(ctx, &zone)
			if err != nil {
				errStr := constants.ReadAccessZoneErrorMsg + "with error: "
				message := GetErrorString(err, errStr)
				return models.AccessZoneResourceModel{}, fmt.Errorf("error finding new access zone after create : %s", message)

			}
			return state, nil
		}
	}

	return models.AccessZoneResourceModel{}, fmt.Errorf("error finding new access zone after create, Unable to create successfully")
}

// ExtractCustomAuthForInput extracts the custom auth provider from actual auth provider for input.
func ExtractCustomAuthForInput(ctx context.Context, stateResponse, plan models.AccessZoneResourceModel) basetypes.ListValue {

	automaticProviderAdded := false
	automaticProviderName := "lsa-local-provider:" + plan.Name.ValueString()
	var customAuthProviders []attr.Value
	for _, v := range plan.CustomAuthProviders.Elements() {
		name := strings.Trim(v.String(), "\"")
		if !strings.Contains(name, ":") {
			name = "lsa-local-provider:" + name
		}
		if name == automaticProviderName {
			automaticProviderAdded = true
		}
		customAuthProviders = append(customAuthProviders, types.StringValue(name))
	}
	if !automaticProviderAdded {
		customAuthProviders = append(customAuthProviders, types.StringValue(automaticProviderName))
	}

	authProviders := stateResponse.AuthProviders.Elements()
	if len(customAuthProviders) != len(authProviders) {
		return stateResponse.AuthProviders
	}
	for i, v := range authProviders {
		if !customAuthProviders[i].Equal(v) {
			return stateResponse.AuthProviders
		}
	}

	return plan.CustomAuthProviders
}

// QueryZoneNameByID returns a specific zone id by name.
func QueryZoneNameByID(ctx context.Context, client *client.Client, zoneID int64) (string, error) {
	zones, err := GetAllAccessZones(ctx, client)
	if err != nil {
		return "", err
	}
	for _, zone := range zones.Zones {
		if int64(*zone.ZoneId) == int64(zoneID) {
			return *zone.Name, nil
		}
	}

	return "", fmt.Errorf("error finding zone name for zone ID %d", zoneID)
}
