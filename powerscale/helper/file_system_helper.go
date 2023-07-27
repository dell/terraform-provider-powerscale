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
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NamespaceMetadata Needed because we have to marshal this manually.
type NamespaceMetadata struct {
	Attrs []Meta `json:"attrs,omitempty"`
}

// Meta Needed because we have to marshal this manually.
type Meta struct {
	Name      string      `json:"name,omitempty"`
	Namespace string      `json:"namespace,omitempty"`
	Value     interface{} `json:"value,omitempty"`
}

// GetDirectoryMetadata returns the metadata of a filesystem.
func GetDirectoryMetadata(ctx context.Context, client *client.Client, directory string) (*powerscale.NamespaceMetadataList, error) {
	metaParam := client.PscaleOpenAPIClient.NamespaceApi.GetDirectoryMetadata(ctx, directory)
	metaParam = metaParam.Metadata(true)
	result, code, err := metaParam.Execute()
	if err != nil {
		message := GetErrorString(err, "")
		if strings.HasPrefix(message, "json: cannot unmarshal") {
			var unmarshaledJSONResponse NamespaceMetadata
			localVarBody, err := io.ReadAll(code.Body)
			if err != nil {
				return result, err
			}
			err = json.Unmarshal(localVarBody, &unmarshaledJSONResponse)
			if err != nil {
				return result, err
			}
			// Need to manually unmarshal because value could be any type
			var newResult []powerscale.NamespaceMetadataListAttrsInner
			for _, value := range unmarshaledJSONResponse.Attrs {
				name := value.Name
				nameSpace := value.Namespace
				switch val := value.Value.(type) {
				case string:
					newResult = append(newResult, powerscale.NamespaceMetadataListAttrsInner{
						Name:      &name,
						Namespace: &nameSpace,
						Value:     &val,
					})
				case int64:
					newVal := fmt.Sprintf("%v", &val)
					newResult = append(newResult, powerscale.NamespaceMetadataListAttrsInner{
						Name:      &name,
						Namespace: &nameSpace,
						Value:     &newVal,
					})
				case int32:
					newVal := fmt.Sprintf("%v", &val)
					newResult = append(newResult, powerscale.NamespaceMetadataListAttrsInner{
						Name:      &name,
						Namespace: &nameSpace,
						Value:     &newVal,
					})
				case bool:
					newVal := fmt.Sprintf("%v", &val)
					newResult = append(newResult, powerscale.NamespaceMetadataListAttrsInner{
						Name:      &name,
						Namespace: &nameSpace,
						Value:     &newVal,
					})
				}
			}

			return &powerscale.NamespaceMetadataList{
				Attrs: newResult,
			}, nil
		}
	}
	return result, err
}

// GetDirectoryACL returns the filesystem acl.
func GetDirectoryACL(ctx context.Context, client *client.Client, directory string) (*powerscale.NamespaceAcl, error) {
	aclParam := client.PscaleOpenAPIClient.NamespaceApi.GetAcl(ctx, directory)
	aclParam = aclParam.Acl(true)
	aclParam = aclParam.Nsaccess(true)
	result, _, err := aclParam.Execute()
	return result, err
}

// GetDirectoryQuota returns the filesystem quota.
func GetDirectoryQuota(ctx context.Context, client *client.Client, directory string) (*powerscale.V12QuotaQuotas, error) {
	quoteParam := client.PscaleOpenAPIClient.QuotaApi.ListQuotav12QuotaQuotas(ctx)
	quoteParam = quoteParam.Path(directory)
	result, _, err := quoteParam.Execute()
	return result, err
}

// GetDirectorySnapshots returns the filesystem snapshots.
func GetDirectorySnapshots(ctx context.Context, client *client.Client) (*powerscale.V1SnapshotSnapshots, error) {
	result, _, err := client.PscaleOpenAPIClient.SnapshotApi.ListSnapshotv1SnapshotSnapshots(ctx).Execute()
	return result, err
}

// FilterPowerScaleSnapshots returns the filtered list of filesystem snapshots.
func FilterPowerScaleSnapshots(unfilteredSnaps *powerscale.V1SnapshotSnapshots, directory string) []powerscale.V1SnapshotSnapshotExtended {
	var snaps []powerscale.V1SnapshotSnapshotExtended
	for _, vsse := range unfilteredSnaps.Snapshots {
		if strings.HasSuffix(vsse.Path, directory) {
			snaps = append(snaps, vsse)
		}
	}
	return snaps
}

