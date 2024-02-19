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
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateLdapProviderResourceState updates resource state.
func UpdateLdapProviderResourceState(ctx context.Context, ldapProviderModel *models.LdapProviderModel, ldapProviderResponse any) (err error) {
	originModel := *ldapProviderModel

	// tls_revocation_check_level and ocsp_server_uris are available in 9.5
	ldapProviderModel.TLSRevocationCheckLevel = types.StringNull()
	ldapProviderModel.OcspServerUris = types.ListNull(types.StringType)

	if err = CopyFields(ctx, ldapProviderResponse, ldapProviderModel); err != nil {
		return
	}

	if len(originModel.FindableGroups.Elements()) != 0 && IsListValueEquals(originModel.FindableGroups, ldapProviderModel.FindableGroups) {
		ldapProviderModel.FindableGroups = originModel.FindableGroups
	}
	if len(originModel.FindableUsers.Elements()) != 0 && IsListValueEquals(originModel.FindableUsers, ldapProviderModel.FindableUsers) {
		ldapProviderModel.FindableUsers = originModel.FindableUsers
	}
	if len(originModel.ListableGroups.Elements()) != 0 && IsListValueEquals(originModel.ListableGroups, ldapProviderModel.ListableGroups) {
		ldapProviderModel.ListableGroups = originModel.ListableGroups
	}
	if len(originModel.ListableUsers.Elements()) != 0 && IsListValueEquals(originModel.ListableUsers, ldapProviderModel.ListableUsers) {
		ldapProviderModel.ListableUsers = originModel.ListableUsers
	}
	if len(originModel.ServerUris.Elements()) != 0 && IsListValueEquals(originModel.ServerUris, ldapProviderModel.ServerUris) {
		ldapProviderModel.ServerUris = originModel.ServerUris
	}
	if len(originModel.UnfindableGroups.Elements()) != 0 && IsListValueEquals(originModel.UnfindableGroups, ldapProviderModel.UnfindableGroups) {
		ldapProviderModel.UnfindableGroups = originModel.UnfindableGroups
	}
	if len(originModel.UnfindableUsers.Elements()) != 0 && IsListValueEquals(originModel.UnfindableUsers, ldapProviderModel.UnfindableUsers) {
		ldapProviderModel.UnfindableUsers = originModel.UnfindableUsers
	}
	if len(originModel.UnlistableGroups.Elements()) != 0 && IsListValueEquals(originModel.UnlistableGroups, ldapProviderModel.UnlistableGroups) {
		ldapProviderModel.UnlistableGroups = originModel.UnlistableGroups
	}
	if len(originModel.UnlistableUsers.Elements()) != 0 && IsListValueEquals(originModel.UnlistableUsers, ldapProviderModel.UnlistableUsers) {
		ldapProviderModel.UnlistableUsers = originModel.UnlistableUsers
	}
	if len(originModel.OcspServerUris.Elements()) != 0 && IsListValueEquals(originModel.OcspServerUris, ldapProviderModel.OcspServerUris) {
		ldapProviderModel.OcspServerUris = originModel.OcspServerUris
	}

	return
}

// UpdateLdapProviderDataSourceState updates datasource state.
func UpdateLdapProviderDataSourceState(ctx context.Context, ldapProviderModel *models.LdapProviderDataSourceModel, ldapProviderListResponse any) (err error) {

	switch v := ldapProviderListResponse.(type) {
	case *powerscale.V16ProvidersLdap:
		s, ok := ldapProviderListResponse.(*powerscale.V16ProvidersLdap)
		if !ok {
			return fmt.Errorf("error pararmeter LdapProviderDataSourceModel - Unexpected type: %T", v)
		}
		ldapProviderModel.LdapProviders = make([]models.LdapProviderDetailModel, 0)
		for _, ldapProvider := range s.GetLdap() {
			var model models.LdapProviderDetailModel
			if err = CopyFields(ctx, ldapProvider, &model); err != nil {
				return
			}
			ldapProviderModel.LdapProviders = append(ldapProviderModel.LdapProviders, model)
		}

	case *powerscale.V11ProvidersLdap:
		s, ok := ldapProviderListResponse.(*powerscale.V11ProvidersLdap)
		if !ok {
			return fmt.Errorf("error pararmeter LdapProviderDataSourceModel - Unexpected type: %T", v)
		}
		ldapProviderModel.LdapProviders = make([]models.LdapProviderDetailModel, 0)
		for _, ldapProvider := range s.GetLdap() {
			model := &models.LdapProviderDetailModel{}
			// tls_revocation_check_level and ocsp_server_uris are available in 9.5
			model.TLSRevocationCheckLevel = types.StringNull()
			model.OcspServerUris = types.ListNull(types.StringType)

			if err = CopyFields(ctx, ldapProvider, model); err != nil {
				return
			}
			ldapProviderModel.LdapProviders = append(ldapProviderModel.LdapProviders, *model)
		}
	default:
		return fmt.Errorf("error pararmeter LdapProviderDataSourceModel - Unexpected type: %T", v)
	}
	return
}

