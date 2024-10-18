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

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema implements resource.Resource.
func (s *synciqPolicyResource) Schema(ctx context.Context, res resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is used to manage the SyncIQ Replication Policy entity of PowerScale Array. " +
			"We can Create, Read, Update and Delete the SyncIQ Replication Policy using this resource. We can also import existing SyncIQ Replication Policy from PowerScale array.",
		Description: "This resource is used to manage the SyncIQ Replication Policy entity of PowerScale Array. " +
			"We can Create, Read, Update and Delete the SyncIQ Replication Policy using this resource. We can also import existing SyncIQ Replication Policy from PowerScale array.",
		Attributes: map[string]schema.Attribute{
			"accelerated_failback": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to true, SyncIQ will perform failback configuration tasks during the next job run, rather than waiting to perform those tasks during the failback process. Performing these tasks ahead of time will increase the speed of failback operations.",
				MarkdownDescription: "If set to true, SyncIQ will perform failback configuration tasks during the next job run, rather than waiting to perform those tasks during the failback process. Performing these tasks ahead of time will increase the speed of failback operations.",
			},
			"action": schema.StringAttribute{
				Required:            true,
				Description:         "If 'copy', source files will be copied to the target cluster.  If 'sync', the target directory will be made an image of the source directory:  Files and directories that have been deleted on the source, have been moved within the target directory, or no longer match the selection criteria will be deleted from the target directory.",
				MarkdownDescription: "If 'copy', source files will be copied to the target cluster.  If 'sync', the target directory will be made an image of the source directory:  Files and directories that have been deleted on the source, have been moved within the target directory, or no longer match the selection criteria will be deleted from the target directory.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"copy",
						"sync",
					),
				},
			},
			"allow_copy_fb": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to true, SyncIQ will allow a policy with copy action failback which is not supported by default.",
				MarkdownDescription: "If set to true, SyncIQ will allow a policy with copy action failback which is not supported by default.",
			},
			"bandwidth_reservation": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The desired bandwidth reservation for this policy in kb/s. This feature will not activate unless a SyncIQ bandwidth rule is in effect.",
				MarkdownDescription: "The desired bandwidth reservation for this policy in kb/s. This feature will not activate unless a SyncIQ bandwidth rule is in effect.",
			},
			"changelist": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, retain previous source snapshot and incremental repstate, both of which are required for changelist creation.",
				MarkdownDescription: "If true, retain previous source snapshot and incremental repstate, both of which are required for changelist creation.",
			},
			"check_integrity": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, the sync target performs cyclic redundancy checks (CRC) on the data as it is received.",
				MarkdownDescription: "If true, the sync target performs cyclic redundancy checks (CRC) on the data as it is received.",
			},
			"cloud_deep_copy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to deny, replicates all CloudPools smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, the job will fail. If set to force, replicates all smartlinks to the target cluster as regular files. If set to allow, SyncIQ will attempt to replicate smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, SyncIQ will replicate the smartlinks as regular files.",
				MarkdownDescription: "If set to deny, replicates all CloudPools smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, the job will fail. If set to force, replicates all smartlinks to the target cluster as regular files. If set to allow, SyncIQ will attempt to replicate smartlinks to the target cluster as smartlinks; if the target cluster does not support the smartlinks, SyncIQ will replicate the smartlinks as regular files.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"deny",
						"allow",
						"force",
					),
				},
			},
			"delete_quotas": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, forcibly remove quotas on the target after they have been removed on the source.",
				MarkdownDescription: "If true, forcibly remove quotas on the target after they have been removed on the source.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User-assigned description of this sync policy.",
				MarkdownDescription: "User-assigned description of this sync policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_file_split": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  If true, the 7.2+ file splitting capability will be disabled.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  If true, the 7.2+ file splitting capability will be disabled.",
			},
			"disable_fofb": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  Enable/disable sync failover/failback.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  Enable/disable sync failover/failback.",
			},
			"disable_quota_tmp_dir": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to true, SyncIQ will not create temporary quota directories to aid in replication to target paths which contain quotas.",
				MarkdownDescription: "If set to true, SyncIQ will not create temporary quota directories to aid in replication to target paths which contain quotas.",
			},
			"disable_stf": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  Enable/disable the 6.5+ STF based data transfer and uses only treewalk.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  Enable/disable the 6.5+ STF based data transfer and uses only treewalk.",
			},
			"enable_hash_tmpdir": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, syncs will use temporary working directory subdirectories to reduce lock contention.",
				MarkdownDescription: "If true, syncs will use temporary working directory subdirectories to reduce lock contention.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, jobs will be automatically run based on this policy, according to its schedule.",
				MarkdownDescription: "If true, jobs will be automatically run based on this policy, according to its schedule.",
			},
			"encryption_cipher_list": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The cipher list (comma separated) being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.",
				MarkdownDescription: "The cipher list (comma separated) being used with encryption. For SyncIQ targets, this list serves as a list of supported ciphers. For SyncIQ sources, the list of ciphers will be attempted to be used in order.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"expected_dataloss": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  Continue sending files even with the corrupted filesystem.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  Continue sending files even with the corrupted filesystem.",
			},
			"file_matching_pattern": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "A file matching pattern, organized as an OR'ed set of AND'ed file criteria, for example ((a AND b) OR (x AND y)) used to define a set of files with specific properties.  Policies of type 'sync' cannot use 'path' or time criteria in their matching patterns, but policies of type 'copy' can use all listed criteria.",
				MarkdownDescription: "A file matching pattern, organized as an OR'ed set of AND'ed file criteria, for example ((a AND b) OR (x AND y)) used to define a set of files with specific properties.  Policies of type 'sync' cannot use 'path' or time criteria in their matching patterns, but policies of type 'copy' can use all listed criteria.",
				Attributes: map[string]schema.Attribute{
					"or_criteria": schema.ListNestedAttribute{
						Required:            true,
						Description:         "An array containing objects with \"and_criteria\" properties, each set of and_criteria will be logically OR'ed together to create the full file matching pattern.",
						MarkdownDescription: "An array containing objects with \"and_criteria\" properties, each set of and_criteria will be logically OR'ed together to create the full file matching pattern.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"and_criteria": schema.ListNestedAttribute{
									Required:            true,
									Description:         "An array containing individual file criterion objects each describing one criterion.  These are logically AND'ed together to form a set of criteria.",
									MarkdownDescription: "An array containing individual file criterion objects each describing one criterion.  These are logically AND'ed together to form a set of criteria.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"attribute_exists": schema.BoolAttribute{
												Optional:            true,
												Description:         "For \"custom_attribute\" type criteria.  The file will match as long as the attribute named by \"field\" exists.  Default is true.",
												MarkdownDescription: "For \"custom_attribute\" type criteria.  The file will match as long as the attribute named by \"field\" exists.  Default is true.",
											},
											"case_sensitive": schema.BoolAttribute{
												Optional:            true,
												Description:         "If true, the value comparison will be case sensitive.  Default is true.",
												MarkdownDescription: "If true, the value comparison will be case sensitive.  Default is true.",
											},
											"field": schema.StringAttribute{
												Optional:            true,
												Description:         "The name of the file attribute to match on (only required if this is a custom_attribute type criterion).  Default is an empty string \"\".",
												MarkdownDescription: "The name of the file attribute to match on (only required if this is a custom_attribute type criterion).  Default is an empty string \"\".",
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 255),
												},
											},
											"operator": schema.StringAttribute{
												Optional: true,
												Description: "How to compare the specified attribute of each file to the specified value." +
													"  Possible values are: `==`, `!=`, `>`, `>=`, `<`, `<=`, `!`.  Default is `==`.",
												MarkdownDescription: "How to compare the specified attribute of each file to the specified value." +
													"  Possible values are: `==`, `!=`, `>`, `>=`, `<`, `<=`, `!`.  Default is `==`.",
												Validators: []validator.String{
													stringvalidator.OneOf(
														"==",
														"!=",
														">",
														">=",
														"<",
														"<=",
														"!",
													),
												},
											},
											"type": schema.StringAttribute{
												Optional: true,
												Description: "The type of this criterion, that is, which file attribute to match on." +
													` Accepted values are "name", "path", "accessed_time", "birth_time", "changed_time", "size", "file_type", "posix_regex_name", "user_name", "user_id", "group_name", "group_id", "no_user", "no_group".`,
												MarkdownDescription: "The type of this criterion, that is, which file attribute to match on." +
													" Accepted values are , `name`, `path`, `accessed_time`, `birth_time`, `changed_time`, `size`, `file_type`, `posix_regex_name`, `user_name`, `user_id`, `group_name`, `group_id`, `no_user`, `no_group`.",
												Validators: []validator.String{
													stringvalidator.OneOf(
														"name",
														"path",
														"accessed_time",
														"accessed_before",
														"accessed_after",
														"birth_time",
														"birth_before",
														"birth_after",
														"changed_time",
														"changed_before",
														"changed_after",
														"size",
														"file_type",
														"posix_regex_name",
														"user_name",
														"user_id",
														"group_name",
														"group_id",
														"no_user",
														"no_group",
													),
												},
											},
											"value": schema.StringAttribute{
												Optional:            true,
												Description:         "The value to compare the specified attribute of each file to.",
												MarkdownDescription: "The value to compare the specified attribute of each file to.",
											},
											"whole_word": schema.BoolAttribute{
												Optional:            true,
												Description:         "If true, the attribute must match the entire word.  Default is true.",
												MarkdownDescription: "If true, the attribute must match the entire word.  Default is true.",
											},
										},
									},
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
				},
			},
			"force_interface": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.  If you enable this option, the net.inet.ip.choose_ifa_by_ipsrc sysctl should be set.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  Determines whether data is sent only through the subnet and pool specified in the \"source_network\" field. This option can be useful if there are multiple interfaces for the given source subnet.  If you enable this option, the net.inet.ip.choose_ifa_by_ipsrc sysctl should be set.",
			},
			"ignore_recursive_quota": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to true, SyncIQ will not check the recursive quota in target paths to aid in replication to target paths which contain no quota but target cluster has lots of quotas.",
				MarkdownDescription: "If set to true, SyncIQ will not check the recursive quota in target paths to aid in replication to target paths which contain no quota but target cluster has lots of quotas.",
			},
			"job_delay": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "If `schedule` is set to `when-source-modified`, the duration to wait after a modification is made before starting a job (default is 0 seconds).",
				MarkdownDescription: "If `schedule` is set to `when-source-modified`, the duration to wait after a modification is made before starting a job (default is 0 seconds).",
			},
			"log_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Severity an event must reach before it is logged." +
					" Accepted values are `fatal`, `error`, `notice`, `info`, `copy`, `debug`, `trace`.",
				MarkdownDescription: "Severity an event must reach before it is logged." +
					" Accepted values are `fatal`, `error`, `notice`, `info`, `copy`, `debug`, `trace`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"fatal",
						"error",
						"notice",
						"info",
						"copy",
						"debug",
						"trace",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_removed_files": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, the system will log any files or directories that are deleted due to a sync.",
				MarkdownDescription: "If true, the system will log any files or directories that are deleted due to a sync.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "User-assigned name of this sync policy.",
				MarkdownDescription: "User-assigned name of this sync policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"ocsp_address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The address of the OCSP responder to which to connect. Set to empty string to disable OCSP.",
				MarkdownDescription: "The address of the OCSP responder to which to connect. Set to empty string to disable OCSP.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"ocsp_issuer_certificate_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the certificate authority that issued the certificate whose revocation status is being checked. Set to empty string to disable certificate verification.",
				MarkdownDescription: "The ID of the certificate authority that issued the certificate whose revocation status is being checked. Set to empty string to disable certificate verification.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Description:         "The password for the target cluster. This field is not readable.",
				MarkdownDescription: "The password for the target cluster. This field is not readable.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Determines the priority level of a policy. Policies with higher priority will have precedence to run over lower priority policies. Valid range is [0, 1]. Default is 0.",
				MarkdownDescription: "Determines the priority level of a policy. Policies with higher priority will have precedence to run over lower priority policies. Valid range is [0, 1]. Default is 0.",
			},
			"report_max_age": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Length of time (in seconds) a policy report will be stored.",
				MarkdownDescription: "Length of time (in seconds) a policy report will be stored.",
			},
			"report_max_count": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Maximum number of policy reports that will be stored on the system.",
				MarkdownDescription: "Maximum number of policy reports that will be stored on the system.",
			},
			"restrict_target_network": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.",
				MarkdownDescription: "If you specify true, and you specify a SmartConnect zone in the \"target_host\" field, replication policies will connect only to nodes in the specified SmartConnect zone.  If you specify false, replication policies are not restricted to specific nodes on the target cluster.",
			},
			"rpo_alert": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "If `schedule` is set to a time/date, an alert is created if the specified RPO for this policy is exceeded. The default value is 0, which will not generate RPO alerts.",
				MarkdownDescription: "If `schedule` is set to a time/date, an alert is created if the specified RPO for this policy is exceeded. The default value is 0, which will not generate RPO alerts.",
			},
			"schedule": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The schedule on which new jobs will be run for this policy.",
				MarkdownDescription: "The schedule on which new jobs will be run for this policy.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_lookup": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Skip DNS lookup of target IPs.",
				MarkdownDescription: "Skip DNS lookup of target IPs.",
			},
			"skip_when_source_unmodified": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true and `schedule` is set to a time/date, the policy will not run if no changes have been made to the contents of the source directory since the last job successfully completed.",
				MarkdownDescription: "If true and `schedule` is set to a time/date, the policy will not run if no changes have been made to the contents of the source directory since the last job successfully completed.",
			},
			"snapshot_sync_existing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, snapshot-triggered syncs will include snapshots taken before policy creation time (requires --schedule when-snapshot-taken).",
				MarkdownDescription: "If true, snapshot-triggered syncs will include snapshots taken before policy creation time (requires --schedule when-snapshot-taken).",
			},
			"snapshot_sync_pattern": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The naming pattern that a snapshot must match to trigger a sync when the schedule is when-snapshot-taken (default is \"*\").",
				MarkdownDescription: "The naming pattern that a snapshot must match to trigger a sync when the schedule is when-snapshot-taken (default is \"*\").",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"source_exclude_directories": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Directories that will be excluded from the sync.  Modifying this field will result in a full synchronization of all data.",
				MarkdownDescription: "Directories that will be excluded from the sync.  Modifying this field will result in a full synchronization of all data.",
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"source_include_directories": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Directories that will be included in the sync.  Modifying this field will result in a full synchronization of all data.",
				MarkdownDescription: "Directories that will be included in the sync.  Modifying this field will result in a full synchronization of all data.",
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"source_network": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Restricts replication policies on the local cluster to running on the specified subnet and pool.",
				MarkdownDescription: "Restricts replication policies on the local cluster to running on the specified subnet and pool.",
				Attributes: map[string]schema.Attribute{
					"pool": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "The pool to restrict replication policies to.",
						MarkdownDescription: "The pool to restrict replication policies to.",
					},
					"subnet": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "The subnet to restrict replication policies to.",
						MarkdownDescription: "The subnet to restrict replication policies to.",
					},
				},
				Validators: []validator.Object{
					objectvalidator.AlsoRequires(path.MatchRelative().AtName("subnet")),
					objectvalidator.AlsoRequires(path.MatchRelative().AtName("pool")),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"source_root_path": schema.StringAttribute{
				Required:            true,
				Description:         "The root directory on the source cluster the files will be synced from.  Modifying this field will result in a full synchronization of all data.",
				MarkdownDescription: "The root directory on the source cluster the files will be synced from.  Modifying this field will result in a full synchronization of all data.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"source_snapshot_archive": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, archival snapshots of the source data will be taken on the source cluster before a sync.",
				MarkdownDescription: "If true, archival snapshots of the source data will be taken on the source cluster before a sync.",
			},
			"source_snapshot_expiration": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The length of time in seconds to keep snapshots on the source cluster.",
				MarkdownDescription: "The length of time in seconds to keep snapshots on the source cluster.",
			},
			"source_snapshot_pattern": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The name pattern for snapshots taken on the source cluster before a sync.",
				MarkdownDescription: "The name pattern for snapshots taken on the source cluster before a sync.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"sync_existing_snapshot_expiration": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If set to true, the expire duration for target archival snapshot is the remaining expire duration of source snapshot, requires --sync-existing-snapshot=true",
				MarkdownDescription: "If set to true, the expire duration for target archival snapshot is the remaining expire duration of source snapshot, requires --sync-existing-snapshot=true",
			},
			"sync_existing_target_snapshot_pattern": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The naming pattern for snapshot on the destination cluster when --sync-existing-snapshot is true",
				MarkdownDescription: "The naming pattern for snapshot on the destination cluster when --sync-existing-snapshot is true",
				Validators: []validator.String{
					// even though openAPI spec allows empty string, actual API does not
					// and its giving a non-user-friendly error
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"target_certificate_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the target cluster certificate being used for encryption. Set to empty string to disable target certificate verification.",
				MarkdownDescription: "The ID of the target cluster certificate being used for encryption. Set to empty string to disable target certificate verification.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"target_compare_initial_sync": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, the target creates diffs against the original sync.",
				MarkdownDescription: "If true, the target creates diffs against the original sync.",
			},
			"target_detect_modifications": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, target cluster will detect if files have been changed on the target by legacy tree walk syncs.",
				MarkdownDescription: "If true, target cluster will detect if files have been changed on the target by legacy tree walk syncs.",
			},
			"target_host": schema.StringAttribute{
				Required:            true,
				Description:         "Hostname or IP address of sync target cluster.  Modifying the target cluster host can result in the policy being unrunnable if the new target does not match the current target association.",
				MarkdownDescription: "Hostname or IP address of sync target cluster.  Modifying the target cluster host can result in the policy being unrunnable if the new target does not match the current target association.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"target_path": schema.StringAttribute{
				Required:            true,
				Description:         "Absolute filesystem path on the target cluster for the sync destination.",
				MarkdownDescription: "Absolute filesystem path on the target cluster for the sync destination.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"target_snapshot_alias": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The alias of the snapshot taken on the target cluster after the sync completes. Do not use the value `DEFAULT`.",
				MarkdownDescription: "The alias of the snapshot taken on the target cluster after the sync completes. Do not use the value `DEFAULT`.",
				Validators: []validator.String{
					// even though openAPI spec allows empty string, actual API does not
					// and its giving a non-user-friendly error
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"target_snapshot_archive": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "If true, archival snapshots of the target data will be taken on the target cluster after successful sync completions.",
				MarkdownDescription: "If true, archival snapshots of the target data will be taken on the target cluster after successful sync completions.",
			},
			"target_snapshot_expiration": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The length of time in seconds to keep snapshots on the target cluster.",
				MarkdownDescription: "The length of time in seconds to keep snapshots on the target cluster.",
			},
			"target_snapshot_pattern": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The name pattern for snapshots taken on the target cluster after the sync completes. Do not use the value `@DEFAULT`.",
				MarkdownDescription: "The name pattern for snapshots taken on the target cluster after the sync completes. Do not use the value `@DEFAULT`.",
				Validators: []validator.String{
					// even though openAPI spec allows empty string, actual API does not
					// and its giving a non-user-friendly error
					stringvalidator.LengthAtLeast(1),
				},
			},
			"workers_per_node": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The number of worker threads on a node performing a sync.",
				MarkdownDescription: "The number of worker threads on a node performing a sync.",
			},
			"conflicted": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "NOTE: This field should not be changed without the help of PowerScale support.  If true, the most recent run of this policy encountered an error and this policy will not start any more scheduled jobs until this field is manually set back to 'false'.",
				MarkdownDescription: "NOTE: This field should not be changed without the help of PowerScale support.  If true, the most recent run of this policy encountered an error and this policy will not start any more scheduled jobs until this field is manually set back to 'false'.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The system ID given to this sync policy.",
				MarkdownDescription: "The system ID given to this sync policy.",
			},
		},
	}
}
