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

// NfsGlobalSettingsModel Specifies the global NFS configuration settings.
type NfsGlobalSettingsModel struct {
	ID types.String `tfsdk:"id"`
	// True if NFSv3 is enabled.
	Nfsv3Enabled types.Bool `tfsdk:"nfsv3_enabled"`
	// True if the RDMA is enabled for NFSv3.
	Nfsv3RdmaEnabled types.Bool `tfsdk:"nfsv3_rdma_enabled"`
	// True if NFSv4 is enabled.
	Nfsv4Enabled types.Bool `tfsdk:"nfsv4_enabled"`
	// Specifies the maximum number of threads in the nfsd thread pool.
	RPCMaxthreads types.Int64 `tfsdk:"rpc_maxthreads"`
	// Specifies the minimum number of threads in the nfsd thread pool.
	RPCMinthreads types.Int64 `tfsdk:"rpc_minthreads"`
	// True if the rquota protocol is enabled.
	RquotaEnabled types.Bool `tfsdk:"rquota_enabled"`
	// True if the NFS service is enabled. When set to false, the NFS service is disabled.
	Service types.Bool `tfsdk:"service"`
}
