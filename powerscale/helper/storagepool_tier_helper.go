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
	"strconv"

	"fmt"
	"math"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

	var storagepoolTier powerscale.V16StoragepoolTierExtended
	for _, v := range storagepoolTiers.Tiers {
		if v.Name == toCreate.Name {
			storagepoolTier = v
		}
	}

	err = CopyFields(ctx, storagepoolTier, state)
	// Read Value from storagepoolTier and set into state
	if storagepoolTier.TransferLimitPct != nil {
		state.TransferLimitPct = types.Int64Value(int64(math.Round(*storagepoolTier.TransferLimitPct)))
	} else {
		state.TransferLimitPct = types.Int64Value(100)
	}
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

	var storagepoolTier powerscale.V16StoragepoolTierExtended
	storagepoolTier = storagepoolTiers.Tiers[0]

	err = CopyFields(ctx, storagepoolTier, state)
	if storagepoolTier.TransferLimitPct != nil {
		state.TransferLimitPct = types.Int64Value(int64(math.Round(*storagepoolTier.TransferLimitPct)))
	} else {
		state.TransferLimitPct = types.Int64Value(100)
	}
	if err != nil {
		diags.AddError(
			"Error copying fields of storagepool Tier resource",
			err.Error(),
		)
		return diags
	}

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

	var storagepoolTier = storagepoolTiers.Tiers[0]

	err = CopyFields(ctx, storagepoolTier, state)
	if storagepoolTier.TransferLimitPct != nil {
		state.TransferLimitPct = types.Int64Value(int64(math.Round(*storagepoolTier.TransferLimitPct)))
	} else {
		state.TransferLimitPct = types.Int64Value(100)
	}
	if err != nil {
		diags.AddError(
			"Error copying fields of storagepool Tier resource",
			err.Error(),
		)
		return diags
	}

	return diags
}
