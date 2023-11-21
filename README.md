<!--
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
-->

# Terraform Provider for Dell Technologies PowerScale
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](about/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-MPL_2.0-blue.svg)](LICENSE)

The Terraform Provider for Dell Technologies (Dell) PowerScale allows Data Center and IT administrators to use Hashicorp Terraform to automate and orchestrate the provisioning and management of Dell PowerScale storage systems.

The Terraform Provider can be used to manage ...

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

| **Terraform Provider** | **PowerScale Version** | **OS**                                | **Terraform**    | **Golang** |
|------------------------|:-----------------------|:--------------------------------------|------------------|------------|
| v1.0.0                 | 9.4 <br> 9.5           | ubuntu22.04 <br> rhel8.x <br> rhel9.x | 1.4.x <br> 1.5.x | 1.20       |

## List of DataSources in Terraform Provider for Dell PowerScale
* Cluster
* Access Zone
* Active Directory
* File System
* NFS Export
* SMB Share
* User
* User Group

## List of Resources in Terraform Provider for Dell PowerScale
* Access Zone
* Active Directory
* File System
* NFS Export
* SMB Share
* User
* User Group

## Installation and execution of Terraform Provider for Dell PowerScale
The installation and execution steps of Terraform Provider for Dell PowerScale can be found [here](about/INSTALLATION.md).

## Releasing, Maintenance and Deprecation

Terraform Provider for Dell Technologies PowerScale follows [Semantic Versioning](https://semver.org/).

New versions will be release regularly if significant changes (bug fix or new feature) are made in the provider.

Released code versions are located on tags in the form of "vx.y.z" where x.y.z corresponds to the version number.
