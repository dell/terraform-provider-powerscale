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
	"fmt"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/joho/godotenv"
	"log"
	"os"
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

	ProviderConfig = fmt.Sprintf(`
		provider "powerscale" {
			username      = "%s"
			password      = "%s"
  			endpoint      = "%s"
  			insecure      = true
		}
	`, username, password, endpoint)
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
