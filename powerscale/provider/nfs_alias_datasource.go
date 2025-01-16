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
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &NfsAliasDataSource{}
	_ datasource.DataSourceWithConfigure = &NfsAliasDataSource{}
)

// NewNfsAliasDataSource returns the NfsAlias data source object.
func NewNfsAliasDataSource() datasource.DataSource {
	return &NfsAliasDataSource{}
}

// NfsAliasDataSource defines the data source implementation.
type NfsAliasDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *NfsAliasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nfs_alias"
}

// Schema describes the data source arguments.
func (d *NfsAliasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This datasource is used to query the existing NFS aliases from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale provides an NFS server so you can share files on your cluster",
		Description: "This datasource is used to query the existing NFS aliases from PowerScale array. " +
			"The information fetched from this datasource can be used for getting the details or for further processing in resource block. " +
			"PowerScale provides an NFS server so you can share files on your cluster",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"nfs_aliases": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of nfs aliases",
				MarkdownDescription: "List of nfs aliases",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Name of NFS Alias.",
							MarkdownDescription: "Name of NFS Alias.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							Description:         "Path of NFS Alias.",
							MarkdownDescription: "Path of NFS Alias.",
							Computed:            true,
						},
						"zone": schema.StringAttribute{
							Description:         "Zone of NFS Alias.",
							MarkdownDescription: "Zone of NFS Alias.",
							Computed:            true,
						},
						"health": schema.StringAttribute{
							Description:         "Health status of NFS Alias.",
							MarkdownDescription: "Health status of NFS Alias.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "ID of NFS Alias.",
							MarkdownDescription: "ID of NFS Alias.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"ids": schema.SetAttribute{
						Description:         "IDs to filter nfs Aliases.",
						MarkdownDescription: "IDs to filter nfs Aliases.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"sort": schema.StringAttribute{
						Description:         "The field that will be used for sorting.",
						MarkdownDescription: "The field that will be used for sorting.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"zone": schema.StringAttribute{
						Description:         "Specifies which access zone to use.",
						MarkdownDescription: "Specifies which access zone to use.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"limit": schema.Int32Attribute{
						Description:         "Return no more than this many results at once (see resume).",
						MarkdownDescription: "Return no more than this many results at once (see resume).",
						Optional:            true,
					},
					"dir": schema.StringAttribute{
						Description:         "The direction of the sort.",
						MarkdownDescription: "The direction of the sort.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"check": schema.BoolAttribute{
						Description:         "Check for conflicts when listing Aliases.",
						MarkdownDescription: "Check for conflicts when listing Aliases.",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Configure configures the resource.
func (d *NfsAliasDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = pscaleClient
}

// Read reads data from the data source.
func (d *NfsAliasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Nfs Alias data source ")
	var aliasPlan models.NfsAliasDatasource
	var aliasState models.NfsAliasDatasource
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &aliasPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	totalNfsAliases, err := helper.ListNFSAliases(ctx, d.client, aliasPlan.NfsAliasesFilter)
	if err != nil {
		resp.Diagnostics.AddError("Error reading nfs alias datasource plan",
			fmt.Sprintf("Could not list nfs aliases with error: %s", err.Error()))
		return
	}
	var exportsIDs []types.String
	var paths []types.String
	if aliasPlan.NfsAliasesFilter != nil {
		exportsIDs = aliasPlan.NfsAliasesFilter.IDs

		filteredAliases, err := helper.FilterAliases(paths, exportsIDs, *totalNfsAliases)
		if err != nil {
			errStr := constants.ListNfsAliasErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error filtering nfs alias",
				message)
			return
		}

		var validAliases []string
		for _, export := range filteredAliases {
			entity := models.NfsAliasDatasourceEntity{}
			err := helper.CopyFields(ctx, export, &entity)
			if err != nil {
				resp.Diagnostics.AddError("Error reading nfs aliases datasource plan",
					fmt.Sprintf("Could not list nfs aliases with error: %s", err.Error()))
				return
			}
			aliasState.NfsAliases = append(aliasState.NfsAliases, entity)

			validAliases = append(validAliases, entity.ID.ValueString())
		}

		if len(aliasState.NfsAliases) < len(aliasPlan.NfsAliasesFilter.IDs) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered NFS Alias id is not a valid powerscale NFS Alias.",
				fmt.Sprintf("Valid NFS Aliases: [%v], filtered list: [%v]", strings.Join(validAliases, " , "), aliasPlan.NfsAliasesFilter.IDs),
			)
		}
	}

	aliasState.ID = types.StringValue("nfs_alias")
	aliasState.NfsAliasesFilter = aliasPlan.NfsAliasesFilter
	resp.Diagnostics.Append(resp.State.Set(ctx, &aliasState)...)
}
