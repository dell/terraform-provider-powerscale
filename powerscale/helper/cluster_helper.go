package helper

import (
	"context"
	powerscale "dell/powerscale-go-client"
	"terraform-provider-powerscale/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetClusterConfigSchema Get cluster config schema
func GetClusterConfigSchema() schema.Attribute {
	config := schema.SingleNestedAttribute{
		MarkdownDescription: "The configuration information of cluster.",
		Description:         "The configuration information of cluster.",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description:         "Customer configurable description.",
				MarkdownDescription: "Customer configurable description.",
				Computed:            true,
			},
			"devices": schema.ListNestedAttribute{
				Description:         "device",
				MarkdownDescription: "device",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"devid": schema.Int64Attribute{
							Description:         "Device ID.",
							MarkdownDescription: "Device ID.",
							Computed:            true,
						},
						"guid": schema.StringAttribute{
							Description:         "Device GUID.",
							MarkdownDescription: "Device GUID.",
							Computed:            true,
						},
						"is_up": schema.BoolAttribute{
							Description:         "If true, this node is online and communicating with the local node and every other node with the is_up property normally",
							MarkdownDescription: "If true, this node is online and communicating with the local node and every other node with the is_up property normally",
							Computed:            true,
						},
						"lnn": schema.Int64Attribute{
							Description:         "Device logical node number.",
							MarkdownDescription: "Device logical node number.",
							Computed:            true,
						},
					},
				},
			},
			"guid": schema.StringAttribute{
				Description:         "Cluster GUID.",
				MarkdownDescription: "Cluster GUID.",
				Computed:            true,
			},

			"join_mode": schema.StringAttribute{
				Description:         "Node join mode: 'manual' or 'secure'.",
				MarkdownDescription: "Node join mode: 'manual' or 'secure'.",
				Computed:            true,
			},
			"local_devid": schema.Int64Attribute{
				Description:         "Device ID of the queried node.",
				MarkdownDescription: "Device ID of the queried node.",
				Computed:            true,
			},
			"local_lnn": schema.Int64Attribute{
				Description:         "Device logical node number of the queried node.",
				MarkdownDescription: "Device logical node number of the queried node.",
				Computed:            true,
			},
			"local_serial": schema.StringAttribute{
				Description:         "Device serial number of the queried node.",
				MarkdownDescription: "Device serial number of the queried node.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Cluster name.",
				MarkdownDescription: "Cluster name.",
				Computed:            true,
			},
			"onefs_version": schema.SingleNestedAttribute{
				Description:         "version",
				MarkdownDescription: "version",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"build": schema.StringAttribute{
						Description:         "OneFS build string.",
						MarkdownDescription: "OneFS build string.",
						Computed:            true,
					},
					"release": schema.StringAttribute{
						Description:         "Kernel release number.",
						MarkdownDescription: "Kernel release number.",
						Computed:            true,
					},
					"revision": schema.StringAttribute{
						Description:         "OneFS build number.",
						MarkdownDescription: "OneFS build number.",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						Description:         "Kernel release type.",
						MarkdownDescription: "Kernel release type.",
						Computed:            true,
					},
					"version": schema.StringAttribute{
						Description:         "Kernel full version information.",
						MarkdownDescription: "Kernel full version information.",
						Computed:            true,
					},
				},
			},
			"timezone": schema.SingleNestedAttribute{
				Description:         "version",
				MarkdownDescription: "version",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"abbreviation": schema.StringAttribute{
						Description:         "Timezone abbreviation.",
						MarkdownDescription: "Timezone abbreviation.",
						Computed:            true,
					},
					"custom": schema.StringAttribute{
						Description:         "Customer timezone information.",
						MarkdownDescription: "Customer timezone information.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description:         "Timezone full name.",
						MarkdownDescription: "Timezone full name.",
						Computed:            true,
					},
					"path": schema.StringAttribute{
						Description:         "Timezone hierarchical name.",
						MarkdownDescription: "Timezone hierarchical name.",
						Computed:            true,
					},
				},
			},
		},
	}
	return config
}

// GetClusterIdentitySchema Get cluster identity schema
func GetClusterIdentitySchema() schema.Attribute {
	identity := schema.SingleNestedAttribute{
		MarkdownDescription: "Unprivileged cluster information for display when logging in.",
		Description:         "Unprivileged cluster information for display when logging in.",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description:         "A description of the cluster.",
				MarkdownDescription: "A description of the cluster.",
				Computed:            true,
			},
			"logon": schema.SingleNestedAttribute{
				Description:         "	//",
				MarkdownDescription: "	//",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"motd": schema.StringAttribute{
						Description:         "The message of the day.",
						MarkdownDescription: "The message of the day.",
						Computed:            true,
					},
					"motd_header": schema.StringAttribute{
						Description:         "The header to the message of the day.",
						MarkdownDescription: "The header to the message of the day.",
						Computed:            true,
					},
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the cluster.",
				MarkdownDescription: "The name of the cluster.",
				Computed:            true,
			},
		},
	}
	return identity
}

