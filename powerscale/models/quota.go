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

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// QuotaResource specifies configuration values for NFS exports.
type QuotaResource struct {
	// Optional attributes and can not update
	// Optional named zone to use for user and group resolution.
	Zone types.String `tfsdk:"zone"`
	// Specifies the persona of the file group.
	Persona types.Object `tfsdk:"persona"`

	// Optional and can update
	// If true, skip child quota's threshold comparison with parent quota path.
	IgnoreLimitChecks types.Bool `tfsdk:"ignore_limit_checks"`
	// Force creation of quotas on the root of /ifs or percent based quotas.
	Force types.Bool `tfsdk:"force"`

	// Required attributes cannot be updated
	// If true, quota governs snapshot data as well as head data.
	IncludeSnapshots types.Bool `tfsdk:"include_snapshots"`
	// The /ifs path governed.
	Path types.String `tfsdk:"path"`
	// The type of quota.
	Type types.String `tfsdk:"type"`

	// Attributes can be updated only after created
	// For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked.
	Linked types.Bool `tfsdk:"linked"`

	// Attributes can be read and updated
	// If true, SMB shares using the quota directory see the quota thresholds as share size.
	Container types.Bool `tfsdk:"container"`
	// True if the quota provides enforcement, otherwise an accounting quota.
	Enforced types.Bool `tfsdk:"enforced"`
	// The thresholds of quota
	Thresholds types.Object `tfsdk:"thresholds"`
	// Thresholds apply on quota accounting metric.
	ThresholdsOn types.String `tfsdk:"thresholds_on"`

	// ReadOnly Attributes
	// Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.
	EfficiencyRatio types.Number `tfsdk:"efficiency_ratio"`
	// The system ID given to the quota.
	ID types.String `tfsdk:"id"`
	// Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.
	Notifications types.String `tfsdk:"notifications"`
	// True if the default resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	Ready types.Bool `tfsdk:"ready"`
	// Represents the ratio of logical space provided to physical data space used. This accounts for compression and data deduplication effects.
	ReductionRatio types.Number `tfsdk:"reduction_ratio"`
	// The usage of quota
	Usage types.Object `tfsdk:"usage"`
}

// QuotaDatasource holds quotas datasource schema attribute details.
type QuotaDatasource struct {
	ID          types.String            `tfsdk:"id"`
	Quotas      []QuotaDatasourceEntity `tfsdk:"quotas"`
	QuotaFilter *QuotaDatasourceFilter  `tfsdk:"filter"`
}

// QuotaDatasourceFilter holds filter conditions.
type QuotaDatasourceFilter struct {
	Enforced            types.Bool   `tfsdk:"enforced"`
	Exceeded            types.Bool   `tfsdk:"exceeded"`
	IncludeSnapshots    types.Bool   `tfsdk:"include_snapshots"`
	Path                types.String `tfsdk:"path"`
	Persona             types.String `tfsdk:"persona"`
	RecursePathChildren types.Bool   `tfsdk:"recurse_path_children"`
	RecursePathParents  types.Bool   `tfsdk:"recurse_path_parents"`
	ReportID            types.String `tfsdk:"report_id"`
	Type                types.String `tfsdk:"type"`
	Zone                types.String `tfsdk:"zone"`
}

// QuotaDatasourceEntity struct for Quota data source model.
type QuotaDatasourceEntity struct {
	// If true, SMB shares using the quota directory see the quota thresholds as share size.
	Container types.Bool `tfsdk:"container"`
	// Represents the ratio of logical space provided to physical space used. This accounts for protection overhead, metadata, and compression ratios for the data.
	EfficiencyRatio types.Number `tfsdk:"efficiency_ratio"`
	// True if the quota provides enforcement, otherwise an accounting quota.
	Enforced types.Bool `tfsdk:"enforced"`
	// The system ID given to the quota.
	ID types.String `tfsdk:"id"`
	// If true, quota governs snapshot data as well as head data.
	IncludeSnapshots types.Bool `tfsdk:"include_snapshots"`
	// For user, group and directory quotas, true if the quota is linked and controlled by a parent default-* quota. Linked quotas cannot be modified until they are unlinked.
	Linked types.Bool `tfsdk:"linked"`
	// Summary of notifications: 'custom' indicates one or more notification rules available from the notifications sub-resource; 'default' indicates system default rules are used; 'disabled' indicates that no notifications will be used for this quota.; 'badmap' indicates that notification rule has problem in rule map.
	Notifications types.String `tfsdk:"notifications"`
	// The /ifs path governed.
	Path types.String `tfsdk:"path"`
	// Specifies the persona of the file group.
	Persona V1AuthAccessAccessItemFileGroup `tfsdk:"persona"`
	// True if the default resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	Ready types.Bool `tfsdk:"ready"`
	// Represents the ratio of logical space provided to physical data space used. This accounts for compression and data deduplication effects.
	ReductionRatio types.Number `tfsdk:"reduction_ratio"`
	// The thresholds of quota
	Thresholds QuotaThreshold `tfsdk:"thresholds"`
	// Thresholds apply on quota accounting metric.
	ThresholdsOn types.String `tfsdk:"thresholds_on"`
	// The type of quota.
	Type types.String `tfsdk:"type"`
	// The usage of quota
	Usage QuotaUsage `tfsdk:"usage"`
}

