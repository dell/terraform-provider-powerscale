package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ClusterDataSource returns the overall cluster details
type ClusterDataSource struct {
	ID               types.String             `tfsdk:"id"`
	Config           *ClusterConfig           `tfsdk:"config"`
	Identity         *ClusterIdentity         `tfsdk:"identity"`
	Nodes            *ClusterNodes            `tfsdk:"nodes"`
	InternalNetworks *ClusterInternalNetworks `tfsdk:"internal_networks"`
}

// ClusterConfig returns the configuration information of cluster.
type ClusterConfig struct {
	Description  types.String               `tfsdk:"description"`
	Devices      []*ClusterConfigDevice     `tfsdk:"devices"`
	GUID         types.String               `tfsdk:"guid"`
	JoinMode     types.String               `tfsdk:"join_mode"`
	LocalDevID   types.Int64                `tfsdk:"local_devid"`
	LocalLnn     types.Int64                `tfsdk:"local_lnn"`
	LocalSerial  types.String               `tfsdk:"local_serial"`
	Name         types.String               `tfsdk:"name"`
	OnefsVersion *ClusterConfigOnefsVersion `tfsdk:"onefs_version"`
	Timezone     *ClusterConfigTimezone     `tfsdk:"timezone"`
}

// ClusterConfigOnefsVersion struct for ClusterConfigOnefsVersion
type ClusterConfigOnefsVersion struct {
	// OneFS build string.
	Build types.String `tfsdk:"build"`
	// Kernel release number.
	Release types.String `tfsdk:"release"`
	// OneFS build number.
	Revision types.String `tfsdk:"revision"`
	// Kernel release type.
	Type types.String `tfsdk:"type"`
	// Kernel full version information.
	Version types.String `tfsdk:"version"`
}

// ClusterConfigTimezone The cluster timezone settings.
type ClusterConfigTimezone struct {
	// Timezone abbreviation.
	Abbreviation types.String `tfsdk:"abbreviation"`
	// Customer timezone information.
	Custom types.String `tfsdk:"custom"`
	// Timezone full name.
	Name types.String `tfsdk:"name"`
	// Timezone hierarchical name.
	Path types.String `tfsdk:"path"`
}

// ClusterConfigDevice refers to device information of a cluster
type ClusterConfigDevice struct {
	// Device ID.
	DevID types.Int64 `tfsdk:"devid"`
	// Device GUID.
	GUID types.String `tfsdk:"guid"`
	// If true, this node is online and communicating with the local node and every other node with the is_up property normally.
	// If false, this node is not currently communicating with the local node or other nodes with the is_up property.
	// It may be shut down, rebooting, disconnected from the backend network, or connected only to other nodes.
	IsUp types.Bool `tfsdk:"is_up"`
	// Device logical node number.
	Lnn types.Int64 `tfsdk:"lnn"`
}

// ClusterIdentity Unprivileged cluster information for display when logging in.
type ClusterIdentity struct {
	// A description of the cluster.
	Description types.String `tfsdk:"description"`
	// Logon information
	Logon ClusterIdentityLogon `tfsdk:"logon"`
	// The name of the cluster.
	Name types.String `tfsdk:"name"`
}

// ClusterIdentityLogon The information displayed when a user logs in to the cluster.
type ClusterIdentityLogon struct {
	// The message of the day.
	Motd types.String `tfsdk:"motd"`
	// The header to the message of the day.
	MotdHeader types.String `tfsdk:"motd_header"`
}

// ClusterNodes struct for ClusterNodes
type ClusterNodes struct {
	// A list of errors encountered by the individual nodes involved in this request, or an empty list if there were no errors.
	Errors []NodeStatusError `tfsdk:"errors"`
	// The responses from the individual nodes involved in this request.
	Nodes []ClusterNode `tfsdk:"nodes"`
	// The total number of nodes responding.
	Total types.Int64 `tfsdk:"total"`
}

