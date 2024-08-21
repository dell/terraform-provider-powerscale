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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &S3GlobalSettingResource{}
	_ resource.ResourceWithImportState = &S3GlobalSettingResource{}
)

// NewS3GlobalSettingResource returns the S3 Global Setting resource object.
func NewS3GlobalSettingResource() resource.Resource {
	return &S3GlobalSettingResource{}
}

// S3GlobalSettingResource defines the resource implementation.
type S3GlobalSettingResource struct {
	client *client.Client
}

// Configure configures the resource.
func (r *S3GlobalSettingResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

// Metadata describes the resource arguments.
func (r *S3GlobalSettingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_global_settings"
}

// Schema describes the resource arguments.
func (r *S3GlobalSettingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the S3 Global Setting entity of PowerScale Array. PowerScale S3 Global Setting map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Global Setting using this resource. We can also import an existing S3 Global Setting from PowerScale array.",
		Description:         "This resource is used to manage the S3 Global Setting entity of PowerScale Array. PowerScale S3 Global Setting map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Global Setting using this resource. We can also import an existing S3 Global Setting from PowerScale array.",
		Attributes:          S3GlobalSettingResourceSchema(),
	}
}

func S3GlobalSettingResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"service": schema.BoolAttribute{
			Description:         "Specifies if the service is enabled.",
			MarkdownDescription: "Specifies if the service is enabled.",
			Optional:            true,
			Computed:            true,
		},
		"https_only": schema.BoolAttribute{
			Description:         "Specifies if the service is HTTPS only.",
			MarkdownDescription: "Specifies if the service is HTTPS only.",
			Optional:            true,
			Computed:            true,
		},
		"http_port": schema.Int64Attribute{
			Description:         "Specifies the HTTP port.",
			MarkdownDescription: "Specifies the HTTP port.",
			Optional:            true,
			Computed:            true,
		},
		"https_port": schema.Int64Attribute{
			Description:         "Specifies the HTTPS port.",
			MarkdownDescription: "Specifies the HTTPS port.",
			Optional:            true,
			Computed:            true,
		},
	}
}

// SetGlobalSetting updates the S3 Global Setting.
func SetGlobalSetting(ctx context.Context, client *client.Client, s3GSPlan models.S3GlobalSettingResource) (models.S3GlobalSettingResource, error) {
	var toUpdate powerscale.V10S3SettingsGlobalSettings
	err := helper.ReadFromState(ctx, &s3GSPlan, &toUpdate)
	if err != nil {
		return models.S3GlobalSettingResource{}, err
	}
	err = helper.UpdateS3GlobalSetting(ctx, client, toUpdate)
	if err != nil {
		return models.S3GlobalSettingResource{}, err
	}
	globalSettings, err := helper.GetS3GlobalSetting(ctx, client)
	if err != nil {
		return models.S3GlobalSettingResource{}, err
	}
	var state models.S3GlobalSettingResource
	err = helper.CopyFieldsToNonNestedModel(ctx, globalSettings.GetSettings(), &state)
	if err != nil {
		return models.S3GlobalSettingResource{}, err
	}
	return state, nil
}

// getGlobalSetting reads the S3 Global Setting.
func GetGlobalSetting(ctx context.Context, client *client.Client, s3GlobalSettingState models.S3GlobalSettingResource) error {
	globalSettings, err := helper.GetS3GlobalSetting(ctx, client)
	if err != nil {
		return err
	}
	err = helper.CopyFieldsToNonNestedModel(ctx, globalSettings.GetSettings(), &s3GlobalSettingState)
	if err != nil {
		return err
	}
	return nil
}

// Create allocates the resource.
func (r *S3GlobalSettingResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "Creating S3 Global Setting")

	var s3GSPlan models.S3GlobalSettingResource
	diags := request.Plan.Get(ctx, &s3GSPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	state, err := SetGlobalSetting(ctx, r.client, s3GSPlan)
	if err != nil {
		response.Diagnostics.AddError("Error creating s3 global setting",
			fmt.Sprintf("Could not create s3 global setting with error: %s", err.Error()),
		)
		return
	}
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Create s3 Global Setting completed")
}

// Read reads data from the resource.
func (r *S3GlobalSettingResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "Reading S3 Global Setting resource")
	var s3GlobalSettingState models.S3GlobalSettingResource
	diags := request.State.Get(ctx, &s3GlobalSettingState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := GetGlobalSetting(ctx, r.client, s3GlobalSettingState)
	if err != nil {
		response.Diagnostics.AddError("Error reading s3 global setting",
			fmt.Sprintf("Could not read s3 global setting with error: %s", err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, &s3GlobalSettingState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Read S3 Global Setting completed")
}

// Update updates the resource state.
func (r *S3GlobalSettingResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating S3 Global Setting")
	var s3GSPlan models.S3GlobalSettingResource
	diags := request.Plan.Get(ctx, &s3GSPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	state, err := SetGlobalSetting(ctx, r.client, s3GSPlan)
	if err != nil {
		response.Diagnostics.AddError("Error updating s3 global setting",
			fmt.Sprintf("Could not update s3 global setting with error: %s", err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Update S3 Global Setting completed")
}

// Delete deletes the resource.
func (r S3GlobalSettingResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting S3 Global Setting")
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "Delete S3 Global Setting completed")
}

// ImportState imports the resource state.
func (r S3GlobalSettingResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var s3GlobalSettingState models.S3GlobalSettingResource
	err := GetGlobalSetting(ctx, r.client, s3GlobalSettingState)
	if err != nil {
		response.Diagnostics.AddError("Error importing s3 global setting",
			fmt.Sprintf("Could not import s3 global setting with error: %s", err.Error()),
		)
		return
	}
	tflog.Info(ctx, " S3 Global Setting Import")
	diag := response.State.Set(ctx, s3GlobalSettingState)
	response.Diagnostics.Append(diag...)
	if response.Diagnostics.HasError() {
		return
	}
}
