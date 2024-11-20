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
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &S3BucketDataSource{}

// NewS3BucketDataSource creates a new s3 bucket.
func NewS3BucketDataSource() datasource.DataSource {
	return &S3BucketDataSource{}
}

// S3BucketDataSource defines the data source implementation.
type S3BucketDataSource struct {
	client *client.Client
}

// Configure configures the resource.
func (d *S3BucketDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Metadata describes the data source arguments.
func (d *S3BucketDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket"
}

// Schema describes the data source arguments.
func (d *S3BucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Access S3 Buckets. This datasource is used to query the existing S3 Bucket from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale S3 Bucket map to the PowerScale file system as base directory for Objects.",
		Description:         "Access S3 Buckets. This datasource is used to query the existing S3 Bucket from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block. PowerScale S3 Bucket map to the PowerScale file system as base directory for Objects.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"s3_buckets": schema.ListNestedAttribute{
				Description:         "List of S3 Buckets",
				MarkdownDescription: "List of S3 Buckets",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"acl": schema.ListNestedAttribute{
							Description:         "Specifies properties for an S3 Access Control Entry.",
							MarkdownDescription: "Specifies properties for an S3 Access Control Entry.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"grantee": schema.SingleNestedAttribute{
										Description:         "Specifies the persona of the file group.",
										MarkdownDescription: "Specifies the persona of the file group.",
										Computed:            true,
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
									"permission": schema.StringAttribute{
										Description:         "Specifies the S3 rights being allowed.",
										MarkdownDescription: "Specifies the S3 rights being allowed.",
										Computed:            true,
									},
								},
							},
						},
						"description": schema.StringAttribute{
							Description:         "Description for this S3 bucket.",
							MarkdownDescription: "Description for this S3 bucket.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Bucket ID.",
							MarkdownDescription: "Bucket ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "Bucket name.",
							MarkdownDescription: "Bucket name.",
							Computed:            true,
						},
						"object_acl_policy": schema.StringAttribute{
							Description:         "Set behavior of modifying object acls",
							MarkdownDescription: "Set behavior of modifying object acls",
							Computed:            true,
						},
						"owner": schema.StringAttribute{
							Description:         "Specifies the name of the owner.",
							MarkdownDescription: "Specifies the name of the owner.",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							Description:         "Path of bucket within /ifs.",
							MarkdownDescription: "Path of bucket within /ifs.",
							Computed:            true,
						},
						"zid": schema.Int64Attribute{
							Description:         "Zone ID",
							MarkdownDescription: "Zone ID",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"zone": schema.StringAttribute{
						Description:         "Specifies which access zone to use.",
						MarkdownDescription: "Specifies which access zone to use.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"owner": schema.StringAttribute{
						Description:         "Specifies the name of the owner.",
						MarkdownDescription: "Specifies the name of the owner.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		},
	}
}

// Read reads data from the data source.
func (d *S3BucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading S3 Bucket data source")
	var plan models.S3BucketDatasource
	var state models.S3BucketDatasource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	buckets, err := helper.ListS3Buckets(ctx, d.client, plan.S3BucketFilter)
	if err != nil {
		errStr := constants.ReadS3BucketErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error getting the list of s3 buckets",
			message,
		)
		return
	}

	var filteredBuckets []models.S3BucketDatasourceEntity
	// Convert from json to terraform model
	for _, bucket := range buckets {
		entity := models.S3BucketDatasourceEntity{}
		err := helper.CopyFields(ctx, bucket, &entity)
		if err != nil {
			resp.Diagnostics.AddError("Error reading s3 bucket datasource plan",
				fmt.Sprintf("Could not list s3 buckets with error: %s", err.Error()))
			return
		}
		filteredBuckets = append(filteredBuckets, entity)
	}

	if filteredBuckets == nil {
		resp.Diagnostics.AddError("Error reading s3 buckets", "No buckets found with the specified filter(s)")
	}

	state.ID = types.StringValue("s3_bucket_datasource")
	state.S3BucketFilter = plan.S3BucketFilter
	state.S3Buckets = filteredBuckets

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
