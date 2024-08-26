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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ClusterTime struct for ClusterTime resource.
type ClusterTime struct {
	ID           types.String `tfsdk:"id"`
	Date         types.String `tfsdk:"date"`
	Time         types.String `tfsdk:"time"`
	TimeMillis   types.Int32  `tfsdk:"time_millis"`
	Abbreviation types.String `tfsdk:"abbreviation"`
	Path         types.String `tfsdk:"path"`
}

type ClusterTimezoneSetting struct {
	Settings *ClusterConfigTimezone `tfsdk:"settings"`
}
