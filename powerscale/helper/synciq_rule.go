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
	"errors"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllSyncIQRules retrieve the cluster information.
func GetAllSyncIQRules(ctx context.Context, client *client.Client) (*powerscale.V3SyncRules, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.ListSyncv3SyncRules(context.Background()).Execute()
	return resp, err
}

// GetSyncIQRuleByID retrieve the cluster information.
func GetSyncIQRuleByID(ctx context.Context, client *client.Client, id string) (*powerscale.V3SyncRulesExtended, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv3SyncRule(context.Background(), id).Execute()
	return resp, err
}

// CreateSyncIQRule creates SyncIQRule.
func CreateSyncIQRule(ctx context.Context, client *client.Client, v3SyncRule powerscale.V3SyncRule) (string, error) {
	respC, _, err := client.PscaleOpenAPIClient.SyncApi.CreateSyncv3SyncRule(context.Background()).V3SyncRule(v3SyncRule).Execute()
	if err != nil {
		return "", err
	}
	return respC.Id, nil
}

// UpdateSyncIQRule updates SyncIQRule.
func UpdateSyncIQRule(ctx context.Context, client *client.Client, id string, v3SyncRule powerscale.V3SyncRuleExtendedExtended) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv3SyncRule(context.Background(), id).V3SyncRule(v3SyncRule).Execute()
	return err
}

// SyncIQRuleDataSourceResponse represents a generic for SyncIQRuleDataSource response.
type SyncIQRuleDataSourceResponse interface {
	powerscale.V3SyncRuleExtended | powerscale.V3SyncRuleExtendedExtendedExtended
}

// NewSyncIQRuleDataSource creates a new SyncIQRuleDataSource from datasource responses.
func NewSyncIQRuleDataSource[V SyncIQRuleDataSourceResponse](ctx context.Context, Rules []V) (*models.SyncIQRuleDataSource, error) {
	var err error
	ret := models.SyncIQRuleDataSource{
		ID:    types.StringValue("dummy"),
		Rules: make([]models.SyncIQRuleModel, len(Rules)),
	}
	for i := range Rules {
		var item models.SyncIQRuleModel
		ierr := CopyFields(ctx, &Rules[i], &item)
		err = errors.Join(err, ierr)
		ret.Rules[i] = item
	}
	if len(ret.Rules) == 1 {
		ret.ID = ret.Rules[0].ID
	}
	return &ret, err
}

// NewSyncIQRuleDataSource creates a new SyncIQRuleDataSource from datasource responses.
func NewSyncIQRuleResource(ctx context.Context, source powerscale.V3SyncRuleExtendedExtendedExtended) (models.SyncIQRuleResource, diag.Diagnostics) {
	ret := models.SyncIQRuleResource{
		Type:        types.StringValue(source.Type),
		Description: types.StringValue(source.Description),
		Enabled:     types.BoolValue(source.Enabled),
		ID:          types.StringValue(source.Id),
		Limit:       types.Int64Value(int64(source.Limit)),
	}
	schedule := models.SyncIQRuleResourceSchedule{
		End:   types.StringPointerValue(source.Schedule.End),
		Begin: types.StringPointerValue(source.Schedule.Begin),
	}

	daysOfWeek := make([]string, 0)
	if source.Schedule.Monday != nil && *source.Schedule.Monday {
		daysOfWeek = append(daysOfWeek, "monday")
	}
	if source.Schedule.Tuesday != nil && *source.Schedule.Tuesday {
		daysOfWeek = append(daysOfWeek, "tuesday")
	}
	if source.Schedule.Wednesday != nil && *source.Schedule.Wednesday {
		daysOfWeek = append(daysOfWeek, "wednesday")
	}
	if source.Schedule.Thursday != nil && *source.Schedule.Thursday {
		daysOfWeek = append(daysOfWeek, "thursday")
	}
	if source.Schedule.Friday != nil && *source.Schedule.Friday {
		daysOfWeek = append(daysOfWeek, "friday")
	}
	if source.Schedule.Saturday != nil && *source.Schedule.Saturday {
		daysOfWeek = append(daysOfWeek, "saturday")
	}
	if source.Schedule.Sunday != nil && *source.Schedule.Sunday {
		daysOfWeek = append(daysOfWeek, "sunday")
	}
	daysOfWeekSet, dgs := types.SetValueFrom(ctx, types.StringType, daysOfWeek)
	schedule.DaysOfWeek = daysOfWeekSet

	scheduleObj, dgsObj := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"begin": types.StringType,
		"end":   types.StringType,
		"days_of_thr_week": types.SetType{
			ElemType: types.StringType,
		},
	}, schedule)
	dgs.Append(dgsObj...)
	ret.Schedule = scheduleObj
	return ret, dgs
}
