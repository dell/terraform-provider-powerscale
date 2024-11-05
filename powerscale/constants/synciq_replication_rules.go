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

package constants

import "fmt"

// SyncIQRuleType is the type of a SyncIQ rule.
type SyncIQRuleType string

const (
	// SyncIQRuleTypeBW is a bandwidth-based SyncIQ rule.
	SyncIQRuleTypeBW SyncIQRuleType = "bandwidth"

	// SyncIQRuleTypeFC is a file count-based SyncIQ rule.
	SyncIQRuleTypeFC SyncIQRuleType = "file_count"

	// SyncIQRuleTypeCPU is a CPU-based SyncIQ rule.
	SyncIQRuleTypeCPU SyncIQRuleType = "cpu"

	// SyncIQRuleTypeWK is a worker-based SyncIQ rule.
	SyncIQRuleTypeWK SyncIQRuleType = "worker"
)

// GetSynciqRuleID gets the ID for a synciq rule for a particular index in a
// particular stack.
func GetSynciqRuleID(i int, ruleType SyncIQRuleType) string {
	idType := map[SyncIQRuleType]string{
		SyncIQRuleTypeBW:  "bw",
		SyncIQRuleTypeFC:  "fc",
		SyncIQRuleTypeCPU: "cpu",
		SyncIQRuleTypeWK:  "wk",
	}
	return fmt.Sprintf("%s-%d", idType[ruleType], i)
}

// SyncIQRule Days of the week
const (
	SyncIQRuleDayMonday    = "monday"
	SyncIQRuleDayTuesday   = "tuesday"
	SyncIQRuleDayWednesday = "wednesday"
	SyncIQRuleDayThursday  = "thursday"
	SyncIQRuleDayFriday    = "friday"
	SyncIQRuleDaySaturday  = "saturday"
	SyncIQRuleDaySunday    = "sunday"
)
