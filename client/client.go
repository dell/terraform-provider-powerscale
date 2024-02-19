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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	SessionEndpoint = "session/1/session"
	AuthType        = "auth_type"
	BasicAuthType   = 0
	SessionAuthType = 1
)

var mutex sync.Mutex

// AuthContextKey define own type for context key to avoid collisions between packages using context.
type AuthContextKey string

// Client type is to hold powerscale client.
type Client struct {
	PscaleOpenAPIClient *powerscale.APIClient
	onefsVersion        *OnefsVersion
	mu                  sync.Mutex
}

// GetOnefsVersion get OneFS version.
func (c *Client) GetOnefsVersion() (*OnefsVersion, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.onefsVersion == nil {
		c.onefsVersion = &OnefsVersion{}
		config, _, err := c.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterConfig(context.Background()).Execute()
		if err != nil {
			return nil, err
		}

		parts := strings.Split(config.OnefsVersion.Release, ".")
		if len(parts) > 2 {
			major, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}
			minor, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, err
			}
			patch, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, err
			}
			c.onefsVersion = &OnefsVersion{major, minor, patch}
		} else {
			return nil, errors.New("Unable to parse OneFS version " + config.OnefsVersion.Release)
		}
	}
	return c.onefsVersion, nil
}

// SetOnefsVersion sets the OneFS version of the client.
// Note: This function is not supposed to be called.
// It is only used for testing.
func (c *Client) SetOnefsVersion(major, minor, patch int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onefsVersion = &OnefsVersion{major, minor, patch}
}

// OnefsVersion present OneFS release version.
type OnefsVersion struct {
	Major, Minor, Patch int
}

func (v OnefsVersion) IsEqualTo(version string) bool {
	parsedVersion, err := parseVersion(version)
	if err != nil {
		return false
	}
	return v.compare(parsedVersion) == 0
}

func (v OnefsVersion) IsLessThan(version string) bool {
	parsedVersion, err := parseVersion(version)
	if err != nil {
		return false
	}
	return v.compare(parsedVersion) < 0
}

func (v OnefsVersion) IsGreaterThan(version string) bool {
	parsedVersion, err := parseVersion(version)
	if err != nil {
		return false
	}
	return v.compare(parsedVersion) > 0
}

func parseVersion(version string) (*OnefsVersion, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid version format")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version")
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version")
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version")
	}

	return &OnefsVersion{Major: major, Minor: minor, Patch: patch}, nil
}

func (v OnefsVersion) compare(other *OnefsVersion) int {
	if v.Major != other.Major {
		return v.Major - other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}
	return v.Patch - other.Patch
}

// String return formatted version string.
func (v OnefsVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// NewClient returns the client.
func NewClient(endpoint string,
	insecure bool,
	user string, pass string, authType, timeout int64) (*Client, error) {
	openAPIClient, err := NewOpenAPIClient(
		context.Background(),
		endpoint,
		insecure,
		user,
		pass,
		authType,
		timeout,
	)
	if err != nil {
		return nil, err
	}

	client := Client{
		PscaleOpenAPIClient: openAPIClient,
	}

	return &client, nil
}

// NewOpenAPIClient returns the OpenApi Client.
func NewOpenAPIClient(ctx context.Context, endpoint string, insecure bool, user string, pass string, authType int64, timeout int64) (*powerscale.APIClient, error) {
	// Setup a User-Agent for your API client (replace the provider name for yours):
	userAgent := "terraform-powerscale-provider/1.0.0"
	jar, err := cookiejar.New(nil)
	if err != nil {
		tflog.Error(ctx, "Got error while creating cookie jar")
	}

	httpclient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Jar:     jar,
	}

	var transport *http.Transport
	if insecure {
		transport = &http.Transport{
			// This is done intentionally if the user sets the skipVerify to true
			/* #nosec */
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				/* #nosec */
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 30,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     90 * time.Second,
		}
	} else {
		// Loading system certs by default if insecure is set to false
		// TODO: Check if we need to remove references to UseCerts from the code
		pool, err := x509.SystemCertPool()
		if err != nil {
			errSysCerts := errors.New("unable to initialize cert pool from system")
			return nil, errSysCerts
		}
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				RootCAs:            pool,
				InsecureSkipVerify: false,
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 30,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     90 * time.Second,
		}
	}

	cfg := powerscale.Configuration{
		HTTPClient:    httpclient,
		DefaultHeader: make(map[string]string),
		UserAgent:     userAgent,
		Debug:         false,
		Servers: powerscale.ServerConfigurations{
			powerscale.ServerConfiguration{
				URL:         endpoint,
				Description: endpoint,
			},
		},
		OperationServers: map[string]powerscale.ServerConfigurations{},
	}
	cfg.DefaultHeader = getHeaders()
	fmt.Printf("config %+v header %+v\n", cfg, cfg.DefaultHeader)

	if authType == BasicAuthType {
		httpclient.Transport = transport
		basicAuth(user, pass, &cfg)
	} else if authType == SessionAuthType {
		ctx = context.WithValue(ctx, AuthContextKey(AuthType), SessionAuthType)
		httpclient.Transport = &TokenTransport{Ctx: ctx, Username: user, Password: pass, RoundTripper: transport}
		err := sessionAuth(ctx, user, pass, &cfg)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Auth type is not valid. Should be 0 or 1. ")
	}

	apiClient := powerscale.NewAPIClient(&cfg)
	if tr, ok := cfg.HTTPClient.Transport.(*TokenTransport); ok {
		tr.Client = apiClient
	}
	return apiClient, nil
}

