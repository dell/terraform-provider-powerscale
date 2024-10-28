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

// NfsAliasResourceModel describes the resource data model.
type NfsAliasResourceModel struct {
	// Query param
	// Specifies which access zone to use.
	Zone   types.String `tfsdk:"zone"`
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Health types.String `tfsdk:"health"`
	Path   types.String `tfsdk:"path"`
}

// NfsAliasDatasource holds nfs aliases datasource schema attribute details.
type NfsAliasDatasource struct {
	ID               types.String               `tfsdk:"id"`
	NfsAliases       []NfsAliasDatasourceEntity `tfsdk:"nfs_aliases"`
	NfsAliasesFilter *NfsAliasDatasourceFilter  `tfsdk:"filter"`
}

// NfsAliasDatasourceFilter holds filter conditions.
type NfsAliasDatasourceFilter struct {
	// supported by api
	Sort  types.String `tfsdk:"sort"`
	Zone  types.String `tfsdk:"zone"`
	Limit types.Int64  `tfsdk:"limit"`
	Check types.Bool   `tfsdk:"check"`
	Dir   types.String `tfsdk:"dir"`
	// custom id & path list
	IDs []types.String `tfsdk:"ids"`
}

// NfsAliasDatasourceEntity Specifies entity values for NFS exports.
type NfsAliasDatasourceEntity struct {
	// Specifies the system-assigned ID for the export. This ID is returned when an export is created through the POST method.
	ID types.String `tfsdk:"id"`
	// Specifies the zone in which the export is valid.
	Zone   types.String `tfsdk:"zone"`
	Name   types.String `tfsdk:"name"`
	Path   types.String `tfsdk:"path"`
	Health types.String `tfsdk:"health"`
}
