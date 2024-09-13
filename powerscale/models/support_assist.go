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

import (
	powerscale "dell/powerscale-go-client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SupportAssistModel represents the model for support assist resource
type SupportAssistModel struct {
	ID                    types.String `tfsdk:"id"`
	EnableDownload        types.Bool   `tfsdk:"enable_download"`
	Contact               types.Object `tfsdk:"contact"`
	Telemetry             types.Object `tfsdk:"telemetry"`
	AutomaticCaseCreation types.Bool   `tfsdk:"automatic_case_creation"`
	Connection            types.Object `tfsdk:"connections"`
	EnableRemoteSupport   types.Bool   `tfsdk:"enable_remote_support"`
	Accepted              types.Bool   `tfsdk:"accepted_terms"`
	SupportassistEnabled  types.Bool   `tfsdk:"supportassist_enabled"`
	AccessKey             types.String `tfsdk:"access_key"`
	Pin                   types.String `tfsdk:"pin"`
}

// V16SupportassistSettingsCustomised represents the customized settings for the support assist
type V16SupportassistSettingsCustomised struct {
	AutomaticCaseCreation *bool                                         `json:"automatic_case_creation,omitempty"`
	Connections           *V16SupportassistSettingsConnection           `json:"connections,omitempty"`
	ConnectionState       *string                                       `json:"connection_state,omitempty"`
	Contact               *powerscale.V16SupportassistSettingsContact   `json:"contact,omitempty"`
	EnableDownload        *bool                                         `json:"enable_download,omitempty"`
	EnableRemoteSupport   *bool                                         `json:"enable_remote_support,omitempty"`
	SupportassistEnabled  bool                                          `json:"supportassist_enabled"`
	Telemetry             *powerscale.V16SupportassistSettingsTelemetry `json:"telemetry,omitempty"`
}

// UpdateSupportassistSettings represents the settings for the support assist
type UpdateSupportassistSettings struct {
	AutomaticCaseCreation *bool                                                  `json:"automatic_case_creation,omitempty"`
	Connection            *powerscale.V16SupportassistSettingsConnectionExtended `json:"connections,omitempty"`
	ConnectionState       *string                                                `json:"connection_state,omitempty"`
	Contact               *powerscale.V16SupportassistSettingsContactExtended    `json:"contact,omitempty"`
	EnableDownload        *bool                                                  `json:"enable_download,omitempty"`
	EnableRemoteSupport   *bool                                                  `json:"enable_remote_support,omitempty"`
	Telemetry             *powerscale.V16SupportassistSettingsTelemetryExtended  `json:"telemetry,omitempty"`
}

// V16SupportassistSettingsConnection represents the settings for the support assist connection
type V16SupportassistSettingsConnection struct {
	GatewayEndpoints []powerscale.V16SupportassistSettingsConnectionGatewayEndpoint `json:"gateway_endpoints,omitempty"`
	Mode             *string                                                        `json:"mode,omitempty"`
	NetworkPools     []string                                                       `json:"network_pools,omitempty"`
}

// Connection represents the settings for the support assist connection
type Connection struct {
	Mode             types.String `tfsdk:"mode"`
	NetworkPools     types.List   `tfsdk:"network_pools"`
	GatewayEndpoints types.List   `tfsdk:"gateway_endpoints"`
}

// GatewayEndpoints represents the settings for the support assist gateway
type GatewayEndpoints struct {
	Port        types.Int64  `tfsdk:"port"`
	ValidateSsl types.Bool   `tfsdk:"validate_ssl"`
	Priority    types.Int64  `tfsdk:"priority"`
	UseProxy    types.Bool   `tfsdk:"use_proxy"`
	Host        types.String `tfsdk:"host"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}
