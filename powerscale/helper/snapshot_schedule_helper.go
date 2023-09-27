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
)

// ListSnapshotSchedules lists the snapshot schedules.
func ListSnapshotSchedules(ctx context.Context, client *client.Client, ssFilter *models.SnapshotScheduleFilter) ([]powerscale.V1SnapshotScheduleExtended, error) {
	listSsParam := client.PscaleOpenAPIClient.SnapshotApi.ListSnapshotv1SnapshotSchedules(ctx)
	if ssFilter != nil {
		if !ssFilter.Sort.IsNull() {
			listSsParam = listSsParam.Sort(ssFilter.Sort.ValueString())
		}
		if !ssFilter.Dir.IsNull() {
			listSsParam = listSsParam.Dir(ssFilter.Dir.ValueString())
		}
		if !ssFilter.Limit.IsNull() {
			listSsParam = listSsParam.Limit(int32(ssFilter.Limit.ValueInt64()))
		}
	}
	snapshotSchedules, _, err := listSsParam.Execute()
	if err != nil {
		return nil, err
	}
	return snapshotSchedules.Schedules, nil
}
