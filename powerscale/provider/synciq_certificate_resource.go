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
	"errors"
	"fmt"
	"strings"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SyncIQPeerCertificateResource{}
	_ resource.ResourceWithConfigure   = &SyncIQPeerCertificateResource{}
	_ resource.ResourceWithImportState = &SyncIQPeerCertificateResource{}
)

// NewSyncIQPeerCertificateResource creates a new resource.
func NewSyncIQPeerCertificateResource() resource.Resource {
	return &SyncIQPeerCertificateResource{}
}

// SyncIQPeerCertificateResource defines the resource implementation.
type SyncIQPeerCertificateResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *SyncIQPeerCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_certificate"
}

// Schema describes the resource arguments.
func (r *SyncIQPeerCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the SyncIQ Peer Certificate entity of PowerScale Array. " +
			"We can Create, Read, Update and Delete the SyncIQ Peer Certificate using this resource. We can also import existing SyncIQ Peer Certificate from PowerScale array.",
		Description: "This resource is used to manage the SyncIQ Peer Certificate entity of PowerScale Array. " +
			"We can Create, Read, Update and Delete the SyncIQ Peer Certificate using this resource. We can also import existing SyncIQ Peer Certificate from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "ID of the SyncIQ Peer certificate.",
				MarkdownDescription: "ID of the SyncIQ Peer certificate.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				Description: "Local path (on the PowerScale filesystem) to the certificate that is to be imported." +
					" This resource will be recreated if the value of this field is changed.",
				MarkdownDescription: "Local path (on the PowerScale filesystem) to the certificate that is to be imported." +
					" This resource will be recreated if the value of this field is changed.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Administrator specified name identifier.",
				MarkdownDescription: "Administrator specified name identifier.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description:         "Description field associated with a certificate provided for administrative convenience.",
				MarkdownDescription: "Description field associated with a certificate provided for administrative convenience.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *SyncIQPeerCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create allocates the resource.
func (r *SyncIQPeerCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan into the model
	var plan models.SyncIQPeerCertificateResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := helper.CreatePeerCert(ctx, r.client, powerscale.V7CertificateAuthorityItem{
		Name:            helper.GetKnownStringPointer(plan.Name),
		Description:     helper.GetKnownStringPointer(plan.Description),
		CertificatePath: plan.Path,
	})

	if err != nil {
		errStr := "Could not create syncIQ Peer Certificate with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Failed to create SyncIQ Peer Certificate", message)
		return
	}

	state, err := r.get(ctx, id, plan.Path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SyncIQ Peer Certificate after create", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// get state.
func (r *SyncIQPeerCertificateResource) get(ctx context.Context, id, path string) (models.SyncIQPeerCertificateResource, error) {
	certResp, err := helper.ReadPeerCert(ctx, r.client, id)
	if err != nil {
		errStr := "Could not read syncIQ Peer Certificate with error: "
		message := helper.GetErrorString(err, errStr)
		return models.SyncIQPeerCertificateResource{}, errors.New(message)
	}
	cert := certResp.Certificates[0]
	state := models.SyncIQPeerCertificateResource{
		ID:          types.StringValue(id),
		Name:        types.StringValue(cert.Name),
		Description: types.StringValue(cert.Description),
		Path:        path,
	}
	return state, nil
}

// Read reads the resource state.
func (r *SyncIQPeerCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var oldState models.SyncIQPeerCertificateResource
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, err := r.get(ctx, oldState.ID.ValueString(), oldState.Path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read SyncIQ Peer Certificate", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource state.
func (r *SyncIQPeerCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.SyncIQPeerCertificateResource
	var oldState models.SyncIQPeerCertificateResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update the certificate
	err := helper.UpdatePeerCert(ctx, r.client, oldState.ID.ValueString(), powerscale.V16CertificatesSyslogIdParams{
		Name:        helper.GetKnownStringPointer(plan.Name),
		Description: helper.GetKnownStringPointer(plan.Description),
	})
	if err != nil {
		errStr := "Could not update syncIQ Peer Certificate with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError("Failed to update SyncIQ Peer Certificate", message)
		return
	}

	// read the certificate
	state, err := r.get(ctx, oldState.ID.ValueString(), oldState.Path)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get SyncIQ Peer Certificate after update", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete deletes the resource.
func (r *SyncIQPeerCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SyncIQPeerCertificateResource

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// delete cert
	err := helper.DeletePeerCert(ctx, r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete SyncIQ Peer Certificate", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource.
func (r *SyncIQPeerCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Error importing syncIQ peer certificate", "Cannot import syncIQ peer certificate with empty ID")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), "/dummy")...)

	// if name is provided, use it
	if name, ok := strings.CutPrefix(req.ID, "name:"); ok {
		if name == "" {
			resp.Diagnostics.AddError("Error importing syncIQ peer certificate", "Cannot import syncIQ peer certificate with empty name, please use ID.")
			return
		}
		config, err := helper.ListPeerCerts(ctx, r.client)
		if err != nil {
			message := helper.GetErrorString(err, "")
			resp.Diagnostics.AddError("Error listing syncIQ peer certificates", message)
			return
		}
		for _, cert := range config.Certificates {
			if cert.Name == name {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), cert.Id)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), cert.Name)...)
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), cert.Description)...)
				return
			}
		}
		resp.Diagnostics.AddError("Could not find syncIQ peer certificate", fmt.Sprintf("Could not find syncIQ peer certificate with name %s", name))
	}

	// else use the ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
