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
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateGroupnetDataSourceState updates datasource state.
func UpdateGroupnetDataSourceState(ctx context.Context, groupnetState *models.GroupnetDataSourceModel, groupnetResponse []powerscale.V10NetworkGroupnetExtended) (err error) {
	for _, groupnet := range groupnetResponse {
		var model models.GroupnetModel

		if err = CopyFields(ctx, groupnet, &model); err != nil {
			return
		}

		model.DNSResolverRotate = types.BoolValue(false)
		if groupnet.HasDnsOptions() {
			for _, option := range groupnet.DnsOptions {
				if option == "rotate" {
					model.DNSResolverRotate = types.BoolValue(true)
					break
				}
			}
		}

		groupnetState.Groupnets = append(groupnetState.Groupnets, model)
	}
	return
}

// UpdateGroupnetResourceState updates resource state.
func UpdateGroupnetResourceState(ctx context.Context, groupnetModel *models.GroupnetModel, groupnetResponse *powerscale.V10NetworkGroupnetExtended) (err error) {
	originModel := *groupnetModel

	if err = CopyFields(ctx, groupnetResponse, groupnetModel); err != nil {
		return
	}

	if strings.Trim(*groupnetResponse.Description, " ") == strings.Trim(originModel.Description.ValueString(), " ") {
		groupnetModel.Description = originModel.Description
	}
	if IsListValueEquals(originModel.DNSSearch, groupnetModel.DNSSearch) {
		groupnetModel.DNSSearch = originModel.DNSSearch
	}
	if IsListValueEquals(originModel.DNSServers, groupnetModel.DNSServers) {
		groupnetModel.DNSServers = originModel.DNSServers
	}
	groupnetModel.DNSResolverRotate = types.BoolValue(false)
	if groupnetResponse.HasDnsOptions() {
		for _, option := range groupnetResponse.DnsOptions {
			if option == "rotate" {
				groupnetModel.DNSResolverRotate = types.BoolValue(true)
				break
			}
		}
	}
	return
}

// UpdateGroupnetImportState updates resource import state.
func UpdateGroupnetImportState(ctx context.Context, groupnetModel *models.GroupnetModel, groupnetResponse *powerscale.V10NetworkGroupnetExtended) (err error) {

	if err = CopyFields(ctx, groupnetResponse, groupnetModel); err != nil {
		return
	}

	groupnetModel.DNSResolverRotate = types.BoolValue(false)
	if groupnetResponse.HasDnsOptions() {
		for _, option := range groupnetResponse.DnsOptions {
			if option == "rotate" {
				groupnetModel.DNSResolverRotate = types.BoolValue(true)
				break
			}
		}
	}

	return
}

// GetAllGroupnets returns all groupnets.
func GetAllGroupnets(ctx context.Context, client *client.Client, state *models.GroupnetDataSourceModel) (groupnets []powerscale.V10NetworkGroupnetExtended, err error) {

	groupnetParams := client.PscaleOpenAPIClient.NetworkApi.ListNetworkv10NetworkGroupnets(ctx)

	if state.Filter != nil {
		if !state.Filter.Sort.IsNull() {
			groupnetParams = groupnetParams.Sort(state.Filter.Sort.ValueString())
		}
		if !state.Filter.Dir.IsNull() {
			groupnetParams = groupnetParams.Dir(state.Filter.Dir.ValueString())
		}
		if !state.Filter.Limit.IsNull() {
			groupnetParams = groupnetParams.Limit(int32(state.Filter.Limit.ValueInt64()))
		}
	}

	result, _, err := groupnetParams.Execute()

	// pagination
	for result.Resume != nil && (state.Filter == nil || state.Filter.Limit.IsNull()) {
		groupnetParams = groupnetParams.Resume(*result.Resume)
		newres, _, err := groupnetParams.Execute()
		if err != nil {
			return result.Groupnets, err
		}
		result.Resume = newres.Resume
		result.Groupnets = append(result.Groupnets, newres.Groupnets...)
	}
	return result.Groupnets, err
}

// GetGroupnet Returns the Groupnet by groupnet name.
func GetGroupnet(ctx context.Context, client *client.Client, groupnetName string) (*powerscale.V10NetworkGroupnetExtended, error) {

	result, _, err := client.PscaleOpenAPIClient.NetworkApi.GetNetworkv10NetworkGroupnet(ctx, groupnetName).Execute()
	if err != nil {
		errStr := constants.ReadUserGroupErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return nil, fmt.Errorf("error getting groupnet: %s", message)
	}
	if len(result.Groupnets) < 1 {
		message := constants.ReadGroupnetErrorMsg + "with error: "
		return nil, fmt.Errorf("got empty groupnet: %s", message)
	}

	return &result.Groupnets[0], err
}

