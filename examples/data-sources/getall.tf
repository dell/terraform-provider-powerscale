/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

# =============================================================================
# POWERSCALE - RETRIEVE ALL DATA SOURCES
# =============================================================================
# This configuration retrieves ALL data from a PowerScale cluster.
# Copy this file to a working directory and create a terraform.tfvars file.
# Use "terraform apply" or "terraform refresh" to populate the state.
# =============================================================================

terraform {
  required_providers {
    powerscale = {
      source = "registry.terraform.io/dell/powerscale"
    }
  }
}

# -----------------------------------------------------------------------------
# VARIABLES - Set these in terraform.tfvars
# -----------------------------------------------------------------------------

variable "powerscale_endpoint" {
  description = "PowerScale API endpoint (e.g., https://192.168.1.100:8080)"
  type        = string
}

variable "powerscale_username" {
  description = "PowerScale API username"
  type        = string
}

variable "powerscale_password" {
  description = "PowerScale API password"
  type        = string
  sensitive   = true
}

variable "powerscale_insecure" {
  description = "Skip SSL verification (true for self-signed certs)"
  type        = bool
  default     = false
}

variable "powerscale_timeout" {
  description = "API request timeout in milliseconds"
  type        = number
  default     = 2000
}

# -----------------------------------------------------------------------------
# PROVIDER CONFIGURATION
# -----------------------------------------------------------------------------

provider "powerscale" {
  endpoint = var.powerscale_endpoint
  username = var.powerscale_username
  password = var.powerscale_password
  insecure = var.powerscale_insecure
  timeout  = var.powerscale_timeout
}

# =============================================================================
# CLUSTER INFORMATION
# =============================================================================

# Cluster configuration: nodes, version, hardware, licenses
data "powerscale_cluster" "all" {}
output "cluster" { value = data.powerscale_cluster.all }

# Cluster email notification settings
data "powerscale_cluster_email" "all" {}
output "cluster_email" { value = data.powerscale_cluster_email.all }

# =============================================================================
# NETWORKING
# =============================================================================

# Groupnets - top-level networking containers
data "powerscale_groupnet" "all" {}
output "groupnets" { value = data.powerscale_groupnet.all }

# Subnets within groupnets
data "powerscale_subnet" "all" {}
output "subnets" { value = data.powerscale_subnet.all }

# Network pools for IP address management
data "powerscale_networkpool" "all" {}
output "networkpools" { value = data.powerscale_networkpool.all }

# Global network settings
data "powerscale_network_settings" "all" {}
output "network_settings" { value = data.powerscale_network_settings.all }

# Network provisioning rules
data "powerscale_network_rule" "all" {}
output "network_rules" { value = data.powerscale_network_rule.all }

# =============================================================================
# ACCESS ZONES & AUTHENTICATION
# =============================================================================

# Access zones - isolated data access control areas
data "powerscale_accesszone" "all" {}
output "accesszones" { value = data.powerscale_accesszone.all }

# ACL settings
data "powerscale_aclsettings" "all" {}
output "aclsettings" { value = data.powerscale_aclsettings.all }

# Active Directory providers
data "powerscale_adsprovider" "all" {}
output "adsproviders" { value = data.powerscale_adsprovider.all }

# LDAP providers
data "powerscale_ldap_provider" "all" {}
output "ldap_providers" { value = data.powerscale_ldap_provider.all }

# User mapping rules
data "powerscale_user_mapping_rules" "all" {}
output "user_mapping_rules" { value = data.powerscale_user_mapping_rules.all }

# =============================================================================
# USERS & GROUPS
# =============================================================================

# Local users
data "powerscale_user" "all" {}
output "users" { value = data.powerscale_user.all }

# Local user groups
data "powerscale_user_group" "all" {}
output "user_groups" { value = data.powerscale_user_group.all }

# Roles
data "powerscale_role" "all" {}
output "roles" { value = data.powerscale_role.all }

# Role privileges
data "powerscale_roleprivilege" "all" {}
output "role_privileges" { value = data.powerscale_roleprivilege.all }

# =============================================================================
# QUOTAS
# =============================================================================

