/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	powerscale "dell/powerscale-go-client"
	"fmt"
	"regexp"
	"terraform-provider-powerscale/client"
	"terraform-provider-powerscale/powerscale/helper"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

var filePoolPolicyMocker *Mocker

func TestAccFilePoolPolicyResource(t *testing.T) {
	var policyResourceName = "powerscale_filepool_policy.policy_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + filePoolPolicyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "tfacc_filePoolPolicy"),
					resource.TestCheckResourceAttr(policyResourceName, "description", "tfacc_filePoolPolicy description"),
					resource.TestCheckResourceAttr(policyResourceName, "apply_order", "1"),
					resource.TestCheckResourceAttr(policyResourceName, "actions.#", "7"),
					resource.TestCheckResourceAttr(policyResourceName, "file_matching_pattern.or_criteria.#", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName: policyResourceName,
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "tfacc_filePoolPolicy", states[0].Attributes["name"])
					assert.Equal(t, "1", states[0].Attributes["apply_order"])
					assert.Equal(t, "tfacc_filePoolPolicy description", states[0].Attributes["description"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + filePoolPolicyResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "tfacc_filePoolPolicy_update"),
					resource.TestCheckResourceAttr(policyResourceName, "description", "tfacc_filePoolPolicy description updated"),
					resource.TestCheckResourceAttr(policyResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(policyResourceName, "actions.0.enable_coalescer_action", "false"),
					resource.TestCheckResourceAttr(policyResourceName, "file_matching_pattern.or_criteria.#", "1"),
					resource.TestCheckResourceAttr(policyResourceName, "file_matching_pattern.or_criteria.0.and_criteria.0.value", "1073741824"),
				),
			},
		},
	})
}

func TestAccFilePoolPolicyResourceErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error - invalid param
			{
				Config:      ProviderConfig + filePoolPolicyInvalidResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error creating File Pool Policy*.`),
			},
			// Create Error - required params
			{
				Config:      ProviderConfig + filePoolPolicyInvalidParamConfig,
				ExpectError: regexp.MustCompile(`.*Error creating File Pool Policy*.`),
			},
			{
				Config: ProviderConfig + filePoolPolicyResourceConfig,
			},
			// Update Error - invalid param
			{
				Config:      ProviderConfig + filePoolPolicyInvalidResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error updating the File Pool Policy resource*.`),
			},
			// Update Error - required param
			{
				Config:      ProviderConfig + filePoolPolicyInvalidParamConfig,
				ExpectError: regexp.MustCompile(`.*file_matching_pattern is required*.`),
			},
			// Update Error - cannot change default flag
			{
				Config:      ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*may not change is_default_policy*.`),
			},
			// Update
			{
				Config: ProviderConfig + filePoolPolicyResourceUpdateConfig,
			},
		},
	})
}

func TestAccFilePoolDefaultPolicyResource(t *testing.T) {
	var policyResourceName = "powerscale_filepool_policy.policy_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + filePoolDefaultPolicyResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "Default policy"),
					resource.TestCheckResourceAttr(policyResourceName, "description", "This policy applies to all files not selected by higher-priority policies."),
					resource.TestCheckResourceAttr(policyResourceName, "actions.#", "6"),
					resource.TestCheckNoResourceAttr(policyResourceName, "file_matching_pattern"),
					resource.TestCheckNoResourceAttr(policyResourceName, "apply_order"),
				),
			},
			// ImportState testing
			{
				ResourceName:  policyResourceName,
				ImportState:   true,
				ImportStateId: "is_default_policy=true",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					assert.Equal(t, "Default policy", states[0].Attributes["name"])
					assert.Equal(t, "This policy applies to all files not selected by higher-priority policies.", states[0].Attributes["description"])
					return nil
				},
			},
			// Update
			{
				Config: ProviderConfig + filePoolDefaultPolicyResourceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(policyResourceName, "name", "Default policy"),
					resource.TestCheckResourceAttr(policyResourceName, "description", "This policy applies to all files not selected by higher-priority policies."),
					resource.TestCheckResourceAttr(policyResourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(policyResourceName, "actions.0.enable_packing_action", "false"),
					resource.TestCheckNoResourceAttr(policyResourceName, "file_matching_pattern"),
					resource.TestCheckNoResourceAttr(policyResourceName, "apply_order"),
				),
			},
			// Update and Read testing
			{
				Config: ProviderConfig + filePoolDefaultPolicyResourceConfig,
			},
		},
	})
}

func TestAccFilePoolDefaultPolicyResourceErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error - invalid name
			{
				Config:      ProviderConfig + filePoolDefaultPolicyInvalidNameConfig,
				ExpectError: regexp.MustCompile(`.*the default policy name should be \"Default policy\"*.`),
			},
			// Create Error - invalid params
			{
				Config:      ProviderConfig + filePoolDefaultPolicyInvalidParamConfig,
				ExpectError: regexp.MustCompile(`.*Error creating File Pool Policy*.`),
			},
			// Create Error - required params
			{
				Config:      ProviderConfig + filePoolDefaultPolicyInvalidActionConfig,
				ExpectError: regexp.MustCompile(`.*default policy not support \"set_cloudpool_policy\" action*.`),
			},
			{
				Config: ProviderConfig + filePoolDefaultPolicyResourceConfig,
			},
			// Update Error - invalid param
			{
				Config:      ProviderConfig + filePoolPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*may not change is_default_policy*.`),
			},
			// Update Error - required param
			{
				Config:      ProviderConfig + filePoolDefaultPolicyInvalidParamConfig,
				ExpectError: regexp.MustCompile(`.*Error updating File Pool Policy*.`),
			},
			// Update
			{
				Config: ProviderConfig + filePoolDefaultPolicyResourceUpdateConfig,
			},
		},
	})
}

func TestAccFilePoolPolicyResourceMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).CreateFilepoolv12FilepoolPolicyExecute).Return(nil, nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Create and Read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).GetFilepoolv12FilepoolPolicyExecute).Return(nil, nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolPolicyResourceConfig,
			},
			// Read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).GetFilepoolv12FilepoolPolicyExecute).Return(nil, nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceActionChangeConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Update Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).UpdateFilepoolv12FilepoolPolicyExecute).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceActionChangeConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Update and Read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.GetFilePoolPolicy).When(func(ctx context.Context, client *client.Client, policyName string) bool {
						return filePoolPolicyMocker.Times() > 1
					}).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceActionChangeConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolPolicyResourceActionChangeConfig,
			},
		},
	})
}

func TestAccFilePoolPolicyResourceDeleteMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolPolicyResourceConfig,
			},
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).DeleteFilepoolv12FilepoolPolicyExecute).When(func(r powerscale.ApiDeleteFilepoolv12FilepoolPolicyRequest) bool {
						return filePoolPolicyMocker.Times() == 1
					}).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolPolicyResourceConfig,
				Destroy:     true,
				ExpectError: regexp.MustCompile("filePoolPolicy mock error"),
			},
		},
	})
}

func TestAccFilePoolPolicyResourceImportMockErr(t *testing.T) {
	var policyResourceName = "powerscale_filepool_policy.policy_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolPolicyResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.GetFilePoolPolicy).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:            ProviderConfig + filePoolPolicyResourceConfig,
				ResourceName:      policyResourceName,
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*filePoolPolicy mock error*.`),
				ImportStateVerify: true,
			},
			// Import and parse Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.UpdateFilePoolPolicyImportState).Return(fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:            ProviderConfig + filePoolPolicyResourceConfig,
				ResourceName:      policyResourceName,
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*filePoolPolicy mock error*.`),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFilePoolDefaultPolicyResourceMockErr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).UpdateFilepoolv4FilepoolDefaultPolicyExecute).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Create and Read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).GetFilepoolv4FilepoolDefaultPolicyExecute).Return(nil, nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Create and Parse Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.UpdateFilePoolDefaultPolicyState).Return(fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error reading Default File Pool Policy Resource*.`),
			},
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolDefaultPolicyResourceConfig,
			},
			// Read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).GetFilepoolv4FilepoolDefaultPolicyExecute).Return(nil, nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Read and Parse Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.UpdateFilePoolDefaultPolicyState).Return(fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ExpectError: regexp.MustCompile(`.*Error reading Default File Pool Policy Resource*.`),
			},
			// Update Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock((*powerscale.FilepoolApiService).UpdateFilepoolv4FilepoolDefaultPolicyExecute).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			// Update and Read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.GetFilePoolDefaultPolicy).When(func(ctx context.Context, client *client.Client) bool {
						return filePoolPolicyMocker.Times() > 1
					}).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:      ProviderConfig + filePoolDefaultPolicyResourceUpdateConfig,
				ExpectError: regexp.MustCompile(`.*filePoolPolicy mock error*.`),
			},
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolDefaultPolicyResourceUpdateConfig,
			},
		},
	})
}

