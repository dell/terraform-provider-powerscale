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

// NamespaceACLResourceModel describes the resource data model.
type NamespaceACLResourceModel struct {
	Namespace     types.String `tfsdk:"namespace"`
	Nsaccess      types.Bool   `tfsdk:"nsaccess"`
	Owner         types.Object `tfsdk:"owner"`
	Group         types.Object `tfsdk:"group"`
	Authoritative types.String `tfsdk:"authoritative"`
	Mode          types.String `tfsdk:"mode"`
	CustomACL     types.List   `tfsdk:"acl_custom"`
	ACL           types.List   `tfsdk:"acl"`
}
