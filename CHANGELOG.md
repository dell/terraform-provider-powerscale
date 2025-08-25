<!--
Copyright (c) 2023-2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# v1.8.0 (Aug 29, 2025)
## Release Summary
This release enhances .

## Bug Fixes
* `powerscale_cluster_email`: fix version check for correct url call

## Enhancements
* [Added support for configuring Security Identifiers (SIDs) in the PowerScale filesystem resource.](https://github.com/dell/dell-terraform-providers/issues/21)

# v1.7.1 (Apr 29, 2025)
## Release Summary
This release addresses bug fixes to improve stability and user experience for Dell PowerScale.

## Bug Fixes
* [Stabilizing the user addition to the role resource.](https://github.com/dell/dell-terraform-providers/issues/24)
* [Reading addition field from networkpool resource.](https://github.com/dell/dell-terraform-providers/issues/23)
* [Handling user present in non-system zone while using namespaceacl resource.](https://github.com/dell/dell-terraform-providers/issues/18)
* [Added support for SID in filters for user data-source.](https://github.com/dell/dell-terraform-providers/issues/21)

# v1.7.0 (Feb 28, 2025)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_storagepool_tier` for reading the Storagepool Tier in PowerScale.

### Resources

* `powerscale_storagepool_tier` for managing the Storagepool Tier entity of PowerScale.

### Others
N/A

## Enhancements
* All existing resources and datasources are qualified against PowerScale v9.10.

## Bug Fixes
* resource/nfs_export: Fixes the ordering issue for clients and case sensitivity issue for zone attribute
* resource/snapshot_schedule: Fixes the import issue


# v1.6.0 (Nov 29, 2024)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_synciq_replication_report` for reading SyncIQ Replication Report in PowerScale.
* `powerscale_nfs_alias` for reading NFS Alias in PowerScale.
* `powerscale_writable_snapshot` for reading Writeable Snapshot in PowerScale.
* `powerscale_synciq_replication_job` for reading SyncIQ Replication Job in PowerScale.


### Resources

* `powerscale_writable_snapshot` for managing Writeable Snapshots in PowerScale.
* `powerscale_snapshot_restore` for Restoring from snapshot in PowerScale.
* `powerscale_nfs_alias` for managing NFS Alias in PowerScale.
* `powerscale_synciq_replication_job` for managing SyncIQ Replication Job in PowerScale.
* `powerscale_synciq_rules` for managing SyncIQ Replication Performance Rules in PowerScale.

### Others
N/A

## Enhancements
* Selected data sources are enhanced to support additional filters.

## Bug Fixes
N/A


# v1.5.0 (Sept 27, 2024)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_synciq_policy` for reading SyncIQ Policy in PowerScale.
* `powerscale_synciq_global_settings` for reading SyncIQ Global Settings in PowerScale.
* `powerscale_synciq_rule` for reading SyncIQ Rule in PowerScale.
* `powerscale_synciq_peer_certificate` for reading SyncIQ Peer Certificate in PowerScale.


### Resources

* `powerscale_cluster_identity` for managing Cluster Idenity in PowerScale.
* `powerscale_cluster_owner` for managing Cluster Owner in PowerScale.
* `powerscale_cluster_time` for managing Cluster Time settings in PowerScale.
* `powerscale_cluster_snmp` for managing Cluster SNMP settings in PowerScale.
* `powerscale_support_assist` for managing Support Assist settings in PowerScale.
* `powerscale_s3_key` for managing S3 key in PowerScale.
* `powerscale_s3_zone_settings` for managing S3 Zone Settings in PowerScale.
* `powerscale_s3_global_settings` for managing S3 Global Settings in PowerScale.
* `powerscale_synciq_policy` for managing SyncIQ Policy in PowerScale.
* `powerscale_synciq_global_settings` for managing SyncIQ Global Settings in PowerScale.
* `powerscale_synciq_peer_certificate` for managing SyncIQ Peer Certificate in PowerScale.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A


# v1.4.0 (Jun 28, 2024)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_s3_bucket` for reading s3 bucket in PowerScale.
* `powerscale_nfs_global_settings` for reading nfs global settings in PowerScale.
* `powerscale_nfs_zone_settings` for reading nfs zone settings in PowerScale.
* `powerscale_smb_share_settings` for reading smb share settings in PowerScale.
* `powerscale_smb_server_settings` for reading smb server settings in PowerScale.
* `powerscale_namespace_acl` for reading namespace acl in PowerScale.


### Resources

* `powerscale_s3_bucket` for managing s3 bucket in PowerScale.
* `powerscale_nfs_global_settings` for managing nfs global settings in PowerScale.
* `powerscale_nfs_zone_settings` for managing nfs zone settings in PowerScale.
* `powerscale_smb_share_settings` for managing smb share settings in PowerScale.
* `powerscale_smb_server_settings` for managing smb server settings in PowerScale.
* `powerscale_namespace_acl` for managing namespace acl in PowerScale.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A

# v1.3.0 (Apr 30, 2024)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_nfs_export_settings` for reading nfs export settings in PowerScale.
* `powerscale_role` for reading role in PowerScale.
* `powerscale_role_privilege` for reading role privilege in PowerScale.
* `powerscale_user_mapping_rules` for reading user mapping rules in PowerScale.


### Resources

* `powerscale_nfs_export_settings` for managing nfs export settings in PowerScale.
* `powerscale_role` for managing role in PowerScale.
* `powerscale_user_mapping_rules` for managing user mapping rules in PowerScale.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A

# v1.2.0 (Mar 28, 2024)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_ldap_provider` for reading LDAP providers in PowerScale.
* `powerscale_network_rule` for reading network rules in PowerScale.
* `powerscale_filepool_policy` for reading file pool policies in PowerScale.
* `powerscale_ntpserver` for reading NTP servers in PowerScale.
* `powerscale_ntpsettings` for reading NTP settings in PowerScale.
* `powerscale_cluster_email` for reading cluster email settings in PowerScale.
* `powerscale_aclsettings` for reading ACL settings in PowerScale.


### Resources

* `powerscale_ldap_provider` for managing LDAP providers in PowerScale.
* `powerscale_network_rule` for managing network rules in PowerScale.
* `powerscale_filepool_policy` for managing file pool policies in PowerScale.
* `powerscale_ntpserver` for managing NTP servers in PowerScale.
* `powerscale_ntpsettings` for managing NTP settings in PowerScale.
* `powerscale_cluster_email` for managing cluster email settings in PowerScale.
* `powerscale_aclsettings` for managing ACL settings in PowerScale.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A

# v1.1.0 (Nov 24, 2023)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_groupnet` for reading groupnets in PowerScale.
* `powerscale_network_settings` for reading network settings in PowerScale.
* `powerscale_networkpool` for reading network pools in PowerScale.
* `powerscale_quota` for reading quotas in PowerScale.
* `powerscale_smartpool_settings` for reading smart pool settings in PowerScale.
* `powerscale_snapshot` for reading snapshots in PowerScale.
* `powerscale_snapshot_schedule` for reading snapshot schedules in PowerScale.
* `powerscale_subnet` for reading subnets in PowerScale.


### Resources

* `powerscale_groupnet` for managing groupnets in PowerScale.
* `powerscale_network_settings` for managing network settings in PowerScale.
* `powerscale_networkpool` for managing network pools in PowerScale.
* `powerscale_quota` for managing quotas in PowerScale.
* `powerscale_smartpool_settings` for managing smart pool settings in PowerScale.
* `powerscale_snapshot` for managing snapshots in PowerScale.
* `powerscale_snapshot_schedule` for managing snapshot schedules in PowerScale.
* `powerscale_subnet` for managing subnets in PowerScale.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A

# v1.0.0 (Sep 27, 2023)
## Release Summary
The release supports resources and data sources mentioned in the Features section for Dell PowerScale.
## Features

### Data Sources:

* `powerscale_accesszone` for reading access zones in PowerScale.
* `powerscale_adsprovider` for reading ADS providers in PowerScale.
* `powerscale_cluster` for reading cluster in PowerScale.
* `powerscale_filesystem` for reading file systems in PowerScale.
* `powerscale_nfs_export` for reading nfs exports in PowerScale.
* `powerscale_smb_share` for reading smb shares in PowerScale.
* `powerscale_user` for reading users in PowerScale.
* `powerscale_user_group` for reading user groups in PowerScale.


### Resources

* `powerscale_accesszone` for managing access zones in PowerScale.
* `powerscale_adsprovider` for managing ADS providers in PowerScale.
* `powerscale_filesystem` for managing file systems in PowerScale.
* `powerscale_nfs_export` for managing nfs exports in PowerScale.
* `powerscale_smb_share` for managing smb shares in PowerScale.
* `powerscale_user` for managing users in PowerScale.
* `powerscale_user_group` for managing user groups in PowerScale.

### Others
N/A

## Enhancements
N/A

## Bug Fixes
N/A