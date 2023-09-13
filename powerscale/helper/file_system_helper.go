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
func UpdateFileSystemResourceState(ctx context.Context, plan *models.FileSystemResource, state *models.FileSystemResource, acl *powerscale.NamespaceAcl, meta *powerscale.NamespaceMetadataList) {

	for _, attribute := range meta.Attrs {
		switch *attribute.Name {
		case "type":
			state.Type = types.StringValue(*attribute.Value)
		case "create_time":
			state.CreationTime = types.StringValue(*attribute.Value)
		}
	}

	if owner, ok := acl.GetOwnerOk(); ok {
		state.Owner.ID = types.StringValue(*owner.Id)
		state.Owner.Name = types.StringValue(*owner.Name)
		state.Owner.Type = types.StringValue(*owner.Type)
	}
	if group, ok := acl.GetGroupOk(); ok {
		state.Group.ID = types.StringValue(*group.Id)
		state.Group.Name = types.StringValue(*group.Name)
		state.Group.Type = types.StringValue(*group.Type)
	}
	if authoritative, ok := acl.GetAuthoritativeOk(); ok {
		state.Authoritative = types.StringValue(*authoritative)
	}
	if mode, ok := acl.GetModeOk(); ok {
		state.Mode = types.StringValue(*mode)
	}
}

// UpdateFileSystemResourcePlanData Updates File System Resource State from plan data.
func UpdateFileSystemResourcePlanData(plan *models.FileSystemResource, state *models.FileSystemResource) {
	state.ID = types.StringValue(GetDirectoryPath(plan.DirectoryPath.ValueString(), plan.Name.ValueString()))
	state.Name = plan.Name
	state.DirectoryPath = plan.DirectoryPath
	state.Overwrite = plan.Overwrite
	state.Recursive = plan.Recursive
	state.AccessControl = plan.AccessControl
	state.QueryZone = plan.QueryZone
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

// UpdateFileSystem Updates the file system attributes.
func UpdateFileSystem(ctx context.Context, client client.Client, dirPath string, plan *models.FileSystemResource, state *models.FileSystemResource) error {

	// Update Owner / Group if modified
	if plan.Owner.Name.ValueString() != state.Owner.Name.ValueString() || plan.Group.Name.ValueString() != state.Group.Name.ValueString() ||
		plan.Owner.ID.ValueString() != state.Owner.ID.ValueString() || plan.Group.ID.ValueString() != state.Group.ID.ValueString() {
		errAuth := ValidateUserAndGroup(ctx, client, plan.Owner, plan.Group, plan.QueryZone.ValueString())
		if errAuth != nil {
			errStr := constants.UpdateFileSystemErrorMsg
			message := GetErrorString(errAuth, errStr)
			return fmt.Errorf(message)
		}
		setACLUpdReq := client.PscaleOpenAPIClient.NamespaceApi.SetAcl(ctx, dirPath)
		setACLUpdReq = setACLUpdReq.Acl(true)

		namespaceUpdateUser := *powerscale.NewNamespaceAcl()
		namespaceUpdateUser.SetAuthoritative(mode)

		owner := *powerscale.NewMemberObject()
		owner.Id = plan.Owner.ID.ValueStringPointer()
		owner.Name = plan.Owner.Name.ValueStringPointer()
		owner.Type = plan.Owner.Type.ValueStringPointer()
		namespaceUpdateUser.SetOwner(owner)

		group := *powerscale.NewMemberObject()
		group.Id = plan.Group.ID.ValueStringPointer()
		group.Name = plan.Group.Name.ValueStringPointer()
		group.Type = plan.Group.Type.ValueStringPointer()
		namespaceUpdateUser.SetGroup(group)

		setACLUpdReq = setACLUpdReq.NamespaceAcl(namespaceUpdateUser)

		_, _, err := setACLUpdReq.Execute()
		if err != nil {
			errStr := constants.UpdateFileSystemErrorMsg + "Error Updating User / Groups for the filesystem with error: "
			message := GetErrorString(err, errStr)
			return fmt.Errorf(message)
		}
	}
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
			errStr := constants.UpdateFileSystemErrorMsg + "Error Updating AccessControl for the filesystem with error: "
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

// ValidateUserAndGroup check if owner/group information is correct.
func ValidateUserAndGroup(ctx context.Context, client client.Client, owner models.MemberObject, group models.MemberObject, accessZone string) error {
	// Validate owner information
	userReq := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthUser(ctx, owner.Name.ValueString())

	// If zone filter is set use it otherwise leave blank and it will use the default zone
	if accessZone != "" {
		userReq = userReq.Zone(accessZone)
	}

	users, _, err := userReq.Execute()
	if err != nil {
		errStr := "unable to validate user information with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf(message)
	}
	user, ok := users.GetUsersOk()
	if ok && len(user) > 0 {
		userEntity := user[0].OnDiskUserIdentity
		if *userEntity.Id != *owner.ID.ValueStringPointer() || *userEntity.Name != *owner.Name.ValueStringPointer() || *userEntity.Type != *owner.Type.ValueStringPointer() {
			return fmt.Errorf("Incorrect owner information. Please make sure owner id, name, and type are valid")
		}
	} else {
		return fmt.Errorf("unable to retrieve user information")
	}

	// Validate group information
	groupReq := client.PscaleOpenAPIClient.AuthApi.GetAuthv1AuthGroup(ctx, group.Name.ValueString())

	// If zone filter is set use it otherwise leave blank and it will use the default zone
	if accessZone != "" {
		groupReq = groupReq.Zone(accessZone)
	}

	groups, _, err := groupReq.Execute()
	if err != nil {
		errStr := "unable to validate group information with error: "
		message := GetErrorString(err, errStr)
		return fmt.Errorf(message)
	}
	grp, okGroup := groups.GetGroupsOk()
	if okGroup && len(grp) > 0 {
		grpEntity := grp[0].Gid
		if *grpEntity.Id != *group.ID.ValueStringPointer() || *grpEntity.Name != *group.Name.ValueStringPointer() || *grpEntity.Type != *group.Type.ValueStringPointer() {
			return fmt.Errorf("incorrect group information. Please make sure group id, name, and type are valid")
		}
	} else {
		return fmt.Errorf("unable to retrieve group information")
	}
	return nil
}
