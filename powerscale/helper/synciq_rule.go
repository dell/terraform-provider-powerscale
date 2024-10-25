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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetAllSyncIQRules retrieve the cluster information.
func GetAllSyncIQRules(ctx context.Context, client *client.Client) (*powerscale.V3SyncRules, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.ListSyncv3SyncRules(context.Background()).Execute()
	if err != nil {
		return resp, err
	}
	// Pagination
	for resp.Resume != "" {
		respAdd, _, errAdd := client.PscaleOpenAPIClient.SyncApi.ListSyncv3SyncRules(context.Background()).Resume(resp.Resume).Execute()
		if errAdd != nil {
			return resp, errAdd
		}
		resp.Resume = respAdd.Resume
		resp.Rules = append(resp.Rules, respAdd.Rules...)
	}
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
func UpdateSyncIQRule(ctx context.Context, client *client.Client, id string, v3SyncRule powerscale.V3SyncRule) error {
	req := powerscale.V3SyncRuleExtendedExtended{
		Description: v3SyncRule.Description,
		Enabled:     v3SyncRule.Enabled,
		Limit:       &v3SyncRule.Limit,
		Schedule:    v3SyncRule.Schedule,
	}
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv3SyncRule(context.Background(), id).V3SyncRule(req).Execute()
	return err
}

// DeleteSyncIQRule deletes SyncIQRule.
func DeleteSyncIQRule(ctx context.Context, client *client.Client, id string) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.DeleteSyncv3SyncRule(context.Background(), id).Execute()
	return err
}

// SyncIQRuleDataSourceResponse represents a generic for SyncIQRuleDataSource response.
type SyncIQRuleDataSourceResponse interface {
	powerscale.V3SyncRuleExtended | powerscale.V3SyncRuleExtendedExtendedExtended
}

// NewSyncIQRuleDataSource creates a new SyncIQRuleDataSource from datasource responses.
func NewSyncIQRuleDataSource[V SyncIQRuleDataSourceResponse](ctx context.Context, Rules []V) (*models.SyncIQRuleDataSource, error) {
	var err error
	dsRules := make([]models.SyncIQRuleModel, len(Rules))
	for i := range Rules {
		var item models.SyncIQRuleModel
		ierr := CopyFields(ctx, &Rules[i], &item)
		err = errors.Join(err, ierr)
		dsRules[i] = item
	}
	if err != nil {
		return nil, err
	}
	ret := models.SyncIQRuleDataSource{
		ID:    types.StringValue("dummy"),
		Rules: dsRules,
	}
	if len(ret.Rules) == 1 {
		ret.ID = ret.Rules[0].ID
	}
	return &ret, nil
}

var syncIQRuleResourceScheduleType = map[string]attr.Type{
	"begin": types.StringType,
	"end":   types.StringType,
	"days_of_week": types.SetType{
		ElemType: types.StringType,
	},
}

var syncIQRuleResourceType = map[string]attr.Type{
	"description": types.StringType,
	"enabled":     types.BoolType,
	"id":          types.StringType,
	"limit":       types.Int32Type,
	"schedule": types.ObjectType{
		AttrTypes: syncIQRuleResourceScheduleType,
	},
}

// NewSyncIQRulesResource creates a new SyncIQRulesResource from resource response.
func NewSyncIQRulesResource(ctx context.Context, source *powerscale.V3SyncRules) (models.SyncIQRulesResource, diag.Diagnostics) {
	var dgs diag.Diagnostics
	bw := make([]models.SyncIQRuleResource, 0)
	for _, item := range source.Rules {
		state, diags := NewSyncIQRuleResource(ctx, item)
		dgs.Append(diags...)
		switch *item.Type {
		case "bandwidth":
			bw = append(bw, state)
		}
	}
	bwList, bwListDgs := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: syncIQRuleResourceType}, bw)
	dgs.Append(bwListDgs...)
	return models.SyncIQRulesResource{
		BandWidthRules: bwList,
		ID:             types.StringValue("all"),
	}, dgs
}

