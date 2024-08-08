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

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Unique identifier of the S3 key.",
			Description:         "Unique identifier of the S3 key.",
		},
		"user": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The username to create the S3 key. Required.",
			Description:         "The username to create the S3 key. Required.",
		},
		"zone": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The zone of the user. Required.",
			Description:         "The zone of the user. Required.",
		},
		"existing_key_expiry_time": schema.Int64Attribute{
			Optional:            true,
			Default:             int64default.StaticInt64(0),
			MarkdownDescription: "The expiry of the old secret key in minutes. Optional, default is 0.",
			Description:         "The expiry of the old secret key in minutes. Optional, default is 0.",
		},
		"access_id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The access id of the key. Computed.",
			Description:         "The access id of the key. Computed.",
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

}

// Read reads data from the resource.
func (r *S3KeyResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
}

// Update updates the resource state.
func (r *S3KeyResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

// Delete deletes the resource.
func (r S3KeyResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
}
