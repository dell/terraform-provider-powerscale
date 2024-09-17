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

package provider

import (
	"context"
	"crypto/tls"
	powerscale "dell/powerscale-go-client"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	"github.com/bytedance/mockey"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerscale": providerserver.NewProtocol6WithError(New("test")()),
}

var powerscaleUsername = ""
var powerscalePassword = ""
var powerscaleEndpoint = ""
var powerScaleSSHIP = ""
var powerscaleSSHPort = "22"
var powerscaleInsecure = false
var ProviderConfig = ""
var SessionAuthProviderConfig = ""
var BasicAuthProviderErrorConfig = ""
var FunctionMocker *mockey.Mocker

func init() {
	err := godotenv.Load("powerscale.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
		return
	}

	powerscaleUsername = os.Getenv("POWERSCALE_USERNAME")
	powerscalePassword = os.Getenv("POWERSCALE_PASSWORD")
	powerscaleEndpoint = os.Getenv("POWERSCALE_ENDPOINT")
	authType := os.Getenv("POWERSCALE_AUTH_TYPE")
	timeout := os.Getenv("POWERSCALE_TIMEOUT")
	insecure := os.Getenv("POWERSCALE_INSECURE")
	powerscaleInsecure = strings.ToLower(insecure) == "true"
	if pscaleSSHPort := os.Getenv("POWERSCALE_SSH_PORT"); len(pscaleSSHPort) > 0 {
		powerscaleSSHPort = pscaleSSHPort
	}
	if len(timeout) == 0 {
		timeout = "2000"
	}
	if len(authType) == 0 {
		authType = "1"
	}

	u, err := url.Parse(powerscaleEndpoint)
	if err != nil {
		log.Fatal("Error parsing POWERSCALE_ENDPOINT:", err.Error())
		return
	}
	powerScaleSSHIP = u.Hostname()

	ProviderConfig = fmt.Sprintf(`
		provider "powerscale" {
			username      = "%s"
			password      = "%s"
  			endpoint      = "%s"
  			insecure      = %s
			auth_type     = %s
			timeout       = %s
		}
	`, powerscaleUsername, powerscalePassword, powerscaleEndpoint, insecure, authType, timeout)

	SessionAuthProviderConfig = fmt.Sprintf(`
		provider "powerscale" {
			username      = "%s"
			password      = "%s"
  			endpoint      = "%s"
  			insecure      = true
			auth_type     = %d
			timeout       = %s
		}
	`, powerscaleUsername, powerscalePassword, powerscaleEndpoint, client.SessionAuthType, timeout)

	BasicAuthProviderErrorConfig = fmt.Sprintf(`
		provider "powerscale" {
			username      = "%s"
			password      = "wrong"
  			endpoint      = "%s"
  			insecure      = true
			auth_type     = %d
			timeout       = %s
		}
	`, powerscaleUsername, powerscaleEndpoint, client.BasicAuthType, timeout)
}

var sweepClient *client.Client

// getClientForRegion returns a common provider client configured for the specified region
func getClientForRegion(_ string) (*client.Client, error) {
	if sweepClient != nil {
		return sweepClient, nil
	}
	client, err := client.NewClient(
		powerscaleEndpoint,
		powerscaleInsecure,
		powerscaleUsername,
		powerscalePassword,
		client.BasicAuthType,
		2000,
	)
	if err != nil {
		return nil, err
	}
	sweepClient = client
	return sweepClient, nil
}

// this is required for initializing sweepers
func TestMain(m *testing.M) {
	resource.TestMain(m)
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
					FunctionMocker = mockey.Mock(client.RequestSession).Return(nil, fmt.Errorf("mock error")).Build()
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
							"Content-Type": []string{"application/json"},
						},
						Close:         true,
						ContentLength: -1,
						Body:          io.NopCloser(strings.NewReader("abcdef")),
					}
					FunctionMocker = mockey.Mock(client.RequestSession).Return(response, nil).Build()
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
			powerscale.ServerConfiguration{
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
}