// NewSyncIQRuleResource creates a new SyncIQRuleResource from resource responses.
func NewSyncIQRuleResource(ctx context.Context, source powerscale.V3SyncRuleExtended) (models.SyncIQRuleResource, diag.Diagnostics) {
	ret := models.SyncIQRuleResource{
		Description: types.StringValue(*source.Description),
		Enabled:     types.BoolValue(*source.Enabled),
		ID:          types.StringValue(*source.Id),
		Limit:       types.Int32Value(*source.Limit),
	}
	schedule := models.SyncIQRuleResourceSchedule{
		End:   source.Schedule.End,
		Begin: source.Schedule.Begin,
	}

	schedule.DaysOfWeek = unmarshalJSONSyncIQRuleschedule(source.Schedule)

	scheduleObj, dgsObj := types.ObjectValueFrom(ctx, syncIQRuleResourceScheduleType, schedule)
	ret.Schedule = scheduleObj
	return ret, dgsObj
}

// unmarshalJSONSyncIQRuleschedule converts V1SyncRuleSchedule to list of days of week
func unmarshalJSONSyncIQRuleschedule(schedule *powerscale.V1SyncRuleSchedule) []string {
	daysOfWeek := make([]string, 0)
	if schedule.Monday != nil && *schedule.Monday {
		daysOfWeek = append(daysOfWeek, "monday")
	}
	if schedule.Tuesday != nil && *schedule.Tuesday {
		daysOfWeek = append(daysOfWeek, "tuesday")
	}
	if schedule.Wednesday != nil && *schedule.Wednesday {
		daysOfWeek = append(daysOfWeek, "wednesday")
	}
	if schedule.Thursday != nil && *schedule.Thursday {
		daysOfWeek = append(daysOfWeek, "thursday")
	}
	if schedule.Friday != nil && *schedule.Friday {
		daysOfWeek = append(daysOfWeek, "friday")
	}
	if schedule.Saturday != nil && *schedule.Saturday {
		daysOfWeek = append(daysOfWeek, "saturday")
	}
	if schedule.Sunday != nil && *schedule.Sunday {
		daysOfWeek = append(daysOfWeek, "sunday")
	}
	return daysOfWeek
}

// marshalJSONSyncIQRuleschedule parses list of days of week and writes to V1SyncRuleSchedule
func marshalJSONSyncIQRuleschedule(daysOfWeek []string, schedule *powerscale.V1SyncRuleSchedule) {
	// set all values to false to start with
	schedule.Monday = New(false)
	schedule.Tuesday = New(false)
	schedule.Wednesday = New(false)
	schedule.Thursday = New(false)
	schedule.Friday = New(false)
	schedule.Saturday = New(false)
	schedule.Sunday = New(false)
	// set specified values to false
	for _, day := range daysOfWeek {
		switch day {
		case "monday":
			schedule.Monday = New(true)
		case "tuesday":
			schedule.Tuesday = New(true)
		case "wednesday":
			schedule.Wednesday = New(true)
		case "thursday":
			schedule.Thursday = New(true)
		case "friday":
			schedule.Friday = New(true)
		case "saturday":
			schedule.Saturday = New(true)
		case "sunday":
			schedule.Sunday = New(true)
		}
	}
}

// GetRequestsFromSynciqRulesResource converts SyncIQRulesResource to SyncIQRulesResourceRequest
func GetRequestsFromSynciqRulesResource(ctx context.Context, source models.SyncIQRulesResource) models.SyncIQRulesResourceRequest {
	ret := models.SyncIQRulesResourceRequest{
		BandWidthRules: make([]models.SyncIQRuleResource, 0),
	}
	source.BandWidthRules.ElementsAs(ctx, &ret.BandWidthRules, true)
	return ret
}

// GetRequestFromSynciqRuleResource creates a new SyncIQRule API request from resource plan.
func GetRequestFromSynciqRuleResource(ctx context.Context, plan models.SyncIQRuleResource, ruleType string) powerscale.V3SyncRule {
	ret := powerscale.V3SyncRule{
		Type:        ruleType,
		Limit:       plan.Limit.ValueInt32(),
		Description: GetKnownStringPointer(plan.Description),
		Enabled:     GetKnownBoolPointer(plan.Enabled),
	}
	if plan.Schedule.IsUnknown() || plan.Schedule.IsNull() {
		return ret
	}
	var schedule models.SyncIQRuleResourceSchedule
	plan.Schedule.As(ctx, &schedule, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	ret.Schedule = &powerscale.V1SyncRuleSchedule{
		Begin: schedule.Begin,
		End:   schedule.End,
	}
	if schedule.DaysOfWeek == nil {
		return ret
	}
	marshalJSONSyncIQRuleschedule(schedule.DaysOfWeek, ret.Schedule)
	return ret
}
