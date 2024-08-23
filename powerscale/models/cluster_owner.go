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

// ClusterOwner struct for ClusterOwner resource.
type ClusterOwner struct {
	ID              types.String `tfsdk:"id"`
	Company         types.String `tfsdk:"company"`
	Location        types.String `tfsdk:"location"`
	PrimaryEmail    types.String `tfsdk:"primary_email"`
	PrimaryName     types.String `tfsdk:"primary_name"`
	PrimaryPhone1   types.String `tfsdk:"primary_phone1"`
	PrimaryPhone2   types.String `tfsdk:"primary_phone2"`
	SecondaryEmail  types.String `tfsdk:"secondary_email"`
	SecondaryName   types.String `tfsdk:"secondary_name"`
	SecondaryPhone1 types.String `tfsdk:"secondary_phone1"`
	SecondaryPhone2 types.String `tfsdk:"secondary_phone2"`
}