// NodeStatusError An object describing a single error.
type NodeStatusError struct {
	// The error code.
	Code types.String `tfsdk:"code"`
	// The field with the error if applicable.
	Field types.String `tfsdk:"field"`
	// Node ID (Device Number) of a node.
	ID types.Int64 `tfsdk:"id"`
	// Logical Node Number (LNN) of a node.
	Lnn types.Int64 `tfsdk:"lnn"`
	// The error message.
	Message types.String `tfsdk:"message"`
	// HTTP Status code returned by this node.
	Status types.Int64 `tfsdk:"status"`
}

// ClusterNode Node information.
type ClusterNode struct {
	// List of the drives in this node.
	Drives []ClusterNodeDrive `tfsdk:"drives"`
	// Error message, if the HTTP status returned from this node was not 200.
	Error types.String `tfsdk:"error"`
	// Cluster node hardware
	Hardware *ClusterNodeHardware `tfsdk:"hardware"`
	// Node ID (Device Number) of a node.
	ID types.Int64 `tfsdk:"id"`
	// Logical Node Number (LNN) of a node.
	Lnn types.Int64 `tfsdk:"lnn"`
	//
	Partitions *ClusterNodePartitions `tfsdk:"partitions"`
	//
	Sensors *ClusterNodeSensors `tfsdk:"sensors"`
	//
	State *ClusterNodeState `tfsdk:"state"`
	//
	Status *ClusterNodeStatus `tfsdk:"status"`
}

// ClusterNodeDrive Drive information.
type ClusterNodeDrive struct {
	// Numerical representation of this drive's bay.
	Baynum types.Int64 `tfsdk:"baynum"`
	// Number of blocks on this drive.
	Blocks types.Int64 `tfsdk:"blocks"`
	// The chassis number which contains this drive.
	Chassis types.Int64 `tfsdk:"chassis"`
	// This drive's device name.
	Devname types.String `tfsdk:"devname"`
	// Drive firmware information
	Firmware *ClusterNodeDriveFirmware `tfsdk:"firmware"`
	// Drive_d's handle representation for this driveIf we fail to retrieve the handle for this drive from drive_d: -1
	Handle types.Int64 `tfsdk:"handle"`
	// String representation of this drive's interface type.
	InterfaceType types.String `tfsdk:"interface_type"`
	// This drive's logical drive number in IFS.
	Lnum types.Int64 `tfsdk:"lnum"`
	// String representation of this drive's physical location.
	Locnstr types.String `tfsdk:"locnstr"`
	// Size of a logical block on this drive.
	LogicalBlockLength types.Int64 `tfsdk:"logical_block_length"`
	// String representation of this drive's media type.
	MediaType types.String `tfsdk:"media_type"`
	// This drive's manufacturer and model.
	Model types.String `tfsdk:"model"`
	// Size of a physical block on this drive.
	PhysicalBlockLength types.Int64 `tfsdk:"physical_block_length"`
	// Indicates whether this drive is physically present in the node.
	Present types.Bool `tfsdk:"present"`
	// This drive's purpose in the DRV state machine.
	Purpose types.String `tfsdk:"purpose"`
	// Description of this drive's purpose.
	PurposeDescription types.String `tfsdk:"purpose_description"`
	// Serial number for this drive.
	Serial types.String `tfsdk:"serial"`
	// This drive's state as presented to the UI.
	UIState types.String `tfsdk:"ui_state"`
	// The drive's 'worldwide name' from its NAA identifiers.
	Wwn types.String `tfsdk:"wwn"`
	// This drive's x-axis grid location.
	XLoc types.Int64 `tfsdk:"x_loc"`
	// This drive's y-axis grid location.
	YLoc types.Int64 `tfsdk:"y_loc"`
}

// ClusterNodeDriveFirmware Drive firmware information.
type ClusterNodeDriveFirmware struct {
	// This drive's current firmware revision
	CurrentFirmware types.String `tfsdk:"current_firmware"`
	// This drive's desired firmware revision.
	DesiredFirmware types.String `tfsdk:"desired_firmware"`
}