func TestAccFilePoolDefaultPolicyResourceImportMockErr(t *testing.T) {
	var policyResourceName = "powerscale_filepool_policy.policy_test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
				},
				Config: ProviderConfig + filePoolDefaultPolicyResourceConfig,
			},
			// Import and read Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.GetFilePoolDefaultPolicy).Return(nil, fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:            ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ResourceName:      policyResourceName,
				ImportStateId:     "is_default_policy=true",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*filePoolPolicy mock error*.`),
				ImportStateVerify: true,
			},
			// Import and parse Error testing
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.UnPatch()
					}
					filePoolPolicyMocker = Mock(helper.UpdateFilePoolDefaultPolicyState).Return(fmt.Errorf("filePoolPolicy mock error")).Build()
				},
				Config:            ProviderConfig + filePoolDefaultPolicyResourceConfig,
				ResourceName:      policyResourceName,
				ImportStateId:     "is_default_policy=true",
				ImportState:       true,
				ExpectError:       regexp.MustCompile(`.*Error reading Default File Pool Policy Resource*.`),
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccFilePoolPolicyReleaseMockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if filePoolPolicyMocker != nil {
						filePoolPolicyMocker.Release()
					}
				},
				Config: ProviderConfig + filePoolPolicyResourceConfig,
			},
		},
	})
}

var filePoolPolicyResourceConfig = `
resource "powerscale_filepool_policy" "policy_test" {
  name = "tfacc_filePoolPolicy"
  file_matching_pattern = {
    or_criteria = [
      {
        and_criteria = [
          {
            operator = ">"
            type = "size"
            units = "B"
            value = "1073741824"
          },
          {
            operator = ">"
            type = "birth_time"
            use_relative_time = true
            value = "20"
          },
          {
            operator = ">"
            type = "metadata_changed_time"
            use_relative_time = false
            value = "1704742200"
          },
          {
            operator = "<"
            type = "accessed_time"
            use_relative_time = true
            value = "20"
          }
        ]
      },
      {
        and_criteria = [
          {
            operator = "<"
            type = "changed_time"
            use_relative_time = false
            value = "1704820500"
          },
          {
            attribute_exists = false
            field = "test"
            type = "custom_attribute"
            value = ""
          },
          {
            operator = "!="
            type = "file_type"
            value = "directory"
          },
          {
            begins_with = false
            case_sensitive = true
            operator = "!="
            type = "path"
            value = "test"
          },
          {
            case_sensitive = true
            operator = "!="
            type = "name"
            value = "test"
          }
        ]
      }
    ]
  }
  actions = [
    {
      data_access_pattern_action = "concurrency"
      action_type = "set_data_access_pattern"
    },
    {
      data_storage_policy_action = {
        ssd_strategy ="metadata"
        storagepool = "anywhere"
      }
      action_type = "apply_data_storage_policy"
    },
    {
      snapshot_storage_policy_action = {
        ssd_strategy = "metadata"
        storagepool = "anywhere"
      }
      action_type = "apply_snapshot_storage_policy"
    },
    {
      requested_protection_action = "default"
      action_type = "set_requested_protection"
    },
    {
      enable_coalescer_action = true
      action_type = "enable_coalescer"
    },
    {
      enable_packing_action = true,
      action_type = "enable_packing"
    },
    {
      action_type = "set_cloudpool_policy"
      cloudpool_policy_action = {
        archive_snapshot_files = true
        cache = {
          expiration = 86400
          read_ahead = "partial"
          type = "cached"
        }
        compression = true
        data_retention = 604800
        encryption = true
        full_backup_retention = 145152000
        incremental_backup_retention = 145152000
        pool = "cloudPool_policy"
        writeback_frequency = 32400
      }
    }
  ]
  description = "tfacc_filePoolPolicy description"
  apply_order = 1
}
`

var filePoolPolicyResourceUpdateConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "tfacc_filePoolPolicy_update"
	file_matching_pattern = {
	  or_criteria = [
		{
		  and_criteria = [
			{
			  operator = ">"
			  type = "size"
			  units = "B"
			  value = "1073741824"
			}
		  ]
		}
	  ]
	}
	actions = [
	  {
		enable_coalescer_action = false
		action_type = "enable_coalescer"
	  }
	]
	description = "tfacc_filePoolPolicy description updated"
	apply_order = 1
  }
`

var filePoolPolicyResourceActionChangeConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "tfacc_filePoolPolicy"
	file_matching_pattern = {
	  or_criteria = [
		{
		  and_criteria = [
			{
			  operator = ">"
			  type = "size"
			  units = "B"
			  value = "1073741824"
			}
		  ]
		}
	  ]
	}
	actions = [
	  {
		enable_coalescer_action = false
		action_type = "enable_coalescer"
	  }
	]
	description = "tfacc_filePoolPolicy action changed"
	apply_order = 1
  }
`

var filePoolPolicyInvalidResourceConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "tfacc_filePoolPolicy"
	file_matching_pattern = {
	  or_criteria = [
		{
		  and_criteria = [
			{
			  operator = ">"
			  type = "size"
			  units = "B"
			  value = "test"
			}
		  ]
		}
	  ]
	}
	actions = [
	  {
		enable_coalescer_action = false
		action_type = "enable_coalescer"
	  }
	]
	apply_order = 10000
  }
`

var filePoolPolicyInvalidParamConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "tfacc_filePoolPolicy"
	actions = [
	  {
		enable_coalescer_action = false
		action_type = "enable_coalescer"
	  }
	]
	apply_order = 10000
  }
`

var filePoolDefaultPolicyResourceConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "Default policy"
	is_default_policy = true
	actions = [
	  {
		data_access_pattern_action = "concurrency"
		action_type                = "set_data_access_pattern"
	  },
	  {
		data_storage_policy_action = {
		  ssd_strategy = "metadata"
		  storagepool  = "anywhere"
		}
		action_type = "apply_data_storage_policy"
	  },
	  {
		snapshot_storage_policy_action = {
		  ssd_strategy = "metadata"
		  storagepool  = "anywhere"
		}
		action_type = "apply_snapshot_storage_policy"
	  },
	  {
		requested_protection_action = "default"
		action_type                 = "set_requested_protection"
	  },
	  {
		enable_coalescer_action = true
		action_type             = "enable_coalescer"
	  },
	  {
		enable_packing_action = true,
		action_type           = "enable_packing"
	  }
	]
  }
`

var filePoolDefaultPolicyResourceUpdateConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "Default policy"
	is_default_policy = true
	actions = [
	  {
		enable_packing_action = false,
		action_type           = "enable_packing"
	  }
	]
  }
`

var filePoolDefaultPolicyInvalidNameConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "tfacc_filePoolPolicy"
	is_default_policy = true
  }
`

var filePoolDefaultPolicyInvalidParamConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "Default policy"
	is_default_policy = true
	actions = [
	  {
		enable_packing_action = false,
		action_type           = "enable_packing"
	  }
	]
	apply_order = 1
  }
`

var filePoolDefaultPolicyInvalidActionConfig = `
resource "powerscale_filepool_policy" "policy_test" {
	name = "Default policy"
	is_default_policy = true
	actions = [
	  {
		enable_packing_action = false,
		action_type           = "enable_packing"
	  },
	  {
		action_type = "set_cloudpool_policy"
		cloudpool_policy_action = {
		  archive_snapshot_files = true
		  cache = {
			expiration = 86400
			read_ahead = "partial"
			type = "cached"
		  }
		  compression = true
		  data_retention = 604800
		  encryption = true
		  full_backup_retention = 145152000
		  incremental_backup_retention = 145152000
		  pool = "cloudPool_policy"
		  writeback_frequency = 32400
		}
	  }
	]
  }
`
