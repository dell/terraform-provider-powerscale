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

package provider

import (
	"context"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ACLSettingsDataSource{}

// NewACLSettingsDataSource creates a new data source.
func NewACLSettingsDataSource() datasource.DataSource {
	return &ACLSettingsDataSource{}
}

// ACLSettingsDataSource defines the data source implementation.
type ACLSettingsDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *ACLSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aclsettings"
}

// Schema describes the data source arguments.
func (d *ACLSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the ACL Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use ACL Settings to manage file and directory permissions, referred to as access rights.",
		Description:         "This datasource is used to query the ACL Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. You can use ACL Settings to manage file and directory permissions, referred to as access rights.",
		Attributes: map[string]schema.Attribute{
			"access": schema.StringAttribute{
				Description:         "Access checks (chmod, chown). Options: unix, windows",
				MarkdownDescription: "Access checks (chmod, chown). Options: unix, windows",
				Computed:            true,
			},
			"calcmode": schema.StringAttribute{
				Description:         "Displayed mode bits. Options: approx, 777",
				MarkdownDescription: "Displayed mode bits. Options: approx, 777",
				Computed:            true,
			},
			"calcmode_group": schema.StringAttribute{
				Description:         "Approximate group mode bits when ACL exists. Options: group_aces, group_only",
				MarkdownDescription: "Approximate group mode bits when ACL exists. Options: group_aces, group_only",
				Computed:            true,
			},
			"calcmode_owner": schema.StringAttribute{
				Description:         "Approximate owner mode bits when ACL exists. Options: owner_aces, owner_only",
				MarkdownDescription: "Approximate owner mode bits when ACL exists. Options: owner_aces, owner_only",
				Computed:            true,
			},
			"calcmode_traverse": schema.StringAttribute{
				Description:         "Require traverse rights in order to traverse directories with existing ACLs. Options: require, ignore",
				MarkdownDescription: "Require traverse rights in order to traverse directories with existing ACLs. Options: require, ignore",
				Computed:            true,
			},
			"chmod": schema.StringAttribute{
				Description:         "chmod on files with existing ACLs. Options: remove, replace, replace_users_and_groups, merge_with_ugo_priority, merge, deny, ignore",
				MarkdownDescription: "chmod on files with existing ACLs. Options: remove, replace, replace_users_and_groups, merge_with_ugo_priority, merge, deny, ignore",
				Computed:            true,
			},
			"chmod_007": schema.StringAttribute{
				Description:         "chmod (007) on files with existing ACLs. Options: default, remove",
				MarkdownDescription: "chmod (007) on files with existing ACLs. Options: default, remove",
				Computed:            true,
			},
			"chmod_inheritable": schema.StringAttribute{
				Description:         "ACLs created on directories by UNIX chmod. Options: yes, no",
				MarkdownDescription: "ACLs created on directories by UNIX chmod. Options: yes, no",
				Computed:            true,
			},
			"chown": schema.StringAttribute{
				Description:         "chown/chgrp on files with existing ACLs. Options: owner_group_and_acl, owner_group_only, ignore",
				MarkdownDescription: "chown/chgrp on files with existing ACLs. Options: owner_group_and_acl, owner_group_only, ignore",
				Computed:            true,
			},
			"create_over_smb": schema.StringAttribute{
				Description:         "ACL creation over SMB. Options: allow, disallow",
				MarkdownDescription: "ACL creation over SMB. Options: allow, disallow",
				Computed:            true,
			},
			"dos_attr": schema.StringAttribute{
				Description:         " Read only DOS attribute. Options: deny_smb, deny_smb_and_nfs",
				MarkdownDescription: " Read only DOS attribute. Options: deny_smb, deny_smb_and_nfs",
				Computed:            true,
			},
			"group_owner_inheritance": schema.StringAttribute{
				Description:         "Group owner inheritance. Options: native, parent, creator",
				MarkdownDescription: "Group owner inheritance. Options: native, parent, creator",
				Computed:            true,
			},
			"rwx": schema.StringAttribute{
				Description:         "Treatment of 'rwx' permissions. Options: retain, full_control",
				MarkdownDescription: "Treatment of 'rwx' permissions. Options: retain, full_control",
				Computed:            true,
			},
			"synthetic_denies": schema.StringAttribute{
				Description:         "Synthetic 'deny' ACEs. Options: none, remove",
				MarkdownDescription: "Synthetic 'deny' ACEs. Options: none, remove",
				Computed:            true,
			},
			"utimes": schema.StringAttribute{
				Description:         "Access check (utimes). Options: only_owner, owner_and_write",
				MarkdownDescription: "Access check (utimes). Options: only_owner, owner_and_write",
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *ACLSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *ACLSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading acl settings data source")

	var state models.ACLSettingsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	aclSettingsResp, err := helper.GetACLSettings(ctx, d.client)

	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the acl settings",
			message,
		)
		return
	}

	aclSettings, err := helper.ACLSettingsDetailMapper(ctx, aclSettingsResp.AclPolicySettings)
	if err != nil {
		errStr := constants.ReadACLSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error mapping the list of acl settings",
			message,
		)
		return
	}

	state = aclSettings

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with reading acl settings data source")
}
