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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// create struct to unmarshall tfsdk schema defined in S3KeyResourceSchema().
type S3KeyResourceData struct {
	AccessID              types.String `tfsdk:"access_id"`
	User                  types.String `tfsdk:"user"`
	Zone                  types.String `tfsdk:"zone"`
	ExistingKeyExpiryTime types.Int64  `tfsdk:"existing_key_expiry_time"`
	SecretKey             types.String `tfsdk:"secret_key"`
	SecretKeyTimestamp    types.Int64  `tfsdk:"secret_key_timestamp"`
	OldSecretKey          types.String `tfsdk:"old_secret_key"`
	OldKeyExpiry          types.Int64  `tfsdk:"old_key_expiry"`
	OldKeyTimestamp       types.Int64  `tfsdk:"old_key_timestamp"`
}
