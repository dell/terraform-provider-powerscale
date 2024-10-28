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
	"strconv"

	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// GetNFSAlias retrieve nfs alias information.
func GetNFSAlias(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) (*powerscale.V2NfsExportsExtended, error) {
	queryParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsExport(ctx, strconv.FormatInt(nfsModel.ID.ValueInt64(), 10))
	if !nfsModel.Zone.IsNull() {
		queryParam = queryParam.Zone(nfsModel.Zone.ValueString())
	}
	if !nfsModel.Scope.IsNull() {
		queryParam = queryParam.Scope(nfsModel.Scope.ValueString())
	}
	exportRes, _, err := queryParam.Execute()
	return exportRes, err
}

// ListNFSAliases list nfs alias entities.
func ListNFSAliases(ctx context.Context, client *client.Client, nfsFilter *models.NfsAliasDatasourceFilter) (*[]powerscale.V15NfsAliasExtended, error) {
	listNfsParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv2NfsAliases(ctx)
	if nfsFilter != nil {
		if !nfsFilter.Zone.IsNull() {
			listNfsParam = listNfsParam.Zone(nfsFilter.Zone.ValueString())
		}
		if !nfsFilter.Sort.IsNull() {
			listNfsParam = listNfsParam.Sort(nfsFilter.Sort.ValueString())
		}
		if !nfsFilter.Dir.IsNull() {
			listNfsParam = listNfsParam.Dir(nfsFilter.Dir.ValueString())
		}
		if !nfsFilter.Check.IsNull() {
			listNfsParam = listNfsParam.Check(nfsFilter.Check.ValueBool())
		}
		if !nfsFilter.Limit.IsNull() {
			listNfsParam = listNfsParam.Limit(int32(nfsFilter.Limit.ValueInt64()))
		}
	}
	NfsAliases, _, err := listNfsParam.Execute()
	if err != nil {
		return nil, err
	}
	totalNfsAliases := NfsAliases.Aliases
	for NfsAliases.Resume != nil && (nfsFilter == nil || nfsFilter.Limit.IsNull()) {
		resumeNfsParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv2NfsAliases(ctx).Resume(*NfsAliases.Resume)
		NfsAliases, _, err = resumeNfsParam.Execute()
		if err != nil {
			return &totalNfsAliases, err
		}
		totalNfsAliases = append(totalNfsAliases, NfsAliases.Aliases...)
	}
	return &totalNfsAliases, nil
}

// FilterAliases list nfs aliases entities.
func FilterAliases(paths []types.String, ids []types.String, exports []powerscale.V15NfsAliasExtended) ([]powerscale.V15NfsAliasExtended, error) {
	// if names are specified filter locally
	if len(paths) == 0 && len(ids) == 0 {
		return exports, nil
	}
	var idFilteredExports []powerscale.V15NfsAliasExtended
	if len(ids) == 0 {
		idFilteredExports = exports
	} else {
		idMap := make(map[string]powerscale.V15NfsAliasExtended)
		for _, export := range exports {
			idMap[*export.Name] = export
		}
		for _, id := range ids {
			if specifiedExport, ok := idMap[id.ValueString()]; ok {
				idFilteredExports = append(idFilteredExports, specifiedExport)
			}
		}
	}
	return idFilteredExports, nil
}