// GetACLKeyObjects create the listObject for ACLs.
func GetACLKeyObjects(listObject []powerscale.AclObject) (types.List, diag.Diagnostics) {
	var keyObjects []attr.Value
	memberType := map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	}
	newType := map[string]attr.Type{
		"accessrights":  types.ListType{ElemType: types.StringType},
		"accesstype":    types.StringType,
		"inherit_flags": types.ListType{ElemType: types.StringType},
		"op":            types.StringType,
		"trustee": types.ObjectType{
			AttrTypes: memberType,
		},
	}

	for _, value := range listObject {
		newMap := make(map[string]attr.Value)
		newMap["accessrights"], _ = types.ListValue(types.StringType, getStringList(value.Accessrights))
		newMap["accesstype"] = types.StringValue(*value.Accesstype)
		newMap["inherit_flags"], _ = types.ListValue(types.StringType, getStringList(value.InheritFlags))
		if value.Op != nil {
			newMap["op"] = types.StringValue(*value.Op)
		} else {
			newMap["op"] = types.StringValue("")
		}
		newMap["trustee"], _ = types.ObjectValue(memberType, getMemberObject(value.Trustee))
		accessObject, _ := types.ObjectValue(newType, newMap)
		keyObjects = append(keyObjects, accessObject)
	}
	return types.ListValue(types.ObjectType{AttrTypes: newType}, keyObjects)
}

func getMemberObject(member *powerscale.MemberObject) map[string]attr.Value {
	newMap := make(map[string]attr.Value)
	if member.Id != nil {
		newMap["id"] = types.StringValue(*member.Id)
	} else {
		newMap["id"] = types.StringValue("")
	}
	if member.Id != nil {
		newMap["name"] = types.StringValue(*member.Name)
	} else {
		newMap["name"] = types.StringValue("")
	}
	if member.Id != nil {
		newMap["type"] = types.StringValue(*member.Type)
	} else {
		newMap["type"] = types.StringValue("")
	}
	return newMap
}

func getStringList(list []string) []attr.Value {
	var keyObjects []attr.Value
	for _, value := range list {
		keyObjects = append(keyObjects, types.StringValue(value))
	}
	return keyObjects
}

func extractMetaModel(ctx context.Context, attr []powerscale.NamespaceMetadataListAttrsInner) ([]models.FileSystemAttribues, error) {
	var metaModel []models.FileSystemAttribues
	for _, nmlai := range attr {
		var destination = models.FileSystemAttribues{}
		err := CopyFields(ctx, nmlai, &destination)
		if err != nil {
			return metaModel, err
		}
		metaModel = append(metaModel, destination)
	}
	return metaModel, nil
}

func extractQuotaModel(ctx context.Context, attr []powerscale.V12QuotaQuotaExtended) ([]models.FileSystemQuota, error) {
	var quotaModel []models.FileSystemQuota
	for _, nmlai := range attr {
		var destination = models.FileSystemQuota{}
		err := CopyFields(ctx, nmlai, &destination)
		if err != nil {
			return quotaModel, err
		}
		destination.ID = types.StringValue(nmlai.Id)
		quotaModel = append(quotaModel, destination)
	}
	return quotaModel, nil
}

func extractSnapshotModel(ctx context.Context, attr []powerscale.V1SnapshotSnapshotExtended) ([]models.FileSystemSnaps, error) {
	var snapModel []models.FileSystemSnaps
	for _, nmlai := range attr {
		var destination = models.FileSystemSnaps{}
		err := CopyFields(ctx, nmlai, &destination)
		if err != nil {
			return snapModel, err
		}
		destination.ID = types.Int64Value(int64(nmlai.Id))
		destination.TargetID = types.Int64Value(int64(nmlai.TargetId))
		snapModel = append(snapModel, destination)
	}
	return snapModel, nil
}

// BuildFilesystemDatasource returns the filesystem datasource fileed.
func BuildFilesystemDatasource(ctx context.Context, state *models.FileSystemDataSourceModel, snap []powerscale.V1SnapshotSnapshotExtended, quota *powerscale.V12QuotaQuotas, acl *powerscale.NamespaceAcl, meta *powerscale.NamespaceMetadataList) error {
	var aclModel models.FileSystemACL
	var quotaModel []models.FileSystemQuota
	var snapModel []models.FileSystemSnaps
	var metaModel []models.FileSystemAttribues

	err := CopyFields(ctx, acl, &aclModel)
	if err != nil {
		return err
	}
	aclModel.ACL, _ = GetACLKeyObjects(acl.Acl)
	metaModel, err = extractMetaModel(ctx, meta.Attrs)
	if err != nil {
		return err
	}
	quotaModel, err = extractQuotaModel(ctx, quota.Quotas)
	if err != nil {
		return err
	}
	snapModel, err = extractSnapshotModel(ctx, snap)
	if err != nil {
		return err
	}
	state.FileSystem = &models.FileSystemDetailModel{
		FileSystemACL:       &aclModel,
		FileSystemAttribues: metaModel,
		FileSystemSnapshots: snapModel,
		FileSystemQuota:     quotaModel,
	}
	return nil
}
