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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetNamespaceACL retrieve Namespace ACL information.
func GetNamespaceACL(ctx context.Context, client *client.Client, model models.NamespaceACLResourceModel) (*powerscale.NamespaceAcl, error) {
	queryParam := client.PscaleOpenAPIClient.NamespaceApi.GetAcl(ctx, model.Namespace.ValueString())
	queryParam = queryParam.Acl(true)
	if !model.Nsaccess.IsNull() {
		queryParam = queryParam.Nsaccess(model.Nsaccess.ValueBool())
	}
	if !model.Zone.IsNull() {
		queryParam = queryParam.Zone(model.Zone.ValueString())
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

// aclTrusteeAttrTypes defines the attr.Type map for ACL trustee objects.
var aclTrusteeAttrTypes = map[string]attr.Type{
	"id":   types.StringType,
	"name": types.StringType,
	"type": types.StringType,
}

// aclEntryAttrTypes defines the attr.Type map for ACL entry objects.
var aclEntryAttrTypes = map[string]attr.Type{
	"accesstype":    types.StringType,
	"accessrights":  types.ListType{ElemType: types.StringType},
	"inherit_flags": types.ListType{ElemType: types.StringType},
	"op":            types.StringType,
	"trustee":       types.ObjectType{AttrTypes: aclTrusteeAttrTypes},
}

// MergeCustomACLWithServerACL merges the user's original acl_custom configuration
// with the server's canonical ACL response. For each ACE, user-specified values
// are preserved (accessrights, inherit_flags, accesstype) while computed-only
// fields that were left unset by the user (op, trustee.name, trustee.type, etc.)
// are filled in from the server response. This avoids Terraform's plan/apply
// consistency check failure when OneFS canonicalizes ACL values.
func MergeCustomACLWithServerACL(userCustomACL, serverACL types.List) (types.List, error) {
	if userCustomACL.IsNull() || userCustomACL.IsUnknown() {
		return serverACL, nil
	}

	userElements := userCustomACL.Elements()
	serverElements := serverACL.Elements()

	mergedElements := make([]attr.Value, len(userElements))

	for i, userElem := range userElements {
		userObj, ok := userElem.(basetypes.ObjectValue)
		if !ok || userObj.IsNull() || userObj.IsUnknown() {
			mergedElements[i] = userElem
			continue
		}

		// If we have a corresponding server ACE at this index, use it to fill unknowns
		var serverObj basetypes.ObjectValue
		hasServerObj := false
		if i < len(serverElements) {
			if sObj, ok := serverElements[i].(basetypes.ObjectValue); ok && !sObj.IsNull() && !sObj.IsUnknown() {
				serverObj = sObj
				hasServerObj = true
			}
		}

		userAttrs := userObj.Attributes()
		mergedAttrs := make(map[string]attr.Value)

		// Copy all user attributes first
		for k, v := range userAttrs {
			mergedAttrs[k] = v
		}

		if hasServerObj {
			serverAttrs := serverObj.Attributes()

			// Fill in unknown scalar fields from server
			for _, field := range []string{"op", "accesstype"} {
				if val, exists := mergedAttrs[field]; exists {
					if strVal, ok := val.(basetypes.StringValue); ok && strVal.IsUnknown() {
						if serverVal, sExists := serverAttrs[field]; sExists {
							mergedAttrs[field] = serverVal
						}
					}
				}
			}

			// Fill in unknown list fields from server
			for _, field := range []string{"accessrights", "inherit_flags"} {
				if val, exists := mergedAttrs[field]; exists {
					if listVal, ok := val.(basetypes.ListValue); ok && listVal.IsUnknown() {
						if serverVal, sExists := serverAttrs[field]; sExists {
							mergedAttrs[field] = serverVal
						}
					}
				}
			}

			// Merge trustee: fill in unknown trustee fields from server
			mergedAttrs["trustee"] = mergeTrustee(mergedAttrs["trustee"], serverAttrs["trustee"])
		}

		mergedObj, diags := types.ObjectValue(aclEntryAttrTypes, mergedAttrs)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: aclEntryAttrTypes}),
				errors.New("failed to merge acl_custom entry with server ACL response")
		}
		mergedElements[i] = mergedObj
	}

	mergedList, diags := types.ListValue(types.ObjectType{AttrTypes: aclEntryAttrTypes}, mergedElements)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: aclEntryAttrTypes}),
			errors.New("failed to build merged acl_custom list")
	}
	return mergedList, nil
}

// mergeTrustee merges user-specified trustee attributes with server trustee values,
// filling in unknown fields from the server response.
func mergeTrustee(userTrustee, serverTrustee attr.Value) attr.Value {
	userObj, ok := userTrustee.(basetypes.ObjectValue)
	if !ok || userObj.IsNull() || userObj.IsUnknown() {
		return serverTrustee
	}

	serverObj, ok := serverTrustee.(basetypes.ObjectValue)
	if !ok || serverObj.IsNull() || serverObj.IsUnknown() {
		return userTrustee
	}

	userAttrs := userObj.Attributes()
	serverAttrs := serverObj.Attributes()
	mergedAttrs := make(map[string]attr.Value)

	for _, field := range []string{"id", "name", "type"} {
		userVal := userAttrs[field]
		if strVal, ok := userVal.(basetypes.StringValue); ok && (strVal.IsUnknown() || strVal.IsNull()) {
			mergedAttrs[field] = serverAttrs[field]
		} else {
			mergedAttrs[field] = userVal
		}
	}

	mergedObj, diags := types.ObjectValue(aclTrusteeAttrTypes, mergedAttrs)
	if diags.HasError() {
		return userTrustee
	}
	return mergedObj
}
