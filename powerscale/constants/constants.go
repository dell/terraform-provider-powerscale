/*
Copyright (c) 2022-2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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

const (

	// ReadSnapshotErrorMessage specifies error details occurred while reading Snapshots.
	ReadSnapshotErrorMessage = "Could not read snapshots "

	// CreateSnapshotErrorMessage specifies error details occurred while create a Snapshot.
	CreateSnapshotErrorMessage = "Could not create snapshot "

	// UpdateSnapshotErrorMessage specifies error details occurred while delete a Snapshot.
	UpdateSnapshotErrorMessage = "Could not update snapshot "

	// DeleteSnapshotErrorMessage specifies error details occurred while delete a Snapshot.
	DeleteSnapshotErrorMessage = "Could not delete snapshot "

	// ReadAccessZoneErrorMsg specifies error details occurred while reading Access Zones.
	ReadAccessZoneErrorMsg = "Could not read access zones "

	// CreateAccessZoneErrorMsg specifies error details occurred while creating an Access Zones.
	CreateAccessZoneErrorMsg = "Could not create access zones "

	// UpdateAccessZoneErrorMsg specifies error details occurred while updating an Access Zones.
	UpdateAccessZoneErrorMsg = "Could not update access zones "

	// DeleteAccessZoneErrorMsg specifies error details occurred while deleting an Access Zones.
	DeleteAccessZoneErrorMsg = "Could not delete access zones "

	// ReadFileSystemErrorMsg specifies error details occurred while reading File Systems.
	ReadFileSystemErrorMsg = "Could not read file systems "

	// CreateFileSystemErrorMsg specifies error details occurred while creating an File Systems.
	CreateFileSystemErrorMsg = "Could not create file systems "

	// UpdateFileSystemErrorMsg specifies error details occurred while updating an File Systems.
	UpdateFileSystemErrorMsg = "Could not update file systems "

	// DeleteFileSystemErrorMsg specifies error details occurred while deleting an File Systems.
	DeleteFileSystemErrorMsg = "Could not delete file systems "

	// ReadAdsProviderErrorMsg specifies error details occurred while reading Ads Providers.
	ReadAdsProviderErrorMsg = "Could not read ads providers "

	// CreateAdsProviderErrorMsg specifies error details occurred while creating an Ads Provider.
	CreateAdsProviderErrorMsg = "Could not create ads providers "

	// UpdateAdsProviderErrorMsg specifies error details occurred while updating an Ads Provider.
	UpdateAdsProviderErrorMsg = "Could not update ads providers "

	// DeleteAdsProviderErrorMsg specifies error details occurred while deleting an Ads Provider.
	DeleteAdsProviderErrorMsg = "Could not delete ads providers "

	// ReadNetworkPoolErrorMsg specifies error details occurred while reading Network Pools.
	ReadNetworkPoolErrorMsg = "Could not read network pools "

	// CreateNetworkPoolErrorMsg specifies error details occurred while creating a Network Pool.
	CreateNetworkPoolErrorMsg = "Could not create network pools "

	// UpdateNetworkPoolErrorMsg specifies error details occurred while updating a Network Pool.
	UpdateNetworkPoolErrorMsg = "Could not update network pools "

	// DeleteNetworkPoolErrorMsg specifies error details occurred while deleting a Network Pool.
	DeleteNetworkPoolErrorMsg = "Could not delete network pools "

	// ReadUserErrorMsg specifies error details occurred while reading Users.
	ReadUserErrorMsg = "Could not read users "

	// CreateUserErrorMsg specifies error details occurred while creating an User.
	CreateUserErrorMsg = "Could not create user "

	// UpdateUserErrorMsg specifies error details occurred while updating an User.
	UpdateUserErrorMsg = "Could not update user "

	// DeleteUserErrorMsg specifies error details occurred while deleting an User.
	DeleteUserErrorMsg = "Could not delete user "

	// ReadRoleErrorMsg specifies error details occurred while reading Roles.
	ReadRoleErrorMsg = "Could not read roles "

	// AddRoleMemberErrorMsg specifies error details occurred while adding member to role.
	AddRoleMemberErrorMsg = "Could not add member to role "

	// DeleteRoleMemberErrorMsg specifies error details occurred while deleting member from role.
	DeleteRoleMemberErrorMsg = "Could not delete member from role "

	// ReadUserGroupErrorMsg specifies error details occurred while reading User Groups.
	ReadUserGroupErrorMsg = "Could not read user groups "

	// CreateUserGroupErrorMsg specifies error details occurred while creating an User Group.
	CreateUserGroupErrorMsg = "Could not create user group "

	// UpdateUserGroupErrorMsg specifies error details occurred while updating an User Group.
	UpdateUserGroupErrorMsg = "Could not update user group "

	// DeleteUserGroupErrorMsg specifies error details occurred while deleting an User Group.
	DeleteUserGroupErrorMsg = "Could not delete user group "

	// ReadUserGroupMemberErrorMsg specifies error details occurred while reading UserGroup Members.
	ReadUserGroupMemberErrorMsg = "Could not read User Group Members "

	// AddUserGroupMemberErrorMsg specifies error details occurred while adding member to user group.
	AddUserGroupMemberErrorMsg = "Could not add member to user group "

	// DeleteUserGroupMemberErrorMsg specifies error details occurred while deleting member from user group.
	DeleteUserGroupMemberErrorMsg = "Could not delete member from user group "

	// GetNfsExportErrorMsg specifies error details occurred while getting nfs export.
	GetNfsExportErrorMsg = "Could not get nfs export "

	// UpdateNfsExportErrorMsg specifies error details occurred while updating nfs export.
	UpdateNfsExportErrorMsg = "Could not update nfs export "

	// DeleteNfsExportErrorMsg specifies error details occurred while deleting nfs export.
	DeleteNfsExportErrorMsg = "Could not delete nfs export "

	// CreateNfsExportErrorMsg specifies error details occurred while creating nfs export.
	CreateNfsExportErrorMsg = "Could not create nfs export "

	// ListNfsExportErrorMsg specifies error details occurred while listing nfs exports.
	ListNfsExportErrorMsg = "Could not list nfs exports "

	// GetSmbShareErrorMsg specifies error details occurred while getting smb share.
	GetSmbShareErrorMsg = "Could not get smb share "

	// UpdateSmbShareErrorMsg specifies error details occurred while updating smb share.
	UpdateSmbShareErrorMsg = "Could not update smb share "

	// DeleteSmbShareErrorMsg specifies error details occurred while deleting smb share.
	DeleteSmbShareErrorMsg = "Could not delete smb share "

	// CreateSmbShareErrorMsg specifies error details occurred while creating smb share.
	CreateSmbShareErrorMsg = "Could not create smb share "

	// ListSmbShareErrorMsg specifies error details occurred while listing smb shares.
	ListSmbShareErrorMsg = "Could not list smb shares "

	// ReadWellKnownErrorMsg specifies error details occurred while reading Well-Knowns.
	ReadWellKnownErrorMsg = "Could not read well-knowns "

	// ReadClusterErrorMsg specifies error details occurred while reading cluster.
	ReadClusterErrorMsg = "Could not read cluster "

	// ListSnapshotSchedulesMsg specifies error details occurred while listing snapshot schedules.
	ListSnapshotSchedulesMsg = "Could not list snapshot schedules "

	// ReadSnapshotScheduleErrorMessage specifies error details occurred while reading Snapshot schedule.
	ReadSnapshotScheduleErrorMessage = "Could not read snapshot schedule"

	// CreateSnapshotScheduleErrorMessage specifies error details occurred while create a Snapshot schedule.
	CreateSnapshotScheduleErrorMessage = "Could not create snapshot schedule"

	// UpdateSnapshotScheduleErrorMessage specifies error details occurred while delete a Snapshot schedule.
	UpdateSnapshotScheduleErrorMessage = "Could not update snapshot schedule"

	// DeleteSnapshotScheduleErrorMessage specifies error details occurred while delete a Snapshot schedule.
	DeleteSnapshotScheduleErrorMessage = "Could not delete snapshot schedule"

	// ListQuotaErrorMsg specifies error details occurred while listing quotas.
	ListQuotaErrorMsg = "Could not list quotas "

	// ReadQuotaErrorMsg specifies error details occurred while reading quotas.
	ReadQuotaErrorMsg = "Could not read quotas "

	// CreateQuotaErrorMsg specifies error details occurred while creating quotas.
	CreateQuotaErrorMsg = "Could not create quotas "

	// UpdateQuotaErrorMsg specifies error details occurred while updating quotas.
	UpdateQuotaErrorMsg = "Could not update quotas "

	// DeleteQuotaErrorMsg specifies error details occurred while deleting quotas.
	DeleteQuotaErrorMsg = "Could not delete quotas "

	// ListSubnetErrorMsg specifies error details occurred while listing subnets.
	ListSubnetErrorMsg = "Could not list subnets "

	// ReadGroupnetErrorMsg specifies error details occurred while reading groupnet.
	ReadGroupnetErrorMsg = "Could not read groupnet "

	// CreateGroupnetErrorMsg specifies error details occurred while creating a groupnet.
	CreateGroupnetErrorMsg = "Could not create groupnet "

	// UpdateGroupnetErrorMsg specifies error details occurred while updating a groupnet.
	UpdateGroupnetErrorMsg = "Could not update groupnet "

	// DeleteGroupnetErrorMsg specifies error details occurred while deleting a groupnet.
	DeleteGroupnetErrorMsg = "Could not delete groupnet "

	// ReadSmartPoolSettingsErrorMsg specifies error details occurred while reading smart pool settings.
	ReadSmartPoolSettingsErrorMsg = "Could not read SmartPool settings "

	// ReadNetworkSettingErrorMsg specifies error details occurred while reading network setting.
	ReadNetworkSettingErrorMsg = "Could not read network setting "

	// UpdateNetworkSettingErrorMsg specifies error details occurred while updating network setting.
	UpdateNetworkSettingErrorMsg = "Could not update network setting "

	// CreateSubnetErrorMsg specifies error details occurred while creating subnet.
	CreateSubnetErrorMsg = "Could not create subnet "

	// GetSubnetErrorMsg specifies error details occurred while getting subnet.
	GetSubnetErrorMsg = "Could not get subnet "

	// UpdateSubnetErrorMsg specifies error details occurred while updating subnet.
	UpdateSubnetErrorMsg = "Could not update subnet "

	// DeleteSubnetErrorMsg specifies error details occurred while deleting subnet.
	DeleteSubnetErrorMsg = "Could not delete subnet "
)
