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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
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
		"access_rights": types.ListType{ElemType: types.StringType},
		"access_type":   types.StringType,
		"inherit_flags": types.ListType{ElemType: types.StringType},
		"op":            types.StringType,
		"trustee": types.ObjectType{
			AttrTypes: memberType,
		},
	}

	for _, value := range listObject {
		newMap := make(map[string]attr.Value)
		newMap["access_rights"], _ = types.ListValue(types.StringType, getStringList(value.Accessrights))
		newMap["access_type"] = types.StringValue(*value.Accesstype)
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

// UpdateFileSystemResourceState Updates File System Resource State.
func UpdateFileSystemResourceState(ctx context.Context, plan *models.FileSystemResource, acl *powerscale.NamespaceAcl, meta *powerscale.NamespaceMetadataList) (diags diag.Diagnostics) {

	for _, attribute := range meta.Attrs {
		switch *attribute.Name {
		case "type":
			plan.Type = types.StringValue(*attribute.Value)
		case "create_time":
			plan.CreationTime = types.StringValue(*attribute.Value)
		}
	}
	if plan.Type.IsUnknown() {
		plan.Type = types.StringNull()
	}
	if plan.CreationTime.IsUnknown() {
		plan.CreationTime = types.StringNull()
	}
	if owner, ok := acl.GetOwnerOk(); ok {
		// setting owner type
		// - if owner type is unknown, will set it to the actual type in response.
		// - else the owner type in plan not equals the actual type value in the response,
		// still keep the plan owner type value but will add a warning to the diagnostics to report actual owner type in response.
		ownerResType := types.StringNull()
		if owner.Type != nil {
			ownerResType = types.StringValue(owner.GetType())
		}
		ownerTypeUnknown := false
		if plan.Owner.Type.IsUnknown() {
			ownerTypeUnknown = true
			plan.Owner.Type = ownerResType
		} else if !ownerResType.Equal(plan.Owner.Type) {
			diags.AddWarning("File System ACL Owner type is invalid", fmt.Sprintf("File System ACL Owner type is invalid, actual type is '%s'", ownerResType.ValueString()))
		}
		if !plan.Owner.ID.IsUnknown() && plan.Owner.Name.IsUnknown() {
			// when owner id is not empty in plan but owner name is empty, and will check if the owner id matches the actual id in response:
			// - if the owner id match, keep the owner id and will set this unknown owner name to the actual name in response.
			// - if the owner id don't match, still keep the plan owner id value but will add a warning to the diagnostics to report actual owner id in response,
			// and revert the user's type to null if the type was empty.
			plan.Owner.Name = types.StringNull()
			ownerResID := types.StringValue(owner.GetId())
			if !ownerResID.Equal(plan.Owner.ID) {
				if ownerTypeUnknown {
					plan.Owner.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Owner ID is invalid", fmt.Sprintf("File System ACL Owner ID is invalid, actual ID is '%s'", ownerResID.ValueString()))
			} else {
				if ownerName, ok := owner.GetNameOk(); ok {
					plan.Owner.Name = types.StringValue(*ownerName)
				}
			}
		} else if plan.Owner.ID.IsUnknown() && !plan.Owner.Name.IsUnknown() {
			// when owner name is not empty in plan but owner id is empty, and will check if the owner name matches the actual name in response:
			// - if the owner name match, keep the owner name and will set this unknown owner id to the actual id in response.
			// - if the owner name don't match,  still keep the plan owner name value but will add a warning to the diagnostics to report actual owner name in response,
			// and revert the user's type to null if the type was empty.
			plan.Owner.ID = types.StringNull()
			ownerResName := types.StringValue(owner.GetName())
			if !ownerResName.Equal(plan.Owner.Name) {
				if ownerTypeUnknown {
					plan.Owner.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Owner Name is invalid", fmt.Sprintf("File System ACL Owner Name is invalid, actual Name is '%s'", ownerResName.ValueString()))
			} else {
				if ownerID, ok := owner.GetIdOk(); ok {
					plan.Owner.ID = types.StringValue(*ownerID)
				}
			}
		} else {
			// when owner name and id are not empty in plan, and will check if the owner name and id match the actual name and id in response:
			// - if the owner name match, keep the owner name; else still keep the plan owner name value but will add a warning to the diagnostics, and revert the user's type to null if the type was empty.
			// - if the owner id match, keep the owner id; else still keep the plan owner id value but will add a warning to the diagnostics, and revert the user's type to null if the type was empty.
			ownerResID := types.StringValue(owner.GetId())
			if !ownerResID.Equal(plan.Owner.ID) {
				if ownerTypeUnknown {
					plan.Owner.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Owner ID is invalid", fmt.Sprintf("File System ACL Owner ID is invalid, actual ID is '%s'", ownerResID.ValueString()))
			}
			ownerResName := types.StringValue(owner.GetName())
			if !ownerResName.Equal(plan.Owner.Name) {
				if ownerTypeUnknown {
					plan.Owner.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Owner Name is invalid", fmt.Sprintf("File System ACL Owner Name is invalid, actual Name is '%s'", ownerResName.ValueString()))
			}
		}
	}
	if group, ok := acl.GetGroupOk(); ok {
		// setting group type
		// - if group type is unknown, will set it to the actual type in response.
		// - else the group type in plan not equals the actual type value in the response,
		// still keep the plan group type value but will add a warning to the diagnostics to report actual group type in response.
		groupResType := types.StringNull()
		if group.Type != nil {
			groupResType = types.StringValue(group.GetType())
		}
		groupTypeUnknown := false
		if plan.Group.Type.IsUnknown() {
			groupTypeUnknown = true
			plan.Group.Type = groupResType
		} else if !groupResType.Equal(plan.Group.Type) {
			diags.AddWarning("File System ACL Group type is invalid", fmt.Sprintf("File System ACL Group type is invalid, actual type is '%s'", groupResType.ValueString()))
		}
		if !plan.Group.ID.IsUnknown() && plan.Group.Name.IsUnknown() {
			// when group id is not empty in plan but group name is empty, and will check if the group id matches the actual id in response:
			// - if the group id match, keep the group id and will set this unknown group name to the actual name in response.
			// - if the group id don't match, still keep the plan group id value but will add a warning to the diagnostics to report actual group id in response,
			// and revert the user's type to null if the type was empty.
			plan.Group.Name = types.StringNull()
			groupResID := types.StringValue(group.GetId())
			if !groupResID.Equal(plan.Group.ID) {
				if groupTypeUnknown {
					plan.Group.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Group ID is invalid", fmt.Sprintf("File System ACL Group ID is invalid, actual ID is '%s'", groupResID.ValueString()))
			} else {
				if groupName, ok := group.GetNameOk(); ok {
					plan.Group.Name = types.StringValue(*groupName)
				}
			}
		} else if plan.Group.ID.IsUnknown() && !plan.Group.Name.IsUnknown() {
			// when group name is not empty in plan but group id is empty, and will check if the group name matches the actual name in response:
			// - if the group name match, keep the group name and will set this unknown group id to the actual id in response.
			// - if the group name don't match,  still keep the plan group name value but will add a warning to the diagnostics to report actual group name in response,
			// and revert the user's type to null if the type was empty.
			plan.Group.ID = types.StringNull()
			groupResName := types.StringValue(group.GetName())
			if !groupResName.Equal(plan.Group.Name) {
				if groupTypeUnknown {
					plan.Group.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Group Name is invalid", fmt.Sprintf("File System ACL Group Name is invalid, actual Name is '%s'", groupResName.ValueString()))
			} else {
				if groupID, ok := group.GetIdOk(); ok {
					plan.Group.ID = types.StringValue(*groupID)
				}
			}
		} else {
			// when group name and id are not empty in plan, and will check if the group name and id match the actual name and id in response:
			// - if the group name match, keep the group name; else still keep the plan group name value but will add a warning to the diagnostics, and revert the user's type to null if the type was empty.
			// - if the group id match, keep the group id; else still keep the plan group id value but will add a warning to the diagnostics, and revert the user's type to null if the type was empty.
			groupResID := types.StringValue(group.GetId())
			if !groupResID.Equal(plan.Group.ID) {
				if groupTypeUnknown {
					plan.Group.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Group ID is invalid", fmt.Sprintf("File System ACL Group ID is invalid, actual ID is '%s'", groupResID.ValueString()))
			}
			groupResName := types.StringValue(group.GetName())
			if !groupResName.Equal(plan.Group.Name) {
				if groupTypeUnknown {
					plan.Group.Type = types.StringNull()
				}
				diags.AddWarning("File System ACL Group Name is invalid", fmt.Sprintf("File System ACL Group Name is invalid, actual Name is '%s'", groupResName.ValueString()))
			}
		}
	}
	if authoritative, ok := acl.GetAuthoritativeOk(); ok {
		plan.Authoritative = types.StringValue(*authoritative)
	} else if plan.Authoritative.IsUnknown() {
		plan.Authoritative = types.StringNull()
	}
	if mode, ok := acl.GetModeOk(); ok {
		plan.Mode = types.StringValue(*mode)
	} else if plan.Mode.IsUnknown() {
		plan.Mode = types.StringNull()
	}
	plan.ID = types.StringValue(GetDirectoryPath(plan.DirectoryPath.ValueString(), plan.Name.ValueString()))
	plan.FullPath = types.StringValue("/" + plan.ID.ValueString())
	return
}

// UpdateFileSystemResourceImportState Updates File System Resource Import State.
func UpdateFileSystemResourceImportState(ctx context.Context, id string, state *models.FileSystemResource, acl *powerscale.NamespaceAcl, meta *powerscale.NamespaceMetadataList) {

	for _, attribute := range meta.Attrs {
		switch *attribute.Name {
		case "type":
			state.Type = types.StringValue(*attribute.Value)
		case "create_time":
			state.CreationTime = types.StringValue(*attribute.Value)
		}
	}

	if owner, ok := acl.GetOwnerOk(); ok {
		if ownerID, ok := owner.GetIdOk(); ok {
			state.Owner.ID = types.StringValue(*ownerID)
		}
		if ownerName, ok := owner.GetNameOk(); ok {
			state.Owner.Name = types.StringValue(*ownerName)
		}
		if ownerType, ok := owner.GetTypeOk(); ok {
			state.Owner.Type = types.StringValue(*ownerType)
		}
	}
	if group, ok := acl.GetGroupOk(); ok {
		if groupID, ok := group.GetIdOk(); ok {
			state.Group.ID = types.StringValue(*groupID)
		}
		if groupName, ok := group.GetNameOk(); ok {
			state.Group.Name = types.StringValue(*groupName)
		}
		if groupType, ok := group.GetTypeOk(); ok {
			state.Group.Type = types.StringValue(*groupType)
		}
	}
	if authoritative, ok := acl.GetAuthoritativeOk(); ok {
		state.Authoritative = types.StringValue(*authoritative)
	}
	if mode, ok := acl.GetModeOk(); ok {
		state.Mode = types.StringValue(*mode)
	}
	state.ID = types.StringValue(id)
	dir, name := filepath.Split(id)
	dir = filepath.Clean(dir)
	name = filepath.Clean(name)
	state.Name = types.StringValue(name)
	state.DirectoryPath = types.StringValue("/" + dir)
	state.Overwrite = types.BoolValue(false)
	state.Recursive = types.BoolValue(true)
	state.AccessControl = state.Mode
}

// GetDirectoryPath Gets the final directory path(dirPath+dirName).
func GetDirectoryPath(dirPath string, dirName string) string {
	directoryPath := filepath.Join(dirPath, dirName)
	directoryPath = filepath.ToSlash(directoryPath)
	directoryPath = strings.TrimLeft(directoryPath, "/")
	return directoryPath
}

const acl = "acl"
const mode = "mode"

// UpdateFileSystemOwnerAndGroup Updates the file system Owner and Group.
func UpdateFileSystemOwnerAndGroup(ctx context.Context, client *client.Client, dirPath string, plan *models.FileSystemResource, state *models.FileSystemResource) error {
	// Update Owner / Group if modified
	if plan.Owner.Name.ValueString() != state.Owner.Name.ValueString() || plan.Group.Name.ValueString() != state.Group.Name.ValueString() ||
		plan.Owner.ID.ValueString() != state.Owner.ID.ValueString() || plan.Group.ID.ValueString() != state.Group.ID.ValueString() {
		setACLUpdReq := client.PscaleOpenAPIClient.NamespaceApi.SetAcl(ctx, dirPath)
		setACLUpdReq = setACLUpdReq.Acl(true)

		namespaceUpdateUser := *powerscale.NewNamespaceAcl()
		namespaceUpdateUser.SetAuthoritative(mode)

		owner := *powerscale.NewMemberObject()
		if !plan.Owner.ID.IsNull() && !plan.Owner.ID.IsUnknown() {
			owner.Id = plan.Owner.ID.ValueStringPointer()
		}
		if !plan.Owner.Name.IsNull() && !plan.Owner.Name.IsUnknown() {
			owner.Name = plan.Owner.Name.ValueStringPointer()
		}
		if !plan.Owner.Type.IsNull() && !plan.Owner.Type.IsUnknown() {
			owner.Type = plan.Owner.Type.ValueStringPointer()
		}
		namespaceUpdateUser.SetOwner(owner)

		group := *powerscale.NewMemberObject()
		if !plan.Group.ID.IsNull() && !plan.Group.ID.IsUnknown() {
			group.Id = plan.Group.ID.ValueStringPointer()
		}
		if !plan.Group.Name.IsNull() && !plan.Group.Name.IsUnknown() {
			group.Name = plan.Group.Name.ValueStringPointer()
		}
		if !plan.Group.Type.IsNull() && !plan.Group.Type.IsUnknown() {
			group.Type = plan.Group.Type.ValueStringPointer()
		}
		namespaceUpdateUser.SetGroup(group)

		setACLUpdReq = setACLUpdReq.NamespaceAcl(namespaceUpdateUser)

		_, _, err := setACLUpdReq.Execute()
		if err != nil {
			errStr := constants.SetFileSystemACLErrorMsg
			message := GetErrorString(err, errStr)
			return fmt.Errorf("error setting the File system Acl for %s: %s", dirPath, message)
		}
	}
	return nil
}

// UpdateFileSystemAccessControl Updates the file system access control.
func UpdateFileSystemAccessControl(ctx context.Context, client *client.Client, dirPath string, plan *models.FileSystemResource, state *models.FileSystemResource) error {
	// Update Access Control if modified and supported
	if !plan.AccessControl.IsNull() && plan.AccessControl.ValueString() != "" && plan.AccessControl.ValueString() != state.AccessControl.ValueString() {

		setACLUpdReq := client.PscaleOpenAPIClient.NamespaceApi.SetAcl(ctx, dirPath)
		setACLUpdReq = setACLUpdReq.Acl(true)
		namespaceUpdateACL := *powerscale.NewNamespaceAcl()

		newMode, newAuthoritative := getNewAccessControlParams(plan.AccessControl.ValueString())

		if newAuthoritative == acl {
			return fmt.Errorf("error updating access control for File System. Modifying through the Provider only supports to POSIX format(authoritative = mode) but new authoritative is: %s", newAuthoritative)
		}
		namespaceUpdateACL.SetAuthoritative(newAuthoritative)
		namespaceUpdateACL.SetMode(newMode)
		setACLUpdReq = setACLUpdReq.NamespaceAcl(namespaceUpdateACL)

		_, _, err := setACLUpdReq.Execute()
		if err != nil {
			errStr := constants.SetFileSystemACLErrorMsg + "Error Updating AccessControl for the filesystem with error: "
			message := GetErrorString(err, errStr)
			return fmt.Errorf(message)
		}
	}
	return nil
}

func getNewAccessControlParams(accessControl string) (string, string) {
	switch accessControl {
	case "private_read":
		return "0550", acl
	case "private":
		return "0770", acl
	case "public_read":
		return "0775", acl
	case "public_read_write":
		return "0777", acl
	case "public":
		return "0777", acl
	default:
		return accessControl, mode
	}
}

// ExecuteCreate executes the create file system request.
func ExecuteCreate(reqCreate powerscale.ApiCreateDirectoryRequest) (map[string]interface{}, *http.Response, error) {
	return reqCreate.Execute()
}

// DeleteFileSystem Deletes a filesystem.
func DeleteFileSystem(ctx context.Context, client *client.Client, dirPath string) error {

	if _, _, err := client.PscaleOpenAPIClient.NamespaceApi.DeleteDirectory(ctx, dirPath).Execute(); err != nil {
		errStr := constants.DeleteFileSystemErrorMsg
		message := GetErrorString(err, errStr)
		return fmt.Errorf("error deleting filesystem - %s : %s", dirPath, message)
	}
	return nil
}
