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
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &S3BucketResource{}
	_ resource.ResourceWithImportState = &S3BucketResource{}
)

// NewS3BucketResource returns the S3 Bucket resource object.
func NewS3BucketResource() resource.Resource {
	return &S3BucketResource{}
}

// S3BucketResource defines the resource implementation.
type S3BucketResource struct {
	client *client.Client
}

// Configure configures the resource.
func (r *S3BucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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
func (r *S3BucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket"
}

// Schema describes the resource arguments.
func (r *S3BucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the S3 Bucket entity of PowerScale Array. PowerScale S3 Bucket map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Bucket using this resource. We can also import an existing S3 Bucket from PowerScale array.",
		Description:         "This resource is used to manage the S3 Bucket entity of PowerScale Array. PowerScale S3 Bucket map to the PowerScale file system as base directory for Objects. We can Create, Update and Delete the S3 Bucket using this resource. We can also import an existing S3 Bucket from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"acl": schema.ListNestedAttribute{
				Description:         "Specifies properties for an S3 Access Control Entry.",
				MarkdownDescription: "Specifies properties for an S3 Access Control Entry.",
				Computed:            true,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"grantee": schema.SingleNestedAttribute{
							Description:         "Specifies the persona of the file group.",
							MarkdownDescription: "Specifies the persona of the file group.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									MarkdownDescription: "Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "Specifies the persona name, which must be combined with a type.",
									MarkdownDescription: "Specifies the persona name, which must be combined with a type.",
									Required:            true,
								},
								"type": schema.StringAttribute{
									Description:         "Specifies the type of persona, which must be combined with a name.",
									MarkdownDescription: "Specifies the type of persona, which must be combined with a name.",
									Required:            true,
								},
							},
						},
						"permission": schema.StringAttribute{
							Description:         "Specifies the S3 rights being allowed.",
							MarkdownDescription: "Specifies the S3 rights being allowed.",
							Required:            true,
						},
					},
				},
			},
			"create_path": schema.BoolAttribute{
				Description:         "Create path if does not exist.",
				MarkdownDescription: "Create path if does not exist.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description for this S3 bucket.",
				MarkdownDescription: "Description for this S3 bucket.",
				Computed:            true,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Bucket ID.",
				MarkdownDescription: "Bucket ID.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Bucket name.",
				MarkdownDescription: "Bucket name.",
				Required:            true,
			},
			"object_acl_policy": schema.StringAttribute{
				Description:         "Set behavior of modifying object acls.",
				MarkdownDescription: "Set behavior of modifying object acls.",
				Computed:            true,
				Optional:            true,
			},
			"owner": schema.StringAttribute{
				Description:         "Specifies the name of the owner.",
				MarkdownDescription: "Specifies the name of the owner.",
				Computed:            true,
				Optional:            true,
			},
			"path": schema.StringAttribute{
				Description:         "Path of bucket within /ifs.",
				MarkdownDescription: "Path of bucket within /ifs.",
				Required:            true,
			},
			"zid": schema.Int64Attribute{
				Description:         "Zone ID.",
				MarkdownDescription: "Zone ID.",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description:         "Zone Name.",
				MarkdownDescription: "Zone Name.",
				Optional:            true,
			},
		},
	}
}

