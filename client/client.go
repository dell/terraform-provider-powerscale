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

package client

import (
	"context"
	"github.com/dell/goisilon"
)

// Client type is to hold powerscale client.
type Client struct {
	PscaleClient *goisilon.Client
}

// NewClient returns the gopowerscale client.
func NewClient(endpoint string,
	insecure bool, verboseLogging uint,
	user, group, pass, volumePath string, volumePathPermissions string, ignoreUnresolvableHosts bool, authType uint8) (*Client, error) {
	pscaleClient, err := goisilon.NewClientWithArgs(
		context.Background(),
		endpoint,
		insecure,
		verboseLogging,
		user,
		group,
		pass,
		volumePath,
		volumePathPermissions,
		ignoreUnresolvableHosts,
		authType,
	)
	if err != nil {
		return nil, err
	}

	client := Client{
		PscaleClient: pscaleClient,
	}
	return &client, nil
}