// ClusterNodeHardware Node hardware identifying information (static).
type ClusterNodeHardware struct {
	// Name of this node's chassis.
	Chassis types.String `tfsdk:"chassis"`
	// Chassis code of this node (1U, 2U, etc.).
	ChassisCode types.String `tfsdk:"chassis_code"`
	// Number of chassis making up this node.
	ChassisCount types.String `tfsdk:"chassis_count"`
	// Class of this node (storage, accelerator, etc.).
	Class types.String `tfsdk:"class"`
	// Node configuration ID.
	ConfigurationID types.String `tfsdk:"configuration_id"`
	// Manufacturer and model of this node's CPU.
	CPU types.String `tfsdk:"cpu"`
	// Manufacturer and model of this node's disk controller.
	DiskController types.String `tfsdk:"disk_controller"`
	// Manufacturer and model of this node's disk expander.
	DiskExpander types.String `tfsdk:"disk_expander"`
	// Family code of this node (X, S, NL, etc.).
	FamilyCode types.String `tfsdk:"family_code"`
	// Manufacturer, model, and device id of this node's flash drive.
	FlashDrive types.String `tfsdk:"flash_drive"`
	// Generation code of this node.
	GenerationCode types.String `tfsdk:"generation_code"`
	// PowerScale hardware generation name.
	Hwgen types.String `tfsdk:"hwgen"`
	// Version of this node's PowerScale Management Board.
	ImbVersion types.String `tfsdk:"imb_version"`
	// Infiniband card type.
	Infiniband types.String `tfsdk:"infiniband"`
	// Version of the LCD panel.
	LcdVersion types.String `tfsdk:"lcd_version"`
	// Manufacturer and model of this node's motherboard.
	Motherboard types.String `tfsdk:"motherboard"`
	// Description of all this node's network interfaces.
	NetInterfaces types.String `tfsdk:"net_interfaces"`
	// Manufacturer and model of this node's NVRAM board.
	Nvram types.String `tfsdk:"nvram"`
	// Description strings for each power supply on this node.
	Powersupplies types.List `tfsdk:"powersupplies"`
	// Number of processors and cores on this node.
	Processor types.String `tfsdk:"processor"`
	// PowerScale product name.
	Product types.String `tfsdk:"product"`
	// Size of RAM in bytes.
	RAM types.Int64 `tfsdk:"ram"`
	// Serial number of this node.
	SerialNumber types.String `tfsdk:"serial_number"`
	// Series of this node (X, I, NL, etc.).
	Series types.String `tfsdk:"series"`
	// Storage class of this node (storage or diskless).
	StorageClass types.String `tfsdk:"storage_class"`
}

// ClusterNodePartitions Node partition information.
type ClusterNodePartitions struct {
	// Count of how many partitions are included.
	Count types.Int64 `tfsdk:"count"`
	// Partition information.
	Partitions []ClusterNodePartition `tfsdk:"partitions"`
}

// ClusterNodePartition Node partition information.
type ClusterNodePartition struct {
	// The block size used for the reported partition information.
	BlockSize types.Int64 `tfsdk:"block_size"`
	// Total blocks on this file system partition.
	Capacity types.Int64 `tfsdk:"capacity"`
	// Comma separated list of devices used for this file system partition.
	ComponentDevices types.String `tfsdk:"component_devices"`
	// Directory on which this partition is mounted.
	MountPoint types.String `tfsdk:"mount_point"`
	// Used blocks on this file system partition, expressed as a percentage.
	PercentUsed types.String `tfsdk:"percent_used"`
	//
	Statfs *ClusterNodePartitionStatfs `tfsdk:"statfs"`
	// Used blocks on this file system partition.
	Used types.Int64 `tfsdk:"used"`
}

