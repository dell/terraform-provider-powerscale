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

// GetNtpServers Get a list of NTP Servers.
func GetNtpServers(ctx context.Context, client *client.Client) (*powerscale.V3NtpServers, error) {
	ntpServerParams := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv3NtpServers(ctx)
	ntpServers, _, err := ntpServerParams.Execute()

	for ntpServers.Resume != nil {
		respAdd, _, errAdd := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv3NtpServers(context.Background()).Resume(*ntpServers.Resume).Execute()
		if errAdd != nil {
			return ntpServers, errAdd
		}
		ntpServers.Resume = respAdd.Resume
		ntpServers.Servers = append(ntpServers.Servers, respAdd.Servers...)
	}

	return ntpServers, err
}

// NtpServerDetailMapper Does the mapping from response to model.
//
//go:noinline
func NtpServerDetailMapper(ctx context.Context, ntpServer *powerscale.V3NtpServerExtended) (models.NtpServerDetailModel, error) {
	model := models.NtpServerDetailModel{}
	err := CopyFields(ctx, ntpServer, &model)
	return model, err
}

// CreateNtpServer Create a NTP Server.
func CreateNtpServer(ctx context.Context, client *client.Client, ntpServer powerscale.V3NtpServer) (string, error) {
	ntpID, _, err := client.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv3NtpServer(ctx).V3NtpServer(ntpServer).Execute()
	if id, ok := ntpID["id"].(string); ok {
		return id, err
	}
	return "", err
}

// GetNtpServer retrieve NTP Server information.
func GetNtpServer(ctx context.Context, client *client.Client, ntpServerModel models.NtpServerResourceModel) (*powerscale.V3NtpServersExtended, error) {
	queryParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv3NtpServer(ctx, ntpServerModel.ID.ValueString())
	ntpServerRes, _, err := queryParam.Execute()
	return ntpServerRes, err
}

// UpdateNtpServer Update a NTP Server.
func UpdateNtpServer(ctx context.Context, client *client.Client, ntpServerID string, ntpServerToUpdate powerscale.V3NtpServerExtendedExtended) error {
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv3NtpServer(ctx, ntpServerID)
	_, err := updateParam.V3NtpServer(ntpServerToUpdate).Execute()
	return err
}

// DeleteNtpServer Delete a NTP Server.
func DeleteNtpServer(ctx context.Context, client *client.Client, ntpServerID string) error {
	_, err := client.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv3NtpServer(ctx, ntpServerID).Execute()
	return err
}

// IsUpdateNtpServerParamInvalid Verify if update params contain params only for creating.
func IsUpdateNtpServerParamInvalid(plan models.NtpServerResourceModel, state models.NtpServerResourceModel) bool {
	return !state.Name.Equal(plan.Name)
}
