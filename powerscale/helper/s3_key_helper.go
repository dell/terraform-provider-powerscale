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

type S3keyHelper struct {
	client                *client.Client
	v10S3KeyId, zone      string
	existingKeyExpiryTime int32
}

// CreateS3Key create s3 bucket.
func (skh *S3keyHelper) CreateS3Key(ctx context.Context) (*powerscale.Createv10S3KeyResponseKeys, error) {
	param := skh.client.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv10S3Key(ctx, skh.v10S3KeyId)
	if skh.existingKeyExpiryTime > 0 {
		ex := powerscale.V10S3Key{
			ExistingKeyExpiryTime: &skh.existingKeyExpiryTime,
		}
		param = param.V10S3Key(ex)
	}
	if len(skh.zone) > 0 {
		param = param.Zone(skh.zone)
	}
	response, _, err := param.Execute()
	return &response.Keys, err
}

// GetS3Key gets S3 Bucket.
func (skh *S3keyHelper) GetS3Key(ctx context.Context) (*powerscale.V10S3KeysKeys, error) {
	param := skh.client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv10S3Key(ctx, skh.v10S3KeyId)
	if len(skh.zone) > 0 {
		param = param.Zone(skh.zone)
	}
	response, _, err := param.Execute()
	return &response.Keys, err
}

// DeleteS3Key delete s3 bucket.
func (skh *S3keyHelper) DeleteS3Key(ctx context.Context) error {
	param := skh.client.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv10S3Key(ctx, skh.v10S3KeyId)
	if len(skh.zone) > 0 {
		param = param.Zone(skh.zone)
	}
	_, err := param.Execute()
	return err
}