# SmartQuotas for storage limits
data "powerscale_quota" "all" {}
output "quotas" { value = data.powerscale_quota.all }

# =============================================================================
# NFS CONFIGURATION
# =============================================================================

# NFS exports
data "powerscale_nfs_export" "all" {}
output "nfs_exports" { value = data.powerscale_nfs_export.all }

# NFS aliases
data "powerscale_nfs_alias" "all" {}
output "nfs_aliases" { value = data.powerscale_nfs_alias.all }

# Global NFS settings
data "powerscale_nfs_global_settings" "all" {}
output "nfs_global_settings" { value = data.powerscale_nfs_global_settings.all }

# NFS export default settings
data "powerscale_nfs_export_settings" "all" {}
output "nfs_export_settings" { value = data.powerscale_nfs_export_settings.all }

# NFS zone settings
data "powerscale_nfs_zone_settings" "all" {}
output "nfs_zone_settings" { value = data.powerscale_nfs_zone_settings.all }

# =============================================================================
# SMB CONFIGURATION
# =============================================================================

# SMB shares
data "powerscale_smb_share" "all" {}
output "smb_shares" { value = data.powerscale_smb_share.all }

# SMB server settings
data "powerscale_smb_server_settings" "all" {}
output "smb_server_settings" { value = data.powerscale_smb_server_settings.all }

# SMB share default settings
data "powerscale_smb_share_settings" "all" {}
output "smb_share_settings" { value = data.powerscale_smb_share_settings.all }

# =============================================================================
# S3 CONFIGURATION
# =============================================================================

# S3 buckets
data "powerscale_s3_bucket" "all" {}
output "s3_buckets" { value = data.powerscale_s3_bucket.all }

# =============================================================================
# SNAPSHOTS
# =============================================================================

# Snapshots
data "powerscale_snapshot" "all" {}
output "snapshots" { value = data.powerscale_snapshot.all }

# Snapshot schedules
data "powerscale_snapshot_schedule" "all" {}
output "snapshot_schedules" { value = data.powerscale_snapshot_schedule.all }

# Writable snapshots
data "powerscale_writable_snapshot" "all" {}
output "writable_snapshots" { value = data.powerscale_writable_snapshot.all }

# =============================================================================
# SYNCIQ (REPLICATION)
# =============================================================================

# SyncIQ global settings
data "powerscale_synciq_global_settings" "all" {}
output "synciq_global_settings" { value = data.powerscale_synciq_global_settings.all }

# SyncIQ policies
data "powerscale_synciq_policy" "all" {}
output "synciq_policies" { value = data.powerscale_synciq_policy.all }

# SyncIQ rules
data "powerscale_synciq_rule" "all" {}
output "synciq_rules" { value = data.powerscale_synciq_rule.all }

# SyncIQ peer certificates
data "powerscale_synciq_peer_certificate" "all" {}
output "synciq_peer_certificates" { value = data.powerscale_synciq_peer_certificate.all }

# SyncIQ replication jobs
data "powerscale_synciq_replication_job" "all" {}
output "synciq_replication_jobs" { value = data.powerscale_synciq_replication_job.all }

# SyncIQ replication reports
data "powerscale_synciq_replication_report" "all" {}
output "synciq_replication_reports" { value = data.powerscale_synciq_replication_report.all }

# =============================================================================
# STORAGE POOLS & TIERS
# =============================================================================

# SmartPools settings
data "powerscale_smartpool_settings" "all" {}
output "smartpool_settings" { value = data.powerscale_smartpool_settings.all }

# Storage pool tiers
data "powerscale_storagepool_tier" "all" {}
output "storagepool_tiers" { value = data.powerscale_storagepool_tier.all }

# File pool policies
data "powerscale_filepool_policy" "all" {}
output "filepool_policies" { value = data.powerscale_filepool_policy.all }

# =============================================================================
# NTP (TIME SYNCHRONIZATION)
# =============================================================================

# NTP servers
data "powerscale_ntpserver" "all" {}
output "ntpservers" { value = data.powerscale_ntpserver.all }

# NTP settings
data "powerscale_ntpsettings" "all" {}
output "ntpsettings" { value = data.powerscale_ntpsettings.all }
