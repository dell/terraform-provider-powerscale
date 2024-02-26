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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// ACLSettingsDetailMapper Does the mapping from response to model.
//
//go:noinline
func ACLSettingsDetailMapper(ctx context.Context, aclSettings *powerscale.V11SettingsAclsAclPolicySettings) (models.ACLSettingsDataSourceModel, error) {
	model := models.ACLSettingsDataSourceModel{}
	err := CopyFields(ctx, aclSettings, &model)
	return model, err
}

// GetACLSettings retrieve ACL Settings information.
func GetACLSettings(ctx context.Context, client *client.Client) (*powerscale.V11SettingsAcls, error) {
	queryParam := client.PscaleOpenAPIClient.AuthApi.GetAuthv11SettingsAcls(ctx)
	aclSettingsRes, _, err := queryParam.Execute()
	return aclSettingsRes, err
}

// UpdateACLSettings Update ACL Settings.
func UpdateACLSettings(ctx context.Context, client *client.Client, aclSettingsToUpdate powerscale.V11SettingsAclsAclPolicySettings) error {
	updateParam := client.PscaleOpenAPIClient.AuthApi.UpdateAuthv11SettingsAcls(ctx)
	_, err := updateParam.V11SettingsAcls(aclSettingsToUpdate).Execute()
	return err
}
