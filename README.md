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

The Terraform Provider can be used to manage access zone, active directory, cluster, user, user group, file system, smb share, nfs export, snapshot, snapshot schedule, quota, groupnet, subnet, network pool, network settings, smart pool settings, ldap providers, network rule, file pool policy, ntp server, ntp settings, cluster email settings, acl settings, nfs export settings, role, user mapping rules, role privilege, s3 bucket, nfs global settings, nfs zone settings, smb share settings, smb server settings, namespace acl, cluster identity, cluster snmp, cluster owner, cluster time, support assist, s3 keys, s3 zone settings, s3 global settings, synciq policies, synciq rules, synciq global settings, synciq peer certificates.

The logged-in user configured in the Terraform provider must possess adequate permissions against the target Dell PowerScale System.

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
* [Documentation](#documentation)
* [New to Terraform?](#new-to-terraform)

## Support
For any Terraform Provider for Dell PowerScale issues, questions or feedback, please follow our [support process](https://github.com/dell/dell-terraform-providers/blob/main/docs/SUPPORT.md)

## License
The Terraform Provider for Dell PowerScale is released and licensed under the MPL-2.0 license. See [LICENSE](LICENSE) for the full terms.

## Prerequisites

| **Terraform Provider** | **PowerScale Version** | **OS**                    | **Terraform**               | **Golang** |
|------------------------|:-----------------------|:--------------------------|-----------------------------|------------|
| v1.6.0                 | 9.5 <br> 9.7 <br> 9.8  | ubuntu22.04 <br>  rhel9.x |    1.8.x <br> 1.9.x         | 1.22       |

## List of DataSources in Terraform Provider for Dell PowerScale
* [Cluster](docs/data-sources/cluster.md)
* [Access Zone](docs/data-sources/accesszone.md)
* [ACL Settings](docs/data-sources/aclsettings.md)
* [Active Directory Service Provider](docs/data-sources/adsprovider.md)
* [Cluster Email Settings](docs/data-sources/cluster_email.md)
* [File Pool Policy](docs/data-sources/filepool_policy.md)
* [File System](docs/data-sources/filesystem.md)
* [Groupnet](docs/data-sources/groupnet.md)
* [LDAP Provider](docs/data-sources/ldap_provider.md)
* [Namespace ACL](docs/data-sources/namespace_acl.md)
* [Network Pool](docs/data-sources/networkpool.md)
* [Network Rule](docs/data-sources/network_rule.md)
* [Network Settings](docs/data-sources/network_settings.md)
* [NFS Export](docs/data-sources/nfs_export.md)
* [NFS Export Settings](docs/data-sources/nfs_export_settings.md)
* [NFS Global Settings](docs/data-sources/nfs_global_settings.md)
* [NFS Zone Settings](docs/data-sources/nfs_zone_settings.md)
* [NTP Server](docs/data-sources/ntpserver.md)
* [NTP Settings](docs/data-sources/ntpsettings.md)
* [Quota](docs/data-sources/quota.md)
* [Role](docs/data-sources/role.md)
* [Role Privilege](docs/data-sources/roleprivilege.md)
* [S3 Bucket](docs/data-sources/s3_bucket.md)
* [Smart Pool Settings](docs/data-sources/smartpool_settings.md)
* [SMB Server Settings](docs/data-sources/smb_server_settings.md)
* [SMB Share](docs/data-sources/smb_share.md)
* [SMB Share Settings](docs/data-sources/smb_share_settings.md)
* [Snapshot](docs/data-sources/snapshot.md)
* [Snapshot Schedule](docs/data-sources/snapshot_schedule.md)
* [Subnet](docs/data-sources/subnet.md)
* [User](docs/data-sources/user.md)
* [User Group](docs/data-sources/user_group.md)
* [User Mapping Rules](docs/data-sources/user_mapping_rules.md)
* [SyncIQ Policy](docs/data-sources/synciq_policy.md)
* [SyncIQ Global Settings](docs/data-sources/synciq_global_settings.md)
* [SyncIQ Rule](docs/data-sources/synciq_rule.md)
* [SyncIQ Peer Certificate](docs/data-sources/synciq_peer_certificate.md)
* [SyncIQ Replication Report](docs/data-sources/synciq_replication_report.md)
* [NFS Alias](docs/data-sources/nfs_alias.md)
* [Writeable Snapshot](docs/data-sources/writable_snapshot.md)
* [SyncIQ Replication Job](docs/data-sources/synciq_replication_job.md)

## List of Resources in Terraform Provider for Dell PowerScale
* [Access Zone](docs/resources/accesszone.md)
* [ACL Settings](docs/resources/aclsettings.md)
* [Active Directory Service Provider](docs/resources/adsprovider.md)
* [Cluster Email Settings](docs/resources/cluster_email.md)
* [File Pool Policy](docs/resources/filepool_policy.md)
* [File System](docs/resources/filesystem.md)
* [Groupnet](docs/resources/groupnet.md)
* [LDAP Provider](docs/resources/ldap_provider.md)
* [Namespace ACL](docs/resources/namespace_acl.md)
* [Network Pool](docs/resources/networkpool.md)
* [Network Rule](docs/resources/network_rule.md)
* [Network Settings](docs/resources/network_settings.md)
* [NFS Export](docs/resources/nfs_export.md)
* [NFS Export Settings](docs/resources/nfs_export_settings.md)
* [NFS Global Settings](docs/resources/nfs_global_settings.md)
* [NFS Zone Settings](docs/resources/nfs_zone_settings.md)
* [NTP Server](docs/resources/ntpserver.md)
* [NTP Settings](docs/resources/ntpsettings.md)
* [Quota](docs/resources/quota.md)
* [Role](docs/resources/role.md)
* [S3 Bucket](docs/resources/s3_bucket.md)
* [Smart Pool Settings](docs/resources/smartpool_settings.md)
* [SMB Server Settings](docs/resources/smb_server_settings.md)
* [SMB Share](docs/resources/smb_share.md)
* [SMB Share Settings](docs/resources/smb_share_settings.md)
* [Snapshot](docs/resources/snapshot.md)
* [Snapshot Schedule](docs/resources/snapshot_schedule.md)
* [Subnet](docs/resources/subnet.md)
* [User](docs/resources/user.md)
* [User Group](docs/resources/user_group.md)
* [User Mapping Rules](docs/resources/user_mapping_rules.md)
* [Cluster Identity](docs/resources/cluster_identity.md)
* [Cluster SNMP](docs/resources/cluster_snmp.md)
* [Cluster Owner](docs/resources/cluster_owner.md)
* [Cluster Time](docs/resources/cluster_time.md)
* [Support Assist](docs/resources/support_assist.md)
* [S3 Key](docs/resources/s3_key.md)
* [S3 Zone Settings](docs/resources/s3_zone_settings.md)
* [S3 Global Settings](docs/resources/s3_global_settings.md)
* [SyncIQ Policy](docs/resources/synciq_policy.md)
* [SyncIQ Global Settings](docs/resources/synciq_global_settings.md)
* [SyncIQ Peer Certificate](docs/resources/synciq_peer_certificate.md)
* [Support Assist](docs/resources/support_assist.md)
* [Writeable Snapshot](docs/resources/writable_snapshot.md)
* [Snapshot Restore](docs/resources/snapshot_restore.md)
* [NFS Alias](docs/resources/nfs_alias.md)
* [SyncIQ Replication Job](docs/resources/synciq_replication_job.md)

## Installation and execution of Terraform Provider for Dell PowerScale

## Installation from public repository

The provider will be fetched from the public repository and installed by Terraform automatically.
Create a file called `main.tf` in your workspace with the following contents

```tf
terraform {
  required_providers {
    powerscale = { 
      version = "1.6.0"
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

## Documentation

For more detailed information, please refer to [Dell Terraform Providers Documentation](https://dell.github.io/terraform-docs/).

## New to Terraform?
**Here are some helpful links to get you started if you are new to terraform before using our provider:**

- Intro to Terraform: https://developer.hashicorp.com/terraform/intro 
- Providers: https://developer.hashicorp.com/terraform/language/providers 
- Resources: https://developer.hashicorp.com/terraform/language/resources
- Datasources: https://developer.hashicorp.com/terraform/language/data-sources
- Import: https://developer.hashicorp.com/terraform/language/import
- Variables: https://developer.hashicorp.com/terraform/language/values/variables
- Modules: https://developer.hashicorp.com/terraform/language/modules
- State: https://developer.hashicorp.com/terraform/language/state
