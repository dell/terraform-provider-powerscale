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

package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CreateQuota creates quota.
func CreateQuota(ctx context.Context, client *client.Client, quota powerscale.V12QuotaQuota, zone string) (*powerscale.CreateResponse, error) {
	param := client.PscaleOpenAPIClient.QuotaApi.CreateQuotav12QuotaQuota(ctx).V12QuotaQuota(quota)
	if zone != "" {
		param = param.Zone(zone)
	}
	response, _, err := param.Execute()
	return response, err
}

// GetQuota gets quota.
func GetQuota(ctx context.Context, client *client.Client, quotaID string, zone string) (*powerscale.V12QuotaQuotasExtended, error) {
	param := client.PscaleOpenAPIClient.QuotaApi.GetQuotav12QuotaQuota(ctx, quotaID)
	if zone != "" {
		param = param.Zone(zone)
	}
	response, _, err := param.ResolveNames(true).Execute()
	return response, err
}

// UpdateQuota updates quota.
func UpdateQuota(ctx context.Context, client *client.Client, quotaID string, updatedQuota powerscale.V12QuotaQuotaExtendedExtended, linked bool) error {
	doLink := false
	if updatedQuota.Linked != nil && linked != *updatedQuota.Linked {
		if *updatedQuota.Linked {
			// if linking quota, need to update rest part at first
			doLink = true
		} else {
			err := LinkQuota(ctx, client, quotaID, !linked)
			if err != nil {
				return err
			}
		}
	}
	updatedQuota.Linked = nil
	param := client.PscaleOpenAPIClient.QuotaApi.UpdateQuotav12QuotaQuota(ctx, quotaID).V12QuotaQuota(updatedQuota)
	_, err := param.Execute()
	if err != nil {
		return err
	}
	if doLink {
		err = LinkQuota(ctx, client, quotaID, !linked)
	}
	return err
}

// LinkQuota link or unlink quota.
func LinkQuota(ctx context.Context, client *client.Client, quotaID string, linked bool) error {
	param := client.PscaleOpenAPIClient.QuotaApi.UpdateQuotav12QuotaQuota(ctx, quotaID).V12QuotaQuota(powerscale.V12QuotaQuotaExtendedExtended{Linked: &linked})
	_, err := param.Execute()
	return err
}

// DeleteQuota deletes quota.
func DeleteQuota(ctx context.Context, client *client.Client, quotaID string) error {
	param := client.PscaleOpenAPIClient.QuotaApi.DeleteQuotav12QuotaQuota(ctx, quotaID)
	_, err := param.Execute()
	return err
}

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

// ValidateQuotaUpdate validates if update params contain params only for creating.
func ValidateQuotaUpdate(plan models.QuotaResource, state models.QuotaResource) error {
	if !plan.Zone.IsNull() && !plan.Zone.Equal(state.Zone) && (plan.Zone.ValueString() != "System" || !state.Zone.IsNull()) {
		return fmt.Errorf("do not update field Zone")
	}
	if !plan.Path.IsNull() && !plan.Path.Equal(state.Path) {
		return fmt.Errorf("do not update field Path")
	}
	if !plan.Type.IsNull() && !plan.Type.Equal(state.Type) {
		return fmt.Errorf("do not update field Type")
	}
	if !plan.IncludeSnapshots.IsNull() && !plan.IncludeSnapshots.Equal(state.IncludeSnapshots) {
		return fmt.Errorf("do not update field IncludeSnapshots")
	}
	if !plan.Persona.IsNull() && !plan.Persona.Equal(state.Persona) {
		return fmt.Errorf("do not update field Persona.ID")
	}
	return nil
}

// IsQuotaParamInvalid Verify if persona and zone params are valid for different quota type.
func IsQuotaParamInvalid(plan models.QuotaResource) error {
	quotaType := plan.Type.ValueString()
	switch quotaType {
	case "user", "group":
		if plan.Persona.IsNull() || plan.Persona.IsUnknown() {
			return fmt.Errorf("\"persona\" is required for %s type", quotaType)
		}
	case "directory":
		if !plan.Persona.IsNull() {
			return fmt.Errorf("\"persona\" is not needed for %s type", quotaType)
		}
		if !plan.Zone.IsNull() {
			return fmt.Errorf("\"zone\" is not needed for %s type", quotaType)
		}
	default:
		return fmt.Errorf("unsupported type: %s", quotaType)
	}
	return nil
}

// RemovePercentThreshold reset percent_advisory and percent_soft data from response data
// Since PowerScale rest api directly use actual number for percentage calculation and accuracy loss might occur
// threshold data in response might be different from request body
// Thus, computed data should not include percent threshold when accurate threshold is given
func RemovePercentThreshold(ctx context.Context, prior, updated types.Object) (types.Object, diag.Diagnostics) {
	// no comparison need in this case
	if prior.IsNull() || updated.IsNull() {
		return updated, diag.Diagnostics{}
	}
	thresholdValue := updated.Attributes()
	priorThresholdValue := prior.Attributes()
	if adv, ok := thresholdValue["percent_advisory"]; !adv.IsNull() && ok {
		if _, ok := priorThresholdValue["percent_advisory"]; ok {
			thresholdValue["percent_advisory"] = priorThresholdValue["percent_advisory"]
		}
	}

	if soft, ok := thresholdValue["percent_soft"]; !soft.IsNull() && ok {
		if _, ok := priorThresholdValue["percent_soft"]; ok {
			thresholdValue["percent_soft"] = priorThresholdValue["percent_soft"]
		}
	}
	return types.ObjectValue(updated.AttributeTypes(ctx), thresholdValue)
}
