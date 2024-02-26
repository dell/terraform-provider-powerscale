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

# Terraform Provider for Dell Technologies PowerScale
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](https://github.com/dell/terraform-provider-powerscale/blob/main/about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-MPL_2.0-blue.svg)](https://github.com/dell/terraform-provider-powerscale/blob/main/LICENSE)

The Terraform Provider for Dell Technologies (Dell) PowerScale allows Data Center and IT administrators to use Hashicorp Terraform to automate and orchestrate the provisioning and management of Dell PowerScale storage systems.

The Terraform Provider can be used to manage access zone, active directory, cluster, user, user group, file system, smb share, nfs export, snapshot, snapshot schedule, quota, groupnet, subnet, network pool, network settings and smart pool settings.

The logged-in user configured in the Terraform provider must possess adequate permissions  against the target Dell PowerScale System

## Table of Contents

* [Code of Conduct](https://github.com/dell/dell-terraform-providers/blob/main/docs/CODE_OF_CONDUCT.md)
* [Maintainer Guide](https://github.com/dell/dell-terraform-providers/blob/main/docs/MAINTAINER_GUIDE.md)
* [Committer Guide](https://github.com/dell/dell-terraform-providers/blob/main/docs/COMMITTER_GUIDE.md)
* [Contributing Guide](https://github.com/dell/dell-terraform-providers/blob/main/docs/CONTRIBUTING.md)
* [List of Adopters](https://github.com/dell/dell-terraform-providers/blob/main/docs/ADOPTERS.md)
* [Security](https://github.com/dell/dell-terraform-providers/blob/main/docs/SECURITY.md)
* [Support](#support)
* [License](#license)
* [Prerequisites](#prerequisites)
* [List of DataSources in Terraform Provider for Dell PowerScale](#list-of-datasources-in-terraform-provider-for-dell-powerscale)
* [List of Resources in Terraform Provider for Dell PowerScale](#list-of-resources-in-terraform-provider-for-dell-powerscale)
* [Releasing, Maintenance and Deprecation](#releasing-maintenance-and-deprecation)

## Support
For any Terraform Provider for Dell PowerScale issues, questions or feedback, please follow our [support process](https://github.com/dell/dell-terraform-providers/blob/main/docs/SUPPORT.md)

## License
The Terraform Provider for Dell PowerScale is released and licensed under the MPL-2.0 license. See [LICENSE](LICENSE) for the full terms.

## Prerequisites

| **Terraform Provider** | **PowerScale Version** | **OS**                    | **Terraform**    | **Golang** |
|------------------------|:-----------------------|:--------------------------|------------------|------------|
| v1.1.0                 | 9.4 <br> 9.5           | ubuntu22.04 <br>  rhel9.x | 1.5.x <br> 1.6.x | 1.21       |

## List of DataSources in Terraform Provider for Dell PowerScale
* Cluster
* Access Zone
* Active Directory
* File System
* Groupnet
* Network Pool
* Network Settings
* NFS Export
* Quota
* Smart Pool Settings
* SMB Share
* Snapshot
* Snapshot Schedule
* Subnet
* User
* User Group

## List of Resources in Terraform Provider for Dell PowerScale
* Access Zone
* Active Directory
* File System
* Groupnet
* Network Pool
* Network Settings
* NFS Export
* Quota
* Smart Pool Settings
* SMB Share
* Snapshot
* Snapshot Schedule
* Subnet
* User
* User Group

## Installation and execution of Terraform Provider for Dell PowerScale

## Installation from public repository

The provider will be fetched from the public repository and installed by Terraform automatically.
Create a file called `main.tf` in your workspace with the following contents

```tf
terraform {
  required_providers {
    powerscale = { 
      version = "1.1.0"
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

## SSL Certificate Verification

For SSL verifcation on RHEL, these steps can be performed:
* Copy the CA certificate to the `/etc/pki/ca-trust/source/anchors` path of the host by any external means.
* Import the SSL certificate to host by running
```
update-ca-trust extract
```

For SSL verification on Ubuntu, these steps can be performed:
* Copy the CA certificate to the `/etc/ssl/certs` path of the host by any external means.
* Import the SSL certificate to host by running:
 ```
  update-ca-certificates
```

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technologies PowerScale follows [Semantic Versioning](https://semver.org/).

New versions will be release regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags in the form of "vx.y.z" where x.y.z corresponds to the version number.
