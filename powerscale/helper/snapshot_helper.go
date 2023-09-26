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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllSnapshots returns the full list of snapshots.
func GetAllSnapshots(ctx context.Context, client *client.Client) ([]powerscale.V1SnapshotSnapshotExtended, error) {
	result, _, err := client.PscaleOpenAPIClient.SnapshotApi.ListSnapshotv1SnapshotSnapshots(ctx).Execute()
	return result.GetSnapshots(), err
}

// SnapshotDetailMapper Does the mapping from response to model.
func SnapshotDetailMapper(ctx context.Context, snap powerscale.V1SnapshotSnapshotExtended) (models.SnapshotDetailModel, error) {
	model := models.SnapshotDetailModel{}
	err := CopyFields(ctx, &snap, &model)
	if err != nil {
		return model, err
	}
	model.ID = types.Int64Value(int64(snap.Id))
	model.TargetID = types.Int64Value(int64(snap.TargetId))
	return model, nil
}
