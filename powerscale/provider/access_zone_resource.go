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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AccessZoneResource{}
var _ resource.ResourceWithImportState = &AccessZoneResource{}

// NewAccessZoneResource creates a new resource.
func NewAccessZoneResource() resource.Resource {
	return &AccessZoneResource{}
}

// AccessZoneResource defines the resource implementation.
type AccessZoneResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *AccessZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accesszone"
}

// Schema describes the resource arguments.
func (r *AccessZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Access Zone resource",

		Attributes: map[string]schema.Attribute{
			"alternate_system_provider": schema.StringAttribute{
				Description:         "Specifies an alternate system provider.",
				MarkdownDescription: "Specifies an alternate system provider.",
				Computed:            true,
			},
			"custom_auth_providers": schema.ListAttribute{
				Description:         "An optional parameter which adds new auth_providers to the access zone",
				MarkdownDescription: "An optional parameter which adds new auth_providers to the access zone",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             listdefault.StaticValue(basetypes.NewListNull(types.StringType)),
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
				Required:            true,
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
				Description:         "Specifies the access zone name.",
				MarkdownDescription: "Specifies the access zone name.",
				Required:            true,
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
				Required:            true,
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
	}
}

// Configure configures the resource.
func (r *AccessZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = pscaleClient
}

// Create allocates the resource.
func (r *AccessZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating access zone")
	var plan *models.AccessZoneResourceModel
	var authProv []string
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {

		return
	}

	resp.Diagnostics.Append(plan.CustomAuthProviders.ElementsAs(ctx, &authProv, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for i, v := range authProv {
		authProv[i] = "lsa-file-provider:" + v
	}

	err := helper.CreateAccessZones(ctx, r.client, authProv, plan)
	if err != nil {
		errStr := constants.CreateAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating access zone",
			message,
		)
		return
	}

	zoneParam, err := helper.GetAllAccessZones(ctx, r.client)
	if err != nil {
		errStr := constants.ReadAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading access zone",
			message,
		)
		return
	}

	state, err := helper.GetSpecificZone(ctx, plan.Name.ValueString(), zoneParam.Zones)

	if err != nil {
		resp.Diagnostics.AddError("Error reading access zone", err.Error())
		return
	}
	state.CustomAuthProviders = plan.CustomAuthProviders
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read reads the resource state.
func (r *AccessZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading access zone")
	var plan *models.AccessZoneResourceModel

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := helper.GetAllAccessZones(ctx, r.client)
	if err != nil {
		errStr := constants.ReadAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading access zone",
			message,
		)
		return
	}

	state, err := helper.GetSpecificZone(ctx, plan.Name.ValueString(), result.Zones)

	if err != nil {
		resp.Diagnostics.AddError("Error reading access zone", err.Error())
		return
	}

	state.CustomAuthProviders = plan.CustomAuthProviders
	// Save updated plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource state Path, Name, AuthProviders.
func (r *AccessZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating access zone")
	var plan *models.AccessZoneResourceModel
	var state *models.AccessZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read the state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {

		return
	}
	// Populate the Edit Parameters
	editValues := powerscale.V3ZoneExtendedExtended{}
	if state.Path != plan.Path {
		editValues.Path = plan.Path.ValueStringPointer()
	}
	if state.Name != plan.Name {
		editValues.Name = plan.Name.ValueStringPointer()
	}
	if len(state.CustomAuthProviders.Elements()) != len(plan.CustomAuthProviders.Elements()) {
		resp.Diagnostics.Append(plan.CustomAuthProviders.ElementsAs(ctx, &editValues.AuthProviders, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i, v := range editValues.AuthProviders {
			editValues.AuthProviders[i] = "lsa-file-provider:" + v
		}
	}
	editParam := r.client.PscaleOpenAPIClient.ZonesApi.UpdateZonesv3Zone(ctx, state.ID.ValueString())
	editParam = editParam.V3Zone(editValues)

	_, err := editParam.Execute()
	if err != nil {
		errStr := constants.UpdateAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error editing access zone",
			message,
		)
		return
	}

	result, err := helper.GetAllAccessZones(ctx, r.client)
	if err != nil {
		errStr := constants.ReadAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading access zone",
			message,
		)
		return
	}

	finalState, err := helper.GetSpecificZone(ctx, plan.Name.ValueString(), result.Zones)
	finalState.CustomAuthProviders = plan.CustomAuthProviders
	if err != nil {
		resp.Diagnostics.AddError("Error reading access zone", err.Error())
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalState)...)
}

// Delete deletes the resource.
func (r *AccessZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting access zone")
	var data *models.AccessZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteParam := r.client.PscaleOpenAPIClient.ZonesApi.DeleteZonesv3Zone(ctx, data.ID.ValueString())
	_, err := deleteParam.Execute()
	if err != nil {
		errStr := constants.DeleteAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of access zones",
			message,
		)
	}
}

// ImportState imports the resource state.
func (r *AccessZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing access zone")
	id := req.ID
	result, err := helper.GetAllAccessZones(ctx, r.client)
	if err != nil {
		errStr := constants.ReadAccessZoneErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading access zone",
			message,
		)
		return
	}
	state, err := helper.GetSpecificZone(ctx, id, result.Zones)
	if err != nil {
		resp.Diagnostics.AddError("Error reading access zone", err.Error())
		return
	}
	customAuthProviders, diags := helper.ExtractCustomAuthForInput(ctx, state.AuthProviders, state.Name.ValueString())

	if diags.HasError() {
		resp.Diagnostics.AddError("Unable to import access zone", "failed to extract custom providers")
	}
	state.CustomAuthProviders = customAuthProviders
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