// GetClusterNodeSchema get cluster node schema
func GetClusterNodeSchema() schema.Attribute {
	attribute := schema.SingleNestedAttribute{
		MarkdownDescription: "IsiClusterNodes struct for IsiClusterNodes",
		Description:         "IsiClusterNodes struct for IsiClusterNodes",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"errors": schema.ListNestedAttribute{
				Description:         "A list of errors encountered by the individual nodes involved in this request, or an empty list if there were no errors.",
				MarkdownDescription: "A list of errors encountered by the individual nodes involved in this request, or an empty list if there were no errors.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							Description:         "The error code.",
							MarkdownDescription: "The error code.",
							Computed:            true,
						},
						"field": schema.StringAttribute{
							Description:         "The field with the error if applicable.",
							MarkdownDescription: "The field with the error if applicable.",
							Computed:            true,
						},
						"id": schema.Int64Attribute{
							Description:         "Node ID (Device Number) of a node.",
							MarkdownDescription: "Node ID (Device Number) of a node.",
							Computed:            true,
						},
						"lnn": schema.Int64Attribute{
							Description:         "Logical Node Number (LNN) of a node.",
							MarkdownDescription: "Logical Node Number (LNN) of a node.",
							Computed:            true,
						},
						"message": schema.StringAttribute{
							Description:         "The error message.",
							MarkdownDescription: "The error message.",
							Computed:            true,
						},
						"status": schema.Int64Attribute{
							Description:         "HTTP Status code returned by this node.",
							MarkdownDescription: "HTTP Status code returned by this node.",
							Computed:            true,
						},
					},
				},
			},
			"nodes": schema.ListNestedAttribute{
				Description:         "The responses from the individual nodes involved in this request.",
				MarkdownDescription: "The responses from the individual nodes involved in this request.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"drives": schema.ListNestedAttribute{
							Description:         "List of the drives in this node.",
							MarkdownDescription: "List of the drives in this node.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"baynum": schema.Int64Attribute{
										Description:         "Numerical representation of this drive's bay.",
										MarkdownDescription: "Numerical representation of this drive's bay.",
										Computed:            true,
									},
									"blocks": schema.Int64Attribute{
										Description:         "Number of blocks on this drive.",
										MarkdownDescription: "Number of blocks on this drive.",
										Computed:            true,
									},
									"chassis": schema.Int64Attribute{
										Description:         "The chassis number which contains this drive.",
										MarkdownDescription: "The chassis number which contains this drive.",
										Computed:            true,
									},
									"devname": schema.StringAttribute{
										Description:         "This drive's device name.",
										MarkdownDescription: "This drive's device name.",
										Computed:            true,
									},
									"firmware": schema.SingleNestedAttribute{
										Description:         "Drive firmware information",
										MarkdownDescription: "Drive firmware information",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"current_firmware": schema.StringAttribute{
												Description:         "This drive's current firmware revision",
												MarkdownDescription: "This drive's current firmware revision",
												Computed:            true,
											},
											"desired_firmware": schema.StringAttribute{
												Description:         "This drive's desired firmware revision.",
												MarkdownDescription: "This drive's desired firmware revision.",
												Computed:            true,
											},
										},
									},
									"handle": schema.Int64Attribute{
										Description:         "Drive_d's handle representation for this driveIf we fail to retrieve the handle for this drive from drive_d: -1",
										MarkdownDescription: "Drive_d's handle representation for this driveIf we fail to retrieve the handle for this drive from drive_d: -1",
										Computed:            true,
									},
									"interface_type": schema.StringAttribute{
										Description:         "String representation of this drive's interface type.",
										MarkdownDescription: "String representation of this drive's interface type.",
										Computed:            true,
									},
									"lnum": schema.Int64Attribute{
										Description:         "This drive's logical drive number in IFS.",
										MarkdownDescription: "This drive's logical drive number in IFS.",
										Computed:            true,
									},
									"locnstr": schema.StringAttribute{
										Description:         "String representation of this drive's physical location.",
										MarkdownDescription: "String representation of this drive's physical location.",
										Computed:            true,
									},
									"logical_block_length": schema.Int64Attribute{
										Description:         "Size of a logical block on this drive.",
										MarkdownDescription: "Size of a logical block on this drive.",
										Computed:            true,
									},
									"media_type": schema.StringAttribute{
										Description:         "String representation of this drive's media type.",
										MarkdownDescription: "String representation of this drive's media type.",
										Computed:            true,
									},
									"model": schema.StringAttribute{
										Description:         "This drive's manufacturer and model.",
										MarkdownDescription: "This drive's manufacturer and model.",
										Computed:            true,
									},
									"physical_block_length": schema.Int64Attribute{
										Description:         "Size of a physical block on this drive.",
										MarkdownDescription: "Size of a physical block on this drive.",
										Computed:            true,
									},
									"present": schema.BoolAttribute{
										Description:         "Indicates whether this drive is physically present in the node.",
										MarkdownDescription: "Indicates whether this drive is physically present in the node.",
										Computed:            true,
									},
									"purpose": schema.StringAttribute{
										Description:         "This drive's purpose in the DRV state machine.",
										MarkdownDescription: "This drive's purpose in the DRV state machine.",
										Computed:            true,
									},
									"purpose_description": schema.StringAttribute{
										Description:         "Description of this drive's purpose.",
										MarkdownDescription: "Description of this drive's purpose.",
										Computed:            true,
									},
									"serial": schema.StringAttribute{
										Description:         "Serial number for this drive.",
										MarkdownDescription: "Serial number for this drive.",
										Computed:            true,
									},
									"ui_state": schema.StringAttribute{
										Description:         "This drive's state as presented to the UI.",
										MarkdownDescription: "This drive's state as presented to the UI.",
										Computed:            true,
									},
									"wwn": schema.StringAttribute{
										Description:         "The drive's 'worldwide name' from its NAA identifiers.",
										MarkdownDescription: "The drive's 'worldwide name' from its NAA identifiers.",
										Computed:            true,
									},
									"x_loc": schema.Int64Attribute{
										Description:         "This drive's x-axis grid location.",
										MarkdownDescription: "This drive's x-axis grid location.",
										Computed:            true,
									},
									"y_loc": schema.Int64Attribute{
										Description:         "This drive's y-axis grid location.",
										MarkdownDescription: "This drive's y-axis grid location.",
										Computed:            true,
									},
								},
							},
						},
						"error": schema.StringAttribute{
							Description:         "Error message, if the HTTP status returned from this node was not 200.",
							MarkdownDescription: "Error message, if the HTTP status returned from this node was not 200.",
							Computed:            true,
						},
						"hardware": schema.SingleNestedAttribute{
							Description:         "	//",
							MarkdownDescription: "	//",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"chassis": schema.StringAttribute{
									Description:         "Name of this node's chassis.",
									MarkdownDescription: "Name of this node's chassis.",
									Computed:            true,
								},
								"chassis_code": schema.StringAttribute{
									Description:         "Chassis code of this node (1U, 2U, etc.).",
									MarkdownDescription: "Chassis code of this node (1U, 2U, etc.).",
									Computed:            true,
								},
								"chassis_count": schema.StringAttribute{
									Description:         "Number of chassis making up this node.",
									MarkdownDescription: "Number of chassis making up this node.",
									Computed:            true,
								},
								"class": schema.StringAttribute{
									Description:         "Class of this node (storage, accelerator, etc.).",
									MarkdownDescription: "Class of this node (storage, accelerator, etc.).",
									Computed:            true,
								},
								"configuration_id": schema.StringAttribute{
									Description:         "Node configuration ID.",
									MarkdownDescription: "Node configuration ID.",
									Computed:            true,
								},
								"cpu": schema.StringAttribute{
									Description:         "Manufacturer and model of this node's CPU.",
									MarkdownDescription: "Manufacturer and model of this node's CPU.",
									Computed:            true,
								},
								"disk_controller": schema.StringAttribute{
									Description:         "Manufacturer and model of this node's disk controller.",
									MarkdownDescription: "Manufacturer and model of this node's disk controller.",
									Computed:            true,
								},
								"disk_expander": schema.StringAttribute{
									Description:         "Manufacturer and model of this node's disk expander.",
									MarkdownDescription: "Manufacturer and model of this node's disk expander.",
									Computed:            true,
								},
								"family_code": schema.StringAttribute{
									Description:         "Family code of this node (X, S, NL, etc.).",
									MarkdownDescription: "Family code of this node (X, S, NL, etc.).",
									Computed:            true,
								},
								"flash_drive": schema.StringAttribute{
									Description:         "Manufacturer, model, and device id of this node's flash drive.",
									MarkdownDescription: "Manufacturer, model, and device id of this node's flash drive.",
									Computed:            true,
								},
								"generation_code": schema.StringAttribute{
									Description:         "Generation code of this node.",
									MarkdownDescription: "Generation code of this node.",
									Computed:            true,
								},
								"hwgen": schema.StringAttribute{
									Description:         "PowerScale hardware generation name.",
									MarkdownDescription: "PowerScale hardware generation name.",
									Computed:            true,
								},
								"imb_version": schema.StringAttribute{
									Description:         "Version of this node's PowerScale Management Board.",
									MarkdownDescription: "Version of this node's PowerScale Management Board.",
									Computed:            true,
								},
								"infiniband": schema.StringAttribute{
									Description:         "Infiniband card type.",
									MarkdownDescription: "Infiniband card type.",
									Computed:            true,
								},
								"lcd_version": schema.StringAttribute{
									Description:         "Version of the LCD panel.",
									MarkdownDescription: "Version of the LCD panel.",
									Computed:            true,
								},
								"motherboard": schema.StringAttribute{
									Description:         "Manufacturer and model of this node's motherboard.",
									MarkdownDescription: "Manufacturer and model of this node's motherboard.",
									Computed:            true,
								},
								"net_interfaces": schema.StringAttribute{
									Description:         "Description of all this node's network interfaces.",
									MarkdownDescription: "Description of all this node's network interfaces.",
									Computed:            true,
								},
								"nvram": schema.StringAttribute{
									Description:         "Manufacturer and model of this node's NVRAM board.",
									MarkdownDescription: "Manufacturer and model of this node's NVRAM board.",
									Computed:            true,
								},
								"powersupplies": schema.ListAttribute{
									Description:         "Description strings for each power supply on this node.",
									MarkdownDescription: "Description strings for each power supply on this node.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"processor": schema.StringAttribute{
									Description:         "Number of processors and cores on this node.",
									MarkdownDescription: "Number of processors and cores on this node.",
									Computed:            true,
								},
								"product": schema.StringAttribute{
									Description:         "PowerScale product name.",
									MarkdownDescription: "PowerScale product name.",
									Computed:            true,
								},
								"ram": schema.Int64Attribute{
									Description:         "Size of RAM in bytes.",
									MarkdownDescription: "Size of RAM in bytes.",
									Computed:            true,
								},
								"serial_number": schema.StringAttribute{
									Description:         "Serial number of this node.",
									MarkdownDescription: "Serial number of this node.",
									Computed:            true,
								},
								"series": schema.StringAttribute{
									Description:         "Series of this node (X, I, NL, etc.).",
									MarkdownDescription: "Series of this node (X, I, NL, etc.).",
									Computed:            true,
								},
								"storage_class": schema.StringAttribute{
									Description:         "Storage class of this node (storage or diskless).",
									MarkdownDescription: "Storage class of this node (storage or diskless).",
									Computed:            true,
								},
							},
						},
						"id": schema.Int64Attribute{
							Description:         "Node ID (Device Number) of a node.",
							MarkdownDescription: "Node ID (Device Number) of a node.",
							Computed:            true,
						},
						"lnn": schema.Int64Attribute{
							Description:         "Logical Node Number (LNN) of a node.",
							MarkdownDescription: "Logical Node Number (LNN) of a node.",
							Computed:            true,
						},
						"partitions": schema.SingleNestedAttribute{
							Description:         "	//",
							MarkdownDescription: "	//",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"count": schema.Int64Attribute{
									Description:         "Count of how many partitions are included.",
									MarkdownDescription: "Count of how many partitions are included.",
									Computed:            true,
								},
								"partitions": schema.ListNestedAttribute{
									Description:         "Partition information.",
									MarkdownDescription: "Partition information.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"block_size": schema.Int64Attribute{
												Description:         "The block size used for the reported partition information.",
												MarkdownDescription: "The block size used for the reported partition information.",
												Computed:            true,
											},
											"capacity": schema.Int64Attribute{
												Description:         "Total blocks on this file system partition.",
												MarkdownDescription: "Total blocks on this file system partition.",
												Computed:            true,
											},
											"component_devices": schema.StringAttribute{
												Description:         "Comma separated list of devices used for this file system partition.",
												MarkdownDescription: "Comma separated list of devices used for this file system partition.",
												Computed:            true,
											},
											"mount_point": schema.StringAttribute{
												Description:         "Directory on which this partition is mounted.",
												MarkdownDescription: "Directory on which this partition is mounted.",
												Computed:            true,
											},
											"percent_used": schema.StringAttribute{
												Description:         "Used blocks on this file system partition, expressed as a percentage.",
												MarkdownDescription: "Used blocks on this file system partition, expressed as a percentage.",
												Computed:            true,
											},
											"statfs": schema.SingleNestedAttribute{
												Description:         "	//",
												MarkdownDescription: "	//",
												Computed:            true,
												Attributes: map[string]schema.Attribute{
													"f_bavail": schema.Int64Attribute{
														Description:         "Free blocks available to non-superuser on this partition.",
														MarkdownDescription: "Free blocks available to non-superuser on this partition.",
														Computed:            true,
													},
													"f_bfree": schema.Int64Attribute{
														Description:         "Free blocks on this partition.",
														MarkdownDescription: "Free blocks on this partition.",
														Computed:            true,
													},
													"f_blocks": schema.Int64Attribute{
														Description:         "Total data blocks on this partition.",
														MarkdownDescription: "Total data blocks on this partition.",
														Computed:            true,
													},
													"f_bsize": schema.Int64Attribute{
														Description:         "Filesystem fragment size; block size in OneFS.",
														MarkdownDescription: "Filesystem fragment size; block size in OneFS.",
														Computed:            true,
													},
													"f_ffree": schema.Int64Attribute{
														Description:         "Free file nodes avail to non-superuser.",
														MarkdownDescription: "Free file nodes avail to non-superuser.",
														Computed:            true,
													},
													"f_files": schema.Int64Attribute{
														Description:         "Total file nodes in filesystem.",
														MarkdownDescription: "Total file nodes in filesystem.",
														Computed:            true,
													},
													"f_flags": schema.Int64Attribute{
														Description:         "Mount exported flags.",
														MarkdownDescription: "Mount exported flags.",
														Computed:            true,
													},
													"f_fstypename": schema.StringAttribute{
														Description:         "File system type name.",
														MarkdownDescription: "File system type name.",
														Computed:            true,
													},
													"f_iosize": schema.Int64Attribute{
														Description:         "Optimal transfer block size.",
														MarkdownDescription: "Optimal transfer block size.",
														Computed:            true,
													},
													"f_mntfromname": schema.StringAttribute{
														Description:         "Names of devices this partition is mounted from.",
														MarkdownDescription: "Names of devices this partition is mounted from.",
														Computed:            true,
													},
													"f_mntonname": schema.StringAttribute{
														Description:         "Directory this partition is mounted to.",
														MarkdownDescription: "Directory this partition is mounted to.",
														Computed:            true,
													},
													"f_namemax": schema.Int64Attribute{
														Description:         "Maximum filename length.",
														MarkdownDescription: "Maximum filename length.",
														Computed:            true,
													},
													"f_owner": schema.Int64Attribute{
														Description:         "UID of user that mounted the filesystem.",
														MarkdownDescription: "UID of user that mounted the filesystem.",
														Computed:            true,
													},
													"f_type": schema.Int64Attribute{
														Description:         "Type of filesystem.",
														MarkdownDescription: "Type of filesystem.",
														Computed:            true,
													},
													"f_version": schema.Int64Attribute{
														Description:         "statfs() structure version number.",
														MarkdownDescription: "statfs() structure version number.",
														Computed:            true,
													},
												},
											},
											"used": schema.Int64Attribute{
												Description:         "Used blocks on this file system partition.",
												MarkdownDescription: "Used blocks on this file system partition.",
												Computed:            true,
											},
										},
									},
								},
							},
						},
						"sensors": schema.SingleNestedAttribute{
							Description:         "	//",
							MarkdownDescription: "	//",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"sensors": schema.ListNestedAttribute{
									Description:         "This node's sensor information.",
									MarkdownDescription: "This node's sensor information.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"count": schema.Int64Attribute{
												Description:         "The count of values in this sensor group.",
												MarkdownDescription: "The count of values in this sensor group.",
												Computed:            true,
											},
											"name": schema.StringAttribute{
												Description:         "The name of this sensor group.",
												MarkdownDescription: "The name of this sensor group.",
												Computed:            true,
											},
											"values": schema.ListNestedAttribute{
												Description:         "The list of specific sensor value info in this sensor group.",
												MarkdownDescription: "The list of specific sensor value info in this sensor group.",
												Computed:            true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"desc": schema.StringAttribute{
															Description:         "The descriptive name of this sensor.",
															MarkdownDescription: "The descriptive name of this sensor.",
															Computed:            true,
														},
														"name": schema.StringAttribute{
															Description:         "The identifier name of this sensor.",
															MarkdownDescription: "The identifier name of this sensor.",
															Computed:            true,
														},
														"units": schema.StringAttribute{
															Description:         "The units of this sensor.",
															MarkdownDescription: "The units of this sensor.",
															Computed:            true,
														},
														"value": schema.StringAttribute{
															Description:         "The value of this sensor.",
															MarkdownDescription: "The value of this sensor.",
															Computed:            true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"state": schema.SingleNestedAttribute{
							Description:         "	//",
							MarkdownDescription: "	//",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"readonly": schema.SingleNestedAttribute{
									Description:         "    //",
									MarkdownDescription: "    //",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"allowed": schema.BoolAttribute{
											Description:         "The current read-only mode allowed status for the node.",
											MarkdownDescription: "The current read-only mode allowed status for the node.",
											Computed:            true,
										},
										"enabled": schema.BoolAttribute{
											Description:         "The current read-only user mode status for the node. NOTE: If read-only mode is currently disallowed for this node, it will remain read/write until read-only mode is allowed again. This value only sets or clears any user-specified requests for read-only mode. If the node has been placed into read-only mode by the system, it will remain in read-only mode until the system conditions which triggered read-only mode have cleared.",
											MarkdownDescription: "The current read-only user mode status for the node. NOTE: If read-only mode is currently disallowed for this node, it will remain read/write until read-only mode is allowed again. This value only sets or clears any user-specified requests for read-only mode. If the node has been placed into read-only mode by the system, it will remain in read-only mode until the system conditions which triggered read-only mode have cleared.",
											Computed:            true,
										},
										"mode": schema.BoolAttribute{
											Description:         "The current read-only mode status for the node.",
											MarkdownDescription: "The current read-only mode status for the node.",
											Computed:            true,
										},
										"status": schema.StringAttribute{
											Description:         "The current read-only mode status description for the node.",
											MarkdownDescription: "The current read-only mode status description for the node.",
											Computed:            true,
										},
										"valid": schema.BoolAttribute{
											Description:         "The read-only state values are valid (False = Error).",
											MarkdownDescription: "The read-only state values are valid (False = Error).",
											Computed:            true,
										},
										"value": schema.Int64Attribute{
											Description:         "The current read-only value (enumerated bitfield) for the node.",
											MarkdownDescription: "The current read-only value (enumerated bitfield) for the node.",
											Computed:            true,
										},
									},
								},
								"servicelight": schema.SingleNestedAttribute{
									Description:         "	//",
									MarkdownDescription: "	//",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Description:         "The node service light state (True = on).",
											MarkdownDescription: "The node service light state (True = on).",
											Computed:            true,
										},
									},
								},
								"smartfail": schema.SingleNestedAttribute{
									Description:         "	//",
									MarkdownDescription: "	//",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"smartfailed": schema.BoolAttribute{
											Description:         "This node is smartfailed (soft_devs).",
											MarkdownDescription: "This node is smartfailed (soft_devs).",
											Computed:            true,
										},
									},
								},
							},
						},
						"status": schema.SingleNestedAttribute{
							Description:         "	//",
							MarkdownDescription: "	//",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"batterystatus": schema.SingleNestedAttribute{
									Description:         "    //",
									MarkdownDescription: "    //",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"last_test_time1": schema.StringAttribute{
											Description:         "The last battery test time for battery 1.",
											MarkdownDescription: "The last battery test time for battery 1.",
											Computed:            true,
										},
										"last_test_time2": schema.StringAttribute{
											Description:         "The last battery test time for battery 2.",
											MarkdownDescription: "The last battery test time for battery 2.",
											Computed:            true,
										},
										"next_test_time1": schema.StringAttribute{
											Description:         "The next checkup for battery 1.",
											MarkdownDescription: "The next checkup for battery 1.",
											Computed:            true,
										},
										"next_test_time2": schema.StringAttribute{
											Description:         "The next checkup for battery 2.",
											MarkdownDescription: "The next checkup for battery 2.",
											Computed:            true,
										},
										"present": schema.BoolAttribute{
											Description:         "Node has battery status.",
											MarkdownDescription: "Node has battery status.",
											Computed:            true,
										},
										"result1": schema.StringAttribute{
											Description:         "The result of the last battery test for battery 1.",
											MarkdownDescription: "The result of the last battery test for battery 1.",
											Computed:            true,
										},
										"result2": schema.StringAttribute{
											Description:         "The result of the last battery test for battery 2.",
											MarkdownDescription: "The result of the last battery test for battery 2.",
											Computed:            true,
										},
										"status1": schema.StringAttribute{
											Description:         "The status of battery 1.",
											MarkdownDescription: "The status of battery 1.",
											Computed:            true,
										},
										"status2": schema.StringAttribute{
											Description:         "The status of battery 2.",
											MarkdownDescription: "The status of battery 2.",
											Computed:            true,
										},
										"supported": schema.BoolAttribute{
											Description:         "Node supports battery status.",
											MarkdownDescription: "Node supports battery status.",
											Computed:            true,
										},
									},
								},
								"capacity": schema.ListNestedAttribute{
									Description:         "Storage capacity of this node.",
									MarkdownDescription: "Storage capacity of this node.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"bytes": schema.Int64Attribute{
												Description:         "Total device storage bytes.",
												MarkdownDescription: "Total device storage bytes.",
												Computed:            true,
											},
											"count": schema.Int64Attribute{
												Description:         "Total device count.",
												MarkdownDescription: "Total device count.",
												Computed:            true,
											},
											"type": schema.StringAttribute{
												Description:         "Device type.",
												MarkdownDescription: "Device type.",
												Computed:            true,
											},
										},
									},
								},
								"cpu": schema.SingleNestedAttribute{
									Description:         "	//",
									MarkdownDescription: "	//",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"model": schema.StringAttribute{
											Description:         "Manufacturer model description of this CPU.",
											MarkdownDescription: "Manufacturer model description of this CPU.",
											Computed:            true,
										},
										"overtemp": schema.StringAttribute{
											Description:         "CPU overtemp state.",
											MarkdownDescription: "CPU overtemp state.",
											Computed:            true,
										},
										"proc": schema.StringAttribute{
											Description:         "Type of processor and core of this CPU.",
											MarkdownDescription: "Type of processor and core of this CPU.",
											Computed:            true,
										},
										"speed_limit": schema.StringAttribute{
											Description:         "CPU throttling (expressed as a percentage).",
											MarkdownDescription: "CPU throttling (expressed as a percentage).",
											Computed:            true,
										},
									},
								},
								"nvram": schema.SingleNestedAttribute{
									Description:         "	//",
									MarkdownDescription: "	//",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"batteries": schema.ListNestedAttribute{
											Description:         "This node's NVRAM battery status information.",
											MarkdownDescription: "This node's NVRAM battery status information.",
											Computed:            true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"color": schema.StringAttribute{
														Description:         "The current status color of the NVRAM battery.",
														MarkdownDescription: "The current status color of the NVRAM battery.",
														Computed:            true,
													},
													"id": schema.Int64Attribute{
														Description:         "Identifying index for the NVRAM battery.",
														MarkdownDescription: "Identifying index for the NVRAM battery.",
														Computed:            true,
													},
													"status": schema.StringAttribute{
														Description:         "The current status message of the NVRAM battery.",
														MarkdownDescription: "The current status message of the NVRAM battery.",
														Computed:            true,
													},
													"voltage": schema.StringAttribute{
														Description:         "The current voltage of the NVRAM battery.",
														MarkdownDescription: "The current voltage of the NVRAM battery.",
														Computed:            true,
													},
												},
											},
										},
										"battery_count": schema.Int64Attribute{
											Description:         "This node's NVRAM battery count. On failure: -1, otherwise 1 or 2.",
											MarkdownDescription: "This node's NVRAM battery count. On failure: -1, otherwise 1 or 2.",
											Computed:            true,
										},
										"charge_status": schema.StringAttribute{
											Description:         "This node's NVRAM battery charge status, as a color.",
											MarkdownDescription: "This node's NVRAM battery charge status, as a color.",
											Computed:            true,
										},
										"charge_status_number": schema.Int64Attribute{
											Description:         "This node's NVRAM battery charge status, as a number. Error or not supported: -1. BR_BLACK: 0. BR_GREEN: 1. BR_YELLOW: 2. BR_RED: 3.",
											MarkdownDescription: "This node's NVRAM battery charge status, as a number. Error or not supported: -1. BR_BLACK: 0. BR_GREEN: 1. BR_YELLOW: 2. BR_RED: 3.",
											Computed:            true,
										},
										"device": schema.StringAttribute{
											Description:         "This node's NVRAM device name with path.",
											MarkdownDescription: "This node's NVRAM device name with path.",
											Computed:            true,
										},
										"present": schema.BoolAttribute{
											Description:         "This node has NVRAM.",
											MarkdownDescription: "This node has NVRAM.",
											Computed:            true,
										},
										"present_flash": schema.BoolAttribute{
											Description:         "This node has NVRAM with flash storage.",
											MarkdownDescription: "This node has NVRAM with flash storage.",
											Computed:            true,
										},
										"present_size": schema.Int64Attribute{
											Description:         "The size of the NVRAM, in bytes.",
											MarkdownDescription: "The size of the NVRAM, in bytes.",
											Computed:            true,
										},
										"present_type": schema.StringAttribute{
											Description:         "This node's NVRAM type.",
											MarkdownDescription: "This node's NVRAM type.",
											Computed:            true,
										},
										"ship_mode": schema.Int64Attribute{
											Description:         "This node's current ship mode state for NVRAM batteries. If not supported or on failure: -1. Disabled: 0. Enabled: 1.",
											MarkdownDescription: "This node's current ship mode state for NVRAM batteries. If not supported or on failure: -1. Disabled: 0. Enabled: 1.",
											Computed:            true,
										},
										"supported": schema.BoolAttribute{
											Description:         "This node supports NVRAM.",
											MarkdownDescription: "This node supports NVRAM.",
											Computed:            true,
										},
										"supported_flash": schema.BoolAttribute{
											Description:         "This node supports NVRAM with flash storage.",
											MarkdownDescription: "This node supports NVRAM with flash storage.",
											Computed:            true,
										},
										"supported_size": schema.Int64Attribute{
											Description:         "The maximum size of the NVRAM, in bytes.",
											MarkdownDescription: "The maximum size of the NVRAM, in bytes.",
											Computed:            true,
										},
										"supported_type": schema.StringAttribute{
											Description:         "This node's supported NVRAM type.",
											MarkdownDescription: "This node's supported NVRAM type.",
											Computed:            true,
										},
									},
								},
								"powersupplies": schema.SingleNestedAttribute{
									Description:         "	//",
									MarkdownDescription: "	//",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"count": schema.Int64Attribute{
											Description:         "Count of how many power supplies are supported.",
											MarkdownDescription: "Count of how many power supplies are supported.",
											Computed:            true,
										},
										"failures": schema.Int64Attribute{
											Description:         "Count of how many power supplies have failed.",
											MarkdownDescription: "Count of how many power supplies have failed.",
											Computed:            true,
										},
										"has_cff": schema.BoolAttribute{
											Description:         "Does this node have a CFF power supply.",
											MarkdownDescription: "Does this node have a CFF power supply.",
											Computed:            true,
										},
										"status": schema.StringAttribute{
											Description:         "A descriptive status string for this node's power supplies.",
											MarkdownDescription: "A descriptive status string for this node's power supplies.",
											Computed:            true,
										},
										"supplies": schema.ListNestedAttribute{
											Description:         "List of this node's power supplies.",
											MarkdownDescription: "List of this node's power supplies.",
											Computed:            true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"chassis": schema.Int64Attribute{
														Description:         "Which node chassis is this power supply in.",
														MarkdownDescription: "Which node chassis is this power supply in.",
														Computed:            true,
													},
													"firmware": schema.StringAttribute{
														Description:         "The current firmware revision of this power supply.",
														MarkdownDescription: "The current firmware revision of this power supply.",
														Computed:            true,
													},
													"good": schema.StringAttribute{
														Description:         "Is this power supply in a failure state.",
														MarkdownDescription: "Is this power supply in a failure state.",
														Computed:            true,
													},
													"id": schema.Int64Attribute{
														Description:         "Identifying index for this power supply.",
														MarkdownDescription: "Identifying index for this power supply.",
														Computed:            true,
													},
													"name": schema.StringAttribute{
														Description:         "Complete identifying string for this power supply.",
														MarkdownDescription: "Complete identifying string for this power supply.",
														Computed:            true,
													},
													"status": schema.StringAttribute{
														Description:         "A descriptive status string for this power supply.",
														MarkdownDescription: "A descriptive status string for this power supply.",
														Computed:            true,
													},
													"type": schema.StringAttribute{
														Description:         "The type of this power supply.",
														MarkdownDescription: "The type of this power supply.",
														Computed:            true,
													},
												},
											},
										},
										"supports_cff": schema.BoolAttribute{
											Description:         "Does this node support CFF power supplies.",
											MarkdownDescription: "Does this node support CFF power supplies.",
											Computed:            true,
										},
									},
								},
								"release": schema.StringAttribute{
									Description:         "OneFS release.",
									MarkdownDescription: "OneFS release.",
									Computed:            true,
								},
								"uptime": schema.Int64Attribute{
									Description:         "Seconds this node has been online.",
									MarkdownDescription: "Seconds this node has been online.",
									Computed:            true,
								},
								"version": schema.StringAttribute{
									Description:         "OneFS version.",
									MarkdownDescription: "OneFS version.",
									Computed:            true,
								},
							},
						},
					},
				},
			},
			"total": schema.Int64Attribute{
				Description:         "The total number of nodes responding.",
				MarkdownDescription: "The total number of nodes responding.",
				Computed:            true,
			},
		},
	}
	return attribute
}