// CreateGroupnet Creates a Groupnet.
func CreateGroupnet(ctx context.Context, client *client.Client, plan *models.GroupnetModel) (diags diag.Diagnostics) {

	createParam := client.PscaleOpenAPIClient.NetworkApi.CreateNetworkv10NetworkGroupnet(ctx)

	body := &powerscale.V10NetworkGroupnet{Name: plan.Name.ValueString()}

	if !plan.AllowWildcardSubdomains.IsNull() {
		body.AllowWildcardSubdomains = plan.AllowWildcardSubdomains.ValueBoolPointer()
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DNSCacheEnabled.IsNull() {
		body.DnsCacheEnabled = plan.DNSCacheEnabled.ValueBoolPointer()
	}
	if !plan.ServerSideDNSSearch.IsNull() {
		body.ServerSideDnsSearch = plan.ServerSideDNSSearch.ValueBoolPointer()
	}
	if !plan.DNSResolverRotate.IsNull() && plan.DNSResolverRotate.ValueBool() {
		body.DnsOptions = append(body.DnsOptions, "rotate")
	}
	if !plan.DNSSearch.IsNull() && len(plan.DNSSearch.Elements()) > 0 {
		var DNSSearchList []string
		if diags = plan.DNSSearch.ElementsAs(ctx, &DNSSearchList, false); diags.HasError() {
			return
		}
		body.DnsSearch = DNSSearchList
	}

	if !plan.DNSServers.IsNull() && len(plan.DNSServers.Elements()) > 0 {
		var DNSServerList []string
		if diags = plan.DNSServers.ElementsAs(ctx, &DNSServerList, false); diags.HasError() {
			return
		}
		body.DnsServers = DNSServerList
	}

	createParam = createParam.V10NetworkGroupnet(*body)
	if _, _, err := createParam.Execute(); err != nil {
		errStr := constants.CreateGroupnetErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(fmt.Sprintf("error creating groupnet: %s", message), err.Error())
	}

	return
}

// UpdateGroupnet Updates a Groupnet parameters.
func UpdateGroupnet(ctx context.Context, client *client.Client, state *models.GroupnetModel, plan *models.GroupnetModel) (diags diag.Diagnostics) {
	updateParam := client.PscaleOpenAPIClient.NetworkApi.UpdateNetworkv10NetworkGroupnet(ctx, state.Name.ValueString())

	body := &powerscale.V10NetworkGroupnetExtendedExtended{}

	if !state.AllowWildcardSubdomains.Equal(plan.AllowWildcardSubdomains) {
		body.AllowWildcardSubdomains = plan.AllowWildcardSubdomains.ValueBoolPointer()
	}
	if !state.DNSCacheEnabled.Equal(plan.DNSCacheEnabled) {
		body.DnsCacheEnabled = plan.DNSCacheEnabled.ValueBoolPointer()
	}
	if !state.ServerSideDNSSearch.Equal(plan.ServerSideDNSSearch) {
		body.ServerSideDnsSearch = plan.ServerSideDNSSearch.ValueBoolPointer()
	}
	if !state.Description.Equal(plan.Description) {
		body.Description = plan.Description.ValueStringPointer()
		if body.Description == nil {
			emptyDescription := ""
			body.Description = &emptyDescription
		}
	}
	if !state.Name.Equal(plan.Name) {
		body.Name = plan.Name.ValueStringPointer()
	}
	if !state.DNSResolverRotate.Equal(plan.DNSResolverRotate) {
		body.DnsOptions = make([]string, 0)
		if !plan.DNSResolverRotate.IsNull() && plan.DNSResolverRotate.ValueBool() {
			body.DnsOptions = append(body.DnsOptions, "rotate")
		}
	}
	if !state.DNSServers.Equal(plan.DNSServers) {
		var DNSServerList []string
		if diags = plan.DNSServers.ElementsAs(ctx, &DNSServerList, false); diags.HasError() {
			return
		}
		body.DnsServers = append(make([]string, 0), DNSServerList...)
	}
	if !state.DNSSearch.Equal(plan.DNSSearch) {
		var DNSSearchList []string
		if diags = plan.DNSSearch.ElementsAs(ctx, &DNSSearchList, false); diags.HasError() {
			return
		}
		body.DnsSearch = append(make([]string, 0), DNSSearchList...)
	}

	updateParam = updateParam.V10NetworkGroupnet(*body)
	if _, err := updateParam.Execute(); err != nil {
		errStr := constants.UpdateGroupnetErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(fmt.Sprintf("error updating groupnet: %s", message), err.Error())
	}

	return
}

// DeleteGroupnet Deletes a Groupnet.
func DeleteGroupnet(ctx context.Context, client *client.Client, groupnetName string) error {
	deleteParam := client.PscaleOpenAPIClient.NetworkApi.DeleteNetworkv10NetworkGroupnet(ctx, groupnetName)

	if _, err := deleteParam.Execute(); err != nil {
		errStr := constants.DeleteGroupnetErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting groupnet - %s : %s", groupnetName, message)
	}
	return nil
}
