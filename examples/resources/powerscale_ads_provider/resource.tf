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
resource "powerscale_adsprovider" "ads_test" {
  #   required
  name     = "ADS.PROVIDER.EXAMPLE.COM"
  user     = "admin"
  password = "password"

  #   optional query parameters
  #   scope = "effective"
  #   check_duplicates = true

  #   optional for creating
  #   dns_domain = "testDNSDomain"
  #   groupnet = "testGroupNet"
  #   instance = "testInstance"
  #   kerberos_hdfs_spn = true
  #   kerberos_nfs_spn = true
  #   machine_account = "testMachineAccount"
  #   organizational_unit = "testOrganizationalUnit"

  #   optional for updating
  #   domain_controller = "testDomainController"
  #   reset_schannel = true
  #   spns = ["testSPN"]

  #   optional for creating and updating
  #   allocate_gids = true
  #   allocate_uids = true
  #   assume_default_domain = true
  #   authentication = true
  #   check_online_interval = 310
  #   controller_time = 1692087697
  #   create_home_directory = true
  #   domain_offline_alerts = true
  #   extra_expected_spns = ["testExtraExpectedSPN"]
  #   findable_groups = ["testFindableGroup"]
  #   findable_users = ["testFindableUser"]
  #   home_directory_template = "testHomeDirectoryTemplate"
  #   ignore_all_trusts = true
  #   ignored_trusted_domains = ["testIgnoredTrustedDomain"]
  #   include_trusted_domains = ["testIncludeTrustedDomain"]
  #   ldap_sign_and_seal = true
  #   login_shell = "testLoginShell"
  #   lookup_domains = ["testLookupDomains"]
  #   lookup_groups = true
  #   lookup_normalize_groups = true
  #   lookup_normalize_users = true
  #   lookup_users = true
  #   machine_password_changes = true
  #   machine_password_lifespan = 2591000
  #   node_dc_affinity = "testNodeDcAffinity"
  #   node_dc_affinity_timeout = 1000000
  #   nss_enumeration = true
  #   restrict_findable = true
  #   rpc_call_timeout = 70
  #   server_retry_limit = 4
  #   sfu_support = "testSfuSupport"
  #   store_sfu_mappings = true
  #   unfindable_groups = ["testUnfindableGroup"]
  #   unfindable_users = ["testUnfindableUser"]
}