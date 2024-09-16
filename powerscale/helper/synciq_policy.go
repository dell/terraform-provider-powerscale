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
	"errors"
	"fmt"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetAllSyncIQPolicies retrieve the cluster information.
func GetAllSyncIQPolicies(ctx context.Context, client *client.Client) (*powerscale.V14SyncPolicies, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.ListSyncv14SyncPolicies(context.Background()).Execute()
	if err != nil {
		return resp, err
	}
	for resp.Resume != nil {
		respAdd, _, errAdd := client.PscaleOpenAPIClient.SyncApi.ListSyncv14SyncPolicies(context.Background()).Resume(*resp.Resume).Execute()
		if errAdd != nil {
			return resp, errAdd
		}
		resp.Resume = respAdd.Resume
		resp.Policies = append(resp.Policies, respAdd.Policies...)
	}
	return resp, err
}

// GetSyncIQPolicyIDByName retrieve the cluster information.
func GetSyncIQPolicyIDByName(ctx context.Context, client *client.Client, name string) (string, error) {
	policies, err := GetAllSyncIQPolicies(ctx, client)
	if err != nil {
		errStr := "Could not get list of SyncIQ policies with error: "
		message := GetErrorString(err, errStr)
		return "", errors.New(message)
	}
	for _, policy := range policies.Policies {
		if policy.Name == name {
			return policy.Id, nil
		}
	}
	return "", fmt.Errorf("policy by name %s not found", name)
}

// GetSyncIQPolicyByID retrieve the cluster information.
func GetSyncIQPolicyByID(ctx context.Context, client *client.Client, id string) (*powerscale.V14SyncPoliciesExtended, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.GetSyncv14SyncPolicy(context.Background(), id).Execute()
	return resp, err
}

// CreateSyncIQPolicy creates the sync iq policy.
func CreateSyncIQPolicy(ctx context.Context, client *client.Client, policy powerscale.V14SyncPolicy) (string, error) {
	resp, _, err := client.PscaleOpenAPIClient.SyncApi.CreateSyncv14SyncPolicy(ctx).V14SyncPolicy(policy).Execute()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// DeleteSyncIQPolicy deletes the sync iq policy.
func DeleteSyncIQPolicy(ctx context.Context, client *client.Client, id string) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.DeleteSyncv14SyncPolicy(ctx, id).Execute()
	return err
}

// UpdateSyncIQPolicy updates the sync iq policy.
func UpdateSyncIQPolicy(ctx context.Context, client *client.Client, id string, policy powerscale.V14SyncPolicyExtendedExtended) error {
	_, err := client.PscaleOpenAPIClient.SyncApi.UpdateSyncv14SyncPolicy(ctx, id).V14SyncPolicy(policy).Execute()
	return err
}

// SyncIQPolicyDataSourceResponse is the union of all response types for syncIQ policy datasource.
type SyncIQPolicyDataSourceResponse interface {
	powerscale.V14SyncPolicyExtended | powerscale.V14SyncPolicyExtendedExtendedExtended
}

// NewSyncIQPolicyDataSource creates a new SyncIQPolicyDataSource from datasource responses.
func NewSyncIQPolicyDataSource[V SyncIQPolicyDataSourceResponse](ctx context.Context, policies []V) (*models.SyncIQPolicyDataSource, error) {
	var err error
	ret := models.SyncIQPolicyDataSource{
		ID:       types.StringValue("dummy"),
		Policies: make([]models.V14SyncPolicyExtendedModel, len(policies)),
	}
	for i := range policies {
		var item models.V14SyncPolicyExtendedModel
		ierr := CopyFields(ctx, &policies[i], &item)
		err = errors.Join(err, ierr)
		ret.Policies[i] = item
	}
	if len(ret.Policies) == 1 {
		ret.ID = ret.Policies[0].ID
	}
	return &ret, err
}

// FillNilPointerWithDefault fills nill pointers with default values.
func FillNilPointerWithDefault[T any](loc *T, def T) *T {
	if loc != nil {
		return loc
	}
	return &def
}

