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
	"regexp"
	"slices"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	powerscale "dell/powerscale-go-client"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &SyncIQPeerCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &SyncIQPeerCertificateDataSource{}
)

// NewSyncIQPeerCertificateDataSource creates a new peer certificate data source.
func NewSyncIQPeerCertificateDataSource() datasource.DataSource {
	return &SyncIQPeerCertificateDataSource{}
}

// SyncIQPeerCertificateDataSource defines the data source implementation.
type SyncIQPeerCertificateDataSource struct {
	client *client.Client
}

// Metadata describes the data source arguments.
func (d *SyncIQPeerCertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synciq_peer_certificate"
}

// Schema describes the data source arguments.
func (d *SyncIQPeerCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This datasource is used to query the existing SyncIQ Peer Certificates from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Description:         "This datasource is used to query the existing SyncIQ Peer Certificates from PowerScale array. The information fetched from this datasource can be used for getting the details or for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"certificates": schema.ListNestedAttribute{
				Description:         "List of certificates fetched.",
				MarkdownDescription: "List of certificates fetched.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Description field associated with a certificate provided for administrative convenience.",
							MarkdownDescription: "Description field associated with a certificate provided for administrative convenience.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(0, 2048),
							},
						},
						"fingerprints": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "A list of zero or more certificate fingerprints which can be used for certificate identification.",
							MarkdownDescription: "A list of zero or more certificate fingerprints which can be used for certificate identification.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed:            true,
										Description:         "Fingerprint hash algorithm",
										MarkdownDescription: "Fingerprint hash algorithm",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 100),
										},
									},
									"value": schema.StringAttribute{
										Computed:            true,
										Description:         "Fingerprint value",
										MarkdownDescription: "Fingerprint value",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 512),
										},
									},
								},
							},
						},
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Unique server certificate identifier.",
							MarkdownDescription: "Unique server certificate identifier.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 512),
							},
						},
						"issuer": schema.StringAttribute{
							Computed:            true,
							Description:         "Certificate issuer field extracted from the certificate.",
							MarkdownDescription: "Certificate issuer field extracted from the certificate.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 2048),
							},
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Administrator specified name identifier.",
							MarkdownDescription: "Administrator specified name identifier.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile("^[a-zA-Z0-9_-]*$"),
									"Only alphanumeric characters, hyphens and underscores are allowed",
								),
								stringvalidator.LengthBetween(0, 128),
							},
						},
						"not_after": schema.Int64Attribute{
							Computed:            true,
							Description:         "Certificate notAfter field extracted from the certificate encoded as a UNIX epoch timestamp.  The certificate is not valid after this timestamp.",
							MarkdownDescription: "Certificate notAfter field extracted from the certificate encoded as a UNIX epoch timestamp.  The certificate is not valid after this timestamp.",
						},
						"not_before": schema.Int64Attribute{
							Computed:            true,
							Description:         "Certificate notBefore field extracted from the certificate encoded as a UNIX epoch timestamp.  The certificate is not valid before this timestamp.",
							MarkdownDescription: "Certificate notBefore field extracted from the certificate encoded as a UNIX epoch timestamp.  The certificate is not valid before this timestamp.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							Description:         "Certificate validity status",
							MarkdownDescription: "Certificate validity status",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"valid",
									"invalid",
									"expired",
									"expiring",
								),
							},
						},
						"subject": schema.StringAttribute{
							Computed:            true,
							Description:         "Certificate subject field extracted from the certificate.",
							MarkdownDescription: "Certificate subject field extracted from the certificate.",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 2048),
							},
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "ID of the SyncIQ Peer Certificate to be fetched. If not provided, all the certificates will be fetched.",
				MarkdownDescription: "ID of the SyncIQ Peer Certificate to be fetched. If not provided, all the certificates will be fetched.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},

		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Description:         "Filters for fetching SyncIQ Peer Certificate.",
				MarkdownDescription: "Filters for fetching SyncIQ Peer Certificate.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Optional:            true,
						Description:         "Name of the SyncIQ Peer Certificate to be fetched.",
						MarkdownDescription: "Name of the SyncIQ Peer Certificate to be fetched.",
					},
				},
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("id")),
				},
			},
		},
	}
}

// Configure configures the data source.
func (d *SyncIQPeerCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the state of the syncIQ peer certificates from PowerScale.
func (d *SyncIQPeerCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read Terraform configuration data into the model
	var data models.PeerCertificateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *models.PeerCertificateDataSourceModel

	var errD error
	if id := data.ID.ValueString(); id == "" {
		config, err := helper.ListPeerCerts(ctx, d.client)
		if err != nil {
			message := helper.GetErrorString(err, "")
			resp.Diagnostics.AddError("Error reading syncIQ peer certificates", message)
			return
		}

		// Apply the Name Filter if it is set
		if data.PeerCertificateFilter != nil && data.PeerCertificateFilter.Name.ValueString() != "" {
			nameFilter := data.PeerCertificateFilter.Name.ValueString()
			fmt.Println("nameFilter:", nameFilter)
			config.Certificates = slices.DeleteFunc(config.Certificates, func(i powerscale.V16CertificatesSyslogCertificate) bool { return i.Name != nameFilter })
			if len(config.Certificates) == 0 {
				resp.Diagnostics.AddError("Error reading syncIQ peer certificate", fmt.Sprintf("Could not find syncIQ peer certificate with name %s", nameFilter))
				return
			}
		}

		state, errD = d.getState(ctx, config)

	} else {
		config, err := helper.ReadPeerCert(ctx, d.client, id)
		if err != nil {
			message := helper.GetErrorString(err, "")
			resp.Diagnostics.AddError("Error reading syncIQ peer certificate", message)
			return
		}
		state, errD = d.getState(ctx, config)
		state.ID = types.StringValue(id)
	}

	if errD != nil {
		resp.Diagnostics.AddError("Failed to map sync peer certificate fields", errD.Error())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (d *SyncIQPeerCertificateDataSource) getState(ctx context.Context, certs interface{}) (*models.PeerCertificateDataSourceModel, error) {
	ret := models.PeerCertificateDataSourceModel{
		ID: types.StringValue("dummy"),
	}
	err := helper.CopyFieldsToNonNestedModel(ctx, certs, &ret)
	return &ret, err
}
