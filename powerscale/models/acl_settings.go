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

// ACLSettingsResourceModel describes the resource data model.
type ACLSettingsResourceModel struct {
	// Access checks (chmod, chown).
	Access types.String `tfsdk:"access"`
	// Displayed mode bits.
	Calcmode types.String `tfsdk:"calcmode"`
	// Approximate group mode bits when ACL exists.
	CalcmodeGroup types.String `tfsdk:"calcmode_group"`
	// Approximate owner mode bits when ACL exists.
	CalcmodeOwner types.String `tfsdk:"calcmode_owner"`
	// Require traverse rights in order to traverse directories with existing ACLs.
	CalcmodeTraverse types.String `tfsdk:"calcmode_traverse"`
	// chmod on files with existing ACLs.
	Chmod types.String `tfsdk:"chmod"`
	// chmod (007) on files with existing ACLs.
	Chmod007 types.String `tfsdk:"chmod_007"`
	// ACLs created on directories by UNIX chmod.
	ChmodInheritable types.String `tfsdk:"chmod_inheritable"`
	// chown/chgrp on files with existing ACLs.
	Chown types.String `tfsdk:"chown"`
	// ACL creation over SMB.
	CreateOverSmb types.String `tfsdk:"create_over_smb"`
	//  Read only DOS attribute.
	DosAttr types.String `tfsdk:"dos_attr"`
	// Group owner inheritance.
	GroupOwnerInheritance types.String `tfsdk:"group_owner_inheritance"`
	// Treatment of 'rwx' permissions.
	Rwx types.String `tfsdk:"rwx"`
	// Synthetic 'deny' ACEs.
	SyntheticDenies types.String `tfsdk:"synthetic_denies"`
	// Access check (utimes)
	Utimes types.String `tfsdk:"utimes"`
}
