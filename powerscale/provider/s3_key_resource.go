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
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &S3KeyResource{}
	// _ resource.ResourceWithImportState = &S3KeyResource{} (can't support import due to secret key not present during read-refresh response)
)

// NewS3KeyResource returns the S3 Bucket resource object.
func NewS3KeyResource() resource.Resource {
	return &S3KeyResource{}
}

// S3KeyResource defines the resource implementation.
type S3KeyResource struct {
	client *client.Client
}

// Configure configures the resource.
func (r *S3KeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (r *S3KeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_key"
}

// Schema describes the resource arguments.
func (r *S3KeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the S3 Key of PowerScale Array. PowerScale S3 keys are used to sign the requests you send to the S3 protocol. We can Create, Update and Delete the S3 Bucket using this resource.",
		Description:         "This resource is used to manage the S3 Key of PowerScale Array. PowerScale S3 keys are used to sign the requests you send to the S3 protocol. We can Create, Update and Delete the S3 Bucket using this resource.",
		Attributes:          S3KeyResourceSchema(),
	}
}

// S3KeyResourceSchema describe s3 key management schema
func S3KeyResourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"access_id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Unique identifier of the S3 key.",
			Description:         "Unique identifier of the S3 key.",
		},
		"user": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The username to create the S3 key. Required.",
			Description:         "The username to create the S3 key. Required.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIfConfigured(),
			},
		},
		"zone": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The zone of the user. Required.",
			Description:         "The zone of the user. Required.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIfConfigured(),
			},
		},
		"existing_key_expiry_time": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "The expiry of the old secret key in minutes. Optional, default is 0. It will be applicable only if old_secret_key is exist.",
			Description:         "The expiry of the old secret key in minutes. Optional, default is 0. It will be applicable only if old_secret_key is exist.",
		},
		"secret_key": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The secret key of the key. Computed.",
			Description:         "The secret key of the key. Computed.",
		},
		"secret_key_timestamp": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The timestamp of the secret key. Computed.",
			Description:         "The timestamp of the secret key. Computed.",
		},
		"old_secret_key": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The secret key of the old key. Computed.",
			Description:         "The secret key of the old key. Computed.",
		},
		"old_key_expiry": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The expiry of the old key. Computed.",
			Description:         "The expiry of the old key. Computed.",
		},
		"old_key_timestamp": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The timestamp of the old key. Computed.",
			Description:         "The timestamp of the old key. Computed.",
		},
	}
}

// Create allocates the resource.
func (r *S3KeyResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var s3key models.S3KeyResourceData
	diags := request.Plan.Get(ctx, &s3key)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	// call create s3key
	resp, err := helper.GenerateS3Key(ctx, r.client, s3key)
	if err != nil {
		response.Diagnostics.AddError("Error creating s3 key ", err.Error())
		return
	}
	helper.CopyFieldsToNonNestedModel(ctx, resp.Keys, &s3key)
	response.State.Set(ctx, s3key)
}

// Read reads data from the resource.
func (r *S3KeyResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {

	var s3key models.S3KeyResourceData
	diags := request.State.Get(ctx, &s3key)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// call get s3key
	resp, err := helper.GetS3Key(ctx, r.client, s3key)
	if err != nil {
		response.Diagnostics.AddError("Error reading s3 key ", err.Error())
		return
	}

	// precheck to invalidate the refresh
	errMsg := "[UNKNOWN KEY] Key Generated Outside of Terraform"
	if resp.Keys.GetOldKeyTimestamp() != int32(s3key.SecretKeyTimestamp.ValueInt64()) {
		s3key.OldSecretKey = types.StringValue(errMsg)
	} else {
		s3key.OldSecretKey = s3key.SecretKey
	}
	if resp.Keys.GetSecretKeyTimestamp() != int32(s3key.SecretKeyTimestamp.ValueInt64()) {
		response.Diagnostics.AddWarning(errMsg, errMsg)
		s3key.SecretKey = types.StringValue(errMsg)
	}
	helper.CopyFieldsToNonNestedModel(ctx, resp.Keys, &s3key)
	response.State.Set(ctx, s3key)
}

// Update updates the resource state.
func (r *S3KeyResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

	var s3key models.S3KeyResourceData
	diags := request.Plan.Get(ctx, &s3key)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var s3KeyState models.S3KeyResourceData
	diags = response.State.Get(ctx, &s3KeyState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	// call update s3key
	resp, err := helper.GenerateS3Key(ctx, r.client, s3key)
	if err != nil {
		response.Diagnostics.AddError("Error updating s3 key ", err.Error())
		return
	}
	if int64(resp.Keys.GetOldKeyTimestamp()) == s3KeyState.SecretKeyTimestamp.ValueInt64() {
		resp.Keys.SetOldSecretKey(s3KeyState.SecretKey.ValueString())
	}
	helper.CopyFieldsToNonNestedModel(ctx, resp.Keys, &s3key)
	response.State.Set(ctx, s3key)

}

// Delete deletes the resource.
func (r S3KeyResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {

	var s3key models.S3KeyResourceData
	diags := request.State.Get(ctx, &s3key)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	// call delete s3key
	err := helper.DeleteS3Key(ctx, r.client, s3key)
	if err != nil {
		response.Diagnostics.AddError("Error deleting s3 key ", err.Error())
		return
	}

	response.State.RemoveResource(ctx)
}