// Create allocates the resource.
func (r *S3BucketResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Info(ctx, "Creating S3 Bucket")

	var bucketPlan models.S3BucketResource
	diags := request.Plan.Get(ctx, &bucketPlan)

	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	bucketeToCreate := powerscale.V10S3Bucket{}
	err := helper.ReadFromState(ctx, bucketPlan, &bucketeToCreate)
	if err != nil {
		response.Diagnostics.AddError("Error creating s3 bucket",
			fmt.Sprintf("Could not read s3 bucket param of Path: %s with error: %s", bucketPlan.Path.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "creating s3 bucket", map[string]interface{}{
		"bucketeToCreate": bucketeToCreate,
	})

	bucket, err := helper.CreateS3Bucket(ctx, r.client, bucketeToCreate, bucketPlan.Zone.ValueString())
	if err != nil {
		errStr := constants.CreateS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating s3 bucket ",
			message)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("s3 bucket %s created", bucket.Id), map[string]interface{}{
		"s3BucketResponse": bucket,
	})

	getBucketResponse, err := helper.GetS3Bucket(ctx, r.client, bucket.Id, bucketPlan.Zone.ValueString())
	if err != nil {
		errStr := constants.GetS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error creating s3 bucket ",
			message)
		return
	}

	// update resource state according to response
	if len(getBucketResponse.Buckets) <= 0 {
		response.Diagnostics.AddError(
			"Error creating s3 bucket",
			fmt.Sprintf("Could not get created s3 bucket state %s with error: s3 bucket not found", bucket.Id),
		)
		return
	}
	createdBucket := getBucketResponse.Buckets[0]
	err = helper.CopyFieldsToNonNestedModel(ctx, createdBucket, &bucketPlan)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating s3 bucket",
			fmt.Sprintf("Could not read s3 bucket struct %s with error: %s", bucket.Id, err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, bucketPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Create s3 bucket completed")
}

// Read reads data from the resource.
func (r *S3BucketResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Info(ctx, "Reading S3 Bucket resource")
	var bucketState models.S3BucketResource
	diags := request.State.Get(ctx, &bucketState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	bucketID := bucketState.ID.ValueString()
	zone := bucketState.Zone.ValueString()
	tflog.Debug(ctx, "calling get S3 Bucket by ID", map[string]interface{}{
		"BucketID": bucketID,
		"Zone":     zone,
	})
	bucketResponse, err := helper.GetS3Bucket(ctx, r.client, bucketID, zone)
	if err != nil {
		errStr := constants.ReadS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error reading s3 bucket ",
			message)
		return
	}

	if len(bucketResponse.Buckets) <= 0 {
		response.Diagnostics.AddError(
			"Error reading s3 bucket",
			fmt.Sprintf("Could not read s3 bucket %s from pscale with error: s3 bucket not found", bucketID),
		)
		return
	}
	tflog.Debug(ctx, "updating read s3 bucket state", map[string]interface{}{
		"S3BucketResponse": bucketResponse,
		"S3BucketState":    bucketState,
	})
	err = helper.CopyFieldsToNonNestedModel(ctx, bucketResponse.Buckets[0], &bucketState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error read s3 bucket",
			fmt.Sprintf("Could not read s3 bucket struct %s with error: %s", bucketID, err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, bucketState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Read S3 Bucket completed")
}

// Update updates the resource state.
func (r *S3BucketResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating S3 Bucket")
	var bucketPlan models.S3BucketResource
	diags := request.Plan.Get(ctx, &bucketPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var bucketState models.S3BucketResource
	diags = response.State.Get(ctx, &bucketState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update s3 bucket", map[string]interface{}{
		"bucketPlan":  bucketPlan,
		"bucketState": bucketState,
	})

	bucketID := bucketState.ID.ValueString()

	// validate update params
	if err := helper.ValidateS3BucketUpdate(bucketPlan, bucketState); err != nil {
		response.Diagnostics.AddError(
			"Error updating s3 bucket",
			fmt.Sprintf("Could not update S3 Bucket %s with error: %s", bucketID, err.Error()),
		)
		return
	}

	var bucketToUpdate powerscale.V10S3BucketExtendedExtended
	// Get param from tf input
	err := helper.ReadFromState(ctx, bucketPlan, &bucketToUpdate)
	if err != nil {
		response.Diagnostics.AddError(
			"Error update s3 bucket",
			fmt.Sprintf("Could not read s3 bucket struct %s with error: %s", bucketID, err.Error()),
		)
		return
	}
	zoneName := bucketState.Zone.ValueString()
	err = helper.UpdateS3Bucket(ctx, r.client, bucketID, zoneName, bucketToUpdate)
	if err != nil {
		errStr := constants.UpdateS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating s3 bucket ",
			message)
		return
	}

	bucketID = bucketState.ID.ValueString()
	tflog.Debug(ctx, "calling get s3 bucket by ID on pscale client", map[string]interface{}{
		"s3BucketID": bucketID,
	})
	updatedBucket, err := helper.GetS3Bucket(ctx, r.client, bucketID, zoneName)
	if err != nil {
		errStr := constants.GetS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error updating s3 bucket ",
			message)
		return
	}

	if len(updatedBucket.Buckets) <= 0 {
		response.Diagnostics.AddError(
			"Error reading s3 bucket",
			fmt.Sprintf("Could not read s3 bucket %s from pscale with error: s3 bucket not found", bucketID),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, updatedBucket.Buckets[0], &bucketState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading s3 bucket",
			fmt.Sprintf("Could not read s3 bucket struct %s with error: %s", bucketID, err.Error()),
		)
		return
	}

	diags = response.State.Set(ctx, bucketState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Update S3 Bucket completed")
}

// Delete deletes the resource.
func (r S3BucketResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting S3 Bucket")
	var bucketState models.S3BucketResource
	diags := request.State.Get(ctx, &bucketState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	bucketID := bucketState.ID.ValueString()
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete s3 bucket on pscale client", map[string]interface{}{
		"s3BucketID": bucketID,
	})
	err := helper.DeleteS3Bucket(ctx, r.client, bucketID, bucketState.Zone.ValueString())
	if err != nil {
		errStr := constants.DeleteS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error deleting s3 bucket ",
			message)
		return
	}
	response.State.RemoveResource(ctx)
	tflog.Info(ctx, "Delete S3 Bucket completed")
}

// ImportState imports the resource state.
func (r S3BucketResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing S3 Bucket resource")
	var zoneName string
	bucketID := request.ID
	// request.ID is form of zoneName:bucketID
	if strings.Contains(request.ID, ":") {
		params := strings.Split(request.ID, ":")
		bucketID = strings.Trim(params[1], " ")
		zoneName = strings.Trim(params[0], " ")
	}

	bucketResponse, err := helper.GetS3Bucket(ctx, r.client, bucketID, zoneName)
	if err != nil {
		errStr := constants.ReadS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		response.Diagnostics.AddError("Error importing s3 bucket ",
			message)
		return
	}

	if len(bucketResponse.Buckets) <= 0 {
		response.Diagnostics.AddError(
			"Error importing s3 bucket",
			fmt.Sprintf("Could not read s3 bucket %s from pscale with error: s3 bucket not found", bucketID),
		)
		return
	}
	var bucketState models.S3BucketResource
	err = helper.CopyFieldsToNonNestedModel(ctx, bucketResponse.Buckets[0], &bucketState)
	if err != nil {
		response.Diagnostics.AddError(
			"Error importing s3 bucket",
			fmt.Sprintf("Could not read s3 bucket struct %s with error: %s", bucketID, err.Error()),
		)
		return
	}
	bucketState.Zone = types.StringNull()
	if len(zoneName) > 0 {
		bucketState.Zone = types.StringValue(zoneName)
	}
	response.Diagnostics.Append(response.State.Set(ctx, bucketState)...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Import S3 Bucket completed")
}
