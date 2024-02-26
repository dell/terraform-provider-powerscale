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
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"
)

// ListNetworkRules list network rules.
func ListNetworkRules(ctx context.Context, client *client.Client, filter *models.NetworkRuleFilterType) ([]powerscale.V3PoolsPoolRulesRule, error) {
	networkRuleParams := client.PscaleOpenAPIClient.NetworkApi.GetNetworkv3NetworkRules(ctx)

	if filter != nil {
		if groupnet := filter.Groupnet.ValueString(); groupnet != "" {
			networkRuleParams = networkRuleParams.Groupnet(groupnet)
		}
		if subnet := filter.Subnet.ValueString(); subnet != "" {
			networkRuleParams = networkRuleParams.Subnet(subnet)
		}
		if pool := filter.Pool.ValueString(); pool != "" {
			networkRuleParams = networkRuleParams.Pool(pool)
		}
	}

	ruleList, _, err := networkRuleParams.Execute()
	if err != nil {
		return nil, err
	}
	rules := ruleList.GetRules()

	// filter rules by filter.Names
	if filter != nil && len(filter.Names) > 0 {
		var filteredRules []powerscale.V3PoolsPoolRulesRule
		for _, rule := range rules {
			for _, name := range filter.Names {
				if name.ValueString() == rule.GetName() {
					filteredRules = append(filteredRules, rule)
				}
			}
		}
		return filteredRules, nil
	}
	return rules, nil
}

// CreateNetworkRule create.
func CreateNetworkRule(ctx context.Context, client *client.Client, groupnet string, subnet string, pool string, ruleToCreate powerscale.V3PoolsPoolRule) (*powerscale.CreateResponse, error) {
	networkRuleID, _, err := client.PscaleOpenAPIClient.NetworkGroupnetsSubnetsApi.CreateNetworkGroupnetsSubnetsv3PoolsPoolRule(ctx, groupnet, subnet, pool).V3PoolsPoolRule(ruleToCreate).Execute()
	return networkRuleID, err
}

// GetNetworkRule get rule.
func GetNetworkRule(ctx context.Context, client *client.Client, ruleName string, groupnet string, subnet string, pool string) (*powerscale.V3PoolsPoolRulesRule, error) {
	rules, _, err := client.PscaleOpenAPIClient.NetworkApi.GetNetworkv3GroupnetsGroupnetSubnetsSubnetPoolsPoolRule(ctx, ruleName, groupnet, subnet, pool).Execute()
	if err != nil {
		return nil, err
	}
	ruleSlice := rules.GetRules()
	if len(ruleSlice) != 1 {
		return nil, fmt.Errorf("error get network rule, %d rules are found with Name: %s", len(ruleSlice), ruleName)
	}
	return &ruleSlice[0], err
}

// UpdateNetworkRule update network rule.
func UpdateNetworkRule(ctx context.Context, client *client.Client, ruleName string, groupnet string, subnet string, pool string, ruleToUpdate powerscale.V3GroupnetsGroupnetSubnetsSubnetPoolsPoolRule) error {
	_, err := client.PscaleOpenAPIClient.NetworkApi.UpdateNetworkv3GroupnetsGroupnetSubnetsSubnetPoolsPoolRule(ctx, ruleName, groupnet, subnet, pool).V3GroupnetsGroupnetSubnetsSubnetPoolsPoolRule(ruleToUpdate).Execute()
	return err
}

// DeleteNetworkRule delete network rule.
func DeleteNetworkRule(ctx context.Context, client *client.Client, ruleName string, groupnet string, subnet string, pool string) error {
	_, err := client.PscaleOpenAPIClient.NetworkApi.DeleteNetworkv3GroupnetsGroupnetSubnetsSubnetPoolsPoolRule(ctx, ruleName, groupnet, subnet, pool).Execute()
	return err
}
