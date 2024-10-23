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
	"encoding/json"
	"fmt"
	"terraform-provider-powerscale/powerscale/constants"
	"terraform-provider-powerscale/powerscale/helper"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SyncIQRuleResource{}
	_ resource.ResourceWithConfigure   = &SyncIQRuleResource{}
	_ resource.ResourceWithImportState = &SyncIQRuleResource{}
)

// NewSyncIQRuleResource creates a new data source.
func NewSyncIQRuleResource() resource.Resource {
	return &SyncIQRuleResource{
		commonResourceConfigurer{
			name: "synciq_rules",
		},
	}
}

// SyncIQRuleResource defines the resource implementation.
type SyncIQRuleResource struct {
	commonResourceConfigurer
}

// Schema describes the data source arguments.
func (d *SyncIQRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = helper.SyncIQRulesResourceSchema(ctx)
}

// Create allocates the resource.
func (d *SyncIQRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SyncIQRulesResource
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 	Read Initial state
	state, dgsR := d.get(ctx)
	resp.Diagnostics.Append(dgsR...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update
	state, dgsR = d.update(ctx, plan, *state)
	resp.Diagnostics.Append(dgsR...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read reads data for the resource.
func (d *SyncIQRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state, dgs := d.get(ctx)
	resp.Diagnostics.Append(dgs...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Get fetches the resource state.
func (d *SyncIQRuleResource) get(ctx context.Context) (*models.SyncIQRulesResource, diag.Diagnostics) {
	// Read Terraform configuration data into the model
	var dgs diag.Diagnostics
	config, err := helper.GetAllSyncIQRules(ctx, d.client)
	if err != nil {
		errStr := constants.ListSynciqRulesMsg + "with error: "
		message := helper.GetErrorString(err, errStr)
		dgs.AddError("Error reading syncIQ rules", message)
		return nil, dgs
	}
	state, diags := helper.NewSyncIQRulesResource(ctx, config)
	dgs.Append(diags...)
	return &state, dgs
}

// Update updates the resource.
func (d *SyncIQRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state models.SyncIQRulesResource
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update
	stateF, dgsR := d.update(ctx, plan, state)
	resp.Diagnostics.Append(dgsR...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateF)...)
}

// Delete implements resource.Resource.
func (d *SyncIQRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

// update updates all supported synciq rule stacks
func (d *SyncIQRuleResource) update(ctx context.Context, plan, state models.SyncIQRulesResource) (*models.SyncIQRulesResource, diag.Diagnostics) {
	// upmarshall
	planReqs, stateReqs := helper.GetRequestsFromSynciqRulesResource(ctx, plan), helper.GetRequestsFromSynciqRulesResource(ctx, state)

	// run update for bandwidth stack
	dgs := d.updateStack(ctx, planReqs.BandWidthRules, stateReqs.BandWidthRules, "bandwidth")
	if dgs.HasError() {
		return nil, dgs
	}

	// return final state
	return d.get(ctx)
}

// updateStack updates any one of the synciq rule stacks
// the stack is identified by the ruleType param
// only bandwidth supported now
func (d *SyncIQRuleResource) updateStack(ctx context.Context, planTfsdk, stateTfsdk []models.SyncIQRuleResource, ruleType string) diag.Diagnostics {
	var dgs diag.Diagnostics

	// convert plan and state to []powerscale.V3SyncRule
	plan, state := make([]powerscale.V3SyncRule, 0, len(planTfsdk)), make([]powerscale.V3SyncRule, 0, len(stateTfsdk))
	for _, j := range planTfsdk {
		plan = append(plan, helper.GetRequestFromSynciqRuleResource(ctx, j, ruleType))
	}
	for _, j := range stateTfsdk {
		state = append(state, helper.GetRequestFromSynciqRuleResource(ctx, j, ruleType))
	}

	if toBeDeleted := len(state) - len(plan); toBeDeleted > 0 {
		// If number of planned rules is less than existing rules, delete the excess rules from last applicable to first
		for i := toBeDeleted - 1; i >= 0; i-- {
			id := d.getID(i, ruleType)
			err := helper.DeleteSyncIQRule(ctx, d.client, id)
			if err != nil {
				errStr := constants.ListSynciqRulesMsg + "with error: "
				message := helper.GetErrorString(err, errStr)
				dgs.AddError("Error deleting syncIQ rule "+id, message)
				return dgs
			}
		}
		state = state[toBeDeleted:]
	} else if toBeAdded := len(plan) - len(state); toBeAdded > 0 {
		// If number of planned rules is more than existing rules, create the excess rules from last applicable to first
		for i := toBeAdded - 1; i >= 0; i-- {
			_, err := helper.CreateSyncIQRule(ctx, d.client, plan[i])
			if err != nil {
				errStr := constants.ListSynciqRulesMsg + "with error: "
				message := helper.GetErrorString(err, errStr)
				dgs.AddError("Error creating syncIQ rule for index "+string(i), message)
				return dgs
			}
		}
		// prepend to state the items from plan[0] to plan[toBeAdded-1]
		state = append(
			append(make([]powerscale.V3SyncRule, 0, len(plan)), plan[:toBeAdded]...),
			state...,
		)
	}

	// now the number of rules are equal
	// update remaining rules to make state consistent with plan
	for i, planItem := range plan {
		stateItem := state[i]
		if d.areRulesEqual(planItem, stateItem) {
			continue
		}
		id := d.getID(i, ruleType)
		err := helper.UpdateSyncIQRule(ctx, d.client, id, planItem)
		if err != nil {
			errStr := constants.ListSynciqRulesMsg + "with error: "
			message := helper.GetErrorString(err, errStr)
			dgs.AddError("Error updating syncIQ rule "+id, message)
			return dgs
		}
	}

	return nil
}

// gets the ID for a synciq rule for a particular index in a particular stack
func (d *SyncIQRuleResource) getID(i int, ruleType string) string {
	var idType string
	if ruleType == "bandwidth" {
		idType = "bw"
	}
	return fmt.Sprintf("%s-%d", idType, i)
}

// checks if two synciq rules are equal by comparing their JSON representations
func (d *SyncIQRuleResource) areRulesEqual(plan, existing powerscale.V3SyncRule) bool {
	planJson, _ := json.Marshal(plan)
	stateJson, _ := json.Marshal(existing)
	return string(planJson) == string(stateJson)
}
