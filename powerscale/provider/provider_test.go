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

package provider

import (
	"context"
	"crypto/tls"
	powerscale "dell/powerscale-go-client"
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"
	"terraform-provider-powerscale/client"
	"testing"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerscale": providerserver.NewProtocol6WithError(New("test")()),
}

var ProviderConfig = ""
var SessionAuthProviderConfig = ""
var FunctionMocker *Mocker

func init() {
	err := godotenv.Load("powerscale.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	username := os.Getenv("POWERSCALE_USERNAME")
	password := os.Getenv("POWERSCALE_PASSWORD")
	endpoint := os.Getenv("POWERSCALE_ENDPOINT")
	authType := os.Getenv("POWERSCALE_AUTH_TYPE")
	timeout := os.Getenv("POWERSCALE_TIMEOUT")
	insecure := os.Getenv("POWERSCALE_INSECURE")
	if len(timeout) == 0 {
		timeout = "2000"
	}

	ProviderConfig = fmt.Sprintf(`
		provider "powerscale" {
			username      = "%s"
			password      = "%s"
  			endpoint      = "%s"
  			insecure      = %s
			auth_type     = %s
			timeout       = %s
		}
	`, username, password, endpoint, insecure, authType, timeout)

	SessionAuthProviderConfig = fmt.Sprintf(`
		provider "powerscale" {
			username      = "%s"
			password      = "%s"
  			endpoint      = "%s"
  			insecure      = true
			auth_type     = %d
			timeout       = %s
		}
	`, username, password, endpoint, client.SessionAuthType, timeout)
}

func testAccPreCheck(t *testing.T) {
	// Check that the required environment variables are set.
	if os.Getenv("POWERSCALE_ENDPOINT") == "" {
		t.Fatal("POWERSCALE_ENDPOINT environment variable not set")
	}
	if os.Getenv("POWERSCALE_USERNAME") == "" {
		t.Fatal("POWERSCALE_USERNAME environment variable not set")
	}
	if os.Getenv("POWERSCALE_PASSWORD") == "" {
		t.Fatal("POWERSCALE_PASSWORD environment variable not set")
	}

	// Before each test clear out the mocker
	if FunctionMocker != nil {
		FunctionMocker.UnPatch()
	}
}

func TestSessionAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: SessionAuthProviderConfig + testAccClusterDataSourceConfig,
				Check:  resource.TestCheckResourceAttr("data.powerscale_cluster.test", "id", "cluster-data-source"),
			},
		},
	})
}

func TestSessionAuthError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					FunctionMocker = Mock(client.RequestSession).Return(nil, fmt.Errorf("mock error")).Build()
				},
				Config:      SessionAuthProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*mock error*.`),
			},
		},
	})
}

func TestSessionAuthUnauthorizedError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					response := &http.Response{
						Status:     "401 Unauthorized",
						StatusCode: 401,
						Proto:      "HTTP/1.0",
						ProtoMajor: 1,
						ProtoMinor: 0,
						Request:    &http.Request{Method: "POST"},
						Header: http.Header{
							"Content-Type": {"application/json"},
						},
						Close:         true,
						ContentLength: -1,
						Body:          io.NopCloser(strings.NewReader("abcdef")),
					}
					FunctionMocker = Mock(client.RequestSession).Return(response, nil).Build()
				},
				Config:      SessionAuthProviderConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*authentication failed*.`),
			},
		},
	})
}

func TestSessionRefresh(t *testing.T) {
	testAccPreCheck(t)
	// Build a powerscale client with session cookie unset, to mock the timeout cookie
	// and see if the request can refresh the cookie automatically
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	httpclient := &http.Client{
		Jar: jar,
	}
	ctx := context.WithValue(context.Background(), client.AuthContextKey(client.AuthType), client.SessionAuthType)
	username := os.Getenv("POWERSCALE_USERNAME")
	password := os.Getenv("POWERSCALE_PASSWORD")
	endpoint := os.Getenv("POWERSCALE_ENDPOINT")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}
	httpclient.Transport = &client.TokenTransport{Ctx: ctx, Username: username, Password: password, RoundTripper: transport}
	configuration := powerscale.Configuration{
		DefaultHeader: make(map[string]string),
		UserAgent:     "terraform-powerscale-provider/1.0.0",
		HTTPClient:    httpclient,
		Servers: powerscale.ServerConfigurations{
			{
				URL:         endpoint,
				Description: endpoint,
			},
		},
		OperationServers: map[string]powerscale.ServerConfigurations{},
	}
	apiClient := powerscale.NewAPIClient(&configuration)
	if tr, ok := configuration.HTTPClient.Transport.(*client.TokenTransport); ok {
		tr.Client = apiClient
	}
	_, resp, err := apiClient.ClusterApi.GetClusterv3ClusterConfig(ctx).Execute()
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Log("succeeded to refresh token and get cluster config")
	}
	if resp.StatusCode == http.StatusUnauthorized {
		t.Errorf("failed to refresh session token")
	}
}
