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

// SyncIQGlobalSettingsModel struct for SyncIQ global settings
type SyncIQGlobalSettingsModel struct {
	PreferredRpoAlert                     types.Int64  `tfsdk:"preferred_rpo_alert"`
	PasswordSet                           types.Bool   `tfsdk:"password_set"`
	Service                               types.String `tfsdk:"service"`
	ClusterCertificateID                  types.String `tfsdk:"cluster_certificate_id"`
	MaxConcurrentJobs                     types.Int64  `tfsdk:"max_concurrent_jobs"`
	OcspIssuerCertificateID               types.String `tfsdk:"ocsp_issuer_certificate_id"`
	ReportMaxCount                        types.Int64  `tfsdk:"report_max_count"`
	ForceInterface                        types.Bool   `tfsdk:"force_interface"`
	OcspAddress                           types.String `tfsdk:"ocsp_address"`
	SourceNetwork                         types.Object `tfsdk:"source_network"`
	RpoAlerts                             types.Bool   `tfsdk:"rpo_alerts"`
	BandwidthReservationReservePercentage types.Int64  `tfsdk:"bandwidth_reservation_reserve_percentage"`
	EncryptionCipherList                  types.String `tfsdk:"encryption_cipher_list"`
	RenegotiationPeriod                   types.Int64  `tfsdk:"renegotiation_period"`
	BandwidthReservationReserveAbsolute   types.Int64  `tfsdk:"bandwidth_reservation_reserve_absolute"`
	ServiceHistoryMaxCount                types.Int64  `tfsdk:"service_history_max_count"`
	UseWorkersPerNode                     types.Bool   `tfsdk:"use_workers_per_node"`
	TwChkptInterval                       types.Int64  `tfsdk:"tw_chkpt_interval"`
	EncryptionRequired                    types.Bool   `tfsdk:"encryption_required"`
	ServiceHistoryMaxAge                  types.Int64  `tfsdk:"service_history_max_age"`
	ReportMaxAge                          types.Int64  `tfsdk:"report_max_age"`
	RestrictTargetNetwork                 types.Bool   `tfsdk:"restrict_target_network"`
	ReportEmail                           types.Set    `tfsdk:"report_email"`
	Password                              types.String `tfsdk:"password"`
}
