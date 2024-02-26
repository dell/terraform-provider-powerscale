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
var _ datasource.DataSource = &AccessZoneDataSource{}

// NewAccessZoneDataSource creates a new data source.
func NewAccessZoneDataSource() datasource.DataSource {
	return &AccessZoneDataSource{}
}

// AccessZoneDataSource defines the data source implementation.
type AccessZoneDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *AccessZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accesszone"
}

// Schema describes the data source arguments.
func (d *AccessZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Access Zone Datasource. This datasource is used to query the existing Access Zone from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Access Zones allow you to isolate data and control who can access data in each zone.",
		Description:         "Access Zone Datasource. This datasource is used to query the existing Access Zone from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale Access Zones allow you to isolate data and control who can access data in each zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"access_zones_details": schema.ListNestedAttribute{
				Description:         "List of AccessZones",
				MarkdownDescription: "List of AccessZones",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"alternate_system_provider": schema.StringAttribute{
							Description:         "Specifies an alternate system provider.",
							MarkdownDescription: "Specifies an alternate system provider.",
							Computed:            true,
						},
						"auth_providers": schema.ListAttribute{
							Description:         "Specifies the list of authentication providers available on this access zone.",
							MarkdownDescription: "Specifies the list of authentication providers available on this access zone.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"cache_entry_expiry": schema.Int64Attribute{
							Description:         "Specifies amount of time in seconds to cache a user/group.",
							MarkdownDescription: "Specifies amount of time in seconds to cache a user/group.",
							Computed:            true,
						},
						"create_path": schema.BoolAttribute{
							Description:         "Determines if a path is created when a path does not exist.",
							MarkdownDescription: "Determines if a path is created when a path does not exist.",
							Computed:            true,
						},
						"groupnet": schema.StringAttribute{
							Description:         "Groupnet identifier",
							MarkdownDescription: "Groupnet identifier",
							Computed:            true,
						},
						"home_directory_umask": schema.Int64Attribute{
							Description:         "Specifies the permissions set on automatically created user home directories.",
							MarkdownDescription: "Specifies the permissions set on automatically created user home directories.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Specifies the system-assigned ID for the access zone. This value is returned when an access zone is created through the POST method",
							MarkdownDescription: "Specifies the system-assigned ID for the access zone. This value is returned when an access zone is created through the POST method",
							Computed:            true,
						},
						"ifs_restricted": schema.ListNestedAttribute{
							Description:         "Specifies a list of users and groups that have read and write access to /ifs.",
							MarkdownDescription: "Specifies a list of users and groups that have read and write access to /ifs.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description:         "Specifies the persona name, which must be combined with a type.",
										MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										Description:         "Specifies the type of persona, which must be combined with a name.",
										MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
										Computed:            true,
									},
								},
							},
						},
						"map_untrusted": schema.StringAttribute{
							Description:         "Maps untrusted domains to this NetBIOS domain during authentication.",
							MarkdownDescription: "Maps untrusted domains to this NetBIOS domain during authentication.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Specifies the access zone name..",
							MarkdownDescription: "Specifies the access zone name.",
							Computed:            true,
						},
						"negative_cache_entry_expiry": schema.Int64Attribute{
							Description:         "Specifies number of seconds the negative cache entry is valid.",
							MarkdownDescription: "Specifies number of seconds the negative cache entry is valid.",
							Computed:            true,
						},
						"netbios_name": schema.StringAttribute{
							Description:         "Specifies the NetBIOS name.",
							MarkdownDescription: "Specifies the NetBIOS name.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							Description:         "Specifies the access zone base directory path.",
							MarkdownDescription: "Specifies the access zone base directory path.",
							Computed:            true,
						},
						"skeleton_directory": schema.StringAttribute{
							Description:         "Specifies the skeleton directory that is used for user home directories.",
							MarkdownDescription: "Specifies the skeleton directory that is used for user home directories.",
							Computed:            true,
						},
						"system": schema.BoolAttribute{
							Description:         "True if the access zone is built-in.",
							MarkdownDescription: "True if the access zone is built-in.",
							Computed:            true,
						},
						"system_provider": schema.StringAttribute{
							Description:         "Specifies the system provider for the access zone.",
							MarkdownDescription: "Specifies the system provider for the access zone.",
							Computed:            true,
						},
						"user_mapping_rules": schema.ListAttribute{
							Description:         "Specifies the current ID mapping rules.",
							MarkdownDescription: "Specifies the current ID mapping rules.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"zone_id": schema.Int64Attribute{
							Description:         "Specifies the access zone ID on the system.",
							MarkdownDescription: "Specifies the access zone ID on the system.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"names": schema.SetAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *AccessZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *AccessZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.AccessZoneDataSourceModel
	var plan models.AccessZoneDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
	result, err := helper.GetAllAccessZones(ctx, d.client)
	if err != nil {
		errStr := constants.ReadAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of access zones",
			message,
		)
		return
	}
	fulldetail := []models.AccessZoneDetailModel{}
	for _, vze := range result.Zones {
		val := vze
		detail, err := helper.AccessZoneDetailMapper(ctx, &val)
		if err != nil {
			errStr := constants.ReadAccessZoneErrorMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			resp.Diagnostics.AddError(
				"Error getting the list of access zones",
				message,
			)
			return
		}
		fulldetail = append(fulldetail, detail)
	}
	var validAccessZones []string
	if plan.AccessZoneFilter != nil && len(plan.AccessZoneFilter.Names) > 0 {
		for _, name := range plan.AccessZoneFilter.Names {
			for _, det := range fulldetail {
				if name.ValueString() == det.Name.ValueString() {
					state.AccessZones = append(state.AccessZones, det)
					validAccessZones = append(validAccessZones, det.ID.ValueString())
				}
			}
		}
		if len(state.AccessZones) != len(plan.AccessZoneFilter.Names) {
			resp.Diagnostics.AddError(
				"Error one or more of the filtered access zone names is not a valid powerscale access zone",
				fmt.Sprintf("Valid access zones [%v] filtered list [%v]", validAccessZones, plan.AccessZoneFilter.Names),
			)
		}
	} else {
		state.AccessZones = append(state.AccessZones, fulldetail...)
	}
	// save into the Terraform state.
	state.ID = types.StringValue("access_zone_datasource")

	tflog.Trace(ctx, "read the Access Zone datasource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
