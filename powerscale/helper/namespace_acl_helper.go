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
)

// GetNamespaceACL retrieve Namespace ACL information.
func GetNamespaceACL(ctx context.Context, client *client.Client, model models.NamespaceACLResourceModel) (*powerscale.NamespaceAcl, error) {
	queryParam := client.PscaleOpenAPIClient.NamespaceApi.GetAcl(ctx, model.Namespace.ValueString())
	queryParam = queryParam.Acl(true)
	if !model.Nsaccess.IsNull() {
		queryParam = queryParam.Nsaccess(model.Nsaccess.ValueBool())
	}
	aclSettingsRes, _, err := queryParam.Execute()
	return aclSettingsRes, err
}

// UpdateNamespaceACL Update Namespace ACL.
func UpdateNamespaceACL(ctx context.Context, client *client.Client, model models.NamespaceACLResourceModel, namespaceACLToUpdate powerscale.NamespaceAcl) error {
	authoritative := "acl"
	updateParam := client.PscaleOpenAPIClient.NamespaceApi.SetAcl(ctx, model.Namespace.ValueString())
	updateParam = updateParam.Acl(true)
	if !model.Nsaccess.IsNull() {
		updateParam = updateParam.Nsaccess(model.Nsaccess.ValueBool())
	}
	if !model.Zone.IsNull() {
		updateParam = updateParam.Zone(model.Zone.ValueString())
	}
	namespaceACLToUpdate.Authoritative = &authoritative
	_, _, err := updateParam.NamespaceAcl(namespaceACLToUpdate).Execute()
	return err
}

// IsMemberShipFormatInvalid Verify if user/group format is correct.
func IsMemberShipFormatInvalid(requestBody *powerscale.NamespaceAcl) bool {
	if requestBody.HasOwner() {
		if requestBody.Owner.HasId() && (requestBody.Owner.HasName() || requestBody.Owner.HasType()) {
			return true
		}
		if !requestBody.Owner.HasId() && (!requestBody.Owner.HasName() || !requestBody.Owner.HasType()) {
			return true
		}
	}
	if requestBody.HasGroup() {
		if requestBody.Group.HasId() && (requestBody.Group.HasName() || requestBody.Group.HasType()) {
			return true
		}
		if !requestBody.Group.HasId() && (!requestBody.Group.HasName() || !requestBody.Group.HasType()) {
			return true
		}
	}
	if requestBody.HasAcl() {
		for _, acl := range requestBody.Acl {
			if acl.HasTrustee() {
				if acl.Trustee.HasId() && (acl.Trustee.HasName() || acl.Trustee.HasType()) {
					return true
				}
				if !acl.Trustee.HasId() && (!acl.Trustee.HasName() || !acl.Trustee.HasType()) {
					return true
				}
			}
		}
	}
	return false
}

// IsACLParamProvided Verify if acl is provided as parameters.
func IsACLParamProvided(requestBody *powerscale.NamespaceAcl) bool {
	if !requestBody.HasAcl() && (requestBody.HasOwner() || requestBody.HasGroup()) {
		return false
	}
	return true
}

// CheckNamespaceACLParam Verify if namespace acl parameters are valid.
func CheckNamespaceACLParam(requestBody *powerscale.NamespaceAcl) error {
	if !IsACLParamProvided(requestBody) {
		return errors.New("should provide acl configuration for initialization or updating")
	}
	if IsMemberShipFormatInvalid(requestBody) {
		return errors.New("should provide either id or name+type for owner, group and trustee")
	}
	return nil
}

// GetNamespaceACLDatasource retrieve Namespace ACL datasource information.
func GetNamespaceACLDatasource(ctx context.Context, client *client.Client, model models.NamespaceACLDataSourceModel) (*powerscale.NamespaceAcl, error) {
	queryParam := client.PscaleOpenAPIClient.NamespaceApi.GetAcl(ctx, model.NamespaceACLFilter.Namespace.ValueString())
	queryParam = queryParam.Acl(true)
	if !model.NamespaceACLFilter.Nsaccess.IsNull() {
		queryParam = queryParam.Nsaccess(model.NamespaceACLFilter.Nsaccess.ValueBool())
	}
	namespaceACLResp, _, err := queryParam.Execute()
	return namespaceACLResp, err
}
