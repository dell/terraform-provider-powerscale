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
	"errors"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllSyncIQPolicies retrieve the cluster information.
func GetAllSyncIQPolicies(ctx context.Context, client *client.Client) (*powerscale.V14SyncPolicies, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.ListSyncv14SyncPolicies(context.Background()).Execute()
	return resp, err
}

// GetSyncIQPolicyByID retrieve the cluster information.
func GetSyncIQPolicyByID(ctx context.Context, client *client.Client, id string) (*powerscale.V14SyncPoliciesExtended, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv14SyncPolicy(context.Background(), id).Execute()
	return resp, err
}

// SyncIQPolicyDataSourceResponse is the union of all response types for syncIQ policy datasource.
type SyncIQPolicyDataSourceResponse interface {
	powerscale.V14SyncPolicyExtended | powerscale.V14SyncPolicyExtendedExtendedExtended
}

// NewSyncIQPolicyDataSource creates a new SyncIQPolicyDataSource from datasource responses.
func NewSyncIQPolicyDataSource[V SyncIQPolicyDataSourceResponse](ctx context.Context, policies []V) (*models.SyncIQPolicyDataSource, error) {
	var err error
	ret := models.SyncIQPolicyDataSource{
		ID:       types.StringValue("dummy"),
		Policies: make([]models.V14SyncPolicyExtendedModel, len(policies)),
	}
	for i := range policies {
		var item models.V14SyncPolicyExtendedModel
		ierr := CopyFields(ctx, &policies[i], &item)
		err = errors.Join(err, ierr)
		ret.Policies[i] = item
	}
	if len(ret.Policies) == 1 {
		ret.ID = ret.Policies[0].ID
	}
	return &ret, err
}
