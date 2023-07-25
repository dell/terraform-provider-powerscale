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
	"crypto/tls"
	"crypto/x509"
	powerscale "dell/powerscale-go-client"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/dell/goisilon"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Client type is to hold powerscale client.
type Client struct {
	PscaleOpenAPIClient *powerscale.APIClient
	PscaleClient        *goisilon.Client
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
	openAPIClient, err := NewOpenAPIClient(
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
		PscaleClient:        pscaleClient,
		PscaleOpenAPIClient: openAPIClient,
	}
	return &client, nil
}

// NewOpenAPIClient returns the OpenApi Client.
func NewOpenAPIClient(ctx context.Context, endpoint string, insecure bool, verboseLogging uint, user, group, pass, volumePath string, volumePathPermissions string, ignoreUnresolvableHosts bool, authType uint8) (*powerscale.APIClient, error) {
	// Setup a User-Agent for your API client (replace the provider name for yours):
	userAgent := "terraform-powermax-provider/1.0.0"
	jar, err := cookiejar.New(nil)
	if err != nil {
		tflog.Error(ctx, "Got error while creating cookie jar")
	}

	httpclient := &http.Client{
		Timeout: 2000 * time.Second,
		Jar:     jar,
	}
	if insecure {
		httpclient.Transport = &http.Transport{
			// This is done intentionally if the user sets the skipVerify to true
			/* #nosec */
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				/* #nosec */
				InsecureSkipVerify: true,
			},
		}
	} else {
		// Loading system certs by default if insecure is set to false
		// TODO: Check if we need to remove references to UseCerts from the code
		pool, err := x509.SystemCertPool()
		if err != nil {
			errSysCerts := errors.New("unable to initialize cert pool from system")
			return nil, errSysCerts
		}
		httpclient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				RootCAs:            pool,
				InsecureSkipVerify: false,
			},
		}
	}
	basicAuthString := basicAuth(user, pass)

	cfg := powerscale.Configuration{
		HTTPClient:    httpclient,
		DefaultHeader: make(map[string]string),
		UserAgent:     userAgent,
		Debug:         false,
		Servers: powerscale.ServerConfigurations{
			{
				URL:         endpoint,
				Description: endpoint,
			},
		},
		OperationServers: map[string]powerscale.ServerConfigurations{},
	}
	cfg.DefaultHeader = getHeaders()
	fmt.Printf("config %+v header %+v", cfg, cfg.DefaultHeader)
	cfg.AddDefaultHeader("Authorization", "Basic "+basicAuthString)
	apiClient := powerscale.NewAPIClient(&cfg)
	return apiClient, nil
}

// Generate the base 64 Authorization string from username / password.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getHeaders() map[string]string {
	header := make(map[string]string)

	header["Content-Type"] = "application/json; charset=utf-8"
	header["Accept"] = "application/json; charset=utf-8"
	return header

}