// Generate the base 64 Authorization string from username / password and add to header.
func basicAuth(username, password string, cfg *powerscale.Configuration) {
	auth := username + ":" + password
	basicAuthString := base64.StdEncoding.EncodeToString([]byte(auth))
	cfg.AddDefaultHeader("Authorization", "Basic "+basicAuthString)
}

func getHeaders() map[string]string {
	header := make(map[string]string)

	header["Content-Type"] = "application/json; charset=utf-8"
	header["Accept"] = "application/json; charset=utf-8"
	return header

}

func sessionAuth(ctx context.Context, user string, pass string, cfg *powerscale.Configuration) error {
	mutex.Lock()
	defer mutex.Unlock()
	host, err := cfg.ServerURLWithContext(ctx, "")
	if err != nil {
		return err
	}
	resp, err := RequestSession(host, user, pass, cfg)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("authentication failed. empty response")
	}
	tflog.Debug(ctx, fmt.Sprintf("response code: %d, response body: %s", resp.StatusCode, resp.Body))
	if resp.Body == nil || resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("authentication failed. response code: %d", resp.StatusCode)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error closing HTTP response: %s", err.Error()))
		}
	}()
	isisessid := getCookie(resp.Cookies(), "isisessid")
	isicsrf := getCookie(resp.Cookies(), "isicsrf")
	if len(isisessid) == 0 || len(isicsrf) == 0 {
		return errors.New("authentication failed. isisessid or isicsrf cookie invalid")
	}
	cfg.AddDefaultHeader("Cookie", fmt.Sprintf("isisessid=%s", isisessid))
	cfg.AddDefaultHeader("X-CSRF-Token", isicsrf)
	cfg.AddDefaultHeader("Referer", host)
	return nil
}

func getCookie(cookies []*http.Cookie, cookieName string) string {
	for _, cookie := range cookies {
		if strings.EqualFold(cookie.Name, cookieName) {
			return cookie.Value
		}
	}
	return ""
}

func RequestSession(host string, user string, pass string, cfg *powerscale.Configuration) (*http.Response, error) {
	sessionUrl := concatUrl(host, SessionEndpoint)
	json := fmt.Sprintf(`{"username":"%s", "password":"%s", "services":["platform", "namespace"]}`, user, pass)

	request, err := http.NewRequest("POST", sessionUrl, strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", cfg.UserAgent)
	resp, err := cfg.HTTPClient.Do(request)
	return resp, err
}

func concatUrl(endpoint string, s string) string {
	endpoint = strings.TrimSuffix(endpoint, "/")
	return fmt.Sprintf("%s/%s", endpoint, s)
}

type TokenTransport struct {
	http.RoundTripper
	Ctx      context.Context
	Username string
	Password string
	Client   *powerscale.APIClient
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.RoundTripper.RoundTrip(req)
	if t.Ctx.Value(AuthContextKey(AuthType)) != SessionAuthType || req.URL.Path == SessionEndpoint {
		return resp, err
	}
	if err != nil {
		return resp, err
	}
	if resp == nil {
		return nil, fmt.Errorf("got empty response for request [%s]", req.URL.Path)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		config := t.Client.GetConfig()
		err := sessionAuth(t.Ctx, t.Username, t.Password, config)
		if err != nil {
			return nil, err
		}
		newReq := req.Clone(req.Context()) // per RoundTrip contract
		for key, value := range config.DefaultHeader {
			newReq.Header.Set(key, value)
		}
		trip, err := t.RoundTripper.RoundTrip(newReq)
		if err != nil {
			return nil, err
		}
		return trip, nil
	}
	return resp, err
}
