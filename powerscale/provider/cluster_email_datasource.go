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

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ClusterEmailDataSource{}
	_ datasource.DataSourceWithConfigure = &ClusterEmailDataSource{}
)

// NewClusterEmailDataSource creates a new cluster email settings data source.
func NewClusterEmailDataSource() datasource.DataSource {
	return &ClusterEmailDataSource{}
}

// ClusterEmailDataSource defines the data source implementation.
type ClusterEmailDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *ClusterEmailDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_email"
}

// Schema describes the data source arguments.
func (d *ClusterEmailDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the Cluster Email Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Description:         "This datasource is used to query the Cluster Email Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier.",
				MarkdownDescription: "Unique identifier.",
				Computed:            true,
			},
			"settings": schema.SingleNestedAttribute{
				Description:         "Cluster email notification settings.",
				MarkdownDescription: "Cluster email notification settings.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"batch_mode": schema.StringAttribute{
						Description:         "This setting determines how notifications will be batched together to be sent by email.  'none' means each notification will be sent separately.  'severity' means notifications of the same severity will be sent together.  'category' means notifications of the same category will be sent together.  'all' means all notifications will be batched together and sent in a single email.",
						MarkdownDescription: "This setting determines how notifications will be batched together to be sent by email.  'none' means each notification will be sent separately.  'severity' means notifications of the same severity will be sent together.  'category' means notifications of the same category will be sent together.  'all' means all notifications will be batched together and sent in a single email.",
						Computed:            true,
					},
					"mail_relay": schema.StringAttribute{
						Description:         "The address of the SMTP server to be used for relaying the notification messages.  An SMTP server is required in order to send notifications.  If this string is empty, no emails will be sent.",
						MarkdownDescription: "The address of the SMTP server to be used for relaying the notification messages.  An SMTP server is required in order to send notifications.  If this string is empty, no emails will be sent.",
						Computed:            true,
					},
					"mail_sender": schema.StringAttribute{
						Description:         "The full email address that will appear as the sender of notification messages.",
						MarkdownDescription: "The full email address that will appear as the sender of notification messages.",
						Computed:            true,
					},
					"mail_subject": schema.StringAttribute{
						Description:         "The subject line for notification messages from this cluster.",
						MarkdownDescription: "The subject line for notification messages from this cluster.",
						Computed:            true,
					},
					"smtp_auth_passwd_set": schema.BoolAttribute{
						Description:         "Indicates if an SMTP authentication password is set.",
						MarkdownDescription: "Indicates if an SMTP authentication password is set.",
						Computed:            true,
					},
					"smtp_auth_security": schema.StringAttribute{
						Description:         "The type of secure communication protocol to use if SMTP is being used.  If 'none', plain text will be used, if 'starttls', the encrypted STARTTLS protocol will be used.",
						MarkdownDescription: "The type of secure communication protocol to use if SMTP is being used.  If 'none', plain text will be used, if 'starttls', the encrypted STARTTLS protocol will be used.",
						Computed:            true,
					},
					"smtp_auth_username": schema.StringAttribute{
						Description:         "Username to authenticate with if SMTP authentication is being used.",
						MarkdownDescription: "Username to authenticate with if SMTP authentication is being used.",
						Computed:            true,
					},
					"smtp_port": schema.Int64Attribute{
						Description:         "The port on the SMTP server to be used for relaying the notification messages.",
						MarkdownDescription: "The port on the SMTP server to be used for relaying the notification messages.",
						Computed:            true,
					},
					"use_smtp_auth": schema.BoolAttribute{
						Description:         "If true, this cluster will send SMTP authentication credentials to the SMTP relay server in order to send its notification emails.  If false, the cluster will attempt to send its notification emails without authentication.",
						MarkdownDescription: "If true, this cluster will send SMTP authentication credentials to the SMTP relay server in order to send its notification emails.  If false, the cluster will attempt to send its notification emails without authentication.",
						Computed:            true,
					},
					"user_template": schema.StringAttribute{
						Description:         "Location of a custom template file that can be used to specify the layout of the notification emails.",
						MarkdownDescription: "Location of a custom template file that can be used to specify the layout of the notification emails.",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *ClusterEmailDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ClusterEmailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Cluster Email Settings data source ")

	var state models.ClusterEmailDataSource

	if resp.Diagnostics.HasError() {
		return
	}

	clusterEmail, err := helper.GetClusterEmail(ctx, d.client)
	if err != nil {
		errStr := constants.ReadClusterEmailSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading cluster email",
			message,
		)
		return
	}

	err = helper.CopyFields(ctx, clusterEmail, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error copying fields of cluster email datasource", err.Error())
		return
	}
	state.ID = types.StringValue("cluster_email")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Cluster Email Settings data source ")
}
