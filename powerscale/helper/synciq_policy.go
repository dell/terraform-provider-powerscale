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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllSyncIQPolicies retrieve the cluster information.
func GetAllSyncIQPolicies(ctx context.Context, client *client.Client) (*powerscale.V14SyncPolicies, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.ListSyncv14SyncPolicies(context.Background()).Execute()
	if err != nil {
		return resp, err
	}
	for resp.Resume != nil {
		respAdd, _, errAdd := client.PscaleOpenAPIClient.SyncApi.ListSyncv14SyncPolicies(context.Background()).Resume(*resp.Resume).Execute()
		if errAdd != nil {
			return resp, errAdd
		}
		resp.Resume = respAdd.Resume
		resp.Policies = append(resp.Policies, respAdd.Policies...)
	}
	return resp, err
}

// GetSyncIQPolicyIDByName retrieve the cluster information.
func GetSyncIQPolicyIDByName(ctx context.Context, client *client.Client, name string) (string, error) {
	policies, err := GetAllSyncIQPolicies(ctx, client)
	if err != nil {
		errStr := "Could not get list of SyncIQ policies with error: "
		message := GetErrorString(err, errStr)
		return "", errors.New(message)
	}
	for _, policy := range policies.Policies {
		if policy.Name == name {
			return policy.Id, nil
		}
	}
	return "", fmt.Errorf("policy by name %s not found", name)
}

// GetSyncIQPolicyByID retrieve the cluster information.
func GetSyncIQPolicyByID(ctx context.Context, client *client.Client, id string) (*powerscale.V14SyncPoliciesExtended, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv14SyncPolicy(context.Background(), id).Execute()
	return resp, err
}

// CreateSyncIQPolicy creates the sync iq policy.
func CreateSyncIQPolicy(ctx context.Context, client *client.Client, policy powerscale.V14SyncPolicy) (string, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.CreateSyncv14SyncPolicy(ctx).V14SyncPolicy(policy).Execute()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// DeleteSyncIQPolicy deletes the sync iq policy.
func DeleteSyncIQPolicy(ctx context.Context, client *client.Client, id string) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.DeleteSyncv14SyncPolicy(ctx, id).Execute()
	return err
}

// UpdateSyncIQPolicy updates the sync iq policy.
func UpdateSyncIQPolicy(ctx context.Context, client *client.Client, id string, policy powerscale.V14SyncPolicyExtendedExtended) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv14SyncPolicy(ctx, id).V14SyncPolicy(policy).Execute()
	return err
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

// NewSynciqpolicyResourceModel creates a new SynciqpolicyResourceModel from resource read response.
func NewSynciqpolicyResourceModel(ctx context.Context, respR *powerscale.V14SyncPoliciesExtended) (models.SynciqpolicyResourceModel, diag.Diagnostics) {
	var state models.SynciqpolicyResourceModel
	var dgs diag.Diagnostics
	source := respR.Policies[0]

	// rpo value of zero and null both mean no rpo alert
	// so better convert null to zero
	if source.RpoAlert == nil {
		source.RpoAlert = new(int32)
	}

	// same with bandwidth reservation
	if source.BandwidthReservation == nil {
		source.BandwidthReservation = new(int32)
	}

	// same for skip lookup
	if source.SkipLookup == nil {
		source.SkipLookup = new(bool)
	}

	if source.FileMatchingPattern != nil && len(source.FileMatchingPattern.OrCriteria) == 0 {
		source.FileMatchingPattern = nil
	}

	if source.SourceNetwork == nil {
		source.SourceNetwork = &powerscale.V1SyncPolicySourceNetwork{}
	}

	// set these lists from null to zero
	if source.SourceIncludeDirectories == nil {
		source.SourceIncludeDirectories = make([]string, 0)
	}

	if source.SourceExcludeDirectories == nil {
		source.SourceExcludeDirectories = make([]string, 0)
	}

	if source.LinkedServicePolicies == nil {
		source.LinkedServicePolicies = make([]string, 0)
	}

	err := CopyFieldsToNonNestedModel(ctx, source, &state)
	if err != nil {
		dgs.AddError(
			"Error copying fields of SyncIQ Policy resource",
			err.Error(),
		)
		return state, dgs
	}

	return state, nil
}
