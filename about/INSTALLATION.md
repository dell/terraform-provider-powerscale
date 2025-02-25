<!--
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
-->

# Installation of Terraform Provider for Dell PowerScale

## Installation from public repository

The provider will be fetched from the public repository and installed by Terraform automatically.
Create a file called `main.tf` in your workspace with the following contents

```tf
terraform {
  required_providers {
    powerscale = { 
      version = "1.7.0"
      source = "registry.terraform.io/dell/powerscale"
    }
  }
}
```
Then, in that workspace, run
```
terraform init
``` 

## Installation from source code

1. Clone this repo
2. In the root of this repo run
```
make install
```
Then follow [installation from public repo](#installation-from-public-repository)

## Usage
Once you have installed the PowerScale provider, you can start using it in your Terraform configuration files. The provider has a number of resources that you can use to manage your PowerScale storage arrays.

For example, you can use the `powerscale_accesszone` resource to create a new access zone:
```terraform
resource "powerscale_accesszone" "zone" {
  # Required fields
  name     = "testAccessZone"
  groupnet = "groupnet0"
  path     = "/ifs"

  # Optional to apply Auth Providers
  custom_auth_providers = ["System"]
}
```

## SSL Certificate Verification

For SSL verifcation on RHEL, below steps can be performed:
 * Copy the CA certificate to the `/etc/pki/ca-trust/source/anchors` path of the host by any external means.
 * Import the SSL certificate to host by running
```
update-ca-trust extract
```
