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

# Available actions: Create, Update (name, path, custom_auth_providers), Delete and Import
# After `terraform apply` of this example file it will create a new access zone with the name set in `name` attribute on the PowerScale

# PowerScale access zones allow you to isolate data and control who can access data in each zone.
resource "powerscale_accesszone" "zone" {

  # Required name of the new access zone
  name = "testAccessZoneSample"

  # Required Groupnet identifier to be assoicated with this access zone
  # Note can not be changed after the access zone is created
  groupnet = "groupnet0"

  # Required Specifies the access zone base directory path
  path = "/ifs"

  # Optional pecifies the list of authentication providers available on this access zone
  # A provider name should be of the form '[provider-type:]provider-name', the provider-type defaults to 'lsa-local-provider'.
  custom_auth_providers = [
    "localProviderName",
    "lsa-local-provider:testAccessZoneSample",
    "lsa-local-provider:localProviderName",
    "lsa-file-provider:fileProviderName",
    "lsa-activedirectory-provider:adsProviderName",
    "lsa-ldap-provider:testProvider",
  ]
}

# After the execution of above resource block, accesszone would have been created on the PowerScale array. For more information, Please check the terraform state file. 