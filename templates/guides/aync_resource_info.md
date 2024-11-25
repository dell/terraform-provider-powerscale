---
# Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://mozilla.org/MPL/2.0/


# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
page_title: "Adding Asynchronous Operations"
title: "Adding Asynchronous Operations"
linkTitle: "Addding Asynchronous Operations"
---

# Terraform Internal Behaviour
Terraform's concurrent execution capabilities allow multiple operations to run in parallel, each thread responsible for managing the lifecycle of a resource. However, even with this parallel execution, each individual thread must still execute its operations in a synchronous manner. By default, Terraform's create, update, and delete operations wait for a resource to reach its expected lifecycle state before proceeding with the next step. 

For instance, when a resource is in the creation phase, Terraform will not advance to the subsequent step until the resource has successfully reached the completion phase. In scenarios where platform operations are performed asynchronously, which can result in extended creation times, Terraform will wait for the resource to reach the completion phase before proceeding, thereby introducing a delay in the overall workflow.

# Consideration for Asynchronous Operations In Terraform
When utilizing Terraform, it is essential to consider the implications of asynchronous operations of the platform. To ensure seamless automation, this guide provides recommendations for resources that support asynchronous operations, enabling users to effectively integrate these features into their infrastructure provisioning workflows.

1. The state file doesn't have the full information of the resource because the resource is not yet created.
2. The resource's stage must be refreshed to get its full information, including whether the resource failed to be created.
3. Terraform does not check the stage of a resource again after initiating its creation or deletion. Failed operations are only shown in subsequent refreshes of the resource.
4. The asynchronous resource must have no dependencies, because the resource is not fully created before another operation begins.

# List of Asynchronous Resources
1. SyncIQ replication Job





