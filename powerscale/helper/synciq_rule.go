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

	"github.com/hashicorp/terraform-plugin-framework/types"
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
