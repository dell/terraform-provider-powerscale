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

// GetClusterEmail retrieve cluster email.
func GetClusterEmail(ctx context.Context, client *client.Client) (*powerscale.V1ClusterEmail, error) {
	clusterEmail, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv1ClusterEmail(ctx).Execute()
	return clusterEmail, err
}

// UpdateClusterEmail update cluster email.
func UpdateClusterEmail(ctx context.Context, client *client.Client, v1ClusterEmail powerscale.V1ClusterEmailExtended) error {
	_, err := client.PscaleOpenAPIClient.ClusterApi.UpdateClusterv1ClusterEmail(ctx).V1ClusterEmail(v1ClusterEmail).Execute()
	return err
}