// GetAllLdapProvidersWithFilter Returns all filtered Ldap Providers based on Onefs version.
func GetAllLdapProvidersWithFilter(ctx context.Context, client *client.Client, filter *models.LdapProviderFilterType) (any, error) {
	onfsVersion, err := client.GetOnefsVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get OneFS version: %v", err)
	}

	if onfsVersion.IsGreaterThan("9.4.0") {
		queryParam := client.PscaleOpenAPIClient.AuthApi.ListAuthv16ProvidersLdap(ctx)
		if filter != nil && filter.Scope.ValueString() != "" {
			queryParam = queryParam.Scope(filter.Scope.ValueString())
		}
		result, _, err := queryParam.Execute()
		if err != nil {
			errStr := constants.ReadLdapProviderErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting list of ldap providers: %s", message)
		}
		return result, err
	}
	queryParam := client.PscaleOpenAPIClient.AuthApi.ListAuthv11ProvidersLdap(ctx)
	if filter != nil && filter.Scope.ValueString() != "" {
		queryParam = queryParam.Scope(filter.Scope.ValueString())
	}
	result, _, err := queryParam.Execute()
	if err != nil {
		errStr := constants.ReadLdapProviderErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting ldap provider: %s", message)
	}
	return result, err
}

// GetLdapProvider Returns the Ldap Provider by ldapProviderID based on Onefs version.
func GetLdapProvider(ctx context.Context, client *client.Client, ldapProviderName, scope string) (any, error) {
	onfsVersion, err := client.GetOnefsVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get OneFS version: %v", err)
	}

	if onfsVersion.IsGreaterThan("9.4.0") {
		queryParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv16ProvidersLdapById(ctx, ldapProviderName)
		if scope != "" {
			queryParam = queryParam.Scope(scope)
		}
		result, _, err := queryParam.Execute()
		if err != nil {
			errStr := constants.ReadLdapProviderErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return nil, fmt.Errorf("error getting ldap provider: %s", message)
		}
		if len(result.Ldap) <= 0 {
			message := constants.ReadLdapProviderErrorMsg + "with error: "
			return nil, fmt.Errorf("got empty ldap provider: %s", message)
		}
		return &result.Ldap[0], err
	}
	queryParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv11ProvidersLdapById(ctx, ldapProviderName)
	if scope != "" {
		queryParam = queryParam.Scope(scope)
	}
	result, _, err := queryParam.Execute()
	if err != nil {
		errStr := constants.ReadLdapProviderErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting ldap provider: %s", message)
	}
	if len(result.Ldap) <= 0 {
		message := constants.ReadLdapProviderErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty ldap provider: %s", message)
	}
	return &result.Ldap[0], err
}