func TestUnauthorizedErrorParse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      BasicAuthProviderErrorConfig + testAccClusterDataSourceConfig,
				ExpectError: regexp.MustCompile(`.*Unauthorized*.`),
			},
		},
	})
}

func TestUnauthorizedErrorParseHTML(t *testing.T) {
	htmlContent := `
		<!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
		<html>
		
		<head>
			<title>401 Unauthorized to access PAPI.</title>
		</head>
		<script type="text/javascript">
			var regex = new RegExp(/http:\/\/([^:]+):([^/]+)/),
					match = regex.exec(window.location.href);
		
				if (match !== null) {
				   window.location = 'https://' + match[1] + ':' + match[2];
				}
		</script>
		
		<body>
			<h1>Unauthorized to access PAPI.</h1>
			<p>Please contact Administrator.</p>
		</body>
		
		</html>
	`
	errorString, err := helper.ParseBody([]byte(htmlContent))
	t.Log(errorString)
	assert.Equal(t, "401 Unauthorized to access PAPI. ", errorString)
	assert.Empty(t, err)
}

func TestOnefsVersion(t *testing.T) {
	testAccPreCheck(t)
	version940 := client.OnefsVersion{Major: 9, Minor: 4, Patch: 0}
	assert.True(t, version940.IsGreaterThan("9.3.0"))
	assert.True(t, version940.IsEqualTo("9.4.0"))
	assert.True(t, version940.IsLessThan("9.5.0"))
	assert.Equal(t, version940.String(), "9.4.0")
}

func TestErrorOnefsVersion(t *testing.T) {
	testAccPreCheck(t)
	version940 := client.OnefsVersion{Major: 9, Minor: 4, Patch: 0}
	assert.False(t, version940.IsGreaterThan("a.b.c.d"))
	assert.False(t, version940.IsGreaterThan("a.3.0"))
	assert.False(t, version940.IsGreaterThan("9.b.0"))
	assert.False(t, version940.IsGreaterThan("9.3.c"))
}

func TestGetOnefsVersionErrorResponse(t *testing.T) {
	testAccPreCheck(t)
	openAPIClient, err := client.NewOpenAPIClient(
		context.Background(),
		"endpoint",
		true,
		"user",
		"pass",
		0,
		300,
	)
	if err != nil {
		assert.Errorf(t, err, "NewOpenAPIClient failed")
	}

	client := client.Client{
		PscaleOpenAPIClient: openAPIClient,
	}
	clusterConfigMocker := Mock(powerscale.ApiGetClusterv3ClusterConfigRequest.Execute).Return(nil, nil, errors.New("mock REST error")).Build()
	defer clusterConfigMocker.UnPatch()
	_, err = client.GetOnefsVersion()
	assert.True(t, err != nil)
}

func TestSetOnefsVersion(t *testing.T) {
	testAccPreCheck(t)
	openAPIClient, err := client.NewOpenAPIClient(
		context.Background(),
		"endpoint",
		true,
		"user",
		"pass",
		0,
		300,
	)
	if err != nil {
		assert.Errorf(t, err, "NewOpenAPIClient failed")
	}

	client := client.Client{
		PscaleOpenAPIClient: openAPIClient,
	}
	client.SetOnefsVersion(9, 4, 0)
	version, _ := client.GetOnefsVersion()
	assert.True(t, version.IsEqualTo("9.4.0"))
}

func TestInsecureClientWithInsecureParam(t *testing.T) {
	testAccPreCheck(t)
	openAPIClient, err := client.NewOpenAPIClient(
		context.Background(),
		"endpoint",
		true,
		"user",
		"pass",
		0,
		300,
	)
	if err != nil {
		assert.Errorf(t, err, "NewOpenAPIClient failed")
	}
	assert.NotNil(t, openAPIClient)

	openAPIClient, err = client.NewOpenAPIClient(
		context.Background(),
		"endpoint",
		false,
		"user",
		"pass",
		0,
		300,
	)
	if err != nil {
		assert.Errorf(t, err, "NewOpenAPIClient failed")
	}
	assert.NotNil(t, openAPIClient)
}
