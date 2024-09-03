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
)

// CreatePeerCert creates a Peer Certificate.
func CreatePeerCert(ctx context.Context, client *client.Client, req powerscale.V7CertificateAuthorityItem) (string, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.CreateSyncv7CertificatesPeerItem(context.Background()).V7CertificatesPeerItem(req).Execute()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// ReadPeerCert reads a Peer Certificate.
func ReadPeerCert(ctx context.Context, client *client.Client, id string) (*powerscale.V16CertificatesSyslogExtended, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv7CertificatesPeerById(context.Background(), id).Execute()
	return resp, err
}

// UpdatePeerCert updates a Peer Certificate.
func UpdatePeerCert(ctx context.Context, client *client.Client, id string, req powerscale.V16CertificatesSyslogIdParams) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv7CertificatesPeerById(context.Background(), id).V7CertificatesPeerIdParams(req).Execute()
	return err
}

// DeletePeerCert deletes a Peer Certificate.
func DeletePeerCert(ctx context.Context, client *client.Client, id string) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.DeleteSyncv7CertificatesPeerById(context.Background(), id).Execute()
	return err
}
