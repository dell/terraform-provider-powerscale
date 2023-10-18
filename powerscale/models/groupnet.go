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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// GroupnetDataSourceModel describes the data source data model.
type GroupnetDataSourceModel struct {
	Groupnets []GroupnetModel     `tfsdk:"groupnets"`
	ID        types.String        `tfsdk:"id"`
	Filter    *GroupnetFilterType `tfsdk:"filter"`
}

// GroupnetModel holds groupnet schema attribute details.
type GroupnetModel struct {
	AllowWildcardSubdomains types.Bool   `tfsdk:"allow_wildcard_subdomains"`
	Description             types.String `tfsdk:"description"`
	DNSCacheEnabled         types.Bool   `tfsdk:"dns_cache_enabled"`
	DNSResolverRotate       types.Bool   `tfsdk:"dns_resolver_rotate"`
	DNSSearch               types.List   `tfsdk:"dns_search"`
	DNSServers              types.List   `tfsdk:"dns_servers"`
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	ServerSideDNSSearch     types.Bool   `tfsdk:"server_side_dns_search"`
	Subnets                 types.List   `tfsdk:"subnets"`
}

// GroupnetFilterType holds filter attribute for groupnet.
type GroupnetFilterType struct {
	Names []types.String `tfsdk:"names"`
}