// ClusterNodePartitionStatfs System partition details as provided by statfs(2).
type ClusterNodePartitionStatfs struct {
	// Free blocks available to non-superuser on this partition.
	FBavail types.Int64 `tfsdk:"f_bavail"`
	// Free blocks on this partition.
	FBfree types.Int64 `tfsdk:"f_bfree"`
	// Total data blocks on this partition.
	FBlocks types.Int64 `tfsdk:"f_blocks"`
	// Filesystem fragment size; block size in OneFS.
	FBsize types.Int64 `tfsdk:"f_bsize"`
	// Free file nodes avail to non-superuser.
	FFfree types.Int64 `tfsdk:"f_ffree"`
	// Total file nodes in filesystem.
	FFiles types.Int64 `tfsdk:"f_files"`
	// Mount exported flags.
	FFlags types.Int64 `tfsdk:"f_flags"`
	// File system type name.
	FFstypename types.String `tfsdk:"f_fstypename"`
	// Optimal transfer block size.
	FIosize types.Int64 `tfsdk:"f_iosize"`
	// Names of devices this partition is mounted from.
	FMntfromname types.String `tfsdk:"f_mntfromname"`
	// Directory this partition is mounted to.
	FMntonname types.String `tfsdk:"f_mntonname"`
	// Maximum filename length.
	FNamemax types.Int64 `tfsdk:"f_namemax"`
	// UID of user that mounted the filesystem.
	FOwner types.Int64 `tfsdk:"f_owner"`
	// Type of filesystem.
	FType types.Int64 `tfsdk:"f_type"`
	// statfs() structure version number.
	FVersion types.Int64 `tfsdk:"f_version"`
}

// ClusterNodeSensors Node sensor information (hardware reported).
type ClusterNodeSensors struct {
	// This node's sensor information.
	Sensors []ClusterNodeSensor `tfsdk:"sensors"`
}

// ClusterNodeSensor Node sensor information.
type ClusterNodeSensor struct {
	// The count of values in this sensor group.
	Count types.Int64 `tfsdk:"count"`
	// The name of this sensor group.
	Name types.String `tfsdk:"name"`
	// The list of specific sensor value info in this sensor group.
	Values []ClusterNodeSensorValue `tfsdk:"values"`
}

// ClusterNodeSensorValue Specific sensor value info.
type ClusterNodeSensorValue struct {
	// The descriptive name of this sensor.
	Desc types.String `tfsdk:"desc"`
	// The identifier name of this sensor.
	Name types.String `tfsdk:"name"`
	// The units of this sensor.
	Units types.String `tfsdk:"units"`
	// The value of this sensor.
	Value types.String `tfsdk:"value"`
}

// ClusterNodeState Node state information (reported and modifiable).
type ClusterNodeState struct {
	//
	Readonly ClusterNodeStateReadonly `tfsdk:"readonly"`
	//
	Servicelight ClusterNodeStateServicelight `tfsdk:"servicelight"`
	//
	Smartfail ClusterNodeStateSmartfail `tfsdk:"smartfail"`
}

// ClusterNodeStateReadonly Node readonly state.
type ClusterNodeStateReadonly struct {
	// The current read-only mode allowed status for the node.
	Allowed types.Bool `tfsdk:"allowed"`
	// The current read-only user mode status for the node. NOTE: If read-only mode is currently disallowed for this node, it will remain read/write until read-only mode is allowed again. This value only sets or clears any user-specified requests for read-only mode. If the node has been placed into read-only mode by the system, it will remain in read-only mode until the system conditions which triggered read-only mode have cleared.
	Enabled types.Bool `tfsdk:"enabled"`
	// The current read-only mode status for the node.
	Mode types.Bool `tfsdk:"mode"`
	// The current read-only mode status description for the node.
	Status types.String `tfsdk:"status"`
	// The read-only state values are valid (False = Error).
	Valid types.Bool `tfsdk:"valid"`
	// The current read-only value (enumerated bitfield) for the node.
	Value types.Int64 `tfsdk:"value"`
}

// ClusterNodeStateServicelight Node service light state.
type ClusterNodeStateServicelight struct {
	// The node service light state (True = on).
	Enabled types.Bool `tfsdk:"enabled"`
}

// ClusterNodeStatus Node status information (hardware reported).
type ClusterNodeStatus struct {
	//
	Batterystatus *ClusterNodeStatusBatterystatus `tfsdk:"batterystatus"`
	// Storage capacity of this node.
	Capacity []IsiClusterNodeStatusCapacityItem `tfsdk:"capacity"`
	//
	CPU *ClusterNodeStatusCPU `tfsdk:"cpu"`
	//
	Nvram *ClusterNodeStatusNvram `tfsdk:"nvram"`
	//
	Powersupplies *ClusterNodeStatusPowersupplies `tfsdk:"powersupplies"`
	// OneFS release.
	Release types.String `tfsdk:"release"`
	// Seconds this node has been online.
	Uptime types.Int64 `tfsdk:"uptime"`
	// OneFS version.
	Version types.String `tfsdk:"version"`
}

