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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// ListQuotas list Quota entities.
func ListQuotas(ctx context.Context, client *client.Client, quotaFilter *models.QuotaDatasourceFilter) ([]powerscale.V12QuotaQuotaExtended, error) {
	listQuotaParam := client.PscaleOpenAPIClient.QuotaApi.ListQuotav12QuotaQuotas(ctx)
	if quotaFilter != nil {
		if !quotaFilter.Type.IsNull() {
			listQuotaParam = listQuotaParam.Type_(quotaFilter.Type.ValueString())
		}
		if !quotaFilter.Path.IsNull() {
			listQuotaParam = listQuotaParam.Path(quotaFilter.Path.ValueString())
		}
		if !quotaFilter.Zone.IsNull() {
			listQuotaParam = listQuotaParam.Zone(quotaFilter.Zone.ValueString())
		}
		if !quotaFilter.Enforced.IsNull() {
			listQuotaParam = listQuotaParam.Enforced(quotaFilter.Enforced.ValueBool())
		}
		if !quotaFilter.Exceeded.IsNull() {
			listQuotaParam = listQuotaParam.Exceeded(quotaFilter.Exceeded.ValueBool())
		}
		if !quotaFilter.IncludeSnapshots.IsNull() {
			listQuotaParam = listQuotaParam.IncludeSnapshots(quotaFilter.IncludeSnapshots.ValueBool())
		}
		if !quotaFilter.Persona.IsNull() {
			listQuotaParam = listQuotaParam.Persona(quotaFilter.Persona.ValueString())
		}
		if !quotaFilter.RecursePathChildren.IsNull() {
			listQuotaParam = listQuotaParam.RecursePathChildren(quotaFilter.RecursePathChildren.ValueBool())
		}
		if !quotaFilter.RecursePathParents.IsNull() {
			listQuotaParam = listQuotaParam.RecursePathParents(quotaFilter.RecursePathParents.ValueBool())
		}
		if !quotaFilter.ReportID.IsNull() {
			listQuotaParam = listQuotaParam.ReportId(quotaFilter.ReportID.ValueString())
		}
	}
	QuotasResponse, _, err := listQuotaParam.Execute()
	if err != nil {
		return nil, err
	}
	totalQuotas := QuotasResponse.Quotas
	for QuotasResponse.Resume != nil {
		resumeQuotaParam := client.PscaleOpenAPIClient.QuotaApi.ListQuotav12QuotaQuotas(ctx).Resume(*QuotasResponse.Resume)
		QuotasResponse, _, err = resumeQuotaParam.Execute()
		if err != nil {
			return totalQuotas, err
		}
		totalQuotas = append(totalQuotas, QuotasResponse.Quotas...)
	}
	return totalQuotas, nil
}
