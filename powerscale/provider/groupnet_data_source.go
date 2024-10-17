/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &GroupnetDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupnetDataSource{}
)

// NewGroupnetDataSource creates a new Groupnet data source.
func NewGroupnetDataSource() datasource.DataSource {
	return &GroupnetDataSource{}
}

// GroupnetDataSource defines the data source implementation.
type GroupnetDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *GroupnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groupnet"
}

// Schema describes the data source arguments.
func (d *GroupnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing Groupnets from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.",
		Description:         "This datasource is used to query the existing Groupnets from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Groupnet sits above subnets and pools and allows separate Access Zones to contain distinct DNS settings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of the groupnet instance.",
				Description:         "Unique identifier of the groupnet instance.",
			},
			"groupnets": schema.ListNestedAttribute{
				MarkdownDescription: "List of groupnets.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "The name of the groupnet.",
							MarkdownDescription: "The name of the groupnet.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Unique Interface ID.",
							MarkdownDescription: "Unique Interface ID.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "A description of the groupnet.",
							MarkdownDescription: "A description of the groupnet.",
							Computed:            true,
						},
						"allow_wildcard_subdomains": schema.BoolAttribute{
							Description:         "If enabled, SmartConnect treats subdomains of known dns zones as the known dns zone. This is required for S3 Virtual Host domains.",
							MarkdownDescription: "If enabled, SmartConnect treats subdomains of known dns zones as the known dns zone. This is required for S3 Virtual Host domains.",
							Computed:            true,
						},
						"dns_cache_enabled": schema.BoolAttribute{
							Description:         "DNS caching is enabled or disabled.",
							MarkdownDescription: "DNS caching is enabled or disabled.",
							Computed:            true,
						},
						"server_side_dns_search": schema.BoolAttribute{
							Description:         "Enable or disable appending nodes DNS search list to client DNS inquiries directed at SmartConnect service IP.",
							MarkdownDescription: "Enable or disable appending nodes DNS search list to client DNS inquiries directed at SmartConnect service IP.",
							Computed:            true,
						},
						"dns_resolver_rotate": schema.BoolAttribute{
							Description:         "Enable or disable DNS resolver rotate.",
							MarkdownDescription: "Enable or disable DNS resolver rotate.",
							Computed:            true,
						},
						"dns_search": schema.ListAttribute{
							Description:         "List of DNS search suffixes.",
							MarkdownDescription: "List of DNS search suffixes.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"dns_servers": schema.ListAttribute{
							Description:         "List of Domain Name Server IP addresses.",
							MarkdownDescription: "List of Domain Name Server IP addresses.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"subnets": schema.ListAttribute{
							Description:         "Name of the subnets in the groupnet.",
							MarkdownDescription: "Name of the subnets in the groupnet.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.GroupnetFilterType{})),
			},
		},
	}
}

// Configure configures the data source.
func (d *GroupnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *GroupnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Groupnet data source ")

	var state models.GroupnetDataSourceModel
	var plan models.GroupnetDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groupnets, err := helper.GetAllGroupnets(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting the list of PowerScale Groupnets.", err.Error())
		return
	}

	// parse groupnet response to state groupnet model
	if err := helper.UpdateGroupnetDataSourceState(ctx, &plan, groupnets.Groupnets); err != nil {
		resp.Diagnostics.AddError("Error reading groupnets datasource plan",
			fmt.Sprintf("Could not list groupnets with error: %s", err.Error()))
		return
	}

	// filter groupnets
	if plan.Filter != nil {
		filterRes, err := helper.GetDataSourceByValue(ctx, *plan.Filter, groupnets.Groupnets)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not list groupnets with error:  %v", plan.Filter), err.Error(),
			)
		}

		groupnets.Groupnets = []powerscale.V10NetworkGroupnetExtended{}
		for _, v := range filterRes {
			groupnetCast := v.(powerscale.V10NetworkGroupnetExtended)
			groupnets.Groupnets = append(groupnets.Groupnets, groupnetCast)
		}
	}

	fulldetail := []models.GroupnetModel{}
	for _, vze := range groupnets.Groupnets {
		val := vze
		detail, err := helper.GroupnetDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadGroupnetErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError("Error reading Groupnet", message)
			return
		}

		fulldetail = append(fulldetail, detail)
	}
	state.Groupnets = fulldetail

	state.ID = types.StringValue("groupnet_datasource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Read Groupnet data source ")
}
