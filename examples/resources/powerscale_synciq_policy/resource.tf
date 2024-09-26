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


# Available actions: Create, Update, Delete and Import.
# For more information, Please check the terraform state file.

resource "powerscale_synciq_peer_certificate" "cert_10_10_10_10" {
  path        = "/ifs/peerCert_10_10_10_10.pem"
  name        = "peerCert_10_10_10_10"
  description = "Certificate for the replication peer cluster 10.10.10.10"
}

# PowerScale Sync IQ Policies can be used to replicate files or directories from one cluster to another. 
resource "powerscale_synciq_policy" "policy" {
  # Required
  name             = "policy1"
  action           = "sync"
  source_root_path = "/ifs"
  target_host      = "10.10.10.10"
  target_path      = "/ifs/policy1Sink"

  # Optional
  description = "Policy 1 description"
  enabled     = true
  file_matching_pattern = {
    or_criteria = [
      {
        and_criteria = [
          {
            type     = "name"
            value    = "tfacc"
            operator = "=="
          },
        ]
      },
      {
        and_criteria = [
          {
            type     = "size"
            value    = "200MB"
            operator = ">"
          }
        ]
      }
    ]
  }

  source_network = {
    pool   = "pool0"
    subnet = "subnet0"
  }

  accelerated_failback                  = false
  bandwidth_reservation                 = 100
  changelist                            = false
  check_integrity                       = true
  cloud_deep_copy                       = "allow"
  delete_quotas                         = true
  disable_quota_tmp_dir                 = true
  enable_hash_tmpdir                    = true
  encryption_cipher_list                = "aes128-ctr,aes192-ctr,aes256-ctr"
  ignore_recursive_quota                = true
  log_level                             = "fatal"
  priority                              = 1
  report_max_age                        = 3600
  report_max_count                      = 100
  rpo_alert                             = 10
  schedule                              = "Every Monday at 2AM"
  skip_when_source_unmodified           = true
  snapshot_sync_existing                = true
  snapshot_sync_pattern                 = "*"
  source_exclude_directories            = ["/ifs/SourceExcludeDir"]
  source_include_directories            = ["/ifs/SourceIncludeDir"]
  source_snapshot_archive               = true
  source_snapshot_expiration            = 8 * 60 * 60 * 24 # 8 days
  source_snapshot_pattern               = ".*"
  sync_existing_snapshot_expiration     = true
  sync_existing_target_snapshot_pattern = "%%{SnapName}-%%{SnapCreateTime}"
  target_certificate_id                 = powerscale_synciq_peer_certificate.cert_10_10_10_10.id
  target_compare_initial_sync           = true
  target_detect_modifications           = true
  target_snapshot_alias                 = "SIQ-%%{SrcCluster}-%%{PolicyName}-latest"
  target_snapshot_archive               = true
  target_snapshot_expiration            = 80 * 60 * 60 * 24 # 80 days
  target_snapshot_pattern               = "SIQ-%%{SrcCluster}-%%{PolicyName}-%Y-%m-%d_%H-%M-%S"
}

# After the execution of above resource block, a Sync IQ Policy would have been cached in terraform state file
# and a Sync IQ Policies would have been created/updated on PowerScale.
# For more information, Please check the terraform state file.

# sheduling a policy when source is modified
resource "powerscale_synciq_policy" "policy_when_source_modified" {
  name             = "policy_when_source_modified"
  enabled          = true
  action           = "sync"
  source_root_path = "/ifs/Source"
  target_host      = "10.10.10.9"
  password         = "W0ulntUWannaKn0w"
  target_path      = "/ifs/Sink2"

  # scheduling
  schedule  = "when-source-modified"
  job_delay = 20 * 60 * 60

  # OCSP settings
  ocsp_address = "10.20.10.9"
  # Installed ca certificates can be listed using "isi certificate authority list"
  # which gives a list of all certificates with their short IDs.
  # Then full ID can be looked up using "isi certificate authority view <short ID>".
  ocsp_issuer_certificate_id = "16af57a9f676b0ab126095aa5ebadef22ab31119d644ac95cd4b93dbf3f26aeb"
}