// ClusterNodeStateSmartfail Node smartfail state.
type ClusterNodeStateSmartfail struct {
	// This node is smartfailed (soft_devs).
	Smartfailed types.Bool `tfsdk:"smartfailed"`
}

// ClusterNodeStatusBatterystatus Battery status information.
type ClusterNodeStatusBatterystatus struct {
	// The last battery test time for battery 1.
	LastTestTime1 types.String `tfsdk:"last_test_time1"`
	// The last battery test time for battery 2.
	LastTestTime2 types.String `tfsdk:"last_test_time2"`
	// The next checkup for battery 1.
	NextTestTime1 types.String `tfsdk:"next_test_time1"`
	// The next checkup for battery 2.
	NextTestTime2 types.String `tfsdk:"next_test_time2"`
	// Node has battery status.
	Present types.Bool `tfsdk:"present"`
	// The result of the last battery test for battery 1.
	Result1 types.String `tfsdk:"result1"`
	// The result of the last battery test for battery 2.
	Result2 types.String `tfsdk:"result2"`
	// The status of battery 1.
	Status1 types.String `tfsdk:"status1"`
	// The status of battery 2.
	Status2 types.String `tfsdk:"status2"`
	// Node supports battery status.
	Supported types.Bool `tfsdk:"supported"`
}

// IsiClusterNodeStatusCapacityItem Node capacity information.
type IsiClusterNodeStatusCapacityItem struct {
	// Total device storage bytes.
	Bytes types.Int64 `tfsdk:"bytes"`
	// Total device count.
	Count types.Int64 `tfsdk:"count"`
	// Device type.
	Type types.String `tfsdk:"type"`
}

// ClusterNodeStatusCPU CPU status information for this node.
type ClusterNodeStatusCPU struct {
	// Manufacturer model description of this CPU.
	Model types.String `tfsdk:"model"`
	// CPU overtemp state.
	Overtemp types.String `tfsdk:"overtemp"`
	// Type of processor and core of this CPU.
	Proc types.String `tfsdk:"proc"`
	// CPU throttling (expressed as a percentage).
	SpeedLimit types.String `tfsdk:"speed_limit"`
}

// ClusterNodeStatusNvram Node NVRAM information.
type ClusterNodeStatusNvram struct {
	// This node's NVRAM battery status information.
	Batteries []IsiClusterNodeStatusNvramBattery `tfsdk:"batteries"`
	// This node's NVRAM battery count. On failure: -1, otherwise 1 or 2.
	BatteryCount types.Int64 `tfsdk:"battery_count"`
	// This node's NVRAM battery charge status, as a color.
	ChargeStatus types.String `tfsdk:"charge_status"`
	// This node's NVRAM battery charge status, as a number. Error or not supported: -1. BR_BLACK: 0. BR_GREEN: 1. BR_YELLOW: 2. BR_RED: 3.
	ChargeStatusNumber types.Int64 `tfsdk:"charge_status_number"`
	// This node's NVRAM device name with path.
	Device types.String `tfsdk:"device"`
	// This node has NVRAM.
	Present types.Bool `tfsdk:"present"`
	// This node has NVRAM with flash storage.
	PresentFlash types.Bool `tfsdk:"present_flash"`
	// The size of the NVRAM, in bytes.
	PresentSize types.Int64 `tfsdk:"present_size"`
	// This node's NVRAM type.
	PresentType types.String `tfsdk:"present_type"`
	// This node's current ship mode state for NVRAM batteries. If not supported or on failure: -1. Disabled: 0. Enabled: 1.
	ShipMode types.Int64 `tfsdk:"ship_mode"`
	// This node supports NVRAM.
	Supported types.Bool `tfsdk:"supported"`
	// This node supports NVRAM with flash storage.
	SupportedFlash types.Bool `tfsdk:"supported_flash"`
	// The maximum size of the NVRAM, in bytes.
	SupportedSize types.Int64 `tfsdk:"supported_size"`
	// This node's supported NVRAM type.
	SupportedType types.String `tfsdk:"supported_type"`
}

