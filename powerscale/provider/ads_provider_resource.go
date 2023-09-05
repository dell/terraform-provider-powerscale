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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AdsProviderResource{}
var _ resource.ResourceWithConfigure = &AdsProviderResource{}
var _ resource.ResourceWithImportState = &AdsProviderResource{}

// NewAdsProviderResource creates a new resource.
func NewAdsProviderResource() resource.Resource {
	return &AdsProviderResource{}
}

// AdsProviderResource defines the resource implementation.
type AdsProviderResource struct {
	client *client.Client
}

// Metadata describes the resource arguments.
func (r *AdsProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adsprovider"
}

// Schema describes the resource arguments.
func (r *AdsProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ADS Provider resource",
		Description:         "ADS Provider resource",
		Attributes: map[string]schema.Attribute{
			"scope": schema.StringAttribute{
				Description:         "When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.",
				MarkdownDescription: "When specified as 'effective', or not specified, all fields are returned. When specified as 'user', only fields with non-default values are shown. When specified as 'default', the original values are returned.",
				Optional:            true,
			},
			"check_duplicates": schema.BoolAttribute{
				Description:         "Check for duplicate SPNs registered in Active Directory.",
				MarkdownDescription: "Check for duplicate SPNs registered in Active Directory.",
				Optional:            true,
			},
			"allocate_gids": schema.BoolAttribute{
				Description:         "Allocates an ID for an unmapped Active Directory (ADS) group. ADS groups without GIDs can be proactively assigned a GID by the ID mapper. If the ID mapper option is disabled, GIDs are not proactively assigned, and when a primary group for a user does not include a GID, the system may allocate one.",
				MarkdownDescription: "Allocates an ID for an unmapped Active Directory (ADS) group. ADS groups without GIDs can be proactively assigned a GID by the ID mapper. If the ID mapper option is disabled, GIDs are not proactively assigned, and when a primary group for a user does not include a GID, the system may allocate one.",
				Optional:            true,
				Computed:            true,
			},
			"allocate_uids": schema.BoolAttribute{
				Description:         "Allocates a user ID for an unmapped Active Directory (ADS) user. ADS users without UIDs can be proactively assigned a UID by the ID mapper. IF the ID mapper option is disabled, UIDs are not proactively assigned, and when an identify for a user does not include a UID, the system may allocate one.",
				MarkdownDescription: "Allocates a user ID for an unmapped Active Directory (ADS) user. ADS users without UIDs can be proactively assigned a UID by the ID mapper. IF the ID mapper option is disabled, UIDs are not proactively assigned, and when an identify for a user does not include a UID, the system may allocate one.",
				Optional:            true,
				Computed:            true,
			},
			"assume_default_domain": schema.BoolAttribute{
				Description:         "Enables lookup of unqualified user names in the primary domain.",
				MarkdownDescription: "Enables lookup of unqualified user names in the primary domain.",
				Optional:            true,
				Computed:            true,
			},
			"authentication": schema.BoolAttribute{
				Description:         "Enables authentication and identity management through the authentication provider.",
				MarkdownDescription: "Enables authentication and identity management through the authentication provider.",
				Optional:            true,
				Computed:            true,
			},
			"check_online_interval": schema.Int64Attribute{
				Description:         "Specifies the time in seconds between provider online checks.",
				MarkdownDescription: "Specifies the time in seconds between provider online checks.",
				Optional:            true,
				Computed:            true,
			},
			"controller_time": schema.Int64Attribute{
				Description:         "Specifies the current time for the domain controllers.",
				MarkdownDescription: "Specifies the current time for the domain controllers.",
				Optional:            true,
				Computed:            true,
			},
			"create_home_directory": schema.BoolAttribute{
				Description:         "Automatically creates a home directory on the first login.",
				MarkdownDescription: "Automatically creates a home directory on the first login.",
				Optional:            true,
				Computed:            true,
			},
			"dns_domain": schema.StringAttribute{
				Description:         "Specifies the DNS search domain. Set this parameter if the DNS search domain has a unique name or address.",
				MarkdownDescription: "Specifies the DNS search domain. Set this parameter if the DNS search domain has a unique name or address.",
				Optional:            true,
			},
			"domain_controller": schema.StringAttribute{
				Description:         "Specifies the domain controller to which the authentication service should send requests",
				MarkdownDescription: "Specifies the domain controller to which the authentication service should send requests",
				Optional:            true,
			},
			"domain_offline_alerts": schema.BoolAttribute{
				Description:         "Sends an alert if the domain goes offline.",
				MarkdownDescription: "Sends an alert if the domain goes offline.",
				Optional:            true,
				Computed:            true,
			},
			"dup_spns": schema.ListAttribute{
				Description:         "Get duplicate SPNs in the provider domain",
				MarkdownDescription: "Get duplicate SPNs in the provider domain",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"extra_expected_spns": schema.ListAttribute{
				Description:         "List of additional SPNs to expect beyond what automatic checking routines might find",
				MarkdownDescription: "List of additional SPNs to expect beyond what automatic checking routines might find",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"findable_groups": schema.ListAttribute{
				Description:         "Sets list of groups that can be resolved.",
				MarkdownDescription: "Sets list of groups that can be resolved.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"findable_users": schema.ListAttribute{
				Description:         "Sets list of users that can be resolved.",
				MarkdownDescription: "Sets list of users that can be resolved.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"forest": schema.StringAttribute{
				Description:         "Specifies the Active Directory forest.",
				MarkdownDescription: "Specifies the Active Directory forest.",
				Computed:            true,
			},
			"groupnet": schema.StringAttribute{
				Description:         "Groupnet identifier.",
				MarkdownDescription: "Groupnet identifier.",
				Optional:            true,
				Computed:            true,
			},
			"home_directory_template": schema.StringAttribute{
				Description:         "Specifies the path to the home directory template.",
				MarkdownDescription: "Specifies the path to the home directory template.",
				Optional:            true,
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				Description:         "Specifies the fully qualified hostname stored in the machine account.",
				MarkdownDescription: "Specifies the fully qualified hostname stored in the machine account.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Specifies the ID of the Active Directory provider instance.",
				MarkdownDescription: "Specifies the ID of the Active Directory provider instance.",
				Optional:            true,
				Computed:            true,
			},
			"ignore_all_trusts": schema.BoolAttribute{
				Description:         "If set to true, ignores all trusted domains.",
				MarkdownDescription: "If set to true, ignores all trusted domains.",
				Optional:            true,
				Computed:            true,
			},
			"ignored_trusted_domains": schema.ListAttribute{
				Description:         "Includes trusted domains when 'ignore_all_trusts' is set to false.",
				MarkdownDescription: "Includes trusted domains when 'ignore_all_trusts' is set to false.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"include_trusted_domains": schema.ListAttribute{
				Description:         "Includes trusted domains when 'ignore_all_trusts' is set to true.",
				MarkdownDescription: "Includes trusted domains when 'ignore_all_trusts' is set to true.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"instance": schema.StringAttribute{
				Description:         "Specifies Active Directory provider instance.",
				MarkdownDescription: "Specifies Active Directory provider instance.",
				Optional:            true,
			},
			"kerberos_hdfs_spn": schema.BoolAttribute{
				Description:         "Determines if connecting through HDFS with Kerberos.",
				MarkdownDescription: "Determines if connecting through HDFS with Kerberos.",
				Optional:            true,
			},
			"kerberos_nfs_spn": schema.BoolAttribute{
				Description:         "Determines if connecting through NFS with Kerberos.",
				MarkdownDescription: "Determines if connecting through NFS with Kerberos.",
				Optional:            true,
			},
			"ldap_sign_and_seal": schema.BoolAttribute{
				Description:         "Enables encryption and signing on LDAP requests.",
				MarkdownDescription: "Enables encryption and signing on LDAP requests.",
				Optional:            true,
				Computed:            true,
			},
			"login_shell": schema.StringAttribute{
				Description:         "Specifies the login shell path.",
				MarkdownDescription: "Specifies the login shell path.",
				Optional:            true,
				Computed:            true,
			},
			"lookup_domains": schema.ListAttribute{
				Description:         "Limits user and group lookups to the specified domains.",
				MarkdownDescription: "Limits user and group lookups to the specified domains.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"lookup_groups": schema.BoolAttribute{
				Description:         "Looks up AD groups in other providers before allocating a group ID.",
				MarkdownDescription: "Looks up AD groups in other providers before allocating a group ID.",
				Optional:            true,
				Computed:            true,
			},
			"lookup_normalize_groups": schema.BoolAttribute{
				Description:         "Normalizes AD group names to lowercase before look up.",
				MarkdownDescription: "Normalizes AD group names to lowercase before look up.",
				Optional:            true,
				Computed:            true,
			},
			"lookup_normalize_users": schema.BoolAttribute{
				Description:         "Normalize AD user names to lowercase before look up.",
				MarkdownDescription: "Normalize AD user names to lowercase before look up.",
				Optional:            true,
				Computed:            true,
			},
			"lookup_users": schema.BoolAttribute{
				Description:         "Looks up AD users in other providers before allocating a user ID.",
				MarkdownDescription: "Looks up AD users in other providers before allocating a user ID.",
				Optional:            true,
				Computed:            true,
			},
			"machine_account": schema.StringAttribute{
				Description:         "Specifies the machine account name when creating a SAM account with Active Directory.",
				MarkdownDescription: "Specifies the machine account name when creating a SAM account with Active Directory.",
				Optional:            true,
				Computed:            true,
			},
			"machine_password_changes": schema.BoolAttribute{
				Description:         "Enables periodic changes of the machine password for security.",
				MarkdownDescription: "Enables periodic changes of the machine password for security.",
				Optional:            true,
				Computed:            true,
			},
			"machine_password_lifespan": schema.Int64Attribute{
				Description:         "Sets maximum age of a password in seconds.",
				MarkdownDescription: "Sets maximum age of a password in seconds.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Specifies the Active Directory provider name.",
				MarkdownDescription: "Specifies the Active Directory provider name.",
				Required:            true,
			},
			"netbios_domain": schema.StringAttribute{
				Description:         "Specifies the NetBIOS domain name associated with the machine account.",
				MarkdownDescription: "Specifies the NetBIOS domain name associated with the machine account.",
				Computed:            true,
			},
			"node_dc_affinity": schema.StringAttribute{
				Description:         "Specifies the domain controller for which the node has affinity.",
				MarkdownDescription: "Specifies the domain controller for which the node has affinity.",
				Optional:            true,
				Computed:            true,
			},
			"node_dc_affinity_timeout": schema.Int64Attribute{
				Description:         "Specifies the timeout for the domain controller for which the local node has affinity.",
				MarkdownDescription: "Specifies the timeout for the domain controller for which the local node has affinity.",
				Optional:            true,
				Computed:            true,
			},
			"nss_enumeration": schema.BoolAttribute{
				Description:         "Enables the Active Directory provider to respond to 'getpwent' and 'getgrent' requests.",
				MarkdownDescription: "Enables the Active Directory provider to respond to 'getpwent' and 'getgrent' requests.",
				Optional:            true,
				Computed:            true,
			},
			"primary_domain": schema.StringAttribute{
				Description:         "Specifies the AD domain to which the provider is joined.",
				MarkdownDescription: "Specifies the AD domain to which the provider is joined.",
				Computed:            true,
			},
			"organizational_unit": schema.StringAttribute{
				Description:         "Specifies the organizational unit.",
				MarkdownDescription: "Specifies the organizational unit.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				Description:         "Specifies the password used during domain join.",
				MarkdownDescription: "Specifies the password used during domain join.",
				Required:            true,
				Sensitive:           true,
			},
			"reset_schannel": schema.BoolAttribute{
				Description:         "Resets the secure channel to the primary domain.",
				MarkdownDescription: "Resets the secure channel to the primary domain.",
				Optional:            true,
			},
			"restrict_findable": schema.BoolAttribute{
				Description:         "Check the provider for filtered lists of findable and unfindable users and groups.",
				MarkdownDescription: "Check the provider for filtered lists of findable and unfindable users and groups.",
				Optional:            true,
				Computed:            true,
			},
			"rpc_call_timeout": schema.Int64Attribute{
				Description:         "The maximum amount of time (in seconds) an RPC call to Active Directory is allowed to take.",
				MarkdownDescription: "The maximum amount of time (in seconds) an RPC call to Active Directory is allowed to take.",
				Optional:            true,
				Computed:            true,
			},
			"server_retry_limit": schema.Int64Attribute{
				Description:         "The number of retries attempted when a call to Active Directory fails due to network error.",
				MarkdownDescription: "The number of retries attempted when a call to Active Directory fails due to network error.",
				Optional:            true,
				Computed:            true,
			},
			"sfu_support": schema.StringAttribute{
				Description:         "Specifies whether to support RFC 2307 attributes on ADS domain controllers.",
				MarkdownDescription: "Specifies whether to support RFC 2307 attributes on ADS domain controllers.",
				Optional:            true,
				Computed:            true,
			},
			"site": schema.StringAttribute{
				Description:         "Specifies the site for the Active Directory.",
				MarkdownDescription: "Specifies the site for the Active Directory.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				Description:         "Specifies the status of the provider.",
				MarkdownDescription: "Specifies the status of the provider.",
				Computed:            true,
			},
			"store_sfu_mappings": schema.BoolAttribute{
				Description:         "Stores SFU mappings permanently in the ID mapper.",
				MarkdownDescription: "Stores SFU mappings permanently in the ID mapper.",
				Optional:            true,
				Computed:            true,
			},
			"system": schema.BoolAttribute{
				Description:         "If set to true, indicates that this provider instance was created by OneFS and cannot be removed.",
				MarkdownDescription: "If set to true, indicates that this provider instance was created by OneFS and cannot be removed.",
				Computed:            true,
			},
			"unfindable_groups": schema.ListAttribute{
				Description:         "Specifies groups that cannot be resolved by the provider.",
				MarkdownDescription: "Specifies groups that cannot be resolved by the provider.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"unfindable_users": schema.ListAttribute{
				Description:         "Specifies users that cannot be resolved by the provider.",
				MarkdownDescription: "Specifies users that cannot be resolved by the provider.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"user": schema.StringAttribute{
				Description:         "Specifies the user name that has permission to join a machine to the given domain.",
				MarkdownDescription: "Specifies the user name that has permission to join a machine to the given domain.",
				Required:            true,
			},
			"zone_name": schema.StringAttribute{
				Description:         "Specifies the name of the access zone in which this provider was created.",
				MarkdownDescription: "Specifies the name of the access zone in which this provider was created.",
				Computed:            true,
			},
			"recommended_spns": schema.ListAttribute{
				Description:         "Configuration recommended SPNs.",
				MarkdownDescription: "Configuration recommended SPNs.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"spns": schema.ListAttribute{
				Description:         "Currently configured SPNs.",
				MarkdownDescription: "Currently configured SPNs.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure configures the resource.
func (r *AdsProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	powerscaleClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = powerscaleClient
}

// Create allocates the resource.
func (r *AdsProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating ads provider")

	var plan models.AdsProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if helper.IsCreateAdsProviderParamInvalid(plan) {
		resp.Diagnostics.AddError(
			"Error creating ads provider",
			"Should not provide parameters for updating",
		)
		return
	}
	adsToCreate := powerscale.V14ProvidersAdsItem{}
	// Get param from tf input
	err := helper.ReadFromState(ctx, plan, &adsToCreate)
	if err != nil {
		errStr := constants.CreateAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ads provider",
			fmt.Sprintf("Could not read ads param with error: %s", message),
		)
		return
	}
	adsID, err := helper.CreateAdsProvider(ctx, r.client, adsToCreate)
	if err != nil {
		errStr := constants.CreateAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ads provider",
			message,
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("ads provider %s created", adsID.Id), map[string]interface{}{
		"adsProviderResponse": adsID,
	})

	plan.ID = types.StringValue(adsID.Id)
	getAdsResponse, err := helper.GetAdsProvider(ctx, r.client, plan)
	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ads provider",
			message,
		)
		return
	}

	// update resource state according to response
	if len(getAdsResponse.Ads) <= 0 {
		resp.Diagnostics.AddError(
			"Error creating ads provider",
			fmt.Sprintf("Could not read created ads provider %s", adsID),
		)
		return
	}

	createdAds := getAdsResponse.Ads[0]
	err = helper.CopyFieldsToNonNestedModel(ctx, createdAds, &plan)
	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error creating ads provider",
			fmt.Sprintf("Could not read ads provider struct %s with error: %s", adsID, message),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "create ads provider completed")
}

// Read reads the resource state.
func (r *AdsProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading ads provider")

	var adsState models.AdsProviderResourceModel
	diags := req.State.Get(ctx, &adsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	adsID := adsState.ID
	tflog.Debug(ctx, "calling get ads provider by ID", map[string]interface{}{
		"adsProviderID": adsID,
	})
	adsResponse, err := helper.GetAdsProvider(ctx, r.client, adsState)
	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading ads provider",
			message,
		)
		return
	}

	if len(adsResponse.Ads) <= 0 {
		resp.Diagnostics.AddError(
			"Error reading ads provider",
			fmt.Sprintf("Could not read ads provider %s from powerscale with error: ads provider not found", adsID),
		)
		return
	}
	tflog.Debug(ctx, "updating read ads provider state", map[string]interface{}{
		"adsProviderResponse": adsResponse,
		"adsProviderState":    adsState,
	})
	err = helper.CopyFieldsToNonNestedModel(ctx, adsResponse.Ads[0], &adsState)
	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error reading ads provider",
			fmt.Sprintf("Could not read ads provider struct %s with error: %s", adsID, message),
		)
		return
	}

	diags = resp.State.Set(ctx, adsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "read ads provider completed")
}

