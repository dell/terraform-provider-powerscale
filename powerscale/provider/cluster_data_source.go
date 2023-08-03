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
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ClusterDataSource{}

// NewClusterDataSource creates a new data source.
func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

// ClusterDataSource defines the data source implementation.
type ClusterDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema describes the data source arguments.
func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The cluster attributes and cluster node information.",
		Description:         "The cluster attributes and cluster node information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier of the cluster.",
				MarkdownDescription: "Unique identifier of the cluster.",
				Computed:            true,
			},
			"config":            helper.GetClusterConfigSchema(),
			"identity":          helper.GetClusterIdentitySchema(),
			"nodes":             helper.GetClusterNodeSchema(),
			"internal_networks": helper.GetClusterInternalNetworksSchema(),
			"acs":               helper.GetClusterAcsSchema(),
		},
	}
}

// Configure configures the data source.
func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ClusterDataSource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read and map cluster config data
	config, err := helper.GetClusterConfig(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster config, got error: %s", err))
		return
	}
	var dataConfig models.ClusterConfig
	err = helper.CopyFields(ctx, config, &dataConfig)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map cluster config fields", err.Error())
		return
	}
	data.Config = &dataConfig

	// Read and map cluster identity data
	identity, err := helper.GetClusterIdentity(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster identity, got error: %s", err))
		return
	}
	var dataIdentity models.ClusterIdentity
	err = helper.CopyFields(ctx, identity, &dataIdentity)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map cluster identity fields", err.Error())
		return
	}
	data.Identity = &dataIdentity

	nodes, err := helper.GetClusterNodes(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster nodes, got error: %s", err))
		return
	}
	var dataNodes models.ClusterNodes
	err = helper.CopyFields(ctx, nodes, &dataNodes)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map cluster nodes", err.Error())
		return
	}
	data.Nodes = &dataNodes

	networks, err := helper.GetClusterInternalNetworks(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster internal networks, got error: %s", err))
	}
	var dataNetworks models.ClusterInternalNetworks
	err = helper.CopyFields(ctx, networks, &dataNetworks)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map cluster internal networks", err.Error())
		return
	}
	data.InternalNetworks = &dataNetworks

	acs, err := helper.ListClusterAcs(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cluster acs, got error: %s", err))
	}
	var dataAcs models.ClusterAcs
	err = helper.CopyFields(ctx, acs, &dataAcs)
	if err != nil {
		resp.Diagnostics.AddError("Failed to map cluster acs", err.Error())
		return
	}
	data.ACS = &dataAcs

	data.ID = types.StringValue("cluster-data-source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