// FillNilPointerWithEmpty fills nill pointers with empty values.
func FillNilPointerWithEmpty[T any](loc *T) *T {
	if loc != nil {
		return loc
	}
	return new(T)
}

// FillNilSlice initializes nill slices with empty list.
func FillNilSlice[T any](in []T) []T {
	if in != nil {
		return in
	}
	return make([]T, 0)
}

// normalizeSyncIQPolicy normalizes the sync iq policy.
// The server (according to its openAPI spec) can return null values for certain fields in scenarios where null means empty or some other default value.
// The resource should interpret these values correctly to ensure consistent behavior.
func normalizeSyncIQPolicy(source *powerscale.V14SyncPolicyExtendedExtendedExtended) {
	// If set to true, SyncIQ will perform failback configuration tasks during the next job run, rather than waiting to perform those tasks during the failback process. Performing these tasks ahead of time will increase the speed of failback operations.
	source.AcceleratedFailback = FillNilPointerWithEmpty(source.AcceleratedFailback)
	// If set to true, SyncIQ will allow a policy with copy action failback which is not supported by default.
	source.AllowCopyFb = FillNilPointerWithEmpty(source.AllowCopyFb)
	// The desired bandwidth reservation for this policy in kb/s. This feature will not activate unless a SyncIQ bandwidth rule is in effect.
	source.BandwidthReservation = FillNilPointerWithEmpty(source.BandwidthReservation)
	// If true, retain previous source snapshot and incremental repstate, both of which are required for changelist creation.
	source.Changelist = FillNilPointerWithEmpty(source.Changelist)
	// If true, the sync target performs cyclic redundancy checks (CRC) on the data as it is received.
	source.CheckIntegrity = FillNilPointerWithEmpty(source.CheckIntegrity)
	// If set to deny, replicates all CloudPools smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, the job will fail. If set to force, replicates all smartlinks to the target cluster as regular files. If set to allow, SyncIQ will attempt to replicate smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, SyncIQ will replicate the smartlinks as regular files.
	source.CloudDeepCopy = FillNilPointerWithEmpty(source.CloudDeepCopy) // empty value means nothing
	// If true, the most recent run of this policy encountered an error and this policy will not start any more scheduled jobs until this field is manually set back to 'false'.
	source.Conflicted = FillNilPointerWithEmpty(source.Conflicted)
	// If true, SyncIQ databases have been mirrored.
	source.DatabaseMirrored = FillNilPointerWithEmpty(source.DatabaseMirrored)
	// If true, forcibly remove quotas on the target after they have been removed on the source.
	source.DeleteQuotas = FillNilPointerWithEmpty(source.DeleteQuotas)
	// User-assigned description of this sync policy.
	source.Description = FillNilPointerWithEmpty(source.Description)
	// If true, the 7.2+ file splitting capability will be disabled.
	source.DisableFileSplit = FillNilPointerWithEmpty(source.DisableFileSplit)
	// Enable/disable sync failover/failback.
	source.DisableFofb = FillNilPointerWithEmpty(source.DisableFofb)
	// If set to true, SyncIQ will not create temporary quota directories to aid in replication to target paths which contain quotas.
	source.DisableQuotaTmpDir = FillNilPointerWithEmpty(source.DisableQuotaTmpDir)
	// Enable/disable the 6.5+ STF based data transfer and uses only treewalk.
	source.DisableStf = FillNilPointerWithEmpty(source.DisableStf)
	// If true, syncs will use temporary working directory subdirectories to reduce lock contention.
	source.EnableHashTmpdir = FillNilPointerWithEmpty(source.EnableHashTmpdir)
	// If true, syncs will be encrypted.
	source.Encrypted = FillNilPointerWithEmpty(source.Encrypted)
	// The cipher list being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.
	source.EncryptionCipherList = FillNilPointerWithEmpty(source.EncryptionCipherList)
	// Continue sending files even with the corrupted filesystem.
	source.ExpectedDataloss = FillNilPointerWithEmpty(source.ExpectedDataloss)

	if source.FileMatchingPattern != nil && len(source.FileMatchingPattern.OrCriteria) == 0 {
		source.FileMatchingPattern = nil
	}

	// Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.  If you enable this option, the net.inet.ip.choose_ifa_by_ipsrc sysctl should be set.
	source.ForceInterface = FillNilPointerWithEmpty(source.ForceInterface)
	// This field is false if the policy is in its initial sync state and true otherwise.  Setting this field to false will reset the policy's sync state.
	source.HasSyncState = FillNilPointerWithEmpty(source.HasSyncState)
	// If set to true, SyncIQ will not check the recursive quota in target paths to aid in replication to target paths which contain no quota but target cluster has lots of quotas.
	source.IgnoreRecursiveQuota = FillNilPointerWithEmpty(source.IgnoreRecursiveQuota)
	// If --schedule is set to When-Source-Modified, the duration to wait after a modification is made before starting a job (default is 0 seconds).
	source.JobDelay = FillNilPointerWithEmpty(source.JobDelay)
	// A list of service replication policies that this data replication policy will be associated with.
	source.LinkedServicePolicies = FillNilSlice(source.LinkedServicePolicies)
	// Severity an event must reach before it is logged.
	source.LogLevel = FillNilPointerWithEmpty(source.LogLevel) // empty value doesnt mean anything
	// If true, the system will log any files or directories that are deleted due to a sync.
	source.LogRemovedFiles = FillNilPointerWithEmpty(source.LogRemovedFiles)
	// The address of the OCSP responder to which to connect.
	source.OcspAddress = FillNilPointerWithEmpty(source.OcspAddress)
	// The ID of the certificate authority that issued the certificate whose revocation status is being checked.
	source.OcspIssuerCertificateId = FillNilPointerWithEmpty(source.OcspIssuerCertificateId)
	// Indicates if a password is set for accessing the target cluster. Password value is not shown with GET.
	source.PasswordSet = FillNilPointerWithEmpty(source.PasswordSet)
	// Determines the priority level of a policy. Policies with higher priority will have precedence to run over lower priority policies. Valid range is [0, 1]. Default is 0.
	source.Priority = FillNilPointerWithEmpty(source.Priority)

	// these two fields should actually never be null
	// Length of time (in seconds) a policy report will be stored.
	source.ReportMaxAge = FillNilPointerWithEmpty(source.ReportMaxAge)
	// Maximum number of policy reports that will be stored on the system.
	source.ReportMaxCount = FillNilPointerWithEmpty(source.ReportMaxCount)

	// If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.
	source.RestrictTargetNetwork = FillNilPointerWithEmpty(source.RestrictTargetNetwork)
	// If --schedule is set to a time/date, an alert is created if the specified RPO for this policy is exceeded. The default value is 0, which will not generate RPO alerts.
	source.RpoAlert = FillNilPointerWithEmpty(source.RpoAlert)
	// The schedule on which new jobs will be run for this policy.
	// Schedule string `json:"schedule"`
	// If true, this is a service replication policy.
	source.ServicePolicy = FillNilPointerWithEmpty(source.ServicePolicy)
	// Skip DNS lookup of target IPs.
	source.SkipLookup = FillNilPointerWithEmpty(source.SkipLookup)
	// If true and --schedule is set to a time/date, the policy will not run if no changes have been made to the contents of the source directory since the last job successfully completed.
	source.SkipWhenSourceUnmodified = FillNilPointerWithEmpty(source.SkipWhenSourceUnmodified)
	// If true, snapshot-triggered syncs will include snapshots taken before policy creation time (requires --schedule when-snapshot-taken).
	source.SnapshotSyncExisting = FillNilPointerWithEmpty(source.SnapshotSyncExisting)
	// The naming pattern that a snapshot must match to trigger a sync when the schedule is when-snapshot-taken (default is \"*\").
	source.SnapshotSyncPattern = FillNilPointerWithDefault(source.SnapshotSyncPattern, "*")
	// The ID of the source cluster certificate being used for encryption.
	// source.SourceCertificateId = FillNilPointerWithEmpty(source.SourceCertificateId)
	// If true, the source root path has been domain marked with a SyncIQ domain.
	source.SourceDomainMarked = FillNilPointerWithEmpty(source.SourceDomainMarked)
	// Directories that will be excluded from the sync.  Modifying this field will result in a full synchronization of all data.
	source.SourceExcludeDirectories = FillNilSlice(source.SourceExcludeDirectories)
	// Directories that will be included in the sync.  Modifying this field will result in a full synchronization of all data.
	source.SourceIncludeDirectories = FillNilSlice(source.SourceIncludeDirectories)

	source.SourceNetwork = FillNilPointerWithEmpty(source.SourceNetwork)

	// If true, archival snapshots of the source data will be taken on the source cluster before a sync.
	source.SourceSnapshotArchive = FillNilPointerWithEmpty(source.SourceSnapshotArchive)
	// The length of time in seconds to keep snapshots on the source cluster.
	// SourceSnapshotExpiration *int32 `json:"source_snapshot_expiration,omitempty"`
	// The name pattern for snapshots taken on the source cluster before a sync.
	source.SourceSnapshotPattern = FillNilPointerWithEmpty(source.SourceSnapshotPattern)
	// If set to true, the expire duration for target archival snapshot is the remaining expire duration of source snapshot, requires --sync-existing-snapshot=true
	source.SyncExistingSnapshotExpiration = FillNilPointerWithEmpty(source.SyncExistingSnapshotExpiration)
	// The naming pattern for snapshot on the destination cluster when --sync-existing-snapshot is true
	source.SyncExistingTargetSnapshotPattern = FillNilPointerWithEmpty(source.SyncExistingTargetSnapshotPattern)
	// The ID of the target cluster certificate being used for encryption.
	source.TargetCertificateId = FillNilPointerWithEmpty(source.TargetCertificateId)
	// If true, the target creates diffs against the original sync.
	source.TargetCompareInitialSync = FillNilPointerWithEmpty(source.TargetCompareInitialSync)
	// If true, target cluster will detect if files have been changed on the target by legacy tree walk syncs.
	source.TargetDetectModifications = FillNilPointerWithEmpty(source.TargetDetectModifications)
	// The alias of the snapshot taken on the target cluster after the sync completes. A value of @DEFAULT will reset this field to the default creation value.
	source.TargetSnapshotAlias = FillNilPointerWithEmpty(source.TargetSnapshotAlias)
	// If true, archival snapshots of the target data will be taken on the target cluster after successful sync completions.
	source.TargetSnapshotArchive = FillNilPointerWithEmpty(source.TargetSnapshotArchive)
	// The length of time in seconds to keep snapshots on the target cluster.
	source.TargetSnapshotExpiration = FillNilPointerWithEmpty(source.TargetSnapshotExpiration)
	// The name pattern for snapshots taken on the target cluster after the sync completes.  A value of @DEFAULT will reset this field to the default creation value.
	source.TargetSnapshotPattern = FillNilPointerWithEmpty(source.TargetSnapshotPattern)
	// The number of worker threads on a node performing a sync.
	source.WorkersPerNode = FillNilPointerWithEmpty(source.WorkersPerNode)
}

// NewSynciqpolicyResourceModel creates a new SynciqpolicyResourceModel from resource read response.
func NewSynciqpolicyResourceModel(ctx context.Context, respR *powerscale.V14SyncPoliciesExtended) (models.SynciqpolicyResourceModel, diag.Diagnostics) {
	var state models.SynciqpolicyResourceModel
	var dgs diag.Diagnostics
	source := respR.Policies[0]
	normalizeSyncIQPolicy(&source)

	err := CopyFieldsToNonNestedModel(ctx, source, &state)
	if err != nil {
		dgs.AddError(
			"Error copying fields of SyncIQ Policy resource",
			err.Error(),
		)
		return state, dgs
	}

	return state, nil
}
