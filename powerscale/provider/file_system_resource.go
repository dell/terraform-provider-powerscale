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
	"regexp"
	"terraform-provider-powerscale/client"

	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FileSystemResource{}
var _ resource.ResourceWithImportState = &FileSystemResource{}

// NewFileSystemResource creates a new data source.
func NewFileSystemResource() resource.Resource {
	return &FileSystemResource{}
}

// FileSystemResource defines the data source implementation.
type FileSystemResource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (r *FileSystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem"
}

// Schema describes the resource arguments.
func (r *FileSystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource is used to manage the FileSystem (Namespace directory) entity of PowerScale Array. We can Create, Update and Delete the FileSystem using this resource. We can also import an existing FileSystem from PowerScale array.",
		Description:         "This resource is used to manage the FileSystem (Namespace directory) entity of PowerScale Array. We can Create, Update and Delete the FileSystem using this resource. We can also import an existing FileSystem from PowerScale array.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "FileSystem identifier. Unique identifier for the FileSystem(Namespace directory)",
				MarkdownDescription: "FileSystem identifier. Unique identifier for the FileSystem(Namespace directory)",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				Description:         "FileSystem directory name",
				MarkdownDescription: "FileSystem directory name",
				Required:            true,
			},
			"full_path": schema.StringAttribute{
				Description:         "The full path of the FileSystem",
				MarkdownDescription: "The full path of the FileSystem",
				Computed:            true,
				Optional:            true,
			},
			"directory_path": schema.StringAttribute{
				Description:         "FileSystem directory path.This specifies the path to the FileSystem(Namespace directory) which we are trying to manage. If no directory path is specified, [/ifs] would be taken by default.",
				MarkdownDescription: "FileSystem directory path.This specifies the path to the FileSystem(Namespace directory) which we are trying to manage. If no directory path is specified, [/ifs] would be taken by default.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("/ifs"),
			},
			"query_zone": schema.StringAttribute{
				Description:         "Specifies the zone that the object belongs to. Optional and will default to the default access zone if one is not set.",
				MarkdownDescription: "Specifies the zone that the object belongs to. Optional and will default to the default access zone if one is not set.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				Description:         "File System Resource type",
				MarkdownDescription: "File System Resource type",
				Computed:            true,
			},
			"creation_time": schema.StringAttribute{
				Description:         "File System Resource Creation time",
				MarkdownDescription: "File System Resource Creation time",
				Computed:            true,
			},
			"owner": schema.SingleNestedAttribute{
				Description:         "The owner of the Filesystem.(Update Supported)",
				MarkdownDescription: "The owner of the Filesystem.(Update Supported)",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "Owner identifier",
						MarkdownDescription: "Owner identifier",
						Optional:            true,
						Computed:            true, Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(UID|SID):`), "must start with 'UID:' or 'SID:'",
							),
						},
					},
					"name": schema.StringAttribute{
						Description:         "Owner name",
						MarkdownDescription: "Owner name",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("id")),
						},
					},
					"type": schema.StringAttribute{
						Description:         "Owner type",
						MarkdownDescription: "Owner type",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("user"),
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"group": schema.SingleNestedAttribute{
				Description:         "The group of the Filesystem.(Update Supported)",
				MarkdownDescription: "The group of the Filesystem.(Update Supported)",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "group identifier",
						MarkdownDescription: "group identifier",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^(GID|SID):`), "must start with 'GID:' or 'SID:'",
							),
						},
					},
					"name": schema.StringAttribute{
						Description:         "group name",
						MarkdownDescription: "group name",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("id")),
						},
					},
					"type": schema.StringAttribute{
						Description:         "group type",
						MarkdownDescription: "group type",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("group"),
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"access_control": schema.StringAttribute{
				Description: `The ACL value for the directory. Users can either provide access rights input such as 'private_read' , 'private' ,
				'public_read', 'public_read_write', 'public' or permissions in POSIX format as '0550', '0770', '0775','0777' or 0700. The Default value is (0700).
				(Update Supported but Modification of ACL is only supported from POSIX to POSIX mode)`,
				MarkdownDescription: `The ACL value for the directory. Users can either provide access rights input such as 'private_read' , 'private' ,
				'public_read', 'public_read_write', 'public' or permissions in POSIX format as '0550', '0770', '0775','0777' or 0700. The Default value is (0700). 
				(Update Supported but Modification of ACL is only supported from POSIX to POSIX mode)`,
				Optional: true,
			},
			"authoritative": schema.StringAttribute{
				Description:         "If the directory has access rights set, then this field returns acl. Otherwise it returns mode.",
				MarkdownDescription: "If the directory has access rights set, then this field returns acl. Otherwise it returns mode.",
				Computed:            true,
			},
			"recursive": schema.BoolAttribute{
				Description:         "Creates intermediate folders recursively when set to true.",
				MarkdownDescription: "Creates intermediate folders recursively when set to true.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
			"overwrite": schema.BoolAttribute{
				Description:         "Deletes and replaces the existing user attributes and ACLs of the directory with user-specified attributes if set to true.",
				MarkdownDescription: "Deletes and replaces the existing user attributes and ACLs of the directory with user-specified attributes if set to true.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"mode": schema.StringAttribute{
				Description:         "Acl mode",
				MarkdownDescription: "Acl mode",
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (r *FileSystemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pscaleClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.client = pscaleClient
}

// Create creates the File system resource.
func (r *FileSystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating File System..")
	var plan models.FileSystemResource

	if resp.Diagnostics.HasError() {
		return
	}
	// Read Terraform plan data into the model
	resp.Diagnostics = append(resp.Diagnostics, req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dirPath := helper.GetDirectoryPath(plan.DirectoryPath.ValueString(), plan.Name.ValueString())

	createReq := r.client.PscaleOpenAPIClient.NamespaceApi.CreateDirectory(ctx, dirPath)

	createReq = createReq.XIsiIfsTargetType("container")
	createReq = createReq.Overwrite(plan.Overwrite.ValueBool())
	createReq = createReq.Recursive(plan.Recursive.ValueBool())
	if !plan.AccessControl.IsNull() && (plan.AccessControl.ValueString() != "") {
		createReq = createReq.XIsiIfsAccessControl(plan.AccessControl.ValueString())
	}

	_, _, errCR := helper.ExecuteCreate(createReq)
	if errCR != nil {
		errStr := constants.CreateFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(errCR, errStr)
		resp.Diagnostics.AddError("Error creating File System", message)
		return
	}

	if err := helper.UpdateFileSystemOwnerAndGroup(ctx, r.client, dirPath, &plan, &models.FileSystemResource{}); err != nil {
		resp.Diagnostics.AddWarning(fmt.Sprintf("Error setting the File system Resource - %s", dirPath), err.Error())
	}

	// Get File system metadata
	meta, err := helper.GetDirectoryMetadata(ctx, r.client, dirPath)
	if err != nil {
		errStr := constants.CreateFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the metadata for the filesystem", message)
		// if err, revert create
		if err = helper.DeleteFileSystem(ctx, r.client, dirPath); err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error deleting filesystem when reverting creation - %s", err.Error()))
		}
		return
	}

	// Get Acl
	acl, err := helper.GetDirectoryACL(ctx, r.client, dirPath)
	if err != nil {
		errStr := constants.CreateFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the acl for the filesystem", message)
		// if err, revert create
		if err = helper.DeleteFileSystem(ctx, r.client, dirPath); err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error deleting filesystem when reverting creation - %s", err.Error()))
		}
		return
	}

	// Update resource state
	resolveUID, _ := helper.ResolveOwnerGroupIdentity(ctx, r.client, plan.Owner.ID.ValueString(),
		plan.Owner.Name.ValueString(), plan.QueryZone.ValueString(), helper.DefaultIfEmpty(plan.Owner.Type.ValueString(), "user"))
	resolveGID, _ := helper.ResolveOwnerGroupIdentity(ctx, r.client, plan.Group.ID.ValueString(),
		plan.Group.Name.ValueString(), plan.QueryZone.ValueString(), helper.DefaultIfEmpty(plan.Group.Type.ValueString(), "group"))
	if diags := helper.UpdateFileSystemResourceState(ctx, &plan, acl, meta, resolveUID, resolveGID); diags.WarningsCount() > 0 {
		resp.Diagnostics.Append(diags...)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Done with Create File System resource")
}

// Read reads data from the resource.
func (r *FileSystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Read File System Resource..")
	var plan models.FileSystemResource

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
	dirPath := helper.GetDirectoryPath(plan.DirectoryPath.ValueString(), plan.Name.ValueString())

	// Get metadata
	meta, err := helper.GetDirectoryMetadata(ctx, r.client, dirPath)

	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the metadata for the filesystem", message)
		return
	}

	// GetAcl
	acl, err := helper.GetDirectoryACL(ctx, r.client, dirPath)
	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the acl for the filesystem", message)
		return
	}

	resolveUID, _ := helper.ResolveOwnerGroupIdentity(ctx, r.client, plan.Owner.ID.ValueString(),
		plan.Owner.Name.ValueString(), plan.QueryZone.ValueString(), helper.DefaultIfEmpty(plan.Owner.Type.ValueString(), "user"))
	resolveGID, _ := helper.ResolveOwnerGroupIdentity(ctx, r.client, plan.Group.ID.ValueString(),
		plan.Group.Name.ValueString(), plan.QueryZone.ValueString(), helper.DefaultIfEmpty(plan.Group.Type.ValueString(), "group"))
	if diags := helper.UpdateFileSystemResourceState(ctx, &plan, acl, meta, resolveUID, resolveGID); diags.WarningsCount() > 0 {
		resp.Diagnostics.Append(diags...)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Read File System Resource Complete.")
}

// Delete deletes the resource.
func (r *FileSystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting File System Resource..")
	var plan models.FileSystemResource

	// Read Terraform prior state plan into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
	dirPath := helper.GetDirectoryPath(plan.DirectoryPath.ValueString(), plan.Name.ValueString())
	if err := helper.DeleteFileSystem(ctx, r.client, dirPath); err != nil {
		resp.Diagnostics.AddError("Error Deleting filesystem", err.Error())
		return
	}
	tflog.Info(ctx, "Delete File system complete")
}

// Update updates the resource state.
func (r *FileSystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating File System.")
	// Read Terraform plan into the model
	var plan models.FileSystemResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state into the model
	var state models.FileSystemResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planDirName := helper.GetDirectoryPath(plan.DirectoryPath.ValueString(), plan.Name.ValueString())
	stateDirName := helper.GetDirectoryPath(state.DirectoryPath.ValueString(), state.Name.ValueString())
	if planDirName != stateDirName {
		resp.Diagnostics.AddError(constants.UpdateFileSystemErrorMsg, "Renaming Directory is not supported")
		return
	}

	if err := helper.UpdateFileSystemOwnerAndGroup(ctx, r.client, planDirName, &plan, &state); err != nil {
		resp.Diagnostics.AddWarning(fmt.Sprintf("Error updating the File system Resource - %s", planDirName), err.Error())
	}

	if err := helper.UpdateFileSystemAccessControl(ctx, r.client, planDirName, &plan, &state); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error updating the File system Resource - %s", planDirName), err.Error())
		return
	}

	// Get metadata
	meta, err := helper.GetDirectoryMetadata(ctx, r.client, planDirName)
	if err != nil {
		errStr := constants.UpdateFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddWarning("Error getting the metadata for the filesystem", message)
		return
	}

	// GetAcl
	acl, err := helper.GetDirectoryACL(ctx, r.client, planDirName)
	if err != nil {
		errStr := constants.UpdateFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the acl for the filesystem", message)
		return
	}

	// copy to model
	resolveUID, _ := helper.ResolveOwnerGroupIdentity(ctx, r.client, plan.Owner.ID.ValueString(),
		plan.Owner.Name.ValueString(), plan.QueryZone.ValueString(), helper.DefaultIfEmpty(plan.Owner.Type.ValueString(), "user"))
	resolveGID, _ := helper.ResolveOwnerGroupIdentity(ctx, r.client, plan.Group.ID.ValueString(),
		plan.Group.Name.ValueString(), plan.QueryZone.ValueString(), helper.DefaultIfEmpty(plan.Group.Type.ValueString(), "group"))
	if diags := helper.UpdateFileSystemResourceState(ctx, &plan, acl, meta, resolveUID, resolveGID); diags.WarningsCount() > 0 {
		resp.Diagnostics.Append(diags...)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Updating File System complete.")
}

// ImportState imports the resource state.
func (r *FileSystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing File System resource")
	var state models.FileSystemResource
	var id = req.ID
	// Get metadata
	meta, err := helper.GetDirectoryMetadata(ctx, r.client, id)
	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the metadata for the filesystem", message)
		return
	}

	// GetAcl
	acl, err := helper.GetDirectoryACL(ctx, r.client, id)
	if err != nil {
		errStr := constants.ReadFileSystemErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Error getting the acl for the filesystem", message)
		return
	}

	// copy to model
	helper.UpdateFileSystemResourceImportState(ctx, id, &state, acl, meta)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Done with Import File System resource")
}
