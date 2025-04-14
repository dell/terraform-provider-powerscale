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
	"fmt"

	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type V16StoragepoolTierExtended struct {
	// The names or IDs of the tier's children.
	Children []string `json:"children,omitempty"`
	// The system ID given to the tier.
	Id int32 `json:"id"`
	// The nodes that are part of this tier.
	Lnns []int32 `json:"lnns"`
	// The tier name.
	Name string `json:"name"`
	// Stop moving files to this tier when this limit is met
	TransferLimitPct int32 `json:"transfer_limit_pct,omitempty"`
	// How the transfer limit value is being applied
	TransferLimitState *string `json:"transfer_limit_state,omitempty"`
}

// GetStoragepoolTier gets storage pool tier.
func GetStoragepoolTier(ctx context.Context, client *client.Client, path string) (*powerscale.V16StoragepoolTiersExtended, error) {
	StoragepoolTier, _, err := client.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv16StoragepoolTier(ctx, path).Execute()
	return StoragepoolTier, err
}

// GetAllStoragepoolTiers returns the full list of storage pool tiers.
func GetAllStoragepoolTiers(ctx context.Context, client *client.Client) ([]powerscale.V16StoragepoolTierExtended, error) {
	StoragepoolTierParams := client.PscaleOpenAPIClient.StoragepoolApi.ListStoragepoolv16StoragepoolTiers(ctx)

	StoragepoolTiers, _, err := StoragepoolTierParams.Execute()
	if err != nil {
		errStr := constants.ReadStoragepoolTiersErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting storagepool tiers: %s", message)
	}
	return StoragepoolTiers.Tiers, nil
}

// CreateStoragepoolTier created the storagepool tier on the array.
func CreateStoragepoolTier(ctx context.Context, client *client.Client, plan models.StoragepoolTierResourceModel, state *models.StoragepoolTierResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var toCreate powerscale.V16StoragepoolTier

	zoneParam := client.PscaleOpenAPIClient.StoragepoolApi.CreateStoragepoolv16StoragepoolTier(ctx)
	err := ReadFromState(ctx, &plan, &toCreate)
	if err != nil {
		errStr := constants.CreateStoragepoolTierErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error creating storagepool tier",
			fmt.Sprintf("Could not create storagepool tier with error: %s", message),
		)
		return diags
	}
	_, _, err2 := zoneParam.V16StoragepoolTier(toCreate).Execute()
	if err2 != nil {
		errStr := constants.CreateStoragepoolTierErrorMsg + "with error: "
		message := GetErrorString(err2, errStr)
		diags.AddError(
			"Error creating storagepool tier",
			message,
		)
		return diags
	}

	checkParam := client.PscaleOpenAPIClient.StoragepoolApi.ListStoragepoolv16StoragepoolTiers(ctx)
	storagepoolTiers, _, err := checkParam.Execute()
	if err != nil {
		errStr := constants.ReadStoragepoolTierErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading storagepool tier",
			message,
		)
		return diags
	}

	var storagepoolTier V16StoragepoolTierExtended
	for _, v := range storagepoolTiers.Tiers {
		if v.Name == toCreate.Name {
			var i int32
			if v.TransferLimitPct != nil {
				i = int32(*v.TransferLimitPct)
			} else {
				i = 100
			}
			storagepoolTier = V16StoragepoolTierExtended{
				Name:               v.Name,
				TransferLimitState: v.TransferLimitState,
				TransferLimitPct:   i,
				Children:           v.Children,
				Id:                 v.Id,
				Lnns:               v.Lnns,
			}
		}
	}

	err = CopyFields(ctx, storagepoolTier, state)
	state.TransferLimitPct = types.Int32Value(storagepoolTier.TransferLimitPct)
	if err != nil {
		diags.AddError(
			"Error copying fields of storagepool tier resource",
			err.Error(),
		)
		return diags
	}
	return diags
}

// UpdateStoragepoolTier updates a particular storagepool tier
func UpdateStoragepoolTier(ctx context.Context, client *client.Client, editValues powerscale.V16StoragepoolTierExtendedExtended, state *models.StoragepoolTierResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	editParam := client.PscaleOpenAPIClient.StoragepoolApi.UpdateStoragepoolv16StoragepoolTier(ctx, strconv.FormatInt(state.Id.ValueInt64(), 10))
	editParam = editParam.V16StoragepoolTier(editValues)

	_, err := editParam.Execute()
	if err != nil {
		errStr := constants.UpdateStoragepoolTierErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error updating storagepool tier",
			message,
		)
	}

	return diags
}

// ReadStoragepoolTier reads a particular storagepool tier from the list of storagepool tiers
func ReadStoragepoolTier(ctx context.Context, client *client.Client, state *models.StoragepoolTierResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	checkParam := client.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv16StoragepoolTier(ctx, strconv.FormatInt(state.Id.ValueInt64(), 10))
	storagepoolTiers, _, err := checkParam.Execute()
	if err != nil {
		errStr := constants.ReadStoragepoolTierErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading storagepool tier",
			message,
		)
		return diags
	}

	var v = storagepoolTiers.Tiers[0]
	var i int32
	if v.TransferLimitPct != nil {
		i = int32(*v.TransferLimitPct)
	} else {
		i = 100
	}
	storagepoolTier := V16StoragepoolTierExtended{
		Name:               v.Name,
		TransferLimitState: v.TransferLimitState,
		TransferLimitPct:   i,
		Children:           v.Children,
		Id:                 v.Id,
		Lnns:               v.Lnns,
	}

	err = CopyFields(ctx, storagepoolTier, state)
	if err != nil {
		diags.AddError(
			"Error copying fields of storagepool Tier resource",
			err.Error(),
		)
		return diags
	}
	state.TransferLimitPct = types.Int32Value(storagepoolTier.TransferLimitPct)
	return diags
}

// GetStoragepoolTierByID retrieves storagepool tier by id.
func GetStoragepoolTierByID(ctx context.Context, client *client.Client, id string, state *models.StoragepoolTierResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	checkParam := client.PscaleOpenAPIClient.StoragepoolApi.GetStoragepoolv16StoragepoolTier(ctx, id)
	storagepoolTiers, _, err := checkParam.Execute()
	if err != nil {
		errStr := constants.ReadStoragepoolTierErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading storagepool tier",
			message,
		)
		return diags
	}

	var v = storagepoolTiers.Tiers[0]
	var i int32
	if v.TransferLimitPct != nil {
		i = int32(*v.TransferLimitPct)
	} else {
		i = 100
	}
	storagepoolTier := V16StoragepoolTierExtended{
		Name:               v.Name,
		TransferLimitState: v.TransferLimitState,
		TransferLimitPct:   i,
		Children:           v.Children,
		Id:                 v.Id,
		Lnns:               v.Lnns,
	}

	err = CopyFields(ctx, storagepoolTier, state)
	state.TransferLimitPct = types.Int32Value(storagepoolTier.TransferLimitPct)
	if err != nil {
		diags.AddError(
			"Error copying fields of storagepool Tier resource",
			err.Error(),
		)
		return diags
	}

	return diags
}
