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
	_ datasource.DataSource              = &NfsGlobalSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &NfsGlobalSettingsDataSource{}
)

// NewNfsGlobalSettingsDataSource creates a new cluster email settings data source.
func NewNfsGlobalSettingsDataSource() datasource.DataSource {
	return &NfsGlobalSettingsDataSource{}
}

// NfsGlobalSettingsDataSource defines the data source implementation.
type NfsGlobalSettingsDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NfsGlobalSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_global_settings"
}

// Schema describes the data source arguments.
func (d *NfsGlobalSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the NFS Global Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Description:         "This datasource is used to query the NFS Global Settings from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Id of NFS Global settings. Readonly. ",
				MarkdownDescription: "Id of NFS Global settings. Readonly. ",
			},
			"nfsv3_enabled": schema.BoolAttribute{
				Description:         "True if NFSv3 is enabled.",
				MarkdownDescription: "True if NFSv3 is enabled.",
				Computed:            true,
			},
			"nfsv3_rdma_enabled": schema.BoolAttribute{
				Description:         "True if the RDMA is enabled for NFSv3.",
				MarkdownDescription: "True if the RDMA is enabled for NFSv3.",
				Computed:            true,
			},
			"nfsv4_enabled": schema.BoolAttribute{
				Description:         "True if NFSv4 is enabled.",
				MarkdownDescription: "True if NFSv4 is enabled.",
				Computed:            true,
			},
			"rpc_maxthreads": schema.Int64Attribute{
				Description:         "Specifies the maximum number of threads in the nfsd thread pool.",
				MarkdownDescription: "Specifies the maximum number of threads in the nfsd thread pool.",
				Computed:            true,
			},
			"rpc_minthreads": schema.Int64Attribute{
				Description:         "Specifies the minimum number of threads in the nfsd thread pool.",
				MarkdownDescription: "Specifies the minimum number of threads in the nfsd thread pool.",
				Computed:            true,
			},
			"rquota_enabled": schema.BoolAttribute{
				Description:         "True if the rquota protocol is enabled.",
				MarkdownDescription: "True if the rquota protocol is enabled.",
				Computed:            true,
			},
			"service": schema.BoolAttribute{
				Description:         "True if the NFS service is enabled. When set to false, the NFS service is disabled.",
				MarkdownDescription: "True if the NFS service is enabled. When set to false, the NFS service is disabled.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *NfsGlobalSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *NfsGlobalSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Nfs Global Settings data source ")

	var settingsState models.NfsGlobalSettingsModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &settingsState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nfsGlobalSettings, err := helper.GetNfsGlobalSettings(ctx, d.client)

	if err != nil {
		errStr := constants.ReadNfsGlobalSettingsErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading nfs global settings",
			message,
		)
		return
	}

	err = helper.CopyFields(ctx, nfsGlobalSettings.GetSettings(), &settingsState)
	if err != nil {
		resp.Diagnostics.AddError("Error copying fields of nfs global settings datasource", err.Error())
		return
	}

	settingsState.ID = types.StringValue("nfs_global_settings")

	resp.Diagnostics.Append(resp.State.Set(ctx, &settingsState)...)
	tflog.Info(ctx, "Done with Read Nfs Global Settings data source ")
}
