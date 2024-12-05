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

type StoragepoolTierResourceModel struct {
	Children types.List `tfsdk:"children"`
	// Nodepools. The names or IDs of the tier's children.
	Id types.Int64 `tfsdk:"id"`
	// The unique identifier of the storagepool tier.
	Lnns types.List `tfsdk:"lnns"`
	// The nodes that are part of this tier.
	Name types.String `tfsdk:"name"`
	// The tier name.
	TransferLimitPct types.Int64 `tfsdk:"transfer_limit_pct"`
	// Stop moving files to this tier when this limit is met.
	TransferLimitState types.String `tfsdk:"transfer_limit_state"`
	// How the transfer limit value is being applied.
}
