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

// GenerateS3Key generates S3 Key.
func GenerateS3Key(ctx context.Context, c *client.Client, state models.S3KeyResourceData) (*powerscale.Createv10S3KeyResponse, error) {
	param := c.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv10S3Key(ctx, state.User.ValueString())
	param = param.Force(true)
	if len(state.Zone.ValueString()) > 0 {
		param = param.Zone(state.Zone.ValueString())
	}
	expiryTime := int32(state.ExistingKeyExpiryTime.ValueInt64())
	eket := powerscale.V10S3Key{
		ExistingKeyExpiryTime: &expiryTime,
	}
	param = param.V10S3Key(eket)
	response, _, err := param.Execute()
	return response, err
}

// GetS3Key gets S3 Key.
func GetS3Key(ctx context.Context, c *client.Client, state models.S3KeyResourceData) (*powerscale.V10S3Keys, error) {
	param := c.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv10S3Key(ctx, state.User.ValueString())
	if len(state.Zone.ValueString()) > 0 {
		param = param.Zone(state.Zone.ValueString())
	}
	response, _, err := param.Execute()
	return response, err
}

// DeleteS3Key delete s3 Key.
func DeleteS3Key(ctx context.Context, c *client.Client, state models.S3KeyResourceData) error {
	param := c.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv10S3Key(ctx, state.User.ValueString())
	if len(state.Zone.ValueString()) > 0 {
		param = param.Zone(state.Zone.ValueString())
	}
	_, err := param.Execute()
	return err
}
