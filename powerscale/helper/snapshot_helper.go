/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllSnapshots returns the full list of snapshots.
func GetAllSnapshots(ctx context.Context, client *client.Client) ([]powerscale.V1SnapshotSnapshotExtended, error) {
	result, _, err := client.PscaleOpenAPIClient.SnapshotApi.ListSnapshotv1SnapshotSnapshots(ctx).Execute()
	return result.GetSnapshots(), err
}

// GetSpecificSnapshot returns a specific snapshot based on the id
func GetSpecificSnapshot(ctx context.Context, client *client.Client, id string) (powerscale.Createv1SnapshotSnapshotResponse, error) {
	snap := powerscale.Createv1SnapshotSnapshotResponse{}
	result, _, err := client.PscaleOpenAPIClient.SnapshotApi.GetSnapshotv1SnapshotSnapshot(ctx, id).Execute()
	if result != nil && len(result.Snapshots) > 0 {
		snap = result.Snapshots[0]
	}
	return snap, err
}

// ModifySnapshot returns the full list of snapshots.
func ModifySnapshot(ctx context.Context, client *client.Client, id string, edit powerscale.V1SnapshotSnapshotExtendedExtended) error {
	updateParam := client.PscaleOpenAPIClient.SnapshotApi.UpdateSnapshotv1SnapshotSnapshot(ctx, id)
	updateParam = updateParam.V1SnapshotSnapshot(edit)
	_, err := updateParam.Execute()
	return err
}

// CreateSnapshot returns the full list of snapshots.
func CreateSnapshot(ctx context.Context, client *client.Client, plan *models.SnapshotDetailModel) (powerscale.Createv1SnapshotSnapshotResponse, error) {
	expire := CalclulateExpire(plan.SetExpires.ValueString())
	nameDefault := time.Now().String()
	// Path should always be set
	// Name should default to current date if unset
	createBody := powerscale.V1SnapshotSnapshot{
		Path: plan.Path.ValueString(),
		Name: &nameDefault,
	}
	// Only set if not dont expire
	if expire != 0 {
		createBody.Expires = &expire
	}

	// Only set if not dont expire
	if plan.Name.ValueString() != "" {
		createBody.Name = plan.Name.ValueStringPointer()
	}
	createParam := client.PscaleOpenAPIClient.SnapshotApi.CreateSnapshotv1SnapshotSnapshot(ctx)
	createParam = createParam.V1SnapshotSnapshot(createBody)
	result, _, err := createParam.Execute()
	return *result, err
}

// SnapshotDetailMapper Does the mapping from response to model.
func SnapshotDetailMapper(ctx context.Context, snap powerscale.V1SnapshotSnapshotExtended) (models.SnapshotDetailModel, error) {
	model := models.SnapshotDetailModel{}
	err := CopyFields(ctx, &snap, &model)
	if err != nil {
		return model, err
	}
	model.ID = types.StringValue(fmt.Sprint(snap.Id))
	model.TargetID = types.Int64Value(int64(snap.TargetId))
	model.SetExpires = types.StringNull()
	return model, nil
}

// SnapshotResourceDetailMapper Does the mapping from response to model.
func SnapshotResourceDetailMapper(ctx context.Context, snap powerscale.Createv1SnapshotSnapshotResponse) (models.SnapshotDetailModel, error) {
	model := models.SnapshotDetailModel{}
	err := CopyFields(ctx, &snap, &model)
	if err != nil {
		return model, err
	}
	model.ID = types.StringValue(fmt.Sprint(snap.Id))
	model.TargetID = types.Int64Value(int64(snap.TargetId))

	return model, nil
}

// CalclulateExpire Calculates the Unix Epic based on 1 day, 1 week or 1 month from the current date and time.
func CalclulateExpire(setExpireValue string) int32 {
	expireTime := time.Now().Unix()
	// 86400 is the Epoch day in seconds
	switch setExpireValue {
	case "Never":
		expireTime = 0
	case "1 Day":
		expireTime = expireTime + 86400
	case "1 Week":
		expireTime = expireTime + (86400 * 7)
	case "1 Month":
		expireTime = expireTime + (86400 * 30)
	}
	return int32(expireTime)
}