// GetClusterInternalNetworksSchema get cluster internal network schema
func GetClusterInternalNetworksSchema() schema.Attribute {
	networks := schema.SingleNestedAttribute{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "V7ClusterInternalNetworks Configuration fields for internal networks.",
		Description:         "V7ClusterInternalNetworks Configuration fields for internal networks.",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"failover_ip_addresses": schema.ListNestedAttribute{
				Description:         "Array of IP address ranges to be used to configure the internal failover network of the OneFS cluster.",
				MarkdownDescription: "Array of IP address ranges to be used to configure the internal failover network of the OneFS cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"high": schema.StringAttribute{
							Description:         "IPv4 address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "IPv4 address in the format: xxx.xxx.xxx.xxx",
							Computed:            true,
						},
						"low": schema.StringAttribute{
							Description:         "IPv4 address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "IPv4 address in the format: xxx.xxx.xxx.xxx",
							Computed:            true,
						},
					},
				},
			},
			"failover_status": schema.StringAttribute{
				Description:         "Status of failover network.",
				MarkdownDescription: "Status of failover network.",
				Computed:            true,
			},
			"int_a_fabric": schema.StringAttribute{
				Description:         "Network fabric used for the primary network int-a.",
				MarkdownDescription: "Network fabric used for the primary network int-a.",
				Computed:            true,
			},
			"int_a_ip_addresses": schema.ListNestedAttribute{
				Description:         "Array of IP address ranges to be used to configure the internal int-a network of the OneFS cluster.",
				MarkdownDescription: "Array of IP address ranges to be used to configure the internal int-a network of the OneFS cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"high": schema.StringAttribute{
							Description:         "IPv4 address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "IPv4 address in the format: xxx.xxx.xxx.xxx",
							Computed:            true,
						},
						"low": schema.StringAttribute{
							Description:         "IPv4 address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "IPv4 address in the format: xxx.xxx.xxx.xxx",
							Computed:            true,
						},
					},
				},
			},
			"int_a_mtu": schema.Int64Attribute{
				Description:         "Maximum Transfer Unit (MTU) of the primary network int-a.",
				MarkdownDescription: "Maximum Transfer Unit (MTU) of the primary network int-a.",
				Computed:            true,
			},
			"int_a_prefix_length": schema.Int64Attribute{
				Description:         "Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.",
				MarkdownDescription: "Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.",
				Computed:            true,
			},
			"int_a_status": schema.StringAttribute{
				Description:         "Status of the primary network int-a.",
				MarkdownDescription: "Status of the primary network int-a.",
				Computed:            true,
			},
			"int_b_fabric": schema.StringAttribute{
				Description:         "Network fabric used for the failover network.",
				MarkdownDescription: "Network fabric used for the failover network.",
				Computed:            true,
			},
			"int_b_ip_addresses": schema.ListNestedAttribute{
				Description:         "Array of IP address ranges to be used to configure the internal int-b network of the OneFS cluster.",
				MarkdownDescription: "Array of IP address ranges to be used to configure the internal int-b network of the OneFS cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"high": schema.StringAttribute{
							Description:         "IPv4 address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "IPv4 address in the format: xxx.xxx.xxx.xxx",
							Computed:            true,
						},
						"low": schema.StringAttribute{
							Description:         "IPv4 address in the format: xxx.xxx.xxx.xxx",
							MarkdownDescription: "IPv4 address in the format: xxx.xxx.xxx.xxx",
							Computed:            true,
						},
					},
				},
			},
			"int_b_mtu": schema.Int64Attribute{
				Description:         "Maximum Transfer Unit (MTU) of the failover network int-b.",
				MarkdownDescription: "Maximum Transfer Unit (MTU) of the failover network int-b.",
				Computed:            true,
			},
			"int_b_prefix_length": schema.Int64Attribute{
				Description:         "Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.",
				MarkdownDescription: "Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.",
				Computed:            true,
			},
		},
	}
	return networks
}

// GetClusterConfig retrieve the cluster information
func GetClusterConfig(ctx context.Context, client *client.Client) (*powerscale.V3ClusterConfig, error) {
	config, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterConfig(ctx).Execute()
	return config, err
}

// GetClusterIdentity retrieve the login information
func GetClusterIdentity(ctx context.Context, client *client.Client) (*powerscale.V1ClusterIdentity, error) {
	identity, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterIdentity(ctx).Execute()
	return identity, err
}

// GetClusterNodes list the nodes on this cluster
func GetClusterNodes(ctx context.Context, client *client.Client) (*powerscale.V3ClusterNodes, error) {
	nodes, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv3ClusterNodes(ctx).Execute()
	return nodes, err
}

// GetClusterInternalNetworks list internal networks settings
func GetClusterInternalNetworks(ctx context.Context, client *client.Client) (*powerscale.V7ClusterInternalNetworks, error) {
	networks, _, err := client.PscaleOpenAPIClient.ClusterApi.GetClusterv7ClusterInternalNetworks(ctx).Execute()
	return networks, err
}
