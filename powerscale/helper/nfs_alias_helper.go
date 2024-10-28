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

package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"

	"fmt"
	// "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	// "github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"
	// "terraform-provider-powerscale/powerscale/models"
)

// UpdateSyncIQGlobalSettings updates the SyncIQ global settings.
func CreateNfsAlias(ctx context.Context, client *client.Client, plan models.NfsAliasResourceModel, state *models.NfsAliasResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var toCreate powerscale.V2NfsAlias

	zoneParam := client.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv2NfsAlias(ctx)
	err := ReadFromState(ctx, &plan, &toCreate)
	if err != nil {
		errStr := constants.CreateNfsAliasErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error creating nfs alias",
			fmt.Sprintf("Could not create nfs alias with error: %s", message),
		)
		return diags
	}

	if !plan.Zone.IsNull() {
		zoneParam = zoneParam.Zone(plan.Zone.ValueString())
	}

	_, _, err2 := zoneParam.V2NfsAlias(toCreate).Execute()
	if err2 != nil {
		errStr := constants.CreateNfsAliasErrorMsg + "with error: "
		message := GetErrorString(err2, errStr)
		diags.AddError(
			"Error creating nfs alias",
			message,
		)
		return diags
	}

	checkParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv2NfsAliases(ctx)
	checkParam = checkParam.Check(true)
	nfsAliases, _, err := checkParam.Execute()
	if err != nil {
		errStr := constants.ReadNfsAliasErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading nfs alias",
			message,
		)
		return diags
	}

	var nfsAlias powerscale.V15NfsAliasExtended
	for _, v := range nfsAliases.Aliases {
		if *v.Name == toCreate.Name {
			nfsAlias = v
		}
	}

	err = CopyFields(ctx, nfsAlias, state)
	if err != nil {
		diags.AddError(
			"Error copying fields of nfs alias resource",
			err.Error(),
		)
		return diags
	}

	return diags
}

func ReadNfsAlias(ctx context.Context, client *client.Client, plan models.NfsAliasResourceModel, state *models.NfsAliasResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	checkParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv2NfsAliases(ctx)
	checkParam = checkParam.Check(true)
	nfsAliases, _, err := checkParam.Execute()
	if err != nil {
		errStr := constants.ReadNfsAliasErrorMsg + "with error: "
		message := GetErrorString(err, errStr)
		diags.AddError(
			"Error reading nfs alias",
			message,
		)
		return diags
	}

	var nfsAlias powerscale.V15NfsAliasExtended
	var found bool
	for _, v := range nfsAliases.Aliases {
		if *v.Name == plan.Name.ValueString() {
			nfsAlias = v
			found = true
			break
		}
	}

	if !found {
		diags.AddError(
			"nfs alias not found",
			fmt.Sprintf("error in finding the nfs alias with name %s", plan.Name.ValueString()),
		)
		return diags
	}

	err = CopyFields(ctx, nfsAlias, state)
	if err != nil {
		diags.AddError(
			"Error copying fields of nfs alias resource",
			err.Error(),
		)
		return diags
	}

	return diags
}
