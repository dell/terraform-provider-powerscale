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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// NetworkRuleDataSourceModel describes the data source model.
type NetworkRuleDataSourceModel struct {
	ID                types.String           `tfsdk:"id"`
	NetworkRules      []V3PoolsPoolRulesRule `tfsdk:"network_rules"`
	NetworkRuleFilter *NetworkRuleFilterType `tfsdk:"filter"`
}

// V3PoolsPoolRulesRule struct for V3PoolsPoolRulesRule.
type V3PoolsPoolRulesRule struct {
	// Description for the provisioning rule.
	Description types.String `tfsdk:"description"`
	// Name of the groupnet this rule belongs to
	Groupnet types.String `tfsdk:"groupnet"`
	// Unique rule ID.
	ID types.String `tfsdk:"id"`
	// Interface name the provisioning rule applies to.
	Iface types.String `tfsdk:"iface"`
	// Name of the provisioning rule.
	Name types.String `tfsdk:"name"`
	// Node type the provisioning rule applies to.
	NodeType types.String `tfsdk:"node_type"`
	// Name of the pool this rule belongs to.
	Pool types.String `tfsdk:"pool"`
	// Name of the subnet this rule belongs to.
	Subnet types.String `tfsdk:"subnet"`
}

// NetworkRuleFilterType describes the filter data model.
type NetworkRuleFilterType struct {
	Names    []types.String `tfsdk:"names"`
	Groupnet types.String   `tfsdk:"groupnet"`
	Subnet   types.String   `tfsdk:"subnet"`
	Pool     types.String   `tfsdk:"pool"`
}
