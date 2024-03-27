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

// GetRolePrivileges Get a list of Role Privileges.
func GetRolePrivileges(ctx context.Context, client *client.Client) (*powerscale.V14AuthPrivileges, error) {
	rolePrivilegeParams := client.PscaleOpenAPIClient.AuthApi.GetAuthv14AuthPrivileges(ctx)
	rolePrivileges, _, err := rolePrivilegeParams.Execute()
	return rolePrivileges, err
}

// RolePrivilegeDetailMapper Does the mapping from response to model.
//
//go:noinline
func RolePrivilegeDetailMapper(ctx context.Context, rolePrivilege *powerscale.V14AuthPrivilege) (models.RolePrivilegeDetailModel, error) {
	model := models.RolePrivilegeDetailModel{}
	err := CopyFields(ctx, rolePrivilege, &model)
	return model, err
}