// CreateLdapProvider Creates a LdapProvider.
func CreateLdapProvider(ctx context.Context, client *client.Client, plan *models.LdapProviderModel) (err error) {
	onfsVersion, err := client.GetOnefsVersion()
	if err != nil {
		return fmt.Errorf("failed to get OneFS version: %v", err)
	}

	if onfsVersion.IsGreaterThan("9.4.0") {
		ldapToCreate := powerscale.V16ProvidersLdapItem{}
		// Get param from tf input
		if err = ReadFromState(ctx, plan, &ldapToCreate); err != nil {
			return
		}
		createParam := client.PscaleOpenAPIClient.AuthApi.CreateAuthv16ProvidersLdapItem(ctx)
		if !plan.IgnoreUnresolvableServerURIs.IsNull() && !plan.IgnoreUnresolvableServerURIs.IsUnknown() {
			createParam = createParam.Force(plan.IgnoreUnresolvableServerURIs.ValueBool())
		}
		createParam = createParam.V16ProvidersLdapItem(ldapToCreate)
		if _, _, err := createParam.Execute(); err != nil {
			errStr := constants.CreateLdapProviderErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return fmt.Errorf("error creating ldap provider: %s", message)
		}
	} else {
		if !plan.TLSRevocationCheckLevel.IsUnknown() || !plan.OcspServerUris.IsUnknown() {
			return fmt.Errorf("error creating ldap provider: %s", "tls_revocation_check_level and ocsp_server_uris not supported for OneFS < 9.4.0")
		}
		ldapToCreate := powerscale.V11ProvidersLdapItem{}
		// Get param from tf input
		if err = ReadFromState(ctx, plan, &ldapToCreate); err != nil {
			return
		}
		createParam := client.PscaleOpenAPIClient.AuthApi.CreateAuthv11ProvidersLdapItem(ctx)
		if !plan.IgnoreUnresolvableServerURIs.IsNull() && !plan.IgnoreUnresolvableServerURIs.IsUnknown() {
			createParam = createParam.Force(plan.IgnoreUnresolvableServerURIs.ValueBool())
		}
		createParam = createParam.V11ProvidersLdapItem(ldapToCreate)
		if _, _, err := createParam.Execute(); err != nil {
			errStr := constants.CreateLdapProviderErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return fmt.Errorf("error creating ldap provider: %s", message)
		}
	}
	return
}

// UpdateLdapProvider Updates a LdapProvider parameters.
func UpdateLdapProvider(ctx context.Context, client *client.Client, state *models.LdapProviderModel, plan *models.LdapProviderModel) (err error) {

	if !plan.Groupnet.IsUnknown() && !state.Groupnet.Equal(plan.Groupnet) {
		return fmt.Errorf("may not change ldap provider's groupnet")
	}

	onfsVersion, err := client.GetOnefsVersion()
	if err != nil {
		return fmt.Errorf("failed to get OneFS version: %v", err)
	}

	if onfsVersion.IsGreaterThan("9.4.0") {
		ldapToUpdate := powerscale.V16ProvidersLdapIdParams{}
		// Get param from tf input
		if err = ReadFromState(ctx, plan, &ldapToUpdate); err != nil {
			return
		}
		updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv16ProvidersLdapById(ctx, state.Name.ValueString())
		if !plan.IgnoreUnresolvableServerURIs.IsNull() && !plan.IgnoreUnresolvableServerURIs.IsUnknown() {
			updateParam = updateParam.Force(plan.IgnoreUnresolvableServerURIs.ValueBool())
		}
		updateParam = updateParam.V16ProvidersLdapIdParams(ldapToUpdate)
		if _, err := updateParam.Execute(); err != nil {
			errStr := constants.UpdateLdapProviderErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return fmt.Errorf("error updating ldap provider: %s", message)
		}
	} else {
		if !plan.TLSRevocationCheckLevel.IsUnknown() || !plan.OcspServerUris.IsUnknown() {
			return fmt.Errorf("error updating ldap provider: %s", "tls_revocation_check_level and ocsp_server_uris not supported for OneFS < 9.4.0")
		}
		ldapToUpdate := powerscale.V11ProvidersLdapIdParams{}
		// Get param from tf input
		if err = ReadFromState(ctx, plan, &ldapToUpdate); err != nil {
			return
		}
		updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv11ProvidersLdapById(ctx, state.Name.ValueString())
		if !plan.IgnoreUnresolvableServerURIs.IsNull() && !plan.IgnoreUnresolvableServerURIs.IsUnknown() {
			updateParam = updateParam.Force(plan.IgnoreUnresolvableServerURIs.ValueBool())
		}
		updateParam = updateParam.V11ProvidersLdapIdParams(ldapToUpdate)
		if _, err := updateParam.Execute(); err != nil {
			errStr := constants.UpdateLdapProviderErrorMsg + "with error: "
			message := GetErrorString(err, errStr)
			return fmt.Errorf("error updating ldap provider: %s", message)
		}
	}
	return
}

// DeleteLdapProvider Deletes a LdapProvider.
func DeleteLdapProvider(ctx context.Context, client *client.Client, ldapProviderName string) error {
	deleteParam := client.PscaleOpenAPIClient.AuthApi.DeleteAuthv11ProvidersLdapById(ctx, ldapProviderName)
	if _, err := deleteParam.Execute(); err != nil {
		errStr := constants.DeleteLdapProviderErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting LdapProvider - %s : %s", ldapProviderName, message)
	}
	return nil
}
