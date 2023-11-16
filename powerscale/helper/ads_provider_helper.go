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

// AdsProviderDetailMapper Does the mapping from response to model.
//
//go:noinline
func AdsProviderDetailMapper(ctx context.Context, adsProvider *powerscale.V14ProvidersAdsAdsItem) (models.AdsProviderDetailModel, error) {
	model := models.AdsProviderDetailModel{}
	err := CopyFields(ctx, adsProvider, &model)
	return model, err
}

// CreateAdsProvider Create an Ads Provider.
func CreateAdsProvider(ctx context.Context, client *client.Client, ads powerscale.V14ProvidersAdsItem) (*powerscale.CreateResponse, error) {
	adsID, _, err := client.PscaleOpenAPIClient.AuthApi.CreateAuthv14ProvidersAdsItem(ctx).V14ProvidersAdsItem(ads).Execute()
	return adsID, err
}

// GetAdsProvider retrieve Ads Provider information.
func GetAdsProvider(ctx context.Context, client *client.Client, adsModel models.AdsProviderResourceModel) (*powerscale.V14ProvidersAdsExtended, error) {
	queryParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv14ProvidersAdsById(ctx, adsModel.ID.ValueString())
	if !adsModel.Scope.IsNull() {
		queryParam = queryParam.Scope(adsModel.Scope.ValueString())
	}
	if !adsModel.CheckDuplicates.IsNull() {
		queryParam = queryParam.CheckDuplicates(adsModel.CheckDuplicates.ValueBool())
	}
	adsRes, _, err := queryParam.Execute()
	return adsRes, err
}

// UpdateAdsProvider Update an Ads Provider.
func UpdateAdsProvider(ctx context.Context, client *client.Client, adsID string, adsToUpdate powerscale.V14ProvidersAdsIdParams) error {
	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv14ProvidersAdsById(ctx, adsID)
	_, err := updateParam.V14ProvidersAdsIdParams(adsToUpdate).Execute()
	return err
}

// DeleteAdsProvider Delete an Ads Provider.
func DeleteAdsProvider(ctx context.Context, client *client.Client, adsID string) error {
	_, err := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv14ProvidersAdsById(ctx, adsID).Execute()
	return err
}

// IsCreateAdsProviderParamInvalid Verify if create params contain params only for updating.
func IsCreateAdsProviderParamInvalid(plan models.AdsProviderResourceModel) bool {
	if !plan.DomainController.IsNull() ||
		!plan.ResetSchannel.IsNull() ||
		!plan.Spns.IsUnknown() {
		return true
	}
	return false
}

// IsUpdateAdsProviderParamInvalid Verify if update params contain params only for creating.
func IsUpdateAdsProviderParamInvalid(plan models.AdsProviderResourceModel, state models.AdsProviderResourceModel) bool {
	if (!plan.DNSDomain.IsNull() && !state.DNSDomain.Equal(plan.DNSDomain)) ||
		(!plan.Instance.IsNull() && !state.Instance.Equal(plan.Instance)) ||
		(!plan.KerberosHdfsSpn.IsNull() && !state.KerberosHdfsSpn.Equal(plan.KerberosHdfsSpn)) ||
		(!plan.KerberosNfsSpn.IsNull() && !state.KerberosNfsSpn.Equal(plan.KerberosNfsSpn)) ||
		(!plan.MachineAccount.IsUnknown() && !state.MachineAccount.Equal(plan.MachineAccount)) ||
		(!plan.OrganizationalUnit.IsNull() && !state.OrganizationalUnit.Equal(plan.OrganizationalUnit)) {
		return true
	}
	return false
}

// IsGroupnetUpdated Verify if groupnet is updated after creation.
func IsGroupnetUpdated(groupnetInPlan string, groupnetInResp string) bool {
	if groupnetInPlan == "" {
		return false
	}
	return groupnetInResp != groupnetInPlan
}