// Update updates the resource state.
func (r *AdsProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating ads provider")

	var adsPlan models.AdsProviderResourceModel
	diags := req.Plan.Get(ctx, &adsPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var adsState models.AdsProviderResourceModel
	diags = resp.State.Get(ctx, &adsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "calling update ads provider", map[string]interface{}{
		"adsPlan":  adsPlan,
		"adsState": adsState,
	})

	if helper.IsUpdateAdsProviderParamInvalid(adsPlan, adsState) {
		resp.Diagnostics.AddError(
			"Error updating ads provider",
			"Should not provide parameters for creating",
		)
		return
	}
	adsID := adsState.ID.ValueString()
	adsPlan.ID = adsState.ID
	var adsToUpdate powerscale.V14ProvidersAdsIdParams
	// Get param from tf input
	err := helper.ReadFromState(ctx, adsPlan, &adsToUpdate)
	if err != nil {
		errStr := constants.UpdateAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ads provider",
			fmt.Sprintf("Could not read ads param with error: %s", message),
		)
		return
	}
	err = helper.UpdateAdsProvider(ctx, r.client, adsID, adsToUpdate)
	if err != nil {
		errStr := constants.UpdateAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ads provider",
			message,
		)
		return
	}

	tflog.Debug(ctx, "calling get ads provider by ID on powerscale client", map[string]interface{}{
		"adsProviderID": adsID,
	})
	updatedAds, err := helper.GetAdsProvider(ctx, r.client, adsPlan)
	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ads provider",
			message,
		)
		return
	}

	if len(updatedAds.Ads) <= 0 {
		resp.Diagnostics.AddError(
			"Error updating ads provider",
			fmt.Sprintf("Could not read updated ads provider %s", adsID),
		)
		return
	}

	err = helper.CopyFieldsToNonNestedModel(ctx, updatedAds.Ads[0], &adsPlan)
	if err != nil {
		errStr := constants.ReadAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error updating ads provider",
			fmt.Sprintf("Could not read ads provider struct %s with error: %s", adsID, message),
		)
		return
	}
	diags = resp.State.Set(ctx, adsPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "update ads provider completed")
}

// Delete deletes the resource.
func (r *AdsProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting ads provider")

	var adsState models.AdsProviderResourceModel
	diags := req.State.Get(ctx, &adsState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	adsID := adsState.ID.ValueString()
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	tflog.Debug(ctx, "calling delete ads provider on powerscale client", map[string]interface{}{
		"adsProviderID": adsID,
	})
	err := helper.DeleteAdsProvider(ctx, r.client, adsID)
	if err != nil {
		errStr := constants.DeleteAdsProviderErrorMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		resp.Diagnostics.AddError(
			"Error deleting ads provider",
			message,
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "delete ads provider completed")
}

// ImportState imports the resource state.
func (r *AdsProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