// ClusterNodeStatusPowersupplies Information about this node's power supplies.
type ClusterNodeStatusPowersupplies struct {
	// Count of how many power supplies are supported.
	Count types.Int64 `tfsdk:"count"`
	// Count of how many power supplies have failed.
	Failures types.Int64 `tfsdk:"failures"`
	// Does this node have a CFF power supply.
	HasCff types.Bool `tfsdk:"has_cff"`
	// A descriptive status string for this node's power supplies.
	Status types.String `tfsdk:"status"`
	// List of this node's power supplies.
	Supplies []IsiClusterNodeStatusPowersuppliesSupply `tfsdk:"supplies"`
	// Does this node support CFF power supplies.
	SupportsCff types.Bool `tfsdk:"supports_cff"`
}

// IsiClusterNodeStatusNvramBattery NVRAM battery status information.
type IsiClusterNodeStatusNvramBattery struct {
	// The current status color of the NVRAM battery.
	Color types.String `tfsdk:"color"`
	// Identifying index for the NVRAM battery.
	ID types.Int64 `tfsdk:"id"`
	// The current status message of the NVRAM battery.
	Status types.String `tfsdk:"status"`
	// The current voltage of the NVRAM battery.
	Voltage types.String `tfsdk:"voltage"`
}

// IsiClusterNodeStatusPowersuppliesSupply Power supply information.
type IsiClusterNodeStatusPowersuppliesSupply struct {
	// Which node chassis is this power supply in.
	Chassis types.Int64 `tfsdk:"chassis"`
	// The current firmware revision of this power supply.
	Firmware types.String `tfsdk:"firmware"`
	// Is this power supply in a failure state.
	Good types.String `tfsdk:"good"`
	// Identifying index for this power supply.
	ID types.Int64 `tfsdk:"id"`
	// Complete identifying string for this power supply.
	Name types.String `tfsdk:"name"`
	// A descriptive status string for this power supply.
	Status types.String `tfsdk:"status"`
	// The type of this power supply.
	Type types.String `tfsdk:"type"`
}

// ClusterInternalNetworks Configuration fields for internal networks.
type ClusterInternalNetworks struct {
	// Array of IP address ranges to be used to configure the internal failover network of the OneFS cluster.
	FailoverIPAddresses []ClusterInternalNetworksFailoverIPAddresse `tfsdk:"failover_ip_addresses"`
	// Status of failover network.
	FailoverStatus types.String `tfsdk:"failover_status"`
	// Network fabric used for the primary network int-a.
	IntAFabric types.String `tfsdk:"int_a_fabric"`
	// Array of IP address ranges to be used to configure the internal int-a network of the OneFS cluster.
	IntAIpAddresses []ClusterInternalNetworksFailoverIPAddresse `tfsdk:"int_a_ip_addresses"`
	// Maximum Transfer Unit (MTU) of the primary network int-a.
	IntAMtu types.Int64 `tfsdk:"int_a_mtu"`
	// Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.
	IntAPrefixLength types.Int64 `tfsdk:"int_a_prefix_length"`
	// Status of the primary network int-a.
	IntAStatus types.String `tfsdk:"int_a_status"`
	// Network fabric used for the failover network.
	IntBFabric types.String `tfsdk:"int_b_fabric"`
	// Array of IP address ranges to be used to configure the internal int-b network of the OneFS cluster.
	IntBIpAddresses []ClusterInternalNetworksFailoverIPAddresse `tfsdk:"int_b_ip_addresses"`
	// Maximum Transfer Unit (MTU) of the failover network int-b.
	IntBMtu types.Int64 `tfsdk:"int_b_mtu"`
	// Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.
	IntBPrefixLength types.Int64 `tfsdk:"int_b_prefix_length"`
}

// ClusterInternalNetworksFailoverIPAddresse Specifies range of IP addresses where 'low' is starting address and 'high' is the end address.' Both 'low' and 'high' addresses are inclusive to the range.
type ClusterInternalNetworksFailoverIPAddresse struct {
	// IPv4 address in the format: xxx.xxx.xxx.xxx
	High types.String `tfsdk:"high"`
	// IPv4 address in the format: xxx.xxx.xxx.xxx
	Low types.String `tfsdk:"low"`
}