// QuotaThreshold struct for Quota threshold.
type QuotaThreshold struct {
	// Usage bytes at which notifications will be sent but writes will not be denied.
	Advisory types.Int64 `tfsdk:"advisory"`
	// True if the advisory threshold has been hit.
	AdvisoryExceeded types.Bool `tfsdk:"advisory_exceeded"`
	// Time at which advisory threshold was hit.
	AdvisoryLastExceeded types.Int64 `tfsdk:"advisory_last_exceeded"`
	// Usage bytes at which further writes will be denied.
	Hard types.Int64 `tfsdk:"hard"`
	// True if the hard threshold has been hit.
	HardExceeded types.Bool `tfsdk:"hard_exceeded"`
	// Time at which hard threshold was hit.
	HardLastExceeded types.Int64 `tfsdk:"hard_last_exceeded"`
	// Advisory threshold as percent of hard threshold. Usage bytes at which notifications will be sent but writes will not be denied.
	PercentAdvisory types.Float64 `tfsdk:"percent_advisory"`
	// Soft threshold as percent of hard threshold. Usage bytes at which notifications will be sent and soft grace time will be started.
	PercentSoft types.Float64 `tfsdk:"percent_soft"`
	// Usage bytes at which notifications will be sent and soft grace time will be started.
	Soft types.Int64 `tfsdk:"soft"`
	// True if the soft threshold has been hit.
	SoftExceeded types.Bool `tfsdk:"soft_exceeded"`
	// Time in seconds after which the soft threshold has been hit before writes will be denied.
	SoftGrace types.Int64 `tfsdk:"soft_grace"`
	// Time at which soft threshold was hit
	SoftLastExceeded types.Int64 `tfsdk:"soft_last_exceeded"`
}

// QuotaUsage struct for QuotaUsage.
type QuotaUsage struct {
	// Bytes used by governed data apparent to application.
	Applogical types.Int64 `tfsdk:"applogical"`
	// True if applogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	ApplogicalReady types.Bool `tfsdk:"applogical_ready"`
	// Bytes used by governed data apparent to filesystem.
	Fslogical types.Int64 `tfsdk:"fslogical"`
	// True if fslogical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	FslogicalReady types.Bool `tfsdk:"fslogical_ready"`
	// Physical data usage adjusted to account for shadow store efficiency
	Fsphysical types.Int64 `tfsdk:"fsphysical"`
	// True if fsphysical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	FsphysicalReady types.Bool `tfsdk:"fsphysical_ready"`
	// Number of inodes (filesystem entities) used by governed data.
	Inodes types.Int64 `tfsdk:"inodes"`
	// True if inodes resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	InodesReady types.Bool `tfsdk:"inodes_ready"`
	// Bytes used for governed data and filesystem overhead.
	Physical types.Int64 `tfsdk:"physical"`
	// Number of physical blocks for file data
	PhysicalData types.Int64 `tfsdk:"physical_data"`
	// True if physical_data resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	PhysicalDataReady types.Bool `tfsdk:"physical_data_ready"`
	// Number of physical blocks for file protection
	PhysicalProtection types.Int64 `tfsdk:"physical_protection"`
	// True if physical_protection resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	PhysicalProtectionReady types.Bool `tfsdk:"physical_protection_ready"`
	// True if physical resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	PhysicalReady types.Bool `tfsdk:"physical_ready"`
	// Number of shadow references (cloned, deduplicated or packed filesystem blocks) used by governed data.
	ShadowRefs types.Int64 `tfsdk:"shadow_refs"`
	// True if shadow_refs resource accounting is accurate on the quota. If false, this quota is waiting on completion of a QuotaScan job.
	ShadowRefsReady types.Bool `tfsdk:"shadow_refs_ready"`
}
