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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// DeleteSmbShare delete smb share.
func DeleteSmbShare(ctx context.Context, client *client.Client, shareID string, zone *string) error {
	param := client.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv7SmbShare(ctx, shareID)
	if zone != nil {
		param = param.Zone(*zone)
	}
	_, err := param.Execute()
	return err
}

// CreateSmbShare create smb share.
func CreateSmbShare(ctx context.Context, client *client.Client, share powerscale.V7SmbShare) (*powerscale.Createv12SmbShareResponse, error) {
	param := client.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv7SmbShare(ctx).V7SmbShare(share)
	if share.Zone != nil {
		param = param.Zone(*(share.Zone))
	}
	shareID, _, err := param.Execute()
	return shareID, err
}

// GetSmbShare get smb share.
func GetSmbShare(ctx context.Context, client *client.Client, shareID string, zone *string) (*powerscale.V7SmbSharesExtended, error) {
	param := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv7SmbShare(ctx, shareID)
	if zone != nil {
		param = param.Zone(*zone)
	}
	response, _, err := param.Execute()
	return response, err
}

// UpdateSmbShare update smb share.
func UpdateSmbShare(ctx context.Context, client *client.Client, shareID string, zone *string, shareToUpdate powerscale.V7SmbShareExtendedExtended) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv7SmbShare(ctx, shareID).V7SmbShare(shareToUpdate)
	if shareToUpdate.Zone != nil {
		updateParam = updateParam.Zone(*zone)
	}
	_, err := updateParam.Execute()
	return err
}

// ListSmbShares update smb share.
func ListSmbShares(ctx context.Context, client *client.Client, smbFilter *models.SmbShareDatasourceFilter) (*[]powerscale.V7SmbShareExtended, error) {
	listSmbParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv7SmbShares(ctx)
	if smbFilter != nil {
		if !smbFilter.Resume.IsNull() {
			listSmbParam = listSmbParam.Resume(smbFilter.Resume.ValueString())
		}
		if !smbFilter.Zone.IsNull() {
			listSmbParam = listSmbParam.Zone(smbFilter.Zone.ValueString())
		}
		if !smbFilter.Scope.IsNull() {
			listSmbParam = listSmbParam.Scope(smbFilter.Scope.ValueString())
		}
		if !smbFilter.Sort.IsNull() {
			listSmbParam = listSmbParam.Sort(smbFilter.Sort.ValueString())
		}
		if !smbFilter.Dir.IsNull() {
			listSmbParam = listSmbParam.Dir(smbFilter.Dir.ValueString())
		}
		if !smbFilter.Limit.IsNull() {
			listSmbParam = listSmbParam.Limit((smbFilter.Limit.ValueInt32()))
		}
		if !smbFilter.Offset.IsNull() {
			listSmbParam = listSmbParam.Offset((smbFilter.Offset.ValueInt32()))
		}
	}
	smbShares, _, err := listSmbParam.Execute()
	if err != nil {
		return nil, err
	}
	totalSmbShares := smbShares.Shares
	for smbShares.Resume != nil && (smbFilter == nil || smbFilter.Limit.IsNull()) {
		resumeSmbParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv7SmbShares(ctx).Resume(*smbShares.Resume)
		smbShares, _, err = resumeSmbParam.Execute()
		if err != nil {
			return &totalSmbShares, err
		}
		totalSmbShares = append(totalSmbShares, smbShares.Shares...)
	}
	return &totalSmbShares, nil
}
