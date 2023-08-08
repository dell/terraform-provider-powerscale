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

package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// GetNFSExport retrieve nfs export information.
func GetNFSExport(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) (*powerscale.V2NfsExportsExtended, error) {
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

// GetNFSExportByID retrieve nfs export information by id.
func GetNFSExportByID(ctx context.Context, client *client.Client, id string) (*powerscale.V2NfsExportsExtended, error) {
	queryParam := client.PscaleOpenAPIClient.ProtocolsApi.GetProtocolsv2NfsExport(ctx, id)
	exportRes, _, err := queryParam.Execute()
	return exportRes, err
}

// CreateNFSExport create nfs export.
func CreateNFSExport(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) (*powerscale.Createv3EventEventResponse, error) {
	nfsExport := powerscale.V2NfsExport{}
	err := ReadFromState(ctx, nfsModel, &nfsExport)
	if err != nil {
		return nil, err
	}
	createParam := client.PscaleOpenAPIClient.ProtocolsApi.CreateProtocolsv2NfsExport(ctx)
	if !nfsModel.Zone.IsNull() {
		createParam = createParam.Zone(nfsModel.Zone.ValueString())
	}
	if !nfsModel.Force.IsNull() {
		createParam = createParam.Force(nfsModel.Force.ValueBool())
	}
	if !nfsModel.IgnoreBadAuth.IsNull() {
		createParam = createParam.IgnoreBadAuth(nfsModel.IgnoreBadAuth.ValueBool())
	}
	if !nfsModel.IgnoreConflicts.IsNull() {
		createParam = createParam.IgnoreConflicts(nfsModel.IgnoreConflicts.ValueBool())
	}
	if !nfsModel.IgnoreUnresolvableHosts.IsNull() {
		createParam = createParam.IgnoreUnresolvableHosts(nfsModel.IgnoreUnresolvableHosts.ValueBool())
	}
	if !nfsModel.IgnoreBadPaths.IsNull() {
		createParam = createParam.IgnoreBadPaths(nfsModel.IgnoreBadPaths.ValueBool())
	}
	evenResp, _, err := createParam.V2NfsExport(nfsExport).Execute()
	return evenResp, err
}

// DeleteNFSExport create nfs export.
func DeleteNFSExport(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) error {
	deleteParam := client.PscaleOpenAPIClient.ProtocolsApi.DeleteProtocolsv2NfsExport(ctx, strconv.FormatInt(nfsModel.ID.ValueInt64(), 10))
	if !nfsModel.Zone.IsNull() {
		deleteParam = deleteParam.Zone(nfsModel.Zone.ValueString())
	}
	_, err := deleteParam.Execute()
	return err
}

// UpdateNFSExport update nfs export config.
func UpdateNFSExport(ctx context.Context, client *client.Client, nfsModel models.NfsExportResource) error {
	nfsExport := powerscale.V2NfsExportExtendedExtended{}
	err := ReadFromState(ctx, nfsModel, &nfsExport)
	if err != nil {
		return err
	}
	updateParam := client.PscaleOpenAPIClient.ProtocolsApi.UpdateProtocolsv2NfsExport(ctx, strconv.FormatInt(nfsModel.ID.ValueInt64(), 10))
	if !nfsModel.Zone.IsNull() {
		updateParam = updateParam.Zone(nfsModel.Zone.ValueString())
	}
	if !nfsModel.Force.IsNull() {
		updateParam = updateParam.Force(nfsModel.Force.ValueBool())
	}
	if !nfsModel.IgnoreBadAuth.IsNull() {
		updateParam = updateParam.IgnoreBadAuth(nfsModel.IgnoreBadAuth.ValueBool())
	}
	if !nfsModel.IgnoreConflicts.IsNull() {
		updateParam = updateParam.IgnoreConflicts(nfsModel.IgnoreConflicts.ValueBool())
	}
	if !nfsModel.IgnoreUnresolvableHosts.IsNull() {
		updateParam = updateParam.IgnoreUnresolvableHosts(nfsModel.IgnoreUnresolvableHosts.ValueBool())
	}
	if !nfsModel.IgnoreBadPaths.IsNull() {
		updateParam = updateParam.IgnoreBadPaths(nfsModel.IgnoreBadPaths.ValueBool())
	}
	_, err = updateParam.V2NfsExport(nfsExport).Execute()
	return err
}

// ListNFSExports list nfs export entities.
func ListNFSExports(ctx context.Context, client *client.Client, nfsFilter *models.NfsExportDatasourceFilter) (*[]powerscale.V2NfsExportExtended, error) {
	listNfsParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv4NfsExports(ctx)
	if nfsFilter != nil {
		listNfsParam.Resume(nfsFilter.Resume.ValueString())
		listNfsParam.Zone(nfsFilter.Zone.ValueString())
		listNfsParam.Scope(nfsFilter.Scope.ValueString())
		listNfsParam.Sort(nfsFilter.Sort.ValueString())
		listNfsParam.Path(nfsFilter.Path.ValueString())
		listNfsParam.Dir(nfsFilter.Dir.ValueString())
		listNfsParam.Check(nfsFilter.Check.ValueBool())
		if !nfsFilter.Limit.IsNull() {
			listNfsParam.Limit(int32(nfsFilter.Limit.ValueInt64()))
		}
		if !nfsFilter.Offset.IsNull() {
			listNfsParam.Offset(int32(nfsFilter.Offset.ValueInt64()))
		}
	}
	NfsExports, _, err := listNfsParam.Execute()
	if err != nil {
		return nil, err
	}
	totalNfsExports := NfsExports.Exports
	for NfsExports.Resume != nil && (nfsFilter == nil || nfsFilter.Limit.IsNull()) {
		resumeNfsParam := client.PscaleOpenAPIClient.ProtocolsApi.ListProtocolsv4NfsExports(ctx).Resume(*NfsExports.Resume)
		NfsExports, _, err = resumeNfsParam.Execute()
		if err != nil {
			return &totalNfsExports, err
		}
		totalNfsExports = append(totalNfsExports, NfsExports.Exports...)
	}
	return &totalNfsExports, nil
}

// FilterExports list nfs export en tities.
func FilterExports(paths []types.String, ids []types.Int64, exports *[]powerscale.V2NfsExportExtended) (*[]powerscale.V2NfsExportExtended, error) {
	// if names are specified filter locally
	if len(paths) == 0 && len(ids) == 0 {
		return exports, nil
	}
	var idFilteredExports []powerscale.V2NfsExportExtended
	idMap := make(map[int64]powerscale.V2NfsExportExtended)
	for _, export := range *exports {
		idMap[*export.Id] = export
	}
	for _, id := range ids {
		if specifiedExport, ok := idMap[id.ValueInt64()]; ok {
			idFilteredExports = append(idFilteredExports, specifiedExport)
		}
	}
	// filter path
	pathMap := make(map[string]bool)
	for _, path := range paths {
		pathMap[path.ValueString()] = true
	}
	var filteredExports []powerscale.V2NfsExportExtended
	for _, export := range idFilteredExports {
		for _, exportPath := range export.Paths {
			if pathMap[exportPath] {
				filteredExports = append(filteredExports, export)
			}
		}
	}

	return &filteredExports, nil
}
